// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package database_test

import (
	"context"
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/lrstanley/spectrograph/internal/database"
	"github.com/lrstanley/spectrograph/internal/database/dbtest"
	"github.com/lrstanley/spectrograph/internal/database/ent"
	"github.com/lrstanley/spectrograph/internal/database/ent/guildevent"
)

// collectStream runs [database.StreamGuildEvents] (as a system admin) for the
// provided duration, returning the IDs of the emitted events, and the error the
// stream ended with.
func collectStream(
	t *testing.T,
	db *ent.Client,
	f *dbtest.Fixtures,
	d time.Duration,
	config database.GuildEventStreamConfig,
) (ids []int, err error) {
	t.Helper()

	ctx, cancel := context.WithTimeout(dbtest.AsUser(context.Background(), f.Admin), d)
	defer cancel()

	config.Interval = 20 * time.Millisecond
	err = database.StreamGuildEvents(ctx, db, config, func(e *ent.GuildEvent) error {
		if e.Edges.Guild == nil || e.Edges.Guild.ID != config.GuildID {
			t.Errorf("event %d: expected guild %d to be loaded, got: %+v", e.ID, config.GuildID, e.Edges.Guild)
		}
		ids = append(ids, e.ID)
		return nil
	})
	if errors.Is(err, context.DeadlineExceeded) {
		err = nil
	}
	return ids, err
}

func TestStreamGuildEvents(t *testing.T) {
	t.Parallel()

	db := dbtest.New(t)
	f := dbtest.NewFixtures(t, db)
	allow := dbtest.AllowAll(context.Background())

	// Older than the history window, but still within the retention period.
	quiet := db.Guild.Create().SetGuildID("2003").SetName("quiet").SetIconURL("").SaveX(allow)
	old := db.GuildEvent.Create().
		SetGuild(quiet).
		SetType(guildevent.TypeINFO).
		SetMessage("old").
		SetCreateTime(time.Now().Add(-3 * time.Hour)).
		SaveX(allow)

	t.Run("history-excludes-old-events", func(t *testing.T) {
		t.Parallel()

		// Without recent events, polling must not fall back to streaming every
		// event within the retention period.
		ids, err := collectStream(t, db, f, 2*time.Second, database.GuildEventStreamConfig{
			GuildID:    quiet.ID,
			HistoryAge: time.Hour,
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(ids) != 0 {
			t.Errorf("event ids = %v, want none", ids)
		}
	})

	t.Run("resume-after-id", func(t *testing.T) {
		t.Parallel()

		ids, err := collectStream(t, db, f, 2*time.Second, database.GuildEventStreamConfig{
			GuildID: f.GuildA.ID,
			AfterID: f.GuildAInfoEvent.ID,
		})
		if err != nil {
			t.Fatal(err)
		}
		if want := []int{f.GuildADebugEvent.ID}; !slices.Equal(ids, want) {
			t.Errorf("event ids = %v, want %v", ids, want)
		}
	})

	t.Run("resume-old", func(t *testing.T) {
		t.Parallel()

		// Resuming isn't limited to the history window.
		ids, err := collectStream(t, db, f, 2*time.Second, database.GuildEventStreamConfig{
			GuildID:    quiet.ID,
			AfterID:    old.ID - 1,
			HistoryAge: time.Hour,
		})
		if err != nil {
			t.Fatal(err)
		}
		if want := []int{old.ID}; !slices.Equal(ids, want) {
			t.Errorf("event ids = %v, want %v", ids, want)
		}
	})

	t.Run("authorize-error", func(t *testing.T) {
		t.Parallel()

		errRevoked := errors.New("revoked")
		var calls int

		ids, err := collectStream(t, db, f, 2*time.Second, database.GuildEventStreamConfig{
			GuildID: f.GuildA.ID,
			Authorize: func(ctx context.Context) (context.Context, error) {
				calls++
				if calls > 2 {
					return nil, errRevoked
				}
				return ctx, nil
			},
		})
		if !errors.Is(err, errRevoked) {
			t.Fatalf("err = %v, want %v", err, errRevoked)
		}
		if want := []int{f.GuildAInfoEvent.ID, f.GuildADebugEvent.ID}; !slices.Equal(ids, want) {
			t.Errorf("event ids = %v, want %v", ids, want)
		}
	})
}

// TestGuildEventNotifier ensures events created in a transaction (as the guild
// logger does) wake streams for the same guild once committed, without waiting
// for the next poll. A notification before the commit would make the stream
// poll too early, and miss the events.
func TestGuildEventNotifier(t *testing.T) {
	t.Parallel()

	db := dbtest.New(t)
	f := dbtest.NewFixtures(t, db)
	allow := dbtest.AllowAll(context.Background())

	var notifier database.GuildEventNotifier
	db.GuildEvent.Use(notifier.Hook())

	notify, unsubscribe := notifier.Subscribe(f.GuildA.ID)
	defer unsubscribe()
	otherNotify, otherUnsubscribe := notifier.Subscribe(f.GuildB.ID)
	defer otherUnsubscribe()

	ctx, cancel := context.WithTimeout(dbtest.AsUser(context.Background(), f.Admin), 5*time.Second)
	defer cancel()

	polled := make(chan struct{}, 10)
	emitted := make(chan int, 10)
	go func() {
		_ = database.StreamGuildEvents(ctx, db, database.GuildEventStreamConfig{
			GuildID:  f.GuildA.ID,
			AfterID:  f.GuildBInfoEvent.ID, // Only stream new events.
			Interval: time.Hour,            // Only notifications can trigger polls.
			Notify:   notify,
			Authorize: func(ctx context.Context) (context.Context, error) {
				polled <- struct{}{}
				return ctx, nil
			},
		}, func(e *ent.GuildEvent) error {
			emitted <- e.ID
			return nil
		})
	}()

	// Wait for the initial poll to finish.
	<-polled
	time.Sleep(200 * time.Millisecond)

	var created []*ent.GuildEvent
	err := database.RunWithTx(allow, db, func(ctx context.Context, tx *ent.Tx) (err error) {
		created, err = tx.GuildEvent.CreateBulk(
			tx.GuildEvent.Create().SetGuildID(f.GuildA.ID).SetType(guildevent.TypeINFO).SetMessage("one"),
			tx.GuildEvent.Create().SetGuildID(f.GuildA.ID).SetType(guildevent.TypeINFO).SetMessage("two"),
		).Save(ctx)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, want := range created {
		select {
		case id := <-emitted:
			if id != want.ID {
				t.Errorf("event id = %d, want %d", id, want.ID)
			}
		case <-ctx.Done():
			t.Fatalf("event %d wasn't streamed after commit", want.ID)
		}
	}

	if len(otherNotify) != 0 {
		t.Error("other guild was notified")
	}
}
