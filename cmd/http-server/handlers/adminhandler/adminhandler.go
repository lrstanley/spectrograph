// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package adminhandler

import (
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/lrstanley/spectrograph/pkg/httpware"
	"github.com/lrstanley/spectrograph/pkg/models"
)

type Handler struct {
	users   models.UserService
	session *scs.SessionManager
}

func New(users models.UserService, session *scs.SessionManager) *Handler {
	return &Handler{
		users:   users,
		session: session,
	}
}

func (h *Handler) Route(r chi.Router) {
	r.Use(httpware.AdminRequired(h.session))

	// r.Get("/redirect", h.getRedirect)
	// r.Get("/callback", h.getCallback)
	// r.Get("/self", h.getSelf)
	// r.Get("/logout", h.getLogout)
}
