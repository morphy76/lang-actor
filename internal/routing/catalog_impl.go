package routing

import (
	"errors"
	"fmt"
	"net/url"
	"sync"

	f "github.com/morphy76/lang-actor/pkg/framework"
	r "github.com/morphy76/lang-actor/pkg/routing"
)

var staticCatalogAssertion r.Catalog = (*catalog)(nil)

type catalog struct {
	lock *sync.Mutex

	actors map[url.URL]f.Transport
}

// Register registers an actor in the catalog.
func (c *catalog) Register(actor f.Transport) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if _, exists := c.actors[actor.Address()]; exists {
		return errors.Join(r.ErrorActorAlreadyRegistered, fmt.Errorf("actor [%s] already registered for scheme [%s]", actor.Address().Host, actor.Address().Scheme))
	}

	c.actors[actor.Address()] = actor

	return nil
}

// Lookup looks up an actor in the catalog by its address.
func (c *catalog) Lookup(address url.URL) (f.Transport, error) {

	rv, found := c.actors[address]
	if !found {
		return nil, errors.Join(r.ErrorActorNotFound, fmt.Errorf("actor [%s] not found for scheme [%s]", address.Host, address.Scheme))
	}

	return rv.(f.Transport), nil
}

// TearDown tears down the catalog and releases any resources.
func (c *catalog) TearDown() {
	c.lock.Lock()
	defer c.lock.Unlock()

	for key := range c.actors {
		delete(c.actors, key)
	}
}

// NewCatalog creates a new catalog instance.
func NewCatalog() r.Catalog {
	return &catalog{
		lock: &sync.Mutex{},

		actors: make(map[url.URL]f.Transport),
	}
}
