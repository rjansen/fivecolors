package api

import (
	"io"
	//"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"farm.e-pedion.com/repo/fivecolors/data"
	//"farm.e-pedion.com/repo/ffm-security/identity"
)

const (
	hydrateSmall = "small"
	hydrateFull  = "full"
)

//NewQueryCardHandler creates a new QueryCardHandler instance
func NewQueryCardHandler() http.Handler {
	return &QueryCardHandler{}
	//return identity.NewCookieAuthenticatedHandler(&GetCardHandler{})
}

//QueryCardHandler is the handler to get and query Cards
type QueryCardHandler struct {
	//    session *identity.Session
}

//func (h *QueryCardHandler) SetSession(session *identity.Session) {
//    h.session = session
//}

func (h *QueryCardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	urlPathParameters := strings.Split(r.URL.Path, "/")
	queryParameters := r.URL.Query()
	log.Printf("GetHandler: URL[%q] PathParameters[%q] QueryParameters[%q]", r.URL.Path, urlPathParameters, queryParameters)

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
			log.Printf("QueryCardHandler.CardReadError: Card.ID[%v] Error[%v]", cardID, convertIDErr)
			http.NotFound(w, r)
			//http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := card.Read(); err != nil {
			log.Printf("QueryCardHandler.CardReadError: Card.ID[%v] Error[%v]", cardID, err)
			http.NotFound(w, r)
			//http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("QueryCardHandler.GotCard: Card.ID[%v] Hydrate[%s]", cardID, hydrate)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		jsonData, err := json.Marshal(card)
		bytesWritten, err := w.Write(jsonData)
		if err != nil {
			log.Printf("QueryCardHandler.WriteResponseError: Card.ID[%v] Error[%v]", cardID, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			log.Printf("QueryCardHandler.JsonWritten: Card.ID[%v] Bytes[%v]", cardID, bytesWritten)
		}
	} else {
		if len(queryParameters) <= 0 {
            //w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		regexCost := queryParameters.Get("rx_cost")
		notRegexCost := queryParameters.Get("nrx_cost")
		regexType := queryParameters.Get("rx_type")
		notRegexType := queryParameters.Get("nrx_type")
		expansionParam := queryParameters.Get("e")
		numberParam := queryParameters.Get("n")
		order := queryParameters.Get("order")

		queryRestrinctions := make(map[string]interface{})
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
				log.Printf("QueryCardHandler.ExpansionParameterError: Parameter[%v] Error[%v]", expansionParam, convertErr)
			}
		}
		if numberParam != "" {
            queryRestrinctions["c.multiverse_number = ?"] = numberParam
		}
		if order != "" {
			order = "c.id"
		}
		card := &data.Card{}
		cards, err := card.Query(queryRestrinctions, order)
		if err != nil {
			log.Printf("QueryCardHandler.CardQueryError: QueryRestrictions[%v] Error[%v]", queryRestrinctions, err)
			//http.NotFound(w, r)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cardsSize := len(cards)
		log.Printf("QueryCardHandler.QueryCards: Cards.Len[%v] Hydrate[%s]", cardsSize, hydrate)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		jsonData, err := json.Marshal(cards)
		bytesWritten, err := w.Write(jsonData)
		if err != nil {
			log.Printf("QueryCardHandler.WriteResponseError: Cards.Len[%v] Error[%v]", cardsSize, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			log.Printf("QueryCardHandler.JsonWritten: Cards.Len[%v] Bytes[%v]", cardsSize, bytesWritten)
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

func (h *QueryExpansionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	urlPathParameters := strings.Split(r.URL.Path, "/")
	queryParameters := r.URL.Query()
	log.Printf("GetHandler: URL[%q] PathParameters[%q] QueryParameters[%q]", r.URL.Path, urlPathParameters, queryParameters)

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
			log.Printf("QueryExpansionHandler.ExpansionReadError: Expansion.ID[%v] Error[%v]", expansionID, convertIDErr)
			http.NotFound(w, r)
			//http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := expansion.Read(); err != nil {
			log.Printf("QueryExpansionHandler.CardReadError: Expansion.ID[%v] Error[%v]", expansionID, err)
			http.NotFound(w, r)
			//http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("QueryExpansionkHandler.GotExpansion: Expansion.ID[%v] Hydrate[%s]", expansionID, hydrate)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		jsonData, err := json.Marshal(expansion)
		bytesWritten, err := w.Write(jsonData)
		if err != nil {
			log.Printf("QueryExpansionHandler.WriteResponseError: Expansion.ID[%v] Error[%v]", expansionID, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			log.Printf("QueryExpansionHandler.JsonWritten: Expansion.ID[%v] Bytes[%v]", expansionID, bytesWritten)
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
			log.Printf("QueryExpansionHandler.ExpansionQueryError: QueryRestrictions[%v] Error[%v]", queryRestrinctions, err)
			//http.NotFound(w, r)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		expansionSize := len(expansions)
		log.Printf("QueryExpansionHandler.QueryExpansions: Expansions.Len[%v] Hydrate[%s]", expansionSize, hydrate)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		jsonData, err := json.Marshal(expansions)
		bytesWritten, err := w.Write(jsonData)
		if err != nil {
			log.Printf("QueryExpansionHandler.WriteResponseError: Expansions.Len[%v] Error[%v]", expansionSize, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			log.Printf("QueryExpansionHandler.JsonWritten: Expansions.Len[%v] Bytes[%v]", expansionSize, bytesWritten)
		}
	}
}

//NewPostInventoryHandler creates a new PostInvetoryHandler instance
func NewPostInventoryHandler() http.Handler {
	return &PostInventoryHandler{}
//return identity.NewCookieAuthenticatedHandler(&PostDeckHandler{})
}

//PostInventoryHandler is the handler to updates Inventory
type PostInventoryHandler struct {
	//    session *identity.Session
}

//func (h *PostInvetoryHandler) SetSession(session *identity.Session) {
//    h.session = session
//}

func (h *PostInventoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	urlPathParameters := strings.Split(r.URL.Path, "/")
	queryParameters := r.URL.Query()
	log.Printf("PostInvetoryHandler: URL[%q] PathParameters[%q] QueryParameters[%q]", r.URL.Path, urlPathParameters, queryParameters)

	inventory := data.Inventory{}
	if err := inventory.Unmarshal(r.Body); err != nil {
		log.Printf("PostInvetoryHandler.UnreadableBodyError: Error[%v]", err)
		http.Error(w, "PostInvetoryHandler.UnreadableBodyError", http.StatusBadRequest)
		return
	}
	if err := inventory.Update(); err != nil {
		log.Printf("PostInvetoryHandler.CreateDeckError: Error[%v]", err)
		http.Error(w, "PostInvetoryHandler.CreateDeckError", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusAccepted)
}

//NewDeckHandler creates a new DeckHandler instance
func NewDeckHandler() http.Handler {
	return &DeckHandler{getHandler: &GetDeckHandler{}, postHandler: &PostDeckHandler{}}
	//return identity.NewCookieAuthenticatedHandler(&GetCardHandler{})
}

//DeckHandler is the handler to get and post Decks
type DeckHandler struct {
	getHandler  http.Handler
	postHandler http.Handler
	//    session *identity.Session
}

func (h *DeckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.getHandler.ServeHTTP(w, r)
	} else if r.Method == "POST" {
		h.postHandler.ServeHTTP(w, r)
	} else {
		http.Error(w, "DeckHandler.MethodNotAllowed: Method="+r.Method, http.StatusInternalServerError)
	}
}

//NewPostDeckHandler creates a new PostDeckHandler instance
//func NewPostDeckHandler() http.Handler {
//	return &PostDeckHandler{}
//return identity.NewCookieAuthenticatedHandler(&PostDeckHandler{})
//}

//PostDeckHandler is the handler to creates Decks
type PostDeckHandler struct {
	//    session *identity.Session
}

//func (h *PostDeckHandler) SetSession(session *identity.Session) {
//    h.session = session
//}

func (h *PostDeckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	urlPathParameters := strings.Split(r.URL.Path, "/")
	queryParameters := r.URL.Query()
	log.Printf("PostDeckdHandler: URL[%q] PathParameters[%q] QueryParameters[%q]", r.URL.Path, urlPathParameters, queryParameters)

	deck := data.Deck{}
	if err := deck.Unmarshal(r.Body); err != nil {
		log.Printf("PostDeckHandler.UnreadableBodyError: Error[%v]", err)
		http.Error(w, "PostDeckHandler.UnreadableBodyError", http.StatusBadRequest)
		return
	}
	if err := deck.Create(); err != nil {
		log.Printf("PostDeckHandler.CreateDeckError: Error[%v]", err)
		http.Error(w, "PostDeckHandler.CreateDeckError", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, strconv.Itoa(deck.ID))
}

//NewGetDeckHandler creates a new GetDeckHandler instance
//func NewGetDeckHandler() http.Handler {
//	return &GetDeckHandler{}
//return identity.NewCookieAuthenticatedHandler(&GetDeckHandler{})
//}

//GetDeckHandler is the handler to get Decks
type GetDeckHandler struct {
	//    session *identity.Session
}

//func (h *GetDeckHandler) SetSession(session *identity.Session) {
//    h.session = session
//}

func (h *GetDeckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	urlPathParameters := strings.Split(r.URL.Path, "/")
	queryParameters := r.URL.Query()
	log.Printf("GetDeckHandler: URL[%q] PathParameters[%q] QueryParameters[%q]", r.URL.Path, urlPathParameters, queryParameters)

	deckID := urlPathParameters[2]
	if deckID == "" {
		http.NotFound(w, r)
		return
	}
	hydrate := queryParameters.Get("hydrate")
	if hydrate == "" {
		hydrate = "small"
	}
	deck := &data.Deck{}
	var convertIDErr error
	deck.ID, convertIDErr = strconv.Atoi(deckID)
	if convertIDErr != nil {
		log.Printf("GetDeckHandler.DeckReadError: Deck.ID[%v] Error[%v]", deckID, convertIDErr)
		http.NotFound(w, r)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := deck.Read(); err != nil {
		log.Printf("GetDeckHandler.DeckReadError: Deck.ID[%v] Error[%v]", deckID, err)
		http.NotFound(w, r)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("GetDeckHandler.GotDeck: Deck.ID[%v] Hydrate[%s]", deckID, hydrate)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	jsonData, err := json.Marshal(deck)
	bytesWritten, err := w.Write(jsonData)
	if err != nil {
		log.Printf("GetDeckHandler.WriteResponseError: Deck.ID[%v] Error[%v]", deckID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		log.Printf("GetDeckHandler.JsonWritten: Deck.ID[%v] Bytes[%v]", deckID, bytesWritten)
	}
}

//NewGetAssetHandler creates a new GetAssetHandler instance
func NewGetAssetHandler() http.Handler {
	return &GetAssetHandler{}
	//return identity.NewCookieAuthenticatedHandler(&GetCardHandler{})
}

//GetAssetHandler is the handler to get Assets
type GetAssetHandler struct {
	//    session *identity.Session
}

//func (h *GetDeckHandler) SetSession(session *identity.Session) {
//    h.session = session
//}

func (h *GetAssetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	urlPathParameters := strings.Split(r.URL.Path, "/")
	queryParameters := r.URL.Query()
	log.Printf("GetAssetHandler: URL[%q] PathParameters[%q] QueryParameters[%q]", r.URL.Path, urlPathParameters, queryParameters)

	assetID := urlPathParameters[2]
	if assetID == "" {
		http.NotFound(w, r)
		return
	}
	asset := &data.Asset{}
	var convertIDErr error
	asset.ID, convertIDErr = strconv.Atoi(assetID)
	if convertIDErr != nil {
		log.Printf("GetAssetHandler.AssetReadError: Asset.ID[%v] Error[%v]", assetID, convertIDErr)
		//http.NotFound(w, r)
		http.Error(w, convertIDErr.Error(), http.StatusInternalServerError)
		return
	}
	if err := asset.Read(); err != nil {
		log.Printf("GetAssetkHandler.AssetReadError: Asset.ID[%v] Error[%v]", assetID, err)
		http.NotFound(w, r)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("GetAssetHandler.GotDeck: Asset.ID[%v]", assetID)
	w.Header().Set("Content-Type", "image/jpeg")
	w.WriteHeader(http.StatusOK)
	bytesWritten, err := w.Write(asset.BinaryData)
	if err != nil {
		log.Printf("GetAssetHandler.WriteResponseError: Asset.ID[%v] Error[%v]", assetID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		log.Printf("GetAssetHandler.BytesWritten: Asset.ID[%v] Bytes[%v]", assetID, bytesWritten)
	}
}
