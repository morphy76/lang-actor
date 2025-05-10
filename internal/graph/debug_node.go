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
		return errors.Join(g.ErrorInvalidRouting, fmt.Errorf("debugNode node [%v] already has a route", r.Address()))
	}

	r.edges[name] = edge{
		Name:        name,
		Destination: destination.Address(),
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
	useDebugNode := &debugNode{
		node: *baseNode,
	}

	taskFn := func(msg f.Message, self f.Actor[string]) (string, error) {
		fmt.Printf("Debug node received message: %+v\n", msg)

		// TODO timeout context between request (not a config message) and response (receiving a config message)

		cfgMessage, ok := msg.(*g.ConfigMessage)
		if ok {
			for key, val := range cfgMessage.Entries {
				fmt.Printf("Debug node received key: %s with value %v\n", key, val)
			}
			useDebugNode.ProceedOnAnyRoute(msg)
		} else {
			requestCfg, err := g.NewConfigMessage(self.Address(), g.Entries)
			if err != nil {
				return self.State(), err
			}
			cfgNodes := useDebugNode.GetResolver().Query("graph", "nodes", "config")
			if len(cfgNodes) == 0 {
				return self.State(), errors.Join(g.ErrorInvalidRouting, fmt.Errorf("no config node found"))
			}
			cfgNodes[0].Deliver(requestCfg)
			fmt.Printf("Debug node requesting config\n")
		}

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

	useDebugNode.actor = debugTask

	return useDebugNode, nil
}
