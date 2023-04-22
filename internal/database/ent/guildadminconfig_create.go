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
	"github.com/lrstanley/spectrograph/internal/database/ent/guild"
	"github.com/lrstanley/spectrograph/internal/database/ent/guildadminconfig"
)

// GuildAdminConfigCreate is the builder for creating a GuildAdminConfig entity.
type GuildAdminConfigCreate struct {
	config
	mutation *GuildAdminConfigMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetCreateTime sets the "create_time" field.
func (gacc *GuildAdminConfigCreate) SetCreateTime(t time.Time) *GuildAdminConfigCreate {
	gacc.mutation.SetCreateTime(t)
	return gacc
}

// SetNillableCreateTime sets the "create_time" field if the given value is not nil.
func (gacc *GuildAdminConfigCreate) SetNillableCreateTime(t *time.Time) *GuildAdminConfigCreate {
	if t != nil {
		gacc.SetCreateTime(*t)
	}
	return gacc
}

// SetUpdateTime sets the "update_time" field.
func (gacc *GuildAdminConfigCreate) SetUpdateTime(t time.Time) *GuildAdminConfigCreate {
	gacc.mutation.SetUpdateTime(t)
	return gacc
}

// SetNillableUpdateTime sets the "update_time" field if the given value is not nil.
func (gacc *GuildAdminConfigCreate) SetNillableUpdateTime(t *time.Time) *GuildAdminConfigCreate {
	if t != nil {
		gacc.SetUpdateTime(*t)
	}
	return gacc
}

// SetEnabled sets the "enabled" field.
func (gacc *GuildAdminConfigCreate) SetEnabled(b bool) *GuildAdminConfigCreate {
	gacc.mutation.SetEnabled(b)
	return gacc
}

// SetNillableEnabled sets the "enabled" field if the given value is not nil.
func (gacc *GuildAdminConfigCreate) SetNillableEnabled(b *bool) *GuildAdminConfigCreate {
	if b != nil {
		gacc.SetEnabled(*b)
	}
	return gacc
}

// SetDefaultMaxChannels sets the "default_max_channels" field.
func (gacc *GuildAdminConfigCreate) SetDefaultMaxChannels(i int) *GuildAdminConfigCreate {
	gacc.mutation.SetDefaultMaxChannels(i)
	return gacc
}

// SetNillableDefaultMaxChannels sets the "default_max_channels" field if the given value is not nil.
func (gacc *GuildAdminConfigCreate) SetNillableDefaultMaxChannels(i *int) *GuildAdminConfigCreate {
	if i != nil {
		gacc.SetDefaultMaxChannels(*i)
	}
	return gacc
}

// SetDefaultMaxClones sets the "default_max_clones" field.
func (gacc *GuildAdminConfigCreate) SetDefaultMaxClones(i int) *GuildAdminConfigCreate {
	gacc.mutation.SetDefaultMaxClones(i)
	return gacc
}

// SetNillableDefaultMaxClones sets the "default_max_clones" field if the given value is not nil.
func (gacc *GuildAdminConfigCreate) SetNillableDefaultMaxClones(i *int) *GuildAdminConfigCreate {
	if i != nil {
		gacc.SetDefaultMaxClones(*i)
	}
	return gacc
}

// SetComment sets the "comment" field.
func (gacc *GuildAdminConfigCreate) SetComment(s string) *GuildAdminConfigCreate {
	gacc.mutation.SetComment(s)
	return gacc
}

// SetNillableComment sets the "comment" field if the given value is not nil.
func (gacc *GuildAdminConfigCreate) SetNillableComment(s *string) *GuildAdminConfigCreate {
	if s != nil {
		gacc.SetComment(*s)
	}
	return gacc
}

// SetGuildID sets the "guild" edge to the Guild entity by ID.
func (gacc *GuildAdminConfigCreate) SetGuildID(id int) *GuildAdminConfigCreate {
	gacc.mutation.SetGuildID(id)
	return gacc
}

// SetGuild sets the "guild" edge to the Guild entity.
func (gacc *GuildAdminConfigCreate) SetGuild(g *Guild) *GuildAdminConfigCreate {
	return gacc.SetGuildID(g.ID)
}

// Mutation returns the GuildAdminConfigMutation object of the builder.
func (gacc *GuildAdminConfigCreate) Mutation() *GuildAdminConfigMutation {
	return gacc.mutation
}

// Save creates the GuildAdminConfig in the database.
func (gacc *GuildAdminConfigCreate) Save(ctx context.Context) (*GuildAdminConfig, error) {
	if err := gacc.defaults(); err != nil {
		return nil, err
	}
	return withHooks[*GuildAdminConfig, GuildAdminConfigMutation](ctx, gacc.sqlSave, gacc.mutation, gacc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (gacc *GuildAdminConfigCreate) SaveX(ctx context.Context) *GuildAdminConfig {
	v, err := gacc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (gacc *GuildAdminConfigCreate) Exec(ctx context.Context) error {
	_, err := gacc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (gacc *GuildAdminConfigCreate) ExecX(ctx context.Context) {
	if err := gacc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (gacc *GuildAdminConfigCreate) defaults() error {
	if _, ok := gacc.mutation.CreateTime(); !ok {
		if guildadminconfig.DefaultCreateTime == nil {
			return fmt.Errorf("ent: uninitialized guildadminconfig.DefaultCreateTime (forgotten import ent/runtime?)")
		}
		v := guildadminconfig.DefaultCreateTime()
		gacc.mutation.SetCreateTime(v)
	}
	if _, ok := gacc.mutation.UpdateTime(); !ok {
		if guildadminconfig.DefaultUpdateTime == nil {
			return fmt.Errorf("ent: uninitialized guildadminconfig.DefaultUpdateTime (forgotten import ent/runtime?)")
		}
		v := guildadminconfig.DefaultUpdateTime()
		gacc.mutation.SetUpdateTime(v)
	}
	if _, ok := gacc.mutation.Enabled(); !ok {
		v := guildadminconfig.DefaultEnabled
		gacc.mutation.SetEnabled(v)
	}
	if _, ok := gacc.mutation.DefaultMaxChannels(); !ok {
		v := guildadminconfig.DefaultDefaultMaxChannels
		gacc.mutation.SetDefaultMaxChannels(v)
	}
	if _, ok := gacc.mutation.DefaultMaxClones(); !ok {
		v := guildadminconfig.DefaultDefaultMaxClones
		gacc.mutation.SetDefaultMaxClones(v)
	}
	if _, ok := gacc.mutation.Comment(); !ok {
		v := guildadminconfig.DefaultComment
		gacc.mutation.SetComment(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (gacc *GuildAdminConfigCreate) check() error {
	if _, ok := gacc.mutation.CreateTime(); !ok {
		return &ValidationError{Name: "create_time", err: errors.New(`ent: missing required field "GuildAdminConfig.create_time"`)}
	}
	if _, ok := gacc.mutation.UpdateTime(); !ok {
		return &ValidationError{Name: "update_time", err: errors.New(`ent: missing required field "GuildAdminConfig.update_time"`)}
	}
	if v, ok := gacc.mutation.DefaultMaxChannels(); ok {
		if err := guildadminconfig.DefaultMaxChannelsValidator(v); err != nil {
			return &ValidationError{Name: "default_max_channels", err: fmt.Errorf(`ent: validator failed for field "GuildAdminConfig.default_max_channels": %w`, err)}
		}
	}
	if v, ok := gacc.mutation.DefaultMaxClones(); ok {
		if err := guildadminconfig.DefaultMaxClonesValidator(v); err != nil {
			return &ValidationError{Name: "default_max_clones", err: fmt.Errorf(`ent: validator failed for field "GuildAdminConfig.default_max_clones": %w`, err)}
		}
	}
	if _, ok := gacc.mutation.GuildID(); !ok {
		return &ValidationError{Name: "guild", err: errors.New(`ent: missing required edge "GuildAdminConfig.guild"`)}
	}
	return nil
}

func (gacc *GuildAdminConfigCreate) sqlSave(ctx context.Context) (*GuildAdminConfig, error) {
	if err := gacc.check(); err != nil {
		return nil, err
	}
	_node, _spec := gacc.createSpec()
	if err := sqlgraph.CreateNode(ctx, gacc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	gacc.mutation.id = &_node.ID
	gacc.mutation.done = true
	return _node, nil
}

func (gacc *GuildAdminConfigCreate) createSpec() (*GuildAdminConfig, *sqlgraph.CreateSpec) {
	var (
		_node = &GuildAdminConfig{config: gacc.config}
		_spec = sqlgraph.NewCreateSpec(guildadminconfig.Table, sqlgraph.NewFieldSpec(guildadminconfig.FieldID, field.TypeInt))
	)
	_spec.OnConflict = gacc.conflict
	if value, ok := gacc.mutation.CreateTime(); ok {
		_spec.SetField(guildadminconfig.FieldCreateTime, field.TypeTime, value)
		_node.CreateTime = value
	}
	if value, ok := gacc.mutation.UpdateTime(); ok {
		_spec.SetField(guildadminconfig.FieldUpdateTime, field.TypeTime, value)
		_node.UpdateTime = value
	}
	if value, ok := gacc.mutation.Enabled(); ok {
		_spec.SetField(guildadminconfig.FieldEnabled, field.TypeBool, value)
		_node.Enabled = value
	}
	if value, ok := gacc.mutation.DefaultMaxChannels(); ok {
		_spec.SetField(guildadminconfig.FieldDefaultMaxChannels, field.TypeInt, value)
		_node.DefaultMaxChannels = value
	}
	if value, ok := gacc.mutation.DefaultMaxClones(); ok {
		_spec.SetField(guildadminconfig.FieldDefaultMaxClones, field.TypeInt, value)
		_node.DefaultMaxClones = value
	}
	if value, ok := gacc.mutation.Comment(); ok {
		_spec.SetField(guildadminconfig.FieldComment, field.TypeString, value)
		_node.Comment = value
	}
	if nodes := gacc.mutation.GuildIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   guildadminconfig.GuildTable,
			Columns: []string{guildadminconfig.GuildColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(guild.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.guild_guild_admin_config = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.GuildAdminConfig.Create().
//		SetCreateTime(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.GuildAdminConfigUpsert) {
//			SetCreateTime(v+v).
//		}).
//		Exec(ctx)
func (gacc *GuildAdminConfigCreate) OnConflict(opts ...sql.ConflictOption) *GuildAdminConfigUpsertOne {
	gacc.conflict = opts
	return &GuildAdminConfigUpsertOne{
		create: gacc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.GuildAdminConfig.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (gacc *GuildAdminConfigCreate) OnConflictColumns(columns ...string) *GuildAdminConfigUpsertOne {
	gacc.conflict = append(gacc.conflict, sql.ConflictColumns(columns...))
	return &GuildAdminConfigUpsertOne{
		create: gacc,
	}
}

type (
	// GuildAdminConfigUpsertOne is the builder for "upsert"-ing
	//  one GuildAdminConfig node.
	GuildAdminConfigUpsertOne struct {
		create *GuildAdminConfigCreate
	}

	// GuildAdminConfigUpsert is the "OnConflict" setter.
	GuildAdminConfigUpsert struct {
		*sql.UpdateSet
	}
)

// SetUpdateTime sets the "update_time" field.
func (u *GuildAdminConfigUpsert) SetUpdateTime(v time.Time) *GuildAdminConfigUpsert {
	u.Set(guildadminconfig.FieldUpdateTime, v)
	return u
}

// UpdateUpdateTime sets the "update_time" field to the value that was provided on create.
func (u *GuildAdminConfigUpsert) UpdateUpdateTime() *GuildAdminConfigUpsert {
	u.SetExcluded(guildadminconfig.FieldUpdateTime)
	return u
}

// SetEnabled sets the "enabled" field.
func (u *GuildAdminConfigUpsert) SetEnabled(v bool) *GuildAdminConfigUpsert {
	u.Set(guildadminconfig.FieldEnabled, v)
	return u
}

// UpdateEnabled sets the "enabled" field to the value that was provided on create.
func (u *GuildAdminConfigUpsert) UpdateEnabled() *GuildAdminConfigUpsert {
	u.SetExcluded(guildadminconfig.FieldEnabled)
	return u
}

// ClearEnabled clears the value of the "enabled" field.
func (u *GuildAdminConfigUpsert) ClearEnabled() *GuildAdminConfigUpsert {
	u.SetNull(guildadminconfig.FieldEnabled)
	return u
}

// SetDefaultMaxChannels sets the "default_max_channels" field.
func (u *GuildAdminConfigUpsert) SetDefaultMaxChannels(v int) *GuildAdminConfigUpsert {
	u.Set(guildadminconfig.FieldDefaultMaxChannels, v)
	return u
}

// UpdateDefaultMaxChannels sets the "default_max_channels" field to the value that was provided on create.
func (u *GuildAdminConfigUpsert) UpdateDefaultMaxChannels() *GuildAdminConfigUpsert {
	u.SetExcluded(guildadminconfig.FieldDefaultMaxChannels)
	return u
}

// AddDefaultMaxChannels adds v to the "default_max_channels" field.
func (u *GuildAdminConfigUpsert) AddDefaultMaxChannels(v int) *GuildAdminConfigUpsert {
	u.Add(guildadminconfig.FieldDefaultMaxChannels, v)
	return u
}

// ClearDefaultMaxChannels clears the value of the "default_max_channels" field.
func (u *GuildAdminConfigUpsert) ClearDefaultMaxChannels() *GuildAdminConfigUpsert {
	u.SetNull(guildadminconfig.FieldDefaultMaxChannels)
	return u
}

// SetDefaultMaxClones sets the "default_max_clones" field.
func (u *GuildAdminConfigUpsert) SetDefaultMaxClones(v int) *GuildAdminConfigUpsert {
	u.Set(guildadminconfig.FieldDefaultMaxClones, v)
	return u
}

// UpdateDefaultMaxClones sets the "default_max_clones" field to the value that was provided on create.
func (u *GuildAdminConfigUpsert) UpdateDefaultMaxClones() *GuildAdminConfigUpsert {
	u.SetExcluded(guildadminconfig.FieldDefaultMaxClones)
	return u
}

// AddDefaultMaxClones adds v to the "default_max_clones" field.
func (u *GuildAdminConfigUpsert) AddDefaultMaxClones(v int) *GuildAdminConfigUpsert {
	u.Add(guildadminconfig.FieldDefaultMaxClones, v)
	return u
}

// ClearDefaultMaxClones clears the value of the "default_max_clones" field.
func (u *GuildAdminConfigUpsert) ClearDefaultMaxClones() *GuildAdminConfigUpsert {
	u.SetNull(guildadminconfig.FieldDefaultMaxClones)
	return u
}

// SetComment sets the "comment" field.
func (u *GuildAdminConfigUpsert) SetComment(v string) *GuildAdminConfigUpsert {
	u.Set(guildadminconfig.FieldComment, v)
	return u
}

// UpdateComment sets the "comment" field to the value that was provided on create.
func (u *GuildAdminConfigUpsert) UpdateComment() *GuildAdminConfigUpsert {
	u.SetExcluded(guildadminconfig.FieldComment)
	return u
}

// ClearComment clears the value of the "comment" field.
func (u *GuildAdminConfigUpsert) ClearComment() *GuildAdminConfigUpsert {
	u.SetNull(guildadminconfig.FieldComment)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.GuildAdminConfig.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *GuildAdminConfigUpsertOne) UpdateNewValues() *GuildAdminConfigUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.CreateTime(); exists {
			s.SetIgnore(guildadminconfig.FieldCreateTime)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.GuildAdminConfig.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *GuildAdminConfigUpsertOne) Ignore() *GuildAdminConfigUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *GuildAdminConfigUpsertOne) DoNothing() *GuildAdminConfigUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the GuildAdminConfigCreate.OnConflict
// documentation for more info.
func (u *GuildAdminConfigUpsertOne) Update(set func(*GuildAdminConfigUpsert)) *GuildAdminConfigUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&GuildAdminConfigUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdateTime sets the "update_time" field.
func (u *GuildAdminConfigUpsertOne) SetUpdateTime(v time.Time) *GuildAdminConfigUpsertOne {
	return u.Update(func(s *GuildAdminConfigUpsert) {
		s.SetUpdateTime(v)
	})
}

// UpdateUpdateTime sets the "update_time" field to the value that was provided on create.
func (u *GuildAdminConfigUpsertOne) UpdateUpdateTime() *GuildAdminConfigUpsertOne {
	return u.Update(func(s *GuildAdminConfigUpsert) {
		s.UpdateUpdateTime()
	})
}

// SetEnabled sets the "enabled" field.
func (u *GuildAdminConfigUpsertOne) SetEnabled(v bool) *GuildAdminConfigUpsertOne {
	return u.Update(func(s *GuildAdminConfigUpsert) {
		s.SetEnabled(v)
	})
}

// UpdateEnabled sets the "enabled" field to the value that was provided on create.
func (u *GuildAdminConfigUpsertOne) UpdateEnabled() *GuildAdminConfigUpsertOne {
	return u.Update(func(s *GuildAdminConfigUpsert) {
		s.UpdateEnabled()
	})
}

// ClearEnabled clears the value of the "enabled" field.
func (u *GuildAdminConfigUpsertOne) ClearEnabled() *GuildAdminConfigUpsertOne {
	return u.Update(func(s *GuildAdminConfigUpsert) {
		s.ClearEnabled()
	})
}

// SetDefaultMaxChannels sets the "default_max_channels" field.
func (u *GuildAdminConfigUpsertOne) SetDefaultMaxChannels(v int) *GuildAdminConfigUpsertOne {
	return u.Update(func(s *GuildAdminConfigUpsert) {
		s.SetDefaultMaxChannels(v)
	})
}

// AddDefaultMaxChannels adds v to the "default_max_channels" field.
func (u *GuildAdminConfigUpsertOne) AddDefaultMaxChannels(v int) *GuildAdminConfigUpsertOne {
	return u.Update(func(s *GuildAdminConfigUpsert) {
		s.AddDefaultMaxChannels(v)
	})
}

// UpdateDefaultMaxChannels sets the "default_max_channels" field to the value that was provided on create.
func (u *GuildAdminConfigUpsertOne) UpdateDefaultMaxChannels() *GuildAdminConfigUpsertOne {
	return u.Update(func(s *GuildAdminConfigUpsert) {
		s.UpdateDefaultMaxChannels()
	})
}

// ClearDefaultMaxChannels clears the value of the "default_max_channels" field.
func (u *GuildAdminConfigUpsertOne) ClearDefaultMaxChannels() *GuildAdminConfigUpsertOne {
	return u.Update(func(s *GuildAdminConfigUpsert) {
		s.ClearDefaultMaxChannels()
	})
}

// SetDefaultMaxClones sets the "default_max_clones" field.
func (u *GuildAdminConfigUpsertOne) SetDefaultMaxClones(v int) *GuildAdminConfigUpsertOne {
	return u.Update(func(s *GuildAdminConfigUpsert) {
		s.SetDefaultMaxClones(v)
	})
}

// AddDefaultMaxClones adds v to the "default_max_clones" field.
func (u *GuildAdminConfigUpsertOne) AddDefaultMaxClones(v int) *GuildAdminConfigUpsertOne {
	return u.Update(func(s *GuildAdminConfigUpsert) {
		s.AddDefaultMaxClones(v)
	})
}

// UpdateDefaultMaxClones sets the "default_max_clones" field to the value that was provided on create.
func (u *GuildAdminConfigUpsertOne) UpdateDefaultMaxClones() *GuildAdminConfigUpsertOne {
	return u.Update(func(s *GuildAdminConfigUpsert) {
		s.UpdateDefaultMaxClones()
	})
}

// ClearDefaultMaxClones clears the value of the "default_max_clones" field.
func (u *GuildAdminConfigUpsertOne) ClearDefaultMaxClones() *GuildAdminConfigUpsertOne {
	return u.Update(func(s *GuildAdminConfigUpsert) {
		s.ClearDefaultMaxClones()
	})
}

// SetComment sets the "comment" field.
func (u *GuildAdminConfigUpsertOne) SetComment(v string) *GuildAdminConfigUpsertOne {
	return u.Update(func(s *GuildAdminConfigUpsert) {
		s.SetComment(v)
	})
}

// UpdateComment sets the "comment" field to the value that was provided on create.
func (u *GuildAdminConfigUpsertOne) UpdateComment() *GuildAdminConfigUpsertOne {
	return u.Update(func(s *GuildAdminConfigUpsert) {
		s.UpdateComment()
	})
}

// ClearComment clears the value of the "comment" field.
func (u *GuildAdminConfigUpsertOne) ClearComment() *GuildAdminConfigUpsertOne {
	return u.Update(func(s *GuildAdminConfigUpsert) {
		s.ClearComment()
	})
}

// Exec executes the query.
func (u *GuildAdminConfigUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for GuildAdminConfigCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *GuildAdminConfigUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *GuildAdminConfigUpsertOne) ID(ctx context.Context) (id int, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *GuildAdminConfigUpsertOne) IDX(ctx context.Context) int {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// GuildAdminConfigCreateBulk is the builder for creating many GuildAdminConfig entities in bulk.
type GuildAdminConfigCreateBulk struct {
	config
	builders []*GuildAdminConfigCreate
	conflict []sql.ConflictOption
}

// Save creates the GuildAdminConfig entities in the database.
func (gaccb *GuildAdminConfigCreateBulk) Save(ctx context.Context) ([]*GuildAdminConfig, error) {
	specs := make([]*sqlgraph.CreateSpec, len(gaccb.builders))
	nodes := make([]*GuildAdminConfig, len(gaccb.builders))
	mutators := make([]Mutator, len(gaccb.builders))
	for i := range gaccb.builders {
		func(i int, root context.Context) {
			builder := gaccb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*GuildAdminConfigMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, gaccb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = gaccb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, gaccb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, gaccb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (gaccb *GuildAdminConfigCreateBulk) SaveX(ctx context.Context) []*GuildAdminConfig {
	v, err := gaccb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (gaccb *GuildAdminConfigCreateBulk) Exec(ctx context.Context) error {
	_, err := gaccb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (gaccb *GuildAdminConfigCreateBulk) ExecX(ctx context.Context) {
	if err := gaccb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.GuildAdminConfig.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.GuildAdminConfigUpsert) {
//			SetCreateTime(v+v).
//		}).
//		Exec(ctx)
func (gaccb *GuildAdminConfigCreateBulk) OnConflict(opts ...sql.ConflictOption) *GuildAdminConfigUpsertBulk {
	gaccb.conflict = opts
	return &GuildAdminConfigUpsertBulk{
		create: gaccb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.GuildAdminConfig.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (gaccb *GuildAdminConfigCreateBulk) OnConflictColumns(columns ...string) *GuildAdminConfigUpsertBulk {
	gaccb.conflict = append(gaccb.conflict, sql.ConflictColumns(columns...))
	return &GuildAdminConfigUpsertBulk{
		create: gaccb,
	}
}

// GuildAdminConfigUpsertBulk is the builder for "upsert"-ing
// a bulk of GuildAdminConfig nodes.
type GuildAdminConfigUpsertBulk struct {
	create *GuildAdminConfigCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.GuildAdminConfig.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *GuildAdminConfigUpsertBulk) UpdateNewValues() *GuildAdminConfigUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.CreateTime(); exists {
				s.SetIgnore(guildadminconfig.FieldCreateTime)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.GuildAdminConfig.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *GuildAdminConfigUpsertBulk) Ignore() *GuildAdminConfigUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *GuildAdminConfigUpsertBulk) DoNothing() *GuildAdminConfigUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the GuildAdminConfigCreateBulk.OnConflict
// documentation for more info.
func (u *GuildAdminConfigUpsertBulk) Update(set func(*GuildAdminConfigUpsert)) *GuildAdminConfigUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&GuildAdminConfigUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdateTime sets the "update_time" field.
func (u *GuildAdminConfigUpsertBulk) SetUpdateTime(v time.Time) *GuildAdminConfigUpsertBulk {
	return u.Update(func(s *GuildAdminConfigUpsert) {
		s.SetUpdateTime(v)
	})
}

// UpdateUpdateTime sets the "update_time" field to the value that was provided on create.
func (u *GuildAdminConfigUpsertBulk) UpdateUpdateTime() *GuildAdminConfigUpsertBulk {
	return u.Update(func(s *GuildAdminConfigUpsert) {
		s.UpdateUpdateTime()
	})
}

// SetEnabled sets the "enabled" field.
func (u *GuildAdminConfigUpsertBulk) SetEnabled(v bool) *GuildAdminConfigUpsertBulk {
	return u.Update(func(s *GuildAdminConfigUpsert) {
		s.SetEnabled(v)
	})
}

// UpdateEnabled sets the "enabled" field to the value that was provided on create.
func (u *GuildAdminConfigUpsertBulk) UpdateEnabled() *GuildAdminConfigUpsertBulk {
	return u.Update(func(s *GuildAdminConfigUpsert) {
		s.UpdateEnabled()
	})
}

// ClearEnabled clears the value of the "enabled" field.
func (u *GuildAdminConfigUpsertBulk) ClearEnabled() *GuildAdminConfigUpsertBulk {
	return u.Update(func(s *GuildAdminConfigUpsert) {
		s.ClearEnabled()
	})
}

// SetDefaultMaxChannels sets the "default_max_channels" field.
func (u *GuildAdminConfigUpsertBulk) SetDefaultMaxChannels(v int) *GuildAdminConfigUpsertBulk {
	return u.Update(func(s *GuildAdminConfigUpsert) {
		s.SetDefaultMaxChannels(v)
	})
}

// AddDefaultMaxChannels adds v to the "default_max_channels" field.
func (u *GuildAdminConfigUpsertBulk) AddDefaultMaxChannels(v int) *GuildAdminConfigUpsertBulk {
	return u.Update(func(s *GuildAdminConfigUpsert) {
		s.AddDefaultMaxChannels(v)
	})
}

// UpdateDefaultMaxChannels sets the "default_max_channels" field to the value that was provided on create.
func (u *GuildAdminConfigUpsertBulk) UpdateDefaultMaxChannels() *GuildAdminConfigUpsertBulk {
	return u.Update(func(s *GuildAdminConfigUpsert) {
		s.UpdateDefaultMaxChannels()
	})
}

// ClearDefaultMaxChannels clears the value of the "default_max_channels" field.
func (u *GuildAdminConfigUpsertBulk) ClearDefaultMaxChannels() *GuildAdminConfigUpsertBulk {
	return u.Update(func(s *GuildAdminConfigUpsert) {
		s.ClearDefaultMaxChannels()
	})
}

// SetDefaultMaxClones sets the "default_max_clones" field.
func (u *GuildAdminConfigUpsertBulk) SetDefaultMaxClones(v int) *GuildAdminConfigUpsertBulk {
	return u.Update(func(s *GuildAdminConfigUpsert) {
		s.SetDefaultMaxClones(v)
	})
}

// AddDefaultMaxClones adds v to the "default_max_clones" field.
func (u *GuildAdminConfigUpsertBulk) AddDefaultMaxClones(v int) *GuildAdminConfigUpsertBulk {
	return u.Update(func(s *GuildAdminConfigUpsert) {
		s.AddDefaultMaxClones(v)
	})
}

// UpdateDefaultMaxClones sets the "default_max_clones" field to the value that was provided on create.
func (u *GuildAdminConfigUpsertBulk) UpdateDefaultMaxClones() *GuildAdminConfigUpsertBulk {
	return u.Update(func(s *GuildAdminConfigUpsert) {
		s.UpdateDefaultMaxClones()
	})
}

// ClearDefaultMaxClones clears the value of the "default_max_clones" field.
func (u *GuildAdminConfigUpsertBulk) ClearDefaultMaxClones() *GuildAdminConfigUpsertBulk {
	return u.Update(func(s *GuildAdminConfigUpsert) {
		s.ClearDefaultMaxClones()
	})
}

// SetComment sets the "comment" field.
func (u *GuildAdminConfigUpsertBulk) SetComment(v string) *GuildAdminConfigUpsertBulk {
	return u.Update(func(s *GuildAdminConfigUpsert) {
		s.SetComment(v)
	})
}

// UpdateComment sets the "comment" field to the value that was provided on create.
func (u *GuildAdminConfigUpsertBulk) UpdateComment() *GuildAdminConfigUpsertBulk {
	return u.Update(func(s *GuildAdminConfigUpsert) {
		s.UpdateComment()
	})
}

// ClearComment clears the value of the "comment" field.
func (u *GuildAdminConfigUpsertBulk) ClearComment() *GuildAdminConfigUpsertBulk {
	return u.Update(func(s *GuildAdminConfigUpsert) {
		s.ClearComment()
	})
}

// Exec executes the query.
func (u *GuildAdminConfigUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the GuildAdminConfigCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for GuildAdminConfigCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *GuildAdminConfigUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
