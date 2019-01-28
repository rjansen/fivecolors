package api

import (
	"errors"

	"github.com/graphql-go/graphql"
	"github.com/rjansen/yggdrasil"
)

var (
	ErrInvalidReference = errors.New("Invalid Schema Reference")
	schemaPath          = yggdrasil.NewPath("/fivecolors/core/api/schema")
)

func Register(roots *yggdrasil.Roots, schema graphql.Schema) error {
	return roots.Register(schemaPath, schema)
}

func Reference(tree yggdrasil.Tree) (graphql.Schema, error) {
	var (
		schema         graphql.Schema
		reference, err = tree.Reference(schemaPath)
	)
	if err != nil {
		return schema, err
	}
	tmpSchema, is := reference.(graphql.Schema)
	if !is {
		return schema, ErrInvalidReference
	}
	schema = tmpSchema
	return schema, nil
}

func MustReference(tree yggdrasil.Tree) graphql.Schema {
	schema, err := Reference(tree)
	if err != nil {
		panic(err)
	}
	return schema
}
