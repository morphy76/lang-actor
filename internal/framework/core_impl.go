package framework

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/google/uuid"
	f "github.com/morphy76/lang-actor/pkg/framework"
)

var staticActorAssertion f.Actor[any] = (*actor[any])(nil)
var staticReceiverAssertion f.Addressable = (*actor[any])(nil)

var DefaultMailboxConfig = f.MailboxConfig{
	Capacity: 100,
	Policy:   f.BackpressurePolicyBlock,
}

type actor[T any] struct {
	lock *sync.Mutex

	status        f.ActorStatus
	stopCompleted chan bool

	ctx       context.Context
	ctxCancel context.CancelFunc

	address       url.URL
	mailbox       chan f.Message
	mailboxConfig f.MailboxConfig
	processingFn  f.ProcessingFn[T]

	parent   f.ActorRef
	children map[url.URL]f.ActorRef

	state T
}

// Address returns the actor's address.
func (a actor[T]) Address() url.URL {
	return a.address
}

// Stop stops the actor.
func (a *actor[T]) Stop() (chan bool, error) {
	for childURL := range a.children {
		a.Crop(childURL)
	}

	if a.status == f.ActorStatusRunning {
		return a.teardown()
	}

	return nil, fmt.Errorf("cannot stop actor: %w", f.ErrorActorNotRunning)
}

// Deliver delivers a message to the actor.
func (a *actor[T]) Deliver(msg f.Message) error {
	if a.status != f.ActorStatusRunning {
		return fmt.Errorf("failed to deliver message: %w", f.ErrorActorNotRunning)
	}

	switch a.mailboxConfig.Policy {
	case f.BackpressurePolicyBlock:
		a.mailbox <- msg
	case f.BackpressurePolicyFail:
		select {
		case a.mailbox <- msg:
		default:
			return fmt.Errorf("mailbox full: message rejected")
		}
	case f.BackpressurePolicyDropNewest:
		select {
		case a.mailbox <- msg:
		default:
			// Mailbox is full, silently drop the message
			// This is intentional as per the policy
		}
	case f.BackpressurePolicyDropOldest:
		if len(a.mailbox) == cap(a.mailbox) && cap(a.mailbox) > 0 {
			select {
			case <-a.mailbox:
			default:
			}
		}
		select {
		case a.mailbox <- msg:
		default:
			return fmt.Errorf("failed to deliver message: mailbox state changed unexpectedly")
		}

	case f.BackpressurePolicyUnbounded:
		a.mailbox <- msg

	default:
		a.mailbox <- msg
	}

	return nil
}

// State returns the actor's state.
func (a actor[T]) State() T {
	return a.state
}

// Status returns, by value, the actor's status.
func (a actor[T]) Status() f.ActorStatus {
	return a.status
}

// Send sends a message to the actor.
func (a actor[T]) Send(msg f.Message, addressable f.Addressable) error {
	return addressable.Deliver(msg)
}

// Append appends a child actor to the actor.
func (a *actor[T]) Append(child f.ActorRef) error {
	if err := a.verifyChildURL(child.Address()); err != nil {
		return fmt.Errorf("invalid child URL: %w", err)
	}
	a.lock.Lock()
	defer a.lock.Unlock()

	if _, ok := a.children[child.Address()]; ok {
		return fmt.Errorf("child already exists: %w", f.ErrorInvalidChildURL)
	}

	a.children[child.Address()] = child

	return nil
}

// Crop removes a child actor from the actor.
func (a *actor[T]) Crop(url url.URL) (f.ActorRef, error) {
	a.lock.Lock()
	defer a.lock.Unlock()

	if child, ok := a.children[url]; !ok {
		return nil, f.ErrorInvalidChildURL
	} else {
		stopCompleted, _ := child.Stop()
		<-stopCompleted
		delete(a.children, url)
		return child, nil
	}
}

// GetParent returns the parent actor of the current actor.
func (a *actor[T]) GetParent() (f.ActorRef, bool) {
	if a.parent == nil {
		return nil, false
	}
	return a.parent, true
}

