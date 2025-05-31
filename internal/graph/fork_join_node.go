package graph

import (
	"fmt"
	"net/url"

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
func NewForkJoingNode[C g.NodeState](forGraph g.Graph, transient bool, processingFns ...f.ProcessingFn[C]) (g.Node, error) {
	address, err := url.Parse("graph://nodes/forkjoin/" + uuid.NewString())
	if err != nil {
		return nil, err
	}

	expectedOutcomes := len(processingFns)
	childOutcomes := make(chan string, expectedOutcomes)

	taskFn := func(msg f.Message, self f.Actor[g.NodeState]) (g.NodeState, error) {

		for _, child := range self.Children() {
			fmt.Printf("Delivering message to child: %v\n", child.Address())
			child.Deliver(msg)
		}

		receivedOutcomes := 0
		for childOutcome := range childOutcomes {
			self.State().GraphState().AppendGraphState(nil, childOutcome)
			receivedOutcomes++
			if receivedOutcomes == expectedOutcomes {
				break
			}
		}

		self.State().Outcome() <- ""
		return self.State(), nil
	}

	baseNode, err := newNode(forGraph, *address, taskFn, true)
	if err != nil {
		return nil, err
	}

	for _, processingFn := range processingFns {
		childState := g.BasicNodeStateBuilder[C](forGraph, baseNode, childOutcomes)
		framework.NewActorWithParent(
			processingFn,
			childState,
			true,
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

	taskFn := func(msg f.Message, self f.Actor[g.NodeState]) (g.NodeState, error) {
		for _, edgeName := range self.State().Routes() {
			self.State().Outcome() <- edgeName
		}
		self.State().Outcome() <- ""
		return self.State(), nil
	}

	baseNode, err := newNode(forGraph, *address, taskFn, true)
	if err != nil {
		return nil, err
	}
	baseNode.multipleOutcomes = true

	return &forkNode{
		node: *baseNode,
	}, nil
}

// NewJoinNode creates a new instance of a join node in the actor graph.
func NewJoinNode(forGraph g.Graph) (g.Node, error) {
	address, err := url.Parse("graph://nodes/join/" + uuid.NewString())
	if err != nil {
		return nil, err
	}

	taskFn := func(msg f.Message, self f.Actor[g.NodeState]) (g.NodeState, error) {
		self.State().Outcome() <- ""
		return self.State(), nil
	}

	baseNode, err := newNode(forGraph, *address, taskFn, true)
	if err != nil {
		return nil, err
	}

	return &joinNode{
		node: *baseNode,
	}, nil
}
