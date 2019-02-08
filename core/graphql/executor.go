package graphql

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/rjansen/yggdrasil"
	"github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/parser"
	"github.com/vektah/gqlparser/validator"
)

type Params struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}

func Execute(tree yggdrasil.Tree, schema graphql.ExecutableSchema, params Params) *graphql.Response {
	response := new(graphql.Response)

	doc, parserErr := parser.ParseQuery(&ast.Source{Input: params.Query})
	if parserErr != nil {
		response.Errors = append(response.Errors, parserErr)
		return response
	}

	validateErrs := validator.Validate(schema.Schema(), doc)
	if validateErrs != nil {
		response.Errors = append(response.Errors, validateErrs...)
		return response
	}

	op := doc.Operations.ForName(params.OperationName)
	vars, varsErr := validator.VariableValues(schema.Schema(), op, params.Variables)
	if varsErr != nil {
		response.Errors = append(response.Errors, varsErr)
		return response
	}

	ctx := graphql.WithRequestContext(
		context.Background(),
		graphql.NewRequestContext(doc, params.Query, vars),
	)
	return schema.Query(ctx, op)
}
