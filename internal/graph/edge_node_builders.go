package graph

import (
	"net/url"

	"github.com/google/uuid"

	f "github.com/morphy76/lang-actor/pkg/framework"
	g "github.com/morphy76/lang-actor/pkg/graph"
)

// NewRootNode creates a new instance of the actor graph.
func NewRootNode(forGraph g.Graph) (g.Node, error) {
	address, err := url.Parse("graph://edge/root/" + uuid.NewString())
	if err != nil {
		return nil, err
	}

	baseNode, err := newNode(forGraph, *address, func(msg f.Message, self f.Actor[g.NodeState]) (g.NodeState, error) {
		self.State().Outcome() <- g.WhateverOutcome
		return self.State(), nil
	}, true)
	if err != nil {
		return nil, err
	}

	return &rootNode{
		node: *baseNode,
	}, nil
}

// NewEndNode creates a new instance of the end node in the actor graph.
func NewEndNode(forGraph g.Graph) (g.Node, error) {
	address, err := url.Parse("graph://edge/end/" + uuid.NewString())
	if err != nil {
		return nil, err
	}

	baseNode, err := newNode(forGraph, *address, func(msg f.Message, self f.Actor[g.NodeState]) (g.NodeState, error) {
		self.State().Outcome() <- g.WhateverOutcome
		return self.State(), nil
	}, true)
	if err != nil {
		return nil, err
	}

	return &endNode{
		node: *baseNode,
	}, nil
}
