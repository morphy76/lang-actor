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

func newNode[T any](task f.Actor[T], address url.URL) *node {
	return &node{
		lock:   &sync.Mutex{},
		routes: make(map[string]route, 0),
		actor:  task,
		name:   fmt.Sprintf("/%s%s", address.Host, address.Path),
	}
}

// NewDebugNode creates a new instance of a debug node in the actor graph.
func NewDebugNode() (g.Node, error) {
	address, err := url.Parse("actor://nodes/debug/" + uuid.NewString())
	if err != nil {
		return nil, err
	}

	baseNode := newNode[string](nil, *address)
	taskFn := func(msg f.Message, self f.Actor[string]) (string, error) {
		fmt.Printf("Debug node received message: %+v\n", msg)
		baseNode.ProceedOnAnyRoute(msg)

		return self.State(), nil
	}

	debugTask, err := framework.NewActor(
		*address,
		taskFn,
		"",
	)
	if err != nil {
		return nil, err
	}

	rv := &debugNode{
		node: *baseNode,
	}
	rv.actor = debugTask

	return rv, nil
}
