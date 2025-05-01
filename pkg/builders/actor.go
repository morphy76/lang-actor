package builders

import (
	"context"
	"errors"
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
//   - parentCtx (context.Context): The parent context for the actor.
//   - address (url.URL): The URL address that specifies the actor's location and protocol.
//   - processingFn (ProcessingFn): The function to process messages sent to the actor.
//   - initialState (ActorState): The initial state of the actor.
//
// Returns:
//   - (Actor): The created Actor instance.
//   - (error): An error if the actor could not be created.
func NewActor[T any](
	parentCtx context.Context,
	address url.URL,
	processingFn f.ProcessingFn[T],
	initialState f.ActorState[T],
) (f.Actor[T], error) {
	rv, err := i.NewActor(parentCtx, address, processingFn, initialState)
	if err != nil {
		return nil, err
	}
	// TODO a better model for the catalog
	actorCatalog, found := parentCtx.Value(f.ActorCatalogContextKey).(map[url.URL]f.Transport)
	if !found {
		return nil, f.ErrActorCatalogNotFound
	}
	if _, found := actorCatalog[address]; found {
		return nil, errors.New("actor already exists")
	}
	actorCatalog[address] = rv.(f.Transport)
	return rv, nil
}
