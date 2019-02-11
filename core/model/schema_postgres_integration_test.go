// +build integration

package model

import (
	stdsql "database/sql"
	"encoding/json"
	"fmt"
	"testing"

	_ "github.com/lib/pq"
	"github.com/rjansen/fivecolors/core/graphql"
	"github.com/rjansen/l"
	"github.com/rjansen/raizel"
	"github.com/rjansen/raizel/sql"
	"github.com/rjansen/yggdrasil"
	"github.com/stretchr/testify/require"
)

type testSchemaPostgres struct {
	name     string
	tree     yggdrasil.Tree
	schema   graphql.Schema
	data     interface{}
	request  graphql.Request
	response *graphql.Response
}

func (scenario *testSchemaPostgres) setup(t *testing.T) {
	var (
		roots           = yggdrasil.NewRoots()
		errLogger       = l.Register(&roots, l.NewZapLoggerDefault())
		driver          = "postgres"
		dsn             = "postgres://postgres:@127.0.0.1:5432/postgres?sslmode=disable"
		sqlDB, errSqlDB = stdsql.Open(driver, dsn)
		db, errDB       = sql.NewDB(sqlDB)
		errRegisterDB   = sql.Register(&roots, db)
	)
	require.Nil(t, errLogger, "setup logger error")
	require.Nil(t, errSqlDB, "new sqldb error")
	require.Nil(t, errDB, "new db error")
	require.Nil(t, errRegisterDB, "setup db error")

	tree := roots.NewTreeDefault()
	schema := NewSchema(tree)
	require.NotNil(t, schema, "new schema error")

	scenario.tree = tree
	scenario.schema = schema
}

func (scenario *testSchemaPostgres) tearDown(*testing.T) {
	if scenario.tree != nil {
		scenario.tree.Close()
	}
}

func TestSchemaPostgres(test *testing.T) {
	scenarios := []testSchemaPostgres{
		{
			name: "Resolves card field successfully",
			data: &struct {
				Card Card `json:"card"`
			}{},
			request: graphql.Request{
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
                    }
				}`,
			},
		},
		{
			name: "Resolves set field successfully",
			data: &struct {
				Set Set `json:"set"`
			}{},
			request: graphql.Request{
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
			request: graphql.Request{
				Query: `{
					cardBy(filter: {
					  name: "mock_cardname"
					  types: "mock_cardtype1 - mock_cardtype2"
					  costs: "1BURW"
					  set: {
						name: "mock_setname"
						alias: "mock_setalias"
					  }
					  rarity: {
						name: Rare
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
			request: graphql.Request{
				Query: `{
					setBy(filter: {
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

				response := graphql.Execute(scenario.tree, scenario.schema, scenario.request)
				require.NotNil(t, response, "schema response invalid")
				// require.Nilf(t, response.Errors, "schema response errors: %+v", response.Errors)
				if len(response.Errors) > 0 {
					require.Lenf(t, response.Errors, 1, "schema response errors: %+v", response.Errors)
					require.Equal(t, raizel.ErrNotFound.Error(), response.Errors[0].Message,
						"schema response unmarshal error")
				}
				err := json.Unmarshal(response.Data, scenario.data)
				require.Nil(t, err, "schema response unmarshal error")
				// require.NotZero(t, scenario.data, "data response invalid")
			},
		)
	}
}
