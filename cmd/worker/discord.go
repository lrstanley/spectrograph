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

	voiceStates *voiceStateTracker
	permissions *permissionsCache

	updateMu sync.RWMutex
	updates  map[disgord.Snowflake]chan *updateEvent
}

func (b *discordBot) setup(wg *sync.WaitGroup) {
	b.updates = make(map[disgord.Snowflake]chan *updateEvent)
	b.voiceStates = NewVoiceStateTracker()
	b.permissions = NewPermissionsCache()

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
			disgord.EvtGuildRoleUpdate, // TODO: Also this.
			disgord.EvtGuildRoleDelete,
			disgord.EvtChannelCreate,
			disgord.EvtChannelUpdate,
			disgord.EvtChannelDelete,
			disgord.EvtVoiceStateUpdate,
		),

		// ShardConfig: disgord.ShardConfig{
		// 	ShardIDs:   []uint{},
		// 	ShardCount: 0,
		// },
	})
	gw := b.client.Gateway().WithContext(b.ctx)
	// TODO: custom logger implementation that contains prefix info?

	b.voiceStates.Register(gw)

	// Register hooks here.
	gw.BotReady(b.botReady)
	gw.GuildRoleUpdate(b.GuildRoleUpdate)
	gw.GuildRoleDelete(b.GuildRoleDelete)
	gw.GuildCreate(b.guildCreate)
	gw.GuildUpdate(b.guildUpdate)
	gw.GuildDelete(b.guildDelete)
	gw.VoiceStateUpdate(b.voiceStateUpdate)
	gw.GuildMemberUpdate(b.guildMemberUpdate)
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

	// Debounce by this timeframe. I.e. if we receive multiple events within
	// the below timeframe, wait until we receive no new events within the
	// timeframe, then send out the last event we received.
	timer := time.NewTimer(5000 * time.Millisecond)

	for {
		select {
		case <-b.ctx.Done():
			return
		case <-timer.C:
			if event != nil {
				b.processUpdateWorker(sess, guildID, event) // TODO: if this panics, catch it.
				event = nil
			}
		case event, ok = <-events:
			if !ok {
				// Assume channel was closed because we were disconnected from
				// the guild.
				return
			}
			timer.Reset(500 * time.Millisecond)
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
		return
	}
	serverOptions, err := svcServers.GetOptions(b.ctx, event.guild.ID.String())
	if err != nil {
		logGuild(logger, event.guild).WithError(err).Error("unable to fetch server options")
		return
	}

	if !serverOptionsAdmin.Enabled || !serverOptions.Enabled {
		logGuild(logger, event.guild).Debug("dropping event, guild disabled via options")
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
	voiceCount := b.voiceStates.UserCount(event.guild.ID)

	// state defines the state of parent channels and their "buckets" of managed
	// channels. See also:
	//   {
	//   	parent.Snowflake: {
	//   		channel1.Name: []<channel>,
	//   		channel2.Name: []<channel>,
	//   	}
	//   }
	state := map[disgord.Snowflake]map[string][]*disgord.Channel{}
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
		}
	}

	// TODO: count how many channel groups we have, as well as how many channels
	// in each group. Check this against our configured limits.

	// TODO: should we make sure we have permissions?
	// TODO: support empty channel being at the top, vs the bottom.
	// TODO: support multiple empty channels? could help in the event of bot
	// issues if there is a "buffer" of empty channels.

	var emptyChannel *disgord.Channel
	var lastOccupiedChannel *disgord.Channel

	toDelete := []*disgord.Channel{}

	for parent := range state {
		for group := range state[parent] {
			emptyChannel = nil
			lastOccupiedChannel = nil

			// Find which channels are empty, and which have users. If there
			// are more than one channels that are empty, mark all subsequent
			// empty channels for deletion.
			for _, channel := range state[parent][group] {
				if voiceCount[channel.ID] == 0 {
					if emptyChannel != nil {
						toDelete = append(toDelete, channel)
						continue
					}

					emptyChannel = channel
				} else {
					lastOccupiedChannel = channel
				}
			}

			// Move empty channel to position after the last occupied channel.
			if emptyChannel != nil && lastOccupiedChannel != nil && lastOccupiedChannel.ID != emptyChannel.ID {
				err = sess.Guild(event.guild.ID).UpdateChannelPositions([]disgord.UpdateGuildChannelPositionsParams{
					disgord.UpdateGuildChannelPositionsParams{ID: lastOccupiedChannel.ID, Position: lastOccupiedChannel.Position},
					disgord.UpdateGuildChannelPositionsParams{ID: emptyChannel.ID, Position: lastOccupiedChannel.Position + 1},
				})
				if err != nil {
					logGuild(logger, event.guild).WithError(err).WithFields(log.Fields{
						"last_occupied_id": lastOccupiedChannel.ID,
						"empty_channel_id": emptyChannel.ID,
					}).Error("unable to reorder empty channel to bottom")
				}
			}

			// If no empty channel, make one, duplicating the config from the first
			// channel in the bucket.
			if emptyChannel == nil {
				channel, err := sess.Guild(event.guild.ID).CreateChannel(state[parent][group][0].Name, &disgord.CreateGuildChannelParams{
					Name:                 state[parent][group][0].Name,
					Type:                 state[parent][group][0].Type,
					Bitrate:              state[parent][group][0].Bitrate,
					UserLimit:            state[parent][group][0].UserLimit,
					RateLimitPerUser:     state[parent][group][0].RateLimitPerUser,
					PermissionOverwrites: state[parent][group][0].PermissionOverwrites,
					ParentID:             state[parent][group][0].ParentID,
					NSFW:                 state[parent][group][0].NSFW,
					Position:             state[parent][group][0].Position + 1,
				})
				if err != nil {
					logGuild(logger, event.guild).WithError(err).WithField("source_channel_id", state[parent][group][0]).Error("unable to create new channel from master channel")
				} else {
					emptyChannel = channel

					// add new channel to state.
					state[parent][group] = append(state[parent][group], channel)
				}
			}
		}
	}

	for _, channel := range toDelete {
		if _, err := sess.Channel(channel.ID).Delete(); err != nil {
			// This can sometimes cause "unknown channel" if executed too fast
			// between other API calls.
			if restErr, ok := err.(*disgord.ErrRest); ok {
				if restErr.Code == 10003 {
					// https://discord.com/developers/docs/topics/opcodes-and-status-codes#json-json-error-codes
					continue
				}
			}

			logGuild(logger, event.guild).WithError(err).WithField("channel_id", channel.ID).Error("unable to remove unneeded empty channel")
		}

		// Remove the deleted channel from state.
		group := state[channel.ParentID][channel.Name]
		for i := 0; i < len(group); i++ {
			if group[i].Compare(channel) {
				copy(group[i:], group[i+1:])
				group[len(group)-1] = nil                                    // remove last element to prevent memory leaking.
				state[channel.ParentID][channel.Name] = group[:len(group)-1] // truncate slice.

				break
			}
		}
	}

	// Loop through all of the channels and make sure their config matches that of the
	// "primary" channel in the list. I.e. change the primary, and the rest should
	// change.
	for parent := range state {
		// parentChannel, err := sess.Channel(parent).Get()
		// if err != nil {
		// 	pretty.Println(err)
		// } else {
		// 	pretty.Println(parentChannel)
		// }
		for group := range state[parent] {
			// Check if it's just one channel.
			if len(state[parent][group]) < 2 {
				continue
			}

			primary := state[parent][group][0]
			var needsUpdate bool

			for _, channel := range state[parent][group] {
				if channel.UserLimit != primary.UserLimit ||
					channel.Bitrate != primary.Bitrate ||
					len(channel.PermissionOverwrites) != len(primary.PermissionOverwrites) {
					needsUpdate = true
				}

				if !needsUpdate {
					for i := 0; i < len(channel.PermissionOverwrites); i++ {
						if channel.PermissionOverwrites[i].Type != primary.PermissionOverwrites[i].Type ||
							channel.PermissionOverwrites[i].ID != primary.PermissionOverwrites[i].ID ||
							channel.PermissionOverwrites[i].Allow != primary.PermissionOverwrites[i].Allow ||
							channel.PermissionOverwrites[i].Deny != primary.PermissionOverwrites[i].Deny {
							needsUpdate = true
							break
						}
					}
				}

				if !needsUpdate {
					continue
				}

				_, err := sess.Channel(channel.ID).UpdateBuilder().
					SetPosition(channel.Position).
					SetUserLimit(primary.UserLimit).
					SetBitrate(primary.Bitrate).
					SetPermissionOverwrites(primary.PermissionOverwrites).Execute()
				if err != nil {
					// TODO: this should be propagated up to the user somehow. events?
					// TODO: should we change the permissions ourselves?
					logGuild(logger, event.guild).WithError(err).WithFields(log.Fields{
						"channel_id": channel.ID,
						"primary_id": primary.ID,
					}).Error("unable to update children channel details based off primary channel")
				}
			}
		}
	}
}

// TODO: function to re-order and/or add permissions specifically for the bot user, into
// the channel permission overrides, ONLY if there is a permission that disallows
// being able to read/update, etc??
func (b *discordBot) changeChannelPermissions(sess *disgord.Session, channel *disgord.Channel) error {
	return nil
}
