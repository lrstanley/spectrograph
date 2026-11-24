//go:build !skippolicy

// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package schema

import (
	"entgo.io/ent"
	"github.com/lrstanley/spectrograph/internal/database/ent/privacy"
)

func (GuildEvent) Policy() ent.Policy {
	return privacy.Policy{
		Mutation: privacy.MutationPolicy{
			AllowSystemAdmin(),
			privacy.AlwaysDenyRule(),
		},
		Query: privacy.QueryPolicy{
			// Must be evaluated before any allow rules, as policy evaluation stops
			// at the first allow decision.
			FilterDebugUnlessSystemAdmin(),
			AllowSystemAdmin(),
			AllowPrivilegedGuildQuery(),
			privacy.AlwaysDenyRule(),
		},
	}
}
