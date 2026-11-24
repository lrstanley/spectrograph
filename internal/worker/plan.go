// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package worker

import (
	"cmp"
	"fmt"
	"log/slog"
	"maps"
	"regexp"
	"slices"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
)

// Common problem fields.
const (
	fieldChannelID  = "channel_id"
	fieldMissing    = "missing"
	fieldPrimaryID  = "primary_id"
	fieldOverwrites = "overwrites"
	fieldHint       = "hint"
)

// overwriteHint explains which permissions the bot can set in permission
// overwrites.
const overwriteHint = "the bot can only set permissions it has itself (in the server or the category), " +
	"and only administrators can set Manage Permissions"

// guildSettings are the effective settings for a guild, after merging the
// user-defined config, the admin config, and system defaults.
type guildSettings struct {
	enabled bool

	// match selects the voice channels which are managed.
	match *regexp.Regexp

	// maxGroups is the max number of managed channel groups per category.
	maxGroups int

	// maxClones is the max number of clones per group (excluding the original).
	maxClones int
}

// snapshot is a point-in-time view of a guild, as seen by the bot.
type snapshot struct {
	guild    discord.Guild
	channels []discord.Channel
	roles    []discord.Role

	// voiceStates are the voice states of users connected to a voice channel
	// in the guild.
	voiceStates []discord.VoiceState

	// bot is the bot's member object in the guild, or nil if unknown.
	bot *discord.Member
}

func (s *snapshot) permissions() *botPermissions {
	channels := make(map[discord.ChannelID]discord.Channel, len(s.channels))
	for _, ch := range s.channels {
		channels[ch.ID] = ch
	}
	return &botPermissions{guild: s.guild, roles: s.roles, member: s.bot, channels: channels}
}

// occupancy returns the number of connected users per voice channel.
func occupancy(states []discord.VoiceState) map[discord.ChannelID]int {
	count := make(map[discord.ChannelID]int, len(states))
	for _, vs := range states {
		if vs.ChannelID.IsValid() {
			count[vs.ChannelID]++
		}
	}
	return count
}

// problem is a condition preventing the worker from fully managing a guild,
// which is reported to the guild admins (deduplicated across passes).
type problem struct {
	level  slog.Level
	msg    string
	fields map[string]any

	// retry is true if the problem is likely transient (e.g. a Discord server
	// error), so the pass should be retried even without further events.
	retry bool

	// forbidden is true if Discord rejected a change as unauthorized.
	forbidden bool
}

// key uniquely identifies the problem, for deduplication.
func (p problem) key() string {
	return fmt.Sprintf("%d|%s|%v", p.level, p.msg, p.fields) // fmt sorts map keys.
}

type createAction struct {
	source   discord.Channel
	position int

	// overwrites are the permission overwrites of the new channel (the
	// source's, as the bot can copy them, see [copiedOverwrites]).
	overwrites []discord.Overwrite
}

type moveAction struct {
	channel  discord.Channel
	position int
}

type syncAction struct {
	channel discord.Channel
	source  discord.Channel

	// overwrites is true if the permission overwrites should be synced too,
	// to permissionOverwrites (the source's, as the bot can copy them, see
	// [copiedOverwrites]).
	overwrites           bool
	permissionOverwrites []discord.Overwrite
}

// plan is the set of changes needed for a guild to converge on its desired
// state. Applying a plan should result in the next plan being empty (modulo
// user activity), so applying the plan from our own resulting events is a no-op.
type plan struct {
	deletes  []discord.Channel
	creates  []createAction
	moves    []moveAction
	syncs    []syncAction
	problems []problem
}

func (p *plan) empty() bool {
	return len(p.deletes) == 0 && len(p.creates) == 0 && len(p.moves) == 0 && len(p.syncs) == 0
}

func (p *plan) addProblem(level slog.Level, msg string, fields map[string]any) {
	p.problems = append(p.problems, problem{level: level, msg: msg, fields: fields})
}

// channelOrder sorts channels the same way Discord displays them (by position,
// then by ID).
func channelOrder(a, b discord.Channel) int {
	return cmp.Or(cmp.Compare(a.Position, b.Position), cmp.Compare(a.ID, b.ID))
}

// channelGroup is a set of voice channels in the same category with the same
// name (an original channel and its clones), in display order.
type channelGroup []discord.Channel

// primary returns the original channel of the group (the oldest, i.e. the one
// with the lowest ID). Clones copy their settings from it, and it is never
// deleted.
func (g channelGroup) primary() discord.Channel {
	return slices.MinFunc(g, func(a, b discord.Channel) int { return cmp.Compare(a.ID, b.ID) })
}

