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

type Handler struct {
	svc rpc.TwirpServer
}

func New() *Handler {
	return &Handler{svc: rpc.NewWorkerServer(&Server{}, twirp.WithServerPathPrefix(rpc.PathPrefix))}
}

func (h *Handler) Route(r chi.Router) {
	r.Mount("/", h.svc)
}

// Validate that the Server interface satisfies the RPC interface.
var _ rpc.Worker = (*Server)(nil)

type Server struct{}

func (s *Server) Health(ctx context.Context, _ *empty.Empty) (*models.WorkerHealth, error) {
	// By the time this endpoint is accessible, we should have initiated all
	// of the necessary background connections/services and be considered healthy.
	return &models.WorkerHealth{Ready: true}, nil
}
