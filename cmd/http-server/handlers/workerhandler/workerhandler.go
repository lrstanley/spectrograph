// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package workerhandler

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/lrstanley/spectrograph/internal/models"
	"github.com/lrstanley/spectrograph/internal/rpc"
	"github.com/twitchtv/twirp"
)

// Validate that the Handler interface satisfies the necessary RPC interface.
var _ rpc.Worker = (*Handler)(nil)

type Handler struct {
	svcRpc     rpc.TwirpServer
	svcServers models.ServerService
}

func New(svcServers models.ServerService) *Handler {
	handler := &Handler{
		svcServers: svcServers,
	}

	handler.svcRpc = rpc.NewWorkerServer(handler, twirp.WithServerPathPrefix(rpc.PathPrefix))
	return handler
}

func (h *Handler) Route(r chi.Router) {
	r.Mount("/", h.svcRpc)
}

func (h *Handler) Health(ctx context.Context, _ *empty.Empty) (*models.WorkerHealth, error) {
	// By the time this endpoint is accessible, we should have initiated all
	// of the necessary background connections/services and be considered healthy.
	return &models.WorkerHealth{Ready: true}, nil
}

func (h *Handler) UpdateServer(ctx context.Context, data *models.ServerDiscordData) (*empty.Empty, error) {
	// json.Default.Encode(server)
	// if err := models.Validate(server); err != nil {
	// 	return nil, rpc.ValidationError(err)
	// }

	server, err := h.svcServers.GetByDiscordID(ctx, data.Id)
	if err != nil {
		if !models.IsNotFound(err) {
			return nil, rpc.WrapError(err)
		}
		err = nil
		server = &models.Server{}
	}

	server.Discord = data

	return &empty.Empty{}, rpc.WrapError(h.svcServers.Upsert(ctx, server))
}

func (h *Handler) UpdateServerStatus(ctx context.Context, status *models.ServerStatus) (*empty.Empty, error) {
	if err := models.Validate(status); err != nil {
		return nil, rpc.WrapError(err)
	}
	return nil, nil
}
