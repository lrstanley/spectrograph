// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package models

// TODO: eventually we should probably JUST support env vars, if this will
// always be a containerized deployment.
// TODO: cut out TLS functionality. assume frontend by reverse proxy? Simplifies
// http servers and similar, by quite a bit.

// LoggerConfig are the flags that define how log entries are processed/returned.
type LoggerConfig struct {
	Quiet  bool   `env:"LOGGING_QUIET"  long:"quiet"  description:"disable logging to stdout (also: see levels)"`
	Level  string `env:"LOGGING_LEVEL"  long:"level"  default:"info" choice:"debug" choice:"info" choice:"warn" choice:"error" choice:"fatal"  description:"logging level"`
	JSON   bool   `env:"LOGGING_JSON"   long:"json"   description:"output logs in JSON format"`
	Pretty bool   `env:"LOGGING_PRETTY" long:"pretty" description:"output logs in a pretty colored format (cannot be easily parsed)"`
}

// FlagsHTTPServer are flags specifically utilized by the HTTP service.
type FlagsHTTPServer struct {
	Debug bool `env:"DEBUG" short:"d" long:"debug" description:"enable debugging (pprof endpoints), CSRF protection, as well as disable caching of templates"`

	// Logging.
	Logger LoggerConfig `group:"Logging Options" namespace:"log"`

	// HTTP.
	HTTP  string `env:"HTTP"         short:"b" long:"http" default:":8080" description:"ip:port pair to bind to" required:"true"`
	Proxy bool   `env:"BEHIND_PROXY" short:"p" long:"behind-proxy" description:"if X-Forwarded-For headers should be trusted"`
	TLS   struct {
		Enabled bool   `env:"TLS_ENABLED" long:"enabled" description:"run tls server rather than standard http"`
		Cert    string `env:"TLS_CERT"    long:"cert"    description:"path to ssl cert file"`
		Key     string `env:"TLS_KEY"     long:"key"     description:"path to ssl key file"`
	} `group:"TLS Options" namespace:"tls"`

	// Authentication.
	Auth struct {
		Github struct {
			ClientID     string   `env:"AUTH_GITHUB_CLIENT_ID"     long:"client-id"     description:"GitHub OAuth Client ID"`
			ClientSecret string   `env:"AUTH_GITHUB_CLIENT_SECRET" long:"client-secret" description:"GitHub OAuth Client Secret"`
			Admins       []string `env:"AUTH_GITHUB_ADMINS"        long:"admins"        description:"user id's of the users you want to be admins"`
		} `group:"GitHub Options" namespace:"github"`
	} `group:"Authentication Options" namespace:"auth"`

	// Databases.
	Migration struct {
		Disabled bool `env:"MIGRATION_DISABLED" long:"disabled" description:"disable database migrations"`
		Purge    bool `env:"MIGRATION_PURGE"    long:"purge" hidden:"true" description:"PURGES ALL DATA ON STARTUP, BE WARNED"`
		Force    bool `env:"MIGRATION_FORCE"    long:"force" description:"force update to version in database (must also specify version)"`
		Version  uint `env:"MIGRATION_VERSION"  long:"version" description:"optional version to migrate the database to"`
	} `group:"Database Migration Options (CAUTION!)" namespace:"migration"`
	Mongo struct {
		DBName string `env:"MONGO_DB_NAME" long:"db-name" default:"spectrograph" description:"database name to use"`
		URI    string `env:"MONGO_URI"     long:"uri"     default:"mongodb://localhost:27017/?maxPoolSize=64" description:"mongodb connection string (see: https://docs.mongodb.com/manual/reference/connection-string/)" env:"DB_URI"`
	} `group:"Database (MongoDB) Options" namespace:"mongo"`
}

// FlagsWorkerServer are flags used by the worker service.
type FlagsWorkerServer struct {
	Debug bool `env:"DEBUG" short:"d" long:"debug" description:"enable debugging"`

	// Logging.
	Logger LoggerConfig `group:"Logging Options" namespace:"log"`
}
