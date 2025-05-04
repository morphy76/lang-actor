package graph

import (
	"net/url"
	"sync"

	"github.com/google/uuid"
	"github.com/morphy76/lang-actor/internal/framework"
	f "github.com/morphy76/lang-actor/pkg/framework"
	g "github.com/morphy76/lang-actor/pkg/graph"
)

func newNode[T any](task f.Actor[T]) *node {
	return &node{
		lock:   &sync.Mutex{},
		routes: make(map[string]route, 0),
		actor:  task,
	}
}

// NewDebugNode creates a new instance of a debug node in the actor graph.
func NewDebugNode() (g.Node, error) {
	address, err := url.Parse("actor://nodes/debug/" + uuid.NewString())
	if err != nil {
		return nil, err
	}
	debugTask, err := framework.NewActor(
		*address,
		func(msg f.Message, self f.Actor[string]) (string, error) {
			return "", nil
		},
		"",
	)
	if err != nil {
		return nil, err
	}

	baseNode := newNode(debugTask)
	return &debugNode{
		node: *baseNode,
	}, nil
}
