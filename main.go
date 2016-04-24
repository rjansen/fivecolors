package main

import (
	"flag"
	"log"
	"net/http"

	"farm.e-pedion.com/repo/fivecolors/api"
	"farm.e-pedion.com/repo/config"
	"farm.e-pedion.com/repo/fivecolors/data"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	log.Println("FivecolorsInitializing")
    config := config.BindConfiguration()
	handlerConfig := config.HandlerConfig
	dbConfig := config.DBConfig
	if err := data.Setup(dbConfig); err != nil {
		log.Fatalf("FivecolorsSetupError: Error[%s]", err.Error())
	}
	defer data.Close()
	flag.Parse()

	http.Handle("/player/", api.NewGetPlayerHandler())
	http.Handle("/card/", api.NewQueryCardHandler())
	http.Handle("/expansion/", api.NewQueryExpansionHandler())
	http.Handle("/asset/", api.NewGetAssetHandler())
	http.Handle("/inventory/", api.NewInventoryHandler())
	http.Handle("/deck/", api.NewDeckHandler())

	log.Printf("FivecolorsStart: Version[%s] HandlerVersion[%s] BindAddress[%s]", config.Version, handlerConfig.Version, handlerConfig.BindAddress)
	err := http.ListenAndServe(handlerConfig.BindAddress, nil)
	if err != nil {
		panic(err)
	}
	log.Printf("FivecolorsStop: Version[%s] HandlerVersion[%s]", config.Version, handlerConfig.Version)
}
