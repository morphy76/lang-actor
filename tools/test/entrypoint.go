package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/morphy76/lang-actor/pkg/builders"
	"github.com/morphy76/lang-actor/pkg/framework"
)

var staticActorStatusAssertion framework.Payload[actorStatus] = (*actorStatus)(nil)

type actorStatus struct {
	processedMessages uint64
}

func (c actorStatus) Cast() actorStatus {
	return c
}

var staticChatMessageAssertion framework.Message[chatMessage] = (*chatMessage)(nil)

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

var echoFn framework.ProcessingFn[chatMessage, actorStatus] = func(msg framework.Message[chatMessage], currentState framework.Payload[actorStatus]) (framework.Payload[actorStatus], error) {
	fmt.Printf("Echo [%s] after [%d] messages\n", msg.Cast().Message, currentState.Cast().processedMessages)
	return actorStatus{processedMessages: currentState.Cast().processedMessages + 1}, nil
}

func main() {

	echoURL, _ := url.Parse("actor://echo")
	echoActor, _ := builders.NewActor(*echoURL, echoFn, actorStatus{})
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
