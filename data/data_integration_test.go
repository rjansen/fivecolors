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
	"github.com/stretchr/testify/assert"
)

var (
	setted      = false
	deckID      int
	inventoryID int
)

func setup() error {
	data.Setup(&config.DBConfig{
		Driver:   "mysql",
		URL:      "tcp(127.0.0.1:3306)/fivecolors",
		Username: "fivecolors",
		Password: "fivecolors",
	})
	setted = true
	return nil
}

func before() error {
	if !setted {
		if err := setup(); err != nil {
			return err
		}
	}
	return nil
}

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
    }
}

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
	assert.NotEqual(t, expansion.IDAsset, "")
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
        assert.NotEqual(t, expansion.IDAsset, "")
    }
}

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

func Test_DeckCreate(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	log.Println("Test_DeckCreate")
	deck := &data.Deck{}
	deck.Name = "Test_DeckCreate"
	deck.Label = "Test_DeckCreate - Tests the Deck Create feature"
	createErr := deck.Create()
	assert.Nil(t, createErr)
	assert.NotEqual(t, deck.ID, 0)
	deckID = deck.ID
}

func Test_DeckRead(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	log.Println("Test_DeckRead")
	deck := &data.Deck{}
	deck.ID = deckID
	readErr := deck.Read()
	assert.Nil(t, readErr)
	assert.Equal(t, deck.ID, deckID)
	assert.NotEqual(t, deck.Name, "")
	assert.NotEqual(t, deck.Label, "")
}

func Test_DeckDelete(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	log.Println("Test_DeckDelete")
	deck := &data.Deck{}
	deck.ID = deckID
	deleteErr := deck.Delete()
	assert.Nil(t, deleteErr)
	readErr := deck.Read()
	assert.NotNil(t, readErr)
}

func Test_InventoryCreate(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	log.Println("Test_InventoryCreate")
	inventory := &data.Inventory{}
	inventory.Name = "Test_InventoryCreate"
	inventory.Label = "Test_InventoryCreate - Tests the Inventory Create feature"
	createErr := inventory.Create()
	assert.Nil(t, createErr)
	assert.NotEqual(t, inventory.ID, 0)
	inventoryID = inventory.ID
}

func Test_InventoryRead(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	log.Println("Test_InventoryRead")
	inventory := &data.Inventory{}
	inventory.ID = inventoryID
	readErr := inventory.Read()
	assert.Nil(t, readErr)
	assert.Equal(t, inventory.ID, inventoryID)
	assert.NotEqual(t, inventory.Name, "")
	assert.NotEqual(t, inventory.Label, "")
}

func Test_InventoryDelete(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	log.Println("Test_InventoryDelete")
	inventory := &data.Inventory{}
	inventory.ID = inventoryID
	deleteErr := inventory.Delete()
	assert.Nil(t, deleteErr)
	readErr := inventory.Read()
	assert.NotNil(t, readErr)
}
