// +build !firestore

package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/rjansen/fivecolors/core/graphql"
	"github.com/rjansen/l"
	"github.com/rjansen/raizel/sql"
	"github.com/rjansen/yggdrasil"
	"github.com/stretchr/testify/require"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type testSchema struct {
	name     string
	tree     yggdrasil.Tree
	dbMock   sqlmock.Sqlmock
	data     interface{}
	columns  []string
	rows     [][]interface{}
	schema   graphql.Schema
	request  graphql.Request
	response *graphql.Response
}

func (scenario *testSchema) setup(t *testing.T) {
	var (
		roots                   = yggdrasil.NewRoots()
		errLogger               = l.Register(&roots, l.NewZapLoggerDefault())
		sqlDB, dbMock, errSqlDB = sqlmock.New()
		db, errDB               = sql.NewDB(sqlDB)
		errRegisterDB           = sql.Register(&roots, db)
	)
	require.Nil(t, errLogger, "setup logger error")
	require.Nil(t, errSqlDB, "new sqldb error")
	require.Nil(t, errDB, "new db error")
	require.Nil(t, errRegisterDB, "setup db error")

	tree := roots.NewTreeDefault()
	schema := NewSchema(tree)
	require.NotNil(t, schema, "new schema error")

	mockRows := sqlmock.NewRows(scenario.columns)
	for _, row := range scenario.rows {
		columns := make([]driver.Value, len(row))
		for index, column := range row {
			columns[index] = column
		}
		mockRows.AddRow(columns...)
	}
	dbMock.ExpectQuery("select").WillReturnRows(mockRows)
	// mock.ExpectQuery(scenario.query).WillReturnError(scenario.err)
	// dbMock.ExpectClose()

	scenario.schema = schema
	scenario.dbMock = dbMock
	scenario.tree = tree
}

func (scenario *testSchema) tearDown(*testing.T) {
	if scenario.tree != nil {
		scenario.tree.Close()
	}
}

func TestSchema(test *testing.T) {
	scenarios := []testSchema{
		{
			name: "Resolves card field successfully",
			data: &struct {
				Card Card `json:"card"`
			}{},
			columns: []string{
				"id", "name", "number_cost", "id_external", "id_rarity", "id_set", "id_asset",
				"rate", "rate_votes", "order_external", "artist", "flavor", "data",
				"types", "costs", "rules", "created_at", "updated_at", "deleted_at",
			},
			rows: [][]interface{}{
				{

					"mock_id", "mock_name", 2.0, "mock_id_external", "mock_id_rarity", "mock_id_set",
					"mock_id_asset", 0.0, 0, "1022A", "Mock Artist", "Mock Flavor",
					Object{"key1": "value1", "key2": 10}, `{"Legendary","Mock"}`, `{"1","B","U"}`,
					`{"rule one bla bla","more bla"}`, time.Now().UTC(), time.Now().UTC(), time.Now().UTC(),
				},
			},
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
			columns: []string{
				"id", "name", "alias", "created_at", "updated_at", "deleted_at",
			},
			rows: [][]interface{}{
				{

					"mock_id", "mock_name", "mock_alias",
					time.Now().UTC(), time.Now().UTC(), time.Now().UTC(),
				},
			},
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
			columns: []string{
				"id", "name", "number_cost", "id_external", "id_rarity", "id_set", "id_asset",
				"rate", "rate_votes", "order_external", "artist", "flavor", "data",
				"types", "costs", "rules", "created_at", "updated_at", "deleted_at",
			},
			rows: [][]interface{}{
				{

					"mock_id", "mock_name", 2.0, "mock_id_external", "mock_id_rarity", "mock_id_set",
					"mock_id_asset", 0.0, 0, "1022A", "Mock Artist", "Mock Flavor",
					Object{"key1": "value1", "key2": 10}, `{"Legendary","Mock"}`, `{"1","B","U"}`,
					`{"rule one bla bla","more bla"}`, time.Now().UTC(), time.Now().UTC(), time.Now().UTC(),
				},
				{

					"mock_id", "mock_name", 2.0, "mock_id_external", "mock_id_rarity", "mock_id_set",
					"mock_id_asset", 0.0, 0, "1022A", "Mock Artist", "Mock Flavor",
					Object{"key1": "value1", "key2": 10}, `{"Legendary","Mock"}`, `{"1","B","U"}`,
					`{"rule one bla bla","more bla"}`, time.Now().UTC(), time.Now().UTC(), time.Now().UTC(),
				},
				{

					"mock_id", "mock_name", 2.0, "mock_id_external", "mock_id_rarity", "mock_id_set",
					"mock_id_asset", 0.0, 0, "1022A", "Mock Artist", "Mock Flavor",
					Object{"key1": "value1", "key2": 10}, `{"Legendary","Mock"}`, `{"1","B","U"}`,
					`{"rule one bla bla","more bla"}`, time.Now().UTC(), time.Now().UTC(), time.Now().UTC(),
				},
				{

					"mock_id", "mock_name", 2.0, "mock_id_external", "mock_id_rarity", "mock_id_set",
					"mock_id_asset", 0.0, 0, "1022A", "Mock Artist", "Mock Flavor",
					Object{"key1": "value1", "key2": 10}, `{"Legendary","Mock"}`, `{"1","B","U"}`,
					`{"rule one bla bla","more bla"}`, time.Now().UTC(), time.Now().UTC(), time.Now().UTC(),
				},
			},
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
			columns: []string{
				"id", "name", "alias", "created_at", "updated_at", "deleted_at",
			},
			rows: [][]interface{}{
				{

					"mock_id", "mock_name", "mock_alias",
					time.Now().UTC(), time.Now().UTC(), time.Now().UTC(),
				},
				{

					"mock_id", "mock_name", "mock_alias",
					time.Now().UTC(), time.Now().UTC(), time.Now().UTC(),
				},
			},
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
				require.Nilf(t, response.Errors, "schema response errors: %+v", response.Errors)
				t.Logf("json data=%q", response.Data)
				err := json.Unmarshal(response.Data, scenario.data)
				require.Nil(t, err, "schema response unmarshal error")
				require.NotZerof(t, scenario.data, "data response invalid: %+v", scenario.data)
				require.Nil(t, scenario.dbMock.ExpectationsWereMet(), "dbmock invalid expectations")
			},
		)
	}
}
