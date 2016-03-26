package data_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	//"errors"
	"testing"

	"farm.e-pedion.com/repo/fivecolors/config"
	"farm.e-pedion.com/repo/fivecolors/data"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

var (
	setted           = false
	minimalDeck      *data.Deck
	fullDeck         *data.Deck
	minimalInventory *data.Inventory
	fullInventory    *data.Inventory
)

func setup() error {
	setupErr := data.Setup(&config.DBConfig{
		Driver:   "mysql",
		URL:      "tcp(127.0.0.1:3306)/fivecolors",
		Username: "fivecolors",
		Password: "fivecolors",
	})
	if setupErr != nil {
		setted = true
	}
	return setupErr
}

func before() error {
	if !setted {
		if err := setup(); err != nil {
			return err
		}
	}
	return nil
}

//Card
func Test_CardRead(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	log.Println("Test_CardRead")
	card := &data.Card{}
	card.ID = 1
	readErr := card.Read()
	assert.Nil(t, readErr)
	assert.Equal(t, card.ID, 1)
	assert.NotEqual(t, card.Name, "")
	assert.NotEqual(t, card.Label, "")
	assert.NotEqual(t, card.TypeLabel, "")
	assert.True(t, card.IDRarity >= 0)
	assert.NotEqual(t, card.Artist, "")
	assert.NotEqual(t, card.IDAsset, "")
	assert.NotEqual(t, card.InventoryCard.InventoryID, 0)
	assert.False(t, card.InventoryCard.Quantity < 0)
}

func Test_CardQuery(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	log.Println("Test_CardQuery")
	cardParameter := make(map[string]interface{}, 2)
	cardParameter["c.id_expansion = ?"] = 19
	card := data.Card{}
	cards, readErr := card.Query(cardParameter, "")
	assert.Nil(t, readErr)
	assert.NotNil(t, cards)
	assert.True(t, len(cards) > 0)
	for _, value := range cards {
		card := value.(data.Card)
		assert.NotEqual(t, card.ID, 0)
		assert.NotEqual(t, card.Name, "")
		assert.NotEqual(t, card.Label, "")
		assert.NotEqual(t, card.TypeLabel, "")
		assert.True(t, card.IDRarity >= 0)
		assert.NotEqual(t, card.Artist, "")
		assert.NotEqual(t, card.IDAsset, "")
		assert.NotEqual(t, card.InventoryCard.InventoryID, 0)
		assert.False(t, card.InventoryCard.Quantity < 0)
	}
}

//Card

//Expansion
func Test_ExpansionRead(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	log.Println("Test_ExpansionRead")
	expansion := &data.Expansion{}
	expansion.ID = 1
	readErr := expansion.Read()
	assert.Nil(t, readErr)
	assert.Equal(t, expansion.ID, 1)
	assert.NotEqual(t, expansion.Name, "")
	assert.NotEqual(t, expansion.Label, "")
	assert.NotEqual(t, expansion.IDAsset, 0)
}

func Test_ExpansionQuery(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	log.Println("Test_ExpansionQuery")
	expansion := &data.Expansion{}
	expansions, readErr := expansion.Query(nil, "")
	assert.Nil(t, readErr)
	assert.NotNil(t, expansions)
	assert.True(t, len(expansions) > 0)
	for _, value := range expansions {
		expansion := value.(data.Expansion)
		assert.NotEqual(t, expansion.ID, 0)
		assert.NotEqual(t, expansion.Name, "")
		assert.NotEqual(t, expansion.Label, "")
		assert.NotEqual(t, expansion.IDAsset, 0)
	}
}

//Expansion

//Asset
func Test_AssetRead(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	log.Println("Test_AssetRead")
	asset := &data.Asset{}
	asset.ID = 1
	readErr := asset.Read()
	assert.Nil(t, readErr)
	assert.Equal(t, asset.ID, 1)
	assert.NotEqual(t, asset.Label, "")
	assert.NotNil(t, asset.BinaryData)
	assetFile := fmt.Sprintf("%v/%v.jpg", os.TempDir(), asset.ID)
	log.Printf("CreatingAssetFile: FilePath=%v", assetFile)
	writeImgErr := ioutil.WriteFile(assetFile, asset.BinaryData, 0644)
	assert.Nil(t, writeImgErr)
}

//Asset

//Inventory
//Inventory Minimal
func Test_InventoryMinimalCreate(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	log.Println("Test_InventoryMinimalCreate")
	inventory := &data.Inventory{}
	inventory.Name = "Test_InventoryMinimalCreate"
	createErr := inventory.Persist()
	assert.Nil(t, createErr)
	assert.NotEqual(t, inventory.ID, 0)
	minimalInventory = inventory
}

