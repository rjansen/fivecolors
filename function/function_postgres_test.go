package function

import (
	"database/sql/driver"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/rjansen/fivecolors/core/api"
	"github.com/rjansen/fivecolors/core/graphql"
	"github.com/rjansen/fivecolors/core/model"
	"github.com/rjansen/l"
	"github.com/rjansen/raizel/sql"
	"github.com/rjansen/yggdrasil"
	"github.com/stretchr/testify/require"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type testPostgresHandler struct {
	name           string
	tree           yggdrasil.Tree
	dbMock         sqlmock.Sqlmock
	data           interface{}
	columns        []string
	rows           [][]interface{}
	body           string
	contentType    string
	request        *http.Request
	response       *httptest.ResponseRecorder
	responseStatus int
}

func (scenario *testPostgresHandler) setup(t *testing.T) {
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

	errSchema := graphql.Register(&roots,
		model.NewSchema(
			model.NewResolver(
				model.NewPostgresQueryResolver(
					roots.NewTreeDefault(),
				),
			),
		),
	)
	require.Nil(t, errSchema, "new schema error")
	tree := roots.NewTreeDefault()

	if len(scenario.columns) > 0 {
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
	}
	r := httptest.NewRequest(
		"POST", "/", strings.NewReader(scenario.body),
	)
	r.Header.Set("content-type", scenario.contentType)

	serverHandler = api.NewGraphQLHandler(tree)
	scenario.dbMock = dbMock
	scenario.tree = tree
	scenario.request = r
	scenario.response = httptest.NewRecorder()
}

func (scenario *testPostgresHandler) tearDown(*testing.T) {
	if scenario.tree != nil {
		scenario.tree.Close()
	}
}

func TestPostgresHandler(test *testing.T) {
	scenarios := []testPostgresHandler{
		{
			name:           "When request body is invalid returns a bad request",
			body:           "<xml><id>xmlid</id></name>Invalid Body</name></xml>",
			contentType:    "application/xml",
			responseStatus: http.StatusBadRequest,
		},
		{
			name: "When request body is graphql executes the query and returns ok with query results",
			columns: []string{
				"id", "name", "number_cost", "id_external", "id_rarity", "id_set", "id_asset",
				"rate", "rate_votes", "order_external", "artist", "flavor", "data",
				"types", "costs", "rules", "created_at", "updated_at", "deleted_at",
			},
			rows: [][]interface{}{
				{

					"mock_id", "mock_name", 2.0, "mock_id_external", "mock_id_rarity", "mock_id_set",
					"mock_id_asset", 0.0, 0, "1022A", "Mock Artist", "Mock Flavor",
					model.Object{"key1": "value1", "key2": 10}, `{"Legendary","Mock"}`, `{"1","B","U"}`,
					`{"rule one bla bla","more bla"}`, time.Now().UTC(), time.Now().UTC(), time.Now().UTC(),
				},
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
			columns: []string{
				"id", "name", "number_cost", "id_external", "id_rarity", "id_set", "id_asset",
				"rate", "rate_votes", "order_external", "artist", "flavor", "data",
				"types", "costs", "rules", "created_at", "updated_at", "deleted_at",
			},
			rows: [][]interface{}{
				{

					"mock_id", "mock_name", 2.0, "mock_id_external", "mock_id_rarity", "mock_id_set",
					"mock_id_asset", 0.0, 0, "1022A", "Mock Artist", "Mock Flavor",
					model.Object{"key1": "value1", "key2": 10}, `{"Legendary","Mock"}`, `{"1","B","U"}`,
					`{"rule one bla bla","more bla"}`, time.Now().UTC(), time.Now().UTC(), time.Now().UTC(),
				},
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
				require.Nil(t, scenario.dbMock.ExpectationsWereMet(), "dbmock invalid expectations")
			},
		)
	}
}
