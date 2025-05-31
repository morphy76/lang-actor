package builders

import (
	"net/url"

	"github.com/google/uuid"

	g "github.com/morphy76/lang-actor/internal/graph"
	"github.com/morphy76/lang-actor/pkg/framework"
	"github.com/morphy76/lang-actor/pkg/graph"
)

// NewGraph creates a new instance of the actor graph.
//
// Type parameters:
//   - T (*any): The type of the initial state of the graph.
//   - C (any): The type of the configurations for the graph.
//
// Parameters:
//   - initialState (T graph.State): The initial state of the graph.
//   - configs (C graph.Configuration): Optional configurations for the graph.
//
// Returns:
//   - (graph.Graph): The created actor graph.
//   - (error): An error if the graph creation fails.
func NewGraph[T graph.State, C graph.Configuration](
	initialState T,
	configs C,
) (graph.Graph, error) {
	return g.NewGraph(uuid.NewString(), initialState, configs)
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
//   - nameParts (...string): Optional parts of the name for the debug node.
//
// Returns:
//   - (graph.Node): The created debug node.
//   - (error): An error if the node creation fails.
func NewDebugNode(forGraph graph.Graph, nameParts ...string) (graph.Node, error) {
	return g.NewDebugNode(forGraph, nameParts...)
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
	return g.NewCustomNode(forGraph, address, taskFn, transient)
}

// NewForkJoingNode creates a new fork-join node for the given graph.
//
// Type parameters:
//   - C (graph.NodeState): The type of the state for the child nodes.
//
// Parameters:
//   - forGraph (graph.Graph): The graph to which the fork-join node belongs.
//   - transient (bool): Whether the node is transient or not.
//   - taskFn (...framework.ProcessingFn[graph.NodeState]): Optional processing functions for the node.
//
// Returns:
//   - (graph.Node): The created fork-join node.
//   - (error): An error if the node creation fails.
func NewForkJoingNode[T graph.NodeState, C graph.NodeState](
	forGraph graph.Graph,
	transient bool,
	taskFn ...framework.ProcessingFn[C],
) (graph.Node, error) {
	return g.NewForkJoingNode[C](forGraph, transient, taskFn...)
}
