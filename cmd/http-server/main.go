// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/apex/log"
	"github.com/golang-migrate/migrate"
	_ "github.com/joho/godotenv/autoload"
	"github.com/lrstanley/spectrograph/internal/database"
	"github.com/lrstanley/spectrograph/internal/models"
	"golang.org/x/oauth2"
)

// Should be auto-injected by build tooling.
const (
	version = "master"
	commit  = "latest"
	date    = "-"
)

var (
	cli    models.FlagsHTTPServer
	logger log.Interface

	svcUsers    models.UserService
	svcSessions scs.Store
	oauthConfig *oauth2.Config
)

func main() {
	_ = models.FlagParse(&cli)
	logger = cli.Logger.Parse(cli.Debug)

	var err error

	if cli.HTTP.BaseURL, err = url.Parse(cli.HTTP.RawBaseURL); err != nil {
		logger.WithError(err).Fatalf("invalid base url provided: %v", cli.HTTP.RawBaseURL)
	}
	if !strings.HasPrefix(cli.HTTP.BaseURL.Scheme, "http") {
		logger.WithError(err).Fatalf("invalid base url provided: %v", cli.HTTP.RawBaseURL)
	}
	cli.HTTP.BaseURL.Path = strings.TrimRight(cli.HTTP.BaseURL.Path, "/")

	oauthConfig = &oauth2.Config{
		ClientID:     cli.Auth.Discord.ID,
		ClientSecret: cli.Auth.Discord.Secret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://discord.com/api/oauth2/authorize",
			TokenURL: "https://discord.com/api/oauth2/token",
		},
		RedirectURL: cli.HTTP.BaseURL.String() + "/auth/callback",
		Scopes: []string{
			"identify", "email", "guilds",
		},
	}

	// Initialize storer/database.
	var store models.Store
	logger.WithFields(log.Fields{
		"dbname": cli.Mongo.DBName,
		"uri":    cli.Mongo.URI,
	}).Info("database params")
	store = database.New(logger)

	logger.Info("initializing connections to database")
	if err = store.Setup(&cli); err != nil {
		logger.WithError(err).Fatal("error initializing database")
	}
	defer store.Close()

	// Start listening for signals here, to prevent corruption during potential database
	// migrations. I.e. even if the OS tries to send a signal, we won't act on it until
	// we have initialized things.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)

	// Setup methods to allow signaling to all children methods that we're stopping.
	ctx, closer := context.WithCancel(context.Background())
	errorChan := make(chan error)
	wg := &sync.WaitGroup{}

	// Initialize migrations.
	if !cli.Migration.Disabled {
		logger.Info("running database migrations")
		if err = store.Migrate(&cli); err != nil {
			if errors.As(err, &migrate.ErrNoChange) {
				logger.Info("database migration: no changes found")
			} else if errors.As(err, &migrate.ErrNilVersion) {
				logger.Info("database migration: no version information in the database")
			} else {
				logger.WithError(err).Fatal("error during migration")
			}
		}
	}

	// Initialize services.
	svcUsers = store.NewUserService()
	svcSessions = store.NewSessionService(ctx, 5*time.Minute)

	// Initialize the http/https server.
	httpServer(ctx, wg, errorChan)

	logger.Info("listening for signal. CTRL+C to quit.")

	go func() {
		for {
			select {
			case <-signals:
				fmt.Println("\nsignal received, shutting down")
			case <-errorChan:
				logger.WithError(err).Error("error received")
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
