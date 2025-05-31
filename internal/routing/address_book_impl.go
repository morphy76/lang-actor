package routing

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"

	c "github.com/morphy76/lang-actor/pkg/common"
	r "github.com/morphy76/lang-actor/pkg/routing"
)

var staticAddressBookAssertion r.AddressBook = (*addressBook)(nil)

type addressBook struct {
	lock *sync.Mutex

	addressables map[url.URL]c.Addressable
}

// Register registers an actor in the addressBook.
func (ab *addressBook) Register(addressable c.Addressable) error {
	ab.lock.Lock()
	defer ab.lock.Unlock()

	if _, exists := ab.addressables[addressable.Address()]; exists {
		return errors.Join(r.ErrorActorAlreadyRegistered, fmt.Errorf("actor [%s] already registered for scheme [%s]", addressable.Address().Host, addressable.Address().Scheme))
	}

	ab.addressables[addressable.Address()] = addressable

	return nil
}

// Lookup looks up an actor in the addressBook by its address.
func (ab *addressBook) Resolve(address url.URL) (c.Addressable, bool) {
	rv, found := ab.addressables[address]
	return rv, found
}

// Query queries the addressBook for actors with a specific scheme and path.
func (ab *addressBook) Query(schema string, pathParts ...string) []c.Addressable {
	ab.lock.Lock()
	defer ab.lock.Unlock()
	rv := make([]c.Addressable, 0, len(ab.addressables))
	for _, addressable := range ab.addressables {
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
func (ab *addressBook) TearDown() {
	ab.lock.Lock()
	defer ab.lock.Unlock()

	for key := range ab.addressables {
		delete(ab.addressables, key)
	}
}

// NewAddressBook creates a new addressBook instance.
func NewAddressBook() r.AddressBook {
	return &addressBook{
		lock: &sync.Mutex{},

		addressables: make(map[url.URL]c.Addressable),
	}
}
