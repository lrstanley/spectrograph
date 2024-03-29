// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.
//
// DO NOT EDIT, CODE GENERATED BY entc.

package ent

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/lrstanley/spectrograph/internal/database/ent/guild"
	"github.com/lrstanley/spectrograph/internal/database/ent/guildadminconfig"
	"github.com/lrstanley/spectrograph/internal/database/ent/guildconfig"
	"github.com/lrstanley/spectrograph/internal/database/ent/guildevent"
	"github.com/lrstanley/spectrograph/internal/database/ent/predicate"
	"github.com/lrstanley/spectrograph/internal/database/ent/user"
)

// GuildQuery is the builder for querying Guild entities.
type GuildQuery struct {
	config
	ctx                  *QueryContext
	order                []guild.OrderOption
	inters               []Interceptor
	predicates           []predicate.Guild
	withGuildConfig      *GuildConfigQuery
	withGuildAdminConfig *GuildAdminConfigQuery
	withGuildEvents      *GuildEventQuery
	withAdmins           *UserQuery
	modifiers            []func(*sql.Selector)
	loadTotal            []func(context.Context, []*Guild) error
	withNamedGuildEvents map[string]*GuildEventQuery
	withNamedAdmins      map[string]*UserQuery
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the GuildQuery builder.
func (gq *GuildQuery) Where(ps ...predicate.Guild) *GuildQuery {
	gq.predicates = append(gq.predicates, ps...)
	return gq
}

// Limit the number of records to be returned by this query.
func (gq *GuildQuery) Limit(limit int) *GuildQuery {
	gq.ctx.Limit = &limit
	return gq
}

// Offset to start from.
func (gq *GuildQuery) Offset(offset int) *GuildQuery {
	gq.ctx.Offset = &offset
	return gq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (gq *GuildQuery) Unique(unique bool) *GuildQuery {
	gq.ctx.Unique = &unique
	return gq
}

// Order specifies how the records should be ordered.
func (gq *GuildQuery) Order(o ...guild.OrderOption) *GuildQuery {
	gq.order = append(gq.order, o...)
	return gq
}

// QueryGuildConfig chains the current query on the "guild_config" edge.
func (gq *GuildQuery) QueryGuildConfig() *GuildConfigQuery {
	query := (&GuildConfigClient{config: gq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := gq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := gq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(guild.Table, guild.FieldID, selector),
			sqlgraph.To(guildconfig.Table, guildconfig.FieldID),
			sqlgraph.Edge(sqlgraph.O2O, false, guild.GuildConfigTable, guild.GuildConfigColumn),
		)
		fromU = sqlgraph.SetNeighbors(gq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryGuildAdminConfig chains the current query on the "guild_admin_config" edge.
func (gq *GuildQuery) QueryGuildAdminConfig() *GuildAdminConfigQuery {
	query := (&GuildAdminConfigClient{config: gq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := gq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := gq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(guild.Table, guild.FieldID, selector),
			sqlgraph.To(guildadminconfig.Table, guildadminconfig.FieldID),
			sqlgraph.Edge(sqlgraph.O2O, false, guild.GuildAdminConfigTable, guild.GuildAdminConfigColumn),
		)
		fromU = sqlgraph.SetNeighbors(gq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryGuildEvents chains the current query on the "guild_events" edge.
func (gq *GuildQuery) QueryGuildEvents() *GuildEventQuery {
	query := (&GuildEventClient{config: gq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := gq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := gq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(guild.Table, guild.FieldID, selector),
			sqlgraph.To(guildevent.Table, guildevent.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, guild.GuildEventsTable, guild.GuildEventsColumn),
		)
		fromU = sqlgraph.SetNeighbors(gq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryAdmins chains the current query on the "admins" edge.
func (gq *GuildQuery) QueryAdmins() *UserQuery {
	query := (&UserClient{config: gq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := gq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := gq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(guild.Table, guild.FieldID, selector),
			sqlgraph.To(user.Table, user.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, guild.AdminsTable, guild.AdminsPrimaryKey...),
		)
		fromU = sqlgraph.SetNeighbors(gq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first Guild entity from the query.
// Returns a *NotFoundError when no Guild was found.
func (gq *GuildQuery) First(ctx context.Context) (*Guild, error) {
	nodes, err := gq.Limit(1).All(setContextOp(ctx, gq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{guild.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (gq *GuildQuery) FirstX(ctx context.Context) *Guild {
	node, err := gq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Guild ID from the query.
// Returns a *NotFoundError when no Guild ID was found.
func (gq *GuildQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = gq.Limit(1).IDs(setContextOp(ctx, gq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{guild.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (gq *GuildQuery) FirstIDX(ctx context.Context) int {
	id, err := gq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Guild entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Guild entity is found.
// Returns a *NotFoundError when no Guild entities are found.
func (gq *GuildQuery) Only(ctx context.Context) (*Guild, error) {
	nodes, err := gq.Limit(2).All(setContextOp(ctx, gq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{guild.Label}
	default:
		return nil, &NotSingularError{guild.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (gq *GuildQuery) OnlyX(ctx context.Context) *Guild {
	node, err := gq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Guild ID in the query.
// Returns a *NotSingularError when more than one Guild ID is found.
// Returns a *NotFoundError when no entities are found.
func (gq *GuildQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = gq.Limit(2).IDs(setContextOp(ctx, gq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{guild.Label}
	default:
		err = &NotSingularError{guild.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (gq *GuildQuery) OnlyIDX(ctx context.Context) int {
	id, err := gq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Guilds.
func (gq *GuildQuery) All(ctx context.Context) ([]*Guild, error) {
	ctx = setContextOp(ctx, gq.ctx, "All")
	if err := gq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*Guild, *GuildQuery]()
	return withInterceptors[[]*Guild](ctx, gq, qr, gq.inters)
}

// AllX is like All, but panics if an error occurs.
func (gq *GuildQuery) AllX(ctx context.Context) []*Guild {
	nodes, err := gq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Guild IDs.
func (gq *GuildQuery) IDs(ctx context.Context) (ids []int, err error) {
	if gq.ctx.Unique == nil && gq.path != nil {
		gq.Unique(true)
	}
	ctx = setContextOp(ctx, gq.ctx, "IDs")
	if err = gq.Select(guild.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (gq *GuildQuery) IDsX(ctx context.Context) []int {
	ids, err := gq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (gq *GuildQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, gq.ctx, "Count")
	if err := gq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, gq, querierCount[*GuildQuery](), gq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (gq *GuildQuery) CountX(ctx context.Context) int {
	count, err := gq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (gq *GuildQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, gq.ctx, "Exist")
	switch _, err := gq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (gq *GuildQuery) ExistX(ctx context.Context) bool {
	exist, err := gq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the GuildQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (gq *GuildQuery) Clone() *GuildQuery {
	if gq == nil {
		return nil
	}
	return &GuildQuery{
		config:               gq.config,
		ctx:                  gq.ctx.Clone(),
		order:                append([]guild.OrderOption{}, gq.order...),
		inters:               append([]Interceptor{}, gq.inters...),
		predicates:           append([]predicate.Guild{}, gq.predicates...),
		withGuildConfig:      gq.withGuildConfig.Clone(),
		withGuildAdminConfig: gq.withGuildAdminConfig.Clone(),
		withGuildEvents:      gq.withGuildEvents.Clone(),
		withAdmins:           gq.withAdmins.Clone(),
		// clone intermediate query.
		sql:  gq.sql.Clone(),
		path: gq.path,
	}
}

// WithGuildConfig tells the query-builder to eager-load the nodes that are connected to
// the "guild_config" edge. The optional arguments are used to configure the query builder of the edge.
func (gq *GuildQuery) WithGuildConfig(opts ...func(*GuildConfigQuery)) *GuildQuery {
	query := (&GuildConfigClient{config: gq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	gq.withGuildConfig = query
	return gq
}

// WithGuildAdminConfig tells the query-builder to eager-load the nodes that are connected to
// the "guild_admin_config" edge. The optional arguments are used to configure the query builder of the edge.
func (gq *GuildQuery) WithGuildAdminConfig(opts ...func(*GuildAdminConfigQuery)) *GuildQuery {
	query := (&GuildAdminConfigClient{config: gq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	gq.withGuildAdminConfig = query
	return gq
}

// WithGuildEvents tells the query-builder to eager-load the nodes that are connected to
// the "guild_events" edge. The optional arguments are used to configure the query builder of the edge.
func (gq *GuildQuery) WithGuildEvents(opts ...func(*GuildEventQuery)) *GuildQuery {
	query := (&GuildEventClient{config: gq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	gq.withGuildEvents = query
	return gq
}

// WithAdmins tells the query-builder to eager-load the nodes that are connected to
// the "admins" edge. The optional arguments are used to configure the query builder of the edge.
func (gq *GuildQuery) WithAdmins(opts ...func(*UserQuery)) *GuildQuery {
	query := (&UserClient{config: gq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	gq.withAdmins = query
	return gq
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
//	client.Guild.Query().
//		GroupBy(guild.FieldCreateTime).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (gq *GuildQuery) GroupBy(field string, fields ...string) *GuildGroupBy {
	gq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &GuildGroupBy{build: gq}
	grbuild.flds = &gq.ctx.Fields
	grbuild.label = guild.Label
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
//	client.Guild.Query().
//		Select(guild.FieldCreateTime).
//		Scan(ctx, &v)
func (gq *GuildQuery) Select(fields ...string) *GuildSelect {
	gq.ctx.Fields = append(gq.ctx.Fields, fields...)
	sbuild := &GuildSelect{GuildQuery: gq}
	sbuild.label = guild.Label
	sbuild.flds, sbuild.scan = &gq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a GuildSelect configured with the given aggregations.
func (gq *GuildQuery) Aggregate(fns ...AggregateFunc) *GuildSelect {
	return gq.Select().Aggregate(fns...)
}

func (gq *GuildQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range gq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, gq); err != nil {
				return err
			}
		}
	}
	for _, f := range gq.ctx.Fields {
		if !guild.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if gq.path != nil {
		prev, err := gq.path(ctx)
		if err != nil {
			return err
		}
		gq.sql = prev
	}
	if guild.Policy == nil {
		return errors.New("ent: uninitialized guild.Policy (forgotten import ent/runtime?)")
	}
	if err := guild.Policy.EvalQuery(ctx, gq); err != nil {
		return err
	}
	return nil
}

func (gq *GuildQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Guild, error) {
	var (
		nodes       = []*Guild{}
		_spec       = gq.querySpec()
		loadedTypes = [4]bool{
			gq.withGuildConfig != nil,
			gq.withGuildAdminConfig != nil,
			gq.withGuildEvents != nil,
			gq.withAdmins != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Guild).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Guild{config: gq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if len(gq.modifiers) > 0 {
		_spec.Modifiers = gq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, gq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := gq.withGuildConfig; query != nil {
		if err := gq.loadGuildConfig(ctx, query, nodes, nil,
			func(n *Guild, e *GuildConfig) { n.Edges.GuildConfig = e }); err != nil {
			return nil, err
		}
	}
	if query := gq.withGuildAdminConfig; query != nil {
		if err := gq.loadGuildAdminConfig(ctx, query, nodes, nil,
			func(n *Guild, e *GuildAdminConfig) { n.Edges.GuildAdminConfig = e }); err != nil {
			return nil, err
		}
	}
	if query := gq.withGuildEvents; query != nil {
		if err := gq.loadGuildEvents(ctx, query, nodes,
			func(n *Guild) { n.Edges.GuildEvents = []*GuildEvent{} },
			func(n *Guild, e *GuildEvent) { n.Edges.GuildEvents = append(n.Edges.GuildEvents, e) }); err != nil {
			return nil, err
		}
	}
	if query := gq.withAdmins; query != nil {
		if err := gq.loadAdmins(ctx, query, nodes,
			func(n *Guild) { n.Edges.Admins = []*User{} },
			func(n *Guild, e *User) { n.Edges.Admins = append(n.Edges.Admins, e) }); err != nil {
			return nil, err
		}
	}
	for name, query := range gq.withNamedGuildEvents {
		if err := gq.loadGuildEvents(ctx, query, nodes,
			func(n *Guild) { n.appendNamedGuildEvents(name) },
			func(n *Guild, e *GuildEvent) { n.appendNamedGuildEvents(name, e) }); err != nil {
			return nil, err
		}
	}
	for name, query := range gq.withNamedAdmins {
		if err := gq.loadAdmins(ctx, query, nodes,
			func(n *Guild) { n.appendNamedAdmins(name) },
			func(n *Guild, e *User) { n.appendNamedAdmins(name, e) }); err != nil {
			return nil, err
		}
	}
	for i := range gq.loadTotal {
		if err := gq.loadTotal[i](ctx, nodes); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (gq *GuildQuery) loadGuildConfig(ctx context.Context, query *GuildConfigQuery, nodes []*Guild, init func(*Guild), assign func(*Guild, *GuildConfig)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[int]*Guild)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
	}
	query.withFKs = true
	query.Where(predicate.GuildConfig(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(guild.GuildConfigColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.guild_guild_config
		if fk == nil {
			return fmt.Errorf(`foreign-key "guild_guild_config" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "guild_guild_config" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}
func (gq *GuildQuery) loadGuildAdminConfig(ctx context.Context, query *GuildAdminConfigQuery, nodes []*Guild, init func(*Guild), assign func(*Guild, *GuildAdminConfig)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[int]*Guild)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
	}
	query.withFKs = true
	query.Where(predicate.GuildAdminConfig(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(guild.GuildAdminConfigColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.guild_guild_admin_config
		if fk == nil {
			return fmt.Errorf(`foreign-key "guild_guild_admin_config" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "guild_guild_admin_config" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}
func (gq *GuildQuery) loadGuildEvents(ctx context.Context, query *GuildEventQuery, nodes []*Guild, init func(*Guild), assign func(*Guild, *GuildEvent)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[int]*Guild)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	query.withFKs = true
	query.Where(predicate.GuildEvent(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(guild.GuildEventsColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.guild_guild_events
		if fk == nil {
			return fmt.Errorf(`foreign-key "guild_guild_events" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "guild_guild_events" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}
func (gq *GuildQuery) loadAdmins(ctx context.Context, query *UserQuery, nodes []*Guild, init func(*Guild), assign func(*Guild, *User)) error {
	edgeIDs := make([]driver.Value, len(nodes))
	byID := make(map[int]*Guild)
	nids := make(map[int]map[*Guild]struct{})
	for i, node := range nodes {
		edgeIDs[i] = node.ID
		byID[node.ID] = node
		if init != nil {
			init(node)
		}
	}
	query.Where(func(s *sql.Selector) {
		joinT := sql.Table(guild.AdminsTable)
		s.Join(joinT).On(s.C(user.FieldID), joinT.C(guild.AdminsPrimaryKey[0]))
		s.Where(sql.InValues(joinT.C(guild.AdminsPrimaryKey[1]), edgeIDs...))
		columns := s.SelectedColumns()
		s.Select(joinT.C(guild.AdminsPrimaryKey[1]))
		s.AppendSelect(columns...)
		s.SetDistinct(false)
	})
	if err := query.prepareQuery(ctx); err != nil {
		return err
	}
	qr := QuerierFunc(func(ctx context.Context, q Query) (Value, error) {
		return query.sqlAll(ctx, func(_ context.Context, spec *sqlgraph.QuerySpec) {
			assign := spec.Assign
			values := spec.ScanValues
			spec.ScanValues = func(columns []string) ([]any, error) {
				values, err := values(columns[1:])
				if err != nil {
					return nil, err
				}
				return append([]any{new(sql.NullInt64)}, values...), nil
			}
			spec.Assign = func(columns []string, values []any) error {
				outValue := int(values[0].(*sql.NullInt64).Int64)
				inValue := int(values[1].(*sql.NullInt64).Int64)
				if nids[inValue] == nil {
					nids[inValue] = map[*Guild]struct{}{byID[outValue]: {}}
					return assign(columns[1:], values[1:])
				}
				nids[inValue][byID[outValue]] = struct{}{}
				return nil
			}
		})
	})
	neighbors, err := withInterceptors[[]*User](ctx, query, qr, query.inters)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected "admins" node returned %v`, n.ID)
		}
		for kn := range nodes {
			assign(kn, n)
		}
	}
	return nil
}

func (gq *GuildQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := gq.querySpec()
	if len(gq.modifiers) > 0 {
		_spec.Modifiers = gq.modifiers
	}
	_spec.Node.Columns = gq.ctx.Fields
	if len(gq.ctx.Fields) > 0 {
		_spec.Unique = gq.ctx.Unique != nil && *gq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, gq.driver, _spec)
}

func (gq *GuildQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(guild.Table, guild.Columns, sqlgraph.NewFieldSpec(guild.FieldID, field.TypeInt))
	_spec.From = gq.sql
	if unique := gq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if gq.path != nil {
		_spec.Unique = true
	}
	if fields := gq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, guild.FieldID)
		for i := range fields {
			if fields[i] != guild.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := gq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := gq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := gq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := gq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (gq *GuildQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(gq.driver.Dialect())
	t1 := builder.Table(guild.Table)
	columns := gq.ctx.Fields
	if len(columns) == 0 {
		columns = guild.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if gq.sql != nil {
		selector = gq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if gq.ctx.Unique != nil && *gq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range gq.predicates {
		p(selector)
	}
	for _, p := range gq.order {
		p(selector)
	}
	if offset := gq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := gq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// WithNamedGuildEvents tells the query-builder to eager-load the nodes that are connected to the "guild_events"
// edge with the given name. The optional arguments are used to configure the query builder of the edge.
func (gq *GuildQuery) WithNamedGuildEvents(name string, opts ...func(*GuildEventQuery)) *GuildQuery {
	query := (&GuildEventClient{config: gq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	if gq.withNamedGuildEvents == nil {
		gq.withNamedGuildEvents = make(map[string]*GuildEventQuery)
	}
	gq.withNamedGuildEvents[name] = query
	return gq
}

// WithNamedAdmins tells the query-builder to eager-load the nodes that are connected to the "admins"
// edge with the given name. The optional arguments are used to configure the query builder of the edge.
func (gq *GuildQuery) WithNamedAdmins(name string, opts ...func(*UserQuery)) *GuildQuery {
	query := (&UserClient{config: gq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	if gq.withNamedAdmins == nil {
		gq.withNamedAdmins = make(map[string]*UserQuery)
	}
	gq.withNamedAdmins[name] = query
	return gq
}

// GuildGroupBy is the group-by builder for Guild entities.
type GuildGroupBy struct {
	selector
	build *GuildQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (ggb *GuildGroupBy) Aggregate(fns ...AggregateFunc) *GuildGroupBy {
	ggb.fns = append(ggb.fns, fns...)
	return ggb
}

// Scan applies the selector query and scans the result into the given value.
func (ggb *GuildGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, ggb.build.ctx, "GroupBy")
	if err := ggb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*GuildQuery, *GuildGroupBy](ctx, ggb.build, ggb, ggb.build.inters, v)
}

func (ggb *GuildGroupBy) sqlScan(ctx context.Context, root *GuildQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(ggb.fns))
	for _, fn := range ggb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*ggb.flds)+len(ggb.fns))
		for _, f := range *ggb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*ggb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := ggb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// GuildSelect is the builder for selecting fields of Guild entities.
type GuildSelect struct {
	*GuildQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (gs *GuildSelect) Aggregate(fns ...AggregateFunc) *GuildSelect {
	gs.fns = append(gs.fns, fns...)
	return gs
}

// Scan applies the selector query and scans the result into the given value.
func (gs *GuildSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, gs.ctx, "Select")
	if err := gs.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*GuildQuery, *GuildSelect](ctx, gs.GuildQuery, gs, gs.inters, v)
}

func (gs *GuildSelect) sqlScan(ctx context.Context, root *GuildQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(gs.fns))
	for _, fn := range gs.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*gs.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := gs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
