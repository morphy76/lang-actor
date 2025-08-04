package graph

import (
	"net/url"
	"sync"

	"github.com/morphy76/lang-actor/internal/framework"
	f "github.com/morphy76/lang-actor/pkg/framework"
	g "github.com/morphy76/lang-actor/pkg/graph"
)

func newNode[T g.NodeRef](
	forGraph g.Graph,
	address url.URL,
	taskFn f.ProcessingFn[T],
	transient bool,
	attrs ...map[string]any,
) (*node, error) {

	actorAddress, err := url.Parse("actor://" + address.Host + address.Path)
	if err != nil {
		return nil, err
	}

	rv := &node{
		lock:             &sync.Mutex{},
		edges:            make(map[string]edge, 0),
		address:          address,
		resolver:         forGraph,
		multipleOutcomes: false,
	}

	actorOutcome := make(chan string, 1)
	useRef := g.BasicNodeRefBuilder[T](forGraph, rv, actorOutcome, attrs...)

	task, err := framework.NewActor(
		*actorAddress,
		taskFn,
		useRef,
		transient,
	)
	if err != nil {
		return nil, err
	}

	if forGraph != nil {
		forGraph.Register(rv)
	}

	rv.actor = task
	rv.nodeRef = useRef

	return rv, nil
}
