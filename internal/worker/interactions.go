// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package worker

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state/store"
	"github.com/diamondburned/arikawa/v3/utils/httputil"
	"github.com/lrstanley/spectrograph/internal/database"
	"github.com/lrstanley/spectrograph/internal/database/ent"
	"github.com/lrstanley/spectrograph/internal/database/ent/guild"
	"github.com/lrstanley/spectrograph/internal/database/ent/guildconfig"
	"github.com/lrstanley/spectrograph/internal/database/ent/guildevent"
	"github.com/lrstanley/spectrograph/internal/database/ent/privacy"
	"github.com/lrstanley/spectrograph/internal/database/ent/user"
	"github.com/lrstanley/spectrograph/internal/guildlogger"
	"github.com/lrstanley/spectrograph/internal/models"
)

const (
	// deferAfter is how long an interaction is handled before it is
	// acknowledged (deferred), and the response sent once ready. Discord
	// requires a response within 3 seconds.
	deferAfter = 2 * time.Second

	// interactionTimeout bounds handling an interaction (Discord accepts
	// follow-ups for 15 minutes after a deferred response).
	interactionTimeout = time.Minute
)

// Error messages shown to users.
const (
	msgNotInGuild    = "Spectrograph's commands only work in a server."
	msgNotAuthorized = "You need the **Manage Server** permission to manage Spectrograph in this server."
	msgBanned        = "You're banned from using Spectrograph."
	msgNotSetUp      = "Spectrograph hasn't finished setting up this server. Try again in a minute."
	msgUnknown       = "This action isn't supported anymore. Run `/" + commandName + " " + subcommandStatus + "` again."
	msgInternal      = "Something went wrong. Try again in a moment."
	msgAdminDisabled = "Spectrograph admins disabled Spectrograph in this server, so it can't be enabled here. " +
		"Contact the Spectrograph admins if you think this is a mistake."
	msgSlowModal = "Loading the settings took too long. Try again."
)

// guildCache is the subset of the state cache used to show the status of a
// guild. Implemented by [store.Cabinet].
type guildCache interface {
	Guild(id discord.GuildID) (*discord.Guild, error)
	Channels(id discord.GuildID) ([]discord.Channel, error)
	Roles(id discord.GuildID) ([]discord.Role, error)
	VoiceStates(id discord.GuildID) ([]discord.VoiceState, error)
	Me() (*discord.User, error)
	Member(guildID discord.GuildID, userID discord.UserID) (*discord.Member, error)
}

// cachedSnapshot builds a snapshot of the guild from the state cache. The bot
// member is nil if it isn't cached.
func cachedSnapshot(cache guildCache, id discord.GuildID) (*snapshot, error) {
	g, err := cache.Guild(id)
	if err != nil {
		return nil, fmt.Errorf("guild: %w", err)
	}

	snap := &snapshot{guild: *g}

	if snap.channels, err = cache.Channels(id); err != nil && !errors.Is(err, store.ErrNotFound) {
		return nil, fmt.Errorf("channels: %w", err)
	}
	if snap.roles, err = cache.Roles(id); err != nil && !errors.Is(err, store.ErrNotFound) {
		return nil, fmt.Errorf("roles: %w", err)
	}
	if snap.voiceStates, err = cache.VoiceStates(id); err != nil && !errors.Is(err, store.ErrNotFound) {
		return nil, fmt.Errorf("voice states: %w", err)
	}

	if me, merr := cache.Me(); merr == nil {
		snap.bot, _ = cache.Member(id, me.ID)
	}
	return snap, nil
}

// interactionResponder sends responses to interactions. Implemented by
// [apiResponder].
type interactionResponder interface {
	// respond sends the initial response.
	respond(ctx context.Context, ev *discord.InteractionEvent, resp *api.InteractionResponse) error
	// editOriginal edits the message of the initial (deferred) response.
	editOriginal(ctx context.Context, ev *discord.InteractionEvent, data *api.InteractionResponseData) error
	// followUp sends a new message, after the initial response.
	followUp(ctx context.Context, ev *discord.InteractionEvent, data *api.InteractionResponseData) error
}

// apiResponder responds to interactions through the Discord REST API.
type apiResponder struct {
	client *api.Client
}

func (r apiResponder) respond(ctx context.Context, ev *discord.InteractionEvent, resp *api.InteractionResponse) error {
	return redactToken(r.client.WithContext(ctx).RespondInteraction(ev.ID, ev.Token, *resp), ev.Token)
}

