package routing

import (
	"errors"
	"net/url"

	"github.com/morphy76/lang-actor/pkg/framework"
)

// ErrorActorAlreadyRegistered is returned when an actor is already registered in the catalog.
var ErrorActorAlreadyRegistered = errors.New("actor already registered")

// ErrorActorNotFound is returned when an actor is not found in the catalog.
var ErrorActorNotFound = errors.New("actor not found")

// AddressBook is an interface that defines the methods for a catalog.
type AddressBook interface {
	// Register registers an actor in the catalog.
	//
	// Parameters:
	//   - actor: The actor to register.
	//
	// Returns:
	//   - error: An error if the registration fails.
	Register(actor framework.Transport) error
	// Lookup looks up an actor in the catalog by its address.
	//
	// Parameters:
	//   - url: The address of the actor to look up.
	//
	// Returns:
	//   - actor: The actor if found, nil otherwise.
	//   - error: An error if the lookup fails.
	Lookup(address url.URL) (framework.Transport, error)
	// Teardown tears down the catalog and releases any resources.
	TearDown()
}
