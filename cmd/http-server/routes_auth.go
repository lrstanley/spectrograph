// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/go-github/github"
	"github.com/lrstanley/pt"
	"github.com/lrstanley/spectrograph/pkg/http/helpers"
	"github.com/lrstanley/spectrograph/pkg/models"
)

// https://github.com/markbates/goth
// https://github.com/volatiletech/authboss

func registerAuthRoutes(r chi.Router) {
	r.Get("/api/v1/auth/github/redirect", func(w http.ResponseWriter, r *http.Request) {
		state := helpers.GenRandString(15)

		if id := session.GetString(r.Context(), "user_id"); id != "" {
			_, err := svcUsers.Get(r.Context(), id)
			if err != nil {
				if models.IsNotFound(err) {
					session.Remove(r.Context(), "user_id")
					goto redirect
				}

				helpers.HTTPError(w, r, http.StatusServiceUnavailable, errors.New("unable to query user"))
				return
			}

			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

	redirect:
		session.Put(r.Context(), "state", state)
		http.Redirect(w, r, oauthConfig.AuthCodeURL(state), http.StatusFound)
	})

	r.Get("/api/v1/auth/github/manage", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://github.com/settings/connections/applications/"+cli.Auth.Github.ClientID, http.StatusFound)
	})

	r.Get("/api/v1/auth/github/callback", func(w http.ResponseWriter, r *http.Request) {
		if !cli.Debug {
			// Only check CSRF tokens if we're out of debug mode.
			state := session.GetString(r.Context(), "state")
			if state == "" {
				w.WriteHeader(http.StatusBadRequest)
				pt.JSON(w, r, pt.M{"error": "Session token not found, possible CSRF (or cookies disabled)? Please try again."})
				return
			}

			session.Remove(r.Context(), "state")

			if state != r.FormValue("state") {
				w.WriteHeader(http.StatusBadRequest)
				pt.JSON(w, r, pt.M{"error": "Session token not found, possible CSRF (or cookies disabled)? Please try again."})
				return
			}
		}
		session.Remove(r.Context(), "state")

		token, err := oauthConfig.Exchange(r.Context(), r.FormValue("code"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			pt.JSON(w, r, pt.M{"error": "Error getting token: " + err.Error()})
			return
		}

		client := github.NewClient(oauthConfig.Client(r.Context(), token))

		guser, _, err := client.Users.Get(r.Context(), "")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			pt.JSON(w, r, pt.M{"error": "Error obtaining user information: " + err.Error()})
			return
		}

		user := &models.User{
			GithubID:       int(guser.GetID()),
			Token:          token.AccessToken,
			AvatarURL:      guser.GetAvatarURL(),
			Username:       guser.GetLogin(),
			Name:           guser.GetName(),
			Email:          guser.GetEmail(),
			AccountCreated: guser.GetCreatedAt().Time,
			AccountUpdated: guser.GetUpdatedAt().Time,
		}

		if err := svcUsers.Upsert(r.Context(), user); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			pt.JSON(w, r, pt.M{"error": "Error writing user to database"})
			logger.WithError(err).Error("error writing user to database")
			return
		}

		// TODO: prevent session fixation: https://github.com/alexedwards/scs#preventing-session-fixation
		session.Put(r.Context(), "user_id", user.ID.Hex())

		w.WriteHeader(http.StatusOK)
		pt.JSON(w, r, pt.M{"authenticated": true, "user": user})
	})

	r.Get("/api/v1/auth/self", func(w http.ResponseWriter, r *http.Request) {
		id := session.GetString(r.Context(), "user_id")
		if id == "" {
			w.WriteHeader(http.StatusUnauthorized)
			pt.JSON(w, r, pt.M{"authenticated": false, "error": http.StatusText(http.StatusUnauthorized)})
			return
		}

		user, err := svcUsers.Get(r.Context(), id)
		if err != nil {
			if models.IsNotFound(err) {
				session.Remove(r.Context(), "user_id")
				w.WriteHeader(http.StatusUnauthorized)
				pt.JSON(w, r, pt.M{"authenticated": false, "error": http.StatusText(http.StatusUnauthorized)})
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			pt.JSON(w, r, pt.M{"authenticated": false, "error": "An internal server error occurred."})
			return
		}

		pt.JSON(w, r, pt.M{"authenticated": true, "user": user})
	})

	r.Get("/api/v1/auth/logout", func(w http.ResponseWriter, r *http.Request) {
		session.Remove(r.Context(), "user_id")

		w.WriteHeader(http.StatusOK)
		pt.JSON(w, r, pt.M{"authenticated": false})
	})
}
