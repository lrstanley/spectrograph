// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.
//
// DO NOT EDIT, CODE GENERATED BY entc.

package ent

import (
	"github.com/lrstanley/spectrograph/internal/database/ent/guild"
	"github.com/lrstanley/spectrograph/internal/database/ent/guildadminconfig"
	"github.com/lrstanley/spectrograph/internal/database/ent/guildconfig"
	"github.com/lrstanley/spectrograph/internal/database/ent/guildevent"
	"github.com/lrstanley/spectrograph/internal/database/ent/predicate"
	"github.com/lrstanley/spectrograph/internal/database/ent/user"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/entql"
	"entgo.io/ent/schema/field"
)

// schemaGraph holds a representation of ent/schema at runtime.
var schemaGraph = func() *sqlgraph.Schema {
	graph := &sqlgraph.Schema{Nodes: make([]*sqlgraph.Node, 5)}
	graph.Nodes[0] = &sqlgraph.Node{
		NodeSpec: sqlgraph.NodeSpec{
			Table:   guild.Table,
			Columns: guild.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: guild.FieldID,
			},
		},
		Type: "Guild",
		Fields: map[string]*sqlgraph.FieldSpec{
			guild.FieldCreateTime:         {Type: field.TypeTime, Column: guild.FieldCreateTime},
			guild.FieldUpdateTime:         {Type: field.TypeTime, Column: guild.FieldUpdateTime},
			guild.FieldGuildID:            {Type: field.TypeString, Column: guild.FieldGuildID},
			guild.FieldName:               {Type: field.TypeString, Column: guild.FieldName},
			guild.FieldFeatures:           {Type: field.TypeJSON, Column: guild.FieldFeatures},
			guild.FieldIconHash:           {Type: field.TypeString, Column: guild.FieldIconHash},
			guild.FieldIconURL:            {Type: field.TypeString, Column: guild.FieldIconURL},
			guild.FieldJoinedAt:           {Type: field.TypeTime, Column: guild.FieldJoinedAt},
			guild.FieldLarge:              {Type: field.TypeBool, Column: guild.FieldLarge},
			guild.FieldMemberCount:        {Type: field.TypeInt, Column: guild.FieldMemberCount},
			guild.FieldOwnerID:            {Type: field.TypeString, Column: guild.FieldOwnerID},
			guild.FieldPermissions:        {Type: field.TypeUint64, Column: guild.FieldPermissions},
			guild.FieldSystemChannelFlags: {Type: field.TypeString, Column: guild.FieldSystemChannelFlags},
		},
	}
	graph.Nodes[1] = &sqlgraph.Node{
		NodeSpec: sqlgraph.NodeSpec{
			Table:   guildadminconfig.Table,
			Columns: guildadminconfig.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: guildadminconfig.FieldID,
			},
		},
		Type: "GuildAdminConfig",
		Fields: map[string]*sqlgraph.FieldSpec{
			guildadminconfig.FieldCreateTime:         {Type: field.TypeTime, Column: guildadminconfig.FieldCreateTime},
			guildadminconfig.FieldUpdateTime:         {Type: field.TypeTime, Column: guildadminconfig.FieldUpdateTime},
			guildadminconfig.FieldEnabled:            {Type: field.TypeBool, Column: guildadminconfig.FieldEnabled},
			guildadminconfig.FieldDefaultMaxChannels: {Type: field.TypeInt, Column: guildadminconfig.FieldDefaultMaxChannels},
			guildadminconfig.FieldDefaultMaxClones:   {Type: field.TypeInt, Column: guildadminconfig.FieldDefaultMaxClones},
			guildadminconfig.FieldComment:            {Type: field.TypeString, Column: guildadminconfig.FieldComment},
		},
	}
	graph.Nodes[2] = &sqlgraph.Node{
		NodeSpec: sqlgraph.NodeSpec{
			Table:   guildconfig.Table,
			Columns: guildconfig.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: guildconfig.FieldID,
			},
		},
		Type: "GuildConfig",
		Fields: map[string]*sqlgraph.FieldSpec{
			guildconfig.FieldCreateTime:       {Type: field.TypeTime, Column: guildconfig.FieldCreateTime},
			guildconfig.FieldUpdateTime:       {Type: field.TypeTime, Column: guildconfig.FieldUpdateTime},
			guildconfig.FieldEnabled:          {Type: field.TypeBool, Column: guildconfig.FieldEnabled},
			guildconfig.FieldDefaultMaxClones: {Type: field.TypeInt, Column: guildconfig.FieldDefaultMaxClones},
			guildconfig.FieldRegexMatch:       {Type: field.TypeString, Column: guildconfig.FieldRegexMatch},
			guildconfig.FieldContactEmail:     {Type: field.TypeString, Column: guildconfig.FieldContactEmail},
		},
	}
	graph.Nodes[3] = &sqlgraph.Node{
		NodeSpec: sqlgraph.NodeSpec{
			Table:   guildevent.Table,
			Columns: guildevent.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: guildevent.FieldID,
			},
		},
		Type: "GuildEvent",
		Fields: map[string]*sqlgraph.FieldSpec{
			guildevent.FieldCreateTime: {Type: field.TypeTime, Column: guildevent.FieldCreateTime},
			guildevent.FieldUpdateTime: {Type: field.TypeTime, Column: guildevent.FieldUpdateTime},
			guildevent.FieldType:       {Type: field.TypeEnum, Column: guildevent.FieldType},
			guildevent.FieldMessage:    {Type: field.TypeString, Column: guildevent.FieldMessage},
			guildevent.FieldMetadata:   {Type: field.TypeJSON, Column: guildevent.FieldMetadata},
		},
	}
	graph.Nodes[4] = &sqlgraph.Node{
		NodeSpec: sqlgraph.NodeSpec{
			Table:   user.Table,
			Columns: user.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: user.FieldID,
			},
		},
		Type: "User",
		Fields: map[string]*sqlgraph.FieldSpec{
			user.FieldCreateTime:    {Type: field.TypeTime, Column: user.FieldCreateTime},
			user.FieldUpdateTime:    {Type: field.TypeTime, Column: user.FieldUpdateTime},
			user.FieldUserID:        {Type: field.TypeString, Column: user.FieldUserID},
			user.FieldAdmin:         {Type: field.TypeBool, Column: user.FieldAdmin},
			user.FieldBanned:        {Type: field.TypeBool, Column: user.FieldBanned},
			user.FieldBanReason:     {Type: field.TypeString, Column: user.FieldBanReason},
			user.FieldUsername:      {Type: field.TypeString, Column: user.FieldUsername},
			user.FieldDiscriminator: {Type: field.TypeString, Column: user.FieldDiscriminator},
			user.FieldEmail:         {Type: field.TypeString, Column: user.FieldEmail},
			user.FieldAvatarHash:    {Type: field.TypeString, Column: user.FieldAvatarHash},
			user.FieldAvatarURL:     {Type: field.TypeString, Column: user.FieldAvatarURL},
			user.FieldLocale:        {Type: field.TypeString, Column: user.FieldLocale},
			user.FieldBot:           {Type: field.TypeBool, Column: user.FieldBot},
			user.FieldSystem:        {Type: field.TypeBool, Column: user.FieldSystem},
			user.FieldMfaEnabled:    {Type: field.TypeBool, Column: user.FieldMfaEnabled},
			user.FieldVerified:      {Type: field.TypeBool, Column: user.FieldVerified},
			user.FieldFlags:         {Type: field.TypeUint64, Column: user.FieldFlags},
			user.FieldPremiumType:   {Type: field.TypeInt, Column: user.FieldPremiumType},
			user.FieldPublicFlags:   {Type: field.TypeUint64, Column: user.FieldPublicFlags},
		},
	}
	graph.MustAddE(
		"guild_config",
		&sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   guild.GuildConfigTable,
			Columns: []string{guild.GuildConfigColumn},
			Bidi:    false,
		},
		"Guild",
		"GuildConfig",
	)
	graph.MustAddE(
		"guild_admin_config",
		&sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   guild.GuildAdminConfigTable,
			Columns: []string{guild.GuildAdminConfigColumn},
			Bidi:    false,
		},
		"Guild",
		"GuildAdminConfig",
	)
	graph.MustAddE(
		"guild_events",
		&sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   guild.GuildEventsTable,
			Columns: []string{guild.GuildEventsColumn},
			Bidi:    false,
		},
		"Guild",
		"GuildEvent",
	)
	graph.MustAddE(
		"admins",
		&sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   guild.AdminsTable,
			Columns: guild.AdminsPrimaryKey,
			Bidi:    false,
		},
		"Guild",
		"User",
	)
	graph.MustAddE(
		"guild",
		&sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   guildadminconfig.GuildTable,
			Columns: []string{guildadminconfig.GuildColumn},
			Bidi:    false,
		},
		"GuildAdminConfig",
		"Guild",
	)
	graph.MustAddE(
		"guild",
		&sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   guildconfig.GuildTable,
			Columns: []string{guildconfig.GuildColumn},
			Bidi:    false,
		},
		"GuildConfig",
		"Guild",
	)
	graph.MustAddE(
		"guild",
		&sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   guildevent.GuildTable,
			Columns: []string{guildevent.GuildColumn},
			Bidi:    false,
		},
		"GuildEvent",
		"Guild",
	)
	graph.MustAddE(
		"user_guilds",
		&sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   user.UserGuildsTable,
			Columns: user.UserGuildsPrimaryKey,
			Bidi:    false,
		},
		"User",
		"Guild",
	)
	graph.MustAddE(
		"banned_users",
		&sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.BannedUsersTable,
			Columns: []string{user.BannedUsersColumn},
			Bidi:    true,
		},
		"User",
		"User",
	)
	graph.MustAddE(
		"banned_by",
		&sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   user.BannedByTable,
			Columns: []string{user.BannedByColumn},
			Bidi:    false,
		},
		"User",
		"User",
	)
	return graph
}()