// redactToken removes the interaction token from err, as arikawa includes the
// request URL (which contains the token) in some errors, and they are logged.
func redactToken(err error, token string) error {
	if err == nil || token == "" {
		return err
	}
	return errors.New(strings.ReplaceAll(err.Error(), token, "[redacted]"))
}

// editOriginal edits the original response. [api.Client.EditInteractionResponse]
// can't set message flags, which are needed to switch a deferred reply to
// components (v2), so the request is made directly.
func (r apiResponder) editOriginal(ctx context.Context, ev *discord.InteractionEvent, data *api.InteractionResponseData) error {
	body := struct {
		Components      *discord.TopLevelComponents `json:"components,omitempty"`
		AllowedMentions *api.AllowedMentions        `json:"allowed_mentions,omitempty"`
		Flags           discord.MessageFlags        `json:"flags,omitempty"`
	}{
		Components:      data.Components,
		AllowedMentions: data.AllowedMentions,
		Flags:           data.Flags & discord.IsComponentsV2, // Only flag allowed when editing.
	}

	return redactToken(r.client.WithContext(ctx).FastRequest(
		"PATCH",
		api.EndpointWebhooks+ev.AppID.String()+"/"+ev.Token+"/messages/@original",
		httputil.WithJSONBody(body),
	), ev.Token)
}

func (r apiResponder) followUp(ctx context.Context, ev *discord.InteractionEvent, data *api.InteractionResponseData) error {
	_, err := r.client.WithContext(ctx).FollowUpInteraction(ev.AppID, ev.Token, *data)
	return redactToken(err, ev.Token)
}

// interactionHandler handles the bot's slash command, and its message
// component and modal interactions.
type interactionHandler struct {
	logger *slog.Logger
	db     *ent.Client
	es     *guildlogger.EventStream

	// dashboardURL is the base URL of the dashboard, linked from replies (if
	// set).
	dashboardURL *url.URL

	// deferAfter is how long to wait for a response before deferring it.
	deferAfter time.Duration
}

// serve handles an interaction, and sends the response.
func (h *interactionHandler) serve(ctx context.Context, r interactionResponder, cache guildCache, ev *discord.InteractionEvent) {
	logger := h.logger.With("interaction_id", ev.ID.String(), "guild_id", ev.GuildID.String(), "user_id", ev.SenderID().String())
	serveInteraction(ctx, logger, r, ev, h.deferAfter, func(ctx context.Context) *api.InteractionResponse {
		return h.handle(ctx, cache, ev)
	})
}

// serveInteraction runs handle, and sends the response it returns (if any). If
// it doesn't return within deferAfter, the interaction is acknowledged first
// (a deferred response, as Discord requires one within 3 seconds), and the
// response is sent once ready.
func serveInteraction(
	ctx context.Context,
	logger *slog.Logger,
	r interactionResponder,
	ev *discord.InteractionEvent,
	deferAfter time.Duration,
	handle func(ctx context.Context) *api.InteractionResponse,
) {
	ctx, cancel := context.WithTimeout(ctx, interactionTimeout)
	defer cancel()

	result := make(chan *api.InteractionResponse, 1)
	go func() { result <- handle(ctx) }()

	timer := time.NewTimer(deferAfter)
	defer timer.Stop()

	select {
	case resp := <-result:
		if resp == nil {
			return
		}
		if err := r.respond(ctx, ev, resp); err != nil {
			logger.WarnContext(ctx, "unable to respond to interaction", "error", err)
		}
		return
	case <-timer.C:
	case <-ctx.Done():
		return
	}

	// Component interactions (and modals submitted from one) update the
	// message the component is attached to, commands reply with a new message.
	// Deferred replies can only be made ephemeral: they are switched to
	// components (v2) when edited.
	updates := ev.Message != nil
	deferred := &api.InteractionResponse{
		Type: api.DeferredMessageInteractionWithSource,
		Data: &api.InteractionResponseData{Flags: discord.EphemeralMessage},
	}
	if updates {
		deferred = &api.InteractionResponse{Type: api.DeferredMessageUpdate}
	}
	if err := r.respond(ctx, ev, deferred); err != nil {
		logger.WarnContext(ctx, "unable to defer interaction response", "error", err)
		return
	}

	var resp *api.InteractionResponse
	select {
	case resp = <-result:
	case <-ctx.Done():
		logger.WarnContext(ctx, "timed out handling interaction", "error", ctx.Err())
		return
	}
	if resp == nil {
		return
	}

	// Modals can only be opened as the initial response.
	if resp.Type == api.ModalResponse {
		resp = errorMessage(msgSlowModal)
	}

	var err error
	if updates && resp.Type != api.UpdateMessage {
		err = r.followUp(ctx, ev, resp.Data)
	} else {
		err = r.editOriginal(ctx, ev, resp.Data)
	}
	if err != nil {
		logger.WarnContext(ctx, "unable to send deferred interaction response", "error", err)
	}
}

