// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package models

import (
	"fmt"

	"github.com/apex/log"
)

// Store provides a generic interface to multiple potential databases/services.
type Store interface {
	// Setup initializes the given Store (e.g. database connections).
	Setup(flags *FlagsHTTPServer) (err error)

	// Close garbage collects the database connections.
	Close() (err error)

	// Migrate initializes go-migrate on the database (if available).
	Migrate(flags *FlagsHTTPServer) (err error)

	NewUserService() UserService
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
