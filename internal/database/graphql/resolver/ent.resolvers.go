package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/lrstanley/spectrograph/internal/database/graphql/gqlhandler"
	"github.com/lrstanley/spectrograph/internal/ent"
)

// Node is the resolver for the node field.
func (r *queryResolver) Node(ctx context.Context, id int) (ent.Noder, error) {
	return ent.FromContext(ctx).Noder(ctx, id)
}

// Nodes is the resolver for the nodes field.
func (r *queryResolver) Nodes(ctx context.Context, ids []int) ([]ent.Noder, error) {
	return ent.FromContext(ctx).Noders(ctx, ids)
}

// Guilds is the resolver for the guilds field.
func (r *queryResolver) Guilds(ctx context.Context, after *ent.Cursor, first *int, before *ent.Cursor, last *int, orderBy *ent.GuildOrder, where *ent.GuildWhereInput) (*ent.GuildConnection, error) {
	return ent.FromContext(ctx).Guild.Query().Paginate(
		ctx, after, first, before, last,
		ent.WithGuildOrder(orderBy),
		ent.WithGuildFilter(where.Filter),
	)
}

// Guildadminconfigs is the resolver for the guildadminconfigs field.
func (r *queryResolver) Guildadminconfigs(ctx context.Context, after *ent.Cursor, first *int, before *ent.Cursor, last *int, where *ent.GuildAdminConfigWhereInput) (*ent.GuildAdminConfigConnection, error) {
	return ent.FromContext(ctx).GuildAdminConfig.Query().Paginate(
		ctx, after, first, before, last,
		ent.WithGuildAdminConfigFilter(where.Filter),
	)
}

// Guildconfigs is the resolver for the guildconfigs field.
func (r *queryResolver) Guildconfigs(ctx context.Context, after *ent.Cursor, first *int, before *ent.Cursor, last *int, where *ent.GuildConfigWhereInput) (*ent.GuildConfigConnection, error) {
	return ent.FromContext(ctx).GuildConfig.Query().Paginate(
		ctx, after, first, before, last,
		ent.WithGuildConfigFilter(where.Filter),
	)
}

// Guildevents is the resolver for the guildevents field.
func (r *queryResolver) Guildevents(ctx context.Context, after *ent.Cursor, first *int, before *ent.Cursor, last *int, orderBy *ent.GuildEventOrder, where *ent.GuildEventWhereInput) (*ent.GuildEventConnection, error) {
	return ent.FromContext(ctx).GuildEvent.Query().Paginate(
		ctx, after, first, before, last,
		ent.WithGuildEventOrder(orderBy),
		ent.WithGuildEventFilter(where.Filter),
	)
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context, after *ent.Cursor, first *int, before *ent.Cursor, last *int, orderBy *ent.UserOrder, where *ent.UserWhereInput) (*ent.UserConnection, error) {
	return ent.FromContext(ctx).User.Query().Paginate(
		ctx, after, first, before, last,
		ent.WithUserOrder(orderBy),
		ent.WithUserFilter(where.Filter),
	)
}

// Query returns gqlhandler.QueryResolver implementation.
func (r *Resolver) Query() gqlhandler.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
