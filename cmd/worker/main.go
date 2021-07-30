// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/apex/log"
	"github.com/jessevdk/go-flags"
	_ "github.com/joho/godotenv/autoload"
	"github.com/lrstanley/spectrograph/internal/database"
	"github.com/lrstanley/spectrograph/internal/models"
	"github.com/lrstanley/spectrograph/internal/reactive"
)

// Should be auto-injected by build tooling.
const (
	version = "master"
	commit  = "latest"
	date    = "-"
)

var (
	cli models.FlagsWorkerServer

	logger log.Interface
)

func main() {
	_, err := flags.Parse(&cli)
	if err != nil {
		if FlagErr, ok := err.(*flags.Error); ok && FlagErr.Type == flags.ErrHelp {
			os.Exit(0)
		}
		os.Exit(1)
	}

	logger = cli.Logger.Parse(cli.Debug).WithFields(log.Fields{
		"build_version": fmt.Sprintf("%s/%s (%s)", version, commit, date),
	})

	// Validate misc. flags.
	if cli.Discord.ShardID+1 > cli.Discord.NumShards {
		logger.Fatal("provided shard id is greater than the total amount of shards")
	}

	if cli.Discord.NumShards < 1 {
		cli.Discord.NumShards = 1
	}

	if cli.Discord.ShardID < 0 {
		cli.Discord.ShardID = 0
	}

	ctx, closer := context.WithCancel(context.Background())
	defer closer()

	errorChan := make(chan error)
	wg := &sync.WaitGroup{}

	campaign := reactive.NewElection(logger, cli.Etcd, false, version)
	defer campaign.Close()
	campaign.Run(ctx, wg)

	// Wait for leader ourselves, or a leader that is compatible with us.
	// This is most critical during initial startup to ensure we're doing
	// db migrations only on the leader.
	_ = campaign.WaitForLeader(ctx, 30*time.Second, 3*time.Second)

	// Start listening for signals here.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)

	// TODO:
	//   - integrate with k8s to auto-scale based off recommendations from gateway?
	//     - https://github.com/hashicorp/go-discover
	//     - I'd assume this would happen on the api server primary?
	//   - if shard id > max concurrent, sleep for time * <id over limit>?

	// Initialize storer/database.
	var store models.Store
	logger.WithFields(log.Fields{
		"dbname": cli.Mongo.DBName,
		"uri":    cli.Mongo.URI,
	}).Info("database params")
	store = database.New(logger)

	logger.Info("initializing connections to database")
	if err = store.Setup(&cli.Mongo); err != nil {
		logger.WithError(err).Fatal("error initializing database")
	}
	defer store.Close()

	discordSetup(ctx, wg, errorChan)

	logger.Info("listening for signals")

	go func() {
		for {
			select {
			case leader := <-campaign.LeaderUpdates:
				if !leader.Leader() && !leader.IsCompatibleVersion() {
					logger.WithError(reactive.ErrIncompatibleVersion).Error("exiting to prevent data loss")
				} else {
					continue
				}
			case <-signals:
				logger.Error("signal received, closing connections")
			case err := <-errorChan:
				if err != nil {
					continue
				}
				logger.WithError(err).Error("error received, closing connections")
			}

			// Signal to exit.
			closer()
		}
	}()

	// Wait for the context to close, and wait for all goroutines/processes to exit.
	<-ctx.Done()
	wg.Wait()

	logger.Info("shutdown complete")
}
