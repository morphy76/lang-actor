package graph_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/morphy76/lang-actor/internal/graph"
	b "github.com/morphy76/lang-actor/pkg/builders"
	f "github.com/morphy76/lang-actor/pkg/framework"
	g "github.com/morphy76/lang-actor/pkg/graph"
)

var staticUUIDGraphStateAssertion g.State = (*uUUIDGraphState)(nil)

type uUUIDGraphState struct {
	uuids []any
}

func (s *uUUIDGraphState) AppendGraphState(purpose any, value any) error {
	s.uuids = append(s.uuids, value.(string))
	return nil
}

func TestForkJoinNode(t *testing.T) {
	t.Log("Fork join test suite")

	t.Run("SimpleForkJoin", func(t *testing.T) {
		t.Log("SimpleForkJoin test case")

		testGraph, err := b.NewGraph(&uUUIDGraphState{
			uuids: []any{},
		}, g.NoConfiguration{})
		if err != nil {
			t.Errorf("Error creating graph: %v", err)
		}

		rootNode, err := graph.NewRootNode(testGraph)
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}

		uuids := []string{uuid.NewString(), uuid.NewString(), uuid.NewString()}
		uuidGenFn := func(i int) f.ProcessingFn[g.NodeState] {
			return func(msg f.Message, self f.Actor[g.NodeState]) (g.NodeState, error) {
				rv := uuids[i]
				t.Logf("Processing UUID: %s", rv)
				self.State().Outcome() <- rv
				return self.State(), nil
			}
		}

		childNode, err := graph.NewForkJoingNode(testGraph, false, uuidGenFn(0), uuidGenFn(1), uuidGenFn(2))
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}

		endNode, err := graph.NewEndNode(testGraph)
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}

		err = rootNode.OneWayRoute("leavingStart", childNode)
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}
		err = childNode.OneWayRoute("rejoining", endNode)
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}

		err = rootNode.Accept(&mockMessage{
			sender: rootNode.Address(),
		})
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}
		t.Log("End node received message, process finished")

		if testGraph.State() == nil {
			t.Errorf("Expected uuids to be generated, but got nil")
		}

		actualUUIDs := testGraph.State().(*uUUIDGraphState).uuids
		if len(actualUUIDs) != len(uuids) {
			t.Errorf("Expected %d UUIDs, but got %d", len(uuids), len(actualUUIDs))
		}

		uuidMap := make(map[any]bool)
		for _, u := range uuids {
			uuidMap[u] = true
		}
		for _, actual := range actualUUIDs {
			if !uuidMap[actual] {
				t.Errorf("Actual UUID %s not found in expected uuids list", actual)
			}
		}

		// Negative test: ensure no extra UUIDs are present
		if len(actualUUIDs) > len(uuids) {
			t.Errorf("Received more UUIDs than expected: got %d, want %d", len(actualUUIDs), len(uuids))
		}
	})
}
