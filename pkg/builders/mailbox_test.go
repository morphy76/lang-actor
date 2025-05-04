package builders_test

import (
	"testing"

	"github.com/morphy76/lang-actor/pkg/builders"
	"github.com/morphy76/lang-actor/pkg/framework"
	"gotest.tools/v3/assert"
)

func TestNewMailboxConfigWithBlockPolicy(t *testing.T) {
	t.Log("NewMailboxConfigWithBlockPolicy test suite")

	t.Run("Creates mailbox with block policy and specified capacity", func(t *testing.T) {
		t.Log("Should create a mailbox configuration with block policy and specified capacity")

		capacity := 200
		config := builders.NewMailboxConfigWithBlockPolicy(capacity)

		assert.Equal(t, config.Capacity, capacity)
		assert.Equal(t, config.Policy, framework.BackpressurePolicyBlock)
	})
}

func TestNewMailboxConfigWithFailPolicy(t *testing.T) {
	t.Log("NewMailboxConfigWithFailPolicy test suite")

	t.Run("Creates mailbox with fail policy and specified capacity", func(t *testing.T) {
		t.Log("Should create a mailbox configuration with fail policy and specified capacity")

		capacity := 150
		config := builders.NewMailboxConfigWithFailPolicy(capacity)

		assert.Equal(t, config.Capacity, capacity)
		assert.Equal(t, config.Policy, framework.BackpressurePolicyFail)
	})
}

func TestNewMailboxConfigWithUnboundedPolicy(t *testing.T) {
	t.Log("NewMailboxConfigWithUnboundedPolicy test suite")

	t.Run("Creates mailbox with fail policy and large capacity", func(t *testing.T) {
		t.Log("Should create a mailbox configuration with fail policy and large capacity for unbounded behavior")

		config := builders.NewMailboxConfigWithUnboundedPolicy()

		assert.Equal(t, config.Capacity, 1000000)
		assert.Equal(t, config.Policy, framework.BackpressurePolicyFail)
	})
}

func TestNewMailboxConfigWithDropNewestPolicy(t *testing.T) {
	t.Log("NewMailboxConfigWithDropNewestPolicy test suite")

	t.Run("Creates mailbox with drop newest policy and specified capacity", func(t *testing.T) {
		t.Log("Should create a mailbox configuration with drop newest policy and specified capacity")

		capacity := 50
		config := builders.NewMailboxConfigWithDropNewestPolicy(capacity)

		assert.Equal(t, config.Capacity, capacity)
		assert.Equal(t, config.Policy, framework.BackpressurePolicyDropNewest)
	})
}

func TestNewMailboxConfigWithDropOldestPolicy(t *testing.T) {
	t.Log("NewMailboxConfigWithDropOldestPolicy test suite")

	t.Run("Creates mailbox with drop oldest policy and specified capacity", func(t *testing.T) {
		t.Log("Should create a mailbox configuration with drop oldest policy and specified capacity")

		capacity := 75
		config := builders.NewMailboxConfigWithDropOldestPolicy(capacity)

		assert.Equal(t, config.Capacity, capacity)
		assert.Equal(t, config.Policy, framework.BackpressurePolicyDropOldest)
	})
}
