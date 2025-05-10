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

	actors map[url.URL]f.Addressable
}

// Register registers an actor in the addressBook.
func (c *addressBook) Register(addressable framework.Addressable) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if _, exists := c.actors[addressable.Address()]; exists {
		return errors.Join(r.ErrorActorAlreadyRegistered, fmt.Errorf("actor [%s] already registered for scheme [%s]", addressable.Address().Host, addressable.Address().Scheme))
	}

	c.actors[addressable.Address()] = addressable

	return nil
}

// Lookup looks up an actor in the addressBook by its address.
func (c *addressBook) Resolve(address url.URL) (f.Addressable, bool) {
	rv, found := c.actors[address]
	return rv, found
}

// TearDown tears down the addressBook and releases any resources.
func (c *addressBook) TearDown() {
	c.lock.Lock()
	defer c.lock.Unlock()

	for key := range c.actors {
		delete(c.actors, key)
	}
}

// NewAddressBook creates a new addressBook instance.
func NewAddressBook() r.AddressBook {
	return &addressBook{
		lock: &sync.Mutex{},

		actors: make(map[url.URL]f.Addressable),
	}
}
