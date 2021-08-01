// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package workerhandler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lrstanley/spectrograph/internal/worker"
	"github.com/twitchtv/twirp"
)

type Handler struct {
	version string
	key     string
	rpc     worker.TwirpServer
}

func New(serverVersion, secretKey string) *Handler {
	return &Handler{
		version: serverVersion,
		key:     secretKey,
		rpc:     worker.NewWorkerServer(&Server{}, twirp.WithServerPathPrefix(worker.PathPrefix)),
	}
}

func (h *Handler) Route(r chi.Router) {
	r.With(h.validate).Mount("/", h.rpc)
}

func (h *Handler) validate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Version checks.
		clientVersion := r.Header.Get("X-Api-Version")
		if clientVersion == "" {
			twerr := twirp.NewError(twirp.FailedPrecondition, "api version not specified")
			twirp.WriteError(w, twerr)
			return
		} else if clientVersion != h.version {
			twerr := twirp.NewError(twirp.FailedPrecondition, fmt.Sprintf("server (%q) and client (%q) version mismatch", h.version, clientVersion))
			twirp.WriteError(w, twerr)
			return
		}

		// Authentication checks.
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
	return &worker.HealthResp{Ready: true}, nil
}
