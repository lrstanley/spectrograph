// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package worker

import (
	"math/bits"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
)

// permissionNames are human-readable names for the permissions the worker
// cares about, used when reporting missing permissions to guild admins.
var permissionNames = []struct {
	perm discord.Permissions
	name string
}{
	{discord.PermissionCreateInstantInvite, "Create Invite"},
	{discord.PermissionAdministrator, "Administrator"},
	{discord.PermissionManageChannels, "Manage Channels"},
	{discord.PermissionPrioritySpeaker, "Priority Speaker"},
	{discord.PermissionStream, "Video"},
	{discord.PermissionViewChannel, "View Channel"},
	{discord.PermissionSendMessages, "Send Messages"},
	{discord.PermissionConnect, "Connect"},
	{discord.PermissionSpeak, "Speak"},
	{discord.PermissionMuteMembers, "Mute Members"},
	{discord.PermissionDeafenMembers, "Deafen Members"},
	{discord.PermissionMoveMembers, "Move Members"},
	{discord.PermissionUseVAD, "Use Voice Activity"},
	// Manage Roles is called Manage Permissions in channel overwrites, the only
	// context the worker reports it in.
	{discord.PermissionManageRoles, "Manage Permissions"},
	{discord.PermissionRequestToSpeak, "Request to Speak"},
	{discord.PermissionStartEmbeddedActivities, "Use Activities"},
	{discord.PermissionUseSoundboard, "Use Soundboard"},
	{discord.PermissionUseExternalSounds, "Use External Sounds"},
}

// permissionString returns a human-readable, comma-separated list of the
// permissions in p. Unknown bits are listed by their bit value.
func permissionString(p discord.Permissions) string {
	names := make([]string, 0, bits.OnesCount64(uint64(p)))
	for _, pn := range permissionNames {
		if p.Has(pn.perm) {
			names = append(names, pn.name)
			p &^= pn.perm
		}
	}
	if p != 0 {
		names = append(names, "0x"+strconv.FormatUint(uint64(p), 16))
	}
	return strings.Join(names, ", ")
}

// botPermissions computes the bot's effective permissions in a guild, using the
// cached guild, roles, and the bot's own member object.
type botPermissions struct {
	guild    discord.Guild
	roles    []discord.Role
	member   *discord.Member
	channels map[discord.ChannelID]discord.Channel
}

// known returns false if the bot's member object is unknown, in which case
// permission checks are skipped and Discord is left to reject what it must.
func (p *botPermissions) known() bool {
	return p.member != nil
}

// guildPermissions returns the bot's guild-level (base) permissions, which
// account for the @everyone role, the bot's roles, guild ownership, and the
// administrator permission (but no channel overwrites).
func (p *botPermissions) guildPermissions() discord.Permissions {
	if p.member == nil {
		return 0
	}
	return discord.CalcOverrides(p.guild, discord.Channel{}, *p.member, p.roles)
}

// isAdmin returns true if the bot bypasses all permission checks in the guild.
func (p *botPermissions) isAdmin() bool {
	if p.member == nil {
		return false
	}
	return p.guild.OwnerID == p.member.User.ID || p.guildPermissions().Has(discord.PermissionAdministrator)
}

// channelPermissions returns the bot's permissions in the given channel. If
// channelID is not valid (e.g. a channel without a category), the guild-level
// permissions are returned.
func (p *botPermissions) channelPermissions(channelID discord.ChannelID) discord.Permissions {
	if p.member == nil {
		return 0
	}

	ch, ok := p.channels[channelID]
	if !channelID.IsValid() || !ok {
		return p.guildPermissions()
	}

	perms := discord.CalcOverrides(p.guild, ch, *p.member, p.roles)

	// Implicit permissions: without View Channel, every other permission in the
	// channel is denied.
	if !perms.Has(discord.PermissionAdministrator) && !perms.Has(discord.PermissionViewChannel) {
		return 0
	}
	return perms
}

// missing returns the permissions in want that the bot lacks in the given
// channel (or the guild, if channelID isn't valid). It returns 0 if the bot's
// permissions are unknown.
func (p *botPermissions) missing(channelID discord.ChannelID, want discord.Permissions) discord.Permissions {
	if !p.known() || p.isAdmin() {
		return 0
	}
	return want &^ p.channelPermissions(channelID)
}

// unsettable returns the permissions the bot can't allow or deny in the
// permission overwrites of a channel created or edited under the given parent.
// Discord only allows setting overwrites for permissions the bot has in the
// guild or parent category, and only administrators may set Manage Roles
// (Manage Permissions) in an overwrite. It returns 0 if the bot's permissions
// are unknown (Discord is left to reject what it must), or if it is an
// administrator.
func (p *botPermissions) unsettable(parentID discord.ChannelID) discord.Permissions {
	if !p.known() || p.isAdmin() {
		return 0
	}

	settable := p.guildPermissions() | p.channelPermissions(parentID)
	return ^(settable &^ discord.PermissionManageRoles)
}
