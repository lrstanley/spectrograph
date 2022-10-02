// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package schema

import (
	"context"
	"strings"

	"github.com/lrstanley/chix"
	"github.com/lrstanley/spectrograph/internal/ent"
	"github.com/lrstanley/spectrograph/internal/ent/guild"
	"github.com/lrstanley/spectrograph/internal/ent/guildadminconfig"
	"github.com/lrstanley/spectrograph/internal/ent/guildconfig"
	"github.com/lrstanley/spectrograph/internal/ent/guildevent"
	"github.com/lrstanley/spectrograph/internal/ent/privacy"
	"github.com/lrstanley/spectrograph/internal/ent/user"
	"github.com/lrstanley/spectrograph/internal/models"
)

func hasRole(ctx context.Context, allowed ...string) bool {
	roles := chix.RolesFromContext(ctx)
	if roles != nil {
		for _, role := range roles {
			for _, allowedRole := range allowed {
				if strings.EqualFold(role, allowedRole) {
					return true
				}
			}
		}
	}

	return false
}

func userID(ctx context.Context) int {
	return chix.IDFromContext[int](ctx)
}

// AllowRoles is a rule that returns Allow decision if the authenticated client
// has ONE of the specified roles.
func AllowRoles(allowed ...string) privacy.QueryMutationRule {
	return privacy.ContextQueryMutationRule(func(ctx context.Context) error {
		if hasRole(ctx, allowed...) {
			return privacy.Allow
		}
		return privacy.Skip
	})
}

func AllowPrivilegedGuildQuery() privacy.QueryRule {
	return privacy.QueryRuleFunc(func(ctx context.Context, query ent.Query) error {
		uid := userID(ctx)

		if uid == 0 {
			return privacy.Skip
		}

		switch q := query.(type) {
		case *ent.GuildQuery:
			q.Where(guild.HasAdminsWith(user.ID(uid)))
		case *ent.GuildEventQuery:
			q.Where(guildevent.HasGuildWith(guild.HasAdminsWith(user.ID(uid))))
		case *ent.GuildConfigQuery:
			q.Where(guildconfig.HasGuildWith(guild.HasAdminsWith(user.ID(uid))))
		case *ent.GuildAdminConfigQuery:
			q.Where(guildadminconfig.HasGuildWith(guild.HasAdminsWith(user.ID(uid))))
		default:
			return privacy.Skip
		}

		return privacy.Allow
	})
}

func AllowPrivilegedGuildMutation() privacy.MutationRule {
	return privacy.MutationRuleFunc(func(ctx context.Context, m ent.Mutation) error {
		uid := userID(ctx)

		if uid == 0 {
			return privacy.Skip
		}

		switch m := m.(type) {
		case *ent.GuildMutation:
			m.Where(guild.HasAdminsWith(user.ID(uid)))
		case *ent.GuildEventMutation:
			m.Where(guildevent.HasGuildWith(guild.HasAdminsWith(user.ID(uid))))
		case *ent.GuildConfigMutation:
			m.Where(guildconfig.HasGuildWith(guild.HasAdminsWith(user.ID(uid))))
		case *ent.GuildAdminConfigMutation:
			m.Where(guildadminconfig.HasGuildWith(guild.HasAdminsWith(user.ID(uid))))
		default:
			return privacy.Skip
		}

		return privacy.Allow
	})
}

func DisallowDebugUnlessAdmin() privacy.QueryRule {
	return privacy.GuildEventQueryRuleFunc(func(ctx context.Context, q *ent.GuildEventQuery) error {
		if !hasRole(ctx, models.RoleSystemAdmin) {
			q.Where(guildevent.TypeNotIn(guildevent.TypeDEBUG))
		}
		return nil
	})
}