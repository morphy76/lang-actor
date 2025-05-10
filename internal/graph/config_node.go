package graph

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/google/uuid"
	"github.com/morphy76/lang-actor/internal/framework"
	f "github.com/morphy76/lang-actor/pkg/framework"
	g "github.com/morphy76/lang-actor/pkg/graph"
	r "github.com/morphy76/lang-actor/pkg/routing"
)

// NewConfigNode creates a new configuration node with the given configuration and graph name.
func NewConfigNode(config map[string]any, graphName string) (g.Node, error) {

	address, err := url.Parse("graph://nodes/config/" + graphName)
	if err != nil {
		return nil, err
	}

	actorAddress, err := url.Parse("actor://" + address.Host + address.Path + "/" + uuid.NewString())
	if err != nil {
		return nil, err
	}

	baseNode := newNode[map[string]any](nil, *address)
	taskFn := func(msg f.Message, self f.Actor[map[string]any]) (map[string]any, error) {

		if msg == nil {
			return nil, fmt.Errorf("message is nil")
		}
		if _, ok := msg.(*g.ConfigMessage); !ok {
			return nil, fmt.Errorf("message is not a config message")
		}

		useMex := msg.(*g.ConfigMessage)
		useState := self.State()

		switch useMex.ConfigMessageType {
		case g.Keys:
			keys := make([]string, 0, len(useState))
			for key := range useState {
				keys = append(keys, key)
			}
			replyMsg := newListKeysResponse(self.Address(), keys)
			addressable, found := baseNode.GetResolver().Resolve(useMex.Sender())
			if !found {
				return nil, fmt.Errorf("addressable not found")
			}

			if err := addressable.Deliver(replyMsg); err != nil {
				return nil, err
			}
		case g.Entries:
			entries := make(map[string]any, len(useState))
			for key, value := range useState {
				entries[key] = value
			}
			replyMsg := newListEntriesResponse(self.Address(), entries)
			addressable, found := baseNode.GetResolver().Resolve(useMex.Sender())
			if !found {
				return nil, fmt.Errorf("addressable not found")
			}

			if err := addressable.Deliver(replyMsg); err != nil {
				return nil, err
			}
		case g.Request:
			if len(useMex.RequestedKeys) == 0 {
				shouldReturn, result, err := sendEmptyResponse(self, baseNode.GetResolver(), useMex)
				if shouldReturn {
					return result, err
				}
			} else if len(useMex.RequestedKeys) == 1 {
				key := useMex.RequestedKeys[0]
				if value, ok := useState[key]; ok {
					replyMsg := newSingleValueResponse(self.Address(), key, value)
					addressable, found := baseNode.GetResolver().Resolve(useMex.Sender())
					if !found {
						return nil, fmt.Errorf("addressable not found")
					}

					if err := addressable.Deliver(replyMsg); err != nil {
						return nil, err
					}
				} else {
					shouldReturn, result, err := sendEmptyResponse(self, baseNode.GetResolver(), useMex)
					if shouldReturn {
						return result, err
					}
				}
			} else {
				multiValueTmp := make(map[string]any, len(useMex.RequestedKeys))
				okys := 0
				for _, key := range useMex.RequestedKeys {
					if value, ok := useState[key]; ok {
						okys++
						multiValueTmp[key] = value
					}
				}
				multiValue := make(map[string]any, okys)
				for key, val := range multiValueTmp {
					multiValue[key] = val
				}
				replyMsg := newMultivalueResponse(self.Address(), multiValue)
				addressable, found := baseNode.GetResolver().Resolve(useMex.Sender())
				if !found {
					return nil, fmt.Errorf("addressable not found")
				}

				if err := addressable.Deliver(replyMsg); err != nil {
					return nil, err
				}
			}
		}

		return self.State(), nil
	}

	debugTask, err := framework.NewActor(
		*actorAddress,
		taskFn,
		config,
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

func newSingleValueResponse(sender url.URL, key string, value any) f.Message {
	return &g.ConfigMessage{
		From:              sender,
		RequestType:       g.Request,
		ConfigMessageType: g.Response,
		RequestedKeys:     []string{key},
		Value:             value,
	}
}

func newListKeysResponse(sender url.URL, configuredKeys []string) f.Message {
	return &g.ConfigMessage{
		From:              sender,
		ConfigMessageType: g.Response,
		RequestType:       g.Keys,
		Keys:              configuredKeys,
	}
}

func newMultivalueResponse(sender url.URL, entries map[string]any) f.Message {
	return &g.ConfigMessage{
		From:              sender,
		ConfigMessageType: g.Response,
		RequestType:       g.Request,
		Entries:           entries,
	}
}

func newListEntriesResponse(sender url.URL, configuredEntries map[string]any) f.Message {
	return &g.ConfigMessage{
		From:              sender,
		ConfigMessageType: g.Response,
		RequestType:       g.Entries,
		Entries:           configuredEntries,
	}
}

func sendEmptyResponse(self f.Actor[map[string]any], resolver r.Resolver, useMex *g.ConfigMessage) (bool, map[string]any, error) {
	emptyPayload := make(map[string]any, 0)
	replyMsg := newMultivalueResponse(self.Address(), emptyPayload)
	addressable, found := resolver.Resolve(useMex.Sender())
	if !found {
		return true, nil, fmt.Errorf("addressable not found")
	}

	if err := addressable.Deliver(replyMsg); err != nil {
		return true, nil, err
	}
	return false, nil, nil
}

var staticConfigAssertion g.Node = (*configNode)(nil)

type configNode struct {
	node
}

// OneWayRoute adds a new possible outgoing route from the node
func (r *configNode) OneWayRoute(name string, destination g.Node) error {
	return errors.Join(g.ErrorInvalidRouting, fmt.Errorf("cannot route [%s] from the config node [%s]", name, r.Name()))
}

// TwoWayRoute adds a new possible outgoing route from the node
func (r *configNode) TwoWayRoute(name string, destination g.Node) error {
	return errors.Join(g.ErrorInvalidRouting, fmt.Errorf("cannot route [%s] from the config node [%s]", name, r.Name()))
}

// ProceedOnFirstRoute proceeds with the first route available
func (r *configNode) ProceedOnFirstRoute(mex f.Message) error {
	return errors.Join(g.ErrorInvalidRouting, fmt.Errorf("cannot route from the config node [%s]", r.Name()))
}
