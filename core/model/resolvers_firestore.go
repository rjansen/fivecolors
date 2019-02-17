// +build firestore

package model

import (
	"context"
	"time"

	"github.com/rjansen/l"
	"github.com/rjansen/yggdrasil"
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
		logger = l.MustReference(r.tree)
		card   Card
		err    error
	)
	logger.Debug("resolve.card.firestore", l.NewValue("arguments", id))

	card = Card{
		ID:            "",
		Name:          "",
		Types:         nil,
		Costs:         nil,
		NumberCost:    0.0,
		IDExternal:    "",
		OrderExternal: nil,
		IDRarity:      "",
		Rarity: Rarity{
			ID:    "",
			Name:  "",
			Alias: "",
		},
		IDSet: "",
		Set: Set{
			ID:    "",
			Name:  "",
			Alias: "",
			Asset: map[string]interface{}{
				RarityNameCommon.String():     "common_set_id_asset",
				RarityNameUncommon.String():   "uncommon_set_id_asset",
				RarityNameRare.String():       "rare_set_id_asset",
				RarityNameMythicRare.String(): "mythicrare_set_id_asset",
			},
			CreatedAt: time.Now().UTC(),
		},
		IDAsset:   "",
		Rules:     nil,
		Rate:      nil,
		RateVotes: nil,
		Artist:    nil,
		Flavor:    nil,
		Data:      nil,
		CreatedAt: time.Now().UTC(),
	}

	logger.Debug("resolve.card.firestore.fetched", l.NewValue("card", card))
	return card, err
}

func (r *queryResolver) CardBy(ctx context.Context, filter CardFilter) ([]Card, error) {
	var (
		logger = l.MustReference(r.tree)
		cards  []Card
		err    error
	)
	logger.Info("resolve.cardby.firestore", l.NewValue("arguments", filter))

	cards = []Card{
		{
			ID:            "",
			Name:          "",
			Types:         nil,
			Costs:         nil,
			NumberCost:    0.0,
			IDExternal:    "",
			OrderExternal: nil,
			IDRarity:      "",
			Rarity: Rarity{
				ID:    "",
				Name:  "",
				Alias: "",
			},
			IDSet: "",
			Set: Set{
				ID:    "",
				Name:  "",
				Alias: "",
				Asset: map[string]interface{}{
					RarityNameCommon.String():     "common_set_id_asset",
					RarityNameUncommon.String():   "uncommon_set_id_asset",
					RarityNameRare.String():       "rare_set_id_asset",
					RarityNameMythicRare.String(): "mythicrare_set_id_asset",
				},
				CreatedAt: time.Now().UTC(),
			},
			IDAsset:   "",
			Rules:     nil,
			Rate:      nil,
			RateVotes: nil,
			Artist:    nil,
			Flavor:    nil,
			Data:      nil,
			CreatedAt: time.Now().UTC(),
		},
	}

	logger.Info("resolve.cardby.firestore.fetched", l.NewValue("cards.len", len(cards)))
	return cards, err
}

func (r *queryResolver) Set(ctx context.Context, id string) (Set, error) {
	var (
		logger = l.MustReference(r.tree)
		set    Set
		err    error
	)
	logger.Debug("resolve.set.firestore", l.NewValue("arguments", id))

	set = Set{
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