// predicateAdder wraps the addPredicate method.
// All update, update-one and query builders implement this interface.
type predicateAdder interface {
	addPredicate(func(s *sql.Selector))
}

// addPredicate implements the predicateAdder interface.
func (gq *GuildQuery) addPredicate(pred func(s *sql.Selector)) {
	gq.predicates = append(gq.predicates, pred)
}

// Filter returns a Filter implementation to apply filters on the GuildQuery builder.
func (gq *GuildQuery) Filter() *GuildFilter {
	return &GuildFilter{config: gq.config, predicateAdder: gq}
}

// addPredicate implements the predicateAdder interface.
func (m *GuildMutation) addPredicate(pred func(s *sql.Selector)) {
	m.predicates = append(m.predicates, pred)
}

// Filter returns an entql.Where implementation to apply filters on the GuildMutation builder.
func (m *GuildMutation) Filter() *GuildFilter {
	return &GuildFilter{config: m.config, predicateAdder: m}
}

// GuildFilter provides a generic filtering capability at runtime for GuildQuery.
type GuildFilter struct {
	predicateAdder
	config
}

// Where applies the entql predicate on the query filter.
func (f *GuildFilter) Where(p entql.P) {
	f.addPredicate(func(s *sql.Selector) {
		if err := schemaGraph.EvalP(schemaGraph.Nodes[0].Type, p, s); err != nil {
			s.AddError(err)
		}
	})
}

