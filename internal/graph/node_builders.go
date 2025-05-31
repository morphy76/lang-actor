package graph

import (
	"net/url"
	"sync"

	"github.com/morphy76/lang-actor/internal/framework"
	f "github.com/morphy76/lang-actor/pkg/framework"
	g "github.com/morphy76/lang-actor/pkg/graph"
)

func newNode(
	forGraph g.Graph,
	address url.URL,
	taskFn f.ProcessingFn[g.NodeState],
	transient bool,
) (*node, error) {

	actorAddress, err := url.Parse("actor://" + address.Host + address.Path)
	if err != nil {
		return nil, err
	}

	actorOutcome := make(chan string, 1)
	useState := &nodeState{
		outcome: actorOutcome,
		graph:   forGraph,
	}
	task, err := framework.NewActor[g.NodeState](
		*actorAddress,
		taskFn,
		useState,
		transient,
	)
	if err != nil {
		return nil, err
	}

	rv := &node{
		lock:         &sync.Mutex{},
		edges:        make(map[string]edge, 0),
		address:      address,
		actor:        task,
		actorOutcome: actorOutcome,
		nodeState:    useState,
		resolver:     forGraph,
	}

	if forGraph != nil {
		forGraph.Register(rv)
	}

	return rv, nil
}
