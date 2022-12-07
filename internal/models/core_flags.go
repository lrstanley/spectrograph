// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package models

import "time"

type Flags struct {
	HTTP          ConfigHTTP     `group:"HTTP Server options"    namespace:"http"     env-namespace:"HTTP"`
	Database      ConfigDatabase `group:"Database options"       namespace:"database" env-namespace:"DATABASE"`
	Discord       ConfigDiscord  `group:"Discord options"        namespace:"discord"  env-namespace:"DISCORD"`
	DefaultWorker ConfigWorker   `group:"Default worker options" namespace:"worker"   env-namespace:"WORKER"`
}

// ConfigDatabase holds the database configuration.
type ConfigDatabase struct {
	URL      string `env:"URL"      long:"url" required:"true" description:"database connection url"`
	Username string `env:"USERNAME" long:"username" description:"database username if not specified via the URL"`
	Password string `env:"PASSWORD" long:"password" description:"database password if not specified via the URL"`
	Database string `env:"DATABASE" long:"database" description:"database name if not specified via the URL"`
}

// ConfigHTTP are configurations specifically utilized by the HTTP service.
type ConfigHTTP struct {
	BaseURL        string   `env:"BASE_URL"        long:"base-url"        default:"http://localhost:8080" description:"base url for the HTTP server"`
	BindAddr       string   `env:"BIND_ADDR"       long:"bind-addr"       default:":8080" required:"true" description:"ip:port pair to bind to"`
	TrustedProxies []string `env:"TRUSTED_PROXIES" long:"trusted-proxies" env-delim:"," description:"CIDR ranges that we trust the X-Forwarded-For header from (addl opts: local, *, cloudflare, and/or custom header to use)"`
	ValidationKey  string   `env:"VALIDATION_KEY"  long:"validation-key"  required:"true" description:"key used to validate session cookies (32 or 64 bytes)"`
	EncryptionKey  string   `env:"ENCRYPTION_KEY"  long:"encryption-key"  required:"true" description:"key used to encrypt session cookies (32 bytes)"`
}

// ConfigDiscord are configurations specifically utilized for interacting with
// Discord.
type ConfigDiscord struct {
	ClientID     string   `env:"CLIENT_ID"     long:"client-id"     required:"true" description:"Discord client ID"`
	ClientSecret string   `env:"CLIENT_SECRET" long:"client-secret" required:"true" description:"Discord client secret"`
	Admins       []string `env:"ADMINS"        long:"admins"        required:"true" env-delim:"," description:"user id's of the users you want to be admins"`
}

type ConfigWorker struct {
	// https://discord.com/developers/docs/topics/gateway#sharding
	// Note: Shard ID 0 will be the only one to receive DMs.
	ShardIDs           []uint        `env:"SHARD_ID"             long:"shard-id"   default:"0" env-delim:"," description:"shard ID of this specific worker (from 0)"`
	NumShards          uint          `env:"NUM_SHARDS"           long:"num-shards" default:"1"            description:"number of total shards"`
	BotToken           string        `env:"BOT_TOKEN"            long:"bot-token"  required:"true"        description:"Discord bot token"`
	EventFlushInterval time.Duration `env:"EVENT_FLUSH_INTERVAL" long:"event-flush-interval" default:"5s" description:"how often to flush events to the database"`
}
