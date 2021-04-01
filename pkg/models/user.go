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

	Discord UserAuthDiscord `bson:"discord" json:"discord"`
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

func (u *User) Public() *UserPublic {
	return &UserPublic{
		ID:             u.ID,
		AccountCreated: u.AccountCreated,
		AccountUpdated: u.AccountUpdated,

		Username:      u.Discord.Username,
		Discriminator: u.Discord.Discriminator,
		Avatar:        u.Discord.Avatar,
	}
}

type UserPublic struct {
	ID             primitive.ObjectID `json:"id"`
	AccountCreated time.Time          `json:"account_created"`
	AccountUpdated time.Time          `json:"account_updated"`

	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
}

// See also: https://discord.com/developers/docs/resources/user#user-object
type UserAuthDiscord struct {
	LastLogin time.Time `bson:"last_login" json:"last_login"`

	// Required dependencies.
	ID            string `bson:"id"            json:"id"            validate:"required"` // the user's id
	Username      string `bson:"username"      json:"username"      validate:"required"` // the user's username, not unique across the platform
	Discriminator string `bson:"discriminator" json:"discriminator" validate:"required"` // the user's 4-digit discord-tag
	Email         string `bson:"email"         json:"email"         validate:"required"` // the user's email
	Avatar        string `bson:"avatar"        json:"avatar"        validate:"required"` // the user's avatar hash

	// Additional parameters provided by the API.
	Locale      string `bson:"locale"       json:"locale"`       // the user's chosen language option
	Bot         bool   `bson:"bot"          json:"bot"`          // whether the user belongs to an OAuth2 application
	System      bool   `bson:"system"       json:"system"`       // whether the user is an Official Discord System user (part of the urgent message system)
	MFAEnabled  bool   `bson:"mfa_enabled"  json:"mfa_enabled"`  // whether the user has two factor enabled on their account
	Verified    bool   `bson:"verified"     json:"verified"`     // whether the email on this account has been verified
	Flags       int    `bson:"flags"        json:"flags"`        // the flags on a user's account
	PremiumType int    `bson:"premium_type" json:"premium_type"` // the type of Nitro subscription on a user's account
	PublicFlags int    `bson:"public_flags" json:"public_flags"` // the public flags on a user's account

	AccessToken  string    `bson:"access_token"  json:"-" validate:"required"`
	RefreshToken string    `bson:"refresh_token" json:"-" validate:"required"`
	ExpiresAt    time.Time `bson:"expires_at"    json:"-" validate:"required"`
}
