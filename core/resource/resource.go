package resource

import (
	"context"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rjansen/fivecolors/core/errors"
	"github.com/rjansen/fivecolors/core/util"
	"github.com/rjansen/fivecolors/core/validator"
	"github.com/rs/zerolog/log"
)

var (
	ErrInit         = errors.New("resource.ErrInit")
	ErrBlankContext = errors.New("resource.ErrBlankContext")
	bindAddress     string
	dataDir         string
	assetServer     http.Handler
	resourceUser    interface{}
)

func Init() error {
	var (
		addr = util.Getenv("RESOURCE_BINDADDRESS", "127.0.0.1:8081")
		dir  = util.Getenv("RESOURCE_DATADIR", "db/asset/data")
	)
	err := validator.Validate(
		validator.ValidateIsBlank(addr),
	)
	if err != nil {
		log.Error().Err(err).Msg("resource.init.validate.err")
		return err
	}
	dataDir = dir
	assetServer = http.StripPrefix("/asset/files/", http.FileServer(http.Dir(dataDir)))
	bindAddress = addr
	resourceUser = struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}{util.NewUUID(), "yuno.gasai"}
	return nil
}

func BindAddress() string {
	return bindAddress
}

func DataDir() string {
	return dataDir
}

func ReadAsset(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	c := util.CreateContext(context.WithValue(r.Context(), "user", resourceUser), time.Minute)
	c.Info().Str("path", r.URL.Path).Interface("user", c.Value("user")).Msg("resource.asset.request.try")
	assetServer.ServeHTTP(w, r)
}
