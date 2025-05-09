package graph

import (
	"fmt"
	"net/url"

	g "github.com/morphy76/lang-actor/pkg/graph"
)

var staticGraphAssertion g.Graph = (*graph)(nil)

type graph struct {
	graphURL   url.URL
	rootNode   g.RootNode
	configNode g.Node
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
