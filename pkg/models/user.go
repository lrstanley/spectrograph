// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/go-playground/validator.v9"
)

type UserService interface {
	Upsert(ctx context.Context, r *User) error
	Get(ctx context.Context, id string) (*User, error)
	List(ctx context.Context) ([]*User, error)
}

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"   json:"id"`
	AccountCreated time.Time          `bson:"account_created" json:"account_created"`
	AccountUpdated time.Time          `bson:"account_updated" json:"account_updated"`

	Discord        UserAuthDiscord     `bson:"discord"         json:"discord"`
	DiscordServers []UserDiscordServer `bson:"discord_servers" json:"discord_servers"`
}

type UserPublic struct {
	ID             primitive.ObjectID `json:"id"`
	AccountCreated time.Time          `json:"account_created"`
	AccountUpdated time.Time          `json:"account_updated"`

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
	ID          string   `bson:"id"              json:"id"`       // guild id.
	Name        string   `bson:"name"            json:"name"`     // guild name (2-100 chars, excl. trailing/leading spaces).
	Owner       bool     `bson:"owner"           json:"owner"`    // true if the user is the owner of the guild
	Features    []string `bson:"features"        json:"features"` // enabled guild features.
	Icon        string   `bson:"icon"            json:"icon"`     // icon hash.
	IconURL     string   `bson:"icon_url"        json:"icon_url"`
	Permissions string   `bson:"permissions_new" json:"permissions_new"` // permissions for the user (excludes overrides).
}

// See also: https://discord.com/developers/docs/resources/user#user-object
type UserAuthDiscord struct {
	LastLogin time.Time `bson:"last_login" json:"last_login"`

	// Required dependencies.
	ID            string `bson:"id"            json:"id"            validate:"required"` // the users id
	Username      string `bson:"username"      json:"username"      validate:"required"` // the users username, not unique across the platform
	Discriminator string `bson:"discriminator" json:"discriminator" validate:"required"` // the users 4-digit discord-tag
	Email         string `bson:"email"         json:"email"         validate:"required"` // the users email
	Avatar        string `bson:"avatar"        json:"avatar"`                            // the users avatar url
	AvatarURL     string `bson:"avatar_url"    json:"avatar_url"    validate:"required"` // the users avatar hash

	// Additional parameters provided by the API.
	Locale      string `bson:"locale"       json:"locale"`       // the users chosen language option
	Bot         bool   `bson:"bot"          json:"bot"`          // whether the user belongs to an OAuth2 application
	System      bool   `bson:"system"       json:"system"`       // whether the user is an Official Discord System user (part of the urgent message system)
	MFAEnabled  bool   `bson:"mfa_enabled"  json:"mfa_enabled"`  // whether the user has two factor enabled on their account
	Verified    bool   `bson:"verified"     json:"verified"`     // whether the email on this account has been verified
	Flags       int    `bson:"flags"        json:"flags"`        // the flags on a users account
	PremiumType int    `bson:"premium_type" json:"premium_type"` // the type of Nitro subscription on a users account
	PublicFlags int    `bson:"public_flags" json:"public_flags"` // the public flags on a users account

	AccessToken  string    `bson:"access_token"  json:"-" validate:"required"`
	RefreshToken string    `bson:"refresh_token" json:"-" validate:"required"`
	ExpiresAt    time.Time `bson:"expires_at"    json:"-" validate:"required"`
}