// WhereID applies the entql int predicate on the id field.
func (f *GuildFilter) WhereID(p entql.IntP) {
	f.Where(p.Field(guild.FieldID))
}

// WhereCreateTime applies the entql time.Time predicate on the create_time field.
func (f *GuildFilter) WhereCreateTime(p entql.TimeP) {
	f.Where(p.Field(guild.FieldCreateTime))
}

// WhereUpdateTime applies the entql time.Time predicate on the update_time field.
func (f *GuildFilter) WhereUpdateTime(p entql.TimeP) {
	f.Where(p.Field(guild.FieldUpdateTime))
}

// WhereGuildID applies the entql string predicate on the guild_id field.
func (f *GuildFilter) WhereGuildID(p entql.StringP) {
	f.Where(p.Field(guild.FieldGuildID))
}

// WhereName applies the entql string predicate on the name field.
func (f *GuildFilter) WhereName(p entql.StringP) {
	f.Where(p.Field(guild.FieldName))
}

// WhereFeatures applies the entql json.RawMessage predicate on the features field.
func (f *GuildFilter) WhereFeatures(p entql.BytesP) {
	f.Where(p.Field(guild.FieldFeatures))
}

// WhereIconHash applies the entql string predicate on the icon_hash field.
func (f *GuildFilter) WhereIconHash(p entql.StringP) {
	f.Where(p.Field(guild.FieldIconHash))
}

