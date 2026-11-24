// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package database

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/lrstanley/spectrograph/internal/database/ent"
	"github.com/lrstanley/spectrograph/internal/database/ent/guild"
	"github.com/lrstanley/spectrograph/internal/database/ent/guildevent"
)

// GuildEventStreamConfig configures [StreamGuildEvents].
type GuildEventStreamConfig struct {
	// GuildID is the (ent) ID of the guild to stream events for.
	GuildID int

	// AfterID, if non-zero, only streams events with an ID greater than this (e.g.
	// when resuming from an SSE Last-Event-ID). Otherwise, the most recent
	// HistoryLimit events (within HistoryAge) are sent first.
	AfterID int

	// HistoryLimit is the max number of historical events to send when the stream
	// starts. Defaults to 100.
	HistoryLimit int

	// HistoryAge is the max age of historical events to send when the stream starts.
	// Defaults to 2 hours.
	HistoryAge time.Duration

	// Interval is how often to poll for new events. Defaults to 3 seconds.
	Interval time.Duration

	// Notify, if set, triggers a poll before the next interval (see
	// [GuildEventNotifier.Subscribe]).
	Notify <-chan struct{}

	// Authorize, if set, is invoked before every query (including the initial
	// one), and returns the context to query with. Streams are long-lived, so
	// this allows re-checking access, and refreshing the user that privacy rules
	// are evaluated against (e.g. if the user is banned, or loses access to the
	// guild, while the stream is open). If it returns an error, the stream ends
	// with that error.
	Authorize func(ctx context.Context) (context.Context, error)
}

const maxEventsPerPoll = 500

// StreamGuildEvents polls the database for new events for a guild, invoking
// emit for each event (in ascending order), until ctx is cancelled or emit
// returns an error. Privacy rules are evaluated using ctx (or the context
// returned by [GuildEventStreamConfig.Authorize]), so DEBUG events will only be
// returned to system admins. Events include their guild edge.
func StreamGuildEvents(
	ctx context.Context,
	db *ent.Client,
	config GuildEventStreamConfig,
	emit func(*ent.GuildEvent) error,
) error {
	if config.HistoryLimit <= 0 {
		config.HistoryLimit = 100
	}
	if config.HistoryAge <= 0 {
		config.HistoryAge = 2 * time.Hour
	}
	if config.Interval <= 0 {
		config.Interval = 3 * time.Second
	}

	lastID := config.AfterID

	poll := func() ([]*ent.GuildEvent, error) {
		qctx := ctx
		if config.Authorize != nil {
			var err error
			if qctx, err = config.Authorize(ctx); err != nil {
				return nil, err
			}
		}

		// The guild is included for parity with the REST endpoints, which
		// eager-load it.
		query := db.GuildEvent.Query().
			Where(
				guildevent.HasGuildWith(guild.ID(config.GuildID)),
				guildevent.CreateTimeGT(time.Now().Add(-MaxGuildEventAge)),
			).
			WithGuild()

		if lastID > 0 {
			events, err := query.
				Where(guildevent.IDGT(lastID)).
				Order(ent.Asc(guildevent.FieldID)).
				Limit(maxEventsPerPoll).
				All(qctx)
			if err != nil {
				return nil, fmt.Errorf("querying guild events: %w", err)
			}
			return events, nil
		}

		events, err := query.
			Where(guildevent.CreateTimeGT(time.Now().Add(-config.HistoryAge))).
			Order(ent.Desc(guildevent.FieldID)).
			Limit(config.HistoryLimit).
			All(qctx)
		if err != nil {
			return nil, fmt.Errorf("querying guild events: %w", err)
		}
		slices.Reverse(events)
		return events, nil
	}

	ticker := time.NewTicker(config.Interval)
	defer ticker.Stop()

	for {
		events, err := poll()
		if err != nil {
			return err
		}

		for _, event := range events {
			lastID = max(lastID, event.ID)
			if eerr := emit(event); eerr != nil {
				return eerr
			}
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		case <-config.Notify:
		}
	}
}

// GuildEventNotifier notifies guild event streams (see
// [GuildEventStreamConfig.Notify]) when new events for their guild are
// committed, so they don't have to wait for the next poll. It only sees events
// written through ent clients that use [GuildEventNotifier.Hook] in this
// process. Events written by other processes (e.g. other replicas, or workers
// running separately) are still only picked up by polling. The zero value is
// ready to use.
type GuildEventNotifier struct {
	mu   sync.Mutex
	subs map[int]map[chan struct{}]struct{} // Guild ID -> subscribers.
}

// Hook returns an ent hook which notifies subscribers about created guild
// events. When created inside a transaction attached to the context (see
// [RunWithTx]), subscribers are notified once it commits.
func (n *GuildEventNotifier) Hook() ent.Hook {
	return ent.ObserveGuildEventID(func(_ context.Context, ev ent.ObserverEvent) error {
		m, ok := ev.Mutation.(*ent.GuildEventMutation)
		if !ok || ev.Op != ent.ObserverOpCreate {
			return nil
		}
		if gid, exists := m.GuildID(); exists {
			n.notify(gid)
		}
		return nil
	})
}

func (n *GuildEventNotifier) notify(guildID int) {
	n.mu.Lock()
	defer n.mu.Unlock()
	for ch := range n.subs[guildID] {
		select {
		case ch <- struct{}{}:
		default: // A notification is already pending.
		}
	}
}

// Subscribe returns a channel which receives a value when new events for the
// guild are created. Notifications are coalesced (the channel has a buffer of
// one). Call the returned function to unsubscribe. Safe to call on a nil
// notifier, which never notifies.
func (n *GuildEventNotifier) Subscribe(guildID int) (notify <-chan struct{}, unsubscribe func()) {
	if n == nil {
		return nil, func() {}
	}

	ch := make(chan struct{}, 1)

	n.mu.Lock()
	if n.subs == nil {
		n.subs = make(map[int]map[chan struct{}]struct{})
	}
	if n.subs[guildID] == nil {
		n.subs[guildID] = make(map[chan struct{}]struct{})
	}
	n.subs[guildID][ch] = struct{}{}
	n.mu.Unlock()

	return ch, func() {
		n.mu.Lock()
		defer n.mu.Unlock()
		delete(n.subs[guildID], ch)
		if len(n.subs[guildID]) == 0 {
			delete(n.subs, guildID)
		}
	}
}
