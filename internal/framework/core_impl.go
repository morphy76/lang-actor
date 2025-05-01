package framework

import (
	"context"
	"net/url"
	"sync"
)

var staticActorAssertion Actor[any] = (*actor[any])(nil)

type actor[T any] struct {
	lock *sync.Mutex

	ctx       context.Context
	ctxCancel context.CancelFunc

	status  ActorStatus
	mailbox chan Message
	address url.URL

	state        Payload[T]
	processingFn ProcessingFn[T]

	stopCompleted chan bool
}

// Address returns the actor's address.
func (a actor[T]) Address() url.URL {
	return a.address
}

// Start starts the actor.
func (a *actor[T]) Start() error {

	if a.status == ActorStatusRunning ||
		a.status == ActorStatusStarting {
		return ErrorActorAlreadyStarted
	}

	if a.status == ActorStatusIdle {
		return a.warmup()
	}

	return ErrorActorNotIdle
}

func (a *actor[T]) warmup() error {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.status = ActorStatusStarting
	a.mailbox = make(chan Message, 100)
	useCtx, useCancelFn := context.WithCancel(context.Background())
	a.ctx = useCtx
	a.ctxCancel = useCancelFn
	go a.consume()

	a.status = ActorStatusRunning

	return nil
}

// Stop stops the actor.
func (a *actor[T]) Stop() (chan bool, error) {

	if a.status == ActorStatusRunning {
		return a.teardown()
	}

	return nil, ErrorActorNotRunning
}

func (a *actor[T]) teardown() (chan bool, error) {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.status = ActorStatusStopping
	if a.ctxCancel != nil {
		a.ctxCancel()
	}

	a.stopCompleted = make(chan bool)

	return a.stopCompleted, nil
}

// Deliver delivers a message to the actor.
func (a *actor[T]) Deliver(msg Message) error {
	if a.status != ActorStatusRunning {
		return ErrorActorNotRunning
	}
	a.mailbox <- msg
	return nil
}

func (a *actor[T]) consume() {
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
					a.status = ActorStatusIdle
					a.stopCompleted <- true
					return
				}
			}
		}
	}
}

func (a *actor[T]) swapState(newState Payload[T]) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.state = newState
}

func (a actor[T]) handleFailure(err error) {
	// TODO: Handle failure
}

// Status returns the actor's status.
func (a actor[T]) Status() ActorStatus {
	return a.status
}

// State returns the actor's state.
func (a actor[T]) State() T {
	return a.state.ToImplementation()
}

// String returns the string representation of the actor's address.
func (a actor[T]) String() string {
	return a.address.String()
}

// NewActor creates a new actor with the given address.
// Supported schemas are:
// - actor://
//
// Parameters:
//   - address (url.URL): The URL address that specifies the actor's location and protocol.
//   - processingFn (ProcessingFn): The function to process messages sent to the actor.
//   - initialState (ActorState): The initial state of the actor.
//
// Returns:
//   - (Actor): The created Actor instance.
//   - (error): An error if the actor could not be created.
func NewActor[T any](
	address url.URL,
	processingFn ProcessingFn[T],
	initialState Payload[T],
) (Actor[T], error) {
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

	return &actor[T]{
		lock: &sync.Mutex{},

		status:  ActorStatusIdle,
		address: address,

		state:        initialState,
		processingFn: processingFn,
	}, nil
}
