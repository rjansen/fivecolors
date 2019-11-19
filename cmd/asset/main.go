package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rjansen/fivecolors/collection/resource"
	"github.com/rjansen/fivecolors/server"
	"github.com/rjansen/l"
	"github.com/rjansen/migi"
	"github.com/rjansen/migi/environment"
)

var (
	version string
)

type options struct {
	version     string
	bindAddress string
	asset       resource.Options
}

func newOptions() options {
	var (
		env     = migi.NewOptions(environment.NewSource())
		options = options{version: version}
	)
	env.StringVar(
		&options.bindAddress, "server_bindaddress", ":8080", "Server bind address, ip:port",
	)
	env.StringVar(
		&options.asset.StripPath, "strip_path", "/files/", "A pattern to be removed from the logical asset path when service is building the physical asset path",
	)
	env.StringVar(
		&options.asset.DataDirectory, "data_dir", "temp/db/asset/data", "Assets data directory",
	)
	env.Load()
	return options
}

func httpRouterHandler(handler http.HandlerFunc) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		handler(w, r)
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "alive")
}

func main() {
	var (
		ctx           = context.Background()
		options       = newOptions()
		assetsHandler = resource.NewHandler(options.asset)
		router        = httprouter.New()
	)

	l.Info(ctx, "asset.router.init")

	router.GET("/healthz", healthCheck)
	router.GET("/files/:assetID", httpRouterHandler(assetsHandler))

	server.Start(ctx, &http.Server{Addr: options.bindAddress, Handler: router})
}
