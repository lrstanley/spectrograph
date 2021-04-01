// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package httpware

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/lrstanley/pt"
)

func AdminRequired(session *scs.SessionManager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			admin := session.GetBool(r.Context(), "isadmin") // TODO: should not be in session?

			if !admin {
				w.WriteHeader(http.StatusForbidden)
				pt.JSON(w, r, pt.M{"error": "Administrators only."})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func AuthRequired(session *scs.SessionManager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := session.GetBool(r.Context(), "isauth")

			if !auth {
				w.WriteHeader(http.StatusUnauthorized)
				pt.JSON(w, r, pt.M{"error": "Authentication required."})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
