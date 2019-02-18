// +build integration,firestore

package model

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/rjansen/fivecolors/core/graphql"
	"github.com/rjansen/l"
	"github.com/rjansen/raizel/firestore"
	"github.com/rjansen/yggdrasil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testProjectID     = "e-pedion"
	testCollectionFmt = "environments/test/%s"
)

type testSchemaFirestore struct {
	name         string
	tree         yggdrasil.Tree
	schema       graphql.Schema
	mockSetup    map[string]interface{}
	mockTearDown []string
	data         interface{}
	request      graphql.Request
	response     *graphql.Response
}

func (scenario *testSchemaFirestore) setup(t *testing.T) {
	var (
		roots     = yggdrasil.NewRoots()
		errLogger = l.Register(&roots, l.NewZapLoggerDefault())
		client    = firestore.NewClient(testProjectID)
		errClient = firestore.Register(&roots, client)
	)
	require.Nil(t, errLogger, "setup logger error")
	require.Nil(t, errClient, "setup client error")

	tree := roots.NewTreeDefault()
	schema := NewSchema(tree)
	require.NotNil(t, schema, "new schema error")

	for key, document := range scenario.mockSetup {
		err := client.Doc(
			fmt.Sprintf(testCollectionFmt, key),
		).Set(
			context.Background(), document,
		)
		assert.Nilf(t, err, "error setup mock data: key=%s document=%+v err=%+v", key, document, err)
	}

	collectionFmt = testCollectionFmt
	scenario.tree = tree
	scenario.schema = schema
}

func (scenario *testSchemaFirestore) tearDown(t *testing.T) {
	if scenario.tree != nil {
		defer scenario.tree.Close()
		var (
			client = firestore.MustReference(scenario.tree)
		)
		for index, key := range scenario.mockTearDown {
			err := client.Doc(
				fmt.Sprintf(testCollectionFmt, key),
			).Delete(
				context.Background(),
			)
			assert.Nilf(t, err, "error terdown mock data: index=%d key=%s err=%+v", index, key, err)
		}
	}
}

