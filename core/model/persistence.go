package model

import (
	"context"
	"database/sql"

	"github.com/rjansen/fivecolors/core/errors"
	"github.com/rjansen/fivecolors/core/util"
	"github.com/rjansen/fivecolors/core/validator"
)

var (
	ErrInit         = errors.New("model.ErrInit")
	ErrNotFound     = errors.New("model.ErrNotFound")
	ErrBlankContext = errors.New("model.ErrBlankContext")
	db              *sql.DB
)

func Init() error {
	var (
		driver = util.Getenv("DB_DRIVER", "postgres")
		dsn    = util.Getenv("DB_DSN", "postgres://postgres:@127.0.0.1:5432/fivecolors?sslmode=disable")
	)
	err := validator.Validate(
		validator.ValidateAll(
			validator.ValidateIsBlank(driver),
			validator.ValidateIsIn(driver, "postgres"),
		),
		validator.ValidateIsBlank(dsn),
	)
	if err != nil {
		return err
	}
	tmpDB, err := sql.Open(driver, dsn)
	if err != nil {
		return errors.ErrorWrap(ErrInit, err)
	}
	db = tmpDB
	return nil
}

type Context struct {
	*util.Context
	*sql.DB
}

func MakeContext() (*Context, error) {
	return NewContext(context.Background())
}

func NewContext(c context.Context) (*Context, error) {
	return CreateContext(util.NewContext(c))
}

func CreateContext(c *util.Context) (*Context, error) {
	c.Debug().Msg("model.newcontext.start")
	if err := db.Ping(); err != nil {
		c.Error().Err(err).Msg("model.newcontext.err")
		return nil, err
	}
	c.Debug().Msg("model.newcontext.dbok")
	return &Context{
		DB:      db,
		Context: c,
	}, nil
}

type Scanner interface {
	Scan(...interface{}) error
}

type Iterator interface {
	Scanner
	Next() bool
}

type BulkIterator func() []interface{}

type QueryResult interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}

type Scannable func(*Context, Scanner) error
type Iterable func(*Context, Iterator) error

func Query(c *Context, q string, i Iterable, p ...interface{}) error {
	if c == nil {
		return ErrBlankContext
	}
	rows, err := c.Query(q, p...)
	if err != nil {
		c.Error().Str("query", q).Interface("params", p).Err(err).Msg("model.query.err")
		return err
	}
	return i(c, rows)
}

func QueryOne(c *Context, q string, s Scannable, p ...interface{}) error {
	if c == nil {
		return ErrBlankContext
	}
	err := s(c, c.QueryRow(q, p...))
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		c.Error().Str("query", q).Interface("params", p).Err(err).Msg("model.queryone.err")
		return err
	}
	return nil
}

func Insert(c *Context, q string, p ...interface{}) error {
	if c == nil {
		return ErrBlankContext
	}
	_, err := c.Exec(q, p...)
	if err != nil {
		c.Error().Str("query", q).Interface("params", p).Err(err).Msg("model.insert.err")
		return errors.Wrap(err, "model.insert.err")
	}
	return nil
}

func Delete(c *Context, q string, p ...interface{}) error {
	if c == nil {
		return ErrBlankContext
	}
	_, err := c.Exec(q, p...)
	if err != nil {
		c.Error().Str("query", q).Interface("params", p).Err(err).Msg("model.delete.err")
		return errors.Wrap(err, "model.delete.err")
	}
	return nil
}

func Bulk(c *Context, q string, p ...[]interface{}) error {
	if c == nil {
		return ErrBlankContext
	}
	tx, err := c.Begin()
	if err != nil {
		return err
	}
	stm, err := tx.Prepare(q)
	if err != nil {
		return err
	}
	for _, v := range p {
		_, err := stm.Exec(v...)
		if err != nil {
			c.Error().Str("query", q).Interface("params", v).Err(err).Msg("model.bulk.err")
			return errors.ErrorWrap(tx.Rollback(), errors.Wrap(err, "model.bulk.err"))
		}
	}
	if err := tx.Commit(); err != nil {
		c.Error().Str("query", q).Int("params", len(p)).Err(err).Msg("model.bulk.commit.err")
		return err
	}
	return nil
}
