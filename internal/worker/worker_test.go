// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package worker

import (
	"context"
	"errors"
	"log/slog"
	"slices"
	"sync"
	"testing"
	"testing/synctest"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state/store/defaultstore"
	"github.com/lrstanley/spectrograph/internal/database/ent"
	"github.com/lrstanley/spectrograph/internal/models"
)

func TestEffectiveSettings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		cfg       ent.GuildConfig
		admin     ent.GuildAdminConfig
		want      guildSettings
		wantRegex string
		wantErr   error
	}{
		{
			name:      "defaults",
			cfg:       ent.GuildConfig{Enabled: true},
			admin:     ent.GuildAdminConfig{Enabled: true},
			want:      guildSettings{enabled: true, maxGroups: models.DefaultMaxChannels, maxClones: models.DefaultMaxClones},
			wantRegex: models.DefaultChannelMatch.String(),
		},
		{
			name:      "user-disabled",
			cfg:       ent.GuildConfig{Enabled: false},
			admin:     ent.GuildAdminConfig{Enabled: true},
			want:      guildSettings{enabled: false, maxGroups: models.DefaultMaxChannels, maxClones: models.DefaultMaxClones},
			wantRegex: models.DefaultChannelMatch.String(),
		},
		{
			name:      "admin-disabled-overrides-user",
			cfg:       ent.GuildConfig{Enabled: true},
			admin:     ent.GuildAdminConfig{Enabled: false},
			want:      guildSettings{enabled: false, maxGroups: models.DefaultMaxChannels, maxClones: models.DefaultMaxClones},
			wantRegex: models.DefaultChannelMatch.String(),
		},
		{
			name:      "user-max-clones",
			cfg:       ent.GuildConfig{Enabled: true, DefaultMaxClones: 5},
			admin:     ent.GuildAdminConfig{Enabled: true},
			want:      guildSettings{enabled: true, maxGroups: models.DefaultMaxChannels, maxClones: 5},
			wantRegex: models.DefaultChannelMatch.String(),
		},
		{
			name:      "admin-caps-user-max-clones",
			cfg:       ent.GuildConfig{Enabled: true, DefaultMaxClones: 15},
			admin:     ent.GuildAdminConfig{Enabled: true, DefaultMaxClones: 3, DefaultMaxChannels: 2},
			want:      guildSettings{enabled: true, maxGroups: 2, maxClones: 3},
			wantRegex: models.DefaultChannelMatch.String(),
		},
		{
			name:      "custom-regex",
			cfg:       ent.GuildConfig{Enabled: true, RegexMatch: `^Gaming \d+$`},
			admin:     ent.GuildAdminConfig{Enabled: true},
			want:      guildSettings{enabled: true, maxGroups: models.DefaultMaxChannels, maxClones: models.DefaultMaxClones},
			wantRegex: `^Gaming \d+$`,
		},
		{
			name:      "invalid-regex-falls-back",
			cfg:       ent.GuildConfig{Enabled: true, RegexMatch: `^Gaming (`},
			admin:     ent.GuildAdminConfig{Enabled: true},
			want:      guildSettings{enabled: true, maxGroups: models.DefaultMaxChannels, maxClones: models.DefaultMaxClones},
			wantRegex: models.DefaultChannelMatch.String(),
			wantErr:   errInvalidRegex,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := effectiveSettings(&tt.cfg, &tt.admin)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("effectiveSettings() error = %v, want %v", err, tt.wantErr)
			}

			if got.match == nil || got.match.String() != tt.wantRegex {
				t.Errorf("match = %v, want %q", got.match, tt.wantRegex)
			}
			got.match = nil
			if got != tt.want {
				t.Errorf("effectiveSettings() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestDebounce(t *testing.T) {
	t.Parallel()

	const quiet, maxDelay = 500 * time.Millisecond, 3 * time.Second

	t.Run("quiet", func(t *testing.T) {
		t.Parallel()
		synctest.Test(t, func(t *testing.T) {
			wake := make(chan struct{}, 1)
			start := time.Now()

			go func() {
				// A burst of events, each resetting the quiet period.
				for range 3 {
					time.Sleep(200 * time.Millisecond)
					wake <- struct{}{}
				}
			}()

			if !debounce(t.Context(), wake, quiet, maxDelay) {
				t.Fatal("debounce() = false, want true")
			}
			if got, want := time.Since(start), 600*time.Millisecond+quiet; got != want {
				t.Errorf("debounced for %v, want %v", got, want)
			}
		})
	})

	t.Run("max-delay", func(t *testing.T) {
		t.Parallel()
		synctest.Test(t, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			var wg sync.WaitGroup
			defer wg.Wait()
			defer cancel()

			wake := make(chan struct{}, 1)
			start := time.Now()

			wg.Go(func() {
				// A continuous stream of events.
				for ctx.Err() == nil {
					time.Sleep(100 * time.Millisecond)
					select {
					case wake <- struct{}{}:
					default:
					}
				}
			})

			if !debounce(ctx, wake, quiet, maxDelay) {
				t.Fatal("debounce() = false, want true")
			}
			if got := time.Since(start); got != maxDelay {
				t.Errorf("debounced for %v, want %v", got, maxDelay)
			}
		})
	})

	t.Run("cancelled", func(t *testing.T) {
		t.Parallel()
		synctest.Test(t, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			cancel()
			if debounce(ctx, make(chan struct{}), quiet, maxDelay) {
				t.Error("debounce() = true, want false")
			}
		})
	})
}

