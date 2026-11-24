// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package handlers

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/lrstanley/chix/v2"
	"github.com/lrstanley/spectrograph/internal/models"
)

const (
	// nextURLCookie is the cookie [chix.UseNextURL] stores the "next" URL in.
	nextURLCookie = "_next"

	// loginPath is the frontend login page, which users are sent to (with a
	// "reason" query parameter) when the OAuth2 flow fails.
	loginPath = "/login"

	// Values of the "reason" query parameter for [loginPath].
	loginReasonBanned = "banned"
	loginReasonError  = "error"
)

// useLoginErrorRedirect redirects to the frontend login page (see [loginPath])
// when the OAuth2 callback fails, rather than rendering a raw JSON error: banned
// users with reason "banned", and any other failure (e.g. the user cancelled on
// Discord, an invalid state, or a Discord API error) with reason "error" and the
// request ID. The error details are only logged. xauth responds through the chix
// error handler, so the error handler is wrapped for these routes.
func useLoginErrorRedirect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg := chix.GetConfig(r.Context())
		handleError := cfg.GetErrorHandler()

		cfg.SetErrorHandler(func(w http.ResponseWriter, r *http.Request, rerr *chix.ResolvedError) {
			switch {
			case errors.Is(rerr, models.ErrUserBanned):
				redirectLogin(w, r, url.Values{"reason": {loginReasonBanned}})
			case strings.HasSuffix(r.URL.Path, "/callback"):
				query := url.Values{"reason": {loginReasonError}}
				if id := chix.GetRequestIDOrHeader(r.Context(), r); id != "" {
					query.Set("request_id", id)
				}
				redirectLogin(w, r, query)
			default:
				handleError(w, r, rerr)
			}
		}).Use()(next).ServeHTTP(w, r)
	})
}

// redirectLogin redirects to [loginPath] with the provided query. The frontend
// may be served from another origin than the API (e.g. the dev server on
// another port), so when the "next" URL (from when the OAuth2 flow started) is
// an absolute URL on the same host, its origin is used.
func redirectLogin(w http.ResponseWriter, r *http.Request, query url.Values) {
	target := loginPath + "?" + query.Encode()
	hostname := (&url.URL{Host: r.Host}).Hostname()

	if c, err := r.Cookie(nextURLCookie); err == nil {
		u, perr := url.Parse(c.Value)
		if perr == nil && (u.Scheme == "http" || u.Scheme == "https") && strings.EqualFold(u.Hostname(), hostname) {
			target = (&url.URL{Scheme: u.Scheme, Host: u.Host}).String() + target
		}

		// Consumed, so it isn't used for the next login.
		http.SetCookie(w, &http.Cookie{ //nolint:gosec // Mirrors the chix cookie being cleared.
			Name:     nextURLCookie,
			Path:     "/",
			Domain:   hostname,
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   r.TLS != nil,
			SameSite: http.SameSiteLaxMode,
		})
	}

	http.Redirect(w, r, target, http.StatusFound)
}

// stripAuthState removes a caller-provided "state" when starting a login.
// gothic uses it as the OAuth state when present, rather than a random one,
// which would let another site log a user in to an account of its choosing
// (login CSRF) by pairing a known state with its own authorization code.
func stripAuthState(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if q := r.URL.Query(); q.Has("state") && !strings.HasSuffix(r.URL.Path, "/callback") {
			q.Del("state")
			r.URL.RawQuery = q.Encode()
		}
		next.ServeHTTP(w, r)
	})
}
