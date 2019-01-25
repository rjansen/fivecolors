package model

import (
	"fmt"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/lib/pq"
	"github.com/rjansen/fivecolors/core/errors"
)

var (
	ErrInvalidState = errors.New("ErrInvalidState")
)

var rarityIDType = graphql.NewEnum(
	graphql.EnumConfig{
		Name: "RarityID",
		Values: graphql.EnumValueConfigMap{
			"Common": &graphql.EnumValueConfig{
				Value: NewRarityID(RarityCommon),
			},
			"Uncommon": &graphql.EnumValueConfig{
				Value: NewRarityID(RarityUncommon),
			},
			"Rare": &graphql.EnumValueConfig{
				Value: NewRarityID(RarityRare),
			},
			"MythcRare": &graphql.EnumValueConfig{
				Value: NewRarityID(RarityMythicRare),
			},
		},
	},
)

var rarityType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Rarity",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"alias": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var setAssetType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "SetAsset",
		Fields: graphql.Fields{
			"Common": &graphql.Field{
				Type: graphql.String,
			},
			"Uncommon": &graphql.Field{
				Type: graphql.String,
			},
			"Rare": &graphql.Field{
				Type: graphql.String,
			},
			"MythcRare": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var setType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Set",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"alias": &graphql.Field{
				Type: graphql.String,
			},
			"asset": &graphql.Field{
				Type: setAssetType,
			},
		},
	},
)

var cardType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Card",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"idExternal": &graphql.Field{
				Type: graphql.String,
			},
			"numberCost": &graphql.Field{
				Type: graphql.Float,
			},
			"idAsset": &graphql.Field{
				Type: graphql.String,
			},
			"types": &graphql.Field{
				Type:    graphql.NewList(graphql.String),
				Resolve: cardResolveTypes,
			},
			"costs": &graphql.Field{
				Type:    graphql.NewList(graphql.String),
				Resolve: cardResolveCosts,
			},
			"rules": &graphql.Field{
				Type:    graphql.NewList(graphql.String),
				Resolve: cardResolveRules,
			},
			"rarity": &graphql.Field{
				Type:    rarityType,
				Resolve: cardResolveRarity,
			},
			"set": &graphql.Field{
				Type:    setType,
				Resolve: cardResolveSet,
			},
		},
	},
)

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"set": &graphql.Field{
				Type: setType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: resolveSet,
			},
			"setBy": &graphql.Field{
				Type: graphql.NewList(setType),
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"alias": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: resolveSet,
			},
			"card": &graphql.Field{
				Type: cardType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: resolveCard,
			},
			"cardBy": &graphql.Field{
				Type: graphql.NewList(cardType),
				Args: graphql.FieldConfigArgument{
					"set": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"name": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"types": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"numberCost": &graphql.ArgumentConfig{
						Type: graphql.Float,
					},
					"rarity": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"costs": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"rules": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: resolveCardBy,
			},
		},
	},
)

func NewSchema() (graphql.Schema, error) {
	return graphql.NewSchema(NewSchemaConfig())
}

func NewSchemaConfig() graphql.SchemaConfig {
	return graphql.SchemaConfig{
		Query: queryType,
	}
}

