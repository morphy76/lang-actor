package graph

import (
	"errors"
	"fmt"

	f "github.com/morphy76/lang-actor/pkg/framework"
	g "github.com/morphy76/lang-actor/pkg/graph"
)

var staticConfigAssertion g.Node = (*configNode)(nil)

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
