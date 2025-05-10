package graph

import (
	"net/url"

	f "github.com/morphy76/lang-actor/pkg/framework"
)

var staticStatusMessageAssertion f.Message = (*StatusMessage[any])(nil)

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
type StatusMessage[T any] struct {
	From              url.URL
	StatusMessageType StatusMessageType

	Value T
}

// Sender returns the sender of the message.
func (c *StatusMessage[T]) Sender() url.URL {
	return c.From
}

// Mutation returns whether the message is a mutation.
func (c *StatusMessage[T]) Mutation() bool {
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
// - A pointer to the created status message.
func NewStatusMessageRequest[T any](
	sender url.URL,
) (*StatusMessage[T], error) {
	return &StatusMessage[T]{
		From:              sender,
		StatusMessageType: StatusRequest,
	}, nil
}

// NewStatusMessageUpdate creates a new status message.
//
// Type parameters:
// - T: The type of the status message.
//
// Parameters:
// - sender: The URL of the sender.
// - value: The value of the status message.
//
// Returns:
// - A pointer to the created status message.
// - An error if the message type is Response.
func NewStatusMessageUpdate[T any](
	sender url.URL,
	value T,
) (*StatusMessage[T], error) {
	return &StatusMessage[T]{
		From:              sender,
		StatusMessageType: StatusUpdate,
		Value:             value,
	}, nil
}
