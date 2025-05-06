package graph_test

import (
	"testing"

	"github.com/morphy76/lang-actor/internal/graph"
)

func TestNewGraph(t *testing.T) {
	t.Log("Graph Builders test suite")

	t.Run("NewGraph", func(t *testing.T) {
		t.Log("NewGraph test case")

		rootNode, err := graph.NewRootNode()
		if err != nil {
			t.Errorf("Expected no error, but got: %v", err)
		}
		if rootNode == nil {
			t.Errorf("Expected a RootNode, but got: %v", rootNode)
		}
	})
}
