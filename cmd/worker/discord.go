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
	"github.com/kr/pretty"
)

var client *disgord.Client

func discordSetup(ctx context.Context, wg *sync.WaitGroup, errs chan<- error) {
	client = disgord.New(disgord.Config{
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
	gw := client.Gateway().WithContext(ctx)
	// TODO: custom logger implementation that contains prefix info?

	// Register hooks here.
	gw.BotReady(botReady)
	gw.GuildCreate(guildCreate)
	gw.GuildUpdate(guildUpdate)
	gw.GuildDelete(guildDelete)

	// TODO: persist this into the db, and/or monitoring?
	gwBot, err := gw.GetBot()
	if err != nil {
		errs <- err
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
		errs <- errors.New("have less than 50% of available sessions remaining, is something broken?")
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
	case <-ctx.Done():
		return
	}

	if err = gw.Connect(); err != nil {
		logger.WithError(err).Error("error connecting to gateway")
		errs <- err
		return
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		<-ctx.Done()

		err = gw.Disconnect()
		if err != nil {
			logger.WithError(err).Error("error disconnecting from gateway")
		}
	}()
}

// botReady is called when the bot successfully connects to the websocket.
func botReady() {
	guilds, err := client.CurrentUser().GetGuilds(&disgord.GetCurrentUserGuildsParams{})
	if err != nil {
		panic(err)
	}
	pretty.Println(guilds)
}

// This event can be sent in three different scenarios:
//   - When a user is initially connecting, to lazily load and backfill
//     information for all unavailable guilds sent in the Ready event. Guilds
//     that are unavailable due to an outage will send a Guild Delete event.
//   - When a Guild becomes available again to the client.
//   - When the current user joins a new Guild.
//   - The inner payload is a guild object, with all the extra fields specified.
func guildCreate(s disgord.Session, h *disgord.GuildCreate) {
	// pretty.Println(h.Guild)
	// guild, err := client.Guild(h.Guild.ID).Get()
	// if err != nil {
	// 	panic(err)
	// }
	// pretty.Println(guild)
}

// Sent when a guild is updated. The inner payload is a guild object.
func guildUpdate(s disgord.Session, h *disgord.GuildUpdate) {
	// pretty.Println(h)
}

// Sent when a guild becomes or was already unavailable due to an outage, or
// when the user leaves or is removed from a guild. The inner payload is an
// unavailable guild object. If the unavailable field is not set, the user
// was removed from the guild.
func guildDelete(s disgord.Session, h *disgord.GuildDelete) {
	if h.UserWasRemoved() {
		// TODO: clean up from db, maybe have it send a notification?
	}
}
