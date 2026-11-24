// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package database_test

import (
	"context"
	"errors"
	"slices"
	"testing"

	"github.com/lrstanley/spectrograph/internal/database/dbtest"
	"github.com/lrstanley/spectrograph/internal/database/ent"
	"github.com/lrstanley/spectrograph/internal/database/ent/guildevent"
	"github.com/lrstanley/spectrograph/internal/database/ent/privacy"
)

func ids[T any](items []*T, fn func(*T) int) []int {
	out := make([]int, 0, len(items))
	for _, item := range items {
		out = append(out, fn(item))
	}
	slices.Sort(out)
	return out
}

func sorted(v ...int) []int {
	slices.Sort(v)
	return v
}

func TestPrivacyQueries(t *testing.T) {
	t.Parallel()

	db := dbtest.New(t)
	f := dbtest.NewFixtures(t, db)

	guildID := func(g *ent.Guild) int { return g.ID }
	userID := func(u *ent.User) int { return u.ID }
	eventID := func(e *ent.GuildEvent) int { return e.ID }

	tests := []struct {
		name       string
		as         *ent.User
		wantDeny   bool
		wantGuilds []int
		wantUsers  []int
		wantEvents []int
	}{
		{
			name:     "unauthenticated",
			as:       nil,
			wantDeny: true,
		},
		{
			name:       "guild-admin-a",
			as:         f.UserA,
			wantGuilds: []int{f.GuildA.ID},
			wantUsers:  []int{f.UserA.ID},
			wantEvents: []int{f.GuildAInfoEvent.ID}, // No DEBUG events.
		},
		{
			name:       "guild-admin-b",
			as:         f.UserB,
			wantGuilds: []int{f.GuildB.ID},
			wantUsers:  []int{f.UserB.ID},
			wantEvents: []int{f.GuildBInfoEvent.ID},
		},
		{
			name:       "system-admin",
			as:         f.Admin,
			wantGuilds: sorted(f.GuildA.ID, f.GuildB.ID),
			wantUsers:  sorted(f.Admin.ID, f.UserA.ID, f.UserB.ID),
			wantEvents: sorted(f.GuildAInfoEvent.ID, f.GuildADebugEvent.ID, f.GuildBInfoEvent.ID),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := dbtest.AsUser(context.Background(), tt.as)

			guilds, err := db.Guild.Query().All(ctx)
			if tt.wantDeny {
				if !errors.Is(err, privacy.Deny) {
					t.Fatalf("expected privacy deny for guilds, got: %v", err)
				}
				if _, err = db.User.Query().All(ctx); !errors.Is(err, privacy.Deny) {
					t.Fatalf("expected privacy deny for users, got: %v", err)
				}
				if _, err = db.GuildEvent.Query().All(ctx); !errors.Is(err, privacy.Deny) {
					t.Fatalf("expected privacy deny for events, got: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("querying guilds: %v", err)
			}
			if got := ids(guilds, guildID); !slices.Equal(got, tt.wantGuilds) {
				t.Errorf("guilds = %v, want %v", got, tt.wantGuilds)
			}

			users, err := db.User.Query().All(ctx)
			if err != nil {
				t.Fatalf("querying users: %v", err)
			}
			if got := ids(users, userID); !slices.Equal(got, tt.wantUsers) {
				t.Errorf("users = %v, want %v", got, tt.wantUsers)
			}

			events, err := db.GuildEvent.Query().All(ctx)
			if err != nil {
				t.Fatalf("querying events: %v", err)
			}
			if got := ids(events, eventID); !slices.Equal(got, tt.wantEvents) {
				t.Errorf("events = %v, want %v", got, tt.wantEvents)
			}

			// Edge traversal should also respect the DEBUG filter.
			if !tt.as.Admin {
				n, cerr := db.GuildEvent.Query().Where(guildevent.TypeEQ(guildevent.TypeDEBUG)).Count(ctx)
				if cerr != nil || n != 0 {
					t.Errorf("expected no DEBUG events, got %d (err: %v)", n, cerr)
				}
			}
		})
	}
}

func TestPrivacyMutations(t *testing.T) {
	t.Parallel()

	db := dbtest.New(t)
	f := dbtest.NewFixtures(t, db)

	ctxA := dbtest.AsUser(context.Background(), f.UserA)
	ctxAdmin := dbtest.AsUser(context.Background(), f.Admin)

	t.Run("guild-config-own-guild", func(t *testing.T) {
		t.Parallel()
		cfg, err := db.GuildConfig.UpdateOneID(f.GuildAConfig.ID).SetRegexMatch("^foo$").Save(ctxA)
		if err != nil {
			t.Fatalf("expected update to succeed: %v", err)
		}
		if cfg.RegexMatch != "^foo$" {
			t.Errorf("regex_match = %q, want ^foo$", cfg.RegexMatch)
		}
	})

	t.Run("guild-config-other-guild", func(t *testing.T) {
		t.Parallel()
		_, err := db.GuildConfig.UpdateOneID(f.GuildBConfig.ID).SetEnabled(false).Save(ctxA)
		if !ent.IsNotFound(err) {
			t.Fatalf("expected not found, got: %v", err)
		}
	})

	t.Run("guild-config-validation", func(t *testing.T) {
		t.Parallel()
		_, err := db.GuildConfig.UpdateOneID(f.GuildAConfig.ID).SetRegexMatch("(").Save(ctxA)
		if !ent.IsValidationError(err) {
			t.Fatalf("expected validation error, got: %v", err)
		}
	})

	t.Run("guild-admin-config-guild-admin", func(t *testing.T) {
		t.Parallel()
		_, err := db.GuildAdminConfig.UpdateOneID(f.GuildAAdminConfig.ID).SetEnabled(false).Save(ctxA)
		if !errors.Is(err, privacy.Deny) {
			t.Fatalf("expected privacy deny, got: %v", err)
		}
	})

	t.Run("guild-admin-config-system-admin", func(t *testing.T) {
		t.Parallel()
		_, err := db.GuildAdminConfig.UpdateOneID(f.GuildBAdminConfig.ID).SetComment("hello").Save(ctxAdmin)
		if err != nil {
			t.Fatalf("expected update to succeed: %v", err)
		}
	})

	t.Run("guild-config-create-denied", func(t *testing.T) {
		t.Parallel()
		_, err := db.GuildConfig.Create().SetGuildID(f.GuildA.ID).Save(ctxA)
		if !errors.Is(err, privacy.Deny) {
			t.Fatalf("expected privacy deny, got: %v", err)
		}
	})

	t.Run("user-self-update-denied", func(t *testing.T) {
		t.Parallel()
		_, err := db.User.UpdateOneID(f.UserA.ID).SetAdmin(true).Save(ctxA)
		if !errors.Is(err, privacy.Deny) {
			t.Fatalf("expected privacy deny, got: %v", err)
		}
	})

	t.Run("user-delete-other-denied", func(t *testing.T) {
		t.Parallel()
		err := db.User.DeleteOneID(f.UserB.ID).Exec(ctxA)
		if !ent.IsNotFound(err) {
			t.Fatalf("expected not found, got: %v", err)
		}
	})

	t.Run("guild-mutation-denied", func(t *testing.T) {
		t.Parallel()
		_, err := db.Guild.UpdateOneID(f.GuildA.ID).SetName("renamed").Save(ctxA)
		if !errors.Is(err, privacy.Deny) {
			t.Fatalf("expected privacy deny, got: %v", err)
		}
	})

	t.Run("banned-user-has-no-access", func(t *testing.T) {
		t.Parallel()
		banned := *f.UserB
		banned.Banned = true
		_, err := db.Guild.Query().All(dbtest.AsUser(context.Background(), &banned))
		if !errors.Is(err, privacy.Deny) {
			t.Fatalf("expected privacy deny, got: %v", err)
		}
	})
}

func TestPrivacyUserDeleteSelf(t *testing.T) {
	t.Parallel()

	db := dbtest.New(t)
	f := dbtest.NewFixtures(t, db)

	if err := db.User.DeleteOneID(f.UserA.ID).Exec(dbtest.AsUser(context.Background(), f.UserA)); err != nil {
		t.Fatalf("expected delete to succeed: %v", err)
	}

	// Guild should still exist, without any admins.
	n, err := db.Guild.Query().QueryAdmins().Count(dbtest.AllowAll(context.Background()))
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Errorf("admins = %d, want 1 (only user B)", n)
	}
}
