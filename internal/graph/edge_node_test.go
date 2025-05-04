package graph_test

import (
	"testing"

	"github.com/morphy76/lang-actor/internal/graph"
)

func TestEdgeNode(t *testing.T) {
	t.Log("Edge nodes test suite")

	t.Run("NewRootNode", func(t *testing.T) {
		t.Log("NewRootNode test case")
		node := graph.NewRootNode()
		if node == nil {
			t.Error("Expected a new root node, but got nil")
		}
	})

	t.Run("NewEndNode", func(t *testing.T) {
		t.Log("NewEndNode test case")
		node := graph.NewEndNode()
		if node == nil {
			t.Error("Expected a new end node, but got nil")
		}
	})
}
