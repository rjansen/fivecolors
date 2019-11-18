package graphql

import (
	"testing"

	collectionmock "github.com/rjansen/fivecolors/collection/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewExecutableSchema(t *testing.T) {
	readerMock := collectionmock.NewReader()
	writerMock := collectionmock.NewWriter()

	resolver := NewResolver(readerMock, writerMock)
	schema := NewExecutableSchema(Config{Resolvers: resolver})
	assert.NotNil(t, schema)
}
