package graph

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/google/uuid"

	f "github.com/morphy76/lang-actor/pkg/framework"
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
		return errors.Join(g.ErrorInvalidRouting, fmt.Errorf("debugNode node [%v] already has a route", r.Address()))
	}

	r.edges[name] = edge{
		Name:        name,
		Destination: destination,
	}

	return nil
}

// NewDebugNode creates a new instance of a debug node in the actor graph.
func NewDebugNode(forGraph g.Graph, nameParts ...string) (g.Node, error) {

	baseName := "graph://nodes/debug/" + uuid.NewString()
	if len(nameParts) > 0 {
		baseName = "graph://nodes/debug/" + nameParts[0] + "/" + uuid.NewString()
		for i := 1; i < len(nameParts); i++ {
			baseName += "/" + nameParts[i]
		}
	}

	address, err := url.Parse(baseName)
	if err != nil {
		return nil, err
	}

	taskFn := func(msg f.Message, self f.Actor[g.NodeRef]) (g.NodeRef, error) {
		fmt.Println("==========================================")
		fmt.Printf("Debug node [%+v] received message [%+v]\n", self.Address(), msg)
		fmt.Println("---------------------------------")
		fmt.Printf("System config [%+v]\n", self.State().GraphConfig())
		fmt.Println("---------------------------------")
		fmt.Printf("Graph status [%+v]\n", self.State().GraphState())
		fmt.Println("==========================================")
		self.State().ProceedOntoRoute() <- g.WhateverOutcome
		return self.State(), nil
	}

	baseNode, err := newNode(forGraph, *address, taskFn)
	if err != nil {
		return nil, err
	}

	return &debugNode{
		node: *baseNode,
	}, nil
}
