package graph_test

import (
	"testing"

	"github.com/morphy76/lang-actor/pkg/builders"
	"gotest.tools/v3/assert"
)

// TestConfigNodeMinimal tests just the creation of a graph with configuration
func TestConfigNodeMinimal(t *testing.T) {
	t.Log("Minimal ConfigNode test")

	// Create root node
	rootNode, err := builders.NewRootNode()
	assert.NilError(t, err, "Failed to create root node")

	// Create simple config
	config := map[string]any{
		"test": true,
	}

	// Create graph which should create a ConfigNode
	graph, err := builders.NewGraph(rootNode, config)
	assert.NilError(t, err, "Failed to create graph")

	// Just verify the graph exists
	assert.Assert(t, graph != nil, "Graph should not be nil")
}
