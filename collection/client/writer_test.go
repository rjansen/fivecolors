package client

import (
	"context"
	"net/http"
	"testing"

	"github.com/rjansen/fivecolors/collection"
	"github.com/stretchr/testify/assert"
)

func TestUpsertSet(t *testing.T) {
	test := httpTest{
		status: http.StatusOK,
		body: `{
			"data": {
				"upsertSet": {
					"id": "myset",
					"name": "My Set",
					"alias": "mys"
				}
			}
		}`,
	}
	test.setup(t)
	defer test.tearDown(t)

	var (
		writer   = NewWriter(test.httpClient)
		ctx      = context.Background()
		setInput = collection.SetInput{
			ID:    "myset",
			Name:  "My Set",
			Alias: "mys",
		}
	)
	set, err := writer.UpsertSet(ctx, setInput)
	assert.NoError(t, err)
	assert.NotNil(t, set)
	assert.Equal(t, set.ID, "myset")
	assert.Equal(t, set.Name, "My Set")
	assert.Equal(t, set.Alias, "mys")
}

func TestUpsertCards(t *testing.T) {
	test := httpTest{
		status: http.StatusOK,
		body: `{
			"data": {
				"upsertCards": {
					"affectedRecords": 3,
					"committedAt": "2019-10-01T00:00:00+00:00"
				}
			}
		}`,
	}
	test.setup(t)
	defer test.tearDown(t)

	var (
		writer     = NewWriter(test.httpClient)
		ctx        = context.Background()
		cardsInput = []collection.CardInput{
			{
				ID:    "mycard1",
				Name:  "My Card 1",
				Types: []string{"type1", "type2", "type3"},
			},
			{
				ID:    "mycard2",
				Name:  "My Card 2",
				Types: []string{"type4", "type5"},
			},
			{
				ID:    "mycard3",
				Name:  "My Card 3",
				Types: []string{"type8", "type9", "type10"},
			},
		}
	)
	cardsResult, err := writer.UpsertCards(ctx, cardsInput)
	assert.NoError(t, err)
	assert.NotNil(t, cardsResult)
	assert.Equal(t, cardsResult.AffectedRecords, len(cardsInput))
	assert.NotZero(t, cardsResult.CommittedAt)
}
