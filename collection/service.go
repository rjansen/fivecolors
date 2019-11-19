package collection

import "context"

type (
	Reader interface {
		Set(ctx context.Context, id string) (*Set, error)
		Card(ctx context.Context, id string) (*Card, error)
		SetBy(ctx context.Context, filter SetFilter) ([]Set, error)
		CardBy(ctx context.Context, filter CardFilter) ([]Card, error)
	}

	Writer interface {
		UpsertSet(ctx context.Context, set SetInput) (*Set, error)
		UpsertCards(ctx context.Context, cards []CardInput) (*UpsertCards, error)
	}
)
