package graph

import (
	g "github.com/morphy76/lang-actor/pkg/graph"
)

var staticStateWrapperAssertion g.State = (*stateWrapper)(nil)

type stateWrapper struct {
	state          g.State
	stateChangesCh chan g.State
}

// MergeChange appends a new state to the graph and notifies the state changes channel.
func (s *stateWrapper) MergeChange(purpose any, value any) error {
	if err := s.state.MergeChange(purpose, value); err != nil {
		return err
	}
	s.stateChangesCh <- s.state
	return nil
}
