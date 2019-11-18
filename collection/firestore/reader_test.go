package firestore

import (
	"context"
	"testing"
	"time"

	"github.com/rjansen/fivecolors/collection"
	firestore "github.com/rjansen/raizel/firestore/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSet(t *testing.T) {
	test := clientTest{
		setupTest: func(t *testing.T, test *clientTest) {
			mockSnapshot := firestore.NewDocumentSnapshotMock()
			mockSnapshot.On("DataTo", mock.Anything).Run(
				func(args mock.Arguments) {
					pointer := args.Get(0)
					if pointer != nil {
						if set, is := pointer.(*collection.Set); is {
							*set = collection.Set{
								ID:        "myset",
								Name:      "My Set",
								Alias:     "mys",
								UpdatedAt: time.Now().UTC(),
							}
						}
					}
				},
			).Return(nil)
			mockDocument := firestore.NewDocumentRefMock()
			mockDocument.On("Get", mock.Anything).Return(mockSnapshot, nil)
			test.client.On("Doc", mock.Anything).Return(mockDocument)
		},
	}
	test.setup(t)
	defer test.tearDown(t)

	var (
		reader = NewReader(test.logger, test.client)
		ctx    = context.Background()
		id     = "myset"
	)

	set, err := reader.Set(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, set.ID, "myset")
	assert.Equal(t, set.Name, "My Set")
	assert.Equal(t, set.Alias, "mys")
	assert.NotZero(t, set.UpdatedAt)
}

func TestCard(t *testing.T) {
	test := clientTest{
		setupTest: func(t *testing.T, test *clientTest) {
			mockSnapshot := firestore.NewDocumentSnapshotMock()
			mockSnapshot.On("DataTo", mock.Anything).Run(
				func(args mock.Arguments) {
					pointer := args.Get(0)
					if pointer != nil {
						if card, is := pointer.(*collection.Card); is {
							*card = collection.Card{
								ID:        "mycard",
								Name:      "My Card",
								Types:     []string{"type1", "type2", "type3", "type4"},
								UpdatedAt: time.Now().UTC(),
							}
						}
					}
				},
			).Return(nil)
			mockDocument := firestore.NewDocumentRefMock()
			mockDocument.On("Get", mock.Anything).Return(mockSnapshot, nil)
			test.client.On("Doc", mock.Anything).Return(mockDocument)
		},
	}
	test.setup(t)
	defer test.tearDown(t)

	var (
		reader = NewReader(test.logger, test.client)
		ctx    = context.Background()
		id     = "mycard"
	)

	card, err := reader.Card(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, card.ID, "mycard")
	assert.Equal(t, card.Name, "My Card")
	assert.Equal(t, card.Types, []string{"type1", "type2", "type3", "type4"})
	assert.NotZero(t, card.UpdatedAt)
}
