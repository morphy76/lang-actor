package framework

import (
	"context"
	"fmt"
	"net/url"
	"sync"

	"github.com/google/uuid"
	c "github.com/morphy76/lang-actor/pkg/common"
	f "github.com/morphy76/lang-actor/pkg/framework"
)

// DefaultMailboxConfig is the default mailbox configuration.
var defaultMailboxConfig = f.MailboxConfig{
	Capacity: 100,
	Policy:   f.BackpressurePolicyBlock,
}

// NewActor creates a new actor with the given address.
func NewActor[T any](
	address url.URL,
	processingFn f.ProcessingFn[T],
	initialState T,
	transient bool,
	mailboxConfig ...f.MailboxConfig,
) (f.Actor[T], error) {
	// Validate the schema
	if address.Scheme != "actor" {
		return nil, f.ErrorInvalidActorAddress
	}

	useCtx, useCancelFn := context.WithCancel(context.Background())

	config := defaultMailboxConfig
	if len(mailboxConfig) > 0 {
		config = mailboxConfig[0]
	}

	var mailbox chan c.Message
	switch config.Policy {
	case f.BackpressurePolicyUnbounded:
		// In Go, we can't truly have an unbounded channel, but we can make it very large
		mailbox = make(chan c.Message, 1000000)
	default:
		capacity := config.Capacity
		if capacity <= 0 {
			capacity = defaultMailboxConfig.Capacity
		}
		mailbox = make(chan c.Message, capacity)
	}

	rv := &actor[T]{
		lock: &sync.Mutex{},

		status:        f.ActorStatusRunning,
		stopCompleted: make(chan bool, 1),

		ctx:       useCtx,
		ctxCancel: useCancelFn,

		address:       address,
		mailbox:       mailbox,
		mailboxConfig: config,
		processingFn:  processingFn,

		children: make(map[url.URL]f.ActorRef),

		state:     initialState,
		transient: transient,
	}
	go rv.consume()

	return rv, nil
}

// NewActorWithParent creates a new actor with the given address and parent actor.
func NewActorWithParent[T any](
	processingFn f.ProcessingFn[T],
	initialState T,
	transient bool,
	parent f.ActorRef,
	mailboxConfig ...f.MailboxConfig,
) (f.Actor[T], error) {
	address, err := url.Parse(fmt.Sprintf(
		"actor://%s%s",
		parent.Address().Host,
		parent.Address().Path+"/"+uuid.NewString(),
	))
	if err != nil {
		return nil, fmt.Errorf("failed to parse actor address: %w", err)
	}

	useCtx, useCancelFn := context.WithCancel(context.Background())

	config := defaultMailboxConfig
	if len(mailboxConfig) > 0 {
		config = mailboxConfig[0]
	}

	var mailbox chan c.Message
	switch config.Policy {
	case f.BackpressurePolicyUnbounded:
		// In Go, we can't truly have an unbounded channel, but we can make it very large
		mailbox = make(chan c.Message, 1000000)
	default:
		capacity := config.Capacity
		if capacity <= 0 {
			capacity = defaultMailboxConfig.Capacity
		}
		mailbox = make(chan c.Message, capacity)
	}

	rv := &actor[T]{
		lock: &sync.Mutex{},

		status:        f.ActorStatusRunning,
		stopCompleted: make(chan bool, 1),

		ctx:       useCtx,
		ctxCancel: useCancelFn,

		address:       *address,
		mailbox:       mailbox,
		mailboxConfig: config,
		processingFn:  processingFn,

		parent:   parent,
		children: make(map[url.URL]f.ActorRef),

		state:     initialState,
		transient: transient,
	}
	go rv.consume()

	if err := parent.Append(rv); err != nil {
		return nil, fmt.Errorf("failed to append child to parent: %w", err)
	}

	return rv, nil
}
