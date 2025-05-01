package framework

import (
	"context"
	"net/url"
	"sync"

	f "github.com/morphy76/lang-actor/pkg/framework"
)

var staticActorAssertion f.Actor[any] = (*actor[any])(nil)
var staticReceiverAssertion f.Transport = (*actor[any])(nil)

type actor[T any] struct {
	parentCtx context.Context
	lock      *sync.Mutex

	ctx       context.Context
	ctxCancel context.CancelFunc

	status  f.ActorStatus
	mailbox chan f.Message
	address url.URL

	state        f.ActorState[T]
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
	return a.state.Cast()
}

// Status returns the actor's status.
func (a actor[T]) Status() f.ActorStatus {
	return a.status
}

// SendFn returns a function to send messages to other actors.
func (a actor[T]) Send(msg f.Message, destination url.URL) error {
	// TODO a better implementation to send messages to other actors
	catalog := a.parentCtx.Value(f.ActorCatalogContextKey).(map[url.URL]f.Transport)
	return catalog[destination].Deliver(msg)
}

func (a *actor[T]) warmup() error {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.status = f.ActorStatusStarting
	a.mailbox = make(chan f.Message, 100)
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
			newState, err := a.processingFn(msg, a.state, a.Send, a.address)
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
					newState, err := a.processingFn(msg, a.state, a.Send, a.address)
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

func (a *actor[T]) swapState(newState f.ActorState[T]) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.state = newState
}

func (a actor[T]) handleFailure(err error) {
	// TODO: Handle failure
}

// NewActor creates a new actor with the given address.
func NewActor[T any](
	parentCtx context.Context,
	address url.URL,
	processingFn f.ProcessingFn[T],
	initialState f.ActorState[T],
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
		parentCtx: parentCtx,
		lock:      &sync.Mutex{},

		status:  f.ActorStatusIdle,
		address: address,

		state:        initialState,
		processingFn: processingFn,
	}, nil
}
