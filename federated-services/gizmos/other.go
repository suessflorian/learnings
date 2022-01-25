package main

type Service struct {
	Schema  string
	Name    string
	Version string
}

func (_ *resolver) Service() Service {
	return Service{
		Schema:  schema,
		Name:    "gizmo-service",
		Version: "v1.0.1",
	}
}