func Test_InventoryMinimalCreatedRead(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	log.Println("Test_InventoryMinimalCreatedRead")
	inventory := &data.Inventory{}
	inventory.ID = minimalInventory.ID
	readErr := inventory.Read()
	assert.Nil(t, readErr)
	assert.Equal(t, inventory.ID, minimalInventory.ID)
	assert.Equal(t, inventory.Name, minimalInventory.Name)
}

func Test_InventoryMinimalUpdate(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	log.Println("Test_InventoryMinimalUpdate")
	inventory := &data.Inventory{}
	inventory.ID = minimalInventory.ID
	inventory.Name = "Test_InventoryMinimalUpdate"
	createErr := inventory.Persist()
	assert.Nil(t, createErr)
	minimalInventory = inventory
}

func Test_InventoryMinimalUpdatedRead(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	log.Println("Test_InventoryMinimalUpdatedRead")
	inventory := &data.Inventory{}
	inventory.ID = minimalInventory.ID
	readErr := inventory.Read()
	assert.Nil(t, readErr)
	assert.Equal(t, inventory.ID, minimalInventory.ID)
	assert.Equal(t, inventory.Name, minimalInventory.Name)
}

func Test_InventoryMinimalDelete(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	log.Println("Test_InventoryDelete")
	inventory := &data.Inventory{}
	inventory.ID = minimalInventory.ID
	deleteErr := inventory.Delete()
	assert.Nil(t, deleteErr)
	readErr := inventory.Read()
	assert.NotNil(t, readErr)
}

//Inventory Minimal

//Inventory Full
func Test_InventoryCreate(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	log.Println("Test_InventoryCreate")
	inventory := &data.Inventory{}
	inventory.Name = "Test_InventoryCreate"
	inventory.Cards = []data.Card{
		data.Card{ID: 1, InventoryCard: data.InventoryCard{Quantity: 2}},
		data.Card{ID: 2, InventoryCard: data.InventoryCard{Quantity: 5}},
		data.Card{ID: 3, InventoryCard: data.InventoryCard{Quantity: 4}},
		data.Card{ID: 4, InventoryCard: data.InventoryCard{Quantity: 4}},
		data.Card{ID: 5, InventoryCard: data.InventoryCard{Quantity: 27}},
	}
	createErr := inventory.Persist()
	assert.Nil(t, createErr)
	assert.NotEqual(t, inventory.ID, 0)
	fullInventory = inventory
}

func Test_InventoryCreatedRead(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	log.Println("Test_InventoryCreatedRead")
	inventory := &data.Inventory{}
	inventory.ID = fullInventory.ID
	readErr := inventory.Read()
	assert.Nil(t, readErr)
	assert.Equal(t, inventory.ID, fullInventory.ID)
	assert.Equal(t, inventory.Name, fullInventory.Name)
    //Read Fully
    readErr = inventory.ReadCards(-1)
    assert.Nil(t, readErr)
	cardsLen := len(inventory.Cards)
	assert.Equal(t, cardsLen, len(fullInventory.Cards))
	for k := 0; k < cardsLen; k++ {
		fullInventoryCard := fullInventory.Cards[k]
		inventoryCard := inventory.Cards[k]
		assert.Equal(t, inventoryCard.ID, fullInventoryCard.ID)
		assert.Equal(t, inventoryCard.InventoryCard.InventoryID, fullInventory.ID)
		assert.Equal(t, inventoryCard.InventoryCard.Quantity, fullInventoryCard.InventoryCard.Quantity)
	}
}

func Test_InventoryUpdate(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	log.Println("Test_InventoryUpdate")
	inventory := &data.Inventory{}
	inventory.ID = fullInventory.ID
	inventory.Name = "Test_InventoryUpdate"
	inventory.Cards = []data.Card{
		//Mock old cards inventory never removes a inserted cards, but change his quantity
        data.Card{ID: 1, InventoryCard: data.InventoryCard{Quantity: 2}},
		data.Card{ID: 2, InventoryCard: data.InventoryCard{Quantity: 5}},
		data.Card{ID: 3, InventoryCard: data.InventoryCard{Quantity: 4}},
		data.Card{ID: 4, InventoryCard: data.InventoryCard{Quantity: 4}},
		data.Card{ID: 5, InventoryCard: data.InventoryCard{Quantity: 27}},
		//New Cards
        data.Card{ID: 530, InventoryCard: data.InventoryCard{Quantity: 2}},
		data.Card{ID: 540, InventoryCard: data.InventoryCard{Quantity: 20}},
		data.Card{ID: 550, InventoryCard: data.InventoryCard{Quantity: 10}},
		data.Card{ID: 1, InventoryCard: data.InventoryCard{Quantity: 1}},
	}
	createErr := inventory.Persist()
	assert.Nil(t, createErr)
	fullInventory = inventory
}

