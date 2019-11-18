package graphql

import (
	"testing"

	collectionmock "github.com/rjansen/fivecolors/collection/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewHandler(t *testing.T) {
	readerMock := collectionmock.NewReader()
	writerMock := collectionmock.NewWriter()

	handler := NewHandler(readerMock, writerMock)
	assert.NotNil(t, handler)
}
