package server

import (
	"net/http"
	"sync"

	"github.com/rjansen/fivecolors/core/api"
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
			serverHandler = api.NewGraphQLHandler(newTree())
		},
	)
	serverHandler(w, r)
}
