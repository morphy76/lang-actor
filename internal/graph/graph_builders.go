package graph

import (
	g "github.com/morphy76/lang-actor/pkg/graph"
)

// NewGraph creates a new instance of the actor graph.
func NewGraph(
	rootNode g.RootNode,
) (g.Graph, error) {

	graph := &graph{
		rootNode: rootNode,
	}

	return graph, nil
}
