package graph

import (
	"errors"
	"fmt"

	c "github.com/morphy76/lang-actor/pkg/common"
	g "github.com/morphy76/lang-actor/pkg/graph"
)

var staticRootNodeAssertion g.RootNode = (*rootNode)(nil)
var staticEndAssertion g.EndNode = (*endNode)(nil)

type rootNode struct {
	node
}

// OneWayRoute adds a new possible outgoing route from the node
func (r *rootNode) OneWayRoute(name string, destination g.Node) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if len(r.edges) > 0 {
		return errors.Join(g.ErrorInvalidRouting, fmt.Errorf("root node [%v] already has a route", r.Address()))
	}

	r.edges[name] = edge{
		Name:        name,
		Destination: destination,
	}

	return nil
}

type endNode struct {
	node
}

// OneWayRoute adds a new possible outgoing route from the node
func (r *endNode) OneWayRoute(name string, destination g.Node) error {
	return errors.Join(g.ErrorInvalidRouting, fmt.Errorf("cannot route [%v] from the end [%v]", name, r.Address()))
}

// ProceedOnAnyRoute proceeds with the first route available
func (r *endNode) ProceedOnAnyRoute(mex c.Message) error {
	return errors.Join(g.ErrorInvalidRouting, fmt.Errorf("cannot route from the end [%v]", r.Address()))
}
