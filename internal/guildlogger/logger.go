// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package guildlogger

import (
	"fmt"

	"github.com/andersfylling/disgord"
	"github.com/apex/log"
)

var _ disgord.Logger = (*discordLogger)(nil)

// discordLogger is a wrapper for apex/log, that can be used with the disgord.Logger
// interface.
type discordLogger struct {
	logger log.Interface
}

func New(l log.Interface) *discordLogger {
	return &discordLogger{logger: l.WithField("src", "disgord")}
}

func (l *discordLogger) Debug(v ...any) {
	l.logger.Debug(fmt.Sprint(v...))
}

func (l *discordLogger) Info(v ...any) {
	l.logger.Info(fmt.Sprint(v...))
}

func (l *discordLogger) Error(v ...any) {
	l.logger.Error(fmt.Sprint(v...))
}
