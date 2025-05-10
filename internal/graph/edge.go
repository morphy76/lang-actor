package graph

import (
	"net/url"
)

type edge struct {
	// Name of the route
	Name string
	// Destination of the route
	Destination url.URL
}
