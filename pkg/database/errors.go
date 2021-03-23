// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package database

import (
	"errors"

	"github.com/lrstanley/spectrograph/pkg/models"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	// ErrInvalidObjectID is returned when the provided ID doesn't match an
	// Object ID.
	ErrInvalidObjectID = &models.ErrClientError{Err: errors.New("invalid ID provided")}
)

// errorWrapper is a useful wrapper for returning customized errors from mongo.
func errorWrapper(err error) error {
	switch err {
	case mongo.ErrNoDocuments:
		return &models.ErrClientError{Err: &models.ErrNotFound{Err: err}}
	default:
	}

	if terr, ok := err.(mongo.WriteException); ok {
		for _, werr := range terr.WriteErrors {
			if werr.Code == 11000 {
				return &models.ErrClientError{Err: &models.ErrDuplicate{Err: err}}
			}
		}
	}
	return err
}

// mongo.WriteException{
// 	WriteConcernError: (*mongo.WriteConcernError)(nil),
// 	WriteErrors:       {
// 		{Index:0, Code:11000, Message:"E11000 duplicate key error collection: spectrograph-dev.scan_requests index: scan_requests_unique dup key: { project_id: \"5ff2922bd091e4695cc98775\", reference: \"refs/heads/master\" }"},
// 	},
// 	Labels: nil,
// }
