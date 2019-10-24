package graphql

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/rjansen/yggdrasil"
	"github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/parser"
	"github.com/vektah/gqlparser/validator"
)

type Schema struct {
	graphql.ExecutableSchema
}

func NewSchema(executableSchema graphql.ExecutableSchema) Schema {
	return Schema{ExecutableSchema: executableSchema}
}

type Request struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}

type Response struct {
	*graphql.Response
}

func newResponse() *Response {
	return &Response{Response: new(graphql.Response)}
}

func Execute(tree yggdrasil.Tree, schema Schema, request Request) *Response {
	response := newResponse()

	doc, parserErr := parser.ParseQuery(&ast.Source{Input: request.Query})
	if parserErr != nil {
		response.Errors = append(response.Errors, parserErr)
		return response
	}

	validateErrs := validator.Validate(schema.Schema(), doc)
	if validateErrs != nil {
		response.Errors = append(response.Errors, validateErrs...)
		return response
	}

	op := doc.Operations.ForName(request.OperationName)
	vars, varsErr := validator.VariableValues(schema.Schema(), op, request.Variables)
	if varsErr != nil {
		response.Errors = append(response.Errors, varsErr)
		return response
	}

	ctx := graphql.WithRequestContext(
		context.Background(),
		graphql.NewRequestContext(doc, request.Query, vars),
	)
	return &Response{Response: schema.Query(ctx, op)}
}
