// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/joho/godotenv/autoload"
	"github.com/lrstanley/recoverer"
	"github.com/lrstanley/spectrograph/pkg/http/helpers"
	lmiddleware "github.com/lrstanley/spectrograph/pkg/http/middleware"
	"golang.org/x/oauth2"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var (
	oauthConfig *oauth2.Config
	session     *scs.SessionManager
)

func httpServer(ctx context.Context, wg *sync.WaitGroup, errors chan<- error) {
	// Initialize OAuth.
	oauthConfig = &oauth2.Config{
		ClientID:     cli.Auth.Github.ClientID,
		ClientSecret: cli.Auth.Github.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
		RedirectURL: "",
		Scopes:      []string{}, // TODO: github scopes?
	}

	// Initialize http error handler with out logger.
	helpers.DefaultHTTPErrorHandler = helpers.NewHTTPErrorHandler(logger)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)

	if cli.Proxy {
		r.Use(middleware.RealIP)
	}

	r.Use(lmiddleware.StructuredLoggerMiddleware(logger, !cli.Debug))
	r.Use(middleware.Compress(5))
	r.Use(recoverer.New(recoverer.Options{
		Logger: &recoverer.LeveledLoggerWriter{Logger: logger},
		Show:   cli.Debug,
		Simple: false,
	}))
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(middleware.StripSlashes)

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

	registerHTTPRoutes(r)

	// Initialize sessions.
	session = scs.New()
	session.ErrorFunc = func(w http.ResponseWriter, r *http.Request, err error) {
		// http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		helpers.HTTPError(w, r, http.StatusInternalServerError, err)
		logger.WithError(err).Error("session error")
	}
	session.Store = svcSessions
	session.IdleTimeout = 7 * 24 * time.Hour
	session.Lifetime = 30 * 24 * time.Hour

	// Setup our http server.
	srv := &http.Server{
		Addr:    cli.HTTP,
		Handler: session.LoadAndSave(r),

		// TODO: lower these at some point.
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.WithField("bind", cli.HTTP).Info("initializing http server")

		var err error

		if cli.TLS.Enabled {
			err = srv.ListenAndServeTLS(cli.TLS.Cert, cli.TLS.Key)
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
