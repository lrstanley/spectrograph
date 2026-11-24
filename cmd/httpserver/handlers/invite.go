// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package handlers

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/lrstanley/spectrograph/internal/models"
)

const discordBotAuthEndpoint = "https://discord.com/oauth2/authorize"

// botPermissions are the permissions requested when inviting the bot. The bot
// never joins voice channels or moves members, it only creates, reorders,
// edits, and deletes voice channels:
//   - View Channel and Manage Channels: required to see and manage channels.
//   - Connect, Speak, Video, and Use Voice Activity: clones copy the permission
//     overwrites of the original channel, and Discord only lets the bot allow or
//     deny permissions it has itself. These are the voice permissions commonly
//     used in overwrites (e.g. denying Connect to @everyone for a private
//     channel). Other allowed permissions are left out of clones, and the worker
//     reports them (and refuses to clone channels denying them).
//
// Manage Roles is intentionally not requested, as it's broad. Without it,
// clones are still created with the original's overwrites, but changes to the
// original's overwrites aren't synced to existing clones (the worker reports
// this, and admins can grant it if they want this).
const botPermissions = models.DiscordPermViewChannel |
	models.DiscordPermManageChannels |
	models.DiscordPermVoiceConnect |
	models.DiscordPermVoiceSpeak |
	models.DiscordPermVoiceStreamVideo |
	models.DiscordPermVoiceUseVAD

// botScopes are the OAuth2 scopes requested when inviting the bot: the bot
// itself, and its slash commands. Discord grants applications.commands along
// with bot nowadays, but requesting it explicitly is harmless, and covers
// older behavior.
const botScopes = "bot applications.commands"

// getBotInvite redirects to the Discord bot authorization flow, optionally
// pre-selecting a guild. See:
// https://discord.com/developers/docs/topics/oauth2#bot-authorization-flow
func (h *handler) getBotInvite(w http.ResponseWriter, r *http.Request) {
	query := url.Values{
		"client_id":   {h.cfg.DiscordClientID},
		"scope":       {botScopes},
		"permissions": {strconv.FormatUint(botPermissions, 10)},
	}

	if guildID := chi.URLParam(r, "guildSnowflake"); guildID != "" {
		query.Set("guild_id", guildID)
		query.Set("disable_guild_select", "true")
	}

	http.Redirect(w, r, discordBotAuthEndpoint+"?"+query.Encode(), http.StatusTemporaryRedirect)
}
