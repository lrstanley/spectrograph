// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package discordapi

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/lrstanley/spectrograph/internal/models"
)

const (
	UserEndpoint   = APIBase + "/users/@me"
	GuildsEndpoint = APIBase + "/users/@me/guilds?limit=200"

	// GIFAvatarPrefix and the below endpoints are documented here:
	// https://discord.com/developers/docs/reference#image-formatting-cdn-endpoints
	GIFAvatarPrefix       = "a_"
	AvatarEndpoint        = "https://cdn.discordapp.com/avatars/%s/%s.%s"     // User id, avatar id, extension.
	DefaultAvatarEndpoint = "https://cdn.discordapp.com/embed/avatars/%d.png" // See [GenerateUserAvatarURL].
	ServerIconEndpoint    = "https://cdn.discordapp.com/icons/%s/%s.%s"       // Guild id, icon id, extension.
)

// FetchUser fetches the current discord user that is authenticated with the token,
// as well as any guilds they are a part of. Make sure to check that the user
// has the administrator permission in the guild, or is an owner of the guild.
func FetchUser(ctx context.Context, token string) (user *models.DiscordUser, guilds []*models.UserGuildResponse, err error) {
	user = &models.DiscordUser{}
	if err = request(ctx, token, http.MethodGet, UserEndpoint, user); err != nil {
		return nil, nil, fmt.Errorf("fetching user info: %w", err)
	}

	if err = request(ctx, token, http.MethodGet, GuildsEndpoint, &guilds); err != nil {
		return nil, nil, fmt.Errorf("fetching guild info: %w", err)
	}

	return user, guilds, nil
}

// UserCanManageGuild returns true if the user is the owner of the guild, or has
// the administrator or manage server permission in the guild.
func UserCanManageGuild(guild *models.UserGuildResponse) bool {
	if guild == nil {
		return false
	}
	return guild.Owner ||
		guild.Permissions.Contains(models.DiscordPermAdministrator) ||
		guild.Permissions.Contains(models.DiscordPermManageServer)
}

func cdnExtension(hash string) string {
	if strings.HasPrefix(hash, GIFAvatarPrefix) {
		return "gif"
	}
	return "png"
}

// GenerateUserAvatarURL returns the users avatar URL. If they don't have an
// avatar, the default avatar endpoint is used. See:
// https://discord.com/developers/docs/reference#image-formatting-cdn-endpoints
func GenerateUserAvatarURL(user *models.DiscordUser) string {
	if user == nil {
		return ""
	}

	if user.Avatar == "" {
		// Legacy users (with a discriminator) use discriminator % 5, migrated users
		// (no discriminator, or "0") use (user_id >> 22) % 6.
		if discriminator, _ := strconv.Atoi(user.Discriminator); discriminator > 0 {
			return fmt.Sprintf(DefaultAvatarEndpoint, discriminator%5)
		}
		id, _ := strconv.ParseUint(user.ID, 10, 64)
		return fmt.Sprintf(DefaultAvatarEndpoint, (id>>22)%6)
	}

	return fmt.Sprintf(AvatarEndpoint, user.ID, user.Avatar, cdnExtension(user.Avatar))
}

// GenerateGuildIconURL generates a discord server icon url from an icon hash.
func GenerateGuildIconURL(id, icon string) string {
	if id == "" || icon == "" {
		return ""
	}
	return fmt.Sprintf(ServerIconEndpoint, id, icon, cdnExtension(icon))
}
