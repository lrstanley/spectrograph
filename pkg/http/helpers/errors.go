// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package helpers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/apex/log"
	"github.com/go-chi/chi/middleware"
	"github.com/lrstanley/pt"
	"github.com/lrstanley/spectrograph/pkg/models"
)

// DefaultHTTPErrorHandler is the default handler (that can be overridden) for
// http errors.
var DefaultHTTPErrorHandler = NewHTTPErrorHandler(nil)

// HTTPError is an http handler that uses the DefaultHTTPErrorHandler's Handle
// method.
func HTTPError(w http.ResponseWriter, r *http.Request, statusCode int, err error) bool {
	return DefaultHTTPErrorHandler.Handle(w, r, statusCode, err)
}

// NewHTTPErrorHandler returns a new HTTPErrorHandler.
func NewHTTPErrorHandler(logger *log.Logger) *HTTPErrorHandler {
	return &HTTPErrorHandler{logger: logger}
}

// HTTPErrorHandler is an http handler that adjusts errors, adds additional
// information, and adjusts formatting depending on the input request.
type HTTPErrorHandler struct {
	logger *log.Logger
}

// Handle handles the error (if any). Handler WILL respond to the request
// with a header and a response, if ther is an error. The return boolean tells
// the caller if the handler has responded to the request or not.
func (h *HTTPErrorHandler) Handle(w http.ResponseWriter, r *http.Request, statusCode int, err error) bool {
	if statusCode == http.StatusNotFound && err == nil {
		err = errors.New("The requested resource was not found")
	}

	if err == nil {
		return false
	}

	if statusCode == http.StatusInternalServerError && models.IsClientError(err) {
		if models.IsNotFound(err) {
			statusCode = http.StatusNotFound
		} else {
			// if it's internal server error, override since we know it's a client error.
			statusCode = http.StatusBadRequest
		}
	}

	w.WriteHeader(statusCode)

	statusText := http.StatusText(statusCode)

	if strings.HasPrefix(r.URL.Path, "/api/") {
		pt.JSON(w, r, pt.M{"error": err.Error(), "type": statusText, "code": statusCode})
	} else {
		http.Error(w, statusText+": "+err.Error(), statusCode)
	}

	if h.logger != nil && statusCode >= 500 {
		h.logger.WithFields(log.Fields{
			"status":      statusCode,
			"status_text": statusText,
			"request_id":  middleware.GetReqID(r.Context()),
		}).Errorf("http error: %v", err)
	}

	return true
}
