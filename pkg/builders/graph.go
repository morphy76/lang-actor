package builders

import (
	"github.com/morphy76/lang-actor/internal/graph"
	g "github.com/morphy76/lang-actor/pkg/graph"
)

// NewGraph creates a new instance of the actor graph.
func NewGraph(
	rootNode g.RootNode,
) (g.Graph, error) {
	return graph.NewGraph(rootNode)
}
