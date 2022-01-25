package main

import (
	_ "embed"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/playground"
	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

//go:embed schema.graphql
var schema string

type gadget struct {
	Name     string
	Inventor string
}

type resolver struct{}

func (_ *resolver) Gadget(args struct {
	Id int32
}) gadget {
	return gadget{
		Name:     "Zapifier",
		Inventor: "Nikola Tesla",
	}
}

func main() {
	schema := graphql.MustParseSchema(schema, new(resolver), graphql.UseFieldResolvers())
	http.Handle("/query", &relay.Handler{Schema: schema})
	http.Handle("/", playground.Handler("gadgets playground", "/query"))
	log.Println("listening on port 3002...")
	log.Fatal(http.ListenAndServe(":3002", nil))
}
