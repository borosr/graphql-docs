package docs

import (
	"testing"

	"github.com/graphql-go/graphql"
)

func TestInit(t *testing.T) {
	s, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query: Query,
		Mutation: Mutation,
		Subscription: Subscription,
	})
	Init(s, Config{pretty: true})
}

var Query = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"graph": &graphql.Field{
			Type:    Graph,
			Args: graphql.FieldConfigArgument{
				"repo": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"filter": &graphql.ArgumentConfig{
					Type: Filter,
				},
			},
		},
	},
})

var Graph = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Graph",
		Fields: graphql.Fields{
			"nodes": &graphql.Field{
				Type: graphql.NewList(Node),
			},
			"edges": &graphql.Field{
				Type: graphql.NewList(Edge),
			},
		},
	},
)

var Edge = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Edge",
		Fields: graphql.Fields{
			"from": &graphql.Field{
				Type: graphql.String,
			},
			"to": &graphql.Field{
				Type: graphql.String,
			},
			"arrows": &graphql.Field{
				Type: graphql.String,
			},
			"attributes": &graphql.Field{
				Type: Attributes,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var Node = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Node",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"label": &graphql.Field{
				Type: graphql.String,
			},
			"attributes": &graphql.Field{
				Type: Attributes,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"weight": &graphql.Field{
				Type: graphql.Int,
			},
			"color": &graphql.Field{
				Type: graphql.String,
			},
			"font": &graphql.Field{
				Name: "Font",
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name: "Color",
					Fields: graphql.Fields{
						"color": &graphql.Field{
							Type: graphql.String,
						},
					},
				}),
			},
		},
	},
)

var Attributes = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Attributes",
		Fields: graphql.Fields{
			"functions": &graphql.Field{
				Type: graphql.NewList(graphql.String),
			},
			"methods": &graphql.Field{
				Type: graphql.NewList(Method),
			},
			"structs": &graphql.Field{
				Type: graphql.NewList(Struct),
			},
		},
	},
)

var Method = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Methods",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"receiver": &graphql.Field{
				Type: graphql.Type(Struct),
			},
		},
	},
)

var Struct = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Structs",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"fields": &graphql.Field{
				Type: graphql.NewList(StructField),
			},
		},
	},
)
var StructField = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "StructField",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"type": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var Filter = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "Filter",
		Fields: graphql.InputObjectConfigFieldMap{
			"value": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"filterIn": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
		},
	},
)

var Mutation = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Mutation",
	Fields:      graphql.Fields{
		"graph": &graphql.Field{
			Type:    Graph,
			Args: graphql.FieldConfigArgument{
				"repo": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"filter": &graphql.ArgumentConfig{
					Type: Filter,
				},
			},
		},
	},
})

var Subscription = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Subscription",
	Fields:      graphql.Fields{
		"graph": &graphql.Field{
			Type:    Graph,
			Args: graphql.FieldConfigArgument{
				"repo": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"filter": &graphql.ArgumentConfig{
					Type: Filter,
				},
			},
		},
	},
})
