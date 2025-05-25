package graph_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/morphy76/lang-actor/internal/graph"
	g "github.com/morphy76/lang-actor/pkg/graph"
)

func TestNewEdgeNode(t *testing.T) {
	t.Log("Edge Node Builders test suite")

	t.Run("NewRootNode", func(t *testing.T) {
		t.Log("NewRootNode test case")
		node, err := graph.NewRootNode(nil)
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}
		if node == nil {
			t.Errorf(errorNewNodeMessage, nil)
		}
		if _, ok := node.(g.RootNode); !ok {
			t.Errorf("Expected a RootNode, but got: %T", node)
		}
	})

	t.Run("NewEndNode", func(t *testing.T) {
		t.Log("NewEndNode test case")
		node, err := graph.NewEndNode(nil)
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}
		if node == nil {
			t.Errorf(errorNewNodeMessage, nil)
		}
		if _, ok := node.(g.EndNode); !ok {
			t.Errorf("Expected an EndNode, but got: %T", node)
		}
	})
}

func TestRootNodeRelationships(t *testing.T) {
	t.Log("RootNodeRelationships test suite")

	t.Run("RootNode can have a single oneway route", func(t *testing.T) {
		t.Log("RootNode single oneway route test case")
		rootNode, err := graph.NewRootNode(nil)
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}

		childNode, err := graph.NewDebugNode(nil)
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}

		err = rootNode.OneWayRoute(uuid.NewString(), childNode)
		if err != nil {
			t.Errorf("Expected no error, but got: %v", err)
		}

		anotherChildNode, err := graph.NewDebugNode(nil)
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}

		err = rootNode.OneWayRoute(uuid.NewString(), anotherChildNode)
		if err == nil {
			t.Errorf("Expected an error when adding a second oneway route to RootNode, but got none")
		}
	})
}

func TestEndNodeRelationships(t *testing.T) {
	t.Log("EndNodeRelationships test suite")

	t.Run("EndNode cannot have routes", func(t *testing.T) {
		t.Log("EndNode route restriction test case")
		endNode, err := graph.NewEndNode(nil)
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}

		childNode, err := graph.NewDebugNode(nil)
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}

		err = endNode.OneWayRoute(uuid.NewString(), childNode)
		if err == nil {
			t.Errorf("Expected an error when adding a oneway route to EndNode, but got none")
		}
	})
}
