// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package worker

import (
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/lrstanley/spectrograph/internal/database/ent"
	"github.com/lrstanley/spectrograph/internal/models"
)

// Views (the messages and modals sent in response to interactions) use
// Discord's components (v2): messages are laid out with containers and text
// displays instead of content and embeds, and modal inputs are wrapped in
// labels. All messages are ephemeral (only shown to the user who interacted).
const viewFlags = discord.EphemeralMessage | discord.IsComponentsV2

// Accent colors of the status view's container, by status.
const (
	colorEnabled  discord.Color = 0x57F287
	colorDisabled discord.Color = 0x99AAB5
	colorBlocked  discord.Color = 0xED4245
)

const (
	// maxTextLength is the max length of the text across all of a message's
	// text displays.
	maxTextLength = 4000

	// maxGroupLines is the max number of managed channel groups listed in the
	// status view.
	maxGroupLines = 15

	// maxGroupsLength is the max length of the list of managed channel
	// groups in the status view (leaving room for the rest of the view).
	maxGroupsLength = 2500

	// maxEventLength is the max length of a single event in the events view.
	maxEventLength = 300

	// maxCloneLimit is the max value of the guild config's default_max_clones.
	maxCloneLimit = 20
)

// guildView is the data shown in the status view.
type guildView struct {
	name   string
	config *ent.GuildConfig
	admin  *ent.GuildAdminConfig

	// settings are the effective settings, and regexErr is set if the
	// configured regex is invalid (and the default is used instead).
	settings guildSettings
	regexErr error

	// dashboardURL is the guild's page on the dashboard, if known.
	dashboardURL string

	// groups are the managed channel groups, and unmanaged the number of
	// matching groups which aren't managed (over the per-category limit).
	groups    []groupView
	unmanaged int

	// missing are the guild-level permissions the bot needs, but lacks.
	missing discord.Permissions
}

// groupView is a managed group of voice channels (an original channel and its
// clones).
type groupView struct {
	category string
	name     string
	channels int
	occupied int
}

// newGroupViews summarizes the managed channel groups of a guild.
func newGroupViews(snap *snapshot, settings guildSettings) (groups []groupView, unmanaged int) {
	names := make(map[discord.ChannelID]string, len(snap.channels))
	for _, ch := range snap.channels {
		if ch.Type == discord.GuildCategory {
			names[ch.ID] = ch.Name
		}
	}
	users := occupancy(snap.voiceStates)

	parents, byParent := groupChannels(snap.channels, settings.match, snap.guild.AFKChannelID)
	for _, parent := range parents {
		gs := byParent[parent]
		if len(gs) > settings.maxGroups {
			unmanaged += len(gs) - settings.maxGroups
			gs = gs[:settings.maxGroups]
		}
		for _, g := range gs {
			gv := groupView{category: names[parent], name: g.primary().Name, channels: len(g)}
			for _, ch := range g {
				if users[ch.ID] > 0 {
					gv.occupied++
				}
			}
			groups = append(groups, gv)
		}
	}
	return groups, unmanaged
}

// markdownEscaper escapes Discord markdown, mentions, and emoji/timestamp
// syntax in user-provided text (e.g. channel names).
var markdownEscaper = strings.NewReplacer(
	`\`, `\\`, "*", `\*`, "_", `\_`, "~", `\~`, "`", "\\`", "|", `\|`, ">", `\>`, "<", `\<`,
	"#", `\#`, "-", `\-`, "[", `\[`, "]", `\]`, "(", `\(`, ")", `\)`, ":", `\:`, "@", `\@`,
)

func escapeMarkdown(s string) string {
	return markdownEscaper.Replace(s)
}

// inlineCode formats s as inline code. Code spans can't escape backticks, so
// spans with backticks use double backticks (which allow single ones).
func inlineCode(s string) string {
	if strings.Contains(s, "`") {
		return "`` " + strings.ReplaceAll(s, "``", "`\u200b`") + " ``"
	}
	return "`" + s + "`"
}

// truncate shortens s to at most n bytes (on a rune boundary), with an
// ellipsis if it was shortened.
func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	const ellipsis = "…"
	s = s[:n-len(ellipsis)]
	for !utf8.ValidString(s) {
		s = s[:len(s)-1]
	}
	return s + ellipsis
}

func plural(n int, singular, pluralForm string) string {
	if n == 1 {
		return "1 " + singular
	}
	return strconv.Itoa(n) + " " + pluralForm
}

