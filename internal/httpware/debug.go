// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package httpware

import (
	"context"
	"net/http"

	"github.com/lrstanley/spectrograph/internal/models"
)

// Debug injects if "debugging" is enabled.
func Debug(debug bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), models.ContextDebug, debug)))
		})
	}
}

// IsDebug returns true if debugging for the server is enabled.
func IsDebug(r *http.Request) bool {
	// If it's not there, return false anyway.
	debug, _ := r.Context().Value(models.ContextDebug).(bool)
	return debug
}