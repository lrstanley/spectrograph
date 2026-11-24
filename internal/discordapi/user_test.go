// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package discordapi

import (
	"testing"

	"github.com/lrstanley/spectrograph/internal/models"
)

func TestGenerateUserAvatarURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		user *models.DiscordUser
		want string
	}{
		{name: "nil", user: nil, want: ""},
		{
			name: "legacy-discriminator",
			user: &models.DiscordUser{ID: "80351110224678912", Discriminator: "1337"},
			want: "https://cdn.discordapp.com/embed/avatars/2.png",
		},
		{
			name: "migrated-zero-discriminator",
			user: &models.DiscordUser{ID: "80351110224678912", Discriminator: "0"},
			want: "https://cdn.discordapp.com/embed/avatars/5.png", // (80351110224678912 >> 22) % 6.
		},
		{
			name: "migrated-no-discriminator",
			user: &models.DiscordUser{ID: "80351110224678912"},
			want: "https://cdn.discordapp.com/embed/avatars/5.png",
		},
		{
			name: "static-avatar",
			user: &models.DiscordUser{ID: "1", Avatar: "abc"},
			want: "https://cdn.discordapp.com/avatars/1/abc.png",
		},
		{
			name: "animated-avatar",
			user: &models.DiscordUser{ID: "1", Avatar: "a_abc"},
			want: "https://cdn.discordapp.com/avatars/1/a_abc.gif",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := GenerateUserAvatarURL(tt.user); got != tt.want {
				t.Errorf("GenerateUserAvatarURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestUserCanManageGuild(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		guild *models.UserGuildResponse
		want  bool
	}{
		{name: "nil", guild: nil, want: false},
		{name: "owner", guild: &models.UserGuildResponse{Owner: true}, want: true},
		{name: "admin", guild: &models.UserGuildResponse{Permissions: models.DiscordPermAdministrator}, want: true},
		{name: "manage-server", guild: &models.UserGuildResponse{Permissions: models.DiscordPermManageServer}, want: true},
		{name: "member", guild: &models.UserGuildResponse{Permissions: models.DiscordPermSendMessages}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := UserCanManageGuild(tt.guild); got != tt.want {
				t.Errorf("UserCanManageGuild() = %v, want %v", got, tt.want)
			}
		})
	}
}
