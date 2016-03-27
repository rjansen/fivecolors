package main

import (
	"flag"
	"log"
	"net/http"

	"farm.e-pedion.com/repo/fivecolors/api"
	"farm.e-pedion.com/repo/fivecolors/config"
	"farm.e-pedion.com/repo/fivecolors/data"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	log.Println("FivecolorsInitializing")
	handlerConfig := config.GetHandlerConfiguration()
	dbConfig := config.GetDBConfiguration()
	if err := data.Setup(dbConfig); err != nil {
		log.Fatalf("FivecolorsSetupError: Error[%s]", err.Error())
	}
	defer data.Close()
	flag.Parse()

	http.Handle("/card/", api.NewQueryCardHandler())
	http.Handle("/expansion/", api.NewQueryExpansionHandler())
	http.Handle("/asset/", api.NewGetAssetHandler())
	http.Handle("/inventory/", api.NewInventoryHandler())
	http.Handle("/deck/", api.NewDeckHandler())

	log.Printf("FivecolorsStart: Version[%s] BindAddress[%s]", handlerConfig.Version, handlerConfig.BindAddress)
	err := http.ListenAndServe(handlerConfig.BindAddress, nil)
	if err != nil {
		panic(err)
	}
	log.Printf("FivecolorsEnd: Version[%s]", handlerConfig.Version)
}
