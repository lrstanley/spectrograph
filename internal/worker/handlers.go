// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package worker

import (
	"fmt"

	"github.com/andersfylling/disgord"
	"github.com/lrstanley/spectrograph/internal/metrics"
)

// botReady is called when the bot successfully connects to the websocket.
func (w *Worker) botReady() {}

// guildMemberUpdate is sent for current-user updates regardless of whether
// the GUILD_MEMBERS intent is set.
func (w *Worker) guildMemberUpdate(s disgord.Session, h *disgord.GuildMemberUpdate) {
	w.permissions.GuildMemberUpdate(s, h)
}

func (w *Worker) GuildRoleUpdate(s disgord.Session, h *disgord.GuildRoleUpdate) {
	w.permissions.GuildRoleChange(s, h.GuildID)
}

func (w *Worker) GuildRoleDelete(s disgord.Session, h *disgord.GuildRoleDelete) {
	w.permissions.GuildRoleChange(s, h.GuildID)
}

// This event can be sent in three different scenarios:
//   - When a user is initially connecting, to lazily load and backfill
//     information for all unavailable guilds sent in the Ready event. Guilds
//     that are unavailable due to an outage will send a Guild Delete event.
//   - When a Guild becomes available again to the client.
//   - When the current user joins a new Guild.
//   - The inner payload is a guild object, with all the extra fields specified.
func (w *Worker) guildCreate(s disgord.Session, h *disgord.GuildCreate) {
	permissions, err := w.permissions.GuildCreate(s, h)
	if err != nil {
		w.es.Guild(h.Guild).WithError(err).Error("unable to fetch permissions")
	}

	if err = w.dbUpdateGuild(h.Guild, permissions, true); err != nil {
		w.es.Guild(h.Guild).WithError(err).Error("error querying db for guild id")
		return
	}

	w.updateMu.Lock()
	if _, ok := w.updates[h.Guild.ID]; !ok {
		ch := make(chan *updateEvent)
		w.updates[h.Guild.ID] = ch

		// Process the first event.
		w.processUpdateWorker(s, h.Guild.ID, &updateEvent{sess: s, event: h, guild: h.Guild})

		// Then start the event watcher for all subsequent events.
		metrics.WorkerGuildCount.WithLabelValues(fmt.Sprintf("%d", h.ShardID)).Inc()
		go w.eventWatcher(s, h.Guild.ID, ch)
	}
	w.updateMu.Unlock()
}

// Sent when a guild is updated. The inner payload is a guild object.
func (w *Worker) guildUpdate(s disgord.Session, h *disgord.GuildUpdate) {
	permissions, err := w.permissions.GuildUpdate(s, h)
	if err != nil {
		w.es.Guild(h.Guild).WithError(err).Error("unable to fetch permissions")
	}

	if err = w.dbUpdateGuild(h.Guild, permissions, true); err != nil {
		w.es.Guild(h.Guild).WithError(err).Error("error querying db for guild id")
		return
	}
}

// Sent when a guild becomes or was already unavailable due to an outage, or
// when the user leaves or is removed from a guild. The inner payload is an
// unavailable guild object. If the unavailable field is not set, the user
// was removed from the guild.
func (w *Worker) guildDelete(s disgord.Session, h *disgord.GuildDelete) {
	w.updateMu.Lock()
	if ch, ok := w.updates[h.UnavailableGuild.ID]; ok {
		metrics.WorkerGuildCount.WithLabelValues(fmt.Sprintf("%d", h.ShardID)).Dec()

		close(ch)

		delete(w.updates, h.UnavailableGuild.ID)
	}
	w.updateMu.Unlock()

	if h.UserWasRemoved() {
		w.es.GuildID(h.UnavailableGuild.ID.String()).Info("bot was removed from guild, disabling via configuration")
		w.permissions.RemoveGuild(h.UnavailableGuild.ID)

		// TODO: maybe have it send a notification?
		if err := w.dbDisableGuild(h.UnavailableGuild.ID.String()); err != nil {
			w.es.GuildID(h.UnavailableGuild.ID.String()).WithError(err).Error("error disabling guild")
		}
	}
}

func (w *Worker) voiceStateUpdate(s disgord.Session, h *disgord.VoiceStateUpdate) {
	w.routeEvent(s, h, h.GuildID)
}

func (w *Worker) channelCreate(s disgord.Session, h *disgord.ChannelCreate) {
	w.routeEvent(s, h, h.Channel.GuildID)
}

func (w *Worker) channelDelete(s disgord.Session, h *disgord.ChannelDelete) {
	w.routeEvent(s, h, h.Channel.GuildID)
}

func (w *Worker) channelUpdate(s disgord.Session, h *disgord.ChannelUpdate) {
	w.routeEvent(s, h, h.Channel.GuildID)
}
