// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package guildlogger

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/lrstanley/spectrograph/internal/database/dbtest"
	"github.com/lrstanley/spectrograph/internal/database/ent/guildevent"
	"github.com/lrstanley/spectrograph/internal/database/ent/privacy"
)

// TestCloseFlushesAfterCancel checks that events logged after the stream's
// context is cancelled (e.g. while the worker shuts down), including events
// for guilds which aren't cached yet, are persisted by Close.
func TestCloseFlushesAfterCancel(t *testing.T) {
	t.Parallel()

	db := dbtest.New(t)
	ctx := privacy.DecisionContext(t.Context(), privacy.Allow)

	for _, gid := range []string{"100", "200"} {
		if err := db.Guild.Create().SetGuildID(gid).SetName("guild " + gid).SetIconURL("").Exec(ctx); err != nil {
			t.Fatal(err)
		}
	}

	sctx, cancel := context.WithCancel(ctx)
	es := NewEventStream(sctx, slog.New(slog.DiscardHandler), db, time.Hour)

	es.GuildID("100").Info("before shutdown")
	cancel()
	// Let the flush triggered by the cancellation happen first.
	time.Sleep(100 * time.Millisecond)

	es.GuildID("100").WithField("k", "v").Warn("during shutdown")
	es.GuildID("200").Error("during shutdown, uncached guild")
	es.Close()
	es.Close() // Idempotent.

	events, err := db.GuildEvent.Query().Order(guildevent.ByID()).All(ctx)
	if err != nil {
		t.Fatal(err)
	}

	var msgs []string
	for _, e := range events {
		msgs = append(msgs, e.Message)
	}
	if len(msgs) != 3 {
		t.Fatalf("persisted events = %q, want 3", msgs)
	}
	if events[1].Metadata["k"] != "v" {
		t.Errorf("metadata = %v, want k=v", events[1].Metadata)
	}
}

// TestFlushesBursts checks that a burst of events larger than the queue (e.g.
// every guild syncing after a reconnect) is flushed early instead of dropped.
func TestFlushesBursts(t *testing.T) {
	t.Parallel()

	db := dbtest.New(t)
	ctx := privacy.DecisionContext(t.Context(), privacy.Allow)

	if err := db.Guild.Create().SetGuildID("100").SetName("guild 100").SetIconURL("").Exec(ctx); err != nil {
		t.Fatal(err)
	}

	es := NewEventStream(ctx, slog.New(slog.DiscardHandler), db, time.Hour)

	const burst = queueSize * 2
	for range burst {
		es.GuildID("100").Info("burst")
		// The producer outpaces a flush, but not a reasonable one.
		time.Sleep(time.Millisecond)
	}
	es.Close()

	n, err := db.GuildEvent.Query().Count(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if n != burst {
		t.Errorf("persisted events = %d, want %d", n, burst)
	}
}
