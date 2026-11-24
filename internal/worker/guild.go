// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package worker

import (
	"context"
	"errors"
	"log/slog"
	"slices"
	"sync"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state/store"
)

const (
	// debounceQuiet is how long a guild must go without events before it is
	// processed, so bursts of events (e.g. many users joining at once, or our
	// own changes) result in a single pass.
	debounceQuiet = 500 * time.Millisecond

	// debounceMax caps how long processing can be delayed by a continuous
	// stream of events.
	debounceMax = 3 * time.Second

	// passTimeout bounds a single pass (DB queries and Discord API calls,
	// including waiting on rate limits).
	passTimeout = 2 * time.Minute

	// minRetryDelay and maxRetryDelay bound how long to wait before retrying a
	// pass which failed with a transient error (e.g. the database being
	// unavailable), as there may be no further events for the guild.
	minRetryDelay = 10 * time.Second
	maxRetryDelay = 5 * time.Minute

	// forbiddenCooldown is how long changes aren't attempted after Discord
	// rejected one for missing permissions (which the bot didn't predict), so a
	// busy guild doesn't make the bot exceed Discord's invalid request limit.
	forbiddenCooldown = time.Minute
)

// passResult is the outcome of a pass. Results are ordered by precedence, so
// the result of several steps is their max.
type passResult int

const (
	passDone     passResult = iota // Wait for the next event.
	passDeferred                   // Changes are paused, retry when the cooldown ends.
	passRetry                      // Retry after a delay (or the next event).
	passExit                       // Stop the actor.
)

// debounce waits until no wakeups are received for quiet, or until max has
// elapsed. Returns false if ctx is cancelled.
func debounce(ctx context.Context, wake <-chan struct{}, quiet, maxDelay time.Duration) bool {
	quietT := time.NewTimer(quiet)
	defer quietT.Stop()
	maxT := time.NewTimer(maxDelay)
	defer maxT.Stop()

	for {
		select {
		case <-ctx.Done():
			return false
		case <-wake:
			quietT.Reset(quiet)
		case <-quietT.C:
			return true
		case <-maxT.C:
			return true
		}
	}
}

// guildActor serializes all processing for a single guild: database syncs,
// and reconciling the guild's voice channels. Gateway handlers only record
// what changed and wake the actor, so they never block the gateway event loop,
// and each pass works from the latest cached state (no events are dropped).
type guildActor struct {
	w     *Worker
	shard *shard
	id    discord.GuildID
	cache *store.Cabinet
	api   *api.Client

	// prev is the previous actor for the same guild (if any), which must exit
	// before this one starts, so DB writes stay ordered (e.g. leave then rejoin).
	prev *guildActor

	wake chan struct{}
	done chan struct{}

	mu      sync.Mutex
	created *gateway.GuildCreateEvent // pending full DB sync.
	updated bool                      // pending partial DB sync.
	left    bool                      // bot was removed from the guild.
	stopped bool                      // guild left or became unavailable.

	// Owned by the run goroutine.
	perms      discord.Permissions
	problems   map[string]problem
	retryDelay time.Duration
	cooldown   time.Time // no changes are attempted until then.
}

func newGuildActor(w *Worker, s *shard, id discord.GuildID, client *api.Client, prev *guildActor) *guildActor {
	return &guildActor{
		w:     w,
		shard: s,
		id:    id,
		cache: s.state.Cabinet,
		api:   client,
		prev:  prev,
		wake:  make(chan struct{}, 1),
		done:  make(chan struct{}),
	}
}

// poke wakes the actor, without blocking.
func (a *guildActor) poke() {
	select {
	case a.wake <- struct{}{}:
	default:
	}
}

func (a *guildActor) guildCreated(ev *gateway.GuildCreateEvent) {
	a.mu.Lock()
	a.created = ev
	a.mu.Unlock()
	a.poke()
}

func (a *guildActor) guildUpdated() {
	a.mu.Lock()
	a.updated = true
	a.mu.Unlock()
	a.poke()
}

// stop stops the actor. If left is true, the bot was removed from the guild,
// and the guild is marked as such in the database before the actor exits.
func (a *guildActor) stop(left bool) {
	a.mu.Lock()
	a.stopped = true
	a.left = a.left || left
	a.mu.Unlock()
	a.poke()
}

