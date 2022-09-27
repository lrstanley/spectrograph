// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package permissions

import (
	"github.com/andersfylling/disgord"
	"github.com/lrstanley/spectrograph/internal/models"
)

type BotPermissions struct {
	member      *disgord.Member
	permissions models.DiscordPermissions
}
