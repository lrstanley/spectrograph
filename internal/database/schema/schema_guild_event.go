// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/lrstanley/entrest"
)

type GuildEvent struct {
	ent.Schema
}

func (GuildEvent) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("type").
			Values(
				"INFO",
				"WARNING",
				"ERROR",
				"DEBUG",
			).
			Annotations(
				entrest.WithSortable(true),
				entrest.WithFilter(entrest.FilterGroupEqualExact|entrest.FilterGroupArray),
			).
			Comment("The type of event that occurred."),
		field.String("message").
			Annotations(entrest.WithFilter(entrest.FilterGroupContains)).
			Comment("The message associated with the event."),
		field.JSON("metadata", map[string]any{}).
			Optional().
			Annotations(entrest.WithSchema(entrest.SchemaObjectAny)).
			Comment("Additional metadata associated with the event."),
	}
}

func (GuildEvent) Mixin() []ent.Mixin {
	return []ent.Mixin{
		MixinTime{},
	}
}

func (GuildEvent) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("guild", Guild.Type).
			Ref("guild_events").
			Unique().
			Required().
			Immutable().
			// Eager-loaded so event lists across guilds (e.g. for system admins) can
			// show the guild. Privacy rules still apply to the eager-loaded guild.
			Annotations(entrest.WithEagerLoad(true)).
			Comment("The guild these events belong to."),
	}
}

func (GuildEvent) Indexes() []ent.Index {
	return []ent.Index{
		// Event lists and streams are almost always for a single guild.
		index.Edges("guild"),
	}
}

func (GuildEvent) Annotations() []schema.Annotation {
	return []schema.Annotation{
		// Events are only created by the worker.
		entrest.WithIncludeOperations(entrest.OperationRead, entrest.OperationList),
		entrest.WithDefaultSort("create_time"),
		entrest.WithDefaultOrder(entrest.OrderDesc),
	}
}
