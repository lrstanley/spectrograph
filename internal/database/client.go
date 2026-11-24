// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

// Package database provides the Postgres-backed ent client, migrations, and
// database-related helpers (transactions, guild event streaming, pruning).
package database

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql/schema"
	"github.com/lrstanley/spectrograph/internal/database/ent"
	"github.com/lrstanley/spectrograph/internal/database/ent/migrate"
	_ "github.com/lrstanley/spectrograph/internal/database/ent/runtime" // Required by ent.
	"github.com/lrstanley/x/entx"
)

//go:generate go run -mod=readonly entc.go

const (
	connectAttempts = 3
	connectBackoff  = 5 * time.Second
)

// Open opens a Postgres connection pool (using pgx, see [entx.Driver] for the
// supported pool options), and returns an ent client. Connecting is retried a
// few times, to account for the database still starting up.
func Open(ctx context.Context, logger *slog.Logger, dsn *url.URL) (*ent.Client, error) {
	if err := entx.ConstrainDialect(dsn, dialect.Postgres); err != nil {
		return nil, fmt.Errorf("database url: %w", err)
	}

	driver, err := entx.Driver(ctx, logger, dsn)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	for attempt := 1; ; attempt++ {
		pctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		err = driver.DB().PingContext(pctx)
		cancel()
		if err == nil {
			break
		}

		logger.ErrorContext(ctx, "failed to connect to database", "attempt", attempt, "error", err)
		if attempt >= connectAttempts {
			_ = driver.Close()
			return nil, fmt.Errorf("connecting to database after %d attempts: %w", attempt, err)
		}

		select {
		case <-ctx.Done():
			_ = driver.Close()
			return nil, ctx.Err()
		case <-time.After(connectBackoff):
		}
	}

	logger.InfoContext(ctx, "connected to database")
	return ent.NewClient(ent.Driver(driver)), nil
}

// Ping checks if the database is reachable.
func Ping(ctx context.Context, db *ent.Client) error {
	if db == nil {
		return errors.New("database client is nil")
	}
	_, err := db.ExecContext(ctx, "SELECT 1")
	return err
}

// Migrate runs the ent auto-migration against the database.
func Migrate(ctx context.Context, logger *slog.Logger, db *ent.Client) error {
	if db == nil {
		return errors.New("database client is nil")
	}

	logger.InfoContext(ctx, "initiating database schema migration")
	if err := db.Schema.Create(
		ctx,
		schema.WithDropColumn(true),
		schema.WithDropIndex(true),
		migrate.WithGlobalUniqueID(true),
	); err != nil {
		return fmt.Errorf("migrating database schema: %w", err)
	}
	logger.InfoContext(ctx, "database schema migration complete")
	return nil
}