func Test_InventoryUpdatedRead(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	log.Println("Test_InventoryUpdatedRead")
	inventory := &data.Inventory{}
	inventory.ID = fullInventory.ID
	readErr := inventory.Read()
	assert.Nil(t, readErr)
	assert.Equal(t, inventory.ID, fullInventory.ID)
	assert.Equal(t, inventory.Name, fullInventory.Name)
    //Read Fully
    readErr = inventory.ReadCards(-1)
    assert.Nil(t, readErr)
	cardsLen := len(inventory.Cards)
	for k := 0; k < cardsLen; k++ {
		inventoryCard := inventory.Cards[k]
		assert.NotEqual(t, inventoryCard.ID, 0)
		assert.Equal(t, inventoryCard.InventoryCard.InventoryID, fullInventory.ID)
		assert.True(t, inventoryCard.InventoryCard.Quantity >= 0)
	}
}

func Test_InventoryDelete(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	log.Println("Test_InventoryDelete")
	inventory := &data.Inventory{}
	inventory.ID = fullInventory.ID
	deleteErr := inventory.Delete()
	assert.Nil(t, deleteErr)
	readErr := inventory.Read()
	assert.NotNil(t, readErr)
}

//Inventory Full
//Inventory

//Deck
//Deck Minimal
func Test_DeckMinimalCreate(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	log.Println("Test_DeckMinimalCreate")
	deck := &data.Deck{}
	deck.Name = "Test_DeckMinimalCreate"
	createErr := deck.Persist()
	assert.Nil(t, createErr)
	assert.NotEqual(t, deck.ID, 0)
	minimalDeck = deck
}

func Test_DeckMinialCreatedRead(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	log.Println("Test_DeckMinialCreatedRead")
	deck := &data.Deck{}
	deck.ID = minimalDeck.ID
	readErr := deck.Read()
	assert.Nil(t, readErr)
}

func Test_DeckMinimalUpdate(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	log.Println("Test_DeckMinimalUpdate")
	deck := &data.Deck{}
	deck.ID = minimalDeck.ID
	deck.Name = "Test_DeckMinimalUpdate"
	createErr := deck.Persist()
	assert.Nil(t, createErr)
	minimalDeck = deck
}

func Test_DeckMinimalUpdatedRead(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	log.Println("Test_DeckMinimalUpdatedRead")
	deck := &data.Deck{}
	deck.ID = minimalDeck.ID
	readErr := deck.Read()
	assert.Nil(t, readErr)
	assert.Equal(t, deck.ID, minimalDeck.ID)
	assert.Equal(t, deck.Name, minimalDeck.Name)
}

func Test_DeckMinialDelete(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	log.Println("Test_DeckMinialDelete")
	deck := &data.Deck{}
	deck.ID = minimalDeck.ID
	deleteErr := deck.Delete()
	assert.Nil(t, deleteErr)
	readErr := deck.Read()
	assert.NotNil(t, readErr)
}

//Deck Minimal

//Deck Full
func Test_DeckCreate(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	log.Println("Test_DeckCreate")
	deck := &data.Deck{}
	deck.Name = "Test_DeckCreate"
	deck.MainCards = []data.Card{
		data.Card{ID: 1, DeckCard: data.DeckCard{BoardID: data.MainBoard, Quantity: 2}},
		data.Card{ID: 2, DeckCard: data.DeckCard{BoardID: data.MainBoard, Quantity: 4}},
		data.Card{ID: 3, DeckCard: data.DeckCard{BoardID: data.MainBoard, Quantity: 10}},
		data.Card{ID: 4, DeckCard: data.DeckCard{BoardID: data.MainBoard, Quantity: 3}},
	}
	deck.SideCards = []data.Card{
		data.Card{ID: 501, DeckCard: data.DeckCard{BoardID: data.SideBoard, Quantity: 5}},
		data.Card{ID: 502, DeckCard: data.DeckCard{BoardID: data.SideBoard, Quantity: 27}},
	}
	createErr := deck.Persist()
	assert.Nil(t, createErr)
	assert.NotEqual(t, deck.ID, 0)
	fullDeck = deck
}

