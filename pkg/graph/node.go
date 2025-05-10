package graph

import (
	"errors"
	"net/url"

	"github.com/morphy76/lang-actor/pkg/common"
	"github.com/morphy76/lang-actor/pkg/framework"
	"github.com/morphy76/lang-actor/pkg/routing"
)

// ErrorInvalidRouting is returned when a routing is invalid.
var ErrorInvalidRouting = errors.New("invalid routing")

// Routable represents a node that can have routes to other nodes.
type Routable interface {
	// SetResolver sets the resolver for the node.
	SetResolver(resolver routing.Resolver)
	// GetResolver returns the resolver for the node.
	GetResolver() routing.Resolver
	// OneWayRoute add a new possible outgoing route from the node
	//
	// Parameters:
	//   - name (string): The name of the route.
	//   - destination (Node): The destination node.
	//
	// Returns:
	//   - (error): An error if the route is invalid.
	OneWayRoute(name string, destination Node) error
	// TwoWayRoute add a new possible outgoing route from the node
	//
	// Parameters:
	//   - name (string): The name of the route.
	//   - destination (Node): The destination node.
	//
	// Returns:
	//   - (error): An error if the route is invalid.
	TwoWayRoute(name string, destination Node) error
	// Edges returns the edges of the node.
	//
	// Parameters:
	//   - includeInverse (bool): Whether to include inverse edges.
	//
	// Returns:
	//   - ([]framework.Addressable): The edges of the node.
	Edges(includeInverse bool) []url.URL
	// ProceedOnAnyRoute proceeds the message on any route.
	//
	// Parameters:
	//   - msg (framework.Message): The message to be sent.
	//
	// Returns:
	//   - (error): An error if the routing is invalid.
	ProceedOnAnyRoute(msg framework.Message) error
}

// Node represents a node in the actor graph.
type Node interface {
	common.Visitable
	Routable
	framework.Addressable

	ActorRef() framework.ActorRef
}

// RootNode represents the root node of the actor graph.
type RootNode interface {
	Node
}

// EndNode represents an end node in the actor graph.
type EndNode interface {
	Node
}

// DebugNode represents a node used for debugging purposes.
type DebugNode interface {
	Node
}
