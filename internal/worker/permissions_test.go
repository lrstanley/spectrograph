// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package worker

import (
	"slices"
	"testing"

	"github.com/diamondburned/arikawa/v3/discord"
)

const (
	testGuildID   discord.GuildID   = 1000
	testBotID     discord.UserID    = 2000
	testOwnerID   discord.UserID    = 2001
	testRoleID    discord.RoleID    = 3000
	testCategory  discord.ChannelID = 4000
	everyoneRole                    = discord.RoleID(testGuildID)
	testOtherUser discord.UserID    = 2002
)

func testPermissions(everyone, role discord.Permissions, channels ...discord.Channel) *botPermissions {
	snap := &snapshot{
		guild: discord.Guild{ID: testGuildID, OwnerID: testOwnerID},
		roles: []discord.Role{
			{ID: everyoneRole, Permissions: everyone},
			{ID: testRoleID, Permissions: role},
		},
		channels: channels,
		bot:      &discord.Member{User: discord.User{ID: testBotID}, RoleIDs: []discord.RoleID{testRoleID}},
	}
	return snap.permissions()
}

func TestGuildPermissions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		everyone discord.Permissions
		role     discord.Permissions
		owner    bool
		want     discord.Permissions
	}{
		{
			name:     "everyone-and-role",
			everyone: discord.PermissionViewChannel | discord.PermissionConnect,
			role:     discord.PermissionManageChannels,
			want:     discord.PermissionViewChannel | discord.PermissionConnect | discord.PermissionManageChannels,
		},
		{
			name:     "everyone-only",
			everyone: discord.PermissionViewChannel,
			want:     discord.PermissionViewChannel,
		},
		{
			name: "administrator",
			role: discord.PermissionAdministrator,
			want: discord.PermissionAll,
		},
		{
			name:  "owner",
			owner: true,
			want:  discord.PermissionAll,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			p := testPermissions(tt.everyone, tt.role)
			if tt.owner {
				p.guild.OwnerID = testBotID
			}

			if got := p.guildPermissions(); got != tt.want {
				t.Errorf("guildPermissions() = %s, want %s", permissionString(got), permissionString(tt.want))
			}
		})
	}
}

func TestGuildPermissionsUnknownMember(t *testing.T) {
	t.Parallel()

	p := &botPermissions{guild: discord.Guild{ID: testGuildID}}
	if p.known() {
		t.Error("known() = true, want false")
	}
	if got := p.missing(0, discord.PermissionManageChannels); got != 0 {
		t.Errorf("missing() = %s, want none (unknown permissions aren't checked)", permissionString(got))
	}
}

func TestChannelPermissions(t *testing.T) {
	t.Parallel()

	base := discord.PermissionViewChannel | discord.PermissionConnect | discord.PermissionManageChannels

	tests := []struct {
		name       string
		overwrites []discord.Overwrite
		want       discord.Permissions
	}{
		{
			name: "no-overwrites",
			want: base,
		},
		{
			name: "everyone-deny",
			overwrites: []discord.Overwrite{
				{ID: discord.Snowflake(everyoneRole), Type: discord.OverwriteRole, Deny: discord.PermissionConnect},
			},
			want: base &^ discord.PermissionConnect,
		},
		{
			name: "everyone-deny-role-allow",
			overwrites: []discord.Overwrite{
				{ID: discord.Snowflake(everyoneRole), Type: discord.OverwriteRole, Deny: discord.PermissionConnect},
				{ID: discord.Snowflake(testRoleID), Type: discord.OverwriteRole, Allow: discord.PermissionConnect},
			},
			want: base,
		},
		{
			name: "member-deny",
			overwrites: []discord.Overwrite{
				{ID: discord.Snowflake(testBotID), Type: discord.OverwriteMember, Deny: discord.PermissionManageChannels},
			},
			want: base &^ discord.PermissionManageChannels,
		},
		{
			name: "other-member-deny",
			overwrites: []discord.Overwrite{
				{ID: discord.Snowflake(testOtherUser), Type: discord.OverwriteMember, Deny: discord.PermissionManageChannels},
			},
			want: base,
		},
		{
			name: "view-denied-implies-nothing",
			overwrites: []discord.Overwrite{
				{ID: discord.Snowflake(everyoneRole), Type: discord.OverwriteRole, Deny: discord.PermissionViewChannel},
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ch := discord.Channel{ID: 5000, GuildID: testGuildID, Type: discord.GuildVoice, Overwrites: tt.overwrites}
			p := testPermissions(discord.PermissionViewChannel|discord.PermissionConnect, discord.PermissionManageChannels, ch)

			if got := p.channelPermissions(ch.ID); got != tt.want {
				t.Errorf("channelPermissions() = %q, want %q", permissionString(got), permissionString(tt.want))
			}
		})
	}
}

