package graph

import (
	"net/url"

	"github.com/google/uuid"
	"github.com/morphy76/lang-actor/internal/framework"
	f "github.com/morphy76/lang-actor/pkg/framework"
	g "github.com/morphy76/lang-actor/pkg/graph"
)

// NewRootNode creates a new instance of the actor graph.
func NewRootNode() (g.Node, error) {
	address, err := url.Parse("actor://edge/root/" + uuid.NewString())
	if err != nil {
		return nil, err
	}
	rootTask, err := framework.NewActor(
		*address,
		func(msg f.Message, self f.Actor[string]) (string, error) {
			return "", nil
		},
		"",
	)
	if err != nil {
		return nil, err
	}

	baseNode := newNode(rootTask, *address)
	return &rootNode{
		node: *baseNode,
	}, nil
}

// NewEndNode creates a new instance of the end node in the actor graph.
func NewEndNode() (g.Node, chan bool, error) {
	address, err := url.Parse("actor://edge/end/" + uuid.NewString())
	if err != nil {
		return nil, nil, err
	}
	endCh := make(chan bool)
	endTask, err := framework.NewActor(
		*address,
		func(msg f.Message, self f.Actor[string]) (string, error) {
			endCh <- true
			return "", nil
		},
		"",
	)
	if err != nil {
		return nil, nil, err
	}

	baseNode := newNode(endTask, *address)
	return &endNode{
		node: *baseNode,
	}, endCh, nil
}
