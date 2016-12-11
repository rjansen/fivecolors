package api

import (
	"io"
	//"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"path"
	"strconv"
	"strings"

	haki "farm.e-pedion.com/repo/context/http"
	"farm.e-pedion.com/repo/fivecolors/config"
	"farm.e-pedion.com/repo/fivecolors/data"
	"farm.e-pedion.com/repo/fivecolors/security"
	"farm.e-pedion.com/repo/logger"
	l "farm.e-pedion.com/repo/logger"
	raizel "farm.e-pedion.com/repo/persistence"
	"farm.e-pedion.com/repo/security/identity"
	"github.com/valyala/fasthttp"
)

const (
	hydrateSmall = "small"
	hydrateFull  = "full"
)

var (
	log logger.Logger
)

func Setup(config config.Configuration) error {
	log = logger.Get()
	return nil
}

//NewQueryCardHandler creates a new QueryCardHandler instance
func NewQueryCardHandler() http.Handler {
	return identity.NewHeaderAuthenticatedHandler(
		security.NewInjectPlayerHandler(&QueryCardHandler{}))
}

//QueryCardHandler is the handler to get and query Cards
type QueryCardHandler struct {
	security.InjectedPlayerHandler
}

func (h QueryCardHandler) HandleRequest(ctx *fasthttp.RequestCtx) {

}

