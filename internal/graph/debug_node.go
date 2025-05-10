package graph

import (
	"encoding/json"
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
		Destination: destination,
	}

	return nil
}

type debugNodeState struct {
	originalMessage f.Message
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

	baseNode := newNode[debugNodeState](nil, *address)
	useDebugNode := &debugNode{
		node: *baseNode,
	}

	taskFn := func(msg f.Message, self f.Actor[debugNodeState]) (debugNodeState, error) {

		// TODO timeout context between request (not a config message) and response (receiving a config message)

		statusResponse, okStatus := msg.(*g.StatusMessage[any])
		configResponse, okConfig := msg.(*g.ConfigMessage)

		if okStatus {
			fmt.Println("==========================================")
			fmt.Printf("Debug node [%+v] received message:\n", useDebugNode.Address())
			fmt.Println("---------------------------------")
			fmt.Println("Status response:")
			jsonStatusResponse, err := json.Marshal(statusResponse)
			if err != nil {
				fmt.Printf("%s\n", err)
			} else {
				fmt.Printf("%s\n", jsonStatusResponse)
			}
			fmt.Println("==========================================")
			err = useDebugNode.ProceedOnAnyRoute(self.State().originalMessage)
			if err != nil {
				return self.State(), err
			}
		} else {
			if okConfig {
				fmt.Println("==========================================")
				fmt.Printf("Debug node [%+v] received message:\n", useDebugNode.Address())
				jsonOriginalMessage, err := json.Marshal(self.State().originalMessage)
				if err != nil {
					fmt.Printf("%s\n", err)
				} else {
					fmt.Printf("%s\n", jsonOriginalMessage)
				}
				fmt.Println("---------------------------------")
				fmt.Println("System config:")
				jsonConfigResponse, err := json.Marshal(configResponse.Entries)
				if err != nil {
					fmt.Printf("%s\n", err)
				} else {
					fmt.Printf("%s\n", jsonConfigResponse)
				}
				fmt.Println("==========================================")
				err = useDebugNode.ProceedOnAnyRoute(self.State().originalMessage)
				if err != nil {
					return self.State(), err
				}

				requestStatus, err := g.NewStatusMessageRequest[any](self.Address())
				if err != nil {
					return self.State(), err
				}
				statusNodes := useDebugNode.GetResolver().Query("graph", "nodes", "status")
				if len(statusNodes) == 0 {
					return self.State(), errors.Join(g.ErrorInvalidRouting, fmt.Errorf("no status node found"))
				}
				statusNodes[0].Deliver(requestStatus)
				return self.State(), nil
			} else {
				requestCfg, err := g.NewConfigMessage(self.Address(), g.ConfigEntries)
				if err != nil {
					return self.State(), err
				}
				cfgNodes := useDebugNode.GetResolver().Query("graph", "nodes", "config")
				if len(cfgNodes) == 0 {
					return self.State(), errors.Join(g.ErrorInvalidRouting, fmt.Errorf("no config node found"))
				}
				cfgNodes[0].Deliver(requestCfg)
				return debugNodeState{
					originalMessage: msg,
				}, nil
			}
		}

		return self.State(), nil
	}

	debugTask, err := framework.NewActor(
		*actorAddress,
		taskFn,
		debugNodeState{},
		false,
	)
	if err != nil {
		return nil, err
	}

	useDebugNode.actor = debugTask

	return useDebugNode, nil
}
