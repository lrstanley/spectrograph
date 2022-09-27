// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package worker

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"sync"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/snowflake/v5"
	"github.com/apex/log"
	"github.com/lrstanley/spectrograph/internal/ent"
	"github.com/lrstanley/spectrograph/internal/ent/privacy"
	"github.com/lrstanley/spectrograph/internal/guildlogger"
	"github.com/lrstanley/spectrograph/internal/metrics"
	"github.com/lrstanley/spectrograph/internal/models"
	"github.com/lrstanley/spectrograph/internal/permissions"
	"github.com/lrstanley/spectrograph/internal/voicestate"
)

type Worker struct {
	ctx    context.Context
	logger log.Interface
	es     *guildlogger.EventStream
	db     *ent.Client

	config models.ConfigWorker
	auth   models.ConfigDiscord

	client      *disgord.Client
	voiceStates *voicestate.Tracker
	permissions *permissions.Cache

	errs chan<- error

	updateMu sync.RWMutex
	updates  map[disgord.Snowflake]chan *updateEvent
}

func New(ctx context.Context, config models.ConfigWorker, auth models.ConfigDiscord) (*Worker, error) {
	ctx = privacy.DecisionContext(ctx, privacy.Allow)
	db := ent.FromContext(ctx)

	if config.NumShards < 1 {
		config.NumShards = 1
	}

	if len(config.ShardIDs) < 1 {
		config.ShardIDs = []uint{0}
	}

	logger := log.FromContext(ctx).WithFields(log.Fields{
		"shard_ids": config.ShardIDs,
		"src":       "worker",
	})

	es := guildlogger.NewEventStream(ctx, logger, db)

	w := &Worker{
		ctx:         ctx,
		logger:      logger,
		es:          es,
		db:          db,
		config:      config,
		auth:        auth,
		voiceStates: voicestate.NewTracker(),
		permissions: permissions.NewCache(ctx, es),
		errs:        make(chan<- error),
		updates:     make(map[disgord.Snowflake]chan *updateEvent),

		client: disgord.New(disgord.Config{
			BotToken:    config.BotToken,
			ProjectName: "spectrograph (https://github.com/lrstanley, https://liam.sh)",
			Presence: &disgord.UpdateStatusPayload{
				Since: nil,
				Game: []*disgord.Activity{
					{Name: "voice chan curator", Type: disgord.ActivityTypeGame},
				},
				Status: disgord.StatusOnline,
				AFK:    false,
			},
			Logger: guildlogger.New(logger),
			RejectEvents: disgord.AllEventsExcept(
				// See: https://github.com/andersfylling/disgord/issues/360#issuecomment-830918707
				disgord.EvtReady,
				disgord.EvtResumed,
				disgord.EvtGuildCreate,
				disgord.EvtGuildUpdate,
				disgord.EvtGuildDelete,
				disgord.EvtGuildRoleCreate,
				disgord.EvtGuildRoleUpdate, // TODO: handle this.
				disgord.EvtGuildRoleDelete,
				disgord.EvtChannelCreate,
				disgord.EvtChannelUpdate,
				disgord.EvtChannelDelete,
				disgord.EvtVoiceStateUpdate,
			),

			ShardConfig: disgord.ShardConfig{
				ShardIDs:   config.ShardIDs,
				ShardCount: config.NumShards,
			},
		}),
	}

	logger.WithField("config", config).Debug("registered worker")

	return w, nil
}

func (w *Worker) Run(ctx context.Context) error {
	ctx = privacy.DecisionContext(ctx, privacy.Allow)

	go w.metricMonitor(ctx)

	gw := w.client.Gateway().WithContext(ctx)

	w.voiceStates.Register(gw)

	// Register hooks here.
	gw.BotReady(w.botReady)

	mw := gw.WithMiddleware(w.metricOperation)
	mw.GuildRoleUpdate(w.GuildRoleUpdate)
	mw.GuildRoleDelete(w.GuildRoleDelete)
	mw.GuildCreate(w.guildCreate)
	mw.GuildUpdate(w.guildUpdate)
	mw.GuildDelete(w.guildDelete)
	mw.VoiceStateUpdate(w.voiceStateUpdate)
	mw.GuildMemberUpdate(w.guildMemberUpdate)
	mw.ChannelCreate(w.channelCreate)
	mw.ChannelDelete(w.channelDelete)
	mw.ChannelUpdate(w.channelUpdate)

	// TODO: persist this into the db, and/or monitoring?
	gwBot, err := gw.GetBot()
	if err != nil {
		return err
	}

	// https://discord.com/developers/docs/topics/gateway#get-gateway-bot
	w.logger.WithFields(log.Fields{
		"identify_total":           gwBot.SessionStartLimit.Total,
		"identify_remaining":       gwBot.SessionStartLimit.Remaining,
		"identify_reset_after_sec": gwBot.SessionStartLimit.ResetAfter / 1000,
		"recommended_shards":       gwBot.Shards,
	}).Info("session start limit information")

	percentRemaining := float64(gwBot.SessionStartLimit.Remaining) / float64(gwBot.SessionStartLimit.Total)

	if percentRemaining < 0.5 {
		return errors.New("have less than 50% of available sessions remaining, is something broken?")
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

	w.logger.WithField("duration", duration).Info("delaying connection to gateway (jitter + session start limits)")

	select {
	case <-time.After(duration):
	case <-ctx.Done():
		return nil
	}

	if err = gw.Connect(); err != nil {
		w.logger.WithError(err).Error("error connecting to gateway")
		return err
	}

	<-ctx.Done()

	err = gw.Disconnect()
	if err != nil {
		w.logger.WithError(err).Error("error disconnecting from gateway")
	}
	return err
}

func (w *Worker) metricMonitor(ctx context.Context) {
	logger := w.logger.WithField("src", "worker-metrics")
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			bot, err := w.client.Gateway().GetBot()
			if err != nil {
				logger.WithError(err).Error("error getting bot information")
				continue
			}

			for _, shard := range w.config.ShardIDs {
				metrics.WorkerSessionsRemainingCount.
					WithLabelValues(fmt.Sprintf("%d", shard)).
					Set(float64(bot.SessionStartLimit.Remaining))

				metrics.WorkerSessionsTotalCount.
					WithLabelValues(fmt.Sprintf("%d", shard)).
					Set(float64(bot.SessionStartLimit.Total))
			}
		case <-ctx.Done():
			return
		}
	}
}

func (w *Worker) metricOperation(v any) any {
	rv := reflect.ValueOf(v)

	if rv.IsNil() {
		return v
	}

	var gid string
	var shardID uint

	if f := reflect.Indirect(rv).FieldByName("GuildID"); f != reflect.ValueOf(nil) {
		if id, ok := f.Interface().(snowflake.Snowflake); ok {
			gid = id.String()
		}
	} else if f := reflect.Indirect(rv).FieldByName("Guild"); f != reflect.ValueOf(nil) {
		if g, ok := f.Interface().(*disgord.Guild); ok {
			gid = g.ID.String()
		}
	}

	if f := reflect.Indirect(rv).FieldByName("ShardID"); f != reflect.ValueOf(nil) {
		if id, ok := f.Interface().(uint); ok {
			shardID = id
		}
	}

	if gid != "" {
		metrics.WorkerEventCount.
			WithLabelValues(
				fmt.Sprintf("%d", shardID),
				reflect.TypeOf(v).String(),
				gid,
			).Inc()
	}

	return v
}
