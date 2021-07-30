// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package models

import (
	"context"
	"fmt"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/apex/log"
)

// Store provides a generic interface to multiple potential databases/services.
type Store interface {
	// Setup initializes the given Store (e.g. database connections).
	Setup(flags *MongoConfig) (err error)

	// Close garbage collects the database connections.
	Close() (err error)

	// Migrate initializes go-migrate on the database (if available).
	Migrate(ctx context.Context, mongoFlags *MongoConfig, migrateFlags *MigrateConfig) (err error)

	NewUserService() UserService
	NewSessionService(ctx context.Context, cleanup time.Duration) scs.Store
}

// MigrateLogger implements the migrate.Logger interface.
type MigrateLogger struct {
	Logger *log.Entry
}

// Printf wraps log.Logger to add a prefix.
func (m *MigrateLogger) Printf(format string, v ...interface{}) {
	m.Logger.WithField("source", "migrator").Info(fmt.Sprintf(format, v...))
}

// Verbose should always be enabled.
func (MigrateLogger) Verbose() bool { return true }
