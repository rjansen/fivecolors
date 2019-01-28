package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type testHandler struct {
	name           string
	body           string
	contentType    string
	request        *http.Request
	response       *httptest.ResponseRecorder
	responseStatus int
}

func (scenario *testHandler) setup(*testing.T) {
	r := httptest.NewRequest(
		"POST", "/query", strings.NewReader(scenario.body),
	)
	r.Header.Set("content-type", scenario.contentType)

	scenario.request = r
	scenario.response = httptest.NewRecorder()
}

func (scenario *testHandler) tearDown(*testing.T) {}

func TestHandler(test *testing.T) {
	scenarios := []testHandler{
		{
			name:           "When request body is invalid returns a bad request",
			body:           "<xml><id>xmlid</id></name>Invalid Body</name></xml>",
			contentType:    "application/xml",
			responseStatus: http.StatusBadRequest,
		},
		{
			name:           "When request body is graphql executes the query and returns ok with query results",
			body:           "{me{tid,user{id,name}}}",
			contentType:    "application/graphql",
			responseStatus: http.StatusOK,
		},
		{
			name:           "When request body is a valid json executes the query and returns ok with query results",
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

				Handler(scenario.response, scenario.request)
				result := scenario.response.Result()
				require.Equal(t, scenario.responseStatus, result.StatusCode, "invalid response statuscode")
			},
		)
	}
}
