package graph

import (
	g "github.com/morphy76/lang-actor/pkg/graph"
)

type edge struct {
	// Name of the route
	Name string
	// Destination of the route
	Destination g.Node
}
