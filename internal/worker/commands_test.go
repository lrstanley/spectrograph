// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package worker

import (
	"encoding/json/v2"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	arikawajson "github.com/diamondburned/arikawa/v3/utils/json"
	"github.com/lrstanley/spectrograph/internal/database/ent"
	"github.com/lrstanley/spectrograph/internal/database/ent/guildevent"
	"github.com/lrstanley/spectrograph/internal/models"
)

// viewText returns the text of all text displays in the components.
func viewText(components []discord.Component) string {
	var b strings.Builder
	var visit func(c discord.Component)
	visit = func(c discord.Component) {
		switch c := c.(type) {
		case *discord.TextDisplayComponent:
			b.WriteString(c.Content + "\n")
		case *discord.ContainerComponent:
			for _, child := range c.Components {
				visit(child)
			}
		case *discord.SectionComponent:
			for _, child := range c.Components {
				visit(child)
			}
		}
	}
	for _, c := range components {
		visit(c)
	}
	return b.String()
}

// viewButtons returns the buttons in the components, by label.
func viewButtons(components []discord.Component) map[string]*discord.ButtonComponent {
	buttons := map[string]*discord.ButtonComponent{}
	for _, c := range components {
		if row, ok := c.(*discord.ActionRowComponent); ok {
			for _, child := range *row {
				if b, isButton := child.(*discord.ButtonComponent); isButton {
					buttons[b.Label] = b
				}
			}
		}
	}
	return buttons
}

// responseComponents returns the top-level components of a response.
func responseComponents(t *testing.T, resp *api.InteractionResponse) []discord.Component {
	t.Helper()
	if resp == nil || resp.Data == nil || resp.Data.Components == nil {
		t.Fatalf("response has no components: %+v", resp)
	}
	components := make([]discord.Component, 0, len(*resp.Data.Components))
	for _, c := range *resp.Data.Components {
		components = append(components, c)
	}
	return components
}

// marshalDiscord marshals v as arikawa does when sending it to Discord.
func marshalDiscord(t *testing.T, v any) map[string]any {
	t.Helper()
	b, err := arikawajson.Marshal(v)
	if err != nil {
		t.Fatalf("marshaling: %v", err)
	}
	var out map[string]any
	if err = json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshaling: %v", err)
	}
	return out
}

func TestCommandDefinitions(t *testing.T) {
	t.Parallel()

	cmds := commandDefinitions()
	if len(cmds) != 1 || cmds[0].Name != commandName {
		t.Fatalf("commands = %+v, want a single %q command", cmds, commandName)
	}

	raw := marshalDiscord(t, cmds[0])
	if raw["default_member_permissions"] != "32" {
		t.Errorf("default_member_permissions = %v, want Manage Server (32)", raw["default_member_permissions"])
	}
	if raw["dm_permission"] != false {
		t.Errorf("dm_permission = %v, want false", raw["dm_permission"])
	}

	var subs []string
	for _, opt := range cmds[0].Options {
		subs = append(subs, opt.Name())
	}
	want := []string{subcommandStatus, subcommandSettings, subcommandEnable, subcommandDisable, subcommandEvents, subcommandHelp}
	if !slices.Equal(subs, want) {
		t.Errorf("subcommands = %v, want %v", subs, want)
	}

	nameRgx := regexp.MustCompile(`^[-_\p{Ll}\p{N}]{1,32}$`)
	for _, opt := range cmds[0].Options {
		sub, ok := opt.(*discord.SubcommandOption)
		if !ok || !nameRgx.MatchString(sub.OptionName) || sub.Description == "" || len(sub.Description) > 100 {
			t.Errorf("invalid subcommand %+v", opt)
		}
	}
}

