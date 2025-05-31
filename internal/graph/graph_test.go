package graph_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/morphy76/lang-actor/internal/graph"
)

type graphState struct {
	stateAsMap map[string]any
}

func (s graphState) AppendGraphState(purpose any, value any) error {
	return nil
}

func TestNewGraph(t *testing.T) {
	t.Log("Graph Builders test suite")

	t.Run("NewGraph", func(t *testing.T) {
		t.Log("NewGraph test case")

		initialStateMap := make(map[string]any)
		initialStateMap["key1"] = "value1"
		initialStateMap["key2"] = 42

		initialState := graphState{stateAsMap: initialStateMap}

		newGraph, err := graph.NewGraph(uuid.NewString(), initialState, make(map[string]any))
		if err != nil {
			t.Errorf("Expected no error, but got: %v", err)
		}
		if newGraph == nil {
			t.Errorf("Expected a Graph, but got: %v", newGraph)
		}
	})
}
