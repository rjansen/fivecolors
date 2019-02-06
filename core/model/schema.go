package model

import (
	"context"
	stdsql "database/sql"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/lib/pq"
	"github.com/rjansen/fivecolors/core/errors"
	"github.com/rjansen/l"
	"github.com/rjansen/raizel"
	"github.com/rjansen/raizel/sql"
	"github.com/rjansen/yggdrasil"
	"github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/parser"
	"github.com/vektah/gqlparser/validator"
)

var (
	ErrInvalidState = errors.New("ErrInvalidState")
)

type Params struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}

func Execute(tree yggdrasil.Tree, params Params) *graphql.Response {
	response := new(graphql.Response)

	doc, parserErr := parser.ParseQuery(&ast.Source{Input: params.Query})
	if parserErr != nil {
		response.Errors = append(response.Errors, parserErr)
		return response
	}

	schema := NewExecutableSchema(Config{Resolvers: NewResolver()})
	validateErrs := validator.Validate(schema.Schema(), doc)
	if validateErrs != nil {
		response.Errors = append(response.Errors, validateErrs...)
		return response
	}

	op := doc.Operations.ForName(params.OperationName)
	vars, varsErr := validator.VariableValues(schema.Schema(), op, params.Variables)
	if varsErr != nil {
		response.Errors = append(response.Errors, varsErr)
		return response
	}

	ctx := graphql.WithRequestContext(
		context.Background(),
		graphql.NewRequestContext(doc, params.Query, vars),
	)
	return schema.Query(ctx, op)
}

func resolveSet(ctx context.Context, id string) (*Set, error) {
	c, err := NewContext(ctx)
	if err != nil {
		return nil, err
	}
	var (
		s = `select id, name, alias, created_at, deleted_at from set where `
		q string
		a []interface{}
	)
	q = s + "id = $1"
	a = append(a, id)
	/*
		var w []string
		for k, v := range p.Args {
			w = append(w, fmt.Sprintf("%s = $%d", k, len(w)+1))
			a = append(a, v)
		}
		q = s + strings.Join(w, " and ")
	*/
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
			set.DeletedAt = &deletedAt.Time
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
	return &set, err
}

func resolveCard(tree yggdrasil.Tree, id string) (*Card, error) {
	var (
		db     = sql.MustReference(tree)
		logger = l.MustReference(tree)
		s      = `select id, name, number_cost, id_external, id_asset,
					id_rarity, id_set, created_at, deleted_at
			 from card where `
		q string
		a []interface{}
		// fieldAlias = map[string]string{
		// 	 "idExternal": "id_external",
		// 	 "idAsset":    "id_asset",
		// }
	)
	q = s + "id = $1"
	a = append(a, id)
	/*
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
	*/
	logger.Debug("schema.resolve.card.sql", l.NewValue("sql", q), l.NewValue("arguments", a))
	var (
		card      Card
		deletedAt pq.NullTime
		row       = db.QueryRow(q, a...)
		err       = row.Scan(
			&card.ID, &card.Name, &card.NumberCost, &card.IDExternal, &card.IDAsset,
			&card.IDRarity, &card.IDSet, &card.CreatedAt, &deletedAt,
		)
	)
	if err != nil {
		if err == stdsql.ErrNoRows {
			return nil, raizel.ErrNotFound
		}
		logger.Error("schema.resolve.card.err_fetch", l.NewValue("error", err))
		return nil, err
	}
	if deletedAt.Valid {
		card.DeletedAt = &deletedAt.Time
	}
	logger.Debug("schema.resolve.card.fetched", l.NewValue("card", card))
	return &card, err
}

func resolveCardBy(tree yggdrasil.Tree, filter CardFilter) ([]Card, error) {
	var (
		db     = sql.MustReference(tree)
		logger = l.MustReference(tree)
		s      = `select id, name, number_cost, id_external, id_asset,
					id_rarity, id_set, created_at, deleted_at
			 from card where `
		q string
		a []interface{}
		w []string
		// fieldAlias = map[string]func(string, int) string{
		// 	"set": func(key string, index int) string {
		// 		return fmt.Sprintf("id_set = $%d", index)
		// 	},
		// 	"name": func(key string, index int) string {
		// 		return fmt.Sprintf("name ~* $%d", index)
		// 	},
		// 	"types": func(key string, index int) string {
		// 		return fmt.Sprintf("exists (select * from unnest(types) as type where type ~* $%d)", index)
		// 	},
		// 	"numberCost": func(key string, index int) string {
		// 		return fmt.Sprintf("number_cost = $%d", index)
		// 	},
		// 	"costs": func(key string, index int) string {
		// 		return fmt.Sprintf("array_to_string(costs, '') ~* $%d", index)
		// 	},
		// }
	)
	/*
		for k, v := range p.Args {
			if alias, ok := fieldAlias[k]; ok {
				w = append(w, alias(k, len(w)+1))
				a = append(a, v)
			}
		}
	*/
	q = s + strings.Join(w, " and ")

	logger.Info("schema.resolve.cardby.sql", l.NewValue("query", q), l.NewValue("arguments", a))
	var (
		cards     = make([]Card, 0, 200)
		rows, err = db.Query(q, a...)
	)
	if err != nil {
		if err == stdsql.ErrNoRows {
			return cards, nil
		}
		logger.Error("schema.resolve.cardby.err_execute_sql", l.NewValue("error", err))
		return nil, err
	}
	for rows.Next() {
		var (
			card      Card
			deletedAt pq.NullTime
		)
		err := rows.Scan(
			&card.ID, &card.Name, &card.NumberCost, &card.IDExternal, &card.IDAsset,
			&card.IDRarity, &card.IDSet, &card.CreatedAt, &deletedAt,
		)
		if err != nil {
			logger.Error("schema.resolve.cardby.err_fetch", l.NewValue("error", err))
			return nil, err
		}
		if deletedAt.Valid {
			card.DeletedAt = &deletedAt.Time
		}
		cards = append(cards, card)
	}
	logger.Info("schema.resolve.cardby.result", l.NewValue("cards.len", len(cards)))
	return cards, err
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
