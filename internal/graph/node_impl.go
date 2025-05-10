package graph

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"

	f "github.com/morphy76/lang-actor/pkg/framework"
	g "github.com/morphy76/lang-actor/pkg/graph"
	r "github.com/morphy76/lang-actor/pkg/routing"
)

var staticNodeAssertion g.Node = (*node)(nil)

type node struct {
	lock *sync.Mutex

	resolver r.Resolver

	edges map[string]edge

	name        string
	actor       f.ActorRef
	addressBook r.AddressBook
}

// Name returns the name of the node
func (r *node) Name() string {
	return r.name
}

// Edges returns the edges of the node
func (r *node) Edges(includeInverse bool) []url.URL {
	edges := make([]url.URL, 0, len(r.edges))
	count := 0
	for edgeName, edge := range r.edges {
		if !includeInverse && strings.Contains(edgeName, "inverse-") {
			continue
		}
		count++
		edges = append(edges, edge.Destination)
	}
	if includeInverse {
		return edges
	}

	rv := make([]url.URL, count)
	for _, edge := range edges {
		rv = append(rv, edge)
	}

	return rv
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
		Destination: destination.Address(),
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
		Destination: destination.Address(),
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
		if route.Destination == addressable.Address() {
			addressable, found := r.GetResolver().Resolve(route.Destination)
			if !found {
				return errors.Join(g.ErrorInvalidRouting, fmt.Errorf("Unknown addres [%s] from node [%s]", route.Destination, r.Name()))
			} else {
				return addressable.Deliver(mex)
			}
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
		addressable, found := r.GetResolver().Resolve(route.Destination)
		if !found {
			return errors.Join(g.ErrorInvalidRouting, fmt.Errorf("Unknown addres [%s] from node [%s]", route.Destination, r.Name()))
		} else {
			return addressable.Deliver(mex)
		}
	}

	return nil
}

// SetResolver sets the resolver for the node
func (r *node) SetResolver(resolver r.Resolver) {
	r.resolver = resolver
}

// GetResolver returns the resolver for the node
func (r *node) GetResolver() r.Resolver {
	return r.resolver
}

// Visit visits the node and applies the given function
func (r *node) Visit(fn g.VisitFn) {

	fn(r)

	for _, edge := range r.edges {
		addressable, found := r.GetResolver().Resolve(edge.Destination)
		if !found {
			if node, ok := addressable.(g.Node); ok {
				node.Visit(fn)
			}
		}
	}
}
