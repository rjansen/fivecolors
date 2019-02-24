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

	collectionFmt = testCollectionFmt
	scenario.tree = tree
	scenario.schema = schema
}

func (scenario *testSchemaFirestore) tearDown(t *testing.T) {
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
					Rarity: &Rarity{
						ID: "mock_rarityid", Name: RarityNameMythicRare, Alias: RarityAliasM,
					},
					Set: &Set{
						ID: "mock_setid", Name: "Set Mock", Alias: "stm",
						CreatedAt: timeNow,
					},
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
					  createdAt
					  updatedAt
					  deletedAt
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
					ID:         "mock_cardid1",
					IDExternal: "mock_cardidexternal1",
					IDAsset:    "mock_cardassetid1",
					Name:       "Card Mock One",
					Types:      []string{"Legendary", "Creature", "Elf"},
					Costs:      []string{"1", "G", "G", "G"},
					NumberCost: 4.0,
					CreatedAt:  timeNow,
					IDRarity:   "mock_rarityid1",
					IDSet:      "mock_setid1",
					Rarity: &Rarity{
						ID: "mock_rarityid1", Name: RarityNameMythicRare, Alias: RarityAliasM,
					},
					Set: &Set{
						ID: "mock_setid1", Name: "Set Mock", Alias: "stm",
						CreatedAt: timeNow,
					},
				},
				"cards/mock_cardid2": Card{
					ID:         "mock_cardid2",
					IDExternal: "mock_cardidexternal2",
					IDAsset:    "mock_cardassetid2",
					Name:       "Card Mock Two",
					Types:      []string{"Legendary", "Instant"},
					Costs:      []string{"1", "B", "B"},
					NumberCost: 3.0,
					CreatedAt:  timeNow,
					IDRarity:   "mock_rarityid1",
					IDSet:      "mock_setid1",
					Rarity: &Rarity{
						ID: "mock_rarityid1", Name: RarityNameMythicRare, Alias: RarityAliasM,
					},
					Set: &Set{
						ID: "mock_setid1", Name: "Set Mock", Alias: "stm",
						CreatedAt: timeNow,
					},
				},
				"cards/mock_cardid3": Card{
					ID:         "mock_cardid3",
					IDExternal: "mock_cardidexternal3",
					IDAsset:    "mock_cardassetid3",
					Name:       "Card Mock Three",
					Types:      []string{"Legendary", "Creature", "Goblin"},
					Costs:      []string{"1", "R", "R", "R"},
					NumberCost: 4.0,
					CreatedAt:  timeNow,
					IDRarity:   "mock_rarityid1",
					IDSet:      "mock_setid1",
					Rarity: &Rarity{
						ID: "mock_rarityid1", Name: RarityNameMythicRare, Alias: RarityAliasM,
					},
					Set: &Set{
						ID: "mock_setid1", Name: "Set Mock", Alias: "stm",
						CreatedAt: timeNow,
					},
				},
			},
			mockTearDown: []string{
				"cards/mock_cardid1", "cards/mock_cardid2", "cards/mock_cardid3",
				"rarities/mock_rarityid1", "sets/mock_setid1",
			},
			data: &struct {
				Card []Card `json:"cardBy"`
			}{},
			request: graphql.Request{
				Query: `{
					cardBy(filter: {
					  name: "Card Mock"
					  types: ["Legendary", "Instant"]
					  costs: ["1", "B", "B", "B"]
					  set: {
						id: "mock_setid1"
						name: "Set Mock"
						alias: "stm"
					  }
					  rarity: {
						id: "mock_rarityid1"
						name: MythicRare
						alias: M
					  }
					}) {
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
			},
		},
		{
			name: "Resolves cardBy field with another filter successfully",
			mockSetup: map[string]interface{}{
				"rarities/mock_rarityid1": Rarity{
					ID: "mock_rarityid1", Name: RarityNameMythicRare, Alias: RarityAliasM,
				},
				"sets/mock_setid1": Set{
					ID: "mock_setid1", Name: "Set Mock", Alias: "stm",
					CreatedAt: timeNow,
				},
				"cards/mock_cardid1": Card{
					ID:         "mock_cardid1",
					IDExternal: "mock_cardidexternal1",
					IDAsset:    "mock_cardassetid1",
					Name:       "Card Mock One",
					Types:      []string{"Legendary", "Creature", "Elf"},
					Costs:      []string{"1", "G", "G", "G"},
					NumberCost: 4.0,
					CreatedAt:  timeNow,
					IDRarity:   "mock_rarityid1",
					IDSet:      "mock_setid1",
					Rarity: &Rarity{
						ID: "mock_rarityid1", Name: RarityNameMythicRare, Alias: RarityAliasM,
					},
					Set: &Set{
						ID: "mock_setid1", Name: "Set Mock", Alias: "stm",
						CreatedAt: timeNow,
					},
				},
				"cards/mock_cardid2": Card{
					ID:         "mock_cardid2",
					IDExternal: "mock_cardidexternal2",
					IDAsset:    "mock_cardassetid2",
					Name:       "Card Mock Two",
					Types:      []string{"Legendary", "Instant"},
					Costs:      []string{"1", "B", "B"},
					NumberCost: 3.0,
					CreatedAt:  timeNow,
					IDRarity:   "mock_rarityid1",
					IDSet:      "mock_setid1",
					Rarity: &Rarity{
						ID: "mock_rarityid1", Name: RarityNameMythicRare, Alias: RarityAliasM,
					},
					Set: &Set{
						ID: "mock_setid1", Name: "Set Mock", Alias: "stm",
						CreatedAt: timeNow,
					},
				},
				"cards/mock_cardid3": Card{
					ID:         "mock_cardid3",
					IDExternal: "mock_cardidexternal3",
					IDAsset:    "mock_cardassetid3",
					Name:       "Card Mock Three",
					Types:      []string{"Legendary", "Creature", "Goblin"},
					Costs:      []string{"1", "R", "R", "R"},
					NumberCost: 4.0,
					CreatedAt:  timeNow,
					IDRarity:   "mock_rarityid1",
					IDSet:      "mock_setid1",
					Rarity: &Rarity{
						ID: "mock_rarityid1", Name: RarityNameMythicRare, Alias: RarityAliasM,
					},
					Set: &Set{
						ID: "mock_setid1", Name: "Set Mock", Alias: "stm",
						CreatedAt: timeNow,
					},
				},
			},
			mockTearDown: []string{
				"cards/mock_cardid1", "cards/mock_cardid2", "cards/mock_cardid3",
				"rarities/mock_rarityid1", "sets/mock_setid1",
			},
			data: &struct {
				Card []Card `json:"cardBy"`
			}{},
			request: graphql.Request{
				Query: `{
					cardBy(filter: {
					  name: "Card Mock"
					  types: ["Legendary", "Instant"]
					  costs: ["1", "B", "B"]
					}) {
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
			},
		},
		/*
			{
				name: "Resolves setBy field successfully",
				mockSetup: map[string]interface{}{
					"sets/mock_setid1": Set{
						ID: "mock_setid1", Name: "Set Mocki One", Alias: "stm1",
						CreatedAt: timeNow,
					},
					"sets/mock_setid2": Set{
						ID: "mock_setid2", Name: "Set Mock Two", Alias: "stm2",
						CreatedAt: timeNow,
					},
					"sets/mock_setid3": Set{
						ID: "mock_setid3", Name: "Set Mock Three", Alias: "stm3",
						CreatedAt: timeNow,
					},
					"sets/mock_setid4": Set{
						ID: "mock_setid4", Name: "Set Mock Four", Alias: "stm4",
						CreatedAt: timeNow,
					},
				},
				mockTearDown: []string{
					"sets/mock_setid1",
					"sets/mock_setid2",
					"sets/mock_setid3",
					"sets/mock_setid4",
				},
				data: &struct {
					Set Set `json:"set"`
				}{},
				request: graphql.Request{
					Query: `{
						setBy(filter: {
						  name: "Set Mock"
						  alias: "stm4"
						}) {
						  id
						  name
						  alias
						  createdAt
						  updatedAt
						  deletedAt
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
				require.NotEmpty(t, scenario.data, "data response len invalid: %+v", scenario.data)
			},
		)
	}
}
