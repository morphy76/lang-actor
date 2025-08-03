package common

// Transport is the interface for the transport layer of the actor model.
type Transport interface {
	// Deliver a message to the actor
	//
	// Parameters:
	//   - msg (any): The message to be delivered.
	//   - from (Addressable): The addressable actor from which the message is sent.
	//
	// Returns:
	//   - (error): An error if the delivery fails, otherwise nil.
	Deliver(msg any, from Addressable) error
	// Send is a function to send messages to other actors.
	//
	// Parameters:
	//   - msg (any): The message to be sent.
	//   - destination (Transport): The addressable actor to which the message is sent.
	//
	// Returns:
	//   - (error): An error if the sending fails, otherwise nil.
	Send(msg any, destination Transport) error
}