// WhereIconURL applies the entql string predicate on the icon_url field.
func (f *GuildFilter) WhereIconURL(p entql.StringP) {
	f.Where(p.Field(guild.FieldIconURL))
}

// WhereJoinedAt applies the entql time.Time predicate on the joined_at field.
func (f *GuildFilter) WhereJoinedAt(p entql.TimeP) {
	f.Where(p.Field(guild.FieldJoinedAt))
}

// WhereLarge applies the entql bool predicate on the large field.
func (f *GuildFilter) WhereLarge(p entql.BoolP) {
	f.Where(p.Field(guild.FieldLarge))
}

// WhereMemberCount applies the entql int predicate on the member_count field.
func (f *GuildFilter) WhereMemberCount(p entql.IntP) {
	f.Where(p.Field(guild.FieldMemberCount))
}

// WhereOwnerID applies the entql string predicate on the owner_id field.
func (f *GuildFilter) WhereOwnerID(p entql.StringP) {
	f.Where(p.Field(guild.FieldOwnerID))
}

// WherePermissions applies the entql uint64 predicate on the permissions field.
func (f *GuildFilter) WherePermissions(p entql.Uint64P) {
	f.Where(p.Field(guild.FieldPermissions))
}

// WhereSystemChannelFlags applies the entql string predicate on the system_channel_flags field.
func (f *GuildFilter) WhereSystemChannelFlags(p entql.StringP) {
	f.Where(p.Field(guild.FieldSystemChannelFlags))
}

// WhereHasGuildConfig applies a predicate to check if query has an edge guild_config.
func (f *GuildFilter) WhereHasGuildConfig() {
	f.Where(entql.HasEdge("guild_config"))
}

// WhereHasGuildConfigWith applies a predicate to check if query has an edge guild_config with a given conditions (other predicates).
func (f *GuildFilter) WhereHasGuildConfigWith(preds ...predicate.GuildConfig) {
	f.Where(entql.HasEdgeWith("guild_config", sqlgraph.WrapFunc(func(s *sql.Selector) {
		for _, p := range preds {
			p(s)
		}
	})))
}

// WhereHasGuildAdminConfig applies a predicate to check if query has an edge guild_admin_config.
func (f *GuildFilter) WhereHasGuildAdminConfig() {
	f.Where(entql.HasEdge("guild_admin_config"))
}

// WhereHasGuildAdminConfigWith applies a predicate to check if query has an edge guild_admin_config with a given conditions (other predicates).
func (f *GuildFilter) WhereHasGuildAdminConfigWith(preds ...predicate.GuildAdminConfig) {
	f.Where(entql.HasEdgeWith("guild_admin_config", sqlgraph.WrapFunc(func(s *sql.Selector) {
		for _, p := range preds {
			p(s)
		}
	})))
}

// WhereHasGuildEvents applies a predicate to check if query has an edge guild_events.
func (f *GuildFilter) WhereHasGuildEvents() {
	f.Where(entql.HasEdge("guild_events"))
}

// WhereHasGuildEventsWith applies a predicate to check if query has an edge guild_events with a given conditions (other predicates).
func (f *GuildFilter) WhereHasGuildEventsWith(preds ...predicate.GuildEvent) {
	f.Where(entql.HasEdgeWith("guild_events", sqlgraph.WrapFunc(func(s *sql.Selector) {
		for _, p := range preds {
			p(s)
		}
	})))
}

// WhereHasAdmins applies a predicate to check if query has an edge admins.
func (f *GuildFilter) WhereHasAdmins() {
	f.Where(entql.HasEdge("admins"))
}

// WhereHasAdminsWith applies a predicate to check if query has an edge admins with a given conditions (other predicates).
func (f *GuildFilter) WhereHasAdminsWith(preds ...predicate.User) {
	f.Where(entql.HasEdgeWith("admins", sqlgraph.WrapFunc(func(s *sql.Selector) {
		for _, p := range preds {
			p(s)
		}
	})))
}

// addPredicate implements the predicateAdder interface.
func (gacq *GuildAdminConfigQuery) addPredicate(pred func(s *sql.Selector)) {
	gacq.predicates = append(gacq.predicates, pred)
}

