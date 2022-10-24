package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"

	"github.com/lrstanley/chix"
	"github.com/lrstanley/spectrograph/internal/ent"
)

// BanUser is the resolver for the banUser field.
func (r *mutationResolver) BanUser(ctx context.Context, id int, reason string) (bool, error) {
	user := chix.IdentFromContext[ent.User](ctx)
	if user == nil {
		return false, errors.New("not authenticated")
	}

	db := ent.FromContext(ctx)

	err := db.User.
		UpdateOneID(id).
		SetBanned(true).
		SetBanReason(reason).
		SetBannedByID(user.ID).
		Exec(ctx)
	if err != nil {
		return false, err
	}

	return true, nil
}

// UnbanUser is the resolver for the unbanUser field.
func (r *mutationResolver) UnbanUser(ctx context.Context, id int) (bool, error) {
	db := ent.FromContext(ctx)

	err := db.User.
		UpdateOneID(id).
		SetBanned(false).
		ClearBanReason().
		ClearBannedBy().
		Exec(ctx)
	if err != nil {
		return false, err
	}

	return true, nil
}

// DeleteAccount is the resolver for the deleteAccount field.
func (r *mutationResolver) DeleteAccount(ctx context.Context, noop *int) (bool, error) {
	uid := chix.IDFromContext[int](ctx)

	if err := ent.FromContext(ctx).User.DeleteOneID(uid).Exec(ctx); err != nil {
		return false, err
	}

	return true, nil
}

// Self is the resolver for the self field.
func (r *queryResolver) Self(ctx context.Context) (*ent.User, error) {
	user := chix.IdentFromContext[ent.User](ctx)
	if user == nil {
		return nil, errors.New("not authenticated")
	}

	return user, nil
}
