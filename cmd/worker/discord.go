// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"regexp"
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

type updateEvent struct {
	sess  disgord.Session
	event interface{}
	guild *disgord.Guild
}

func (b *discordBot) routeEvent(sess disgord.Session, event interface{}, guildID disgord.Snowflake) {
	update := &updateEvent{sess: sess, event: event}

	var err error
	update.guild, err = update.sess.Guild(guildID).Get()
	if err != nil {
		logGuild(logger, guildID).WithError(err).Error("unable to fetch guild for update event")
		return
	}

	b.updateMu.Lock()
	ch, ok := b.updates[guildID]
	if !ok {
		logGuild(logger, update.guild).Debug("dropping update event due to untracked guild")
		return
	}

	select {
	case ch <- update:
		b.updateMu.Unlock()
	default:
		b.updateMu.Unlock()
		// TODO: better way to handle us triggering subsequent events, as we
		// go through and make changes? i.e. us moving channels around will
		// cause events. Could we check the author?
		//
		// Alternatively, could we keep the last message, but do a wait and retry?
		logGuild(logger, update.guild).Debug("dropping event, as already processing event")
	}
}

func (b *discordBot) eventWatcher(sess disgord.Session, guildID disgord.Snowflake, events <-chan *updateEvent) {
	logGuild(logger, guildID).Debug("starting worker")
	defer logGuild(logger, guildID).Debug("closing worker")

	var event *updateEvent
	var ok bool

	for {
		select {
		case <-b.ctx.Done():
			return
		case event, ok = <-events:
			if !ok {
				// Assume channel was closed because we were disconnected from
				// the guild.
				return
			}

			b.processUpdateWorker(sess, guildID, event)
		}
	}
}

func (b *discordBot) processUpdateWorker(sess disgord.Session, guildID disgord.Snowflake, event *updateEvent) {
	logGuild(logger, event.guild).Info("processing event")
	fmt.Printf("%#v\n", event)

	// Get server options.
	serverOptionsAdmin, err := svcServers.GetOptionsAdmin(b.ctx, event.guild.ID.String())
	if err != nil {
		logGuild(logger, event.guild).WithError(err).Error("unable to fetch server options (admin)")
	}
	serverOptions, err := svcServers.GetOptions(b.ctx, event.guild.ID.String())
	if err != nil {
		logGuild(logger, event.guild).WithError(err).Error("unable to fetch server options")
	}

	if !serverOptionsAdmin.Enabled || !serverOptions.Enabled {
		return
	}

	// TODO: make a cache for this (with key just being the regex string).
	// Probably TTL it or key off the guild ID so it can't be abused.
	rgx, err := regexp.Compile(serverOptions.RegexMatch)
	if err != nil {
		logGuild(logger, event.guild).WithError(err).WithField("regex", serverOptions.RegexMatch).Error("unable to parse regex") // TODO
		return
	}

	// Get number of users in each voice channel.
	voiceCount := map[disgord.Snowflake]int{}
	for _, state := range event.guild.VoiceStates {
		voiceCount[state.ChannelID]++
	}

	// var hasEmptyChannel bool
	// var emptyChannelID string
	// var lastOccupiedChannel *disgord.Channel

	// state defines the state of parent channels and their "buckets" of managed
	// channels. See also:
	//   {
	//   	"<parent-id>": { // Could be nil or similar?
	//   		"Channel group 1": []<channel>, // The key being the results of the regex.
	//   		"Channel group 2": []<channel>,
	//   	}
	//   }
	state := map[disgord.Snowflake]map[string][]*disgord.Channel{}

	// TODO: count how many channel groups we have, as well as how many channels
	// in each group. Check this against our configured limits.

	for _, channel := range event.guild.Channels {
		if channel.Type != disgord.ChannelTypeGuildVoice {
			continue
		}

		if ok := rgx.MatchString(channel.Name); !ok {
			continue
		}

		if _, ok := state[channel.ParentID]; !ok {
			state[channel.ParentID] = map[string][]*disgord.Channel{}
		}

		if _, ok := state[channel.ParentID][channel.Name]; !ok {
			state[channel.ParentID][channel.Name] = []*disgord.Channel{channel}
		} else {
			state[channel.ParentID][channel.Name] = append(state[channel.ParentID][channel.Name], channel)
			// TODO: also check if the channel properties match that of the first entry.
			// If not, change them (assuming it's something supported by the server options).
		}
	}
}
