// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package worker

import (
	"github.com/andersfylling/disgord"
	"github.com/lrstanley/spectrograph/internal/database"
	"github.com/lrstanley/spectrograph/internal/discordapi"
	"github.com/lrstanley/spectrograph/internal/ent"
	"github.com/lrstanley/spectrograph/internal/ent/guild"
	"github.com/lrstanley/spectrograph/internal/ent/guildadminconfig"
	"github.com/lrstanley/spectrograph/internal/ent/guildconfig"
	"github.com/lrstanley/spectrograph/internal/models"
)

// dbUpdateGuild will update the guild in the database with the latest information
// from Discord. It create the guild, guild config, and guild admin config if they
// do not exist.
func (w *Worker) dbUpdateGuild(g *disgord.Guild, permissions models.DiscordPermissions, enabled bool) (err error) {
	tx, err := w.db.Tx(w.ctx)
	if err != nil {
		return err
	}

	q := w.db.Guild.Create().
		SetGuildID(g.ID.String()).
		SetName(g.Name).
		SetFeatures(g.Features).
		SetIconHash(g.Icon).
		SetIconURL(discordapi.GenerateGuildIconURL(g.ID.String(), g.Icon)).
		SetJoinedAt(g.JoinedAt.Time).
		SetLarge(g.Large).
		SetMemberCount(int(g.MemberCount)).
		SetOwnerID(g.OwnerID.String()).
		SetSystemChannelFlags(g.SystemChannelID.String())

	if permissions > 0 {
		q = q.SetPermissions(uint64(permissions))
	}

	var gid int
	gid, err = q.OnConflictColumns(guild.FieldGuildID).UpdateNewValues().ID(w.ctx)
	if err != nil {
		return database.Commit(tx, err)
	}

	err = w.db.GuildConfig.Create().
		SetGuildID(gid).
		OnConflictColumns(guildconfig.GuildColumn).UpdateNewValues().
		Exec(w.ctx)
	if err != nil {
		return database.Commit(tx, err)
	}

	err = w.db.GuildAdminConfig.Create().
		SetGuildID(gid).
		SetEnabled(enabled).
		OnConflictColumns(guildadminconfig.GuildColumn).UpdateNewValues().
		Exec(w.ctx)
	if err != nil {
		return database.Commit(tx, err)
	}

	return database.Commit(tx, nil)
}

// dbDisableGuild will disable the guild in the database (typically done if the bot
// gets removed from the guild).
func (w *Worker) dbDisableGuild(gid string) (err error) {
	_, err = w.db.GuildConfig.Update().
		Where(guildconfig.HasGuildWith(guild.GuildID(gid))).
		SetEnabled(false).
		Save(w.ctx)

	if err == nil {
		_, err = w.db.Guild.Update().
			Where(guild.GuildID(gid)).
			ClearJoinedAt().
			Save(w.ctx)
	}

	return err
}

// dbGetSettings will return the guild config and guild admin config for the
// given guild ID.
func (w *Worker) dbGetSettings(gid string) (config *ent.GuildConfig, adminconfig *ent.GuildAdminConfig, err error) {
	config, err = w.db.GuildConfig.Query().
		Where(guildconfig.HasGuildWith(guild.GuildID(gid))).
		Only(w.ctx)
	if err != nil {
		return nil, nil, err
	}

	adminconfig, err = w.db.GuildAdminConfig.Query().
		Where(guildadminconfig.HasGuildWith(guild.GuildID(gid))).
		Only(w.ctx)
	if err != nil {
		return nil, nil, err
	}

	return config, adminconfig, nil
}

// dbIsGuildEnabled will return true if the guild is enabled in the database
// (respecting the guild admin config first).
func (w *Worker) dbIsGuildEnabled(gid string) bool {
	config, adminconfig, err := w.dbGetSettings(gid)
	if err != nil {
		w.es.GuildID(gid).Error("failed to get guild settings")
		return false
	}

	if !adminconfig.Enabled {
		return false
	}

	return config.Enabled
}

// dbGetMaxClones will return the max clones of a channel, for the given guild ID.
func (w *Worker) dbGetMaxClones(gid string) int {
	config, adminconfig, err := w.dbGetSettings(gid)
	if err != nil {
		w.es.GuildID(gid).WithError(err).Error("failed to get guild settings")
		return models.DefaultMaxClones
	}

	if config.DefaultMaxClones == 0 {
		config.DefaultMaxClones = models.DefaultMaxClones
	}

	if adminconfig.DefaultMaxClones == 0 {
		adminconfig.DefaultMaxClones = models.DefaultMaxClones
	}

	if config.DefaultMaxClones < adminconfig.DefaultMaxClones {
		return config.DefaultMaxClones
	}

	return adminconfig.DefaultMaxClones
}

// dbGetMaxChannels will return the max channels that can be monitored, for the
// given guild ID.
func (w *Worker) dbGetMaxChannels(gid string) int {
	_, adminconfig, err := w.dbGetSettings(gid)
	if err != nil {
		w.es.GuildID(gid).WithError(err).Error("failed to get guild settings")
		return models.DefaultMaxChannels
	}

	if adminconfig.DefaultMaxChannels == 0 {
		return models.DefaultMaxChannels
	}

	return adminconfig.DefaultMaxChannels
}
