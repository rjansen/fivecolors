package graphql

import (
	"testing"

	collectionmock "github.com/rjansen/fivecolors/collection/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewResolver(t *testing.T) {
	readerMock := collectionmock.NewReader()
	writerMock := collectionmock.NewWriter()

	resolver := NewResolver(readerMock, writerMock)
	assert.NotNil(t, resolver)
	assert.NotNil(t, resolver.Query())
	assert.NotNil(t, resolver.Mutation())
}
