// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/lrstanley/chix/v2"
)

const (
	// staticFallback is served for routes which don't match a file, so the SPA
	// router can handle them.
	staticFallback = "index.html"

	// staticAssetsDir holds the build output with content-hashed file names, which
	// never change, so they can be cached forever.
	staticAssetsDir = "_nuxt/"

	cacheImmutable       = "public, max-age=31536000, immutable"
	cacheRevalidate      = "no-cache"
	staticAllowedMethods = "GET, HEAD"
)

// staticHandler serves the embedded SPA build (see [newStaticHandler]).
type staticHandler struct {
	fsys fs.FS
	// etags maps file paths to their ETag. Embedded files have no modification
	// time, so the content hash is used for conditional requests instead.
	etags map[string]string
}

// newStaticHandler returns a catch-all handler (for use with chi's NotFound)
// which serves the SPA build in fsys:
//
//   - Requests under the API base path get a JSON 404.
//   - Only GET and HEAD are allowed.
//   - Existing files are served as-is, with an ETag. Hashed assets (see
//     [staticAssetsDir]) are cached forever, and everything else, including
//     HTML, must be revalidated.
//   - Missing files that look like non-HTML assets (by extension), or are under
//     [staticAssetsDir], get a 404. Anything else gets [staticFallback], so the
//     SPA router handles it.
func newStaticHandler(fsys fs.FS) (*staticHandler, error) {
	h := &staticHandler{fsys: fsys, etags: map[string]string{}}

	err := fs.WalkDir(fsys, ".", func(name string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		f, err := fsys.Open(name)
		if err != nil {
			return err
		}
		defer f.Close()

		sum := sha256.New()
		if _, err = io.Copy(sum, f); err != nil {
			return fmt.Errorf("hashing %q: %w", name, err)
		}
		// Weak, as the response may be compressed.
		h.etags[name] = `W/"` + hex.EncodeToString(sum.Sum(nil)[:16]) + `"`
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("indexing static files: %w", err)
	}

	if _, ok := h.etags[staticFallback]; !ok {
		return nil, fmt.Errorf("static files are missing %q (build the frontend first)", staticFallback)
	}
	return h, nil
}

func (h *staticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if base := chix.GetConfig(r.Context()).GetAPIBasePath(); strings.HasPrefix(r.URL.Path, base) {
		chix.ErrorWithCode(w, r, http.StatusNotFound, errors.New("resource not found"))
		return
	}

	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", staticAllowedMethods)
		chix.ErrorWithCode(w, r, http.StatusMethodNotAllowed)
		return
	}

	name := strings.TrimPrefix(path.Clean("/"+r.URL.Path), "/")
	if name == "" {
		name = staticFallback
	} else if _, ok := h.etags[name]; !ok {
		// Directories with an index (e.g. pre-rendered routes).
		if _, ok = h.etags[path.Join(name, "index.html")]; ok {
			name = path.Join(name, "index.html")
		}
	}

	if _, ok := h.etags[name]; !ok {
		ext, _, _ := strings.Cut(mime.TypeByExtension(path.Ext(name)), ";")
		if strings.HasPrefix(name, staticAssetsDir) || (ext != "" && ext != "text/html") {
			chix.ErrorWithCode(w, r, http.StatusNotFound, errors.New("resource not found"))
			return
		}
		name = staticFallback
	}

	f, err := h.fsys.Open(name)
	if err != nil {
		chix.Error(w, r, err)
		return
	}
	defer f.Close()

	content, ok := f.(io.ReadSeeker)
	if !ok {
		chix.Error(w, r, fmt.Errorf("static file %q isn't seekable", name))
		return
	}

	w.Header().Set("ETag", h.etags[name])
	if strings.HasPrefix(name, staticAssetsDir) {
		w.Header().Set("Cache-Control", cacheImmutable)
	} else {
		w.Header().Set("Cache-Control", cacheRevalidate)
	}

	// Handles HEAD, conditional (If-None-Match) and range requests, and the
	// content type (by extension).
	http.ServeContent(w, r, name, time.Time{}, content)
}
