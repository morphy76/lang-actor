package graph

import (
	"sync"

	"github.com/morphy76/lang-actor/pkg/graph"
)

// NewNode creates a new instance of a node in the actor graph.
func NewNode() graph.Node {
	return &node{
		lock:  &sync.Mutex{},
		nodes: make([]graph.Node, 0),
	}
}
