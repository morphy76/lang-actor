package graph

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/google/uuid"
	"github.com/morphy76/lang-actor/internal/framework"
	f "github.com/morphy76/lang-actor/pkg/framework"
	g "github.com/morphy76/lang-actor/pkg/graph"
	r "github.com/morphy76/lang-actor/pkg/routing"
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

		cfgMessage, ok := msg.(*configMessage)
		if ok {
			for key, val := range cfgMessage.entries {
				fmt.Printf("Debug node received key: %s with value %v\n", key, val)
			}
			baseNode.ProceedOnAnyRoute(msg)
		} else {
			_, err := newConfigMessage(self.Address(), entries)
			if err != nil {
				return self.State(), err
			}
			// TODO every node has to access the address book to send a message
		}

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

var staticConfigMessageAssertion f.Message = (*configMessage)(nil)

type configMessageType int8

const (
	keys configMessageType = iota
	entries
	request
	response
)

type configMessage struct {
	sender            url.URL
	configMessageType configMessageType
	requestedKeys     []string

	requestType configMessageType
	value       any
	keys        []string
	entries     map[string]any
}

// Sender returns the sender of the message.
func (c *configMessage) Sender() url.URL {
	return c.sender
}

// Mutation returns whether the message is a mutation.
func (c *configMessage) Mutation() bool {
	return false
}

func newConfigMessage(sender url.URL, configMessageType configMessageType, requestedKey ...string) (*configMessage, error) {
	if configMessageType == response {
		return nil, fmt.Errorf("TODO: cannot create a config message with response type")
	}

	return &configMessage{
		sender:            sender,
		configMessageType: configMessageType,
		requestedKeys:     requestedKey,
	}, nil
}

func newSingleValueResponse(sender url.URL, key string, value any) f.Message {
	return &configMessage{
		sender:            sender,
		requestType:       request,
		configMessageType: response,
		requestedKeys:     []string{key},
		value:             value,
	}
}

func newListKeysResponse(sender url.URL, configuredKeys []string) f.Message {
	return &configMessage{
		sender:            sender,
		configMessageType: response,
		requestType:       keys,
		keys:              configuredKeys,
	}
}

func newMultivalueResponse(sender url.URL, entries map[string]any) f.Message {
	return &configMessage{
		sender:            sender,
		configMessageType: response,
		requestType:       request,
		entries:           entries,
	}
}

func newListEntriesResponse(sender url.URL, configuredEntries map[string]any) f.Message {
	return &configMessage{
		sender:            sender,
		configMessageType: response,
		requestType:       entries,
		entries:           configuredEntries,
	}
}

func newConfigNode(config map[string]any, address url.URL, addressBook r.AddressBook) (g.Node, error) {

	baseNode := newNode[map[string]any](nil, address)
	baseAddressBook := addressBook
	taskFn := func(msg f.Message, self f.Actor[map[string]any]) (map[string]any, error) {

		if msg == nil {
			return nil, fmt.Errorf("message is nil")
		}
		if _, ok := msg.(*configMessage); !ok {
			return nil, fmt.Errorf("message is not a config message")
		}

		useMex := msg.(*configMessage)
		useState := self.State()

		switch useMex.configMessageType {
		case keys:
			keys := make([]string, 0, len(useState))
			for key := range useState {
				keys = append(keys, key)
			}
			replyMsg := newListKeysResponse(self.Address(), keys)
			addressable, err := baseAddressBook.Lookup(useMex.Sender())
			if err != nil {
				return nil, err
			}
			if err := addressable.Deliver(replyMsg); err != nil {
				return nil, err
			}
		case entries:
			entries := make(map[string]any, len(useState))
			for key, value := range useState {
				entries[key] = value
			}
			replyMsg := newListEntriesResponse(self.Address(), entries)
			addressable, err := baseAddressBook.Lookup(useMex.Sender())
			if err != nil {
				return nil, err
			}
			if err := addressable.Deliver(replyMsg); err != nil {
				return nil, err
			}
		case request:
			if len(useMex.requestedKeys) == 0 {
				shouldReturn, result, err := sendEmptyResponse(self, baseAddressBook, useMex)
				if shouldReturn {
					return result, err
				}
			} else if len(useMex.requestedKeys) == 1 {
				key := useMex.requestedKeys[0]
				if value, ok := useState[key]; ok {
					replyMsg := newSingleValueResponse(self.Address(), key, value)
					addressable, err := baseAddressBook.Lookup(useMex.Sender())
					if err != nil {
						return nil, err
					}
					if err := addressable.Deliver(replyMsg); err != nil {
						return nil, err
					}
				} else {
					shouldReturn, result, err := sendEmptyResponse(self, baseAddressBook, useMex)
					if shouldReturn {
						return result, err
					}
				}
			} else {
				multiValueTmp := make(map[string]any, len(useMex.requestedKeys))
				okys := 0
				for _, key := range useMex.requestedKeys {
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
				addressable, err := baseAddressBook.Lookup(useMex.Sender())
				if err != nil {
					return nil, err
				}
				if err := addressable.Deliver(replyMsg); err != nil {
					return nil, err
				}
			}
		}

		return self.State(), nil
	}

	debugTask, err := framework.NewActor(
		address,
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

func sendEmptyResponse(self f.Actor[map[string]any], addressBook r.AddressBook, useMex *configMessage) (bool, map[string]any, error) {
	emptyPayload := make(map[string]any, 0)
	replyMsg := newMultivalueResponse(self.Address(), emptyPayload)
	addressable, err := addressBook.Lookup(useMex.Sender())
	if err != nil {
		return true, nil, err
	}
	if err := addressable.Deliver(replyMsg); err != nil {
		return true, nil, err
	}
	return false, nil, nil
}
