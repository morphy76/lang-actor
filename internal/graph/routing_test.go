package graph_test

import (
	"net/url"
	"testing"

	"github.com/morphy76/lang-actor/internal/graph"
	f "github.com/morphy76/lang-actor/pkg/framework"
)

var staticMockMessageAssertion f.Message = (*mockMessage)(nil)

type mockMessage struct {
	sender url.URL
}

func (m *mockMessage) Sender() url.URL {
	return m.sender
}
func (m *mockMessage) Mutation() bool {
	return false
}

func TestSimpleGraph(t *testing.T) {
	t.Log("Simple Graph test suite")

	t.Run("SimpleGraph", func(t *testing.T) {
		t.Log("SimpleGraph test case")

		rootNode, err := graph.NewRootNode()
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}

		childNode, err := graph.NewDebugNode()
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}

		endNode, endCh, err := graph.NewEndNode()
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}

		err = rootNode.OneWayRoute("leavingStart", childNode)
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}
		err = childNode.OneWayRoute("leavingDebug", endNode)
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}

		err = rootNode.ProceedOnAnyRoute(&mockMessage{
			sender: rootNode.Address(),
		})
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}
		<-endCh
		t.Log("End node received message, process finished")
	})
}
