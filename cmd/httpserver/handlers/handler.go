// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

// Package handlers mounts the HTTP API (generated entrest endpoints, and custom
// endpoints), authentication, and health routes.
package handlers

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/lrstanley/chix/v2"
	"github.com/lrstanley/chix/xauth/v2"
	"github.com/lrstanley/spectrograph/internal/database"
	"github.com/lrstanley/spectrograph/internal/database/ent"
	"github.com/lrstanley/spectrograph/internal/database/ent/rest"
	"github.com/lrstanley/spectrograph/internal/database/schema"
	"github.com/lrstanley/spectrograph/internal/models"
)

// Config configures the API handlers.
type Config struct {
	DB *ent.Client

	// AuthService is the authentication service used for the OAuth2 flow.
	AuthService xauth.Service[ent.User, int]

	// SessionStorage is the session storage used for authentication.
	SessionStorage sessions.Store

	// DiscordClientID is used for bot invite links.
	DiscordClientID string

	// BasePath is the path where the API is mounted (e.g. "/-").
	BasePath string

	// Shutdown is cancelled when the server is shutting down, which terminates
	// long-lived (SSE) streams.
	Shutdown context.Context //nolint:containedctx

	// GuildEvents, if set, wakes guild event streams when new events are
	// created in this process. Streams still poll (see StreamPollInterval), for
	// events created by other processes.
	GuildEvents *database.GuildEventNotifier

	// StreamPollInterval is how often guild event streams poll for new events,
	// and re-validate access. Defaults to 3 seconds.
	StreamPollInterval time.Duration
}

type handler struct {
	cfg *Config
	srv *rest.Server
}

// New returns the API handlers. Panics if the generated REST server can't be
// initialized.
func New(cfg Config) *handler { //nolint:revive // unexported return is intentional.
	h := &handler{cfg: &cfg}

	var err error
	h.srv, err = rest.NewServer(cfg.DB, &rest.ServerConfig{
		BasePath: cfg.BasePath,
		// Route errors through chix (and [models.ErrorResolver]), so all API errors
		// share the same response format (entrest's default handler can't include
		// "field", and uses different status codes and masking).
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, _ rest.Operation, err error) {
			chix.Error(w, r, err)
		},
	})
	if err != nil {
		panic(err)
	}
	return h
}

// Route mounts all API routes onto the provided router.
func (h *handler) Route(r chi.Router) {
	r.Get("/healthy", h.getHealthy)
	r.Get("/invite", h.getBotInvite)
	r.Get("/invite/{guildSnowflake:[0-9]+}", h.getBotInvite)

	// Mount these separately, so they don't require authentication.
	r.Get("/docs", h.srv.Docs)
	r.Get("/openapi.json", h.srv.Spec)

	r.With(useLoginErrorRedirect, stripAuthState).Mount("/auth", xauth.NewGothHandler(&xauth.GothConfig[ent.User, int]{
		Service:             h.cfg.AuthService,
		SessionStorage:      h.cfg.SessionStorage,
		DisableSelfEndpoint: true, // See "/self" instead.
	}))

	r.Group(func(r chi.Router) {
		r.Use(xauth.UseAuthRequired[ent.User](), requireNotBanned)

		r.Get("/self", h.getSelf)
		r.Delete("/self", h.deleteSelf)

		r.With(requireSystemAdmin).Post("/users/{id}/ban", h.postBanUser)
		r.With(requireSystemAdmin).Post("/users/{id}/unban", h.postUnbanUser)

		r.Get("/guilds/{id}/guild-events/stream", h.streamGuildEvents())

		r.With(rejectAggregateSorts).Route("/", h.srv.Handler) // Catch-all for generated endpoints.
	})
}

// requireNotBanned is a middleware which rejects banned users with a 403 (and
// [models.ErrUserBanned]), so the frontend can tell them apart from signed out
// users. Privacy rules also treat banned users as unauthenticated.
func requireNotBanned(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if u := xauth.IdentFromContext[ent.User](r.Context()); u != nil && u.Banned {
			chix.ErrorWithCode(w, r, http.StatusForbidden, models.ErrUserBanned)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// rejectAggregateSorts rejects sorting by edge aggregates (e.g. "admins.count"),
// which entrest generates for every edge. Aggregates are computed without the
// edge's privacy rules, so they could reveal e.g. how many DEBUG events, or other
// admins, a guild has. The sort field is also decoded from "Sorted.sort".
func rejectAggregateSorts(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for key, values := range r.URL.Query() {
			if key != "sort" && !strings.HasSuffix(key, ".sort") {
				continue
			}
			for _, s := range values {
				if strings.HasSuffix(s, ".count") || strings.HasSuffix(s, ".sum") {
					chix.ErrorWithCode(w, r, http.StatusBadRequest, errors.New("sorting by edge counts or sums isn't supported"))
					return
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}

// requireSystemAdmin is a middleware which requires the authenticated user to be
// a system admin. Note that privacy rules are still enforced at the database
// layer, this simply allows failing fast.
func requireSystemAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !schema.IsSystemAdmin(r.Context()) {
			chix.ErrorWithCode(w, r, http.StatusForbidden, errors.New("system admin privileges required"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *handler) getHealthy(w http.ResponseWriter, r *http.Request) {
	if err := database.Ping(r.Context(), h.cfg.DB); err != nil {
		chix.Error(w, r, err)
		return
	}
	chix.JSON(w, r, http.StatusOK, chix.M{"status": "ok"})
}
