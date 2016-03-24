package main

import (
    "log"
    "flag"
    "net/http"
    "farm.e-pedion.com/repo/fivecolors/config"
    "farm.e-pedion.com/repo/fivecolors/api"
    "farm.e-pedion.com/repo/fivecolors/data"
)

func main() {
    log.Println("FivecolorsInitializing")
    handlerConfig := config.GetHandlerConfiguration()
    dbConfig := config.GetDBConfiguration()
    data.Setup(dbConfig)
    flag.Parse()

    http.Handle("/card", api.NewQueryCardHandler())
    http.Handle("/expansion", api.NewQueryExpansionHandler())
    http.Handle("/inventory", api.NewPostInventoryHandler())
    http.Handle("/deck", api.NewDeckHandler())
    http.Handle("/asset/", api.NewGetAssetHandler())
    //http.Handle("/asset", handler.NewAssetdHandler())
    
    log.Printf("FivecolorsStart: Version[%s] BindAddress[%s]", handlerConfig.Version, handlerConfig.BindAddress)
    err := http.ListenAndServe(handlerConfig.BindAddress, nil)
    if err != nil {
        panic(err)
    }
    log.Printf("FivecolorsEnd: Version[%s]", handlerConfig.Version)
}