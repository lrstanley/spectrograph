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

// serverService satisfies the models.ServerService interface.
type serverService struct {
	store *mongoStore
}

// Ensure struct matches necessary interface.
var _ models.ServerService = (*serverService)(nil)

// NewServerService returns a new ServerService that satisfies the models.ServerService interface.
func (s *mongoStore) NewServerService() models.ServerService {
	return &serverService{store: s}
}

func (s *serverService) Upsert(ctx context.Context, server *models.Server) (err error) {
	if server.Id == "" {
		server.Id = primitive.NewObjectID().Hex()
	}

	if err = models.Validate(server); err != nil {
		return errorWrapper(err)
	}

	_, err = s.store.servers.UpdateOne(
		ctx,
		bson.M{"_id": server.Id}, bson.M{"$set": server},
		options.Update().SetUpsert(true),
	)
	return errorWrapper(err)
}

func (s *serverService) Get(ctx context.Context, id string) (server *models.Server, err error) {
	err = s.store.servers.FindOne(ctx, bson.M{"_id": id}).Decode(&server)
	return server, errorWrapper(err)
}

func (s *serverService) GetByDiscordID(ctx context.Context, id string) (server *models.Server, err error) {
	err = s.store.servers.FindOne(ctx, bson.M{"discord.id": id}).Decode(&server)
	return server, errorWrapper(err)
}

func (s *serverService) List(ctx context.Context, opts *models.ServerListOpts) (servers []*models.Server, err error) {
	if opts == nil {
		opts = &models.ServerListOpts{}
	}

	var cursor *mongo.Cursor

	filter := bson.M{}

	if opts.OwnerID != "" {
		var user models.User

		err := s.store.users.FindOne(
			ctx,
			bson.M{"_id": opts.OwnerID},
			options.FindOne().SetProjection(bson.M{"discord_servers.id": 1}),
		).Decode(&user)
		if err != nil {
			return nil, errorWrapper(err)
		}

		serverIds := []string{}
		for _, server := range user.DiscordServers {
			serverIds = append(serverIds, server.ID)
		}
		filter["discord.id"] = bson.M{"$in": serverIds}
	}

	cursor, err = s.store.servers.Find(ctx, filter)
	if err != nil {
		return nil, errorWrapper(err)
	}
	defer cursor.Close(ctx)

	return servers, errorWrapper(cursor.All(ctx, &servers))
}