func TestParseAction(t *testing.T) {
	t.Parallel()

	command := func(sub string, opts ...discord.CommandInteractionOption) *discord.InteractionEvent {
		return &discord.InteractionEvent{Data: &discord.CommandInteraction{
			Name:    commandName,
			Options: []discord.CommandInteractionOption{{Type: discord.SubcommandOptionType, Name: sub, Options: opts}},
		}}
	}
	button := func(id discord.ComponentID) *discord.InteractionEvent {
		return &discord.InteractionEvent{Data: &discord.ButtonInteraction{CustomID: id}, Message: &discord.Message{}}
	}
	count := func(n string) discord.CommandInteractionOption {
		return discord.CommandInteractionOption{Type: discord.IntegerOptionType, Name: optionEventCount, Value: arikawajson.Raw(n)}
	}

	tests := []struct {
		name   string
		ev     *discord.InteractionEvent
		want   action
		wantOK bool
	}{
		{"status", command(subcommandStatus), action{kind: actionStatus}, true},
		{"settings", command(subcommandSettings), action{kind: actionSettings}, true},
		{"enable", command(subcommandEnable), action{kind: actionEnable}, true},
		{"disable", command(subcommandDisable), action{kind: actionDisable}, true},
		{"help", command(subcommandHelp), action{kind: actionHelp}, true},
		{"events-default", command(subcommandEvents), action{kind: actionEvents, count: defaultEventCount}, true},
		{"events-count", command(subcommandEvents, count("3")), action{kind: actionEvents, count: 3}, true},
		{"events-count-capped", command(subcommandEvents, count("500")), action{kind: actionEvents, count: maxEventCount}, true},
		{"unknown-subcommand", command("nope"), action{}, true},
		{"other-command", &discord.InteractionEvent{Data: &discord.CommandInteraction{Name: "other"}}, action{}, false},
		{"refresh-button", button(customIDRefresh), action{kind: actionStatus, update: true}, true},
		{"edit-button", button(customIDEditSettings), action{kind: actionSettings, update: true}, true},
		{"enable-button", button(customIDEnable), action{kind: actionEnable, update: true}, true},
		{"disable-button", button(customIDDisable), action{kind: actionDisable, update: true}, true},
		{"unknown-own-button", button(customIDPrefix + "nope"), action{update: true}, true},
		{"other-button", button("other:thing"), action{}, false},
		{"ping", &discord.InteractionEvent{Data: &discord.PingInteraction{}}, action{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, ok := parseAction(tt.ev)
			if ok != tt.wantOK || fmt.Sprint(got) != fmt.Sprint(tt.want) {
				t.Errorf("parseAction() = %+v, %v, want %+v, %v", got, ok, tt.want, tt.wantOK)
			}
		})
	}
}

// TestParseModalSubmit decodes a modal submission (as sent by Discord, with
// inputs wrapped in labels) through arikawa, and checks the submitted values.
func TestParseModalSubmit(t *testing.T) {
	t.Parallel()

	payload := `{
		"type": 5, "id": "11", "application_id": "12", "token": "tok", "version": 1,
		"guild_id": "1000",
		"member": {"user": {"id": "2001", "username": "owner"}, "roles": [], "permissions": "8"},
		"message": {"id": "13", "channel_id": "14", "content": ""},
		"data": {"custom_id": "spectrograph:settings", "components": [
			{"type": 18, "id": 1, "component": {"type": 4, "id": 2, "custom_id": "spectrograph:regex-match", "value": "^foo$"}},
			{"type": 18, "id": 3, "component": {"type": 3, "id": 4, "custom_id": "spectrograph:max-clones", "values": ["5"]}},
			{"type": 18, "id": 5, "component": {"type": 4, "id": 6, "custom_id": "spectrograph:contact-email", "value": "a@example.com"}}
		]}
	}`

	var ev discord.InteractionEvent
	if err := arikawajson.Unmarshal([]byte(payload), &ev); err != nil {
		t.Fatalf("unmarshaling: %v", err)
	}

	act, ok := parseAction(&ev)
	if !ok || act.kind != actionSaveSettings || !act.update {
		t.Fatalf("parseAction() = %+v, %v, want an update saving the settings", act, ok)
	}

	want := map[discord.ComponentID]string{
		customIDRegexMatch:   "^foo$",
		customIDMaxClones:    "5",
		customIDContactEmail: "a@example.com",
	}
	if fmt.Sprint(act.values) != fmt.Sprint(want) {
		t.Errorf("values = %v, want %v", act.values, want)
	}
}

