// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"context"
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/apex/log"
)

func discordSetup(ctx context.Context, wg *sync.WaitGroup, errs chan<- error) {
	// TODO: add delay interval for sharding connections?

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
	})

	wg.Add(1)
	go func() {
		defer wg.Done()

		gw := client.Gateway()
		// TODO: register things here.
		// TODO: custom logger implementation that contains prefix info?

		var attempts int

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			attempts++

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

			if attempts > 5 && (float64(gwBot.SessionStartLimit.Remaining)/float64(gwBot.SessionStartLimit.Total)) < 0.75 {
				errs <- errors.New("have less than 75% of available sessions remaining (and failed to connect more than 5 times), is something broken?")
				return
			}

			// Add connect jitter delay.
			var duration time.Duration
			if attempts > 1 {
				duration = time.Duration(attempts-1) * 45 * time.Second
			} else {
				duration = time.Duration(rand.Intn(5)) * 5 * time.Second
			}

			logger.WithField("duration", duration).Info("delaying connection to gateway (jitter)")

			select {
			case <-time.After(duration):
			case <-ctx.Done():
				return
			}

			return

			// if err = gw.WithContext(ctx).Connect(); err != nil {
			// 	logger.WithError(err).Error("error connecting to gateway")
			// 	continue
			// }

			// Assume we connected successfully.
			attempts = 0
		}
	}()
}
