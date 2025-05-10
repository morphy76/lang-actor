package graph

import (
	"fmt"
	"net/url"

	f "github.com/morphy76/lang-actor/pkg/framework"
	g "github.com/morphy76/lang-actor/pkg/graph"
)

var staticGraphAssertion g.Graph = (*graph)(nil)

type graph struct {
	resolvables map[url.URL]f.Addressable
	graphURL    url.URL
	rootNode    g.RootNode
	configNode  g.Node
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
	_, found := g.resolvables[addressable.Address()]
	if found {
		return fmt.Errorf("addressable already registered")
	}
	g.resolvables[addressable.Address()] = addressable
	return nil
}

// Resolve resolves the given URL to a framework.Addressable.
func (g *graph) Resolve(address url.URL) (f.Addressable, bool) {
	rv, found := g.resolvables[address]
	return rv, found
}
