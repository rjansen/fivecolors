package model

import (
	"context"
	"fmt"

	"github.com/rjansen/l"
	"github.com/rjansen/raizel/firestore"
	"github.com/rjansen/yggdrasil"
)

var (
	collectionFmt = "environments/development/%s"
)

type firestoreQueryResolver struct {
	tree yggdrasil.Tree
}

func NewFirestoreQueryResolver(tree yggdrasil.Tree) QueryResolver {
	return &firestoreQueryResolver{
		tree: tree,
	}
}

func (r *firestoreQueryResolver) Card(ctx context.Context, id string) (Card, error) {
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

func (r *firestoreQueryResolver) CardBy(ctx context.Context, filter CardFilter) ([]Card, error) {
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

func (r *firestoreQueryResolver) Set(ctx context.Context, id string) (Set, error) {
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

func (r *firestoreQueryResolver) SetBy(ctx context.Context, filter SetFilter) ([]Set, error) {
	var (
		client    = firestore.MustReference(r.tree)
		logger    = l.MustReference(r.tree)
		setsRef   = fmt.Sprintf(collectionFmt, "sets")
		sets      []Set
		setsQuery firestore.Query = client.Collection(setsRef)
	)
	logger.Info("resolve.setby.firestore", l.NewValue("arguments", filter))

	if filter.Name != nil {
		setsQuery = setsQuery.Where("Name", ">=", *filter.Name)
	}
	if filter.Alias != nil {
		setsQuery = setsQuery.Where("Alias", "==", *filter.Alias)
	}
	setsQuery = setsQuery.OrderBy("Name", firestore.Asc)

	logger.Info("resolve.setdby.firestore.query", l.NewValue("query", setsQuery))
	documents, err := setsQuery.Documents(ctx).GetAll()
	if err != nil {
		logger.Error("resolve.setby.firestore.query_err", l.NewValue("error", err))
		return sets, err
	}
	for index, document := range documents {
		var set Set
		err := document.DataTo(&set)
		if err != nil {
			logger.Error(
				"resolve.setby.firestore.fetch_err",
				l.NewValue("index", index),
				l.NewValue("error", err),
			)
			return sets, err
		}
		sets = append(sets, set)
	}

	logger.Info("resolve.setby.firestore.fetched", l.NewValue("sets.len", len(sets)))
	return sets, err
}
