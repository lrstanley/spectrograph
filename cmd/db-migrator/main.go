// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/apex/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-migrate/migrate"
	"github.com/jessevdk/go-flags"
	_ "github.com/joho/godotenv/autoload"
	"github.com/lrstanley/pt"
	"github.com/lrstanley/spectrograph/internal/database"
	"github.com/lrstanley/spectrograph/internal/httpware"
	"github.com/lrstanley/spectrograph/internal/models"
)

// Should be auto-injected by build tooling.
const (
	version = "master"
	commit  = "latest"
	date    = "-"
)

var (
	cli    models.FlagsMigratorServer
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

	// Start listening for signals here, to prevent corruption during database
	// migrations. I.e. even if the OS tries to send a signal, we won't act
	// on it until we have initialized things.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)

	// Setup methods to allow signaling to all children methods that we're stopping.
	ctx, closer := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		// Initialize migrations.
		logger.Info("running database migrations")
		if err = store.Migrate(ctx, &cli.Mongo, &cli.Migration); err != nil {
			if errors.As(err, &migrate.ErrNoChange) {
				logger.Info("database migration: no changes found")
			} else if errors.As(err, &migrate.ErrNilVersion) {
				logger.Info("database migration: no version information in the database")
			} else {
				logger.WithError(err).Fatal("error during migration")
			}
		}

		logger.Info("migrations (if any) are complete")

		if !cli.StayAlive {
			logger.Info("stay-alive not set, closing")
			closer()
		}
	}()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(httpware.StructuredLogger(logger, nil, !cli.Debug))

	if cli.HTTP.Proxy {
		r.Use(middleware.RealIP)
	}

	r.Use(middleware.Timeout(5 * time.Second))
	r.Use(middleware.GetHead)

	r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		pt.JSON(w, r, pt.M{
			"healthy":          true,
			"migration_config": cli.Migration,
		})
	})

	// Setup our http server that's used for healthchecks.
	srv := &http.Server{
		Addr:    cli.HTTP.BindAddr,
		Handler: r,

		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if cli.StayAlive {
		go func() {
			logger.WithField("bind", cli.HTTP.BindAddr).Info("initializing http server")

			if err := srv.ListenAndServe(); err != nil {
				logger.WithError(err).Error("server closed")
			}

			closer()
		}()
	}

	// Wait until we receive any signals, then close.
	select {
	case <-signals:
		logger.Info("signal received, exiting")
	case <-ctx.Done():
		logger.Info("context closed, exiting")
	}

	closer()
	wg.Wait()
	logger.Info("shutdown complete")
}
