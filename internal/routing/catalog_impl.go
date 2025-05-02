package routing

import (
	"errors"
	"fmt"
	"net/url"
	"sync"

	f "github.com/morphy76/lang-actor/pkg/framework"
	r "github.com/morphy76/lang-actor/pkg/routing"
)

var staticAddressBookAssertion r.AddressBook = (*addressBook)(nil)

type addressBook struct {
	lock *sync.Mutex

	actors map[url.URL]f.Addressable
}

// Register registers an actor in the addressBook.
func (c *addressBook) Register(actor f.Addressable) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if _, exists := c.actors[actor.Address()]; exists {
		return errors.Join(r.ErrorActorAlreadyRegistered, fmt.Errorf("actor [%s] already registered for scheme [%s]", actor.Address().Host, actor.Address().Scheme))
	}

	c.actors[actor.Address()] = actor

	return nil
}

// Lookup looks up an actor in the addressBook by its address.
func (c *addressBook) Lookup(address url.URL) (f.Addressable, error) {

	rv, found := c.actors[address]
	if !found {
		return nil, errors.Join(r.ErrorActorNotFound, fmt.Errorf("actor [%s] not found for scheme [%s]", address.Host, address.Scheme))
	}

	return rv.(f.Addressable), nil
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
