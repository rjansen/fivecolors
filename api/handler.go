package api

import (
	"io"
	//"bytes"
	// "encoding/json"
	// "errors"
	// "fmt"
	"net/http"
	"path"
	"strconv"
	// "strings"

	// "github.com/rjansen/fivecolors/config"
	"github.com/rjansen/fivecolors/data"
	haki "github.com/rjansen/haki/http"
	// "github.com/rjansen/fivecolors/security"
	"github.com/rjansen/l"
	"github.com/rjansen/raizel"
	// "github.com/rjansen/avalon/identity"
	// "github.com/valyala/fasthttp"
)

//NewGetPlayerHandler creates a new GetPlayerHandler instance
// func NewGetPlayerHandler() http.Handler {
// 	return identity.NewHeaderAuthenticatedHandler(&GetPlayerHandler{})
// }

//GetPlayerHandler is the handler to get Players
// type GetPlayerHandler struct {
// 	identity.AuthenticatedHandler
// }

// func (h GetPlayerHandler) HandleRequest(ctx *fasthttp.RequestCtx) {

// }

// func (h *GetPlayerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	urlPathParameters := strings.Split(r.URL.Path, "/")
// 	queryParameters := r.URL.Query()
// 	l.Infof("GetPlayerHandler: URL[%q] PathParameters[%q] QueryParameters[%q]", r.URL.Path, urlPathParameters, queryParameters)

