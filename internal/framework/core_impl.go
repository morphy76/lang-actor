package framework

import (
	"context"
	"net/url"
	"sync"

	f "github.com/morphy76/lang-actor/pkg/framework"
)

var staticActorAssertion f.Actor[any] = (*actor[any])(nil)
var staticActorViewAssertion f.ActorView[any] = (*actorView[any])(nil)
var staticReceiverAssertion f.Addressable = (*actor[any])(nil)

type actor[T any] struct {
	lock *sync.Mutex

	ctx       context.Context
	ctxCancel context.CancelFunc

	status   f.ActorStatus
	mailbox  chan f.Message
	address  url.URL
	parent   f.ActorRef
	children map[url.URL]f.ActorRef

	state        T
	processingFn f.ProcessingFn[T]

	stopCompleted chan bool
}

// Address returns the actor's address.
func (a actor[T]) Address() url.URL {
	return a.address
}

// Start starts the actor.
func (a *actor[T]) Start() error {

	if a.status == f.ActorStatusRunning ||
		a.status == f.ActorStatusStarting {
		return f.ErrorActorAlreadyStarted
	}

	if a.status == f.ActorStatusIdle {
		return a.warmup()
	}

	return f.ErrorActorNotIdle
}

// Stop stops the actor.
func (a *actor[T]) Stop() (chan bool, error) {
	for childURL := range a.children {
		a.Crop(childURL)
	}

	if a.status == f.ActorStatusRunning {
		return a.teardown()
	}

	return nil, f.ErrorActorNotRunning
}

// Deliver delivers a message to the actor.
func (a *actor[T]) Deliver(msg f.Message) error {
	if a.status != f.ActorStatusRunning {
		return f.ErrorActorNotRunning
	}
	a.mailbox <- msg
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
		return err
	}
	a.lock.Lock()
	defer a.lock.Unlock()

	if _, ok := a.children[child.Address()]; ok {
		return f.ErrorInvalidChildURL
	}
	a.children[child.Address()] = child

	return child.Start()
}

// Crop removes a child actor from the actor.
func (a *actor[T]) Crop(url url.URL) error {
	if err := a.verifyChildURL(url); err != nil {
		return err
	}
	a.lock.Lock()
	defer a.lock.Unlock()

	if child, ok := a.children[url]; !ok {
		return f.ErrorInvalidChildURL
	} else {
		go child.(f.Controllable).Stop()
	}
	delete(a.children, url)
	return nil
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

func (a *actor[T]) warmup() error {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.status = f.ActorStatusStarting
	a.mailbox = make(chan f.Message, 100)
	a.children = make(map[url.URL]f.ActorRef)
	useCtx, useCancelFn := context.WithCancel(context.Background())
	a.ctx = useCtx
	a.ctxCancel = useCancelFn
	go a.consume()

	a.status = f.ActorStatusRunning

	return nil
}

func (a *actor[T]) teardown() (chan bool, error) {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.status = f.ActorStatusStopping
	if a.ctxCancel != nil {
		a.ctxCancel()
	}

	a.stopCompleted = make(chan bool)

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
				default:
					a.status = f.ActorStatusIdle
					a.stopCompleted <- true
					return
				}
			}
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
}

type actorView[T any] struct {
	actor f.Actor[T]
}

// Addess returns the actor's address.
func (a *actorView[T]) Address() url.URL {
	return a.actor.Address()
}

// State returns the actor's state.
func (a actorView[T]) State() T {
	return a.actor.State()
}

// Send sends a message to the actor.
func (a *actorView[T]) Send(msg f.Message, addressable f.Addressable) error {
	return a.actor.Send(msg, addressable)
}

// Deliver delivers a message to the actor.
func (a *actorView[T]) Deliver(msg f.Message) error {
	return a.actor.Deliver(msg)
}

// GetParent returns the parent actor of the current actor.
func (a *actorView[T]) GetParent() (f.ActorRef, bool) {
	return a.actor.GetParent()
}

// NewActor creates a new actor with the given address.
func NewActor[T any](
	address url.URL,
	processingFn f.ProcessingFn[T],
	initialState T,
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

	return &actor[T]{
		lock: &sync.Mutex{},

		status:  f.ActorStatusIdle,
		address: address,

		state:        initialState,
		processingFn: processingFn,
	}, nil
}

// NewActorWithParent creates a new actor with the given address and parent actor.
func NewActorWithParent[T any](
	address url.URL,
	processingFn f.ProcessingFn[T],
	initialState T,
	parent f.ActorRef,
) (f.Actor[T], error) {
	actor, err := NewActor(address, processingFn, initialState)
	if err != nil {
		return nil, err
	}

	if err := parent.Append(actor); err != nil {
		return nil, err
	}

	return actor, nil
}