func TestSyncGuildCache(t *testing.T) {
	t.Parallel()

	cab := defaultstore.New()

	initial := &gateway.GuildCreateEvent{
		ID: testGuildID, Roles: []discord.Role{{ID: everyoneRole}, {ID: testRoleID}},
		Channels:    []discord.Channel{voice(1, "Gaming +", 1), voice(2, "Gaming +", 2)},
		VoiceStates: []discord.VoiceState{{UserID: 1, ChannelID: 1}, {UserID: 2, ChannelID: 2}},
	}

	for _, ch := range initial.Channels {
		_ = cab.ChannelSet(&ch, false)
	}
	for _, r := range initial.Roles {
		_ = cab.RoleSet(testGuildID, &r, false)
	}
	for _, vs := range initial.VoiceStates {
		vs.GuildID = testGuildID
		_ = cab.VoiceStateSet(testGuildID, &vs, false)
	}

	// The guild becomes available again, after channel 2 was deleted, the
	// role was deleted, and user 2 disconnected.
	syncGuildCache(cab, &gateway.GuildCreateEvent{
		ID: testGuildID, Roles: initial.Roles[:1],
		Channels:    initial.Channels[:1],
		VoiceStates: initial.VoiceStates[:1],
	})

	channels, _ := cab.Channels(testGuildID)
	if ids := channelIDs(channels); !slices.Equal(ids, []discord.ChannelID{1}) {
		t.Errorf("channels = %v, want [1]", ids)
	}

	roles, _ := cab.Roles(testGuildID)
	if len(roles) != 1 || roles[0].ID != everyoneRole {
		t.Errorf("roles = %v, want only @everyone", roles)
	}

	states, _ := cab.VoiceStates(testGuildID)
	if len(states) != 1 || states[0].UserID != 1 {
		t.Errorf("voice states = %v, want only user 1", states)
	}

	removeGuildCache(cab, testGuildID)

	channels, _ = cab.Channels(testGuildID)
	roles, _ = cab.Roles(testGuildID)
	states, _ = cab.VoiceStates(testGuildID)
	if len(channels)+len(roles)+len(states) > 0 {
		t.Errorf("cache not empty after removal: channels=%v roles=%v states=%v", channels, roles, states)
	}
}

