package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rjansen/abend"
	"github.com/rjansen/l"
)

const (
	defaultMethod = http.MethodPost
	defaultURL    = "http://localhost:8080/query"
)

var (
	method        = defaultMethod
	url           = defaultURL
	newRequest    = http.NewRequest
	marshal       = json.Marshal
	defaultDecode = func(r io.Reader, d interface{}) error {
		return json.NewDecoder(r).Decode(d)
	}
	decode = defaultDecode
)

type (
	Object = map[string]interface{}

	Request struct {
		Query     string `json:"query"`
		Variables Object `json:"variables"`
	}

	Response struct {
		Errors []Object `json:"errors"`
	}
)

func NewHTTPClient() *http.Client {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: false,
	}

	return &http.Client{Transport: tr}
}

func newHTTPRequest(method, url string, request Request) (*http.Request, error) {
	requestBytes, err := marshal(request)
	if err != nil {
		return nil, fmt.Errorf("err_marshall_request: %w", err)
	}

	httpRequest, err := newRequest(method, url, bytes.NewBuffer(requestBytes))
	if err != nil {
		return nil, err
	}
	httpRequest.Header.Set("Content-Type", "application/json")
	return httpRequest, nil
}

type graphqlClient struct {
	method, url string
	httpClient  *http.Client
}

func newGraphQLClient(method, url string, httpClient *http.Client) *graphqlClient {
	return &graphqlClient{method: method, url: url, httpClient: httpClient}
}

func (c *graphqlClient) Do(ctx context.Context, request Request, data interface{}) error {
	l.Debug(ctx, "graphqlclient.request", l.NewValue("method", c.method), l.NewValue("url", c.url), l.NewValue("request", request))
	httpResponse, err := c.httpDo(ctx, request)
	if err != nil {
		return err
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return newGraphQLError(httpResponse)
	}

	err = decode(httpResponse.Body, data)
	if err != nil {
		return fmt.Errorf("err_unmarshall_response: %w", err)
	}
	l.Debug(ctx, "graphqlclient.response", l.NewValue("method", method), l.NewValue("url", url), l.NewValue("response", data), l.NewValue("err", err))

	return nil
}

func (c *graphqlClient) httpDo(ctx context.Context, request Request) (*http.Response, error) {
	httpRequest, err := newHTTPRequest(c.method, c.url, request)
	if err != nil {
		return nil, err
	}

	httpResponse, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}

	return httpResponse, nil
}

func newGraphQLError(httpResponse *http.Response) error {
	var (
		errs = []error{
			fmt.Errorf("err_graphql_status_code: %d", httpResponse.StatusCode),
		}
		response Response
	)

	err := decode(httpResponse.Body, &response)
	if err != nil {
		errs = append(errs, fmt.Errorf("err_unmarshall_error_response: %w", err))
	}

	for _, graphqlErr := range response.Errors {
		message, is := graphqlErr["message"].(string)
		if !is {
			message = fmt.Sprintf("%v", graphqlErr["message"])
		}
		errs = append(errs, errors.New(message))
	}

	return fmt.Errorf("err_graphql: %w", abend.List(errs))
}
