package graph

import (
	"errors"
	"fmt"

	g "github.com/morphy76/lang-actor/pkg/graph"
)

var staticDebugAssertion g.DebugNode = (*debugNode)(nil)

type debugNode struct {
	node
}

// OneWayRoute adds a new possible outgoing route from the node
func (r *debugNode) OneWayRoute(name string, destination g.Node) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if len(r.edges) > 0 {
		return errors.Join(g.ErrorInvalidRouting, fmt.Errorf("debugNode node [%s] already has a route", r.Name()))
	}

	r.edges[name] = edge{
		Name:        name,
		Destination: destination,
	}

	return nil
}
