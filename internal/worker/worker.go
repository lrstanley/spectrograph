// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

// Package worker implements the Discord bot, which monitors voice channels and
// manages (clones, reorders, removes) matching channels.
package worker

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/session"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/state/store/defaultstore"
	"github.com/diamondburned/arikawa/v3/utils/handler"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/lrstanley/spectrograph/internal/database/ent"
	"github.com/lrstanley/spectrograph/internal/database/ent/guild"
	"github.com/lrstanley/spectrograph/internal/database/ent/guildadminconfig"
	"github.com/lrstanley/spectrograph/internal/database/ent/guildconfig"
	"github.com/lrstanley/spectrograph/internal/database/ent/predicate"
	"github.com/lrstanley/spectrograph/internal/database/ent/privacy"
	"github.com/lrstanley/spectrograph/internal/guildlogger"
	"github.com/lrstanley/spectrograph/internal/models"
)

// intents are the gateway intents used by the worker. Neither is privileged:
//   - Guilds: guild create/update/delete, roles (for permission checks),
//     channels, and the bot's own member updates.
//   - GuildVoiceStates: voice states in guild create events, and voice state
//     updates (to know which channels are occupied).
//
// Interactions (slash commands, buttons, and modals) are always sent, without
// an intent.
const intents = gateway.IntentGuilds | gateway.IntentGuildVoiceStates

const (
	// maxConcurrentPasses bounds how many guilds are processed concurrently
	// (e.g. after connecting, when every guild is synced to the database).
	maxConcurrentPasses = 10

	initialBackoff = time.Second
	maxBackoff     = 2 * time.Minute
)

// Worker is the Discord bot. See [New].
type Worker struct {
	logger *slog.Logger
	es     *guildlogger.EventStream
	db     *ent.Client
	config models.ConfigWorker

	// client is the REST client shared by all shards, so they share rate
	// limits.
	client *api.Client
	shards []*shard

	passSem chan struct{}

	// minRetry and maxRetry bound the delay before retrying failed passes.
	minRetry, maxRetry time.Duration
	// cooldown is how long changes are paused after Discord rejects them.
	cooldown time.Duration

	// interactions handles slash commands, buttons, and modals.
	interactions *interactionHandler

	// commandsGuildID is the guild to register commands in (instead of
	// globally), if valid.
	commandsGuildID discord.GuildID
	commandsOnce    sync.Once

	mu     sync.Mutex
	runCtx context.Context //nolint:containedctx // set by Run, used to start guild actors from gateway handlers.
	actors map[discord.GuildID]*guildActor
	wg     sync.WaitGroup
}

// shard is a single gateway connection (and its state cache).
type shard struct {
	id    int
	state *state.State
}

// Option configures a [Worker].
type Option func(*Worker)

// WithDashboardURL sets the base URL of the dashboard, which is linked from
// the replies to the bot's slash commands.
func WithDashboardURL(u *url.URL) Option {
	return func(w *Worker) { w.interactions.dashboardURL = u }
}

