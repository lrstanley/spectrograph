// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package schema

import (
	"errors"
	"net/mail"
	"regexp"

	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/lrstanley/entrest"
)

type GuildConfig struct {
	ent.Schema
}

func (GuildConfig) Fields() []ent.Field {
	return []ent.Field{
		field.Bool("enabled").
			Optional().
			Default(true).
			Comment("True if the guild should be monitored/acted upon."),
		field.Int("default_max_clones").
			Optional().
			Default(0).
			Min(0).
			Max(20).
			Comment("Default max cloned channels for the guild (0 uses the system default)."),
		field.String("regex_match").
			Optional().
			Default("").
			MaxLen(100).
			Validate(validateRegex).
			Comment("Regex match for channel names (empty uses the system default)."),
		field.String("contact_email").
			Optional().
			Default("").
			MaxLen(255).
			Validate(validateEmail).
			Comment("Contact email for the guild."),
	}
}

func (GuildConfig) Mixin() []ent.Mixin {
	return []ent.Mixin{
		MixinTime{},
	}
}

func (GuildConfig) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("guild", Guild.Type).
			Ref("guild_config").
			Unique().
			Required().
			Immutable().
			Comment("The guild these settings belong to."),
	}
}

func (GuildConfig) Annotations() []schema.Annotation {
	return []schema.Annotation{
		// Created by the worker when the bot joins a guild.
		entrest.WithIncludeOperations(entrest.OperationRead, entrest.OperationUpdate),
	}
}

func validateRegex(input string) error {
	if input == "" {
		return nil
	}
	_, err := regexp.Compile(input)
	return err
}

func validateEmail(input string) error {
	if input == "" {
		return nil
	}
	addr, err := mail.ParseAddress(input)
	if err != nil || addr.Address != input {
		return errors.New("invalid email address")
	}
	return nil
}