func TestCanManageGuild(t *testing.T) {
	t.Parallel()

	const (
		adminRole  discord.RoleID = 10
		manageRole discord.RoleID = 11
		plainRole  discord.RoleID = 12
	)

	g := &discord.Guild{ID: testGuildID, OwnerID: testOwnerID}
	roles := []discord.Role{
		{ID: everyoneRole, Permissions: discord.PermissionViewChannel},
		{ID: adminRole, Permissions: discord.PermissionAdministrator},
		{ID: manageRole, Permissions: discord.PermissionManageGuild},
		{ID: plainRole, Permissions: discord.PermissionManageChannels | discord.PermissionManageRoles},
	}

	tests := []struct {
		name   string
		member discord.Member
		want   bool
	}{
		{"owner", discord.Member{User: discord.User{ID: testOwnerID}}, true},
		{"administrator", discord.Member{User: discord.User{ID: testOtherUser}, RoleIDs: []discord.RoleID{adminRole}}, true},
		{"manage-server", discord.Member{User: discord.User{ID: testOtherUser}, RoleIDs: []discord.RoleID{plainRole, manageRole}}, true},
		{"other-permissions", discord.Member{User: discord.User{ID: testOtherUser}, RoleIDs: []discord.RoleID{plainRole}}, false},
		{"no-roles", discord.Member{User: discord.User{ID: testOtherUser}}, false},
		{"unknown-role", discord.Member{User: discord.User{ID: testOtherUser}, RoleIDs: []discord.RoleID{999}}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := canManageGuild(g, roles, &tt.member); got != tt.want {
				t.Errorf("canManageGuild() = %v, want %v", got, tt.want)
			}
		})
	}

	// Manage Server granted to @everyone.
	everyoneManage := []discord.Role{{ID: everyoneRole, Permissions: discord.PermissionManageGuild}}
	if !canManageGuild(g, everyoneManage, &discord.Member{User: discord.User{ID: testOtherUser}}) {
		t.Error("canManageGuild() = false for Manage Server granted to @everyone, want true")
	}
}

func TestSettingsModal(t *testing.T) {
	t.Parallel()

	cfg := &ent.GuildConfig{RegexMatch: "^foo$", DefaultMaxClones: 5, ContactEmail: "a@example.com"}
	resp := settingsModal(cfg, &ent.GuildAdminConfig{Enabled: true, DefaultMaxClones: 8})
	if resp.Type != api.ModalResponse {
		t.Fatalf("type = %v, want a modal", resp.Type)
	}

	raw := marshalDiscord(t, resp.Data)
	if raw["custom_id"] != string(customIDSettingsModal) || raw["title"] == "" {
		t.Errorf("modal custom_id/title = %v/%v", raw["custom_id"], raw["title"])
	}

	components, _ := raw["components"].([]any)
	if len(components) != 3 {
		t.Fatalf("components = %v, want 3 labels", components)
	}

	inputs := map[string]map[string]any{}
	for _, c := range components {
		label, _ := c.(map[string]any)
		if label["type"] != float64(discord.LabelComponentType) {
			t.Errorf("top-level component type = %v, want a label", label["type"])
		}
		if desc, _ := label["description"].(string); len(desc) > 100 {
			t.Errorf("label description is too long (%d): %q", len(desc), desc)
		}
		if l, _ := label["label"].(string); l == "" || len(l) > 45 {
			t.Errorf("invalid label %q", l)
		}
		input, _ := label["component"].(map[string]any)
		inputs[input["custom_id"].(string)] = input
	}

	if v := inputs[string(customIDRegexMatch)]["value"]; v != "^foo$" {
		t.Errorf("regex value = %v, want prefilled", v)
	}
	if v := inputs[string(customIDContactEmail)]["value"]; v != "a@example.com" {
		t.Errorf("email value = %v, want prefilled", v)
	}

	sel := inputs[string(customIDMaxClones)]
	options, _ := sel["options"].([]any)
	if len(options) != maxCloneLimit+1 || len(options) > 25 {
		t.Fatalf("max clones options = %d, want %d (Discord allows up to 25)", len(options), maxCloneLimit+1)
	}
	for _, o := range options {
		opt, _ := o.(map[string]any)
		if (opt["value"] == "5") != (opt["default"] == true) {
			t.Errorf("option %v: default = %v, want only the current value (5) selected", opt["value"], opt["default"])
		}
	}
}

func testGuildView() *guildView {
	cfg := &ent.GuildConfig{Enabled: true}
	admin := &ent.GuildAdminConfig{Enabled: true}
	settings, _ := effectiveSettings(cfg, admin)
	return &guildView{
		name:         "Test *guild*",
		config:       cfg,
		admin:        admin,
		settings:     settings,
		dashboardURL: "https://example.com/dashboard/guild/1",
	}
}

