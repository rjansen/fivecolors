package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"bitbucket.org/raphaeljansen/fivecolors-etl/config"
	"bitbucket.org/raphaeljansen/fivecolors-etl/resource"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

var (
	version string
)

func init() {
	err := config.Init()
	if err != nil {
		panic(err)
	}
	log.Logger = log.With().Str("version", version).Logger()
	log.Info().Msg("asset.resource.init.try")
	err = resource.Init()
	if err != nil {
		panic(err)
	}
	log.Info().Msg("asset.initialized")
}

func healthCheck(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "alive")
}

func main() {
	log.Info().Msg("asset.start")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	log.Info().Msg("asset.router.init")
	router := httprouter.New()
	router.GET("/asset/healthcheck", healthCheck)
	router.GET("/asset/files/:assetID", resource.ReadAsset)

	log.Info().Str("address", resource.BindAddress()).Msg("asset.create")
	server := &http.Server{
		Addr:    resource.BindAddress(),
		Handler: router,
	}

	go func() {
		log.Info().Str("address", resource.BindAddress()).Msg("asset.listening")

		if err := server.ListenAndServe(); err != nil {
			log.Fatal().Err(err).Str("address", resource.BindAddress()).Msg("asset.listen.err")
		}
	}()

	<-stop

	log.Info().Str("address", resource.BindAddress()).Msg("asset.shutdown")
	shutDownCtx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	server.Shutdown(shutDownCtx)
	log.Info().Str("address", resource.BindAddress()).Msg("asset.gracefully.stoped")
	log.Info().Msg("asset.end")
}
