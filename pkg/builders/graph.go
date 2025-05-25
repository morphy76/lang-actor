package builders

import (
	"github.com/google/uuid"
	g "github.com/morphy76/lang-actor/internal/graph"
	"github.com/morphy76/lang-actor/pkg/graph"
)

// NewGraph creates a new instance of the actor graph.
//
// Parameters:
//   - initialState (map[string]any): The initial state of the graph.
//   - configs (map[string]any): Optional configurations for the graph.
//
// Returns:
//   - (graph.Graph): The created actor graph.
//   - (error): An error if the graph creation fails.
func NewGraph(
	initialState map[string]any,
	configs ...map[string]any,
) (graph.Graph, error) {
	mergedConfigs := make(map[string]any)
	for _, config := range configs {
		for key, value := range config {
			mergedConfigs[key] = value
		}
	}
	return g.NewGraph(uuid.NewString(), initialState, mergedConfigs)
}

// NewRootNode creates a new instance of the root node.
//
// Parameters:
//   - forGraph (graph.Graph): The graph to which the root node belongs.
//
// Returns:
//   - (graph.Node): The created root node.
//   - (error): An error if the node creation fails.
func NewRootNode(forGraph graph.Graph) (graph.Node, error) {
	return g.NewRootNode(forGraph)
}

// NewDebugNode creates a new instance of the debug node.
//
// Parameters:
//   - forGraph (graph.Graph): The graph to which the root node belongs.
//
// Returns:
//   - (graph.Node): The created debug node.
//   - (error): An error if the node creation fails.
func NewDebugNode(forGraph graph.Graph) (graph.Node, error) {
	return g.NewDebugNode(forGraph)
}

// NewEndNode creates a new instance of the end node.
//
// Parameters:
//   - forGraph (graph.Graph): The graph to which the root node belongs.
//
// Returns:
//   - (graph.Node): The created end node.
//   - (error): An error if the node creation fails.
func NewEndNode(forGraph graph.Graph) (graph.Node, error) {
	return g.NewEndNode(forGraph)
}