// New creates a new worker. All database operations performed by the worker
// bypass privacy rules. ctx should be the long-lived application context, as it
// is also used for the guild event stream.
//
// New registers hooks on db, which wake the worker when guild configs are
// updated (e.g. from the dashboard), so changes apply right away.
func New(ctx context.Context, logger *slog.Logger, db *ent.Client, config models.ConfigWorker, opts ...Option) (*Worker, error) {
	token := strings.TrimSpace(strings.TrimPrefix(config.BotToken, "Bot "))
	if token == "" {
		return nil, errors.New("worker bot token is required when the worker is enabled")
	}

	if config.NumShards < 1 {
		config.NumShards = 1
	}

	if len(config.ShardIDs) < 1 {
		config.ShardIDs = []uint{0}
	}

	for _, id := range config.ShardIDs {
		if id >= config.NumShards {
			return nil, fmt.Errorf("shard id %d is out of range (num shards: %d)", id, config.NumShards)
		}
	}

	var commandsGuildID discord.GuildID
	if config.CommandsGuildID != "" {
		sf, err := discord.ParseSnowflake(config.CommandsGuildID)
		if err != nil || !sf.IsValid() {
			return nil, fmt.Errorf("invalid commands guild id %q", config.CommandsGuildID)
		}
		commandsGuildID = discord.GuildID(sf)
	}

	ctx = privacy.DecisionContext(ctx, privacy.Allow)
	logger = logger.With("src", "worker", "shard_ids", config.ShardIDs)
	setLogger(logger.With("src", "arikawa"))

	w := &Worker{
		logger:          logger,
		es:              guildlogger.NewEventStream(ctx, logger, db, config.EventFlushInterval),
		db:              db,
		config:          config,
		client:          api.NewClient("Bot " + token),
		passSem:         make(chan struct{}, maxConcurrentPasses),
		minRetry:        minRetryDelay,
		maxRetry:        maxRetryDelay,
		cooldown:        forbiddenCooldown,
		actors:          make(map[discord.GuildID]*guildActor),
		commandsGuildID: commandsGuildID,
	}
	w.interactions = &interactionHandler{
		logger:     logger.With("component", "interactions"),
		db:         db,
		es:         w.es,
		deferAfter: deferAfter,
	}
	for _, opt := range opts {
		opt(w)
	}
	w.watchConfigs()

	// The identify rate limiters are shared by all shards (as they are
	// per-bot), and are updated from the session start limits when connecting.
	id := gateway.DefaultIdentifier("Bot " + token)
	id.Intents = option.NewUint(uint(intents))
	id.Presence = &gateway.UpdatePresenceCommand{
		Status:     discord.OnlineStatus,
		Activities: []discord.Activity{{Name: "voice chan curator", Type: discord.GameActivity}},
	}

	for _, shardID := range config.ShardIDs {
		sid := id
		sid.SetShard(int(shardID), int(config.NumShards)) //nolint:gosec // validated above, shard counts are small.

		st := state.NewFromSession(session.NewCustom(sid, w.client, handler.New()), defaultstore.New())
		st.StateLog = func(err error) {
			logger.DebugContext(context.Background(), "state cache error", "shard_id", shardID, "error", err)
		}

		s := &shard{id: int(shardID), state: st} //nolint:gosec // validated above.
		w.registerHandlers(s)
		w.shards = append(w.shards, s)
	}

	logger.DebugContext(ctx, "registered worker", "num_shards", config.NumShards)
	return w, nil
}

// Run connects to the Discord gateway, and blocks until ctx is cancelled (in
// which case the gateway connections are closed gracefully, and nil is
// returned), or a shard fails with an unrecoverable error (e.g. an invalid
// token, or disallowed intents). Pending guild events are flushed to the
// database before it returns. Run must only be called once.
func (w *Worker) Run(ctx context.Context) error {
	ctx = privacy.DecisionContext(ctx, privacy.Allow)

	// Flush the guild events logged while shutting down, once the actors have
	// exited (and before the caller closes the database).
	defer w.es.Close()

	if err := w.checkGateway(ctx); err != nil {
		if ctx.Err() != nil {
			return nil
		}
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	w.mu.Lock()
	w.runCtx = ctx
	w.mu.Unlock()

	var swg sync.WaitGroup
	errs := make([]error, len(w.shards))
	for i, s := range w.shards {
		swg.Go(func() {
			if err := w.runShard(ctx, s); err != nil {
				errs[i] = err
				cancel()
			}
		})
	}
	swg.Wait()

	// Gateways are closed, so no new actors can be started.
	cancel()
	w.wg.Wait()

	return errors.Join(errs...)
}

// checkGateway validates the bot token, and logs the session start limits.
// Session start (identify) limits are enforced by arikawa, using the limits
// returned by Discord when each shard connects. Transient errors are retried.
func (w *Worker) checkGateway(ctx context.Context) error {
	backoff := initialBackoff
	for {
		bot, err := w.client.WithContext(ctx).BotURL()
		if err == nil {
			attrs := []any{"recommended_shards", bot.Shards, "num_shards", w.config.NumShards}
			if bot.StartLimit != nil {
				attrs = append(attrs,
					"identify_total", bot.StartLimit.Total,
					"identify_remaining", bot.StartLimit.Remaining,
					"identify_reset_after", bot.StartLimit.ResetAfter.Duration(),
					"identify_max_concurrency", bot.StartLimit.MaxConcurrency,
				)
				if bot.StartLimit.Remaining < len(w.shards) {
					w.logger.WarnContext(ctx, "session start limit nearly exhausted, connecting will wait for it to reset", attrs...)
				}
			}
			w.logger.InfoContext(ctx, "connecting to discord gateway", attrs...)
			return nil
		}

		if isUnauthorized(err) {
			return fmt.Errorf("discord rejected the bot token: %w", err)
		}

		w.logger.WarnContext(ctx, "unable to query discord gateway, retrying", "error", err, "backoff", backoff)
		if !sleep(ctx, backoff) {
			return ctx.Err()
		}
		backoff = min(backoff*2, maxBackoff)
	}
}

// runShard keeps the shard connected until ctx is cancelled, or an
// unrecoverable error occurs. Reconnects (and resumes) after transient
// websocket errors are handled by arikawa, this only handles the initial
// connection failing, or the gateway giving up.
func (w *Worker) runShard(ctx context.Context, s *shard) error {
	logger := w.logger.With("shard_id", s.id)
	backoff := initialBackoff

	for {
		err := s.state.Open(ctx)
		if err == nil {
			logger.InfoContext(ctx, "connected to discord gateway")
			backoff = initialBackoff
			err = s.state.Wait(ctx)
		}

		if ctx.Err() != nil {
			logger.InfoContext(ctx, "disconnected from discord gateway")
			return nil //nolint:nilerr // shutting down, errors from closing are expected.
		}

		if s.state.GatewayOpts().ErrorIsFatalClose(err) || isUnauthorized(err) {
			return fmt.Errorf("shard %d: unrecoverable gateway error: %w", s.id, err)
		}

		logger.WarnContext(ctx, "gateway connection failed, reconnecting", "error", err, "backoff", backoff)

		// Reset the session, otherwise the next Open returns the last error.
		_ = s.state.Close()

		if !sleep(ctx, backoff) {
			return nil //nolint:nilerr // shutting down.
		}
		backoff = min(backoff*2, maxBackoff)
	}
}

// sleep sleeps for d, returning false if ctx is cancelled first.
func sleep(ctx context.Context, d time.Duration) bool {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-t.C:
		return true
	}
}

