package graph

import (
	"errors"

	f "github.com/morphy76/lang-actor/pkg/framework"
)

// ErrorInvalidRouting is returned when a routing is invalid.
var ErrorInvalidRouting = errors.New("invalid routing")

// Routable represents a node that can have routes to other nodes.
type Routable interface {
	// OneWayRoute add a new possible outgoing route from the node
	OneWayRoute(name string, destination Node) error
	// TwoWayRoute add a new possible outgoing route from the node
	TwoWayRoute(name string, destination Node) error
	// RouteNames returns the names of the routes from this node
	RouteNames() []string
}

// Node represents a node in the actor graph.
type Node interface {
	Routable
	f.Addressable
	// Name returns the name of the node.
	Name() string
	// ProceedOnFirstRoute proceeds with the first route available.
	ProceedOnFirstRoute(msg f.Message) error
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
