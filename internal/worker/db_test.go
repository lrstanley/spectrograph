// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package worker

import (
	"context"
	"testing"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/lrstanley/spectrograph/internal/database/dbtest"
	"github.com/lrstanley/spectrograph/internal/database/ent"
	"github.com/lrstanley/spectrograph/internal/database/ent/guild"
	"github.com/lrstanley/spectrograph/internal/database/ent/privacy"
	"github.com/lrstanley/spectrograph/internal/models"
)

func testGuildCreate(joined time.Time) *gateway.GuildCreateEvent {
	return &gateway.GuildCreateEvent{
		ID:                 testGuildID,
		Name:               "test guild",
		Icon:               "a_abcdef",
		OwnerID:            testOwnerID,
		Features:           []discord.GuildFeature{"COMMUNITY"},
		SystemChannelFlags: 5,
		SystemChannelID:    1234567,
		Joined:             discord.NewTimestamp(joined),
		Large:              true,
		MemberCount:        42,
	}
}

func getGuild(ctx context.Context, t *testing.T, db *ent.Client) *ent.Guild {
	t.Helper()
	g, err := db.Guild.Query().
		Where(guild.GuildID(testGuildID.String())).
		WithGuildConfig().
		WithGuildAdminConfig().
		Only(ctx)
	if err != nil {
		t.Fatalf("querying guild: %v", err)
	}
	return g
}

func TestDBGuildLifecycle(t *testing.T) {
	t.Parallel()

	db := dbtest.New(t)
	ctx := privacy.DecisionContext(t.Context(), privacy.Allow)
	joined := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)

	if err := dbUpsertGuild(ctx, db, testGuildCreate(joined), discord.PermissionManageChannels); err != nil {
		t.Fatalf("dbUpsertGuild() = %v", err)
	}

	g := getGuild(ctx, t, db)
	if g.Name != "test guild" || g.MemberCount != 42 || !g.Large || g.OwnerID != testOwnerID.String() ||
		g.IconHash != "a_abcdef" || g.IconURL == "" || len(g.Features) != 1 ||
		g.JoinedAt == nil || !g.JoinedAt.Equal(joined) ||
		g.Permissions != uint64(discord.PermissionManageChannels) {
		t.Errorf("unexpected guild after create: %+v", g)
	}
	if g.SystemChannelFlags != "5" {
		t.Errorf("system channel flags = %q, want %q (flags, not the system channel ID)", g.SystemChannelFlags, "5")
	}
	if g.Edges.GuildConfig == nil || g.Edges.GuildAdminConfig == nil {
		t.Fatal("guild configs not created")
	}

	// User and admin choices must survive the bot reconnecting.
	err := db.GuildConfig.UpdateOne(g.Edges.GuildConfig).SetRegexMatch(`^x$`).SetEnabled(false).Exec(ctx)
	if err != nil {
		t.Fatal(err)
	}
	err = db.GuildAdminConfig.UpdateOne(g.Edges.GuildAdminConfig).SetEnabled(false).SetDefaultMaxClones(2).Exec(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if err = dbUpsertGuild(ctx, db, testGuildCreate(joined), discord.PermissionAdministrator); err != nil {
		t.Fatalf("dbUpsertGuild() (again) = %v", err)
	}

	g = getGuild(ctx, t, db)
	if g.Edges.GuildConfig.Enabled || g.Edges.GuildConfig.RegexMatch != `^x$` {
		t.Errorf("guild config overwritten: %+v", g.Edges.GuildConfig)
	}
	if g.Edges.GuildAdminConfig.Enabled || g.Edges.GuildAdminConfig.DefaultMaxClones != 2 {
		t.Errorf("guild admin config overwritten: %+v", g.Edges.GuildAdminConfig)
	}
	if g.Permissions != uint64(discord.PermissionAdministrator) {
		t.Errorf("permissions = %d, want %d", g.Permissions, discord.PermissionAdministrator)
	}

	// Guild updates don't carry the join time, size, or member count.
	update := testGuildCreate(joined).Guild
	update.Name = "renamed"
	if err = dbUpdateGuild(ctx, db, &update); err != nil {
		t.Fatalf("dbUpdateGuild() = %v", err)
	}

	g = getGuild(ctx, t, db)
	if g.Name != "renamed" || g.MemberCount != 42 || !g.Large || g.JoinedAt == nil {
		t.Errorf("unexpected guild after update: %+v", g)
	}

	if err = dbUpdateGuildPermissions(ctx, db, testGuildID, discord.PermissionViewChannel); err != nil {
		t.Fatalf("dbUpdateGuildPermissions() = %v", err)
	}
	if g = getGuild(ctx, t, db); g.Permissions != uint64(discord.PermissionViewChannel) {
		t.Errorf("permissions = %d, want %d", g.Permissions, discord.PermissionViewChannel)
	}

	settings, err := dbGetSettings(ctx, db, testGuildID)
	if err != nil {
		t.Fatalf("dbGetSettings() = %v", err)
	}
	if settings.enabled || settings.maxClones != 2 || settings.maxGroups != models.DefaultMaxChannels || settings.match.String() != `^x$` {
		t.Errorf("unexpected settings: %+v", settings)
	}

	if err = dbGuildRemoved(ctx, db, testGuildID); err != nil {
		t.Fatalf("dbGuildRemoved() = %v", err)
	}

	g = getGuild(ctx, t, db)
	if g.JoinedAt != nil || g.Permissions != 0 {
		t.Errorf("guild not marked as removed: %+v", g)
	}
	if g.Edges.GuildConfig == nil || g.Edges.GuildConfig.RegexMatch != `^x$` {
		t.Errorf("guild config not kept after removal: %+v", g.Edges.GuildConfig)
	}
}

func TestDBGetSettingsMissing(t *testing.T) {
	t.Parallel()

	db := dbtest.New(t)
	ctx := privacy.DecisionContext(t.Context(), privacy.Allow)

	if _, err := dbGetSettings(ctx, db, testGuildID); err == nil {
		t.Error("dbGetSettings() = nil error for unknown guild")
	}

	// Guilds created through the auth flow have no configs until the bot joins.
	err := db.Guild.Create().SetGuildID(testGuildID.String()).SetName("test guild").SetIconURL("").Exec(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = dbGetSettings(ctx, db, testGuildID); err == nil {
		t.Error("dbGetSettings() = nil error for guild without configs")
	}
}
