package server

import (
	"net/http"
	"sync"
)

var (
	once          sync.Once
	serverHandler http.HandlerFunc
)

/*
func newTree() yggdrasil.Tree {
	schema := mockschema.New()
	roots := yggdrasil.NewRoots()
	err := migi.Register(&roots, migi.NewOptions(migi.NewEnvironmentSource()))
	if err != nil {
		panic(err)
	}
	err = l.Register(&roots, l.NewZapLoggerDefault())
	if err != nil {
		panic(err)
	}
	err = graphql.Register(&roots, schema)
	if err != nil {
		panic(err)
	}
	return roots.NewTreeDefault()
}
*/

func Handler(w http.ResponseWriter, r *http.Request) {
	// once.Do(
	// 	func() {
	// 		serverHandler = api.NewGraphQLHandler(newTree())
	// 	},
	// )
	// serverHandler(w, r)
}
