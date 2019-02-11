package model

import (
	"github.com/rjansen/fivecolors/core/errors"
	"github.com/rjansen/fivecolors/core/graphql"
	"github.com/rjansen/yggdrasil"
)

var (
	ErrInvalidState = errors.New("ErrInvalidState")
)

func NewSchema(tree yggdrasil.Tree) graphql.Schema {
	return graphql.NewSchema(
		NewExecutableSchema(
			Config{
				Resolvers: NewResolver(tree),
			},
		),
	)
}

/*
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
*/
