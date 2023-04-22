// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"context"

	"github.com/apex/log"
	"github.com/lrstanley/chix"
	"github.com/lrstanley/clix"
	"github.com/lrstanley/spectrograph/internal/database"
	"github.com/lrstanley/spectrograph/internal/database/ent"
	_ "github.com/lrstanley/spectrograph/internal/database/ent/runtime"
	"github.com/lrstanley/spectrograph/internal/models"
	"github.com/lrstanley/spectrograph/internal/worker"
)

var (
	db     *ent.Client
	logger log.Interface

	cli = &clix.CLI[models.Flags]{
		Links: clix.GithubLinks("github.com/lrstanley/spectrograph", "master", "https://liam.sh"),
	}
)

func main() {
	cli.Parse()
	logger = cli.Logger

	db = database.Open(context.Background(), logger, cli.Flags.Database)
	defer db.Close()

	ctx := ent.NewContext(log.NewContext(context.Background(), logger), db)

	database.Migrate(ctx, logger)

	w, err := worker.New(ctx, cli.Flags.DefaultWorker, cli.Flags.Discord)
	if err != nil {
		logger.WithError(err).Fatal("failed to create worker")
	}

	if err := chix.RunCtx(
		ctx,
		httpServer(ctx),
		w.Run,
		database.CronGuildEvents,
	); err != nil {
		logger.WithError(err).Fatal("shutting down")
	}
}
