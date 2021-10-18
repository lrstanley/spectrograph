// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"fmt"

	"github.com/andersfylling/disgord"
	"github.com/apex/log"
	"github.com/lrstanley/spectrograph/internal/models"
)

// LoggerApex is a wrapper for apex/log, that can be used with the disgord.Logger
// interface.
type LoggerApex struct {
	logger log.Interface
}

var _ disgord.Logger = (*LoggerApex)(nil)

func (l *LoggerApex) Debug(v ...interface{}) {
	l.logger.WithField("src", "disgord").Debug(fmt.Sprint(v...))
}

func (l *LoggerApex) Info(v ...interface{}) {
	l.logger.WithField("src", "disgord").Info(fmt.Sprint(v...))
}

func (l *LoggerApex) Error(v ...interface{}) {
	l.logger.WithField("src", "disgord").Error(fmt.Sprint(v...))
}

// Additional logging helpers below.

func logGuild(l log.Interface, v interface{}) log.Interface {
	switch guild := v.(type) {
	case *disgord.Guild:
		return l.WithFields(log.Fields{
			"guild_id":     guild.ID.String(),
			"guild_name":   guild.Name,
			"guild_owner":  guild.OwnerID.String(),
			"guild_region": guild.Region,
		})
	case *models.Server:
		return l.WithFields(log.Fields{
			"guild_id":     guild.Discord.ID,
			"guild_name":   guild.Discord.Name,
			"guild_owner":  guild.Discord.OwnerID,
			"guild_region": guild.Discord.Region,
		})
	case *models.ServerDiscordData:
		return l.WithFields(log.Fields{
			"guild_id":     guild.ID,
			"guild_name":   guild.Name,
			"guild_owner":  guild.OwnerID,
			"guild_region": guild.Region,
		})
	case disgord.Snowflake:
		return l.WithField("guild_id", guild.String())
	default:
		return l.WithField("guild_id", v)
	}
}