func TestStatusComponents(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		modify      func(v *guildView)
		wantText    []string
		wantToggle  string
		wantEnabled bool // The toggle button isn't disabled.
		wantAccent  discord.Color
	}{
		{
			name:        "enabled",
			wantText:    []string{"**Status:** enabled", "`^.* \\+$` (default)", "20 (default)", "not set", `Test \*guild\*`, "No voice channels match"},
			wantToggle:  "Disable",
			wantEnabled: true,
			wantAccent:  colorEnabled,
		},
		{
			name:        "disabled",
			modify:      func(v *guildView) { v.config.Enabled = false },
			wantText:    []string{"**Status:** disabled, so"},
			wantToggle:  "Enable",
			wantEnabled: true,
			wantAccent:  colorDisabled,
		},
		{
			// Guild admins can't re-enable a guild disabled by system admins.
			name: "admin-disabled",
			modify: func(v *guildView) {
				v.config.Enabled = false
				v.admin.Enabled = false
			},
			wantText:    []string{"disabled by Spectrograph admins"},
			wantToggle:  "Enable",
			wantEnabled: false,
			wantAccent:  colorBlocked,
		},
		{
			name: "custom-settings",
			modify: func(v *guildView) {
				v.config.RegexMatch = "^Game`s$"
				v.config.DefaultMaxClones = 15
				v.config.ContactEmail = "a_b@example.com"
				v.admin.DefaultMaxClones = 5
				v.settings, _ = effectiveSettings(v.config, v.admin)
			},
			wantText:    []string{"``", "5 (15 set, limited by Spectrograph admins)", `a\_b\@example.com`},
			wantToggle:  "Disable",
			wantEnabled: true,
			wantAccent:  colorEnabled,
		},
		{
			name: "groups",
			modify: func(v *guildView) {
				v.groups = []groupView{
					{name: "Gaming +", channels: 3, occupied: 2},
					{name: "Chill +", category: "Lounge", channels: 1},
				}
				v.unmanaged = 2
				v.missing = discord.PermissionManageChannels
			},
			wantText: []string{
				"Gaming +: 2 clones (2 in use)",
				"Chill + in Lounge: 0 clones (0 in use)",
				"2 channels matching the pattern are over the per-category limit",
				"**Missing permissions:** Spectrograph needs Manage Channels",
			},
			wantToggle:  "Disable",
			wantEnabled: true,
			wantAccent:  colorEnabled,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := testGuildView()
			if tt.modify != nil {
				tt.modify(v)
			}
			components := statusComponents(v, "Settings saved.")

			got := viewText(components)
			for _, want := range append(tt.wantText, "Settings saved.") {
				if !strings.Contains(got, want) {
					t.Errorf("text doesn't contain %q:\n%s", want, got)
				}
			}
			if len(got) > maxTextLength {
				t.Errorf("text is %d characters, want at most %d", len(got), maxTextLength)
			}

			buttons := viewButtons(components)
			toggle := buttons[tt.wantToggle]
			if toggle == nil || toggle.Disabled == tt.wantEnabled {
				t.Errorf("toggle button %q = %+v, want enabled=%v", tt.wantToggle, toggle, tt.wantEnabled)
			}
			if buttons["Edit settings"] == nil || buttons["Refresh"] == nil {
				t.Errorf("buttons = %v, want edit settings and refresh buttons", buttons)
			}
			if b := buttons["Open dashboard"]; b == nil || b.LinkURL() != v.dashboardURL {
				t.Errorf("dashboard button = %+v, want a link to %q", b, v.dashboardURL)
			}

			var container *discord.ContainerComponent
			for _, c := range components {
				if cc, ok := c.(*discord.ContainerComponent); ok {
					container = cc
				}
			}
			if container == nil || container.AccentColor != tt.wantAccent {
				t.Errorf("container = %+v, want accent %v", container, tt.wantAccent)
			}

			// Responses must be ephemeral, and use components (v2).
			resp := message(components...)
			if resp.Data.Flags != discord.EphemeralMessage|discord.IsComponentsV2 {
				t.Errorf("flags = %v, want ephemeral and components v2", resp.Data.Flags)
			}
			marshalDiscord(t, resp)
		})
	}
}

