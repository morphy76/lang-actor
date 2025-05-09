package builders

import (
	"net/url"

	i "github.com/morphy76/lang-actor/internal/framework"
	f "github.com/morphy76/lang-actor/pkg/framework"
)

// NewActor creates a new actor with the given address.
// Supported schemas are:
// - actor://
//
// Type Parameters:
//   - T: The type of the actor state.
//
// Parameters:
//   - address (url.URL): The URL address that specifies the actor's location and protocol.
//   - processingFn (ProcessingFn): The function to process messages sent to the actor.
//   - initialState (T): The initial state of the actor.
//   - mailboxConfig (MailboxConfig): Optional configuration for the actor's mailbox.
//
// Returns:
//   - (Actor): The created Actor instance.
//   - (error): An error if the actor could not be created.
func NewActor[T any](
	address url.URL,
	processingFn f.ProcessingFn[T],
	initialState T,
	mailboxConfig ...f.MailboxConfig,
) (f.Actor[T], error) {
	return i.NewActor(address, processingFn, initialState, mailboxConfig...)
}

// SpawnChild creates a new child actor with the given processing function and initial state.
//
// Type Parameters:
//   - T: The type of the actor state.
//
// Parameters:
//   - parent (Actor): The parent actor that will spawn the child actor.
//   - processingFn (ProcessingFn): The function to process messages sent to the child actor.
//   - initialState (T): The initial state of the child actor.
//   - mailboxConfig (MailboxConfig): Optional configuration for the child actor's mailbox.
//
// Returns:
//   - (Actor): The created child Actor instance.
//   - (error): An error if the child actor could not be created.
func SpawnChild[T any](
	parent f.ActorRef,
	processingFn f.ProcessingFn[T],
	initialState T,
	mailboxConfig ...f.MailboxConfig,
) (f.Actor[T], error) {
	child, err := i.NewActorWithParent(processingFn, initialState, parent, mailboxConfig...)
	if err != nil {
		return nil, err
	}

	return child, nil
}
