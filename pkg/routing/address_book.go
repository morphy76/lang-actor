package routing

import (
	"errors"
	"net/url"

	"github.com/morphy76/lang-actor/pkg/common"
)

// ErrorActorAlreadyRegistered is returned when an actor is already registered in the catalog.
var ErrorActorAlreadyRegistered = errors.New("actor already registered")

// ErrorActorNotFound is returned when an actor is not found in the catalog.
var ErrorActorNotFound = errors.New("actor not found")

// Resolver is an interface for resolving addresses to framework.Addressable.
type Resolver interface {
	// Register registers the given URL with the provided Addressable.
	//
	// Parameters:
	//   - addressable (common.Addressable): The Addressable to associate with the URL.
	//
	// Returns:
	//   - (error): An error if the registration fails.
	Register(addressable common.Addressable) error
	// Resolve resolves the given URL to a framework.Addressable.
	//
	// Parameters:
	//   - url (url.URL): The URL to resolve.
	//
	// Returns:
	//   - (common.Addressable): The resolved Addressable.
	//   - (bool): A boolean indicating whether the resolution was successful.
	Resolve(address url.URL) (common.Addressable, bool)
	// Query queries the catalog for Addressables matching the given schema and path parts.
	//
	// Parameters:
	//   - schema (string): The schema to match.
	//   - pathParts (...string): The path parts to match.
	//
	// Returns:
	//   - ([]common.Addressable): A slice of Addressables matching the query.
	Query(schema string, pathParts ...string) []common.Addressable
}

// AddressBook is an interface that defines the methods for a catalog.
type AddressBook interface {
	Resolver
	// Teardown tears down the catalog and releases any resources.
	TearDown()
}
