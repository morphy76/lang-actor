package graph

import (
	"github.com/morphy76/lang-actor/pkg/routing"
)

// Configuration defines the interface for a graph configuration.
type Configuration interface {
}

// State defines the interface for managing state within a graph.
type State interface {
}

// Graph represents the actor, runnable, graph.
type Graph interface {
	routing.Resolver
	// Config retrieves the configuration of the graph.
	//
	// Returns:
	//   - Configuration: The configuration of the graph.
	Config() Configuration
	// State retrieves the current state of the graph.
	//
	// Returns:
	//   - State: The current state of the graph.
	State() State
	// UpdateState updates the state of the graph.
	//
	// Parameters:
	//   - state State: The new state to set for the graph.
	//
	// Returns:
	//   - error: An error if the update fails, nil otherwise.
	UpdateState(state State) error
}