func resolveSet(p graphql.ResolveParams) (interface{}, error) {
	c, err := NewContext(p.Context)
	if err != nil {
		return nil, err
	}
	var (
		s = `select id, name, alias, created_at, deleted_at from set where `
		q string
		a []interface{}
	)
	if p.Info.FieldName == "set" {
		q = s + "id = $1"
		a = append(a, p.Args["id"])
	} else {
		var w []string
		for k, v := range p.Args {
			w = append(w, fmt.Sprintf("%s = $%d", k, len(w)+1))
			a = append(a, v)
		}
		q = s + strings.Join(w, " and ")
	}
	var set Set
	fetchSet := func(m *Context, s Scanner) error {
		m.Info().Msg("api.query.fetch.set")
		var deletedAt pq.NullTime
		err := s.Scan(&set.ID, &set.Name, &set.Alias, &set.CreatedAt, &deletedAt)
		if err != nil {
			m.Error().Err(err).Msg("api.query.fetch.set.err")
			return err
		}
		if deletedAt.Valid {
			set.DeletedAt = deletedAt.Time
		}
		m.Info().Interface("set", set).Msg("api.query.fetched.set")
		return nil
	}
	c.Info().Str("query", q).Msg("api.query.set.try")
	err = QueryOne(c, q, fetchSet, a...)
	if err != nil {
		if err == ErrNotFound {
			return nil, nil
		}
		c.Error().Err(err).Str("query", q).Msg("api.query.set.err")
		return nil, err
	}
	c.Info().Err(err).Interface("set", set).Msg("api.query.set.result")
	return set, err
}

func resolveCard(p graphql.ResolveParams) (interface{}, error) {
	c, err := NewContext(p.Context)
	if err != nil {
		return nil, err
	}
	var (
		s = `select id, name, number_cost, id_external, id_asset,
					id_rarity, id_set, created_at, deleted_at
			 from card where `
		q          string
		a          []interface{}
		fieldAlias = map[string]string{
			"idExternal": "id_external",
			"idAsset":    "id_asset",
		}
	)
	if p.Info.FieldName == "card" {
		q = s + "id = $1"
		a = append(a, p.Args["id"])
	} else {
		var w []string
		for k, v := range p.Args {
			if alias, ok := fieldAlias[k]; ok {
				w = append(w, fmt.Sprintf("%s = $%d", alias, len(w)+1))
			} else {
				w = append(w, fmt.Sprintf("%s = $%d", k, len(w)+1))
			}
			a = append(a, v)
		}
		q = s + strings.Join(w, " and ")
	}
	var card Card
	fetchCard := func(m *Context, s Scanner) error {
		m.Info().Msg("api.query.fetch.card")
		var deletedAt pq.NullTime
		err := s.Scan(
			&card.ID, &card.Name, &card.NumberCost, &card.IDExternal, &card.IDAsset,
			&card.IDRarity, &card.IDSet, &card.CreatedAt, &deletedAt,
		)
		if err != nil {
			m.Error().Err(err).Msg("api.query.fetch.card.err")
			return err
		}
		if deletedAt.Valid {
			card.DeletedAt = deletedAt.Time
		}
		m.Info().Interface("card", card).Msg("api.query.fetched.card")
		return nil
	}
	c.Info().Str("query", q).Msg("api.query.card.try")
	err = QueryOne(c, q, fetchCard, a...)
	if err != nil {
		if err == ErrNotFound {
			return nil, nil
		}
		c.Error().Err(err).Str("query", q).Msg("api.query.card.err")
		return nil, err
	}
	c.Info().Err(err).Interface("card", card).Msg("api.query.card.result")
	return card, err
}

