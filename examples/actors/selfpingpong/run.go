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

type actorState struct {
	processedMessages uint64
}

type chatMessage struct {
	cancelFn  context.CancelFunc
	stopAfter int
}

var pingPongFn framework.ProcessingFn[actorState] = func(
	msg framework.Message,
	self framework.Actor[actorState],
) (actorState, error) {
	var useMsg chatMessage = msg.Payload().(chatMessage)

	fmt.Println("-----------------------------------")
	fmt.Printf("I'm [%s] and I'm rocessing message from [%s]\n", self.Address().Host, msg.Sender().Host)

	if useMsg.stopAfter < int(self.State().processedMessages) {
		fmt.Println("Current state:", self.State().processedMessages)
		fmt.Println("Stopping after:", useMsg.stopAfter)
		fmt.Println("Cancelling the actor")
		useMsg.cancelFn()
		fmt.Println("====================================")
		return self.State(), nil
	}

	content := chatMessage{
		stopAfter: useMsg.stopAfter,
		cancelFn:  useMsg.cancelFn,
	}
	fmt.Println("Sending message to:", msg.Sender().Host)

	self.Send(content, self)
	fmt.Println("-----------------------------------")
	return actorState{processedMessages: self.State().processedMessages + 1}, nil
}

func main() {

	ctx, cancelFn := context.WithCancel(context.Background())

	pingURL, _ := url.Parse("actor://ping")
	pingActor, err := builders.NewActor(*pingURL, pingPongFn, actorState{})
	if err != nil {
		fmt.Println("Error creating actor:", err)
		return
	}

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
	}
	pingActor.Deliver(initialMessage, pingActor)

	<-ctx.Done()
}
