// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/lrstanley/spectrograph/internal/database/ent"
	"github.com/lrstanley/spectrograph/internal/database/ent/guildevent"
	"github.com/lrstanley/spectrograph/internal/database/ent/privacy"
)

// MaxGuildEventAge is the maximum age of guild events before they are pruned.
const MaxGuildEventAge = 7 * 24 * time.Hour

// PruneGuildEvents removes guild events older than [MaxGuildEventAge]. Intended
// to be run on an interval.
func PruneGuildEvents(db *ent.Client) func(ctx context.Context, logger *slog.Logger) error {
	return func(ctx context.Context, logger *slog.Logger) error {
		count, err := db.GuildEvent.Delete().
			Where(guildevent.CreateTimeLT(time.Now().Add(-MaxGuildEventAge))).
			Exec(privacy.DecisionContext(ctx, privacy.Allow))
		if err != nil {
			return fmt.Errorf("pruning old guild events: %w", err)
		}

		if count > 0 {
			logger.InfoContext(ctx, "pruned old guild events", "count", count)
		}
		return nil
	}
}
