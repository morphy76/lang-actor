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
	currentState framework.ActorState[actorState],
	sendFn framework.SendFn,
	me url.URL,
) (framework.ActorState[actorState], error) {
	var useMsg chatMessage = msg.(chatMessage)

	fmt.Println("-----------------------------------")
	fmt.Printf("I'm [%s] and I'm rocessing message from [%s]\n", me.Host, msg.Sender().Host)

	if useMsg.stopAfter < int(currentState.Cast().processedMessages) {
		fmt.Println("Current state:", currentState.Cast().processedMessages)
		fmt.Println("Stopping after:", useMsg.stopAfter)
		fmt.Println("Cancelling the actor")
		useMsg.cancelFn()
		fmt.Println("====================================")
		return currentState, nil
	}

	content := chatMessage{
		sender:    me,
		stopAfter: useMsg.stopAfter,
		cancelFn:  useMsg.cancelFn,
	}
	fmt.Println("Sending message to:", msg.Sender().Host)
	sendFn(content, msg.Sender())
	fmt.Println("-----------------------------------")
	return actorState{processedMessages: currentState.Cast().processedMessages + 1}, nil
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

	pongURL, _ := url.Parse("actor://pong")
	pongActor, err := builders.NewActor(ctx, *pongURL, pingPongFn, actorState{})
	if err != nil {
		fmt.Println("Error creating actor:", err)
		return
	}

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
