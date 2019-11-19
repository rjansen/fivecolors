package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rjansen/fivecolors/collection"
	collectionfirestore "github.com/rjansen/fivecolors/collection/firestore"
	collectiongraphql "github.com/rjansen/fivecolors/collection/graphql"
	"github.com/rjansen/fivecolors/server"
	"github.com/rjansen/l"
	"github.com/rjansen/migi"
	"github.com/rjansen/migi/environment"
	"github.com/rjansen/raizel/firestore"
)

var (
	version string
)

type options struct {
	version     string
	bindAddress string
	projectID   string
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
		&options.projectID, "project_id", "fivecolors", "GCP project identifier",
	)
	err := env.Load()
	if err != nil {
		panic(err)
	}
	return options
}

func newLogger(options options) l.Logger {
	return l.LoggerDefault
}

func newFirestoreClient(options options) firestore.Client {
	return firestore.NewClient(options.projectID)
}

func newCollectionServices(options options) (collection.Reader, collection.Writer) {
	logger := newLogger(options)
	client := newFirestoreClient(options)
	return collectionfirestore.NewReader(logger, client), collectionfirestore.NewWriter(logger, client)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "alive")
}

func main() {
	var (
		ctx            = context.Background()
		options        = newOptions()
		graphqlHandler = collectiongraphql.NewHandler(newCollectionServices(options))
		mux            = http.NewServeMux()
	)

	mux.HandleFunc("/healthz", healthCheck)
	mux.Handle("/query", graphqlHandler)

	server.Start(ctx, &http.Server{Addr: options.bindAddress, Handler: mux})
}
