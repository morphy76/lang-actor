package graph

import (
	r "github.com/morphy76/lang-actor/pkg/routing"
)

// Graph represents the actor, runnable, graph.
type Graph interface {
	r.Resolver
	// TODO Accept accepts a todo item.
	Accept(todo any) error
}
