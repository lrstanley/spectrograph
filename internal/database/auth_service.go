// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package database

import (
	"context"
	"fmt"

	"github.com/lrstanley/chix"
	"github.com/lrstanley/spectrograph/internal/database/ent"
	"github.com/lrstanley/spectrograph/internal/database/ent/guild"
	"github.com/lrstanley/spectrograph/internal/database/ent/privacy"
	"github.com/lrstanley/spectrograph/internal/database/ent/user"
	"github.com/lrstanley/spectrograph/internal/discordapi"
	"github.com/lrstanley/spectrograph/internal/metrics"
	"github.com/lrstanley/spectrograph/internal/models"
	"github.com/markbates/goth"
	"golang.org/x/exp/slices"
)

// Validate authService implements chix.AuthService.
var _ chix.AuthService[ent.User, int] = (*authService)(nil)

func NewAuthService(db *ent.Client, adminIDs []string) *authService {
	return &authService{
		db:       db,
		adminIDs: adminIDs,
	}
}

type authService struct {
	db       *ent.Client
	adminIDs []string
}

func (s *authService) Get(ctx context.Context, id int) (*ent.User, error) {
	return s.db.User.Get(privacy.DecisionContext(ctx, privacy.Allow), id)
}

func (s *authService) Set(ctx context.Context, guser *goth.User) (id int, err error) {
	ctx = privacy.DecisionContext(ctx, privacy.Allow)

	metrics.AuthCount.Inc()

	db, err := s.db.Tx(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to start transaction: %w", err)
	}

	id, err = s.set(ctx, db, guser)

	err = Commit(db, err)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *authService) set(ctx context.Context, db *ent.Tx, guser *goth.User) (id int, err error) {
	duser, guilds, err := discordapi.FetchUser(ctx, guser.AccessToken)
	if err != nil {
		return 0, fmt.Errorf("error fetching user info from discord api: %w", err)
	}

	userBuilder := db.User.Create().
		SetUserID(duser.ID.String()).
		SetUsername(duser.Username).
		SetDiscriminator(duser.Discriminator.String()).
		SetEmail(duser.Email).
		SetAvatarHash(duser.Avatar).
		SetAvatarURL(discordapi.GenerateUserAvatarURL(duser)).
		SetLocale(duser.Locale).
		SetBot(duser.Bot).
		SetSystem(duser.System).
		SetMfaEnabled(duser.MFAEnabled).
		SetVerified(duser.Verified).
		SetFlags(uint64(duser.Flags)).
		SetPremiumType(int(duser.PremiumType)).
		SetPublicFlags(uint64(duser.PublicFlags))

	for _, admin := range s.adminIDs {
		if admin == duser.ID.String() {
			userBuilder.SetAdmin(true)
			break
		}
	}

	// Update user.
	uid, err := userBuilder.OnConflictColumns(user.FieldUserID).UpdateNewValues().ID(ctx)
	if err != nil {
		return 0, fmt.Errorf("error creating user in db: %w", err)
	}

	// Get current users guilds.
	currentGuilds, err := s.db.Guild.Query().
		Where(guild.HasAdminsWith(user.ID(uid))).
		Select(guild.FieldGuildID).
		Strings(ctx)
	if err != nil {
		return 0, fmt.Errorf("error fetching current user guilds from db: %w", err)
	}

	guildBuilders := make([]*ent.GuildCreate, 0, len(guilds))

	// Generate builders to update/create all user guilds.
	for _, guild := range guilds {
		isGuildAdmin := discordapi.UserHasAdmin(duser, guild)

		// Check if the user is an admin or an owner (and thus has permissions to add the bot connection).
		if !guild.Owner && !isGuildAdmin {
			continue
		}

		q := db.Guild.Create().
			SetGuildID(guild.ID).
			SetName(guild.Name).
			SetFeatures(guild.Features).
			SetIconHash(guild.Icon).
			SetIconURL(discordapi.GenerateGuildIconURL(guild.ID, guild.Icon)).
			SetPermissions(uint64(guild.Permissions))

		if guild.Owner {
			q.SetOwnerID(duser.ID.String())
		}

		if !slices.Contains(currentGuilds, guild.ID) {
			q.AddAdminIDs(uid)
		}

		guildBuilders = append(guildBuilders, q)
	}

	// Bulk update user guilds.
	err = db.Guild.CreateBulk(guildBuilders...).
		OnConflictColumns(guild.FieldGuildID).UpdateNewValues().Exec(ctx)
	if err != nil {
		return 0, err
	}

	// Remove any guild associations that the user is no longer in (or an admin of).
	for _, guildID := range currentGuilds {
		var matches bool
		for _, guild := range guilds {
			if guild.ID == guildID {
				matches = true
				break
			}
		}

		if matches {
			continue
		}

		// Remove the user from the guild.
		err = db.Guild.Update().
			Where(guild.GuildID(guildID)).
			RemoveAdminIDs(uid).Exec(ctx)

		if err != nil {
			return 0, fmt.Errorf("error removing user from guild (as they no longer have sufficient permission): %w", err)
		}
	}

	// Check user status.
	var u *ent.User
	u, err = db.User.Get(ctx, uid)
	if err != nil {
		return 0, fmt.Errorf("error fetching user status: %w", err)
	}

	if u.Banned {
		return 0, &models.ErrUserBanned{User: u}
	}

	return uid, nil
}

func (s *authService) Roles(ctx context.Context, id int) (roles []string, err error) {
	ctx = privacy.DecisionContext(ctx, privacy.Allow)

	var u *ent.User

	u, err = s.db.User.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error fetching user from db: %w", err)
	}

	roles = []string{models.RoleUser}

	if u.Admin {
		roles = append(roles, models.RoleSystemAdmin)
	}

	return roles, nil
}
