package graph

import (
	g "github.com/morphy76/lang-actor/pkg/graph"
)

var staticRootNodeAssertion g.RootNode = (*rootNode)(nil)
var staticEndAssertion g.EndNode = (*endNode)(nil)

type rootNode struct {
	node
}

type endNode struct {
	node
}
