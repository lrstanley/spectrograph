// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/lrstanley/chix"
	"github.com/lrstanley/spectrograph/internal/database"
	"github.com/lrstanley/spectrograph/internal/database/graphql"
	"github.com/lrstanley/spectrograph/internal/ent"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/discord"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//go:generate sh -c "mkdir -vp public/dist;touch public/dist/index.html"
//go:embed all:public/dist
var staticFS embed.FS

const (
	discordBotAuthEndpoint = "https://discord.com/oauth2/authorize?client_id=%s&scope=bot&permissions=1049616"
)

func httpServer(ctx context.Context) *http.Server {
	chix.DefaultAPIPrefix = "/-/"
	r := chi.NewRouter()

	goth.UseProviders(
		discord.New(
			cli.Flags.Discord.ClientID,
			cli.Flags.Discord.ClientSecret,
			cli.Flags.HTTP.BaseURL+"/-/auth/providers/discord/callback",
			"identify", "guilds", "email",
		),
	)

	auth := chix.NewAuthHandler[ent.User, int](
		database.NewAuthService(db, cli.Flags.Discord.Admins),
		cli.Flags.HTTP.ValidationKey,
		cli.Flags.HTTP.EncryptionKey,
	)

	if len(cli.Flags.HTTP.TrustedProxies) > 0 {
		r.Use(chix.UseRealIPCLIOpts(cli.Flags.HTTP.TrustedProxies))
	}

	// Core middeware.
	r.Use(
		chix.UseDebug(cli.Debug),
		chix.UseContextIP,
		middleware.RequestID,
		chix.UseStructuredLogger(logger),
		chix.UsePrometheus,
		chix.Recoverer,
		middleware.Maybe(middleware.StripSlashes, func(r *http.Request) bool {
			return !strings.HasPrefix(r.URL.Path, "/debug/")
		}),
		middleware.Compress(5),
		chix.UseNextURL,
	)

	// Security related.
	if !cli.Debug {
		r.Use(middleware.SetHeader("Strict-Transport-Security", "max-age=31536000"))
	}
	r.Use(
		cors.AllowAll().Handler,
		chix.UseHeaders(map[string]string{
			"Content-Security-Policy": "default-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data: *; connect-src 'self' *; media-src 'none'; object-src 'none'; child-src 'none'; frame-src 'none'; worker-src 'self'",
			"X-Frame-Options":         "DENY",
			"X-Content-Type-Options":  "nosniff",
			"Referrer-Policy":         "strict-origin",
			"Permissions-Policy":      "clipboard-write=(self)",
		}),
		auth.AddToContext,
		httprate.LimitByIP(200, 1*time.Minute),
	)

	// Misc.
	r.Use(chix.UseSecurityTxt(&chix.SecurityConfig{
		ExpiresIn: 182 * 24 * time.Hour,
		Contacts: []string{
			"mailto:me@liamstanley.io",
			"https://liam.sh/chat",
			"https://github.com/lrstanley",
		},
		KeyLinks:  []string{"https://github.com/lrstanley.gpg"},
		Languages: []string{"en"},
	}))

	r.Mount("/-/graphql", graphql.New(db, cli))
	r.With(middleware.SetHeader(
		"Content-Security-Policy",
		"default-src * 'unsafe-inline' 'unsafe-eval' data: blob:; ",
	)).Mount("/-/playground", playground.Handler("GraphQL playground", "/-/graphql"))
	r.Mount("/-/auth", auth)
	r.With(chix.UsePrivateIP).Mount("/metrics", promhttp.Handler())

	if cli.Debug {
		r.With(chix.UsePrivateIP).Mount("/debug", middleware.Profiler())
	}

	// Regular routes.
	r.Get("/-/invite", getAuthorizeBot)
	r.Get("/-/invite/{id:[a-zA-Z0-9]{3,}}", getAuthorizeBot)

	r.NotFound(chix.UseStatic(ctx, &chix.Static{
		FS:         staticFS,
		CatchAll:   true,
		AllowLocal: cli.Debug,
		Path:       "public/dist",
		SPA:        true,
		Headers: map[string]string{
			"Vary":          "Accept-Encoding",
			"Cache-Control": "public, max-age=7776000",
		},
	}).ServeHTTP)

	return &http.Server{
		Addr:    cli.Flags.HTTP.BindAddr,
		Handler: r,

		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

func getAuthorizeBot(w http.ResponseWriter, r *http.Request) {
	// https://discord.com/developers/docs/topics/oauth2#bot-authorization-flow
	url := fmt.Sprintf(discordBotAuthEndpoint, cli.Flags.Discord.ClientID)
	if guildID := chi.URLParam(r, "id"); guildID != "" {
		url += "&guild_id=" + guildID + "&disable_guild_select=true"
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
