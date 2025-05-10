package graph

import (
	"errors"
	"fmt"
	"net/url"
	"sync"

	f "github.com/morphy76/lang-actor/pkg/framework"
	g "github.com/morphy76/lang-actor/pkg/graph"
	r "github.com/morphy76/lang-actor/pkg/routing"
)

var staticNodeAssertion g.Node = (*node)(nil)

type node struct {
	lock *sync.Mutex

	edges map[string]edge

	name        string
	actor       f.ActorRef
	addressBook r.AddressBook
}

// Name returns the name of the node
func (r *node) Name() string {
	return r.name
}

// RouteNames returns the names of all possible routes from the node
func (r *node) RouteNames() []string {
	r.lock.Lock()
	defer r.lock.Unlock()

	names := make([]string, 0, len(r.edges))
	for name := range r.edges {
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

	r.edges[name] = edge{
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

	r.edges[name] = edge{
		Name:        name,
		Destination: destination,
	}

	var meAsNode g.Node = r
	return destination.OneWayRoute("inverse-"+name, meAsNode)
}

// DestinationAddress returns the address of the destination node
func (r *node) Address() url.URL {
	return r.actor.Address()
}

// Deliver delivers a message to the node
func (r *node) Deliver(mex f.Message) error {
	return r.actor.Deliver(mex)
}

// Send sends a message to the node
func (r *node) Send(mex f.Message, addressable f.Addressable) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	for _, route := range r.edges {
		if route.Destination.Address() == addressable.Address() {
			return route.Destination.Deliver(mex)
		}
	}
	destinationAddress := fmt.Sprintf("%s%s", addressable.Address().Host, addressable.Address().Path)
	return errors.Join(g.ErrorInvalidRouting, fmt.Errorf("cannot route message to [%s] from node [%s]", destinationAddress, r.Name()))
}

// ProceedOnAnyRoute proceeds with the first route available
func (r *node) ProceedOnAnyRoute(mex f.Message) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if len(r.edges) == 0 {
		return errors.Join(g.ErrorInvalidRouting, fmt.Errorf("node [%s] has no routes", r.Name()))
	}

	for _, route := range r.edges {
		return route.Destination.Deliver(mex)
	}

	return nil
}
