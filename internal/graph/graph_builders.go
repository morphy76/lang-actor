package graph

import (
	"net/url"
	"sync"

	"github.com/morphy76/lang-actor/internal/routing"
	c "github.com/morphy76/lang-actor/pkg/common"
	g "github.com/morphy76/lang-actor/pkg/graph"
)

// NewGraph creates a new instance of the actor graph.
func NewGraph[T g.State, C g.Configuration](
	graphName string,
	initialState T,
	config C,
) (g.Graph, error) {

	graphURL, err := url.Parse("graph://" + graphName)
	if err != nil {
		return nil, err
	}

	stateChangedCh := make(chan g.State, 100000000) // Buffered channel to avoid blocking
	useState := NewStateWrapper(initialState, stateChangedCh)

	graph := &graph{
		lock: &sync.Mutex{},

		resolvables: make(map[url.URL]*c.Addressable),
		graphURL:    *graphURL,
		config:      config,
		state:       useState,
		addressBook: routing.NewAddressBook(),

		stateChangedCh: stateChangedCh,
	}

	return graph, nil
}
