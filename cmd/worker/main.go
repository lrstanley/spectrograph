// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/apex/log"
	"github.com/jessevdk/go-flags"
	_ "github.com/joho/godotenv/autoload"
	"github.com/lrstanley/spectrograph/internal/models"
)

var (
	// For use with goreleaser, it will auto-inject version/commit/date/etc.
	version = "master"
	commit  = "latest"
	date    = "-"

	cli models.FlagsWorkerServer

	logger *log.Logger
)

func main() {
	_, err := flags.Parse(&cli)
	if err != nil {
		if FlagErr, ok := err.(*flags.Error); ok && FlagErr.Type == flags.ErrHelp {
			os.Exit(0)
		}
		os.Exit(1)
	}

	logger = cli.Logger.Parse(cli.Debug)

	// Start listening for signals here, to prevent corruption during potential database
	// migrations.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)

	for range signals {
		logger.Info("received signal, closing connections")

		logger.Info("done cleaning up; exiting")
		os.Exit(1)
	}
}
