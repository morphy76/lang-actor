package graph_test

import (
	"testing"

	"github.com/google/uuid"

	"github.com/morphy76/lang-actor/internal/graph"
	g "github.com/morphy76/lang-actor/pkg/graph"
)

const errorNewNodeMessage = "Expected no error, but got: %v"

func TestNewNode(t *testing.T) {
	t.Log("Node Builders test suite")

	t.Run("NewDebugNode", func(t *testing.T) {
		t.Log("NewDebugNode test case")
		node, err := graph.NewDebugNode(nil)
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}
		if node == nil {
			t.Errorf(errorNewNodeMessage, nil)
		}
		if _, ok := node.(g.DebugNode); !ok {
			t.Errorf("Expected a DebugNode, but got: %T", node)
		}
	})
}

func TestNodeRelationships(t *testing.T) {
	t.Log("NodeRelationships test suite")

	t.Run("Add node route", func(t *testing.T) {
		t.Log("Add node route test case")
		startNode, err := graph.NewRootNode(nil)
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}

		endNode, err := graph.NewEndNode(nil)
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}

		childNode1, err := graph.NewDebugNode(nil)
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}

		childNode2, err := graph.NewDebugNode(nil)
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}

		startNode.OneWayRoute(uuid.NewString(), childNode1)
		childNode2.OneWayRoute(uuid.NewString(), endNode)
	})
}
