// +build firestore

package model

import (
	"context"
	"fmt"
	"time"

	"github.com/rjansen/l"
	"github.com/rjansen/raizel/firestore"
	"github.com/rjansen/yggdrasil"
)

var (
	collectionFmt = "environments/development/%s"
)

type Resolver struct {
	tree yggdrasil.Tree
}

func NewResolver(tree yggdrasil.Tree) *Resolver {
	return &Resolver{
		tree: tree,
	}
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{Resolver: r}
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Card(ctx context.Context, id string) (Card, error) {
	var (
		client  = firestore.MustReference(r.tree)
		logger  = l.MustReference(r.tree)
		cardRef = fmt.Sprintf(collectionFmt, fmt.Sprintf("cards/%s", id))
		card    Card
	)
	logger.Debug("resolve.card.firestore", l.NewValue("arguments", id))

	document, err := client.Doc(cardRef).Get(ctx)
	if err != nil {
		return card, err
	}

	err = document.DataTo(&card)
	if err != nil {
		return card, err
	}

	logger.Debug("resolve.card.firestore.fetched", l.NewValue("card", card))
	return card, err
}

func (r *queryResolver) CardBy(ctx context.Context, filter CardFilter) ([]Card, error) {
	var (
		client     = firestore.MustReference(r.tree)
		logger     = l.MustReference(r.tree)
		cardsRef   = fmt.Sprintf(collectionFmt, "cards")
		cards      []Card
		cardsQuery firestore.Query = client.Collection(cardsRef)
	)
	logger.Info("resolve.cardby.firestore", l.NewValue("arguments", filter))

	if filter.Name != nil {
		cardsQuery = cardsQuery.Where("Name", ">=", *filter.Name)
	}
	for index, typeName := range filter.TypesObject {
		cardsQuery = cardsQuery.Where(fmt.Sprintf("TypesObject.%d", index), "==", typeName)
	}
	for index, cost := range filter.CostsObject {
		cardsQuery = cardsQuery.Where(fmt.Sprintf("CostsObject.%d", index), "==", cost)
	}
	for _, idSet := range filter.IDSets {
		cardsQuery = cardsQuery.Where("IDSet", "==", idSet)
	}
	if filter.IDRarity != nil {
		cardsQuery = cardsQuery.Where("IDRarity", "==", *filter.IDRarity)
	}
	if filter.Set != nil {
		var (
			setsRef                   = fmt.Sprintf(collectionFmt, "sets")
			setsQuery firestore.Query = client.Collection(setsRef)
		)
		if filter.Set.Name != nil {
			setsQuery = setsQuery.Where("Name", ">=", *filter.Set.Name)
		}
		// if filter.Set.Alias != nil {
		// 	setsQuery = setsQuery.Where("Alias", ">=", *filter.Set.Alias)
		// }
		documents, err := setsQuery.Documents(ctx).GetAll()
		if err != nil {
			logger.Error("resolve.cardby.firestore.set.query_err", l.NewValue("error", err))
			return cards, err
		}
		for index, document := range documents {
			var set Set
			err := document.DataTo(&set)
			if err != nil {
				logger.Error(
					"resolve.cardby.firestore.set.fetch_err",
					l.NewValue("index", index), l.NewValue("error", err),
				)
				return cards, err
			}
			cardsQuery = cardsQuery.Where("IDSet", "==", set.ID)
		}
	}
	if filter.Rarity != nil {
		var (
			raritiesRef                   = fmt.Sprintf(collectionFmt, "rarities")
			raritiesQuery firestore.Query = client.Collection(raritiesRef)
		)
		if filter.Rarity.Name != nil {
			raritiesQuery = raritiesQuery.Where("Name", ">=", *filter.Rarity.Name)
		}
		// if filter.Rarity.Alias != nil {
		// 	raritiesQuery = raritiesQuery.Where("Alias", ">=", *filter.Rarity.Alias)
		// }
		documents, err := raritiesQuery.Documents(ctx).GetAll()
		if err != nil {
			logger.Error("resolve.cardby.firestore.rarity.query_err", l.NewValue("error", err))
			return cards, err
		}
		for index, document := range documents {
			var rarity Rarity
			err := document.DataTo(&rarity)
			if err != nil {
				logger.Error(
					"resolve.cardby.firestore.rarity.fetch_err",
					l.NewValue("index", index),
					l.NewValue("error", err),
				)
				return cards, err
			}
			cardsQuery = cardsQuery.Where("IDRarity", "==", rarity.ID)
		}
	}
	logger.Info("resolve.cardby.firestore.query", l.NewValue("query", cardsQuery))
	documents, err := cardsQuery.Documents(ctx).GetAll()
	if err != nil {
		logger.Error("resolve.cardby.firestore.query_err", l.NewValue("error", err))
		return cards, err
	}
	for index, document := range documents {
		var card Card
		err := document.DataTo(&card)
		if err != nil {
			logger.Error(
				"resolve.cardby.firestore.fetch_err",
				l.NewValue("index", index),
				l.NewValue("error", err),
			)
			return cards, err
		}
		cards = append(cards, card)
	}

	logger.Info("resolve.cardby.firestore.fetched", l.NewValue("cards.len", len(cards)))
	return cards, err
}

func (r *queryResolver) Set(ctx context.Context, id string) (Set, error) {
	var (
		client = firestore.MustReference(r.tree)
		logger = l.MustReference(r.tree)
		setRef = fmt.Sprintf(collectionFmt, fmt.Sprintf("sets/%s", id))
		set    Set
	)
	logger.Debug("resolve.set.firestore", l.NewValue("arguments", id))

	document, err := client.Doc(setRef).Get(ctx)
	if err != nil {
		return set, err
	}

	err = document.DataTo(&set)
	if err != nil {
		return set, err
	}

	logger.Debug("resolve.set.firestore.fetched", l.NewValue("set", set))
	return set, err
}

func (r *queryResolver) SetBy(ctx context.Context, filter SetFilter) ([]Set, error) {
	var (
		logger = l.MustReference(r.tree)
		sets   []Set
		err    error
	)

	logger.Info("resolve.setby.firestore", l.NewValue("arguments", filter))

	sets = []Set{
		{
			ID:    "",
			Name:  "",
			Alias: "",
			Asset: map[string]interface{}{
				"": nil,
			},
			Cards:     nil,
			CreatedAt: time.Time{},
			UpdatedAt: &time.Time{},
			DeletedAt: &time.Time{},
		},
	}

	logger.Info("resolve.setby.firestore.fetched", l.NewValue("sets.len", len(sets)))
	return sets, err
}
