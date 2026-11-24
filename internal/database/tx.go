// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/lrstanley/spectrograph/internal/database/ent"
)

// RunWithTx runs fn within a transaction, committing if fn returns nil, and
// rolling back otherwise. The transaction is also attached to the context passed
// to fn (see [ent.TxFromContext]), so mutation observers (see
// [GuildEventNotifier]) are only notified once it commits.
func RunWithTx(ctx context.Context, db *ent.Client, fn func(ctx context.Context, tx *ent.Tx) error) error {
	tx, err := db.Tx(ctx)
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}

	if err = fn(ent.NewTxContext(ctx, tx), tx); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			return errors.Join(err, fmt.Errorf("rolling back transaction: %w", rerr))
		}
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	return nil
}
