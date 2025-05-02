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
//   - initialState (ActorState): The initial state of the actor.
//
// Returns:
//   - (Actor): The created Actor instance.
//   - (error): An error if the actor could not be created.
func NewActor[T any](
	address url.URL,
	processingFn f.ProcessingFn[T],
	initialState f.ActorState[T],
) (f.Actor[T], error) {
	return i.NewActor(address, processingFn, initialState)
}
