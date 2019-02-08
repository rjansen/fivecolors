package mockschema

import (
	"github.com/99designs/gqlgen/graphql"
)

func New() graphql.ExecutableSchema {
	return NewExecutableSchema(
		Config{
			Resolvers: NewResolver(),
		},
	)
}
