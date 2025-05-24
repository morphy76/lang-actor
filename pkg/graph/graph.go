package graph

import (
	"github.com/morphy76/lang-actor/pkg/common"
	"github.com/morphy76/lang-actor/pkg/routing"
)

// GraphConfiguration defines the interface for a graph configuration.
type GraphConfiguration interface {
	// Keys returns the keys of the graph configuration.
	//
	// Returns:
	//   - []string: A slice of keys in the graph configuration.
	Keys() []string
	// Value retrieves the value associated with the given key in the graph configuration.
	//
	// Parameters:
	//   - key (string): The key for which to retrieve the value.
	//
	// Returns:
	//   - any: The value associated with the key, or nil if the key does not exist.
	//   - bool: A boolean indicating whether the key exists in the configuration.
	Value(key string) (any, bool)
}

// GraphState defines the interface for managing state within a graph.
type GraphState interface {
	// Set stores a value in the graph state associated with the given key.
	//
	// Parameters:
	//   - key (string): The key under which to store the value.
	//   - value (any): The value to store in the graph state.
	Set(key string, value any)
	// Value retrieves the value associated with the given key in the graph state.
	//
	// Parameters:
	//   - key (string): The key for which to retrieve the value.
	//
	// Returns:
	//   - any: The value associated with the key, or nil if the key does not exist.
	//   - bool: A boolean indicating whether the key exists in the state.
	Value(key string) (any, bool)
	// Keys retrieves all keys present in the graph state.
	//
	// Returns:
	//   - []string: A slice of keys in the graph state.
	Keys() []string
}

// Graph represents the actor, runnable, graph.
type Graph interface {
	routing.Resolver
	common.MessageHandler
	// Config retrieves the configuration of the graph.
	//
	// Returns:
	//   - GraphConfiguration: The configuration of the graph.
	Config() GraphConfiguration
	// State retrieves the current state of the graph.
	//
	// Returns:
	//   - GraphState: The current state of the graph.
	State() GraphState
}

// GraphAware defines the interface for components that can be configured and have a state within a graph.
type GraphAware interface {
	// SetConfig sets the configuration for the graph-aware component.
	//
	// Parameters:
	//   - config (GraphConfiguration): The configuration to set for the component.
	SetConfig(config GraphConfiguration)
	// SetState sets the state for the graph-aware component.
	//
	// Parameters:
	//   - state (GraphState): The state to set for the component.
	SetState(state GraphState)
	// GetConfig retrieves the configuration of the graph-aware component.
	//
	// Returns:
	//   - GraphConfiguration: The configuration of the component.
	GetConfig() GraphConfiguration
	// GetState retrieves the state of the graph-aware component.
	//
	// Returns:
	//   - GraphState: The state of the component.
	GetState() GraphState
}
