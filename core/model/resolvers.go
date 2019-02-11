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
	/*
	   types varchar(128)[] not null,
	   costs varchar(5)[] not null default '{}',
	   rules varchar(512)[] not null default '{}',
	*/
	var (
		db     = sql.MustReference(r.tree)
		logger = l.MustReference(r.tree)
		q      = `select id, name, number_cost, id_external, id_rarity, id_set, id_asset,
                         rate, rate_votes, order_external, artist, flavor, data,
					     created_at, updated_at, deleted_at
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
		&card.CreatedAt, &updatedAt, &deletedAt,
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

func (r *queryResolver) CardBy(ctx context.Context, filter CardFilter) ([]Card, error) {
	var (
		db     = sql.MustReference(r.tree)
		logger = l.MustReference(r.tree)
		s      = `select id, name, number_cost, id_external, id_rarity, id_set, id_asset,
                         rate, rate_votes, order_external, artist, flavor, data,
					     created_at, updated_at, deleted_at
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
		a = append(a, *filter.Types)
	}
	if filter.Costs != nil {
		w = append(w, fmt.Sprintf("array_to_string(costs, '') ~* $%d", len(a)+1))
		a = append(a, *filter.Costs)
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
			&card.CreatedAt, &updatedAt, &deletedAt,
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

func (r *queryResolver) Set(ctx context.Context, id string) (Set, error) {
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

func (r *queryResolver) SetBy(ctx context.Context, filter SetFilter) ([]Set, error) {
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
