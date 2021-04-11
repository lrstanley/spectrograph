// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package httpware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/apex/log"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lrstanley/pt"
	"github.com/lrstanley/spectrograph/pkg/models"
)

// HandleError handles the error (if any). Handler WILL respond to the request
// with a header and a response, if ther is an error. The return boolean tells
// the caller if the handler has responded to the request or not.
func HandleError(w http.ResponseWriter, r *http.Request, statusCode int, err error) bool {
	// TODO: err being passed in should be interface{} so we can support passing
	// in things other than errors (like strings!).
	if statusCode == http.StatusNotFound && err == nil {
		err = errors.New("the requested resource was not found")
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
		pt.JSON(w, r, pt.M{
			"error":      err.Error(),
			"type":       statusText,
			"code":       statusCode,
			"request_id": middleware.GetReqID(r.Context()),
		})
	} else {
		http.Error(w, fmt.Sprintf(
			"%s: %s (id: %s)",
			statusText,
			err.Error(),
			middleware.GetReqID(r.Context()),
		), statusCode)
	}

	if statusCode >= 500 {
		log.FromContext(r.Context()).WithFields(log.Fields{
			"status":      statusCode,
			"status_text": statusText,
		}).Errorf("http error: %v", err)
	}

	return true
}
