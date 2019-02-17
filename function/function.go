package server

import (
	stdsql "database/sql"
	"net/http"
	"sync"

	"github.com/rjansen/fivecolors/core/api"
	"github.com/rjansen/fivecolors/core/graphql"
	"github.com/rjansen/l"
	"github.com/rjansen/migi"
	"github.com/rjansen/raizel/sql"
	"github.com/rjansen/yggdrasil"
)

var (
	once          sync.Once
	serverHandler http.HandlerFunc
)

type options struct {
	driver string
	dsn    string
}

func newOptions() options {
	var (
		env     = migi.NewOptions(migi.NewEnvironmentSource())
		options options
	)
	env.StringVar(
		&options.driver, "raizel_driver", "postgres", "Raizel database driver",
	)
	env.StringVar(
		&options.dsn,
		"raizel_dsn",
		"postgres://postgres:@127.0.0.1:5432/postgres?sslmode=disable",
		"Raizel datasource name format",
	)
	env.Parse()
	return options
}

func newTree() yggdrasil.Tree {
	var (
		options   = newOptions()
		logger    = l.NewZapLoggerDefault()
		roots     = yggdrasil.NewRoots()
		err       error
		sqlDriver string
		sqlDsn    string
	)

	err = l.Register(&roots, logger)
	if err != nil {
		panic(err)
	}
	sqlDB, err := stdsql.Open(options.driver, options.dsn)
	if err != nil {
		panic(err)
	}
	db, err := sql.NewDB(sqlDB)
	if err != nil {
		panic(err)
	}
	err = sql.Register(&roots, db)
	if err != nil {
		panic(err)
	}
	err = graphql.Register(&roots, NewSchema(roots.NewTreeDefault()))
	if err != nil {
		panic(err)
	}
	return roots.NewTreeDefault()
}

func Handler(w http.ResponseWriter, r *http.Request) {
	once.Do(
		func() {
			serverHandler = api.NewGraphQLHandler(newTree())
		},
	)
	serverHandler(w, r)
}