// actor returns the actor for the guild, or nil if the guild isn't tracked.
func (w *Worker) actor(id discord.GuildID) *guildActor {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.actors[id]
}

// startActor returns the running actor for the guild, starting a new one if
// needed (e.g. the guild became available, or the bot joined it). Returns nil
// if the worker isn't running.
func (w *Worker) startActor(s *shard, id discord.GuildID) *guildActor {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.runCtx == nil || w.runCtx.Err() != nil {
		return nil
	}

	prev := w.actors[id]
	if prev != nil && !prev.isStopped() {
		return prev
	}

	return w.launchActor(s, id, prev)
}

// guildLeft stops the guild's actor, marking the guild as removed in the
// database. If the guild has no running actor (e.g. it was unavailable when the
// bot was removed), one is started just to do so.
func (w *Worker) guildLeft(s *shard, id discord.GuildID) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.runCtx == nil || w.runCtx.Err() != nil {
		return
	}
	w.leaveLocked(s, id)
}

// leaveLocked is [Worker.guildLeft], with w.mu held and the worker running.
func (w *Worker) leaveLocked(s *shard, id discord.GuildID) {
	a := w.actors[id]
	if a == nil || a.isStopped() {
		a = w.launchActor(s, id, a)
	}
	a.stop(true)
}

