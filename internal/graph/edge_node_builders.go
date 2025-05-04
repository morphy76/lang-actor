package graph

import (
	g "github.com/morphy76/lang-actor/pkg/graph"
)

// NewRootNode creates a new instance of the actor graph.
func NewRootNode() g.RootNode {
	return &rootNode{}
}

// NewEndNode creates a new instance of the end node in the actor graph.
func NewEndNode() g.EndNode {
	return &endNode{}
}
