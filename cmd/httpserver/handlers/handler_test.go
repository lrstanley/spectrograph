// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package handlers

import (
	"bufio"
	"context"
	"encoding/json/v2"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/lrstanley/chix/v2"
	"github.com/lrstanley/chix/xauth/v2"
	"github.com/lrstanley/spectrograph/internal/auth"
	"github.com/lrstanley/spectrograph/internal/database/dbtest"
	"github.com/lrstanley/spectrograph/internal/database/ent"
	"github.com/lrstanley/spectrograph/internal/models"
)

const testUserHeader = "X-Test-User"

type testEnv struct {
	t   *testing.T
	db  *ent.Client
	f   *dbtest.Fixtures
	srv *httptest.Server
}

// testEnvOption customizes the handler config used by [newTestEnv].
type testEnvOption func(db *ent.Client, cfg *Config)

func newTestEnv(t *testing.T, opts ...testEnvOption) *testEnv {
	t.Helper()

	db := dbtest.New(t)
	f := dbtest.NewFixtures(t, db)

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	cfg := Config{
		DB:              db,
		AuthService:     auth.NewService(db, nil, nil),
		SessionStorage:  xauth.NewCookieStore(xauth.GenerateAuthKey(), xauth.GenerateEncryptionKey()),
		DiscordClientID: "123456",
		BasePath:        "/-",
		Shutdown:        ctx,
	}
	for _, opt := range opts {
		opt(db, &cfg)
	}

	r := chi.NewRouter()
	r.Use(
		chix.NewConfig().
			SetAPIBasePath("/-/").
			SetLogger(slog.New(slog.DiscardHandler)).
			SetMaskPrivateErrors(true).
			AddErrorResolvers(models.ErrorResolver).
			SetErrorHandler(models.ErrorHandler).
			Use(),
		chix.UseRequestID(),
		// Inject authentication through a header, rather than session cookies.
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if id, err := strconv.Atoi(r.Header.Get(testUserHeader)); err == nil {
					u := db.User.GetX(dbtest.AllowAll(r.Context()), id)
					r = r.WithContext(dbtest.AsUser(r.Context(), u))
				}
				next.ServeHTTP(w, r)
			})
		},
		// Session cookie authentication (see [testEnv.sessionCookie]), a no-op if
		// the header above was used.
		xauth.UseAuthContext(cfg.AuthService),
	)

	r.Route("/-", New(cfg).Route)

	srv := httptest.NewServer(r)
	t.Cleanup(srv.Close)

	return &testEnv{t: t, db: db, f: f, srv: srv}
}

