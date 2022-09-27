// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package schema

import (
	"regexp"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/lrstanley/spectrograph/internal/ent/privacy"
	"github.com/lrstanley/spectrograph/internal/models"
)

type GuildConfig struct {
	ent.Schema
}

func (GuildConfig) Fields() []ent.Field {
	return []ent.Field{
		field.Bool("enabled").Optional().Default(true).Comment("True if the guild should be monitored/acted upon."),
		field.Int("default_max_clones").Optional().Default(0).Comment("Default max cloned channels for the guild."),
		field.String("regex_match").Optional().Default("").MaxLen(100).Validate(func(input string) error {
			if input == "" {
				return nil
			}

			_, err := regexp.Compile(input)
			return err
		}).Comment("Regex match for channel names."),
		field.String("contact_email").Optional().Default("").MaxLen(255).Comment("Contact email for the guild."),
	}
}

func (GuildConfig) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

func (GuildConfig) Policy() ent.Policy {
	return privacy.Policy{
		Mutation: privacy.MutationPolicy{
			AllowRoles(models.RoleSystemAdmin),
			privacy.OnMutationOperation(AllowPrivilegedGuildMutation(), ent.OpUpdate),
			privacy.AlwaysDenyRule(),
		},
		Query: privacy.QueryPolicy{
			AllowRoles(models.RoleSystemAdmin),
			AllowPrivilegedGuildQuery(),
			privacy.AlwaysDenyRule(),
		},
	}
}

func (GuildConfig) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("guild", Guild.Type).Ref("guild_config").
			Unique().
			Required().
			Immutable().
			Comment("The guild these settings belong to."),
	}
}

func (GuildConfig) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.RelayConnection(),
		entgql.QueryField(),
		entgql.Mutations(entgql.MutationUpdate()),
	}
}
