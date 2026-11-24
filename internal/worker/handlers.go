// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package worker

import (
	"errors"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/state/store"
	"github.com/diamondburned/arikawa/v3/utils/ws"
)

// registerHandlers registers the gateway event handlers for the shard. The
// handlers are synchronous (run in the gateway event loop, after the state
// cache has been updated), so they must not block: they only record what
// changed, and wake the guild's actor.
func (w *Worker) registerHandlers(s *shard) {
	st := s.state

	st.AddSyncHandler(func(ev *ws.BackgroundErrorEvent) {
		// Discord regularly adds gateway events that arikawa doesn't know about
		// yet (e.g. VOICE_CHANNEL_START_TIME_UPDATE). They're skipped, and none
		// of them matter to the worker.
		if isUnknownEvent(ev.Err) {
			w.logger.Debug("skipping unknown gateway event", "shard_id", s.id, "error", ev.Err)
			return
		}
		w.logger.Warn("gateway background error", "shard_id", s.id, "error", ev.Err)
	})

	st.AddSyncHandler(func(ev *gateway.ReadyEvent) {
		w.logger.Info(
			"gateway ready",
			"shard_id", s.id,
			"user", ev.User.Tag(),
			"guilds", len(ev.Guilds),
		)
		w.shardReady(s, ev)
		w.shardReadyCommands(s, ev.Application.ID)
	})

	// Slash commands, buttons, and modals.
	st.AddSyncHandler(func(ev *gateway.InteractionCreateEvent) {
		w.interactionCreated(s, &ev.InteractionEvent)
	})

	// Sent when the bot connects (once per guild), when a guild becomes
	// available again after an outage, and when the bot joins a guild. See
	// [state.GuildReadyEvent], [state.GuildAvailableEvent], and
	// [state.GuildJoinEvent].
	st.AddSyncHandler(func(ev *gateway.GuildCreateEvent) {
		if ev.Unavailable {
			return
		}
		syncGuildCache(st.Cabinet, ev)
		if a := w.startActor(s, ev.ID); a != nil {
			a.guildCreated(ev)
		}
	})

	st.AddSyncHandler(func(ev *gateway.GuildUpdateEvent) {
		if a := w.actor(ev.ID); a != nil {
			a.guildUpdated()
		}
	})

	st.AddSyncHandler(func(ev *state.GuildUnavailableEvent) {
		if a := w.actor(ev.ID); a != nil {
			a.stop(false)
		}
	})

	// Sent when the bot leaves or is removed from a guild, or the guild is
	// deleted.
	st.AddSyncHandler(func(ev *state.GuildLeaveEvent) {
		removeGuildCache(st.Cabinet, ev.ID)
		w.guildLeft(s, ev.ID)
	})

	// Events which (may) change the desired state of a guild's channels, or
	// the bot's permissions.
	st.AddSyncHandler(func(ev *gateway.VoiceStateUpdateEvent) { w.pokeGuild(ev.GuildID) })
	st.AddSyncHandler(func(ev *gateway.ChannelCreateEvent) { w.pokeGuild(ev.GuildID) })
	st.AddSyncHandler(func(ev *gateway.ChannelUpdateEvent) { w.pokeGuild(ev.GuildID) })
	st.AddSyncHandler(func(ev *gateway.ChannelDeleteEvent) { w.pokeGuild(ev.GuildID) })
	st.AddSyncHandler(func(ev *gateway.GuildRoleCreateEvent) { w.pokeGuild(ev.GuildID) })
	st.AddSyncHandler(func(ev *gateway.GuildRoleUpdateEvent) { w.pokeGuild(ev.GuildID) })
	st.AddSyncHandler(func(ev *gateway.GuildRoleDeleteEvent) { w.pokeGuild(ev.GuildID) })

	// Sent for the bot's own member, even without the guild members intent.
	st.AddSyncHandler(func(ev *gateway.GuildMemberUpdateEvent) { w.pokeGuild(ev.GuildID) })
}

// syncGuildCache makes the cached channels, roles, members, and voice states
// of the guild match the guild create event. The state cache only adds entries
// on guild create, without updating or removing existing ones, so without this,
// anything changed while the guild was unavailable would be stale (e.g. a user
// who moved to another channel would still appear in the old one, so the
// channel they are in could be deleted), and anything removed would linger
// (e.g. a stale voice state would keep a channel occupied).
func syncGuildCache(cab *store.Cabinet, ev *gateway.GuildCreateEvent) {
	channels := make(map[discord.ChannelID]struct{}, len(ev.Channels)+len(ev.Threads))
	for _, ch := range ev.Channels {
		channels[ch.ID] = struct{}{}
	}
	for _, ch := range ev.Threads {
		channels[ch.ID] = struct{}{}
	}
	if cached, err := cab.Channels(ev.ID); err == nil {
		for _, ch := range cached {
			if _, ok := channels[ch.ID]; !ok {
				_ = cab.ChannelRemove(&ch)
			}
		}
	}

	roles := make(map[discord.RoleID]struct{}, len(ev.Roles))
	for _, r := range ev.Roles {
		roles[r.ID] = struct{}{}
		_ = cab.RoleSet(ev.ID, &r, true)
	}
	if cached, err := cab.Roles(ev.ID); err == nil {
		for _, r := range cached {
			if _, ok := roles[r.ID]; !ok {
				_ = cab.RoleRemove(ev.ID, r.ID)
			}
		}
	}

	// Without the guild members intent, only some members are sent (including
	// the bot itself), so members which aren't in the event are kept.
	for _, m := range ev.Members {
		_ = cab.MemberSet(ev.ID, &m, true)
	}

	users := make(map[discord.UserID]struct{}, len(ev.VoiceStates))
	for _, vs := range ev.VoiceStates {
		users[vs.UserID] = struct{}{}
		vs.GuildID = ev.ID
		_ = cab.VoiceStateSet(ev.ID, &vs, true)
	}
	if cached, err := cab.VoiceStates(ev.ID); err == nil {
		for _, vs := range cached {
			if _, ok := users[vs.UserID]; !ok {
				_ = cab.VoiceStateRemove(ev.ID, vs.UserID)
			}
		}
	}
}

// removeGuildCache removes everything cached for the guild (the state cache
// only removes the guild itself).
func removeGuildCache(cab *store.Cabinet, id discord.GuildID) {
	if channels, err := cab.Channels(id); err == nil {
		for _, ch := range channels {
			_ = cab.ChannelRemove(&ch)
		}
	}
	if roles, err := cab.Roles(id); err == nil {
		for _, r := range roles {
			_ = cab.RoleRemove(id, r.ID)
		}
	}
	if states, err := cab.VoiceStates(id); err == nil {
		for _, vs := range states {
			_ = cab.VoiceStateRemove(id, vs.UserID)
		}
	}
	if members, err := cab.Members(id); err == nil {
		for _, m := range members {
			_ = cab.MemberRemove(id, m.User.ID)
		}
	}
}

// isUnknownEvent returns true if err is from a gateway event that arikawa
// doesn't know how to decode. arikawa creates [ws.UnknownEventError] as a
// value, while its own [ws.IsUnknownEvent] only matches a pointer, so both are
// checked.
func isUnknownEvent(err error) bool {
	var unknown ws.UnknownEventError
	var unknownPtr *ws.UnknownEventError
	return errors.As(err, &unknown) || errors.As(err, &unknownPtr)
}
