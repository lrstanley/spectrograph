// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/alexedwards/scs/v2"
	"github.com/apex/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/joho/godotenv/autoload"
	"github.com/lrstanley/recoverer"
	"github.com/lrstanley/spectrograph/cmd/http-server/handlers/authhandler"
	"github.com/lrstanley/spectrograph/pkg/httpware"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var (
	session *scs.SessionManager
)

func httpServer(ctx context.Context, wg *sync.WaitGroup, errors chan<- error) {
	// Initialize sessions.
	session = scs.New()
	session.ErrorFunc = func(w http.ResponseWriter, r *http.Request, err error) {
		httpware.HandleError(w, r, http.StatusInternalServerError, err)
		log.FromContext(r.Context()).WithError(err).Error("session error")
	}
	session.Store = svcSessions
	session.IdleTimeout = 7 * 24 * time.Hour
	session.Lifetime = 30 * 24 * time.Hour
	session.Cookie.HttpOnly = true
	session.Cookie.Path = cli.HTTP.BaseURL.Path
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteStrictMode

	if cli.HTTP.BaseURL.Scheme == "https" {
		session.Cookie.Secure = true
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(httpware.StructuredLogger(logger, !cli.Debug))
	r.Use(httpware.Debug(cli.Debug))

	if cli.HTTP.Proxy {
		r.Use(middleware.RealIP)
	}

	r.Use(middleware.Compress(5))
	r.Use(recoverer.New(recoverer.Options{
		Logger: &recoverer.LeveledLoggerWriter{Logger: logger},
		Show:   cli.Debug,
		Simple: false,
	}))
	// TODO: throttle or httprate: https://github.com/go-chi/httprate
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(middleware.StripSlashes)

	// Bind/mount routes here.
	r.Mount("/static/dist", http.StripPrefix("/static/dist", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("Cache-Control", "public, max-age=7776000")
		http.FileServer(rice.MustFindBox("public/dist").HTTPBox()).ServeHTTP(w, r)
	})))
	r.Mount("/static", http.StripPrefix("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("Cache-Control", "public, max-age=7776000")
		http.FileServer(rice.MustFindBox("public/static").HTTPBox()).ServeHTTP(w, r)
	})))

	if cli.Debug {
		r.Mount("/debug", middleware.Profiler())
	}

	// Because it's Vue, serve the index.html when possible.
	r.Get("/", serveIndex)
	r.NotFound(serveIndex)
	r.Route("/api/auth", authhandler.New(svcUsers, oauthConfig, session).Route)
	registerAdminRoutes(r)

	// Setup our http server.
	srv := &http.Server{
		Addr:    cli.HTTP.BindAddr,
		Handler: session.LoadAndSave(r),

		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.WithField("bind", cli.HTTP.BindAddr).Info("initializing http server")

		var err error

		if cli.HTTP.TLS.Enabled {
			err = srv.ListenAndServeTLS(cli.HTTP.TLS.Cert, cli.HTTP.TLS.Key)
		} else {
			err = srv.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			errors <- fmt.Errorf("http error: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()

		logger.Info("requesting http server to shutdown")
		if err := srv.Shutdown(context.Background()); err != nil && err != http.ErrServerClosed {
			errors <- fmt.Errorf("unable to shutdown http server: %v", err)
		}
	}()
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
