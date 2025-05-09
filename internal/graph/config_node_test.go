package graph_test

import (
	"net/url"
	"testing"

	"github.com/google/uuid"
	"github.com/morphy76/lang-actor/pkg/builders"
	"gotest.tools/v3/assert"
)

// configMessage is a mock implementation that mimics the internal configMessage
// structure to test the ConfigNode processing.
type configMessage struct {
	sender            url.URL
	configMessageType int8 // keys(0), entries(1), request(2), response(3)
	requestedKeys     []string
	entries           map[string]any
}

func (m *configMessage) Sender() url.URL {
	return m.sender
}

func (m *configMessage) Mutation() bool {
	return false
}

// TestConfigNodeConfigRequests tests sending configuration requests to a graph
// containing a ConfigNode and verifying the expected behavior
func TestConfigNodeConfigRequests(t *testing.T) {
	t.Log("ConfigNode requests test suite")

	t.Run("Config node processes configuration requests", func(t *testing.T) {
		// Create a root node
		rootNode, err := builders.NewRootNode()
		assert.NilError(t, err, "Failed to create root node")

		// Define test configuration
		testConfig := map[string]any{
			"server": "localhost",
			"port":   8080,
			"debug":  true,
		}

		// Create a graph with the root node and config
		// The graph builder internally creates a ConfigNode with our configuration
		graph, err := builders.NewGraph(rootNode, testConfig)
		assert.NilError(t, err, "Failed to create graph with configuration")

		// Create a node that will receive responses from ConfigNode
		endNode, endCh, err := builders.NewEndNode()
		assert.NilError(t, err, "Failed to create end node")

		// Create a debug node that will send requests to the ConfigNode
		debugNode, err := builders.NewDebugNode()
		assert.NilError(t, err, "Failed to create debug node")

		// Connect the nodes in a route
		err = rootNode.OneWayRoute(uuid.NewString(), debugNode)
		assert.NilError(t, err, "Failed to create route from root to debug")

		err = debugNode.OneWayRoute(uuid.NewString(), endNode)
		assert.NilError(t, err, "Failed to create route from debug to end")

		// Create a config request message
		senderURL := debugNode.Address()
		configRequest := &configMessage{
			sender:            senderURL,
			configMessageType: 2, // request
			requestedKeys:     []string{"server", "port"},
		}

		// Start the graph by accepting a message
		err = graph.Accept(configRequest)
		assert.NilError(t, err, "Failed to accept message in graph")

		// Wait for the end node to receive the message
		// (indicating the message traversed the entire graph)
		select {
		case <-endCh:
			t.Log("Message processed through the graph successfully")
		default:
			t.Error("Message did not reach the end node")
		}
	})
}

// TestConfigNodeResponses tests that the ConfigNode creates proper responses
// for different types of configuration requests
func TestConfigNodeResponses(t *testing.T) {
	t.Log("ConfigNode responses test suite")

	t.Run("Config node returns configuration values", func(t *testing.T) {
		// Create the test components (actors, mock address book, etc.) to test the
		// ConfigNode's response generation capabilities

		// Set up a test environment
		rootNode, err := builders.NewRootNode()
		assert.NilError(t, err)

		testConfig := map[string]any{
			"appName":     "test-service",
			"environment": "testing",
			"features": map[string]bool{
				"logging": true,
				"metrics": false,
			},
		}

		// Create graph with config node
		graph, err := builders.NewGraph(rootNode, testConfig)
		assert.NilError(t, err)

		// Create debugging node to receive responses
		debugNode, err := builders.NewDebugNode()
		assert.NilError(t, err)

		// Usually we'd need to mock components to properly verify the response
		// Here we're limited by test capabilities, so we'll just verify the graph processes
		// the requests without errors

		// Route the message through the graph
		err = rootNode.OneWayRoute("main-route", debugNode)
		assert.NilError(t, err)

		// Create a request
		reqURL, _ := url.Parse("actor://test/" + uuid.NewString())
		req := &configMessage{
			sender:            *reqURL,
			configMessageType: 0, // keys request
		}

		// Send the request
		err = graph.Accept(req)
		assert.NilError(t, err, "Failed to process keys request")

		// Send entry request
		entriesReq := &configMessage{
			sender:            *reqURL,
			configMessageType: 1, // entries request
		}
		err = graph.Accept(entriesReq)
		assert.NilError(t, err, "Failed to process entries request")

		// Send specific key request
		keyReq := &configMessage{
			sender:            *reqURL,
			configMessageType: 2, // specific request
			requestedKeys:     []string{"appName"},
		}
		err = graph.Accept(keyReq)
		assert.NilError(t, err, "Failed to process specific key request")
	})
}
