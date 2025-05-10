package graph

import (
	"net/url"

	"github.com/morphy76/lang-actor/internal/routing"

	c "github.com/morphy76/lang-actor/pkg/common"
	f "github.com/morphy76/lang-actor/pkg/framework"
	g "github.com/morphy76/lang-actor/pkg/graph"
)

// NewGraph creates a new instance of the actor graph.
func NewGraph(
	graphName string,
	rootNode g.RootNode,
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

	graph := &graph{
		resolvables: make(map[url.URL]*f.Addressable),
		graphURL:    *graphURL,
		rootNode:    rootNode,
		configNode:  configNode,
		addressBook: routing.NewAddressBook(),
	}

	var registerFn c.VisitFn = func(visitable c.Visitable) {
		node, ok := visitable.(g.Node)
		if !ok {
			return
		}

		graph.Register(node)
		node.SetResolver(graph)
	}
	rootNode.Visit(registerFn)
	configNode.Visit(registerFn)
	configNode.SetResolver(graph)

	// TODO it should send an init message to the root node which propagates to all its edges recursively

	return graph, nil
}
