package main

import (
	"net/http"

	"farm.e-pedion.com/repo/fivecolors/api"
	"farm.e-pedion.com/repo/fivecolors/config"
	"farm.e-pedion.com/repo/fivecolors/data"
	"farm.e-pedion.com/repo/logger"
	"farm.e-pedion.com/repo/security/identity"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var err error
	err = config.Setup()
	if err != nil {
		panic(err)
	}
	configuration := config.Get()

	err = logger.Setup(configuration.Logger)
	if err != nil {
		panic(err)
	}
	log := logger.GetLogger()

	if err = data.Setup(configuration.DB); err != nil {
		log.Panic("FivecolorsSetupError", logger.Error(err))
	}
	defer data.Close()

	if err = identity.Setup(); err != nil {
		log.Panic("IdentitySetupError", logger.Error(err))
	}

	if err = api.Setup(*configuration); err != nil {
		log.Panic("ApiSetupError", logger.Error(err))
	}

	http.Handle("/player/", api.NewGetPlayerHandler())
	http.Handle("/card/", api.NewQueryCardHandler())
	http.Handle("/expansion/", api.NewQueryExpansionHandler())
	http.Handle("/asset/", api.NewGetAssetHandler())
	http.Handle("/inventory/", api.NewInventoryHandler())
	http.Handle("/deck/", api.NewDeckHandler())

	log.Info("FivecolorsStart",
		logger.String("Version", configuration.Version),
		logger.String("HandlerVersion", configuration.Handler.Version),
		logger.String("BindAddress", configuration.Handler.BindAddress),
	)
	err = http.ListenAndServe(configuration.Handler.BindAddress, nil)
	if err != nil {
		panic(err)
	}
	log.Info("FivecolorsStop",
		logger.String("Version", configuration.Version),
		logger.String("HandlerVersion", configuration.Handler.Version),
	)
}
