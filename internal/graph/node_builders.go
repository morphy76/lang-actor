package graph

import (
	"net/url"
	"sync"

	"github.com/google/uuid"

	"github.com/morphy76/lang-actor/internal/framework"
	f "github.com/morphy76/lang-actor/pkg/framework"
	g "github.com/morphy76/lang-actor/pkg/graph"
)

// NewNode creates a new instance of a node in the graph with the given address and actor.
func NewNode(
	forGraph g.Graph,
	address url.URL,
	taskFn f.ProcessingFn[g.NodeState],
	transient bool,
) (*node, error) {
	actorAddress, err := url.Parse("actor://" + address.Host + address.Path + "/" + uuid.NewString())
	if err != nil {
		return nil, err
	}

	actorOutcome := make(chan string, 1)
	nodeState := g.NodeState{
		Outcome: actorOutcome,
	}
	task, err := framework.NewActor(
		*actorAddress,
		taskFn,
		nodeState,
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
		nodeState:    nodeState,
	}

	if forGraph != nil {
		forGraph.Register(rv)
		rv.SetResolver(forGraph)
		rv.SetConfig(forGraph.Config())
		rv.SetState(forGraph.State())
	}

	return rv, nil
}