func Test_DeckCreatedRead(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	log.Printf("Test_DeckCreatedRead: MainBoard=%v SideBoard=%v", data.MainBoard, data.SideBoard)
	deck := &data.Deck{}
	deck.ID = fullDeck.ID
	readErr := deck.Read()
	assert.Nil(t, readErr)
	assert.Equal(t, deck.ID, fullDeck.ID)
	assert.Equal(t, deck.Name, fullDeck.Name)

	mainCardsLen := len(deck.MainCards)
	assert.Equal(t, mainCardsLen, len(fullDeck.MainCards))
	for k := 0; k < mainCardsLen; k++ {
		fullDeckMainCard := fullDeck.MainCards[k]
		deckMainCard := deck.MainCards[k]
		assert.Equal(t, deckMainCard.ID, fullDeckMainCard.ID)
		assert.Equal(t, deckMainCard.DeckCard.DeckID, fullDeck.ID)
		assert.Equal(t, deckMainCard.DeckCard.BoardID, data.MainBoard)
		assert.Equal(t, deckMainCard.DeckCard.Quantity, fullDeckMainCard.DeckCard.Quantity)
	}
	sideCardsLen := len(deck.SideCards)
	assert.Equal(t, sideCardsLen, len(fullDeck.SideCards))
	for k := 0; k < sideCardsLen; k++ {
		fullDeckSideCard := fullDeck.SideCards[k]
		deckSideCard := deck.SideCards[k]
		assert.Equal(t, deckSideCard.ID, fullDeckSideCard.ID)
		assert.Equal(t, deckSideCard.DeckCard.DeckID, fullDeck.ID)
		assert.Equal(t, deckSideCard.DeckCard.BoardID, data.SideBoard)
		assert.Equal(t, deckSideCard.DeckCard.Quantity, fullDeckSideCard.DeckCard.Quantity)
	}
}

func Test_DeckUpdate(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	log.Println("Test_DeckUpdate")
	deck := &data.Deck{}
	deck.ID = fullDeck.ID
	deck.Name = "Test_DeckUpdate"
	deck.MainCards = []data.Card{
		data.Card{ID: 1, DeckCard: data.DeckCard{BoardID: data.MainBoard, Quantity: 2}},
		data.Card{ID: 2, DeckCard: data.DeckCard{BoardID: data.MainBoard, Quantity: 4}},
		data.Card{ID: 20, DeckCard: data.DeckCard{BoardID: data.MainBoard, Quantity: 3}},
		data.Card{ID: 33, DeckCard: data.DeckCard{BoardID: data.MainBoard, Quantity: 6}},
		data.Card{ID: 600, DeckCard: data.DeckCard{BoardID: data.MainBoard, Quantity: 34}},
		data.Card{ID: 701, DeckCard: data.DeckCard{BoardID: data.MainBoard, Quantity: 47}},
	}
	deck.SideCards = []data.Card{
		data.Card{ID: 51, DeckCard: data.DeckCard{BoardID: data.SideBoard, Quantity: 5}},
		data.Card{ID: 52, DeckCard: data.DeckCard{BoardID: data.SideBoard, Quantity: 27}},
	}
	createErr := deck.Persist()
	assert.Nil(t, createErr)
	fullDeck = deck
}

func Test_DeckUpdatedRead(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	log.Println("Test_DeckUpdatedRead")
	deck := &data.Deck{}
	deck.ID = fullDeck.ID
	readErr := deck.Read()
	assert.Nil(t, readErr)
	assert.Equal(t, deck.ID, fullDeck.ID)
	assert.Equal(t, deck.Name, fullDeck.Name)

	mainCardsLen := len(deck.MainCards)
	assert.Equal(t, mainCardsLen, len(fullDeck.MainCards))
	for k := 0; k < mainCardsLen; k++ {
		fullDeckMainCard := fullDeck.MainCards[k]
		deckMainCard := deck.MainCards[k]
		assert.Equal(t, deckMainCard.ID, fullDeckMainCard.ID)
		assert.Equal(t, deckMainCard.DeckCard.DeckID, fullDeck.ID)
		assert.Equal(t, deckMainCard.DeckCard.BoardID, data.MainBoard)
		assert.Equal(t, deckMainCard.DeckCard.Quantity, fullDeckMainCard.DeckCard.Quantity)
	}
	sideCardsLen := len(deck.SideCards)
	assert.Equal(t, sideCardsLen, len(fullDeck.SideCards))
	for k := 0; k < sideCardsLen; k++ {
		fullDeckSideCard := fullDeck.SideCards[k]
		deckSideCard := deck.SideCards[k]
		assert.Equal(t, deckSideCard.ID, fullDeckSideCard.ID)
		assert.Equal(t, deckSideCard.DeckCard.DeckID, fullDeck.ID)
		assert.Equal(t, deckSideCard.DeckCard.BoardID, data.SideBoard)
		assert.Equal(t, deckSideCard.DeckCard.Quantity, fullDeckSideCard.DeckCard.Quantity)
	}
}

func Test_DeckDelete(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	log.Println("Test_DeckDelete")
	deck := &data.Deck{}
	deck.ID = fullDeck.ID
	deleteErr := deck.Delete()
	assert.Nil(t, deleteErr)
	readErr := deck.Read()
	assert.NotNil(t, readErr)
}

//Deck Full
//Deck
