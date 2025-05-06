package graph

import (
	g "github.com/morphy76/lang-actor/pkg/graph"
)

var staticGraphAssertion g.Graph = (*graph)(nil)

type graph struct {
	rootNode g.RootNode
}
