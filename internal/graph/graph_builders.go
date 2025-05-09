package graph

import (
	"fmt"
	"net/url"

	"github.com/morphy76/lang-actor/internal/routing"
	g "github.com/morphy76/lang-actor/pkg/graph"
)

// NewGraph creates a new instance of the actor graph.
func NewGraph(
	graphName string,
	rootNode g.RootNode,
	configs map[string]any,
) (g.Graph, error) {

	graphURL, err := url.Parse("graph://" + graphName)
	if err != nil {
		return nil, err
	}

	configURL, err := url.Parse("graph://" + graphName + "/config")
	if err != nil {
		return nil, err
	}

	addressBook := routing.NewAddressBook()
	rootNode.Visit(func(node g.Node) {
		addressBook.Register(node.ActorRef())
		node.SetAddressBook(addressBook)
	}, true)
	fmt.Printf("TODO AAAAAAAAA BLOCKED!!! Address book: %v\n", addressBook)

	configNode, err := newConfigNode(configs, *configURL, addressBook)
	if err != nil {
		return nil, err
	}

	graph := &graph{
		graphURL:   *graphURL,
		rootNode:   rootNode,
		configNode: configNode,
	}

	return graph, nil
}
