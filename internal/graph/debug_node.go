package graph

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/google/uuid"
	"github.com/morphy76/lang-actor/internal/framework"
	f "github.com/morphy76/lang-actor/pkg/framework"
	g "github.com/morphy76/lang-actor/pkg/graph"
)

var staticDebugAssertion g.DebugNode = (*debugNode)(nil)

type debugNode struct {
	node
}

// OneWayRoute adds a new possible outgoing route from the node
func (r *debugNode) OneWayRoute(name string, destination g.Node) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if len(r.edges) > 0 {
		return errors.Join(g.ErrorInvalidRouting, fmt.Errorf("debugNode node [%s] already has a route", r.Name()))
	}

	r.edges[name] = edge{
		Name:        name,
		Destination: destination,
	}

	return nil
}

// NewDebugNode creates a new instance of a debug node in the actor graph.
func NewDebugNode() (g.Node, error) {

	address, err := url.Parse("graph://nodes/debug/" + uuid.NewString())
	if err != nil {
		return nil, err
	}

	actorAddress, err := url.Parse("actor://" + address.Host + address.Path + "/" + uuid.NewString())
	if err != nil {
		return nil, err
	}

	baseNode := newNode[string](nil, *address)
	taskFn := func(msg f.Message, self f.Actor[string]) (string, error) {
		fmt.Printf("Debug node received message: %+v\n", msg)

		cfgMessage, ok := msg.(*g.ConfigMessage)
		if ok {
			for key, val := range cfgMessage.Entries {
				fmt.Printf("Debug node received key: %s with value %v\n", key, val)
			}
		} else {
			_, err := g.NewConfigMessage(self.Address(), g.Entries)
			if err != nil {
				return self.State(), err
			}
			// TODO every node has to access the address book to send a message
		}

		baseNode.ProceedOnAnyRoute(msg)

		return self.State(), nil
	}

	debugTask, err := framework.NewActor(
		*actorAddress,
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
