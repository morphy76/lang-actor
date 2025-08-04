package graph

import (
	"testing"
	"time"

	g "github.com/morphy76/lang-actor/pkg/graph"
)

// mockState is a mock implementation of the g.State interface for testing.
type mockState struct {
	mergeChangeCalled bool
	mergeChangeError  error
	purposeReceived   any
	valueReceived     any
}

func (m *mockState) MergeChange(purpose any, value any) error {
	m.mergeChangeCalled = true
	m.purposeReceived = purpose
	m.valueReceived = value
	return m.mergeChangeError
}

func TestNewStateWrapper(t *testing.T) {
	t.Log("StateWrapper test suite")

	t.Run("NewStateWrapper creates valid wrapper", func(t *testing.T) {
		t.Log("NewStateWrapper creates valid wrapper test case")

		// Arrange
		mockState := &mockState{}
		stateChangesCh := make(chan g.State, 1)

		// Act
		wrapper := &stateWrapper{
			state:          mockState,
			stateChangesCh: stateChangesCh,
		}

		// Assert
		if wrapper == nil {
			t.Errorf("Expected a StateWrapper, but got nil")
		}

		// Verify that the wrapper implements the State interface
		var _ g.State = wrapper
	})

	t.Run("MergeChange proxies to underlying state", func(t *testing.T) {
		t.Log("MergeChange proxies to underlying state test case")

		// Arrange
		mockState := &mockState{}
		stateChangesCh := make(chan g.State, 1)
		wrapper := &stateWrapper{
			state:          mockState,
			stateChangesCh: stateChangesCh,
		}

		purpose := "test-purpose"
		value := "test-value"

		// Act
		err := wrapper.MergeChange(purpose, value)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, but got: %v", err)
		}

		if !mockState.mergeChangeCalled {
			t.Errorf("Expected MergeChange to be called on underlying state")
		}

		if mockState.purposeReceived != purpose {
			t.Errorf("Expected purpose %v, but got %v", purpose, mockState.purposeReceived)
		}

		if mockState.valueReceived != value {
			t.Errorf("Expected value %v, but got %v", value, mockState.valueReceived)
		}
	})

	t.Run("MergeChange notifies state changes channel on success", func(t *testing.T) {
		t.Log("MergeChange notifies state changes channel on success test case")

		// Arrange
		mockState := &mockState{}
		stateChangesCh := make(chan g.State, 1)
		wrapper := &stateWrapper{
			state:          mockState,
			stateChangesCh: stateChangesCh,
		}

		purpose := "test-purpose"
		value := "test-value"

		// Act
		err := wrapper.MergeChange(purpose, value)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, but got: %v", err)
		}

		// Check that a notification was sent to the channel
		select {
		case receivedState := <-stateChangesCh:
			if receivedState != mockState {
				t.Errorf("Expected to receive the underlying state in the channel, but got a different state")
			}
		case <-time.After(100 * time.Millisecond):
			t.Errorf("Expected to receive a state change notification, but none was received")
		}
	})

	t.Run("MergeChange propagates error from underlying state", func(t *testing.T) {
		t.Log("MergeChange propagates error from underlying state test case")

		// Arrange
		expectedError := &mockError{message: "test error"}
		mockState := &mockState{mergeChangeError: expectedError}
		stateChangesCh := make(chan g.State, 1)
		wrapper := &stateWrapper{
			state:          mockState,
			stateChangesCh: stateChangesCh,
		}

		purpose := "test-purpose"
		value := "test-value"

		// Act
		err := wrapper.MergeChange(purpose, value)

		// Assert
		if err != expectedError {
			t.Errorf("Expected error %v, but got %v", expectedError, err)
		}

		if !mockState.mergeChangeCalled {
			t.Errorf("Expected MergeChange to be called on underlying state")
		}

		// Check that no notification was sent to the channel on error
		select {
		case <-stateChangesCh:
			t.Errorf("Expected no state change notification on error, but one was received")
		case <-time.After(50 * time.Millisecond):
			// This is expected - no notification should be sent on error
		}
	})

	t.Run("MergeChange works with nil values", func(t *testing.T) {
		t.Log("MergeChange works with nil values test case")

		// Arrange
		mockState := &mockState{}
		stateChangesCh := make(chan g.State, 1)
		wrapper := &stateWrapper{
			state:          mockState,
			stateChangesCh: stateChangesCh,
		}

		// Act
		err := wrapper.MergeChange(nil, nil)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, but got: %v", err)
		}

		if !mockState.mergeChangeCalled {
			t.Errorf("Expected MergeChange to be called on underlying state")
		}

		if mockState.purposeReceived != nil {
			t.Errorf("Expected purpose to be nil, but got %v", mockState.purposeReceived)
		}

		if mockState.valueReceived != nil {
			t.Errorf("Expected value to be nil, but got %v", mockState.valueReceived)
		}

		// Check that a notification was sent to the channel
		select {
		case <-stateChangesCh:
			// Expected
		case <-time.After(100 * time.Millisecond):
			t.Errorf("Expected to receive a state change notification, but none was received")
		}
	})

	t.Run("Multiple MergeChange calls send multiple notifications", func(t *testing.T) {
		t.Log("Multiple MergeChange calls send multiple notifications test case")

		// Arrange
		mockState := &mockState{}
		stateChangesCh := make(chan g.State, 3) // Buffer for multiple notifications
		wrapper := &stateWrapper{
			state:          mockState,
			stateChangesCh: stateChangesCh,
		}

		// Act
		err1 := wrapper.MergeChange("purpose1", "value1")
		err2 := wrapper.MergeChange("purpose2", "value2")
		err3 := wrapper.MergeChange("purpose3", "value3")

		// Assert
		if err1 != nil || err2 != nil || err3 != nil {
			t.Errorf("Expected no errors, but got: %v, %v, %v", err1, err2, err3)
		}

		// Check that three notifications were sent
		notificationCount := 0
		for i := 0; i < 3; i++ {
			select {
			case <-stateChangesCh:
				notificationCount++
			case <-time.After(100 * time.Millisecond):
				break
			}
		}

		if notificationCount != 3 {
			t.Errorf("Expected 3 notifications, but received %d", notificationCount)
		}
	})
}

// mockError is a simple error implementation for testing.
type mockError struct {
	message string
}

func (e *mockError) Error() string {
	return e.message
}
