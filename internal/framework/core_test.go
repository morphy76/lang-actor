package framework_test

import (
	"net/url"
	"testing"

	"github.com/morphy76/lang-actor/internal/framework"
	f "github.com/morphy76/lang-actor/pkg/framework"
	"gotest.tools/v3/assert"
)

const actorURI = "actor://example"

type noState struct {
}

func TestNewActor(t *testing.T) {
	t.Log("NewActor test suite")

	var initialState = noState{}
	var nullProcessingFn f.ProcessingFn[noState] = func(msg f.Message, actor f.Actor[noState]) (noState, error) {
		return noState{}, nil
	}

	t.Run("Valid schema", func(t *testing.T) {
		t.Log("Should create an actor with a valid schema")

		address, err := url.Parse(actorURI)
		assert.NilError(t, err)

		actor, err := framework.NewActor(*address, nullProcessingFn, initialState)
		assert.NilError(t, err)
		assert.Equal(t, actor.Address().Scheme, address.Scheme)
	})

	t.Run("Invalid schema", func(t *testing.T) {
		t.Log("Should return an error for an unsupported schema")

		address, err := url.Parse("http://example")
		assert.NilError(t, err)

		actor, err := framework.NewActor(*address, nullProcessingFn, initialState)
		assert.Assert(t, actor == nil)
		assert.ErrorContains(t, err, "invalid actor address")
	})
}
func TestActorLifecycle(t *testing.T) {
	t.Log("Actor lifecycle test suite")

	var initialState = noState{}
	var nullProcessingFn f.ProcessingFn[noState] = func(msg f.Message, actor f.Actor[noState]) (noState, error) {
		return noState{}, nil
	}

	t.Run("Start actor successfully", func(t *testing.T) {
		t.Log("Should start the actor successfully when idle")

		address, err := url.Parse(actorURI)
		assert.NilError(t, err)

		actor, err := framework.NewActor(*address, nullProcessingFn, initialState)
		assert.NilError(t, err)

		assert.Equal(t, actor.Status(), f.ActorStatusRunning)
	})

	t.Run("Stop actor successfully", func(t *testing.T) {
		t.Log("Should stop the actor successfully when running")

		address, err := url.Parse(actorURI)
		assert.NilError(t, err)

		actor, err := framework.NewActor(*address, nullProcessingFn, initialState)
		assert.NilError(t, err)

		stopCompleted, err := actor.Stop()
		assert.NilError(t, err)

		<-stopCompleted

		assert.Equal(t, actor.Status(), f.ActorStatusIdle)
	})
}

var staticMockMessageAssertion f.Message = (*mockMessage)(nil)

type mockMessage struct {
	sender   url.URL
	mutation bool
}

func (m mockMessage) Sender() url.URL {
	return m.sender
}

func (m mockMessage) Mutation() bool {
	return m.mutation
}

type mockActorState struct {
	processed bool
}

func TestActorMessageDelivery(t *testing.T) {
	t.Log("Actor message delivery test suite")

	const actorURI = "actor://example"

	t.Run("Deliver message successfully", func(t *testing.T) {
		t.Log("Should deliver a message to the actor and invoke the processing function")

		address, err := url.Parse(actorURI)
		assert.NilError(t, err)

		var messageProcessed *bool = new(bool)
		*messageProcessed = false
		var spyFn f.ProcessingFn[noState] = func(msg f.Message, actor f.Actor[noState]) (noState, error) {
			*messageProcessed = true
			return noState{}, nil
		}

		actor, err := framework.NewActor(*address, spyFn, noState{})
		assert.NilError(t, err)

		message := &mockMessage{sender: *address, mutation: false}
		err = actor.Deliver(message)
		assert.NilError(t, err)

		stopCompleted, err := actor.Stop()
		assert.NilError(t, err)

		<-stopCompleted

		assert.Assert(t, *messageProcessed)
	})

	t.Run("Update actor state on mutation message", func(t *testing.T) {
		t.Log("Should update the actor's state when a mutation message is delivered")

		var initialState = mockActorState{processed: false}
		var spyFn f.ProcessingFn[mockActorState] = func(msg f.Message, actor f.Actor[mockActorState]) (mockActorState, error) {
			return mockActorState{processed: true}, nil
		}

		address, err := url.Parse(actorURI)
		assert.NilError(t, err)

		actor, err := framework.NewActor(*address, spyFn, initialState)
		assert.NilError(t, err)

		message := &mockMessage{sender: *address, mutation: true}
		err = actor.Deliver(message)
		assert.NilError(t, err)

		stopCompleted, err := actor.Stop()
		assert.NilError(t, err)

		<-stopCompleted

		assert.Assert(t, actor.State().processed)
	})
}
