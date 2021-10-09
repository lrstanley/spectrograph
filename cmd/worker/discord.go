// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"context"
	"errors"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/apex/log"
	"github.com/lrstanley/spectrograph/internal/discordapi"
	"github.com/lrstanley/spectrograph/internal/models"
)

type discordBot struct {
	ctx    context.Context
	errs   chan<- error
	client *disgord.Client

	updateMu sync.RWMutex
	updates  map[disgord.Snowflake]chan *updateEvent
}

func (b *discordBot) setup(wg *sync.WaitGroup) {
	b.updates = make(map[disgord.Snowflake]chan *updateEvent)

	b.client = disgord.New(disgord.Config{
		BotToken:    cli.Discord.BotToken,
		ProjectName: "spectrograph (https://github.com/lrstanley, https://liam.sh)",
		Presence: &disgord.UpdateStatusPayload{
			Since: nil,
			Game: []*disgord.Activity{
				{Name: "voice chan curator", Type: disgord.ActivityTypeGame},
			},
			Status: disgord.StatusOnline,
			AFK:    false,
		},
		Logger: &LoggerApex{logger: logger},
		RejectEvents: disgord.AllEventsExcept(
			// See: https://github.com/andersfylling/disgord/issues/360#issuecomment-830918707
			disgord.EvtReady,
			disgord.EvtResumed,
			disgord.EvtGuildCreate,
			disgord.EvtGuildUpdate,
			disgord.EvtGuildDelete,
			disgord.EvtGuildRoleCreate,
			disgord.EvtGuildRoleUpdate,
			disgord.EvtGuildRoleDelete,
			disgord.EvtChannelCreate,
			disgord.EvtChannelUpdate,
			disgord.EvtChannelDelete,
			disgord.EvtVoiceServerUpdate,
			disgord.EvtVoiceStateUpdate,
		),

		ShardConfig: disgord.ShardConfig{
			ShardIDs:   []uint{},
			ShardCount: 0,
		},
	})
	gw := b.client.Gateway().WithContext(b.ctx)
	// TODO: custom logger implementation that contains prefix info?

	// Register hooks here.
	gw.BotReady(b.botReady)
	gw.GuildCreate(b.guildCreate)
	gw.GuildUpdate(b.guildUpdate)
	gw.GuildDelete(b.guildDelete)
	gw.VoiceStateUpdate(b.voiceStateUpdate)
	gw.ChannelCreate(b.channelCreate)
	gw.ChannelDelete(b.channelDelete)
	gw.ChannelUpdate(b.channelUpdate)

	// TODO: persist this into the db, and/or monitoring?
	gwBot, err := gw.GetBot()
	if err != nil {
		b.errs <- err
		return
	}

	// https://discord.com/developers/docs/topics/gateway#get-gateway-bot
	logger.WithFields(log.Fields{
		"identify_total":           gwBot.SessionStartLimit.Total,
		"identify_remaining":       gwBot.SessionStartLimit.Remaining,
		"identify_reset_after_sec": gwBot.SessionStartLimit.ResetAfter / 1000,
		"recommended_shards":       gwBot.Shards,
	}).Info("session start limit information")

	percentRemaining := float64(gwBot.SessionStartLimit.Remaining) / float64(gwBot.SessionStartLimit.Total)

	if percentRemaining < 0.5 {
		b.errs <- errors.New("have less than 50% of available sessions remaining, is something broken?")
		return
	}

	// Add connect jitter delay.
	var duration time.Duration
	if percentRemaining <= 0.9 {
		// Simple exponential backoff using the session start information.
		// ((1-0.90)*100)^2 where 0.9 is 90%.
		// 90% remaining: 100s
		// 80% remaining: 400s
		// 50% remaining: 2500s
		// ...etc.
		duration = time.Duration(math.Pow(2, (1.0-percentRemaining)*100.0)) * time.Second
	} else {
		duration = time.Duration(rand.Intn(5)) * 5 * time.Second
	}

	logger.WithField("duration", duration).Info("delaying connection to gateway (jitter + session start limits)")

	select {
	case <-time.After(duration):
	case <-b.ctx.Done():
		return
	}

	if err = gw.Connect(); err != nil {
		logger.WithError(err).Error("error connecting to gateway")
		b.errs <- err
		return
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		<-b.ctx.Done()

		err = gw.Disconnect()
		if err != nil {
			logger.WithError(err).Error("error disconnecting from gateway")
		}
	}()
}

