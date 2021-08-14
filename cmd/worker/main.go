// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/apex/log"
	"github.com/golang/protobuf/ptypes/empty"
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
	cli       models.FlagsWorkerServer
	rpcWorker rpc.Worker
	logger    log.Interface
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

	rpcWorker = rpc.NewWorkerClient(cli.API.URI, rpc.NewHTTPClient(30*time.Second, map[string]string{
		"X-Api-Version": version,
		"X-Api-Key":     cli.API.Key,
		"X-Shard-Id":    strconv.Itoa(cli.Discord.ShardID),
	}))

	if !checkAPIHealth(ctx) {
		closer()
		wg.Wait()
		os.Exit(1)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(150 * time.Second):
				if !checkAPIHealth(ctx) {
					closer()
					wg.Wait()
					os.Exit(1)
				}
			}
		}
	}()

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

func checkAPIHealth(ctx context.Context) (healthy bool) {
	resp, err := rpcWorker.Health(ctx, &empty.Empty{})
	if err != nil {
		logger.WithError(err).Error("healthcheck: failed while waiting for api server to respond (multiple attempts)")
		return false
	}
	if !resp.Ready {
		logger.Error("healthcheck: api server is not ready yet")
		return false
	}

	logger.Info("healthcheck: api server is responsive")
	return true
}