// groupChannels groups the voice channels matching rgx, by category (parent ID)
// and name. Categories are returned sorted by ID, and the groups of each
// category from the oldest to the newest (by the ID of their original channel),
// so the groups kept when over the per-category limit don't change as channels
// are reordered. The channels of each group are in display order. Channels in
// skip (e.g. the AFK channel) are never managed.
func groupChannels(channels []discord.Channel, rgx *regexp.Regexp, skip ...discord.ChannelID) (parents []discord.ChannelID, groups map[discord.ChannelID][]channelGroup) {
	type key struct {
		parent discord.ChannelID
		name   string
	}

	byKey := map[key]channelGroup{}
	for _, ch := range channels {
		if ch.Type != discord.GuildVoice || !rgx.MatchString(ch.Name) || slices.Contains(skip, ch.ID) {
			continue
		}
		k := key{parent: ch.ParentID, name: ch.Name}
		byKey[k] = append(byKey[k], ch)
	}

	groups = make(map[discord.ChannelID][]channelGroup)
	for k, g := range byKey {
		slices.SortFunc(g, channelOrder)
		groups[k.parent] = append(groups[k.parent], g)
	}

	for _, gs := range groups {
		slices.SortFunc(gs, func(a, b channelGroup) int { return cmp.Compare(a.primary().ID, b.primary().ID) })
	}

	return slices.Sorted(maps.Keys(groups)), groups
}

// planGuild computes the changes needed so that every managed group has
// exactly one empty channel, positioned after the last occupied channel, and
// so that all clones match the settings of their original channel.
func planGuild(snap *snapshot, settings guildSettings) *plan {
	gp := &groupPlanner{
		plan:      &plan{},
		snap:      snap,
		perms:     snap.permissions(),
		users:     occupancy(snap.voiceStates),
		maxClones: settings.maxClones,
	}

	// The AFK channel is never managed: users idling in it aren't using it, and
	// clones of it wouldn't be AFK channels.
	parents, groups := groupChannels(snap.channels, settings.match, snap.guild.AFKChannelID)
	for _, parent := range parents {
		gs := groups[parent]
		if len(gs) > settings.maxGroups {
			gp.plan.addProblem(slog.LevelWarn, "too many matching channels in category, only managing the first ones", map[string]any{
				"category_id": parent.String(),
				"channels":    len(gs),
				"limit":       settings.maxGroups,
			})
			gs = gs[:settings.maxGroups]
		}

		for _, g := range gs {
			gp.planGroup(g)
		}
	}

	return gp.plan
}

// groupPlanner plans the changes for the groups of a guild.
type groupPlanner struct {
	plan      *plan
	snap      *snapshot
	perms     *botPermissions
	users     map[discord.ChannelID]int
	maxClones int
}

