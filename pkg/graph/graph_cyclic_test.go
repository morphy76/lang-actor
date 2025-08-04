package graph_test

import (
	"net/url"
	"testing"

	"github.com/google/uuid"

	b "github.com/morphy76/lang-actor/pkg/builders"
	f "github.com/morphy76/lang-actor/pkg/framework"
	g "github.com/morphy76/lang-actor/pkg/graph"
)

func NewCounterNode(forGraph g.Graph) (g.Node, error) {

	address, err := url.Parse("graph://nodes/counter/" + uuid.NewString())
	if err != nil {
		return nil, err
	}

	taskFn := func(msg f.Message, self f.Actor[g.NodeRef]) (g.NodeRef, error) {
		cfg, okCfg := self.State().GraphConfig().(graphConfig)
		graphState, okState := self.State().GraphState().(*graphState)

		if !okCfg || !okState {
			// TODO handle functional error
			self.State().ProceedOntoRoute() <- "leavingCounter"
		}

		if graphState.Counter < cfg.Iterations {
			self.State().GraphState().AppendGraphState(nil, nil)
			self.State().ProceedOntoRoute() <- "iterate"
		} else {
			self.State().ProceedOntoRoute() <- "leavingCounter"
		}

		return self.State(), nil
	}

	return b.NewCustomNode(
		forGraph,
		address,
		taskFn,
		false,
	)
}

var staticGraphConfigAssetion g.Configuration = (*graphConfig)(nil)
var staticGraphStateAssetion g.State = (*graphState)(nil)

type graphState struct {
	Counter int
}

func (s *graphState) AppendGraphState(purpose any, value any) error {
	s.Counter++
	return nil
}

type graphConfig struct {
	Iterations int
}

func TestNewCyclicGraph(t *testing.T) {
	t.Log("Cyclic Graph test suite")

	t.Run("NewCyclicGraph", func(t *testing.T) {
		t.Log("NewCyclicGraph test case")

		state := graphState{
			Counter: 0,
		}

		cfg := graphConfig{
			Iterations: 10,
		}

		graph, err := b.NewGraph(
			&state,
			cfg,
		)
		if err != nil {
			t.Errorf("Error creating graph: %v\n", err)
		}

		rootNode, err := b.NewRootNode(graph)
		if err != nil {
			t.Errorf("Error creating root node: %v\n", err)
		}

		upstreamDebugNode, err := b.NewDebugNode(graph, "upstream")
		if err != nil {
			t.Errorf("Error creating upstream debug node: %v\n", err)
		}

		downstreamDebugNode, err := b.NewDebugNode(graph, "downstream")
		if err != nil {
			t.Errorf("Error creating downstream debug node: %v\n", err)
		}

		counterNode, err := NewCounterNode(graph)
		if err != nil {
			t.Errorf("Error creating counter node: %v\n", err)
		}

		endNode, err := b.NewEndNode(graph)
		if err != nil {
			t.Errorf("Error creating end node: %v\n", err)
		}

		err = rootNode.OneWayRoute("leavingStart", upstreamDebugNode)
		if err != nil {
			t.Errorf("Error creating route from root to child: %v\n", err)
		}
		err = upstreamDebugNode.OneWayRoute("debug", counterNode)
		if err != nil {
			t.Errorf("Error creating route from debug node to child: %v\n", err)
		}
		err = counterNode.OneWayRoute("iterate", counterNode)
		if err != nil {
			t.Errorf("Error creating route from child to itself: %v\n", err)
		}
		err = counterNode.OneWayRoute("leavingCounter", downstreamDebugNode)
		if err != nil {
			t.Errorf("Error creating route from child to debug node: %v\n", err)
		}
		err = downstreamDebugNode.OneWayRoute("debug", endNode)
		if err != nil {
			t.Errorf("Error creating route from debug node to end node: %v\n", err)
		}

		err = rootNode.Accept(uuid.NewString())
		if err != nil {
			t.Errorf("Error accepting message: %v\n", err)
		}

		currentConfig, ok := graph.Config().(graphConfig)
		if !ok {
			t.Errorf("Expected graph config to be of type graphConfig, got %T", graph.Config())
		}

		currentGraphState, ok := graph.State().(*graphState)
		if !ok {
			t.Errorf("Expected graph state to be of type graphState, got %T", graph.State())
		}
		if currentGraphState.Counter != currentConfig.Iterations {
			t.Errorf("Expected counter to be %d, got %d", currentConfig.Iterations, currentGraphState.Counter)
		}
		t.Log("Cyclic Graph test case completed successfully")
	})
}
