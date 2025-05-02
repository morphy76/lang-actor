package main

import (
	"bufio"
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/morphy76/lang-actor/pkg/builders"
	"github.com/morphy76/lang-actor/pkg/framework"
	"github.com/morphy76/lang-actor/pkg/routing"
)

var staticActorStatusAssertion framework.ActorState[actorState] = (*actorState)(nil)

type actorState struct {
	processedMessages uint64
}

func (c actorState) Cast() actorState {
	return c
}

var staticChatMessageAssertion framework.Message = (*chatMessage)(nil)

type chatMessage struct {
	cancelFn  context.CancelFunc
	sender    url.URL
	stopAfter int
}

func (m chatMessage) Sender() url.URL {
	return m.sender
}

func (m chatMessage) Mutation() bool {
	return true
}

func (m chatMessage) Cast() chatMessage {
	return m
}

func getPingPongFn(addressBook routing.AddressBook) framework.ProcessingFn[actorState] {
	return func(
		msg framework.Message,
		self framework.ActorView[actorState],
	) (framework.ActorState[actorState], error) {
		var useMsg chatMessage = msg.(chatMessage)

		fmt.Println("-----------------------------------")
		fmt.Printf("I'm [%s] and I'm rocessing message from [%s]\n", self.Address().Host, msg.Sender().Host)

		if useMsg.stopAfter < int(self.State().Cast().processedMessages) {
			fmt.Println("Current state:", self.State().Cast().processedMessages)
			fmt.Println("Stopping after:", useMsg.stopAfter)
			fmt.Println("Cancelling the actor")
			useMsg.cancelFn()
			fmt.Println("====================================")
			return self.State(), nil
		}

		content := chatMessage{
			sender:    self.Address(),
			stopAfter: useMsg.stopAfter,
			cancelFn:  useMsg.cancelFn,
		}
		fmt.Println("Sending message to:", msg.Sender().Host)

		transport, _ := addressBook.Lookup(msg.Sender())

		self.Send(content, transport)
		fmt.Println("-----------------------------------")
		return actorState{processedMessages: self.State().Cast().processedMessages + 1}, nil
	}
}

func main() {

	addressBook := builders.NewAddressBook()
	defer addressBook.TearDown()

	pingPongFn := getPingPongFn(addressBook)

	ctx, cancelFn := context.WithCancel(context.Background())

	pingURL, _ := url.Parse("actor://ping")
	pingActor, err := builders.NewActor(*pingURL, pingPongFn, actorState{})
	if err != nil {
		fmt.Println("Error creating actor:", err)
		return
	}

	pingActor.Start()
	defer func() {
		done, _ := pingActor.Stop()
		<-done
	}()

	pongURL, _ := url.Parse("actor://pong")
	pongActor, err := builders.NewActor(*pongURL, pingPongFn, actorState{})
	if err != nil {
		fmt.Println("Error creating actor:", err)
		return
	}

	addressBook.Register(pingActor)
	addressBook.Register(pongActor)

	pongActor.Start()
	defer func() {
		done, _ := pongActor.Stop()
		<-done
	}()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Press a enter to start")

	reader.ReadString('\n')
	initialMessage := chatMessage{
		stopAfter: 5,
		cancelFn:  cancelFn,
		sender:    *pongURL,
	}
	pingActor.Deliver(initialMessage)

	<-ctx.Done()
}