func (a *guildActor) isStopped() bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.stopped
}

func (a *guildActor) run(ctx context.Context) {
	defer close(a.done)

	if a.prev != nil {
		select {
		case <-a.prev.done:
			a.prev = nil
		case <-ctx.Done():
			return
		}
	}

	var retry <-chan time.Time
	for {
		select {
		case <-ctx.Done():
			return
		case <-a.wake:
		case <-retry:
		}
		retry = nil

		if !debounce(ctx, a.wake, debounceQuiet, debounceMax) {
			return
		}

		select {
		case a.w.passSem <- struct{}{}:
		case <-ctx.Done():
			return
		}
		res := a.pass(ctx)
		<-a.w.passSem

		switch res {
		case passExit:
			return
		case passRetry:
			a.retryDelay = min(max(a.retryDelay*2, a.w.minRetry), a.w.maxRetry)
			retry = time.After(a.retryDelay)
		case passDeferred:
			// There may be no further events. Fires immediately if the cooldown
			// ended since the pass.
			retry = time.After(time.Until(a.cooldown))
		case passDone:
			a.retryDelay = 0
		}
	}
}

// pass processes pending changes.
func (a *guildActor) pass(parent context.Context) passResult {
	ctx, cancel := context.WithTimeout(parent, passTimeout)
	defer cancel()

	a.mu.Lock()
	created, updated, left, stopped := a.created, a.updated, a.left, a.stopped
	a.created, a.updated = nil, false
	a.mu.Unlock()

	es := a.w.es.GuildID(a.id.String())

	// requeue restores the pending guild create and update events, so they are
	// retried.
	requeue := func() {
		a.mu.Lock()
		if a.created == nil {
			a.created = created
		}
		a.updated = a.updated || updated
		a.mu.Unlock()
	}

	if left {
		if err := dbGuildRemoved(ctx, a.w.db, a.id); err != nil {
			if parent.Err() != nil {
				return passExit // Shutting down, handled on the next start.
			}
			a.w.logger.ErrorContext(ctx, "unable to mark guild as removed, retrying", "guild_id", a.id.String(), "error", err)
			return passRetry
		}
		es.Info("bot was removed from guild")
		return passExit
	}

	if stopped {
		return passExit
	}

	snap, err := a.snapshot(ctx)
	if err != nil {
		// e.g. the cache was reset (new gateway session), and the guild
		// create event hasn't been received yet (which triggers a new pass).
		requeue()
		a.w.logger.ErrorContext(ctx, "unable to load guild from cache", "guild_id", a.id.String(), "error", err)
		return passDone
	}
	perms := snap.permissions().guildPermissions()

	res := passDone
	if created != nil {
		if err = dbUpsertGuild(ctx, a.w.db, created, perms); err != nil {
			requeue()
			if parent.Err() != nil {
				return passExit
			}
			a.w.logger.ErrorContext(ctx, "unable to sync guild to database, retrying", "guild_id", a.id.String(), "error", err)
			return passRetry
		}
		a.perms = perms
		es.Debug("synced guild")
	}
	// The cache holds the latest guild, which may be newer than the create
	// event (when both were received since the last pass).
	if updated {
		if err = dbUpdateGuild(ctx, a.w.db, &snap.guild); err != nil {
			a.mu.Lock()
			a.updated = true
			a.mu.Unlock()
			if parent.Err() != nil {
				return passExit
			}
			a.w.logger.ErrorContext(ctx, "unable to update guild in database", "guild_id", a.id.String(), "error", err)
			res = passRetry
		}
	}

	if perms != a.perms {
		if err = dbUpdateGuildPermissions(ctx, a.w.db, a.id, perms); err != nil {
			if parent.Err() != nil {
				return passExit
			}
			a.w.logger.ErrorContext(ctx, "unable to update bot permissions in database", "guild_id", a.id.String(), "error", err)
			res = passRetry
		} else {
			es.WithFields(map[string]any{
				"old_permissions": uint64(a.perms),
				"new_permissions": uint64(perms),
			}).Debug("bot permissions changed")
			a.perms = perms
		}
	}

	return max(res, a.reconcile(parent, ctx, snap))
}

