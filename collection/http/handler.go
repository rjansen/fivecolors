package http

import (
	"net/http"

	"github.com/99designs/gqlgen/handler"
	"github.com/rjansen/fivecolors/collection/graphql"
)

func New() http.Handler {
	return handler.GraphQL(graphql.NewExecutableSchema(graphql.Config{Resolvers: &graphql.Resolver{}}))
}
