package graph

import (
	"errors"
	"fmt"

	f "github.com/morphy76/lang-actor/pkg/framework"
	g "github.com/morphy76/lang-actor/pkg/graph"
)

type debugNode struct {
	node
}

// OneWayRoute adds a new possible outgoing route from the node
func (r *debugNode) OneWayRoute(name string, destination g.Node) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if len(r.routes) > 0 {
		return errors.Join(g.ErrorInvalidRouting, fmt.Errorf("debugNode node [%s] already has a route", r.Name()))
	}

	r.routes[name] = route{
		Name:        name,
		Destination: destination,
	}

	return nil
}

type configNode struct {
	node
}

// OneWayRoute adds a new possible outgoing route from the node
func (r *configNode) OneWayRoute(name string, destination g.Node) error {
	return errors.Join(g.ErrorInvalidRouting, fmt.Errorf("cannot route [%s] from the config node [%s]", name, r.Name()))
}

// TwoWayRoute adds a new possible outgoing route from the node
func (r *configNode) TwoWayRoute(name string, destination g.Node) error {
	return errors.Join(g.ErrorInvalidRouting, fmt.Errorf("cannot route [%s] from the config node [%s]", name, r.Name()))
}

// ProceedOnFirstRoute proceeds with the first route available
func (r *configNode) ProceedOnFirstRoute(mex f.Message) error {
	return errors.Join(g.ErrorInvalidRouting, fmt.Errorf("cannot route from the config node [%s]", r.Name()))
}
