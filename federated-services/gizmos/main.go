package main

import (
	_ "embed"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

type gizmo struct {
	Name      string
	Craziness int32
}

type resolver struct{}

func (_ *resolver) Gizmos() []gizmo {
	return []gizmo{
		{
			Name:      "Lightning Belt",
			Craziness: 12,
		},
	}
}

//go:embed schema.graphql
var schema string

func main() {
	schema := graphql.MustParseSchema(schema, new(resolver), graphql.UseFieldResolvers())
	http.Handle("/query", &relay.Handler{Schema: schema})
	http.Handle("/", playground.Handler("gizmos playground", "/query"))
	log.Println("listening on port 3001...")
	log.Fatal(http.ListenAndServe(":3001", nil))
}
