// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/lrstanley/chix/v2"
	"github.com/lrstanley/chix/xauth/v2"
	"github.com/lrstanley/spectrograph/internal/database"
	"github.com/lrstanley/spectrograph/internal/database/ent"
	"github.com/lrstanley/spectrograph/internal/database/ent/guild"
	"github.com/lrstanley/spectrograph/internal/models"
)

type contextKeyGuildID struct{}

// streamGuildEvents streams guild events over SSE. Access to the guild is
// validated before the stream starts (privacy rules restrict the guild query to
// guilds the user can manage), and again before every poll. Privacy rules are
// also applied to the event queries themselves (e.g. DEBUG events are only
// visible to system admins). Streams poll on an interval, and also poll as soon
// as [Config.GuildEvents] reports a new event created in this process.
func (h *handler) streamGuildEvents() http.HandlerFunc {
	sse := chix.UseSSE(&chix.SSEConfig{
		HeartbeatInterval: 25 * time.Second,
		Shutdown:          h.cfg.Shutdown,
		ProducerFn: func(ctx context.Context, _ *http.Request, events chan<- *chix.SSEEvent) error {
			gid, _ := ctx.Value(contextKeyGuildID{}).(int)

			// Invalid (non-numeric) IDs are ignored, and the stream starts fresh.
			afterID, _ := strconv.Atoi(chix.GetSSELastEventID(ctx))

			uid := xauth.IdentFromContext[ent.User](ctx).ID

			notify, unsubscribe := h.cfg.GuildEvents.Subscribe(gid)
			defer unsubscribe()

			return database.StreamGuildEvents(ctx, h.cfg.DB, database.GuildEventStreamConfig{
				GuildID:  gid,
				AfterID:  afterID,
				Interval: h.cfg.StreamPollInterval,
				Notify:   notify,
				// Re-validate access on every poll, using the current user state, so
				// banned users (or users who lost access to the guild) stop receiving
				// events.
				Authorize: func(ctx context.Context) (context.Context, error) {
					u, err := h.cfg.AuthService.Get(ctx, uid)
					if err != nil {
						return nil, err
					}
					if u.Banned {
						return nil, models.ErrUserBanned
					}
					ctx = xauth.OverrideContextAuth(ctx, u.ID, u)
					if _, err = h.cfg.DB.Guild.Query().Where(guild.ID(gid)).OnlyID(ctx); err != nil {
						return nil, err
					}
					return ctx, nil
				},
			}, func(event *ent.GuildEvent) error {
				events <- &chix.SSEEvent{ID: strconv.Itoa(event.ID), Data: event}
				return nil
			})
		},
	})

	return func(w http.ResponseWriter, r *http.Request) {
		gid, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || gid < 1 {
			chix.ErrorWithCode(w, r, http.StatusBadRequest, errors.New("invalid guild ID"))
			return
		}

		// Heartbeats commit the response immediately, so access must be validated
		// before the stream starts.
		_, err = h.cfg.DB.Guild.Query().Where(guild.ID(gid)).OnlyID(r.Context())
		if err != nil {
			chix.Error(w, r, err)
			return
		}

		sse.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), contextKeyGuildID{}, gid)))
	}
}
