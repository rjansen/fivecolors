package resource

import (
	"net/http"
)

type Options struct {
	StripPath     string
	DataDirectory string
}

func NewHandler(options Options) http.HandlerFunc {
	/*
		"/asset/files/",
		http.FileServer(http.Dir("db/asset/data")),
	*/
	resourceServer := http.StripPrefix(
		options.StripPath,
		http.FileServer(
			http.Dir(options.DataDirectory),
		),
	)
	return func(w http.ResponseWriter, r *http.Request) {
		resourceServer.ServeHTTP(w, r)
	}
}
