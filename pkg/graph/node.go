package graph

import (
	"errors"

	"github.com/morphy76/lang-actor/pkg/common"
	"github.com/morphy76/lang-actor/pkg/routing"
)

// ErrorInvalidRouting is returned when a routing is invalid.
var ErrorInvalidRouting = errors.New("invalid routing")

// Routable represents a node that can have routes to other nodes.
type Routable interface {
	// SetResolver sets the resolver for the node.
	//
	// Parameters:
	//   - resolver (routing.Resolver): The resolver to be set.
	SetResolver(resolver routing.Resolver)
	// GetResolver returns the resolver for the node.
	//
	// Returns:
	//   - (routing.Resolver): The resolver for the node.
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
	// Edges returns the edges of the node.
	//
	// Returns:
	//   - ([]common.Addressable): The edges of the node.
	Edges() []common.Addressable
	// EdgeNames returns the names of the edges of the node.
	//
	// Returns:
	//   - ([]string): The names of the edges of the node.
	EdgeNames() []string
	// ProceedOnAnyRoute proceeds the message on any route.
	//
	// Parameters:
	//   - msg (any): The message to be sent.
	//
	// Returns:
	//   - (error): An error if the routing is invalid.
	ProceedOnAnyRoute(msg any) error
	// ProceedOnRoute proceeds the message on a specific route.
	//
	// Parameters:
	//   - name (string): The name of the route.
	//   - msg (any): The message to be sent.
	//
	// Returns:
	//   - (error): An error if the routing is invalid.
	ProceedOnRoute(name string, msg any) error
}

// Stateful represents a graph that can maintain its state and configuration.
type Stateful interface {
	// GraphConfig returns the configuration of the graph.
	GraphConfig() Configuration
	// GraphState returns the current state of the graph.
	GraphState() State
	// UpdateGraphState updates the state of the graph.
	UpdateGraphState(state State) error
}

// Connected represents a node that is connected to other nodes in the actor graph.
type Connected interface {
	// Routes returns the routes of the node.
	Routes() []string
}

// WithMutableAttributes represents a node that can have mutable attributes.
type WithMutableAttributes interface {
	// SetAttribute sets an attribute for the node.
	//
	// Parameters:
	//   - key (string): The key of the attribute.
	//   - value (any): The value of the attribute.
	SetAttribute(key string, value any)
	// GetAttribute retrieves an attribute for the node.
	//
	// Parameters:
	//   - key (string): The key of the attribute.
	// Returns:
	//   - (any): The value of the attribute, or nil if not found.
	//   - (bool): True if the attribute exists, false otherwise.
	GetAttribute(key string) (any, bool)
}

// NodeState holds the state of a node in the actor graph, including its configuration and current state.
type NodeState interface {
	Stateful
	Connected
	WithMutableAttributes
	// Outcome returns a channel that will receive the outcome of the node's processing.
	Outcome() chan string
}

// Node represents a node in the actor graph.
type Node interface {
	common.Addressable
	Routable
	// Accept processes a message sent to the node.
	Accept(msg any) error
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

// ForkJoinNode represents a node that can handle fork-join patterns in the actor graph.
type ForkJoinNode interface {
	Node
}

// ForkNode is a node that can handle forking in the actor graph.
type ForkNode interface {
	Node
}

// JoinNode is a node that can handle joining in the actor graph.
type JoinNode interface {
	Node
}
