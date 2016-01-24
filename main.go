package main

import(
    "io"
    "log"
    "fmt"
    "flag"
    "crypto/rand"
    "net/http"
    "golang.org/x/net/websocket"
    "farm.e-pedion.com/repo/fivecolors/data"
)

var (
    bindAddress string
    mtgDir string
)

// newUUID generates a random UUID according to RFC 4122
func newUUID() (string, error) {
    uuid := make([]byte, 16)
    n, err := io.ReadFull(rand.Reader, uuid)
    if n != len(uuid) || err != nil {
        return "", err
    }
    // variant bits; see section 4.1.1
    uuid[8] = uuid[8]&^0xc0 | 0x80
    // version 4 (pseudo-random); see section 4.1.3
    uuid[6] = uuid[6]&^0xf0 | 0x40
    return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

func fiveColorsHandler(ws *websocket.Conn) {
    log.Println("ConnectionBegin")
    var command data.Command
    for {
        websocket.JSON.Receive(ws, &command)
        log.Printf("CommandReceived: Cmd[%+v]", command)
        command.ProcessId, _ = newUUID()
        websocket.JSON.Send(ws, command)
    }
    log.Println("ConnectionEnd")
}

func main() {
    flag.StringVar(&bindAddress, "bind", "127.0.0.1:7000", "Bind address for the http interface")
    flag.StringVar(&mtgDir, "mtg_dir", "/Users/raphaeljansen/Storage/fivecolors/mtg/", "FiveColors MTG database path")
    flag.Parse()
    log.Printf("DatabasePath2: Database[MTG] Bind[%s] Path[%s]", bindAddress, mtgDir)

    http.Handle("/fivecolors", websocket.Handler(fiveColorsHandler))
    http.Handle("/mtg/", http.StripPrefix("/mtg/", http.FileServer(http.Dir(mtgDir))))
    log.Printf("%s-ServerStarted: BindAddress[%s]", "0.0.1-100999", bindAddress)
    err := http.ListenAndServe(bindAddress, nil)
    if err != nil {
        panic(err)
    }
}