// actionKind is what an interaction asks for.
type actionKind int

const (
	actionUnknown actionKind = iota
	actionStatus
	actionSettings
	actionSaveSettings
	actionEnable
	actionDisable
	actionEvents
	actionHelp
)

// action is a parsed interaction.
type action struct {
	kind actionKind

	// update is true if the response should update the message the
	// interaction's component is attached to, rather than send a new message.
	update bool

	// count is the number of events to show (for actionEvents).
	count int

	// values are the submitted modal values (for actionSaveSettings).
	values map[discord.ComponentID]string
}

// parseAction returns what the interaction asks for. It returns ok=false for
// interactions which aren't for the bot's command or components.
func parseAction(ev *discord.InteractionEvent) (act action, ok bool) {
	switch data := ev.Data.(type) {
	case *discord.CommandInteraction:
		if data.Name != commandName || len(data.Options) == 0 {
			return action{}, false
		}
		sub := data.Options[0]
		switch sub.Name {
		case subcommandStatus:
			act.kind = actionStatus
		case subcommandSettings:
			act.kind = actionSettings
		case subcommandEnable:
			act.kind = actionEnable
		case subcommandDisable:
			act.kind = actionDisable
		case subcommandHelp:
			act.kind = actionHelp
		case subcommandEvents:
			act.kind = actionEvents
			act.count = defaultEventCount
			if opt := sub.Options.Find(optionEventCount); opt.Name != "" {
				if n, err := opt.IntValue(); err == nil {
					act.count = int(min(max(n, 1), maxEventCount))
				}
			}
		}
		return act, true
	case *discord.ButtonInteraction:
		if !strings.HasPrefix(string(data.CustomID), customIDPrefix) {
			return action{}, false
		}
		act.update = true
		switch data.CustomID {
		case customIDRefresh:
			act.kind = actionStatus
		case customIDEditSettings:
			act.kind = actionSettings
		case customIDEnable:
			act.kind = actionEnable
		case customIDDisable:
			act.kind = actionDisable
		}
		return act, true
	case *discord.ModalInteraction:
		if !strings.HasPrefix(string(data.CustomID), customIDPrefix) {
			return action{}, false
		}
		if data.CustomID == customIDSettingsModal {
			act.kind = actionSaveSettings
			act.update = ev.Message != nil
			act.values = modalValues(data.Components)
		}
		return act, true
	default:
		return action{}, false
	}
}

// canManageGuild returns true if the member can manage Spectrograph in the
// guild, the same as who the dashboard considers a guild admin: the guild
// owner, and members with the Administrator or Manage Server permission.
// Guild admins can change who can use the command (see [commandPermissions]),
// so this is checked for every interaction.
//
// The member (and their roles) are the ones sent with the interaction, and the
// role permissions are from the state cache (kept up to date by role events).
// Discord also sends the member's resolved permissions, but arikawa doesn't
// decode them.
func canManageGuild(g *discord.Guild, roles []discord.Role, member *discord.Member) bool {
	perms := discord.CalcOverrides(*g, discord.Channel{}, *member, roles)
	return g.OwnerID == member.User.ID ||
		perms.Has(discord.PermissionAdministrator) ||
		perms.Has(discord.PermissionManageGuild)
}

