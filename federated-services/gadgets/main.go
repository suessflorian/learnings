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

var gadgets = []gadget{
	{
		Id:       graphql.ID("0"),
		Name:     "Zapifier",
		Inventor: "Nikola Tesla",
	},
	{
		Id:       graphql.ID("1"),
		Name:     "Blitzer",
		Inventor: "Elon Musk",
	},
}

type gadget struct {
	Id       graphql.ID
	Name     string
	Inventor string
}

type resolver struct{}

func (_ *resolver) Gadget(args struct {
	Id graphql.ID
}) *gadget {
	for _, gadget := range gadgets {
		if gadget.Id == args.Id {
			return &gadget
		}
	}

	return nil
}

type gizmo struct {
	Id     graphql.ID
	Gadget gadget
}

func (_ *resolver) Gizmo(args struct {
	Id graphql.ID
}) *gizmo {
	for _, gadget := range gadgets {
		if gadget.Id == args.Id {
			return &gizmo{
				Id:     args.Id,
				Gadget: gadget,
			}
		}
	}
	return nil
}

func main() {
	schema := graphql.MustParseSchema(schema, new(resolver), graphql.UseFieldResolvers())
	http.Handle("/query", &relay.Handler{Schema: schema})
	http.Handle("/", playground.Handler("gadgets playground", "/query"))
	log.Println("listening on port 3002...")
	log.Fatal(http.ListenAndServe(":3002", nil))
}
