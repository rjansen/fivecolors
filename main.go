package main

import (
	"net/http"

	"github.com/rjansen/fivecolors/api"
	// "github.com/rjansen/fivecolors/data"
	"github.com/rjansen/fivecolors/config"
	"github.com/rjansen/l"
	// "github.com/rjansen/migi"
	raizelSQL "github.com/rjansen/raizel/sql"
	//"github.com/rjansen/avalon/identity"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var err error
	if err = config.Setup(); err != nil {
		panic(err)
	}

	if err = l.Setup(&config.Value.L); err != nil {
		panic(err)
	}

	// if err = data.Setup(configuration.DB); err != nil {
	// 	l.Panic("FivecolorsSetupError", l.Err(err))
	// }
	// defer data.Close()

	if err = raizelSQL.Setup(&config.Value.Raizel); err != nil {
		l.Panic("5colors.RaizelSetupError", l.Err(err))
	}

	// if err = identity.Setup(&configuration.Identity); err != nil {
	// 	l.Panic("IdentitySetupError", l.Err(err))
	// }

	// if err = api.Setup(*configuration); err != nil {
	// 	l.Panic("ApiSetupError", l.Err(err))
	// }

	// http.Handle("/identity/", security.NewIdentityHandler())
	// http.Handle("/player/", api.NewGetPlayerHandler())
	// http.Handle("/card/", api.NewQueryCardHandler())
	// http.Handle("/expansion/", api.NewQueryExpansionHandler())
	// http.Handle("/asset/", api.NewGetAssetHandler())
	// http.Handle("/inventory/", api.NewInventoryHandler())
	// http.Handle("/deck/", api.NewDeckHandler())
	http.HandleFunc("/api/deck/", api.NewAnonDeckHandler())

	l.Info("FivecolorsStart",
		l.String("Version", config.Value.Version),
		l.String("HandlerVersion", config.Value.Handler.Version),
		l.String("BindAddress", config.Value.Handler.BindAddress()),
	)
	err = http.ListenAndServe(config.Value.Handler.BindAddress(), nil)
	if err != nil {
		panic(err)
	}
	l.Info("FivecolorsStop",
		l.String("Version", config.Value.Version),
		l.String("HandlerVersion", config.Value.Handler.Version),
	)
}
