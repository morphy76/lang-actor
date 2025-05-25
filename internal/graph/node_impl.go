package graph

import (
	"errors"
	"fmt"
	"net/url"
	"sync"

	c "github.com/morphy76/lang-actor/pkg/common"
	f "github.com/morphy76/lang-actor/pkg/framework"
	g "github.com/morphy76/lang-actor/pkg/graph"
	r "github.com/morphy76/lang-actor/pkg/routing"
)

var staticNodeAssertion g.Node = (*node)(nil)
var staticNodeMessageAssertion f.Message = (*nodeMessage)(nil)

type node struct {
	lock *sync.Mutex

	resolver r.Resolver

	edges map[string]edge

	address url.URL
	actor   f.ActorRef

	actorOutcome chan string

	nodeState g.NodeState
}

type nodeMessage struct {
	sender  url.URL
	payload any
}

// Sender returns the sender of the message
func (m nodeMessage) Sender() url.URL {
	return m.sender
}

// Mutation returns false, indicating that this message is not a mutation
func (m nodeMessage) Mutation() bool {
	return false
}

// Edges returns the edges of the node
func (r *node) Edges() []c.Addressable {
	rv := make([]c.Addressable, 0, len(r.edges))
	for _, edge := range r.edges {
		rv = append(rv, edge.Destination)
	}

	return rv
}

// OneWayRoute adds a new possible outgoing route from the node
func (r *node) OneWayRoute(name string, destination g.Node) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if _, ok := destination.(*rootNode); ok {
		return errors.Join(g.ErrorInvalidRouting, fmt.Errorf("cannot route [%s] from node [%v] to root node", name, r.Address()))
	}

	r.edges[name] = edge{
		Name:        name,
		Destination: destination,
	}

	return nil
}

// DestinationAddress returns the address of the destination node
func (r *node) Address() url.URL {
	return r.address
}

// ProceedOnAnyRoute proceeds with the first route available
func (r *node) ProceedOnAnyRoute(mex c.Message) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if len(r.edges) == 0 {
		return errors.Join(g.ErrorInvalidRouting, fmt.Errorf("node [%v] has no routes", r.Address()))
	}

	for _, route := range r.edges {
		resolver := r.GetResolver()
		addressable, found := resolver.Resolve(route.Destination.Address())
		if !found {
			return errors.Join(g.ErrorInvalidRouting, fmt.Errorf("Unknown address [%v] from node [%v]", route.Destination, r.Address()))
		} else {
			handler, ok := addressable.(c.MessageHandler)
			if !ok {
				return fmt.Errorf("destination [%v] is not a message handler", addressable.Address())
			}
			return handler.Accept(mex)
		}
	}

	return nil
}

// ProceedOnRoute proceeds the message on a specific route
func (r *node) ProceedOnRoute(name string, mex c.Message) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	edge, ok := r.edges[name]
	if !ok {
		return errors.Join(g.ErrorInvalidRouting, fmt.Errorf("node [%v] has no route named [%s]", r.Address(), name))
	}
	resolver := r.GetResolver()
	addressable, found := resolver.Resolve(edge.Destination.Address())
	if !found {
		return errors.Join(g.ErrorInvalidRouting, fmt.Errorf("Unknown address [%v] from node [%v]", edge.Destination.Address(), r.Address()))
	} else {
		handler, ok := addressable.(c.MessageHandler)
		if !ok {
			return fmt.Errorf("destination [%v] is not a message handler", addressable.Address())
		}
		return handler.Accept(mex)
	}
}

func (r *node) Accept(message c.Message) error {
	_, ok := message.(f.Message)
	if ok {
		if err := r.actor.Deliver(message); err != nil {
			return fmt.Errorf("failed to deliver message to node [%v]: %w", r.Address(), err)
		}
	} else {
		actorMessage := &nodeMessage{
			sender:  r.actor.Address(),
			payload: message,
		}
		if err := r.actor.Deliver(actorMessage); err != nil {
			return fmt.Errorf("failed to deliver message to node [%v]: %w", r.Address(), err)
		}
	}

	// TODO implement a timeout for the outcome channel
	outcome := <-r.actorOutcome
	if outcome != "" {
		r.ProceedOnRoute(outcome, message)
	} else {
		r.ProceedOnAnyRoute(message)
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

// SetConfig sets the configuration for the graph-aware component.
func (r *node) SetConfig(config g.GraphConfiguration) {
	r.nodeState.GraphConfig = config
}

// SetState sets the state for the graph-aware component.
func (r *node) SetState(state g.GraphState) {
	r.nodeState.GraphState = state
}

// GetConfig retrieves the configuration of the graph-aware component.
func (r *node) GetConfig() g.GraphConfiguration {
	return r.nodeState.GraphConfig
}

// GetState retrieves the state of the graph-aware component.
func (r *node) GetState() g.GraphState {
	return r.nodeState.GraphState
}
