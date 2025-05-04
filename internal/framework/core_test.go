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
		assert.Equal(t, actor.Address().Scheme, "BBB")
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

func TestBackpressurePolicies(t *testing.T) {
	t.Log("Backpressure policies test suite")

	var initialState = noState{}
	var nullProcessingFn f.ProcessingFn[noState] = func(msg f.Message, actor f.Actor[noState]) (noState, error) {
		return noState{}, nil
	}

	t.Run("Default block policy", func(t *testing.T) {
		t.Log("Should create an actor with the default block policy")

		address, err := url.Parse(actorURI)
		assert.NilError(t, err)

		actor, err := framework.NewActor(*address, nullProcessingFn, initialState)
		assert.NilError(t, err)

		// Verify we can deliver a message
		message := &mockMessage{sender: *address, mutation: false}
		err = actor.Deliver(message)
		assert.NilError(t, err)

		stopCompleted, err := actor.Stop()
		assert.NilError(t, err)
		<-stopCompleted
	})

	t.Run("Fail policy", func(t *testing.T) {
		t.Log("Should create an actor with fail policy that rejects messages when mailbox is full")

		address, err := url.Parse(actorURI)
		assert.NilError(t, err)

		// Create a mailbox config with tiny capacity and fail policy
		config := f.MailboxConfig{
			Capacity: 1,
			Policy:   f.BackpressurePolicyFail,
		}

		// Create a slow processing function that blocks
		blockCh := make(chan struct{})
		slowProcessingFn := func(msg f.Message, actor f.Actor[noState]) (noState, error) {
			<-blockCh // Block until test releases
			return noState{}, nil
		}

		actor, err := framework.NewActor(*address, slowProcessingFn, initialState, config)
		assert.NilError(t, err)

		// First message should be accepted
		message1 := &mockMessage{sender: *address, mutation: false}
		err = actor.Deliver(message1)
		assert.NilError(t, err)

		// Second message should be rejected since mailbox is full and we're using fail policy
		message2 := &mockMessage{sender: *address, mutation: false}
		err = actor.Deliver(message2)
		assert.ErrorContains(t, err, "mailbox full")

		// Release the processing function
		close(blockCh)

		stopCompleted, err := actor.Stop()
		assert.NilError(t, err)
		<-stopCompleted
	})

	t.Run("Drop newest policy", func(t *testing.T) {
		t.Log("Should create an actor with drop newest policy that silently drops new messages when mailbox is full")

		address, err := url.Parse(actorURI)
		assert.NilError(t, err)

		// Create a mailbox config with tiny capacity and drop newest policy
		config := f.MailboxConfig{
			Capacity: 1,
			Policy:   f.BackpressurePolicyDropNewest,
		}

		// Create a block channel to control message processing
		blockCh := make(chan struct{})

		// Track which messages were processed
		processedMessages := make(map[string]bool)
		processingFn := func(msg f.Message, actor f.Actor[noState]) (noState, error) {
			idMsg, ok := msg.(*mockMessageWithID)
			if ok {
				processedMessages[idMsg.id] = true
			}
			<-blockCh // Block until test releases
			return noState{}, nil
		}

		actor, err := framework.NewActor(*address, processingFn, initialState, config)
		assert.NilError(t, err)

		// First message should be accepted
		message1 := &mockMessageWithID{mockMessage: mockMessage{sender: *address}, id: "msg1"}
		err = actor.Deliver(message1)
		assert.NilError(t, err)

		// Second message should be silently dropped
		message2 := &mockMessageWithID{mockMessage: mockMessage{sender: *address}, id: "msg2"}
		err = actor.Deliver(message2)
		assert.NilError(t, err) // No error, but message should be dropped

		// Release processing
		close(blockCh)

		stopCompleted, err := actor.Stop()
		assert.NilError(t, err)
		<-stopCompleted

		// Verify only first message was processed
		assert.Equal(t, true, processedMessages["msg1"])
		assert.Equal(t, false, processedMessages["msg2"])
	})

	t.Run("Drop oldest policy", func(t *testing.T) {
		t.Log("Should create an actor with drop oldest policy that discards oldest messages when mailbox is full")

		address, err := url.Parse(actorURI)
		assert.NilError(t, err)

		// Create a mailbox config with tiny capacity and drop oldest policy
		config := f.MailboxConfig{
			Capacity: 1,
			Policy:   f.BackpressurePolicyDropOldest,
		}

		// Track which messages were received
		receiveChannel := make(chan string, 3) // Buffer to prevent blocking
		processingFn := func(msg f.Message, actor f.Actor[noState]) (noState, error) {
			idMsg, ok := msg.(*mockMessageWithID)
			if ok {
				receiveChannel <- idMsg.id
			}
			return noState{}, nil
		}

		actor, err := framework.NewActor(*address, processingFn, initialState, config)
		assert.NilError(t, err)

		// Deliver first message - should be accepted
		message1 := &mockMessageWithID{mockMessage: mockMessage{sender: *address}, id: "msg1"}
		err = actor.Deliver(message1)
		assert.NilError(t, err)

		// Deliver second message - should replace the first due to drop oldest policy
		message2 := &mockMessageWithID{mockMessage: mockMessage{sender: *address}, id: "msg2"}
		err = actor.Deliver(message2)
		assert.NilError(t, err)

		// Wait for processing to complete
		stopCompleted, err := actor.Stop()
		assert.NilError(t, err)
		<-stopCompleted

		// Check which messages were processed
		close(receiveChannel)
		processedMsgs := make([]string, 0)
		for id := range receiveChannel {
			processedMsgs = append(processedMsgs, id)
		}

		// Only the newest message (msg2) should be processed, as it replaced msg1
		// in the mailbox when it was added
		assert.Equal(t, 1, len(processedMsgs))
		if len(processedMsgs) > 0 {
			assert.Equal(t, "msg2", processedMsgs[0])
		}
	})

	t.Run("Unbounded policy", func(t *testing.T) {
		t.Log("Should create an actor with unbounded policy that accepts many messages")

		address, err := url.Parse(actorURI)
		assert.NilError(t, err)

		// Create a mailbox config with unbounded policy
		config := f.MailboxConfig{
			Policy: f.BackpressurePolicyUnbounded,
		}

		// Keep track of processed messages
		processedCount := 0
		counterFn := func(msg f.Message, actor f.Actor[noState]) (noState, error) {
			processedCount++
			return noState{}, nil
		}

		actor, err := framework.NewActor(*address, counterFn, initialState, config)
		assert.NilError(t, err)

		// Deliver a significant number of messages (less than our "unbounded" buffer)
		messagesToSend := 500
		for i := 0; i < messagesToSend; i++ {
			message := &mockMessage{sender: *address, mutation: false}
			err = actor.Deliver(message)
			assert.NilError(t, err)
		}

		// Stop the actor which will process all remaining messages
		stopCompleted, err := actor.Stop()
		assert.NilError(t, err)
		<-stopCompleted

		// Verify all messages were processed
		assert.Equal(t, messagesToSend, processedCount)
	})

	t.Run("Custom capacity", func(t *testing.T) {
		t.Log("Should create an actor with custom mailbox capacity")

		address, err := url.Parse(actorURI)
		assert.NilError(t, err)

		// Create a mailbox config with custom capacity of 5 and fail policy
		config := f.MailboxConfig{
			Capacity: 5,
			Policy:   f.BackpressurePolicyFail,
		}

		// Create a block channel to control message processing
		blockCh := make(chan struct{})

		processingFn := func(msg f.Message, actor f.Actor[noState]) (noState, error) {
			<-blockCh // Block until test releases
			return noState{}, nil
		}

		actor, err := framework.NewActor(*address, processingFn, initialState, config)
		assert.NilError(t, err)

		// Should accept 5 messages (matching our capacity)
		for i := 0; i < 5; i++ {
			message := &mockMessage{sender: *address, mutation: false}
			err = actor.Deliver(message)
			assert.NilError(t, err)
		}

		// Sixth message should fail
		message := &mockMessage{sender: *address, mutation: false}
		err = actor.Deliver(message)
		assert.ErrorContains(t, err, "mailbox full")

		// Release processing
		close(blockCh)

		stopCompleted, err := actor.Stop()
		assert.NilError(t, err)
		<-stopCompleted
	})
}

// mockMessageWithID extends mockMessage to include an ID for tracking
type mockMessageWithID struct {
	mockMessage
	id string
}