func TestUnsettable(t *testing.T) {
	t.Parallel()

	category := discord.Channel{
		ID:      testCategory,
		GuildID: testGuildID,
		Type:    discord.GuildCategory,
		Overwrites: []discord.Overwrite{
			{ID: discord.Snowflake(testRoleID), Type: discord.OverwriteRole, Allow: discord.PermissionSpeak},
		},
	}

	// Permissions referenced by overwrites in these tests.
	const referenced = discord.PermissionConnect | discord.PermissionSpeak | discord.PermissionStream |
		discord.PermissionManageRoles | discord.PermissionManageChannels | discord.PermissionViewChannel

	tests := []struct {
		name   string
		role   discord.Permissions
		parent discord.ChannelID
		want   discord.Permissions // Within referenced.
	}{
		{
			name: "only-own-permissions",
			role: discord.PermissionConnect | discord.PermissionManageChannels,
			want: discord.PermissionSpeak | discord.PermissionStream | discord.PermissionManageRoles,
		},
		{
			name:   "granted-by-category",
			role:   discord.PermissionConnect | discord.PermissionManageChannels,
			parent: testCategory,
			want:   discord.PermissionStream | discord.PermissionManageRoles,
		},
		{
			name: "manage-roles-requires-admin",
			role: discord.PermissionManageRoles | discord.PermissionConnect | discord.PermissionSpeak |
				discord.PermissionStream | discord.PermissionManageChannels,
			want: discord.PermissionManageRoles,
		},
		{
			name: "admin-bypass",
			role: discord.PermissionAdministrator,
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			p := testPermissions(discord.PermissionViewChannel, tt.role, category)
			if got := p.unsettable(tt.parent) & referenced; got != tt.want {
				t.Errorf("unsettable() = %q, want %q", permissionString(got), permissionString(tt.want))
			}
		})
	}
}

func TestUnsettableUnknownMember(t *testing.T) {
	t.Parallel()

	p := &botPermissions{guild: discord.Guild{ID: testGuildID}}
	if got := p.unsettable(0); got != 0 {
		t.Errorf("unsettable() = %q, want none (unknown permissions aren't checked)", permissionString(got))
	}
}

func TestCopyOverwrites(t *testing.T) {
	t.Parallel()

	const unsettable = discord.PermissionManageRoles | discord.PermissionStream

	everyone := discord.Overwrite{ID: discord.Snowflake(everyoneRole), Type: discord.OverwriteRole, Deny: discord.PermissionConnect}
	botRole := discord.Overwrite{
		ID:    discord.Snowflake(testRoleID),
		Type:  discord.OverwriteRole,
		Allow: discord.PermissionViewChannel | discord.PermissionManageChannels | discord.PermissionManageRoles,
	}
	streamOnly := discord.Overwrite{ID: 5, Type: discord.OverwriteMember, Allow: discord.PermissionStream}
	denyStream := discord.Overwrite{ID: 6, Type: discord.OverwriteRole, Allow: discord.PermissionConnect, Deny: discord.PermissionStream}

	t.Run("allow-stripped", func(t *testing.T) {
		t.Parallel()

		c := copyOverwrites([]discord.Overwrite{everyone, botRole, streamOnly}, unsettable)

		want := []discord.Overwrite{everyone, {
			ID:    botRole.ID,
			Type:  botRole.Type,
			Allow: discord.PermissionViewChannel | discord.PermissionManageChannels,
		}} // streamOnly is left empty, so it's dropped.
		if !slices.Equal(c.overwrites, want) {
			t.Errorf("overwrites = %+v, want %+v", c.overwrites, want)
		}

		wantStripped := []discord.Overwrite{
			{ID: botRole.ID, Type: botRole.Type, Allow: discord.PermissionManageRoles},
			{ID: streamOnly.ID, Type: streamOnly.Type, Allow: discord.PermissionStream},
		}
		if !slices.Equal(c.stripped, wantStripped) {
			t.Errorf("stripped = %+v, want %+v", c.stripped, wantStripped)
		}
		if len(c.blocked) > 0 {
			t.Errorf("blocked = %+v, want none", c.blocked)
		}
		if !c.matches(c.overwrites) {
			t.Error("copies made with the copied overwrites don't match them")
		}
		if !c.matches([]discord.Overwrite{everyone, botRole, streamOnly}) {
			t.Error("copies which also allow the stripped permissions don't match")
		}
		// e.g. Video was revoked for @everyone on the original channel, but the
		// bot can't revoke it on copies.
		revoked := slices.Clone(c.overwrites)
		revoked[0].Allow |= discord.PermissionStream
		if c.matches(revoked) {
			t.Error("copies allowing unsettable permissions the original doesn't allow match")
		}
	})

	t.Run("deny-blocked", func(t *testing.T) {
		t.Parallel()

		c := copyOverwrites([]discord.Overwrite{everyone, denyStream}, unsettable)
		want := []discord.Overwrite{{ID: denyStream.ID, Type: denyStream.Type, Deny: discord.PermissionStream}}
		if !slices.Equal(c.blocked, want) {
			t.Errorf("blocked = %+v, want %+v", c.blocked, want)
		}
		if len(c.stripped) > 0 {
			t.Errorf("stripped = %+v, want none", c.stripped)
		}
	})

	t.Run("nothing-unsettable", func(t *testing.T) {
		t.Parallel()

		src := []discord.Overwrite{everyone, botRole, denyStream}
		c := copyOverwrites(src, 0)
		if !slices.Equal(c.overwrites, src) || len(c.stripped) > 0 || len(c.blocked) > 0 {
			t.Errorf("copyOverwrites() = %+v, want the overwrites unchanged", c)
		}
	})
}

func TestPermissionString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   discord.Permissions
		want string
	}{
		{0, ""},
		{discord.PermissionManageChannels, "Manage Channels"},
		{discord.PermissionManageChannels | discord.PermissionConnect, "Manage Channels, Connect"},
		{discord.PermissionConnect | discord.PermissionKickMembers, "Connect, 0x2"},
	}

	for _, tt := range tests {
		if got := permissionString(tt.in); got != tt.want {
			t.Errorf("permissionString(%d) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
