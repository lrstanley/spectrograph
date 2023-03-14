// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.
//
// DO NOT EDIT, CODE GENERATED BY entc.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/lrstanley/spectrograph/internal/ent/guildadminconfig"
	"github.com/lrstanley/spectrograph/internal/ent/predicate"
)

// GuildAdminConfigUpdate is the builder for updating GuildAdminConfig entities.
type GuildAdminConfigUpdate struct {
	config
	hooks    []Hook
	mutation *GuildAdminConfigMutation
}

// Where appends a list predicates to the GuildAdminConfigUpdate builder.
func (gacu *GuildAdminConfigUpdate) Where(ps ...predicate.GuildAdminConfig) *GuildAdminConfigUpdate {
	gacu.mutation.Where(ps...)
	return gacu
}

// SetUpdateTime sets the "update_time" field.
func (gacu *GuildAdminConfigUpdate) SetUpdateTime(t time.Time) *GuildAdminConfigUpdate {
	gacu.mutation.SetUpdateTime(t)
	return gacu
}

// SetEnabled sets the "enabled" field.
func (gacu *GuildAdminConfigUpdate) SetEnabled(b bool) *GuildAdminConfigUpdate {
	gacu.mutation.SetEnabled(b)
	return gacu
}

// SetNillableEnabled sets the "enabled" field if the given value is not nil.
func (gacu *GuildAdminConfigUpdate) SetNillableEnabled(b *bool) *GuildAdminConfigUpdate {
	if b != nil {
		gacu.SetEnabled(*b)
	}
	return gacu
}

// ClearEnabled clears the value of the "enabled" field.
func (gacu *GuildAdminConfigUpdate) ClearEnabled() *GuildAdminConfigUpdate {
	gacu.mutation.ClearEnabled()
	return gacu
}

// SetDefaultMaxChannels sets the "default_max_channels" field.
func (gacu *GuildAdminConfigUpdate) SetDefaultMaxChannels(i int) *GuildAdminConfigUpdate {
	gacu.mutation.ResetDefaultMaxChannels()
	gacu.mutation.SetDefaultMaxChannels(i)
	return gacu
}

// SetNillableDefaultMaxChannels sets the "default_max_channels" field if the given value is not nil.
func (gacu *GuildAdminConfigUpdate) SetNillableDefaultMaxChannels(i *int) *GuildAdminConfigUpdate {
	if i != nil {
		gacu.SetDefaultMaxChannels(*i)
	}
	return gacu
}

// AddDefaultMaxChannels adds i to the "default_max_channels" field.
func (gacu *GuildAdminConfigUpdate) AddDefaultMaxChannels(i int) *GuildAdminConfigUpdate {
	gacu.mutation.AddDefaultMaxChannels(i)
	return gacu
}

// ClearDefaultMaxChannels clears the value of the "default_max_channels" field.
func (gacu *GuildAdminConfigUpdate) ClearDefaultMaxChannels() *GuildAdminConfigUpdate {
	gacu.mutation.ClearDefaultMaxChannels()
	return gacu
}

// SetDefaultMaxClones sets the "default_max_clones" field.
func (gacu *GuildAdminConfigUpdate) SetDefaultMaxClones(i int) *GuildAdminConfigUpdate {
	gacu.mutation.ResetDefaultMaxClones()
	gacu.mutation.SetDefaultMaxClones(i)
	return gacu
}

// SetNillableDefaultMaxClones sets the "default_max_clones" field if the given value is not nil.
func (gacu *GuildAdminConfigUpdate) SetNillableDefaultMaxClones(i *int) *GuildAdminConfigUpdate {
	if i != nil {
		gacu.SetDefaultMaxClones(*i)
	}
	return gacu
}

// AddDefaultMaxClones adds i to the "default_max_clones" field.
func (gacu *GuildAdminConfigUpdate) AddDefaultMaxClones(i int) *GuildAdminConfigUpdate {
	gacu.mutation.AddDefaultMaxClones(i)
	return gacu
}

// ClearDefaultMaxClones clears the value of the "default_max_clones" field.
func (gacu *GuildAdminConfigUpdate) ClearDefaultMaxClones() *GuildAdminConfigUpdate {
	gacu.mutation.ClearDefaultMaxClones()
	return gacu
}

// SetComment sets the "comment" field.
func (gacu *GuildAdminConfigUpdate) SetComment(s string) *GuildAdminConfigUpdate {
	gacu.mutation.SetComment(s)
	return gacu
}

// SetNillableComment sets the "comment" field if the given value is not nil.
func (gacu *GuildAdminConfigUpdate) SetNillableComment(s *string) *GuildAdminConfigUpdate {
	if s != nil {
		gacu.SetComment(*s)
	}
	return gacu
}

