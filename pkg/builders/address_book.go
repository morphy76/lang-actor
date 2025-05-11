package builders

import (
	r "github.com/morphy76/lang-actor/internal/routing"
	"github.com/morphy76/lang-actor/pkg/routing"
)

// NewAddressBook creates a new actor catalog.
//
// Returns:
//   - (routing.AddressBook): The created AddressBook instance.
func NewAddressBook() routing.AddressBook {
	return r.NewAddressBook()
}
