package graph_test

import (
	"net/url"
	"testing"

	"github.com/google/uuid"
	"github.com/morphy76/lang-actor/internal/framework"
	"github.com/morphy76/lang-actor/internal/graph"
	"github.com/morphy76/lang-actor/internal/routing"
	f "github.com/morphy76/lang-actor/pkg/framework"
	g "github.com/morphy76/lang-actor/pkg/graph"
)

func TestStatusNode(t *testing.T) {
	t.Log("StatusNode test suite")

	type testStatus struct {
		Value     string
		Count     int
		IsEnabled bool
	}

	// Define initial test status
	initialStatus := testStatus{
		Value:     uuid.NewString(),
		Count:     42,
		IsEnabled: true,
	}

	t.Run("Status node creation", func(t *testing.T) {
		_, err := graph.NewStatusNode(initialStatus, uuid.NewString())
		if err != nil {
			t.Fatalf("Failed to create status node: %v", err)
		}
	})

	t.Run("Status node responds to status requests", func(t *testing.T) {
		t.Log("Status node should respond to status requests")

		responseCh := make(chan testStatus)
		clientURL, err := url.Parse("actor://client")
		if err != nil {
			t.Fatalf("Failed to parse client URL: %v", err)
		}

		clientActor, err := framework.NewActor(
			*clientURL,
			func(msg f.Message, self f.Actor[chan testStatus]) (chan testStatus, error) {
				if useMex, ok := msg.(*g.StatusMessage[testStatus]); ok {
					self.State() <- useMex.Value
				}
				return self.State(), nil
			},
			responseCh,
			true,
		)
		if err != nil {
			t.Fatalf("Failed to create client actor: %v", err)
		}

		statusNode, err := graph.NewStatusNode(initialStatus, uuid.NewString())
		if err != nil {
			t.Fatalf("Failed to create status node: %v", err)
		}

		addressBook := routing.NewAddressBook()
		addressBook.Register(clientActor)
		addressBook.Register(statusNode)
		statusNode.SetResolver(addressBook)

		// Create a request message
		request := &g.StatusMessage[testStatus]{
			From:              *clientURL,
			StatusMessageType: g.StatusRequest,
		}

		// Send the request
		err = statusNode.Deliver(request)
		if err != nil {
			t.Fatalf("Failed to deliver status request message: %v", err)
		}

		// Wait for response
		response := <-responseCh

		// Verify response
		if response.Value != initialStatus.Value {
			t.Errorf("Expected Value %s, got %s", initialStatus.Value, response.Value)
		}
		if response.Count != initialStatus.Count {
			t.Errorf("Expected Count %d, got %d", initialStatus.Count, response.Count)
		}
		if response.IsEnabled != initialStatus.IsEnabled {
			t.Errorf("Expected IsEnabled %t, got %t", initialStatus.IsEnabled, response.IsEnabled)
		}
	})

	t.Run("Status node handles status updates", func(t *testing.T) {
		t.Log("Status node should update its state when receiving an update message")

		responseCh := make(chan testStatus)
		clientURL, err := url.Parse("actor://client")
		if err != nil {
			t.Fatalf("Failed to parse client URL: %v", err)
		}

		clientActor, err := framework.NewActor(
			*clientURL,
			func(msg f.Message, self f.Actor[chan testStatus]) (chan testStatus, error) {
				if useMex, ok := msg.(*g.StatusMessage[testStatus]); ok {
					self.State() <- useMex.Value
				}
				return self.State(), nil
			},
			responseCh,
			true,
		)
		if err != nil {
			t.Fatalf("Failed to create client actor: %v", err)
		}

		statusNode, err := graph.NewStatusNode(initialStatus, uuid.NewString())
		if err != nil {
			t.Fatalf("Failed to create status node: %v", err)
		}

		addressBook := routing.NewAddressBook()
		addressBook.Register(clientActor)
		addressBook.Register(statusNode)
		statusNode.SetResolver(addressBook)

		// Updated status
		updatedStatus := testStatus{
			Value:     "updated-" + uuid.NewString(),
			Count:     99,
			IsEnabled: false,
		}

		// Create an update message
		updateMsg := &g.StatusMessage[testStatus]{
			From:              *clientURL,
			StatusMessageType: g.StatusUpdate,
			Value:             updatedStatus,
		}

		// Send the update
		err = statusNode.Deliver(updateMsg)
		if err != nil {
			t.Fatalf("Failed to deliver status update message: %v", err)
		}

		// Now send a request to check if the status was updated
		requestMsg := &g.StatusMessage[testStatus]{
			From:              *clientURL,
			StatusMessageType: g.StatusRequest,
		}

		err = statusNode.Deliver(requestMsg)
		if err != nil {
			t.Fatalf("Failed to deliver status request message: %v", err)
		}

		// Wait for response
		response := <-responseCh

		// Verify updated status
		if response.Value != updatedStatus.Value {
			t.Errorf("Expected updated Value %s, got %s", updatedStatus.Value, response.Value)
		}
		if response.Count != updatedStatus.Count {
			t.Errorf("Expected updated Count %d, got %d", updatedStatus.Count, response.Count)
		}
		if response.IsEnabled != updatedStatus.IsEnabled {
			t.Errorf("Expected updated IsEnabled %t, got %t", updatedStatus.IsEnabled, response.IsEnabled)
		}
	})

	t.Run("Status node routing errors", func(t *testing.T) {
		t.Log("Status node should reject routing operations")

		statusNode, err := graph.NewStatusNode(initialStatus, uuid.NewString())
		if err != nil {
			t.Fatalf("Failed to create status node: %v", err)
		}

		dummyNode, err := graph.NewStatusNode("dummy", uuid.NewString())
		if err != nil {
			t.Fatalf("Failed to create dummy node: %v", err)
		}

		// Try one-way route (should fail)
		err = statusNode.OneWayRoute("test", dummyNode)
		if err == nil {
			t.Errorf("Expected error for OneWayRoute, got nil")
		}

		// Try two-way route (should fail)
		err = statusNode.TwoWayRoute("test", dummyNode)
		if err == nil {
			t.Errorf("Expected error for TwoWayRoute, got nil")
		}

		// Try ProceedOnAnyRoute (should fail)
		msg, err := g.NewStatusMessageRequest[any](dummyNode.Address())
		if err != nil {
			t.Fatalf("Failed to create status message: %v", err)
		}
		err = statusNode.ProceedOnAnyRoute(msg)
		if err == nil {
			t.Errorf("Expected error for ProceedOnAnyRoute, got nil")
		}
	})
}
