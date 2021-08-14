// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package models

import (
	"context"
	"time"

	"gopkg.in/go-playground/validator.v9"
)

type UserService interface {
	Upsert(ctx context.Context, r *User) error
	Get(ctx context.Context, id string) (*User, error)
	List(ctx context.Context) ([]*User, error)
}

type User struct {
	ID             string    `bson:"_id,omitempty"   json:"id"`
	AccountCreated time.Time `bson:"account_created" json:"account_created"`
	AccountUpdated time.Time `bson:"account_updated" json:"account_updated"`

	Discord        UserAuthDiscord     `bson:"discord"         json:"discord"`
	DiscordServers []UserDiscordServer `bson:"discord_servers" json:"discord_servers"`
}

type UserPublic struct {
	ID             string    `json:"id"`
	AccountCreated time.Time `json:"account_created"`
	AccountUpdated time.Time `json:"account_updated"`

	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
	AvatarURL     string `json:"avatar_url"`

	Servers []UserDiscordServer `json:"servers"`
}

func (u *User) Public() *UserPublic {
	return &UserPublic{
		ID:             u.ID,
		AccountCreated: u.AccountCreated,
		AccountUpdated: u.AccountUpdated,

		Username:      u.Discord.Username,
		Discriminator: u.Discord.Discriminator,
		Avatar:        u.Discord.Avatar,
		AvatarURL:     u.Discord.AvatarURL,

		Servers: u.DiscordServers,
	}
}

func (r *User) Validate() error {
	err := validate.Struct(r)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			panic(err)
		}
		return err
	}
	return nil
}

// {
//     "features": [
//         "WELCOME_SCREEN_ENABLED", "NEWS", "COMMUNITY"
//     ],
//     "icon":"3a0892e2c181bd2fe877e6c4341d163e",
//     "id":"679506910449500196",
//     "name":"bytecord",
//     "owner":true,
//     "permissions":2.147483647e+09,
//     "permissions_new":"8589934591"
// }
// https://discord.com/developers/docs/resources/guild#guild-object
type UserDiscordServer struct {
	ID          string             `bson:"id"              json:"id"`       // Guild id.
	Name        string             `bson:"name"            json:"name"`     // Guild name (2-100 chars, excl. trailing/leading spaces).
	Owner       bool               `bson:"owner"           json:"owner"`    // True if the user is the owner of the guild
	Features    []string           `bson:"features"        json:"features"` // Enabled guild features.
	Icon        string             `bson:"icon"            json:"icon"`     // Icon hash.
	IconURL     string             `bson:"icon_url"        json:"icon_url"`
	Permissions DiscordPermissions `bson:"permissions_new" json:"permissions_new"` // Permissions for the user (excludes overrides).

	Admin bool `bson:"admin" json:"admin"` // Pulled from permissions.
}

// See also: https://discord.com/developers/docs/resources/user#user-object
type UserAuthDiscord struct {
	LastLogin time.Time `bson:"last_login" json:"last_login"`

	// Required dependencies.
	ID            string `bson:"id"            json:"id"            validate:"required"` // The users id.
	Username      string `bson:"username"      json:"username"      validate:"required"` // The users username, not unique across the platform.
	Discriminator string `bson:"discriminator" json:"discriminator" validate:"required"` // The users 4-digit discord-tag.
	Email         string `bson:"email"         json:"email"         validate:"required"` // The users email.
	Avatar        string `bson:"avatar"        json:"avatar"`                            // The users avatar url.
	AvatarURL     string `bson:"avatar_url"    json:"avatar_url"    validate:"required"` // The users avatar hash.

	// Additional parameters provided by the API.
	Locale      string `bson:"locale"       json:"locale"`       // The users chosen language option.
	Bot         bool   `bson:"bot"          json:"bot"`          // Whether the user belongs to an OAuth2 application.
	System      bool   `bson:"system"       json:"system"`       // Whether the user is an Official Discord System user (part of the urgent message system).
	MFAEnabled  bool   `bson:"mfa_enabled"  json:"mfa_enabled"`  // Whether the user has two factor enabled on their account.
	Verified    bool   `bson:"verified"     json:"verified"`     // Whether the email on this account has been verified.
	Flags       int    `bson:"flags"        json:"flags"`        // The flags on a users account.
	PremiumType int    `bson:"premium_type" json:"premium_type"` // The type of Nitro subscription on a users account.
	PublicFlags int    `bson:"public_flags" json:"public_flags"` // The public flags on a users account.

	AccessToken  string    `bson:"access_token"  json:"-" validate:"required"`
	RefreshToken string    `bson:"refresh_token" json:"-" validate:"required"`
	ExpiresAt    time.Time `bson:"expires_at"    json:"-" validate:"required"`
}
