// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"github.com/andersfylling/disgord"
	"github.com/lrstanley/spectrograph/internal/discordapi"
	"github.com/lrstanley/spectrograph/internal/models"
)

// botReady is called when the bot successfully connects to the websocket.
func (b *discordBot) botReady() {}

// guildMemberUpdate is sent for current-user updates regardless of whether
// the GUILD_MEMBERS intent is set.
func (b *discordBot) guildMemberUpdate(s disgord.Session, h *disgord.GuildMemberUpdate) {
	b.permissions.guildMemberUpdate(s, h)
}

func (b *discordBot) GuildRoleUpdate(s disgord.Session, h *disgord.GuildRoleUpdate) {
	b.permissions.guildRoleChange(s, h.GuildID)
}

func (b *discordBot) GuildRoleDelete(s disgord.Session, h *disgord.GuildRoleDelete) {
	b.permissions.guildRoleChange(s, h.GuildID)
}

// This event can be sent in three different scenarios:
//   - When a user is initially connecting, to lazily load and backfill
//     information for all unavailable guilds sent in the Ready event. Guilds
//     that are unavailable due to an outage will send a Guild Delete event.
//   - When a Guild becomes available again to the client.
//   - When the current user joins a new Guild.
//   - The inner payload is a guild object, with all the extra fields specified.
func (b *discordBot) guildCreate(s disgord.Session, h *disgord.GuildCreate) {
	permissions, err := b.permissions.guildCreate(s, h)
	if err != nil {
		logGuild(logger, h.Guild).WithError(err).Error("unable to fetch permissions")
	}

	server, err := svcServers.Get(b.ctx, h.Guild.ID.String())
	if err != nil {
		if !models.IsNotFound(err) {
			logGuild(logger, h.Guild).WithError(err).Error("error querying db for guild id")
			return
		}
		err = nil
		server = &models.Server{}
	}

	server.Discord = &models.ServerDiscordData{
		ID:                 h.Guild.ID.String(),
		Name:               h.Guild.Name,
		Features:           h.Guild.Features,
		Icon:               h.Guild.Icon,
		IconUrl:            discordapi.GenerateGuildIconURL(h.Guild.ID.String(), h.Guild.Icon),
		JoinedAt:           h.Guild.JoinedAt.Time,
		Large:              h.Guild.Large,
		MemberCount:        int64(h.Guild.MemberCount),
		OwnerID:            h.Guild.OwnerID.String(),
		Permissions:        models.DiscordPermissions(permissions),
		Region:             h.Guild.Region,
		SystemChannelFlags: h.Guild.SystemChannelID.String(),
	}

	if err = svcServers.Upsert(b.ctx, server); err != nil {
		logGuild(logger, server).WithError(err).Error("unable to update server in db")
	}

	b.updateMu.Lock()
	if _, ok := b.updates[h.Guild.ID]; !ok {
		ch := make(chan *updateEvent)
		b.updates[h.Guild.ID] = ch

		// Process the first event.
		b.processUpdateWorker(s, h.Guild.ID, &updateEvent{sess: s, event: h, guild: h.Guild})

		// Then start the event watcher for all subsequent events.
		go b.eventWatcher(s, h.Guild.ID, ch)
	}
	b.updateMu.Unlock()
}

// Sent when a guild is updated. The inner payload is a guild object.
func (b *discordBot) guildUpdate(s disgord.Session, h *disgord.GuildUpdate) {
	// TODO: Crawl through roles and grab our permissions from it.
	// pretty.Println(h)
}

// Sent when a guild becomes or was already unavailable due to an outage, or
// when the user leaves or is removed from a guild. The inner payload is an
// unavailable guild object. If the unavailable field is not set, the user
// was removed from the guild.
func (b *discordBot) guildDelete(s disgord.Session, h *disgord.GuildDelete) {
	b.updateMu.Lock()
	if ch, ok := b.updates[h.UnavailableGuild.ID]; ok {
		close(ch)

		delete(b.updates, h.UnavailableGuild.ID)
	}
	b.updateMu.Unlock()

	if h.UserWasRemoved() {
		// TODO: clean up from db, maybe have it send a notification?
		b.permissions.removeGuild(h.UnavailableGuild.ID)
	}
}

func (b *discordBot) voiceStateUpdate(s disgord.Session, h *disgord.VoiceStateUpdate) {
	b.routeEvent(s, h, h.GuildID)
}

func (b *discordBot) channelCreate(s disgord.Session, h *disgord.ChannelCreate) {
	b.routeEvent(s, h, h.Channel.GuildID)
}

func (b *discordBot) channelDelete(s disgord.Session, h *disgord.ChannelDelete) {
	b.routeEvent(s, h, h.Channel.GuildID)
}

func (b *discordBot) channelUpdate(s disgord.Session, h *disgord.ChannelUpdate) {
	b.routeEvent(s, h, h.Channel.GuildID)
}
