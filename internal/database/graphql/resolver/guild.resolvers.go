package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/lrstanley/spectrograph/internal/ent"
)

// UpdateGuildConfig is the resolver for the updateGuildConfig field.
func (r *mutationResolver) UpdateGuildConfig(ctx context.Context, id int, input ent.UpdateGuildConfigInput) (*ent.GuildConfig, error) {
	return ent.FromContext(ctx).GuildConfig.UpdateOneID(id).SetInput(input).Save(ctx)
}

// UpdateGuildAdminConfig is the resolver for the updateGuildAdminConfig field.
func (r *mutationResolver) UpdateGuildAdminConfig(ctx context.Context, id int, input ent.UpdateGuildAdminConfigInput) (*ent.GuildAdminConfig, error) {
	return ent.FromContext(ctx).GuildAdminConfig.UpdateOneID(id).SetInput(input).Save(ctx)
}