// Filter returns a Filter implementation to apply filters on the GuildAdminConfigQuery builder.
func (gacq *GuildAdminConfigQuery) Filter() *GuildAdminConfigFilter {
	return &GuildAdminConfigFilter{config: gacq.config, predicateAdder: gacq}
}

// addPredicate implements the predicateAdder interface.
func (m *GuildAdminConfigMutation) addPredicate(pred func(s *sql.Selector)) {
	m.predicates = append(m.predicates, pred)
}

// Filter returns an entql.Where implementation to apply filters on the GuildAdminConfigMutation builder.
func (m *GuildAdminConfigMutation) Filter() *GuildAdminConfigFilter {
	return &GuildAdminConfigFilter{config: m.config, predicateAdder: m}
}

// GuildAdminConfigFilter provides a generic filtering capability at runtime for GuildAdminConfigQuery.
type GuildAdminConfigFilter struct {
	predicateAdder
	config
}

// Where applies the entql predicate on the query filter.
func (f *GuildAdminConfigFilter) Where(p entql.P) {
	f.addPredicate(func(s *sql.Selector) {
		if err := schemaGraph.EvalP(schemaGraph.Nodes[1].Type, p, s); err != nil {
			s.AddError(err)
		}
	})
}

// WhereID applies the entql int predicate on the id field.
func (f *GuildAdminConfigFilter) WhereID(p entql.IntP) {
	f.Where(p.Field(guildadminconfig.FieldID))
}

// WhereCreateTime applies the entql time.Time predicate on the create_time field.
func (f *GuildAdminConfigFilter) WhereCreateTime(p entql.TimeP) {
	f.Where(p.Field(guildadminconfig.FieldCreateTime))
}

// WhereUpdateTime applies the entql time.Time predicate on the update_time field.
func (f *GuildAdminConfigFilter) WhereUpdateTime(p entql.TimeP) {
	f.Where(p.Field(guildadminconfig.FieldUpdateTime))
}

// WhereEnabled applies the entql bool predicate on the enabled field.
func (f *GuildAdminConfigFilter) WhereEnabled(p entql.BoolP) {
	f.Where(p.Field(guildadminconfig.FieldEnabled))
}

// WhereDefaultMaxChannels applies the entql int predicate on the default_max_channels field.
func (f *GuildAdminConfigFilter) WhereDefaultMaxChannels(p entql.IntP) {
	f.Where(p.Field(guildadminconfig.FieldDefaultMaxChannels))
}

// WhereDefaultMaxClones applies the entql int predicate on the default_max_clones field.
func (f *GuildAdminConfigFilter) WhereDefaultMaxClones(p entql.IntP) {
	f.Where(p.Field(guildadminconfig.FieldDefaultMaxClones))
}

// WhereComment applies the entql string predicate on the comment field.
func (f *GuildAdminConfigFilter) WhereComment(p entql.StringP) {
	f.Where(p.Field(guildadminconfig.FieldComment))
}

// WhereHasGuild applies a predicate to check if query has an edge guild.
func (f *GuildAdminConfigFilter) WhereHasGuild() {
	f.Where(entql.HasEdge("guild"))
}

// WhereHasGuildWith applies a predicate to check if query has an edge guild with a given conditions (other predicates).
func (f *GuildAdminConfigFilter) WhereHasGuildWith(preds ...predicate.Guild) {
	f.Where(entql.HasEdgeWith("guild", sqlgraph.WrapFunc(func(s *sql.Selector) {
		for _, p := range preds {
			p(s)
		}
	})))
}

// addPredicate implements the predicateAdder interface.
func (gcq *GuildConfigQuery) addPredicate(pred func(s *sql.Selector)) {
	gcq.predicates = append(gcq.predicates, pred)
}

// Filter returns a Filter implementation to apply filters on the GuildConfigQuery builder.
func (gcq *GuildConfigQuery) Filter() *GuildConfigFilter {
	return &GuildConfigFilter{config: gcq.config, predicateAdder: gcq}
}

// addPredicate implements the predicateAdder interface.
func (m *GuildConfigMutation) addPredicate(pred func(s *sql.Selector)) {
	m.predicates = append(m.predicates, pred)
}

// Filter returns an entql.Where implementation to apply filters on the GuildConfigMutation builder.
func (m *GuildConfigMutation) Filter() *GuildConfigFilter {
	return &GuildConfigFilter{config: m.config, predicateAdder: m}
}