// botReady is called when the bot successfully connects to the websocket.
func (b *discordBot) botReady() {}

// This event can be sent in three different scenarios:
//   - When a user is initially connecting, to lazily load and backfill
//     information for all unavailable guilds sent in the Ready event. Guilds
//     that are unavailable due to an outage will send a Guild Delete event.
//   - When a Guild becomes available again to the client.
//   - When the current user joins a new Guild.
//   - The inner payload is a guild object, with all the extra fields specified.
func (b *discordBot) guildCreate(s disgord.Session, h *disgord.GuildCreate) {
	b.updateMu.Lock()
	if _, ok := b.updates[update.guildID]; !ok {
		b.updates[update.guildID] = make(chan *updateEvent)
		// TODO: kick off goroutine.
	}
	b.updateMu.Unlock()

	var permissions uint64
	var err error
	var botMember *disgord.Member

	// The bot should be in the list of members returned during the guild create
	// message (even if no other users are listed due to not having the
	// permissions). Find this so we can understand what permissions we have.
	if bot, err := b.client.CurrentUser().Get(); err != nil {
		logger.WithError(err).Error("unable to fetch bot information during guild create event")
	} else {
		for _, member := range h.Guild.Members {
			if member.UserID.String() == bot.ID.String() {
				botMember = member
				break
			}
		}
	}

	// If we found ourselves, iterate through our roles, and combine the
	// permissions to understand what permissions we have.
	if botMember != nil {
		for _, role := range h.Guild.Roles {
			if !role.Managed {
				continue
			}

			var matches bool
			for _, id := range botMember.Roles {
				if role.ID.String() == id.String() {
					matches = true
					break
				}
			}

			if matches {
				permissions |= uint64(role.Permissions)
			}
		}
	}

	server, err := svcServers.GetByDiscordID(b.ctx, h.Guild.ID.String())
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
	}
}

func (b *discordBot) voiceStateUpdate(s disgord.Session, h *disgord.VoiceStateUpdate) {
	processUpdate(s)
}

func (b *discordBot) channelCreate(s disgord.Session, h *disgord.ChannelCreate) {

}

func (b *discordBot) channelDelete(s disgord.Session, h *disgord.ChannelDelete) {

}

func (b *discordBot) channelUpdate(s disgord.Session, h *disgord.ChannelUpdate) {

}

type updateEvent struct {
	sess    disgord.Session
	guildID disgord.Snowflake
	event   interface{}

	guild *disgord.Guild
}

func (b *discordBot) processUpdate(update *updateEvent) {
	var err error
	if update.guild == nil {
		update.guild, err = update.sess.Guild(update.guildID).Get()
		if err != nil {
			logGuild(logger, update.guildID).WithError(err).Error("unable to fetch guild for update event")
		}
	}

	b.updateMu.Lock()
	ch, ok := b.updates[update.guildID]

	if !ok {
		logGuild(logger, update.guild).Warn("dropping update vent due to missing channel")
		return
	}

	select {
	case ch <- update:
		b.updateMu.Unlock()
	default:
		b.updateMu.Unlock()
		logGuild(logger, update.guild).Warn("dropping update event due to full channel")
	}
}

func (b *discordBot) processUpdateWorker(ctx context.Context, events <-chan *updateEvent) {
	// TODO: make sure there is logic when channel gets closed
}