func resolveCardBy(p graphql.ResolveParams) (interface{}, error) {
	c, err := NewContext(p.Context)
	if err != nil {
		return nil, err
	}
	var (
		s = `select id, name, number_cost, id_external, id_asset,
					id_rarity, id_set, created_at, deleted_at
			 from card where `
		q          string
		a          []interface{}
		fieldAlias = map[string]func(string, int) string{
			"set": func(key string, index int) string {
				return fmt.Sprintf("id_set = $%d", index)
			},
			"name": func(key string, index int) string {
				return fmt.Sprintf("name ~* $%d", index)
			},
			"types": func(key string, index int) string {
				return fmt.Sprintf("exists (select * from unnest(types) as type where type ~* $%d)", index)
			},
			"numberCost": func(key string, index int) string {
				return fmt.Sprintf("number_cost = $%d", index)
			},
			"costs": func(key string, index int) string {
				return fmt.Sprintf("array_to_string(costs, '') ~* $%d", index)
			},
		}
	)
	var w []string
	for k, v := range p.Args {
		if alias, ok := fieldAlias[k]; ok {
			w = append(w, alias(k, len(w)+1))
			a = append(a, v)
		}
	}
	q = s + strings.Join(w, " and ")

	var cards = make([]Card, 0, 200)
	fetchCards := func(m *Context, s Iterator) error {
		m.Info().Msg("api.query.fetch.cardby")
		for s.Next() {
			var (
				card      Card
				deletedAt pq.NullTime
			)
			err := s.Scan(
				&card.ID, &card.Name, &card.NumberCost, &card.IDExternal, &card.IDAsset,
				&card.IDRarity, &card.IDSet, &card.CreatedAt, &deletedAt,
			)
			if err != nil {
				m.Error().Err(err).Msg("api.query.fetch.card.err")
				return err
			}
			if deletedAt.Valid {
				card.DeletedAt = deletedAt.Time
			}
			m.Debug().Interface("card", card).Msg("api.query.read.cardby.item")
			cards = append(cards, card)
		}
		m.Info().Int("cardBy.len", len(cards)).Msg("api.query.fetched.cardby")
		return nil
	}

	c.Info().Str("query", q).Msg("api.query.cardby.try")
	err = Query(c, q, fetchCards, a...)
	if err != nil {
		if err == ErrNotFound {
			return nil, nil
		}
		c.Error().Err(err).Str("query", q).Msg("api.query.cardby.err")
		return nil, err
	}
	c.Info().Err(err).Interface("cardBy.len", cards).Msg("api.query.cardby.result")
	return cards, err
}

func cardResolveRarity(p graphql.ResolveParams) (interface{}, error) {
	card, isCard := p.Source.(Card)
	if !isCard {
		return nil, errors.Wrap(ErrInvalidState, "Source must be a model.Card to fetch rarity")
	}
	c, err := NewContext(p.Context)
	if err != nil {
		return nil, err
	}
	var (
		q = `select id, name, alias, created_at from rarity where id = $1`
	)
	var rarity Rarity
	fetchRarity := func(m *Context, s Scanner) error {
		m.Info().Msg("api.query.fetch.rarity")
		err := s.Scan(&rarity.ID, &rarity.Name, &rarity.Alias, &rarity.CreatedAt)
		if err != nil {
			m.Error().Err(err).Msg("api.query.fetch.rarity.err")
			return err
		}
		m.Info().Interface("rarity", rarity).Msg("api.query.fetched.rarity")
		return nil
	}
	c.Info().Str("query", q).Msg("api.query.rarity.try")
	err = QueryOne(c, q, fetchRarity, card.IDRarity)
	if err != nil {
		if err == ErrNotFound {
			return nil, nil
		}
		c.Error().Err(err).Str("query", q).Msg("api.query.rarity.err")
		return nil, err
	}
	c.Info().Err(err).Interface("rarity", rarity).Msg("api.query.rarity.result")
	return rarity, err
}

func cardResolveSet(p graphql.ResolveParams) (interface{}, error) {
	card, isCard := p.Source.(Card)
	if !isCard {
		return nil, errors.Wrap(ErrInvalidState, "Source must be a model.Card to fetch costs")
	}
	c, err := NewContext(p.Context)
	if err != nil {
		return nil, err
	}
	var (
		q = `select id, name, alias, created_at from set where id = $1`
	)
	var set Set
	fetchSet := func(m *Context, s Scanner) error {
		m.Info().Msg("api.query.fetch.set")
		err := s.Scan(&set.ID, &set.Name, &set.Alias, &set.CreatedAt)
		if err != nil {
			m.Error().Err(err).Msg("api.query.fetch.set.err")
			return err
		}
		m.Info().Interface("set", set).Msg("api.query.fetched.set")
		return nil
	}
	c.Info().Str("query", q).Msg("api.query.set.try")
	err = QueryOne(c, q, fetchSet, card.IDSet)
	if err != nil {
		if err == ErrNotFound {
			return nil, nil
		}
		c.Error().Err(err).Str("query", q).Msg("api.query.set.err")
		return nil, err
	}
	c.Info().Err(err).Interface("set", set).Msg("api.query.set.result")
	return set, err
}

