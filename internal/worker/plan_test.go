// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

//nolint:exhaustive // map literals keyed by channel ID (which has a NullChannelID constant) aren't enums.
package worker

import (
	"fmt"
	"log/slog"
	"maps"
	"regexp"
	"slices"
	"testing"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/lrstanley/spectrograph/internal/models"
)

func voice(id discord.ChannelID, name string, pos int) discord.Channel {
	return discord.Channel{
		ID:             id,
		GuildID:        testGuildID,
		Type:           discord.GuildVoice,
		Name:           name,
		Position:       pos,
		VoiceBitrate:   64000,
		VoiceUserLimit: 5,
	}
}

func inCategory(ch discord.Channel, parent discord.ChannelID) discord.Channel {
	ch.ParentID = parent
	return ch
}

// usersIn returns voice states for the given number of users per channel.
func usersIn(counts map[discord.ChannelID]int) []discord.VoiceState {
	var states []discord.VoiceState
	uid := discord.UserID(100)
	for _, id := range slices.Sorted(maps.Keys(counts)) {
		for range counts[id] {
			uid++
			states = append(states, discord.VoiceState{GuildID: testGuildID, ChannelID: id, UserID: uid})
		}
	}
	return states
}

func defaultSettings() guildSettings {
	return guildSettings{
		enabled:   true,
		match:     models.DefaultChannelMatch,
		maxGroups: models.DefaultMaxChannels,
		maxClones: models.DefaultMaxClones,
	}
}

// testSnapshot returns a snapshot where the bot has the given guild-level
// permissions.
func testSnapshot(perms discord.Permissions, channels []discord.Channel, users map[discord.ChannelID]int) *snapshot {
	return &snapshot{
		guild:       discord.Guild{ID: testGuildID, OwnerID: testOwnerID},
		roles:       []discord.Role{{ID: everyoneRole}, {ID: testRoleID, Permissions: perms}},
		channels:    channels,
		voiceStates: usersIn(users),
		bot:         &discord.Member{User: discord.User{ID: testBotID}, RoleIDs: []discord.RoleID{testRoleID}},
	}
}

const botPerms = discord.PermissionViewChannel | discord.PermissionManageChannels | discord.PermissionConnect

func channelIDs(chs []discord.Channel) []discord.ChannelID {
	ids := make([]discord.ChannelID, 0, len(chs))
	for _, ch := range chs {
		ids = append(ids, ch.ID)
	}
	return ids
}

// summary is a comparable summary of a plan.
type summary struct {
	deletes  []discord.ChannelID
	creates  []string // "<source id>@<position>"
	moves    []string // "<id>@<position>"
	syncs    []string // "<id><-<source id>[+overwrites]"
	problems []string
}

func summarize(p *plan) summary {
	s := summary{deletes: channelIDs(p.deletes)}
	for _, c := range p.creates {
		s.creates = append(s.creates, fmt.Sprintf("%d@%d", c.source.ID, c.position))
	}
	for _, m := range p.moves {
		s.moves = append(s.moves, fmt.Sprintf("%d@%d", m.channel.ID, m.position))
	}
	for i := range p.syncs {
		sy := &p.syncs[i]
		v := fmt.Sprintf("%d<-%d", sy.channel.ID, sy.source.ID)
		if sy.overwrites {
			v += "+overwrites"
		}
		s.syncs = append(s.syncs, v)
	}
	for _, pr := range p.problems {
		s.problems = append(s.problems, pr.msg)
	}
	return s
}

func (s summary) String() string {
	return fmt.Sprintf("deletes=%v creates=%v moves=%v syncs=%v problems=%q", s.deletes, s.creates, s.moves, s.syncs, s.problems)
}

