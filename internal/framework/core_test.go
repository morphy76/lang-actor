package framework_test

import (
	"net/url"
	"testing"

	"github.com/morphy76/lang-actor/internal/framework"
	"gotest.tools/v3/assert"
)

const actorURI = "actor://example"

var staticNoStateAssertion framework.Payload[noState] = (*noState)(nil)

type noState struct {
}

func (n noState) ToImplementation() noState {
	return n
}

func TestNewActor(t *testing.T) {
	t.Log("NewActor test suite")

	var initialState = noState{}
	var nullProcessingFn framework.ProcessingFn[noState] = func(msg framework.Message, currentState framework.Payload[noState]) (framework.Payload[noState], error) {
		return &noState{}, nil
	}

	t.Run("Valid schema", func(t *testing.T) {
		t.Log("Should create an actor with a valid schema")

		address, err := url.Parse(actorURI)
		assert.NilError(t, err)

		actor, err := framework.NewActor(*address, nullProcessingFn, initialState)
		assert.NilError(t, err)
		assert.Equal(t, actor.String(), actorURI)
	})

	t.Run("Invalid schema", func(t *testing.T) {
		t.Log("Should return an error for an unsupported schema")

		address, err := url.Parse("http://example")
		assert.NilError(t, err)

		actor, err := framework.NewActor(*address, nullProcessingFn, initialState)
		assert.Assert(t, actor == nil)
		assert.ErrorContains(t, err, "unsupported schema")
	})
}
func TestActorLifecycle(t *testing.T) {
	t.Log("Actor lifecycle test suite")

	var initialState = noState{}
	var nullProcessingFn framework.ProcessingFn[noState] = func(msg framework.Message, currentState framework.Payload[noState]) (framework.Payload[noState], error) {
		return &noState{}, nil
	}

	t.Run("Start actor successfully", func(t *testing.T) {
		t.Log("Should start the actor successfully when idle")

		address, err := url.Parse(actorURI)
		assert.NilError(t, err)

		actor, err := framework.NewActor(*address, nullProcessingFn, initialState)
		assert.NilError(t, err)

		err = actor.Start()
		assert.NilError(t, err)
		assert.Equal(t, actor.Status(), framework.ActorStatusRunning)
	})

	t.Run("Start actor when already running", func(t *testing.T) {
		t.Log("Should return an error when starting an already running actor")

		address, err := url.Parse(actorURI)
		assert.NilError(t, err)

		actor, err := framework.NewActor(*address, nullProcessingFn, initialState)
		assert.NilError(t, err)

		err = actor.Start()
		assert.NilError(t, err)

		err = actor.Start()
		assert.ErrorContains(t, err, "actor already started")
	})

	t.Run("Stop actor successfully", func(t *testing.T) {
		t.Log("Should stop the actor successfully when running")

		address, err := url.Parse(actorURI)
		assert.NilError(t, err)

		actor, err := framework.NewActor(*address, nullProcessingFn, initialState)
		assert.NilError(t, err)

		err = actor.Start()
		assert.NilError(t, err)

		stopCompleted, err := actor.Stop()
		assert.NilError(t, err)

		<-stopCompleted

		assert.Equal(t, actor.Status(), framework.ActorStatusIdle)
	})

	t.Run("Stop actor when not running", func(t *testing.T) {
		t.Log("Should return an error when stopping an actor that is not running")

		address, err := url.Parse(actorURI)
		assert.NilError(t, err)

		actor, err := framework.NewActor(*address, nullProcessingFn, initialState)
		assert.NilError(t, err)

		_, err = actor.Stop()
		assert.ErrorContains(t, err, "actor not running")
	})
}

var staticMockMessageAssertion framework.Message = (*mockMessage)(nil)

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

var staticMockActorStateAssertion framework.Payload[mockActorState] = (*mockActorState)(nil)

type mockActorState struct {
	processed bool
}

func (m mockActorState) ToImplementation() mockActorState {
	return m
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
		var spyFn framework.ProcessingFn[noState] = func(msg framework.Message, currentState framework.Payload[noState]) (framework.Payload[noState], error) {
			*messageProcessed = true
			return noState{}, nil
		}

		actor, err := framework.NewActor(*address, spyFn, noState{})
		assert.NilError(t, err)

		err = actor.Start()
		assert.NilError(t, err)

		message := &mockMessage{sender: *address, mutation: false}
		err = actor.Deliver(message)
		assert.NilError(t, err)

		stopCompleted, err := actor.Stop()
		assert.NilError(t, err)

		<-stopCompleted

		assert.Assert(t, *messageProcessed)
	})

	t.Run("Deliver message when actor is not running", func(t *testing.T) {
		t.Log("Should return an error when delivering a message to an actor that is not running")

		var messageProcessed bool
		var spyFn framework.ProcessingFn[noState] = func(msg framework.Message, currentState framework.Payload[noState]) (framework.Payload[noState], error) {
			messageProcessed = true
			return noState{}, nil
		}

		address, err := url.Parse(actorURI)
		assert.NilError(t, err)

		actor, err := framework.NewActor(*address, spyFn, noState{})
		assert.NilError(t, err)

		message := &mockMessage{sender: *address, mutation: false}
		err = actor.Deliver(message)
		assert.ErrorContains(t, err, "actor not running")
		assert.Assert(t, !messageProcessed)
	})

	t.Run("Update actor state on mutation message", func(t *testing.T) {
		t.Log("Should update the actor's state when a mutation message is delivered")

		var initialState = mockActorState{processed: false}
		var spyFn framework.ProcessingFn[mockActorState] = func(msg framework.Message, currentState framework.Payload[mockActorState]) (framework.Payload[mockActorState], error) {
			return mockActorState{processed: true}, nil
		}

		address, err := url.Parse(actorURI)
		assert.NilError(t, err)

		actor, err := framework.NewActor(*address, spyFn, initialState)
		assert.NilError(t, err)

		err = actor.Start()
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
