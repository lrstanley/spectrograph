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
	"github.com/lrstanley/spectrograph/internal/database/ent/guildconfig"
)

// GuildConfigCreate is the builder for creating a GuildConfig entity.
type GuildConfigCreate struct {
	config
	mutation *GuildConfigMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetCreateTime sets the "create_time" field.
func (gcc *GuildConfigCreate) SetCreateTime(t time.Time) *GuildConfigCreate {
	gcc.mutation.SetCreateTime(t)
	return gcc
}

// SetNillableCreateTime sets the "create_time" field if the given value is not nil.
func (gcc *GuildConfigCreate) SetNillableCreateTime(t *time.Time) *GuildConfigCreate {
	if t != nil {
		gcc.SetCreateTime(*t)
	}
	return gcc
}

// SetUpdateTime sets the "update_time" field.
func (gcc *GuildConfigCreate) SetUpdateTime(t time.Time) *GuildConfigCreate {
	gcc.mutation.SetUpdateTime(t)
	return gcc
}

// SetNillableUpdateTime sets the "update_time" field if the given value is not nil.
func (gcc *GuildConfigCreate) SetNillableUpdateTime(t *time.Time) *GuildConfigCreate {
	if t != nil {
		gcc.SetUpdateTime(*t)
	}
	return gcc
}

// SetEnabled sets the "enabled" field.
func (gcc *GuildConfigCreate) SetEnabled(b bool) *GuildConfigCreate {
	gcc.mutation.SetEnabled(b)
	return gcc
}

// SetNillableEnabled sets the "enabled" field if the given value is not nil.
func (gcc *GuildConfigCreate) SetNillableEnabled(b *bool) *GuildConfigCreate {
	if b != nil {
		gcc.SetEnabled(*b)
	}
	return gcc
}

// SetDefaultMaxClones sets the "default_max_clones" field.
func (gcc *GuildConfigCreate) SetDefaultMaxClones(i int) *GuildConfigCreate {
	gcc.mutation.SetDefaultMaxClones(i)
	return gcc
}

// SetNillableDefaultMaxClones sets the "default_max_clones" field if the given value is not nil.
func (gcc *GuildConfigCreate) SetNillableDefaultMaxClones(i *int) *GuildConfigCreate {
	if i != nil {
		gcc.SetDefaultMaxClones(*i)
	}
	return gcc
}

// SetRegexMatch sets the "regex_match" field.
func (gcc *GuildConfigCreate) SetRegexMatch(s string) *GuildConfigCreate {
	gcc.mutation.SetRegexMatch(s)
	return gcc
}

// SetNillableRegexMatch sets the "regex_match" field if the given value is not nil.
func (gcc *GuildConfigCreate) SetNillableRegexMatch(s *string) *GuildConfigCreate {
	if s != nil {
		gcc.SetRegexMatch(*s)
	}
	return gcc
}

// SetContactEmail sets the "contact_email" field.
func (gcc *GuildConfigCreate) SetContactEmail(s string) *GuildConfigCreate {
	gcc.mutation.SetContactEmail(s)
	return gcc
}

// SetNillableContactEmail sets the "contact_email" field if the given value is not nil.
func (gcc *GuildConfigCreate) SetNillableContactEmail(s *string) *GuildConfigCreate {
	if s != nil {
		gcc.SetContactEmail(*s)
	}
	return gcc
}

// SetGuildID sets the "guild" edge to the Guild entity by ID.
func (gcc *GuildConfigCreate) SetGuildID(id int) *GuildConfigCreate {
	gcc.mutation.SetGuildID(id)
	return gcc
}

// SetGuild sets the "guild" edge to the Guild entity.
func (gcc *GuildConfigCreate) SetGuild(g *Guild) *GuildConfigCreate {
	return gcc.SetGuildID(g.ID)
}

// Mutation returns the GuildConfigMutation object of the builder.
func (gcc *GuildConfigCreate) Mutation() *GuildConfigMutation {
	return gcc.mutation
}

// Save creates the GuildConfig in the database.
func (gcc *GuildConfigCreate) Save(ctx context.Context) (*GuildConfig, error) {
	if err := gcc.defaults(); err != nil {
		return nil, err
	}
	return withHooks[*GuildConfig, GuildConfigMutation](ctx, gcc.sqlSave, gcc.mutation, gcc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (gcc *GuildConfigCreate) SaveX(ctx context.Context) *GuildConfig {
	v, err := gcc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (gcc *GuildConfigCreate) Exec(ctx context.Context) error {
	_, err := gcc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (gcc *GuildConfigCreate) ExecX(ctx context.Context) {
	if err := gcc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (gcc *GuildConfigCreate) defaults() error {
	if _, ok := gcc.mutation.CreateTime(); !ok {
		if guildconfig.DefaultCreateTime == nil {
			return fmt.Errorf("ent: uninitialized guildconfig.DefaultCreateTime (forgotten import ent/runtime?)")
		}
		v := guildconfig.DefaultCreateTime()
		gcc.mutation.SetCreateTime(v)
	}
	if _, ok := gcc.mutation.UpdateTime(); !ok {
		if guildconfig.DefaultUpdateTime == nil {
			return fmt.Errorf("ent: uninitialized guildconfig.DefaultUpdateTime (forgotten import ent/runtime?)")
		}
		v := guildconfig.DefaultUpdateTime()
		gcc.mutation.SetUpdateTime(v)
	}
	if _, ok := gcc.mutation.Enabled(); !ok {
		v := guildconfig.DefaultEnabled
		gcc.mutation.SetEnabled(v)
	}
	if _, ok := gcc.mutation.DefaultMaxClones(); !ok {
		v := guildconfig.DefaultDefaultMaxClones
		gcc.mutation.SetDefaultMaxClones(v)
	}
	if _, ok := gcc.mutation.RegexMatch(); !ok {
		v := guildconfig.DefaultRegexMatch
		gcc.mutation.SetRegexMatch(v)
	}
	if _, ok := gcc.mutation.ContactEmail(); !ok {
		v := guildconfig.DefaultContactEmail
		gcc.mutation.SetContactEmail(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (gcc *GuildConfigCreate) check() error {
	if _, ok := gcc.mutation.CreateTime(); !ok {
		return &ValidationError{Name: "create_time", err: errors.New(`ent: missing required field "GuildConfig.create_time"`)}
	}
	if _, ok := gcc.mutation.UpdateTime(); !ok {
		return &ValidationError{Name: "update_time", err: errors.New(`ent: missing required field "GuildConfig.update_time"`)}
	}
	if v, ok := gcc.mutation.DefaultMaxClones(); ok {
		if err := guildconfig.DefaultMaxClonesValidator(v); err != nil {
			return &ValidationError{Name: "default_max_clones", err: fmt.Errorf(`ent: validator failed for field "GuildConfig.default_max_clones": %w`, err)}
		}
	}
	if v, ok := gcc.mutation.RegexMatch(); ok {
		if err := guildconfig.RegexMatchValidator(v); err != nil {
			return &ValidationError{Name: "regex_match", err: fmt.Errorf(`ent: validator failed for field "GuildConfig.regex_match": %w`, err)}
		}
	}
	if v, ok := gcc.mutation.ContactEmail(); ok {
		if err := guildconfig.ContactEmailValidator(v); err != nil {
			return &ValidationError{Name: "contact_email", err: fmt.Errorf(`ent: validator failed for field "GuildConfig.contact_email": %w`, err)}
		}
	}
	if _, ok := gcc.mutation.GuildID(); !ok {
		return &ValidationError{Name: "guild", err: errors.New(`ent: missing required edge "GuildConfig.guild"`)}
	}
	return nil
}

func (gcc *GuildConfigCreate) sqlSave(ctx context.Context) (*GuildConfig, error) {
	if err := gcc.check(); err != nil {
		return nil, err
	}
	_node, _spec := gcc.createSpec()
	if err := sqlgraph.CreateNode(ctx, gcc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	gcc.mutation.id = &_node.ID
	gcc.mutation.done = true
	return _node, nil
}

func (gcc *GuildConfigCreate) createSpec() (*GuildConfig, *sqlgraph.CreateSpec) {
	var (
		_node = &GuildConfig{config: gcc.config}
		_spec = sqlgraph.NewCreateSpec(guildconfig.Table, sqlgraph.NewFieldSpec(guildconfig.FieldID, field.TypeInt))
	)
	_spec.OnConflict = gcc.conflict
	if value, ok := gcc.mutation.CreateTime(); ok {
		_spec.SetField(guildconfig.FieldCreateTime, field.TypeTime, value)
		_node.CreateTime = value
	}
	if value, ok := gcc.mutation.UpdateTime(); ok {
		_spec.SetField(guildconfig.FieldUpdateTime, field.TypeTime, value)
		_node.UpdateTime = value
	}
	if value, ok := gcc.mutation.Enabled(); ok {
		_spec.SetField(guildconfig.FieldEnabled, field.TypeBool, value)
		_node.Enabled = value
	}
	if value, ok := gcc.mutation.DefaultMaxClones(); ok {
		_spec.SetField(guildconfig.FieldDefaultMaxClones, field.TypeInt, value)
		_node.DefaultMaxClones = value
	}
	if value, ok := gcc.mutation.RegexMatch(); ok {
		_spec.SetField(guildconfig.FieldRegexMatch, field.TypeString, value)
		_node.RegexMatch = value
	}
	if value, ok := gcc.mutation.ContactEmail(); ok {
		_spec.SetField(guildconfig.FieldContactEmail, field.TypeString, value)
		_node.ContactEmail = value
	}
	if nodes := gcc.mutation.GuildIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   guildconfig.GuildTable,
			Columns: []string{guildconfig.GuildColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(guild.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.guild_guild_config = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.GuildConfig.Create().
//		SetCreateTime(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.GuildConfigUpsert) {
//			SetCreateTime(v+v).
//		}).
//		Exec(ctx)
func (gcc *GuildConfigCreate) OnConflict(opts ...sql.ConflictOption) *GuildConfigUpsertOne {
	gcc.conflict = opts
	return &GuildConfigUpsertOne{
		create: gcc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.GuildConfig.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (gcc *GuildConfigCreate) OnConflictColumns(columns ...string) *GuildConfigUpsertOne {
	gcc.conflict = append(gcc.conflict, sql.ConflictColumns(columns...))
	return &GuildConfigUpsertOne{
		create: gcc,
	}
}

type (
	// GuildConfigUpsertOne is the builder for "upsert"-ing
	//  one GuildConfig node.
	GuildConfigUpsertOne struct {
		create *GuildConfigCreate
	}

	// GuildConfigUpsert is the "OnConflict" setter.
	GuildConfigUpsert struct {
		*sql.UpdateSet
	}
)

// SetUpdateTime sets the "update_time" field.
func (u *GuildConfigUpsert) SetUpdateTime(v time.Time) *GuildConfigUpsert {
	u.Set(guildconfig.FieldUpdateTime, v)
	return u
}

// UpdateUpdateTime sets the "update_time" field to the value that was provided on create.
func (u *GuildConfigUpsert) UpdateUpdateTime() *GuildConfigUpsert {
	u.SetExcluded(guildconfig.FieldUpdateTime)
	return u
}

// SetEnabled sets the "enabled" field.
func (u *GuildConfigUpsert) SetEnabled(v bool) *GuildConfigUpsert {
	u.Set(guildconfig.FieldEnabled, v)
	return u
}

// UpdateEnabled sets the "enabled" field to the value that was provided on create.
func (u *GuildConfigUpsert) UpdateEnabled() *GuildConfigUpsert {
	u.SetExcluded(guildconfig.FieldEnabled)
	return u
}

// ClearEnabled clears the value of the "enabled" field.
func (u *GuildConfigUpsert) ClearEnabled() *GuildConfigUpsert {
	u.SetNull(guildconfig.FieldEnabled)
	return u
}

// SetDefaultMaxClones sets the "default_max_clones" field.
func (u *GuildConfigUpsert) SetDefaultMaxClones(v int) *GuildConfigUpsert {
	u.Set(guildconfig.FieldDefaultMaxClones, v)
	return u
}

// UpdateDefaultMaxClones sets the "default_max_clones" field to the value that was provided on create.
func (u *GuildConfigUpsert) UpdateDefaultMaxClones() *GuildConfigUpsert {
	u.SetExcluded(guildconfig.FieldDefaultMaxClones)
	return u
}

// AddDefaultMaxClones adds v to the "default_max_clones" field.
func (u *GuildConfigUpsert) AddDefaultMaxClones(v int) *GuildConfigUpsert {
	u.Add(guildconfig.FieldDefaultMaxClones, v)
	return u
}

// ClearDefaultMaxClones clears the value of the "default_max_clones" field.
func (u *GuildConfigUpsert) ClearDefaultMaxClones() *GuildConfigUpsert {
	u.SetNull(guildconfig.FieldDefaultMaxClones)
	return u
}

// SetRegexMatch sets the "regex_match" field.
func (u *GuildConfigUpsert) SetRegexMatch(v string) *GuildConfigUpsert {
	u.Set(guildconfig.FieldRegexMatch, v)
	return u
}

// UpdateRegexMatch sets the "regex_match" field to the value that was provided on create.
func (u *GuildConfigUpsert) UpdateRegexMatch() *GuildConfigUpsert {
	u.SetExcluded(guildconfig.FieldRegexMatch)
	return u
}

// ClearRegexMatch clears the value of the "regex_match" field.
func (u *GuildConfigUpsert) ClearRegexMatch() *GuildConfigUpsert {
	u.SetNull(guildconfig.FieldRegexMatch)
	return u
}

// SetContactEmail sets the "contact_email" field.
func (u *GuildConfigUpsert) SetContactEmail(v string) *GuildConfigUpsert {
	u.Set(guildconfig.FieldContactEmail, v)
	return u
}

// UpdateContactEmail sets the "contact_email" field to the value that was provided on create.
func (u *GuildConfigUpsert) UpdateContactEmail() *GuildConfigUpsert {
	u.SetExcluded(guildconfig.FieldContactEmail)
	return u
}

// ClearContactEmail clears the value of the "contact_email" field.
func (u *GuildConfigUpsert) ClearContactEmail() *GuildConfigUpsert {
	u.SetNull(guildconfig.FieldContactEmail)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.GuildConfig.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *GuildConfigUpsertOne) UpdateNewValues() *GuildConfigUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.CreateTime(); exists {
			s.SetIgnore(guildconfig.FieldCreateTime)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.GuildConfig.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *GuildConfigUpsertOne) Ignore() *GuildConfigUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *GuildConfigUpsertOne) DoNothing() *GuildConfigUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the GuildConfigCreate.OnConflict
// documentation for more info.
func (u *GuildConfigUpsertOne) Update(set func(*GuildConfigUpsert)) *GuildConfigUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&GuildConfigUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdateTime sets the "update_time" field.
func (u *GuildConfigUpsertOne) SetUpdateTime(v time.Time) *GuildConfigUpsertOne {
	return u.Update(func(s *GuildConfigUpsert) {
		s.SetUpdateTime(v)
	})
}

// UpdateUpdateTime sets the "update_time" field to the value that was provided on create.
func (u *GuildConfigUpsertOne) UpdateUpdateTime() *GuildConfigUpsertOne {
	return u.Update(func(s *GuildConfigUpsert) {
		s.UpdateUpdateTime()
	})
}

// SetEnabled sets the "enabled" field.
func (u *GuildConfigUpsertOne) SetEnabled(v bool) *GuildConfigUpsertOne {
	return u.Update(func(s *GuildConfigUpsert) {
		s.SetEnabled(v)
	})
}

// UpdateEnabled sets the "enabled" field to the value that was provided on create.
func (u *GuildConfigUpsertOne) UpdateEnabled() *GuildConfigUpsertOne {
	return u.Update(func(s *GuildConfigUpsert) {
		s.UpdateEnabled()
	})
}

// ClearEnabled clears the value of the "enabled" field.
func (u *GuildConfigUpsertOne) ClearEnabled() *GuildConfigUpsertOne {
	return u.Update(func(s *GuildConfigUpsert) {
		s.ClearEnabled()
	})
}

// SetDefaultMaxClones sets the "default_max_clones" field.
func (u *GuildConfigUpsertOne) SetDefaultMaxClones(v int) *GuildConfigUpsertOne {
	return u.Update(func(s *GuildConfigUpsert) {
		s.SetDefaultMaxClones(v)
	})
}

// AddDefaultMaxClones adds v to the "default_max_clones" field.
func (u *GuildConfigUpsertOne) AddDefaultMaxClones(v int) *GuildConfigUpsertOne {
	return u.Update(func(s *GuildConfigUpsert) {
		s.AddDefaultMaxClones(v)
	})
}

// UpdateDefaultMaxClones sets the "default_max_clones" field to the value that was provided on create.
func (u *GuildConfigUpsertOne) UpdateDefaultMaxClones() *GuildConfigUpsertOne {
	return u.Update(func(s *GuildConfigUpsert) {
		s.UpdateDefaultMaxClones()
	})
}

// ClearDefaultMaxClones clears the value of the "default_max_clones" field.
func (u *GuildConfigUpsertOne) ClearDefaultMaxClones() *GuildConfigUpsertOne {
	return u.Update(func(s *GuildConfigUpsert) {
		s.ClearDefaultMaxClones()
	})
}

// SetRegexMatch sets the "regex_match" field.
func (u *GuildConfigUpsertOne) SetRegexMatch(v string) *GuildConfigUpsertOne {
	return u.Update(func(s *GuildConfigUpsert) {
		s.SetRegexMatch(v)
	})
}

// UpdateRegexMatch sets the "regex_match" field to the value that was provided on create.
func (u *GuildConfigUpsertOne) UpdateRegexMatch() *GuildConfigUpsertOne {
	return u.Update(func(s *GuildConfigUpsert) {
		s.UpdateRegexMatch()
	})
}

// ClearRegexMatch clears the value of the "regex_match" field.
func (u *GuildConfigUpsertOne) ClearRegexMatch() *GuildConfigUpsertOne {
	return u.Update(func(s *GuildConfigUpsert) {
		s.ClearRegexMatch()
	})
}

// SetContactEmail sets the "contact_email" field.
func (u *GuildConfigUpsertOne) SetContactEmail(v string) *GuildConfigUpsertOne {
	return u.Update(func(s *GuildConfigUpsert) {
		s.SetContactEmail(v)
	})
}

// UpdateContactEmail sets the "contact_email" field to the value that was provided on create.
func (u *GuildConfigUpsertOne) UpdateContactEmail() *GuildConfigUpsertOne {
	return u.Update(func(s *GuildConfigUpsert) {
		s.UpdateContactEmail()
	})
}

// ClearContactEmail clears the value of the "contact_email" field.
func (u *GuildConfigUpsertOne) ClearContactEmail() *GuildConfigUpsertOne {
	return u.Update(func(s *GuildConfigUpsert) {
		s.ClearContactEmail()
	})
}

// Exec executes the query.
func (u *GuildConfigUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for GuildConfigCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *GuildConfigUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *GuildConfigUpsertOne) ID(ctx context.Context) (id int, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *GuildConfigUpsertOne) IDX(ctx context.Context) int {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// GuildConfigCreateBulk is the builder for creating many GuildConfig entities in bulk.
type GuildConfigCreateBulk struct {
	config
	builders []*GuildConfigCreate
	conflict []sql.ConflictOption
}

// Save creates the GuildConfig entities in the database.
func (gccb *GuildConfigCreateBulk) Save(ctx context.Context) ([]*GuildConfig, error) {
	specs := make([]*sqlgraph.CreateSpec, len(gccb.builders))
	nodes := make([]*GuildConfig, len(gccb.builders))
	mutators := make([]Mutator, len(gccb.builders))
	for i := range gccb.builders {
		func(i int, root context.Context) {
			builder := gccb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*GuildConfigMutation)
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
					_, err = mutators[i+1].Mutate(root, gccb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = gccb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, gccb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, gccb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (gccb *GuildConfigCreateBulk) SaveX(ctx context.Context) []*GuildConfig {
	v, err := gccb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (gccb *GuildConfigCreateBulk) Exec(ctx context.Context) error {
	_, err := gccb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (gccb *GuildConfigCreateBulk) ExecX(ctx context.Context) {
	if err := gccb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.GuildConfig.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.GuildConfigUpsert) {
//			SetCreateTime(v+v).
//		}).
//		Exec(ctx)
func (gccb *GuildConfigCreateBulk) OnConflict(opts ...sql.ConflictOption) *GuildConfigUpsertBulk {
	gccb.conflict = opts
	return &GuildConfigUpsertBulk{
		create: gccb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.GuildConfig.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (gccb *GuildConfigCreateBulk) OnConflictColumns(columns ...string) *GuildConfigUpsertBulk {
	gccb.conflict = append(gccb.conflict, sql.ConflictColumns(columns...))
	return &GuildConfigUpsertBulk{
		create: gccb,
	}
}

// GuildConfigUpsertBulk is the builder for "upsert"-ing
// a bulk of GuildConfig nodes.
type GuildConfigUpsertBulk struct {
	create *GuildConfigCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.GuildConfig.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *GuildConfigUpsertBulk) UpdateNewValues() *GuildConfigUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.CreateTime(); exists {
				s.SetIgnore(guildconfig.FieldCreateTime)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.GuildConfig.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *GuildConfigUpsertBulk) Ignore() *GuildConfigUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *GuildConfigUpsertBulk) DoNothing() *GuildConfigUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the GuildConfigCreateBulk.OnConflict
// documentation for more info.
func (u *GuildConfigUpsertBulk) Update(set func(*GuildConfigUpsert)) *GuildConfigUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&GuildConfigUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdateTime sets the "update_time" field.
func (u *GuildConfigUpsertBulk) SetUpdateTime(v time.Time) *GuildConfigUpsertBulk {
	return u.Update(func(s *GuildConfigUpsert) {
		s.SetUpdateTime(v)
	})
}

// UpdateUpdateTime sets the "update_time" field to the value that was provided on create.
func (u *GuildConfigUpsertBulk) UpdateUpdateTime() *GuildConfigUpsertBulk {
	return u.Update(func(s *GuildConfigUpsert) {
		s.UpdateUpdateTime()
	})
}

// SetEnabled sets the "enabled" field.
func (u *GuildConfigUpsertBulk) SetEnabled(v bool) *GuildConfigUpsertBulk {
	return u.Update(func(s *GuildConfigUpsert) {
		s.SetEnabled(v)
	})
}

// UpdateEnabled sets the "enabled" field to the value that was provided on create.
func (u *GuildConfigUpsertBulk) UpdateEnabled() *GuildConfigUpsertBulk {
	return u.Update(func(s *GuildConfigUpsert) {
		s.UpdateEnabled()
	})
}

// ClearEnabled clears the value of the "enabled" field.
func (u *GuildConfigUpsertBulk) ClearEnabled() *GuildConfigUpsertBulk {
	return u.Update(func(s *GuildConfigUpsert) {
		s.ClearEnabled()
	})
}

// SetDefaultMaxClones sets the "default_max_clones" field.
func (u *GuildConfigUpsertBulk) SetDefaultMaxClones(v int) *GuildConfigUpsertBulk {
	return u.Update(func(s *GuildConfigUpsert) {
		s.SetDefaultMaxClones(v)
	})
}

// AddDefaultMaxClones adds v to the "default_max_clones" field.
func (u *GuildConfigUpsertBulk) AddDefaultMaxClones(v int) *GuildConfigUpsertBulk {
	return u.Update(func(s *GuildConfigUpsert) {
		s.AddDefaultMaxClones(v)
	})
}

// UpdateDefaultMaxClones sets the "default_max_clones" field to the value that was provided on create.
func (u *GuildConfigUpsertBulk) UpdateDefaultMaxClones() *GuildConfigUpsertBulk {
	return u.Update(func(s *GuildConfigUpsert) {
		s.UpdateDefaultMaxClones()
	})
}

// ClearDefaultMaxClones clears the value of the "default_max_clones" field.
func (u *GuildConfigUpsertBulk) ClearDefaultMaxClones() *GuildConfigUpsertBulk {
	return u.Update(func(s *GuildConfigUpsert) {
		s.ClearDefaultMaxClones()
	})
}

// SetRegexMatch sets the "regex_match" field.
func (u *GuildConfigUpsertBulk) SetRegexMatch(v string) *GuildConfigUpsertBulk {
	return u.Update(func(s *GuildConfigUpsert) {
		s.SetRegexMatch(v)
	})
}

// UpdateRegexMatch sets the "regex_match" field to the value that was provided on create.
func (u *GuildConfigUpsertBulk) UpdateRegexMatch() *GuildConfigUpsertBulk {
	return u.Update(func(s *GuildConfigUpsert) {
		s.UpdateRegexMatch()
	})
}

// ClearRegexMatch clears the value of the "regex_match" field.
func (u *GuildConfigUpsertBulk) ClearRegexMatch() *GuildConfigUpsertBulk {
	return u.Update(func(s *GuildConfigUpsert) {
		s.ClearRegexMatch()
	})
}

// SetContactEmail sets the "contact_email" field.
func (u *GuildConfigUpsertBulk) SetContactEmail(v string) *GuildConfigUpsertBulk {
	return u.Update(func(s *GuildConfigUpsert) {
		s.SetContactEmail(v)
	})
}

// UpdateContactEmail sets the "contact_email" field to the value that was provided on create.
func (u *GuildConfigUpsertBulk) UpdateContactEmail() *GuildConfigUpsertBulk {
	return u.Update(func(s *GuildConfigUpsert) {
		s.UpdateContactEmail()
	})
}

// ClearContactEmail clears the value of the "contact_email" field.
func (u *GuildConfigUpsertBulk) ClearContactEmail() *GuildConfigUpsertBulk {
	return u.Update(func(s *GuildConfigUpsert) {
		s.ClearContactEmail()
	})
}

// Exec executes the query.
func (u *GuildConfigUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the GuildConfigCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for GuildConfigCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *GuildConfigUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
