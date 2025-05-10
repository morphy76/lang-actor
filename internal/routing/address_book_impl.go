package routing

import (
	"errors"
	"fmt"
	"net/url"
	"sync"

	"github.com/morphy76/lang-actor/pkg/framework"
	f "github.com/morphy76/lang-actor/pkg/framework"
	r "github.com/morphy76/lang-actor/pkg/routing"
)

var staticAddressBookAssertion r.AddressBook = (*addressBook)(nil)

type addressBook struct {
	lock *sync.Mutex

	addressables map[url.URL]f.Addressable
}

// Register registers an actor in the addressBook.
func (c *addressBook) Register(addressable framework.Addressable) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if _, exists := c.addressables[addressable.Address()]; exists {
		return errors.Join(r.ErrorActorAlreadyRegistered, fmt.Errorf("actor [%s] already registered for scheme [%s]", addressable.Address().Host, addressable.Address().Scheme))
	}

	c.addressables[addressable.Address()] = addressable

	return nil
}

// Lookup looks up an actor in the addressBook by its address.
func (c *addressBook) Resolve(address url.URL) (f.Addressable, bool) {
	rv, found := c.addressables[address]
	return rv, found
}

func (c *addressBook) Query(schema string, pathParts ...string) []f.Addressable {
	c.lock.Lock()
	defer c.lock.Unlock()
	rv := make([]f.Addressable, 0, len(c.addressables))
	for _, addressable := range c.addressables {
		if addressable.Address().Scheme == schema {
			if len(pathParts) == 0 {
				rv = append(rv, addressable)
				continue
			}

			path := addressable.Address().Path
			for _, part := range pathParts {
				if path == part {
					rv = append(rv, addressable)
					break
				}
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
func NewAddressBook() r.AddressBook {
	return &addressBook{
		lock: &sync.Mutex{},

		addressables: make(map[url.URL]f.Addressable),
	}
}
