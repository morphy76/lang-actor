package graph

import (
	"net/url"

	"github.com/morphy76/lang-actor/internal/routing"

	c "github.com/morphy76/lang-actor/pkg/common"
	f "github.com/morphy76/lang-actor/pkg/framework"
	g "github.com/morphy76/lang-actor/pkg/graph"
)

// NewGraph creates a new instance of the actor graph.
func NewGraph[T any](
	graphName string,
	rootNode g.RootNode,
	initialStatus T,
	configs map[string]any,
) (g.Graph, error) {

	graphURL, err := url.Parse("graph://" + graphName)
	if err != nil {
		return nil, err
	}

	configNode, err := NewConfigNode(configs, graphName)
	if err != nil {
		return nil, err
	}

	statusNode, err := NewStatusNode(initialStatus, graphName)
	if err != nil {
		return nil, err
	}

	graph := &graph{
		resolvables: make(map[url.URL]*f.Addressable),
		graphURL:    *graphURL,
		rootNode:    rootNode,
		configNode:  configNode,
		statusNode:  statusNode,
		addressBook: routing.NewAddressBook(),
	}

	var registerFn c.VisitFn = func(visitable c.Visitable) {
		addressable, ok := visitable.(f.Addressable)
		if !ok {
			return
		}

		graph.Register(addressable)

		routable, ok := visitable.(g.Routable)
		if ok {
			routable.SetResolver(graph)
		}
	}
	rootNode.Visit(registerFn)
	configNode.Visit(registerFn)
	configNode.SetResolver(graph)
	statusNode.Visit(registerFn)
	statusNode.SetResolver(graph)

	// TODO it should send an init message to the root node which propagates to all its edges recursively

	return graph, nil
}
