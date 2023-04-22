// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package guildlogger

import (
	"context"
	"fmt"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/apex/log"
	"github.com/lrstanley/spectrograph/internal/cache"
	"github.com/lrstanley/spectrograph/internal/database"
	"github.com/lrstanley/spectrograph/internal/database/ent"
	"github.com/lrstanley/spectrograph/internal/database/ent/guild"
	"github.com/lrstanley/spectrograph/internal/database/ent/guildevent"
)

type event struct {
	guild *ent.Guild
	event *ent.GuildEvent
}

type EventStream struct {
	logger        log.Interface
	db            *ent.Client
	gc            *cache.Cache[string, *ent.Guild]
	workerCh      chan *event
	flushInterval time.Duration

	// Internally used for chaining.
	base *EventStream

	logSend    bool
	logFields  log.Fields
	logGuildID string
	logError   error
}

// NewEventStream returns a new EventStream for streaming guild logs/events to the
// database with appropriate annotations. The events will be flushed at the interval
// defined by flushInterval.
//
// When the context provided to the event stream is cancelled, the event stream
// will close its worker and any pending events will be dropped (too lazy to bother).
func NewEventStream(ctx context.Context, l log.Interface, db *ent.Client, flushInterval time.Duration) *EventStream {
	es := &EventStream{
		logger: l,
		db:     db,
		gc: cache.New[string, *ent.Guild](1000, 10*time.Minute, func(gid string) (*ent.Guild, bool) {
			g, err := db.Guild.Query().Where(guild.GuildID(gid)).Only(ctx)
			if err != nil {
				l.WithError(err).Error("failed to fill guild cache from db")
				return nil, false
			}

			return g, true
		}),
		workerCh:      make(chan *event, 100),
		flushInterval: flushInterval,
	}

	// Spin up worker to send events to the database.
	go es.worker(ctx)

	return es
}

func (es *EventStream) worker(ctx context.Context) {
	ticker := time.NewTicker(es.flushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			var events []*event

			for {
				select {
				case evt, ok := <-es.workerCh:
					if !ok {
						goto next
					}

					events = append(events, evt)
				default:
					goto next
				}
			}

		next:
			if len(events) == 0 {
				continue
			}

			tx, err := es.db.Tx(ctx)
			if err != nil {
				es.logger.WithError(err).Error("failed to flush logs to db, dropping events")
				continue
			}

			var bulk []*ent.GuildEventCreate

			for _, e := range events {
				bulk = append(
					bulk,
					tx.GuildEvent.Create().
						SetGuildID(e.guild.ID).
						SetType(e.event.Type).
						SetMessage(e.event.Message).
						SetMetadata(e.event.Metadata),
				)
			}

			_, err = tx.GuildEvent.CreateBulk(bulk...).Save(ctx)

			err = database.Commit(tx, err)
			if err != nil {
				es.logger.WithError(err).Error("failed to flush logs to db, dropping events")
				continue
			}
		}
	}
}

// clone returns a new EventStream with the same base as the current one, so
// additional fields can be annotated onto the event.
func (es *EventStream) clone() *EventStream {
	base := es.base

	if es.base == nil {
		base = es
	}

	return &EventStream{
		logger: es.logger,
		db:     es.db,
		gc:     es.gc,

		logFields:  es.logFields,
		logGuildID: es.logGuildID,
		logError:   es.logError,

		base: base,
	}
}

// Guild annotates a new event with the guild information.
func (es *EventStream) Guild(dg *disgord.Guild) *EventStream {
	evt := es.clone()
	evt.logGuildID = dg.ID.String()
	return evt
}

// GuildID annotates a new event with the guild information.
func (es *EventStream) GuildID(gid string) *EventStream {
	evt := es.clone()
	evt.logGuildID = gid
	return evt
}

// WithError annotates a new event with the provided error.
func (es *EventStream) WithError(err error) *EventStream {
	evt := es.clone()
	evt.logError = err
	return evt
}

// WithField annotates a new event with the provided field.
func (es *EventStream) WithField(key string, value any) *EventStream {
	evt := es.clone()

	if evt.logFields == nil {
		evt.logFields = log.Fields{key: value}
	} else {
		evt.logFields[key] = value
	}

	return evt
}

// WithFields annotates a new event with the provided fields.
func (es *EventStream) WithFields(fields log.Fields) *EventStream {
	evt := es.clone()

	if evt.logFields == nil {
		evt.logFields = fields
	} else {
		for k, v := range fields {
			evt.logFields[k] = v
		}
	}

	return evt
}

// sendEvent writes the event to the EventStream queue to be sent to the db.
func (es *EventStream) sendEvent(eventType guildevent.Type, msg string) {
	if es.logGuildID == "" {
		return
	}

	if es.logError != nil {
		msg = fmt.Sprintf("%s: %s", msg, es.logError)
	}

	g := es.gc.Get(es.logGuildID)
	if g == nil {
		es.logger.Error("failed to get guild from cache when sending event")
		return
	}

	evt := &event{
		guild: g,
		event: &ent.GuildEvent{
			Type:     eventType,
			Message:  msg,
			Metadata: es.logFields,
		},
	}

	if es.base != nil {
		es.base.workerCh <- evt
	} else {
		es.workerCh <- evt
	}
}

// buildLog builds a log entry with the current EventStream's fields.
func (es *EventStream) buildLog() log.Interface {
	var l log.Interface

	if es.logFields != nil {
		l = es.logger.WithFields(es.logFields)
	} else {
		l = es.logger
	}

	if es.logGuildID != "" {
		l = l.WithField("guild_id", es.logGuildID)
	}

	if es.logError != nil {
		l = l.WithError(es.logError)
	}

	if g := es.gc.Get(es.logGuildID); g != nil {
		l = l.WithField("guild_name", g.Name)
	}

	return l
}

func (es *EventStream) Info(msg string) {
	es.sendEvent(guildevent.TypeINFO, msg)
	es.buildLog().Info(msg)
}

func (es *EventStream) Infof(msg string, v ...any) {
	msg = fmt.Sprintf(msg, v...)
	es.sendEvent(guildevent.TypeINFO, msg)
	es.buildLog().Info(msg)
}

func (es *EventStream) Warn(msg string) {
	es.sendEvent(guildevent.TypeWARNING, msg)
	es.buildLog().Warn(msg)
}

func (es *EventStream) Warnf(msg string, v ...any) {
	msg = fmt.Sprintf(msg, v...)
	es.sendEvent(guildevent.TypeWARNING, msg)
	es.buildLog().Warn(msg)
}

func (es *EventStream) Error(msg string) {
	es.sendEvent(guildevent.TypeERROR, msg)
	es.buildLog().Error(msg)
}

func (es *EventStream) Errorf(msg string, v ...any) {
	msg = fmt.Sprintf(msg, v...)
	es.sendEvent(guildevent.TypeERROR, msg)
	es.buildLog().Error(msg)
}

func (es *EventStream) Debug(msg string) {
	es.sendEvent(guildevent.TypeDEBUG, msg)
	es.buildLog().Debug(msg)
}

func (es *EventStream) Debugf(msg string, v ...any) {
	msg = fmt.Sprintf(msg, v...)
	es.sendEvent(guildevent.TypeDEBUG, msg)
	es.buildLog().Debug(msg)
}
