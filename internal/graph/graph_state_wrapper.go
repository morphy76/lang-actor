package graph

import (
	"errors"
	"sync"

	g "github.com/morphy76/lang-actor/pkg/graph"
)

var staticStateWrapperAssertion g.State = (*stateWrapper)(nil)

// NewStateWrapper creates a new state wrapper for the graph state.
// Parameters:
//   - state (g.State): The initial state of the graph.
//   - stateChangesCh (chan g.State): A channel to notify state changes.
//
// // Returns:
//   - (*stateWrapper): The created state wrapper instance.
//   - (error): An error if the state wrapper could not be created.
func NewStateWrapper(state g.State, stateChangesCh chan g.State) (*stateWrapper, error) {
	if state == nil {
		return nil, errors.New("state cannot be nil")
	}
	if stateChangesCh == nil {
		return nil, errors.New("stateChangesCh cannot be nil")
	}
	return &stateWrapper{
		lock: &sync.Mutex{},

		state:          state,
		stateChangesCh: stateChangesCh,
	}, nil
}

type stateWrapper struct {
	lock *sync.Mutex

	state          g.State
	stateChangesCh chan g.State
}

// MergeChange appends a new state to the graph and notifies the state changes channel.
func (s *stateWrapper) MergeChange(purpose any, value any) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if err := s.state.MergeChange(purpose, value); err != nil {
		return err
	}

	select {
	case s.stateChangesCh <- s.state:
	default:
		// Channel is full, skip notification rather than blocking
		// TODO In production, this could be logged as a warning
	}

	return nil
}

// Unwrap retrieves the underlying, non-proxy, state.
func (s *stateWrapper) Unwrap() g.State {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.state
}

// TODO
func (s *stateWrapper) ReadAttribute(name string) any {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.state.ReadAttribute(name)
}
