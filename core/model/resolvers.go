package model

import (
	"context"
	"time"
)

type Resolver struct{}

func NewResolver() *Resolver {
	return new(Resolver)
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Set(ctx context.Context, id string) (Set, error) {
	return Set{
		ID:    "setone",
		Name:  "Set One",
		Alias: "DES",
		Asset: map[string]interface{}{
			RarityIdCommon.String():     "https://assets.fivecolors.mock/rarity/common",
			RarityIdUncommon.String():   "https://assets.fivecolors.mock/rarity/uncommon",
			RarityIdRare.String():       "https://assets.fivecolors.mock/rarity/rare",
			RarityIdMythicRare.String(): "https://assets.fivecolors.mock/rarity/mythicrare",
		},
		Cards: []Card{
			{
				ID:            "cardone",
				Name:          "Card One",
				Types:         []string{},
				Costs:         []string{},
				NumberCost:    1,
				IDExternal:    "cardoneexternal",
				OrderExternal: nil,
				IDRarity:      RarityIdCommon.String(),
				Rarity: Rarity{
					ID:    RarityIdCommon.String(),
					Name:  RarityIdCommon.String(),
					Alias: RarityAliasC.String(),
				},
				IDSet: "setone",
				Set: Set{
					ID:    "setone",
					Name:  "Set One",
					Alias: "DES",
					Asset: Object{
						RarityIdCommon.String():     "https://assets.fivecolors.mock/rarity/common",
						RarityIdUncommon.String():   "https://assets.fivecolors.mock/rarity/uncommon",
						RarityIdRare.String():       "https://assets.fivecolors.mock/rarity/rare",
						RarityIdMythicRare.String(): "https://assets.fivecolors.mock/rarity/mythicrare",
					},
				},

				IDAsset:   "https://assets.fivecolors.mock/card/cardone",
				Rules:     []string{},
				Rate:      nil,
				RateVotes: nil,
				Artist:    nil,
				Flavor:    nil,
				Data: &Object{
					"customkey1": "customdata1",
					"customkey2": "customdata2",
					"customkey3": "customdata3",
				},
				CreatedAt: time.Now().UTC(),
				UpdatedAt: nil,
				DeletedAt: nil,
			},
			{
				ID:            "cardtwo",
				Name:          "Card Two",
				Types:         []string{},
				Costs:         []string{},
				NumberCost:    1,
				IDExternal:    "cardtwoexternal",
				OrderExternal: nil,
				IDRarity:      RarityIdUncommon.String(),
				Rarity: Rarity{
					ID:    RarityIdUncommon.String(),
					Name:  RarityIdUncommon.String(),
					Alias: RarityAliasU.String(),
				},
				IDSet: "setone",
				Set: Set{
					ID:    "setone",
					Name:  "Set One",
					Alias: "DES",
					Asset: Object{
						RarityIdCommon.String():     "https://assets.fivecolors.mock/rarity/common",
						RarityIdUncommon.String():   "https://assets.fivecolors.mock/rarity/uncommon",
						RarityIdRare.String():       "https://assets.fivecolors.mock/rarity/rare",
						RarityIdMythicRare.String(): "https://assets.fivecolors.mock/rarity/mythicrare",
					},
				},
				IDAsset:   "https://assets.fivecolors.mock/card/cardtwo",
				Rules:     []string{},
				Rate:      nil,
				RateVotes: nil,
				Artist:    nil,
				Flavor:    nil,
				Data: &Object{
					"customkey1": "customdata1",
					"customkey2": "customdata2",
					"customkey3": "customdata3",
				},
				CreatedAt: time.Now().UTC(),
				UpdatedAt: nil,
				DeletedAt: nil,
			},
		},
	}, nil
}
func (r *queryResolver) Card(ctx context.Context, id string) (Card, error) {
	return Card{
		ID:            "cardone",
		Name:          "Card One",
		Types:         []string{},
		Costs:         []string{},
		NumberCost:    1,
		IDExternal:    "cardoneexternal",
		OrderExternal: nil,
		IDRarity:      RarityIdCommon.String(),
		Rarity: Rarity{
			ID:    RarityIdCommon.String(),
			Name:  RarityIdCommon.String(),
			Alias: RarityAliasC.String(),
		},
		IDSet: "setone",
		Set: Set{
			ID:    "setone",
			Name:  "Set One",
			Alias: "DES",
			Asset: Object{
				RarityIdCommon.String():     "https://assets.fivecolors.mock/rarity/common",
				RarityIdUncommon.String():   "https://assets.fivecolors.mock/rarity/uncommon",
				RarityIdRare.String():       "https://assets.fivecolors.mock/rarity/rare",
				RarityIdMythicRare.String(): "https://assets.fivecolors.mock/rarity/mythicrare",
			},
		},

		IDAsset:   "https://assets.fivecolors.mock/card/cardone",
		Rules:     []string{},
		Rate:      nil,
		RateVotes: nil,
		Artist:    nil,
		Flavor:    nil,
		Data: &Object{
			"customkey1": "customdata1",
			"customkey2": "customdata2",
			"customkey3": "customdata3",
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: nil,
		DeletedAt: nil,
	}, nil
}

func (r *queryResolver) SetBy(ctx context.Context, filter SetFilter) ([]Set, error) {
	return []Set{
		{
			ID:    "setone",
			Name:  "Set One",
			Alias: "STO",
			Asset: Object{
				RarityIdCommon.String():     "https://assets.fivecolors.mock/rarity/common",
				RarityIdUncommon.String():   "https://assets.fivecolors.mock/rarity/uncommon",
				RarityIdRare.String():       "https://assets.fivecolors.mock/rarity/rare",
				RarityIdMythicRare.String(): "https://assets.fivecolors.mock/rarity/mythicrare",
			},
		},
		{
			ID:    "settwo",
			Name:  "Set Two",
			Alias: "STT",
			Asset: Object{
				RarityIdCommon.String():     "https://assets.fivecolors.mock/rarity/common",
				RarityIdUncommon.String():   "https://assets.fivecolors.mock/rarity/uncommon",
				RarityIdRare.String():       "https://assets.fivecolors.mock/rarity/rare",
				RarityIdMythicRare.String(): "https://assets.fivecolors.mock/rarity/mythicrare",
			},
		},
	}, nil

}
func (r *queryResolver) CardBy(ctx context.Context, filter CardFilter) ([]Card, error) {
	return []Card{
		{
			ID:            "cardone",
			Name:          "Card One",
			Types:         []string{},
			Costs:         []string{},
			NumberCost:    1,
			IDExternal:    "cardoneexternal",
			OrderExternal: nil,
			IDRarity:      RarityIdCommon.String(),
			Rarity: Rarity{
				ID:    RarityIdCommon.String(),
				Name:  RarityIdCommon.String(),
				Alias: RarityAliasC.String(),
			},
			IDSet: "setone",
			Set: Set{
				ID:    "setone",
				Name:  "Set One",
				Alias: "DES",
				Asset: Object{
					RarityIdCommon.String():     "https://assets.fivecolors.mock/rarity/common",
					RarityIdUncommon.String():   "https://assets.fivecolors.mock/rarity/uncommon",
					RarityIdRare.String():       "https://assets.fivecolors.mock/rarity/rare",
					RarityIdMythicRare.String(): "https://assets.fivecolors.mock/rarity/mythicrare",
				},
			},

			IDAsset:   "https://assets.fivecolors.mock/card/cardone",
			Rules:     []string{},
			Rate:      nil,
			RateVotes: nil,
			Artist:    nil,
			Flavor:    nil,
			Data: &Object{
				"customkey1": "customdata1",
				"customkey2": "customdata2",
				"customkey3": "customdata3",
			},
			CreatedAt: time.Now().UTC(),
			UpdatedAt: nil,
			DeletedAt: nil,
		},
		{
			ID:            "cardtwo",
			Name:          "Card Two",
			Types:         []string{},
			Costs:         []string{},
			NumberCost:    1,
			IDExternal:    "cardtwoexternal",
			OrderExternal: nil,
			IDRarity:      RarityIdUncommon.String(),
			Rarity: Rarity{
				ID:    RarityIdUncommon.String(),
				Name:  RarityIdUncommon.String(),
				Alias: RarityAliasU.String(),
			},
			IDSet: "setone",
			Set: Set{
				ID:    "setone",
				Name:  "Set One",
				Alias: "DES",
				Asset: Object{
					RarityIdCommon.String():     "https://assets.fivecolors.mock/rarity/common",
					RarityIdUncommon.String():   "https://assets.fivecolors.mock/rarity/uncommon",
					RarityIdRare.String():       "https://assets.fivecolors.mock/rarity/rare",
					RarityIdMythicRare.String(): "https://assets.fivecolors.mock/rarity/mythicrare",
				},
			},
			IDAsset:   "https://assets.fivecolors.mock/card/cardtwo",
			Rules:     []string{},
			Rate:      nil,
			RateVotes: nil,
			Artist:    nil,
			Flavor:    nil,
			Data: &Object{
				"customkey1": "customdata1",
				"customkey2": "customdata2",
				"customkey3": "customdata3",
			},
			CreatedAt: time.Now().UTC(),
			UpdatedAt: nil,
			DeletedAt: nil,
		},
	}, nil
}
