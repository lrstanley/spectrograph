// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package logging

import (
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/discard"
	"github.com/apex/log/handlers/json"
	"github.com/apex/log/handlers/logfmt"
	"github.com/apex/log/handlers/text"
	"github.com/lrstanley/spectrograph/pkg/models"
)

// ParseConfig parses
func ParseConfig(flags models.LoggerConfig, debug bool) *log.Logger {
	logger := &log.Logger{}

	if debug {
		logger.Level = log.DebugLevel
	} else {
		logger.Level = log.MustParseLevel(flags.Level)
	}

	if flags.Quiet {
		logger.Handler = discard.New()
	} else if flags.JSON {
		logger.Handler = json.New(os.Stdout)
	} else if flags.Pretty {
		logger.Handler = text.New(os.Stdout)
	} else {
		logger.Handler = logfmt.New(os.Stdout)
	}

	// Set global options as well, just in case.
	log.SetLevel(logger.Level)
	log.SetHandler(logger.Handler)

	return logger
}