// handle handles an interaction, and returns the response (nil if the
// interaction isn't for the bot).
func (h *interactionHandler) handle(ctx context.Context, cache guildCache, ev *discord.InteractionEvent) *api.InteractionResponse {
	act, ok := parseAction(ev)
	if !ok {
		return nil
	}
	if act.kind == actionUnknown {
		return errorMessage(msgUnknown)
	}

	if !ev.GuildID.IsValid() || ev.Member == nil {
		return errorMessage(msgNotInGuild)
	}

	cached, err := cache.Guild(ev.GuildID)
	if err != nil {
		// Interactions are sent to the shard receiving the guild's events, so
		// it's cached unless it just became available.
		return errorMessage(msgNotSetUp)
	}
	roles, err := cache.Roles(ev.GuildID)
	if err != nil && !errors.Is(err, store.ErrNotFound) {
		return h.internalError(ctx, ev, "unable to load guild roles", err)
	}
	if !canManageGuild(cached, roles, ev.Member) {
		return errorMessage(msgNotAuthorized)
	}

	// The member is a guild admin, so privacy rules (which only know about
	// users logged in to the dashboard) are bypassed from here on.
	ctx = privacy.DecisionContext(ctx, privacy.Allow)

	banned, err := h.db.User.Query().
		Where(user.UserID(ev.Member.User.ID.String()), user.Banned(true)).
		Exist(ctx)
	if err != nil {
		return h.internalError(ctx, ev, "unable to check if user is banned", err)
	}
	if banned {
		return errorMessage(msgBanned)
	}

	if act.kind == actionHelp {
		return message(helpComponents(h.guildURL(0))...)
	}

	g, err := h.db.Guild.Query().
		Where(guild.GuildID(ev.GuildID.String())).
		WithGuildConfig().
		WithGuildAdminConfig().
		Only(ctx)
	switch {
	case ent.IsNotFound(err):
		return errorMessage(msgNotSetUp)
	case err != nil:
		return h.internalError(ctx, ev, "unable to load guild", err)
	case g.Edges.GuildConfig == nil || g.Edges.GuildAdminConfig == nil:
		return errorMessage(msgNotSetUp)
	}

	switch act.kind {
	case actionStatus:
		return h.status(cache, ev, g, "", act.update)
	case actionSettings:
		return settingsModal(g.Edges.GuildConfig, g.Edges.GuildAdminConfig)
	case actionEnable, actionDisable:
		return h.setEnabled(ctx, cache, ev, g, act)
	case actionSaveSettings:
		return h.saveSettings(ctx, cache, ev, g, act)
	case actionEvents:
		return h.events(ctx, ev, g, act.count)
	case actionUnknown, actionHelp: // Handled above.
	}
	return errorMessage(msgUnknown)
}

// internalError logs err, and returns a generic error message.
func (h *interactionHandler) internalError(ctx context.Context, ev *discord.InteractionEvent, msg string, err error) *api.InteractionResponse {
	h.logger.ErrorContext(ctx, "unable to handle interaction",
		"step", msg,
		"interaction_id", ev.ID.String(),
		"guild_id", ev.GuildID.String(),
		"error", err,
	)
	return errorMessage(msgInternal)
}

// guildURL returns the URL of the guild's dashboard page (or of the dashboard
// itself, if id is 0), or an empty string if the dashboard URL isn't known.
func (h *interactionHandler) guildURL(id int) string {
	if h.dashboardURL == nil {
		return ""
	}
	if id == 0 {
		return h.dashboardURL.JoinPath("dashboard").String()
	}
	return h.dashboardURL.JoinPath("dashboard", "guild", strconv.Itoa(id)).String()
}

// status returns the status view of the guild, with an optional notice.
func (h *interactionHandler) status(cache guildCache, ev *discord.InteractionEvent, g *ent.Guild, notice string, updateMessage bool) *api.InteractionResponse {
	cfg, admin := g.Edges.GuildConfig, g.Edges.GuildAdminConfig

	settings, err := effectiveSettings(cfg, admin)
	v := &guildView{
		name:         g.Name,
		config:       cfg,
		admin:        admin,
		settings:     settings,
		regexErr:     err,
		dashboardURL: h.guildURL(g.ID),
	}

	if snap, serr := cachedSnapshot(cache, ev.GuildID); serr == nil {
		v.name = snap.guild.Name
		v.groups, v.unmanaged = newGroupViews(snap, settings)
		v.missing = snap.permissions().missing(0, discord.PermissionViewChannel|discord.PermissionManageChannels)
	}

	if updateMessage {
		return update(statusComponents(v, notice)...)
	}
	return message(statusComponents(v, notice)...)
}

// changeLogger returns the guild's event logger, annotated with the user who
// made a change.
func (h *interactionHandler) changeLogger(ev *discord.InteractionEvent, fields map[string]any) *guildlogger.EventStream {
	return h.es.GuildID(ev.GuildID.String()).
		WithFields(fields).
		WithFields(map[string]any{"user_id": ev.Member.User.ID.String(), "username": ev.Member.User.Username})
}

