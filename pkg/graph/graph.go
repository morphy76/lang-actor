package graph

import (
	"github.com/morphy76/lang-actor/pkg/routing"
)

const (
	// SkipOutcome is a special outcome that indicates the graph should skip to the next step without any action.
	SkipOutcome = "/dev/null"
	// WhateverOutcome is a special outcome that indicates the graph can proceed with any action.
	WhateverOutcome = ""
)

var staticNoConfiguration Configuration = (*NoConfiguration)(nil)

// NoConfiguration is an empty struct that implements the Configuration interface.
type NoConfiguration struct{}

// Configuration defines the interface for a graph configuration.
type Configuration interface {
}

// State defines the interface for managing state within a graph.
type State interface {
	// MergeChange appends a new state to the graph.
	//
	// Parameters:
	//   - purpose (any): The purpose of the state.
	//   - value (any): The value of the state.
	//
	// Returns:
	//   - error: An error if the append operation fails, nil otherwise.
	MergeChange(purpose any, value any) error

	// Unwrap retrieves the underlying, non-proxy, state.
	//
	// Returns:
	//   - State: The underlying state.
	Unwrap() State

	// TODO
	ReadAttribute(name string) any
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
	// StateChangedCh returns a channel that is closed when the state of the graph changes.
	StateChangedCh() <-chan State
}
