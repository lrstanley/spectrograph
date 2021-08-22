// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package serverhandler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kr/pretty"
	"github.com/lrstanley/pt"
	"github.com/lrstanley/spectrograph/internal/httpware"
	"github.com/lrstanley/spectrograph/internal/models"
)

type Handler struct {
	servers models.ServerService
}

func New(servers models.ServerService) *Handler {
	return &Handler{
		servers: servers,
	}
}

func (h *Handler) Route(r chi.Router) {
	r.Get("/servers", h.list)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	user := httpware.MustGetUser(r)

	servers, err := h.servers.List(r.Context(), &models.ServerListOpts{
		OwnerID: user.ID,
	})
	if err != nil {
		httpware.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	pretty.Print(servers)

	pt.JSON(w, r, pt.M{"servers": servers})
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	user := httpware.MustGetUser(r)
	id := chi.URLParam(r, "id")
	if id == "" {
		httpware.Error(w, r, http.StatusBadRequest, errors.New("no id provided"))
		return
	}

	// Check if the ID in mention is in the users guild list.
	if !user.HasPermissions(id) {
		httpware.Error(w, r, http.StatusForbidden, errors.New("access denied"))
		return
	}

	server, err := h.servers.Get(r.Context(), id)
	if err != nil {
		if models.IsNotFound(err) {
			httpware.Error(w, r, http.StatusNotFound, nil)
			return
		}

		httpware.Error(w, r, http.StatusInternalServerError, err)
	}

	w.WriteHeader(http.StatusOK)
	pt.JSON(w, r, server)
}
