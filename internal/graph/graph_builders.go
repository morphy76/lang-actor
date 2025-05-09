package graph

import (
	"net/url"

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

	configURL, err := url.Parse("graph://" + graphName + "/config")
	if err != nil {
		return nil, err
	}

	configNode, err := newConfigNode(configs, *configURL)
	if err != nil {
		return nil, err
	}

	graph := &graph{
		graphURL:   *graphURL,
		rootNode:   rootNode,
		configNode: configNode,
	}

	return graph, nil
}
