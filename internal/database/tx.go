// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package database

import (
	"context"
	"fmt"

	"github.com/apex/log"
	"github.com/lrstanley/spectrograph/internal/ent"
)

// Commit will commit the transaction, unless the provided error is non-nil, in
// which case it will rollback the transaction and return a wrapped error.
func Commit(tx *ent.Tx, err error) error {
	if err != nil {
		return Rollback(tx, err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Rollback will rollback the transaction and return a wrapped error of the original
// error.
func Rollback(tx *ent.Tx, err error) error {
	if rerr := tx.Rollback(); rerr != nil {
		err = fmt.Errorf("%w: %v", err, rerr)
	}

	return err
}

type TxFn func(ctx context.Context, logger log.Interface, db *ent.Tx) error

func RunWithTx(ctx context.Context, logger log.Interface, db *ent.Client, fn TxFn) error {
	tx, err := db.Tx(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	return Commit(tx, fn(ctx, logger, tx))
}