// snapshot builds a snapshot of the guild from the cache.
func (a *guildActor) snapshot(ctx context.Context) (*snapshot, error) {
	snap, err := cachedSnapshot(a.cache, a.id)
	if err != nil {
		return nil, err
	}
	if snap.bot == nil {
		snap.bot = a.fetchBotMember(ctx)
	}
	return snap, nil
}

// fetchBotMember fetches the bot's member object in the guild, and caches it.
// It is normally cached from the guild create event (which always includes the
// bot itself), and kept up to date by guild member update events for the bot
// (which are sent without the privileged guild members intent), so this is
// only a fallback.
func (a *guildActor) fetchBotMember(ctx context.Context) *discord.Member {
	me, err := a.cache.Me()
	if err != nil {
		return nil
	}

	m, err := a.api.WithContext(ctx).Member(a.id, me.ID)
	if err != nil {
		a.w.logger.WarnContext(ctx, "unable to fetch bot member", "guild_id", a.id.String(), "error", err)
		return nil
	}
	_ = a.cache.MemberSet(a.id, m, false)
	return m
}

// reconcile plans and applies the changes needed for the guild's voice
// channels. parent is the actor's context: once it is cancelled (shutting
// down), the resulting errors aren't reported to the guild.
func (a *guildActor) reconcile(parent, ctx context.Context, snap *snapshot) passResult {
	settings, err := dbGetSettings(ctx, a.w.db, a.id)

	var problems []problem
	switch {
	case errors.Is(err, errInvalidRegex):
		problems = append(problems, problem{level: slog.LevelWarn, msg: err.Error()})
	case err != nil:
		if parent.Err() != nil {
			return passExit
		}
		a.w.logger.ErrorContext(ctx, "unable to fetch guild settings, retrying", "guild_id", a.id.String(), "error", err)
		return passRetry
	}

	if !settings.enabled {
		a.report(append(problems, problem{level: slog.LevelInfo, msg: "guild is disabled, not managing channels"}))
		return passDone
	}

	p := planGuild(snap, settings)
	problems = append(problems, p.problems...)

	res := passDone
	switch {
	case p.empty():
	case time.Now().Before(a.cooldown):
		// Keep the API errors from before the cooldown, so they aren't reported
		// again once it ends.
		for _, prev := range a.problems {
			if prev.forbidden {
				problems = append(problems, prev)
			}
		}
		res = passDeferred
	default:
		a.w.es.GuildID(a.id.String()).WithFields(map[string]any{
			"deletes": len(p.deletes),
			"creates": len(p.creates),
			"moves":   len(p.moves),
			"syncs":   len(p.syncs),
		}).Debug("applying channel changes")

		applied := apply(a.api.WithContext(ctx), a.cache, a.id, p)
		problems = append(problems, applied...)

		if slices.ContainsFunc(applied, func(p problem) bool { return p.forbidden }) {
			a.cooldown = time.Now().Add(a.w.cooldown)
			res = passDeferred
		}
		if slices.ContainsFunc(applied, func(p problem) bool { return p.retry }) {
			res = passRetry
		}
	}

	if parent.Err() != nil {
		// Shutting down, so API calls may have been cancelled.
		return passExit
	}

	a.report(problems)
	return res
}

// report logs problems which weren't present in the previous pass, so
// persistent problems (e.g. missing permissions) don't flood the guild's event
// log on every voice state change. Problems that are resolved are logged too.
func (a *guildActor) report(problems []problem) {
	current := make(map[string]problem, len(problems))
	for _, p := range problems {
		current[p.key()] = p
	}

	es := a.w.es.GuildID(a.id.String())
	for k, p := range current {
		if _, ok := a.problems[k]; ok {
			continue
		}

		evt := es.WithFields(p.fields)
		switch {
		case p.level >= slog.LevelError:
			evt.Error(p.msg)
		case p.level >= slog.LevelWarn:
			evt.Warn(p.msg)
		default:
			evt.Info(p.msg)
		}
	}

	for k, p := range a.problems {
		// Only plan problems (warnings) are reported as resolved, as API errors
		// are often transient.
		if _, ok := current[k]; !ok && p.level == slog.LevelWarn {
			es.WithFields(p.fields).Info("resolved: " + p.msg)
		}
	}

	a.problems = current
}