// GuildConfigFilter provides a generic filtering capability at runtime for GuildConfigQuery.
type GuildConfigFilter struct {
	predicateAdder
	config
}

// Where applies the entql predicate on the query filter.
func (f *GuildConfigFilter) Where(p entql.P) {
	f.addPredicate(func(s *sql.Selector) {
		if err := schemaGraph.EvalP(schemaGraph.Nodes[2].Type, p, s); err != nil {
			s.AddError(err)
		}
	})
}

// WhereID applies the entql int predicate on the id field.
func (f *GuildConfigFilter) WhereID(p entql.IntP) {
	f.Where(p.Field(guildconfig.FieldID))
}

// WhereCreateTime applies the entql time.Time predicate on the create_time field.
func (f *GuildConfigFilter) WhereCreateTime(p entql.TimeP) {
	f.Where(p.Field(guildconfig.FieldCreateTime))
}

// WhereUpdateTime applies the entql time.Time predicate on the update_time field.
func (f *GuildConfigFilter) WhereUpdateTime(p entql.TimeP) {
	f.Where(p.Field(guildconfig.FieldUpdateTime))
}

// WhereEnabled applies the entql bool predicate on the enabled field.
func (f *GuildConfigFilter) WhereEnabled(p entql.BoolP) {
	f.Where(p.Field(guildconfig.FieldEnabled))
}

// WhereDefaultMaxClones applies the entql int predicate on the default_max_clones field.
func (f *GuildConfigFilter) WhereDefaultMaxClones(p entql.IntP) {
	f.Where(p.Field(guildconfig.FieldDefaultMaxClones))
}

// WhereRegexMatch applies the entql string predicate on the regex_match field.
func (f *GuildConfigFilter) WhereRegexMatch(p entql.StringP) {
	f.Where(p.Field(guildconfig.FieldRegexMatch))
}

// WhereContactEmail applies the entql string predicate on the contact_email field.
func (f *GuildConfigFilter) WhereContactEmail(p entql.StringP) {
	f.Where(p.Field(guildconfig.FieldContactEmail))
}

// WhereHasGuild applies a predicate to check if query has an edge guild.
func (f *GuildConfigFilter) WhereHasGuild() {
	f.Where(entql.HasEdge("guild"))
}

// WhereHasGuildWith applies a predicate to check if query has an edge guild with a given conditions (other predicates).
func (f *GuildConfigFilter) WhereHasGuildWith(preds ...predicate.Guild) {
	f.Where(entql.HasEdgeWith("guild", sqlgraph.WrapFunc(func(s *sql.Selector) {
		for _, p := range preds {
			p(s)
		}
	})))
}

// addPredicate implements the predicateAdder interface.
func (geq *GuildEventQuery) addPredicate(pred func(s *sql.Selector)) {
	geq.predicates = append(geq.predicates, pred)
}

// Filter returns a Filter implementation to apply filters on the GuildEventQuery builder.
func (geq *GuildEventQuery) Filter() *GuildEventFilter {
	return &GuildEventFilter{config: geq.config, predicateAdder: geq}
}

// addPredicate implements the predicateAdder interface.
func (m *GuildEventMutation) addPredicate(pred func(s *sql.Selector)) {
	m.predicates = append(m.predicates, pred)
}

// Filter returns an entql.Where implementation to apply filters on the GuildEventMutation builder.
func (m *GuildEventMutation) Filter() *GuildEventFilter {
	return &GuildEventFilter{config: m.config, predicateAdder: m}
}

// GuildEventFilter provides a generic filtering capability at runtime for GuildEventQuery.
type GuildEventFilter struct {
	predicateAdder
	config
}

// Where applies the entql predicate on the query filter.
func (f *GuildEventFilter) Where(p entql.P) {
	f.addPredicate(func(s *sql.Selector) {
		if err := schemaGraph.EvalP(schemaGraph.Nodes[3].Type, p, s); err != nil {
			s.AddError(err)
		}
	})
}

// WhereID applies the entql int predicate on the id field.
func (f *GuildEventFilter) WhereID(p entql.IntP) {
	f.Where(p.Field(guildevent.FieldID))
}

