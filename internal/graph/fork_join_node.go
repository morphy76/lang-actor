package graph

import (
	"fmt"
	"net/url"

	"github.com/google/uuid"

	"github.com/morphy76/lang-actor/internal/framework"
	f "github.com/morphy76/lang-actor/pkg/framework"
	g "github.com/morphy76/lang-actor/pkg/graph"
)

var staticForkJoinAssertion g.ForkJoinNode = (*forkJoin)(nil)

type forkJoin struct {
	node
}

// NewForkJoingNode creates a new fork-join node for the given graph.
func NewForkJoingNode[C g.NodeState](forGraph g.Graph, transient bool, processingFns ...f.ProcessingFn[C]) (g.Node, error) {

	address, err := url.Parse("graph://nodes/fork/" + uuid.NewString())
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
		childState := g.BasicNodeStateBuilder[C](forGraph, childOutcomes)
		framework.NewActorWithParent(
			processingFn,
			childState,
			true,
			baseNode.actor,
		)
	}

	return &forkJoin{
		node: *baseNode,
	}, nil
}
