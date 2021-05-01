// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package models

import (
	"fmt"
	"net/url"
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/discard"
	"github.com/apex/log/handlers/json"
	"github.com/apex/log/handlers/logfmt"
	"github.com/apex/log/handlers/text"
	"github.com/jessevdk/go-flags"
)

// TODO: eventually we should probably JUST support env vars, if this will
// always be a containerized deployment.

// LoggerConfig are the flags that define how log entries are processed/returned.
type LoggerConfig struct {
	Quiet  bool   `env:"QUIET"  long:"quiet"  description:"disable logging to stdout (also: see levels)"`
	Level  string `env:"LEVEL"  long:"level"  default:"info" choice:"debug" choice:"info" choice:"warn" choice:"error" choice:"fatal"  description:"logging level"`
	JSON   bool   `env:"JSON"   long:"json"   description:"output logs in JSON format"`
	Pretty bool   `env:"PRETTY" long:"pretty" description:"output logs in a pretty colored format (cannot be easily parsed)"`
}

// ParseLoggerConfig parses LoggerConfig and adjusts the structured logger as
// necessary.
func (c LoggerConfig) Parse(debug bool) log.Interface {
	logger := &log.Logger{}

	if debug {
		logger.Level = log.DebugLevel
	} else {
		logger.Level = log.MustParseLevel(c.Level)
	}

	if c.Quiet {
		logger.Handler = discard.New()
	} else if c.JSON {
		logger.Handler = json.New(os.Stdout)
	} else if c.Pretty {
		logger.Handler = text.New(os.Stdout)
	} else {
		logger.Handler = logfmt.New(os.Stdout)
	}

	// Set global options as well, just in case.
	log.SetLevel(logger.Level)
	log.SetHandler(logger.Handler)

	return logger
}

// FlagsHTTPServer are flags specifically utilized by the HTTP service.
type FlagsHTTPServer struct {
	Debug bool `env:"DEBUG" long:"debug" description:"enable debugging (pprof endpoints), CSRF protection, as well as disable caching of templates"`

	// Logging.
	Logger LoggerConfig `group:"Logging Options" namespace:"log" env-namespace:"LOG"`

	// HTTP.
	// TODO: if more than one service makes an http service, hoist this out.
	HTTP struct {
		BindAddr   string `env:"BIND_ADDR"    long:"bind-addr"    default:":8080" required:"true" description:"ip:port pair to bind to"`
		Proxy      bool   `env:"BEHIND_PROXY" long:"behind-proxy" description:"if X-Forwarded-For headers should be trusted"`
		RawBaseURL string `env:"BASE_URL"     long:"base-url"     required:"true" description:"public url (including scheme and port as necessary)"`
		BaseURL    *url.URL

		// Example way of doing this in most linux environments:
		//   tr -dc A-Za-z0-9 </dev/urandom | head -c 32;echo
		//   tr -dc 'A-Za-z0-9!#$%&*+\-./:;=<>?@[]^_~' </dev/urandom | head -c 32;echo
		SessionKeys []string `env:"SESSION_KEYS" long:"session-keys" required:"true" env-delim:"," description:"key pairs used to encrypt cookie session data (can specify multiple: hash,block,hash,block) -- hash should be at least 32 bytes long, block keys should be 16 bytes (AES-128) or 32 bytes (AES-256) long"`

		TLS struct {
			Enabled bool   `env:"ENABLED" long:"enabled" description:"run tls server rather than standard http"`
			Cert    string `env:"CERT"    long:"cert"    description:"path to ssl cert file"`
			Key     string `env:"KEY"     long:"key"     description:"path to ssl key file"`
		} `group:"TLS Options" namespace:"tls" env-namespace:"TLS"`
	} `group:"HTTP Server options" namespace:"http" env-namespace:"HTTP"`

	// Authentication.
	Auth struct {
		Discord struct {
			ID     string   `env:"ID"     long:"ID"     required:"true" description:"Discord oauth2 ID"`
			Secret string   `env:"SECRET" long:"secret" required:"true" description:"Discord oauth2 secret"`
			Admins []string `env:"ADMINS" long:"admins" required:"true" env-delim:"," description:"user id's of the users you want to be admins"`
		} `group:"Discord Options" namespace:"discord" env-namespace:"DISCORD"`
	} `group:"Authentication Options" namespace:"auth" env-namespace:"AUTH"`

	// Databases.
	Migration struct {
		Disabled bool `env:"DISABLED" long:"disabled" description:"disable database migrations"`
		Purge    bool `env:"PURGE"    long:"purge"    hidden:"true" description:"PURGES ALL DATA ON STARTUP, BE WARNED"`
		Force    bool `env:"FORCE"    long:"force"    description:"force update to version in database (must also specify version)"`
		Version  uint `env:"VERSION"  long:"version"  description:"optional version to migrate the database to"`
	} `group:"Database Migration Options (CAUTION!)" namespace:"migration" env-namespace:"MIGRATION"`

	Mongo struct {
		DBName string `env:"DB_NAME" long:"db-name" default:"spectrograph" description:"database name to use"`
		URI    string `env:"URI"     long:"uri"     default:"mongodb://localhost:27017/?maxPoolSize=64" description:"mongodb connection string (see: https://docs.mongodb.com/manual/reference/connection-string/)"`
	} `group:"Database (MongoDB) Options" namespace:"mongo" env-namespace:"MONGO"`
}

// FlagsWorkerServer are flags used by the worker service.
type FlagsWorkerServer struct {
	Debug bool `env:"DEBUG" long:"debug" description:"enable debugging"`

	// Authentication.
	Auth struct {
		Discord struct {
			ID     string `env:"ID"     long:"ID"     required:"true" description:"Discord oauth2 ID"`
			Secret string `env:"SECRET" long:"secret" required:"true" description:"Discord oauth2 secret"`
		} `group:"Discord Options" namespace:"discord" env-namespace:"DISCORD"`
	} `group:"Authentication Options" namespace:"auth" env-namespace:"AUTH"`

	Discord struct {
		// https://discord.com/developers/docs/topics/gateway#sharding
		// Note: Shard ID 0 will be the only one to receive DMs.
		ShardID   int `env:"SHARD_ID" description:"shard ID of this specific worker (from 0)"`
		NumShards int `env:"NUM_SHARDS" description:"number of total shards"`
	} `group:"Discord Options" namespace:"discord" env-namespace:"DISCORD"`

	// Logging.
	Logger LoggerConfig `group:"Logging Options" namespace:"log" env-namespace:"LOG"`
}

func FlagParse(data interface{}) (args []string) {
	parser := flags.NewParser(data, flags.HelpFlag|flags.PrintErrors|flags.PassDoubleDash)
	parser.NamespaceDelimiter = "."
	parser.EnvNamespaceDelimiter = "_"

	var err error

	args, err = parser.Parse()
	if err != nil {
		if FlagErr, ok := err.(*flags.Error); ok && FlagErr.Type == flags.ErrHelp {
			os.Exit(0)
		}

		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return args
}
