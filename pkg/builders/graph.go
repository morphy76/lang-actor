package builders

import (
	g "github.com/morphy76/lang-actor/internal/graph"
	"github.com/morphy76/lang-actor/pkg/graph"
)

// StartGraph starts the actor graph.
func StartGraph() (graph.RootNode, error) {
	return g.NewRootNode(), nil
}
