package common

import "net/url"

// Addressable is the interface for the routing layer of the actor model.
type Addressable interface {
	// Actor URI
	//
	// Returns:
	//   - (url.URL): The URL of the actor.
	Address() url.URL
}
