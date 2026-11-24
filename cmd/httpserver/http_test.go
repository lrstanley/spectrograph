// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/lrstanley/chix/xauth/v2"
	"github.com/lrstanley/clix/v2"
	"github.com/lrstanley/spectrograph/internal/database/dbtest"
	"github.com/lrstanley/spectrograph/internal/database/ent"
	"github.com/lrstanley/spectrograph/internal/models"
	"github.com/markbates/goth/gothic"
)

func testFlags(baseURL string) models.Flags {
	u, err := url.Parse(baseURL)
	if err != nil {
		panic(err)
	}
	return models.Flags{
		HTTP: models.ConfigHTTP{
			BaseURL:       u,
			ValidationKey: xauth.GenerateAuthKey(),
			EncryptionKey: xauth.GenerateEncryptionKey(),
		},
		Discord: models.ConfigDiscord{ClientID: "123", ClientSecret: "secret"},
	}
}

func TestConfigHTTPValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in      string
		wantErr bool
	}{
		{in: "http://localhost:8080"},
		{in: "https://example.com/"},
		{in: "example.com", wantErr: true},
		{in: "ftp://example.com", wantErr: true},
		{in: "https://", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			t.Parallel()

			flags := testFlags(tt.in)
			if err := flags.HTTP.Validate(); (err != nil) != tt.wantErr {
				t.Fatalf("Validate() with base url %q, error = %v, wantErr %v", tt.in, err, tt.wantErr)
			}
		})
	}
}

func TestNewSessionStore(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		baseURL    string
		wantSecure bool
	}{
		{name: "http", baseURL: "http://localhost:8080", wantSecure: false},
		{name: "https", baseURL: "https://example.com", wantSecure: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			store, err := newSessionStore(testFlags(tt.baseURL).HTTP)
			if err != nil {
				t.Fatal(err)
			}

			req := httptest.NewRequest(http.MethodGet, tt.baseURL, http.NoBody)
			rec := httptest.NewRecorder()

			sess, err := store.New(req, "_gothic_session")
			if err != nil {
				t.Fatal(err)
			}
			sess.Values["k"] = "v"
			if err = sess.Save(req, rec); err != nil {
				t.Fatal(err)
			}

			header := rec.Header().Get("Set-Cookie")
			attrs := strings.Split(header, "; ")[1:]

			if got := containsFold(attrs, "Secure"); got != tt.wantSecure {
				t.Errorf("Secure = %v, want %v: %s", got, tt.wantSecure, header)
			}
			// Browsers reject Partitioned cookies without Secure, and the session
			// cookie is always first-party.
			if containsFold(attrs, "Partitioned") {
				t.Errorf("unexpected Partitioned attribute: %s", header)
			}
			if !containsFold(attrs, "HttpOnly") || !containsFold(attrs, "SameSite=Lax") {
				t.Errorf("expected HttpOnly and SameSite=Lax: %s", header)
			}
		})
	}
}

func containsFold(attrs []string, want string) bool {
	for _, a := range attrs {
		if strings.EqualFold(a, want) {
			return true
		}
	}
	return false
}

func TestNewSessionStoreInvalidKeys(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		validationKey string
		encryptionKey string
	}{
		{name: "not-hex", validationKey: "zz", encryptionKey: xauth.GenerateEncryptionKey()},
		{name: "short-validation", validationKey: "abcd", encryptionKey: xauth.GenerateEncryptionKey()},
		{name: "bad-encryption-length", validationKey: xauth.GenerateAuthKey(), encryptionKey: xauth.GenerateAuthKey()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg := testFlags("http://localhost:8080").HTTP
			cfg.ValidationKey, cfg.EncryptionKey = tt.validationKey, tt.encryptionKey
			_, err := newSessionStore(cfg)
			if err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

// TestNextURLOpenRedirect ensures the "next" URL used after logging in/out
// can't redirect to another host (browsers treat backslashes like forward
// slashes, so these must be escaped, which net/url does).
func TestNextURLOpenRedirect(t *testing.T) {
	t.Parallel()

	flags := testFlags("http://localhost:8080")
	srv, err := httpServer(
		context.Background(),
		slog.New(slog.DiscardHandler),
		nil, // Logout doesn't touch the database.
		nil,
		&clix.CLI[models.Flags]{Flags: &flags},
	)
	if err != nil {
		t.Fatal(err)
	}

	ts := httptest.NewServer(srv.Handler)
	t.Cleanup(ts.Close)

	client := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}}

	tests := []struct {
		name string
		next string
		want string
	}{
		{name: "relative", next: "/dashboard", want: "/dashboard"},
		{name: "backslash", next: `/\evil.example.com`, want: "/%5Cevil.example.com"},
		{name: "double-backslash", next: `\\evil.example.com`, want: "/-/auth/%5C%5Cevil.example.com"},
		{name: "triple-slash", next: "///evil.example.com", want: "/evil.example.com"},
		{name: "protocol-relative", next: "//evil.example.com", want: "/"},
		{name: "other-host", next: "https://evil.example.com/", want: "/"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req, rerr := http.NewRequestWithContext(
				context.Background(),
				http.MethodGet,
				ts.URL+"/-/auth/logout?next="+url.QueryEscape(tt.next),
				http.NoBody,
			)
			if rerr != nil {
				t.Fatal(rerr)
			}

			resp, rerr := client.Do(req)
			if rerr != nil {
				t.Fatal(rerr)
			}
			resp.Body.Close()

			if loc := resp.Header.Get("Location"); loc != tt.want {
				t.Errorf("location = %q, want %q", loc, tt.want)
			}
		})
	}
}

