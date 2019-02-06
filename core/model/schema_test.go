package model

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"github.com/rjansen/l"
	"github.com/rjansen/yggdrasil"
	"github.com/stretchr/testify/require"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type testSchema struct {
	name   string
	tree   yggdrasil.Tree
	dbMock sqlmock.Sqlmock
	data   interface{}
	params Params
	result *graphql.Response
}

func (scenario *testSchema) setup(t *testing.T) {
	var (
		roots = yggdrasil.NewRoots()
		// db, dbMock, errSqlMock = sqlmock.New()
		// errDB                  = sql.Register(&roots, db)
		errLogger = l.Register(&roots, l.NewZapLoggerDefault())
	)
	require.Nil(t, errLogger, "setup logger error")
	// require.Nil(t, errDB, "setup db error")

	tree := roots.NewTreeDefault()
	// schema, err := NewSchema(tree)
	// require.Nil(t, err, "setup schema error")

	// mockRows := sqlmock.NewRows(scenario.columns)
	// for _, row := range scenario.rows {
	// 	columns := make([]driver.Value, len(row))
	// 	for index, column := range row {
	// 		columns[index] = column
	// 	}
	// 	mockRows.AddRow(columns...)
	// }
	// mock.ExpectQuery(scenario.query).WillReturnRows(mockRows)
	// // mock.ExpectQuery(scenario.query).WillReturnError(scenario.err)
	// mock.ExpectClose()

	scenario.tree = tree
	// scenario.dbMock = dbMock
}

func (scenario *testSchema) tearDown(*testing.T) {}

func TestSchema(test *testing.T) {
	scenarios := []testSchema{
		{
			name: "Resolves card field successfully",
			data: &struct {
				Card Card `json:"card"`
			}{},
			params: Params{
				Query: `{
					card(id: "mock_cardid") {
					  id
					  name
					  types
					  costs
					  numberCost
				      idAsset
					  data
					  createdAt
					  updatedAt
					  deletedAt
                      rarity {
                        id
                        name
                        alias
                      }
                      set {
                        id
                        name
                        alias
                      }
                    }
				}`,
			},
		},
		{
			name: "Resolves set field successfully",
			data: &struct {
				Set Set `json:"set"`
			}{},
			params: Params{
				Query: `{
					set(id: "mock_setid") {
					  id
					  name
					  alias
					  cards {
					    id
					    name
					    types
					    costs
					    numberCost
					    idAsset
					    data
					    createdAt
					    updatedAt
					    deletedAt
					    rarity {
					      id
					      name
					      alias
						}
					  }
					}
				}`,
			},
		},
		{
			name: "Resolves cardBy field successfully",
			data: &struct {
				Card []Card `json:"card"`
			}{},
			params: Params{
				Query: `{
					cardBy(filter: {
					  id: "mock_cardid"
					  name: "mock_cardname"
					  types: ["mock_cardtype1", "mock_cardtype2"]
					  costs: ["1", "B", "U", "R", "W"]
					  set: {
						id: "mock_setid"
						name: "mock_setname"
						alias: "mock_setalias"
					  }
					  rarity: {
						id: Rare
						alias: R
					  }
					}) {
					  id
					  name
					  types
					  costs
					  numberCost
					  idAsset
					  data
					  createdAt
					  updatedAt
					  deletedAt
					  rarity {
					    id
					    name
					    alias
					  }
					  set {
					    id
					    name
					    alias
					  }
					}
				}`,
			},
		},
		{
			name: "Resolves setBy field successfully",
			data: &struct {
				Set []Set `json:"set"`
			}{},
			params: Params{
				Query: `{
					setBy(filter: {
					  id: "mock_setid"
					  name: "mock_setname"
					  alias: "mock_setalias"
					}) {
					  id
					  name
					  alias
					  cards {
					    id
				  	    name
					    types
					    costs
					    numberCost
					    idAsset
					    data
					    createdAt
					    updatedAt
					    deletedAt
					    rarity {
					      id
					      name
					      alias
					    }
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

				response := Execute(scenario.tree, scenario.params)
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
