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
		node, err := graph.NewDebugNode()
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
		startNode, err := graph.NewRootNode()
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}

		endNode, err := graph.NewEndNode()
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}

		childNode1, err := graph.NewDebugNode()
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}

		childNode2, err := graph.NewDebugNode()
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}

		startNode.OneWayRoute(uuid.NewString(), childNode1)
		childNode1.TwoWayRoute(uuid.NewString(), childNode2)
		childNode2.OneWayRoute(uuid.NewString(), endNode)
	})

	t.Run("Add and verify a two way route", func(t *testing.T) {
		t.Log("Add and verify a two way route test case")

		childNode1, err := graph.NewDebugNode()
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}

		childNode2, err := graph.NewDebugNode()
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}

		err = childNode1.TwoWayRoute(uuid.NewString(), childNode2)
		if err != nil {
			t.Errorf("Expected no error, but got: %v", err)
		}

		routes1 := childNode1.RouteNames()
		routes2 := childNode2.RouteNames()

		if len(routes1) != 1 {
			t.Errorf("Expected 1 route, but got: %d", len(routes1))
		}
		if len(routes2) != 1 {
			t.Errorf("Expected 1 route, but got: %d", len(routes2))
		}
		if routes2[0] != "inverse-"+routes1[0] {
			t.Errorf("Expected route2 name to be 'inverse-' of route1 name, but got: %s and %s", routes1[0], routes2[0])
		}
	})
}
