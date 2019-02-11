package graphql

import (
	"errors"

	"github.com/rjansen/yggdrasil"
)

var (
	ErrInvalidReference = errors.New("Invalid Schema Reference")
	schemaPath          = yggdrasil.NewPath("/fivecolors/core/graphql/schema")
)

func Register(roots *yggdrasil.Roots, schema Schema) error {
	return roots.Register(schemaPath, schema)
}

func Reference(tree yggdrasil.Tree) (Schema, error) {
	var (
		schema         Schema
		reference, err = tree.Reference(schemaPath)
	)
	if err != nil {
		return schema, err
	}
	tmpSchema, is := reference.(Schema)
	if !is {
		return schema, ErrInvalidReference
	}
	schema = tmpSchema
	return schema, nil
}

func MustReference(tree yggdrasil.Tree) Schema {
	schema, err := Reference(tree)
	if err != nil {
		panic(err)
	}
	return schema
}