func (e *testEnv) do(ctx context.Context, method, path string, as *ent.User, body string) *http.Response {
	e.t.Helper()

	var rbody io.Reader = http.NoBody
	if body != "" {
		rbody = strings.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, e.srv.URL+path, rbody)
	if err != nil {
		e.t.Fatal(err)
	}
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if as != nil {
		req.Header.Set(testUserHeader, strconv.Itoa(as.ID))
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

func decode[T any](t *testing.T, resp *http.Response) T {
	t.Helper()
	defer resp.Body.Close()

	var v T
	if err := json.UnmarshalRead(resp.Body, &v); err != nil {
		t.Fatalf("decoding response: %v", err)
	}
	return v
}

// checkFieldError checks that the response is a validation error for field,
// without leaking raw ent errors.
func checkFieldError(field, wantErr string) func(t *testing.T, resp *http.Response) {
	return func(t *testing.T, resp *http.Response) {
		t.Helper()

		body := decode[models.ErrorBody](t, resp)
		if body.Field != field {
			t.Errorf("field = %q, want %q", body.Field, field)
		}
		if body.Error != wantErr {
			t.Errorf("error = %q, want %q", body.Error, wantErr)
		}
		if body.Code != http.StatusBadRequest || body.RequestID == "" || body.Timestamp == "" {
			t.Errorf("unexpected error body: %+v", body)
		}
	}
}

func TestHandlers(t *testing.T) {
	t.Parallel()

	e := newTestEnv(t)
	f := e.f
	ctx := context.Background()

	tests := []struct {
		name       string
		method     string
		path       string
		as         *ent.User
		body       string
		wantStatus int
		check      func(t *testing.T, resp *http.Response)
	}{
		{name: "healthy", method: http.MethodGet, path: "/-/healthy", wantStatus: http.StatusOK},
		{name: "openapi", method: http.MethodGet, path: "/-/openapi.json", wantStatus: http.StatusOK},
		{
			name: "invite", method: http.MethodGet, path: "/-/invite", wantStatus: http.StatusTemporaryRedirect,
			check: func(t *testing.T, resp *http.Response) {
				loc := resp.Header.Get("Location")
				if !strings.HasPrefix(loc, "https://discord.com/oauth2/authorize?") || !strings.Contains(loc, "client_id=123456") {
					t.Errorf("unexpected redirect: %s", loc)
				}
			},
		},
		{
			name: "invite-guild", method: http.MethodGet, path: "/-/invite/2001", wantStatus: http.StatusTemporaryRedirect,
			check: func(t *testing.T, resp *http.Response) {
				if loc := resp.Header.Get("Location"); !strings.Contains(loc, "guild_id=2001") {
					t.Errorf("unexpected redirect: %s", loc)
				}
			},
		},
		{name: "self-unauthenticated", method: http.MethodGet, path: "/-/self", wantStatus: http.StatusUnauthorized},
		{
			name: "unknown-endpoint", method: http.MethodGet, path: "/-/nope", as: f.UserA, wantStatus: http.StatusNotFound,
			check: func(t *testing.T, resp *http.Response) {
				if body := decode[chix.DefaultErrorBody](t, resp); body.Code != http.StatusNotFound {
					t.Errorf("unexpected error body: %+v", body)
				}
			},
		},
		{
			name: "self", method: http.MethodGet, path: "/-/self", as: f.UserA, wantStatus: http.StatusOK,
			check: func(t *testing.T, resp *http.Response) {
				if u := decode[ent.User](t, resp); u.ID != f.UserA.ID {
					t.Errorf("self = %d, want %d", u.ID, f.UserA.ID)
				}
			},
		},
		{name: "guilds-unauthenticated", method: http.MethodGet, path: "/-/guilds", wantStatus: http.StatusUnauthorized},
		{
			name: "guilds-guild-admin", method: http.MethodGet, path: "/-/guilds", as: f.UserA, wantStatus: http.StatusOK,
			check: func(t *testing.T, resp *http.Response) {
				page := decode[struct {
					Content []*ent.Guild `json:"content"`
				}](t, resp)
				if len(page.Content) != 1 || page.Content[0].ID != f.GuildA.ID {
					t.Errorf("unexpected guilds: %+v", page.Content)
				}
				if page.Content[0].Edges.GuildConfig == nil {
					t.Error("expected guild_config to be eager-loaded")
				}
			},
		},
		{
			name: "guild-other", method: http.MethodGet, path: "/-/guilds/" + strconv.Itoa(f.GuildB.ID),
			as: f.UserA, wantStatus: http.StatusNotFound,
		},
		{
			name: "users-guild-admin", method: http.MethodGet, path: "/-/users", as: f.UserA, wantStatus: http.StatusOK,
			check: func(t *testing.T, resp *http.Response) {
				page := decode[struct {
					Content []*ent.User `json:"content"`
				}](t, resp)
				if len(page.Content) != 1 || page.Content[0].ID != f.UserA.ID {
					t.Errorf("unexpected users: %+v", page.Content)
				}
			},
		},
		{
			name: "guild-events-no-debug", method: http.MethodGet,
			path: "/-/guilds/" + strconv.Itoa(f.GuildA.ID) + "/guild-events", as: f.UserA, wantStatus: http.StatusOK,
			check: func(t *testing.T, resp *http.Response) {
				page := decode[struct {
					Content []*ent.GuildEvent `json:"content"`
				}](t, resp)
				if len(page.Content) != 1 || page.Content[0].ID != f.GuildAInfoEvent.ID {
					t.Fatalf("unexpected events: %+v", page.Content)
				}
				if g := page.Content[0].Edges.Guild; g == nil || g.ID != f.GuildA.ID {
					t.Errorf("expected guild to be eager-loaded, got: %+v", g)
				}
			},
		},
		{
			name: "guild-events-guild-admin", method: http.MethodGet, path: "/-/guild-events", as: f.UserB,
			wantStatus: http.StatusOK,
			check: func(t *testing.T, resp *http.Response) {
				page := decode[struct {
					Content []*ent.GuildEvent `json:"content"`
				}](t, resp)
				if len(page.Content) != 1 || page.Content[0].ID != f.GuildBInfoEvent.ID {
					t.Fatalf("unexpected events: %+v", page.Content)
				}
				if g := page.Content[0].Edges.Guild; g == nil || g.ID != f.GuildB.ID {
					t.Errorf("expected guild to be eager-loaded, got: %+v", g)
				}
			},
		},
		{
			name: "guild-events-system-admin", method: http.MethodGet, path: "/-/guild-events", as: f.Admin,
			wantStatus: http.StatusOK,
			check: func(t *testing.T, resp *http.Response) {
				page := decode[struct {
					Content []*ent.GuildEvent `json:"content"`
				}](t, resp)
				want := map[int]int{
					f.GuildAInfoEvent.ID:  f.GuildA.ID,
					f.GuildADebugEvent.ID: f.GuildA.ID,
					f.GuildBInfoEvent.ID:  f.GuildB.ID,
				}
				if len(page.Content) != len(want) {
					t.Fatalf("unexpected events: %+v", page.Content)
				}
				for _, ev := range page.Content {
					if g := ev.Edges.Guild; g == nil || g.ID != want[ev.ID] {
						t.Errorf("event %d: expected guild %d to be eager-loaded, got: %+v", ev.ID, want[ev.ID], g)
					}
				}
			},
		},
		{
			name: "update-guild-config", method: http.MethodPatch,
			path: "/-/guild-configs/" + strconv.Itoa(f.GuildAConfig.ID), as: f.UserA,
			body: `{"regex_match": "^bar$"}`, wantStatus: http.StatusOK,
			check: func(t *testing.T, resp *http.Response) {
				if cfg := decode[ent.GuildConfig](t, resp); cfg.RegexMatch != "^bar$" {
					t.Errorf("regex_match = %q", cfg.RegexMatch)
				}
			},
		},
		{
			name: "update-guild-config-invalid-regex", method: http.MethodPatch,
			path: "/-/guild-configs/" + strconv.Itoa(f.GuildAConfig.ID), as: f.UserA,
			body: `{"regex_match": "("}`, wantStatus: http.StatusBadRequest,
			check: checkFieldError("regex_match", "regex_match: error parsing regexp: missing closing ): `(`"),
		},
		{
			name: "update-guild-config-invalid-email", method: http.MethodPatch,
			path: "/-/guild-configs/" + strconv.Itoa(f.GuildAConfig.ID), as: f.UserA,
			body: `{"contact_email": "nope"}`, wantStatus: http.StatusBadRequest,
			check: checkFieldError("contact_email", "contact_email: invalid email address"),
		},
		{
			name: "update-guild-config-out-of-range", method: http.MethodPatch,
			path: "/-/guild-configs/" + strconv.Itoa(f.GuildAConfig.ID), as: f.UserA,
			body: `{"default_max_clones": 21}`, wantStatus: http.StatusBadRequest,
			check: checkFieldError("default_max_clones", "default_max_clones: value is out of range"),
		},
		{
			name: "update-guild-admin-config-out-of-range", method: http.MethodPatch,
			path: "/-/guild-admin-configs/" + strconv.Itoa(f.GuildAAdminConfig.ID), as: f.Admin,
			body: `{"default_max_channels": -1}`, wantStatus: http.StatusBadRequest,
			check: checkFieldError("default_max_channels", "default_max_channels: value is out of range"),
		},
		{
			name: "update-guild-config-unknown-field", method: http.MethodPatch,
			path: "/-/guild-configs/" + strconv.Itoa(f.GuildAConfig.ID), as: f.UserA,
			body: `{"nope": true}`, wantStatus: http.StatusBadRequest,
			check: func(t *testing.T, resp *http.Response) {
				if body := decode[chix.DefaultErrorBody](t, resp); !strings.HasPrefix(body.Error, "bad request: ") || body.RequestID == "" {
					t.Errorf("unexpected error body: %+v", body)
				}
			},
		},
		{
			name: "method-not-allowed", method: http.MethodDelete,
			path: "/-/guild-configs/" + strconv.Itoa(f.GuildAConfig.ID), as: f.UserA,
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name: "body-too-large", method: http.MethodPatch,
			path: "/-/guild-configs/" + strconv.Itoa(f.GuildAConfig.ID), as: f.UserA,
			body:       `{"comment": "` + strings.Repeat("a", 9<<20) + `"}`,
			wantStatus: http.StatusRequestEntityTooLarge,
		},
		{
			name: "update-guild-config-other", method: http.MethodPatch,
			path: "/-/guild-configs/" + strconv.Itoa(f.GuildBConfig.ID), as: f.UserA,
			body: `{"enabled": false}`, wantStatus: http.StatusNotFound,
		},
		{
			name: "update-guild-admin-config-guild-admin", method: http.MethodPatch,
			path: "/-/guild-admin-configs/" + strconv.Itoa(f.GuildAAdminConfig.ID), as: f.UserA,
			body: `{"enabled": false}`, wantStatus: http.StatusForbidden,
		},
		{
			name: "update-guild-admin-config-system-admin", method: http.MethodPatch,
			path: "/-/guild-admin-configs/" + strconv.Itoa(f.GuildAAdminConfig.ID), as: f.Admin,
			body: `{"comment": "ok"}`, wantStatus: http.StatusOK,
		},
		{
			name: "ban-not-admin", method: http.MethodPost,
			path: "/-/users/" + strconv.Itoa(f.UserB.ID) + "/ban", as: f.UserA,
			body: `{"reason": "spam"}`, wantStatus: http.StatusForbidden,
		},
		{
			name: "ban-missing-reason", method: http.MethodPost,
			path: "/-/users/" + strconv.Itoa(f.UserB.ID) + "/ban", as: f.Admin,
			body: `{}`, wantStatus: http.StatusBadRequest,
		},
		{
			name: "ban-self", method: http.MethodPost,
			path: "/-/users/" + strconv.Itoa(f.Admin.ID) + "/ban", as: f.Admin,
			body: `{"reason": "spam"}`, wantStatus: http.StatusBadRequest,
		},
		{
			name: "stream-other-guild", method: http.MethodGet,
			path: "/-/guilds/" + strconv.Itoa(f.GuildB.ID) + "/guild-events/stream", as: f.UserA,
			wantStatus: http.StatusNotFound,
		},
		{
			name: "stream-unauthenticated", method: http.MethodGet,
			path:       "/-/guilds/" + strconv.Itoa(f.GuildA.ID) + "/guild-events/stream",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			resp := e.do(ctx, tt.method, tt.path, tt.as, tt.body)
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				b, _ := io.ReadAll(resp.Body)
				t.Fatalf("status = %d, want %d: %s", resp.StatusCode, tt.wantStatus, b)
			}
			if tt.check != nil {
				tt.check(t, resp)
			}
		})
	}
}

