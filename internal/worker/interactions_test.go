// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package worker

import (
	"context"
	"log/slog"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state/store"
	"github.com/diamondburned/arikawa/v3/state/store/defaultstore"
	"github.com/lrstanley/spectrograph/internal/database/dbtest"
	"github.com/lrstanley/spectrograph/internal/database/ent"
	"github.com/lrstanley/spectrograph/internal/database/ent/guild"
	"github.com/lrstanley/spectrograph/internal/database/ent/guildevent"
	"github.com/lrstanley/spectrograph/internal/database/ent/privacy"
	"github.com/lrstanley/spectrograph/internal/models"
)

// fakeResponder records interaction responses.
type fakeResponder struct {
	mu        sync.Mutex
	responses []*api.InteractionResponse
	edits     []*api.InteractionResponseData
	followUps []*api.InteractionResponseData
}

func (f *fakeResponder) respond(_ context.Context, _ *discord.InteractionEvent, resp *api.InteractionResponse) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.responses = append(f.responses, resp)
	return nil
}

func (f *fakeResponder) editOriginal(_ context.Context, _ *discord.InteractionEvent, data *api.InteractionResponseData) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.edits = append(f.edits, data)
	return nil
}

func (f *fakeResponder) followUp(_ context.Context, _ *discord.InteractionEvent, data *api.InteractionResponseData) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.followUps = append(f.followUps, data)
	return nil
}

func mustParseURL(t *testing.T, raw string) *url.URL {
	t.Helper()
	u, err := url.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	return u
}

func itoa(i int) string { return strconv.Itoa(i) }

func dataText(data *api.InteractionResponseData) string {
	if data == nil || data.Components == nil {
		return ""
	}
	components := make([]discord.Component, 0, len(*data.Components))
	for _, c := range *data.Components {
		components = append(components, c)
	}
	return viewText(components)
}

func TestServeInteraction(t *testing.T) {
	t.Parallel()

	status := update(text("status"))
	slow := func(resp *api.InteractionResponse) func(context.Context) *api.InteractionResponse {
		return func(context.Context) *api.InteractionResponse {
			time.Sleep(50 * time.Millisecond)
			return resp
		}
	}

	tests := []struct {
		name          string
		component     bool // Interaction from a message component.
		handle        func(context.Context) *api.InteractionResponse
		wantResponses []api.InteractionResponseType
		wantEdit      string
		wantFollowUp  string
	}{
		{
			name:          "fast",
			handle:        func(context.Context) *api.InteractionResponse { return status },
			wantResponses: []api.InteractionResponseType{api.UpdateMessage},
		},
		{
			name:   "not-handled",
			handle: func(context.Context) *api.InteractionResponse { return nil },
		},
		{
			name:          "slow-command",
			handle:        slow(message(text("status"))),
			wantResponses: []api.InteractionResponseType{api.DeferredMessageInteractionWithSource},
			wantEdit:      "status",
		},
		{
			name:          "slow-component-update",
			component:     true,
			handle:        slow(status),
			wantResponses: []api.InteractionResponseType{api.DeferredMessageUpdate},
			wantEdit:      "status",
		},
		{
			// Errors don't replace the message the component is attached to.
			name:          "slow-component-error",
			component:     true,
			handle:        slow(errorMessage("boom")),
			wantResponses: []api.InteractionResponseType{api.DeferredMessageUpdate},
			wantFollowUp:  "boom",
		},
		{
			name:          "slow-modal",
			handle:        slow(settingsModal(&ent.GuildConfig{}, &ent.GuildAdminConfig{})),
			wantResponses: []api.InteractionResponseType{api.DeferredMessageInteractionWithSource},
			wantEdit:      msgSlowModal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ev := &discord.InteractionEvent{ID: 1, Token: "token"}
			if tt.component {
				ev.Message = &discord.Message{ID: 2}
			}

			r := &fakeResponder{}
			serveInteraction(t.Context(), slog.New(slog.DiscardHandler), r, ev, 10*time.Millisecond, tt.handle)

			var types []api.InteractionResponseType
			for _, resp := range r.responses {
				types = append(types, resp.Type)
			}
			if len(types) != len(tt.wantResponses) || (len(types) > 0 && types[0] != tt.wantResponses[0]) {
				t.Errorf("responses = %v, want %v", types, tt.wantResponses)
			}
			if len(r.responses) > 0 && r.responses[0].Type == api.DeferredMessageInteractionWithSource &&
				r.responses[0].Data.Flags&discord.EphemeralMessage == 0 {
				t.Error("deferred response isn't ephemeral")
			}

			if (len(r.edits) > 0) != (tt.wantEdit != "") || (tt.wantEdit != "" && !strings.Contains(dataText(r.edits[0]), tt.wantEdit)) {
				t.Errorf("edits = %d, want one containing %q", len(r.edits), tt.wantEdit)
			}
			if (len(r.followUps) > 0) != (tt.wantFollowUp != "") || (tt.wantFollowUp != "" && !strings.Contains(dataText(r.followUps[0]), tt.wantFollowUp)) {
				t.Errorf("follow-ups = %d, want one containing %q", len(r.followUps), tt.wantFollowUp)
			}
		})
	}
}

