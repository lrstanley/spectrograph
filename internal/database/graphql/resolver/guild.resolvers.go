package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"time"

	"github.com/apex/log"
	"github.com/lrstanley/spectrograph/internal/database/graphql/gqlhandler"
	"github.com/lrstanley/spectrograph/internal/ent"
	"github.com/lrstanley/spectrograph/internal/ent/guildevent"
)

// UpdateGuildConfig is the resolver for the updateGuildConfig field.
func (r *mutationResolver) UpdateGuildConfig(ctx context.Context, id int, input ent.UpdateGuildConfigInput) (*ent.GuildConfig, error) {
	return ent.FromContext(ctx).GuildConfig.UpdateOneID(id).SetInput(input).Save(ctx)
}

// UpdateGuildAdminConfig is the resolver for the updateGuildAdminConfig field.
func (r *mutationResolver) UpdateGuildAdminConfig(ctx context.Context, id int, input ent.UpdateGuildAdminConfigInput) (*ent.GuildAdminConfig, error) {
	return ent.FromContext(ctx).GuildAdminConfig.UpdateOneID(id).SetInput(input).Save(ctx)
}

// GuildEventAdded is the resolver for the guildEventAdded field.
func (r *subscriptionResolver) GuildEventAdded(ctx context.Context, input ent.GuildEventWhereInput) (<-chan *ent.GuildEvent, error) {
	db := ent.FromContext(ctx)

	q, err := input.P()
	if err != nil {
		return nil, err
	}

	ch := make(chan *ent.GuildEvent)
	oldestTimestamp := time.Now().Add(-(time.Hour * 2))
	lastEventID := 0

	fn := func() error {
		events, err := db.GuildEvent.Query().
			Where(
				guildevent.And(q, guildevent.IDGT(lastEventID), guildevent.CreateTimeGT(oldestTimestamp)),
			).All(ctx) // .Order(ent.Asc(guildevent.FieldCreateTime))
		if err != nil {
			return err
		}

		for _, event := range events {
			lastEventID = event.ID
			ch <- event
		}

		return nil
	}

	timer := time.NewTimer(4 * time.Second)
	go func() {
		if err = fn(); err != nil {
			log.FromContext(ctx).WithError(err).Error("failed to query guild events")
			return
		}

		for {
			select {
			case <-ctx.Done():
				log.FromContext(ctx).Info("closing guild event subscription")
				return
			case <-timer.C:
				if err = fn(); err != nil {
					log.FromContext(ctx).WithError(err).Error("failed to query guild events")
					return
				}
			}
		}
	}()

	return ch, nil
}

// Subscription returns gqlhandler.SubscriptionResolver implementation.
func (r *Resolver) Subscription() gqlhandler.SubscriptionResolver { return &subscriptionResolver{r} }

type subscriptionResolver struct{ *Resolver }
