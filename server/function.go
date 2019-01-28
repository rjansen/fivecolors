package server

import (
	"net/http"
	"sync"

	"github.com/rjansen/fivecolors/core/api"
	"github.com/rjansen/fivecolors/core/config"
	"github.com/rjansen/fivecolors/core/model"
	"github.com/rs/zerolog/log"
)

var (
	version       string
	once          sync.Once
	serverHandler http.HandlerFunc
)

func newTree(schema graphql.Schema) yggdrasil.Tree {
	var (
		roots = yggdrasil.NewRoots()
		err = api.Register(&roots, schema)
	)
	if err != nil {
		panic(err)
	}
	return roots.NewTreeDefault()
}

func Handler(w http.ResponseWriter, r *http.Request) {
	once.Do(
		func() {
			serverHandler = api.NewGraphQLHandler(
				newTree(api.NewMockSchema())
			)
		},
	)
	serverHandler(w, r)
}
