package graph_test

import (
	"testing"

	"github.com/morphy76/lang-actor/internal/graph"
)

func TestNewNode(t *testing.T) {
	t.Log("NewNode test suite")

	t.Run("NewNode", func(t *testing.T) {
		t.Log("NewNode test case")
		node := graph.NewNode()
		if node == nil {
			t.Error("Expected a new node, but got nil")
		}
	})
}

func TestNodeRelationships(t *testing.T) {
	t.Log("NodeRelationships test suite")

	t.Run("Append a child", func(t *testing.T) {
		t.Log("Append test case")
		node := graph.NewNode()
		if node == nil {
			t.Error("Expected a new node, but got nil")
		}

		childNode := graph.NewNode()
		if childNode == nil {
			t.Error("Expected a new child node, but got nil")
		}
		node.Append(childNode)
	})
}
