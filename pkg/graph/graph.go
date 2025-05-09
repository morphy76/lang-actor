package graph

import (
	"errors"

	f "github.com/morphy76/lang-actor/pkg/framework"
	r "github.com/morphy76/lang-actor/pkg/routing"
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
	// ProceedOnAnyRoute proceeds the message on any route.
	ProceedOnAnyRoute(msg f.Message) error
}

// Node represents a node in the actor graph.
type Node interface {
	Routable
	f.Addressable
	// Name returns the name of the node.
	Name() string
	// ActorRef returns the actor reference of the actor.
	ActorRef() f.ActorRef
	// SetAddressBook sets the address book for the node.
	SetAddressBook(addressBook r.AddressBook)
	// Visit visits the node and its children.
	//
	// Parameters:
	//   - visitFn: the function to call for each node
	//   - recursive: whether to visit the node recursively
	Visit(visitFn VisitFn, recursive bool)
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

// Graph represents the actor, runnable, graph.
type Graph interface {
	// TODO Accept accepts a todo item.
	Accept(todo any) error
}

// VisitFn is a function that visits a node in the graph.
//
// Parameters:
//   - node: the node to visit
type VisitFn func(node Node)
