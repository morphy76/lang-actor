package framework

import (
	"context"
	"net/url"
	"sync"

	f "github.com/morphy76/lang-actor/pkg/framework"
)

var staticActorAssertion f.Actor[any, any] = (*actor[any, any])(nil)

type actor[M any, T any] struct {
	lock *sync.Mutex

	ctx       context.Context
	ctxCancel context.CancelFunc

	status  f.ActorStatus
	mailbox chan f.Message[M]
	address url.URL

	state        f.Payload[T]
	processingFn f.ProcessingFn[M, T]

	stopCompleted chan bool
}

// Address returns the actor's address.
func (a actor[M, T]) Address() url.URL {
	return a.address
}

// Start starts the actor.
func (a *actor[M, T]) Start() error {

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
func (a *actor[M, T]) Stop() (chan bool, error) {

	if a.status == f.ActorStatusRunning {
		return a.teardown()
	}

	return nil, f.ErrorActorNotRunning
}

// Deliver delivers a message to the actor.
func (a *actor[M, T]) Deliver(msg f.Message[M]) error {
	if a.status != f.ActorStatusRunning {
		return f.ErrorActorNotRunning
	}
	a.mailbox <- msg
	return nil
}

// State returns the actor's state.
func (a actor[M, T]) State() T {
	return a.state.Cast()
}

// Status returns the actor's status.
func (a actor[M, T]) Status() f.ActorStatus {
	return a.status
}

func (a *actor[M, T]) warmup() error {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.status = f.ActorStatusStarting
	a.mailbox = make(chan f.Message[M], 100)
	useCtx, useCancelFn := context.WithCancel(context.Background())
	a.ctx = useCtx
	a.ctxCancel = useCancelFn
	go a.consume()

	a.status = f.ActorStatusRunning

	return nil
}

func (a *actor[M, T]) teardown() (chan bool, error) {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.status = f.ActorStatusStopping
	if a.ctxCancel != nil {
		a.ctxCancel()
	}

	a.stopCompleted = make(chan bool)

	return a.stopCompleted, nil
}

func (a *actor[M, T]) consume() {
	for {
		select {
		case msg := <-a.mailbox:
			newState, err := a.processingFn(msg, a.state)
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
					newState, err := a.processingFn(msg, a.state)
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

func (a *actor[M, T]) swapState(newState f.Payload[T]) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.state = newState
}

func (a actor[M, T]) handleFailure(err error) {
	// TODO: Handle failure
}

// NewActor creates a new actor with the given address.
func NewActor[M any, T any](
	address url.URL,
	processingFn f.ProcessingFn[M, T],
	initialState f.Payload[T],
) (f.Actor[M, T], error) {
	// TODO, future schema support:
	// - actor+http:// to dispatch messages over HTTP
	// - actor+https:// to dispatch messages over HTTPS
	// - actor+unix:// to dispatch messages over Unix domain sockets
	// - actor+tcp:// to dispatch messages over TCP
	// - actor+udp:// to dispatch messages over UDP
	// Validate the schema
	if address.Scheme != "actor" {
		return nil, &url.Error{
			Op:  "parse",
			URL: address.String(),
			Err: &url.Error{
				Op:  "unsupported schema",
				URL: address.String(),
				Err: nil,
			},
		}
	}

	return &actor[M, T]{
		lock: &sync.Mutex{},

		status:  f.ActorStatusIdle,
		address: address,

		state:        initialState,
		processingFn: processingFn,
	}, nil
}
