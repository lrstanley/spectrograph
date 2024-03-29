// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.
//
// DO NOT EDIT, CODE GENERATED BY entc.

package ent

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/lrstanley/spectrograph/internal/database/ent/guild"
	"github.com/lrstanley/spectrograph/internal/database/ent/guildevent"
)

// GuildEvent is the model entity for the GuildEvent schema.
type GuildEvent struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// The type of event that occurred.
	Type guildevent.Type `json:"type,omitempty"`
	// The message associated with the event.
	Message string `json:"message,omitempty"`
	// Additional metadata associated with the event.
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the GuildEventQuery when eager-loading is set.
	Edges              GuildEventEdges `json:"edges"`
	guild_guild_events *int
	selectValues       sql.SelectValues
}

// GuildEventEdges holds the relations/edges for other nodes in the graph.
type GuildEventEdges struct {
	// The guild these events belong to.
	Guild *Guild `json:"guild,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
	// totalCount holds the count of the edges above.
	totalCount [1]map[string]int
}

// GuildOrErr returns the Guild value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e GuildEventEdges) GuildOrErr() (*Guild, error) {
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
func (*GuildEvent) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case guildevent.FieldMetadata:
			values[i] = new([]byte)
		case guildevent.FieldID:
			values[i] = new(sql.NullInt64)
		case guildevent.FieldType, guildevent.FieldMessage:
			values[i] = new(sql.NullString)
		case guildevent.FieldCreateTime, guildevent.FieldUpdateTime:
			values[i] = new(sql.NullTime)
		case guildevent.ForeignKeys[0]: // guild_guild_events
			values[i] = new(sql.NullInt64)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the GuildEvent fields.
func (ge *GuildEvent) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case guildevent.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			ge.ID = int(value.Int64)
		case guildevent.FieldCreateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field create_time", values[i])
			} else if value.Valid {
				ge.CreateTime = value.Time
			}
		case guildevent.FieldUpdateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field update_time", values[i])
			} else if value.Valid {
				ge.UpdateTime = value.Time
			}
		case guildevent.FieldType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field type", values[i])
			} else if value.Valid {
				ge.Type = guildevent.Type(value.String)
			}
		case guildevent.FieldMessage:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field message", values[i])
			} else if value.Valid {
				ge.Message = value.String
			}
		case guildevent.FieldMetadata:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field metadata", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &ge.Metadata); err != nil {
					return fmt.Errorf("unmarshal field metadata: %w", err)
				}
			}
		case guildevent.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field guild_guild_events", value)
			} else if value.Valid {
				ge.guild_guild_events = new(int)
				*ge.guild_guild_events = int(value.Int64)
			}
		default:
			ge.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the GuildEvent.
// This includes values selected through modifiers, order, etc.
func (ge *GuildEvent) Value(name string) (ent.Value, error) {
	return ge.selectValues.Get(name)
}

// QueryGuild queries the "guild" edge of the GuildEvent entity.
func (ge *GuildEvent) QueryGuild() *GuildQuery {
	return NewGuildEventClient(ge.config).QueryGuild(ge)
}

// Update returns a builder for updating this GuildEvent.
// Note that you need to call GuildEvent.Unwrap() before calling this method if this GuildEvent
// was returned from a transaction, and the transaction was committed or rolled back.
func (ge *GuildEvent) Update() *GuildEventUpdateOne {
	return NewGuildEventClient(ge.config).UpdateOne(ge)
}

// Unwrap unwraps the GuildEvent entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (ge *GuildEvent) Unwrap() *GuildEvent {
	_tx, ok := ge.config.driver.(*txDriver)
	if !ok {
		panic("ent: GuildEvent is not a transactional entity")
	}
	ge.config.driver = _tx.drv
	return ge
}

// String implements the fmt.Stringer.
func (ge *GuildEvent) String() string {
	var builder strings.Builder
	builder.WriteString("GuildEvent(")
	builder.WriteString(fmt.Sprintf("id=%v, ", ge.ID))
	builder.WriteString("create_time=")
	builder.WriteString(ge.CreateTime.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("update_time=")
	builder.WriteString(ge.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("type=")
	builder.WriteString(fmt.Sprintf("%v", ge.Type))
	builder.WriteString(", ")
	builder.WriteString("message=")
	builder.WriteString(ge.Message)
	builder.WriteString(", ")
	builder.WriteString("metadata=")
	builder.WriteString(fmt.Sprintf("%v", ge.Metadata))
	builder.WriteByte(')')
	return builder.String()
}

// GuildEvents is a parsable slice of GuildEvent.
type GuildEvents []*GuildEvent
