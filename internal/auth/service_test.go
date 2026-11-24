// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package auth_test

import (
	"context"
	"errors"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/lrstanley/spectrograph/internal/auth"
	"github.com/lrstanley/spectrograph/internal/database/dbtest"
	"github.com/lrstanley/spectrograph/internal/database/ent/guild"
	"github.com/lrstanley/spectrograph/internal/database/ent/user"
	"github.com/lrstanley/spectrograph/internal/models"
	"github.com/markbates/goth"
)

type fakeDiscord struct {
	user   *models.DiscordUser
	guilds []*models.UserGuildResponse
}

func (f *fakeDiscord) fetch(_ context.Context, _ string) (*models.DiscordUser, []*models.UserGuildResponse, error) {
	return f.user, f.guilds, nil
}

func TestServiceSet(t *testing.T) {
	t.Parallel()

	db := dbtest.New(t)
	ctx := context.Background()
	allow := dbtest.AllowAll(ctx)

	discord := &fakeDiscord{
		user: &models.DiscordUser{
			ID:            "5001",
			Username:      "tester",
			GlobalName:    "Tester",
			Discriminator: "0",
			Email:         "tester@example.com",
		},
		guilds: []*models.UserGuildResponse{
			{ID: "6001", Name: "owned", Owner: true},
			{ID: "6002", Name: "managed", Permissions: models.DiscordPermManageServer},
			{ID: "6003", Name: "member", Permissions: models.DiscordPermSendMessages},
		},
	}

	svc := auth.NewService(db, []string{"5001"}, discord.fetch)

	id, err := svc.Set(ctx, &goth.User{AccessToken: "token"})
	if err != nil {
		t.Fatalf("Set() error: %v", err)
	}

	u, err := svc.Get(ctx, id)
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if !u.Admin || u.GlobalName != "Tester" || u.Discriminator != "0" {
		t.Errorf("unexpected user: %+v", u)
	}

	guildIDs := func() []string {
		t.Helper()
		gids, gerr := db.Guild.Query().
			Where(guild.HasAdminsWith(user.ID(id))).
			Select(guild.FieldGuildID).
			Strings(allow)
		if gerr != nil {
			t.Fatal(gerr)
		}
		slices.Sort(gids)
		return gids
	}

	if got := guildIDs(); !slices.Equal(got, []string{"6001", "6002"}) {
		t.Errorf("guilds = %v, want [6001 6002]", got)
	}

	// Lose manage permissions on one guild, and admin privileges, then re-auth.
	discord.guilds = discord.guilds[:1]
	svc = auth.NewService(db, nil, discord.fetch)

	if _, err = svc.Set(ctx, &goth.User{AccessToken: "token"}); err != nil {
		t.Fatalf("Set() error: %v", err)
	}

	if got := guildIDs(); !slices.Equal(got, []string{"6001"}) {
		t.Errorf("guilds = %v, want [6001]", got)
	}

	if u, err = svc.Get(ctx, id); err != nil || u.Admin {
		t.Errorf("expected admin to be revoked (err: %v)", err)
	}

	// Banned users should not be able to authenticate, and the ban reason (only
	// meant for other admins) shouldn't leak through the error.
	db.User.UpdateOneID(id).SetBanned(true).SetBanReason("spam").ExecX(allow)

	_, err = svc.Set(ctx, &goth.User{AccessToken: "token"})
	if !errors.Is(err, models.ErrUserBanned) {
		t.Errorf("Set() expected banned error, got: %v", err)
	}
	if err != nil && strings.Contains(err.Error(), "spam") {
		t.Errorf("Set() error leaks the ban reason: %v", err)
	}

	// Existing sessions of banned users still resolve (so the API can respond
	// with 403 rather than 401), with Banned set.
	if u, err = svc.Get(ctx, id); err != nil || !u.Banned {
		t.Errorf("Get() = %+v, %v; want banned user", u, err)
	}
}

// TestServiceSetPreservesGuildState ensures logging in only updates the
// guild fields sourced from the users OAuth2 guild list, and doesn't reset the
// fields maintained by the worker (bot permissions, member count, owner, etc).
func TestServiceSetPreservesGuildState(t *testing.T) {
	t.Parallel()

	db := dbtest.New(t)
	ctx := context.Background()
	allow := dbtest.AllowAll(ctx)

	joined := time.Now().Add(-time.Hour).Truncate(time.Second)
	existing := db.Guild.Create().
		SetGuildID("7001").
		SetName("old name").
		SetIconURL("").
		SetOwnerID("9999").
		SetPermissions(models.DiscordPermManageChannels).
		SetMemberCount(42).
		SetLarge(true).
		SetJoinedAt(joined).
		SaveX(allow)

	discord := &fakeDiscord{
		user: &models.DiscordUser{ID: "5002", Username: "tester"},
		guilds: []*models.UserGuildResponse{
			{ID: "7001", Name: "new name", Icon: "abc", Permissions: models.DiscordPermAdministrator},
			{ID: "7002", Name: "owned", Owner: true},
		},
	}

	svc := auth.NewService(db, nil, discord.fetch)
	if _, err := svc.Set(ctx, &goth.User{AccessToken: "token"}); err != nil {
		t.Fatalf("Set() error: %v", err)
	}

	g := db.Guild.GetX(allow, existing.ID)
	if g.Name != "new name" || g.IconHash != "abc" || g.IconURL == "" {
		t.Errorf("expected discord-sourced fields to be updated: %+v", g)
	}
	if g.Permissions != models.DiscordPermManageChannels || g.MemberCount != 42 || !g.Large {
		t.Errorf("worker-maintained fields were reset: %+v", g)
	}
	if g.JoinedAt == nil || !g.JoinedAt.Equal(joined) {
		t.Errorf("joined_at = %v, want %v", g.JoinedAt, joined)
	}
	if g.OwnerID != "9999" {
		t.Errorf("owner_id = %q, want %q (user isn't the owner)", g.OwnerID, "9999")
	}

	owned := db.Guild.Query().Where(guild.GuildID("7002")).OnlyX(allow)
	if owned.OwnerID != "5002" {
		t.Errorf("owner_id = %q, want %q", owned.OwnerID, "5002")
	}
}