// interactionTest is a guild synced to the database, with a cached state, and
// a worker to handle interactions for it.
type interactionTest struct {
	ctx   context.Context //nolint:containedctx
	db    *ent.Client
	w     *Worker
	cache *store.Cabinet
}

const (
	testManageRole discord.RoleID = 3001
	testAdminUser  discord.UserID = 2003
)

func newInteractionTest(t *testing.T) *interactionTest {
	t.Helper()

	db := dbtest.New(t)
	ctx, cancel := context.WithCancel(privacy.DecisionContext(t.Context(), privacy.Allow))
	t.Cleanup(cancel)

	w := newTestWorker(ctx, t, db, models.ConfigWorker{})

	create := testGuildCreate(time.Now())
	if err := dbUpsertGuild(ctx, db, create, botPerms); err != nil {
		t.Fatal(err)
	}

	cache := defaultstore.New()
	_ = cache.MyselfSet(discord.User{ID: testBotID}, false)
	_ = cache.GuildSet(&create.Guild, false)
	for _, r := range []discord.Role{
		{ID: everyoneRole, Permissions: discord.PermissionViewChannel},
		{ID: testRoleID, Permissions: botPerms},
		{ID: testManageRole, Permissions: discord.PermissionManageGuild},
	} {
		_ = cache.RoleSet(testGuildID, &r, false)
	}
	_ = cache.MemberSet(testGuildID, &discord.Member{User: discord.User{ID: testBotID}, RoleIDs: []discord.RoleID{testRoleID}}, false)
	for _, ch := range []discord.Channel{voice(1, "Gaming +", 1), voice(2, "Gaming +", 2)} {
		_ = cache.ChannelSet(&ch, false)
	}
	_ = cache.VoiceStateSet(testGuildID, &discord.VoiceState{GuildID: testGuildID, ChannelID: 1, UserID: 500}, false)

	return &interactionTest{ctx: ctx, db: db, w: w, cache: cache}
}

// event returns an interaction from a guild admin (with the Manage Server
// permission through a role).
func (it *interactionTest) event(data discord.InteractionData) *discord.InteractionEvent {
	return &discord.InteractionEvent{
		ID:      1,
		AppID:   2,
		Token:   "token",
		GuildID: testGuildID,
		Data:    data,
		Member: &discord.Member{
			User:    discord.User{ID: testAdminUser, Username: "alice"},
			RoleIDs: []discord.RoleID{testManageRole},
		},
	}
}

func command(sub string) *discord.CommandInteraction {
	return &discord.CommandInteraction{
		Name:    commandName,
		Options: []discord.CommandInteractionOption{{Type: discord.SubcommandOptionType, Name: sub}},
	}
}

// handle handles the interaction, and returns the response and its text.
func (it *interactionTest) handle(t *testing.T, ev *discord.InteractionEvent) (resp *api.InteractionResponse, text string) {
	t.Helper()
	resp = it.w.interactions.handle(it.ctx, it.cache, ev)
	if resp == nil {
		t.Fatal("handle() = nil, want a response")
	}
	if resp.Type == api.ModalResponse {
		return resp, ""
	}
	return resp, viewText(responseComponents(t, resp))
}

func (it *interactionTest) config(t *testing.T) *ent.GuildConfig {
	t.Helper()
	return getGuild(it.ctx, t, it.db).Edges.GuildConfig
}

