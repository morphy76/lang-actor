package builders

import (
	"github.com/google/uuid"
	"github.com/morphy76/lang-actor/internal/graph"
	g "github.com/morphy76/lang-actor/pkg/graph"
)

// NewGraph creates a new instance of the actor graph.
func NewGraph(
	rootNode g.RootNode,
	configs ...map[string]any,
) (g.Graph, error) {
	mergedConfigs := make(map[string]any)
	for _, config := range configs {
		for key, value := range config {
			mergedConfigs[key] = value
		}
	}
	return graph.NewGraph(uuid.NewString(), rootNode, mergedConfigs)
}

// NewRootNode creates a new instance of the root node.
func NewRootNode() (g.Node, error) {
	return graph.NewRootNode()
}

// NewDebugNode creates a new instance of the debug node.
func NewDebugNode() (g.Node, error) {
	return graph.NewDebugNode()
}

// NewEndNode creates a new instance of the end node.
func NewEndNode() (g.Node, chan bool, error) {
	return graph.NewEndNode()
}