// 	session := h.GetSession()
// 	if strings.TrimSpace(session.Username) == "" {
// 		http.NotFound(w, r)
// 		return
// 	}
// 	player, err := data.GetPlayer(session.Username)
// 	if err != nil {
// 		l.Errorf("GetPlayerHandler.GetPlayerError: Session.ID[%v] Player.Username[%v] Error[%v]", session.ID, session.Username, err.Error())
// 		http.NotFound(w, r)
// 		//http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json; charset=utf-8")
// 	w.WriteHeader(http.StatusOK)
// 	jsonData, err := json.Marshal(player)
// 	bytesWritten, err := w.Write(jsonData)
// 	if err != nil {
// 		l.Errorf("GetPlayerHandler.WriteResponseError: Player[%v] Error[%v]", player, err)
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	} else {
// 		l.Debugf("GetPlayerHandler.JsonWritten: Player[%v] Bytes[%v]", player, bytesWritten)
// 	}

// }

//DeleteInventoryHandler is the handler to get Decks
// type DeleteInventoryHandler struct {
// 	security.InjectedPlayerHandler
// }

// func (h DeleteInventoryHandler) HandleRequest(ctx *fasthttp.RequestCtx) {

// }

// func (h *DeleteInventoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	urlPathParameters := strings.Split(r.URL.Path, "/")
// 	queryParameters := r.URL.Query()
// 	l.Infof("DeleteInventoryHandler: URL[%q] PathParameters[%q] QueryParameters[%q]", r.URL.Path, urlPathParameters, queryParameters)

// 	inventoryID := urlPathParameters[2]
// 	if inventoryID == "" {
// 		http.NotFound(w, r)
// 		return
// 	}
// 	inventory := &data.Inventory{}
// 	var convertIDErr error
// 	inventory.ID, convertIDErr = strconv.Atoi(inventoryID)
// 	if convertIDErr != nil {
// 		l.Errorf("DeleteInventoryHandler.DeckDeleteError: Inventory.ID[%v] Error[%v]", inventoryID, convertIDErr)
// 		http.NotFound(w, r)
// 		//http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	inventory.IDPlayer = h.GetPlayer().ID
// 	if err := inventory.Delete(); err != nil {
// 		l.Errorf("DeleteInventoryHandler.DeckDeleteError: Inventory.ID[%v] Error[%v]", inventoryID, err)
// 		http.NotFound(w, r)
// 		//http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	l.Debugf("DeleteInventoryHandler.DeletedDeck: Deck.ID[%v]", inventoryID)
// 	w.Header().Set("Content-Type", "application/json; charset=utf-8")
// 	w.WriteHeader(http.StatusOK)
// }

// func (h DeckHandler) Delete(player *data.Player, w http.ResponseWriter, r *http.Request) error {
// 	urlPathParameters := strings.Split(r.URL.Path, "/")
// 	queryParameters := r.URL.Query()
// 	l.Infof("DeleteDeckHandler: URL[%q] PathParameters[%q] QueryParameters[%q]", r.URL.Path, urlPathParameters, queryParameters)

// 	deckID := urlPathParameters[2]
// 	if deckID == "" {
// 		http.NotFound(w, r)
// 		return errors.New(http.StatusText(http.StatusNotFound))
// 	}
// 	deck := &data.Deck{}
// 	var convertIDErr error
// 	deck.ID, convertIDErr = strconv.Atoi(deckID)
// 	if convertIDErr != nil {
// 		l.Errorf("DeleteDeckHandler.DeckDeleteError: Deck.ID[%v] Error[%v]", deckID, convertIDErr)
// 		http.NotFound(w, r)
// 		//http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return errors.New(http.StatusText(http.StatusNotFound))
// 	}
// 	if err := deck.Delete(); err != nil {
// 		l.Errorf("DeleteDeckHandler.DeckDeleteError: Deck.ID[%v] Error[%v]", deckID, err)
// 		http.NotFound(w, r)
// 		//http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return errors.New(http.StatusText(http.StatusNotFound))
// 	}
// 	l.Debugf("DeleteDeckHandler.DeletedDeck: Deck.ID[%v]", deckID)
// 	w.Header().Set("Content-Type", "application/json; charset=utf-8")
// 	w.WriteHeader(http.StatusOK)
// 	return nil
// }

//NewAnonCardHandler creates a new unauthorized cardHandler instance
func NewAnonCardHandler() http.HandlerFunc {
	var cardHandler CardHandler
	return haki.Handler(haki.Log(haki.Error(cardHandler.ServeHTTP)))
}

type CardHandler struct{}

func (h CardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	basePath, lastPath := path.Split(r.URL.Path)
	l.Info("DeckHandler.ServeHTTP",
		l.String("MethodPath", r.Method),
		l.String("Path", r.URL.Path),
		l.String("BasePath", basePath),
		l.String("LastPath", lastPath),
	)
	if r.Method == "GET" {
		if lastPath == "" {
			return h.Query(w, r)
		}
		return h.Read(w, r)
	}
	return haki.Status(w, http.StatusMethodNotAllowed)
}

func (h CardHandler) Read(w http.ResponseWriter, r *http.Request) error {
	readParameter := path.Base(r.URL.Path)
	l.Info("CardHandler.Read",
		l.String("URI", r.URL.Path),
		l.String("parameter", readParameter),
	)
	var card data.Card
	var err error
	if id, atoirErr := strconv.Atoi(readParameter); atoirErr == nil {
		card.ID = id
		err = raizel.Execute(card.ReadByID)
	} else {
		card.Name = readParameter
		err = raizel.Execute(card.ReadByName)
	}
	if err != nil {
		if err == raizel.ErrNotFound {
			return haki.Status(w, http.StatusNotFound)
		}
		l.Error("CardHandler.ReadErr", l.String("parameter", readParameter), l.Err(err))
		return haki.Status(w, http.StatusBadRequest)
	}
	if err != nil {
		return haki.Status(w, http.StatusBadRequest)
	}
	return haki.JSON(w, http.StatusOK, card)
}

func (h CardHandler) Query(w http.ResponseWriter, r *http.Request) error {
	queryParameters := r.URL.Query()
	l.Info("CardHandler.Query",
		l.String("URI", r.URL.Path),
		l.Struct("QueryParameters", queryParameters),
	)
	if len(queryParameters) <= 0 {
		return haki.Status(w, http.StatusBadRequest)
	}

	var (
		card      data.Card
		cardQuery data.CardQuery
	)
	cardQuery.Hydrate = queryParameters.Get("hydrate")
	cardQuery.RegexName = queryParameters.Get("rx_name")
	cardQuery.RegexCost = queryParameters.Get("rx_cost")
	cardQuery.NotRegexCost = queryParameters.Get("nrx_cost")
	cardQuery.RegexType = queryParameters.Get("rx_type")
	cardQuery.NotRegexType = queryParameters.Get("nrx_type")
	cardQuery.IDExpansion = queryParameters.Get("e")
	cardQuery.Number = queryParameters.Get("n")
	cardQuery.InventoryQtd = queryParameters.Get("q")
	cardQuery.Order = queryParameters.Get("order")

	err := raizel.ExecuteWith(card.Query, &cardQuery)
	if err != nil {
		l.Error("CardHandler.QueryErr",
			l.Struct("QueryParameters", queryParameters),
			l.Err(err),
		)
		return haki.Err(w, err)
	}
	cardsSize := len(cardQuery.Result)
	l.Info("CardHandler.QueryResult",
		l.Int("Cards.Len", cardsSize),
		l.String("Hydrate", cardQuery.Hydrate),
	)
	return haki.JSON(w, http.StatusOK, cardQuery.Result)
}

func NewAnonTokenHandler() http.HandlerFunc {
	var tokenHandler TokenHandler
	return haki.Handler(haki.Log(haki.Error(tokenHandler.ServeHTTP)))
}

type TokenHandler struct{}

func (h TokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	basePath, lastPath := path.Split(r.URL.Path)
	l.Info("TokenHandler.ServeHTTP",
		l.String("MethodPath", r.Method),
		l.String("Path", r.URL.Path),
		l.String("BasePath", basePath),
		l.String("LastPath", lastPath),
	)
	if r.Method == "GET" {
		if lastPath == "" {
			return h.Query(w, r)
		}
		return h.Read(w, r)
	}
	return haki.Status(w, http.StatusMethodNotAllowed)
}

func (h TokenHandler) Read(w http.ResponseWriter, r *http.Request) error {
	readParameter := path.Base(r.URL.Path)
	l.Info("TokenHandler.Read",
		l.String("URI", r.URL.Path),
		l.String("parameter", readParameter),
	)
	var token data.Token
	var err error
	if id, atoirErr := strconv.Atoi(readParameter); atoirErr == nil {
		token.ID = id
		err = raizel.Execute(token.ReadByID)
	} else {
		token.Name = readParameter
		err = raizel.Execute(token.ReadByName)
	}
	if err != nil {
		if err == raizel.ErrNotFound {
			return haki.Status(w, http.StatusNotFound)
		}
		l.Error("TokenHandler.ReadErr", l.String("parameter", readParameter), l.Err(err))
		return haki.Status(w, http.StatusBadRequest)
	}
	if err != nil {
		return haki.Status(w, http.StatusBadRequest)
	}
	return haki.JSON(w, http.StatusOK, token)
}

func (h TokenHandler) Query(w http.ResponseWriter, r *http.Request) error {
	queryParameters := r.URL.Query()
	l.Info("TokenHandler.Query",
		l.String("URI", r.URL.Path),
		l.Struct("QueryParameters", queryParameters),
	)

	var (
		token      data.Token
		tokenQuery data.TokenQuery
	)
	tokenQuery.Hydrate = queryParameters.Get("hydrate")
	tokenQuery.RegexName = queryParameters.Get("rx_name")
	tokenQuery.RegexType = queryParameters.Get("rx_type")
	tokenQuery.NotRegexType = queryParameters.Get("nrx_type")
	tokenQuery.IDExpansion = queryParameters.Get("e")
	tokenQuery.Order = queryParameters.Get("order")

	err := raizel.ExecuteWith(token.Query, &tokenQuery)
	if err != nil {
		l.Error("TokenHandler.QueryErr",
			l.Struct("QueryParameters", queryParameters),
			l.Err(err),
		)
		return haki.Err(w, err)
	}
	tokensSize := len(tokenQuery.Result)
	l.Info("TokenHandler.QueryResult",
		l.Int("Tokens.Len", tokensSize),
		l.String("Hydrate", tokenQuery.Hydrate),
	)
	return haki.JSON(w, http.StatusOK, tokenQuery.Result)
}

//NewAnonDeckHandler creates a new unauthorized deckHandler instance
func NewAnonDeckHandler() http.HandlerFunc {
	var deckHandler DeckHandler
	return haki.Handler(haki.Log(haki.Error(deckHandler.ServeHTTP)))
}

type DeckHandler struct{}

func (h DeckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	basePath, lastPath := path.Split(r.URL.Path)
	l.Info("DeckHandler.ServeHTTP",
		l.String("MethodPath", r.Method),
		l.String("Path", r.URL.Path),
		l.String("BasePath", basePath),
		l.String("LastPath", lastPath),
	)
	if r.Method == "GET" {
		if lastPath == "" {
			return h.Query(w, r)
		}
		return h.Read(w, r)
	} else if r.Method == "POST" {
		return h.Persist(w, r)
	}
	return haki.Status(w, http.StatusMethodNotAllowed)
}

func (h DeckHandler) Persist(w http.ResponseWriter, r *http.Request) error {
	queryParameters := r.URL.Query()
	l.Info("DeckHandler.Persist",
		l.String("URI", r.URL.Path),
		l.Struct("QueryParameters", queryParameters),
	)
	var err error
	var deck data.Deck
	if err = haki.ReadJSON(r, &deck); err != nil {
		return haki.Err(w, err)
	}
	isCreateRequest := deck.ID == 0
	if err = raizel.Execute(deck.Persist); err != nil {
		return haki.Err(w, err)
	}
	if isCreateRequest {
		if err = haki.Status(w, http.StatusCreated); err != nil {
			return err
		}
		io.WriteString(w, strconv.Itoa(deck.ID))
	}
	return haki.Status(w, http.StatusAccepted)
}

func (h DeckHandler) Read(w http.ResponseWriter, r *http.Request) error {
	readParameter := path.Base(r.URL.Path)
	l.Info("DeckHandler.Read",
		l.String("URI", r.URL.Path),
		l.Struct("ReadParameter", readParameter),
	)

	var deck data.Deck
	var err error
	if id, atoirErr := strconv.Atoi(readParameter); atoirErr == nil {
		deck.ID = id
		err = raizel.Execute(deck.ReadByID)
	} else {
		deck.Name = readParameter
		err = raizel.Execute(deck.ReadByName)
	}
	if err != nil {
		return haki.Status(w, http.StatusNotFound)
	}
	return haki.JSON(w, http.StatusOK, deck)
}

func (h DeckHandler) Query(w http.ResponseWriter, r *http.Request) error {
	queryParameters := r.URL.Query()
	l.Infof("DeckHandler.Query",
		l.String("URI", r.URL.Path),
		l.Struct("QueryParameters", queryParameters),
	)
	var (
		deck      data.Deck
		deckQuery data.DeckQuery
	)
	deckQuery.RegexName = queryParameters.Get("rx_name")
	err := raizel.ExecuteWith(deck.Query, &deckQuery)
	if err != nil {
		return haki.Err(w, err)
	}
	return haki.JSON(w, http.StatusOK, deckQuery.Result)
}

//NewAnonExpansionHandler creates a new unauthorized expansionHandler instance
func NewAnonExpansionHandler() http.HandlerFunc {
	var expansionHandler ExpansionHandler
	return haki.Handler(haki.Log(haki.Error(expansionHandler.ServeHTTP)))
}

type ExpansionHandler struct{}

func (h ExpansionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	basePath, lastPath := path.Split(r.URL.Path)
	l.Info("ExpansionHandler.ServeHTTP",
		l.String("MethodPath", r.Method),
		l.String("Path", r.URL.Path),
		l.String("BasePath", basePath),
		l.String("LastPath", lastPath),
	)
	if r.Method == "GET" {
		if lastPath == "" {
			return h.Query(w, r)
		}
		return h.Read(w, r)
	}
	return haki.Status(w, http.StatusMethodNotAllowed)
}

func (h ExpansionHandler) Read(w http.ResponseWriter, r *http.Request) error {
	readParameter := path.Base(r.URL.Path)
	l.Info("ExpansionHandler.Read",
		l.String("URI", r.URL.Path),
		l.Struct("ReadParameter", readParameter),
	)
	var (
		expansion data.Expansion
		err       error
	)
	if id, atoirErr := strconv.Atoi(readParameter); atoirErr == nil {
		expansion.ID = id
		err = raizel.Execute(expansion.ReadByID)
	} else {
		expansion.Name = readParameter
		err = raizel.Execute(expansion.ReadByName)
	}
	if err != nil {
		return haki.Status(w, http.StatusNotFound)
	}
	return haki.JSON(w, http.StatusOK, expansion)
}

func (h ExpansionHandler) Query(w http.ResponseWriter, r *http.Request) error {
	queryParameters := r.URL.Query()
	l.Info("ExpansionHandler.Query",
		l.String("URI", r.URL.Path),
		l.Struct("QueryParameters", queryParameters),
	)
	var (
		expansion    data.Expansion
		queryBuilder data.ExpansionQuery
	)
	queryBuilder.Hydrate = queryParameters.Get("hydrate")
	queryBuilder.RegexName = queryParameters.Get("rx_name")
	queryBuilder.Order = queryParameters.Get("order")
	err := raizel.ExecuteWith(expansion.Query, &queryBuilder)
	if err != nil {
		return haki.Err(w, err)
	}
	expansionSize := len(queryBuilder.Result)
	l.Debug("ExpansionHandler.QueryResult",
		l.Int("Expansions.Len", expansionSize),
		l.String("Hydrate", queryBuilder.Hydrate),
	)
	return haki.JSON(w, http.StatusOK, queryBuilder.Result)
}

//NewAnonInventoryHandler creates a new DeckHandler instance
func NewAnonInventoryHandler() http.HandlerFunc {
	var inventoryHandler InventoryHandler
	return haki.Handler(haki.Log(haki.Error(inventoryHandler.ServeHTTP)))
}

type InventoryHandler struct{}

func (h InventoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	basePath, lastPath := path.Split(r.URL.Path)
	l.Info("InventoryHandler.ServeHTTP",
		l.String("MethodPath", r.Method),
		l.String("Path", r.URL.Path),
		l.String("BasePath", basePath),
		l.String("LastPath", lastPath),
	)
	if r.Method == "POST" || r.Method == "PUT" {
		return h.Persist(w, r)
	}
	return haki.Status(w, http.StatusMethodNotAllowed)
}

func (h InventoryHandler) Persist(w http.ResponseWriter, r *http.Request) error {
	readParameter := path.Base(r.URL.Path)
	l.Info("InventoryHandler.Persist",
		l.String("URI", r.URL.Path),
		l.String("ReadParameters", readParameter),
	)
	var inventory data.Inventory
	if err := haki.ReadJSON(r, &inventory); err != nil {
		return haki.Err(w, err)
	}
	//Fixed to zero for anonymous inventory
	inventory.ID = 0
	inventory.IDPlayer = 0
	if err := raizel.Execute(inventory.Persist); err != nil {
		return haki.Err(w, err)
	}
	return haki.Status(w, http.StatusAccepted)
}
