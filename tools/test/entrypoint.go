package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/morphy76/lang-actor/internal/framework"
)

var staticActorStatusAssertion framework.Payload[actorStatus] = (*actorStatus)(nil)

type actorStatus struct {
	processedMessages uint64
}

func (c actorStatus) ToImplementation() actorStatus {
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

func (m chatMessage) ToImplementation() chatMessage {
	return m
}

func main() {

	var echoFn framework.ProcessingFn[chatMessage, actorStatus] = func(msg framework.Message[chatMessage], currentState framework.Payload[actorStatus]) (framework.Payload[actorStatus], error) {
		fmt.Printf("Echo: %s\n", msg.ToImplementation().Message)
		return actorStatus{processedMessages: currentState.ToImplementation().processedMessages + 1}, nil
	}

	echoURL, _ := url.Parse("actor://echo")
	echoActor, _ := framework.NewActor(*echoURL, echoFn, actorStatus{})
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
