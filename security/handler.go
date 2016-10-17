package security

import (
	"farm.e-pedion.com/repo/fivecolors/data"
	"farm.e-pedion.com/repo/security/identity"
	"github.com/valyala/fasthttp"
	"log"
	"net/http"
)

type PlayerInjectableHandler interface {
	identity.AuthenticatableHandler
	GetPlayer() *data.Player
	SetPlayer(player *data.Player)
}

type InjectedPlayerHandler struct {
	identity.AuthenticatedHandler
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
	serveUnauthorizedResult := func(w http.ResponseWriter, r *http.Request) {
		log.Println("security.UnauthorizedRequest: Message[401 StatusUnauthorized]")
		w.WriteHeader(http.StatusUnauthorized)
	}
	session := handler.GetSession()
	player := &data.Player{}
	if err := player.FillFromSession(session); err != nil {
		log.Printf("security.InjectPlayerHandler.ErrorFillingSession: Err=%v", err.Error())
		serveUnauthorizedResult(w, r)
	} else {
		handler.PlayerInjectableHandler.SetPlayer(player)
		handler.PlayerInjectableHandler.ServeHTTP(w, r)
	}
}

func (handler InjectPlayerHandler) HandleRequest(ctx *fasthttp.RequestCtx) {

}
