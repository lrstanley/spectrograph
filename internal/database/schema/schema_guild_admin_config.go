// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/lrstanley/entrest"
)

type GuildAdminConfig struct {
	ent.Schema
}

func (GuildAdminConfig) Fields() []ent.Field {
	return []ent.Field{
		field.Bool("enabled").
			Optional().
			Default(true).
			Comment("True if the guild should be monitored/acted upon (overrides user-defined settings)."),
		field.Int("default_max_channels").
			Optional().
			Default(0).
			Min(0).
			Comment("Default max channels for the guild (overrides user-defined settings, 0 uses the system default)."),
		field.Int("default_max_clones").
			Optional().
			Default(0).
			Min(0).
			Comment("Default max clones for the guild (overrides user-defined settings, 0 uses the system default)."),
		field.String("comment").
			Optional().
			Default("").
			Comment("Admin comment for the guild."),
	}
}

func (GuildAdminConfig) Mixin() []ent.Mixin {
	return []ent.Mixin{
		MixinTime{},
	}
}

func (GuildAdminConfig) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("guild", Guild.Type).
			Ref("guild_admin_config").
			Unique().
			Required().
			Immutable().
			Comment("The guild these settings belong to."),
	}
}

func (GuildAdminConfig) Annotations() []schema.Annotation {
	return []schema.Annotation{
		// Created by the worker when the bot joins a guild. Updates are restricted
		// to system admins through privacy policies.
		entrest.WithIncludeOperations(entrest.OperationRead, entrest.OperationUpdate),
	}
}