func TestSchemaFirestore(test *testing.T) {
	timeNow := time.Now().UTC()
	scenarios := []testSchemaFirestore{
		{
			name: "Resolves card field successfully",
			mockSetup: map[string]interface{}{
				"rarities/mock_rarityid": Rarity{
					ID: "mock_rarityid", Name: RarityNameMythicRare, Alias: RarityAliasM,
				},
				"sets/mock_setid": Set{
					ID: "mock_setid", Name: "Set Mock", Alias: "stm",
					CreatedAt: timeNow,
				},
				"cards/mock_cardid": Card{
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
				},
			},
			mockTearDown: []string{
				"cards/mock_cardid", "sets/mock_setid", "rarities/mock_rarityid",
			},
			data: &struct {
				Card Card `json:"card"`
			}{},
			request: graphql.Request{
				Query: `{
					card(id: "mock_cardid") {
					  id
					  name
					  types
					  costs
					  numberCost
				      idAsset
					  data
					  createdAt
					  updatedAt
					  deletedAt
                    }
				}`,
			},
		},
		{
			name: "Resolves set field successfully",
			mockSetup: map[string]interface{}{
				"sets/mock_setid": Set{
					ID: "mock_setid", Name: "Set Mock", Alias: "stm",
					CreatedAt: timeNow,
				},
			},
			mockTearDown: []string{
				"sets/mock_setid",
			},
			data: &struct {
				Set Set `json:"set"`
			}{},
			request: graphql.Request{
				Query: `{
					set(id: "mock_setid") {
					  id
					  name
					  alias
					  cards {
						id
						name
						types
						costs
						numberCost
						idAsset
						data
						createdAt
						updatedAt
						deletedAt
						rarity {
						  id
						  name
						  alias
						}
					  }
					}
				}`,
			},
		},
		{
			name: "Resolves cardBy field successfully",
			mockSetup: map[string]interface{}{
				"rarities/mock_rarityid1": Rarity{
					ID: "mock_rarityid1", Name: RarityNameMythicRare, Alias: RarityAliasM,
				},
				"sets/mock_setid1": Set{
					ID: "mock_setid1", Name: "Set Mock", Alias: "stm",
					CreatedAt: timeNow,
				},
				"cards/mock_cardid1": Card{
					ID:          "mock_cardid1",
					IDExternal:  "mock_cardidexternal1",
					IDAsset:     "mock_cardassetid1",
					Name:        "Card Mock One",
					TypesObject: Object{"0": "Legendary", "1": "Creature", "2": "Elf"},
					CostsObject: &Object{"0": "1", "1": "G", "2": "G", "3": "G"},
					NumberCost:  4.0,
					CreatedAt:   timeNow,
					IDRarity:    "mock_rarityid1",
					IDSet:       "mock_setid1",
				},
				"cards/mock_cardid2": Card{
					ID:          "mock_cardid2",
					IDExternal:  "mock_cardidexternal2",
					IDAsset:     "mock_cardassetid2",
					Name:        "Card Mock Two",
					TypesObject: Object{"0": "Legendary", "1": "Instant"},
					CostsObject: &Object{"0": "1", "1": "B", "2": "B"},
					NumberCost:  3.0,
					CreatedAt:   timeNow,
					IDRarity:    "mock_rarityid1",
					IDSet:       "mock_setid1",
				},
				"cards/mock_cardid3": Card{
					ID:          "mock_cardid3",
					IDExternal:  "mock_cardidexternal3",
					IDAsset:     "mock_cardassetid3",
					Name:        "Card Mock Three",
					TypesObject: Object{"0": "Legendary", "1": "Creature", "2": "Goblin"},
					CostsObject: &Object{"0": "1", "1": "R", "2": "R", "3": "R"},
					NumberCost:  4.0,
					CreatedAt:   timeNow,
					IDRarity:    "mock_rarityid1",
					IDSet:       "mock_setid1",
				},
			},
			// mockTearDown: []string{
			// 	"cards/mock_cardid1", "cards/mock_cardid2", "cards/mock_cardid3",
			// 	"rarities/mock_rarityid1", "sets/mock_setid1",
			// },
			data: &struct {
				Card []Card `json:"cardBy"`
			}{},
			request: graphql.Request{
				Query: `{
					cardBy(filter: {
					  name: "Card Mock"
					  typesObject: ["Legendary", "Creature"]
					  costsObject: ["1", "G"]
					}) {
					  id
					  name
					  typesObject
					  costsObject
					  numberCost
					  idAsset
					  idSet
					  idRarity
					  data
					  createdAt
					  updatedAt
					  deletedAt
					}
				}`,
			},
		},
		/*
			{
				name: "Resolves setBy field successfully",
				mockSetup: []dataGenerator{
					{
						sql:       "insert into set (id, name, alias, created_at) values ($1, $2, $3, $4)",
						arguments: []interface{}{"mock_setid1", "Set Mock One", "stm1", timeNow},
					},
					{
						sql:       "insert into set (id, name, alias, created_at) values ($1, $2, $3, $4)",
						arguments: []interface{}{"mock_setid2", "Set Mock Two", "stm2", timeNow},
					},
					{
						sql:       "insert into set (id, name, alias, created_at) values ($1, $2, $3, $4)",
						arguments: []interface{}{"mock_setid3", "Set Mock Three", "stm3", timeNow},
					},
				},
				mockTearDown: []dataGenerator{
					{sql: "delete from set where id = $1", arguments: []interface{}{"mock_setid1"}},
					{sql: "delete from set where id = $1", arguments: []interface{}{"mock_setid2"}},
					{sql: "delete from set where id = $1", arguments: []interface{}{"mock_setid3"}},
				},

				data: &struct {
					Set []Set `json:"setBy"`
				}{},
				request: graphql.Request{
					Query: `{
						setBy(filter: {
						  name: "Set Mock"
						  alias: "stm"
						}) {
						  id
						  name
						  alias
						  cards {
							id
							name
							types
							costs
							numberCost
							idAsset
							data
							createdAt
							updatedAt
							deletedAt
							rarity {
							  id
							  name
							  alias
							}
						  }
						}
					}`,
				},
			},
		*/
	}

	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)
				defer scenario.tearDown(t)

				response := graphql.Execute(scenario.tree, scenario.schema, scenario.request)
				require.NotNil(t, response, "schema response invalid")
				require.Nilf(t, response.Errors, "schema response errors: %+v", response.Errors)
				require.NotNilf(t, response.Data, "schema response nil: %+v", response.Data)
				t.Logf("%s", response.Data)
				err := json.Unmarshal(response.Data, scenario.data)
				require.Nil(t, err, "schema response unmarshal error")
				require.NotZerof(t, scenario.data, "data response invalid: %+v", scenario.data)
				require.NotEmpty(t, scenario.data, "data response len invalid: %+v", scenario.data)
			},
		)
	}
}
