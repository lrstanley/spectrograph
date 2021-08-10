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
	"github.com/lrstanley/spectrograph/internal/models"
	"github.com/lrstanley/spectrograph/internal/rpc"
)

// Should be auto-injected by build tooling.
const (
	version = "master"
	commit  = "latest"
	date    = "-"
)

var (
	cli    models.FlagsWorkerServer
	api    rpc.Worker
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
		"shard_id":      cli.Discord.ShardID,
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

	// Start listening for signals here.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)

	// TODO:
	//   - integrate with k8s to auto-scale based off recommendations from gateway?
	//     - https://github.com/hashicorp/go-discover
	//     - I'd assume this would happen on the api server primary?
	//   - if shard id > max concurrent, sleep for time * <id over limit>?

	ctx, closer := context.WithCancel(context.Background())
	defer closer()

	errorChan := make(chan error)
	wg := &sync.WaitGroup{}

	// Wait for initial health check from the api server before we continue.
	api = rpc.NewWorkerClient(cli.RPC.URI, cli.RPC.SecretKey, version, cli.Discord.ShardID, 10*time.Second, 5)

	if resp, err := api.Health(ctx, &rpc.NoArgs{}); err != nil {
		logger.WithError(err).Fatal("failed while waiting for rpc server to respond")
	} else {
		logger.WithField("health", resp.Ready).Info("rpc server is responsive")
	}

	discordSetup(ctx, wg, errorChan)

	logger.Info("listening for signals")

	go func() {
		for {
			select {
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
