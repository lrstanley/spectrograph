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

func FetchUser(client *http.Client, token *oauth2.Token) (user *models.UserAuthDiscord, servers []*models.UserDiscordServer, err error) {
	user = &models.UserAuthDiscord{}
	servers = []*models.UserDiscordServer{}

	// Fetch user details.
	_, err = handleRequest(client, token, "GET", UserEndpoint, nil, &user)
	if err != nil {
		return nil, nil, fmt.Errorf("error fetching user info: %w", err)
	}

	user.LastLogin = time.Now()
	user.AccessToken = token.AccessToken
	user.RefreshToken = token.RefreshToken
	user.ExpiresAt = token.Expiry

	// Properly parse out the discord avatar. If they don't have an avatar,
	// use the default avatar endpoint.
	if user.Avatar == "" {
		discriminator, _ := strconv.Atoi(user.Discriminator)
		user.AvatarURL = fmt.Sprintf(DefaultAvatarEndpoint, discriminator%5)
	} else {
		extension := "png"
		if len(user.Avatar) >= len(GIFAvatarPrefix) &&
			user.Avatar[0:len(GIFAvatarPrefix)] == GIFAvatarPrefix {
			extension = "gif"
		}

		user.AvatarURL = fmt.Sprintf(AvatarEndpoint, user.ID, user.Avatar, extension)
	}

	// Fetch guild details.
	rawServers := []*models.UserDiscordServer{}
	_, err = handleRequest(client, token, "GET", GuildsEndpoint, nil, &rawServers)
	if err != nil {
		return nil, nil, fmt.Errorf("error fetching guild info: %w", err)
	}

	for i := range rawServers {
		// Check if they have the admin permission bit.
		rawServers[i].Admin = rawServers[i].Permissions.Contains(models.DiscordPermAdministrator)

		if !rawServers[i].Owner && !rawServers[i].Admin {
			// Ignore servers that they're not an owner of.
			continue
		}

		rawServers[i].IconURL = GenerateGuildIconURL(rawServers[i].ID, rawServers[i].Icon)
		servers = append(servers, rawServers[i])
	}

	sort.SliceStable(servers, func(i, j int) bool {
		return strings.ToLower(servers[i].Name) < strings.ToLower(servers[j].Name)
	})

	return user, servers, nil
}

// GenerateGuildIconURL generates a discord server icon url from an icon hash.
func GenerateGuildIconURL(guildID, iconHash string) string {
	if guildID == "" || iconHash == "" {
		return ""
	}

	extension := "png"

	if len(iconHash) >= len(GIFAvatarPrefix) &&
		iconHash[0:len(GIFAvatarPrefix)] == GIFAvatarPrefix {
		extension = "gif"
	}

	return fmt.Sprintf(
		ServerIconEndpoint,
		guildID,
		iconHash,
		extension,
	)
}
