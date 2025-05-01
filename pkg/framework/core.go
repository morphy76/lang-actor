package framework

import (
	"errors"
	"net/url"
)

// ErrActorCatalogNotFound is returned when the actor catalog is not found.
var ErrActorCatalogNotFound = errors.New("actor catalog not found")

// ErrorInvalidActorAddress is returned when an actor address is invalid.
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
//   - T: The type of the actor state.
type Actor[T any] interface {
	ActorView[T]
	// Start the actor
	Start() error
	// Stop the actor
	Stop() (chan bool, error)
	// Status of the actor
	Status() ActorStatus
}

// ActorView is the interface for the actor view.
//
// Type Parameters:
//   - T: The type of the actor state.
type ActorView[T any] interface {
	Transport
	// State of the actor
	State() T
	// Send a message to another actor
	TransportByAddress(address url.URL) (Transport, error)
}

// Transport is the interface for the transport layer of the actor model.
type Transport interface {
	Addressable
	// Deliver a message to the actor
	Deliver(msg Message) error
	// Send is a function to send messages to other actors.
	Send(msg Message, transport Transport) error
}

// Addressable is the interface for addressable entities.
type Addressable interface {
	// Actor URI
	Address() url.URL
}

// ActorState represents the state of an actor.
//
// Type Parameters:
//   - T: The type of the actor state.
type ActorState[T any] interface {
	// Cast returns the exact type of the struct implementing this interface.
	Cast() T
}

// Message is the interface for messages sent to actors.
//
// Type Parameters:
//   - T: The type of the message.
type Message interface {
	// Sender returns the URL of the sender.
	Sender() url.URL
	// Mutation returns true if the message is an actor mutation.
	Mutation() bool
}

// ProcessingFn defines a generic function type for processing messages within an actor system.
//
// It takes a message of type M and the current actor state of type T, and returns the updated
// state of type T along with an error, if any.
//
// If the message is a mutation, the state is updated and returned. If the message is not a mutation,
// the state remains unchanged regardless by what's returned by the processing funcion.
//
// Type Parameters:
//   - T: The type of the actor state.
//
// Parameters:
//   - msg: The message of type T to be processed.
//   - self: The actor view of type T that is processing the message.
//
// Returns:
//   - Payload[T]: The updated state of the actor after processing the message.
//   - error: An error if the processing fails, otherwise nil.
type ProcessingFn[T any] func(
	msg Message,
	self ActorView[T],
) (ActorState[T], error)

// SendFn defines a function type for sending messages to actors.
//
// It takes a message of type M and returns an error if the sending fails.
//
// This function is typically used within the actor system to send messages
// between actors.
//
// Type Parameters:
//   - M: The type of the message.
//
// Parameters:
//   - msg: The message of type M to be sent.
//   - destination: The URL of the destination actor to which the message is sent.
//
// Returns:
//   - error: An error if the sending fails, otherwise nil.
type SendFn func(msg Message, destination url.URL) error
