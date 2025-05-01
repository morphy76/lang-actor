package main

import (
	"bufio"
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

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
	sender  url.URL
	Message string
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

var echoFn framework.ProcessingFn[actorState] = func(
	msg framework.Message,
	actor framework.ActorView[actorState],
) (framework.ActorState[actorState], error) {
	var useMsg chatMessage = msg.(chatMessage)
	fmt.Printf("Echo [%s] after [%d] messages\n", useMsg.Message, actor.State().Cast().processedMessages)
	return actorState{processedMessages: actor.State().Cast().processedMessages + 1}, nil
}

func main() {

	actorCatalog := make(map[url.URL]framework.Transport)
	defer func() {
		for k := range actorCatalog {
			delete(actorCatalog, k)
		}
	}()

	ctx := context.WithValue(context.Background(), framework.ActorCatalogContextKey, actorCatalog)

	echoURL, _ := url.Parse("actor://echo")
	echoActor, err := builders.NewActor(ctx, *echoURL, echoFn, actorState{})
	if err != nil {
		fmt.Println("Error creating actor:", err)
		return
	}

	echoActor.Start()
	defer func() {
		done, _ := echoActor.Stop()
		<-done
	}()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Send messages to get an echo, send 'exit' to quit.")
	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "exit" {
			break
		}
		emitMessage := chatMessage{
			sender:  *echoURL,
			Message: input,
		}
		echoActor.Deliver(emitMessage)
	}
}
