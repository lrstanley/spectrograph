// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package handlers

import (
	"context"
	"errors"
	"io"
	"maps"
	"net/http"
	"net/http/httptest"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/lrstanley/chix/v2"
	"github.com/lrstanley/spectrograph/internal/auth"
	"github.com/lrstanley/spectrograph/internal/database/dbtest"
	"github.com/lrstanley/spectrograph/internal/database/ent"
	"github.com/lrstanley/spectrograph/internal/models"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/discord"
)

// sessionCookie returns a session cookie which authenticates as u, like the one
// set after a successful OAuth2 flow.
func (e *testEnv) sessionCookie(u *ent.User) *http.Cookie {
	e.t.Helper()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, e.srv.URL, http.NoBody)

	// "_auth" is the session key used by xauth.
	if err := gothic.StoreInSession("_auth", strconv.Itoa(u.ID), req, rec); err != nil {
		e.t.Fatal(err)
	}

	cookies := rec.Result().Cookies()
	if len(cookies) != 1 {
		e.t.Fatalf("expected 1 session cookie, got %d", len(cookies))
	}
	return cookies[0]
}

// get sends a GET request with the provided cookies, without following
// redirects.
func (e *testEnv) get(ctx context.Context, path string, cookies ...*http.Cookie) *http.Response {
	e.t.Helper()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, e.srv.URL+path, http.NoBody)
	if err != nil {
		e.t.Fatal(err)
	}
	for _, c := range cookies {
		req.AddCookie(c)
	}

	client := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}}

	resp, err := client.Do(req)
	if err != nil {
		e.t.Fatal(err)
	}
	return resp
}

func banUser(t *testing.T, db *ent.Client, u *ent.User) {
	t.Helper()
	db.User.UpdateOneID(u.ID).
		SetBanned(true).
		SetBanReason("secret reason").
		ExecX(dbtest.AllowAll(context.Background()))
}

func TestBannedSession(t *testing.T) {
	t.Parallel()

	e := newTestEnv(t)
	ctx := context.Background()
	banUser(t, e.db, e.f.UserB)

	tests := []struct {
		name       string
		path       string
		as         *ent.User
		wantStatus int
	}{
		{name: "self", path: "/-/self", as: e.f.UserA, wantStatus: http.StatusOK},
		{name: "self-banned", path: "/-/self", as: e.f.UserB, wantStatus: http.StatusForbidden},
		{name: "guilds-banned", path: "/-/guilds", as: e.f.UserB, wantStatus: http.StatusForbidden},
		{
			name: "stream-banned", path: "/-/guilds/" + strconv.Itoa(e.f.GuildB.ID) + "/guild-events/stream",
			as: e.f.UserB, wantStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resp := e.get(ctx, tt.path, e.sessionCookie(tt.as))
			if resp.StatusCode != tt.wantStatus {
				b, _ := io.ReadAll(resp.Body)
				resp.Body.Close()
				t.Fatalf("status = %d, want %d: %s", resp.StatusCode, tt.wantStatus, b)
			}

			if tt.wantStatus != http.StatusForbidden {
				resp.Body.Close()
				return
			}

			body := decode[chix.DefaultErrorBody](t, resp)
			if body.Error != models.ErrUserBanned.Error() {
				t.Errorf("error = %q, want %q", body.Error, models.ErrUserBanned.Error())
			}
		})
	}
}

