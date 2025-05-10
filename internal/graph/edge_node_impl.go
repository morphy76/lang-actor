package graph

import (
	"errors"
	"fmt"

	f "github.com/morphy76/lang-actor/pkg/framework"
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
		return errors.Join(g.ErrorInvalidRouting, fmt.Errorf("root node [%s] already has a route", r.Name()))
	}

	r.edges[name] = edge{
		Name:        name,
		Destination: destination,
	}

	return nil
}

// TwoWayRoute adds a new possible outgoing route from the node
func (r *rootNode) TwoWayRoute(name string, destination g.Node) error {
	return errors.Join(g.ErrorInvalidRouting, fmt.Errorf("root node [%s] cannot have a two way route", r.Name()))
}

type endNode struct {
	node
}

// OneWayRoute adds a new possible outgoing route from the node
func (r *endNode) OneWayRoute(name string, destination g.Node) error {
	return errors.Join(g.ErrorInvalidRouting, fmt.Errorf("cannot route [%s] from the end [%s]", name, r.Name()))
}

// TwoWayRoute adds a new possible outgoing route from the node
func (r *endNode) TwoWayRoute(name string, destination g.Node) error {
	return errors.Join(g.ErrorInvalidRouting, fmt.Errorf("cannot route [%s] from the end [%s]", name, r.Name()))
}

// ProceedOnFirstRoute proceeds with the first route available
func (r *endNode) ProceedOnFirstRoute(mex f.Message) error {
	return errors.Join(g.ErrorInvalidRouting, fmt.Errorf("cannot route from the end [%s]", r.Name()))
}
