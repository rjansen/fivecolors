package model

import (
	"context"
	stdsql "database/sql"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/rjansen/l"
	"github.com/rjansen/raizel"
	"github.com/rjansen/raizel/sql"
	"github.com/rjansen/yggdrasil"
)

type postgresQueryResolver struct {
	tree yggdrasil.Tree
}

func NewPostgresQueryResolver(tree yggdrasil.Tree) QueryResolver {
	return &postgresQueryResolver{
		tree: tree,
	}
}

func (r *postgresQueryResolver) Card(ctx context.Context, id string) (Card, error) {
	var (
		db     = sql.MustReference(r.tree)
		logger = l.MustReference(r.tree)
		q      = `select id, name, number_cost, id_external, id_rarity, id_set, id_asset,
                         rate, rate_votes, order_external, artist, flavor, data,
					     types, costs, rules, created_at, updated_at, deleted_at
				  from card where id = $1`
		card                          Card
		deletedAt, updatedAt          pq.NullTime
		orderExternal, artist, flavor stdsql.NullString
	)
	logger.Debug("schema.resolve.card.sql", l.NewValue("sql", q), l.NewValue("arguments", id))

	row := db.QueryRow(q, id)
	err := row.Scan(
		&card.ID, &card.Name, &card.NumberCost, &card.IDExternal, &card.IDRarity, &card.IDSet,
		&card.IDAsset, &card.Rate, &card.RateVotes, &orderExternal, &artist, &flavor, &card.Data,
		pq.Array(&card.Types), pq.Array(&card.Costs), pq.Array(&card.Rules), &card.CreatedAt,
		&updatedAt, &deletedAt,
	)
	if err != nil {
		if err == stdsql.ErrNoRows {
			return card, raizel.ErrNotFound
		}
		logger.Error("schema.resolve.card.err_fetch", l.NewValue("error", err))
		return card, err
	}
	if orderExternal.Valid {
		card.OrderExternal = &orderExternal.String
	}
	if artist.Valid {
		card.Artist = &artist.String
	}
	if flavor.Valid {
		card.Flavor = &artist.String
	}
	if updatedAt.Valid {
		card.UpdatedAt = &updatedAt.Time
	}
	if deletedAt.Valid {
		card.DeletedAt = &deletedAt.Time
	}
	logger.Debug("schema.resolve.card.fetched", l.NewValue("card", card))
	return card, err
}

func (r *postgresQueryResolver) CardBy(ctx context.Context, filter CardFilter) ([]Card, error) {
	var (
		db     = sql.MustReference(r.tree)
		logger = l.MustReference(r.tree)
		s      = `select id, name, number_cost, id_external, id_rarity, id_set, id_asset,
                         rate, rate_votes, order_external, artist, flavor, data,
					     types, costs, rules, created_at, updated_at, deleted_at
				  from card where `
		q string
		a []interface{}
		w []string
	)
	if filter.Name != nil {
		w = append(w, fmt.Sprintf("name ~* $%d", len(a)+1))
		a = append(a, *filter.Name)
	}
	if filter.Types != nil {
		w = append(w, fmt.Sprintf("exists (select 1 from unnest(types) as type where type ~* $%d)", len(a)+1))
		a = append(a, pq.Array(filter.Types))
	}
	if filter.Costs != nil {
		w = append(w, fmt.Sprintf("array_to_string(costs, '') ~* $%d", len(a)+1))
		a = append(a, pq.Array(filter.Costs))
	}
	if filter.Set != nil {
		var (
			setq = `exists (select 1 from set where id = card.id_set and %s)`
			setw []string
			seta []interface{}
		)
		if filter.Set.Name != nil {
			setw = append(setw, fmt.Sprintf("name ~* $%d", len(a)+len(seta)+1))
			seta = append(seta, *filter.Set.Name)
		}
		if filter.Set.Alias != nil {
			setw = append(setw, fmt.Sprintf("alias ~* $%d", len(a)+len(seta)+1))
			seta = append(seta, *filter.Set.Alias)
		}
		w = append(w, fmt.Sprintf(setq, strings.Join(setw, " and ")))
		a = append(a, seta...)
	}
	if filter.Rarity != nil {
		var (
			rarityq = `exists (select 1 from rarity where id = card.id_rarity and %s)`
			rarityw []string
			raritya []interface{}
		)
		if filter.Rarity.Name != nil {
			rarityw = append(rarityw, fmt.Sprintf("name ~* $%d", len(a)+len(raritya)+1))
			raritya = append(raritya, *filter.Rarity.Name)
		}
		if filter.Rarity.Alias != nil {
			rarityw = append(rarityw, fmt.Sprintf("alias ~* $%d", len(a)+len(raritya)+1))
			raritya = append(raritya, *filter.Rarity.Alias)
		}
		w = append(w, fmt.Sprintf(rarityq, strings.Join(rarityw, " and ")))
		a = append(a, raritya...)
	}
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
			card                          Card
			deletedAt, updatedAt          pq.NullTime
			orderExternal, artist, flavor stdsql.NullString
		)
		err := rows.Scan(
			&card.ID, &card.Name, &card.NumberCost, &card.IDExternal, &card.IDRarity, &card.IDSet,
			&card.IDAsset, &card.Rate, &card.RateVotes, &orderExternal, &artist, &flavor, &card.Data,
			pq.Array(&card.Types), pq.Array(&card.Costs), pq.Array(&card.Rules), &card.CreatedAt,
			&updatedAt, &deletedAt,
		)
		if err != nil {
			logger.Error("schema.resolve.cardby.err_fetch", l.NewValue("error", err))
			return nil, err
		}
		if orderExternal.Valid {
			card.OrderExternal = &orderExternal.String
		}
		if artist.Valid {
			card.Artist = &artist.String
		}
		if flavor.Valid {
			card.Flavor = &artist.String
		}
		if updatedAt.Valid {
			card.UpdatedAt = &updatedAt.Time
		}
		if deletedAt.Valid {
			card.DeletedAt = &deletedAt.Time
		}
		cards = append(cards, card)
	}
	logger.Info("schema.resolve.cardby.result", l.NewValue("cards.len", len(cards)))
	return cards, err
}

