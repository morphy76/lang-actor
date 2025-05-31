package graph

var staticNodeStateAssertion NodeState = (*BasicNodeState)(nil)

// BasicNodeStateBuilder creates a new instance of BasicNodeState for the specified graph.
//
// Type Parameters:
//   - T (NodeState): The type of the node state to be created.
//
// Parameters:
//   - forGraph (Graph): The graph to which the node state belongs.
//   - owner (Node): The owner node of the state.
//   - outcome (chan string): A channel to receive the outcome of the node's processing.
//
// Returns:
//   - (NodeState): A new instance of BasicNodeState.
func BasicNodeStateBuilder[T NodeState](forGraph Graph, owner Node, outcome chan string) T {
	var rv NodeState = &BasicNodeState{
		outcome: outcome,
		graph:   forGraph,
		owner:   owner,
	}
	return rv.(T)
}

// BasicNodeState implements the NodeState interface for managing the state of a graph node.
type BasicNodeState struct {
	outcome chan string
	graph   Graph
	owner   Node
}

// Outcome returns the outcome channel for the node state.
//
// Returns:
//   - (chan string): A channel that receives the outcome of the node's processing.
func (ns *BasicNodeState) Outcome() chan string {
	return ns.outcome
}

// GraphConfig returns the configuration of the graph associated with the node state.
//
// Returns:
//   - (Configuration): The configuration of the graph, or nil if the graph is not set.
func (ns *BasicNodeState) GraphConfig() Configuration {
	if ns.graph == nil {
		return nil
	}
	return ns.graph.Config()
}

// GraphState returns the state of the graph associated with the node state.
//
// Returns:
//   - (State): The state of the graph, or nil if the graph is not set.
func (ns *BasicNodeState) GraphState() State {
	if ns.graph == nil {
		return nil
	}
	return ns.graph.State()
}

// UpdateGraphState initializes the state for the node state.
//
// Parameters:
//   - state (State): The new state to set for the
//
// Returns:
//   - (error): An error if the update fails, nil otherwise.
func (ns *BasicNodeState) UpdateGraphState(state State) error {
	return ns.graph.UpdateState(state)
}

// Routes returns the routes of the node state.
func (ns *BasicNodeState) Routes() []string {
	return ns.owner.EdgeNames()
}
