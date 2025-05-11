package graph

import (
	"github.com/morphy76/lang-actor/pkg/framework"
)

type edge struct {
	// Name of the route
	Name string
	// Destination of the route
	Destination framework.Addressable
}
