package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/rjansen/fivecolors/core/errors"
	"github.com/rjansen/fivecolors/core/util"
	"github.com/rjansen/fivecolors/core/validator"
	"github.com/rs/zerolog/log"
)

var (
	ErrInit         = errors.New("api.ErrInit")
	ErrBlankContext = errors.New("api.ErrBlankContext")
	bindAddress     string
	apiSchema       graphql.Schema
	apiUser         interface{}
)

func Init(config graphql.SchemaConfig) error {
	var (
		// TODO: Move this flag to server context
		addr = util.Getenv("SERVER_BINDADDRESS", "127.0.0.1:8080")
	)
	err := validator.Validate(
		validator.ValidateIsBlank(addr),
	)
	if err != nil {
		log.Error().Err(err).Msg("api.init.validate.err")
		return err
	}
	s, err := graphql.NewSchema(config)
	if err != nil {
		log.Error().Err(err).Msg("api.init.schema.err")
		return err
	}
	apiSchema = s
	bindAddress = addr
	apiUser = struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}{util.NewUUID(), "yuno.gasai"}
	return nil
}

func BindAddress() string {
	return bindAddress
}

func Schema() graphql.Schema {
	return apiSchema
}

type query struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}

func GraphQL(w http.ResponseWriter, r *http.Request) {
	c := util.CreateContext(context.WithValue(r.Context(), "user", apiUser), time.Minute)
	c.Info().Str("tid", c.Value("tid").(string)).Interface("user", c.Value("user")).Msg("api.grapql.request.try")
	var q query
	if r.Method == http.MethodGet {
		q = query{Query: r.URL.Query().Get("query")}
	} else {
		switch r.Header.Get("Content-Type") {
		case "application/graphql":
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				c.Error().Err(err).Msg("graphql.request.body.err")
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			q = query{Query: string(body)}
		case "application/json":
			fallthrough
		default:
			err := json.NewDecoder(r.Body).Decode(&q)
			if err != nil {
				c.Error().Err(err).Msg("graphql.request.body.err")
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
		}
	}
	if err := validator.IsBlank(q.Query); err != nil {
		c.Error().Err(err).Msg("api.grapql.request.blank.err")
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	w.Header().Set("content-type", "application/json")
	c.Info().Str("query", q.Query).Msg("api.grapql.request")
	result := graphql.Do(graphql.Params{
		Schema:        apiSchema,
		RequestString: q.Query,
		Context:       c,
	})
	c.Info().Interface("result", result).Str("query", q.Query).Msg("api.grapql.response")
	if len(result.Errors) > 0 {
		c.Error().Interface("err", result.Errors).Msg("api.graphql.request.err")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(result)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