func TestOccupancy(t *testing.T) {
	t.Parallel()

	got := occupancy([]discord.VoiceState{
		{UserID: 1, ChannelID: 10},
		{UserID: 2, ChannelID: 10},
		{UserID: 3, ChannelID: 11},
		{UserID: 4}, // Not connected.
	})

	want := map[discord.ChannelID]int{10: 2, 11: 1}
	if !maps.Equal(got, want) {
		t.Errorf("occupancy() = %v, want %v", got, want)
	}
}

func TestGroupChannels(t *testing.T) {
	t.Parallel()

	channels := []discord.Channel{
		voice(3, "Gaming +", 3),
		voice(1, "Gaming +", 1),
		voice(2, "Chill +", 2),
		voice(8, "Alpha +", 0),
		inCategory(voice(4, "Gaming +", 0), testCategory),
		voice(5, "Not managed", 5),
		{ID: 6, Type: discord.GuildText, Name: "text +"},
		{ID: 7, Type: discord.GuildStageVoice, Name: "stage +"},
	}

	tests := []struct {
		name  string
		match *regexp.Regexp
		want  map[discord.ChannelID][][]discord.ChannelID
	}{
		{
			name:  "default-regex",
			match: models.DefaultChannelMatch,
			want: map[discord.ChannelID][][]discord.ChannelID{
				0:            {{1, 3}, {2}, {8}},
				testCategory: {{4}},
			},
		},
		{
			name:  "custom-regex",
			match: regexp.MustCompile(`^Gaming`),
			want: map[discord.ChannelID][][]discord.ChannelID{
				0:            {{1, 3}},
				testCategory: {{4}},
			},
		},
		{
			name:  "unanchored-regex",
			match: regexp.MustCompile(`managed`),
			want: map[discord.ChannelID][][]discord.ChannelID{
				0: {{5}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			parents, groups := groupChannels(channels, tt.match)
			if !slices.IsSorted(parents) || len(parents) != len(tt.want) {
				t.Fatalf("parents = %v, want sorted keys of %v", parents, tt.want)
			}

			for parent, want := range tt.want {
				var got [][]discord.ChannelID
				for _, g := range groups[parent] {
					got = append(got, channelIDs(g))
				}
				if fmt.Sprint(got) != fmt.Sprint(want) {
					t.Errorf("groups[%d] = %v, want %v", parent, got, want)
				}
			}
		})
	}
}

const strippedMsg = "clones leave out permissions allowed by the original channel's overwrites, as the bot can't grant them"

// botRoleOverwrites allow Manage Channels, View Channel, and Manage
// Permissions to the bot's own role.
var botRoleOverwrites = []discord.Overwrite{{
	ID:    discord.Snowflake(testRoleID),
	Type:  discord.OverwriteRole,
	Allow: discord.PermissionViewChannel | discord.PermissionManageChannels | discord.PermissionManageRoles,
}}

func TestPlanGuild(t *testing.T) {
	t.Parallel()

	denyConnect := []discord.Overwrite{
		{ID: discord.Snowflake(everyoneRole), Type: discord.OverwriteRole, Deny: discord.PermissionConnect},
		{ID: discord.Snowflake(testRoleID), Type: discord.OverwriteRole, Allow: discord.PermissionConnect | discord.PermissionStream},
	}

	withOverwrites := func(ch discord.Channel, ows []discord.Overwrite) discord.Channel {
		ch.Overwrites = ows
		return ch
	}

	tests := []struct {
		name     string
		perms    discord.Permissions
		afk      discord.ChannelID
		settings func(*guildSettings)
		channels []discord.Channel
		users    map[discord.ChannelID]int
		want     summary
	}{
		{
			name:     "single-empty-channel",
			channels: []discord.Channel{voice(1, "Gaming +", 1)},
		},
		{
			name:     "unmatched-channel-occupied",
			channels: []discord.Channel{voice(1, "Gaming", 1)},
			users:    map[discord.ChannelID]int{1: 3},
		},
		{
			name:     "single-occupied-channel-clones",
			channels: []discord.Channel{voice(1, "Gaming +", 1), voice(2, "Other", 2)},
			users:    map[discord.ChannelID]int{1: 1},
			want:     summary{creates: []string{"1@2"}},
		},
		{
			name:     "all-occupied-clones-from-original-after-last",
			channels: []discord.Channel{voice(5, "Gaming +", 3), voice(1, "Gaming +", 1), voice(3, "Gaming +", 2)},
			users:    map[discord.ChannelID]int{1: 1, 3: 2, 5: 1},
			want:     summary{creates: []string{"1@4"}},
		},
		{
			name:     "one-empty-after-occupied-is-stable",
			channels: []discord.Channel{voice(1, "Gaming +", 1), voice(2, "Gaming +", 2)},
			users:    map[discord.ChannelID]int{1: 1},
		},
		{
			name:     "extra-empty-channels-removed-keeping-oldest",
			channels: []discord.Channel{voice(1, "Gaming +", 1), voice(2, "Gaming +", 2), voice(3, "Gaming +", 3), voice(4, "Gaming +", 4)},
			users:    map[discord.ChannelID]int{2: 1},
			want:     summary{deletes: []discord.ChannelID{3, 4}, moves: []string{"1@3"}},
		},
		{
			name:     "all-empty-keeps-original",
			channels: []discord.Channel{voice(2, "Gaming +", 1), voice(1, "Gaming +", 2)},
			want:     summary{deletes: []discord.ChannelID{2}},
		},
		{
			name:     "empty-original-moved-after-occupied-clone",
			channels: []discord.Channel{voice(1, "Gaming +", 1), voice(2, "Gaming +", 2)},
			users:    map[discord.ChannelID]int{2: 1},
			want:     summary{moves: []string{"1@3"}},
		},
		{
			name:     "same-position-ordered-by-id",
			channels: []discord.Channel{voice(1, "Gaming +", 1), voice(2, "Gaming +", 1)},
			users:    map[discord.ChannelID]int{2: 1},
			want:     summary{moves: []string{"1@2"}},
		},
		{
			name: "groups-are-independent",
			channels: []discord.Channel{
				voice(1, "Gaming +", 1),
				voice(2, "Chill +", 2),
				inCategory(voice(3, "Gaming +", 1), testCategory),
			},
			users: map[discord.ChannelID]int{1: 1, 3: 1},
			want:  summary{creates: []string{"1@2", "3@2"}},
		},
		{
			name:     "clone-limit-reached",
			settings: func(s *guildSettings) { s.maxClones = 1 },
			channels: []discord.Channel{voice(1, "Gaming +", 1), voice(2, "Gaming +", 2)},
			users:    map[discord.ChannelID]int{1: 1, 2: 1},
			want:     summary{problems: []string{"all channels are occupied, but the clone limit has been reached"}},
		},
		{
			name:     "over-clone-limit-still-cleans-up",
			settings: func(s *guildSettings) { s.maxClones = 1 },
			channels: []discord.Channel{voice(1, "Gaming +", 1), voice(2, "Gaming +", 2), voice(3, "Gaming +", 3), voice(4, "Gaming +", 4)},
			users:    map[discord.ChannelID]int{1: 1},
			want:     summary{deletes: []discord.ChannelID{3, 4}},
		},
		{
			name:     "group-limit-per-category",
			settings: func(s *guildSettings) { s.maxGroups = 1 },
			channels: []discord.Channel{voice(1, "A +", 1), voice(2, "B +", 2), inCategory(voice(3, "C +", 1), testCategory)},
			users:    map[discord.ChannelID]int{1: 1, 2: 1, 3: 1},
			want: summary{
				creates:  []string{"1@2", "3@2"},
				problems: []string{"too many matching channels in category, only managing the first ones"},
			},
		},
		{
			// The AFK channel (1) is left alone, other matching channels with
			// the same name are managed without it.
			name:     "afk-channel-not-managed",
			afk:      1,
			channels: []discord.Channel{voice(1, "AFK +", 1), voice(2, "AFK +", 2)},
			users:    map[discord.ChannelID]int{1: 1, 2: 1},
			want:     summary{creates: []string{"2@3"}},
		},
		{
			// The managed groups are the oldest ones, so which groups are
			// managed doesn't change as channels are reordered (here, group A's
			// original was moved after its clone, which is after group B).
			name:     "group-limit-keeps-oldest-groups",
			settings: func(s *guildSettings) { s.maxGroups = 1 },
			channels: []discord.Channel{voice(1, "A +", 3), voice(10, "A +", 2), voice(2, "B +", 2)},
			users:    map[discord.ChannelID]int{1: 1, 10: 1},
			want: summary{
				creates:  []string{"1@4"},
				problems: []string{"too many matching channels in category, only managing the first ones"},
			},
		},
		{
			name:  "clone-settings-synced",
			perms: botPerms | discord.PermissionManageRoles,
			channels: []discord.Channel{
				voice(1, "Gaming +", 1),
				func() discord.Channel { ch := voice(2, "Gaming +", 2); ch.VoiceUserLimit = 10; return ch }(),
				withOverwrites(voice(3, "Gaming +", 3), denyConnect),
			},
			users: map[discord.ChannelID]int{1: 1, 2: 1},
			want:  summary{syncs: []string{"2<-1", "3<-1+overwrites"}},
		},
		{
			name:  "overwrite-order-ignored",
			perms: botPerms | discord.PermissionStream,
			channels: []discord.Channel{
				withOverwrites(voice(1, "Gaming +", 1), denyConnect),
				withOverwrites(voice(2, "Gaming +", 2), []discord.Overwrite{denyConnect[1], denyConnect[0]}),
			},
			users: map[discord.ChannelID]int{1: 1},
		},
		{
			name:     "deleted-channels-not-synced",
			channels: []discord.Channel{voice(1, "Gaming +", 1), withOverwrites(voice(2, "Gaming +", 2), denyConnect), voice(3, "Gaming +", 3)},
			want:     summary{deletes: []discord.ChannelID{2, 3}},
		},
		{
			name:     "missing-manage-channels",
			perms:    discord.PermissionViewChannel,
			channels: []discord.Channel{voice(1, "Gaming +", 1), voice(2, "Gaming +", 2), voice(3, "Gaming +", 3)},
			users:    map[discord.ChannelID]int{2: 1},
			want: summary{problems: []string{
				"missing permissions to remove unneeded empty channel",
				"missing permissions to reorder empty channel",
			}},
		},
		{
			name:     "missing-manage-channels-clone",
			perms:    discord.PermissionViewChannel,
			channels: []discord.Channel{voice(1, "Gaming +", 1)},
			users:    map[discord.ChannelID]int{1: 1},
			want:     summary{problems: []string{"missing permissions to clone channel"}},
		},
		{
			// The Video allow can't be copied (the bot doesn't have it), so it is
			// left out, and the clone is still made.
			name:     "unsettable-allow-stripped-clone",
			channels: []discord.Channel{withOverwrites(voice(1, "Gaming +", 1), denyConnect)},
			users:    map[discord.ChannelID]int{1: 1},
			want: summary{
				creates:  []string{"1@2"},
				problems: []string{strippedMsg},
			},
		},
		{
			// Leaving out a deny would grant access, so the channel isn't cloned.
			name: "unsettable-deny-refuses-clone",
			channels: []discord.Channel{withOverwrites(voice(1, "Gaming +", 1), []discord.Overwrite{
				{ID: discord.Snowflake(everyoneRole), Type: discord.OverwriteRole, Deny: discord.PermissionStream},
			})},
			users: map[discord.ChannelID]int{1: 1},
			want:  summary{problems: []string{"unable to clone channel, as its overwrites deny permissions the bot can't deny"}},
		},
		{
			name:     "overwrite-permissions-granted-clone",
			perms:    botPerms | discord.PermissionStream,
			channels: []discord.Channel{withOverwrites(voice(1, "Gaming +", 1), denyConnect)},
			users:    map[discord.ChannelID]int{1: 1},
			want:     summary{creates: []string{"1@2"}},
		},
		{
			// The original channel allows Manage Permissions to the bot's own
			// role, which only administrators can set in an overwrite (even
			// with Manage Roles in the server).
			name:     "bot-role-manage-permissions-stripped-clone",
			perms:    botPerms | discord.PermissionManageRoles,
			channels: []discord.Channel{withOverwrites(voice(1, "Gaming +", 1), botRoleOverwrites)},
			users:    map[discord.ChannelID]int{1: 1},
			want: summary{
				creates:  []string{"1@2"},
				problems: []string{strippedMsg},
			},
		},
		{
			// A clone made without the stripped permissions isn't out of sync.
			name:  "stripped-clone-not-resynced",
			perms: botPerms | discord.PermissionManageRoles,
			channels: []discord.Channel{
				withOverwrites(voice(1, "Gaming +", 1), botRoleOverwrites),
				withOverwrites(voice(2, "Gaming +", 2), []discord.Overwrite{{
					ID:    discord.Snowflake(testRoleID),
					Type:  discord.OverwriteRole,
					Allow: discord.PermissionViewChannel | discord.PermissionManageChannels,
				}}),
			},
			users: map[discord.ChannelID]int{1: 1},
			want:  summary{problems: []string{strippedMsg}},
		},
		{
			// Syncing a clone leaves out the same permissions.
			name:  "stripped-clone-synced-without-unsettable",
			perms: botPerms | discord.PermissionManageRoles,
			channels: []discord.Channel{
				withOverwrites(voice(1, "Gaming +", 1), botRoleOverwrites),
				voice(2, "Gaming +", 2),
			},
			users: map[discord.ChannelID]int{1: 1},
			want: summary{
				syncs:    []string{"2<-1+overwrites"},
				problems: []string{strippedMsg},
			},
		},
		{
			name:  "unsettable-deny-refuses-overwrite-sync",
			perms: botPerms | discord.PermissionManageRoles,
			channels: []discord.Channel{
				withOverwrites(voice(1, "Gaming +", 1), []discord.Overwrite{
					{ID: discord.Snowflake(everyoneRole), Type: discord.OverwriteRole, Deny: discord.PermissionStream},
				}),
				voice(2, "Gaming +", 2),
			},
			users: map[discord.ChannelID]int{1: 1},
			want: summary{problems: []string{
				"unable to sync clone permission overwrites, as the original channel's overwrites deny permissions the bot can't deny",
			}},
		},
		{
			name: "missing-manage-roles-syncs-settings-only",
			channels: []discord.Channel{
				voice(1, "Gaming +", 1),
				func() discord.Channel {
					ch := withOverwrites(voice(2, "Gaming +", 2), denyConnect[:1])
					ch.VoiceBitrate = 96000
					return ch
				}(),
			},
			users: map[discord.ChannelID]int{1: 1},
			want: summary{
				syncs:    []string{"2<-1"},
				problems: []string{"missing permissions to sync clone permission overwrites with original channel"},
			},
		},
		{
			name:  "manage-roles-syncs-overwrites",
			perms: botPerms | discord.PermissionManageRoles,
			channels: []discord.Channel{
				voice(1, "Gaming +", 1),
				withOverwrites(voice(2, "Gaming +", 2), denyConnect[:1]),
			},
			users: map[discord.ChannelID]int{1: 1},
			want:  summary{syncs: []string{"2<-1+overwrites"}},
		},
		{
			name:     "administrator-bypasses-checks",
			perms:    discord.PermissionAdministrator,
			channels: []discord.Channel{withOverwrites(voice(1, "Gaming +", 1), denyConnect)},
			users:    map[discord.ChannelID]int{1: 1},
			want:     summary{creates: []string{"1@2"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			settings := defaultSettings()
			if tt.settings != nil {
				tt.settings(&settings)
			}

			perms := tt.perms
			if perms == 0 {
				perms = botPerms
			}

			snap := testSnapshot(perms, tt.channels, tt.users)
			snap.guild.AFKChannelID = tt.afk
			got := summarize(planGuild(snap, settings))
			if got.String() != tt.want.String() {
				t.Errorf("planGuild():\n got: %s\nwant: %s", got, tt.want)
			}
		})
	}
}

// TestPlanConvergence applies plans to a fake guild through a series of user
// activity, and checks that a single pass converges (the next plan is empty),
// and that each group ends up with exactly one empty channel, after all
// occupied channels.
func TestPlanConvergence(t *testing.T) {
	t.Parallel()

	fake := newFakeDiscord(voice(1, "Gaming +", 1), voice(50, "Other", 2), inCategory(voice(60, "Chill +", 1), testCategory))
	settings := defaultSettings()

	steps := []map[discord.ChannelID]int{
		{},
		{1: 1},
		{1: 1, 9001: 2},
		{1: 1, 9001: 2, 9002: 1, 60: 1},
		{9001: 2, 9002: 1},
		{9002: 1},
		{1: 1},
		{},
	}

	for i, users := range steps {
		snap := testSnapshot(botPerms, fake.channels, users)
		p := planGuild(snap, settings)
		if problems := apply(fake, &fakeStore{}, testGuildID, p); len(problems) > 0 {
			t.Fatalf("step %d: apply() problems: %v", i, problems)
		}

		snap = testSnapshot(botPerms, fake.channels, users)
		if next := planGuild(snap, settings); !next.empty() {
			t.Fatalf("step %d: plan did not converge, next plan: %s", i, summarize(next))
		}

		_, groups := groupChannels(fake.channels, settings.match)
		for _, gs := range groups {
			for _, g := range gs {
				var empty int
				for j, ch := range g {
					if users[ch.ID] > 0 {
						continue
					}
					empty++
					if j != len(g)-1 {
						t.Errorf("step %d: empty channel %d is not last in group %v", i, ch.ID, channelIDs(g))
					}
				}
				if empty != 1 {
					t.Errorf("step %d: group %v has %d empty channels, want 1", i, channelIDs(g), empty)
				}
			}
		}
	}

	// The original channels should never be deleted.
	for _, id := range []discord.ChannelID{1, 50, 60} {
		if fake.index(id) < 0 {
			t.Errorf("channel %d was deleted", id)
		}
	}
}

func TestOverwritesEqual(t *testing.T) {
	t.Parallel()

	a := discord.Overwrite{ID: 1, Type: discord.OverwriteRole, Allow: discord.PermissionConnect}
	b := discord.Overwrite{ID: 2, Type: discord.OverwriteMember, Deny: discord.PermissionSpeak}

	aStream := a
	aStream.Allow |= discord.PermissionStream
	empty := discord.Overwrite{ID: 3, Type: discord.OverwriteRole}

	tests := []struct {
		name string
		x, y []discord.Overwrite
		want bool
	}{
		{"both-empty", nil, []discord.Overwrite{}, true},
		{"same-order", []discord.Overwrite{a, b}, []discord.Overwrite{a, b}, true},
		{"different-order", []discord.Overwrite{a, b}, []discord.Overwrite{b, a}, true},
		{"different-length", []discord.Overwrite{a}, []discord.Overwrite{a, b}, false},
		{"different-allow", []discord.Overwrite{a}, []discord.Overwrite{{ID: 1, Type: discord.OverwriteRole}}, false},
		{"different-type", []discord.Overwrite{a}, []discord.Overwrite{{ID: 1, Type: discord.OverwriteMember, Allow: discord.PermissionConnect}}, false},
		{"empty-overwrite-ignored", []discord.Overwrite{a, empty}, []discord.Overwrite{a}, true},
		{"extra-allow", []discord.Overwrite{aStream}, []discord.Overwrite{a}, false},
	}

	for _, tt := range tests {
		if got := overwritesEqual(tt.x, tt.y); got != tt.want {
			t.Errorf("%s: overwritesEqual() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

// TestPlanStrippedOverwritesConverge checks that a clone made without the
// permissions the bot can't set is created with the rest of the overwrites,
// and that the next plan is empty (no re-sync loop), even when the bot could
// sync overwrites.
func TestPlanStrippedOverwritesConverge(t *testing.T) {
	t.Parallel()

	everyoneDeny := discord.Overwrite{ID: discord.Snowflake(everyoneRole), Type: discord.OverwriteRole, Deny: discord.PermissionConnect}
	primary := voice(1, "Gaming +", 1)
	primary.Overwrites = append([]discord.Overwrite{everyoneDeny}, botRoleOverwrites...)

	fake := newFakeDiscord(primary)
	settings := defaultSettings()
	users := map[discord.ChannelID]int{1: 1}
	perms := botPerms | discord.PermissionManageRoles

	p := planGuild(testSnapshot(perms, fake.channels, users), settings)
	if len(p.creates) != 1 {
		t.Fatalf("plan = %s, want a single create", summarize(p))
	}

	want := []discord.Overwrite{everyoneDeny, {
		ID:    discord.Snowflake(testRoleID),
		Type:  discord.OverwriteRole,
		Allow: discord.PermissionViewChannel | discord.PermissionManageChannels,
	}}
	if len(p.problems) != 1 || p.problems[0].level != slog.LevelInfo ||
		p.problems[0].fields[fieldOverwrites] != "role 3000: Manage Permissions" {
		t.Errorf("problems = %+v, want an info problem naming the stripped permission", p.problems)
	}
	if got := p.creates[0].overwrites; !slices.Equal(got, want) {
		t.Errorf("create overwrites = %+v, want %+v", got, want)
	}

	if problems := apply(fake, &fakeStore{}, testGuildID, p); len(problems) > 0 {
		t.Fatalf("apply() problems: %v", problems)
	}

	for i := range 3 {
		next := planGuild(testSnapshot(perms, fake.channels, users), settings)
		if !next.empty() {
			t.Fatalf("pass %d: plan did not converge, next plan: %s", i, summarize(next))
		}
		if got := summarize(next).problems; !slices.Equal(got, []string{strippedMsg}) {
			t.Errorf("pass %d: problems = %q, want only the stripped permissions notice", i, got)
		}
	}
}

func TestDescribeOverwrites(t *testing.T) {
	t.Parallel()

	snap := testSnapshot(botPerms, nil, nil)
	snap.roles[1].Name = "Spectrograph"
	snap.bot.User.Username = "spectrograph"
	snap.voiceStates = []discord.VoiceState{{
		UserID: testOtherUser,
		Member: &discord.Member{User: discord.User{ID: testOtherUser, Username: "alice"}},
	}}

	got := snap.describeOverwrites([]discord.Overwrite{
		{ID: discord.Snowflake(everyoneRole), Type: discord.OverwriteRole, Deny: discord.PermissionConnect},
		{ID: discord.Snowflake(testRoleID), Type: discord.OverwriteRole, Allow: discord.PermissionManageRoles},
		{ID: 9999, Type: discord.OverwriteRole, Allow: discord.PermissionStream},
		{ID: discord.Snowflake(testBotID), Type: discord.OverwriteMember, Allow: discord.PermissionSpeak},
		{ID: discord.Snowflake(testOtherUser), Type: discord.OverwriteMember, Deny: discord.PermissionStream},
		{ID: 8888, Type: discord.OverwriteMember, Deny: discord.PermissionStream},
	})

	want := "@everyone: Connect; @Spectrograph: Manage Permissions; role 9999: Video; " +
		"@spectrograph: Speak; @alice: Video; member 8888: Video"
	if got != want {
		t.Errorf("describeOverwrites() =\n %q, want\n %q", got, want)
	}
}
