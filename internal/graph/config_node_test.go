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

func TestConfigNode(t *testing.T) {
	t.Log("ConfigNode test suite")

	randomKey1 := uuid.NewString()
	randomKey2 := uuid.NewString()
	randomKey3 := uuid.NewString()
	randomValue1 := uuid.NewString()
	randomValue2 := uuid.NewString()
	randomValue3 := uuid.NewString()

	// Define test configuration
	testConfig := map[string]any{
		randomKey1: randomValue1,
		randomKey2: randomValue2,
		randomKey3: randomValue3,
	}

	t.Run("Config node creation", func(t *testing.T) {
		_, err := graph.NewConfigNode(testConfig, uuid.NewString())
		if err != nil {
			t.Fatalf("Failed to create config node: %v", err)
		}
	})

	t.Run("Config node keys request to get a keys response", func(t *testing.T) {
		t.Log("Config node keys request to get a keys response")

		responseCh := make(chan []string)
		clientURL, err := url.Parse("actor://client")
		clientActor, err := framework.NewActor(
			*clientURL,
			func(msg f.Message, self f.Actor[chan []string]) (chan []string, error) {
				useMex := msg.(*g.ConfigMessage)
				self.State() <- useMex.Keys
				return self.State(), nil
			},
			responseCh,
		)

		configNode, err := graph.NewConfigNode(testConfig, uuid.NewString())
		if err != nil {
			t.Fatalf("Failed to create config node: %v", err)
		}

		addressBook := routing.NewAddressBook()
		addressBook.Register(clientActor)
		addressBook.Register(configNode)
		configNode.SetResolver(addressBook)

		request, err := g.NewConfigMessage(*clientURL, g.Keys)
		configNode.Deliver(request)

		cfgKeys := <-responseCh

		if len(cfgKeys) != len(testConfig) {
			t.Fatalf("Expected %d keys, got %d", len(testConfig), len(cfgKeys))
		}

		for _, key := range cfgKeys {
			if _, ok := testConfig[key]; !ok {
				t.Fatalf("Key %s not found in config", key)
			}
		}
	})

	t.Run("Config node entries request to get all entries", func(t *testing.T) {
		t.Log("Config node entries request to get all entries")

		responseCh := make(chan map[string]any)
		clientURL, err := url.Parse("actor://client")
		clientActor, err := framework.NewActor(
			*clientURL,
			func(msg f.Message, self f.Actor[chan map[string]any]) (chan map[string]any, error) {
				useMex := msg.(*g.ConfigMessage)
				self.State() <- useMex.Entries
				return self.State(), nil
			},
			responseCh,
		)

		configNode, err := graph.NewConfigNode(testConfig, uuid.NewString())
		if err != nil {
			t.Fatalf("Failed to create config node: %v", err)
		}

		addressBook := routing.NewAddressBook()
		addressBook.Register(clientActor)
		addressBook.Register(configNode)
		configNode.SetResolver(addressBook)

		request, err := g.NewConfigMessage(*clientURL, g.Entries)
		configNode.Deliver(request)

		entries := <-responseCh

		if len(entries) != len(testConfig) {
			t.Fatalf("Expected %d entries, got %d", len(testConfig), len(entries))
		}

		for key, value := range entries {
			expectedValue, ok := testConfig[key]
			if !ok {
				t.Fatalf("Key %s not found in config", key)
			}
			if value != expectedValue {
				t.Fatalf("Expected value %v for key %s, got %v", expectedValue, key, value)
			}
		}
	})

	t.Run("Config node request to get a single value", func(t *testing.T) {
		t.Log("Config node request to get a single value")

		responseCh := make(chan any)
		clientURL, err := url.Parse("actor://client")
		clientActor, err := framework.NewActor(
			*clientURL,
			func(msg f.Message, self f.Actor[chan any]) (chan any, error) {
				useMex := msg.(*g.ConfigMessage)
				self.State() <- useMex.Value
				return self.State(), nil
			},
			responseCh,
		)

		configNode, err := graph.NewConfigNode(testConfig, uuid.NewString())
		if err != nil {
			t.Fatalf("Failed to create config node: %v", err)
		}

		addressBook := routing.NewAddressBook()
		addressBook.Register(clientActor)
		addressBook.Register(configNode)
		configNode.SetResolver(addressBook)

		request, err := g.NewConfigMessage(*clientURL, g.Request, randomKey1)
		configNode.Deliver(request)

		result := <-responseCh

		if result != randomValue1 {
			t.Fatalf("Expected value %v for key %s, got %v", randomValue1, randomKey1, result)
		}
	})

	t.Run("Config node request to get multiple values", func(t *testing.T) {
		t.Log("Config node request to get multiple values")

		responseCh := make(chan map[string]any)
		clientURL, err := url.Parse("actor://client")
		clientActor, err := framework.NewActor(
			*clientURL,
			func(msg f.Message, self f.Actor[chan map[string]any]) (chan map[string]any, error) {
				useMex := msg.(*g.ConfigMessage)
				self.State() <- useMex.Entries
				return self.State(), nil
			},
			responseCh,
		)

		configNode, err := graph.NewConfigNode(testConfig, uuid.NewString())
		if err != nil {
			t.Fatalf("Failed to create config node: %v", err)
		}

		addressBook := routing.NewAddressBook()
		addressBook.Register(clientActor)
		addressBook.Register(configNode)
		configNode.SetResolver(addressBook)

		request, err := g.NewConfigMessage(*clientURL, g.Request, randomKey1, randomKey2)
		configNode.Deliver(request)

		result := <-responseCh

		if len(result) != 2 {
			t.Fatalf("Expected 2 entries, got %d", len(result))
		}

		value1, ok := result[randomKey1]
		if !ok {
			t.Fatalf("Key %s not found in response", randomKey1)
		}
		if value1 != randomValue1 {
			t.Fatalf("Expected value %v for key %s, got %v", randomValue1, randomKey1, value1)
		}

		value2, ok := result[randomKey2]
		if !ok {
			t.Fatalf("Key %s not found in response", randomKey2)
		}
		if value2 != randomValue2 {
			t.Fatalf("Expected value %v for key %s, got %v", randomValue2, randomKey2, value2)
		}
	})

	t.Run("Config node request with unknown key", func(t *testing.T) {
		t.Log("Config node request with unknown key")

		responseCh := make(chan map[string]any)
		clientURL, err := url.Parse("actor://client")
		clientActor, err := framework.NewActor(
			*clientURL,
			func(msg f.Message, self f.Actor[chan map[string]any]) (chan map[string]any, error) {
				useMex := msg.(*g.ConfigMessage)
				self.State() <- useMex.Entries
				return self.State(), nil
			},
			responseCh,
		)

		configNode, err := graph.NewConfigNode(testConfig, uuid.NewString())
		if err != nil {
			t.Fatalf("Failed to create config node: %v", err)
		}

		addressBook := routing.NewAddressBook()
		addressBook.Register(clientActor)
		addressBook.Register(configNode)
		configNode.SetResolver(addressBook)

		unknownKey := "unknown-key-" + uuid.NewString()
		request, err := g.NewConfigMessage(*clientURL, g.Request, unknownKey)
		configNode.Deliver(request)

		result := <-responseCh

		if len(result) != 0 {
			t.Fatalf("Expected 0 entries for unknown key, got %d", len(result))
		}
	})
}
