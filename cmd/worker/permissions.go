// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"errors"
	"sync"

	"github.com/andersfylling/disgord"
	"github.com/apex/log"
	"github.com/kr/pretty"
	"github.com/lrstanley/spectrograph/internal/models"
)

// TODO: sub-package.

type BotPermissions struct {
	member      *disgord.Member
	permissions models.DiscordPermissions
}

// NewPermissionsCache returns a new permissions cache.
func NewPermissionsCache() *permissionsCache {
	return &permissionsCache{
		db: make(map[disgord.Snowflake]*BotPermissions),
	}
}

type permissionsCache struct {
	mu sync.RWMutex
	db map[disgord.Snowflake]*BotPermissions
}

func (pc *permissionsCache) guildCreate(sess disgord.Session, event *disgord.GuildCreate) (models.DiscordPermissions, error) {
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

	permissions := mergeRolePermissions(botMember, event.Guild.Roles)

	pc.mu.Lock()
	defer pc.mu.Unlock()

	pc.db[event.Guild.ID] = &BotPermissions{
		member:      botMember,
		permissions: permissions,
	}

	return permissions, nil
}

// TODO: this doesn't seem t otrigger when we add or remove roles from ourselves.
func (pc *permissionsCache) guildMemberUpdate(sess disgord.Session, event *disgord.GuildMemberUpdate) {
	guild, err := sess.Guild(event.GuildID).Get()
	if err != nil {
		logGuild(logger, event.GuildID).WithError(err).Error("unable to fetch guild for member update")
		return
	}

	pc.mu.Lock()
	defer pc.mu.Unlock()

	permissions := mergeRolePermissions(event.Member, guild.Roles)

	if previous, ok := pc.db[event.GuildID]; ok {
		if previous.permissions == permissions {
			return
		}

		// Permissions have changed.
		// TODO: should we trigger an event of some kind to log to the db?
		logGuild(logger, guild).WithFields(log.Fields{
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

func (pc *permissionsCache) guildRoleChange(sess disgord.Session, guildID disgord.Snowflake) {
	// Check if it's a role ID that has been given to us.
	pc.mu.Lock()
	defer pc.mu.Unlock()

	perms, ok := pc.db[guildID]
	if !ok {
		return
	}

	roles, err := sess.Guild(guildID).GetRoles()
	if err != nil {
		logGuild(logger, guildID).WithError(err).Error("unable to fetch roles for role-update permissions checks")
		return
	}

	pc.db[guildID].permissions = mergeRolePermissions(perms.member, roles)
}

func (pc *permissionsCache) removeGuild(guildID disgord.Snowflake) {
	if guildID.IsZero() {
		return
	}

	pc.mu.Lock()
	defer pc.mu.Unlock()

	delete(pc.db, guildID)
}

func (pc *permissionsCache) Get(guildID disgord.Snowflake) models.DiscordPermissions {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	if perms, ok := pc.db[guildID]; ok {
		return perms.permissions
	}

	return 0
}

// mergeRolePermissions converts the role id's that it has permission for,
// and generates a permission bit based of all access.
func mergeRolePermissions(member *disgord.Member, roles []*disgord.Role) (perms models.DiscordPermissions) {
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

	pretty.Println(perms)
	return perms
}
