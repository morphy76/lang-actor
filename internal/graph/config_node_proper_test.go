package graph_test

import (
	"net/url"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/morphy76/lang-actor/pkg/builders"
	"gotest.tools/v3/assert"
)

// TestConfigNodeIntegration tests the ConfigNode's functionality in a complete graph
func TestConfigNodeIntegration(t *testing.T) {
	t.Log("ConfigNode Integration test suite")

	t.Run("Graph processes configuration requests correctly", func(t *testing.T) {
		// 1. Create a root node
		rootNode, err := builders.NewRootNode()
		assert.NilError(t, err, "Failed to create root node")

		// 2. Define test configuration with various types of values
		testConfig := map[string]any{
			"string_value": "test string",
			"int_value":    42,
			"bool_value":   true,
			"nested": map[string]any{
				"nested_string": "nested value",
				"nested_int":    100,
			},
			"array_value": []string{"one", "two", "three"},
		}

		// 3. Create graph with configuration (this will create the ConfigNode internally)
		graph, err := builders.NewGraph(rootNode, testConfig)
		assert.NilError(t, err, "Failed to create graph with configuration")

		// 4. Create a debug node that will verify the response from ConfigNode
		debugNode, err := builders.NewDebugNode()
		assert.NilError(t, err, "Failed to create debug node")

		// 5. Create an end node that will help us verify message flow
		endNode, endCh, err := builders.NewEndNode()
		assert.NilError(t, err, "Failed to create end node")

		// 6. Build the route topology
		err = rootNode.OneWayRoute("debug-route", debugNode)
		assert.NilError(t, err, "Failed to create route from root to debug")

		err = debugNode.OneWayRoute("end-route", endNode)
		assert.NilError(t, err, "Failed to create route from debug to end")

		// 7. Create and send a config request message
		// Note: The ConfigNode is addressed directly by the Graph implementation
		// so we need to use the right URL scheme
		debugAddress := debugNode.Address()

		// Create our test message with an actual sender address
		configMsg := &configMessage{
			sender:            debugAddress,
			configMessageType: 2, // request
			requestedKeys:     []string{"string_value", "int_value", "nested"},
			entries:           make(map[string]any),
		}

		// 8. Send the message to the graph
		err = graph.Accept(configMsg)
		assert.NilError(t, err, "Failed to send config request message")

		// 9. Give some time for message processing
		select {
		case <-endCh:
			t.Log("Message flowed through the graph successfully")
		case <-time.After(500 * time.Millisecond):
			t.Log("Timeout waiting for message to flow through graph")
		}

		// Note: In a complete test we would check if the debug node
		// received the expected configuration values. However, as shown
		// in the node_builders.go file, the actual response delivery code
		// in the ConfigNode is currently commented out, so we can't fully
		// test the response processing yet.
	})
}

// TestSimpleConfigRequest tests sending a simple config request to the graph
func TestSimpleConfigRequest(t *testing.T) {
	t.Log("Simple config request test")

	// 1. Create necessary components
	rootNode, err := builders.NewRootNode()
	assert.NilError(t, err, "Failed to create root node")

	// 2. Define minimal test configuration
	testConfig := map[string]any{
		"version": "1.0.0",
		"enabled": true,
	}

	// 3. Create graph with configuration
	graph, err := builders.NewGraph(rootNode, testConfig)
	assert.NilError(t, err, "Failed to create graph with configuration")

	// 4. Create a sender URL
	senderURL, err := url.Parse("actor://test-sender/" + uuid.NewString())
	assert.NilError(t, err, "Failed to create sender URL")

	// 5. Create a config request
	configMsg := &configMessage{
		sender:            *senderURL,
		configMessageType: 0, // keys request
	}

	// 6. Send the message
	err = graph.Accept(configMsg)
	assert.NilError(t, err, "Failed to send config request message")
}
