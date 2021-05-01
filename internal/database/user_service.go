// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package database

import (
	"context"
	"time"

	"github.com/lrstanley/spectrograph/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// userService satisfies the models.UserService interface.
type userService struct {
	store *mongoStore
}

// Ensure struct matches necessary interface.
var _ models.UserService = (*userService)(nil)

// NewUserService returns a new userService that satisfies the models.UserService interface.
func (s *mongoStore) NewUserService() models.UserService {
	return &userService{store: s}
}

func (s *userService) Upsert(ctx context.Context, user *models.User) (err error) {
	if err = user.Validate(); err != nil {
		return err
	}

	var res models.User

	// TODO: if multi-auth, do "$or".
	err = s.store.user.FindOne(ctx,
		bson.M{"discord.id": user.Discord.ID},
		options.FindOne().SetProjection(bson.M{"_id": 1}), // TODO: bson.D{}?
	).Decode(&res)

	if err == nil {
		user.ID = res.ID
	} else if err == mongo.ErrNoDocuments {
		user.ID = primitive.NewObjectID()
		user.AccountCreated = time.Now()
	} else {
		return err
	}

	_, err = s.store.user.UpdateOne(
		ctx,
		bson.M{"discord.id": user.Discord.ID}, bson.M{"$set": user},
		options.Update().SetUpsert(true),
	)
	return errorWrapper(err)
}

func (s *userService) Get(ctx context.Context, id string) (user *models.User, err error) {
	var oid primitive.ObjectID
	if oid, err = primitive.ObjectIDFromHex(id); err != nil {
		return nil, ErrInvalidObjectID
	}

	err = s.store.user.FindOne(ctx, bson.M{"_id": oid}).Decode(&user)
	return user, errorWrapper(err)
}

func (s *userService) List(ctx context.Context) (users []*models.User, err error) {
	var cursor *mongo.Cursor

	cursor, err = s.store.user.Find(ctx, bson.M{})
	if err != nil {
		return nil, errorWrapper(err)
	}
	defer cursor.Close(ctx)

	return users, errorWrapper(cursor.All(ctx, &users))
}