func text(content string) *discord.TextDisplayComponent {
	return &discord.TextDisplayComponent{Content: content}
}

func separator() *discord.SeparatorComponent {
	return &discord.SeparatorComponent{Divider: true, Spacing: discord.SeparatorComponentSpacingSmallPadding}
}

// message returns an ephemeral message response with the given components.
func message(components ...discord.Component) *api.InteractionResponse {
	return &api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: &api.InteractionResponseData{
			Components:      discord.ComponentsPtr(components...),
			Flags:           viewFlags,
			AllowedMentions: &api.AllowedMentions{},
		},
	}
}

// update returns a response which updates the message the component
// interaction is attached to.
func update(components ...discord.Component) *api.InteractionResponse {
	resp := message(components...)
	resp.Type = api.UpdateMessage
	return resp
}

// errorMessage returns an ephemeral message with the given (markdown) text,
// and optionally buttons.
func errorMessage(msg string, buttons ...discord.InteractiveComponent) *api.InteractionResponse {
	components := []discord.Component{&discord.ContainerComponent{
		AccentColor: colorBlocked,
		Components:  []discord.Component{text(msg)},
	}}
	if len(buttons) > 0 {
		row := discord.ActionRowComponent(buttons)
		components = append(components, &row)
	}
	return message(components...)
}

func editSettingsButton() *discord.ButtonComponent {
	return &discord.ButtonComponent{Style: discord.PrimaryButtonStyle(), Label: "Edit settings", CustomID: customIDEditSettings}
}

func dashboardButton(url string) *discord.ButtonComponent {
	return &discord.ButtonComponent{Style: discord.LinkButtonStyle(url), Label: "Open dashboard"}
}

// statusComponents returns the status view of a guild: its status, settings,
// managed channels, and buttons to change them. notice, if set, is shown
// above it (e.g. the outcome of a change).
func statusComponents(v *guildView, notice string) []discord.Component {
	var status string
	accent := colorEnabled
	switch {
	case !v.admin.Enabled:
		status = "**Status:** disabled by Spectrograph admins, so channels are left as they are. " +
			"Contact the Spectrograph admins if you think this is a mistake."
		accent = colorBlocked
	case !v.config.Enabled:
		status = "**Status:** disabled, so channels are left as they are."
		accent = colorDisabled
	default:
		status = "**Status:** enabled, managing voice channels."
	}

	pattern := inlineCode(v.settings.match.String())
	switch {
	case v.regexErr != nil:
		pattern = inlineCode(v.config.RegexMatch) + " is invalid, so the default " + pattern + " is used"
	case v.config.RegexMatch == "":
		pattern += " (default)"
	}

	clones := strconv.Itoa(v.settings.maxClones)
	switch {
	case v.config.DefaultMaxClones > 0 && v.settings.maxClones < v.config.DefaultMaxClones:
		clones += fmt.Sprintf(" (%d set, limited by Spectrograph admins)", v.config.DefaultMaxClones)
	case v.config.DefaultMaxClones == 0:
		clones += " (default)"
	}

	email := "not set"
	if v.config.ContactEmail != "" {
		email = escapeMarkdown(v.config.ContactEmail)
	}

	settings := fmt.Sprintf(
		"**Channel name pattern:** %s\n**Max clones per channel:** %s\n**Max managed channels per category:** %d\n**Contact email:** %s",
		pattern, clones, v.settings.maxGroups, email,
	)

	var channels strings.Builder
	channels.WriteString("**Managed voice channels**\n")
	if len(v.groups) == 0 {
		channels.WriteString("No voice channels match the pattern. Rename a voice channel to match it, " +
			"e.g. " + inlineCode("Gaming +") + " for the default pattern.")
	}
	for i, g := range v.groups {
		line := "- " + escapeMarkdown(truncate(g.name, 100))
		if g.category != "" {
			line += " in " + escapeMarkdown(truncate(g.category, 100))
		}
		line += fmt.Sprintf(": %s (%d in use)\n", plural(g.channels-1, "clone", "clones"), g.occupied)

		if i == maxGroupLines || channels.Len()+len(line) > maxGroupsLength {
			fmt.Fprintf(&channels, "…and %d more.\n", len(v.groups)-i)
			break
		}
		channels.WriteString(line)
	}
	if v.unmanaged > 0 {
		fmt.Fprintf(&channels, "%s matching the pattern %s over the per-category limit, so not managed.",
			plural(v.unmanaged, "channel", "channels"), map[bool]string{true: "is", false: "are"}[v.unmanaged == 1])
	}

	inner := []discord.Component{
		text("## Spectrograph in " + escapeMarkdown(truncate(v.name, 100)) + "\n" + status),
		separator(),
		text(settings),
		separator(),
		text(strings.TrimSpace(channels.String())),
	}
	if v.missing != 0 {
		inner = append(inner, separator(), text(fmt.Sprintf(
			"**Missing permissions:** Spectrograph needs %s in this server to manage voice channels. "+
				"See %s for other problems.",
			permissionString(v.missing), inlineCode("/"+commandName+" "+subcommandEvents),
		)))
	}

	toggle := &discord.ButtonComponent{Style: discord.SecondaryButtonStyle(), Label: "Disable", CustomID: customIDDisable}
	if !v.config.Enabled {
		toggle = &discord.ButtonComponent{
			Style:    discord.SuccessButtonStyle(),
			Label:    "Enable",
			CustomID: customIDEnable,
			Disabled: !v.admin.Enabled,
		}
	}
	row := discord.ActionRowComponent{
		toggle,
		editSettingsButton(),
		&discord.ButtonComponent{Style: discord.SecondaryButtonStyle(), Label: "Refresh", CustomID: customIDRefresh},
	}
	if v.dashboardURL != "" {
		row = append(row, dashboardButton(v.dashboardURL))
	}

	var components []discord.Component
	if notice != "" {
		components = append(components, text(notice))
	}
	return append(components, &discord.ContainerComponent{AccentColor: accent, Components: inner}, &row)
}

