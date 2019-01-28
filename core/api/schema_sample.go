package api

import (
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/rjansen/fivecolors/core/util"
	"github.com/rjansen/raizel/mock"
)

var userType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var meResponseType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "MeResponse",
		Fields: graphql.Fields{
			"tid": &graphql.Field{
				Type: graphql.String,
			},
			"user": &graphql.Field{
				Type: userType,
			},
		},
	},
)

var anyScalarType = graphql.NewScalar(
	graphql.ScalarConfig{
		Name:        "AnyScalarType",
		Description: "The `AnyScalarType` scalar type represents any value.",
		Serialize: func(value interface{}) interface{} {
			return value
		},
		ParseValue: func(value interface{}) interface{} {
			return value
		},
		ParseLiteral: func(valueAST ast.Value) interface{} {
			return valueAST.GetValue()
		},
	},
)

var objectType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Object",
		Fields: graphql.Fields{
			"key": &graphql.Field{
				Type: graphql.String,
			},
			"value": &graphql.Field{
				Type: anyScalarType,
			},
		},
	},
)

var mockEntityType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "MockEntity",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"string": &graphql.Field{
				Type: graphql.String,
			},
			"integer": &graphql.Field{
				Type: graphql.Int,
			},
			"float": &graphql.Field{
				Type: graphql.Float,
			},
			"date_time": &graphql.Field{
				Type: graphql.DateTime,
			},
			"boolean": &graphql.Field{
				Type: graphql.Boolean,
			},
			/*
				"object": &graphql.Field{
					Type: graphql.NewList(objectType),
				},
			*/
		},
	},
)

var mockEntityResponseType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "MockEntityResponse",
		Fields: graphql.Fields{
			"tid": &graphql.Field{
				Type: graphql.String,
			},
			"entity": &graphql.Field{
				Type: mockEntityType,
			},
		},
	},
)

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "Query",
		Description: "A simple query schema for samples and tests purposes",
		Interfaces:  nil,
		Fields: graphql.Fields{
			"me": &graphql.Field{
				Type: meResponseType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return map[string]interface{}{
						"tid": util.NewUUID(),
						"user": map[string]interface{}{
							"id":   "fivecolorsd",
							"name": "Fivecolors D",
						},
					}, nil
				},
			},
			"mockEntity": &graphql.Field{
				Type: mockEntityResponseType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return map[string]interface{}{
						"tid": util.NewUUID(),
						"entity": mock.MockEntity{
							ID:       util.NewUUID(),
							String:   "string field",
							Integer:  999,
							Float:    999.99,
							DateTime: time.Now().UTC(),
							Boolean:  false,
							Object: mock.Object{
								"key_string":   "string value",
								"key_integer":  int64(999),
								"key_float":    float64(999.99),
								"key_boolean":  true,
								"key_datetime": time.Now().UTC(),
							},
						},
					}, nil
				},
			},
		},
	},
)

func NewMockSchema() (graphql.Schema, error) {
	return graphql.NewSchema(
		graphql.SchemaConfig{
			Query: queryType,
		},
	)
}
