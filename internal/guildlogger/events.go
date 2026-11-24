// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

// Package guildlogger provides a logger which, in addition to logging through
// [log/slog], persists guild-specific events to the database, so they can be
// shown to guild admins.
package guildlogger

import (
	"context"
	"fmt"
	"log/slog"
	"maps"
	"sync"
	"time"

	"github.com/lrstanley/spectrograph/internal/database"
	"github.com/lrstanley/spectrograph/internal/database/ent"
	"github.com/lrstanley/spectrograph/internal/database/ent/guild"
	"github.com/lrstanley/spectrograph/internal/database/ent/guildevent"
	"github.com/lrstanley/spectrograph/internal/database/ent/privacy"
	"github.com/lrstanley/x/sync/cache"
	"github.com/lrstanley/x/sync/cache/policy/lru"
)

const (
	queueSize      = 500
	flushTimeout   = 10 * time.Second
	fillTimeout    = 5 * time.Second
	guildCacheSize = 1000
	guildCacheTTL  = 10 * time.Minute
)

type event struct {
	guildID  int
	typ      guildevent.Type
	message  string
	metadata map[string]any
}

// EventStream is a chainable logger, which logs through [log/slog], and also
// persists events associated with a guild to the database (flushed on an
// interval). The zero value is not usable, see [NewEventStream].
type EventStream struct {
	logger *slog.Logger
	db     *ent.Client
	gc     *cache.Cache[string, *ent.Guild] // Discord guild ID -> guild.
	fill   context.Context                  //nolint:containedctx // Used to load guilds into gc.
	queue  chan *event
	full   chan struct{} // signaled when the queue is half full, to flush early.

	stop     chan struct{} // closed by Close.
	stopOnce *sync.Once
	done     chan struct{} // closed when the flush loop exits.

	fields  map[string]any
	guildID string
	err     error
}

// NewEventStream returns a new EventStream for streaming guild logs/events to the
// database with appropriate annotations. The events will be flushed at the interval
// defined by flushInterval, and once more when ctx is cancelled, or when
// [EventStream.Close] is called.
func NewEventStream(ctx context.Context, logger *slog.Logger, db *ent.Client, flushInterval time.Duration) *EventStream {
	ctx = privacy.DecisionContext(ctx, privacy.Allow)

	es := &EventStream{
		logger: logger,
		db:     db,
		gc: cache.New(
			ctx,
			cache.WithLRU[string, *ent.Guild](lru.WithCapacity(guildCacheSize)),
			cache.WithDefaultEntryOptions[string, *ent.Guild](cache.WithExpiration(guildCacheTTL)),
		),
		// Not cancelled with ctx, so events logged while shutting down (before
		// Close) are still persisted.
		fill:     context.WithoutCancel(ctx),
		queue:    make(chan *event, queueSize),
		full:     make(chan struct{}, 1),
		stop:     make(chan struct{}),
		stopOnce: &sync.Once{},
		done:     make(chan struct{}),
	}

	go es.worker(ctx, flushInterval)

	return es
}

func (es *EventStream) worker(ctx context.Context, flushInterval time.Duration) {
	defer close(es.done)

	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			es.flush(ctx)
			return
		case <-es.stop:
			es.flush(ctx)
			return
		case <-ticker.C:
			es.flush(ctx)
		case <-es.full:
			es.flush(ctx)
		}
	}
}

// Close stops flushing events on an interval, and flushes all queued events,
// including those queued after the context passed to [NewEventStream] was
// cancelled. It blocks until the events are flushed, so it should be called
// before the database is closed, once nothing else logs events (events logged
// afterwards are only logged through slog).
func (es *EventStream) Close() {
	es.stopOnce.Do(func() { close(es.stop) })
	<-es.done
	es.flush(privacy.DecisionContext(context.Background(), privacy.Allow))
}

// flush writes all queued events to the database. Events are drained from the
// queue first, so the write isn't cancelled with ctx (e.g. when shutting down
// while flushing), only bounded by a timeout.
func (es *EventStream) flush(ctx context.Context) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), flushTimeout)
	defer cancel()

	var events []*event

drain:
	for {
		select {
		case evt := <-es.queue:
			events = append(events, evt)
		default:
			break drain
		}
	}

	if len(events) == 0 {
		return
	}

	err := database.RunWithTx(ctx, es.db, func(ctx context.Context, tx *ent.Tx) error {
		bulk := make([]*ent.GuildEventCreate, 0, len(events))
		for _, e := range events {
			bulk = append(
				bulk,
				tx.GuildEvent.Create().
					SetGuildID(e.guildID).
					SetType(e.typ).
					SetMessage(e.message).
					SetMetadata(e.metadata),
			)
		}
		return tx.GuildEvent.CreateBulk(bulk...).Exec(ctx)
	})
	if err != nil {
		es.logger.ErrorContext(ctx, "failed to flush guild events to db, dropping events", "count", len(events), "error", err)
	}
}

