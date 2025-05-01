package builders

import (
	"context"
	"net/url"

	i "github.com/morphy76/lang-actor/internal/framework"
	f "github.com/morphy76/lang-actor/pkg/framework"
	r "github.com/morphy76/lang-actor/pkg/routing"
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

	actorCatalog, found := parentCtx.Value(r.ActorCatalogContextKey).(r.Catalog)
	if !found {
		return nil, f.ErrActorCatalogNotFound
	}

	err = actorCatalog.Register(rv)
	if err != nil {
		return nil, err
	}

	return rv, nil
}
