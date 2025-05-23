package builders

import (
	"github.com/google/uuid"
	g "github.com/morphy76/lang-actor/internal/graph"
	"github.com/morphy76/lang-actor/pkg/graph"
)

// NewGraph creates a new instance of the actor graph.
//
// Type parameters:
//   - T: The type of the root node.
//
// Parameters:
//   - rootNode (graph.RootNode): The root node of the graph.
//   - configs (map[string]any): Optional configurations for the graph.
//
// Returns:
//   - (graph.Graph): The created actor graph.
//   - (error): An error if the graph creation fails.
func NewGraph[T any](
	rootNode graph.RootNode,
	initialStatus T,
	configs ...map[string]any,
) (graph.Graph, error) {
	mergedConfigs := make(map[string]any)
	for _, config := range configs {
		for key, value := range config {
			mergedConfigs[key] = value
		}
	}
	return g.NewGraph(uuid.NewString(), rootNode, initialStatus, mergedConfigs)
}

// NewRootNode creates a new instance of the root node.
//
// Returns:
//   - (graph.Node): The created root node.
//   - (error): An error if the node creation fails.
func NewRootNode() (graph.Node, error) {
	return g.NewRootNode()
}

// NewDebugNode creates a new instance of the debug node.
//
// Returns:
//   - (graph.Node): The created debug node.
//   - (error): An error if the node creation fails.
func NewDebugNode() (graph.Node, error) {
	return g.NewDebugNode()
}

// NewEndNode creates a new instance of the end node.
//
// Returns:
//   - (graph.Node): The created end node.
//   - (chan bool): A channel to signal the end of processing.
//   - (error): An error if the node creation fails.
func NewEndNode() (graph.Node, chan bool, error) {
	return g.NewEndNode()
}
