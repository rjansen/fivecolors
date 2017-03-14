package main

import (
	"github.com/rjansen/fivecolors/api"
	"github.com/rjansen/fivecolors/config"
	"github.com/rjansen/fivecolors/mtgo"
	"github.com/rjansen/l"
	raizelSQL "github.com/rjansen/raizel/sql"
	"net/http"
	"os"
	//"github.com/rjansen/avalon/identity"
	// _ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func init() {
	var err error
	if err = config.Setup(); err != nil {
		panic(err)
	}

	if err = l.Setup(&config.Value.L); err != nil {
		panic(err)
	}

	if err = raizelSQL.Setup(&config.Value.Raizel); err != nil {
		l.Panic("5colors.RaizelSetupError", l.Err(err))
	}

	// if err = identity.Setup(&configuration.Identity); err != nil {
	// 	l.Panic("IdentitySetupError", l.Err(err))
	// }
}

func main() {
	// if true {
	// 	var (
	// 		// mtgoDeckPath       = "/Users/raphaeljansen/Storage/fivecolors/mtgo/Jeskai Control.csv"
	// 		mtgoCollectionPath = "/Users/raphaeljansen/Storage/fivecolors/mtgo/MTGOCollection-20170205.csv"
	// 	)
	// 	collectionReader, err := os.Open(mtgoCollectionPath)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	if err := mtgo.ImportCollection(collectionReader); err != nil {
	// 		panic(err)
	// 	}

	// 	// mtgoDeckReader, err := os.Open(mtgoDeckPath)
	// 	// if err != nil {
	// 	// 	panic(err)
	// 	// }
	// 	// if err := mtgo.ImportDeck(mtgoDeckReader); err != nil {
	// 	// 	panic(err)
	// 	// }
	// 	return
	// }
	// http.Handle("/identity/", security.NewIdentityHandler())
	http.HandleFunc("/api/players/", api.NewAnonPlayerHandler())
	http.HandleFunc("/api/cards/", api.NewAnonCardHandler())
	http.HandleFunc("/api/tokens/", api.NewAnonTokenHandler())
	http.HandleFunc("/api/decks/", api.NewAnonDeckHandler())
	http.HandleFunc("/api/expansions/", api.NewAnonExpansionHandler())
	http.HandleFunc("/api/inventories/", api.NewAnonInventoryHandler())
	http.Handle("/api/assets/",
		http.StripPrefix("/api/assets/",
			http.FileServer(http.Dir(config.Value.AssetDir)),
		),
	)
	http.Handle("/",
		http.FileServer(http.Dir(config.Value.WebDir)),
	)

	l.Info("FivecolorsStart",
		l.String("Version", config.Value.Version),
		l.String("HandlerVersion", config.Value.Handler.Version),
		l.String("BindAddress", config.Value.Handler.BindAddress()),
	)
	err := http.ListenAndServe(config.Value.Handler.BindAddress(), nil)
	if err != nil {
		panic(err)
	}
	l.Info("FivecolorsStop",
		l.String("Version", config.Value.Version),
		l.String("HandlerVersion", config.Value.Handler.Version),
	)
}
