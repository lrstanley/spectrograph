// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package discordapi

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/andersfylling/disgord"
	"github.com/lrstanley/spectrograph/internal/models"
)

const (
	UserEndpoint   = "https://discord.com/api/users/@me"
	GuildsEndpoint = "https://discord.com/api/users/@me/guilds?limit=200"

	// https://discord.com/developers/docs/reference#image-formatting-cdn-endpoints
	GIFAvatarPrefix       = "a_"
	AvatarEndpoint        = "https://cdn.discordapp.com/avatars/%s/%s.%s"     // user id, avatar id, extension.
	DefaultAvatarEndpoint = "https://cdn.discordapp.com/embed/avatars/%d.png" // user-discriminator modulo 5 (Test#1337 % 5 == 2).
	ServerIconEndpoint    = "https://cdn.discordapp.com/icons/%s/%s.%s"       // guild id, icon id, extension.
)

// FetchUser fetches the current discord user that is authenticated with the token,
// as well as any guilds that are a part of. Make sure to check that the user
// has the administrator permission in the guild, or is an owner of the guild.
func FetchUser(ctx context.Context, token string) (user *disgord.User, guilds []*models.UserGuildResponse, err error) {
	user = &disgord.User{}
	guilds = []*models.UserGuildResponse{}

	// Fetch user details.
	_, err = handleRequest(ctx, token, "GET", UserEndpoint, http.NoBody, user)
	if err != nil {
		return nil, nil, fmt.Errorf("error fetching user info: %w", err)
	}

	// Fetch guild details.
	_, err = handleRequest(ctx, token, "GET", GuildsEndpoint, nil, &guilds)
	if err != nil {
		return nil, nil, fmt.Errorf("error fetching guild info: %w", err)
	}

	return user, guilds, nil
}

// UserHasAdmin checks if a user has the administrator permission in a guild.
func UserHasAdmin(user *disgord.User, guild *models.UserGuildResponse) bool {
	if user == nil {
		return false
	}

	if guild.Permissions.Contains(models.DiscordPermAdministrator) || guild.Permissions.Contains(models.DiscordPermManageServer) {
		return true
	}

	return false
}

// GenerateUserAvatarURL parses out the discord avatar. If they don't have an avatar,
// use the default avatar endpoint.
func GenerateUserAvatarURL(user *disgord.User) string {
	if user == nil {
		return ""
	}

	if user.Avatar == "" {
		discriminator, _ := strconv.Atoi(user.Discriminator.String())
		return fmt.Sprintf(DefaultAvatarEndpoint, discriminator%5)
	}

	extension := "png"
	if len(user.Avatar) >= len(GIFAvatarPrefix) &&
		user.Avatar[0:len(GIFAvatarPrefix)] == GIFAvatarPrefix {
		extension = "gif"
	}

	return fmt.Sprintf(AvatarEndpoint, user.ID, user.Avatar, extension)
}

// GenerateGuildIconURL generates a discord server icon url from an icon hash.
func GenerateGuildIconURL(id, icon string) string {
	if id == "" || icon == "" {
		return ""
	}

	extension := "png"

	if len(icon) >= len(GIFAvatarPrefix) &&
		icon[0:len(GIFAvatarPrefix)] == GIFAvatarPrefix {
		extension = "gif"
	}

	return fmt.Sprintf(
		ServerIconEndpoint,
		id,
		icon,
		extension,
	)
}
