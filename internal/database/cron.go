// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package database

import (
	"context"
	"errors"
	"time"

	"github.com/apex/log"
	"github.com/lrstanley/spectrograph/internal/ent"
	"github.com/lrstanley/spectrograph/internal/ent/guildevent"
	"github.com/lrstanley/spectrograph/internal/ent/privacy"
)

const (
	oldestGuildEvent = 7 * 24 * time.Hour
)

func CronGuildEvents(ctx context.Context) error {
	logger := log.FromContext(ctx).WithField("src", "db-cron-guild-events")
	db := ent.FromContext(privacy.DecisionContext(ctx, privacy.Allow))
	if db == nil {
		return errors.New("database not found in context")
	}

	var count int
	var err error

	timer := time.NewTicker(10 * time.Minute)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-timer.C:
			if count, err = db.GuildEvent.Delete().Where(
				guildevent.CreateTimeLT(time.Now().Add(-oldestGuildEvent)),
			).Exec(ctx); err != nil {
				logger.WithError(err).Error("failed to delete old guild events")
			}

			if count > 0 {
				logger.WithField("count", count).Info("pruned old guild events")
			}
		}
	}
}