func (r *postgresQueryResolver) Set(ctx context.Context, id string) (Set, error) {
	var (
		db     = sql.MustReference(r.tree)
		logger = l.MustReference(r.tree)
		q      = `select id, name, alias, created_at, updated_at, deleted_at
		          from set where id = $1`
		set                  Set
		deletedAt, updatedAt pq.NullTime
	)
	logger.Debug("schema.resolve.set.sql", l.NewValue("sql", q), l.NewValue("arguments", id))
	row := db.QueryRow(q, id)
	err := row.Scan(
		&set.ID, &set.Name, &set.Alias, &set.CreatedAt, &updatedAt, &deletedAt,
	)
	if err != nil {
		if err == stdsql.ErrNoRows {
			return set, raizel.ErrNotFound
		}
		logger.Error("schema.resolve.set.err_fetch", l.NewValue("error", err))
		return set, err
	}
	if updatedAt.Valid {
		set.UpdatedAt = &updatedAt.Time
	}
	if deletedAt.Valid {
		set.DeletedAt = &deletedAt.Time
	}
	logger.Debug("schema.resolve.set.fetched", l.NewValue("set", set))
	return set, err
}

func (r *postgresQueryResolver) SetBy(ctx context.Context, filter SetFilter) ([]Set, error) {
	var (
		db     = sql.MustReference(r.tree)
		logger = l.MustReference(r.tree)
		s      = `select id, name, alias, created_at, updated_at, deleted_at
		          from set where `
		q string
		a []interface{}
		w []string
	)
	if filter.Name != nil {
		w = append(w, fmt.Sprintf("name ~* $%d", len(a)+1))
		a = append(a, *filter.Name)
	}
	if filter.Alias != nil {
		w = append(w, fmt.Sprintf("alias ~* $%d", len(a)+1))
		a = append(a, *filter.Alias)
	}
	q = s + strings.Join(w, " and ")

	logger.Info("schema.resolve.setby.sql", l.NewValue("query", q), l.NewValue("arguments", a))
	var (
		sets      = make([]Set, 0, 50)
		rows, err = db.Query(q, a...)
	)
	if err != nil {
		if err == stdsql.ErrNoRows {
			return sets, nil
		}
		logger.Error("schema.resolve.setdby.err_execute_sql", l.NewValue("error", err))
		return nil, err
	}
	for rows.Next() {
		var (
			set                  Set
			deletedAt, updatedAt pq.NullTime
		)
		err := rows.Scan(
			&set.ID, &set.Name, &set.Alias, &set.CreatedAt, &updatedAt, &deletedAt,
		)
		if err != nil {
			logger.Error("schema.resolve.setby.err_fetch", l.NewValue("error", err))
			return nil, err
		}
		if updatedAt.Valid {
			set.UpdatedAt = &updatedAt.Time
		}
		if deletedAt.Valid {
			set.DeletedAt = &deletedAt.Time
		}
		sets = append(sets, set)
	}
	logger.Info("schema.resolve.setby.result", l.NewValue("sets.len", len(sets)))
	return sets, err
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
