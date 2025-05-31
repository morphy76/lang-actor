package graph_test

import (
	"net/url"
	"testing"
	"time"

	"github.com/morphy76/lang-actor/internal/graph"
	b "github.com/morphy76/lang-actor/pkg/builders"
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

	cfg := make(map[string]any)
	cfg["type"] = "test"
	cfg["version"] = "1.0.0"
	cfg["description"] = "test description"
	cfg["author"] = "test author"
	cfg["ts"] = time.Now()

	t.Run("SimpleGraph", func(t *testing.T) {
		t.Log("SimpleGraph test case")

		addressBook := b.NewAddressBook()

		rootNode, err := graph.NewRootNode(nil)
		rootNode.SetResolver(addressBook)
		addressBook.Register(rootNode)
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}

		childNode, err := graph.NewDebugNode(nil)
		childNode.SetResolver(addressBook)
		addressBook.Register(childNode)
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}

		endNode, err := graph.NewEndNode(nil)
		endNode.SetResolver(addressBook)
		addressBook.Register(endNode)
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

		err = rootNode.Accept(&mockMessage{
			sender: rootNode.Address(),
		})
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}
		t.Log("End node received message, process finished")
	})
}
