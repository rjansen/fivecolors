package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type httpMock struct {
	status int
	body   string
}

func newHTTPServerMock(status int, body string) *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(status)
				w.Write([]byte(body))
			},
		),
	)
}

func TestNewHTTPRequestError(t *testing.T) {
	type (
		expected struct {
			err string
		}

		test struct {
			name        string
			method, url string
			request     Request
			expected    expected
			setup       func(t *testing.T)
			tearDown    func(t *testing.T)
		}
	)

	tests := []test{
		{
			name:    "when marshal returns error",
			method:  http.MethodGet,
			url:     "http://localhost/ping",
			request: Request{},
			expected: expected{
				err: "err_marshall_request: err_mock_marshal",
			},
			setup: func(t *testing.T) {
				marshal = func(interface{}) ([]byte, error) {
					return nil, errors.New("err_mock_marshal")
				}
			},
			tearDown: func(t *testing.T) {
				marshal = json.Marshal
			},
		},
		{
			name:    "when new request returns error",
			method:  http.MethodGet,
			url:     "http://localhost/ping",
			request: Request{},
			expected: expected{
				err: "err_mock_newrequest",
			},
			setup: func(t *testing.T) {
				newRequest = func(string, string, io.Reader) (*http.Request, error) {
					return nil, errors.New("err_mock_newrequest")
				}
			},
			tearDown: func(t *testing.T) {
				newRequest = http.NewRequest
			},
		},
	}

	for index, test := range tests {
		t.Run(
			fmt.Sprintf("%d-%s", index, test.name),
			func(t *testing.T) {
				test.setup(t)
				defer test.tearDown(t)

				httpRequest, err := newHTTPRequest(test.method, test.url, test.request)
				assert.EqualError(t, err, test.expected.err)
				assert.Nil(t, httpRequest)
			},
		)
	}
}

func TestHttpClient(t *testing.T) {
	type (
		expected struct {
			response string
			err      string
		}

		test struct {
			name     string
			httpMock httpMock
			ctx      context.Context
			request  Request
			expected expected
		}
	)

	tests := []test{
		{
			name:     "when server returns some graphql errors",
			httpMock: httpMock{status: http.StatusOK, body: `{"errors": [{"message": "err_mock"}]}`},
			ctx:      context.Background(),
			request:  Request{Query: `query {data(id: "xpto"){id,name,alias}}`},
			expected: expected{
				response: `{"errors":[{"message": "err_mock"}]}`,
			},
		},
		{
			name:     "when request query is invalid returns bad request",
			httpMock: httpMock{status: http.StatusBadRequest, body: `{"errors": [{"message": "err_mock_bad_request"}]}`},
			ctx:      context.Background(),
			request:  Request{Query: `invalid`},
			expected: expected{
				err: "err_graphql: errors.List{err_graphql_status_code: 400, err_mock_bad_request}",
			},
		},
		{
			name:     "when server does not returns status OK",
			httpMock: httpMock{status: http.StatusInternalServerError, body: `{"errors": [{"message": "err_mock_1"}, {"message": "err_mock_2"}, {"message": "err_mock_3"}]}`},
			ctx:      context.Background(),
			request:  Request{Query: `invalid`},
			expected: expected{
				err: "err_graphql: errors.List{err_graphql_status_code: 500, err_mock_1, err_mock_2, err_mock_3}",
			},
		},
		{
			name:     "when server returns a non string error message",
			httpMock: httpMock{status: http.StatusInternalServerError, body: `{"errors": [{"message": 555}]}`},
			ctx:      context.Background(),
			request:  Request{Query: `query {data(id: "xpto"){id,name,alias}}`},
			expected: expected{
				err: "err_graphql: errors.List{err_graphql_status_code: 500, 555}",
			},
		},
		{
			name:     "when server returns a invalid json body",
			httpMock: httpMock{status: http.StatusOK, body: `invalid`},
			ctx:      context.Background(),
			request:  Request{Query: `query {data(id: "xpto"){id,name,alias}}`},
			expected: expected{
				err: "err_unmarshall_response: invalid character 'i' looking for beginning of value",
			},
		},
	}

	for index, test := range tests {
		t.Run(
			fmt.Sprintf("%d-%s", index, test.name),
			func(t *testing.T) {
				httpServer := newHTTPServerMock(test.httpMock.status, test.httpMock.body)
				defer httpServer.Close()

				graphqlClient := newGraphQLClient(method, httpServer.URL, NewHTTPClient())
				var response json.RawMessage
				err := graphqlClient.Do(test.ctx, test.request, &response)
				if err != nil {
					assert.EqualError(t, err, test.expected.err)
				}
				if response != nil {
					assert.JSONEq(t, test.expected.response, string(response))
				}
			},
		)
	}
}
