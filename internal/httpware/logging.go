// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package httpware

import (
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/apex/log"
	"github.com/go-chi/chi/v5/middleware"
)

// StructuredLogger wraps each request and writes a log entry with
// extra info. StructuredLogger also injects a logger into the request
// context that can be used by children middleware business logic.
func StructuredLogger(logger log.Interface, session *scs.SessionManager, private bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			logEntry := logger.WithField("src", "httphandler")

			// RequestID middleware must be loaded before this is loaded into
			// the chain.
			if id := middleware.GetReqID(r.Context()); id != "" {
				logEntry = logEntry.WithField("request_id", id)
			}

			wrappedWriter := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			start := time.Now()
			defer func() {
				finish := time.Now()

				if !private {
					logEntry = logEntry.WithField("remote_ip", r.RemoteAddr)
				}

				if workerShardID := r.Header.Get("X-Shard-Id"); workerShardID != "" {
					logEntry = logEntry.WithField("shard_id", workerShardID)
				}

				if workerVersion := r.Header.Get("X-Api-Version"); workerVersion != "" {
					logEntry = logEntry.WithField("worker_version", workerVersion)
				}

				authed, userId := IsAuthed(session, r)
				if authed {
					logEntry = logEntry.WithField("user_id", userId)
				}

				logEntry.WithFields(log.Fields{
					"host":        r.Host,
					"proto":       r.Proto,
					"method":      r.Method,
					"path":        r.URL.Path,
					"user_agent":  r.Header.Get("User-Agent"),
					"status":      wrappedWriter.Status(),
					"duration_ms": float64(finish.Sub(start).Nanoseconds()) / 1000000.0,
					"bytes_in":    r.Header.Get("Content-Length"),
					"bytes_out":   wrappedWriter.BytesWritten(),

					"authed": authed,
				}).Info("handled request")
			}()

			next.ServeHTTP(wrappedWriter, r.WithContext(log.NewContext(r.Context(), logEntry)))
		}
		return http.HandlerFunc(fn)
	}
}
