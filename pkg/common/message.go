package common

// Message is the interface for messages in the actor model.
type Message interface {
}

// MessageHandler is the interface for handling messages in the actor model.
type MessageHandler interface {
	// Accept processes a message.
	//
	// Parameters:
	//   - msg (Message): The message to be processed.
	//
	// Returns:
	//   - (error): An error if the processing fails, otherwise nil.
	Accept(msg Message) error
}
