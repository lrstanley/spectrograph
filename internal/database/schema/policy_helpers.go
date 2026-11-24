//go:build !skippolicy

// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package schema

import (
	"context"

	"github.com/lrstanley/chix/xauth/v2"
	"github.com/lrstanley/spectrograph/internal/database/ent"
	"github.com/lrstanley/spectrograph/internal/database/ent/guild"
	"github.com/lrstanley/spectrograph/internal/database/ent/guildadminconfig"
	"github.com/lrstanley/spectrograph/internal/database/ent/guildconfig"
	"github.com/lrstanley/spectrograph/internal/database/ent/guildevent"
	"github.com/lrstanley/spectrograph/internal/database/ent/privacy"
	"github.com/lrstanley/spectrograph/internal/database/ent/user"
)

// authUserID returns the authenticated (and not banned) user ID from the context,
// or 0 if the request is unauthenticated.
func authUserID(ctx context.Context) int {
	u := xauth.IdentFromContext[ent.User](ctx)
	if u == nil || u.Banned {
		return 0
	}
	return u.ID
}

// IsSystemAdmin returns true if the authenticated user in the context is a
// system (spectrograph) admin.
func IsSystemAdmin(ctx context.Context) bool {
	u := xauth.IdentFromContext[ent.User](ctx)
	return u != nil && u.Admin && !u.Banned
}

// AllowSystemAdmin is a rule that returns an allow decision if the authenticated
// user is a system admin.
func AllowSystemAdmin() privacy.QueryMutationRule {
	return privacy.ContextQueryMutationRule(func(ctx context.Context) error {
		if IsSystemAdmin(ctx) {
			return privacy.Allow
		}
		return privacy.Skip
	})
}

// AllowPrivilegedGuildQuery restricts guild-related queries to guilds that the
// authenticated user is an admin (or owner) of.
func AllowPrivilegedGuildQuery() privacy.QueryRule {
	return privacy.QueryRuleFunc(func(ctx context.Context, query ent.Query) error {
		uid := authUserID(ctx)
		if uid == 0 {
			return privacy.Skip
		}

		isAdmin := guild.HasAdminsWith(user.ID(uid))

		switch q := query.(type) {
		case *ent.GuildQuery:
			q.Where(isAdmin)
		case *ent.GuildEventQuery:
			q.Where(guildevent.HasGuildWith(isAdmin))
		case *ent.GuildConfigQuery:
			q.Where(guildconfig.HasGuildWith(isAdmin))
		case *ent.GuildAdminConfigQuery:
			q.Where(guildadminconfig.HasGuildWith(isAdmin))
		default:
			return privacy.Skip
		}

		return privacy.Allow
	})
}

// AllowPrivilegedGuildMutation restricts guild-related mutations to guilds that
// the authenticated user is an admin (or owner) of.
func AllowPrivilegedGuildMutation() privacy.MutationRule {
	return privacy.MutationRuleFunc(func(ctx context.Context, m ent.Mutation) error {
		uid := authUserID(ctx)
		if uid == 0 {
			return privacy.Skip
		}

		isAdmin := guild.HasAdminsWith(user.ID(uid))

		switch m := m.(type) {
		case *ent.GuildMutation:
			m.Where(isAdmin)
		case *ent.GuildEventMutation:
			m.Where(guildevent.HasGuildWith(isAdmin))
		case *ent.GuildConfigMutation:
			m.Where(guildconfig.HasGuildWith(isAdmin))
		case *ent.GuildAdminConfigMutation:
			m.Where(guildadminconfig.HasGuildWith(isAdmin))
		default:
			return privacy.Skip
		}

		return privacy.Allow
	})
}

// FilterDebugUnlessSystemAdmin filters out DEBUG guild events, unless the
// authenticated user is a system admin. This rule never returns a decision, so
// it must be placed before any allow rules.
func FilterDebugUnlessSystemAdmin() privacy.QueryRule {
	return privacy.GuildEventQueryRuleFunc(func(ctx context.Context, q *ent.GuildEventQuery) error {
		if !IsSystemAdmin(ctx) {
			q.Where(guildevent.TypeNEQ(guildevent.TypeDEBUG))
		}
		return privacy.Skip
	})
}

// AllowUserQuerySelf restricts user queries to the authenticated user.
func AllowUserQuerySelf() privacy.QueryRule {
	return privacy.UserQueryRuleFunc(func(ctx context.Context, q *ent.UserQuery) error {
		uid := authUserID(ctx)
		if uid == 0 {
			return privacy.Skip
		}

		q.Where(user.ID(uid))
		return privacy.Allow
	})
}

// AllowUserMutateSelf restricts user mutations to the authenticated user.
func AllowUserMutateSelf() privacy.MutationRule {
	return privacy.UserMutationRuleFunc(func(ctx context.Context, m *ent.UserMutation) error {
		uid := authUserID(ctx)
		if uid == 0 {
			return privacy.Skip
		}

		m.Where(user.ID(uid))
		return privacy.Allow
	})
}
