// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"context"
	"encoding/json/v2"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/go-chi/chi/v5"
	"github.com/lrstanley/chix/v2"
	"github.com/lrstanley/clix/v2"
	"github.com/lrstanley/spectrograph/internal/models"
)

// newStaticTestRouter mounts the static handler as a catch-all, like
// [httpServer] does.
func newStaticTestRouter(static http.Handler) http.Handler {
	r := chi.NewRouter()
	r.Use(chix.NewConfig().SetAPIBasePath("/-/").Use())
	r.Get("/-/healthy", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	r.NotFound(static.ServeHTTP)
	return r
}

func TestStaticHandler(t *testing.T) {
	t.Parallel()

	files := fstest.MapFS{
		"index.html":         {Data: []byte("<!doctype html><title>index</title>")},
		"200.html":           {Data: []byte("<!doctype html><title>200</title>")},
		"favicon.png":        {Data: []byte("png")},
		"_nuxt/Bx1y2z3.js":   {Data: []byte("console.log(1)")},
		"_nuxt/entry.Ab.css": {Data: []byte("body{}")},
		"docs/index.html":    {Data: []byte("<!doctype html><title>docs</title>")},
	}

	static, err := newStaticHandler(files)
	if err != nil {
		t.Fatal(err)
	}

	r := newStaticTestRouter(static)

	tests := []struct {
		name        string
		method      string
		path        string
		header      map[string]string
		wantStatus  int
		wantBody    string // Prefix, if set.
		wantCache   string
		wantType    string // Prefix, if set.
		wantJSONErr bool
	}{
		{name: "root", path: "/", wantStatus: http.StatusOK, wantBody: "<!doctype html><title>index", wantCache: cacheRevalidate, wantType: "text/html"},
		{name: "spa-route", path: "/dashboard/guild/1", wantStatus: http.StatusOK, wantBody: "<!doctype html><title>index", wantCache: cacheRevalidate, wantType: "text/html"},
		{name: "spa-route-html-ext", path: "/foo.html", wantStatus: http.StatusOK, wantBody: "<!doctype html><title>index", wantCache: cacheRevalidate},
		{name: "html-file", path: "/200.html", wantStatus: http.StatusOK, wantBody: "<!doctype html><title>200", wantCache: cacheRevalidate},
		{name: "index-html", path: "/index.html", wantStatus: http.StatusOK, wantBody: "<!doctype html><title>index", wantCache: cacheRevalidate},
		{name: "dir-index", path: "/docs", wantStatus: http.StatusOK, wantBody: "<!doctype html><title>docs", wantCache: cacheRevalidate},
		{name: "asset", path: "/_nuxt/Bx1y2z3.js", wantStatus: http.StatusOK, wantBody: "console.log(1)", wantCache: cacheImmutable, wantType: "text/javascript"},
		{name: "asset-css", path: "/_nuxt/entry.Ab.css", wantStatus: http.StatusOK, wantCache: cacheImmutable, wantType: "text/css"},
		{name: "other-file", path: "/favicon.png", wantStatus: http.StatusOK, wantCache: cacheRevalidate, wantType: "image/png"},
		{name: "head-root", method: http.MethodHead, path: "/", wantStatus: http.StatusOK, wantCache: cacheRevalidate, wantType: "text/html"},
		{name: "head-asset", method: http.MethodHead, path: "/_nuxt/Bx1y2z3.js", wantStatus: http.StatusOK, wantCache: cacheImmutable},
		{name: "head-spa-route", method: http.MethodHead, path: "/dashboard", wantStatus: http.StatusOK, wantCache: cacheRevalidate},
		{name: "missing-asset", path: "/_nuxt/missing.js", wantStatus: http.StatusNotFound},
		{name: "missing-asset-no-ext", path: "/_nuxt/missing", wantStatus: http.StatusNotFound},
		{name: "missing-file", path: "/missing.png", wantStatus: http.StatusNotFound},
		{name: "post", method: http.MethodPost, path: "/", wantStatus: http.StatusMethodNotAllowed},
		{name: "api-404", path: "/-/nope", wantStatus: http.StatusNotFound, wantJSONErr: true},
		{name: "api-404-post", method: http.MethodPost, path: "/-/nope", wantStatus: http.StatusNotFound, wantJSONErr: true},
		{name: "traversal", path: "/../../etc/passwd", wantStatus: http.StatusOK, wantBody: "<!doctype html><title>index", wantCache: cacheRevalidate},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			method := tt.method
			if method == "" {
				method = http.MethodGet
			}
			req := httptest.NewRequest(method, tt.path, http.NoBody)
			for k, v := range tt.header {
				req.Header.Set(k, v)
			}
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)
			resp := rec.Result()
			body, _ := io.ReadAll(resp.Body)

			if resp.StatusCode != tt.wantStatus {
				t.Fatalf("status = %d, want %d: %s", resp.StatusCode, tt.wantStatus, body)
			}
			if tt.wantBody != "" && !strings.HasPrefix(string(body), tt.wantBody) {
				t.Errorf("body = %q, want prefix %q", body, tt.wantBody)
			}
			if method == http.MethodHead && len(body) != 0 {
				t.Errorf("HEAD returned a body: %q", body)
			}
			if got := resp.Header.Get("Cache-Control"); got != tt.wantCache && tt.wantStatus == http.StatusOK {
				t.Errorf("Cache-Control = %q, want %q", got, tt.wantCache)
			}
			if tt.wantStatus == http.StatusOK && resp.Header.Get("ETag") == "" {
				t.Error("expected an ETag")
			}
			if tt.wantStatus != http.StatusOK && resp.Header.Get("Cache-Control") == cacheImmutable {
				t.Error("error responses must not be cached as immutable")
			}
			if got := resp.Header.Get("Content-Type"); tt.wantType != "" && !strings.HasPrefix(got, tt.wantType) {
				t.Errorf("Content-Type = %q, want prefix %q", got, tt.wantType)
			}
			if tt.wantStatus == http.StatusMethodNotAllowed && resp.Header.Get("Allow") != "GET, HEAD" {
				t.Errorf("Allow = %q", resp.Header.Get("Allow"))
			}
			if tt.wantJSONErr {
				var v chix.DefaultErrorBody
				if jerr := json.Unmarshal(body, &v); jerr != nil || v.Code != http.StatusNotFound {
					t.Errorf("expected a JSON 404 error, got %q (%v)", body, jerr)
				}
			}
		})
	}
}

