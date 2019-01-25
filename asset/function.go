package asset

import (
	"net/http"
	"sync"

	"github.com/rjansen/fivecolors/core/config"
	"github.com/rjansen/fivecolors/core/resource"
	"github.com/rs/zerolog/log"
)

var (
	once sync.Once
)

func setup() {
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

func Handler(w http.ResponseWriter, r *http.Request) {
	once.Do(
		func() {
			setup()
		},
	)
	resource.ReadAsset(w, r)
}
