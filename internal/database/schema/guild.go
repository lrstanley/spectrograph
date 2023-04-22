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
	"github.com/lrstanley/spectrograph/internal/database/ent/privacy"
	"github.com/lrstanley/spectrograph/internal/models"
)

type Guild struct {
	ent.Schema
}

func (Guild) Fields() []ent.Field {
	return []ent.Field{
		field.String("guild_id").Unique().Immutable().Comment("Guild id."),
		field.String("name").Annotations(
			entgql.OrderField("NAME"),
		).MinLen(2).MaxLen(100).Comment("Guild name (2-100 chars, excl. trailing/leading spaces)."),
		field.Strings("features").Optional().Default([]string{}).Comment("Enabled guild features."),
		field.String("icon_hash").Optional().MaxLen(2048).Comment("Icon hash."),
		field.String("icon_url").MaxLen(2048),
		field.Time("joined_at").Optional().Nillable().Annotations(
			entgql.OrderField("JOINED_AT"),
		).Comment("When the bot joined the guild."),
		field.Bool("large").Optional().Default(false).Comment("True if the guild is considered large (according to Discord standards)."),
		field.Int("member_count").Optional().Default(0).Comment("Total number of members in the guild."),
		field.String("owner_id").Optional().Comment("Discord snowflake ID of the user that owns the guild."),
		field.Uint64("permissions").Optional().Default(0).Annotations(
			entgql.Type("Uint64"),
		).Comment("Permissions of the bot on this guild (excludes overrides)."),
		field.String("system_channel_flags").Optional().MaxLen(32).Comment("System channel flags."),
	}
}

func (Guild) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

func (Guild) Policy() ent.Policy {
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

func (Guild) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("guild_config", GuildConfig.Type).Unique(),
		edge.To("guild_admin_config", GuildAdminConfig.Type).Unique(),
		edge.To("guild_events", GuildEvent.Type),

		edge.From("admins", User.Type).Ref("user_guilds").
			Annotations(entgql.RelayConnection()).
			Comment("The users that are an admin (or owner) of this guild."),
	}
}

func (Guild) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.RelayConnection(),
		entgql.QueryField(),
	}
}
