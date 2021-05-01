// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package discordapi

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/lrstanley/spectrograph/internal/models"
	"golang.org/x/oauth2"
)

const (
	UserEndpoint   = "https://discord.com/api/users/@me"
	GuildsEndpoint = "https://discord.com/api/users/@me/guilds"

	// https://discord.com/developers/docs/reference#image-formatting-cdn-endpoints
	GIFAvatarPrefix       = "a_"
	AvatarEndpoint        = "https://cdn.discordapp.com/avatars/%s/%s.%s"     // user id, avatar id, extension.
	DefaultAvatarEndpoint = "https://cdn.discordapp.com/embed/avatars/%d.png" // user-discriminator modulo 5 (Test#1337 % 5 == 2).
	ServerIconEndpoint    = "https://cdn.discordapp.com/icons/%s/%s.%s"       // guild id, icon id, extension.
)

func FetchUser(client *http.Client, token *oauth2.Token) (user *models.User, err error) {
	user = &models.User{}

	// Fetch user details.
	_, err = handleRequest(client, token, "GET", UserEndpoint, nil, &user.Discord)
	if err != nil {
		return nil, fmt.Errorf("error fetching user info: %w", err)
	}

	user.AccountUpdated = time.Now()
	user.Discord.LastLogin = time.Now()
	user.Discord.AccessToken = token.AccessToken
	user.Discord.RefreshToken = token.RefreshToken
	user.Discord.ExpiresAt = token.Expiry

	// Properly parse out the discord avatar. If they don't have an avatar,
	// use the default avatar endpoint.
	if user.Discord.Avatar == "" {
		discriminator, _ := strconv.Atoi(user.Discord.Discriminator)
		user.Discord.AvatarURL = fmt.Sprintf(DefaultAvatarEndpoint, discriminator%5)
	} else {
		extension := "png"
		if len(user.Discord.Avatar) >= len(GIFAvatarPrefix) &&
			user.Discord.Avatar[0:len(GIFAvatarPrefix)] == GIFAvatarPrefix {
			extension = "gif"
		}

		user.Discord.AvatarURL = fmt.Sprintf(AvatarEndpoint, user.Discord.ID, user.Discord.Avatar, extension)
	}

	// Fetch guild details.
	servers := []models.UserDiscordServer{}
	_, err = handleRequest(client, token, "GET", GuildsEndpoint, nil, &servers)
	if err != nil {
		return nil, fmt.Errorf("error fetching guild info: %w", err)
	}

	var extension string
	for i := range servers {
		// Check if they have the admin permission bit.
		servers[i].Admin = servers[i].Permissions.Contains(models.DiscordPermAdministrator)

		if !servers[i].Owner && !servers[i].Admin {
			// Ignore servers that they're not an owner of.
			continue
		}

		// TODO: default server icon?
		if servers[i].Icon != "" {
			extension = "png"

			if len(servers[i].Icon) >= len(GIFAvatarPrefix) &&
				servers[i].Icon[0:len(GIFAvatarPrefix)] == GIFAvatarPrefix {
				extension = "gif"
			}

			servers[i].IconURL = fmt.Sprintf(
				ServerIconEndpoint,
				servers[i].ID,
				servers[i].Icon,
				extension,
			)
		}

		user.DiscordServers = append(user.DiscordServers, servers[i])
	}

	sort.SliceStable(user.DiscordServers, func(i, j int) bool {
		return strings.ToLower(user.DiscordServers[i].Name) < strings.ToLower(user.DiscordServers[j].Name)
	})

	return user, nil
}
