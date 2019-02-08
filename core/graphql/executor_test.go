package graphql

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/99designs/gqlgen/graphql"
	mockschema "github.com/rjansen/fivecolors/core/graphql/mockschema"
	"github.com/rjansen/l"
	"github.com/rjansen/yggdrasil"
	"github.com/stretchr/testify/require"
)

type testExecute struct {
	name   string
	tree   yggdrasil.Tree
	data   interface{}
	schema graphql.ExecutableSchema
	params Params
	result *graphql.Response
}

func (scenario *testExecute) setup(t *testing.T) {
	var (
		roots     = yggdrasil.NewRoots()
		errLogger = l.Register(&roots, l.NewZapLoggerDefault())
		schema    = mockschema.New()
	)
	require.Nil(t, errLogger, "setup logger error")
	scenario.tree = roots.NewTreeDefault()
	scenario.schema = schema
}

func (scenario *testExecute) tearDown(*testing.T) {}

func TestExecute(test *testing.T) {
	scenarios := []testExecute{
		{
			name: "Resolves me field successfully",
			data: &struct {
				Me mockschema.MeResponse `json:"me"`
			}{},
			params: Params{
				Query: `{
					me {
					  tid
					  user {
						id
						name
					  }
                    }
				}`,
			},
		},
		{
			name: "Resolves mockEntity field successfully",
			data: &struct {
				MockEntityResponse mockschema.MockEntityResponse `json:"mockEntity"`
			}{},
			params: Params{
				Query: `{
					mockEntity {
					  tid
					  entity {
						id
						string
						integer
						float
						boolean
						dateTime
						object
					  }
					}
				}`,
			},
		},
	}

	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)
				defer scenario.tearDown(t)

				response := Execute(scenario.tree, scenario.schema, scenario.params)
				require.NotNil(t, response, "schema response invalid")
				require.Nil(t, response.Errors, "schema response errors")
				t.Logf("json data=%q", response.Data)
				err := json.Unmarshal(response.Data, scenario.data)
				require.Nil(t, err, "schema response unmarshal error")
				require.NotZero(t, scenario.data, "data response invalid")
			},
		)
	}
}