// waitForEvent waits for a guild event containing msg to be persisted.
func (it *interactionTest) waitForEvent(t *testing.T, msg string) *ent.GuildEvent {
	t.Helper()
	var ev *ent.GuildEvent
	waitFor(t, 5*time.Second, "event "+msg, func() bool {
		ev, _ = it.db.GuildEvent.Query().
			Where(guildevent.HasGuildWith(guild.GuildID(testGuildID.String())), guildevent.MessageContains(msg)).
			First(it.ctx)
		return ev != nil
	})
	return ev
}

func TestInteractionAuthorization(t *testing.T) {
	t.Parallel()
	it := newInteractionTest(t)

	t.Run("not-in-guild", func(t *testing.T) {
		t.Parallel()

		ev := it.event(command(subcommandStatus))
		ev.GuildID, ev.Member, ev.User = 0, nil, &discord.User{ID: testAdminUser}
		if _, got := it.handle(t, ev); !strings.Contains(got, msgNotInGuild) {
			t.Errorf("text = %q, want %q", got, msgNotInGuild)
		}
	})

	t.Run("not-guild-admin", func(t *testing.T) {
		t.Parallel()

		ev := it.event(command(subcommandDisable))
		ev.Member.RoleIDs = nil
		if _, got := it.handle(t, ev); !strings.Contains(got, msgNotAuthorized) {
			t.Errorf("text = %q, want %q", got, msgNotAuthorized)
		}
		if !it.config(t).Enabled {
			t.Error("guild was disabled by a member who isn't a guild admin")
		}
	})

	t.Run("owner", func(t *testing.T) {
		t.Parallel()

		ev := it.event(command(subcommandStatus))
		ev.Member = &discord.Member{User: discord.User{ID: testOwnerID, Username: "owner"}}
		if _, got := it.handle(t, ev); !strings.Contains(got, "**Status:** enabled") {
			t.Errorf("text = %q, want the status", got)
		}
	})

	t.Run("banned", func(t *testing.T) {
		t.Parallel()

		err := it.db.User.Create().
			SetUserID(testOtherUser.String()).
			SetUsername("banned").
			SetAvatarURL("https://example.com/a.png").
			SetBanned(true).
			Exec(it.ctx)
		if err != nil {
			t.Fatal(err)
		}

		ev := it.event(command(subcommandDisable))
		ev.Member.User.ID = testOtherUser
		if _, got := it.handle(t, ev); !strings.Contains(got, msgBanned) {
			t.Errorf("text = %q, want %q", got, msgBanned)
		}
		if !it.config(t).Enabled {
			t.Error("guild was disabled by a banned user")
		}
	})
}

func TestInteractionStatus(t *testing.T) {
	t.Parallel()
	it := newInteractionTest(t)
	it.w.interactions.dashboardURL = mustParseURL(t, "https://example.com/")

	resp, got := it.handle(t, it.event(command(subcommandStatus)))
	if resp.Type != api.MessageInteractionWithSource || resp.Data.Flags != viewFlags {
		t.Errorf("response = %v (flags %v), want an ephemeral message", resp.Type, resp.Data.Flags)
	}
	for _, want := range []string{"Spectrograph in test guild", "**Status:** enabled", "Gaming +: 1 clone (1 in use)"} {
		if !strings.Contains(got, want) {
			t.Errorf("text doesn't contain %q:\n%s", want, got)
		}
	}

	g := getGuild(it.ctx, t, it.db)
	buttons := viewButtons(responseComponents(t, resp))
	if b := buttons["Open dashboard"]; b == nil || !strings.HasSuffix(b.LinkURL(), "/dashboard/guild/"+itoa(g.ID)) {
		t.Errorf("dashboard button = %+v", b)
	}

	// The refresh button updates the message.
	ev := it.event(&discord.ButtonInteraction{CustomID: customIDRefresh})
	ev.Message = &discord.Message{ID: 5}
	if resp, _ = it.handle(t, ev); resp.Type != api.UpdateMessage {
		t.Errorf("refresh response = %v, want a message update", resp.Type)
	}
}

