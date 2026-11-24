// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package worker

import (
	"context"
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
)

// commandName is the name of the slash command, which has a subcommand per
// action (e.g. "/spectrograph status").
const commandName = "spectrograph"

// Subcommands of [commandName].
const (
	subcommandStatus   = "status"
	subcommandSettings = "settings"
	subcommandEnable   = "enable"
	subcommandDisable  = "disable"
	subcommandEvents   = "events"
	subcommandHelp     = "help"
)

const (
	// optionEventCount is the option of the events subcommand for the number
	// of events to show.
	optionEventCount = "count"

	defaultEventCount = 10
	maxEventCount     = 25
)

// customIDPrefix prefixes the custom IDs of the bot's message components and
// modals, so interactions for other apps' components are ignored.
const customIDPrefix = commandName + ":"

// Custom IDs of message components and modals.
const (
	customIDEnable        discord.ComponentID = customIDPrefix + "enable"
	customIDDisable       discord.ComponentID = customIDPrefix + "disable"
	customIDEditSettings  discord.ComponentID = customIDPrefix + "edit-settings"
	customIDRefresh       discord.ComponentID = customIDPrefix + "refresh"
	customIDSettingsModal discord.ComponentID = customIDPrefix + "settings"

	// Inputs of the settings modal.
	customIDRegexMatch   discord.ComponentID = customIDPrefix + "regex-match"
	customIDMaxClones    discord.ComponentID = customIDPrefix + "max-clones"
	customIDContactEmail discord.ComponentID = customIDPrefix + "contact-email"
)

// commandPermissions are the default permissions needed to see and use the
// command, matching who the dashboard considers a guild admin (the owner,
// administrators, and members with Manage Server). Guild admins can change
// who can use the command (in Server Settings > Integrations), so permissions
// are also checked when handling each interaction (see [authorizeMember]).
const commandPermissions = discord.PermissionManageGuild

// commandDefinitions returns the application commands of the bot.
func commandDefinitions() []api.CreateCommandData {
	perms := commandPermissions
	return []api.CreateCommandData{{
		Name:                     commandName,
		Description:              "Manage Spectrograph, which clones voice channels as they fill up",
		Type:                     discord.ChatInputCommand,
		DefaultMemberPermissions: &perms,
		NoDMPermission:           true,
		Options: discord.CommandOptions{
			&discord.SubcommandOption{
				OptionName:  subcommandStatus,
				Description: "Show Spectrograph's status, settings, and managed channels in this server",
			},
			&discord.SubcommandOption{
				OptionName:  subcommandSettings,
				Description: "Edit the channel name pattern, max clones, and contact email",
			},
			&discord.SubcommandOption{
				OptionName:  subcommandEnable,
				Description: "Start managing voice channels in this server",
			},
			&discord.SubcommandOption{
				OptionName:  subcommandDisable,
				Description: "Stop managing voice channels in this server (channels are left as they are)",
			},
			&discord.SubcommandOption{
				OptionName:  subcommandEvents,
				Description: "Show recent events, such as problems managing channels",
				Options: []discord.CommandOptionValue{
					&discord.IntegerOption{
						OptionName:  optionEventCount,
						Description: fmt.Sprintf("Number of events to show (default %d)", defaultEventCount),
						Min:         option.NewInt(1),
						Max:         option.NewInt(maxEventCount),
					},
				},
			},
			&discord.SubcommandOption{
				OptionName:  subcommandHelp,
				Description: "Explain how Spectrograph works, and how to set it up",
			},
		},
	}}
}

// registerCommands registers the application commands with Discord: in
// guildID only if it is valid (which updates instantly, for development),
// otherwise globally. Commands are bulk overwritten, which is idempotent:
// unchanged commands are left as they are. Transient errors are retried until
// ctx is cancelled.
func (w *Worker) registerCommands(ctx context.Context, appID discord.AppID, guildID discord.GuildID) {
	backoff := initialBackoff
	for {
		var cmds []discord.Command
		var err error
		if guildID.IsValid() {
			cmds, err = w.client.WithContext(ctx).BulkOverwriteGuildCommands(appID, guildID, commandDefinitions())
		} else {
			cmds, err = w.client.WithContext(ctx).BulkOverwriteCommands(appID, commandDefinitions())
		}
		if err == nil {
			names := make([]string, 0, len(cmds))
			for _, c := range cmds {
				names = append(names, "/"+c.Name)
			}
			w.logger.InfoContext(ctx, "registered application commands",
				"commands", strings.Join(names, ", "),
				"guild_id", guildID.String(), // Empty for global commands.
			)
			return
		}

		if ctx.Err() != nil {
			return
		}
		if !isTransient(err) {
			w.logger.ErrorContext(ctx, "unable to register application commands", "guild_id", guildID.String(), "error", err)
			return
		}

		w.logger.WarnContext(ctx, "unable to register application commands, retrying", "error", err, "backoff", backoff)
		if !sleep(ctx, backoff) {
			return
		}
		backoff = min(backoff*2, maxBackoff)
	}
}
