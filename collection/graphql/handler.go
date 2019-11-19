package graphql

import (
	"net/http"

	"github.com/99designs/gqlgen/handler"
	"github.com/rjansen/fivecolors/collection"
)

func NewHandler(reader collection.Reader, writer collection.Writer) http.Handler {
	resolver := NewResolver(reader, writer)
	schema := NewExecutableSchema(Config{Resolvers: resolver})
	return handler.GraphQL(schema)
}
