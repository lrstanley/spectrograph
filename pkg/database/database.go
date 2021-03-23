// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package database

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/apex/log"
	"github.com/golang-migrate/migrate/v4"
	mongomigrate "github.com/golang-migrate/migrate/v4/database/mongodb"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
	"github.com/lrstanley/spectrograph/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// Commands used during migration:
//   https://docs.mongodb.com/manual/reference/command/create/

const (
	defaultMaxPoolSize = 32
)

// New returns a mongo Store implementation.
func New(logger *log.Logger) *mongoStore {
	return &mongoStore{
		log: logger.WithFields(log.Fields{
			"source":   "database",
			"database": "mongo",
		}),
	}
}

type mongoStore struct {
	db     *mongo.Database
	client *mongo.Client
	log    *log.Entry

	// Pre-initialized collection configs.
	project     *mongo.Collection
	scanRequest *mongo.Collection
	provider    *mongo.Collection
	user        *mongo.Collection
	session     *mongo.Collection
}

func (s *mongoStore) Setup(flags *models.FlagsHTTPServer) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	opts := &options.ClientOptions{}

	opts.SetAppName("spectrograph")
	// TODO: this makes things very slow. do we really need it?
	// opts.SetCompressors([]string{"zstd", "zlib", "snappy"})
	opts.SetConnectTimeout(15 * time.Second)
	opts.SetDirect(true)
	opts.SetMaxConnIdleTime(120 * time.Second)
	opts.SetMinPoolSize(5)
	opts.SetRetryReads(true)
	opts.SetRetryWrites(true)
	opts.ApplyURI(flags.Mongo.URI)

	if *opts.MaxPoolSize < *opts.MinPoolSize {
		opts.SetMaxPoolSize(defaultMaxPoolSize)
	}

	s.client, err = mongo.Connect(ctx, opts)
	if err != nil {
		return fmt.Errorf("unable to connect to mongodb: %v", err)
	}

	// Ping the primary to make sure it's online.
	if err = s.client.Ping(ctx, readpref.Primary()); err != nil {
		return fmt.Errorf("unable to ping mongodb primary: %v", err)
	}

	// Anything that has a preference on write majority...
	wcMajority := writeconcern.New(writeconcern.WMajority(), writeconcern.WTimeout(10*time.Second))
	wcMajorityCollectionOpts := options.Collection().SetWriteConcern(wcMajority)

	s.db = s.client.Database(flags.Mongo.DBName)
	s.project = s.db.Collection("projects")
	s.scanRequest = s.db.Collection("scan_requests", wcMajorityCollectionOpts)
	s.provider = s.db.Collection("providers")
	s.user = s.db.Collection("users")
	s.session = s.db.Collection("sessions")

	return nil
}

func (s *mongoStore) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.client.Disconnect(ctx); err != nil {
		err = fmt.Errorf("unable to close mongodb connection: %v", err)
		s.log.Error(err.Error())
		return err
	}
	return nil
}

func (s *mongoStore) Migrate(flags *models.FlagsHTTPServer) error {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	// Access ricebox bundled migrations.
	migrationAssets, err := rice.FindBox("migrations/")
	if err != nil {
		return err
	}

	// Walk through migrations and build a file list to pass into bindata.Resource.
	migrations := []string{}
	err = migrationAssets.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".json") {
			return err
		}

		migrations = append(migrations, path)
		return nil
	})
	if err != nil {
		return err
	}

	source, err := bindata.WithInstance(bindata.Resource(migrations, migrationAssets.Bytes))
	if err != nil {
		return err
	}

	// Check if one of the nodes we connected to is in a replicaset.
	replicaSet := true
	cursor := s.client.Database("admin").RunCommand(ctx, bson.D{{"replSetGetStatus", 1}})
	var result bson.M
	if err = cursor.Decode(&result); err != nil {
		if cerr, ok := err.(mongo.CommandError); ok {
			if cerr.Name != "NoReplicationEnabled" {
				return err
			}
			replicaSet = false
		} else {
			return err
		}
	}

	destination, err := mongomigrate.WithInstance(s.client, &mongomigrate.Config{
		DatabaseName:    flags.Mongo.DBName,
		TransactionMode: replicaSet,
	})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("bindata", source, flags.Mongo.DBName, destination)
	if err == nil {
		// TODO: this logic can be moved to the main package because everything after
		// this point has potentially duplicated logic.
		m.Log = &models.MigrateLogger{Logger: s.log}

		// Do they want to purge data during startup?
		if flags.Migration.Purge {
			s.log.Info("migration: purge requested, dropping")
			err = m.Drop()
			if err != nil {
				// Don't exit if there is a purge error (could just be no changes found).
				s.log.WithError(err).Error("migration: purge errored")
			} else {
				s.log.Info("migration: purge complete")
			}
		}

		// Obey user-provided migration. Since this is likely related to a corruption or
		// downgrade, don't allow the app to continue the startup, and make sure to quit.
		if flags.Migration.Version != 0 {
			s.log.WithField("version", flags.Migration.Version).Info("database migration to specific version requested")
			if flags.Migration.Force {
				err = m.Force(int(flags.Migration.Version))
			} else {
				err = m.Migrate(flags.Migration.Version)
			}

			if err == nil {
				s.log.Info("successfully migrated to requested version. since this is a manual request, exiting gracefully to prevent corruption.")
				os.Exit(0)
			}
		} else {
			err = m.Up()
		}
	}

	return err
}
