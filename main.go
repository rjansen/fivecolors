package main

import(
    "log"
    "flag"
    "net/http"
)

var (
    mtgDir string
)

func main() {
    flag.StringVar(&mtgDir, "mtg_dir", "/Users/raphaeljansen/Storage/fivecolors/mtg/", "FiveColors MTG database path")
    flag.Parse()
    log.Printf("DatabasePath2: Database[MTG] Path[%s]", mtgDir)
    
    http.Handle("/mtg/", http.StripPrefix("/mtg/", http.FileServer(http.Dir(mtgDir))))
    bindAddress := "127.0.0.1:7000"
    log.Printf("%s-ServerStarted: BindAddress[%s]", "0.0.1-100999", bindAddress)
    err := http.ListenAndServe(bindAddress, nil)
    if err != nil {
        panic(err)
    }
}
