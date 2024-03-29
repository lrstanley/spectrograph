// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.
//
// DO NOT EDIT, CODE GENERATED BY entc.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/lrstanley/spectrograph/internal/database/ent/guild"
	"github.com/lrstanley/spectrograph/internal/database/ent/guildadminconfig"
)

// GuildAdminConfig is the model entity for the GuildAdminConfig schema.
type GuildAdminConfig struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// True if the guild should be monitored/acted upon (overrides user-defined settings).
	Enabled bool `json:"enabled,omitempty"`
	// Default max channels for the guild (overrides user-defined settings).
	DefaultMaxChannels int `json:"default_max_channels,omitempty"`
	// Default max clones for the guild (overrides user-defined settings).
	DefaultMaxClones int `json:"default_max_clones,omitempty"`
	// Admin comment for the guild.
	Comment string `json:"comment,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the GuildAdminConfigQuery when eager-loading is set.
	Edges                    GuildAdminConfigEdges `json:"edges"`
	guild_guild_admin_config *int
	selectValues             sql.SelectValues
}

// GuildAdminConfigEdges holds the relations/edges for other nodes in the graph.
type GuildAdminConfigEdges struct {
	// The guild these settings belong to.
	Guild *Guild `json:"guild,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
	// totalCount holds the count of the edges above.
	totalCount [1]map[string]int
}

// GuildOrErr returns the Guild value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e GuildAdminConfigEdges) GuildOrErr() (*Guild, error) {
	if e.loadedTypes[0] {
		if e.Guild == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: guild.Label}
		}
		return e.Guild, nil
	}
	return nil, &NotLoadedError{edge: "guild"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*GuildAdminConfig) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case guildadminconfig.FieldEnabled:
			values[i] = new(sql.NullBool)
		case guildadminconfig.FieldID, guildadminconfig.FieldDefaultMaxChannels, guildadminconfig.FieldDefaultMaxClones:
			values[i] = new(sql.NullInt64)
		case guildadminconfig.FieldComment:
			values[i] = new(sql.NullString)
		case guildadminconfig.FieldCreateTime, guildadminconfig.FieldUpdateTime:
			values[i] = new(sql.NullTime)
		case guildadminconfig.ForeignKeys[0]: // guild_guild_admin_config
			values[i] = new(sql.NullInt64)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the GuildAdminConfig fields.
func (gac *GuildAdminConfig) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case guildadminconfig.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			gac.ID = int(value.Int64)
		case guildadminconfig.FieldCreateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field create_time", values[i])
			} else if value.Valid {
				gac.CreateTime = value.Time
			}
		case guildadminconfig.FieldUpdateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field update_time", values[i])
			} else if value.Valid {
				gac.UpdateTime = value.Time
			}
		case guildadminconfig.FieldEnabled:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field enabled", values[i])
			} else if value.Valid {
				gac.Enabled = value.Bool
			}
		case guildadminconfig.FieldDefaultMaxChannels:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field default_max_channels", values[i])
			} else if value.Valid {
				gac.DefaultMaxChannels = int(value.Int64)
			}
		case guildadminconfig.FieldDefaultMaxClones:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field default_max_clones", values[i])
			} else if value.Valid {
				gac.DefaultMaxClones = int(value.Int64)
			}
		case guildadminconfig.FieldComment:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field comment", values[i])
			} else if value.Valid {
				gac.Comment = value.String
			}
		case guildadminconfig.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field guild_guild_admin_config", value)
			} else if value.Valid {
				gac.guild_guild_admin_config = new(int)
				*gac.guild_guild_admin_config = int(value.Int64)
			}
		default:
			gac.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the GuildAdminConfig.
// This includes values selected through modifiers, order, etc.
func (gac *GuildAdminConfig) Value(name string) (ent.Value, error) {
	return gac.selectValues.Get(name)
}

// QueryGuild queries the "guild" edge of the GuildAdminConfig entity.
func (gac *GuildAdminConfig) QueryGuild() *GuildQuery {
	return NewGuildAdminConfigClient(gac.config).QueryGuild(gac)
}

// Update returns a builder for updating this GuildAdminConfig.
// Note that you need to call GuildAdminConfig.Unwrap() before calling this method if this GuildAdminConfig
// was returned from a transaction, and the transaction was committed or rolled back.
func (gac *GuildAdminConfig) Update() *GuildAdminConfigUpdateOne {
	return NewGuildAdminConfigClient(gac.config).UpdateOne(gac)
}

// Unwrap unwraps the GuildAdminConfig entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (gac *GuildAdminConfig) Unwrap() *GuildAdminConfig {
	_tx, ok := gac.config.driver.(*txDriver)
	if !ok {
		panic("ent: GuildAdminConfig is not a transactional entity")
	}
	gac.config.driver = _tx.drv
	return gac
}

// String implements the fmt.Stringer.
func (gac *GuildAdminConfig) String() string {
	var builder strings.Builder
	builder.WriteString("GuildAdminConfig(")
	builder.WriteString(fmt.Sprintf("id=%v, ", gac.ID))
	builder.WriteString("create_time=")
	builder.WriteString(gac.CreateTime.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("update_time=")
	builder.WriteString(gac.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("enabled=")
	builder.WriteString(fmt.Sprintf("%v", gac.Enabled))
	builder.WriteString(", ")
	builder.WriteString("default_max_channels=")
	builder.WriteString(fmt.Sprintf("%v", gac.DefaultMaxChannels))
	builder.WriteString(", ")
	builder.WriteString("default_max_clones=")
	builder.WriteString(fmt.Sprintf("%v", gac.DefaultMaxClones))
	builder.WriteString(", ")
	builder.WriteString("comment=")
	builder.WriteString(gac.Comment)
	builder.WriteByte(')')
	return builder.String()
}

// GuildAdminConfigs is a parsable slice of GuildAdminConfig.
type GuildAdminConfigs []*GuildAdminConfig