// clone returns a copy of the EventStream, so additional fields can be annotated
// onto the event, without modifying the parent.
func (es *EventStream) clone() *EventStream {
	evt := *es
	evt.fields = maps.Clone(es.fields)
	return &evt
}

// GuildID annotates a new event with the guild (Discord snowflake) ID. Events
// without a guild ID are only logged, not persisted.
func (es *EventStream) GuildID(gid string) *EventStream {
	evt := es.clone()
	evt.guildID = gid
	return evt
}

// WithError annotates a new event with the provided error.
func (es *EventStream) WithError(err error) *EventStream {
	evt := es.clone()
	evt.err = err
	return evt
}

// WithField annotates a new event with the provided field.
func (es *EventStream) WithField(key string, value any) *EventStream {
	evt := es.clone()
	if evt.fields == nil {
		evt.fields = make(map[string]any, 1)
	}
	evt.fields[key] = value
	return evt
}

// WithFields annotates a new event with the provided fields.
func (es *EventStream) WithFields(fields map[string]any) *EventStream {
	evt := es.clone()
	if evt.fields == nil {
		evt.fields = make(map[string]any, len(fields))
	}
	maps.Copy(evt.fields, fields)
	return evt
}

// sendEvent queues the event to be persisted to the database.
func (es *EventStream) sendEvent(typ guildevent.Type, msg string) {
	if es.guildID == "" {
		return
	}

	if es.err != nil {
		msg = fmt.Sprintf("%s: %s", msg, es.err)
	}

	g := es.guild()
	if g == nil {
		return
	}

	select {
	case es.queue <- &event{guildID: g.ID, typ: typ, message: msg, metadata: maps.Clone(es.fields)}:
		// Flush bursts (e.g. every guild syncing after a reconnect) early, rather
		// than dropping events until the next interval.
		if len(es.queue) >= queueSize/2 {
			select {
			case es.full <- struct{}{}:
			default:
			}
		}
	default:
		es.logger.Warn("guild event queue full, dropping event", "guild_id", es.guildID)
	}
}

// guild returns the (cached) guild for the event's guild ID, or nil if it can't
// be loaded.
func (es *EventStream) guild() *ent.Guild {
	g, err := es.gc.GetOrLoad(es.fill, es.guildID, func(ctx context.Context, gid string) (*ent.Guild, error) {
		ctx, cancel := context.WithTimeout(ctx, fillTimeout)
		defer cancel()
		return es.db.Guild.Query().Where(guild.GuildID(gid)).Only(ctx)
	})
	if err != nil {
		es.logger.Error("failed to load guild for event", "guild_id", es.guildID, "error", err)
		return nil
	}
	return g
}

// log logs the message through slog, with the EventStream's annotations.
func (es *EventStream) log(level slog.Level, msg string) {
	attrs := make([]slog.Attr, 0, len(es.fields)+3)
	for k, v := range es.fields {
		attrs = append(attrs, slog.Any(k, v))
	}

	if es.guildID != "" {
		attrs = append(attrs, slog.String("guild_id", es.guildID))
		if g := es.guild(); g != nil {
			attrs = append(attrs, slog.String("guild_name", g.Name))
		}
	}

	if es.err != nil {
		attrs = append(attrs, slog.Any("error", es.err))
	}

	es.logger.LogAttrs(context.Background(), level, msg, attrs...) //nolint:sloglint // dynamic by design.
}

func (es *EventStream) Info(msg string) {
	es.sendEvent(guildevent.TypeINFO, msg)
	es.log(slog.LevelInfo, msg)
}

func (es *EventStream) Infof(format string, v ...any) {
	es.Info(fmt.Sprintf(format, v...))
}

func (es *EventStream) Warn(msg string) {
	es.sendEvent(guildevent.TypeWARNING, msg)
	es.log(slog.LevelWarn, msg)
}

func (es *EventStream) Warnf(format string, v ...any) {
	es.Warn(fmt.Sprintf(format, v...))
}

func (es *EventStream) Error(msg string) {
	es.sendEvent(guildevent.TypeERROR, msg)
	es.log(slog.LevelError, msg)
}

func (es *EventStream) Errorf(format string, v ...any) {
	es.Error(fmt.Sprintf(format, v...))
}

func (es *EventStream) Debug(msg string) {
	es.sendEvent(guildevent.TypeDEBUG, msg)
	es.log(slog.LevelDebug, msg)
}

func (es *EventStream) Debugf(format string, v ...any) {
	es.Debug(fmt.Sprintf(format, v...))
}