func TestInteractionEnableDisable(t *testing.T) {
	t.Parallel()
	it := newInteractionTest(t)

	// Disable from a button, which updates the status message.
	ev := it.event(&discord.ButtonInteraction{CustomID: customIDDisable})
	ev.Message = &discord.Message{ID: 5}
	resp, got := it.handle(t, ev)
	if resp.Type != api.UpdateMessage || !strings.Contains(got, "Spectrograph is now disabled.") {
		t.Errorf("response = %v %q, want the updated status", resp.Type, got)
	}
	if it.config(t).Enabled {
		t.Error("guild wasn't disabled")
	}
	logged := it.waitForEvent(t, "guild disabled via Discord by alice")
	if logged.Type != guildevent.TypeINFO || logged.Metadata["user_id"] != testAdminUser.String() {
		t.Errorf("event = %+v, want an info event with the user", logged)
	}

	if _, got = it.handle(t, it.event(command(subcommandDisable))); !strings.Contains(got, "already disabled") {
		t.Errorf("text = %q, want already disabled", got)
	}

	// Guilds disabled by system admins can't be enabled by guild admins.
	g := getGuild(it.ctx, t, it.db)
	if err := it.db.GuildAdminConfig.UpdateOne(g.Edges.GuildAdminConfig).SetEnabled(false).Exec(it.ctx); err != nil {
		t.Fatal(err)
	}
	if _, got = it.handle(t, it.event(command(subcommandEnable))); !strings.Contains(got, msgAdminDisabled) {
		t.Errorf("text = %q, want %q", got, msgAdminDisabled)
	}
	if it.config(t).Enabled {
		t.Error("guild disabled by admins was enabled")
	}

	if err := it.db.GuildAdminConfig.UpdateOne(g.Edges.GuildAdminConfig).SetEnabled(true).Exec(it.ctx); err != nil {
		t.Fatal(err)
	}
	if _, got = it.handle(t, it.event(command(subcommandEnable))); !strings.Contains(got, "Spectrograph is now enabled.") {
		t.Errorf("text = %q, want enabled", got)
	}
	if !it.config(t).Enabled {
		t.Error("guild wasn't enabled")
	}
}

func TestInteractionSettings(t *testing.T) {
	t.Parallel()
	it := newInteractionTest(t)

	resp, _ := it.handle(t, it.event(command(subcommandSettings)))
	if resp.Type != api.ModalResponse {
		t.Fatalf("response = %v, want a modal", resp.Type)
	}

	submit := func(values map[discord.ComponentID]string) (*api.InteractionResponse, string) {
		var components discord.TopLevelComponents
		for id, v := range values {
			var input discord.Component = &discord.TextInputComponent{CustomID: id, Value: v}
			if id == customIDMaxClones {
				input = &discord.StringSelectComponent{CustomID: id, Values: []string{v}}
			}
			components = append(components, &discord.LabelComponent{Component: input})
		}
		return it.handle(t, it.event(&discord.ModalInteraction{CustomID: customIDSettingsModal, Components: components}))
	}

	tests := []struct {
		name   string
		values map[discord.ComponentID]string
		want   string
	}{
		{"invalid-regex", map[discord.ComponentID]string{customIDRegexMatch: "^(foo"}, "**Channel name pattern:** error parsing regexp"},
		{"regex-too-long", map[discord.ComponentID]string{customIDRegexMatch: strings.Repeat("a", 101)}, "**Channel name pattern:** value is too long"},
		{"invalid-email", map[discord.ComponentID]string{customIDContactEmail: "nope"}, "**Contact email:** invalid email address"},
		{"invalid-max-clones", map[discord.ComponentID]string{customIDMaxClones: "21"}, "**Max clones per channel:** choose a number from 0 to 20"},
		{"no-changes", map[discord.ComponentID]string{customIDRegexMatch: "", customIDMaxClones: "0", customIDContactEmail: ""}, "No changes to save."},
	}
	// Not subtests, as the config is checked once all are submitted.
	for _, tt := range tests {
		errResp, got := submit(tt.values)
		if !strings.Contains(got, tt.want) {
			t.Errorf("%s: text = %q, want %q", tt.name, got, tt.want)
		}
		if tt.name != "no-changes" && viewButtons(responseComponents(t, errResp))["Edit settings"] == nil {
			t.Errorf("%s: error has no button to edit the settings again", tt.name)
		}
	}

	cfg := it.config(t)
	if cfg.RegexMatch != "" || cfg.DefaultMaxClones != 0 || cfg.ContactEmail != "" {
		t.Fatalf("config = %+v, want unchanged after invalid submissions", cfg)
	}

	_, got := submit(map[discord.ComponentID]string{
		customIDRegexMatch:   " ^Game.*$ ",
		customIDMaxClones:    "5",
		customIDContactEmail: "a@example.com",
	})
	if !strings.Contains(got, "Settings saved.") || !strings.Contains(got, "`^Game.*$`") {
		t.Errorf("text = %q, want the saved settings", got)
	}

	cfg = it.config(t)
	if cfg.RegexMatch != "^Game.*$" || cfg.DefaultMaxClones != 5 || cfg.ContactEmail != "a@example.com" {
		t.Errorf("config = %+v, want the submitted settings", cfg)
	}

	logged := it.waitForEvent(t, "settings updated via Discord by alice")
	if logged.Metadata["regex_match"] != "^Game.*$" || logged.Metadata["contact_email"] != "a@example.com" {
		t.Errorf("event metadata = %v, want the changes", logged.Metadata)
	}
}

