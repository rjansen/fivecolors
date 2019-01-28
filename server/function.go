package server

import (
	"net/http"
	"sync"

	"github.com/rjansen/fivecolors/core/api"
	"github.com/rjansen/fivecolors/core/config"
	"github.com/rjansen/fivecolors/core/model"
	"github.com/rs/zerolog/log"
)

var (
	version       string
	once          sync.Once
	serverHandler http.HandlerFunc
)

func setup() {
	err := config.Init()
	if err != nil {
		panic(err)
	}
	log.Logger = log.With().Str("version", version).Logger()
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

func Handler(w http.ResponseWriter, r *http.Request) {
	once.Do(
		func() {
			setup()
			serverHandler = api.GraphQL
		},
	)
	serverHandler(w, r)
}