func TestStatusComponentsManyGroups(t *testing.T) {
	t.Parallel()

	v := testGuildView()
	for i := range 40 {
		v.groups = append(v.groups, groupView{name: strings.Repeat("x", 200) + strconv.Itoa(i), category: strings.Repeat("y", 200), channels: 1})
	}

	got := viewText(statusComponents(v, ""))
	if len(got) > maxTextLength {
		t.Errorf("text is %d characters, want at most %d", len(got), maxTextLength)
	}
	if !strings.Contains(got, "more.") {
		t.Errorf("text doesn't mention the groups left out:\n%s", got)
	}
}

func TestNewGroupViews(t *testing.T) {
	t.Parallel()

	settings := defaultSettings()
	settings.maxGroups = 1

	snap := testSnapshot(botPerms, []discord.Channel{
		{ID: testCategory, GuildID: testGuildID, Type: discord.GuildCategory, Name: "Games"},
		inCategory(voice(1, "Gaming +", 1), testCategory),
		inCategory(voice(2, "Gaming +", 2), testCategory),
		inCategory(voice(3, "Other +", 3), testCategory),
		voice(4, "Chill +", 1),
		voice(5, "Not managed", 2),
	}, map[discord.ChannelID]int{1: 2, 4: 1}) //nolint:exhaustive // keyed by channel ID, not an enum.

	groups, unmanaged := newGroupViews(snap, settings)
	want := []groupView{
		{name: "Chill +", channels: 1, occupied: 1},
		{category: "Games", name: "Gaming +", channels: 2, occupied: 1},
	}
	if !slices.Equal(groups, want) || unmanaged != 1 {
		t.Errorf("newGroupViews() = %+v, %d, want %+v, 1", groups, unmanaged, want)
	}
}

func TestEventComponents(t *testing.T) {
	t.Parallel()

	now := time.Unix(1_700_000_000, 0)
	events := []*ent.GuildEvent{
		{
			Type:       guildevent.TypeWARNING,
			Message:    "missing permissions to clone channel",
			CreateTime: now,
			Metadata:   map[string]any{fieldChannelID: "123", fieldMissing: "Manage Channels"},
		},
		{Type: guildevent.TypeINFO, Message: "settings updated via Discord by a_b", CreateTime: now},
	}

	got := viewText(eventComponents("Test", events, ""))
	for _, want := range []string{
		"<t:1700000000:R> **WARNING** missing permissions to clone channel (<#123>, missing: Manage Channels)",
		`settings updated via Discord by a\_b`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("text doesn't contain %q:\n%s", want, got)
		}
	}

	if got = viewText(eventComponents("Test", nil, "")); !strings.Contains(got, "No events") {
		t.Errorf("empty events text = %q", got)
	}

	// Long events are truncated, and the total fits in a message.
	long := make([]*ent.GuildEvent, maxEventCount)
	for i := range long {
		long[i] = &ent.GuildEvent{Type: guildevent.TypeERROR, Message: strings.Repeat("e", 1000), CreateTime: now}
	}
	if got = viewText(eventComponents("Test", long, "")); len(got) > maxTextLength {
		t.Errorf("text is %d characters, want at most %d", len(got), maxTextLength)
	}
}

func TestHelpComponents(t *testing.T) {
	t.Parallel()

	components := helpComponents("https://example.com/dashboard")
	got := viewText(components)
	for _, sub := range []string{subcommandStatus, subcommandSettings, subcommandEnable, subcommandDisable, subcommandEvents} {
		if !strings.Contains(got, "/"+commandName+" "+sub) {
			t.Errorf("help doesn't mention %q", sub)
		}
	}
	if viewButtons(components)["Open dashboard"] == nil {
		t.Error("help has no dashboard button")
	}
}

func TestMarkdownHelpers(t *testing.T) {
	t.Parallel()

	if got := escapeMarkdown("<@123> *hi* `x` _y_ #z"); got != `\<\@123\> \*hi\* \`+"`x\\`"+` \_y\_ \#z` {
		t.Errorf("escapeMarkdown() = %q", got)
	}
	if got := inlineCode(models.DefaultChannelMatch.String()); got != "`^.* \\+$`" {
		t.Errorf("inlineCode() = %q", got)
	}
	if got := inlineCode("a`b"); got != "`` a`b ``" {
		t.Errorf("inlineCode() = %q", got)
	}
	if got := truncate("héllo", 4); got != "…" && len(got) > 4 {
		t.Errorf("truncate() = %q, want at most 4 bytes", got)
	}
	if got := truncate("hello", 10); got != "hello" {
		t.Errorf("truncate() = %q, want unchanged", got)
	}
}
