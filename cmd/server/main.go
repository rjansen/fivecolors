package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/rjansen/fivecolors/core/api"
	"github.com/rjansen/fivecolors/core/config"
	"github.com/rjansen/fivecolors/core/model"
	"github.com/rs/zerolog/log"
)

func init() {
	err := config.Init()
	if err != nil {
		panic(err)
	}
	log.Info().Msg("server.init.model.try")
	err = model.Init()
	if err != nil {
		panic(err)
	}
	log.Info().Msg("server.init.api.try")
	err = api.Init(model.NewSchemaConfig())
	if err != nil {
		panic(err)
	}
	log.Info().Msg("server.initialized")
}

func httpRouterHandler(handler http.HandlerFunc) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		handler(w, r)
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "alive")
}

func main() {
	log.Info().Msg("server.start")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	log.Info().Msg("server.router.init")
	router := httprouter.New()
	router.GET("/api/healthcheck", httpRouterHandler(healthCheck))
	router.GET("/api/query", httpRouterHandler(api.GraphQL))
	router.POST("/api/query", httpRouterHandler(api.GraphQL))

	log.Info().Str("address", api.BindAddress()).Msg("server.create")
	server := &http.Server{
		Addr:    api.BindAddress(),
		Handler: router,
	}

	go func() {
		log.Info().Str("address", api.BindAddress()).Msg("server.listening")

		if err := server.ListenAndServe(); err != nil {
			log.Fatal().Err(err).Str("address", api.BindAddress()).Msg("server.listen.err")
		}
	}()

	<-stop

	log.Info().Str("address", api.BindAddress()).Msg("server.shutdown")
	shutDownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(shutDownCtx)
	log.Info().Str("address", api.BindAddress()).Msg("server.gracefully.stoped")
	log.Info().Msg("server.end")
}
