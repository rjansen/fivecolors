package function

import (
	stdsql "database/sql"
	"fmt"
	"net/http"
	"sync"

	"github.com/rjansen/fivecolors/core/api"
	"github.com/rjansen/fivecolors/core/graphql"
	"github.com/rjansen/fivecolors/core/model"
	"github.com/rjansen/l"
	"github.com/rjansen/migi"
	"github.com/rjansen/raizel/firestore"
	"github.com/rjansen/raizel/sql"
	"github.com/rjansen/yggdrasil"
)

var (
	once          sync.Once
	serverHandler http.HandlerFunc
)

type options struct {
	projectID string
	dataStore string
	driver    string
	dsn       string
}

func newOptions() options {
	var (
		env     = migi.NewOptions(migi.NewEnvironmentSource())
		options options
	)
	env.StringVar(
		&options.projectID, "project_id", "project-id", "GCP project identifier",
	)
	env.StringVar(
		&options.dataStore, "data_store", "firestore", "Persistence data store",
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

func newTree(options options) yggdrasil.Tree {
	var (
		logger = l.NewZapLoggerDefault()
		roots  = yggdrasil.NewRoots()
		err    error
	)

	err = l.Register(&roots, logger)
	if err != nil {
		panic(err)
	}

	err = graphql.Register(&roots, newSchema(&roots, options))
	if err != nil {
		panic(err)
	}
	return roots.NewTreeDefault()
}

func newSchema(roots *yggdrasil.Roots, options options) graphql.Schema {
	var queryResolver model.QueryResolver
	switch options.dataStore {
	case "firestore":
		err := firestore.Register(roots, newFirestoreClient(options))
		if err != nil {
			panic(err)
		}
		queryResolver = model.NewFirestoreQueryResolver(
			roots.NewTreeDefault(),
		)
	case "postgres":
		err := sql.Register(roots, newSqlDB(options))
		if err != nil {
			panic(err)
		}
		queryResolver = model.NewPostgresQueryResolver(
			roots.NewTreeDefault(),
		)
	default:
		panic(
			fmt.Sprintf("invalid_datastore: datastore=%s valid_values=[sql,firestore]", options.dataStore),
		)
	}
	return model.NewSchema(
		model.NewResolver(queryResolver),
	)
}

func newFirestoreClient(options options) firestore.Client {
	return firestore.NewClient(options.projectID)
}

func newSqlDB(options options) sql.DB {
	sqlDB, err := stdsql.Open(options.driver, options.dsn)
	if err != nil {
		panic(err)
	}
	db, err := sql.NewDB(sqlDB)
	if err != nil {
		panic(err)
	}
	return db
}

func Handler(w http.ResponseWriter, r *http.Request) {
	once.Do(
		func() {
			if serverHandler == nil {
				serverHandler = api.NewGraphQLHandler(
					newTree(newOptions()),
				)
			}
		},
	)
	serverHandler(w, r)
}
