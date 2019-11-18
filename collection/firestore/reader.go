package firestore

import (
	"context"
	"fmt"

	"github.com/rjansen/fivecolors/collection"
	"github.com/rjansen/l"
	"github.com/rjansen/raizel/firestore"
)

var (
	collectionFmt = "environments/development/%s"
)

type reader struct {
	logger l.Logger
	client firestore.Client
}

func NewReader(logger l.Logger, client firestore.Client) collection.Reader {
	return &reader{logger: logger, client: client}
}

func (r *reader) Card(ctx context.Context, id string) (*collection.Card, error) {
	var (
		cardRef = fmt.Sprintf(collectionFmt, fmt.Sprintf("cards/%s", id))
		card    collection.Card
	)
	r.logger.Debug(ctx, "resolve.card.firestore", l.NewValue("arguments", id))

	document, err := r.client.Doc(cardRef).Get(ctx)
	if err != nil {
		return nil, err
	}

	err = document.DataTo(&card)
	if err != nil {
		return nil, err
	}

	r.logger.Debug(ctx, "resolve.card.firestore.fetched", l.NewValue("card", card))
	return &card, err
}

func (r *reader) CardBy(ctx context.Context, filter collection.CardFilter) ([]collection.Card, error) {
	var (
		cardsRef   = fmt.Sprintf(collectionFmt, "cards")
		cards      []collection.Card
		cardsQuery firestore.Query = r.client.Collection(cardsRef)
	)
	r.logger.Info(ctx, "resolve.cardby.firestore", l.NewValue("arguments", filter))

	if filter.Name != nil {
		cardsQuery = cardsQuery.Where("Name", ">=", *filter.Name)
	}
	if len(filter.Types) > 0 {
		cardsQuery = cardsQuery.Where("Types", "==", filter.Types)
	}
	if len(filter.Costs) > 0 {
		cardsQuery = cardsQuery.Where("Costs", "==", filter.Costs)
	}
	if filter.Set != nil {
		if filter.Set.ID != nil {
			cardsQuery = cardsQuery.Where("Set.ID", "==", *filter.Set.ID)
		}
		if filter.Set.Name != nil {
			cardsQuery = cardsQuery.Where("Set.Name", "==", *filter.Set.Name)
		}
		if filter.Set.Alias != nil {
			cardsQuery = cardsQuery.Where("Set.Alias", "==", *filter.Set.Alias)
		}
	}
	if filter.Rarity != nil {
		if filter.Rarity.ID != nil {
			cardsQuery = cardsQuery.Where("Rarity.ID", "==", *filter.Rarity.ID)
		}
		if filter.Rarity.Name != nil {
			cardsQuery = cardsQuery.Where("Rarity.Name", "==", *filter.Rarity.Name)
		}
		if filter.Rarity.Alias != nil {
			cardsQuery = cardsQuery.Where("Rarity.Alias", "==", *filter.Rarity.Alias)
		}
	}
	cardsQuery = cardsQuery.OrderBy("Name", firestore.Asc)
	r.logger.Info(ctx, "resolve.cardby.firestore.query", l.NewValue("query", cardsQuery))
	documents, err := cardsQuery.Documents(ctx).GetAll()
	if err != nil {
		r.logger.Error(ctx, "resolve.cardby.firestore.query_err", l.NewValue("error", err))
		return cards, err
	}
	for index, document := range documents {
		var card collection.Card
		err := document.DataTo(&card)
		if err != nil {
			r.logger.Error(
				ctx,
				"resolve.cardby.firestore.fetch_err",
				l.NewValue("index", index),
				l.NewValue("error", err),
			)
			return cards, err
		}
		cards = append(cards, card)
	}

	r.logger.Info(ctx, "resolve.cardby.firestore.fetched", l.NewValue("cards.len", len(cards)))
	return cards, err
}

func (r *reader) Set(ctx context.Context, id string) (*collection.Set, error) {
	var (
		setRef = fmt.Sprintf(collectionFmt, fmt.Sprintf("sets/%s", id))
		set    collection.Set
	)
	r.logger.Debug(ctx, "resolve.set.firestore", l.NewValue("arguments", id))

	document, err := r.client.Doc(setRef).Get(ctx)
	if err != nil {
		return nil, err
	}

	err = document.DataTo(&set)
	if err != nil {
		return nil, err
	}

	r.logger.Debug(ctx, "resolve.set.firestore.fetched", l.NewValue("set", set))
	return &set, err
}

func (r *reader) SetBy(ctx context.Context, filter collection.SetFilter) ([]collection.Set, error) {
	var (
		setsRef   = fmt.Sprintf(collectionFmt, "sets")
		sets      []collection.Set
		setsQuery firestore.Query = r.client.Collection(setsRef)
	)
	r.logger.Info(ctx, "resolve.setby.firestore", l.NewValue("arguments", filter))

	if filter.Name != nil {
		setsQuery = setsQuery.Where("Name", ">=", *filter.Name)
	}
	if filter.Alias != nil {
		setsQuery = setsQuery.Where("Alias", "==", *filter.Alias)
	}
	setsQuery = setsQuery.OrderBy("Name", firestore.Asc)

	r.logger.Info(ctx, "resolve.setdby.firestore.query", l.NewValue("query", setsQuery))
	documents, err := setsQuery.Documents(ctx).GetAll()
	if err != nil {
		r.logger.Error(ctx, "resolve.setby.firestore.query_err", l.NewValue("error", err))
		return sets, err
	}
	for index, document := range documents {
		var set collection.Set
		err := document.DataTo(&set)
		if err != nil {
			r.logger.Error(
				ctx,
				"resolve.setby.firestore.fetch_err",
				l.NewValue("index", index),
				l.NewValue("error", err),
			)
			return sets, err
		}
		sets = append(sets, set)
	}

	r.logger.Info(ctx, "resolve.setby.firestore.fetched", l.NewValue("sets.len", len(sets)))
	return sets, err
}
