package graph

import (
	"sync"

	g "github.com/morphy76/lang-actor/pkg/graph"
)

var staticNodeAssertion g.Node = (*node)(nil)

type node struct {
	lock *sync.Mutex

	nodes []g.Node
}

// Append adds a new node to the graph.
func (n *node) Append(node g.Node) {
	n.lock.Lock()
	defer n.lock.Unlock()
	n.nodes = append(n.nodes, node)
}
