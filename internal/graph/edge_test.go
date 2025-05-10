package graph_test

import (
	"net/url"
	"testing"
	"time"

	"github.com/morphy76/lang-actor/internal/graph"
	"github.com/morphy76/lang-actor/pkg/builders"
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

		addressBook := builders.NewAddressBook()

		cfgNode, err := graph.NewConfigNode(cfg, "testGraph")
		if err != nil {
			t.Errorf("Error creating config node: %v", err)
		}
		cfgNode.SetResolver(addressBook)
		addressBook.Register(cfgNode)
		addressBook.Register(cfgNode.ActorRef())

		rootNode, err := graph.NewRootNode()
		rootNode.SetResolver(addressBook)
		addressBook.Register(rootNode)
		addressBook.Register(rootNode.ActorRef())
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}

		childNode, err := graph.NewDebugNode()
		childNode.SetResolver(addressBook)
		addressBook.Register(childNode)
		addressBook.Register(childNode.ActorRef())
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}

		endNode, endCh, err := graph.NewEndNode()
		endNode.SetResolver(addressBook)
		addressBook.Register(endNode)
		addressBook.Register(endNode.ActorRef())
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
