package security

import (
	// "github.com/rjansen/avalon/identity"
	"github.com/rjansen/fivecolors/data"
	// "github.com/rjansen/l"
	// "github.com/valyala/fasthttp"
	"net/http"
	// "strings"
)

type PlayerInjectableHandler interface {
	// identity.AuthenticatableHandler
	GetPlayer() *data.Player
	SetPlayer(player *data.Player)
}

type InjectedPlayerHandler struct {
	// identity.AuthenticatedHandler
	player *data.Player
}

func (p *InjectedPlayerHandler) SetPlayer(player *data.Player) {
	p.player = player
}

func (p *InjectedPlayerHandler) GetPlayer() *data.Player {
	return p.player
}

func NewInjectPlayerHandler(handler PlayerInjectableHandler) identity.AuthenticatableHandler {
	return &InjectPlayerHandler{PlayerInjectableHandler: handler}
}

type InjectPlayerHandler struct {
	PlayerInjectableHandler
}

func (handler *InjectPlayerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// serveUnauthorizedResult := func(w http.ResponseWriter, r *http.Request) {
	// 	log.Println("security.UnauthorizedRequest: Message[401 StatusUnauthorized]")
	// 	w.WriteHeader(http.StatusUnauthorized)
	// }
	// session := handler.GetSession()
	// player := &data.Player{}
	// if err := player.FillFromSession(session); err != nil {
	// 	log.Printf("security.InjectPlayerHandler.ErrorFillingSession: Err=%v", err.Error())
	// 	serveUnauthorizedResult(w, r)
	// } else {
	// 	handler.PlayerInjectableHandler.SetPlayer(player)
	// 	handler.PlayerInjectableHandler.ServeHTTP(w, r)
	// }
}

func (handler InjectPlayerHandler) HandleRequest(ctx *fasthttp.RequestCtx) {

}

//NewIdentityHandler creates a new GetPlayerHandler instance
func NewIdentityHandler() http.Handler {
	return identity.NewHeaderAuthenticatedHandler(&IdentityHandler{})
}

//IdentityHandler is the handler to get Players
type IdentityHandler struct {
	identity.AuthenticatedHandler
}

func (h IdentityHandler) HandleRequest(ctx *fasthttp.RequestCtx) {

}

func (h *IdentityHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// urlPathParameters := strings.Split(r.URL.Path, "/")
	// queryParameters := r.URL.Query()
	// logger.Infof("IdentityHandler: URL[%q] PathParameters[%q] QueryParameters[%q]", r.URL.Path, urlPathParameters, queryParameters)

	// session := h.GetSession()
	// if strings.TrimSpace(session.Username) == "" {
	// 	http.NotFound(w, r)
	// 	return
	// }
	// player, err := data.GetPlayer(session.Username)
	// if err != nil {
	// 	logger.Errorf("IdentityHandler.GetPlayerError: Session.ID[%v] Player.Username[%v] Error[%v]", session.ID, session.Username, err.Error())
	// 	http.NotFound(w, r)
	// 	//http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// session.Data = player.SessionData()

	// w.Header().Set("Content-Type", "application/octet-stream")
	// w.WriteHeader(http.StatusOK)

	// jwtBytes, err := identity.Serialize(*session)
	// bytesWritten, err := w.Write(jwtBytes)
	// if err != nil {
	// 	logger.Errorf("IdentityHandler.WriteResponseError: Player[%v] Error[%v]", player, err)
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// } else {
	// 	logger.Debugf("IdentityHandler.JwtWritten: Player[%v] Bytes[%v]", player, bytesWritten)
	// }

}