func (gp *groupPlanner) planGroup(g channelGroup) {
	p, perms, users := gp.plan, gp.perms, gp.users
	primary := g.primary()

	// Clones are made with the permission overwrites of the original channel,
	// minus the allowed permissions the bot can't set (so clones are at most
	// more restrictive). This is reported for as long as it applies, rather
	// than whenever a clone is made, so it's only logged once.
	copied := copyOverwrites(primary.Overwrites, perms.unsettable(primary.ParentID))
	if len(copied.stripped) > 0 {
		p.addProblem(slog.LevelInfo, "clones leave out permissions allowed by the original channel's overwrites, as the bot can't grant them", map[string]any{
			fieldChannelID:  primary.ID.String(),
			fieldOverwrites: gp.snap.describeOverwrites(copied.stripped),
			fieldHint:       overwriteHint,
		})
	}

	// The empty channel to keep (the oldest, so the original channel is kept if
	// it's empty), and the last occupied channel in display order.
	var keep, lastOccupied *discord.Channel
	for i := range g {
		ch := &g[i]
		if users[ch.ID] > 0 {
			lastOccupied = ch
			continue
		}
		if keep == nil || ch.ID < keep.ID {
			keep = ch
		}
	}

	kept := make([]discord.Channel, 0, len(g))
	for _, ch := range g {
		if keep == nil || users[ch.ID] > 0 || ch.ID == keep.ID {
			kept = append(kept, ch)
			continue
		}

		// Extra empty channel.
		if missing := perms.missing(ch.ID, discord.PermissionManageChannels); missing != 0 {
			p.addProblem(slog.LevelWarn, "missing permissions to remove unneeded empty channel", map[string]any{
				fieldChannelID: ch.ID.String(),
				fieldMissing:   permissionString(missing),
			})
			kept = append(kept, ch)
			continue
		}
		p.deletes = append(p.deletes, ch)
	}

	switch {
	case keep != nil && lastOccupied != nil && channelOrder(*keep, *lastOccupied) < 0:
		// Move the empty channel after the last occupied channel.
		if missing := perms.missing(keep.ID, discord.PermissionManageChannels); missing != 0 {
			p.addProblem(slog.LevelWarn, "missing permissions to reorder empty channel", map[string]any{
				fieldChannelID: keep.ID.String(),
				fieldMissing:   permissionString(missing),
			})
			break
		}
		p.moves = append(p.moves, moveAction{channel: *keep, position: lastOccupied.Position + 1})
	case keep == nil && len(g) > gp.maxClones:
		p.addProblem(slog.LevelWarn, "all channels are occupied, but the clone limit has been reached", map[string]any{
			fieldChannelID: primary.ID.String(),
			"clones":       len(g) - 1,
			"limit":        gp.maxClones,
		})
	case keep == nil:
		if missing := perms.missing(primary.ParentID, discord.PermissionManageChannels); missing != 0 {
			p.addProblem(slog.LevelWarn, "missing permissions to clone channel", map[string]any{
				fieldChannelID: primary.ID.String(),
				fieldMissing:   permissionString(missing),
			})
			break
		}
		if len(copied.blocked) > 0 {
			p.addProblem(slog.LevelWarn, "unable to clone channel, as its overwrites deny permissions the bot can't deny", map[string]any{
				fieldChannelID:  primary.ID.String(),
				fieldOverwrites: gp.snap.describeOverwrites(copied.blocked),
				fieldHint:       overwriteHint,
			})
			break
		}
		p.creates = append(p.creates, createAction{
			source:     primary,
			position:   g[len(g)-1].Position + 1,
			overwrites: copied.overwrites,
		})
	}

	gp.planSyncs(primary, kept, &copied)
}

// planSyncs plans syncing the settings and permission overwrites of the clones
// in kept with their original channel.
func (gp *groupPlanner) planSyncs(primary discord.Channel, kept []discord.Channel, copied *copiedOverwrites) {
	p, perms := gp.plan, gp.perms

	for _, ch := range kept {
		if ch.ID == primary.ID {
			continue
		}

		settingsDiffer := ch.VoiceUserLimit != primary.VoiceUserLimit || ch.VoiceBitrate != primary.VoiceBitrate
		overwritesDiffer := !copied.matches(ch.Overwrites)
		if !settingsDiffer && !overwritesDiffer {
			continue
		}

		if missing := perms.missing(ch.ID, discord.PermissionManageChannels); missing != 0 {
			p.addProblem(slog.LevelWarn, "missing permissions to sync clone with original channel", map[string]any{
				fieldChannelID: ch.ID.String(),
				fieldMissing:   permissionString(missing),
			})
			continue
		}

		if overwritesDiffer {
			if missing := perms.missing(ch.ID, discord.PermissionManageRoles); missing != 0 {
				p.addProblem(slog.LevelWarn, "missing permissions to sync clone permission overwrites with original channel", map[string]any{
					fieldChannelID: ch.ID.String(),
					fieldMissing:   permissionString(missing),
				})
				overwritesDiffer = false
			} else if len(copied.blocked) > 0 {
				p.addProblem(slog.LevelWarn, "unable to sync clone permission overwrites, as the original channel's overwrites deny permissions the bot can't deny", map[string]any{
					fieldChannelID:  ch.ID.String(),
					fieldPrimaryID:  primary.ID.String(),
					fieldOverwrites: gp.snap.describeOverwrites(copied.blocked),
					fieldHint:       overwriteHint,
				})
				overwritesDiffer = false
			}
		}

		if settingsDiffer || overwritesDiffer {
			sync := syncAction{channel: ch, source: primary, overwrites: overwritesDiffer}
			if overwritesDiffer {
				sync.permissionOverwrites = copied.overwrites
			}
			p.syncs = append(p.syncs, sync)
		}
	}
}

