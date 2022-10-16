// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.
//
// DO NOT EDIT, CODE GENERATED BY entc.

package ent

import (
	"context"
	"database/sql/driver"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/99designs/gqlgen/graphql"
	"github.com/lrstanley/spectrograph/internal/ent/guild"
	"github.com/lrstanley/spectrograph/internal/ent/user"
)

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (gu *GuildQuery) CollectFields(ctx context.Context, satisfies ...string) (*GuildQuery, error) {
	fc := graphql.GetFieldContext(ctx)
	if fc == nil {
		return gu, nil
	}
	if err := gu.collectField(ctx, graphql.GetOperationContext(ctx), fc.Field, nil, satisfies...); err != nil {
		return nil, err
	}
	return gu, nil
}

func (gu *GuildQuery) collectField(ctx context.Context, op *graphql.OperationContext, field graphql.CollectedField, path []string, satisfies ...string) error {
	path = append([]string(nil), path...)
	for _, field := range graphql.CollectFields(op, field.Selections, satisfies) {
		switch field.Name {
		case "guildConfig":
			var (
				alias = field.Alias
				path  = append(path, alias)
				query = &GuildConfigQuery{config: gu.config}
			)
			if err := query.collectField(ctx, op, field, path, satisfies...); err != nil {
				return err
			}
			gu.withGuildConfig = query
		case "guildAdminConfig":
			var (
				alias = field.Alias
				path  = append(path, alias)
				query = &GuildAdminConfigQuery{config: gu.config}
			)
			if err := query.collectField(ctx, op, field, path, satisfies...); err != nil {
				return err
			}
			gu.withGuildAdminConfig = query
		case "guildEvents":
			var (
				alias = field.Alias
				path  = append(path, alias)
				query = &GuildEventQuery{config: gu.config}
			)
			if err := query.collectField(ctx, op, field, path, satisfies...); err != nil {
				return err
			}
			gu.WithNamedGuildEvents(alias, func(wq *GuildEventQuery) {
				*wq = *query
			})
		case "admins":
			var (
				alias = field.Alias
				path  = append(path, alias)
				query = &UserQuery{config: gu.config}
			)
			args := newUserPaginateArgs(fieldArgs(ctx, new(UserWhereInput), path...))
			if err := validateFirstLast(args.first, args.last); err != nil {
				return fmt.Errorf("validate first and last in path %q: %w", path, err)
			}
			pager, err := newUserPager(args.opts)
			if err != nil {
				return fmt.Errorf("create new pager in path %q: %w", path, err)
			}
			if query, err = pager.applyFilter(query); err != nil {
				return err
			}
			ignoredEdges := !hasCollectedField(ctx, append(path, edgesField)...)
			if hasCollectedField(ctx, append(path, totalCountField)...) || hasCollectedField(ctx, append(path, pageInfoField)...) {
				hasPagination := args.after != nil || args.first != nil || args.before != nil || args.last != nil
				if hasPagination || ignoredEdges {
					query := query.Clone()
					gu.loadTotal = append(gu.loadTotal, func(ctx context.Context, nodes []*Guild) error {
						ids := make([]driver.Value, len(nodes))
						for i := range nodes {
							ids[i] = nodes[i].ID
						}
						var v []struct {
							NodeID int `sql:"guild_id"`
							Count  int `sql:"count"`
						}
						query.Where(func(s *sql.Selector) {
							joinT := sql.Table(guild.AdminsTable)
							s.Join(joinT).On(s.C(user.FieldID), joinT.C(guild.AdminsPrimaryKey[0]))
							s.Where(sql.InValues(joinT.C(guild.AdminsPrimaryKey[1]), ids...))
							s.Select(joinT.C(guild.AdminsPrimaryKey[1]), sql.Count("*"))
							s.GroupBy(joinT.C(guild.AdminsPrimaryKey[1]))
						})
						if err := query.Select().Scan(ctx, &v); err != nil {
							return err
						}
						m := make(map[int]int, len(v))
						for i := range v {
							m[v[i].NodeID] = v[i].Count
						}
						for i := range nodes {
							n := m[nodes[i].ID]
							if nodes[i].Edges.totalCount[3] == nil {
								nodes[i].Edges.totalCount[3] = make(map[string]int)
							}
							nodes[i].Edges.totalCount[3][alias] = n
						}
						return nil
					})
				} else {
					gu.loadTotal = append(gu.loadTotal, func(_ context.Context, nodes []*Guild) error {
						for i := range nodes {
							n := len(nodes[i].Edges.Admins)
							if nodes[i].Edges.totalCount[3] == nil {
								nodes[i].Edges.totalCount[3] = make(map[string]int)
							}
							nodes[i].Edges.totalCount[3][alias] = n
						}
						return nil
					})
				}
			}
			if ignoredEdges || (args.first != nil && *args.first == 0) || (args.last != nil && *args.last == 0) {
				continue
			}

			query = pager.applyCursors(query, args.after, args.before)
			if limit := paginateLimit(args.first, args.last); limit > 0 {
				modify := limitRows(guild.AdminsPrimaryKey[1], limit, pager.orderExpr(args.last != nil))
				query.modifiers = append(query.modifiers, modify)
			} else {
				query = pager.applyOrder(query, args.last != nil)
			}
			path = append(path, edgesField, nodeField)
			if field := collectedField(ctx, path...); field != nil {
				if err := query.collectField(ctx, op, *field, path, satisfies...); err != nil {
					return err
				}
			}
			gu.WithNamedAdmins(alias, func(wq *UserQuery) {
				*wq = *query
			})
		}
	}
	return nil
}