func (h *QueryCardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	urlPathParameters := strings.Split(r.URL.Path, "/")
	queryParameters := r.URL.Query()
	log.Infof("api.QueryCardHandler: URL[%q] PathParameters[%q] QueryParameters[%q]", r.URL.Path, urlPathParameters, queryParameters)

	var cardID string
	if len(urlPathParameters) >= 3 {
		cardID = urlPathParameters[2]
	} else {
		cardID = ""
	}
	hydrate := queryParameters.Get("hydrate")
	if hydrate == "" {
		hydrate = "small"
	}
	if cardID != "" {
		card := &data.Card{}
		var convertIDErr error
		card.ID, convertIDErr = strconv.Atoi(cardID)
		if convertIDErr != nil {
			log.Errorf("QueryCardHandler.CardReadError: Card.ID[%v] Error[%v]", cardID, convertIDErr)
			http.NotFound(w, r)
			//http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		card.InventoryCard = data.InventoryCard{IDInventory: h.GetPlayer().IDInventory}
		if err := card.Read(); err != nil {
			log.Errorf("QueryCardHandler.CardReadError: Card.ID[%v] Error[%v]", cardID, err)
			http.NotFound(w, r)
			//http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Debugf("QueryCardHandler.GotCard: Card.ID[%v] Hydrate[%s]", cardID, hydrate)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		jsonData, err := json.Marshal(card)
		bytesWritten, err := w.Write(jsonData)
		if err != nil {
			log.Errorf("QueryCardHandler.WriteResponseError: Card.ID[%v] Error[%v]", cardID, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			log.Debugf("QueryCardHandler.JsonWritten: Card.ID[%v] Bytes[%v]", cardID, bytesWritten)
		}
	} else {
		if len(queryParameters) <= 0 {
			//w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		regexName := queryParameters.Get("rx_name")
		regexCost := queryParameters.Get("rx_cost")
		notRegexCost := queryParameters.Get("nrx_cost")
		regexType := queryParameters.Get("rx_type")
		notRegexType := queryParameters.Get("nrx_type")
		expansionParam := queryParameters.Get("e")
		numberParam := queryParameters.Get("n")
		stockQtdParam := queryParameters.Get("q")
		order := queryParameters.Get("order")

		queryRestrinctions := make(map[string]interface{})
		if regexName != "" {
			queryRestrinctions["c.name regexp ?"] = regexName
		}
		if regexCost != "" {
			queryRestrinctions["c.manacost_label regexp ?"] = regexCost
		}
		if notRegexCost != "" {
			queryRestrinctions["c.manacost_label not regexp ?"] = notRegexCost
		}
		if regexType != "" {
			queryRestrinctions["c.type_label regexp ?"] = regexType
		}
		if notRegexType != "" {
			queryRestrinctions["c.type_label not regexp ?"] = notRegexType
		}
		if expansionParam != "" {
			if IDExpansion, convertErr := strconv.Atoi(expansionParam); convertErr == nil {
				queryRestrinctions["c.id_expansion = ?"] = IDExpansion
			} else {
				log.Errorf("QueryCardHandler.ExpansionParameterError: Parameter[%v] Error[%v]", expansionParam, convertErr)
			}
		}
		if numberParam != "" {
			queryRestrinctions["c.multiverse_number = ?"] = numberParam
		}
		if stockQtdParam != "" {
			queryRestrinctions["coalesce(i.quantity, 0) >= ?"] = stockQtdParam
		}
		card := &data.Card{}
		card.InventoryCard = data.InventoryCard{IDInventory: h.GetPlayer().IDInventory}
		cards, err := card.Query(queryRestrinctions, order)
		if err != nil {
			log.Errorf("QueryCardHandler.CardQueryError: QueryRestrictions[%v] Error[%v]", queryRestrinctions, err)
			//http.NotFound(w, r)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cardsSize := len(cards)
		log.Debugf("QueryCardHandler.QueryCards: Cards.Len[%v] Hydrate[%s]", cardsSize, hydrate)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		jsonData, err := json.Marshal(cards)
		bytesWritten, err := w.Write(jsonData)
		if err != nil {
			log.Errorf("QueryCardHandler.WriteResponseError: Cards.Len[%v] Error[%v]", cardsSize, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			log.Debugf("QueryCardHandler.JsonWritten: Cards.Len[%v] Bytes[%v]", cardsSize, bytesWritten)
		}
	}
}

//NewQueryExpansionHandler creates a new QueryExpansionHandler instance
func NewQueryExpansionHandler() http.Handler {
	return &QueryExpansionHandler{}
	//return identity.NewCookieAuthenticatedHandler(&GetExpansionHandler{})
}

//QueryExpansionHandler is the handler to get and query Expansions
type QueryExpansionHandler struct {
	//    session *identity.Session
}

//func (h *QueryExpansionHandler) SetSession(session *identity.Session) {
//    h.session = session
//}

func (h QueryExpansionHandler) HandleRequest(ctx *fasthttp.RequestCtx) {

}

func (h *QueryExpansionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	urlPathParameters := strings.Split(r.URL.Path, "/")
	queryParameters := r.URL.Query()
	log.Infof("GetHandler: URL[%q] PathParameters[%q] QueryParameters[%q]", r.URL.Path, urlPathParameters, queryParameters)

	var expansionID string
	if len(urlPathParameters) >= 3 {
		expansionID = urlPathParameters[2]
	} else {
		expansionID = ""
	}
	hydrate := queryParameters.Get("hydrate")
	if hydrate == "" {
		hydrate = "small"
	}
	if expansionID != "" {
		expansion := &data.Expansion{}
		var convertIDErr error
		expansion.ID, convertIDErr = strconv.Atoi(expansionID)
		if convertIDErr != nil {
			log.Errorf("QueryExpansionHandler.ExpansionReadError: Expansion.ID[%v] Error[%v]", expansionID, convertIDErr)
			http.NotFound(w, r)
			//http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := expansion.Read(); err != nil {
			log.Errorf("QueryExpansionHandler.CardReadError: Expansion.ID[%v] Error[%v]", expansionID, err)
			http.NotFound(w, r)
			//http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Debugf("QueryExpansionkHandler.GotExpansion: Expansion.ID[%v] Hydrate[%s]", expansionID, hydrate)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		jsonData, err := json.Marshal(expansion)
		bytesWritten, err := w.Write(jsonData)
		if err != nil {
			log.Errorf("QueryExpansionHandler.WriteResponseError: Expansion.ID[%v] Error[%v]", expansionID, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			log.Debugf("QueryExpansionHandler.JsonWritten: Expansion.ID[%v] Bytes[%v]", expansionID, bytesWritten)
		}
	} else {
		regexName := queryParameters.Get("rx_name")
		order := queryParameters.Get("order")

		queryRestrinctions := make(map[string]interface{})
		if regexName != "" {
			queryRestrinctions["e.name regexp ?"] = regexName
		}
		if order != "" {
			order = "c.id"
		}
		expansion := &data.Expansion{}
		expansions, err := expansion.Query(queryRestrinctions, order)
		if err != nil {
			log.Errorf("QueryExpansionHandler.ExpansionQueryError: QueryRestrictions[%v] Error[%v]", queryRestrinctions, err)
			//http.NotFound(w, r)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		expansionSize := len(expansions)
		log.Debugf("QueryExpansionHandler.QueryExpansions: Expansions.Len[%v] Hydrate[%s]", expansionSize, hydrate)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		jsonData, err := json.Marshal(expansions)
		bytesWritten, err := w.Write(jsonData)
		if err != nil {
			log.Errorf("QueryExpansionHandler.WriteResponseError: Expansions.Len[%v] Error[%v]", expansionSize, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			log.Debugf("QueryExpansionHandler.JsonWritten: Expansions.Len[%v] Bytes[%v]", expansionSize, bytesWritten)
		}
	}
}

//NewGetAssetHandler creates a new GetAssetHandler instance
func NewGetAssetHandler() http.Handler {
	return &GetAssetHandler{}
}

//GetAssetHandler is the handler to get Assets
type GetAssetHandler struct {
	//    session *identity.Session
}

//func (h *QueryExpansionHandler) SetSession(session *identity.Session) {
//    h.session = session
//}

func (h GetAssetHandler) HandleRequest(ctx *fasthttp.RequestCtx) {

}

func (h *GetAssetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	urlPathParameters := strings.Split(r.URL.Path, "/")
	queryParameters := r.URL.Query()
	log.Infof("GetAssetHandler: URL[%q] PathParameters[%q] QueryParameters[%q]", r.URL.Path, urlPathParameters, queryParameters)

	assetID := urlPathParameters[2]
	if assetID == "" {
		http.NotFound(w, r)
		return
	}
	asset := &data.Asset{}
	var convertIDErr error
	asset.ID, convertIDErr = strconv.Atoi(assetID)
	if convertIDErr != nil {
		log.Errorf("GetAssetHandler.AssetReadError: Asset.ID[%v] Error[%v]", assetID, convertIDErr)
		//http.NotFound(w, r)
		http.Error(w, convertIDErr.Error(), http.StatusInternalServerError)
		return
	}
	if err := asset.Read(); err != nil {
		log.Errorf("GetAssetkHandler.AssetReadError: Asset.ID[%v] Error[%v]", assetID, err)
		http.NotFound(w, r)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Debugf("GetAssetHandler.GotDeck: Asset.ID[%v]", assetID)
	w.Header().Set("Content-Type", "image/jpeg")
	w.WriteHeader(http.StatusOK)
	bytesWritten, err := w.Write(asset.BinaryData)
	if err != nil {
		log.Errorf("GetAssetHandler.WriteResponseError: Asset.ID[%v] Error[%v]", assetID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		log.Debugf("GetAssetHandler.BytesWritten: Asset.ID[%v] Bytes[%v]", assetID, bytesWritten)
	}
}

//NewGetPlayerHandler creates a new GetPlayerHandler instance
func NewGetPlayerHandler() http.Handler {
	return identity.NewHeaderAuthenticatedHandler(&GetPlayerHandler{})
}

//GetPlayerHandler is the handler to get Players
type GetPlayerHandler struct {
	identity.AuthenticatedHandler
}

func (h GetPlayerHandler) HandleRequest(ctx *fasthttp.RequestCtx) {

}

func (h *GetPlayerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	urlPathParameters := strings.Split(r.URL.Path, "/")
	queryParameters := r.URL.Query()
	log.Infof("GetPlayerHandler: URL[%q] PathParameters[%q] QueryParameters[%q]", r.URL.Path, urlPathParameters, queryParameters)

	session := h.GetSession()
	if strings.TrimSpace(session.Username) == "" {
		http.NotFound(w, r)
		return
	}
	player, err := data.GetPlayer(session.Username)
	if err != nil {
		log.Errorf("GetPlayerHandler.GetPlayerError: Session.ID[%v] Player.Username[%v] Error[%v]", session.ID, session.Username, err.Error())
		http.NotFound(w, r)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	jsonData, err := json.Marshal(player)
	bytesWritten, err := w.Write(jsonData)
	if err != nil {
		log.Errorf("GetPlayerHandler.WriteResponseError: Player[%v] Error[%v]", player, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		log.Debugf("GetPlayerHandler.JsonWritten: Player[%v] Bytes[%v]", player, bytesWritten)
	}

}

//NewInventoryHandler creates a new DeckHandler instance
func NewInventoryHandler() http.Handler {
	return identity.NewHeaderAuthenticatedHandler(
		security.NewInjectPlayerHandler(&InventoryHandler{
			postHandler:   &PostInventoryHandler{},
			deleteHandler: &DeleteInventoryHandler{},
		}))
}

//InventoryHandler is the handler to get and post Decks
type InventoryHandler struct {
	security.InjectedPlayerHandler
	postHandler   security.PlayerInjectableHandler
	deleteHandler security.PlayerInjectableHandler
}

func (h InventoryHandler) HandleRequest(ctx *fasthttp.RequestCtx) {

}

func (h *InventoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session := h.GetSession()
	player := h.GetPlayer()
	if r.Method == "POST" {
		h.postHandler.SetSession(session)
		h.postHandler.SetPlayer(player)
		h.postHandler.ServeHTTP(w, r)
	} else if r.Method == "DELETE" {
		h.deleteHandler.SetSession(session)
		h.deleteHandler.SetPlayer(player)
		h.deleteHandler.ServeHTTP(w, r)
	} else {
		http.Error(w, "InventoryHandler.MethodNotAllowed: Method="+r.Method, http.StatusMethodNotAllowed)
	}
}

//PostInventoryHandler is the handler to updates Inventory
type PostInventoryHandler struct {
	security.InjectedPlayerHandler
}

func (h PostInventoryHandler) HandleRequest(ctx *fasthttp.RequestCtx) {

}

func (h *PostInventoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	urlPathParameters := strings.Split(r.URL.Path, "/")
	queryParameters := r.URL.Query()
	log.Infof("PostInvetoryHandler: URL[%q] PathParameters[%q] QueryParameters[%q]", r.URL.Path, urlPathParameters, queryParameters)

	inventory := data.Inventory{}
	if err := inventory.Unmarshal(r.Body); err != nil {
		log.Errorf("PostInvetoryHandler.UnreadableBodyError: Error[%v]", err)
		http.Error(w, "PostInvetoryHandler.UnreadableBodyError", http.StatusBadRequest)
		return
	}
	inventory.ID = h.GetPlayer().IDInventory
	inventory.IDPlayer = h.GetPlayer().ID
	isCreateRequest := inventory.ID == 0
	if err := inventory.Persist(); err != nil {
		log.Errorf("PostInvetoryHandler.CreateDeckError: Error[%v]", err)
		http.Error(w, "PostInvetoryHandler.CreateDeckError", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if isCreateRequest {
		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, strconv.Itoa(inventory.ID))
	} else {
		w.WriteHeader(http.StatusAccepted)
	}
}

//DeleteInventoryHandler is the handler to get Decks
type DeleteInventoryHandler struct {
	security.InjectedPlayerHandler
}

func (h DeleteInventoryHandler) HandleRequest(ctx *fasthttp.RequestCtx) {

}

func (h *DeleteInventoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	urlPathParameters := strings.Split(r.URL.Path, "/")
	queryParameters := r.URL.Query()
	log.Infof("DeleteInventoryHandler: URL[%q] PathParameters[%q] QueryParameters[%q]", r.URL.Path, urlPathParameters, queryParameters)

	inventoryID := urlPathParameters[2]
	if inventoryID == "" {
		http.NotFound(w, r)
		return
	}
	inventory := &data.Inventory{}
	var convertIDErr error
	inventory.ID, convertIDErr = strconv.Atoi(inventoryID)
	if convertIDErr != nil {
		log.Errorf("DeleteInventoryHandler.DeckDeleteError: Inventory.ID[%v] Error[%v]", inventoryID, convertIDErr)
		http.NotFound(w, r)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	inventory.IDPlayer = h.GetPlayer().ID
	if err := inventory.Delete(); err != nil {
		log.Errorf("DeleteInventoryHandler.DeckDeleteError: Inventory.ID[%v] Error[%v]", inventoryID, err)
		http.NotFound(w, r)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Debugf("DeleteInventoryHandler.DeletedDeck: Deck.ID[%v]", inventoryID)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

//NewDeckHandler creates a new DeckHandler instance
func NewDeckHandler() http.HandlerFunc {
	var deckHandler DeckHandler
	return identity.Authorize(deckHandler.ServeHTTP)
}

//DeckHandler is the handler to get and post Decks
type DeckHandler struct{}

func (h DeckHandler) ServeHTTP(session *identity.Session, w http.ResponseWriter, r *http.Request) error {
	player := new(data.Player)
	if err := player.FillFromSession(session); err != nil {
		return err
	}
	if r.Method == "GET" {
		return h.Query(player, w, r)
	} else if r.Method == "POST" {
		return h.Persist(player, w, r)
	} else if r.Method == "DELETE" {
		return h.Delete(player, w, r)
	}
	http.Error(w, "DeckHandler.MethodNotAllowed: Method="+r.Method, http.StatusMethodNotAllowed)
	return errors.New(http.StatusText(http.StatusMethodNotAllowed))
}

func (h DeckHandler) Persist(player *data.Player, w http.ResponseWriter, r *http.Request) error {
	urlPathParameters := strings.Split(r.URL.Path, "/")
	queryParameters := r.URL.Query()
	log.Infof("PostDeckdHandler: URL[%q] PathParameters[%q] QueryParameters[%q]", r.URL.Path, urlPathParameters, queryParameters)

	deck := data.Deck{}
	if err := deck.Unmarshal(r.Body); err != nil {
		log.Errorf("PostDeckHandler.UnreadableBodyError: Error[%v]", err)
		http.Error(w, "PostDeckHandler.UnreadableBodyError", http.StatusBadRequest)
		return err
	}
	deck.IDPlayer = player.ID
	isCreateRequest := deck.ID == 0
	if err := deck.Persist(); err != nil {
		log.Errorf("PostDeckHandler.CreateDeckError: Error[%v]", err)
		http.Error(w, "PostDeckHandler.CreateDeckError", http.StatusBadRequest)
		return err
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if isCreateRequest {
		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, strconv.Itoa(deck.ID))
	} else {
		w.WriteHeader(http.StatusAccepted)
	}
	return nil
}

func (h DeckHandler) Query(player *data.Player, w http.ResponseWriter, r *http.Request) error {
	urlPathParameters := strings.Split(r.URL.Path, "/")
	queryParameters := r.URL.Query()
	log.Infof("QueryDeckHandler: URL[%q] PathParameters[%q] QueryParameters[%q]", r.URL.Path, urlPathParameters, queryParameters)

	var deckID string
	if len(urlPathParameters) >= 3 {
		deckID = urlPathParameters[2]
	} else {
		deckID = ""
	}
	hydrate := queryParameters.Get("hydrate")
	if hydrate == "" {
		hydrate = "small"
	}
	if deckID != "" {
		deck := &data.Deck{
			IDPlayer:    player.ID,
			IDInventory: player.IDInventory,
		}
		var convertIDErr error
		deck.ID, convertIDErr = strconv.Atoi(deckID)
		if convertIDErr != nil {
			log.Errorf("QueryDeckHandler.DeckReadError: Deck.ID[%v] Error[%v]", deckID, convertIDErr)
			http.NotFound(w, r)
			//http.Error(w, err.Error(), http.StatusInternalServerError)
			return convertIDErr
		}
		if err := deck.Read(); err != nil {
			log.Errorf("QueryDeckHandler.DeckReadError: Deck.ID[%v] Error[%v]", deckID, err)
			http.NotFound(w, r)
			//http.Error(w, err.Error(), http.StatusInternalServerError)
			return err
		}
		log.Debugf("QueryDeckHandler.GotDeck: Deck.ID[%v] Hydrate[%s]", deckID, hydrate)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		jsonData, err := json.Marshal(deck)
		bytesWritten, err := w.Write(jsonData)
		if err != nil {
			log.Errorf("QueryDeckHandler.WriteResponseError: Deck.ID[%v] Error[%v]", deckID, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			log.Debugf("QueryDeckHandler.JsonWritten: Deck.ID[%v] Bytes[%v]", deckID, bytesWritten)
		}
	} else {
		regexName := queryParameters.Get("rx_name")
		order := queryParameters.Get("order")

		queryRestrinctions := make(map[string]interface{})
		if regexName != "" {
			queryRestrinctions["d.name regexp ?"] = regexName
		}
		if order != "" {
			order = "d.name"
		}
		deck := &data.Deck{
			IDPlayer:    player.ID,
			IDInventory: player.IDInventory,
		}
		decks, err := deck.Query(queryRestrinctions, order)
		if err != nil {
			log.Errorf("QueryDeckHandler.DeckReadError: Deck.ID[%v] Error[%v]", deckID, err)
			http.NotFound(w, r)
			//http.Error(w, err.Error(), http.StatusInternalServerError)
			return err
		}
		decksSize := len(decks)
		log.Debugf("QueryDeckHandler.QueryDecks: Decks.Len[%v] Hydrate[%s]", decksSize, hydrate)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		jsonData, err := json.Marshal(decks)
		bytesWritten, err := w.Write(jsonData)
		if err != nil {
			log.Errorf("QueryDeckHandler.WriteResponseError: Deck.ID[%v] Error[%v]", deckID, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			log.Debugf("QueryDeckHandler.JsonWritten: Deck.ID[%v] Bytes[%v]", deckID, bytesWritten)
		}
	}
	return nil
}

func (h DeckHandler) Delete(player *data.Player, w http.ResponseWriter, r *http.Request) error {
	urlPathParameters := strings.Split(r.URL.Path, "/")
	queryParameters := r.URL.Query()
	log.Infof("DeleteDeckHandler: URL[%q] PathParameters[%q] QueryParameters[%q]", r.URL.Path, urlPathParameters, queryParameters)

	deckID := urlPathParameters[2]
	if deckID == "" {
		http.NotFound(w, r)
		return errors.New(http.StatusText(http.StatusNotFound))
	}
	deck := &data.Deck{}
	var convertIDErr error
	deck.ID, convertIDErr = strconv.Atoi(deckID)
	if convertIDErr != nil {
		log.Errorf("DeleteDeckHandler.DeckDeleteError: Deck.ID[%v] Error[%v]", deckID, convertIDErr)
		http.NotFound(w, r)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return errors.New(http.StatusText(http.StatusNotFound))
	}
	if err := deck.Delete(); err != nil {
		log.Errorf("DeleteDeckHandler.DeckDeleteError: Deck.ID[%v] Error[%v]", deckID, err)
		http.NotFound(w, r)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return errors.New(http.StatusText(http.StatusNotFound))
	}
	log.Debugf("DeleteDeckHandler.DeletedDeck: Deck.ID[%v]", deckID)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return nil
}

//NewAnonDeckHandler creates a new unauthorized deckHandler instance
func NewAnonDeckHandler() http.HandlerFunc {
	var deckHandler AnonDeckHandler
	return haki.Handler(haki.Log(haki.Error(deckHandler.ServeHTTP)))
}

//AnonDeckHandler is the unsecure handler to get and post Decks
type AnonDeckHandler struct{}

func (h AnonDeckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		isQuery := path.Base(path.Dir(r.URL.Path)) == "query"
		if isQuery {
			return h.QueryByName(w, r)
		}
		return h.Read(w, r)
	} else if r.Method == "POST" {
		return h.Persist(w, r)
	}
	return haki.Status(w, http.StatusMethodNotAllowed)
}

func (h AnonDeckHandler) Persist(w http.ResponseWriter, r *http.Request) error {
	queryParameters := r.URL.Query()
	log.Infof("PostDeckdHandler: URI[%q] Parameters[%q]", r.URL.Path, queryParameters)

	var err error
	var deck data.Deck
	if err = haki.ReadJSON(r, &deck); err != nil {
		return haki.Err(w, err)
	}
	isCreateRequest := deck.ID == 0
	if err = raizel.Execute(deck.PersistV2); err != nil {
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

func (h AnonDeckHandler) Read(w http.ResponseWriter, r *http.Request) error {
	readParameter := path.Base(r.URL.Path)
	queryParameters := r.URL.Query()
	l.Infof("AnonDeckHandler.Read: URI[%q] ReadParameter[%q] Parameters[%q]", r.URL.Path, readParameter, queryParameters)

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
		return haki.Err(w, err)
	}
	return haki.JSON(w, http.StatusOK, deck)
}

func (h AnonDeckHandler) QueryByName(w http.ResponseWriter, r *http.Request) error {
	queryName := path.Base(r.URL.Path)
	queryParameters := r.URL.Query()
	l.Infof("AnonDeckHandler.QueryByName: URI[%q] Query[%s] Parameters[%q]", r.URL.Path, queryName, queryParameters)

	if queryName == "" {
		return haki.Status(w, http.StatusOK)
	}

	var decks []data.Deck
	err := raizel.Execute(
		func(c raizel.Client) error {
			queryDeck := data.Deck{Name: queryName}
			tmp, err := queryDeck.QueryByName(c)
			if err == nil {
				decks = tmp
			}
			return err
		},
	)
	if err != nil {
		return haki.Err(w, err)
	}
	return haki.JSON(w, http.StatusOK, decks)
}