func TestStaticHandlerConditional(t *testing.T) {
	t.Parallel()

	static, err := newStaticHandler(fstest.MapFS{
		"index.html":  {Data: []byte("<!doctype html>")},
		"_nuxt/a1.js": {Data: []byte("a")},
		"_nuxt/b2.js": {Data: []byte("b")},
	})
	if err != nil {
		t.Fatal(err)
	}
	r := newStaticTestRouter(static)

	for _, path := range []string{"/", "/dashboard", "/_nuxt/a1.js"} {
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, http.NoBody))
		etag := rec.Header().Get("ETag")
		if etag == "" {
			t.Fatalf("%s: missing ETag", path)
		}

		req := httptest.NewRequest(http.MethodGet, path, http.NoBody)
		req.Header.Set("If-None-Match", etag)
		rec = httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		if rec.Code != http.StatusNotModified {
			t.Errorf("%s: status = %d, want %d", path, rec.Code, http.StatusNotModified)
		}
	}

	a, b := httptest.NewRecorder(), httptest.NewRecorder()
	r.ServeHTTP(a, httptest.NewRequest(http.MethodGet, "/_nuxt/a1.js", http.NoBody))
	r.ServeHTTP(b, httptest.NewRequest(http.MethodGet, "/_nuxt/b2.js", http.NoBody))
	if a.Header().Get("ETag") == b.Header().Get("ETag") {
		t.Error("different files must have different ETags")
	}
}

func TestStaticHandlerMissingIndex(t *testing.T) {
	t.Parallel()

	if _, err := newStaticHandler(fstest.MapFS{"favicon.png": {Data: []byte("png")}}); err == nil {
		t.Fatal("expected an error without index.html")
	}
}

// TestServerStatic checks static file routing through the full router (e.g.
// that chi routes HEAD requests to the static handler).
func TestServerStatic(t *testing.T) {
	t.Parallel()

	flags := testFlags("http://localhost:8080")
	srv, err := httpServer(context.Background(), slog.New(slog.DiscardHandler), nil, nil, &clix.CLI[models.Flags]{Flags: &flags})
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		method     string
		path       string
		wantStatus int
		wantCache  string
		wantJSON   bool
	}{
		{method: http.MethodGet, path: "/", wantStatus: http.StatusOK, wantCache: cacheRevalidate},
		{method: http.MethodHead, path: "/", wantStatus: http.StatusOK, wantCache: cacheRevalidate},
		{method: http.MethodHead, path: "/dashboard", wantStatus: http.StatusOK, wantCache: cacheRevalidate},
		// Unknown API routes are handled by the API router (which requires a
		// login), not the static handler, and stay JSON.
		{method: http.MethodGet, path: "/-/nope", wantStatus: http.StatusUnauthorized, wantJSON: true},
		{method: http.MethodDelete, path: "/", wantStatus: http.StatusMethodNotAllowed},
	}

	for _, tt := range tests {
		t.Run(tt.method+tt.path, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			srv.Handler.ServeHTTP(rec, httptest.NewRequest(tt.method, tt.path, http.NoBody))

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d: %s", rec.Code, tt.wantStatus, rec.Body)
			}
			if tt.wantCache != "" && rec.Header().Get("Cache-Control") != tt.wantCache {
				t.Errorf("Cache-Control = %q, want %q", rec.Header().Get("Cache-Control"), tt.wantCache)
			}
			if tt.wantJSON && !strings.HasPrefix(rec.Header().Get("Content-Type"), "application/json") {
				t.Errorf("Content-Type = %q, want JSON", rec.Header().Get("Content-Type"))
			}
		})
	}
}
