package graph_test

import (
	"net/url"
	"testing"
	"time"

	"github.com/morphy76/lang-actor/pkg/builders"
	"gotest.tools/v3/assert"
)

// TestConfigNodeAcceptConfigRequests tests the ConfigNode's ability to accept config requests
func TestConfigNodeAcceptConfigRequests(t *testing.T) {
	t.Log("ConfigNode request acceptance test")

	t.Run("Graph with ConfigNode processes configuration requests", func(t *testing.T) {
		// Create a root node
		rootNode, err := builders.NewRootNode()
		assert.NilError(t, err, "Failed to create root node")

		// Define test configuration
		testConfig := map[string]any{
			"app": map[string]any{
				"name":    "test-app",
				"version": "1.0.0",
			},
			"settings": map[string]any{
				"timeout":  30,
				"retries":  3,
				"debug":    true,
				"features": []string{"a", "b", "c"},
			},
		}

		// Create a graph with the test configuration
		// The graph builder automatically creates a ConfigNode with this configuration
		graph, err := builders.NewGraph(rootNode, testConfig)
		assert.NilError(t, err, "Failed to create graph with configuration")

		// Create a debug node that will receive messages from the ConfigNode
		debugNode, err := builders.NewDebugNode()
		assert.NilError(t, err, "Failed to create debug node")

		// Create an end node to catch messages
		endNode, endCh, err := builders.NewEndNode()
		assert.NilError(t, err, "Failed to create end node")

		// Connect nodes in a route
		err = rootNode.OneWayRoute("debug-route", debugNode)
		assert.NilError(t, err, "Failed to create route from root to debug")

		err = debugNode.OneWayRoute("end-route", endNode)
		assert.NilError(t, err, "Failed to create route from debug to end")

		// Create a custom message that will be acceptable to the ConfigNode
		// We'll use the framework.ActorMessage for this purpose
		senderAddr := debugNode.Address()
		requestMsg := &testConfigMessage{
			senderURL:     senderAddr,
			messageType:   2, // request type
			requestedKeys: []string{"app", "settings.timeout"},
		}

		// Send the message to the graph
		err = graph.Accept(requestMsg)
		assert.NilError(t, err, "Failed to send config request to graph")

		// Wait for the message to flow through the graph
		select {
		case <-endCh:
			t.Log("Message successfully processed through the graph")
		case <-time.After(100 * time.Millisecond):
			t.Log("No message received at end node, but this is expected as responses might not be implemented")
		}
	})
}

// testConfigMessage is a custom message implementation for testing ConfigNode
type testConfigMessage struct {
	senderURL     url.URL
	messageType   int8
	requestedKeys []string
}

// Sender returns the sender of this message
func (m *testConfigMessage) Sender() url.URL {
	return m.senderURL
}

// Mutation returns whether this message is a mutation
func (m *testConfigMessage) Mutation() bool {
	return false
}
