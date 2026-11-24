//go:build !skippolicy

// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package schema

import (
	"entgo.io/ent"
	"github.com/lrstanley/spectrograph/internal/database/ent/privacy"
)

func (User) Policy() ent.Policy {
	return privacy.Policy{
		Mutation: privacy.MutationPolicy{
			AllowSystemAdmin(),
			privacy.OnMutationOperation(AllowUserMutateSelf(), ent.OpDelete|ent.OpDeleteOne),
			privacy.AlwaysDenyRule(),
		},
		Query: privacy.QueryPolicy{
			AllowSystemAdmin(),
			AllowUserQuerySelf(),
			privacy.AlwaysDenyRule(),
		},
	}
}
