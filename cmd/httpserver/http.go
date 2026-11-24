// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"context"
	"embed"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/gorilla/sessions"
	"github.com/lrstanley/chix/v2"
	"github.com/lrstanley/chix/xauth/v2"
	"github.com/lrstanley/clix/v2"
	"github.com/lrstanley/spectrograph/cmd/httpserver/handlers"
	"github.com/lrstanley/spectrograph/internal/auth"
	"github.com/lrstanley/spectrograph/internal/database"
	"github.com/lrstanley/spectrograph/internal/database/ent"
	"github.com/lrstanley/spectrograph/internal/models"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/discord"
)

//go:generate sh -c "mkdir -vp public/dist; test -f public/dist/index.html || touch public/dist/index.html"
//go:embed all:public/dist
var staticFS embed.FS

const apiBasePath = "/-"

// devOrigins are the frontend dev server's origins, trusted in debug mode.
var devOrigins = []string{"http://localhost:3000", "http://localhost:8081"}

func httpServer(
	ctx context.Context,
	logger *slog.Logger,
	db *ent.Client,
	guildEvents *database.GuildEventNotifier,
	cli *clix.CLI[models.Flags],
) (*http.Server, error) {
	sessionStore, err := newSessionStore(cli.Flags.HTTP)
	if err != nil {
		return nil, err
	}

	r := chi.NewRouter()

	r.Use(
		chix.NewConfig().
			SetAPIBasePath(apiBasePath + "/").
			SetLogger(logger).
			SetMaskPrivateErrors(true).
			AddErrorResolvers(models.ErrorResolver).
			SetErrorHandler(models.ErrorHandler).
			Use(),
	)

	goth.UseProviders(
		discord.New(
			cli.Flags.Discord.ClientID,
			cli.Flags.Discord.ClientSecret,
			cli.Flags.HTTP.BaseURL.JoinPath(apiBasePath, "auth/providers/discord/callback").String(),
			discord.ScopeIdentify,
			discord.ScopeGuilds,
			discord.ScopeEmail,
		),
	)

	authSvc := auth.NewService(db, cli.Flags.Discord.Admins, nil)

	// Rejects cross-origin requests with unsafe methods (e.g. from another site
	// using the session cookie), except from the frontend dev server.
	csrf := http.NewCrossOriginProtection()
	if cli.Debug {
		for _, origin := range devOrigins {
			if oerr := csrf.AddTrustedOrigin(origin); oerr != nil {
				return nil, fmt.Errorf("adding trusted origin: %w", oerr)
			}
		}
	}

	if len(cli.Flags.HTTP.TrustedProxies) > 0 {
		r.Use(chix.UseRealIPStringOpts(cli.Flags.HTTP.TrustedProxies))
	}

	r.Use(
		chix.UseContextIP(),
		chix.UseRequestID(),
		chix.UseStripSlashes(),
		chix.UseStructuredLogger(chix.DefaultLogConfig()),
		// Compression buffers responses, which breaks SSE streams.
		chix.UseIfFunc(middleware.Compress(5), func(r *http.Request) bool {
			return !strings.HasSuffix(r.URL.Path, "/stream")
		}),
		chix.UseNextURL(),
		chix.UseIf(!cli.Debug, chix.UseHeaders(map[string]string{
			"Strict-Transport-Security": "max-age=31536000",
		})),
		chix.UseHeaders(map[string]string{
			"X-Frame-Options":        "DENY",
			"X-Content-Type-Options": "nosniff",
			"Referrer-Policy":        "strict-origin-when-cross-origin",
			"Permissions-Policy":     "clipboard-write=(self)",
		}),
		// Allow the frontend dev server to make credentialed requests in debug mode.
		chix.UseIf(cli.Debug, chix.UseCrossOriginResourceSharing(&chix.CORSConfig{
			AllowedOrigins: slices.Clone(devOrigins), // Normalized in place by chix.
			AllowedMethods: []string{
				http.MethodGet,
				http.MethodHead,
				http.MethodPost,
				http.MethodPut,
				http.MethodPatch,
				http.MethodDelete,
				http.MethodOptions,
			},
			AllowCredentials: true,
		})),
		csrf.Handler,
		chix.UseSecurityText(chix.SecurityTextConfig{
			ExpiresIn: 182 * 24 * time.Hour,
			Contacts: []string{
				"mailto:liam@liam.sh",
				"https://liam.sh/chat",
				"https://github.com/lrstanley",
			},
			KeyLinks:  []string{"https://github.com/lrstanley.gpg"},
			Languages: []string{"en"},
		}),
	)

	api := handlers.New(handlers.Config{
		DB:              db,
		AuthService:     authSvc,
		SessionStorage:  sessionStore,
		DiscordClientID: cli.Flags.Discord.ClientID,
		BasePath:        apiBasePath,
		Shutdown:        ctx,
		GuildEvents:     guildEvents,
	})
	r.Route(apiBasePath, func(r chi.Router) {
		r.Use(
			// Keyed off the IP resolved by chix (which respects trusted proxies).
			// Only for the API, as a single SPA page load fetches dozens of assets,
			// and before loading the session, which queries the database.
			httprate.LimitBy(200, 1*time.Minute, func(r *http.Request) (string, error) {
				return httprate.CanonicalizeIP(chix.GetContextIP(r.Context()).String()), nil
			}),
			xauth.UseAuthContext(authSvc),
		)
		api.Route(r)
	})

	if cli.Debug {
		r.With(chix.UsePrivateIP()).Mount("/debug", middleware.Profiler())
	}

	publicFS, err := fs.Sub(staticFS, "public/dist")
	if err != nil {
		return nil, fmt.Errorf("opening embedded static files: %w", err)
	}
	static, err := newStaticHandler(publicFS)
	if err != nil {
		return nil, err
	}
	r.NotFound(static.ServeHTTP)

	return &http.Server{
		Addr:         cli.Flags.HTTP.BindAddr,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}, nil
}

// newSessionStore returns the encrypted cookie session store. Session cookies
// are only marked Secure when the base URL uses https, so local http
// development still works.
func newSessionStore(cfg models.ConfigHTTP) (*sessions.CookieStore, error) {
	validationKey, err := hex.DecodeString(cfg.ValidationKey)
	if err != nil || len(validationKey) < 32 {
		return nil, errors.New("validation key must be a hex-encoded key of at least 32 bytes")
	}
	encryptionKey, err := hex.DecodeString(cfg.EncryptionKey)
	if err != nil || (len(encryptionKey) != 16 && len(encryptionKey) != 24 && len(encryptionKey) != 32) {
		return nil, errors.New("encryption key must be a hex-encoded 16, 24, or 32 byte key")
	}

	store := xauth.NewCookieStore(cfg.ValidationKey, cfg.EncryptionKey)
	store.Options.Secure = cfg.BaseURL.Scheme == "https"
	// xauth sets Partitioned (CHIPS), which browsers reject on cookies without
	// Secure, and which only matters for cookies used in third-party contexts.
	// The session cookie is always first-party (the app can't be framed).
	store.Options.Partitioned = false
	return store, nil
}
