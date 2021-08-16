// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package database

import (
	"errors"

	"github.com/lrstanley/spectrograph/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	// ErrInvalidObjectID is returned when the provided ID doesn't match an
	// Object ID.
	ErrInvalidObjectID = &models.ErrClientError{Err: errors.New("invalid ID provided")}
)

// errorWrapper is a useful wrapper for returning customized errors from mongo.
func errorWrapper(err error) error {
	// Compare values.
	switch err {
	case mongo.ErrNoDocuments:
		return &models.ErrClientError{Err: &models.ErrNotFound{Err: err}}
	default:
	}

	// Compare types.
	switch v := err.(type) {
	case *models.ErrClientError:
		return err
	case *models.ErrDuplicate, *models.ErrNotFound, *models.ErrValidationFailed:
		return &models.ErrClientError{Err: v}
	case mongo.WriteException:
		for _, werr := range v.WriteErrors {
			if werr.Code == 11000 {
				return &models.ErrClientError{Err: &models.ErrDuplicate{Err: err}}
			}
		}
	}
	return err
}
