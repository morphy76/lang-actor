package builders_test

import (
	"testing"

	"github.com/morphy76/lang-actor/pkg/builders"
	"gotest.tools/v3/assert"
)

func TestGraphBuilders(t *testing.T) {
	t.Log("TestGraphBuilders test suite")

	t.Run("StartGraph", func(t *testing.T) {
		t.Log("Should start the actor graph")
		_, err := builders.StartGraph()
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		assert.NilError(t, err)
		t.Log("Actor graph started successfully")
	})
}
