package graphql

import (
	"errors"

	"github.com/99designs/gqlgen/graphql"
	"github.com/rjansen/yggdrasil"
)

var (
	ErrInvalidReference = errors.New("Invalid Schema Reference")
	schemaPath          = yggdrasil.NewPath("/fivecolors/core/graphql/schema")
)

func Register(roots *yggdrasil.Roots, schema graphql.ExecutableSchema) error {
	return roots.Register(schemaPath, schema)
}

func Reference(tree yggdrasil.Tree) (graphql.ExecutableSchema, error) {
	var (
		schema         graphql.ExecutableSchema
		reference, err = tree.Reference(schemaPath)
	)
	if err != nil {
		return schema, err
	}
	tmpSchema, is := reference.(graphql.ExecutableSchema)
	if !is {
		return schema, ErrInvalidReference
	}
	schema = tmpSchema
	return schema, nil
}

func MustReference(tree yggdrasil.Tree) graphql.ExecutableSchema {
	schema, err := Reference(tree)
	if err != nil {
		panic(err)
	}
	return schema
}
