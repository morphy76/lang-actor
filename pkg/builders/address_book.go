package builders

import (
	"github.com/morphy76/lang-actor/internal/routing"
	r "github.com/morphy76/lang-actor/pkg/routing"
)

// NewAddressBook creates a new actor catalog.
//
// Returns:
//   - (AddressBook): The created actor catalog.
func NewAddressBook() r.AddressBook {
	return routing.NewAddressBook()
}
