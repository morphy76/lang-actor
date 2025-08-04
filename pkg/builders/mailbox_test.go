package builders_test

import (
	"testing"

	"github.com/morphy76/lang-actor/pkg/builders"
	"github.com/morphy76/lang-actor/pkg/framework"
	"gotest.tools/v3/assert"
)

func TestNewMailboxConfig(t *testing.T) {
	t.Log("TestNewMailboxConfig test suite")

	t.Run("Creates mailbox with block policy and specified capacity", func(t *testing.T) {
		t.Log("Should create a mailbox configuration with block policy and specified capacity")

		capacity := 200
		config := builders.NewMailboxConfigWithBlockPolicy(capacity)

		assert.Equal(t, config.Capacity, capacity)
		assert.Equal(t, config.Policy, framework.BackpressurePolicyBlock)
	})

	t.Run("Creates mailbox with fail policy and specified capacity", func(t *testing.T) {
		t.Log("Should create a mailbox configuration with fail policy and specified capacity")

		capacity := 150
		config := builders.NewMailboxConfigWithFailPolicy(capacity)

		assert.Equal(t, config.Capacity, capacity)
		assert.Equal(t, config.Policy, framework.BackpressurePolicyFail)
	})

	t.Run("Creates mailbox with fail policy and large capacity", func(t *testing.T) {
		t.Log("Should create a mailbox configuration with unbounded policy")

		config := builders.NewMailboxConfigWithUnboundedPolicy()

		assert.Equal(t, config.Capacity, 0) // Capacity is ignored for unbounded policy
		assert.Equal(t, config.Policy, framework.BackpressurePolicyUnbounded)
	})

	t.Run("Creates mailbox with drop newest policy and specified capacity", func(t *testing.T) {
		t.Log("Should create a mailbox configuration with drop newest policy and specified capacity")

		capacity := 50
		config := builders.NewMailboxConfigWithDropNewestPolicy(capacity)

		assert.Equal(t, config.Capacity, capacity)
		assert.Equal(t, config.Policy, framework.BackpressurePolicyDropNewest)
	})

	t.Run("Creates mailbox with drop oldest policy and specified capacity", func(t *testing.T) {
		t.Log("Should create a mailbox configuration with drop oldest policy and specified capacity")

		capacity := 75
		config := builders.NewMailboxConfigWithDropOldestPolicy(capacity)

		assert.Equal(t, config.Capacity, capacity)
		assert.Equal(t, config.Policy, framework.BackpressurePolicyDropOldest)
	})
}
