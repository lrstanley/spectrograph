// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ServerService interface {
	Upsert(ctx context.Context, r *Server) error
	Get(ctx context.Context, id string) (*Server, error)
	List(ctx context.Context) ([]*Server, error)
}
type Server struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Created time.Time          `bson:"created"       json:"created"`
	Updated time.Time          `bson:"updated"       json:"updated"`
	Discord ServerDiscordData  `bson:"discord"       json:"discord"`
}

// TODO: auto-generate status if last status message is greater than X period
// of time?
type ServerStatusMessage struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Timestamp time.Time          `bson:"timestamp"     json:"timestamp"`
	Type      string             `bson:"type"          json:"type"`
	Message   string             `bson:"message"       json:"message"`
	Healthy   bool               `bson:"healthy"       json:"healthy"`
	Available bool               `bson:"available"     json:"available"`
}

type ServerDiscordData struct {
	ID          string `bson:"id"              json:"id"`          // Guild id.
	Name        string `bson:"name"            json:"name"`        // Guild name (2-100 chars, excl. trailing/leading spaces).
	Description string `bson:"description"     json:"description"` // Guild description.

	Features        []string           `bson:"features"         json:"features"`         // Enabled guild features.
	Icon            string             `bson:"icon"             json:"icon"`             // Icon hash.
	IconURL         string             `bson:"icon_url"         json:"icon_url"`         // We generate this.
	JoinedAt        time.Time          `bson:"joined_at"        json:"joined_at"`        // When the bot joined the guild.
	Large           bool               `bson:"large"            json:"large"`            // If the guild is considered large (to Discord standards).
	MemberCount     int                `bson:"member_count"     json:"member_count"`     // Total members in this guild.
	OwnerID         string             `bson:"owner_id"         json:"owner_id"`         // User ID of the owner.
	Permissions     DiscordPermissions `bson:"permissions"      json:"permissions"`      // Permissions of the bot on the server.
	PreferredLocale string             `bson:"preferred_locale" json:"preferred_locale"` // Preferred locale.
	Region          string             `bson:"region"           json:"region"`           // Voice region (deprecated?).

	PublicUpdatesChannelID string `bson:"public_updates_channel_id" json:"public_updates_channel_id"`
	SystemChannelID        string `bson:"system_channel_id" json:"system_channel_id"`
	SystemChannelFlags     string `bson:"system_channel_flags" json:"system_channel_flags"`
}
