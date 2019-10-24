package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/rjansen/fivecolors/core/resource"
	"github.com/rjansen/l"
	"github.com/rjansen/migi"
	"github.com/rjansen/yggdrasil"
)

var (
	version string
)

type options struct {
	bindAddress string
}

func newOptions() options {
	var (
		env     = migi.NewOptions(migi.NewEnvironmentSource())
		options options
	)
	env.StringVar(
		&options.bindAddress, "server_bindaddress", ":8080", "Server bind address, ip:port",
	)
	env.Parse()
	return options
}

func newTree(options options) yggdrasil.Tree {
	var (
		logger = l.NewZapLoggerDefault()
		roots  = yggdrasil.NewRoots()
		err    error
	)

	err = l.Register(&roots, logger)
	if err != nil {
		panic(err)
	}
	return roots.NewTreeDefault()
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
		options = newOptions()
		tree    = newTree(options)
		logger  = l.MustReference(tree)
		router  = httprouter.New()
	)

	logger.Info("asset.router.init")

	router.GET("/asset/healthcheck", healthCheck)
	router.GET("/asset/files/:assetID", httpRouterHandler(resource.ReadAsset))

	server := &http.Server{
		Addr:    options.bindAddress,
		Handler: router,
	}

	logger.Info("asset.router.created")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)

	go func() {
		logger.Info("asset.starting", l.NewValue("address", options.bindAddress))

		if err := server.ListenAndServe(); err != nil {
			logger.Error(
				"asset.err", l.NewValue("error", err), l.NewValue("address", options.bindAddress),
			)
		}
	}()

	<-stop

	shutDownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	logger.Info("asset.shutdown")
	server.Shutdown(shutDownCtx)
	logger.Info("asset.shutdown.gracefully")
}
