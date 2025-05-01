package main

import (
	"bufio"
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/morphy76/lang-actor/pkg/builders"
	"github.com/morphy76/lang-actor/pkg/framework"
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

var pingPongFn framework.ProcessingFn[actorState] = func(
	msg framework.Message,
	actor framework.ActorView[actorState],
) (framework.ActorState[actorState], error) {
	var useMsg chatMessage = msg.(chatMessage)

	fmt.Println("-----------------------------------")
	fmt.Printf("I'm [%s] and I'm rocessing message from [%s]\n", actor.Address().Host, msg.Sender().Host)

	if useMsg.stopAfter < int(actor.State().Cast().processedMessages) {
		fmt.Println("Current state:", actor.State().Cast().processedMessages)
		fmt.Println("Stopping after:", useMsg.stopAfter)
		fmt.Println("Cancelling the actor")
		useMsg.cancelFn()
		fmt.Println("====================================")
		return actor.State(), nil
	}

	content := chatMessage{
		sender:    actor.Address(),
		stopAfter: useMsg.stopAfter,
		cancelFn:  useMsg.cancelFn,
	}
	fmt.Println("Sending message to:", msg.Sender().Host)
	actor.Send(content, msg.Sender())
	fmt.Println("-----------------------------------")
	return actorState{processedMessages: actor.State().Cast().processedMessages + 1}, nil
}

func main() {

	actorCatalog := make(map[url.URL]framework.Transport)
	defer func() {
		for k := range actorCatalog {
			delete(actorCatalog, k)
		}
	}()

	ctx, cancelFn := context.WithCancel(context.WithValue(context.Background(), framework.ActorCatalogContextKey, actorCatalog))

	pingURL, _ := url.Parse("actor://ping")
	pingActor, err := builders.NewActor(ctx, *pingURL, pingPongFn, actorState{})
	if err != nil {
		fmt.Println("Error creating actor:", err)
		return
	}

	pingActor.Start()
	defer func() {
		done, _ := pingActor.Stop()
		<-done
	}()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Press a enter to start")

	reader.ReadString('\n')
	initialMessage := chatMessage{
		stopAfter: 5,
		cancelFn:  cancelFn,
		sender:    *pingURL,
	}
	pingActor.Deliver(initialMessage)

	<-ctx.Done()
}
