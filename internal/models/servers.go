// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package models

import (
	"context"
	"time"

	"github.com/jinzhu/copier"
)

type ServerService interface {
	Upsert(ctx context.Context, r *Server) error
	Get(ctx context.Context, id string) (*Server, error)
	GetOptionsAdmin(ctx context.Context, id string) (*ServerOptionsAdmin, error)
	GetOptions(ctx context.Context, id string) (*ServerOptions, error)
	List(ctx context.Context, opts *ServerListOpts) ([]*Server, error)
}

// Server represents a Discord server/guild that we're connected to (or we
// have connected to in the past).
type Server struct {
	ID      string    `json:"id"      bson:"_id"     validate:"required"`
	Created time.Time `json:"created" bson:"created" validate:"required"`
	Updated time.Time `json:"updated" bson:"updated" validate:"required"`

	OptionsAdmin *ServerOptionsAdmin `json:"options_admin" bson:"options_admin" validate:"required"`
	Options      *ServerOptions      `json:"options"       bson:"options"       validate:"required"`

	Discord *ServerDiscordData `json:"discord" bson:"discord" validate:"required"`
	Events  []*ServerEvent     `json:"events"  bson:"events"`
}

func (s *Server) Public() (p *ServerPublic) {
	p = &ServerPublic{}
	err := copier.CopyWithOption(p, s, copier.Option{DeepCopy: true})
	if err != nil {
		panic(err)
	}
	return p
}

type ServerPublic struct {
	ID      string             `json:"id"`
	Created time.Time          `json:"created"`
	Updated time.Time          `json:"updated"`
	Discord *ServerDiscordData `json:"discord"`
}

func (s *Server) Validate() error {
	if s.Created.IsZero() {
		s.Created = time.Now()
	}

	if s.Updated.IsZero() {
		s.Updated = time.Now()
	}

	if s.OptionsAdmin == nil {
		s.OptionsAdmin = &ServerOptionsAdmin{
			Enabled: true,
		}
	}

	if s.Options == nil {
		s.Options = &ServerOptions{
			Enabled:    true,
			RegexMatch: "^(.*) \\+$",
		}
	}

	return errUseBuiltinValidator
}

// ServerOptionsAdmin are settings that will be applied by admins.
type ServerOptionsAdmin struct {
	Enabled            bool   `json:"enabled"              bson:"enabled"`
	Comment            string `json:"comment"              bson:"comment"`
	DefaultMaxChannels int    `json:"default_max_channels" bson:"default_max_channels" validate:"gte=0"`
	DefaultMaxClones   int    `json:"default_max_clones"   bson:"default_max_clones"   validate:"gte=0"`
}

// ServerOptions are settings that will be applied by users.
type ServerOptions struct {
	Enabled          bool   `json:"enabled"            bson:"enabled"`
	DefaultMaxClones int    `json:"default_max_clones" bson:"default_max_clones" validate:"gte=0"`
	RegexMatch       string `json:"regex_match"        bson:"regex_match"`
	ContactEmail     string `json:"contact_email"      bson:"contact_email"      validate:"omitempty,email"`
}

// ServerDiscordData represents the discord guild information returned by the
// gateway.
type ServerDiscordData struct {
	// Guild ID.
	ID string `json:"id" bson:"id" validate:"required"`
	// Guild name (2-100 chars, excl. trailing/leading spaces).
	Name string `json:"name" bson:"name" validate:"required"`
	// Enabled guild features.
	Features []string `json:"features" bson:"features"`
	// Icon hash.
	Icon string `json:"icon" bson:"icon"`
	// This is something we generate.
	IconUrl string `json:"icon_url" bson:"icon_url" validate:"omitempty,required_with=Icon,url"`
	// When the bot joined the guild.
	JoinedAt time.Time `json:"joined_at" bson:"joined_at"`
	// If the guild is considered large (to Discord standards).
	Large bool `json:"large" bson:"large"`
	// Total members in this guild.
	MemberCount int64 `json:"member_count" bson:"member_count" validate:"gte=0,required"`
	// User ID of the owner.
	OwnerID string `json:"owner_id" bson:"owner_id"`
	// Permissions of the bot on the server.
	Permissions DiscordPermissions `json:"permissions" bson:"permissions" validate:"required"`
	// Voice region (deprecated?).
	Region             string `json:"region" bson:"region"`
	SystemChannelFlags string `json:"system_channel_flags" bson:"system_channel_flags"`
}

// ServerEvent is a health status update event telling us if the server/guild
// is online & available, and if we're connected to it.
//
// TODO: auto-generate status if last status message is greater than X period
// of time?
type ServerEvent struct {
	// id of server to update status for.
	ID        string    `json:"id"        bson:"_id"       validate:"required"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp" validate:"required"`
	Type      string    `json:"type"      bson:"type"      validate:"required"`
	Message   string    `json:"message"   bson:"message"   validate:"required"`
	Healthy   bool      `json:"healthy"   bson:"healthy"   validate:"required"`
	Available bool      `json:"available" bson:"available" validate:"required"`
}

type ServerListOpts struct {
	OwnerID string
}
