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
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	GithubID int                `bson:"github_id" validate:"required" json:"github_id"`
	Token    string             `bson:"token" validate:"required" json:"-"`

	AvatarURL string `bson:"avatar_url" validate:"required" json:"avatar_url"`
	Username  string `bson:"username" validate:"required" json:"username"`
	Name      string `bson:"name" validate:"required" json:"name"`
	Email     string `bson:"email" validate:"required,email" json:"-"`

	AccountCreated time.Time `bson:"account_created" json:"-"`
	AccountUpdated time.Time `bson:"account_updated" json:"-"`
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
