package framework

import (
	"errors"
	"net/url"

	"github.com/morphy76/lang-actor/pkg/common"
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

// BackpressurePolicy defines how an actor's mailbox handles pressure when reaching capacity
type BackpressurePolicy int

const (
	// BackpressurePolicyBlock causes message sends to wait when the mailbox is full
	BackpressurePolicyBlock BackpressurePolicy = iota

	// BackpressurePolicyFail causes message sends to fail immediately when the mailbox is full
	BackpressurePolicyFail

	// BackpressurePolicyUnbounded means the mailbox has no capacity limit
	BackpressurePolicyUnbounded

	// BackpressurePolicyDropNewest rejects new messages when the mailbox is full
	BackpressurePolicyDropNewest

	// BackpressurePolicyDropOldest discards oldest messages to make room for new ones
	BackpressurePolicyDropOldest
)

// MailboxConfig defines configuration options for an actor's mailbox
type MailboxConfig struct {
	// Capacity defines the maximum number of messages the mailbox can hold
	// Ignored when using BackpressurePolicyUnbounded
	Capacity int
	// Policy defines how the mailbox handles pressure when reaching capacity
	Policy BackpressurePolicy
}

// Addressable is the interface for the routing layer of the actor model.
type Addressable interface {
	// Actor URI
	//
	// Returns:
	//   - (url.URL): The URL of the actor.
	Address() url.URL
	// Deliver a message to the actor
	//
	// Parameters:
	//   - msg (Message): The message to be delivered.
	//
	// Returns:
	//   - (error): An error if the delivery fails, otherwise nil.
	Deliver(msg Message) error
	// Send is a function to send messages to other actors.
	//
	// Parameters:
	//   - msg (Message): The message to be sent.
	//   - addressable (Addressable): The addressable actor to which the message is sent.
	//
	// Returns:
	//   - (error): An error if the sending fails, otherwise nil.
	Send(msg Message, addressable Addressable) error
}

// Controllable is the interface for controllable actors.
type Controllable interface {
	// Stop the actor
	//
	// Returns:
	//   - (chan bool): A channel that is closed when the actor is stopped.
	//   - (error): An error if the stopping fails, otherwise nil.
	Stop() (chan bool, error)
	// Status of the actor
	//
	// Returns:
	//   - (ActorStatus): The status of the actor.
	Status() ActorStatus
}

// Controller is the interface for the controller of the actor model.
type Controller interface {
	// Append a child actor
	//
	// Parameters:
	//   - child (ActorRef): The child actor to be appended.
	//
	// Returns:
	//   - (error): An error if the appending fails, otherwise nil.
	Append(child ActorRef) error
	// Crop a child actor
	//
	// Parameters:
	//   - child (url.URL): The URL of the child actor to be cropped.
	//
	// Returns:
	//   - (ActorRef): The cropped child actor.
	//   - (error): An error if the cropping fails, otherwise nil.
	Crop(url.URL) (ActorRef, error)
}

// Relationable is the interface for the relationable actors.
type Relationable interface {
	// GetParent returns the parent actor of the current actor.
	//
	// Returns:
	//   - (ActorRef): The parent actor.
	//   - (bool): A boolean indicating whether the parent actor exists.
	GetParent() (ActorRef, bool)
}

// ActorRef is the interface for the actor reference.
type ActorRef interface {
	common.Visitable
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
	//
	// Returns:
	//   - (T): The state of the actor.
	State() T
}

// Message is the interface for messages sent to actors.
type Message interface {
	// Sender returns the URL of the sender.
	//
	// Returns:
	//   - (url.URL): The URL of the sender.
	Sender() url.URL
	// Mutation returns true if the message is an actor mutation.
	//
	// Returns:
	//   - (bool): A boolean indicating whether the message is a mutation.
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
//   - msg: The message of type M to be processed.
//   - self: The actor instance that is processing the message.
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
// Parameters:
//   - msg: The message of type M to be sent.
//   - destination: The URL of the destination actor to which the message is sent.
//
// Returns:
//   - error: An error if the sending fails, otherwise nil.
type SendFn func(msg Message, destination url.URL) error