// TestAuthCallbackRedirects ensures failed OAuth2 callbacks redirect to the
// login page (rather than showing a raw JSON error), with the reason, and
// without leaking error details. Not parallel, as goth providers and
// gothic.CompleteUserAuth are global.
func TestAuthCallbackRedirects(t *testing.T) { //nolint:paralleltest
	var (
		bannedSnowflake string
		completeErr     error // Returned by gothic.CompleteUserAuth.
		fetchErr        error // Returned when fetching the user from Discord.
	)

	e := newTestEnv(t, func(db *ent.Client, cfg *Config) {
		cfg.AuthService = auth.NewService(db, nil, func(context.Context, string) (*models.DiscordUser, []*models.UserGuildResponse, error) {
			if fetchErr != nil {
				return nil, nil, fetchErr
			}
			return &models.DiscordUser{ID: bannedSnowflake, Username: "banned"}, nil, nil
		})
	})
	bannedSnowflake = e.f.UserB.UserID
	banUser(t, e.db, e.f.UserB)

	goth.UseProviders(discord.New("client-id", "client-secret", e.srv.URL+"/-/auth/providers/discord/callback"))
	origComplete := gothic.CompleteUserAuth
	gothic.CompleteUserAuth = func(http.ResponseWriter, *http.Request) (goth.User, error) {
		if completeErr != nil {
			return goth.User{}, completeErr
		}
		return goth.User{Provider: "discord", AccessToken: "token"}, nil
	}
	t.Cleanup(func() {
		gothic.CompleteUserAuth = origComplete
		goth.ClearProviders()
	})

	host := strings.TrimPrefix(e.srv.URL, "http://")
	hostname, _, _ := strings.Cut(host, ":")

	tests := []struct {
		name        string
		next        string // Value of the "_next" cookie, set by chix.UseNextURL.
		completeErr error
		fetchErr    error
		wantPath    string // Location, without the query.
		wantReason  string
	}{
		{name: "banned-no-next", wantPath: "/login", wantReason: "banned"},
		{
			name:       "banned-next-same-host",
			next:       "http://" + hostname + ":8081/dashboard?foo=bar",
			wantPath:   "http://" + hostname + ":8081/login",
			wantReason: "banned",
		},
		{name: "banned-next-other-host", next: "https://evil.example.com/dashboard", wantPath: "/login", wantReason: "banned"},
		{name: "banned-next-relative", next: "/dashboard", wantPath: "/login", wantReason: "banned"},
		{
			name:        "cancelled",
			completeErr: errors.New("oauth2: \"access_denied\" \"The resource owner or authorization server denied the request\""),
			wantPath:    "/login",
			wantReason:  "error",
		},
		{
			name:        "cancelled-next-same-host",
			next:        "http://" + hostname + ":8081/dashboard",
			completeErr: errors.New("state token mismatch"),
			wantPath:    "http://" + hostname + ":8081/login",
			wantReason:  "error",
		},
		{
			name:       "discord-api-error",
			fetchErr:   errors.New("discord api: 502 secret internal details"),
			wantPath:   "/login",
			wantReason: "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			completeErr, fetchErr = tt.completeErr, tt.fetchErr

			var cookies []*http.Cookie
			if tt.next != "" {
				cookies = append(cookies, &http.Cookie{Name: "_next", Value: tt.next})
			}

			resp := e.get(context.Background(), "/-/auth/providers/discord/callback?state=x&code=y", cookies...)
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusFound {
				b, _ := io.ReadAll(resp.Body)
				t.Fatalf("status = %d, want %d: %s", resp.StatusCode, http.StatusFound, b)
			}

			loc, err := url.Parse(resp.Header.Get("Location"))
			if err != nil {
				t.Fatal(err)
			}
			query := loc.Query()
			loc.RawQuery = ""
			if loc.String() != tt.wantPath {
				t.Errorf("location = %q, want %q", loc, tt.wantPath)
			}
			if got := query.Get("reason"); got != tt.wantReason {
				t.Errorf("reason = %q, want %q", got, tt.wantReason)
			}

			wantKeys := []string{"reason"}
			if tt.wantReason == "error" {
				wantKeys = append(wantKeys, "request_id")
				if query.Get("request_id") == "" {
					t.Error("expected request_id")
				}
			}
			if got := slices.Sorted(maps.Keys(query)); !slices.Equal(got, wantKeys) {
				t.Errorf("query keys = %v, want %v (no error details)", got, wantKeys)
			}

			if tt.next != "" && !slices.ContainsFunc(resp.Cookies(), func(c *http.Cookie) bool {
				return c.Name == "_next" && c.MaxAge < 0
			}) {
				t.Error("expected the _next cookie to be cleared")
			}
		})
	}
}
