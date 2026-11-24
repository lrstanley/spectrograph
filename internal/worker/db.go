// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package worker

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/lrstanley/spectrograph/internal/database"
	"github.com/lrstanley/spectrograph/internal/database/ent"
	"github.com/lrstanley/spectrograph/internal/database/ent/guild"
	"github.com/lrstanley/spectrograph/internal/database/ent/guildadminconfig"
	"github.com/lrstanley/spectrograph/internal/database/ent/guildconfig"
	"github.com/lrstanley/spectrograph/internal/discordapi"
	"github.com/lrstanley/spectrograph/internal/models"
)

// guildFeatures converts Discord guild features to strings.
func guildFeatures(g *discord.Guild) []string {
	features := make([]string, 0, len(g.Features))
	for _, f := range g.Features {
		features = append(features, string(f))
	}
	return features
}

// dbUpsertGuild creates or updates the guild in the database from a guild
// create event, and ensures the guild's config and admin config exist. Existing
// configs are left untouched, so user and admin choices (e.g. a guild disabled
// by an admin) are preserved.
func dbUpsertGuild(ctx context.Context, db *ent.Client, ev *gateway.GuildCreateEvent, perms discord.Permissions) error {
	return database.RunWithTx(ctx, db, func(ctx context.Context, tx *ent.Tx) error {
		gid, err := tx.Guild.Create().
			SetGuildID(ev.ID.String()).
			SetName(ev.Name).
			SetFeatures(guildFeatures(&ev.Guild)).
			SetIconHash(ev.Icon).
			SetIconURL(discordapi.GenerateGuildIconURL(ev.ID.String(), ev.Icon)).
			SetJoinedAt(ev.Joined.Time()).
			SetLarge(ev.Large).
			SetMemberCount(int(ev.MemberCount)). //nolint:gosec // member counts are far below MaxInt.
			SetOwnerID(ev.OwnerID.String()).
			SetSystemChannelFlags(strconv.FormatUint(uint64(ev.SystemChannelFlags), 10)).
			SetPermissions(uint64(perms)).
			OnConflictColumns(guild.FieldGuildID).
			UpdateNewValues().
			ID(ctx)
		if err != nil {
			return fmt.Errorf("upserting guild: %w", err)
		}

		err = tx.GuildConfig.Create().
			SetGuildID(gid).
			OnConflictColumns(guildconfig.GuildColumn).
			Ignore().
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("upserting guild config: %w", err)
		}

		err = tx.GuildAdminConfig.Create().
			SetGuildID(gid).
			OnConflictColumns(guildadminconfig.GuildColumn).
			Ignore().
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("upserting guild admin config: %w", err)
		}
		return nil
	})
}

// dbUpdateGuild updates the guild metadata in the database, from a guild update
// event (which lacks the join time, size, and member count of a guild create
// event, so those are left untouched).
func dbUpdateGuild(ctx context.Context, db *ent.Client, g *discord.Guild) error {
	err := db.Guild.Update().
		Where(guild.GuildID(g.ID.String())).
		SetName(g.Name).
		SetFeatures(guildFeatures(g)).
		SetIconHash(g.Icon).
		SetIconURL(discordapi.GenerateGuildIconURL(g.ID.String(), g.Icon)).
		SetOwnerID(g.OwnerID.String()).
		SetSystemChannelFlags(strconv.FormatUint(uint64(g.SystemChannelFlags), 10)).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("updating guild: %w", err)
	}
	return nil
}

// dbUpdateGuildPermissions updates the bot's permissions in the guild.
func dbUpdateGuildPermissions(ctx context.Context, db *ent.Client, gid discord.GuildID, perms discord.Permissions) error {
	err := db.Guild.Update().
		Where(guild.GuildID(gid.String())).
		SetPermissions(uint64(perms)).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("updating guild permissions: %w", err)
	}
	return nil
}

// dbGuildRemoved marks the guild as no longer joined (typically done when the
// bot is removed from the guild). The guild configs are kept, so the guild
// works as before if the bot is invited again.
func dbGuildRemoved(ctx context.Context, db *ent.Client, gid discord.GuildID) error {
	err := db.Guild.Update().
		Where(guild.GuildID(gid.String())).
		ClearJoinedAt().
		SetPermissions(0).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("marking guild as removed: %w", err)
	}
	return nil
}

// dbJoinedGuilds returns the IDs of the guilds the bot is in, according to the
// database.
func dbJoinedGuilds(ctx context.Context, db *ent.Client) ([]discord.GuildID, error) {
	raw, err := db.Guild.Query().Where(guild.JoinedAtNotNil()).Select(guild.FieldGuildID).Strings(ctx)
	if err != nil {
		return nil, fmt.Errorf("querying joined guilds: %w", err)
	}

	ids := make([]discord.GuildID, 0, len(raw))
	for _, v := range raw {
		sf, perr := discord.ParseSnowflake(v)
		if perr != nil {
			continue
		}
		ids = append(ids, discord.GuildID(sf))
	}
	return ids, nil
}

// errInvalidRegex is returned (joined with the effective settings) when the
// guild's configured regex is invalid, in which case the default is used.
var errInvalidRegex = errors.New("invalid channel regex, using the default")

// dbGetSettings returns the effective settings for the given guild. If the
// configured regex is invalid, the settings use the default regex, and an
// error wrapping errInvalidRegex is returned alongside them.
func dbGetSettings(ctx context.Context, db *ent.Client, gid discord.GuildID) (guildSettings, error) {
	g, err := db.Guild.Query().
		Where(guild.GuildID(gid.String())).
		WithGuildConfig().
		WithGuildAdminConfig().
		Only(ctx)
	if err != nil {
		return guildSettings{}, fmt.Errorf("fetching guild settings: %w", err)
	}

	if g.Edges.GuildConfig == nil || g.Edges.GuildAdminConfig == nil {
		return guildSettings{}, errors.New("guild config not found")
	}

	return effectiveSettings(g.Edges.GuildConfig, g.Edges.GuildAdminConfig)
}

// effectiveSettings merges the user-defined config, the admin config, and the
// system defaults. Admin limits cap user-defined limits, and zero values use
// the system defaults.
func effectiveSettings(cfg *ent.GuildConfig, admin *ent.GuildAdminConfig) (guildSettings, error) {
	s := guildSettings{
		enabled:   cfg.Enabled && admin.Enabled,
		match:     models.DefaultChannelMatch,
		maxGroups: orDefault(admin.DefaultMaxChannels, models.DefaultMaxChannels),
		maxClones: min(
			orDefault(cfg.DefaultMaxClones, models.DefaultMaxClones),
			orDefault(admin.DefaultMaxClones, models.DefaultMaxClones),
		),
	}

	if cfg.RegexMatch != "" {
		rgx, err := regexp.Compile(cfg.RegexMatch)
		if err != nil {
			return s, fmt.Errorf("%w: %w", errInvalidRegex, err)
		}
		s.match = rgx
	}

	return s, nil
}

// orDefault returns v if positive, otherwise def.
func orDefault(v, def int) int {
	if v > 0 {
		return v
	}
	return def
}