type guildPaginateArgs struct {
	first, last   *int
	after, before *Cursor
	opts          []GuildPaginateOption
}

func newGuildPaginateArgs(rv map[string]interface{}) *guildPaginateArgs {
	args := &guildPaginateArgs{}
	if rv == nil {
		return args
	}
	if v := rv[firstField]; v != nil {
		args.first = v.(*int)
	}
	if v := rv[lastField]; v != nil {
		args.last = v.(*int)
	}
	if v := rv[afterField]; v != nil {
		args.after = v.(*Cursor)
	}
	if v := rv[beforeField]; v != nil {
		args.before = v.(*Cursor)
	}
	if v, ok := rv[orderByField]; ok {
		switch v := v.(type) {
		case map[string]interface{}:
			var (
				err1, err2 error
				order      = &GuildOrder{Field: &GuildOrderField{}}
			)
			if d, ok := v[directionField]; ok {
				err1 = order.Direction.UnmarshalGQL(d)
			}
			if f, ok := v[fieldField]; ok {
				err2 = order.Field.UnmarshalGQL(f)
			}
			if err1 == nil && err2 == nil {
				args.opts = append(args.opts, WithGuildOrder(order))
			}
		case *GuildOrder:
			if v != nil {
				args.opts = append(args.opts, WithGuildOrder(v))
			}
		}
	}
	if v, ok := rv[whereField].(*GuildWhereInput); ok {
		args.opts = append(args.opts, WithGuildFilter(v.Filter))
	}
	return args
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (gac *GuildAdminConfigQuery) CollectFields(ctx context.Context, satisfies ...string) (*GuildAdminConfigQuery, error) {
	fc := graphql.GetFieldContext(ctx)
	if fc == nil {
		return gac, nil
	}
	if err := gac.collectField(ctx, graphql.GetOperationContext(ctx), fc.Field, nil, satisfies...); err != nil {
		return nil, err
	}
	return gac, nil
}

func (gac *GuildAdminConfigQuery) collectField(ctx context.Context, op *graphql.OperationContext, field graphql.CollectedField, path []string, satisfies ...string) error {
	path = append([]string(nil), path...)
	for _, field := range graphql.CollectFields(op, field.Selections, satisfies) {
		switch field.Name {
		case "guild":
			var (
				alias = field.Alias
				path  = append(path, alias)
				query = &GuildQuery{config: gac.config}
			)
			if err := query.collectField(ctx, op, field, path, satisfies...); err != nil {
				return err
			}
			gac.withGuild = query
		}
	}
	return nil
}

type guildadminconfigPaginateArgs struct {
	first, last   *int
	after, before *Cursor
	opts          []GuildAdminConfigPaginateOption
}

func newGuildAdminConfigPaginateArgs(rv map[string]interface{}) *guildadminconfigPaginateArgs {
	args := &guildadminconfigPaginateArgs{}
	if rv == nil {
		return args
	}
	if v := rv[firstField]; v != nil {
		args.first = v.(*int)
	}
	if v := rv[lastField]; v != nil {
		args.last = v.(*int)
	}
	if v := rv[afterField]; v != nil {
		args.after = v.(*Cursor)
	}
	if v := rv[beforeField]; v != nil {
		args.before = v.(*Cursor)
	}
	if v, ok := rv[whereField].(*GuildAdminConfigWhereInput); ok {
		args.opts = append(args.opts, WithGuildAdminConfigFilter(v.Filter))
	}
	return args
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (gc *GuildConfigQuery) CollectFields(ctx context.Context, satisfies ...string) (*GuildConfigQuery, error) {
	fc := graphql.GetFieldContext(ctx)
	if fc == nil {
		return gc, nil
	}
	if err := gc.collectField(ctx, graphql.GetOperationContext(ctx), fc.Field, nil, satisfies...); err != nil {
		return nil, err
	}
	return gc, nil
}

func (gc *GuildConfigQuery) collectField(ctx context.Context, op *graphql.OperationContext, field graphql.CollectedField, path []string, satisfies ...string) error {
	path = append([]string(nil), path...)
	for _, field := range graphql.CollectFields(op, field.Selections, satisfies) {
		switch field.Name {
		case "guild":
			var (
				alias = field.Alias
				path  = append(path, alias)
				query = &GuildQuery{config: gc.config}
			)
			if err := query.collectField(ctx, op, field, path, satisfies...); err != nil {
				return err
			}
			gc.withGuild = query
		}
	}
	return nil
}

type guildconfigPaginateArgs struct {
	first, last   *int
	after, before *Cursor
	opts          []GuildConfigPaginateOption
}

func newGuildConfigPaginateArgs(rv map[string]interface{}) *guildconfigPaginateArgs {
	args := &guildconfigPaginateArgs{}
	if rv == nil {
		return args
	}
	if v := rv[firstField]; v != nil {
		args.first = v.(*int)
	}
	if v := rv[lastField]; v != nil {
		args.last = v.(*int)
	}
	if v := rv[afterField]; v != nil {
		args.after = v.(*Cursor)
	}
	if v := rv[beforeField]; v != nil {
		args.before = v.(*Cursor)
	}
	if v, ok := rv[whereField].(*GuildConfigWhereInput); ok {
		args.opts = append(args.opts, WithGuildConfigFilter(v.Filter))
	}
	return args
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (ge *GuildEventQuery) CollectFields(ctx context.Context, satisfies ...string) (*GuildEventQuery, error) {
	fc := graphql.GetFieldContext(ctx)
	if fc == nil {
		return ge, nil
	}
	if err := ge.collectField(ctx, graphql.GetOperationContext(ctx), fc.Field, nil, satisfies...); err != nil {
		return nil, err
	}
	return ge, nil
}

func (ge *GuildEventQuery) collectField(ctx context.Context, op *graphql.OperationContext, field graphql.CollectedField, path []string, satisfies ...string) error {
	path = append([]string(nil), path...)
	for _, field := range graphql.CollectFields(op, field.Selections, satisfies) {
		switch field.Name {
		case "guild":
			var (
				alias = field.Alias
				path  = append(path, alias)
				query = &GuildQuery{config: ge.config}
			)
			if err := query.collectField(ctx, op, field, path, satisfies...); err != nil {
				return err
			}
			ge.withGuild = query
		}
	}
	return nil
}

type guildeventPaginateArgs struct {
	first, last   *int
	after, before *Cursor
	opts          []GuildEventPaginateOption
}

func newGuildEventPaginateArgs(rv map[string]interface{}) *guildeventPaginateArgs {
	args := &guildeventPaginateArgs{}
	if rv == nil {
		return args
	}
	if v := rv[firstField]; v != nil {
		args.first = v.(*int)
	}
	if v := rv[lastField]; v != nil {
		args.last = v.(*int)
	}
	if v := rv[afterField]; v != nil {
		args.after = v.(*Cursor)
	}
	if v := rv[beforeField]; v != nil {
		args.before = v.(*Cursor)
	}
	if v, ok := rv[orderByField]; ok {
		switch v := v.(type) {
		case map[string]interface{}:
			var (
				err1, err2 error
				order      = &GuildEventOrder{Field: &GuildEventOrderField{}}
			)
			if d, ok := v[directionField]; ok {
				err1 = order.Direction.UnmarshalGQL(d)
			}
			if f, ok := v[fieldField]; ok {
				err2 = order.Field.UnmarshalGQL(f)
			}
			if err1 == nil && err2 == nil {
				args.opts = append(args.opts, WithGuildEventOrder(order))
			}
		case *GuildEventOrder:
			if v != nil {
				args.opts = append(args.opts, WithGuildEventOrder(v))
			}
		}
	}
	if v, ok := rv[whereField].(*GuildEventWhereInput); ok {
		args.opts = append(args.opts, WithGuildEventFilter(v.Filter))
	}
	return args
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (u *UserQuery) CollectFields(ctx context.Context, satisfies ...string) (*UserQuery, error) {
	fc := graphql.GetFieldContext(ctx)
	if fc == nil {
		return u, nil
	}
	if err := u.collectField(ctx, graphql.GetOperationContext(ctx), fc.Field, nil, satisfies...); err != nil {
		return nil, err
	}
	return u, nil
}

func (u *UserQuery) collectField(ctx context.Context, op *graphql.OperationContext, field graphql.CollectedField, path []string, satisfies ...string) error {
	path = append([]string(nil), path...)
	for _, field := range graphql.CollectFields(op, field.Selections, satisfies) {
		switch field.Name {
		case "userGuilds":
			var (
				alias = field.Alias
				path  = append(path, alias)
				query = &GuildQuery{config: u.config}
			)
			args := newGuildPaginateArgs(fieldArgs(ctx, new(GuildWhereInput), path...))
			if err := validateFirstLast(args.first, args.last); err != nil {
				return fmt.Errorf("validate first and last in path %q: %w", path, err)
			}
			pager, err := newGuildPager(args.opts)
			if err != nil {
				return fmt.Errorf("create new pager in path %q: %w", path, err)
			}
			if query, err = pager.applyFilter(query); err != nil {
				return err
			}
			ignoredEdges := !hasCollectedField(ctx, append(path, edgesField)...)
			if hasCollectedField(ctx, append(path, totalCountField)...) || hasCollectedField(ctx, append(path, pageInfoField)...) {
				hasPagination := args.after != nil || args.first != nil || args.before != nil || args.last != nil
				if hasPagination || ignoredEdges {
					query := query.Clone()
					u.loadTotal = append(u.loadTotal, func(ctx context.Context, nodes []*User) error {
						ids := make([]driver.Value, len(nodes))
						for i := range nodes {
							ids[i] = nodes[i].ID
						}
						var v []struct {
							NodeID int `sql:"user_id"`
							Count  int `sql:"count"`
						}
						query.Where(func(s *sql.Selector) {
							joinT := sql.Table(user.UserGuildsTable)
							s.Join(joinT).On(s.C(guild.FieldID), joinT.C(user.UserGuildsPrimaryKey[1]))
							s.Where(sql.InValues(joinT.C(user.UserGuildsPrimaryKey[0]), ids...))
							s.Select(joinT.C(user.UserGuildsPrimaryKey[0]), sql.Count("*"))
							s.GroupBy(joinT.C(user.UserGuildsPrimaryKey[0]))
						})
						if err := query.Select().Scan(ctx, &v); err != nil {
							return err
						}
						m := make(map[int]int, len(v))
						for i := range v {
							m[v[i].NodeID] = v[i].Count
						}
						for i := range nodes {
							n := m[nodes[i].ID]
							if nodes[i].Edges.totalCount[0] == nil {
								nodes[i].Edges.totalCount[0] = make(map[string]int)
							}
							nodes[i].Edges.totalCount[0][alias] = n
						}
						return nil
					})
				} else {
					u.loadTotal = append(u.loadTotal, func(_ context.Context, nodes []*User) error {
						for i := range nodes {
							n := len(nodes[i].Edges.UserGuilds)
							if nodes[i].Edges.totalCount[0] == nil {
								nodes[i].Edges.totalCount[0] = make(map[string]int)
							}
							nodes[i].Edges.totalCount[0][alias] = n
						}
						return nil
					})
				}
			}
			if ignoredEdges || (args.first != nil && *args.first == 0) || (args.last != nil && *args.last == 0) {
				continue
			}

			query = pager.applyCursors(query, args.after, args.before)
			if limit := paginateLimit(args.first, args.last); limit > 0 {
				modify := limitRows(user.UserGuildsPrimaryKey[0], limit, pager.orderExpr(args.last != nil))
				query.modifiers = append(query.modifiers, modify)
			} else {
				query = pager.applyOrder(query, args.last != nil)
			}
			path = append(path, edgesField, nodeField)
			if field := collectedField(ctx, path...); field != nil {
				if err := query.collectField(ctx, op, *field, path, satisfies...); err != nil {
					return err
				}
			}
			u.WithNamedUserGuilds(alias, func(wq *GuildQuery) {
				*wq = *query
			})
		}
	}
	return nil
}

type userPaginateArgs struct {
	first, last   *int
	after, before *Cursor
	opts          []UserPaginateOption
}

func newUserPaginateArgs(rv map[string]interface{}) *userPaginateArgs {
	args := &userPaginateArgs{}
	if rv == nil {
		return args
	}
	if v := rv[firstField]; v != nil {
		args.first = v.(*int)
	}
	if v := rv[lastField]; v != nil {
		args.last = v.(*int)
	}
	if v := rv[afterField]; v != nil {
		args.after = v.(*Cursor)
	}
	if v := rv[beforeField]; v != nil {
		args.before = v.(*Cursor)
	}
	if v, ok := rv[orderByField]; ok {
		switch v := v.(type) {
		case map[string]interface{}:
			var (
				err1, err2 error
				order      = &UserOrder{Field: &UserOrderField{}}
			)
			if d, ok := v[directionField]; ok {
				err1 = order.Direction.UnmarshalGQL(d)
			}
			if f, ok := v[fieldField]; ok {
				err2 = order.Field.UnmarshalGQL(f)
			}
			if err1 == nil && err2 == nil {
				args.opts = append(args.opts, WithUserOrder(order))
			}
		case *UserOrder:
			if v != nil {
				args.opts = append(args.opts, WithUserOrder(v))
			}
		}
	}
	if v, ok := rv[whereField].(*UserWhereInput); ok {
		args.opts = append(args.opts, WithUserFilter(v.Filter))
	}
	return args
}

const (
	afterField     = "after"
	firstField     = "first"
	beforeField    = "before"
	lastField      = "last"
	orderByField   = "orderBy"
	directionField = "direction"
	fieldField     = "field"
	whereField     = "where"
)

func fieldArgs(ctx context.Context, whereInput interface{}, path ...string) map[string]interface{} {
	fc := graphql.GetFieldContext(ctx)
	if fc == nil {
		return nil
	}
	oc := graphql.GetOperationContext(ctx)
	for _, name := range path {
		var field *graphql.CollectedField
		for _, f := range graphql.CollectFields(oc, fc.Field.Selections, nil) {
			if f.Alias == name {
				field = &f
				break
			}
		}
		if field == nil {
			return nil
		}
		cf, err := fc.Child(ctx, *field)
		if err != nil {
			args := field.ArgumentMap(oc.Variables)
			return unmarshalArgs(ctx, whereInput, args)
		}
		fc = cf
	}
	return fc.Args
}

// unmarshalArgs allows extracting the field arguments from their raw representation.
func unmarshalArgs(ctx context.Context, whereInput interface{}, args map[string]interface{}) map[string]interface{} {
	for _, k := range []string{firstField, lastField} {
		v, ok := args[k]
		if !ok {
			continue
		}
		i, err := graphql.UnmarshalInt(v)
		if err == nil {
			args[k] = &i
		}
	}
	for _, k := range []string{beforeField, afterField} {
		v, ok := args[k]
		if !ok {
			continue
		}
		c := &Cursor{}
		if c.UnmarshalGQL(v) == nil {
			args[k] = c
		}
	}
	if v, ok := args[whereField]; ok && whereInput != nil {
		if err := graphql.UnmarshalInputFromContext(ctx, v, whereInput); err == nil {
			args[whereField] = whereInput
		}
	}

	return args
}

func limitRows(partitionBy string, limit int, orderBy ...sql.Querier) func(s *sql.Selector) {
	return func(s *sql.Selector) {
		d := sql.Dialect(s.Dialect())
		s.SetDistinct(false)
		with := d.With("src_query").
			As(s.Clone()).
			With("limited_query").
			As(
				d.Select("*").
					AppendSelectExprAs(
						sql.RowNumber().PartitionBy(partitionBy).OrderExpr(orderBy...),
						"row_number",
					).
					From(d.Table("src_query")),
			)
		t := d.Table("limited_query").As(s.TableName())
		*s = *d.Select(s.UnqualifiedColumns()...).
			From(t).
			Where(sql.LTE(t.C("row_number"), limit)).
			Prefix(with)
	}
}