// shardReady handles a new gateway session for the shard (i.e. not a resumed
// one). The ready event lists every guild of the shard, and no guild delete
// event is sent for guilds the bot was removed from while it was disconnected,
// so guilds which aren't listed are marked as removed: guilds with an actor
// (when reconnecting) right away, and guilds still marked as joined in the
// database (e.g. after a restart) in the background.
func (w *Worker) shardReady(s *shard, ev *gateway.ReadyEvent) {
	handled := make(map[discord.GuildID]struct{}, len(ev.Guilds))
	for i := range ev.Guilds {
		handled[ev.Guilds[i].ID] = struct{}{}
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	if w.runCtx == nil || w.runCtx.Err() != nil {
		return
	}

	for id, a := range w.actors {
		if _, ok := handled[id]; ok || a.shard != s {
			continue
		}
		handled[id] = struct{}{}
		w.leaveLocked(s, id)
	}

	ctx := w.runCtx
	w.wg.Go(func() { w.pruneRemovedGuilds(ctx, s, handled) })
}

// pruneRemovedGuilds marks guilds of the shard which are joined according to
// the database, but which weren't in the shard's ready event (or otherwise
// handled), as removed.
func (w *Worker) pruneRemovedGuilds(ctx context.Context, s *shard, handled map[discord.GuildID]struct{}) {
	ids, err := dbJoinedGuilds(ctx, w.db)
	if err != nil {
		if ctx.Err() == nil {
			w.logger.ErrorContext(ctx, "unable to check for removed guilds", "shard_id", s.id, "error", err)
		}
		return
	}

	for _, id := range ids {
		if _, ok := handled[id]; ok || guildShard(id, w.config.NumShards) != s.id {
			continue
		}

		w.mu.Lock()
		// Guilds with an actor were joined since the ready event.
		if ctx.Err() == nil && w.actors[id] == nil {
			w.leaveLocked(s, id)
		}
		w.mu.Unlock()
	}
}

// guildShard returns the shard ID which receives events for the guild. See:
// https://discord.com/developers/docs/events/gateway#sharding
func guildShard(id discord.GuildID, numShards uint) int {
	return int((uint64(id) >> 22) % uint64(max(numShards, 1))) //nolint:gosec // shard counts are small.
}

// launchActor starts a new actor for the guild, replacing prev (if any), which
// the new actor waits on before starting. w.mu must be held.
func (w *Worker) launchActor(s *shard, id discord.GuildID, prev *guildActor) *guildActor {
	a := newGuildActor(w, s, id, w.client, prev)
	w.actors[id] = a

	ctx := w.runCtx
	w.wg.Go(func() {
		a.run(ctx)

		w.mu.Lock()
		if w.actors[id] == a {
			delete(w.actors, id)
		}
		w.mu.Unlock()
	})
	return a
}

// pokeGuild wakes the guild's actor (if any), to reconcile the guild.
func (w *Worker) pokeGuild(id discord.GuildID) {
	if a := w.actor(id); a != nil {
		a.poke()
	}
}

// watchConfigs registers hooks which wake the guild's actor when its config or
// admin config is updated, so changes (from the dashboard, or the bot's slash
// commands) apply right away, rather than on the guild's next gateway event.
// The hooks only see changes made in this process (through db): when the HTTP
// server runs separately, changes apply on the guild's next event.
func (w *Worker) watchConfigs() {
	if w.db == nil {
		return
	}
	w.db.GuildConfig.Use(ent.ObserveGuildConfigID(w.configChanged(func(id int) predicate.Guild {
		return guild.HasGuildConfigWith(guildconfig.ID(id))
	})))
	w.db.GuildAdminConfig.Use(ent.ObserveGuildAdminConfigID(w.configChanged(func(id int) predicate.Guild {
		return guild.HasGuildAdminConfigWith(guildadminconfig.ID(id))
	})))
}

// configChanged returns an observer which wakes the actor of the guild matching
// owner(id), for updated configs. Errors are only logged, as they would fail
// the update otherwise.
func (w *Worker) configChanged(owner func(id int) predicate.Guild) ent.ObserverHandler {
	return func(ctx context.Context, ev ent.ObserverEvent) error {
		id, ok := ev.ID.(int)
		if ev.Op != ent.ObserverOpUpdate || !ok {
			return nil
		}

		ctx, cancel := context.WithTimeout(privacy.DecisionContext(context.WithoutCancel(ctx), privacy.Allow), 5*time.Second)
		defer cancel()

		gid, err := w.db.Guild.Query().Where(owner(id)).Select(guild.FieldGuildID).String(ctx)
		if err != nil {
			w.logger.WarnContext(ctx, "unable to look up guild of updated config", "type", ev.Type, "id", id, "error", err)
			return nil
		}
		if sf, perr := discord.ParseSnowflake(gid); perr == nil {
			w.pokeGuild(discord.GuildID(sf))
		}
		return nil
	}
}

// shardReadyCommands registers the slash commands in the background, the first
// time shard 0 is ready (not on every reconnect), if enabled.
func (w *Worker) shardReadyCommands(s *shard, appID discord.AppID) {
	if s.id != 0 || !w.config.RegisterCommands || !appID.IsValid() {
		return
	}

	w.commandsOnce.Do(func() {
		w.mu.Lock()
		defer w.mu.Unlock()
		if w.runCtx == nil || w.runCtx.Err() != nil {
			return
		}
		ctx := w.runCtx
		w.wg.Go(func() { w.registerCommands(ctx, appID, w.commandsGuildID) })
	})
}

// interactionCreated handles an interaction in the background, as responding
// needs DB queries and API requests.
func (w *Worker) interactionCreated(s *shard, ev *discord.InteractionEvent) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.runCtx == nil || w.runCtx.Err() != nil {
		return
	}

	ctx := w.runCtx
	w.wg.Go(func() {
		w.interactions.serve(ctx, apiResponder{client: w.client}, s.state.Cabinet, ev)
	})
}
