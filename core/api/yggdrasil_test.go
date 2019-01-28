package api

import (
	"fmt"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/rjansen/yggdrasil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustGetSchema(schema graphql.Schema, err error) graphql.Schema {
	if err != nil {
		panic(err)
	}
	return schema
}

type testRegister struct {
	name   string
	schema graphql.Schema
	err    error
}

func TestRegister(test *testing.T) {
	scenarios := []testRegister{
		{
			name:   "Register the Repository reference",
			schema: mustGetSchema(newMockSchema()),
		},
	}

	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				roots := yggdrasil.NewRoots()
				err := Register(&roots, scenario.schema)
				assert.Equal(t, scenario.err, err)

				tree := roots.NewTreeDefault()
				schema, err := tree.Reference(schemaPath)

				require.Nil(t, err, "tree reference error")
				require.Exactly(t, scenario.schema, schema, "schema reference")
			},
		)
	}
}

type testReference struct {
	name       string
	references map[yggdrasil.Path]yggdrasil.Reference
	tree       yggdrasil.Tree
	err        error
}

func (scenario *testReference) setup(t *testing.T) {
	roots := yggdrasil.NewRoots()
	for path, reference := range scenario.references {
		err := roots.Register(path, reference)
		assert.Nil(t, err, "register error")
	}
	scenario.tree = roots.NewTreeDefault()
}

func TestReference(test *testing.T) {
	scenarios := []testReference{
		{
			name: "Access the Repository Reference",
			references: map[yggdrasil.Path]yggdrasil.Reference{
				schemaPath: yggdrasil.NewReference(mustGetSchema(newMockSchema())),
			},
		},
		{
			name: "When Repository was not register returns path not found",
			err:  yggdrasil.ErrPathNotFound,
		},
		{
			name: "When a invalid Repository was register returns invalid reference error",
			references: map[yggdrasil.Path]yggdrasil.Reference{
				schemaPath: yggdrasil.NewReference(new(struct{})),
			},
			err: ErrInvalidReference,
		},
	}

	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)

				_, err := Reference(scenario.tree)
				assert.Equal(t, scenario.err, err, "reference error")
				if scenario.err != nil {
					assert.PanicsWithValue(t, scenario.err,
						func() {
							_ = MustReference(scenario.tree)
						},
					)
				} else {
					assert.NotPanics(t,
						func() {
							_ = MustReference(scenario.tree)
						},
					)
				}
			},
		)
	}
}
