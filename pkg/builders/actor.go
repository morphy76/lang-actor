package builders

import (
	"net/url"

	f "github.com/morphy76/lang-actor/internal/framework"
	"github.com/morphy76/lang-actor/pkg/framework"
)

// NewActor creates a new actor with the given address. The actor state is always updated by message processing.
// Supported schemas are:
// - actor://
//
// Type Parameters:
//   - T: The type of the actor state.
//
// Parameters:
//   - address (url.URL): The address of the actor.
//   - processingFn (framework.ProcessingFn): The function to process messages sent to the actor.
//   - initialState (T): The initial state of the actor.
//   - mailboxConfig (framework.MailboxConfig): Optional configuration for the actor's mailbox.
//
// Returns:
//   - (framework.Actor): The created Actor instance.
//   - (error): An error if the actor could not be created.
func NewActor[T any](
	address url.URL,
	processingFn framework.ProcessingFn[T],
	initialState T,
	mailboxConfig ...framework.MailboxConfig,
) (framework.Actor[T], error) {
	return f.NewActor(address, processingFn, initialState, true, mailboxConfig...)
}

// NewMutableActor creates a new actor with the given address alwways mutable.
// Supported schemas are:
// - actor://
//
// Type Parameters:
//   - T: The type of the actor state.
//
// Parameters:
//   - address (url.URL): The address of the actor.
//   - processingFn (framework.ProcessingFn): The function to process messages sent to the actor.
//   - initialState (T): The initial state of the actor.
//   - mailboxConfig (framework.MailboxConfig): Optional configuration for the actor's mailbox.
//
// Returns:
//   - (framework.Actor): The created Actor instance.
//   - (error): An error if the actor could not be created.
func NewMutableActor[T any](
	address url.URL,
	processingFn framework.ProcessingFn[T],
	initialState T,
	mailboxConfig ...framework.MailboxConfig,
) (framework.Actor[T], error) {
	return f.NewActor(address, processingFn, initialState, false, mailboxConfig...)
}

// SpawnChild creates a new child actor with the given processing function and initial state.
//
// Type Parameters:
//   - T: The type of the actor state.
//
// Parameters:
//   - parent (framework.ActorRef): The reference to the parent actor.
//   - processingFn (framework.ProcessingFn): The function to process messages sent to the child actor.
//   - initialState (T): The initial state of the child actor.
//   - mailboxConfig (framework.MailboxConfig): Optional configuration for the child's mailbox.
//
// Returns:
//   - (framework.Actor): The created child Actor instance.
//   - (error): An error if the child actor could not be created.
func SpawnChild[T any](
	parent framework.ActorRef,
	processingFn framework.ProcessingFn[T],
	initialState T,
	mailboxConfig ...framework.MailboxConfig,
) (framework.Actor[T], error) {
	child, err := f.NewActorWithParent(processingFn, initialState, true, parent, mailboxConfig...)
	if err != nil {
		return nil, err
	}

	return child, nil
}
