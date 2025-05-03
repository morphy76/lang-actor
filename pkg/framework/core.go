package framework

import (
	"errors"
	"net/url"
)

// ErrorInvalidActorAddress is returned when an actor address is invalid.
var ErrorInvalidActorAddress = errors.New("invalid actor address")

// ErrorActorNotRunning is returned when an actor is not running.
var ErrorActorNotRunning = errors.New("actor not running")

// ErrorInvalidChildURL is returned when a child URL is invalid.
var ErrorInvalidChildURL = errors.New("invalid child URL")

// ActorStatus represents the status of an actor.
type ActorStatus int8

const (
	// ActorStatusIdle indicates that the actor is idle.
	ActorStatusIdle ActorStatus = iota
	// ActorStatusRunning indicates that the actor is running.
	ActorStatusRunning
)

// Addressable is the interface for the routing layer of the actor model.
type Addressable interface {
	// Actor URI
	Address() url.URL
	// Deliver a message to the actor
	Deliver(msg Message) error
	// Send is a function to send messages to other actors.
	Send(msg Message, addressable Addressable) error
}

// Controllable is the interface for controllable actors.
type Controllable interface {
	// Stop the actor
	Stop() (chan bool, error)
	// Status of the actor
	Status() ActorStatus
}

// Controller is the interface for the controller of the actor model.
type Controller interface {
	// Append a child actor
	Append(child ActorRef) error
	// Crop a child actor
	Crop(url.URL) (ActorRef, error)
}

// Relationable is the interface for the relationable actors.
type Relationable interface {
	// GetParent returns the parent actor of the current actor.
	GetParent() (ActorRef, bool)
}

// ActorRef is the interface for the actor reference.
type ActorRef interface {
	Controllable
	Controller
	Addressable
	Relationable
}

// Actor is part of the actor model framework underlying lang-actor.
//
// Type Parameters:
//   - T: The type of the actor state.
type Actor[T any] interface {
	ActorRef
	// State of the actor
	State() T
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
	self Actor[T],
) (T, error)

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
