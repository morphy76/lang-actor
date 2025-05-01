package framework

import (
	"errors"
	"net/url"
)

var ErrorInvalidActorAddress = errors.New("invalid actor address")

// ErrorActorNotStarted is returned when an actor is not started.
var ErrorActorNotStarted = errors.New("actor not started")

// ErrorActorAlreadyStarted is returned when an actor is already started.
var ErrorActorAlreadyStarted = errors.New("actor already started")

// ErrorActorNotIdle is returned when an actor is not idle.
var ErrorActorNotIdle = errors.New("actor not idle")

// ErrorActorNotRunning is returned when an actor is not running.
var ErrorActorNotRunning = errors.New("actor not running")

// ActorStatus represents the status of an actor.
type ActorStatus int8

const (
	// ActorStatusIdle indicates that the actor is idle.
	ActorStatusIdle ActorStatus = iota
	// ActorStatusStarting indicates that the actor is starting.
	ActorStatusStarting
	// ActorStatusRunning indicates that the actor is running.
	ActorStatusRunning
	// ActorStatusStopping indicates that the actor is stopping.
	ActorStatusStopping
)

// Actor is part of the actor model framework underlying lang-actor.
//
// Type Parameters:
//   - M: The type of the message.
//   - T: The type of the actor state.
type Actor[M any, T any] interface {
	// Actor URI
	Address() url.URL
	// Start the actor
	Start() error
	// Stop the actor
	Stop() (chan bool, error)
	// Status of the actor
	Status() ActorStatus
	// State of the actor
	State() T
	// Deliver a message to the actor
	Deliver(msg Message[M]) error
}

// Payload represents the state of an actor.
//
// Type Parameters:
//   - T: The type of the actor state.
type Payload[T any] interface {
	// Cast returns the exact type of the struct implementing this interface.
	Cast() T
}

// Message is the interface for messages sent to actors.
//
// Type Parameters:
//   - T: The type of the message.
type Message[T any] interface {
	// Sender returns the URL of the sender.
	Sender() url.URL
	// Mutation returns true if the message is an actor mutation.
	Mutation() bool
	// ToImplementation returns the exact type of the struct implementing this interface.
	Cast() T
}

// ProcessingFn defines a generic function type for processing messages within an actor system.
// It takes a message of type T and the current actor state of type S, and returns the updated
// state of type S along with an error, if any.
//
// Type Parameters:
//   - M: The type of the message.
//   - T: The type of the actor state.
//
// Parameters:
//   - msg: The message of type T to be processed.
//   - currenttate: The current state of the actor of type S.
//
// Returns:
//   - Payload[T]: The updated state of the actor after processing the message.
//   - error: An error if the processing fails, otherwise nil.
type ProcessingFn[M any, T any] func(msg Message[M], currentState Payload[T]) (Payload[T], error)
