package graph_test

import (
	"net/url"
	"testing"

	"github.com/google/uuid"
	"github.com/morphy76/lang-actor/internal/graph"
	b "github.com/morphy76/lang-actor/pkg/builders"
	f "github.com/morphy76/lang-actor/pkg/framework"
	g "github.com/morphy76/lang-actor/pkg/graph"
)

const errorAddressMessage = "Error generating address: %v"

func TestForkAndThenJoinNode(t *testing.T) {
	t.Log("Fork and then join test suite")

	t.Run("SimpleForkAndThenJoin", func(t *testing.T) {
		t.Log("SimpleForkAndThenJoin test case")

		testGraph, err := b.NewGraph(&uUUIDGraphState{
			uuids: []string{},
		}, g.NoConfiguration{})
		if err != nil {
			t.Errorf("Error creating graph: %v", err)
		}

		rootNode, err := graph.NewRootNode(testGraph)
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}

		forkNode, err := graph.NewForkNode(testGraph)
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}

		uuids := []string{uuid.NewString(), uuid.NewString(), uuid.NewString()}
		uuidGenFn := func(i int) f.ProcessingFn[g.NodeState] {
			return func(msg f.Message, self f.Actor[g.NodeState]) (g.NodeState, error) {
				rv := uuids[i]
				t.Logf("Processing UUID: %s", rv)
				self.State().GraphState().AppendGraphState(nil, rv)
				self.State().Outcome() <- g.WhateverOutcome
				return self.State(), nil
			}
		}
		addr1, err := genAddress(1)
		if err != nil {
			t.Errorf(errorAddressMessage, err)
		}
		branch1, err := graph.NewCustomNode(testGraph, addr1, uuidGenFn(0), true)
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}
		addr2, err := genAddress(2)
		if err != nil {
			t.Errorf(errorAddressMessage, err)
		}
		branch2, err := graph.NewCustomNode(testGraph, addr2, uuidGenFn(1), true)
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}
		addr3, err := genAddress(3)
		if err != nil {
			t.Errorf(errorAddressMessage, err)
		}
		branch3, err := graph.NewCustomNode(testGraph, addr3, uuidGenFn(2), true)
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}

		joinNode, err := graph.NewJoinNode(testGraph, forkNode)
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}

		debugNode1, err := graph.NewDebugNode(testGraph, "debug1")
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}

		debugNode2, err := graph.NewDebugNode(testGraph, "debug2")
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}

		debugNode3, err := graph.NewDebugNode(testGraph, "debug3")
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}

		finalDebugNode, err := graph.NewDebugNode(testGraph, "final")
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}

		endNode, err := graph.NewEndNode(testGraph)
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}

		err = rootNode.OneWayRoute("leavingStart", forkNode)
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}
		err = forkNode.OneWayRoute("branch1", branch1)
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}
		err = forkNode.OneWayRoute("branch2", branch2)
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}
		err = forkNode.OneWayRoute("branch3", branch3)
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}
		err = branch1.OneWayRoute("toJoin", debugNode1)
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}
		err = branch2.OneWayRoute("toJoin", debugNode2)
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}
		err = branch3.OneWayRoute("toJoin", debugNode3)
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}
		err = debugNode1.OneWayRoute("toJoin", joinNode)
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}
		err = debugNode2.OneWayRoute("toJoin", joinNode)
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}
		err = debugNode3.OneWayRoute("toJoin", joinNode)
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}
		err = joinNode.OneWayRoute("rejoining", finalDebugNode)
		if err != nil {
			t.Errorf(errorNewNodeMessage, err)
		}
		err = finalDebugNode.OneWayRoute("toEnd", endNode)
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

		uuidMap := make(map[string]bool)
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

func genAddress(i int) (*url.URL, error) {
	addr := "graph://nodes/forkjoin/test-" + uuid.NewString()
	address, err := url.Parse(addr)
	if err != nil {
		return &url.URL{}, err
	}
	return address, nil
}
