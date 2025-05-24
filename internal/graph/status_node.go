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

// NewStatusNode creates a new status node with the given configuration and graph name.
func NewStatusNode(initialStatus interface{}, graphName string) (g.Node, error) {

	address, err := url.Parse("graph://nodes/status/" + graphName)
	if err != nil {
		return nil, err
	}

	actorAddress, err := url.Parse("actor://" + address.Host + address.Path + "/" + uuid.NewString())
	if err != nil {
		return nil, err
	}

	baseNode := NewNode(*address)
	useNode := &statusNode{
		node: *baseNode,
	}

	taskFn := func(msg f.Message, self f.Actor[interface{}]) (interface{}, error) {
		if msg == nil {
			return self.State(), fmt.Errorf("message is nil")
		}
		if _, ok := msg.(*g.StatusMessage); !ok {
			fmt.Printf("message type: %T\n", msg)
			return self.State(), fmt.Errorf("message is not a status message")
		}
		useMex := msg.(*g.StatusMessage)
		switch useMex.StatusMessageType {
		case g.StatusRequest:
			replyMsg := g.NewStatusMessageResponse(self.Address(), self.State())
			addressable, found := useNode.GetResolver().Resolve(useMex.Sender())
			if !found {
				return self.State(), fmt.Errorf("addressable not found")
			}
			if err := addressable.Deliver(&replyMsg); err != nil {
				return self.State(), err
			}
		case g.StatusUpdate:
			useMexWithPayload, ok := msg.(*g.StatusMessage)
			if !ok {
				return self.State(), fmt.Errorf("message is not a status message with payload")
			}
			return useMexWithPayload.Value, nil
		}

		return self.State(), nil
	}

	statusTask, err := framework.NewActor(
		*actorAddress,
		taskFn,
		initialStatus,
		false,
	)
	if err != nil {
		return nil, err
	}

	useNode.actor = statusTask

	return useNode, nil
}

var staticStatusAssertion g.Node = (*statusNode)(nil)

type statusNode struct {
	node
}

// OneWayRoute adds a new possible outgoing route from the node
func (r *statusNode) OneWayRoute(name string, destination g.Node) error {
	return errors.Join(g.ErrorInvalidRouting, fmt.Errorf("cannot route [%s] from the status node [%v]", name, r.Address()))
}

// TwoWayRoute adds a new possible outgoing route from the node
func (r *statusNode) TwoWayRoute(name string, destination g.Node) error {
	return errors.Join(g.ErrorInvalidRouting, fmt.Errorf("cannot route [%s] from the status node [%v]", name, r.Address()))
}

// ProceedOnAnyRoute proceeds with any route available
func (r *statusNode) ProceedOnAnyRoute(mex f.Message) error {
	return errors.Join(g.ErrorInvalidRouting, fmt.Errorf("cannot route from the status node [%v]", r.Address()))
}
