package graph

// Graph represents the actor, runnable, graph.
type Graph interface {
	// TODO Accept accepts a todo item.
	Accept(todo any) error
}
