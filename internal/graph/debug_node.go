package graph

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/google/uuid"
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
		Destination: destination,
	}

	return nil
}

// NewDebugNode creates a new instance of a debug node in the actor graph.
func NewDebugNode(forGraph g.Graph, nameParts ...string) (g.Node, error) {

	baseName := "graph://nodes/debug/" + uuid.NewString()
	if len(nameParts) > 0 {
		baseName = "graph://nodes/debug/" + nameParts[0] + "/" + uuid.NewString()
		for i := 1; i < len(nameParts); i++ {
			baseName += "/" + nameParts[i]
		}
	}

	address, err := url.Parse(baseName)
	if err != nil {
		return nil, err
	}

	taskFn := func(msg f.Message, self f.Actor[g.NodeState]) (g.NodeState, error) {
		fmt.Println("==========================================")
		fmt.Printf("Debug node [%+v] received message:\n", self.Address())
		jsonOriginalMessage, err := json.Marshal(msg)
		if err != nil {
			fmt.Printf("%s\n", err)
		} else {
			fmt.Printf("%s\n", jsonOriginalMessage)
		}
		fmt.Println("---------------------------------")
		fmt.Println("System config:")
		jsonConfigResponse, err := json.Marshal(self.State().GraphConfig)
		if err != nil {
			fmt.Printf("%s\n", err)
		} else {
			fmt.Printf("%s\n", jsonConfigResponse)
		}
		fmt.Println("---------------------------------")
		fmt.Println("Graph status:")
		jsonStateResponse, err := json.Marshal(self.State().GraphState)
		if err != nil {
			fmt.Printf("%s\n", err)
		} else {
			fmt.Printf("%s\n", jsonStateResponse)
		}
		fmt.Println("==========================================")
		self.State().Outcome <- ""
		return self.State(), nil
	}

	baseNode, err := newNode(forGraph, *address, taskFn, true)
	if err != nil {
		return nil, err
	}

	return &debugNode{
		node: *baseNode,
	}, nil
}
