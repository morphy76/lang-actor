package graph

var staticNodeRefAssertion NodeRef = (*BasicNodeRef)(nil)

// BasicNodeRefBuilder creates a new instance of BasicNodeRef for the specified graph.
//
// Type Parameters:
//   - T (NodeRef): The type of the node state to be created.
//
// Parameters:
//   - forGraph (Graph): The graph to which the node state belongs.
//   - owner (Node): The owner node of the state.
//   - routeFollow (chan string): A channel to receive the route to follow for processing messages.
//   - attrs (...map[string]any): Optional attributes for the node state.
//
// Returns:
//   - (NodeRef): A new instance of BasicNodeRef.
func BasicNodeRefBuilder[T NodeRef](forGraph Graph, owner Node, routeFollow chan string, attrs ...map[string]any) T {
	var rv NodeRef = &BasicNodeRef{
		routeFollow: routeFollow,
		graph:       forGraph,
		owner:       owner,
		attrs:       make(map[string]any),
	}
	for _, attr := range attrs {
		for k, v := range attr {
			rv.(*BasicNodeRef).attrs[k] = v
		}
	}
	return rv.(T)
}

// BasicNodeRef implements the NodeRef interface for managing the state of a graph node.
type BasicNodeRef struct {
	routeFollow chan string
	graph       Graph
	owner       Node
	attrs       map[string]any
}

// ProceedOntoRoute returns the outcome channel for the node state.
//
// Returns:
//   - (chan string): A channel that receives the route to follow.
func (ns *BasicNodeRef) ProceedOntoRoute() chan string {
	return ns.routeFollow
}

// GraphConfig returns the configuration of the graph associated with the node state.
//
// Returns:
//   - (Configuration): The configuration of the graph, or nil if the graph is not set.
func (ns *BasicNodeRef) GraphConfig() Configuration {
	if ns.graph == nil {
		return nil
	}
	return ns.graph.Config()
}

// GraphState returns the state of the graph associated with the node state.
//
// Returns:
//   - (State): The state of the graph, or nil if the graph is not set.
func (ns *BasicNodeRef) GraphState() State {
	if ns.graph == nil {
		return nil
	}
	return ns.graph.State()
}

// Routes returns the routes of the node state.
func (ns *BasicNodeRef) Routes() []string {
	return ns.owner.EdgeNames()
}

// SetAttribute sets an attribute for the node state.
func (ns *BasicNodeRef) SetAttribute(key string, value any) {
	ns.attrs[key] = value
}

// GetAttribute retrieves an attribute for the node state.
func (ns *BasicNodeRef) GetAttribute(key string) (any, bool) {
	val, found := ns.attrs[key]
	return val, found
}
