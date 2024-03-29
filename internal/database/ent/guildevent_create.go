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
	"github.com/lrstanley/spectrograph/internal/database/ent/guildevent"
)

// GuildEventCreate is the builder for creating a GuildEvent entity.
type GuildEventCreate struct {
	config
	mutation *GuildEventMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetCreateTime sets the "create_time" field.
func (gec *GuildEventCreate) SetCreateTime(t time.Time) *GuildEventCreate {
	gec.mutation.SetCreateTime(t)
	return gec
}

// SetNillableCreateTime sets the "create_time" field if the given value is not nil.
func (gec *GuildEventCreate) SetNillableCreateTime(t *time.Time) *GuildEventCreate {
	if t != nil {
		gec.SetCreateTime(*t)
	}
	return gec
}

// SetUpdateTime sets the "update_time" field.
func (gec *GuildEventCreate) SetUpdateTime(t time.Time) *GuildEventCreate {
	gec.mutation.SetUpdateTime(t)
	return gec
}

// SetNillableUpdateTime sets the "update_time" field if the given value is not nil.
func (gec *GuildEventCreate) SetNillableUpdateTime(t *time.Time) *GuildEventCreate {
	if t != nil {
		gec.SetUpdateTime(*t)
	}
	return gec
}

// SetType sets the "type" field.
func (gec *GuildEventCreate) SetType(gu guildevent.Type) *GuildEventCreate {
	gec.mutation.SetType(gu)
	return gec
}

// SetMessage sets the "message" field.
func (gec *GuildEventCreate) SetMessage(s string) *GuildEventCreate {
	gec.mutation.SetMessage(s)
	return gec
}

// SetMetadata sets the "metadata" field.
func (gec *GuildEventCreate) SetMetadata(m map[string]interface{}) *GuildEventCreate {
	gec.mutation.SetMetadata(m)
	return gec
}

// SetGuildID sets the "guild" edge to the Guild entity by ID.
func (gec *GuildEventCreate) SetGuildID(id int) *GuildEventCreate {
	gec.mutation.SetGuildID(id)
	return gec
}

// SetGuild sets the "guild" edge to the Guild entity.
func (gec *GuildEventCreate) SetGuild(g *Guild) *GuildEventCreate {
	return gec.SetGuildID(g.ID)
}

// Mutation returns the GuildEventMutation object of the builder.
func (gec *GuildEventCreate) Mutation() *GuildEventMutation {
	return gec.mutation
}

// Save creates the GuildEvent in the database.
func (gec *GuildEventCreate) Save(ctx context.Context) (*GuildEvent, error) {
	if err := gec.defaults(); err != nil {
		return nil, err
	}
	return withHooks[*GuildEvent, GuildEventMutation](ctx, gec.sqlSave, gec.mutation, gec.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (gec *GuildEventCreate) SaveX(ctx context.Context) *GuildEvent {
	v, err := gec.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (gec *GuildEventCreate) Exec(ctx context.Context) error {
	_, err := gec.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (gec *GuildEventCreate) ExecX(ctx context.Context) {
	if err := gec.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (gec *GuildEventCreate) defaults() error {
	if _, ok := gec.mutation.CreateTime(); !ok {
		if guildevent.DefaultCreateTime == nil {
			return fmt.Errorf("ent: uninitialized guildevent.DefaultCreateTime (forgotten import ent/runtime?)")
		}
		v := guildevent.DefaultCreateTime()
		gec.mutation.SetCreateTime(v)
	}
	if _, ok := gec.mutation.UpdateTime(); !ok {
		if guildevent.DefaultUpdateTime == nil {
			return fmt.Errorf("ent: uninitialized guildevent.DefaultUpdateTime (forgotten import ent/runtime?)")
		}
		v := guildevent.DefaultUpdateTime()
		gec.mutation.SetUpdateTime(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (gec *GuildEventCreate) check() error {
	if _, ok := gec.mutation.CreateTime(); !ok {
		return &ValidationError{Name: "create_time", err: errors.New(`ent: missing required field "GuildEvent.create_time"`)}
	}
	if _, ok := gec.mutation.UpdateTime(); !ok {
		return &ValidationError{Name: "update_time", err: errors.New(`ent: missing required field "GuildEvent.update_time"`)}
	}
	if _, ok := gec.mutation.GetType(); !ok {
		return &ValidationError{Name: "type", err: errors.New(`ent: missing required field "GuildEvent.type"`)}
	}
	if v, ok := gec.mutation.GetType(); ok {
		if err := guildevent.TypeValidator(v); err != nil {
			return &ValidationError{Name: "type", err: fmt.Errorf(`ent: validator failed for field "GuildEvent.type": %w`, err)}
		}
	}
	if _, ok := gec.mutation.Message(); !ok {
		return &ValidationError{Name: "message", err: errors.New(`ent: missing required field "GuildEvent.message"`)}
	}
	if _, ok := gec.mutation.GuildID(); !ok {
		return &ValidationError{Name: "guild", err: errors.New(`ent: missing required edge "GuildEvent.guild"`)}
	}
	return nil
}

func (gec *GuildEventCreate) sqlSave(ctx context.Context) (*GuildEvent, error) {
	if err := gec.check(); err != nil {
		return nil, err
	}
	_node, _spec := gec.createSpec()
	if err := sqlgraph.CreateNode(ctx, gec.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	gec.mutation.id = &_node.ID
	gec.mutation.done = true
	return _node, nil
}

func (gec *GuildEventCreate) createSpec() (*GuildEvent, *sqlgraph.CreateSpec) {
	var (
		_node = &GuildEvent{config: gec.config}
		_spec = sqlgraph.NewCreateSpec(guildevent.Table, sqlgraph.NewFieldSpec(guildevent.FieldID, field.TypeInt))
	)
	_spec.OnConflict = gec.conflict
	if value, ok := gec.mutation.CreateTime(); ok {
		_spec.SetField(guildevent.FieldCreateTime, field.TypeTime, value)
		_node.CreateTime = value
	}
	if value, ok := gec.mutation.UpdateTime(); ok {
		_spec.SetField(guildevent.FieldUpdateTime, field.TypeTime, value)
		_node.UpdateTime = value
	}
	if value, ok := gec.mutation.GetType(); ok {
		_spec.SetField(guildevent.FieldType, field.TypeEnum, value)
		_node.Type = value
	}
	if value, ok := gec.mutation.Message(); ok {
		_spec.SetField(guildevent.FieldMessage, field.TypeString, value)
		_node.Message = value
	}
	if value, ok := gec.mutation.Metadata(); ok {
		_spec.SetField(guildevent.FieldMetadata, field.TypeJSON, value)
		_node.Metadata = value
	}
	if nodes := gec.mutation.GuildIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   guildevent.GuildTable,
			Columns: []string{guildevent.GuildColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(guild.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.guild_guild_events = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.GuildEvent.Create().
//		SetCreateTime(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.GuildEventUpsert) {
//			SetCreateTime(v+v).
//		}).
//		Exec(ctx)
func (gec *GuildEventCreate) OnConflict(opts ...sql.ConflictOption) *GuildEventUpsertOne {
	gec.conflict = opts
	return &GuildEventUpsertOne{
		create: gec,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.GuildEvent.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (gec *GuildEventCreate) OnConflictColumns(columns ...string) *GuildEventUpsertOne {
	gec.conflict = append(gec.conflict, sql.ConflictColumns(columns...))
	return &GuildEventUpsertOne{
		create: gec,
	}
}

type (
	// GuildEventUpsertOne is the builder for "upsert"-ing
	//  one GuildEvent node.
	GuildEventUpsertOne struct {
		create *GuildEventCreate
	}

	// GuildEventUpsert is the "OnConflict" setter.
	GuildEventUpsert struct {
		*sql.UpdateSet
	}
)

// SetUpdateTime sets the "update_time" field.
func (u *GuildEventUpsert) SetUpdateTime(v time.Time) *GuildEventUpsert {
	u.Set(guildevent.FieldUpdateTime, v)
	return u
}

// UpdateUpdateTime sets the "update_time" field to the value that was provided on create.
func (u *GuildEventUpsert) UpdateUpdateTime() *GuildEventUpsert {
	u.SetExcluded(guildevent.FieldUpdateTime)
	return u
}

// SetType sets the "type" field.
func (u *GuildEventUpsert) SetType(v guildevent.Type) *GuildEventUpsert {
	u.Set(guildevent.FieldType, v)
	return u
}

// UpdateType sets the "type" field to the value that was provided on create.
func (u *GuildEventUpsert) UpdateType() *GuildEventUpsert {
	u.SetExcluded(guildevent.FieldType)
	return u
}

// SetMessage sets the "message" field.
func (u *GuildEventUpsert) SetMessage(v string) *GuildEventUpsert {
	u.Set(guildevent.FieldMessage, v)
	return u
}

// UpdateMessage sets the "message" field to the value that was provided on create.
func (u *GuildEventUpsert) UpdateMessage() *GuildEventUpsert {
	u.SetExcluded(guildevent.FieldMessage)
	return u
}

// SetMetadata sets the "metadata" field.
func (u *GuildEventUpsert) SetMetadata(v map[string]interface{}) *GuildEventUpsert {
	u.Set(guildevent.FieldMetadata, v)
	return u
}

// UpdateMetadata sets the "metadata" field to the value that was provided on create.
func (u *GuildEventUpsert) UpdateMetadata() *GuildEventUpsert {
	u.SetExcluded(guildevent.FieldMetadata)
	return u
}

// ClearMetadata clears the value of the "metadata" field.
func (u *GuildEventUpsert) ClearMetadata() *GuildEventUpsert {
	u.SetNull(guildevent.FieldMetadata)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.GuildEvent.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *GuildEventUpsertOne) UpdateNewValues() *GuildEventUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.CreateTime(); exists {
			s.SetIgnore(guildevent.FieldCreateTime)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.GuildEvent.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *GuildEventUpsertOne) Ignore() *GuildEventUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *GuildEventUpsertOne) DoNothing() *GuildEventUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the GuildEventCreate.OnConflict
// documentation for more info.
func (u *GuildEventUpsertOne) Update(set func(*GuildEventUpsert)) *GuildEventUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&GuildEventUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdateTime sets the "update_time" field.
func (u *GuildEventUpsertOne) SetUpdateTime(v time.Time) *GuildEventUpsertOne {
	return u.Update(func(s *GuildEventUpsert) {
		s.SetUpdateTime(v)
	})
}

// UpdateUpdateTime sets the "update_time" field to the value that was provided on create.
func (u *GuildEventUpsertOne) UpdateUpdateTime() *GuildEventUpsertOne {
	return u.Update(func(s *GuildEventUpsert) {
		s.UpdateUpdateTime()
	})
}

// SetType sets the "type" field.
func (u *GuildEventUpsertOne) SetType(v guildevent.Type) *GuildEventUpsertOne {
	return u.Update(func(s *GuildEventUpsert) {
		s.SetType(v)
	})
}

// UpdateType sets the "type" field to the value that was provided on create.
func (u *GuildEventUpsertOne) UpdateType() *GuildEventUpsertOne {
	return u.Update(func(s *GuildEventUpsert) {
		s.UpdateType()
	})
}

// SetMessage sets the "message" field.
func (u *GuildEventUpsertOne) SetMessage(v string) *GuildEventUpsertOne {
	return u.Update(func(s *GuildEventUpsert) {
		s.SetMessage(v)
	})
}

// UpdateMessage sets the "message" field to the value that was provided on create.
func (u *GuildEventUpsertOne) UpdateMessage() *GuildEventUpsertOne {
	return u.Update(func(s *GuildEventUpsert) {
		s.UpdateMessage()
	})
}

// SetMetadata sets the "metadata" field.
func (u *GuildEventUpsertOne) SetMetadata(v map[string]interface{}) *GuildEventUpsertOne {
	return u.Update(func(s *GuildEventUpsert) {
		s.SetMetadata(v)
	})
}

// UpdateMetadata sets the "metadata" field to the value that was provided on create.
func (u *GuildEventUpsertOne) UpdateMetadata() *GuildEventUpsertOne {
	return u.Update(func(s *GuildEventUpsert) {
		s.UpdateMetadata()
	})
}

// ClearMetadata clears the value of the "metadata" field.
func (u *GuildEventUpsertOne) ClearMetadata() *GuildEventUpsertOne {
	return u.Update(func(s *GuildEventUpsert) {
		s.ClearMetadata()
	})
}

// Exec executes the query.
func (u *GuildEventUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for GuildEventCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *GuildEventUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *GuildEventUpsertOne) ID(ctx context.Context) (id int, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *GuildEventUpsertOne) IDX(ctx context.Context) int {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// GuildEventCreateBulk is the builder for creating many GuildEvent entities in bulk.
type GuildEventCreateBulk struct {
	config
	builders []*GuildEventCreate
	conflict []sql.ConflictOption
}

// Save creates the GuildEvent entities in the database.
func (gecb *GuildEventCreateBulk) Save(ctx context.Context) ([]*GuildEvent, error) {
	specs := make([]*sqlgraph.CreateSpec, len(gecb.builders))
	nodes := make([]*GuildEvent, len(gecb.builders))
	mutators := make([]Mutator, len(gecb.builders))
	for i := range gecb.builders {
		func(i int, root context.Context) {
			builder := gecb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*GuildEventMutation)
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
					_, err = mutators[i+1].Mutate(root, gecb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = gecb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, gecb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, gecb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (gecb *GuildEventCreateBulk) SaveX(ctx context.Context) []*GuildEvent {
	v, err := gecb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (gecb *GuildEventCreateBulk) Exec(ctx context.Context) error {
	_, err := gecb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (gecb *GuildEventCreateBulk) ExecX(ctx context.Context) {
	if err := gecb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.GuildEvent.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.GuildEventUpsert) {
//			SetCreateTime(v+v).
//		}).
//		Exec(ctx)
func (gecb *GuildEventCreateBulk) OnConflict(opts ...sql.ConflictOption) *GuildEventUpsertBulk {
	gecb.conflict = opts
	return &GuildEventUpsertBulk{
		create: gecb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.GuildEvent.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (gecb *GuildEventCreateBulk) OnConflictColumns(columns ...string) *GuildEventUpsertBulk {
	gecb.conflict = append(gecb.conflict, sql.ConflictColumns(columns...))
	return &GuildEventUpsertBulk{
		create: gecb,
	}
}

// GuildEventUpsertBulk is the builder for "upsert"-ing
// a bulk of GuildEvent nodes.
type GuildEventUpsertBulk struct {
	create *GuildEventCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.GuildEvent.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *GuildEventUpsertBulk) UpdateNewValues() *GuildEventUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.CreateTime(); exists {
				s.SetIgnore(guildevent.FieldCreateTime)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.GuildEvent.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *GuildEventUpsertBulk) Ignore() *GuildEventUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *GuildEventUpsertBulk) DoNothing() *GuildEventUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the GuildEventCreateBulk.OnConflict
// documentation for more info.
func (u *GuildEventUpsertBulk) Update(set func(*GuildEventUpsert)) *GuildEventUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&GuildEventUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdateTime sets the "update_time" field.
func (u *GuildEventUpsertBulk) SetUpdateTime(v time.Time) *GuildEventUpsertBulk {
	return u.Update(func(s *GuildEventUpsert) {
		s.SetUpdateTime(v)
	})
}

// UpdateUpdateTime sets the "update_time" field to the value that was provided on create.
func (u *GuildEventUpsertBulk) UpdateUpdateTime() *GuildEventUpsertBulk {
	return u.Update(func(s *GuildEventUpsert) {
		s.UpdateUpdateTime()
	})
}

// SetType sets the "type" field.
func (u *GuildEventUpsertBulk) SetType(v guildevent.Type) *GuildEventUpsertBulk {
	return u.Update(func(s *GuildEventUpsert) {
		s.SetType(v)
	})
}

// UpdateType sets the "type" field to the value that was provided on create.
func (u *GuildEventUpsertBulk) UpdateType() *GuildEventUpsertBulk {
	return u.Update(func(s *GuildEventUpsert) {
		s.UpdateType()
	})
}

// SetMessage sets the "message" field.
func (u *GuildEventUpsertBulk) SetMessage(v string) *GuildEventUpsertBulk {
	return u.Update(func(s *GuildEventUpsert) {
		s.SetMessage(v)
	})
}

// UpdateMessage sets the "message" field to the value that was provided on create.
func (u *GuildEventUpsertBulk) UpdateMessage() *GuildEventUpsertBulk {
	return u.Update(func(s *GuildEventUpsert) {
		s.UpdateMessage()
	})
}

// SetMetadata sets the "metadata" field.
func (u *GuildEventUpsertBulk) SetMetadata(v map[string]interface{}) *GuildEventUpsertBulk {
	return u.Update(func(s *GuildEventUpsert) {
		s.SetMetadata(v)
	})
}

// UpdateMetadata sets the "metadata" field to the value that was provided on create.
func (u *GuildEventUpsertBulk) UpdateMetadata() *GuildEventUpsertBulk {
	return u.Update(func(s *GuildEventUpsert) {
		s.UpdateMetadata()
	})
}

// ClearMetadata clears the value of the "metadata" field.
func (u *GuildEventUpsertBulk) ClearMetadata() *GuildEventUpsertBulk {
	return u.Update(func(s *GuildEventUpsert) {
		s.ClearMetadata()
	})
}

// Exec executes the query.
func (u *GuildEventUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the GuildEventCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for GuildEventCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *GuildEventUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
