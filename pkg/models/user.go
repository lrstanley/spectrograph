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

	Discord struct {
		LastLogin time.Time `bson:"last_login" json:"last_login"`

		ID           string    `bson:"id"            json:"id"         validate:"required"`
		Email        string    `bson:"email"         json:"email"      validate:"required"`
		Name         string    `bson:"name"          json:"name"       validate:"required"`
		NameID       string    `bson:"name_id"       json:"name_id"    validate:"required"`
		AvatarURL    string    `bson:"avatar_url"    json:"avatar_url" validate:"required"`
		AccessToken  string    `bson:"access_token"  json:"-"          validate:"required"`
		RefreshToken string    `bson:"refresh_token" json:"-"          validate:"required"`
		ExpiresAt    time.Time `bson:"expires_at"    json:"-"          validate:"required"`

		RawData map[string]interface{} `bson:"raw_data" json:"-"`
	} `bson:"discord" json:"discord"`
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
