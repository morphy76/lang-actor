package graph

import (
	"fmt"
	"net/url"

	f "github.com/morphy76/lang-actor/pkg/framework"
	g "github.com/morphy76/lang-actor/pkg/graph"
	r "github.com/morphy76/lang-actor/pkg/routing"
)

var staticGraphAssertion g.Graph = (*graph)(nil)

type graph struct {
	resolvables map[url.URL]*f.Addressable
	graphURL    url.URL
	rootNode    g.RootNode
	configNode  g.Node
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

// Accept accepts a todo item and proceeds it on the root node.
func (g *graph) Accept(todo any) error {
	if g.rootNode == nil {
		// TODO: handle error
		return fmt.Errorf("TODO error")
	}

	mex := &acceptedMessage{
		sender: g.graphURL,
	}

	if err := g.rootNode.ProceedOnAnyRoute(mex); err != nil {
		return err
	}

	return nil
}

// Register registers the given URL with the provided Addressable.
func (g *graph) Register(addressable f.Addressable) error {
	return g.addressBook.Register(addressable)
}

// Resolve resolves the given URL to a framework.Addressable.
func (g *graph) Resolve(address url.URL) (f.Addressable, bool) {
	return g.addressBook.Resolve(address)
}

// Query queries the address book for the given schema and path parts.
func (g *graph) Query(schema string, pathParts ...string) []f.Addressable {
	return g.addressBook.Query(schema, pathParts...)
}
