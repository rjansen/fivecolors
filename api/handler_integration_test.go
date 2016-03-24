package api_test

import (
	"bytes"
	"strconv"
	"testing"
	//"errors"
	"net/http"
	"net/http/httptest"

	"farm.e-pedion.com/repo/fivecolors/config"
	"farm.e-pedion.com/repo/fivecolors/data"

	"farm.e-pedion.com/repo/fivecolors/api"
	"github.com/stretchr/testify/assert"
)

var (
	setted = false
	rec    *httptest.ResponseRecorder
	deckID string
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

func Test_PostInventory(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	inventoryJSON := `{
		"id": 1000,
		"cards": [
            {"id": 2764, "index": "86", "inventoryCard": {"quantity": 1} },
            {"id": 2884, "index": "90a", "inventoryCard": {"quantity": 15} },
            {"id": 2791, "index": "2", "inventoryCard": {"quantity": 1} },
            {"id": 2863, "index": "130", "inventoryCard": {"quantity": 20} }
        ]
	}`
	req, err := http.NewRequest("POST", "http://mockrequest.com/inventory/", bytes.NewBufferString(inventoryJSON))
	assert.Nil(t, err)

	postInventoryHandler := api.NewPostInventoryHandler()
	postInventoryHandler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusAccepted, rec.Code)
	assert.Equal(t, "application/json; charset=utf-8", rec.Header().Get("Content-Type"))
}

func Test_PostDeck(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	deckJSON := `{
		"name": "Test_PostDeck",
		"label": "Test_PostDeck - Handler Integration Test"
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
