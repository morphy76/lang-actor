package graph

import (
	"net/url"

	c "github.com/morphy76/lang-actor/pkg/common"
	g "github.com/morphy76/lang-actor/pkg/graph"
	r "github.com/morphy76/lang-actor/pkg/routing"
)

var staticGraphAssertion g.Graph = (*graph)(nil)

type graph struct {
	resolvables map[url.URL]*c.Addressable
	graphURL    url.URL
	config      g.Configuration
	status      g.State
	addressBook r.AddressBook
}

type acceptedMessage struct {
	sender url.URL
}

func (m *acceptedMessage) Sender() url.URL {
	return m.sender
}
func (m *acceptedMessage) Mutation() bool {
	return false
}

// Register registers the given URL with the provided Addressable.
func (g *graph) Register(addressable c.Addressable) error {
	return g.addressBook.Register(addressable)
}

// Resolve resolves the given URL to a framework.Addressable.
func (g *graph) Resolve(address url.URL) (c.Addressable, bool) {
	return g.addressBook.Resolve(address)
}

// Query queries the address book for the given schema and path parts.
func (g *graph) Query(schema string, pathParts ...string) []c.Addressable {
	return g.addressBook.Query(schema, pathParts...)
}

// State returns the current state of the graph.
func (g *graph) State() g.State {
	return g.status
}

// Config returns the configuration of the graph.
func (g *graph) Configuration() g.Configuration {
	return g.config
}
