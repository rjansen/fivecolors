package client

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type httpTest struct {
	status     int
	body       string
	httpServer *httptest.Server
	httpClient *http.Client
}

func (test *httpTest) setup(t *testing.T) {
	test.httpServer = newHTTPServerMock(test.status, test.body)
	assert.NotNil(t, test.httpServer)
	test.httpClient = test.httpServer.Client()
	assert.NotNil(t, test.httpClient)
	url = test.httpServer.URL

}

func (test *httpTest) tearDown(t *testing.T) {
	test.httpServer.Close()
	url = defaultURL
}
