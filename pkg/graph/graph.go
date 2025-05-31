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
	Configuration() Configuration
	// State retrieves the current state of the graph.
	//
	// Returns:
	//   - State: The current state of the graph.
	State() State
}

// WithNodeState defines the interface for components that can be configured and have a state within a graph.
type WithNodeState interface {
	// SetConfig sets the configuration for the graph-aware component.
	//
	// Parameters:
	//   - config (GraphConfiguration): The configuration to set for the component.
	SetConfig(config Configuration)
	// SetState sets the state for the graph-aware component.
	//
	// Parameters:
	//   - state (GraphState): The state to set for the component.
	SetState(state State)
}
