package routing_test

import (
	"net/url"
	"testing"

	"github.com/morphy76/lang-actor/internal/routing"
	f "github.com/morphy76/lang-actor/pkg/framework"
	"gotest.tools/v3/assert"
)

var staticMockActorAssertion f.Actor[any] = (*mockActor)(nil)

type mockActor struct {
	address url.URL
}

func (m *mockActor) Address() url.URL {
	return m.address
}

func (m *mockActor) Start() error {
	return nil
}

func (m *mockActor) Stop() (chan bool, error) {
	ch := make(chan bool)
	close(ch)
	return ch, nil
}

func (m *mockActor) Deliver(msg f.Message) error {
	return nil
}

func (m *mockActor) Status() f.ActorStatus {
	return f.ActorStatusIdle
}

func (m *mockActor) State() any {
	return nil
}

func (m *mockActor) Send(msg f.Message, transport f.Transport) error {
	return nil
}

const actorURI = "actor://example"

func TestAddressBookRegister(t *testing.T) {
	t.Log("AddressBook Register test suite")

	t.Run("Register a new actor successfully", func(t *testing.T) {
		t.Log("Should register a new actor successfully")

		addressBook := routing.NewAddressBook()
		address, err := url.Parse(actorURI)
		assert.NilError(t, err)

		actor := &mockActor{address: *address}
		err = addressBook.Register(actor)
		assert.NilError(t, err)
	})

	t.Run("Register an actor that is already registered", func(t *testing.T) {
		t.Log("Should return an error when registering an actor that is already registered")

		addressBook := routing.NewAddressBook()
		address, err := url.Parse(actorURI)
		assert.NilError(t, err)

		actor := &mockActor{address: *address}
		err = addressBook.Register(actor)
		assert.NilError(t, err)

		err = addressBook.Register(actor)
		assert.ErrorContains(t, err, "actor already registered")
	})
}

func TestAddressBookLookup(t *testing.T) {
	t.Log("AddressBook Lookup test suite")

	t.Run("Lookup an actor successfully", func(t *testing.T) {
		t.Log("Should find an actor by its address")

		addressBook := routing.NewAddressBook()
		address, err := url.Parse(actorURI)
		assert.NilError(t, err)

		actor := &mockActor{address: *address}
		err = addressBook.Register(actor)
		assert.NilError(t, err)

		foundActor, err := addressBook.Lookup(*address)
		assert.NilError(t, err)
		assert.Equal(t, foundActor.Address(), actor.Address())
	})

	t.Run("Lookup an actor that does not exist", func(t *testing.T) {
		t.Log("Should return false when looking up an actor that does not exist")

		addressBook := routing.NewAddressBook()
		address, err := url.Parse("actor://nonexistent")
		assert.NilError(t, err)

		_, err = addressBook.Lookup(*address)
		assert.ErrorContains(t, err, "actor not found")
	})
}
