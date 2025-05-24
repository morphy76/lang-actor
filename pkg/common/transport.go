package common

// Transport is the interface for the transport layer of the actor model.
type Transport interface {
	// Deliver a message to the actor
	//
	// Parameters:
	//   - msg (Message): The message to be delivered.
	//
	// Returns:
	//   - (error): An error if the delivery fails, otherwise nil.
	Deliver(msg Message) error
	// Send is a function to send messages to other actors.
	//
	// Parameters:
	//   - msg (Message): The message to be sent.
	//   - destination (Transport): The addressable actor to which the message is sent.
	//
	// Returns:
	//   - (error): An error if the sending fails, otherwise nil.
	Send(msg Message, destination Transport) error
}
