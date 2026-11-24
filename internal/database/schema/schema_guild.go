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

type Guild struct {
	ent.Schema
}

func (Guild) Fields() []ent.Field {
	return []ent.Field{
		field.String("guild_id").
			Unique().
			Immutable().
			Annotations(
				entrest.WithReadOnly(true),
				entrest.WithFilter(entrest.FilterGroupEqualExact|entrest.FilterGroupArray),
			).
			Comment("Discord snowflake ID of the guild."),
		field.String("name").
			MinLen(2).
			MaxLen(100).
			Annotations(
				entrest.WithReadOnly(true),
				entrest.WithSortable(true),
				entrest.WithFilter(entrest.FilterGroupEqual),
			).
			Comment("Guild name (2-100 chars, excl. trailing/leading spaces)."),
		field.Strings("features").
			Optional().
			Default([]string{}).
			Annotations(entrest.WithReadOnly(true)).
			Comment("Enabled guild features."),
		field.String("icon_hash").
			Optional().
			MaxLen(2048).
			Annotations(entrest.WithReadOnly(true)).
			Comment("Icon hash."),
		field.String("icon_url").
			MaxLen(2048).
			Annotations(entrest.WithReadOnly(true)).
			Comment("Icon URL (empty if the guild has no icon)."),
		field.Time("joined_at").
			Optional().
			Nillable().
			Annotations(
				entrest.WithReadOnly(true),
				entrest.WithSortable(true),
				entrest.WithFilter(entrest.FilterGroupNil|entrest.FilterGroupLength),
			).
			Comment("When the bot joined the guild (null if the bot is not in the guild)."),
		field.Bool("large").
			Optional().
			Default(false).
			Annotations(entrest.WithReadOnly(true)).
			Comment("True if the guild is considered large (according to Discord standards)."),
		field.Int("member_count").
			Optional().
			Default(0).
			Annotations(
				entrest.WithReadOnly(true),
				entrest.WithSortable(true),
			).
			Comment("Total number of members in the guild."),
		field.String("owner_id").
			Optional().
			Annotations(
				entrest.WithReadOnly(true),
				entrest.WithFilter(entrest.FilterGroupEqualExact),
			).
			Comment("Discord snowflake ID of the user that owns the guild."),
		field.Uint64("permissions").
			Optional().
			Default(0).
			Annotations(entrest.WithReadOnly(true)).
			Comment("Permissions of the bot on this guild (excludes overrides)."),
		field.String("system_channel_flags").
			Optional().
			MaxLen(32).
			Annotations(entrest.WithReadOnly(true)).
			Comment("System channel flags."),
	}
}

func (Guild) Mixin() []ent.Mixin {
	return []ent.Mixin{
		MixinTime{},
	}
}

func (Guild) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("guild_config", GuildConfig.Type).
			Unique().
			Annotations(
				entrest.WithEagerLoad(true),
				entrest.WithOperationID(entrest.OperationRead, "getGuildConfigForGuild"),
			).
			Comment("User-defined configuration for the guild."),
		edge.To("guild_admin_config", GuildAdminConfig.Type).
			Unique().
			Annotations(
				entrest.WithEagerLoad(true),
				entrest.WithOperationID(entrest.OperationRead, "getGuildAdminConfigForGuild"),
			).
			Comment("Admin-defined configuration for the guild (overrides user-defined configuration)."),
		edge.To("guild_events", GuildEvent.Type).
			Annotations(entrest.WithOperationID(entrest.OperationList, "listGuildEventsForGuild")).
			Comment("Events (logs) associated with the guild."),
		edge.From("admins", User.Type).
			Ref("user_guilds").
			Comment("The users that are an admin (or owner) of this guild."),
	}
}

func (Guild) Annotations() []schema.Annotation {
	return []schema.Annotation{
		// Guilds are only created/updated by the worker or during authentication.
		entrest.WithIncludeOperations(entrest.OperationRead, entrest.OperationList),
	}
}