// copiedOverwrites are the permission overwrites of a channel, as the bot can
// copy them to another channel (a clone). See [copyOverwrites].
type copiedOverwrites struct {
	// overwrites are the overwrites to set on the copy: the original ones,
	// without the allowed permissions the bot can't set (which only makes the
	// copy more restrictive). Overwrites left empty are dropped.
	overwrites []discord.Overwrite

	// stripped are the overwrites whose allowed permissions were left out of
	// the copy, with only those permissions (in Allow).
	stripped []discord.Overwrite

	// blocked are the overwrites which deny permissions the bot can't set,
	// with only those permissions (in Deny). Overwrites with such denies can't
	// be copied, as leaving them out would grant access that the original
	// channel doesn't.
	blocked []discord.Overwrite
}

// copyOverwrites prepares the overwrites src to be copied to another channel,
// where the bot can't set the permissions in unsettable (see
// [botPermissions.unsettable]).
func copyOverwrites(src []discord.Overwrite, unsettable discord.Permissions) copiedOverwrites {
	var c copiedOverwrites
	if unsettable == 0 {
		c.overwrites = src
		return c
	}

	c.overwrites = make([]discord.Overwrite, 0, len(src))
	for _, ow := range src {
		if deny := ow.Deny & unsettable; deny != 0 {
			c.blocked = append(c.blocked, discord.Overwrite{ID: ow.ID, Type: ow.Type, Deny: deny})
		}
		if allow := ow.Allow & unsettable; allow != 0 {
			c.stripped = append(c.stripped, discord.Overwrite{ID: ow.ID, Type: ow.Type, Allow: allow})
			ow.Allow &^= allow
			if ow.Allow == 0 && ow.Deny == 0 {
				continue
			}
		}
		c.overwrites = append(c.overwrites, ow)
	}
	return c
}

// overwriteTarget returns a human-readable name for the role or member an
// overwrite applies to, using the cached roles and members (falling back to
// the ID).
func (s *snapshot) overwriteTarget(ow discord.Overwrite) string {
	if ow.Type == discord.OverwriteRole {
		if discord.GuildID(ow.ID) == s.guild.ID {
			return "@everyone"
		}
		for _, r := range s.roles {
			if discord.Snowflake(r.ID) == ow.ID && r.Name != "" {
				return "@" + r.Name
			}
		}
		return "role " + ow.ID.String()
	}

	if s.bot != nil && discord.Snowflake(s.bot.User.ID) == ow.ID && s.bot.User.Username != "" {
		return "@" + s.bot.User.Username
	}
	for _, vs := range s.voiceStates {
		if discord.Snowflake(vs.UserID) == ow.ID && vs.Member != nil && vs.Member.User.Username != "" {
			return "@" + vs.Member.User.Username
		}
	}
	return "member " + ow.ID.String()
}

// describeOverwrites returns a human-readable list of the overwrites, and
// their allowed and denied permissions, e.g. "@everyone: Connect; @Mods:
// Manage Permissions".
func (s *snapshot) describeOverwrites(ows []discord.Overwrite) string {
	parts := make([]string, 0, len(ows))
	for _, ow := range ows {
		parts = append(parts, s.overwriteTarget(ow)+": "+permissionString(ow.Allow|ow.Deny))
	}
	return strings.Join(parts, "; ")
}

// matches returns true if ows (the overwrites of a copy) are equivalent to the
// overwrites to set on copies. Permissions stripped from an overwrite are
// ignored when the copy's matching overwrite allows them (e.g. a copy made by
// someone else), as the bot can't remove them, but extra unsettable
// permissions that the original doesn't allow still make the copy differ.
func (c *copiedOverwrites) matches(ows []discord.Overwrite) bool {
	if len(c.stripped) > 0 {
		adjusted := slices.Clone(ows)
		for i := range adjusted {
			for _, st := range c.stripped {
				if adjusted[i].ID == st.ID && adjusted[i].Type == st.Type {
					adjusted[i].Allow &^= st.Allow
				}
			}
		}
		ows = adjusted
	}
	return overwritesEqual(ows, c.overwrites)
}

// overwritesEqual returns true if both sets of permission overwrites are
// equivalent: ignoring order, and overwrites without any effect (nothing
// allowed or denied).
func overwritesEqual(a, b []discord.Overwrite) bool {
	type key struct {
		id  discord.Snowflake
		typ discord.OverwriteType
	}
	type perms struct {
		allow, deny discord.Permissions
	}

	normalize := func(ows []discord.Overwrite) map[key]perms {
		m := make(map[key]perms, len(ows))
		for _, ow := range ows {
			if ow.Allow == 0 && ow.Deny == 0 {
				continue
			}
			m[key{ow.ID, ow.Type}] = perms{allow: ow.Allow, deny: ow.Deny}
		}
		return m
	}

	return maps.Equal(normalize(a), normalize(b))
}
