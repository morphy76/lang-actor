package graph

import (
	"net/url"

	f "github.com/morphy76/lang-actor/pkg/framework"
)

var staticStatusMessageAssertion f.Message = (*StatusMessage)(nil)

// StatusMessageType is the type of status message.
type StatusMessageType int8

const (
	// StatusRequest is the type of status message that requests the status.
	StatusRequest = iota
	// StatusResponse is the type of status message that responds to a request.
	StatusResponse
	// StatusUpdate is the type of status message that updates the status.
	StatusUpdate
)

// StatusMessage is the graph status access message.
type StatusMessage struct {
	From              url.URL
	StatusMessageType StatusMessageType

	Value interface{}
}

// Sender returns the sender of the message.
func (c *StatusMessage) Sender() url.URL {
	return c.From
}

// Mutation returns whether the message is a mutation.
func (c *StatusMessage) Mutation() bool {
	return true
}

// NewStatusMessage creates a new status message.
//
// Type parameters:
// - T: The type of the status message.
//
// Parameters:
// - sender: The URL of the sender.
//
// Returns:
// - The created status message.
func NewStatusMessageRequest(
	sender url.URL,
) StatusMessage {
	return StatusMessage{
		From:              sender,
		StatusMessageType: StatusRequest,
	}
}

// NewStatusMessageUpdate creates a new status message.
//
// Parameters:
// - sender: The URL of the sender.
// - value: The value of the status message.
//
// Returns:
// - The created status message.
func NewStatusMessageUpdate(
	sender url.URL,
	value interface{},
) StatusMessage {
	return StatusMessage{
		From:              sender,
		StatusMessageType: StatusUpdate,
		Value:             value,
	}
}

// NewStatusMessageResponse creates a new status message.
//
// Parameters:
// - sender: The URL of the sender.
// - value: The value of the status message.
//
// Returns:
// - The created status message.
func NewStatusMessageResponse(
	sender url.URL,
	value interface{},
) StatusMessage {
	return StatusMessage{
		From:              sender,
		StatusMessageType: StatusResponse,
		Value:             value,
	}
}
