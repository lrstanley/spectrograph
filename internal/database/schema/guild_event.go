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

type GuildEvent struct {
	ent.Schema
}

func (GuildEvent) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("type").Values(
			"INFO",
			"WARNING",
			"ERROR",
			"DEBUG",
		).Comment("The type of event that occurred."),
		field.String("message").Comment("The message associated with the event."),
		field.JSON("metadata", map[string]any{}).Optional().Annotations(
			entgql.Type("Map"),
		).Comment("Additional metadata associated with the event."),
	}
}

func (GuildEvent) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

func (GuildEvent) Policy() ent.Policy {
	return privacy.Policy{
		Mutation: privacy.MutationPolicy{
			AllowRoles(models.RoleSystemAdmin),
			privacy.AlwaysDenyRule(),
		},
		Query: privacy.QueryPolicy{
			AllowRoles(models.RoleSystemAdmin),
			AllowPrivilegedGuildQuery(),
			DisallowDebugUnlessAdmin(),
			privacy.AlwaysDenyRule(),
		},
	}
}

func (GuildEvent) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("guild", Guild.Type).Ref("guild_events").
			Unique().
			Required().
			Immutable().
			Comment("The guild these events belong to."),
	}
}

func (GuildEvent) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.RelayConnection(),
		entgql.QueryField(),
	}
}
