package graph

import (
	c "github.com/morphy76/lang-actor/pkg/common"
)

type edge struct {
	// Name of the route
	Name string
	// Destination of the route
	Destination c.Addressable
}
