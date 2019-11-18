package firestore

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/rjansen/fivecolors/collection"
	firestore "github.com/rjansen/raizel/firestore/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpsertSet(t *testing.T) {
	test := clientTest{
		setupTest: func(t *testing.T, test *clientTest) {
			mockDocument := firestore.NewDocumentRefMock()
			mockDocument.On("Set", mock.Anything, mock.Anything, mock.Anything).Run(
				func(args mock.Arguments) {
					docRaw := args.Get(1)
					if docRaw != nil {
						reflectTime := reflect.ValueOf(time.Now().UTC())
						reflect.ValueOf(docRaw).Elem().FieldByName("UpdatedAt").Set(reflectTime)
					}
				},
			).Return(nil)
			test.client.On("Doc", mock.Anything).Return(mockDocument)
		},
	}
	test.setup(t)
	defer test.tearDown(t)

	var (
		writer = NewWriter(test.logger, test.client)
		ctx    = context.Background()
		set    = collection.SetInput{
			ID:    "myset",
			Name:  "My Set",
			Alias: "mys",
		}
	)

	result, err := writer.UpsertSet(ctx, set)
	assert.NoError(t, err)
	assert.Equal(t, result.ID, "myset")
	assert.Equal(t, result.Name, "My Set")
	assert.Equal(t, result.Alias, "mys")
	assert.NotZero(t, result.UpdatedAt)
}

func TestUpsertCards(t *testing.T) {
	test := clientTest{
		setupTest: func(t *testing.T, test *clientTest) {
			mockDocument := firestore.NewDocumentRefMock()
			test.client.On("Doc", mock.Anything).Return(mockDocument).Times(5)
			mockBatch := firestore.NewWriteBatchMock()
			mockBatch.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(mockBatch).Times(5)
			mockBatch.On("Commit", mock.Anything).Return(nil)
			test.client.On("Batch").Return(mockBatch)
		},
	}
	test.setup(t)
	defer test.tearDown(t)

	var (
		writer = NewWriter(test.logger, test.client)
		ctx    = context.Background()
		cards  = []collection.CardInput{
			{
				ID:    "mycard1",
				Name:  "My Card 1",
				Types: []string{"type1", "type2", "type3"},
			},
			{
				ID:    "mycard2",
				Name:  "My Card 2",
				Types: []string{"type1", "type2", "type3"},
			},
			{
				ID:    "mycard3",
				Name:  "My Card 3",
				Types: []string{"type1", "type2", "type3"},
			},
			{
				ID:    "mycard4",
				Name:  "My Card 4",
				Types: []string{"type1", "type2", "type3"},
			},
			{
				ID:    "mycard5",
				Name:  "My Card 5",
				Types: []string{"type1", "type2", "type3"},
			},
		}
	)

	result, err := writer.UpsertCards(ctx, cards)
	assert.NoError(t, err)
	assert.Equal(t, result.AffectedRecords, len(cards))
	assert.NotZero(t, result.CommittedAt)
}
