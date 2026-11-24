// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

// Package models contains shared configuration, Discord API models, and errors.
package models

import (
	"cmp"
	"fmt"
	"net/url"
	"time"
)

type Flags struct {
	HTTP     ConfigHTTP     `embed:"" group:"HTTP Server options" prefix:"http."     envprefix:"HTTP_"`
	Database ConfigDatabase `embed:"" group:"Database options"    prefix:"database." envprefix:"DATABASE_"`
	Discord  ConfigDiscord  `embed:"" group:"Discord options"     prefix:"discord."  envprefix:"DISCORD_"`
	Worker   ConfigWorker   `embed:"" group:"Worker options"      prefix:"worker."   envprefix:"WORKER_"`
}

// ConfigDatabase holds the database configuration.
type ConfigDatabase struct {
	URL      *url.URL `env:"URL"      name:"url"      required:"" help:"database connection url (e.g. postgres://user:pass@host:5432/db)"`
	Username string   `env:"USERNAME" name:"username" help:"database username if not specified via the URL"`
	Password string   `env:"PASSWORD" name:"password" help:"database password if not specified via the URL"`
	Database string   `env:"DATABASE" name:"database" help:"database name if not specified via the URL"`
}

// AfterApply merges the explicitly provided credentials and database name into
// [ConfigDatabase.URL]. Called by kong after parsing.
func (c *ConfigDatabase) AfterApply() error {
	if c.URL == nil {
		return nil
	}
	if c.Database != "" {
		c.URL.Path = "/" + c.Database
	}
	if c.Username != "" || c.Password != "" {
		password, _ := c.URL.User.Password()
		c.URL.User = url.UserPassword(cmp.Or(c.Username, c.URL.User.Username()), cmp.Or(c.Password, password))
	}
	return nil
}

// ConfigHTTP are configurations specifically utilized by the HTTP service.
type ConfigHTTP struct {
	BaseURL        *url.URL `env:"BASE_URL"        name:"base-url"        default:"http://localhost:8080" help:"base url for the HTTP server (used for OAuth callbacks)"`
	BindAddr       string   `env:"BIND_ADDR"       name:"bind-addr"       default:":8080" required:"" help:"ip:port pair to bind to"`
	TrustedProxies []string `env:"TRUSTED_PROXIES" name:"trusted-proxies" help:"CIDR ranges that we trust the X-Forwarded-For header from (addl opts: local, *, cloudflare, and/or custom header to use)"`
	ValidationKey  string   `env:"VALIDATION_KEY"  name:"validation-key"  required:"" help:"hex-encoded key used to validate session cookies (32 or 64 bytes)"`
	EncryptionKey  string   `env:"ENCRYPTION_KEY"  name:"encryption-key"  required:"" help:"hex-encoded key used to encrypt session cookies (32 bytes)"`
}

// Validate ensures [ConfigHTTP.BaseURL] is an absolute http(s) URL. Called by
// kong after parsing.
func (c *ConfigHTTP) Validate() error {
	if c.BaseURL == nil || (c.BaseURL.Scheme != "http" && c.BaseURL.Scheme != "https") || c.BaseURL.Host == "" {
		return fmt.Errorf("base url %q must be an absolute http or https URL", c.BaseURL)
	}
	return nil
}

// ConfigDiscord are configurations specifically utilized for interacting with
// Discord.
type ConfigDiscord struct {
	ClientID     string   `env:"CLIENT_ID"     name:"client-id"     required:"" help:"Discord client ID"`
	ClientSecret string   `env:"CLIENT_SECRET" name:"client-secret" required:"" help:"Discord client secret"`
	Admins       []string `env:"ADMINS"        name:"admins"        help:"Discord user IDs of the users you want to be system admins"`
}

// ConfigWorker are configurations specifically utilized by the Discord bot worker.
type ConfigWorker struct {
	Enabled bool `env:"ENABLED" name:"enabled" default:"true" negatable:"" help:"run the Discord bot worker (disable to only run the HTTP server)"`

	// https://discord.com/developers/docs/topics/gateway#sharding
	// Note: Shard ID 0 will be the only one to receive DMs.
	ShardIDs           []uint        `env:"SHARD_ID"             name:"shard-id"             default:"0"  help:"shard ID(s) of this specific worker (from 0)"`
	NumShards          uint          `env:"NUM_SHARDS"           name:"num-shards"           default:"1"  help:"number of total shards"`
	BotToken           string        `env:"BOT_TOKEN"            name:"bot-token"                         help:"Discord bot token (required when the worker is enabled)"`
	EventFlushInterval time.Duration `env:"EVENT_FLUSH_INTERVAL" name:"event-flush-interval" default:"5s" help:"how often to flush guild events to the database"`

	// Slash commands are registered once per start, by the worker running shard
	// 0 (the only one guaranteed to exist).
	RegisterCommands bool   `env:"REGISTER_COMMANDS" name:"register-commands" default:"true" negatable:"" help:"register the bot's slash commands with Discord at startup (only done by the worker running shard 0)"`
	CommandsGuildID  string `env:"COMMANDS_GUILD_ID" name:"commands-guild-id" help:"register the slash commands in this guild only (updates instantly, for development), instead of globally (which can take a while to update)"`
}