// WhereCreateTime applies the entql time.Time predicate on the create_time field.
func (f *GuildEventFilter) WhereCreateTime(p entql.TimeP) {
	f.Where(p.Field(guildevent.FieldCreateTime))
}

// WhereUpdateTime applies the entql time.Time predicate on the update_time field.
func (f *GuildEventFilter) WhereUpdateTime(p entql.TimeP) {
	f.Where(p.Field(guildevent.FieldUpdateTime))
}

// WhereType applies the entql string predicate on the type field.
func (f *GuildEventFilter) WhereType(p entql.StringP) {
	f.Where(p.Field(guildevent.FieldType))
}

// WhereMessage applies the entql string predicate on the message field.
func (f *GuildEventFilter) WhereMessage(p entql.StringP) {
	f.Where(p.Field(guildevent.FieldMessage))
}

// WhereMetadata applies the entql json.RawMessage predicate on the metadata field.
func (f *GuildEventFilter) WhereMetadata(p entql.BytesP) {
	f.Where(p.Field(guildevent.FieldMetadata))
}

// WhereHasGuild applies a predicate to check if query has an edge guild.
func (f *GuildEventFilter) WhereHasGuild() {
	f.Where(entql.HasEdge("guild"))
}

// WhereHasGuildWith applies a predicate to check if query has an edge guild with a given conditions (other predicates).
func (f *GuildEventFilter) WhereHasGuildWith(preds ...predicate.Guild) {
	f.Where(entql.HasEdgeWith("guild", sqlgraph.WrapFunc(func(s *sql.Selector) {
		for _, p := range preds {
			p(s)
		}
	})))
}

// addPredicate implements the predicateAdder interface.
func (uq *UserQuery) addPredicate(pred func(s *sql.Selector)) {
	uq.predicates = append(uq.predicates, pred)
}

// Filter returns a Filter implementation to apply filters on the UserQuery builder.
func (uq *UserQuery) Filter() *UserFilter {
	return &UserFilter{config: uq.config, predicateAdder: uq}
}

// addPredicate implements the predicateAdder interface.
func (m *UserMutation) addPredicate(pred func(s *sql.Selector)) {
	m.predicates = append(m.predicates, pred)
}

// Filter returns an entql.Where implementation to apply filters on the UserMutation builder.
func (m *UserMutation) Filter() *UserFilter {
	return &UserFilter{config: m.config, predicateAdder: m}
}

// UserFilter provides a generic filtering capability at runtime for UserQuery.
type UserFilter struct {
	predicateAdder
	config
}

// Where applies the entql predicate on the query filter.
func (f *UserFilter) Where(p entql.P) {
	f.addPredicate(func(s *sql.Selector) {
		if err := schemaGraph.EvalP(schemaGraph.Nodes[4].Type, p, s); err != nil {
			s.AddError(err)
		}
	})
}

// WhereID applies the entql int predicate on the id field.
func (f *UserFilter) WhereID(p entql.IntP) {
	f.Where(p.Field(user.FieldID))
}

// WhereCreateTime applies the entql time.Time predicate on the create_time field.
func (f *UserFilter) WhereCreateTime(p entql.TimeP) {
	f.Where(p.Field(user.FieldCreateTime))
}

// WhereUpdateTime applies the entql time.Time predicate on the update_time field.
func (f *UserFilter) WhereUpdateTime(p entql.TimeP) {
	f.Where(p.Field(user.FieldUpdateTime))
}

// WhereUserID applies the entql string predicate on the user_id field.
func (f *UserFilter) WhereUserID(p entql.StringP) {
	f.Where(p.Field(user.FieldUserID))
}

// WhereAdmin applies the entql bool predicate on the admin field.
func (f *UserFilter) WhereAdmin(p entql.BoolP) {
	f.Where(p.Field(user.FieldAdmin))
}

// WhereBanned applies the entql bool predicate on the banned field.
func (f *UserFilter) WhereBanned(p entql.BoolP) {
	f.Where(p.Field(user.FieldBanned))
}

// WhereBanReason applies the entql string predicate on the ban_reason field.
func (f *UserFilter) WhereBanReason(p entql.StringP) {
	f.Where(p.Field(user.FieldBanReason))
}

// WhereUsername applies the entql string predicate on the username field.
func (f *UserFilter) WhereUsername(p entql.StringP) {
	f.Where(p.Field(user.FieldUsername))
}