func TestBanUnbanUser(t *testing.T) {
	t.Parallel()

	e := newTestEnv(t)
	ctx := context.Background()
	path := "/-/users/" + strconv.Itoa(e.f.UserB.ID)

	resp := e.do(ctx, http.MethodPost, path+"/ban", e.f.Admin, `{"reason": "spam"}`)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("ban status = %d", resp.StatusCode)
	}
	if u := decode[ent.User](t, resp); !u.Banned || u.BanReason != "spam" {
		t.Errorf("unexpected user after ban: %+v", u)
	}

	resp = e.do(ctx, http.MethodPost, path+"/unban", e.f.Admin, "")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unban status = %d", resp.StatusCode)
	}
	if u := decode[ent.User](t, resp); u.Banned || u.BanReason != "" {
		t.Errorf("unexpected user after unban: %+v", u)
	}
}

// TestRejectAggregateSorts ensures edge aggregates can't be used for sorting,
// as they're computed without the edge's privacy rules (e.g. DEBUG events).
func TestRejectAggregateSorts(t *testing.T) {
	t.Parallel()

	e := newTestEnv(t)
	ctx := context.Background()

	for _, path := range []string{
		"/-/guilds?sort=guild_events.count",
		"/-/guilds?sort=admins.count",
		"/-/guilds?Sorted.sort=guild_events.count",
		"/-/users?sort=user_guilds.member_count.sum",
	} {
		resp := e.do(ctx, http.MethodGet, path, e.f.UserA, "")
		resp.Body.Close()
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("%s: status = %d, want %d", path, resp.StatusCode, http.StatusBadRequest)
		}
	}

	resp := e.do(ctx, http.MethodGet, "/-/guilds?sort=name", e.f.UserA, "")
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("sort by name: status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func TestDeleteSelf(t *testing.T) {
	t.Parallel()

	e := newTestEnv(t)
	ctx := context.Background()

	resp := e.do(ctx, http.MethodDelete, "/-/self", e.f.UserA, "")
	resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusNoContent)
	}

	if _, err := e.db.User.Get(dbtest.AllowAll(ctx), e.f.UserA.ID); !ent.IsNotFound(err) {
		t.Errorf("expected user to be deleted, got: %v", err)
	}
	if _, err := e.db.User.Get(dbtest.AllowAll(ctx), e.f.UserB.ID); err != nil {
		t.Errorf("expected other user to still exist, got: %v", err)
	}
}

