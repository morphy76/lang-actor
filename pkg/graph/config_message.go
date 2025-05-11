package graph

import (
	"fmt"
	"net/url"

	f "github.com/morphy76/lang-actor/pkg/framework"
)

var staticConfigMessageAssertion f.Message = (*ConfigMessage)(nil)

// ConfigMessageType is the type of configuration message.
type ConfigMessageType int8

const (
	// ConfigKeys is the type of configuration message that requests the keys.
	ConfigKeys ConfigMessageType = iota
	// ConfigEntries is the type of configuration message that requests the entries.
	ConfigEntries
	// ConfigRequest is the type of configuration message that requests a value.
	ConfigRequest
	// ConfigResponse is the type of configuration message that responds to a request.
	ConfigResponse
)

// ConfigMessage is a configuration message.
type ConfigMessage struct {
	From              url.URL
	ConfigMessageType ConfigMessageType
	RequestedKeys     []string

	RequestType ConfigMessageType
	Value       any
	Keys        []string
	Entries     map[string]any
}

// Sender returns the sender of the message.
func (c *ConfigMessage) Sender() url.URL {
	return c.From
}

// Mutation returns whether the message is a mutation.
func (c *ConfigMessage) Mutation() bool {
	return false
}

// NewConfigMessage creates a new configuration message.
//
// Parameters:
// - sender: The URL of the sender.
// - ConfigMessageType: The type of configuration message.
// - requestedKey: The keys requested in the message (optional).
// - value: The value of the message (optional).
//
// Returns:
// - A pointer to the created configuration message.
// - An error if the message type is Response.
func NewConfigMessage(
	sender url.URL,
	ConfigMessageType ConfigMessageType,
	requestedKey ...string,
) (*ConfigMessage, error) {
	if ConfigMessageType == ConfigResponse {
		return nil, fmt.Errorf("TODO: cannot create a config message with response type")
	}

	return &ConfigMessage{
		From:              sender,
		ConfigMessageType: ConfigMessageType,
		RequestedKeys:     requestedKey,
	}, nil
}
