package builders

import "github.com/morphy76/lang-actor/pkg/framework"

// NewMailboxConfigWithBlockPolicy creates a mailbox configuration with a block backpressure policy.
func NewMailboxConfigWithBlockPolicy(capacity int) framework.MailboxConfig {
	return framework.MailboxConfig{
		Capacity: capacity,
		Policy:   framework.BackpressurePolicyBlock,
	}
}

// NewMailboxConfigWithFailPolicy creates a mailbox configuration with a block backpressure policy.
func NewMailboxConfigWithFailPolicy(capacity int) framework.MailboxConfig {
	return framework.MailboxConfig{
		Capacity: capacity,
		Policy:   framework.BackpressurePolicyFail,
	}
}

// NewMailboxConfigWithUnboundedPolicy creates a mailbox configuration with an unbounded backpressure policy.
func NewMailboxConfigWithUnboundedPolicy() framework.MailboxConfig {
	return framework.MailboxConfig{
		Capacity: 1000000,
		Policy:   framework.BackpressurePolicyFail,
	}
}

// NewMailboxConfigWithDropNewestPolicy creates a mailbox configuration with a drop backpressure policy.
func NewMailboxConfigWithDropNewestPolicy(capacity int) framework.MailboxConfig {
	return framework.MailboxConfig{
		Capacity: capacity,
		Policy:   framework.BackpressurePolicyDropNewest,
	}
}

// NewMailboxConfigWithDropOldestPolicy creates a mailbox configuration with a drop backpressure policy.
func NewMailboxConfigWithDropOldestPolicy(capacity int) framework.MailboxConfig {
	return framework.MailboxConfig{
		Capacity: capacity,
		Policy:   framework.BackpressurePolicyDropOldest,
	}
}