func (a *actor[T]) verifyChildURL(url url.URL) error {
	if url.Scheme != a.address.Scheme || url.Host != a.address.Host {
		return f.ErrorInvalidChildURL
	}

	parentPath := a.address.Path
	childPath := url.Path

	if len(childPath) <= len(parentPath) || childPath[:len(parentPath)] != parentPath {
		return f.ErrorInvalidChildURL
	}

	remainingPath := childPath[len(parentPath):]
	if len(remainingPath) == 0 || remainingPath[0] != '/' {
		return f.ErrorInvalidChildURL
	}

	remainingPath = remainingPath[1:]
	if len(remainingPath) == 0 || remainingPath[0] == '/' || len(remainingPath) != len(url.Path[len(parentPath)+1:]) {
		return f.ErrorInvalidChildURL
	}

	return nil
}

func (a *actor[T]) teardown() (chan bool, error) {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.ctxCancel()

	return a.stopCompleted, nil
}

func (a *actor[T]) consume() {
	for {
		select {
		case msg := <-a.mailbox:
			newState, err := a.processingFn(msg, a)
			if err != nil {
				a.handleFailure(err)
			}

			if msg.Mutation() {
				a.swapState(newState)
			}
		case <-a.ctx.Done():
			cleanupTimeout := time.After(5 * time.Second)
		drainLoop:
			for {
				select {
				case msg := <-a.mailbox:
					newState, err := a.processingFn(msg, a)
					if err != nil {
						a.handleFailure(err)
					}

					if msg.Mutation() {
						a.swapState(newState)
					}
				case <-cleanupTimeout:
					a.status = f.ActorStatusIdle
					a.stopCompleted <- true
					return
				default:
					a.status = f.ActorStatusIdle
					a.stopCompleted <- true
					break drainLoop
				}
			}
			return
		}
	}
}

func (a *actor[T]) swapState(newState T) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.state = newState
}

func (a actor[T]) handleFailure(err error) {
	// TODO: Handle failure
	fmt.Println("handleFailure", err)
}

// NewActor creates a new actor with the given address.
func NewActor[T any](
	address url.URL,
	processingFn f.ProcessingFn[T],
	initialState T,
	mailboxConfig ...f.MailboxConfig,
) (f.Actor[T], error) {
	// TODO, future schema support:
	// - actor+http:// to dispatch messages over HTTP
	// - actor+https:// to dispatch messages over HTTPS
	// - actor+unix:// to dispatch messages over Unix domain sockets
	// - actor+tcp:// to dispatch messages over TCP
	// - actor+udp:// to dispatch messages over UDP
	// Validate the schema
	if address.Scheme != "actor" {
		return nil, f.ErrorInvalidActorAddress
	}

	useCtx, useCancelFn := context.WithCancel(context.Background())

	config := DefaultMailboxConfig
	if len(mailboxConfig) > 0 {
		config = mailboxConfig[0]
	}

	var mailbox chan f.Message
	switch config.Policy {
	case f.BackpressurePolicyUnbounded:
		// In Go, we can't truly have an unbounded channel, but we can make it very large
		mailbox = make(chan f.Message, 1000000)
	default:
		capacity := config.Capacity
		if capacity <= 0 {
			capacity = DefaultMailboxConfig.Capacity
		}
		mailbox = make(chan f.Message, capacity)
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

		state: initialState,
	}
	go rv.consume()

	return rv, nil
}

// NewActorWithParent creates a new actor with the given address and parent actor.
func NewActorWithParent[T any](
	processingFn f.ProcessingFn[T],
	initialState T,
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

	config := DefaultMailboxConfig
	if len(mailboxConfig) > 0 {
		config = mailboxConfig[0]
	}

	var mailbox chan f.Message
	switch config.Policy {
	case f.BackpressurePolicyUnbounded:
		// In Go, we can't truly have an unbounded channel, but we can make it very large
		mailbox = make(chan f.Message, 1000000)
	default:
		capacity := config.Capacity
		if capacity <= 0 {
			capacity = DefaultMailboxConfig.Capacity
		}
		mailbox = make(chan f.Message, capacity)
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

		state: initialState,
	}
	go rv.consume()

	if err := parent.Append(rv); err != nil {
		return nil, fmt.Errorf("failed to append child to parent: %w", err)
	}

	return rv, nil
}
