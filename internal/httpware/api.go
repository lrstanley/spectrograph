// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package httpware

import (
	"errors"
	"fmt"
	"net/http"
)

func APIVersionMatch(version string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientVersion := r.Header.Get("X-Api-Version")
			if clientVersion == "" {
				HandleError(w, r, http.StatusPreconditionFailed, errors.New("api version not specified"))
				return
			} else if clientVersion != version {
				HandleError(w, r, http.StatusPreconditionFailed, fmt.Errorf("server (%q) and client (%q) version mismatch", version, clientVersion))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func APIKeyRequired(keys []string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, key := range keys {
				if r.Header.Get("X-Api-Key") == key {
					next.ServeHTTP(w, r)
					return
				}
			}

			HandleError(w, r, http.StatusUnauthorized, errors.New("invalid token provided"))
		})
	}
}
