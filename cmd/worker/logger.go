// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"fmt"

	"github.com/andersfylling/disgord"
	"github.com/apex/log"
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
