// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"errors"
	"net/http"
	"strings"

	rice "github.com/GeertJohan/go.rice"
	"github.com/go-chi/chi/v5"
	"github.com/lrstanley/spectrograph/pkg/httpware"
)

func registerHTTPRoutes(r chi.Router) {
	// Because it's Vue, serve the index.html when possible.
	r.Get("/", serveIndex)
	r.NotFound(serveIndex)
	registerAuthRoutes(r)
	registerAdminRoutes(r)

	// TODO: similar?
	//  - https://github.com/eamonnmcevoy/go_rest_api/blob/master/pkg/server/server.go#L17-L25
	//  - https://github.com/eamonnmcevoy/go_rest_api/blob/master/pkg/server/user_router.go
	//
	// TODO: ^ the above links for some reason have the http prefix stripped?
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api/") {
		httpware.HandleError(w, r, http.StatusNotFound, nil)
		return
	}
	if r.Method != http.MethodGet {
		httpware.HandleError(w, r, http.StatusMethodNotAllowed, errors.New(http.StatusText(http.StatusMethodNotAllowed)))
		return
	}
	w.Write(rice.MustFindBox("public/dist").MustBytes("index.html"))
}
