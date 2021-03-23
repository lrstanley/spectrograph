// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package middleware

import (
	"net/http"
	"time"

	"github.com/apex/log"
	"github.com/go-chi/chi/middleware"
)

// StructuredLoggerMiddleware wraps each request and writes a log entry with extra info.
func StructuredLoggerMiddleware(logger *log.Logger, private bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			logEntry := log.NewEntry(logger)

			wrappedWriter := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			start := time.Now()
			defer func() {
				finish := time.Now()

				if id := middleware.GetReqID(r.Context()); id != "" {
					logEntry = logEntry.WithField("request_id", id)
				}

				if !private {
					logEntry = logEntry.WithField("remote_ip", r.RemoteAddr)
				}

				logEntry.WithFields(log.Fields{
					"host":       r.Host,
					"proto":      r.Proto,
					"method":     r.Method,
					"path":       r.URL.Path,
					"user_agent": r.Header.Get("User-Agent"),
					"status":     wrappedWriter.Status(),
					"latency_ms": float64(finish.Sub(start).Nanoseconds()) / 1000000.0,
					"bytes_in":   r.Header.Get("Content-Length"),
					"bytes_out":  wrappedWriter.BytesWritten(),
				}).Info("handled request")
			}()

			next.ServeHTTP(wrappedWriter, r)
		}
		return http.HandlerFunc(fn)
	}
}
