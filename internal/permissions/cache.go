// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package permissions

import (
	"context"
	"errors"
	"sync"

	"github.com/andersfylling/disgord"
	"github.com/apex/log"
	"github.com/lrstanley/spectrograph/internal/guildlogger"
	"github.com/lrstanley/spectrograph/internal/models"
)

type Cache struct {
	es *guildlogger.EventStream

	mu sync.RWMutex
	db map[disgord.Snowflake]*BotPermissions
}

// NewPermissionsCache returns a new permissions cache.
func NewCache(ctx context.Context, es *guildlogger.EventStream) *Cache {
	return &Cache{
		es: es,
		db: make(map[disgord.Snowflake]*BotPermissions),
	}
}

func (pc *Cache) GuildCreate(sess disgord.Session, event *disgord.GuildCreate) (models.DiscordPermissions, error) {
	var botMember *disgord.Member

	// The bot should be in the list of members returned during the guild create
	// message (even if no other users are listed due to not having the
	// permissions). Find this so we can understand what permissions we have.
	bot, err := sess.CurrentUser().Get()
	if err != nil {
		return 0, err
	}

	for _, member := range event.Guild.Members {
		if member.UserID.String() == bot.ID.String() {
			botMember = member
			break
		}
	}
	if botMember == nil {
		return 0, errors.New("no member found matching ID of bot")
	}

	permissions := MergeRolePermissions(botMember, event.Guild.Roles)

	pc.mu.Lock()
	defer pc.mu.Unlock()

	pc.db[event.Guild.ID] = &BotPermissions{
		member:      botMember,
		permissions: permissions,
	}

	return permissions, nil
}

func (pc *Cache) GuildUpdate(sess disgord.Session, event *disgord.GuildUpdate) (models.DiscordPermissions, error) {
	return pc.GuildCreate(sess, &disgord.GuildCreate{Guild: event.Guild, ShardID: event.ShardID})
}

// TODO: this doesn't seem to trigger when we add or remove roles from ourselves.
func (pc *Cache) GuildMemberUpdate(sess disgord.Session, event *disgord.GuildMemberUpdate) {
	guild, err := sess.Guild(event.GuildID).Get()
	if err != nil {
		pc.es.GuildID(event.GuildID.String()).WithError(err).Error("unable to fetch guild for member-update permissions checks")
		return
	}

	pc.mu.Lock()
	defer pc.mu.Unlock()

	permissions := MergeRolePermissions(event.Member, guild.Roles)

	if previous, ok := pc.db[event.GuildID]; ok {
		if previous.permissions == permissions {
			return
		}

		// Permissions have changed.
		// TODO: should we trigger an event of some kind to log to the db?
		pc.es.Guild(guild).WithFields(log.Fields{
			"old_permissions": previous.permissions,
			"new_permissions": permissions,
			"member_id":       event.Member.UserID,
		})
	}

	pc.db[event.GuildID] = &BotPermissions{
		member:      event.Member,
		permissions: permissions,
	}
}

func (pc *Cache) GuildRoleChange(sess disgord.Session, guildID disgord.Snowflake) {
	// Check if it's a role ID that has been given to us.
	pc.mu.Lock()
	defer pc.mu.Unlock()

	perms, ok := pc.db[guildID]
	if !ok {
		return
	}

	roles, err := sess.Guild(guildID).GetRoles()
	if err != nil {
		pc.es.GuildID(guildID.String()).WithError(err).Error("unable to fetch roles for role-update permissions checks")
		return
	}

	pc.db[guildID].permissions = MergeRolePermissions(perms.member, roles)
}

func (pc *Cache) RemoveGuild(guildID disgord.Snowflake) {
	if guildID.IsZero() {
		return
	}

	pc.mu.Lock()
	delete(pc.db, guildID)
	pc.mu.Unlock()
}

func (pc *Cache) Get(guildID disgord.Snowflake) models.DiscordPermissions {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	if perms, ok := pc.db[guildID]; ok {
		return perms.permissions
	}

	return 0
}

// MergeRolePermissions converts the role id's that it has permission for,
// and generates a permission bit based of all access.
func MergeRolePermissions(member *disgord.Member, roles []*disgord.Role) (perms models.DiscordPermissions) {
	var ok bool

	idmap := make(map[disgord.Snowflake]struct{}, len(member.Roles))
	for _, role := range member.Roles {
		idmap[role] = struct{}{}
	}

	// Go through roles in reverse order.
	for i := len(roles) - 1; i >= 0; i-- {
		if _, ok = idmap[roles[i].ID]; !ok {
			continue
		}

		perms |= models.DiscordPermissions(roles[i].Permissions)
	}

	return perms
}
