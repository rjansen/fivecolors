package api

import "github.com/graphql-go/graphql"

var responseType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Response",
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

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"me": &graphql.Field{
				Type: responseType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return map[string]interface{}{
						"tid":  p.Context.Value("tid"),
						"user": p.Context.Value("user"),
					}, nil
				},
			},
		},
	},
)

func newSchema() (graphql.Schema, error) {
	return graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
	})
}