// settingsModal returns the modal to edit the guild's settings, prefilled with
// the current ones.
func settingsModal(cfg *ent.GuildConfig, admin *ent.GuildAdminConfig) *api.InteractionResponse {
	options := make([]discord.SelectOption, 0, maxCloneLimit+1)
	options = append(options, discord.SelectOption{
		Label:   fmt.Sprintf("Default (%d)", models.DefaultMaxClones),
		Value:   "0",
		Default: cfg.DefaultMaxClones == 0,
	})
	for i := 1; i <= maxCloneLimit; i++ {
		options = append(options, discord.SelectOption{
			Label:   strconv.Itoa(i),
			Value:   strconv.Itoa(i),
			Default: cfg.DefaultMaxClones == i,
		})
	}

	clonesDescription := "How many clones of a voice channel can exist at once."
	if admin.DefaultMaxClones > 0 {
		clonesDescription = fmt.Sprintf("How many clones of a voice channel can exist at once (at most %d).", admin.DefaultMaxClones)
	}

	return &api.InteractionResponse{
		Type: api.ModalResponse,
		Data: &api.InteractionResponseData{
			CustomID: option.NewNullableString(string(customIDSettingsModal)),
			Title:    option.NewNullableString("Spectrograph settings"),
			Components: discord.ComponentsPtr(
				&discord.LabelComponent{
					Label:       "Channel name pattern",
					Description: "A regular expression matching the voice channels to manage. Leave empty for the default.",
					Component: &discord.TextInputComponent{
						CustomID:     customIDRegexMatch,
						Style:        discord.TextInputShortStyle,
						LengthLimits: [2]int{0, 100},
						Value:        cfg.RegexMatch,
						Placeholder:  models.DefaultChannelMatch.String(),
					},
				},
				&discord.LabelComponent{
					Label:       "Max clones per channel",
					Description: clonesDescription,
					Component: &discord.StringSelectComponent{
						CustomID:    customIDMaxClones,
						Options:     options,
						ValueLimits: [2]int{1, 1},
						Required:    true,
					},
				},
				&discord.LabelComponent{
					Label:       "Contact email",
					Description: "Optional, for Spectrograph admins to reach you about this server.",
					Component: &discord.TextInputComponent{
						CustomID:     customIDContactEmail,
						Style:        discord.TextInputShortStyle,
						LengthLimits: [2]int{0, 255},
						Value:        cfg.ContactEmail,
						Placeholder:  "name@example.com",
					},
				},
			),
		},
	}
}

