package builders

import (
	"net/url"

	g "github.com/morphy76/lang-actor/internal/graph"
	"github.com/morphy76/lang-actor/pkg/framework"
	"github.com/morphy76/lang-actor/pkg/graph"
)

// NewCustomNode creates a new instance of a custom node.
//
// Parameters:
//   - forGraph (graph.Graph): The graph to which the custom node belongs.
//   - address (*url.URL): The URL address of the node.
//   - taskFn (framework.ProcessingFn[NodeState]): The processing function for the node.
//   - transient (bool): Whether the node is transient or not.
//
// Returns:
//   - (graph.Node): The created custom node.
//   - (error): An error if the node creation fails.
func NewCustomNode(
	forGraph graph.Graph,
	address *url.URL,
	taskFn framework.ProcessingFn[graph.NodeState],
	transient bool,
) (graph.Node, error) {
	return g.NewNode(forGraph, *address, taskFn, transient)
}