// WhereDiscriminator applies the entql string predicate on the discriminator field.
func (f *UserFilter) WhereDiscriminator(p entql.StringP) {
	f.Where(p.Field(user.FieldDiscriminator))
}

// WhereEmail applies the entql string predicate on the email field.
func (f *UserFilter) WhereEmail(p entql.StringP) {
	f.Where(p.Field(user.FieldEmail))
}

// WhereAvatarHash applies the entql string predicate on the avatar_hash field.
func (f *UserFilter) WhereAvatarHash(p entql.StringP) {
	f.Where(p.Field(user.FieldAvatarHash))
}

// WhereAvatarURL applies the entql string predicate on the avatar_url field.
func (f *UserFilter) WhereAvatarURL(p entql.StringP) {
	f.Where(p.Field(user.FieldAvatarURL))
}

// WhereLocale applies the entql string predicate on the locale field.
func (f *UserFilter) WhereLocale(p entql.StringP) {
	f.Where(p.Field(user.FieldLocale))
}

// WhereBot applies the entql bool predicate on the bot field.
func (f *UserFilter) WhereBot(p entql.BoolP) {
	f.Where(p.Field(user.FieldBot))
}

// WhereSystem applies the entql bool predicate on the system field.
func (f *UserFilter) WhereSystem(p entql.BoolP) {
	f.Where(p.Field(user.FieldSystem))
}

// WhereMfaEnabled applies the entql bool predicate on the mfa_enabled field.
func (f *UserFilter) WhereMfaEnabled(p entql.BoolP) {
	f.Where(p.Field(user.FieldMfaEnabled))
}

// WhereVerified applies the entql bool predicate on the verified field.
func (f *UserFilter) WhereVerified(p entql.BoolP) {
	f.Where(p.Field(user.FieldVerified))
}

// WhereFlags applies the entql uint64 predicate on the flags field.
func (f *UserFilter) WhereFlags(p entql.Uint64P) {
	f.Where(p.Field(user.FieldFlags))
}

// WherePremiumType applies the entql int predicate on the premium_type field.
func (f *UserFilter) WherePremiumType(p entql.IntP) {
	f.Where(p.Field(user.FieldPremiumType))
}

// WherePublicFlags applies the entql uint64 predicate on the public_flags field.
func (f *UserFilter) WherePublicFlags(p entql.Uint64P) {
	f.Where(p.Field(user.FieldPublicFlags))
}

// WhereHasUserGuilds applies a predicate to check if query has an edge user_guilds.
func (f *UserFilter) WhereHasUserGuilds() {
	f.Where(entql.HasEdge("user_guilds"))
}

// WhereHasUserGuildsWith applies a predicate to check if query has an edge user_guilds with a given conditions (other predicates).
func (f *UserFilter) WhereHasUserGuildsWith(preds ...predicate.Guild) {
	f.Where(entql.HasEdgeWith("user_guilds", sqlgraph.WrapFunc(func(s *sql.Selector) {
		for _, p := range preds {
			p(s)
		}
	})))
}

// WhereHasBannedUsers applies a predicate to check if query has an edge banned_users.
func (f *UserFilter) WhereHasBannedUsers() {
	f.Where(entql.HasEdge("banned_users"))
}

// WhereHasBannedUsersWith applies a predicate to check if query has an edge banned_users with a given conditions (other predicates).
func (f *UserFilter) WhereHasBannedUsersWith(preds ...predicate.User) {
	f.Where(entql.HasEdgeWith("banned_users", sqlgraph.WrapFunc(func(s *sql.Selector) {
		for _, p := range preds {
			p(s)
		}
	})))
}

// WhereHasBannedBy applies a predicate to check if query has an edge banned_by.
func (f *UserFilter) WhereHasBannedBy() {
	f.Where(entql.HasEdge("banned_by"))
}

// WhereHasBannedByWith applies a predicate to check if query has an edge banned_by with a given conditions (other predicates).
func (f *UserFilter) WhereHasBannedByWith(preds ...predicate.User) {
	f.Where(entql.HasEdgeWith("banned_by", sqlgraph.WrapFunc(func(s *sql.Selector) {
		for _, p := range preds {
			p(s)
		}
	})))
}