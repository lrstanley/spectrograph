// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package models

import (
	"errors"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lrstanley/chix/v2"
	"github.com/lrstanley/spectrograph/internal/database/ent"
	"github.com/lrstanley/spectrograph/internal/database/ent/privacy"
	"github.com/lrstanley/spectrograph/internal/database/ent/rest"
)

// ErrUserBanned is returned when a banned user attempts to authenticate, or uses
// an existing session.
var ErrUserBanned = errors.New("user is banned")

// FieldError is a client-facing error about the value of a single field (or a
// database constraint, when Field is empty). API error responses include Field
// as "field" (see [ErrorHandler]), so clients can show the error next to the
// matching input.
type FieldError struct {
	// Field is the JSON name of the field (e.g. "regex_match"), if known.
	Field string
	// Reason is a short, client-safe description of the problem.
	Reason string

	cause error
}

// Error returns "<field>: <reason>", or only the reason if the field is unknown.
func (e *FieldError) Error() string {
	if e.Field == "" {
		return e.Reason
	}
	return e.Field + ": " + e.Reason
}

// Unwrap returns the original error.
func (e *FieldError) Unwrap() error { return e.cause }

// Public implements [chix.ExposableError], as the message is client-safe.
func (e *FieldError) Public() bool { return true }

// validationReasons replaces messages from ent's built-in validators with
// friendlier ones.
var validationReasons = map[string]string{
	"value out of range":                             "value is out of range",
	"value is greater than the required length":      "value is too long",
	"value is less than the required length":         "value is too short",
	"value does not match validation":                "value has an invalid format",
	"value is greater than the required rune length": "value is too long",
	"value is less than the required rune length":    "value is too short",
}

// ClientError returns a client-safe replacement for ent validation errors, and
// database constraint and data errors (which otherwise include raw ent or
// Postgres messages), or nil if err isn't one of them.
func ClientError(err error) *FieldError {
	if ve, ok := errors.AsType[*ent.ValidationError](err); ok {
		// ent wraps validator errors as: validator failed for field "Type.name": <err>.
		reason := "value is invalid"
		if inner := errors.Unwrap(errors.Unwrap(ve)); inner != nil {
			reason = inner.Error()
			if friendly, found := validationReasons[reason]; found {
				reason = friendly
			}
		} else if strings.Contains(ve.Error(), "missing required") {
			reason = "value is required"
		}
		return &FieldError{Field: ve.Name, Reason: reason, cause: err}
	}

	pge, _ := errors.AsType[*pgconn.PgError](err)

	if ent.IsConstraintError(err) {
		fe := &FieldError{Reason: "value violates a constraint", cause: err}
		if pge != nil {
			fe.Field = pge.ColumnName // Only set for some errors, like not-null violations.
			switch pge.Code {
			case pgerrcode.UniqueViolation:
				fe.Reason = "already exists"
			case pgerrcode.ForeignKeyViolation:
				fe.Reason = "references something that doesn't exist"
			case pgerrcode.NotNullViolation:
				fe.Reason = "value is required"
			}
		}
		return fe
	}

	if pge != nil && pgerrcode.IsDataException(pge.Code) {
		return &FieldError{Field: pge.ColumnName, Reason: "value is invalid", cause: err}
	}
	return nil
}

// errorStatus returns the HTTP status code that should be used for the provided
// error, or 0 if the error isn't known.
func errorStatus(err error) int {
	switch {
	case err == nil:
		return 0
	case rest.IsEndpointNotFound(err):
		return http.StatusNotFound
	case rest.IsMethodNotAllowed(err):
		return http.StatusMethodNotAllowed
	case rest.IsBadRequest(err), rest.IsInvalidID(err):
		return http.StatusBadRequest
	case errors.Is(err, privacy.Deny), errors.Is(err, ErrUserBanned):
		return http.StatusForbidden
	case ent.IsNotFound(err):
		return http.StatusNotFound
	case ent.IsValidationError(err):
		return http.StatusBadRequest
	case ent.IsNotSingular(err):
		return http.StatusConflict
	}

	pge, _ := errors.AsType[*pgconn.PgError](err)

	switch {
	case ent.IsConstraintError(err):
		if pge != nil && pge.Code == pgerrcode.UniqueViolation {
			return http.StatusConflict
		}
		return http.StatusBadRequest
	case pge == nil:
		return 0
	case pgerrcode.IsDataException(pge.Code):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

// ErrorResolver is a [chix.ErrorResolverFn] which resolves database, privacy,
// and generated REST server (request parsing, unknown endpoint) errors to
// appropriate status codes. Validation, constraint and data
// errors are replaced with a client-safe [FieldError].
func ErrorResolver(oerr *chix.ResolvedError) *chix.ResolvedError {
	if oerr == nil {
		return nil
	}

	var code int
	var cause error
	for _, err := range slices.Concat(oerr.Errs, []error{oerr.Err}) {
		if code = errorStatus(err); code != 0 {
			cause = err
			break
		}
	}
	if code == 0 {
		return oerr
	}

	visibility := chix.ErrorMasked
	if code < http.StatusInternalServerError {
		visibility = chix.ErrorPublic
	}

	if fe := ClientError(cause); fe != nil && code < http.StatusInternalServerError {
		return &chix.ResolvedError{Err: fe, StatusCode: code, Visibility: visibility}
	}

	return &chix.ResolvedError{
		Err:        oerr.Err,
		Errs:       oerr.Errs,
		StatusCode: code,
		Visibility: visibility,
	}
}

// ErrorBody is the JSON body of API error responses. It matches
// [chix.DefaultErrorBody], plus the field name for [FieldError] errors.
type ErrorBody struct {
	Error     string   `json:"error"`
	Errors    []string `json:"errors,omitempty"`
	Type      string   `json:"type"`
	Code      int      `json:"code"`
	RequestID string   `json:"request_id,omitempty"`
	Timestamp string   `json:"timestamp"`
	// Field is the name of the invalid field (for validation errors), if known.
	Field string `json:"field,omitempty"`
}

// ErrorHandler is a [chix.ErrorHandler] which adds "field" to API error responses
// for public [FieldError] errors, and otherwise uses [chix.DefaultErrorHandler].
func ErrorHandler(w http.ResponseWriter, r *http.Request, rerr *chix.ResolvedError) {
	cfg := chix.GetConfig(r.Context())
	fe, ok := errors.AsType[*FieldError](rerr.Err)

	if !ok || fe.Field == "" || !rerr.Public() || !cfg.GetMaskPrivateErrors() ||
		cfg.GetAPIBasePath() == "" || !strings.HasPrefix(r.URL.Path, cfg.GetAPIBasePath()) {
		chix.DefaultErrorHandler(w, r, rerr)
		return
	}

	chix.JSON(w, r, rerr.StatusCode, ErrorBody{
		Error:     fe.Error(),
		Type:      http.StatusText(rerr.StatusCode),
		Code:      rerr.StatusCode,
		RequestID: chix.GetRequestIDOrHeader(r.Context(), r),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Field:     fe.Field,
	})
}
