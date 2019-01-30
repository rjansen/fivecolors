package model

import (
	"database/sql/driver"
	"fmt"
	"net/http"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/rjansen/l"
	"github.com/rjansen/raizel/sql"
	"github.com/rjansen/yggdrasil"
	"github.com/stretchr/testify/require"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type testSchema struct {
	name   string
	tree   yggdrasil.Tree
	dbMock sqlmock.Sqlmock
	schema graphql.Schema
	query  string
	result *graphql.Result
}

func (scenario *testSchema) setup(t *testing.T) {
	var (
		roots                  = yggdrasil.NewRoots()
		db, dbMock, errSqlMock = sqlmock.New()
		errDB                  = sql.Register(&roots, db)
		errLogger              = l.Register(&roots, l.NewZapLoggerDefault())
	)
	require.Nil(t, errLogger, "setup logger error")
	require.Nil(t, errDB, "setup db error")

	tree := roots.NewTreeDefault()
	schema, err := NewSchema(tree)
	require.Nil(t, err, "setup schema error")

	mockRows := sqlmock.NewRows(scenario.columns)
	for _, row := range scenario.rows {
		columns := make([]driver.Value, len(row))
		for index, column := range row {
			columns[index] = column
		}
		mockRows.AddRow(columns...)
	}
	mock.ExpectQuery(scenario.query).WillReturnRows(mockRows)
	// mock.ExpectQuery(scenario.query).WillReturnError(scenario.err)
	mock.ExpectClose()

	scenario.tree = tree
	scenario.dbMock = dbMock
}

func (scenario *testSchema) tearDown(*testing.T) {}

func TestSchema(test *testing.T) {
	scenarios := []testSchema{
		{
			name:    "When handler receives a HEAD request returns method not allowed",
			columns: []string{"id", "text"},
			rows: [][]interface{}{
				{1, "mock one"},
				{2, "mock two"},
			},
		},
		{
			name:           "When a GET request has a blank query parameter returns bad request",
			method:         http.MethodGet,
			path:           "/query?q=",
			responseStatus: http.StatusBadRequest,
		},
		{
			name:           "When a POST content type is invalid returns bad request",
			method:         http.MethodPost,
			path:           "/query",
			body:           "<xml><id>xmlid</id></name>Invalid Body</name></xml>",
			contentType:    "application/xml",
			responseStatus: http.StatusBadRequest,
		},
		{
			name:           "When a POST request body has a invalid graphql content returns internal server error",
			method:         http.MethodPost,
			path:           "/query",
			body:           "<xml><id>xmlid</id></name>Invalid Body</name></xml>",
			contentType:    "application/graphql",
			responseStatus: http.StatusInternalServerError,
		},
		{
			name:           "When a POST request body has a invalid json content returns bad request",
			method:         http.MethodPost,
			path:           "/query",
			body:           "<xml><id>xmlid</id></name>Invalid Body</name></xml>",
			contentType:    "application/json",
			responseStatus: http.StatusBadRequest,
		},
		{
			name:           "When a POST request body has a invalid graphql query returns internal server error",
			method:         http.MethodPost,
			path:           "/query",
			body:           `{"query": "<xml><id>xmlid</id></name>Invalid Body</name></xml>"}`,
			contentType:    "application/json",
			responseStatus: http.StatusInternalServerError,
		},
		{
			name:           "When a POST request body is graphql executes the query and returns ok with query results",
			method:         http.MethodPost,
			path:           "/query",
			body:           "{me{tid,user{id,name}}}",
			contentType:    "application/graphql",
			responseStatus: http.StatusOK,
		},
		{
			name:           "When a POST request body is graphql executes the query and returns ok with query results",
			method:         http.MethodPost,
			path:           "/query",
			body:           "{mockEntity{tid,entity{id,string,float,integer,date_time}}}",
			contentType:    "application/graphql",
			responseStatus: http.StatusOK,
		},
		{
			name:           "When a POST request body is a valid json executes the query and returns ok with query results",
			method:         http.MethodPost,
			path:           "/query",
			body:           `{"query": "{me{tid,user{id,name}}}"}`,
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

				graphql := NewGraphQLHandler(scenario.tree)
				graphql(scenario.response, scenario.request)
				result := scenario.response.Result()
				/*
					var (
						resultMap map[string]interface{}
						errDecode = json.NewDecoder(result.Body).Decode(&resultMap)
					)
					require.Equal(t, map[string]interface{}{}, resultMap, "result invalid instance")
					require.Nil(t, errDecode, "response body invalid")
					require.NotZero(t, resultMap, "result invalid")
				*/
				require.Equal(t, scenario.responseStatus, result.StatusCode, "invalid response statuscode")
			},
		)
	}
}
