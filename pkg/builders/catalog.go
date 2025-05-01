package builders

import (
	"github.com/morphy76/lang-actor/internal/routing"
	r "github.com/morphy76/lang-actor/pkg/routing"
)

// NewActorCatalog creates a new actor catalog.
//
// Returns:
//   - (Catalog): The created actor catalog.
func NewActorCatalog() r.Catalog {
	return routing.NewCatalog()
}
