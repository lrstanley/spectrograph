// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package authhandler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/lrstanley/pt"
	"github.com/lrstanley/spectrograph/internal/discordapi"
	"github.com/lrstanley/spectrograph/internal/httpware"
	"github.com/lrstanley/spectrograph/internal/models"
	"github.com/lrstanley/spectrograph/internal/util"
	"golang.org/x/oauth2"
)

const (
	discordBotAuthEndpoint = "https://discord.com/oauth2/authorize?client_id=%s&scope=bot&permissions=1049616"
)

type Handler struct {
	users   models.UserService
	servers models.ServerService
	session *scs.SessionManager
	config  *oauth2.Config
}

func New(users models.UserService, servers models.ServerService, config *oauth2.Config, session *scs.SessionManager) *Handler {
	return &Handler{
		users:   users,
		servers: servers,
		config:  config,
		session: session,
	}
}

func (h *Handler) Route(r chi.Router) {
	r.Get("/bot-authorize", h.getAuthorizeBot)
	r.Get("/redirect", h.getRedirect)
	r.Get("/callback", h.getCallback)
	r.With(httpware.AuthRequired(h.session)).Get("/self", h.getSelf)
	r.Get("/logout", h.getLogout)
}

func (h *Handler) getAuthorizeBot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, fmt.Sprintf(discordBotAuthEndpoint, h.config.ClientID), http.StatusTemporaryRedirect)
}

func (h *Handler) getRedirect(w http.ResponseWriter, r *http.Request) {
	if authed, _ := httpware.IsAuthed(h.session, r); authed {
		httpware.Error(w, r, http.StatusBadRequest, errors.New("already authenticated"))
		return
	}

	state := util.GenRandString(15)
	h.session.Put(r.Context(), "state", state)
	pt.JSON(w, r, pt.M{"auth_redirect": h.config.AuthCodeURL(
		state,
		// Provide AuthCodeOptions.
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("prompt", "none"),
	)})
}

func (h *Handler) getCallback(w http.ResponseWriter, r *http.Request) {
	if !httpware.IsDebug(r) {
		// Only check CSRF tokens if we're out of debug mode.
		state := h.session.GetString(r.Context(), "state")
		if state == "" {
			httpware.Error(w, r, http.StatusBadRequest, errors.New("session token not found, possible CSRF (or cookies disabled)? please try again"))
			return
		}

		h.session.Remove(r.Context(), "state")

		if state != r.FormValue("state") {
			httpware.Error(w, r, http.StatusBadRequest, errors.New("session token not found, possible CSRF (or cookies disabled)? please try again"))
			return
		}
	}
	h.session.Remove(r.Context(), "state")

	token, err := h.config.Exchange(r.Context(), r.FormValue("code"))
	if err != nil {
		httpware.Error(w, r, http.StatusUnauthorized, fmt.Errorf("error getting token: %w", err))
		return
	}
	client := h.config.Client(r.Context(), token)

	discordUser, discordServers, err := discordapi.FetchUser(client, token)
	if err != nil {
		// TODO: return statusCode so we can do unauthorized or similar?
		httpware.Error(w, r, http.StatusInternalServerError, fmt.Errorf("discord: %w", err))
		return
	}

	user, err := h.users.Get(r.Context(), discordUser.ID)
	if err != nil {
		if !models.IsNotFound(err) {
			httpware.Error(w, r, http.StatusInternalServerError, err)
			return
		}
		user = &models.User{}
	}

	user.Discord = discordUser
	user.JoinedServers = discordServers

	if err := h.users.Upsert(r.Context(), user); err != nil {
		httpware.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	// Prevent session fixation (and change in auth, or privilege, make a new
	// session token).
	if err = h.session.RenewToken(r.Context()); err != nil {
		httpware.Error(w, r, http.StatusInternalServerError, err)
		return
	}
	h.session.Put(r.Context(), models.SessionUserIDKey, user.ID)

	w.WriteHeader(http.StatusOK)
	pt.JSON(w, r, pt.M{"authenticated": true, "user": user.Public()})
}

func (h *Handler) getSelf(w http.ResponseWriter, r *http.Request) {
	user := httpware.MustGetUser(r)

	servers, err := h.servers.List(r.Context(), &models.ServerListOpts{
		OwnerID: user.ID,
	})
	if err != nil {
		httpware.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	serversPublic := make([]*models.ServerPublic, len(servers))
	for i, server := range servers {
		serversPublic[i] = server.Public()
	}

	pt.JSON(w, r, pt.M{"authenticated": true, "user": user.Public(), "servers": serversPublic})
}

func (h *Handler) getLogout(w http.ResponseWriter, r *http.Request) {
	h.session.Remove(r.Context(), models.SessionUserIDKey)

	// TODO: https://discord.com/api/oauth2/token/revoke ? only do if removing
	// the account?

	w.WriteHeader(http.StatusOK)
	pt.JSON(w, r, pt.M{"authenticated": false})
}

// func authRefreshToken(ctx context.Context, config *oauth2.Config, refreshToken string) (*oauth2.Token, error) {
// 	token := &oauth2.Token{RefreshToken: refreshToken}
// 	ts := config.TokenSource(ctx, token)

// 	return ts.Token()
// }
