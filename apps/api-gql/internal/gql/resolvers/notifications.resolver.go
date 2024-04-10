package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"fmt"

	"github.com/twirapp/twir/apps/api-gql/internal/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/graph"
)

// CreateNotification is the resolver for the createNotification field.
func (r *mutationResolver) CreateNotification(ctx context.Context, text string, userID *string) (*gqlmodel.Notification, error) {
	panic(fmt.Errorf("not implemented: CreateNotification - createNotification"))
}

// Notifications is the resolver for the notifications field.
func (r *queryResolver) Notifications(ctx context.Context, userID string) ([]gqlmodel.Notification, error) {
	panic(fmt.Errorf("not implemented: Notifications - notifications"))
}

// NewNotification is the resolver for the newNotification field.
func (r *subscriptionResolver) NewNotification(ctx context.Context) (<-chan *gqlmodel.Notification, error) {
	channel := make(chan *gqlmodel.Notification)

	go func() {

	}()

	return channel, nil
}

// Subscription returns graph.SubscriptionResolver implementation.
func (r *Resolver) Subscription() graph.SubscriptionResolver { return &subscriptionResolver{r} }

type subscriptionResolver struct{ *Resolver }
