// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package workerhandler

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lrstanley/spectrograph/internal/worker"
	"github.com/twitchtv/twirp"
)

type Handler struct {
	key string
	rpc worker.TwirpServer
}

func New(secretKey string) *Handler {
	return &Handler{
		key: secretKey,
		rpc: worker.NewWorkerServer(&Server{}, twirp.WithServerPathPrefix(worker.PathPrefix)),
	}
}

func (h *Handler) Route(r chi.Router) {
	r.With(h.keyHeaderRequired).Mount("/", h.rpc)
}

func (h *Handler) keyHeaderRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Api-Key") != h.key {
			twerr := twirp.NewError(twirp.Unauthenticated, "invalid token provided")
			twirp.WriteError(w, twerr)
			return
		}

		next.ServeHTTP(w, r)
	})
}

type Server struct{}

func (s *Server) Health(ctx context.Context, _ *worker.NoArgs) (*worker.HealthResp, error) {
	return &worker.HealthResp{Version: "v1.1.1", Ready: true}, nil
}