// modalValues returns the values of the inputs of a submitted modal, by custom
// ID: the text of text inputs, and the first selected value of selects.
func modalValues(components discord.TopLevelComponents) map[discord.ComponentID]string {
	values := make(map[discord.ComponentID]string)

	var visit func(c discord.Component)
	visit = func(c discord.Component) {
		switch c := c.(type) {
		case *discord.LabelComponent:
			visit(c.Component)
		case *discord.ActionRowComponent:
			for _, child := range *c {
				visit(child)
			}
		case *discord.TextInputComponent:
			values[c.CustomID] = c.Value
		case *discord.StringSelectComponent:
			if len(c.Values) > 0 {
				values[c.CustomID] = c.Values[0]
			}
		}
	}

	for _, c := range components {
		visit(c)
	}
	return values
}

// eventComponents returns the view of the most recent events of a guild
// (newest first).
func eventComponents(guildName string, events []*ent.GuildEvent, dashboardURL string) []discord.Component {
	header := "## Recent events in " + escapeMarkdown(truncate(guildName, 100))
	if len(events) == 0 {
		header += "\nNo events in the last week."
	}

	budget := maxTextLength - len(header) - 100
	var lines []string
	for _, ev := range events {
		line := truncate(formatEvent(ev), maxEventLength)
		if budget -= len(line) + 1; budget < 0 {
			break
		}
		lines = append(lines, line)
	}

	inner := []discord.Component{text(header)}
	if len(lines) > 0 {
		inner = append(inner, separator(), text(strings.Join(lines, "\n")))
	}

	components := []discord.Component{&discord.ContainerComponent{AccentColor: colorDisabled, Components: inner}}
	if dashboardURL != "" {
		row := discord.ActionRowComponent{dashboardButton(dashboardURL)}
		components = append(components, &row)
	}
	return components
}

// eventChannelFields are event fields holding a channel ID, shown as channel
// mentions.
var eventChannelFields = []string{fieldChannelID, fieldPrimaryID, "category_id"}

// formatEvent formats a guild event as a line of markdown, with a relative
// timestamp, the event type, message, and fields.
func formatEvent(ev *ent.GuildEvent) string {
	line := fmt.Sprintf("<t:%d:R> **%s** %s", ev.CreateTime.Unix(), ev.Type, escapeMarkdown(ev.Message))

	details := make([]string, 0, len(ev.Metadata))
	for _, k := range slices.Sorted(maps.Keys(ev.Metadata)) {
		v := fmt.Sprint(ev.Metadata[k])
		if slices.Contains(eventChannelFields, k) {
			if sf, err := discord.ParseSnowflake(v); err == nil && sf.IsValid() {
				details = append(details, "<#"+sf.String()+">")
				continue
			}
		}
		details = append(details, escapeMarkdown(k)+": "+escapeMarkdown(truncate(v, 120)))
	}
	if len(details) > 0 {
		line += " (" + strings.Join(details, ", ") + ")"
	}
	return line
}

// helpComponents returns the help view.
func helpComponents(dashboardURL string) []discord.Component {
	cmd := func(sub string) string { return inlineCode("/" + commandName + " " + sub) }

	inner := []discord.Component{
		text("## How Spectrograph works\n" +
			"Spectrograph manages voice channels whose name matches a pattern (by default, names ending in " +
			inlineCode(" +") + ", e.g. " + inlineCode("Gaming +") + "). When every channel with that name is in use, " +
			"it clones the original channel (with its permissions, bitrate, and user limit), and it removes extra " +
			"empty clones, so there's always exactly one free channel."),
		separator(),
		text("## Commands\n" +
			"- " + cmd(subcommandStatus) + ": status, settings, and managed channels, with buttons to change them.\n" +
			"- " + cmd(subcommandSettings) + ": edit the channel name pattern, max clones, and contact email.\n" +
			"- " + cmd(subcommandEnable) + " / " + cmd(subcommandDisable) + ": start or stop managing channels.\n" +
			"- " + cmd(subcommandEvents) + ": recent events, such as missing permissions.\n" +
			"Only members with **Manage Server** (or Administrator) can use them, and the same settings can be " +
			"changed on the dashboard."),
		separator(),
		text("## Permissions\n" +
			"Spectrograph needs **View Channel** and **Manage Channels**. Clones copy the original channel's " +
			"permission overwrites, but Discord only lets the bot set permissions it has itself: allowed " +
			"permissions it can't set are left out of clones, and channels denying permissions it can't set aren't " +
			"cloned. Give it **Manage Roles** to keep existing clones' permissions in sync."),
	}

	components := []discord.Component{&discord.ContainerComponent{AccentColor: colorEnabled, Components: inner}}
	if dashboardURL != "" {
		row := discord.ActionRowComponent{dashboardButton(dashboardURL)}
		components = append(components, &row)
	}
	return components
}