func TestInteractionEvents(t *testing.T) {
	t.Parallel()
	it := newInteractionTest(t)

	g := getGuild(it.ctx, t, it.db)
	for _, e := range []struct {
		typ guildevent.Type
		msg string
	}{
		{guildevent.TypeDEBUG, "debug only"},
		{guildevent.TypeWARNING, "missing permissions to clone channel"},
		{guildevent.TypeINFO, "newest event"},
	} {
		if err := it.db.GuildEvent.Create().SetGuild(g).SetType(e.typ).SetMessage(e.msg).Exec(it.ctx); err != nil {
			t.Fatal(err)
		}
	}

	ev := it.event(&discord.CommandInteraction{
		Name: commandName,
		Options: []discord.CommandInteractionOption{{
			Type: discord.SubcommandOptionType,
			Name: subcommandEvents,
		}},
	})
	_, got := it.handle(t, ev)
	if strings.Contains(got, "debug only") {
		t.Errorf("debug events are shown:\n%s", got)
	}
	newest, older := strings.Index(got, "newest event"), strings.Index(got, "missing permissions to clone channel")
	if newest < 0 || older < 0 || newest > older {
		t.Errorf("events aren't listed newest first:\n%s", got)
	}

	act, _ := parseAction(ev)
	act.count = 1
	resp := it.w.interactions.events(it.ctx, ev, g, act.count)
	if got = viewText(responseComponents(t, resp)); strings.Contains(got, "missing permissions") {
		t.Errorf("count isn't applied:\n%s", got)
	}
}

func TestInteractionNotSetUp(t *testing.T) {
	t.Parallel()
	it := newInteractionTest(t)

	ev := it.event(command(subcommandStatus))
	ev.GuildID = 999
	g := discord.Guild{ID: 999, OwnerID: testAdminUser}
	_ = it.cache.GuildSet(&g, false)

	if _, got := it.handle(t, ev); !strings.Contains(got, msgNotSetUp) {
		t.Errorf("text = %q, want %q", got, msgNotSetUp)
	}
}

// TestConfigChangeWakesActor checks that updating a guild's config (e.g. from
// the dashboard) wakes the guild's actor.
func TestConfigChangeWakesActor(t *testing.T) {
	t.Parallel()
	it := newInteractionTest(t)

	a := newGuildActor(it.w, it.w.shards[0], testGuildID, it.w.client, nil)
	it.w.mu.Lock()
	it.w.actors[testGuildID] = a
	it.w.mu.Unlock()

	g := getGuild(it.ctx, t, it.db)
	updates := []func() error{
		func() error {
			return it.db.GuildConfig.UpdateOne(g.Edges.GuildConfig).SetRegexMatch("^x$").Exec(it.ctx)
		},
		func() error {
			return it.db.GuildAdminConfig.UpdateOne(g.Edges.GuildAdminConfig).SetDefaultMaxClones(3).Exec(it.ctx)
		},
	}
	for i, upd := range updates {
		if err := upd(); err != nil {
			t.Fatal(err)
		}
		select {
		case <-a.wake:
		default:
			t.Errorf("update %d didn't wake the actor", i)
		}
	}
}