func TestNewValidation(t *testing.T) {
	t.Parallel()

	logger := slog.New(slog.DiscardHandler)

	tests := []struct {
		name    string
		config  models.ConfigWorker
		wantErr bool
	}{
		{name: "no-token", config: models.ConfigWorker{}, wantErr: true},
		{name: "bot-prefix-only", config: models.ConfigWorker{BotToken: "Bot "}, wantErr: true},
		{name: "shard-out-of-range", config: models.ConfigWorker{BotToken: "x", NumShards: 2, ShardIDs: []uint{2}}, wantErr: true},
		{name: "defaults", config: models.ConfigWorker{BotToken: "x", EventFlushInterval: time.Second}},
		{name: "multiple-shards", config: models.ConfigWorker{BotToken: "x", NumShards: 4, ShardIDs: []uint{1, 3}, EventFlushInterval: time.Second}},
		{name: "commands-guild", config: models.ConfigWorker{BotToken: "x", CommandsGuildID: "123456789012345678", EventFlushInterval: time.Second}},
		{name: "invalid-commands-guild", config: models.ConfigWorker{BotToken: "x", CommandsGuildID: "my-guild"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w, err := New(t.Context(), logger, nil, tt.config)
			if (err != nil) != tt.wantErr {
				t.Fatalf("New() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}

			want := max(len(tt.config.ShardIDs), 1)
			if len(w.shards) != want {
				t.Errorf("shards = %d, want %d", len(w.shards), want)
			}
			for i, s := range w.shards {
				if !s.state.HasIntents(intents) {
					t.Errorf("shard %d missing intents", i)
				}
			}
		})
	}
}

// TestGuildAvailableRefreshesCache checks that when a guild becomes available
// again (after an outage), the cache reflects the guild create event, rather
// than what was cached before the outage. The state cache doesn't overwrite
// existing roles, members, or voice states on guild create.
func TestGuildAvailableRefreshesCache(t *testing.T) {
	t.Parallel()

	w, err := New(t.Context(), slog.New(slog.DiscardHandler), nil, models.ConfigWorker{BotToken: "x", EventFlushInterval: time.Second})
	if err != nil {
		t.Fatal(err)
	}
	st := w.shards[0].state
	dispatch := func(ev any) { st.Session.Handler.Call(ev) }

	guildCreate := func(rolePerms discord.Permissions, botRoles []discord.RoleID, userChannel discord.ChannelID) *gateway.GuildCreateEvent {
		ev := testGuildCreate(time.Now())
		ev.Roles = []discord.Role{{ID: everyoneRole}, {ID: testRoleID, Permissions: rolePerms}}
		ev.Members = []discord.Member{{User: discord.User{ID: testBotID}, RoleIDs: botRoles}}
		ev.Channels = []discord.Channel{voice(1, "Gaming +", 1), voice(2, "Gaming +", 2)}
		ev.VoiceStates = []discord.VoiceState{{UserID: testOtherUser, ChannelID: userChannel}}
		return ev
	}

	dispatch(&gateway.ReadyEvent{User: discord.User{ID: testBotID, Bot: true}})
	dispatch(guildCreate(discord.PermissionViewChannel, nil, 1))

	// While the guild was unavailable, the user moved to channel 2, the role
	// was granted Manage Channels, and the bot was given the role.
	dispatch(&gateway.GuildDeleteEvent{ID: testGuildID, Unavailable: true})
	dispatch(guildCreate(botPerms, []discord.RoleID{testRoleID}, 2))

	vs, err := st.Cabinet.VoiceState(testGuildID, testOtherUser)
	if err != nil || vs.ChannelID != 2 {
		t.Errorf("voice state = %+v (err %v), want user in channel 2", vs, err)
	}

	role, err := st.Cabinet.Role(testGuildID, testRoleID)
	if err != nil || role.Permissions != botPerms {
		t.Errorf("role = %+v (err %v), want permissions %v", role, err, botPerms)
	}

	m, err := st.Cabinet.Member(testGuildID, testBotID)
	if err != nil || !slices.Equal(m.RoleIDs, []discord.RoleID{testRoleID}) {
		t.Errorf("bot member = %+v (err %v), want roles [%v]", m, err, testRoleID)
	}
}