// ClearComment clears the value of the "comment" field.
func (gacu *GuildAdminConfigUpdate) ClearComment() *GuildAdminConfigUpdate {
	gacu.mutation.ClearComment()
	return gacu
}

// Mutation returns the GuildAdminConfigMutation object of the builder.
func (gacu *GuildAdminConfigUpdate) Mutation() *GuildAdminConfigMutation {
	return gacu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (gacu *GuildAdminConfigUpdate) Save(ctx context.Context) (int, error) {
	if err := gacu.defaults(); err != nil {
		return 0, err
	}
	return withHooks[int, GuildAdminConfigMutation](ctx, gacu.sqlSave, gacu.mutation, gacu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (gacu *GuildAdminConfigUpdate) SaveX(ctx context.Context) int {
	affected, err := gacu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (gacu *GuildAdminConfigUpdate) Exec(ctx context.Context) error {
	_, err := gacu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (gacu *GuildAdminConfigUpdate) ExecX(ctx context.Context) {
	if err := gacu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (gacu *GuildAdminConfigUpdate) defaults() error {
	if _, ok := gacu.mutation.UpdateTime(); !ok {
		if guildadminconfig.UpdateDefaultUpdateTime == nil {
			return fmt.Errorf("ent: uninitialized guildadminconfig.UpdateDefaultUpdateTime (forgotten import ent/runtime?)")
		}
		v := guildadminconfig.UpdateDefaultUpdateTime()
		gacu.mutation.SetUpdateTime(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (gacu *GuildAdminConfigUpdate) check() error {
	if v, ok := gacu.mutation.DefaultMaxChannels(); ok {
		if err := guildadminconfig.DefaultMaxChannelsValidator(v); err != nil {
			return &ValidationError{Name: "default_max_channels", err: fmt.Errorf(`ent: validator failed for field "GuildAdminConfig.default_max_channels": %w`, err)}
		}
	}
	if v, ok := gacu.mutation.DefaultMaxClones(); ok {
		if err := guildadminconfig.DefaultMaxClonesValidator(v); err != nil {
			return &ValidationError{Name: "default_max_clones", err: fmt.Errorf(`ent: validator failed for field "GuildAdminConfig.default_max_clones": %w`, err)}
		}
	}
	if _, ok := gacu.mutation.GuildID(); gacu.mutation.GuildCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "GuildAdminConfig.guild"`)
	}
	return nil
}

func (gacu *GuildAdminConfigUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := gacu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(guildadminconfig.Table, guildadminconfig.Columns, sqlgraph.NewFieldSpec(guildadminconfig.FieldID, field.TypeInt))
	if ps := gacu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := gacu.mutation.UpdateTime(); ok {
		_spec.SetField(guildadminconfig.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := gacu.mutation.Enabled(); ok {
		_spec.SetField(guildadminconfig.FieldEnabled, field.TypeBool, value)
	}
	if gacu.mutation.EnabledCleared() {
		_spec.ClearField(guildadminconfig.FieldEnabled, field.TypeBool)
	}
	if value, ok := gacu.mutation.DefaultMaxChannels(); ok {
		_spec.SetField(guildadminconfig.FieldDefaultMaxChannels, field.TypeInt, value)
	}
	if value, ok := gacu.mutation.AddedDefaultMaxChannels(); ok {
		_spec.AddField(guildadminconfig.FieldDefaultMaxChannels, field.TypeInt, value)
	}
	if gacu.mutation.DefaultMaxChannelsCleared() {
		_spec.ClearField(guildadminconfig.FieldDefaultMaxChannels, field.TypeInt)
	}
	if value, ok := gacu.mutation.DefaultMaxClones(); ok {
		_spec.SetField(guildadminconfig.FieldDefaultMaxClones, field.TypeInt, value)
	}
	if value, ok := gacu.mutation.AddedDefaultMaxClones(); ok {
		_spec.AddField(guildadminconfig.FieldDefaultMaxClones, field.TypeInt, value)
	}
	if gacu.mutation.DefaultMaxClonesCleared() {
		_spec.ClearField(guildadminconfig.FieldDefaultMaxClones, field.TypeInt)
	}
	if value, ok := gacu.mutation.Comment(); ok {
		_spec.SetField(guildadminconfig.FieldComment, field.TypeString, value)
	}
	if gacu.mutation.CommentCleared() {
		_spec.ClearField(guildadminconfig.FieldComment, field.TypeString)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, gacu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{guildadminconfig.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	gacu.mutation.done = true
	return n, nil
}

// GuildAdminConfigUpdateOne is the builder for updating a single GuildAdminConfig entity.
type GuildAdminConfigUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *GuildAdminConfigMutation
}

// SetUpdateTime sets the "update_time" field.
func (gacuo *GuildAdminConfigUpdateOne) SetUpdateTime(t time.Time) *GuildAdminConfigUpdateOne {
	gacuo.mutation.SetUpdateTime(t)
	return gacuo
}

// SetEnabled sets the "enabled" field.
func (gacuo *GuildAdminConfigUpdateOne) SetEnabled(b bool) *GuildAdminConfigUpdateOne {
	gacuo.mutation.SetEnabled(b)
	return gacuo
}

// SetNillableEnabled sets the "enabled" field if the given value is not nil.
func (gacuo *GuildAdminConfigUpdateOne) SetNillableEnabled(b *bool) *GuildAdminConfigUpdateOne {
	if b != nil {
		gacuo.SetEnabled(*b)
	}
	return gacuo
}

// ClearEnabled clears the value of the "enabled" field.
func (gacuo *GuildAdminConfigUpdateOne) ClearEnabled() *GuildAdminConfigUpdateOne {
	gacuo.mutation.ClearEnabled()
	return gacuo
}

// SetDefaultMaxChannels sets the "default_max_channels" field.
func (gacuo *GuildAdminConfigUpdateOne) SetDefaultMaxChannels(i int) *GuildAdminConfigUpdateOne {
	gacuo.mutation.ResetDefaultMaxChannels()
	gacuo.mutation.SetDefaultMaxChannels(i)
	return gacuo
}

// SetNillableDefaultMaxChannels sets the "default_max_channels" field if the given value is not nil.
func (gacuo *GuildAdminConfigUpdateOne) SetNillableDefaultMaxChannels(i *int) *GuildAdminConfigUpdateOne {
	if i != nil {
		gacuo.SetDefaultMaxChannels(*i)
	}
	return gacuo
}

// AddDefaultMaxChannels adds i to the "default_max_channels" field.
func (gacuo *GuildAdminConfigUpdateOne) AddDefaultMaxChannels(i int) *GuildAdminConfigUpdateOne {
	gacuo.mutation.AddDefaultMaxChannels(i)
	return gacuo
}

// ClearDefaultMaxChannels clears the value of the "default_max_channels" field.
func (gacuo *GuildAdminConfigUpdateOne) ClearDefaultMaxChannels() *GuildAdminConfigUpdateOne {
	gacuo.mutation.ClearDefaultMaxChannels()
	return gacuo
}

// SetDefaultMaxClones sets the "default_max_clones" field.
func (gacuo *GuildAdminConfigUpdateOne) SetDefaultMaxClones(i int) *GuildAdminConfigUpdateOne {
	gacuo.mutation.ResetDefaultMaxClones()
	gacuo.mutation.SetDefaultMaxClones(i)
	return gacuo
}

// SetNillableDefaultMaxClones sets the "default_max_clones" field if the given value is not nil.
func (gacuo *GuildAdminConfigUpdateOne) SetNillableDefaultMaxClones(i *int) *GuildAdminConfigUpdateOne {
	if i != nil {
		gacuo.SetDefaultMaxClones(*i)
	}
	return gacuo
}

// AddDefaultMaxClones adds i to the "default_max_clones" field.
func (gacuo *GuildAdminConfigUpdateOne) AddDefaultMaxClones(i int) *GuildAdminConfigUpdateOne {
	gacuo.mutation.AddDefaultMaxClones(i)
	return gacuo
}

// ClearDefaultMaxClones clears the value of the "default_max_clones" field.
func (gacuo *GuildAdminConfigUpdateOne) ClearDefaultMaxClones() *GuildAdminConfigUpdateOne {
	gacuo.mutation.ClearDefaultMaxClones()
	return gacuo
}

// SetComment sets the "comment" field.
func (gacuo *GuildAdminConfigUpdateOne) SetComment(s string) *GuildAdminConfigUpdateOne {
	gacuo.mutation.SetComment(s)
	return gacuo
}

// SetNillableComment sets the "comment" field if the given value is not nil.
func (gacuo *GuildAdminConfigUpdateOne) SetNillableComment(s *string) *GuildAdminConfigUpdateOne {
	if s != nil {
		gacuo.SetComment(*s)
	}
	return gacuo
}

// ClearComment clears the value of the "comment" field.
func (gacuo *GuildAdminConfigUpdateOne) ClearComment() *GuildAdminConfigUpdateOne {
	gacuo.mutation.ClearComment()
	return gacuo
}

// Mutation returns the GuildAdminConfigMutation object of the builder.
func (gacuo *GuildAdminConfigUpdateOne) Mutation() *GuildAdminConfigMutation {
	return gacuo.mutation
}

// Where appends a list predicates to the GuildAdminConfigUpdate builder.
func (gacuo *GuildAdminConfigUpdateOne) Where(ps ...predicate.GuildAdminConfig) *GuildAdminConfigUpdateOne {
	gacuo.mutation.Where(ps...)
	return gacuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (gacuo *GuildAdminConfigUpdateOne) Select(field string, fields ...string) *GuildAdminConfigUpdateOne {
	gacuo.fields = append([]string{field}, fields...)
	return gacuo
}

// Save executes the query and returns the updated GuildAdminConfig entity.
func (gacuo *GuildAdminConfigUpdateOne) Save(ctx context.Context) (*GuildAdminConfig, error) {
	if err := gacuo.defaults(); err != nil {
		return nil, err
	}
	return withHooks[*GuildAdminConfig, GuildAdminConfigMutation](ctx, gacuo.sqlSave, gacuo.mutation, gacuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (gacuo *GuildAdminConfigUpdateOne) SaveX(ctx context.Context) *GuildAdminConfig {
	node, err := gacuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (gacuo *GuildAdminConfigUpdateOne) Exec(ctx context.Context) error {
	_, err := gacuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (gacuo *GuildAdminConfigUpdateOne) ExecX(ctx context.Context) {
	if err := gacuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (gacuo *GuildAdminConfigUpdateOne) defaults() error {
	if _, ok := gacuo.mutation.UpdateTime(); !ok {
		if guildadminconfig.UpdateDefaultUpdateTime == nil {
			return fmt.Errorf("ent: uninitialized guildadminconfig.UpdateDefaultUpdateTime (forgotten import ent/runtime?)")
		}
		v := guildadminconfig.UpdateDefaultUpdateTime()
		gacuo.mutation.SetUpdateTime(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (gacuo *GuildAdminConfigUpdateOne) check() error {
	if v, ok := gacuo.mutation.DefaultMaxChannels(); ok {
		if err := guildadminconfig.DefaultMaxChannelsValidator(v); err != nil {
			return &ValidationError{Name: "default_max_channels", err: fmt.Errorf(`ent: validator failed for field "GuildAdminConfig.default_max_channels": %w`, err)}
		}
	}
	if v, ok := gacuo.mutation.DefaultMaxClones(); ok {
		if err := guildadminconfig.DefaultMaxClonesValidator(v); err != nil {
			return &ValidationError{Name: "default_max_clones", err: fmt.Errorf(`ent: validator failed for field "GuildAdminConfig.default_max_clones": %w`, err)}
		}
	}
	if _, ok := gacuo.mutation.GuildID(); gacuo.mutation.GuildCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "GuildAdminConfig.guild"`)
	}
	return nil
}

func (gacuo *GuildAdminConfigUpdateOne) sqlSave(ctx context.Context) (_node *GuildAdminConfig, err error) {
	if err := gacuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(guildadminconfig.Table, guildadminconfig.Columns, sqlgraph.NewFieldSpec(guildadminconfig.FieldID, field.TypeInt))
	id, ok := gacuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "GuildAdminConfig.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := gacuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, guildadminconfig.FieldID)
		for _, f := range fields {
			if !guildadminconfig.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != guildadminconfig.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := gacuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := gacuo.mutation.UpdateTime(); ok {
		_spec.SetField(guildadminconfig.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := gacuo.mutation.Enabled(); ok {
		_spec.SetField(guildadminconfig.FieldEnabled, field.TypeBool, value)
	}
	if gacuo.mutation.EnabledCleared() {
		_spec.ClearField(guildadminconfig.FieldEnabled, field.TypeBool)
	}
	if value, ok := gacuo.mutation.DefaultMaxChannels(); ok {
		_spec.SetField(guildadminconfig.FieldDefaultMaxChannels, field.TypeInt, value)
	}
	if value, ok := gacuo.mutation.AddedDefaultMaxChannels(); ok {
		_spec.AddField(guildadminconfig.FieldDefaultMaxChannels, field.TypeInt, value)
	}
	if gacuo.mutation.DefaultMaxChannelsCleared() {
		_spec.ClearField(guildadminconfig.FieldDefaultMaxChannels, field.TypeInt)
	}
	if value, ok := gacuo.mutation.DefaultMaxClones(); ok {
		_spec.SetField(guildadminconfig.FieldDefaultMaxClones, field.TypeInt, value)
	}
	if value, ok := gacuo.mutation.AddedDefaultMaxClones(); ok {
		_spec.AddField(guildadminconfig.FieldDefaultMaxClones, field.TypeInt, value)
	}
	if gacuo.mutation.DefaultMaxClonesCleared() {
		_spec.ClearField(guildadminconfig.FieldDefaultMaxClones, field.TypeInt)
	}
	if value, ok := gacuo.mutation.Comment(); ok {
		_spec.SetField(guildadminconfig.FieldComment, field.TypeString, value)
	}
	if gacuo.mutation.CommentCleared() {
		_spec.ClearField(guildadminconfig.FieldComment, field.TypeString)
	}
	_node = &GuildAdminConfig{config: gacuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, gacuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{guildadminconfig.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	gacuo.mutation.done = true
	return _node, nil
}
