package graph_test

import (
	"net/url"
	"testing"

	"github.com/google/uuid"
	"github.com/morphy76/lang-actor/internal/graph"
	"github.com/morphy76/lang-actor/pkg/builders"
	f "github.com/morphy76/lang-actor/pkg/framework"
	g "github.com/morphy76/lang-actor/pkg/graph"
)

func NewCounterNode(forGraph g.Graph) (g.Node, error) {

	address, err := url.Parse("graph://nodes/counter/" + uuid.NewString())
	if err != nil {
		return nil, err
	}

	taskFn := func(msg f.Message, self f.Actor[g.NodeState]) (g.NodeState, error) {
		cfg, okCfg := self.State().GraphConfig.(graphConfig)
		graphState, okState := self.State().GraphState.(graphState)

		if !okCfg || !okState {
			// TODO handle functional error
			self.State().Outcome <- "leavingCounter"
		}

		if graphState.Counter < cfg.Iterations {
			graphState.Counter++
			self.State().Outcome <- "iterate"
		} else {
			self.State().Outcome <- "leavingCounter"
		}

		return self.State(), nil
	}

	return graph.NewCustomNode(
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

		graph, err := builders.NewGraph(
			state,
			cfg,
		)
		if err != nil {
			t.Errorf("Error creating graph: %v\n", err)
		}

		rootNode, err := builders.NewRootNode(graph)
		if err != nil {
			t.Errorf("Error creating root node: %v\n", err)
		}

		upstreamDebugNode, err := builders.NewDebugNode(graph, "upstream")
		if err != nil {
			t.Errorf("Error creating upstream debug node: %v\n", err)
		}

		downstreamDebugNode, err := builders.NewDebugNode(graph, "downstream")
		if err != nil {
			t.Errorf("Error creating downstream debug node: %v\n", err)
		}

		counterNode, err := NewCounterNode(graph)
		if err != nil {
			t.Errorf("Error creating counter node: %v\n", err)
		}

		endNode, err := builders.NewEndNode(graph)
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

		if state.Counter != cfg.Iterations {
			t.Errorf("Expected counter to be %d, got %d", cfg.Iterations, state.Counter)
		}
		t.Log("Cyclic Graph test case completed successfully")
	})
}
