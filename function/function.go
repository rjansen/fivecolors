package server

import (
	"net/http"
	"sync"

	"github.com/99designs/gqlgen/handler"
	"github.com/rjansen/fivecolors/core/api"
	"github.com/rjansen/fivecolors/core/model"
	"github.com/rjansen/l"
	"github.com/rjansen/migi"
	"github.com/rjansen/yggdrasil"
)

var (
	once          sync.Once
	serverHandler http.HandlerFunc
)

func newTree() yggdrasil.Tree {
	schema, err := api.NewMockSchema()
	if err != nil {
		panic(err)
	}
	roots := yggdrasil.NewRoots()
	err = migi.Register(&roots, migi.NewOptions(migi.NewEnvironmentSource()))
	if err != nil {
		panic(err)
	}
	err = l.Register(&roots, l.NewZapLoggerDefault())
	if err != nil {
		panic(err)
	}
	err = api.Register(&roots, schema)
	if err != nil {
		panic(err)
	}
	return roots.NewTreeDefault()
}

func Handler(w http.ResponseWriter, r *http.Request) {
	once.Do(
		func() {
			// serverHandler = api.NewGraphQLHandler(newTree())
			serverHandler = handler.GraphQL(
				model.NewExecutableSchema(
					model.Config{
						Resolvers: model.NewResolver(),
					},
				),
			)
		},
	)
	serverHandler(w, r)
}