// setEnabled enables or disables the bot in the guild. Guilds disabled by
// system admins can't be enabled.
func (h *interactionHandler) setEnabled(ctx context.Context, cache guildCache, ev *discord.InteractionEvent, g *ent.Guild, act action) *api.InteractionResponse {
	enable := act.kind == actionEnable
	cfg := g.Edges.GuildConfig

	if enable && !g.Edges.GuildAdminConfig.Enabled {
		return errorMessage(msgAdminDisabled)
	}

	state := "disabled"
	if enable {
		state = "enabled"
	}
	if cfg.Enabled == enable {
		return h.status(cache, ev, g, "Spectrograph is already "+state+".", act.update)
	}

	cfg, err := h.db.GuildConfig.UpdateOne(cfg).SetEnabled(enable).Save(ctx)
	if err != nil {
		return h.internalError(ctx, ev, "unable to update guild config", err)
	}
	g.Edges.GuildConfig = cfg

	// The guild's actor is woken by the config observer (see
	// [Worker.watchConfigs]), the same as for changes made on the dashboard.
	h.changeLogger(ev, map[string]any{guildconfig.FieldEnabled: enable}).
		Info("guild " + state + " via Discord by " + ev.Member.User.Username)

	return h.status(cache, ev, g, "Spectrograph is now "+state+".", act.update)
}

// settingsFieldLabels are the labels of the settings modal's inputs, by field.
var settingsFieldLabels = map[string]string{
	guildconfig.FieldRegexMatch:       "Channel name pattern",
	guildconfig.FieldDefaultMaxClones: "Max clones per channel",
	guildconfig.FieldContactEmail:     "Contact email",
}

// saveSettings saves the settings submitted with the settings modal. Values
// are validated by the guild config schema, the same as on the dashboard.
func (h *interactionHandler) saveSettings(ctx context.Context, cache guildCache, ev *discord.InteractionEvent, g *ent.Guild, act action) *api.InteractionResponse {
	cfg := g.Edges.GuildConfig
	changes := map[string]any{}
	upd := h.db.GuildConfig.UpdateOne(cfg)

	if v, ok := act.values[customIDRegexMatch]; ok && strings.TrimSpace(v) != cfg.RegexMatch {
		v = strings.TrimSpace(v)
		upd.SetRegexMatch(v)
		changes[guildconfig.FieldRegexMatch] = v
	}

	if v, ok := act.values[customIDMaxClones]; ok {
		n, err := strconv.Atoi(v)
		if err != nil || n < 0 || n > maxCloneLimit {
			return errorMessage(
				fmt.Sprintf("Couldn't save the settings. **%s:** choose a number from 0 to %d.", settingsFieldLabels[guildconfig.FieldDefaultMaxClones], maxCloneLimit),
				editSettingsButton(),
			)
		}
		if n != cfg.DefaultMaxClones {
			upd.SetDefaultMaxClones(n)
			changes[guildconfig.FieldDefaultMaxClones] = n
		}
	}

	if v, ok := act.values[customIDContactEmail]; ok && strings.TrimSpace(v) != cfg.ContactEmail {
		v = strings.TrimSpace(v)
		upd.SetContactEmail(v)
		changes[guildconfig.FieldContactEmail] = v
	}

	if len(changes) == 0 {
		return h.status(cache, ev, g, "No changes to save.", act.update)
	}

	cfg, err := upd.Save(ctx)
	if fe := models.ClientError(err); fe != nil {
		label := settingsFieldLabels[fe.Field]
		if label == "" {
			label = "Settings"
		}
		return errorMessage(
			fmt.Sprintf("Couldn't save the settings. **%s:** %s.", label, escapeMarkdown(fe.Reason)),
			editSettingsButton(),
		)
	}
	if err != nil {
		return h.internalError(ctx, ev, "unable to update guild config", err)
	}
	g.Edges.GuildConfig = cfg

	h.changeLogger(ev, changes).Info("settings updated via Discord by " + ev.Member.User.Username)
	return h.status(cache, ev, g, "Settings saved.", act.update)
}

// events returns the guild's most recent events (excluding debug events).
func (h *interactionHandler) events(ctx context.Context, ev *discord.InteractionEvent, g *ent.Guild, count int) *api.InteractionResponse {
	events, err := h.db.GuildEvent.Query().
		Where(
			guildevent.HasGuildWith(guild.ID(g.ID)),
			guildevent.TypeNEQ(guildevent.TypeDEBUG),
			guildevent.CreateTimeGT(time.Now().Add(-database.MaxGuildEventAge)),
		).
		Order(ent.Desc(guildevent.FieldID)).
		Limit(count).
		All(ctx)
	if err != nil {
		return h.internalError(ctx, ev, "unable to load guild events", err)
	}
	return message(eventComponents(g.Name, events, h.guildURL(g.ID))...)
}
