// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.
//
// DO NOT EDIT, CODE GENERATED BY entc.

package user

import (
	"time"

	"entgo.io/ent"
)

const (
	// Label holds the string label denoting the user type in the database.
	Label = "user"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreateTime holds the string denoting the create_time field in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the update_time field in the database.
	FieldUpdateTime = "update_time"
	// FieldUserID holds the string denoting the user_id field in the database.
	FieldUserID = "user_id"
	// FieldAdmin holds the string denoting the admin field in the database.
	FieldAdmin = "admin"
	// FieldUsername holds the string denoting the username field in the database.
	FieldUsername = "username"
	// FieldDiscriminator holds the string denoting the discriminator field in the database.
	FieldDiscriminator = "discriminator"
	// FieldEmail holds the string denoting the email field in the database.
	FieldEmail = "email"
	// FieldAvatarHash holds the string denoting the avatar_hash field in the database.
	FieldAvatarHash = "avatar_hash"
	// FieldAvatarURL holds the string denoting the avatar_url field in the database.
	FieldAvatarURL = "avatar_url"
	// FieldLocale holds the string denoting the locale field in the database.
	FieldLocale = "locale"
	// FieldBot holds the string denoting the bot field in the database.
	FieldBot = "bot"
	// FieldSystem holds the string denoting the system field in the database.
	FieldSystem = "system"
	// FieldMfaEnabled holds the string denoting the mfa_enabled field in the database.
	FieldMfaEnabled = "mfa_enabled"
	// FieldVerified holds the string denoting the verified field in the database.
	FieldVerified = "verified"
	// FieldFlags holds the string denoting the flags field in the database.
	FieldFlags = "flags"
	// FieldPremiumType holds the string denoting the premium_type field in the database.
	FieldPremiumType = "premium_type"
	// FieldPublicFlags holds the string denoting the public_flags field in the database.
	FieldPublicFlags = "public_flags"
	// EdgeUserGuilds holds the string denoting the user_guilds edge name in mutations.
	EdgeUserGuilds = "user_guilds"
	// Table holds the table name of the user in the database.
	Table = "users"
	// UserGuildsTable is the table that holds the user_guilds relation/edge. The primary key declared below.
	UserGuildsTable = "user_user_guilds"
	// UserGuildsInverseTable is the table name for the Guild entity.
	// It exists in this package in order to avoid circular dependency with the "guild" package.
	UserGuildsInverseTable = "guilds"
)

// Columns holds all SQL columns for user fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldUserID,
	FieldAdmin,
	FieldUsername,
	FieldDiscriminator,
	FieldEmail,
	FieldAvatarHash,
	FieldAvatarURL,
	FieldLocale,
	FieldBot,
	FieldSystem,
	FieldMfaEnabled,
	FieldVerified,
	FieldFlags,
	FieldPremiumType,
	FieldPublicFlags,
}

var (
	// UserGuildsPrimaryKey and UserGuildsColumn2 are the table columns denoting the
	// primary key for the user_guilds relation (M2M).
	UserGuildsPrimaryKey = []string{"user_id", "guild_id"}
)

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

// Note that the variables below are initialized by the runtime
// package on the initialization of the application. Therefore,
// it should be imported in the main as follows:
//
//	import _ "github.com/lrstanley/spectrograph/internal/ent/runtime"
var (
	Hooks  [1]ent.Hook
	Policy ent.Policy
	// DefaultCreateTime holds the default value on creation for the "create_time" field.
	DefaultCreateTime func() time.Time
	// DefaultUpdateTime holds the default value on creation for the "update_time" field.
	DefaultUpdateTime func() time.Time
	// UpdateDefaultUpdateTime holds the default value on update for the "update_time" field.
	UpdateDefaultUpdateTime func() time.Time
	// DefaultAdmin holds the default value on creation for the "admin" field.
	DefaultAdmin bool
	// DiscriminatorValidator is a validator for the "discriminator" field. It is called by the builders before save.
	DiscriminatorValidator func(string) error
	// EmailValidator is a validator for the "email" field. It is called by the builders before save.
	EmailValidator func(string) error
	// AvatarHashValidator is a validator for the "avatar_hash" field. It is called by the builders before save.
	AvatarHashValidator func(string) error
	// AvatarURLValidator is a validator for the "avatar_url" field. It is called by the builders before save.
	AvatarURLValidator func(string) error
	// LocaleValidator is a validator for the "locale" field. It is called by the builders before save.
	LocaleValidator func(string) error
	// DefaultBot holds the default value on creation for the "bot" field.
	DefaultBot bool
	// DefaultSystem holds the default value on creation for the "system" field.
	DefaultSystem bool
	// DefaultMfaEnabled holds the default value on creation for the "mfa_enabled" field.
	DefaultMfaEnabled bool
	// DefaultVerified holds the default value on creation for the "verified" field.
	DefaultVerified bool
	// DefaultFlags holds the default value on creation for the "flags" field.
	DefaultFlags uint64
	// DefaultPremiumType holds the default value on creation for the "premium_type" field.
	DefaultPremiumType int
	// DefaultPublicFlags holds the default value on creation for the "public_flags" field.
	DefaultPublicFlags uint64
)
