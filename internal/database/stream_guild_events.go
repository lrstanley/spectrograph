// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package database

import (
	"context"
	"time"

	"github.com/apex/log"
	"github.com/lrstanley/spectrograph/internal/ent"
	"github.com/lrstanley/spectrograph/internal/ent/guildevent"
	"github.com/lrstanley/spectrograph/internal/ent/predicate"
)

type guildEventStream struct {
	db     *ent.Client
	logger log.Interface
	inputs predicate.GuildEvent
	ch     chan *ent.GuildEvent

	interval   time.Duration
	maxHistory time.Time
	lastEvent  int
}

func NewGuildEventStream(
	ctx context.Context,
	inputs predicate.GuildEvent,
	interval,
	oldest time.Duration,
) <-chan *ent.GuildEvent {
	es := &guildEventStream{
		db:     ent.FromContext(ctx),
		logger: log.FromContext(ctx).WithField("src", "guild-event-stream"),
		inputs: inputs,
		ch:     make(chan *ent.GuildEvent),

		interval:   interval,
		maxHistory: time.Now().Add(-oldest),
		lastEvent:  -1,
	}

	if es.db == nil {
		panic("no ent client found in context")
	}

	go es.watcher(ctx)

	return es.ch
}

func (es *guildEventStream) watcher(ctx context.Context) {
	defer close(es.ch)

	defer es.logger.Info("closing guild event subscription")
	es.logger.Info("beginning guild event subscription")

	var err error

	// Initial population before the timer starts.
	if err = es.query(ctx); err != nil {
		es.logger.WithError(err).Error("failed to query guild events")
		return
	}

	timer := time.NewTicker(es.interval)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			if err = es.query(ctx); err != nil {
				es.logger.WithError(err).Error("failed to query guild events")
				return
			}
		}
	}
}

func (es *guildEventStream) query(ctx context.Context) error {
	events, err := es.db.GuildEvent.Query().Where(
		es.inputs,
		guildevent.IDGT(es.lastEvent),
		guildevent.CreateTimeGT(es.maxHistory),
	).All(ctx)
	if err != nil {
		return err
	}

	for _, event := range events {
		if event.ID > es.lastEvent {
			es.lastEvent = event.ID
		}

		es.ch <- event
	}

	return nil
}
