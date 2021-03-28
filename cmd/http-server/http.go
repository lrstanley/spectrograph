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
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/joho/godotenv/autoload"
	"github.com/lrstanley/recoverer"
	"github.com/lrstanley/spectrograph/pkg/http/helpers"
	lmiddleware "github.com/lrstanley/spectrograph/pkg/http/middleware"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func httpServer(ctx context.Context, wg *sync.WaitGroup, errors chan<- error) {
	// Initialize http error handler with out logger.
	helpers.DefaultHTTPErrorHandler = helpers.NewHTTPErrorHandler(logger)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)

	if cli.HTTP.Proxy {
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

	// Setup our http server.
	srv := &http.Server{
		Addr:    cli.HTTP.BindAddr,
		Handler: r,

		// TODO: lower these at some point.
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
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
