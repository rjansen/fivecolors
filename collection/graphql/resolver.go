package graphql

import (
	"context"

	"github.com/rjansen/fivecolors/collection"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct{}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) UpsertSet(ctx context.Context, set collection.SetInput) (*collection.Set, error) {
	panic("not implemented")
}
func (r *mutationResolver) UpsertCards(ctx context.Context, cards []*collection.CardInput) ([]*collection.Card, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Set(ctx context.Context, id string) (*collection.Set, error) {
	panic("not implemented")
}
func (r *queryResolver) Card(ctx context.Context, id string) (*collection.Card, error) {
	panic("not implemented")
}
func (r *queryResolver) SetBy(ctx context.Context, filter collection.SetFilter) ([]*collection.Set, error) {
	panic("not implemented")
}
func (r *queryResolver) CardBy(ctx context.Context, filter collection.CardFilter) ([]*collection.Card, error) {
	panic("not implemented")
}