// TestSessionCookieCleared ensures logging out and deleting the account expire
// the session cookie with the same attributes it was set with at login, as
// browsers otherwise keep it. Not parallel: gothic's session store is global,
// and set by the first server created, so this runs before the parallel tests
// to use the https store.
func TestSessionCookieCleared(t *testing.T) { //nolint:paralleltest
	db := dbtest.New(t)
	f := dbtest.NewFixtures(t, db)

	flags := testFlags("https://spectrograph.example.com")
	srv, err := httpServer(context.Background(), slog.New(slog.DiscardHandler), db, nil, &clix.CLI[models.Flags]{Flags: &flags})
	if err != nil {
		t.Fatal(err)
	}
	ts := httptest.NewServer(srv.Handler)
	t.Cleanup(ts.Close)

	client := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}}

	// loginCookie returns the session cookie set after a successful login (xauth
	// stores the user ID with gothic.StoreInSession).
	loginCookie := func(t *testing.T, u *ent.User) *http.Cookie {
		t.Helper()
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, ts.URL, http.NoBody)
		if serr := gothic.StoreInSession("_auth", strconv.Itoa(u.ID), req, rec); serr != nil {
			t.Fatal(serr)
		}
		cookies := rec.Result().Cookies()
		if len(cookies) != 1 || cookies[0].Name != gothic.SessionName {
			t.Fatalf("unexpected login cookies: %v", cookies)
		}
		return cookies[0]
	}

	tests := []struct {
		name       string
		method     string
		path       string
		as         *ent.User
		wantStatus int
	}{
		{name: "logout", method: http.MethodGet, path: "/-/auth/logout?next=/", as: f.UserA, wantStatus: http.StatusFound},
		{name: "delete-self", method: http.MethodDelete, path: "/-/self", as: f.UserB, wantStatus: http.StatusNoContent},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			login := loginCookie(t, tt.as)
			if !login.Secure || !login.HttpOnly || login.SameSite != http.SameSiteLaxMode || login.Path != "/" ||
				login.Domain != "" || login.Partitioned || login.MaxAge <= 0 {
				t.Errorf("unexpected login cookie attributes: %s", login)
			}

			req, rerr := http.NewRequestWithContext(context.Background(), tt.method, ts.URL+tt.path, http.NoBody)
			if rerr != nil {
				t.Fatal(rerr)
			}
			req.AddCookie(&http.Cookie{Name: login.Name, Value: login.Value})

			resp, rerr := client.Do(req)
			if rerr != nil {
				t.Fatal(rerr)
			}
			resp.Body.Close()
			if resp.StatusCode != tt.wantStatus {
				t.Fatalf("status = %d, want %d", resp.StatusCode, tt.wantStatus)
			}

			var cleared *http.Cookie
			for _, c := range resp.Cookies() {
				if c.Name == login.Name {
					cleared = c
				}
			}
			if cleared == nil {
				t.Fatalf("session cookie wasn't cleared: %v", resp.Header.Values("Set-Cookie"))
			}

			if cleared.MaxAge >= 0 {
				t.Errorf("cookie isn't expired: %s", cleared)
			}
			// Browsers only replace (and so expire) a cookie with the same name,
			// path, and domain. Secure, SameSite and Partitioned must also match,
			// or browsers may reject the update.
			if cleared.Path != login.Path || cleared.Domain != login.Domain || cleared.Secure != login.Secure ||
				cleared.HttpOnly != login.HttpOnly || cleared.SameSite != login.SameSite ||
				cleared.Partitioned != login.Partitioned {
				t.Errorf("cleared cookie attributes differ:\n  login:   %s\n  cleared: %s", login, cleared)
			}
		})
	}
}

// newTestServer returns a test server for the full HTTP handler, without a
// database (for requests which don't need one).
func newTestServer(t *testing.T) (*httptest.Server, *http.Client) {
	t.Helper()

	flags := testFlags("http://localhost:8080")
	srv, err := httpServer(context.Background(), slog.New(slog.DiscardHandler), nil, nil, &clix.CLI[models.Flags]{Flags: &flags})
	if err != nil {
		t.Fatal(err)
	}
	ts := httptest.NewServer(srv.Handler)
	t.Cleanup(ts.Close)

	return ts, &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}}
}

// TestCrossOriginProtection ensures requests with unsafe methods from other
// origins (which would carry the session cookie) are rejected before reaching
// the API.
func TestCrossOriginProtection(t *testing.T) {
	t.Parallel()

	ts, client := newTestServer(t)

	tests := []struct {
		site       string
		wantStatus int
	}{
		{site: "cross-site", wantStatus: http.StatusForbidden},
		{site: "same-site", wantStatus: http.StatusForbidden},
		// Not authenticated, but past the cross-origin check.
		{site: "same-origin", wantStatus: http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.site, func(t *testing.T) {
			t.Parallel()

			req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, ts.URL+"/-/users/1/unban", http.NoBody)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Sec-Fetch-Site", tt.site)

			resp, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("status = %d, want %d", resp.StatusCode, tt.wantStatus)
			}
		})
	}
}

// TestAuthIgnoresProvidedState ensures logins always use a random OAuth state,
// rather than one provided by the caller (see stripAuthState).
func TestAuthIgnoresProvidedState(t *testing.T) {
	t.Parallel()

	ts, client := newTestServer(t)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, ts.URL+"/-/auth/providers/discord?state=attacker", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	loc, err := url.Parse(resp.Header.Get("Location"))
	if err != nil || loc.Host != "discord.com" {
		t.Fatalf("location = %q, want a redirect to discord.com", resp.Header.Get("Location"))
	}
	if state := loc.Query().Get("state"); state == "" || state == "attacker" {
		t.Errorf("state = %q, want a random state", state)
	}
}
