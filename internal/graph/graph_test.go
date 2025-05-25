package graph_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/morphy76/lang-actor/internal/graph"
)

func TestNewGraph(t *testing.T) {
	t.Log("Graph Builders test suite")

	t.Run("NewGraph", func(t *testing.T) {
		t.Log("NewGraph test case")

		initialState := make(map[string]any)
		initialState["key1"] = "value1"
		initialState["key2"] = 42

		newGraph, err := graph.NewGraph(uuid.NewString(), initialState, make(map[string]any))
		if err != nil {
			t.Errorf("Expected no error, but got: %v", err)
		}
		if newGraph == nil {
			t.Errorf("Expected a Graph, but got: %v", newGraph)
		}
	})
}
