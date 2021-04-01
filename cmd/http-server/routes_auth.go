// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/lrstanley/pt"
	"github.com/lrstanley/spectrograph/pkg/httpware"
	"github.com/lrstanley/spectrograph/pkg/models"
	"github.com/lrstanley/spectrograph/pkg/util"
	"golang.org/x/oauth2"
)

// {
//     "authenticated":true,
//     "user":{
//         "RawData":{
//             "avatar":"TRUNCATED",
//             "discriminator":"0001",
//             "email":"me@liamstanley.io",
//             "flags":0,
//             "id":"212002249445081090",
//             "locale":"en-US",
//             "mfa_enabled":false,
//             "premium_type":1,
//             "public_flags":0,
//             "username":"/home/liam",
//             "verified":true
//         },
//         "Provider":"discord",
//         "Email":"TRUNCATED",
//         "Name":"/home/liam",
//         "FirstName":"",
//         "LastName":"",
//         "NickName":"",
//         "Description":"",
//         "UserID":"212002249445081090",
//         "AvatarURL":"TRUNCATED",
//         "Location":"",
//         "AccessToken":"TRUNCATED",
//         "AccessTokenSecret":"",
//         "RefreshToken":"TRUNCATED",
//         "ExpiresAt":"2021-04-04T06:17:10.4241762Z",
//         "IDToken":""
//     }
// }

var (
	oauthConfig *oauth2.Config
)

const (
	discordUserEndpoint   = "https://discord.com/api/users/@me"
	discordGuildsEndpoint = "https://discord.com/api/users/@me/guilds"
)

func registerAuthRoutes(r chi.Router) {
	// Initialize OAuth.
	oauthConfig = &oauth2.Config{
		ClientID:     cli.Auth.Discord.ID,
		ClientSecret: cli.Auth.Discord.Secret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://discord.com/api/oauth2/authorize",
			TokenURL: "https://discord.com/api/oauth2/token",
		},
		RedirectURL: cli.HTTP.BaseURL.String() + "/api/v1/auth/callback",
		Scopes: []string{
			"identify", "email", "guilds",
		},
	}

	r.Get("/api/v1/auth/redirect", authRedirect)
	r.Get("/api/v1/auth/callback", authCallback)
	r.Get("/api/v1/auth/self", authSelf)
	r.Get("/api/v1/auth/logout", authLogout)
}

func authRedirect(w http.ResponseWriter, r *http.Request) {
	if id := session.GetString(r.Context(), "user_id"); id != "" {
		_, err := svcUsers.Get(r.Context(), id)
		if err != nil {
			if models.IsNotFound(err) {
				session.Remove(r.Context(), "user_id")
				goto redirect
			}

			httpware.HandleError(w, r, http.StatusServiceUnavailable, errors.New("unable to query user"))
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

redirect:
	state := util.GenRandString(15)
	session.Put(r.Context(), "state", state)
	http.Redirect(w, r, oauthConfig.AuthCodeURL(
		state,
		// Provide AuthCodeOptions.
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("prompt", "none"),
	), http.StatusFound)
}

func authCallback(w http.ResponseWriter, r *http.Request) {
	if !cli.Debug {
		// Only check CSRF tokens if we're out of debug mode.
		state := session.GetString(r.Context(), "state")
		if state == "" {
			httpware.HandleError(w, r, http.StatusBadRequest, errors.New("session token not found, possible CSRF (or cookies disabled)? Please try again"))
			return
		}

		session.Remove(r.Context(), "state")

		if state != r.FormValue("state") {
			httpware.HandleError(w, r, http.StatusBadRequest, errors.New("session token not found, possible CSRF (or cookies disabled)? Please try again"))
			return
		}
	}
	session.Remove(r.Context(), "state")

	token, err := oauthConfig.Exchange(r.Context(), r.FormValue("code"))
	if err != nil {
		httpware.HandleError(w, r, http.StatusBadRequest, fmt.Errorf("error getting token: %w", err))
		return
	}

	client := oauthConfig.Client(r.Context(), token)

	req, err := http.NewRequest("GET", discordUserEndpoint, nil)
	if err != nil {
		httpware.HandleError(w, r, http.StatusInternalServerError, err)
		return
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	resp, err := client.Do(req)
	if err != nil {
		httpware.HandleError(w, r, http.StatusInternalServerError, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		httpware.HandleError(w, r, http.StatusInternalServerError, fmt.Errorf("discord responded with %d when trying to fetch user information", resp.StatusCode))
		return
	}

	user := &models.User{
		AccountUpdated: time.Now(),
	}
	err = json.NewDecoder(resp.Body).Decode(&user.Discord)
	if err != nil {
		httpware.HandleError(w, r, http.StatusInternalServerError, errors.New("received an invalid response from Discord"))
		return
	}

	user.Discord.LastLogin = time.Now()
	user.Discord.AccessToken = token.AccessToken
	user.Discord.RefreshToken = token.RefreshToken
	user.Discord.ExpiresAt = token.Expiry

	if err := svcUsers.Upsert(r.Context(), user); err != nil {
		httpware.HandleError(w, r, http.StatusInternalServerError, err)
		return
	}

	// TODO: prevent session fixation: https://github.com/alexedwards/scs#preventing-session-fixation
	session.Put(r.Context(), "user_id", user.ID.Hex())

	w.WriteHeader(http.StatusOK)
	pt.JSON(w, r, pt.M{"authenticated": true, "user": user.Public()})
}

func authSelf(w http.ResponseWriter, r *http.Request) {
	id := session.GetString(r.Context(), "user_id")
	if id == "" {
		httpware.HandleError(w, r, http.StatusUnauthorized, errors.New("not logged in"))
		return
	}

	user, err := svcUsers.Get(r.Context(), id)
	if err != nil {
		if models.IsNotFound(err) {
			session.Remove(r.Context(), "user_id")
			httpware.HandleError(w, r, http.StatusUnauthorized, err)
			return
		}
		httpware.HandleError(w, r, http.StatusInternalServerError, err)
		return
	}

	pt.JSON(w, r, pt.M{"authenticated": true, "user": user.Public()})
}

func authLogout(w http.ResponseWriter, r *http.Request) {
	session.Remove(r.Context(), "user_id")

	// TODO: https://discord.com/api/oauth2/token/revoke ? only do if removing
	// the account?

	w.WriteHeader(http.StatusOK)
	pt.JSON(w, r, pt.M{"authenticated": false})
}

// func authRefreshToken(ctx context.Context, config *oauth2.Config, refreshToken string) (*oauth2.Token, error) {
// 	token := &oauth2.Token{RefreshToken: refreshToken}
// 	ts := config.TokenSource(ctx, token)

// 	return ts.Token()
// }