// readSSEIDs reads SSE event IDs from the stream until ctx is done.
func readSSEIDs(resp *http.Response) (ids []int) {
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		if v, ok := strings.CutPrefix(scanner.Text(), "id: "); ok {
			id, _ := strconv.Atoi(v)
			ids = append(ids, id)
		}
	}
	return ids
}

func TestStreamGuildEvents(t *testing.T) {
	t.Parallel()

	e := newTestEnv(t)
	f := e.f
	path := "/-/guilds/" + strconv.Itoa(f.GuildA.ID) + "/guild-events/stream"

	tests := []struct {
		name string
		as   *ent.User
		want []int
	}{
		{name: "guild-admin", as: f.UserA, want: []int{f.GuildAInfoEvent.ID}},
		{name: "system-admin", as: f.Admin, want: []int{f.GuildAInfoEvent.ID, f.GuildADebugEvent.ID}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
			defer cancel()

			resp := e.do(ctx, http.MethodGet, path, tt.as, "")
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("status = %d", resp.StatusCode)
			}
			if ct := resp.Header.Get("Content-Type"); ct != "text/event-stream" {
				t.Errorf("content-type = %q", ct)
			}

			if got := readSSEIDs(resp); !slices.Equal(got, tt.want) {
				t.Errorf("event ids = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestStreamGuildEventsRevalidates ensures open streams stop once the user loses
// access (banned, or no longer a guild admin), rather than streaming events
// based on the user state from when the stream started.
func TestStreamGuildEventsRevalidates(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		revoke func(t *testing.T, e *testEnv)
	}{
		{
			name:   "banned",
			revoke: func(t *testing.T, e *testEnv) { t.Helper(); banUser(t, e.db, e.f.UserA) },
		},
		{
			name: "guild-access-removed",
			revoke: func(t *testing.T, e *testEnv) {
				t.Helper()
				e.db.Guild.UpdateOne(e.f.GuildA).
					RemoveAdmins(e.f.UserA).
					ExecX(dbtest.AllowAll(context.Background()))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := newTestEnv(t, func(_ *ent.Client, cfg *Config) {
				cfg.StreamPollInterval = 50 * time.Millisecond
			})

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp := e.do(ctx, http.MethodGet, "/-/guilds/"+strconv.Itoa(e.f.GuildA.ID)+"/guild-events/stream", e.f.UserA, "")
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("status = %d", resp.StatusCode)
			}

			ids := make(chan []int, 1)
			go func() { ids <- readSSEIDs(resp) }()

			// Let the initial history through, then revoke access and add an event,
			// which must not be streamed.
			time.Sleep(200 * time.Millisecond)
			tt.revoke(t, e)
			e.db.GuildEvent.Create().
				SetGuild(e.f.GuildA).
				SetType("INFO").
				SetMessage("after revoke").
				ExecX(dbtest.AllowAll(ctx))

			select {
			case got := <-ids:
				if want := []int{e.f.GuildAInfoEvent.ID}; !slices.Equal(got, want) {
					t.Errorf("event ids = %v, want %v", got, want)
				}
			case <-ctx.Done():
				t.Fatal("stream wasn't closed after access was revoked")
			}
		})
	}
}
