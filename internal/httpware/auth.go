// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package httpware

import (
	"context"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/apex/log"
	"github.com/lrstanley/pt"
	"github.com/lrstanley/spectrograph/internal/models"
)

func AdminRequired(session *scs.SessionManager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			admin := session.GetBool(r.Context(), models.SessionAdminKey) // TODO: should not be in session?

			if !admin {
				w.WriteHeader(http.StatusForbidden)
				pt.JSON(w, r, pt.M{"error": "administrators only"})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func AuthRequired(session *scs.SessionManager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if authed, _ := IsAuthed(session, r); !authed {
				w.WriteHeader(http.StatusUnauthorized)
				pt.JSON(w, r, pt.M{"error": "authentication required"})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// ContextUser is used to embed the user into the request context. Use GetUser
// to pull the user information back out of the request context.
func ContextUser(session *scs.SessionManager, svcUsers models.UserService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user_id := session.GetString(r.Context(), models.SessionUserIDKey)

			// Assume they're not authenticated.
			if user_id == "" {
				next.ServeHTTP(w, r)
				return
			}

			// Try to get the user from the DB.
			user, err := svcUsers.Get(r.Context(), user_id)
			if err != nil {
				log.FromContext(r.Context()).WithError(err).WithField("user_id", user).Error("unable to obtain user information")

				if models.IsNotFound(err) {
					// Has a valid session but no matching user, wipe their
					// session data (logout).
					session.Remove(r.Context(), models.SessionUserIDKey)
					next.ServeHTTP(w, r)
					return
				}
				HandleError(w, r, http.StatusInternalServerError, err)
				return
			}

			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), models.ContextUser, user)))
		})
	}
}

func GetUser(r *http.Request) (user *models.User) {
	if r == nil {
		return
	}

	user, _ = r.Context().Value(models.ContextUser).(*models.User)
	return user
}

func IsAuthed(session *scs.SessionManager, r *http.Request) (bool, string) {
	user_id := session.GetString(r.Context(), models.SessionUserIDKey)

	return user_id != "", user_id
}