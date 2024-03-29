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

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("user_id").Unique().Immutable(),
		field.Bool("admin").Optional().Comment("Whether or not the user is a spectrograph admin."),
		field.Bool("banned").Optional().Comment("Whether or not the user is banned from using the service."),
		field.String("ban_reason").Optional().Comment("Reason for the user being banned (if any)."),
		field.String("username").Annotations(
			entgql.OrderField("USERNAME"),
		).Comment("The users username, not unique across the platform."),
		field.String("discriminator").MaxLen(4).Annotations(
			entgql.OrderField("DISCRIMINATOR"),
		).Comment("The users 4-digit discord-tag."),
		field.String("email").MaxLen(320).Annotations(
			entgql.OrderField("EMAIL"),
		).Comment("The users email address."),
		field.String("avatar_hash").Optional().MaxLen(2048).Comment("The users avatar hash."),
		field.String("avatar_url").MaxLen(2048).Comment("The users avatar URL (generated if no avatar present)."),

		// Additional fields provided by querying discord directly.
		field.String("locale").Optional().MaxLen(10).Comment("The users chosen language option."),
		field.Bool("bot").Optional().Comment("Whether the user belongs to an OAuth2 application."),
		field.Bool("system").Optional().Comment("Whether the user is an Official Discord System user (part of the urgent message system)."),
		field.Bool("mfa_enabled").Optional().Comment("Whether the user has two factor enabled on their account."),
		field.Bool("verified").Optional().Comment("Whether the email on this account has been verified."),
		field.Uint64("flags").Optional().Annotations(
			entgql.Type("Uint64"),
		).Comment("The flags on a users account."),
		field.Int("premium_type").Optional().Comment("The type of Nitro subscription on a users account."),
		field.Uint64("public_flags").Optional().Annotations(
			entgql.Type("Uint64"),
		).Comment("The public flags on a users account."),
	}
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

func (User) Policy() ent.Policy {
	return privacy.Policy{
		Mutation: privacy.MutationPolicy{
			AllowRoles(models.RoleSystemAdmin),
			privacy.OnMutationOperation(AllowUserMutateSelf(), ent.OpDelete|ent.OpDeleteOne),
			privacy.AlwaysDenyRule(),
		},
		Query: privacy.QueryPolicy{
			AllowRoles(models.RoleSystemAdmin),
			AllowUserQuerySelf(),
			privacy.AlwaysDenyRule(),
		},
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user_guilds", Guild.Type).Annotations(
			entgql.RelayConnection(),
		).Comment("Guilds that the user is either owner or admin of (and thus can add the connection to the bot)."),
		edge.To("banned_users", User.Type).Comment("Users that were banned by this user."),
		edge.From("banned_by", User.Type).Unique().Ref("banned_users").Comment("User that banned this user."),
	}
}

func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.RelayConnection(),
		entgql.QueryField(),
	}
}
