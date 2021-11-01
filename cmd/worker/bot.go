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
