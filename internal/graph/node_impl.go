package graph

import (
	"errors"
	"fmt"
	"sync"

	f "github.com/morphy76/lang-actor/pkg/framework"
	g "github.com/morphy76/lang-actor/pkg/graph"
)

var staticNodeAssertion g.Node = (*node)(nil)
var staticDebugAssertion g.DebugNode = (*debugNode)(nil)

type node struct {
	lock *sync.Mutex

	routes map[string]route

	actor f.ActorRef
}

// Name returns the name of the node
func (r *node) Name() string {
	return fmt.Sprintf("%s%s", r.actor.Address().Host, r.actor.Address().Path)
}

// RouteNames returns the names of all possible routes from the node
func (r *node) RouteNames() []string {
	r.lock.Lock()
	defer r.lock.Unlock()

	names := make([]string, 0, len(r.routes))
	for name := range r.routes {
		names = append(names, name)
	}

	return names
}

// OneWayRoute adds a new possible outgoing route from the node
func (r *node) OneWayRoute(name string, destination g.Node) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if _, ok := destination.(*rootNode); ok {
		return errors.Join(g.ErrorInvalidRouting, fmt.Errorf("cannot route [%s] from node [%s] to root node", name, r.Name()))
	}

	r.routes[name] = route{
		Name:        name,
		Destination: destination,
	}

	return nil
}

// TwoWayRoute adds a new possible outgoing route from the node
func (r *node) TwoWayRoute(name string, destination g.Node) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if _, ok := destination.(*rootNode); ok {
		return errors.Join(g.ErrorInvalidRouting, fmt.Errorf("cannot route [%s] from node [%s] to root node", name, r.Name()))
	}

	if _, ok := destination.(*endNode); ok {
		return errors.Join(g.ErrorInvalidRouting, fmt.Errorf("cannot route [%s] from node [%s] from end node", name, r.Name()))
	}

	r.routes[name] = route{
		Name:        name,
		Destination: destination,
	}

	var meAsNode g.Node = r
	return destination.OneWayRoute("inverse-"+name, meAsNode)
}

type debugNode struct {
	node
}
