// +build integration

package function

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/rjansen/fivecolors/core/api"
	"github.com/rjansen/fivecolors/core/model"
	"github.com/rjansen/raizel/firestore"
	"github.com/rjansen/yggdrasil"
	"github.com/stretchr/testify/require"
)

const (
	testProjectID     = "e-pedion"
	testCollectionFmt = "environments/test/%s"
)

type testFirestoreIntegrationHandler struct {
	name           string
	tree           yggdrasil.Tree
	mockSetup      map[string]interface{}
	mockTearDown   []string
	body           string
	contentType    string
	request        *http.Request
	response       *httptest.ResponseRecorder
	responseStatus int
}

func (scenario *testFirestoreIntegrationHandler) setup(t *testing.T) {
	var (
		options = options{
			projectID: testProjectID,
			dataStore: "firestore",
		}
		tree   = newTree(options)
		client = firestore.MustReference(tree)
	)
	require.NotNil(t, tree, "tree invalida instance")
	require.NotNil(t, client, "client invalid instance")

	if len(scenario.mockSetup) > 0 {
		batch := client.Batch()
		for key, document := range scenario.mockSetup {
			batch.Set(
				client.Doc(
					fmt.Sprintf(testCollectionFmt, key),
				),
				document,
			)
		}
		err := batch.Commit(context.Background())
		require.Nilf(t, err, "error setup mock data: data=%s err=%+v", scenario.mockSetup, err)
	}

	r := httptest.NewRequest(
		"POST", "/", strings.NewReader(scenario.body),
	)
	r.Header.Set("content-type", scenario.contentType)

	serverHandler = api.NewGraphQLHandler(tree)
	scenario.tree = tree
	scenario.request = r
	scenario.response = httptest.NewRecorder()
}

func (scenario *testFirestoreIntegrationHandler) tearDown(t *testing.T) {
	if scenario.tree != nil {
		defer scenario.tree.Close()
		if len(scenario.mockTearDown) > 0 {
			var (
				client = firestore.MustReference(scenario.tree)
				batch  = client.Batch()
			)
			for _, key := range scenario.mockTearDown {
				batch.Delete(client.Doc(fmt.Sprintf(testCollectionFmt, key)))
			}
			err := batch.Commit(context.Background())
			require.Nilf(t, err, "error terdown mock data: data=%s err=%+v", scenario.mockSetup, err)
		}
	}
}

func TestFirestoreIntegrationHandler(test *testing.T) {
	timeNow := time.Now().UTC()
	scenarios := []testFirestoreIntegrationHandler{
		{
			name:           "When request body is invalid returns a bad request",
			body:           "<xml><id>xmlid</id></name>Invalid Body</name></xml>",
			contentType:    "application/xml",
			responseStatus: http.StatusBadRequest,
		},
		{
			name: "When request body is graphql executes the query and returns ok with query results",
			mockSetup: map[string]interface{}{
				"rarities/mock_rarityid": model.Rarity{
					ID: "mock_rarityid", Name: model.RarityNameMythicRare, Alias: model.RarityAliasM,
				},
				"sets/mock_setid": model.Set{
					ID: "mock_setid", Name: "Set Mock", Alias: "stm",
					CreatedAt: timeNow,
				},
				"cards/mock_cardid": model.Card{
					ID:         "mock_cardid",
					IDExternal: "mock_cardidexternal",
					IDAsset:    "mock_cardassetid",
					Name:       "Card Mock",
					Types:      []string{"Legendary", "Creature", "Goblin"},
					Costs:      []string{"1", "R", "R", "R"},
					NumberCost: 4.0,
					CreatedAt:  timeNow,
					IDRarity:   "mock_rarityid",
					IDSet:      "mock_setid",
					Rarity: &model.Rarity{
						ID: "mock_rarityid", Name: model.RarityNameMythicRare, Alias: model.RarityAliasM,
					},
					Set: &model.Set{
						ID: "mock_setid", Name: "Set Mock", Alias: "stm",
						CreatedAt: timeNow,
					},
				},
			},
			mockTearDown: []string{
				"cards/mock_cardid", "sets/mock_setid", "rarities/mock_rarityid",
			},
			body: `{
				card(id: "mock_cardid") {
				  id
				  name
				  types
				  costs
				  numberCost
				  idAsset
				  idSet
				  set{
					id
					name
					alias
				  }
				  idRarity
				  rarity{
					id
					name
					alias
				  }
				  data
				  createdAt
				  updatedAt
				  deletedAt
				}
			}`,
			contentType:    "application/graphql",
			responseStatus: http.StatusOK,
		},
		{
			name: "When request body is a valid json executes the query and returns ok with query results",
			mockSetup: map[string]interface{}{
				"rarities/mock_rarityid": model.Rarity{
					ID: "mock_rarityid", Name: model.RarityNameMythicRare, Alias: model.RarityAliasM,
				},
				"sets/mock_setid": model.Set{
					ID: "mock_setid", Name: "Set Mock", Alias: "stm",
					CreatedAt: timeNow,
				},
				"cards/mock_cardid": model.Card{
					ID:         "mock_cardid",
					IDExternal: "mock_cardidexternal",
					IDAsset:    "mock_cardassetid",
					Name:       "Card Mock",
					Types:      []string{"Legendary", "Creature", "Goblin"},
					Costs:      []string{"1", "R", "R", "R"},
					NumberCost: 4.0,
					CreatedAt:  timeNow,
					IDRarity:   "mock_rarityid",
					IDSet:      "mock_setid",
					Rarity: &model.Rarity{
						ID: "mock_rarityid", Name: model.RarityNameMythicRare, Alias: model.RarityAliasM,
					},
					Set: &model.Set{
						ID: "mock_setid", Name: "Set Mock", Alias: "stm",
						CreatedAt: timeNow,
					},
				},
			},
			mockTearDown: []string{
				"cards/mock_cardid", "sets/mock_setid", "rarities/mock_rarityid",
			},
			body: `{
				"query": "{card(id: \"mock_cardid\") {id,name,types,costs,numberCost,idAsset,idSet,set{id,name,alias},idRarity,rarity{id,name,alias},data,createdAt,updatedAt,deletedAt}}"
			}`,
			contentType:    "application/json",
			responseStatus: http.StatusOK,
		},
	}

	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)
				defer scenario.tearDown(t)

				Handler(scenario.response, scenario.request)
				result := scenario.response.Result()
				require.Equal(t, scenario.responseStatus, result.StatusCode, "invalid response statuscode")
			},
		)
	}
}
