// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/lrstanley/chix/v2"
	"github.com/lrstanley/clix/v2"
	"github.com/lrstanley/spectrograph/internal/database"
	"github.com/lrstanley/spectrograph/internal/models"
	"github.com/lrstanley/spectrograph/internal/worker"
	"github.com/lrstanley/x/sync/scheduler"
)

func main() {
	clix.NewWithDefaults(
		clix.WithAppInfo[models.Flags](clix.AppInfo{
			Links: clix.GithubLinks("github.com/lrstanley/spectrograph", "master", "https://liam.sh"),
		}),
	).Execute(run)
}

func run(ctx context.Context, cli *clix.CLI[models.Flags], logger *slog.Logger) error {
	db, err := database.Open(ctx, logger, cli.Flags.Database.URL)
	if err != nil {
		return fmt.Errorf("opening database: %w", err)
	}
	defer db.Close() //nolint:errcheck

	if merr := database.Migrate(ctx, logger, db); merr != nil {
		return merr
	}

	// Wakes guild event streams as soon as the (in-process) worker logs events.
	guildEvents := &database.GuildEventNotifier{}
	db.GuildEvent.Use(guildEvents.Hook())

	srv, err := httpServer(ctx, logger, db, guildEvents, cli)
	if err != nil {
		return err
	}

	jobs := []scheduler.Job{
		scheduler.NewCron("prune-guild-events", scheduler.JobLoggerFunc(database.PruneGuildEvents(db))).
			WithInterval(10 * time.Minute).
			WithLogger(logger),
	}

	if cli.Flags.Worker.Enabled {
		w, werr := worker.New(ctx, logger, db, cli.Flags.Worker, worker.WithDashboardURL(cli.Flags.HTTP.BaseURL))
		if werr != nil {
			return fmt.Errorf("creating worker: %w", werr)
		}
		jobs = append(jobs, scheduler.JobFunc(w.Run))
	} else {
		logger.WarnContext(ctx, "discord worker is disabled, only running the http server")
	}

	return chix.Run(ctx, logger, srv, jobs...)
}
