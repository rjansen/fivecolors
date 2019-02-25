// +build integration

package function

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/rjansen/fivecolors/core/api"
	"github.com/rjansen/fivecolors/core/model"
	"github.com/rjansen/raizel/sql"
	"github.com/rjansen/yggdrasil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type dataGenerator struct {
	sql       string
	arguments []interface{}
}

type testPostgresIntegrationHandler struct {
	name           string
	tree           yggdrasil.Tree
	mockSetup      []dataGenerator
	mockTearDown   []dataGenerator
	body           string
	contentType    string
	request        *http.Request
	response       *httptest.ResponseRecorder
	responseStatus int
}

func (scenario *testPostgresIntegrationHandler) setup(t *testing.T) {
	var (
		options = options{
			projectID: "project-id",
			dataStore: "postgres",
			driver:    "postgres",
			dsn:       "postgres://postgres:@127.0.0.1:5432/postgres?sslmode=disable",
		}
		tree = newTree(options)
		db   = sql.MustReference(tree)
	)
	require.NotNil(t, tree, "tree invalid instance")
	require.NotNil(t, db, "db invalid instance")

	for index, dataGen := range scenario.mockSetup {
		_, err := db.Exec(dataGen.sql, dataGen.arguments...)
		assert.Nilf(t, err, "error setup mock data: index=%d err=%+v", index, err)
	}

	r := httptest.NewRequest(
		"POST", "/", strings.NewReader(scenario.body),
	)
	r.Header.Set("content-type", scenario.contentType)

	serverHandler = api.NewGraphQLHandler(tree)
	scenario.tree = tree
	scenario.request = r
	scenario.response = httptest.NewRecorder()
}

func (scenario *testPostgresIntegrationHandler) tearDown(t *testing.T) {
	if scenario.tree != nil {
		defer scenario.tree.Close()
		var (
			db = sql.MustReference(scenario.tree)
		)
		for index, dataGen := range scenario.mockTearDown {
			_, err := db.Exec(dataGen.sql, dataGen.arguments...)
			assert.Nilf(t, err, "error terdown mock data: index=%d err=%+v", index, err)
		}
	}
}

func TestPostgresIntegrationHandler(test *testing.T) {
	timeNow := time.Now().UTC()
	scenarios := []testPostgresIntegrationHandler{
		{
			name:           "When request body is invalid returns a bad request",
			body:           "<xml><id>xmlid</id></name>Invalid Body</name></xml>",
			contentType:    "application/xml",
			responseStatus: http.StatusBadRequest,
		},
		{
			name: "When request body is graphql executes the query and returns ok with query results",
			mockSetup: []dataGenerator{
				{
					sql:       "insert into rarity (id, name, alias, created_at) values ($1, $2, $3, $4)",
					arguments: []interface{}{"mock_rarityid", model.RarityNameMythicRare, model.RarityAliasM, timeNow},
				},
				{
					sql:       "insert into set (id, name, alias, created_at) values ($1, $2, $3, $4)",
					arguments: []interface{}{"mock_setid", "Set Mock", "stm", timeNow},
				},
				{
					sql: `insert into card (
							id, id_external, id_asset, name, types, costs, number_cost, order_external, created_at,
							id_rarity, id_set
						  ) values(
							  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
						  )`,
					arguments: []interface{}{
						"mock_cardid",
						"mock_cardidexternal",
						"mock_cardassetid",
						"Card Mock",
						pq.Array([]string{"Legendary", "Creature", "Goblin"}),
						pq.Array([]string{"1", "R", "R", "R"}),
						4.0,
						"1A",
						timeNow,
						"mock_rarityid",
						"mock_setid",
					},
				},
			},
			mockTearDown: []dataGenerator{
				{sql: "delete from card where id = $1", arguments: []interface{}{"mock_cardid"}},
				{sql: "delete from rarity where id = $1", arguments: []interface{}{"mock_rarityid"}},
				{sql: "delete from set where id = $1", arguments: []interface{}{"mock_setid"}},
			},
			body: `{
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
			contentType:    "application/graphql",
			responseStatus: http.StatusOK,
		},
		{
			name: "When request body is a valid json executes the query and returns ok with query results",
			mockSetup: []dataGenerator{
				{
					sql:       "insert into rarity (id, name, alias, created_at) values ($1, $2, $3, $4)",
					arguments: []interface{}{"mock_rarityid", model.RarityNameMythicRare, model.RarityAliasM, timeNow},
				},
				{
					sql:       "insert into set (id, name, alias, created_at) values ($1, $2, $3, $4)",
					arguments: []interface{}{"mock_setid", "Set Mock", "stm", timeNow},
				},
				{
					sql: `insert into card (
							id, id_external, id_asset, name, types, costs, number_cost, order_external, created_at,
							id_rarity, id_set
						  ) values(
							  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
						  )`,
					arguments: []interface{}{
						"mock_cardid",
						"mock_cardidexternal",
						"mock_cardassetid",
						"Card Mock",
						pq.Array([]string{"Legendary", "Creature", "Goblin"}),
						pq.Array([]string{"1", "R", "R", "R"}),
						4.0,
						"1A",
						timeNow,
						"mock_rarityid",
						"mock_setid",
					},
				},
			},
			mockTearDown: []dataGenerator{
				{sql: "delete from card where id = $1", arguments: []interface{}{"mock_cardid"}},
				{sql: "delete from rarity where id = $1", arguments: []interface{}{"mock_rarityid"}},
				{sql: "delete from set where id = $1", arguments: []interface{}{"mock_setid"}},
			},
			body: `{
				"query": "{card(id: \"mock_cardid\") {id,name,types,costs,numberCost,idAsset,data,createdAt,updatedAt,deletedAt}}"
			}`,
			contentType:    "application/json",
			responseStatus: http.StatusOK,
		},
	}

	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)
				defer scenario.tearDown(t)

				Handler(scenario.response, scenario.request)
				result := scenario.response.Result()
				require.Equal(t, scenario.responseStatus, result.StatusCode, "invalid response statuscode")
			},
		)
	}
}
