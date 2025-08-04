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

var echoFn framework.ProcessingFn[actorState] = func(
	msg framework.Message,
	actor framework.Actor[actorState],
) (actorState, error) {
	useMsg, ok := msg.Payload().(string)
	if !ok {
		return actorState{}, fmt.Errorf("expected string payload, got %T", msg.Payload())
	}
	fmt.Printf("Echo [%s] after [%d] messages\n", useMsg, actor.State().processedMessages)
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
		echoActor.Deliver(input, nil)
	}
}
