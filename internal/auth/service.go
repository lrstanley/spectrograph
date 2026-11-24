// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

// Package auth implements the xauth service used for Discord OAuth2
// authentication.
package auth

import (
	"context"
	"fmt"
	"slices"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/lrstanley/chix/xauth/v2"
	"github.com/lrstanley/spectrograph/internal/database"
	"github.com/lrstanley/spectrograph/internal/database/ent"
	"github.com/lrstanley/spectrograph/internal/database/ent/guild"
	"github.com/lrstanley/spectrograph/internal/database/ent/privacy"
	"github.com/lrstanley/spectrograph/internal/database/ent/user"
	"github.com/lrstanley/spectrograph/internal/discordapi"
	"github.com/lrstanley/spectrograph/internal/models"
	"github.com/markbates/goth"
)

var _ xauth.Service[ent.User, int] = (*Service)(nil)

// UserFetcher fetches the Discord user (and their guilds) associated with an
// OAuth2 access token. See [discordapi.FetchUser].
type UserFetcher func(ctx context.Context, token string) (*models.DiscordUser, []*models.UserGuildResponse, error)

// Service implements [xauth.Service] for Discord OAuth2 authentication.
type Service struct {
	db       *ent.Client
	adminIDs []string
	fetch    UserFetcher
}

// NewService returns a new [Service]. adminIDs are the Discord user IDs
// of users which should be system admins. If fetch is nil, [discordapi.FetchUser]
// is used.
func NewService(db *ent.Client, adminIDs []string, fetch UserFetcher) *Service {
	if fetch == nil {
		fetch = discordapi.FetchUser
	}
	return &Service{
		db:       db,
		adminIDs: adminIDs,
		fetch:    fetch,
	}
}

// Get returns the user with the provided ID. Banned users are still returned
// (with [ent.User.Banned] set), so the HTTP layer can tell them apart from
// signed out users and respond with [models.ErrUserBanned]. Privacy rules treat
// banned users as unauthenticated, so they can't access anything regardless.
func (s *Service) Get(ctx context.Context, id int) (*ent.User, error) {
	return s.db.User.Get(privacy.DecisionContext(ctx, privacy.Allow), id)
}

// Set creates or updates the user (and their associated guilds) after a
// successful OAuth2 flow, returning the users ID.
func (s *Service) Set(ctx context.Context, guser *goth.User) (id int, err error) {
	ctx = privacy.DecisionContext(ctx, privacy.Allow)

	duser, guilds, err := s.fetch(ctx, guser.AccessToken)
	if err != nil {
		return 0, fmt.Errorf("fetching user info from discord api: %w", err)
	}

	err = database.RunWithTx(ctx, s.db, func(ctx context.Context, tx *ent.Tx) error {
		id, err = s.set(ctx, tx, duser, guilds)
		return err
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Service) set(
	ctx context.Context,
	tx *ent.Tx,
	duser *models.DiscordUser,
	guilds []*models.UserGuildResponse,
) (int, error) {
	uid, err := tx.User.Create().
		SetUserID(duser.ID).
		SetAdmin(slices.Contains(s.adminIDs, duser.ID)).
		SetUsername(duser.Username).
		SetGlobalName(duser.GlobalName).
		SetDiscriminator(duser.Discriminator).
		SetEmail(duser.Email).
		SetAvatarHash(duser.Avatar).
		SetAvatarURL(discordapi.GenerateUserAvatarURL(duser)).
		SetLocale(duser.Locale).
		SetBot(duser.Bot).
		SetSystem(duser.System).
		SetMfaEnabled(duser.MFAEnabled).
		SetVerified(duser.Verified).
		SetFlags(duser.Flags).
		SetPremiumType(duser.PremiumType).
		SetPublicFlags(duser.PublicFlags).
		OnConflictColumns(user.FieldUserID).
		UpdateNewValues().
		ID(ctx)
	if err != nil {
		return 0, fmt.Errorf("upserting user: %w", err)
	}

	u, err := tx.User.Get(ctx, uid)
	if err != nil {
		return 0, fmt.Errorf("fetching user: %w", err)
	}
	if u.Banned {
		// The ban reason is only meant for other admins, so it isn't included.
		return 0, models.ErrUserBanned
	}

	// Guilds the user is currently associated with.
	currentGuilds, err := tx.Guild.Query().
		Where(guild.HasAdminsWith(user.ID(uid))).
		Select(guild.FieldGuildID).
		Strings(ctx)
	if err != nil {
		return 0, fmt.Errorf("fetching current user guilds: %w", err)
	}

	builders := make([]*ent.GuildCreate, 0, len(guilds))
	managed := make([]string, 0, len(guilds))

	for _, g := range guilds {
		// Only associate guilds the user can manage (and thus can add the bot to).
		if !discordapi.UserCanManageGuild(g) {
			continue
		}
		managed = append(managed, g.ID)

		q := tx.Guild.Create().
			SetGuildID(g.ID).
			SetName(g.Name).
			SetFeatures(g.Features).
			SetIconHash(g.Icon).
			SetIconURL(discordapi.GenerateGuildIconURL(g.ID, g.Icon))

		if g.Owner {
			q.SetOwnerID(duser.ID)
		}

		if !slices.Contains(currentGuilds, g.ID) {
			q.AddAdminIDs(uid)
		}

		builders = append(builders, q)
	}

	if len(builders) > 0 {
		err = tx.Guild.CreateBulk(builders...).
			OnConflict(
				entsql.ConflictColumns(guild.FieldGuildID),
				entsql.ResolveWith(resolveUserGuildConflict),
			).
			Exec(ctx)
		if err != nil {
			return 0, fmt.Errorf("upserting user guilds: %w", err)
		}
	}

	// Remove associations with guilds that the user is no longer in, or can no
	// longer manage.
	var stale []string
	for _, gid := range currentGuilds {
		if !slices.Contains(managed, gid) {
			stale = append(stale, gid)
		}
	}

	if len(stale) > 0 {
		err = tx.Guild.Update().
			Where(guild.GuildIDIn(stale...)).
			RemoveAdminIDs(uid).
			Exec(ctx)
		if err != nil {
			return 0, fmt.Errorf("removing stale user guild associations: %w", err)
		}
	}

	return uid, nil
}

// resolveUserGuildConflict only updates the guild fields which are sourced from
// the users OAuth2 guild list. The rest (bot permissions, member count, join
// time, etc) are maintained by the worker, and the insert would otherwise reset
// them to their defaults. The owner is only set when the user is the owner (the
// bulk insert uses NULL for the other rows).
func resolveUserGuildConflict(u *entsql.UpdateSet) {
	for _, c := range []string{
		guild.FieldName,
		guild.FieldFeatures,
		guild.FieldIconHash,
		guild.FieldIconURL,
		guild.FieldUpdateTime,
	} {
		u.SetExcluded(c)
	}
	u.Set(guild.FieldOwnerID, entsql.Expr(fmt.Sprintf(
		"COALESCE(%s, %s)",
		entsql.Dialect(u.Dialect()).Table("excluded").C(guild.FieldOwnerID),
		u.Table().C(guild.FieldOwnerID),
	)))
}
