// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package database

import (
	"context"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/apex/log"
	"github.com/lrstanley/spectrograph/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// sessionService satisfies the scs.Store interface.
type sessionService struct {
	store *mongoStore

	cleanup time.Duration
	ctx     context.Context
}

// NewSessionService returns a new sessionService that satisfies the scs.Store
// interface. This method also spins up a goroutine that (at a given interval) removes
// old/expired sessions from the database.
func (s *mongoStore) NewSessionService(ctx context.Context, cleanup time.Duration) scs.Store {
	svc := &sessionService{store: s, cleanup: cleanup, ctx: ctx}
	go svc.expiredCleanupWorker()

	return svc
}

func (s *sessionService) expiredCleanupWorker() {
	logger := s.store.log.WithField("interval", s.cleanup)
	logger.Info("initializing expired session worker")
	var trace *log.Entry
	var err error
	for {
		select {
		case <-s.ctx.Done():
			logger.Info("closing expired session worker")
			return
		case <-time.After(s.cleanup):
			trace = logger.Trace("running expired session checks")

			// TODO: do cleanup
			// https://github.com/alexedwards/scs/blob/master/mysqlstore/mysqlstore.go
		}

		trace.Stop(&err)
	}
}

// Delete should remove the session token and corresponding data from the
// session store. If the token does not exist then Delete should be a no-op
// and return nil (not an error).
func (s *sessionService) Delete(token string) (err error) {
	_, err = s.store.session.DeleteOne(context.Background(), bson.M{"token": token})
	if err == mongo.ErrNoDocuments {
		return nil
	}
	return errorWrapper(err)
}

// Find should return the data for a session token from the store. If the
// session token is not found or is expired, the found return value should
// be false (and the err return value should be nil). Similarly, tampered
// or malformed tokens should result in a found return value of false and a
// nil err value. The err return value should be used for system errors only.
func (s *sessionService) Find(token string) (b []byte, found bool, err error) {
	var session *models.Session
	err = s.store.session.FindOne(context.Background(), bson.M{"token": token}).Decode(&session)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, false, nil
		}
		return nil, false, err
	}

	// Check if expired.
	if time.Until(session.Expires) < 0 {
		return nil, false, nil
	}

	return session.Data, true, nil
}

// Commit should add the session token and data to the store, with the given
// expiry time. If the session token already exists, then the data and
// expiry time should be overwritten.
func (s *sessionService) Commit(token string, b []byte, expiry time.Time) (err error) {
	session := &models.Session{
		Token:   token,
		Data:    b,
		Expires: expiry,
	}

	_, err = s.store.session.UpdateOne(context.TODO(),
		bson.M{"token": token},
		bson.M{"$set": session},
		options.Update().SetUpsert(true),
	)

	return err
}
