package graph

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/google/uuid"

	"github.com/morphy76/lang-actor/internal/framework"
	f "github.com/morphy76/lang-actor/pkg/framework"
	g "github.com/morphy76/lang-actor/pkg/graph"
)

var staticForkJoinNodeAssertion g.ForkJoinNode = (*forkJoinNode)(nil)
var staticForkNodeAssertion g.ForkNode = (*forkNode)(nil)
var staticJoinNodeAssertion g.JoinNode = (*joinNode)(nil)

type forkJoinNode struct {
	node
}

type forkNode struct {
	node
}

type joinNode struct {
	node
}

// NewForkJoingNode creates a new fork-join node for the given graph.
func NewForkJoingNode[C g.NodeRef](forGraph g.Graph, transient bool, processingFns ...f.ProcessingFn[C]) (g.Node, error) {
	address, err := url.Parse("graph://nodes/forkjoin/" + uuid.NewString())
	if err != nil {
		return nil, err
	}

	expectedOutcomes := len(processingFns)
	childOutcomes := make(chan string, expectedOutcomes)

	taskFn := func(msg f.Message, self f.Actor[g.NodeRef]) (g.NodeRef, error) {

		for _, child := range self.Children() {
			child.Deliver(msg, self)
		}

		receivedOutcomes := 0
		for childOutcome := range childOutcomes {
			self.State().GraphState().MergeChange(nil, childOutcome)
			receivedOutcomes++
			if receivedOutcomes == expectedOutcomes {
				break
			}
		}

		self.State().ProceedOntoRoute() <- g.WhateverOutcome
		return self.State(), nil
	}

	baseNode, err := newNode(forGraph, *address, taskFn)
	if err != nil {
		return nil, err
	}

	for _, processingFn := range processingFns {
		childState := g.BasicNodeRefBuilder[C](forGraph, baseNode, childOutcomes)
		framework.NewActorWithParent(
			processingFn,
			childState,
			baseNode.actor,
		)
	}

	return &forkJoinNode{
		node: *baseNode,
	}, nil
}

// NewForkNode creates a new instance of a fork node in the actor graph.
func NewForkNode(forGraph g.Graph) (g.Node, error) {
	address, err := url.Parse("graph://nodes/fork/" + uuid.NewString())
	if err != nil {
		return nil, err
	}

	taskFn := func(msg f.Message, self f.Actor[g.NodeRef]) (g.NodeRef, error) {
		var wg sync.WaitGroup
		for _, edgeName := range self.State().Routes() {
			wg.Add(1)
			go func() {
				defer wg.Done()
				self.State().ProceedOntoRoute() <- edgeName
			}()
		}
		wg.Wait()
		self.State().ProceedOntoRoute() <- g.SkipOutcome
		return self.State(), nil
	}

	baseNode, err := newNode(forGraph, *address, taskFn)
	if err != nil {
		return nil, err
	}
	baseNode.multipleOutcomes = true

	return &forkNode{
		node: *baseNode,
	}, nil
}

// NewJoinNode creates a new instance of a join node in the actor graph.
func NewJoinNode(forGraph g.Graph, forkNode g.Node) (g.Node, error) {
	address, err := url.Parse("graph://nodes/join/" + uuid.NewString())
	if err != nil {
		return nil, err
	}

	taskFn := func(msg f.Message, self f.Actor[g.NodeRef]) (g.NodeRef, error) {
		received, _ := self.State().GetAttribute("received")

		inbounds := len(forkNode.EdgeNames())
		if inbounds == 0 {
			return nil, fmt.Errorf("join node must have at least one inbound edge, but got %d", inbounds)
		}

		if received.(int) < inbounds-1 {
			self.State().SetAttribute("received", received.(int)+1)
			self.State().ProceedOntoRoute() <- g.SkipOutcome
			return self.State(), nil
		}

		self.State().ProceedOntoRoute() <- g.WhateverOutcome
		self.State().SetAttribute("received", 0)
		return self.State(), nil
	}

	attrs := make(map[string]any)
	attrs["received"] = 0

	baseNode, err := newNode(forGraph, *address, taskFn, attrs)
	if err != nil {
		return nil, err
	}
	baseNode.multipleOutcomes = true

	return &joinNode{
		node: *baseNode,
	}, nil
}
