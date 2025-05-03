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

type actorState struct {
	processedMessages uint64
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

var echoFn framework.ProcessingFn[actorState] = func(
	msg framework.Message,
	actor framework.Actor[actorState],
) (actorState, error) {
	var useMsg chatMessage = msg.(chatMessage)
	fmt.Printf("Echo [%s] after [%d] messages\n", useMsg.Message, actor.State().processedMessages)
	return actorState{processedMessages: actor.State().processedMessages + 1}, nil
}

func main() {

	echoURL, _ := url.Parse("actor://echo")
	echoActor, err := builders.NewActor(*echoURL, echoFn, actorState{})
	if err != nil {
		fmt.Println("Error creating actor:", err)
		return
	}

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
