// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	AuthCount = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "spectrograph_auth_count_total",
			Help: "The total number of authentications",
		},
	)

	WorkerGuildCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "spectrograph_worker_guild_count",
			Help: "The number of guilds the worker is currently in",
		},
		[]string{"shard_id"},
	)

	WorkerEventCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "spectrograph_worker_event_count_total",
			Help: "The total number of worker events",
		},
		[]string{"shard_id", "event_type", "guild_id"},
	)

	WorkerSessionsRemainingCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "spectrograph_worker_sessions_remaining_count",
			Help: "The number of sessions remaining for a worker",
		},
		[]string{"shard_id"},
	)

	WorkerSessionsTotalCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "spectrograph_worker_sessions_total_count",
			Help: "The number of total sessions allowed for this shard ID",
		},
		[]string{"shard_id"},
	)
)