func cardResolveTypes(p graphql.ResolveParams) (interface{}, error) {
	card, isCard := p.Source.(Card)
	if !isCard {
		return nil, errors.Wrap(ErrInvalidState, "Source must be a model.Card to fetch costs")
	}
	c, err := NewContext(p.Context)
	if err != nil {
		return nil, err
	}
	var (
		q = `select types from card where id = $1`
	)
	var types []string
	fetchTypes := func(m *Context, s Scanner) error {
		m.Info().Msg("api.query.fetch.types")
		err := s.Scan(pq.Array(&types))
		if err != nil {
			m.Error().Err(err).Msg("api.query.fetch.types.err")
			return err
		}
		m.Info().Interface("types", types).Msg("api.query.fetched.types")
		return nil
	}
	c.Info().Str("query", q).Msg("api.query.types.try")
	err = QueryOne(c, q, fetchTypes, card.ID)
	if err != nil {
		if err == ErrNotFound {
			return nil, nil
		}
		c.Error().Err(err).Str("query", q).Msg("api.query.types.err")
		return nil, err
	}
	c.Info().Err(err).Interface("types", types).Msg("api.query.types.result")
	return types, err
}

func cardResolveCosts(p graphql.ResolveParams) (interface{}, error) {
	card, isCard := p.Source.(Card)
	if !isCard {
		return nil, errors.Wrap(ErrInvalidState, "Source must be a model.Card to fetch costs")
	}
	c, err := NewContext(p.Context)
	if err != nil {
		return nil, err
	}
	var (
		q = `select costs from card where id = $1`
	)
	var costs []string
	fetchCosts := func(m *Context, s Scanner) error {
		m.Info().Msg("api.query.fetch.costs")
		err := s.Scan(pq.Array(&costs))
		if err != nil {
			m.Error().Err(err).Msg("api.query.fetch.costs.err")
			return err
		}
		m.Info().Interface("costs", costs).Msg("api.query.fetched.costs")
		return nil
	}
	c.Info().Str("query", q).Msg("api.query.costs.try")
	err = QueryOne(c, q, fetchCosts, card.ID)
	if err != nil {
		if err == ErrNotFound {
			return nil, nil
		}
		c.Error().Err(err).Str("query", q).Msg("api.query.costs.err")
		return nil, err
	}
	c.Info().Err(err).Interface("costs", costs).Msg("api.query.costs.result")
	return costs, err
}

func cardResolveRules(p graphql.ResolveParams) (interface{}, error) {
	card, isCard := p.Source.(Card)
	if !isCard {
		return nil, errors.Wrap(ErrInvalidState, "Source must be a model.Card to fetch rules")
	}
	c, err := NewContext(p.Context)
	if err != nil {
		return nil, err
	}
	var (
		q = `select rules from card where id = $1`
	)
	var rules []string
	fetchRules := func(m *Context, s Scanner) error {
		m.Info().Msg("api.query.fetch.rules")
		err := s.Scan(pq.Array(&rules))
		if err != nil {
			m.Error().Err(err).Msg("api.query.fetch.rules.err")
			return err
		}
		m.Info().Interface("rules", rules).Msg("api.query.fetched.rules")
		return nil
	}
	c.Info().Str("query", q).Msg("api.query.rules.try")
	err = QueryOne(c, q, fetchRules, card.ID)
	if err != nil {
		if err == ErrNotFound {
			return nil, nil
		}
		c.Error().Err(err).Str("query", q).Msg("api.query.rules.err")
		return nil, err
	}
	c.Info().Err(err).Interface("rules", rules).Msg("api.query.rules.result")
	return rules, err
}
