package api_test

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"
	//"errors"
	"net/http"
	"net/http/httptest"

	"farm.e-pedion.com/repo/fivecolors/config"
	"farm.e-pedion.com/repo/fivecolors/data"

	"farm.e-pedion.com/repo/fivecolors/api"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

var (
	setted      = false
	rec         *httptest.ResponseRecorder
	deckID      string
	inventoryID string
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
	rec = httptest.NewRecorder()
	return nil
}

func Test_GetCard(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	cardID := "1"
	req, err := http.NewRequest("GET", "http://mockrequest.com/card/"+cardID, bytes.NewBufferString(""))
	assert.Nil(t, err)

	getCardHandler := api.NewQueryCardHandler()
	getCardHandler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json; charset=utf-8", rec.Header().Get("Content-Type"))
	assert.Contains(t, rec.Body.String(), cardID)
}

func Test_QueryCardByColor(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	//"rx_cost", "nrx_cost", "rx_type", "nrx_type", "order"
	cardQuery := "rx_cost=Black|Red|White"
	req, err := http.NewRequest("GET", "http://mockrequest.com/card/?"+cardQuery, bytes.NewBufferString(""))
	assert.Nil(t, err)

	getCardHandler := api.NewQueryCardHandler()
	getCardHandler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json; charset=utf-8", rec.Header().Get("Content-Type"))
	assert.NotEqual(t, rec.Body.String(), "")
}

func Test_QueryCardByColorOrderByName(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	//"rx_cost", "nrx_cost", "rx_type", "nrx_type", "order"
	cardQuery := "rx_cost=Black|Red|White&order=name"
	req, err := http.NewRequest("GET", "http://mockrequest.com/card/?"+cardQuery, bytes.NewBufferString(""))
	assert.Nil(t, err)

	getCardHandler := api.NewQueryCardHandler()
	getCardHandler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json; charset=utf-8", rec.Header().Get("Content-Type"))
	assert.NotEqual(t, rec.Body.String(), "")
}

func Test_QueryCardWithoutParameters(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	req, err := http.NewRequest("GET", "http://mockrequest.com/card", bytes.NewBufferString(""))
	assert.Nil(t, err)

	getCardHandler := api.NewQueryCardHandler()
	getCardHandler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func Test_QueryExpansionWithoutParameters(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	req, err := http.NewRequest("GET", "http://mockrequest.com/expansion", bytes.NewBufferString(""))
	assert.Nil(t, err)

	getExpansionHandler := api.NewQueryExpansionHandler()
	getExpansionHandler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json; charset=utf-8", rec.Header().Get("Content-Type"))
	assert.NotEqual(t, rec.Body.String(), "")
}

func Test_PostToCreateInventory(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	inventoryJSON := `{
		"name": "Test_PostToCreateInventory",
		"cards": [
            {"id": 2764, "inventoryCard": {"quantity": 1} },
            {"id": 2884, "inventoryCard": {"quantity": 15} },
            {"id": 2791, "inventoryCard": {"quantity": 1} },
            {"id": 2863, "inventoryCard": {"quantity": 20} }
        ]
	}`
	req, err := http.NewRequest("POST", "http://mockrequest.com/inventory/", bytes.NewBufferString(inventoryJSON))
	assert.Nil(t, err)

	postInventoryHandler := api.NewInventoryHandler()
	postInventoryHandler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "application/json; charset=utf-8", rec.Header().Get("Content-Type"))
	body := rec.Body.String()
	assert.NotEqual(t, body, "")
	intInventoryID, convertIDErr := strconv.Atoi(body)
	assert.Nil(t, convertIDErr)
	assert.NotEqual(t, intInventoryID, 0)
	inventoryID = body
}

func Test_PostToUpdateInventory(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	inventoryJSON := fmt.Sprintf(`{
        "id": %v,
		"name": "Test_PostToUpdateInventory",
		"cards": [
            {"id": 1, "inventoryCard": {"quantity": 4} },
            {"id": 3, "inventoryCard": {"quantity": 5} },
            {"id": 5, "inventoryCard": {"quantity": 10} },
            {"id": 6, "inventoryCard": {"quantity": 2} }
        ]
	}`, inventoryID)
	req, err := http.NewRequest("POST", "http://mockrequest.com/inventory/", bytes.NewBufferString(inventoryJSON))
	assert.Nil(t, err)

	postInventoryHandler := api.NewInventoryHandler()
	postInventoryHandler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusAccepted, rec.Code)
	assert.Equal(t, "application/json; charset=utf-8", rec.Header().Get("Content-Type"))
	body := rec.Body.String()
	assert.Equal(t, body, "")
}

