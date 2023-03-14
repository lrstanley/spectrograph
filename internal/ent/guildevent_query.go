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
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/lrstanley/spectrograph/internal/ent/guild"
	"github.com/lrstanley/spectrograph/internal/ent/guildevent"
	"github.com/lrstanley/spectrograph/internal/ent/predicate"
)

// GuildEventQuery is the builder for querying GuildEvent entities.
type GuildEventQuery struct {
	config
	ctx        *QueryContext
	order      []OrderFunc
	inters     []Interceptor
	predicates []predicate.GuildEvent
	withGuild  *GuildQuery
	withFKs    bool
	modifiers  []func(*sql.Selector)
	loadTotal  []func(context.Context, []*GuildEvent) error
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the GuildEventQuery builder.
func (geq *GuildEventQuery) Where(ps ...predicate.GuildEvent) *GuildEventQuery {
	geq.predicates = append(geq.predicates, ps...)
	return geq
}

// Limit the number of records to be returned by this query.
func (geq *GuildEventQuery) Limit(limit int) *GuildEventQuery {
	geq.ctx.Limit = &limit
	return geq
}

// Offset to start from.
func (geq *GuildEventQuery) Offset(offset int) *GuildEventQuery {
	geq.ctx.Offset = &offset
	return geq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (geq *GuildEventQuery) Unique(unique bool) *GuildEventQuery {
	geq.ctx.Unique = &unique
	return geq
}

// Order specifies how the records should be ordered.
func (geq *GuildEventQuery) Order(o ...OrderFunc) *GuildEventQuery {
	geq.order = append(geq.order, o...)
	return geq
}

// QueryGuild chains the current query on the "guild" edge.
func (geq *GuildEventQuery) QueryGuild() *GuildQuery {
	query := (&GuildClient{config: geq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := geq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := geq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(guildevent.Table, guildevent.FieldID, selector),
			sqlgraph.To(guild.Table, guild.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, guildevent.GuildTable, guildevent.GuildColumn),
		)
		fromU = sqlgraph.SetNeighbors(geq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first GuildEvent entity from the query.
// Returns a *NotFoundError when no GuildEvent was found.
func (geq *GuildEventQuery) First(ctx context.Context) (*GuildEvent, error) {
	nodes, err := geq.Limit(1).All(setContextOp(ctx, geq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{guildevent.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (geq *GuildEventQuery) FirstX(ctx context.Context) *GuildEvent {
	node, err := geq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first GuildEvent ID from the query.
// Returns a *NotFoundError when no GuildEvent ID was found.
func (geq *GuildEventQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = geq.Limit(1).IDs(setContextOp(ctx, geq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{guildevent.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (geq *GuildEventQuery) FirstIDX(ctx context.Context) int {
	id, err := geq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single GuildEvent entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one GuildEvent entity is found.
// Returns a *NotFoundError when no GuildEvent entities are found.
func (geq *GuildEventQuery) Only(ctx context.Context) (*GuildEvent, error) {
	nodes, err := geq.Limit(2).All(setContextOp(ctx, geq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{guildevent.Label}
	default:
		return nil, &NotSingularError{guildevent.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (geq *GuildEventQuery) OnlyX(ctx context.Context) *GuildEvent {
	node, err := geq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only GuildEvent ID in the query.
// Returns a *NotSingularError when more than one GuildEvent ID is found.
// Returns a *NotFoundError when no entities are found.
func (geq *GuildEventQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = geq.Limit(2).IDs(setContextOp(ctx, geq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{guildevent.Label}
	default:
		err = &NotSingularError{guildevent.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (geq *GuildEventQuery) OnlyIDX(ctx context.Context) int {
	id, err := geq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of GuildEvents.
func (geq *GuildEventQuery) All(ctx context.Context) ([]*GuildEvent, error) {
	ctx = setContextOp(ctx, geq.ctx, "All")
	if err := geq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*GuildEvent, *GuildEventQuery]()
	return withInterceptors[[]*GuildEvent](ctx, geq, qr, geq.inters)
}

// AllX is like All, but panics if an error occurs.
func (geq *GuildEventQuery) AllX(ctx context.Context) []*GuildEvent {
	nodes, err := geq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of GuildEvent IDs.
func (geq *GuildEventQuery) IDs(ctx context.Context) (ids []int, err error) {
	if geq.ctx.Unique == nil && geq.path != nil {
		geq.Unique(true)
	}
	ctx = setContextOp(ctx, geq.ctx, "IDs")
	if err = geq.Select(guildevent.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (geq *GuildEventQuery) IDsX(ctx context.Context) []int {
	ids, err := geq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (geq *GuildEventQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, geq.ctx, "Count")
	if err := geq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, geq, querierCount[*GuildEventQuery](), geq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (geq *GuildEventQuery) CountX(ctx context.Context) int {
	count, err := geq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (geq *GuildEventQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, geq.ctx, "Exist")
	switch _, err := geq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (geq *GuildEventQuery) ExistX(ctx context.Context) bool {
	exist, err := geq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the GuildEventQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (geq *GuildEventQuery) Clone() *GuildEventQuery {
	if geq == nil {
		return nil
	}
	return &GuildEventQuery{
		config:     geq.config,
		ctx:        geq.ctx.Clone(),
		order:      append([]OrderFunc{}, geq.order...),
		inters:     append([]Interceptor{}, geq.inters...),
		predicates: append([]predicate.GuildEvent{}, geq.predicates...),
		withGuild:  geq.withGuild.Clone(),
		// clone intermediate query.
		sql:  geq.sql.Clone(),
		path: geq.path,
	}
}

// WithGuild tells the query-builder to eager-load the nodes that are connected to
// the "guild" edge. The optional arguments are used to configure the query builder of the edge.
func (geq *GuildEventQuery) WithGuild(opts ...func(*GuildQuery)) *GuildEventQuery {
	query := (&GuildClient{config: geq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	geq.withGuild = query
	return geq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		CreateTime time.Time `json:"create_time,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.GuildEvent.Query().
//		GroupBy(guildevent.FieldCreateTime).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (geq *GuildEventQuery) GroupBy(field string, fields ...string) *GuildEventGroupBy {
	geq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &GuildEventGroupBy{build: geq}
	grbuild.flds = &geq.ctx.Fields
	grbuild.label = guildevent.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		CreateTime time.Time `json:"create_time,omitempty"`
//	}
//
//	client.GuildEvent.Query().
//		Select(guildevent.FieldCreateTime).
//		Scan(ctx, &v)
func (geq *GuildEventQuery) Select(fields ...string) *GuildEventSelect {
	geq.ctx.Fields = append(geq.ctx.Fields, fields...)
	sbuild := &GuildEventSelect{GuildEventQuery: geq}
	sbuild.label = guildevent.Label
	sbuild.flds, sbuild.scan = &geq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a GuildEventSelect configured with the given aggregations.
func (geq *GuildEventQuery) Aggregate(fns ...AggregateFunc) *GuildEventSelect {
	return geq.Select().Aggregate(fns...)
}

func (geq *GuildEventQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range geq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, geq); err != nil {
				return err
			}
		}
	}
	for _, f := range geq.ctx.Fields {
		if !guildevent.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if geq.path != nil {
		prev, err := geq.path(ctx)
		if err != nil {
			return err
		}
		geq.sql = prev
	}
	if guildevent.Policy == nil {
		return errors.New("ent: uninitialized guildevent.Policy (forgotten import ent/runtime?)")
	}
	if err := guildevent.Policy.EvalQuery(ctx, geq); err != nil {
		return err
	}
	return nil
}

func (geq *GuildEventQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*GuildEvent, error) {
	var (
		nodes       = []*GuildEvent{}
		withFKs     = geq.withFKs
		_spec       = geq.querySpec()
		loadedTypes = [1]bool{
			geq.withGuild != nil,
		}
	)
	if geq.withGuild != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, guildevent.ForeignKeys...)
	}
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*GuildEvent).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &GuildEvent{config: geq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if len(geq.modifiers) > 0 {
		_spec.Modifiers = geq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, geq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := geq.withGuild; query != nil {
		if err := geq.loadGuild(ctx, query, nodes, nil,
			func(n *GuildEvent, e *Guild) { n.Edges.Guild = e }); err != nil {
			return nil, err
		}
	}
	for i := range geq.loadTotal {
		if err := geq.loadTotal[i](ctx, nodes); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (geq *GuildEventQuery) loadGuild(ctx context.Context, query *GuildQuery, nodes []*GuildEvent, init func(*GuildEvent), assign func(*GuildEvent, *Guild)) error {
	ids := make([]int, 0, len(nodes))
	nodeids := make(map[int][]*GuildEvent)
	for i := range nodes {
		if nodes[i].guild_guild_events == nil {
			continue
		}
		fk := *nodes[i].guild_guild_events
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(guild.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "guild_guild_events" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}

func (geq *GuildEventQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := geq.querySpec()
	if len(geq.modifiers) > 0 {
		_spec.Modifiers = geq.modifiers
	}
	_spec.Node.Columns = geq.ctx.Fields
	if len(geq.ctx.Fields) > 0 {
		_spec.Unique = geq.ctx.Unique != nil && *geq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, geq.driver, _spec)
}

func (geq *GuildEventQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(guildevent.Table, guildevent.Columns, sqlgraph.NewFieldSpec(guildevent.FieldID, field.TypeInt))
	_spec.From = geq.sql
	if unique := geq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if geq.path != nil {
		_spec.Unique = true
	}
	if fields := geq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, guildevent.FieldID)
		for i := range fields {
			if fields[i] != guildevent.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := geq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := geq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := geq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := geq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (geq *GuildEventQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(geq.driver.Dialect())
	t1 := builder.Table(guildevent.Table)
	columns := geq.ctx.Fields
	if len(columns) == 0 {
		columns = guildevent.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if geq.sql != nil {
		selector = geq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if geq.ctx.Unique != nil && *geq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range geq.predicates {
		p(selector)
	}
	for _, p := range geq.order {
		p(selector)
	}
	if offset := geq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := geq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// GuildEventGroupBy is the group-by builder for GuildEvent entities.
type GuildEventGroupBy struct {
	selector
	build *GuildEventQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (gegb *GuildEventGroupBy) Aggregate(fns ...AggregateFunc) *GuildEventGroupBy {
	gegb.fns = append(gegb.fns, fns...)
	return gegb
}

// Scan applies the selector query and scans the result into the given value.
func (gegb *GuildEventGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, gegb.build.ctx, "GroupBy")
	if err := gegb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*GuildEventQuery, *GuildEventGroupBy](ctx, gegb.build, gegb, gegb.build.inters, v)
}

func (gegb *GuildEventGroupBy) sqlScan(ctx context.Context, root *GuildEventQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(gegb.fns))
	for _, fn := range gegb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*gegb.flds)+len(gegb.fns))
		for _, f := range *gegb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*gegb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := gegb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// GuildEventSelect is the builder for selecting fields of GuildEvent entities.
type GuildEventSelect struct {
	*GuildEventQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (ges *GuildEventSelect) Aggregate(fns ...AggregateFunc) *GuildEventSelect {
	ges.fns = append(ges.fns, fns...)
	return ges
}

// Scan applies the selector query and scans the result into the given value.
func (ges *GuildEventSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, ges.ctx, "Select")
	if err := ges.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*GuildEventQuery, *GuildEventSelect](ctx, ges.GuildEventQuery, ges, ges.inters, v)
}

func (ges *GuildEventSelect) sqlScan(ctx context.Context, root *GuildEventQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(ges.fns))
	for _, fn := range ges.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*ges.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := ges.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
