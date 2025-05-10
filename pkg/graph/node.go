package graph

import (
	"errors"

	"github.com/morphy76/lang-actor/pkg/framework"
)

// ErrorInvalidRouting is returned when a routing is invalid.
var ErrorInvalidRouting = errors.New("invalid routing")

// Routable represents a node that can have routes to other nodes.
type Routable interface {
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
	// RouteNames returns the names of the routes from this node
	//
	// Returns:
	//   - ([]string): A slice of strings containing the names of the routes.
	RouteNames() []string
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
	Routable
	framework.Addressable
	// Name returns the name of the node.
	//
	// Returns:
	//   - (string): The name of the node.
	Name() string
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
