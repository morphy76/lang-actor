package bah_test

import (
	"testing"

	"github.com/morphy76/lang-actor/internal/bah"
	"gotest.tools/v3/assert"
)

func TestHeadersSuite(t *testing.T) {
	t.Log("Bah test suite")

	t.Run("Bah", func(t *testing.T) {
		t.Log("Simple bah test")

		given := 0
		expected := bah.Bah(given)

		assert.Equal(t, expected, 1, "Bah should return 1 when given 0")
	})
}
