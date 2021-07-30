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

func discordSetup(ctx context.Context, wg *sync.WaitGroup, errs chan<- error) {
	client := disgord.New(disgord.Config{
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
		RejectEvents: []string{
			// TODO: disgord.AllEventsExcept(..)
			// See: https://github.com/andersfylling/disgord/issues/360#issuecomment-830918707
			disgord.EvtTypingStart,
			disgord.EvtPresenceUpdate,
			disgord.EvtGuildMemberAdd,
			disgord.EvtGuildMemberUpdate,
			disgord.EvtGuildMemberRemove,
			disgord.EvtGuildBanAdd,
			disgord.EvtGuildBanRemove,
			disgord.EvtGuildEmojisUpdate,
			disgord.EvtGuildIntegrationsUpdate,
			disgord.EvtWebhooksUpdate,
			disgord.EvtInviteCreate,
			disgord.EvtInviteDelete,
			disgord.EvtMessageCreate,
			disgord.EvtMessageDelete,
			disgord.EvtMessageDeleteBulk,
			disgord.EvtMessageReactionAdd,
			disgord.EvtMessageReactionRemove,
			disgord.EvtMessageReactionRemoveAll,
			disgord.EvtMessageReactionRemoveEmoji,
			disgord.EvtMessageUpdate,
		},
		ShardConfig: disgord.ShardConfig{
			ShardIDs:   []uint{},
			ShardCount: 0,
		},
	})
	gw := client.Gateway().WithContext(ctx)
	// TODO: register things here.
	// TODO: custom logger implementation that contains prefix info?

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

	// TODO: GuildDelete event for when we disconnect from a guild of some kind..?
	gw.GuildCreate(func(s disgord.Session, h *disgord.GuildCreate) {
		// time.Sleep(10 * time.Second)
		// pretty.Println(h)
		pretty.Println(s.GetPermissions())

		// TODO: why???
		guilds, err := s.CurrentUser().GetGuilds(&disgord.GetCurrentUserGuildsParams{})
		if err != nil {
			panic(err)
		}
		pretty.Println(guilds)

		// guild, err := client.Guild(h.Guild.ID).Get()
		// if err != nil {
		// 	panic(err)
		// }
		// pretty.Println(guild)
	})

	// gw.BotReady()

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
