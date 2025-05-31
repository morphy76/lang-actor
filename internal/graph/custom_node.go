package graph

import (
	"net/url"

	f "github.com/morphy76/lang-actor/pkg/framework"
	g "github.com/morphy76/lang-actor/pkg/graph"
)

// NewCustomNode creates a new instance of a custom node.
func NewCustomNode(
	forGraph g.Graph,
	address *url.URL,
	taskFn f.ProcessingFn[g.NodeState],
	transient bool,
) (g.Node, error) {
	return newNode(forGraph, *address, taskFn, transient)
}
