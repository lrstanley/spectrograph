// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package models

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/lrstanley/chix"
	"github.com/lrstanley/spectrograph/internal/ent"
)

func init() {
	chix.AddErrorResolver(func(err error) (status int) {
		if ent.IsConstraintError(err) || ent.IsValidationError(err) {
			return http.StatusBadRequest
		}

		if ent.IsNotFound(err) {
			return http.StatusNotFound
		}

		if IsDatabaseError(err) {
			if code := unwrapDBError(err); code != "" {
				if pgerrcode.IsDataException(code) {
					return http.StatusBadRequest
				}
			}
			return http.StatusInternalServerError
		}

		var e *ErrUserBanned
		if errors.As(err, &e) {
			return http.StatusForbidden
		}

		return 0
	})
}

func unwrapDBError(err error) string {
	var pge *pgconn.PgError
	if errors.As(err, &pge) {
		return pge.Code
	}
	return ""
}

func IsDatabaseError(err error) bool {
	return unwrapDBError(err) != ""
}

type ErrUserBanned struct {
	User *ent.User
}

func (e ErrUserBanned) Error() string {
	return fmt.Sprintf("user is currently banned -- ban reason: %s", e.User.BanReason)
}
