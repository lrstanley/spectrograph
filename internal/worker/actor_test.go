// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package worker

import (
	"context"
	"testing"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/lrstanley/spectrograph/internal/database/dbtest"
	"github.com/lrstanley/spectrograph/internal/database/ent"
	"github.com/lrstanley/spectrograph/internal/database/ent/guild"
	"github.com/lrstanley/spectrograph/internal/database/ent/guildevent"
	"github.com/lrstanley/spectrograph/internal/database/ent/privacy"
	"github.com/lrstanley/spectrograph/internal/models"
)

// waitFor polls fn until it returns true, or fails the test after timeout.
func waitFor(t *testing.T, timeout time.Duration, what string, fn func() bool) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for !fn() {
		if time.Now().After(deadline) {
			t.Fatalf("timed out waiting for %s", what)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

// TestActorLifecycle drives the worker with simulated gateway events (through
// the real state cache and handlers, without connecting to Discord), and
// checks the guild is synced to, and removed from, the database. The guild's
// channels need no changes, so no Discord API calls are made.
func TestActorLifecycle(t *testing.T) {
	t.Parallel()

	db := dbtest.New(t)
	ctx, cancel := context.WithCancel(privacy.DecisionContext(t.Context(), privacy.Allow))
	defer cancel()

	w := newTestWorker(ctx, t, db, models.ConfigWorker{})
	defer cancel()

	s := w.shards[0]
	dispatch := func(ev any) { s.state.Session.Handler.Call(ev) }

	getGuild := func() *ent.Guild {
		g, _ := db.Guild.Query().Where(guild.GuildID(testGuildID.String())).WithGuildAdminConfig().Only(ctx)
		return g
	}

	dispatch(&gateway.ReadyEvent{User: discord.User{ID: testBotID, Bot: true}})

	create := testGuildCreate(time.Now())
	create.Roles = []discord.Role{{ID: everyoneRole}, {ID: testRoleID, Permissions: botPerms}}
	create.Members = []discord.Member{{User: discord.User{ID: testBotID}, RoleIDs: []discord.RoleID{testRoleID}}}
	create.Channels = []discord.Channel{voice(1, "Gaming +", 1)}
	dispatch(create)

	waitFor(t, 5*time.Second, "guild to be synced", func() bool {
		g := getGuild()
		return g != nil && g.JoinedAt != nil && g.Edges.GuildAdminConfig != nil
	})

	if g := getGuild(); g.Permissions != uint64(botPerms) {
		t.Errorf("permissions = %q, want %q", permissionString(discord.Permissions(g.Permissions)), permissionString(botPerms))
	}

	// Disable the guild, and trigger a pass (a user joins the only channel,
	// which would otherwise be cloned).
	if err := db.GuildAdminConfig.UpdateOne(getGuild().Edges.GuildAdminConfig).SetEnabled(false).Exec(ctx); err != nil {
		t.Fatal(err)
	}
	dispatch(&gateway.VoiceStateUpdateEvent{GuildID: testGuildID, ChannelID: 1, UserID: 500})

	waitFor(t, 5*time.Second, "disabled event to be logged", func() bool {
		n, _ := db.GuildEvent.Query().
			Where(guildevent.HasGuildWith(guild.GuildID(testGuildID.String())), guildevent.MessageContains("guild is disabled")).
			Count(ctx)
		return n == 1
	})

	// Role changes update the stored bot permissions.
	newPerms := botPerms | discord.PermissionManageRoles
	dispatch(&gateway.GuildRoleUpdateEvent{GuildID: testGuildID, Role: discord.Role{ID: testRoleID, Permissions: newPerms}})
	waitFor(t, 5*time.Second, "permissions to be updated", func() bool {
		return getGuild().Permissions == uint64(newPerms)
	})

	// The bot is removed from the guild.
	dispatch(&gateway.GuildDeleteEvent{ID: testGuildID})
	waitFor(t, 5*time.Second, "guild to be marked as removed", func() bool {
		g := getGuild()
		return g.JoinedAt == nil && g.Permissions == 0
	})
	waitFor(t, 5*time.Second, "actor to exit", func() bool { return w.actor(testGuildID) == nil })

	// Rejoining starts a new actor.
	dispatch(testGuildCreate(time.Now()))
	waitFor(t, 5*time.Second, "guild to be re-synced", func() bool { return getGuild().JoinedAt != nil })

	// Admin choices survive the bot rejoining.
	if getGuild().Edges.GuildAdminConfig.Enabled {
		t.Error("admin config was re-enabled after rejoining")
	}
}

// TestReadyRemovesMissingGuilds checks that guilds which the bot was removed
// from while it was disconnected (so no guild delete event is received) are
// marked as removed when the shard is ready, both for guilds known from the
// database (e.g. after a restart), and for guilds with a running actor (after
// reconnecting with a new session). Guilds of other shards are left alone.
func TestReadyRemovesMissingGuilds(t *testing.T) {
	t.Parallel()

	db := dbtest.New(t)
	ctx, cancel := context.WithCancel(privacy.DecisionContext(t.Context(), privacy.Allow))
	defer cancel()

	w := newTestWorker(ctx, t, db, models.ConfigWorker{NumShards: 2, ShardIDs: []uint{0}})
	defer cancel()

	s := w.shards[0]
	dispatch := func(ev any) { s.state.Session.Handler.Call(ev) }

	const (
		kept         discord.GuildID = testGuildID // Shard 0.
		removed      discord.GuildID = 2000        // Shard 0, removed while offline.
		removedLive  discord.GuildID = 3000        // Shard 0, removed while reconnecting.
		otherShard   discord.GuildID = 1 << 22     // Shard 1.
		stillPresent discord.GuildID = 4000        // Shard 0, in the second ready event.
	)

	for _, id := range []discord.GuildID{kept, removed, removedLive, otherShard, stillPresent} {
		ev := testGuildCreate(time.Now())
		ev.ID = id
		if err := dbUpsertGuild(ctx, db, ev, 0); err != nil {
			t.Fatal(err)
		}
	}

	joined := func(id discord.GuildID) bool {
		g, gerr := db.Guild.Query().Where(guild.GuildID(id.String())).Only(ctx)
		if gerr != nil {
			t.Fatal(gerr)
		}
		return g.JoinedAt != nil
	}

	ready := func(ids ...discord.GuildID) *gateway.ReadyEvent {
		ev := &gateway.ReadyEvent{User: discord.User{ID: testBotID, Bot: true}}
		for _, id := range ids {
			ev.Guilds = append(ev.Guilds, gateway.GuildCreateEvent{ID: id, Unavailable: true})
		}
		return ev
	}

	// Initial connection: "removed" isn't in the ready event.
	dispatch(ready(kept, removedLive, stillPresent))
	for _, id := range []discord.GuildID{kept, removedLive, stillPresent} {
		ev := testGuildCreate(time.Now())
		ev.ID = id
		dispatch(ev)
	}

	waitFor(t, 5*time.Second, "removed guild to be marked as removed", func() bool { return !joined(removed) })

	// Reconnecting with a new session: "removedLive" isn't in the ready event.
	waitFor(t, 5*time.Second, "actor to start", func() bool { return w.actor(removedLive) != nil })
	dispatch(ready(kept, stillPresent))
	waitFor(t, 5*time.Second, "live guild to be marked as removed", func() bool { return !joined(removedLive) })
	waitFor(t, 5*time.Second, "actor to exit", func() bool { return w.actor(removedLive) == nil })

	for _, id := range []discord.GuildID{kept, otherShard, stillPresent} {
		if !joined(id) {
			t.Errorf("guild %d was marked as removed", id)
		}
	}
}

// TestActorRetriesTransientErrors checks that a pass which fails with a
// transient error (here, the guild's config is missing from the database) is
// retried without further events for the guild.
func TestActorRetriesTransientErrors(t *testing.T) {
	t.Parallel()

	db := dbtest.New(t)
	ctx, cancel := context.WithCancel(privacy.DecisionContext(t.Context(), privacy.Allow))
	defer cancel()

	w := newTestWorker(ctx, t, db, models.ConfigWorker{})
	w.minRetry, w.maxRetry = 100*time.Millisecond, 200*time.Millisecond
	defer cancel()

	s := w.shards[0]
	dispatch := func(ev any) { s.state.Session.Handler.Call(ev) }

	dispatch(&gateway.ReadyEvent{User: discord.User{ID: testBotID, Bot: true}})
	create := testGuildCreate(time.Now())
	create.Members = []discord.Member{{User: discord.User{ID: testBotID}}}
	dispatch(create)

	var g *ent.Guild
	waitFor(t, 5*time.Second, "guild to be synced", func() bool {
		g, _ = db.Guild.Query().Where(guild.GuildID(testGuildID.String())).WithGuildConfig().Only(ctx)
		return g != nil && g.Edges.GuildConfig != nil
	})

	// The guild config goes missing, and a pass fails.
	if err := db.GuildConfig.DeleteOne(g.Edges.GuildConfig).Exec(ctx); err != nil {
		t.Fatal(err)
	}
	dispatch(&gateway.VoiceStateUpdateEvent{GuildID: testGuildID, ChannelID: 1, UserID: 500})
	time.Sleep(debounceQuiet + 200*time.Millisecond)

	// Once the config is back, the pass is retried (without further events).
	if err := db.GuildConfig.Create().SetGuildID(g.ID).SetEnabled(false).Exec(ctx); err != nil {
		t.Fatal(err)
	}
	waitFor(t, 10*time.Second, "pass to be retried", func() bool {
		n, _ := db.GuildEvent.Query().Where(guildevent.MessageContains("guild is disabled")).Count(ctx)
		return n == 1
	})
}

// TestActorForbiddenCooldown checks that after Discord rejects a change for
// missing permissions (which the bot didn't predict), changes aren't attempted
// again on every event, as invalid requests count towards Discord's limit
// (after which the bot is temporarily banned).
func TestActorForbiddenCooldown(t *testing.T) {
	t.Parallel()

	db := dbtest.New(t)
	ctx, cancel := context.WithCancel(privacy.DecisionContext(t.Context(), privacy.Allow))
	defer cancel()

	rt := &forbiddenTransport{}
	w := newTestWorkerTransport(ctx, t, db, models.ConfigWorker{}, rt)
	w.cooldown = 3 * time.Second
	defer cancel()

	s := w.shards[0]
	dispatch := func(ev any) { s.state.Session.Handler.Call(ev) }

	dispatch(&gateway.ReadyEvent{User: discord.User{ID: testBotID, Bot: true}})
	create := testGuildCreate(time.Now())
	create.Roles = []discord.Role{{ID: everyoneRole}, {ID: testRoleID, Permissions: botPerms}}
	create.Members = []discord.Member{{User: discord.User{ID: testBotID}, RoleIDs: []discord.RoleID{testRoleID}}}
	create.Channels = []discord.Channel{voice(1, "Gaming +", 1)}
	create.VoiceStates = []discord.VoiceState{{UserID: 500, ChannelID: 1}}
	dispatch(create)

	waitFor(t, 5*time.Second, "clone to be attempted", func() bool { return rt.requests.Load() == 1 })

	// More users join, each would otherwise retry the clone.
	for uid := range discord.UserID(3) {
		dispatch(&gateway.VoiceStateUpdateEvent{GuildID: testGuildID, ChannelID: 1, UserID: 600 + uid})
		time.Sleep(debounceQuiet + 200*time.Millisecond)
	}

	if n := rt.requests.Load(); n != 1 {
		t.Errorf("requests = %d, want 1", n)
	}

	// Without further events, the clone is retried once the cooldown ends.
	waitFor(t, 5*time.Second, "clone to be retried after the cooldown", func() bool { return rt.requests.Load() == 2 })

	n, err := db.GuildEvent.Query().Where(guildevent.MessageContains("unable to clone channel")).Count(ctx)
	if err != nil || n != 1 {
		t.Errorf("clone errors logged = %d (err %v), want 1", n, err)
	}
}
