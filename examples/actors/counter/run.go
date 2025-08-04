package main

import (
	"fmt"
	"net/url"
	"time"

	"github.com/morphy76/lang-actor/pkg/builders"
	"github.com/morphy76/lang-actor/pkg/framework"
)

// Define actor state
type counterState struct {
	count int
}

// Define a message type
type incrementMessage struct {
	amount int
}

// Define message processing function
var counterFn framework.ProcessingFn[counterState] = func(
	msg framework.Message,
	self framework.Actor[counterState],
) (counterState, error) {
	if incMsg, ok := msg.Payload().(incrementMessage); ok {
		newCount := self.State().count + incMsg.amount
		fmt.Printf("Counter incremented to: %d\n", newCount)
		return counterState{count: newCount}, nil
	}
	return self.State(), nil
}

func main() {
	// Create actor
	actorURL, _ := url.Parse("actor://counter")
	counterActor, err := builders.NewActor(*actorURL, counterFn, counterState{count: 0})
	if err != nil {
		fmt.Println("Error creating actor:", err)
		return
	}

	// Ensure actor stops when program ends
	defer func() {
		done, _ := counterActor.Stop()
		<-done
	}()

	// Send messages to the actor
	for i := 1; i <= 5; i++ {
		msg := incrementMessage{
			amount: i,
		}
		counterActor.Deliver(msg, nil)
		time.Sleep(500 * time.Millisecond)
	}
}
