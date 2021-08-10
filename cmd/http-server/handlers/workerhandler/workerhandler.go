// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package workerhandler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Route(r chi.Router) {
	r.Get("/health", h.getHealth)
}

func (h *Handler) getHealth(w http.ResponseWriter, r *http.Request) {
	// By the time this endpoint is accessible, we should have initiated all
	// of the necessary background connections/services and be considered healthy.
	w.WriteHeader(http.StatusOK)
}
