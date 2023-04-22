// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package schema

import (
	"errors"
	"regexp"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/lrstanley/spectrograph/internal/database/ent/privacy"
	"github.com/lrstanley/spectrograph/internal/models"
)

var emailRegex = regexp.MustCompile("^(?:(?:(?:(?:[a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(?:\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|(?:(?:\\x22)(?:(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(?:\\x20|\\x09)+)?(?:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(\\x20|\\x09)+)?(?:\\x22))))@(?:(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$")

type GuildConfig struct {
	ent.Schema
}

func (GuildConfig) Fields() []ent.Field {
	return []ent.Field{
		field.Bool("enabled").Optional().Default(true).Comment("True if the guild should be monitored/acted upon."),
		field.Int("default_max_clones").Optional().Default(0).Min(0).Max(20).Comment("Default max cloned channels for the guild."),
		field.String("regex_match").Optional().Default("").MaxLen(100).Validate(func(input string) error {
			if input == "" {
				return nil
			}

			_, err := regexp.Compile(input)
			return err
		}).Comment("Regex match for channel names."),
		field.String("contact_email").Optional().Default("").MaxLen(255).Validate(func(email string) error {
			if email != "" && !emailRegex.MatchString(email) {
				return errors.New("invalid email address")
			}

			return nil
		}).Comment("Contact email for the guild."),
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
