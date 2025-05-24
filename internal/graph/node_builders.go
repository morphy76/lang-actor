package graph

import (
	"net/url"
	"sync"

	"github.com/morphy76/lang-actor/pkg/framework"
)

// NewNode creates a new instance of a node in the graph with the given address.
func NewNode(address url.URL) *node {
	return &node{
		lock:    &sync.Mutex{},
		edges:   make(map[string]edge, 0),
		address: address,
	}
}

// NewNodeWithActor creates a new instance of a node in the graph with the given address and actor.
func NewNodeWithActor[T any](address url.URL, actor framework.Actor[T]) *node {
	return &node{
		lock:    &sync.Mutex{},
		edges:   make(map[string]edge, 0),
		address: address,
		actor:   actor,
	}
}