func Test_DeleteInventory(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}

	req, err := http.NewRequest("DELETE", "http://mockrequest.com/inventory/"+inventoryID, bytes.NewBufferString(""))
	assert.Nil(t, err)

	inventoryHandler := api.NewInventoryHandler()
	inventoryHandler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json; charset=utf-8", rec.Header().Get("Content-Type"))
	body := rec.Body.String()
	assert.Equal(t, body, "")
}

func Test_PostToCreateDeck(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	deckJSON := `{
		"name": "Test_PostToCreateDeck",
        "mainCards": [
            {"id": 2001, "deckCard": {"idBoard": 1, "quantity": 4} },
            {"id": 2002, "deckCard": {"idBoard": 1, "quantity": 4} },
            {"id": 2003, "deckCard": {"idBoard": 1, "quantity": 2} }
        ],
        "sideCards": [
            {"id": 30, "deckCard": {"idBoard": 2, "quantity": 13} },
            {"id": 50, "deckCard": {"idBoard": 2, "quantity": 42} },
            {"id": 180, "deckCard": {"idBoard": 2, "quantity": 27} }
        ]
	}`
	req, err := http.NewRequest("POST", "http://mockrequest.com/deck/", bytes.NewBufferString(deckJSON))
	assert.Nil(t, err)

	postDeckHandler := api.NewDeckHandler()
	postDeckHandler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "application/json; charset=utf-8", rec.Header().Get("Content-Type"))
	body := rec.Body.String()
	assert.NotEqual(t, body, "")
	intDeckID, convertIDErr := strconv.Atoi(body)
	assert.Nil(t, convertIDErr)
	assert.NotEqual(t, intDeckID, 0)
	deckID = body
}

func Test_GetDeck(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	req, err := http.NewRequest("GET", "http://mockrequest.com/deck/"+deckID, bytes.NewBufferString(""))
	assert.Nil(t, err)

	getDeckHandler := api.NewDeckHandler()
	getDeckHandler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json; charset=utf-8", rec.Header().Get("Content-Type"))
	assert.Contains(t, rec.Body.String(), deckID)
}

func Test_PostToUpdateDeck(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	deckJSON := fmt.Sprintf(`{
        "id": %v,
		"name": "Test_PostToUpdateDeck",
        "mainCards": [
            {"id": 2764, "deckCard": {"idBoard": 1, "quantity": 4} },
            {"id": 2884, "deckCard": {"idBoard": 1, "quantity": 4} },
            {"id": 2791, "deckCard": {"idBoard": 1, "quantity": 2} },
            {"id": 2863, "deckCard": {"idBoard": 1, "quantity": 1} }
        ],
        "sideCards": [
            {"id": 3, "deckCard": {"idBoard": 2, "quantity": 13} },
            {"id": 5, "deckCard": {"idBoard": 2, "quantity": 42} },
            {"id": 1980, "deckCard": {"idBoard": 2, "quantity": 27} }
        ]
	}`, deckID)

	req, err := http.NewRequest("POST", "http://mockrequest.com/deck/", bytes.NewBufferString(deckJSON))
	assert.Nil(t, err)

	postDeckHandler := api.NewDeckHandler()
	postDeckHandler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusAccepted, rec.Code)
	assert.Equal(t, "application/json; charset=utf-8", rec.Header().Get("Content-Type"))
	body := rec.Body.String()
	assert.Equal(t, body, "")
}

func Test_DeleteDeck(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}

	req, err := http.NewRequest("DELETE", "http://mockrequest.com/deck/"+deckID, bytes.NewBufferString(""))
	assert.Nil(t, err)

	postDeckHandler := api.NewDeckHandler()
	postDeckHandler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json; charset=utf-8", rec.Header().Get("Content-Type"))
	body := rec.Body.String()
	assert.Equal(t, body, "")
}
