// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package database

import (
	"context"

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
	if user.ID == "" {
		user.ID = primitive.NewObjectID().Hex()
	}

	if err = models.Validate(user); err != nil {
		return err
	}

	_, err = s.store.users.UpdateOne(
		ctx,
		bson.M{"discord.id": user.Discord.ID}, bson.M{"$set": user},
		options.Update().SetUpsert(true),
	)
	return errorWrapper(err)
}

func (s *userService) Get(ctx context.Context, id string) (user *models.User, err error) {
	err = s.store.users.FindOne(ctx, bson.M{"$or": []bson.M{
		bson.M{"_id": id},
		bson.M{"discord.id": id},
	}}).Decode(&user)
	return user, errorWrapper(err)
}

func (s *userService) List(ctx context.Context) (users []*models.User, err error) {
	var cursor *mongo.Cursor

	cursor, err = s.store.users.Find(ctx, bson.M{})
	if err != nil {
		return nil, errorWrapper(err)
	}
	defer cursor.Close(ctx)

	return users, errorWrapper(cursor.All(ctx, &users))
}
