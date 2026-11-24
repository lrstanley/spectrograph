// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

// Package dbtest provides helpers for tests which require a Postgres database.
package dbtest

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"sync"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/lrstanley/chix/xauth/v2"
	"github.com/lrstanley/spectrograph/internal/database"
	"github.com/lrstanley/spectrograph/internal/database/ent"
	"github.com/lrstanley/spectrograph/internal/database/ent/guildevent"
	"github.com/lrstanley/spectrograph/internal/database/ent/privacy"
)

// EnvDatabaseURL is the environment variable containing the Postgres URL to use
// for tests. The user must be able to create and drop databases. Tests which
// require a database are skipped if unset.
const EnvDatabaseURL = "TEST_DATABASE_URL"

// New creates a new, empty, and migrated database (dropped when the test
// completes), and returns an ent client connected to it. Skips the test if
// [EnvDatabaseURL] isn't set.
func New(tb testing.TB) *ent.Client {
	tb.Helper()

	base := os.Getenv(EnvDatabaseURL)
	if base == "" {
		tb.Skipf("%s not set, skipping database test", EnvDatabaseURL)
	}

	uri, err := url.Parse(base)
	if err != nil {
		tb.Fatalf("parsing %s: %v", EnvDatabaseURL, err)
	}

	ctx := context.Background()

	admin, err := pgx.Connect(ctx, base)
	if err != nil {
		tb.Fatalf("connecting to test database: %v", err)
	}
	tb.Cleanup(func() { _ = admin.Close(context.Background()) })

	b := make([]byte, 6)
	_, _ = rand.Read(b)
	name := "spectrograph_test_" + hex.EncodeToString(b)

	if _, err = admin.Exec(ctx, "CREATE DATABASE "+name); err != nil {
		tb.Fatalf("creating test database: %v", err)
	}
	tb.Cleanup(func() {
		_, _ = admin.Exec(context.Background(), "DROP DATABASE IF EXISTS "+name+" WITH (FORCE)")
	})

	uri.Path = "/" + name

	logger := slog.New(slog.DiscardHandler)
	db, err := database.Open(ctx, logger, uri)
	if err != nil {
		tb.Fatalf("opening test database: %v", err)
	}
	tb.Cleanup(func() { _ = db.Close() })

	// ent's migration mutates the package-level table definitions, so concurrent
	// migrations (from parallel tests) race.
	migrateMu.Lock()
	err = database.Migrate(ctx, logger, db)
	migrateMu.Unlock()
	if err != nil {
		tb.Fatalf("migrating test database: %v", err)
	}

	return db
}

var migrateMu sync.Mutex

// AsUser returns a context authenticated as the provided user (or
// unauthenticated if nil).
func AsUser(ctx context.Context, u *ent.User) context.Context {
	if u == nil {
		return ctx
	}
	return xauth.OverrideContextAuth(ctx, u.ID, u)
}

// AllowAll returns a context which bypasses all privacy rules.
func AllowAll(ctx context.Context) context.Context {
	return privacy.DecisionContext(ctx, privacy.Allow)
}

// Fixtures holds a common set of entities for testing privacy rules.
type Fixtures struct {
	Admin  *ent.User // System admin, not associated with any guilds.
	UserA  *ent.User // Admin of GuildA.
	UserB  *ent.User // Admin of GuildB.
	GuildA *ent.Guild
	GuildB *ent.Guild

	GuildAConfig      *ent.GuildConfig
	GuildBConfig      *ent.GuildConfig
	GuildAAdminConfig *ent.GuildAdminConfig
	GuildBAdminConfig *ent.GuildAdminConfig

	GuildAInfoEvent  *ent.GuildEvent
	GuildADebugEvent *ent.GuildEvent
	GuildBInfoEvent  *ent.GuildEvent
}

// NewFixtures creates a common set of entities for testing privacy rules.
func NewFixtures(tb testing.TB, db *ent.Client) *Fixtures {
	tb.Helper()

	ctx := AllowAll(context.Background())
	f := &Fixtures{}

	user := func(id string, admin bool) *ent.User {
		return db.User.Create().
			SetUserID(id).
			SetUsername("user-" + id).
			SetEmail(id + "@example.com").
			SetAvatarURL("https://cdn.discordapp.com/embed/avatars/0.png").
			SetAdmin(admin).
			SaveX(ctx)
	}

	f.Admin = user("1000", true)
	f.UserA = user("1001", false)
	f.UserB = user("1002", false)

	guild := func(id string, admin *ent.User) (*ent.Guild, *ent.GuildConfig, *ent.GuildAdminConfig) {
		g := db.Guild.Create().
			SetGuildID(id).
			SetName("guild-" + id).
			SetIconURL("").
			AddAdmins(admin).
			SaveX(ctx)
		return g,
			db.GuildConfig.Create().SetGuild(g).SaveX(ctx),
			db.GuildAdminConfig.Create().SetGuild(g).SaveX(ctx)
	}

	f.GuildA, f.GuildAConfig, f.GuildAAdminConfig = guild("2001", f.UserA)
	f.GuildB, f.GuildBConfig, f.GuildBAdminConfig = guild("2002", f.UserB)

	event := func(g *ent.Guild, typ guildevent.Type) *ent.GuildEvent {
		return db.GuildEvent.Create().
			SetGuild(g).
			SetType(typ).
			SetMessage(fmt.Sprintf("%s event for %s", typ, g.Name)).
			SaveX(ctx)
	}

	f.GuildAInfoEvent = event(f.GuildA, guildevent.TypeINFO)
	f.GuildADebugEvent = event(f.GuildA, guildevent.TypeDEBUG)
	f.GuildBInfoEvent = event(f.GuildB, guildevent.TypeINFO)

	return f
}
