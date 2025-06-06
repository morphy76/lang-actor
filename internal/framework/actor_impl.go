package framework

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"time"

	c "github.com/morphy76/lang-actor/pkg/common"
	f "github.com/morphy76/lang-actor/pkg/framework"
)

var staticActorAssertion f.Actor[any] = (*actor[any])(nil)

type actor[T any] struct {
	lock *sync.Mutex

	status        f.ActorStatus
	stopCompleted chan bool

	ctx       context.Context
	ctxCancel context.CancelFunc

	address       url.URL
	mailbox       chan c.Message
	mailboxConfig f.MailboxConfig
	processingFn  f.ProcessingFn[T]

	parent   f.ActorRef
	children map[url.URL]f.ActorRef

	state     T
	transient bool
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
func (a *actor[T]) Deliver(msg c.Message) error {
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
func (a actor[T]) Send(msg c.Message, destination c.Transport) error {
	return destination.Deliver(msg)
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

// Children returns the children of the actor.
func (a *actor[T]) Children() []f.ActorRef {
	children := make([]f.ActorRef, 0, len(a.children))
	for _, child := range a.children {
		children = append(children, child)
	}
	return children
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
			useMessage, ok := msg.(f.Message)
			if ok {
				newState, err := a.processingFn(useMessage, a)
				if err != nil {
					a.handleFailure(err)
				}

				if useMessage.Mutation() || !a.transient {
					a.swapState(newState)
				}
			}
		case <-a.ctx.Done():
			cleanupTimeout := time.After(5 * time.Second)
		drainLoop:
			for {
				select {
				case msg := <-a.mailbox:
					useMessage, ok := msg.(f.Message)
					if ok {
						newState, err := a.processingFn(useMessage, a)
						if err != nil {
							a.handleFailure(err)
						}

						if useMessage.Mutation() {
							a.swapState(newState)
						}
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
	fmt.Printf("(%v) handleFailure: %s\n", a.address, err)
}
