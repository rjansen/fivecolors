package client

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	test := httpTest{
		status: http.StatusOK,
		body: `{
			"data": {
				"set": {
					"id": "myset",
					"name": "My Set",
					"alias": "mys"
				}
			}
		}`,
	}
	test.setup(t)
	defer test.tearDown(t)

	reader := NewReader(test.httpClient)
	ctx := context.Background()
	set, err := reader.Set(ctx, "myset")
	assert.NoError(t, err)
	assert.NotNil(t, set)
	assert.Equal(t, set.ID, "myset")
	assert.Equal(t, set.Name, "My Set")
	assert.Equal(t, set.Alias, "mys")
}

func TestCard(t *testing.T) {
	test := httpTest{
		status: http.StatusOK,
		body: `{
			"data": {
				"card": {
					"id": "mycard",
					"name": "My Card",
					"types": ["type1", "type2", "type3"]
				}
			}
		}`,
	}
	test.setup(t)
	defer test.tearDown(t)

	reader := NewReader(test.httpClient)
	ctx := context.Background()
	card, err := reader.Card(ctx, "mycard")
	assert.NoError(t, err)
	assert.NotNil(t, card)
	assert.Equal(t, card.ID, "mycard")
	assert.Equal(t, card.Name, "My Card")
	assert.Equal(t, card.Types, []string{"type1", "type2", "type3"})
}
