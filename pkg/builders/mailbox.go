package builders

import "github.com/morphy76/lang-actor/pkg/framework"

// NewMailboxConfigWithBlockPolicy creates a mailbox configuration with a block backpressure policy.
//
// Parameters:
//   - capacity (int): The capacity of the mailbox.
//
// Returns:
//   - (framework.MailboxConfig): The created MailboxConfig instance.
func NewMailboxConfigWithBlockPolicy(capacity int) framework.MailboxConfig {
	return framework.MailboxConfig{
		Capacity: capacity,
		Policy:   framework.BackpressurePolicyBlock,
	}
}

// NewMailboxConfigWithFailPolicy creates a mailbox configuration with a block backpressure policy.
//
// Parameters:
//   - capacity (int): The capacity of the mailbox.
//
// Returns:
//   - (framework.MailboxConfig): The created MailboxConfig instance.
func NewMailboxConfigWithFailPolicy(capacity int) framework.MailboxConfig {
	return framework.MailboxConfig{
		Capacity: capacity,
		Policy:   framework.BackpressurePolicyFail,
	}
}

// NewMailboxConfigWithUnboundedPolicy creates a mailbox configuration with an unbounded backpressure policy.
//
// Returns:
//   - (framework.MailboxConfig): The created MailboxConfig instance.
func NewMailboxConfigWithUnboundedPolicy() framework.MailboxConfig {
	return framework.MailboxConfig{
		Capacity: 0, // Capacity is ignored for unbounded policy
		Policy:   framework.BackpressurePolicyUnbounded,
	}
}

// NewMailboxConfigWithDropNewestPolicy creates a mailbox configuration with a drop backpressure policy.
//
// Parameters:
//   - capacity (int): The capacity of the mailbox.
//
// Returns:
//   - (framework.MailboxConfig): The created MailboxConfig instance.
func NewMailboxConfigWithDropNewestPolicy(capacity int) framework.MailboxConfig {
	return framework.MailboxConfig{
		Capacity: capacity,
		Policy:   framework.BackpressurePolicyDropNewest,
	}
}

// NewMailboxConfigWithDropOldestPolicy creates a mailbox configuration with a drop backpressure policy.
//
// Parameters:
//   - capacity (int): The capacity of the mailbox.
//
// Returns:
//   - (framework.MailboxConfig): The created MailboxConfig instance.
func NewMailboxConfigWithDropOldestPolicy(capacity int) framework.MailboxConfig {
	return framework.MailboxConfig{
		Capacity: capacity,
		Policy:   framework.BackpressurePolicyDropOldest,
	}
}
