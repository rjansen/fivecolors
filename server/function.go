package main

import (
	"net/http"
	"sync"

	"github.com/rjansen/fivecolors/core/api"
)

var (
	once sync.Once
)

func Handler(w http.ResponseWriter, r *http.Request) {
	once.Do(
		func() {
			setup()
		},
	)
	api.GraphQL(w, r)
}
