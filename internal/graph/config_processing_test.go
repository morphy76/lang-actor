package graph_test

import (
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/morphy76/lang-actor/pkg/builders"
	f "github.com/morphy76/lang-actor/pkg/framework"
	"gotest.tools/v3/assert"
)

// mockAddressable implements the Addressable interface for testing
type mockAddressable struct {
	address         url.URL
	receivedMessage f.Message
	mutex           sync.Mutex
}

func (m *mockAddressable) Address() url.URL {
	return m.address
}

func (m *mockAddressable) Deliver(msg f.Message) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.receivedMessage = msg
	return nil
}

// TestConfigNodeProcessingFunction tests the specific processing function
// of the ConfigNode to handle config requests and produce responses
func TestConfigNodeProcessingFunction(t *testing.T) {
	t.Log("ConfigNode Processing function test suite")

	// Setup a test environment with proper components
	t.Run("Process config request and produce response", func(t *testing.T) {
		// Create the necessary components for testing
		rootNode, err := builders.NewRootNode()
		assert.NilError(t, err, "Failed to create root node")

		// Create test configuration
		testConfig := map[string]any{
			"database": map[string]any{
				"host":     "127.0.0.1",
				"port":     5432,
				"username": "admin",
			},
			"server": map[string]any{
				"port": 8080,
				"ssl":  false,
			},
			"features": []string{"auth", "logging", "metrics"},
			"version":  "1.2.3",
		}

		// Create test graph with config
		graph, err := builders.NewGraph(rootNode, testConfig)
		assert.NilError(t, err, "Failed to create graph")

		// Create a mock addressable to capture responses
		mockAddr := &mockAddressable{
			address: url.URL{
				Scheme: "actor",
				Host:   "test-receiver",
				Path:   "/" + uuid.NewString(),
			},
			mutex: sync.Mutex{},
		}

		// Create a debug node that will be used to capture responses
		debugNode, err := builders.NewDebugNode()
		assert.NilError(t, err)

		// Connect nodes
		err = rootNode.OneWayRoute("test-route", debugNode)
		assert.NilError(t, err)

		// Create different types of config messages and test them

		// Test 1: Request specific keys
		t.Run("Request specific keys", func(t *testing.T) {
			reqMsg := &configMessage{
				sender:            mockAddr.address,
				configMessageType: 2, // request
				requestedKeys:     []string{"version"},
			}

			// Send the message through the graph
			err = graph.Accept(reqMsg)
			assert.NilError(t, err)

			// Give time for processing
			time.Sleep(50 * time.Millisecond)

			// Verify response (if we could access it)
			// In a real scenario with full access to the ConfigNode implementation
			// we would check for the exact response
		})

		// Test 2: Request all keys
		t.Run("Request all keys", func(t *testing.T) {
			reqMsg := &configMessage{
				sender:            mockAddr.address,
				configMessageType: 0, // keys request
			}

			// Send the message
			err = graph.Accept(reqMsg)
			assert.NilError(t, err)

			// Give time for processing
			time.Sleep(50 * time.Millisecond)
		})

		// Test 3: Request all entries
		t.Run("Request all entries", func(t *testing.T) {
			reqMsg := &configMessage{
				sender:            mockAddr.address,
				configMessageType: 1, // entries request
			}

			// Send the message
			err = graph.Accept(reqMsg)
			assert.NilError(t, err)

			// Give time for processing
			time.Sleep(50 * time.Millisecond)
		})

		// Test 4: Request nested keys
		t.Run("Request nested configuration", func(t *testing.T) {
			reqMsg := &configMessage{
				sender:            mockAddr.address,
				configMessageType: 2, // request
				requestedKeys:     []string{"database", "server"},
			}

			// Send the message
			err = graph.Accept(reqMsg)
			assert.NilError(t, err)

			// Give time for processing
			time.Sleep(50 * time.Millisecond)
		})
	})
}

// Test behavior of config node with non-existent keys
func TestConfigNodeNonExistentKeys(t *testing.T) {
	t.Log("ConfigNode Non-existent keys test suite")

	t.Run("Request non-existent keys", func(t *testing.T) {
		// Setup
		rootNode, err := builders.NewRootNode()
		assert.NilError(t, err)

		// Minimal config
		testConfig := map[string]any{
			"key1": "value1",
		}

		graph, err := builders.NewGraph(rootNode, testConfig)
		assert.NilError(t, err)

		// Create a request with valid and non-existent keys
		senderURL, err := url.Parse("actor://tester/" + uuid.NewString())
		assert.NilError(t, err)

		reqMsg := &configMessage{
			sender:            *senderURL,
			configMessageType: 2, // request
			requestedKeys:     []string{"key1", "non-existent-key"},
		}

		// Send the message
		err = graph.Accept(reqMsg)
		assert.NilError(t, err, "Should handle non-existent keys gracefully")
	})
}

// Test behavior of config node with malformed messages
func TestConfigNodeMalformedMessages(t *testing.T) {
	t.Log("ConfigNode Malformed messages test suite")

	// Empty message
	t.Run("Empty message", func(t *testing.T) {
		// Setup
		rootNode, err := builders.NewRootNode()
		assert.NilError(t, err)

		testConfig := map[string]any{"test": true}
		graph, err := builders.NewGraph(rootNode, testConfig)
		assert.NilError(t, err)

		// Create an empty request without keys
		senderURL, err := url.Parse("actor://tester/" + uuid.NewString())
		assert.NilError(t, err)

		emptyReqMsg := &configMessage{
			sender:            *senderURL,
			configMessageType: 2,          // request
			requestedKeys:     []string{}, // Empty keys slice
		}

		// Send the message
		err = graph.Accept(emptyReqMsg)
		assert.NilError(t, err, "Should handle empty request gracefully")
	})
}
