package routing

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/morphy76/lang-actor/pkg/common"
	"github.com/morphy76/lang-actor/pkg/routing"
)

var staticAddressBookAssertion routing.AddressBook = (*addressBook)(nil)

type addressBook struct {
	lock *sync.Mutex

	addressables map[url.URL]common.Addressable
}

// Register registers an actor in the addressBook.
func (c *addressBook) Register(addressable common.Addressable) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if _, exists := c.addressables[addressable.Address()]; exists {
		return errors.Join(routing.ErrorActorAlreadyRegistered, fmt.Errorf("actor [%s] already registered for scheme [%s]", addressable.Address().Host, addressable.Address().Scheme))
	}

	c.addressables[addressable.Address()] = addressable

	return nil
}

// Lookup looks up an actor in the addressBook by its address.
func (c *addressBook) Resolve(address url.URL) (common.Addressable, bool) {
	rv, found := c.addressables[address]
	return rv, found
}

// Query queries the addressBook for actors with a specific scheme and path.
func (c *addressBook) Query(schema string, pathParts ...string) []common.Addressable {
	c.lock.Lock()
	defer c.lock.Unlock()
	rv := make([]common.Addressable, 0, len(c.addressables))
	for _, addressable := range c.addressables {
		if addressable.Address().Scheme == schema {
			if len(pathParts) == 0 {
				rv = append(rv, addressable)
				continue
			}

			addressableParts := strings.Split(addressable.Address().Path, "/")
			addressableParts[0] = addressable.Address().Host
			allMatch := true
			for idx, part := range pathParts {
				allMatch = allMatch && addressableParts[idx] == part
				if !allMatch {
					break
				}
			}
			if allMatch {
				rv = append(rv, addressable)
			}
		}
	}
	return rv
}

// TearDown tears down the addressBook and releases any resources.
func (c *addressBook) TearDown() {
	c.lock.Lock()
	defer c.lock.Unlock()

	for key := range c.addressables {
		delete(c.addressables, key)
	}
}

// NewAddressBook creates a new addressBook instance.
func NewAddressBook() routing.AddressBook {
	return &addressBook{
		lock: &sync.Mutex{},

		addressables: make(map[url.URL]common.Addressable),
	}
}
