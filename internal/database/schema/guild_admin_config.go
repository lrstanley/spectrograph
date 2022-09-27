// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/lrstanley/spectrograph/internal/ent/privacy"
	"github.com/lrstanley/spectrograph/internal/models"
)

type GuildAdminConfig struct {
	ent.Schema
}

func (GuildAdminConfig) Fields() []ent.Field {
	return []ent.Field{
		field.Bool("enabled").Optional().Default(true).Comment("True if the guild should be monitored/acted upon (overrides user-defined settings)."),
		field.Int("default_max_channels").Optional().Default(0).Comment("Default max channels for the guild (overrides user-defined settings)."),
		field.Int("default_max_clones").Optional().Default(0).Comment("Default max clones for the guild (overrides user-defined settings)."),
		field.String("comment").Optional().Default("").Comment("Admin comment for the guild."),
	}
}

func (GuildAdminConfig) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

func (GuildAdminConfig) Policy() ent.Policy {
	return privacy.Policy{
		Mutation: privacy.MutationPolicy{
			AllowRoles(models.RoleSystemAdmin),
			privacy.AlwaysDenyRule(),
		},
		Query: privacy.QueryPolicy{
			AllowRoles(models.RoleSystemAdmin),
			AllowPrivilegedGuildQuery(),
			privacy.AlwaysDenyRule(),
		},
	}
}

func (GuildAdminConfig) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("guild", Guild.Type).Ref("guild_admin_config").
			Unique().
			Required().
			Immutable().
			Comment("The guild these settings belong to."),
	}
}

func (GuildAdminConfig) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.RelayConnection(),
		entgql.QueryField(),
		entgql.Mutations(entgql.MutationUpdate()),
	}
}
