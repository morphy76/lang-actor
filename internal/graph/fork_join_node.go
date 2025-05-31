package graph

import (
	"net/url"

	"github.com/google/uuid"

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

	taskFn := func(msg f.Message, self f.Actor[g.NodeState]) (g.NodeState, error) {

		// TODO
		self.State().Outcome() <- ""
		return self.State(), nil
	}

	baseNode, err := newNode(forGraph, *address, taskFn, true)
	if err != nil {
		return nil, err
	}

	return &forkJoin{
		node: *baseNode,
	}, nil
}
