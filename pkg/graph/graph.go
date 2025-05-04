package graph

// Node represents a node in the actor graph.
type Node interface {
	// Append adds a new node to the graph.
	//
	// Parameters:
	//   - node: The node to append.
	Append(node Node)
}

// RootNode represents the root node of the actor graph.
type RootNode interface {
	Node
}

// EndNode represents an end node in the actor graph.
type EndNode interface {
	Node
}
