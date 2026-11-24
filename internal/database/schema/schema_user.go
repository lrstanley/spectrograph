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

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("user_id").
			Unique().
			Immutable().
			Annotations(
				entrest.WithReadOnly(true),
				entrest.WithFilter(entrest.FilterGroupEqualExact|entrest.FilterGroupArray),
			).
			Comment("Discord snowflake ID of the user."),
		field.Bool("admin").
			Optional().
			Annotations(
				entrest.WithReadOnly(true),
				entrest.WithFilter(entrest.FilterEQ),
			).
			Comment("Whether or not the user is a spectrograph admin."),
		field.Bool("banned").
			Optional().
			Annotations(
				entrest.WithReadOnly(true),
				entrest.WithFilter(entrest.FilterEQ),
			).
			Comment("Whether or not the user is banned from using the service."),
		field.String("ban_reason").
			Optional().
			Annotations(entrest.WithReadOnly(true)).
			Comment("Reason for the user being banned (if any)."),
		field.String("username").
			Annotations(
				entrest.WithReadOnly(true),
				entrest.WithSortable(true),
				entrest.WithFilter(entrest.FilterGroupEqual|entrest.FilterGroupArray),
			).
			Comment("The users username, unique across the platform (for migrated users)."),
		field.String("global_name").
			Optional().
			MaxLen(100).
			Annotations(
				entrest.WithReadOnly(true),
				entrest.WithSortable(true),
				entrest.WithFilter(entrest.FilterGroupEqual),
			).
			Comment("The users display name, if set."),
		field.String("discriminator").
			Optional().
			MaxLen(4).
			Annotations(entrest.WithReadOnly(true)).
			Comment("The users legacy 4-digit discord-tag (\"0\" or empty for migrated users)."),
		field.String("email").
			Optional().
			MaxLen(320).
			Annotations(
				entrest.WithReadOnly(true),
				entrest.WithSortable(true),
				entrest.WithFilter(entrest.FilterGroupEqual),
			).
			Comment("The users email address."),
		field.String("avatar_hash").
			Optional().
			MaxLen(2048).
			Annotations(entrest.WithReadOnly(true)).
			Comment("The users avatar hash."),
		field.String("avatar_url").
			MaxLen(2048).
			Annotations(entrest.WithReadOnly(true)).
			Comment("The users avatar URL (generated if no avatar present)."),

		// Additional fields provided by querying discord directly.
		field.String("locale").
			Optional().
			MaxLen(10).
			Annotations(entrest.WithReadOnly(true)).
			Comment("The users chosen language option."),
		field.Bool("bot").
			Optional().
			Annotations(entrest.WithReadOnly(true)).
			Comment("Whether the user belongs to an OAuth2 application."),
		field.Bool("system").
			Optional().
			Annotations(entrest.WithReadOnly(true)).
			Comment("Whether the user is an Official Discord System user (part of the urgent message system)."),
		field.Bool("mfa_enabled").
			Optional().
			Annotations(entrest.WithReadOnly(true)).
			Comment("Whether the user has two factor enabled on their account."),
		field.Bool("verified").
			Optional().
			Annotations(entrest.WithReadOnly(true)).
			Comment("Whether the email on this account has been verified."),
		field.Uint64("flags").
			Optional().
			Annotations(entrest.WithReadOnly(true)).
			Comment("The flags on a users account."),
		field.Int("premium_type").
			Optional().
			Annotations(entrest.WithReadOnly(true)).
			Comment("The type of Nitro subscription on a users account."),
		field.Uint64("public_flags").
			Optional().
			Annotations(entrest.WithReadOnly(true)).
			Comment("The public flags on a users account."),
	}
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		MixinTime{},
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user_guilds", Guild.Type).
			Annotations(entrest.WithOperationID(entrest.OperationList, "listGuildsForUser")).
			Comment("Guilds that the user is either owner or admin of (and thus can add the connection to the bot)."),
		edge.To("banned_users", User.Type).
			Annotations(entrest.WithEdgeEndpoint(false)).
			Comment("Users that were banned by this user."),
		edge.From("banned_by", User.Type).
			Unique().
			Ref("banned_users").
			Comment("User that banned this user."),
	}
}

func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		// Users are only created/updated through authentication, and ban/unban/delete
		// are handled through custom endpoints.
		entrest.WithIncludeOperations(entrest.OperationRead, entrest.OperationList),
	}
}
