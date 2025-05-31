package graph

import (
	"net/url"

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

	graph := &graph{
		resolvables: make(map[url.URL]*c.Addressable),
		graphURL:    *graphURL,
		config:      config,
		status:      initialState,
		addressBook: routing.NewAddressBook(),
	}

	return graph, nil
}
