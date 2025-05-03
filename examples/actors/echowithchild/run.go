// With Full Vibes (Github Copilot using Claude 3.7 Sonnet)
// prompt:
// Replicate the echo example in a new example called echowithchild where:
// - the echo message is prepared by the child
// - uses the request/reply pattern
// - the echo message is printed by the main actor
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

// State for the main actor
type mainActorState struct {
	processedMessages uint64
}

// State for the child actor
type childActorState struct {
	processedRequests uint64
}

// Message types for our request-reply pattern
type messageType int

const (
	messageTypeRequest messageType = iota
	messageTypeReply
)

// Echo message that will be used for both request and reply
type echoMessage struct {
	sender    url.URL
	mexType   messageType
	content   string
	decorated string // The prepared/decorated message (filled by child)
}

func (m echoMessage) Sender() url.URL {
	return m.sender
}

func (m echoMessage) Mutation() bool {
	return true
}

// Main actor processing function
var mainActorFn framework.ProcessingFn[mainActorState] = func(
	msg framework.Message,
	actor framework.Actor[mainActorState],
) (mainActorState, error) {
	useMsg, ok := msg.(echoMessage)
	if !ok {
		return actor.State(), fmt.Errorf("unexpected message type")
	}

	switch useMsg.mexType {
	case messageTypeRequest:
		// Create child actor for message preparation if it doesn't exist
		childActor, err := builders.SpawnChild(actor, childActorFn, childActorState{})
		if err != nil {
			return actor.State(), fmt.Errorf("error creating child actor: %w", err)
		}

		// Send request to child actor
		requestMsg := echoMessage{
			sender:  actor.Address(),
			mexType: messageTypeRequest,
			content: useMsg.content,
		}
		err = actor.Send(requestMsg, childActor)
		if err != nil {
			return actor.State(), fmt.Errorf("error sending message to child: %w", err)
		}
		return actor.State(), nil
	case messageTypeReply:
		// Display the prepared echo message received from child
		fmt.Printf("Echo: [%s] (processed %d messages)\n",
			useMsg.decorated,
			actor.State().processedMessages,
		)
		return mainActorState{processedMessages: actor.State().processedMessages + 1}, nil
	}
	return actor.State(), nil
}

// Child actor processing function
var childActorFn framework.ProcessingFn[childActorState] = func(
	msg framework.Message,
	actor framework.Actor[childActorState],
) (childActorState, error) {
	useMsg, ok := msg.(echoMessage)
	if !ok {
		return actor.State(), fmt.Errorf("unexpected message type")
	}

	if useMsg.mexType == messageTypeRequest {
		// Prepare the decorated message
		decoratedContent := fmt.Sprintf("✨ %s ✨", strings.ToUpper(useMsg.content))

		// Create reply message
		replyMsg := echoMessage{
			sender:    actor.Address(),
			mexType:   messageTypeReply,
			content:   useMsg.content,
			decorated: decoratedContent,
		}

		// Send reply to parent
		if parent, found := actor.GetParent(); found {
			err := actor.Send(replyMsg, parent)
			if err != nil {
				return actor.State(), fmt.Errorf("error sending reply to parent: %w", err)
			}
		}
	}

	// Update state with incremented request count
	return childActorState{processedRequests: actor.State().processedRequests + 1}, nil
}

func main() {
	// Create the main echo actor
	echoURL, _ := url.Parse("actor://echo")
	mainActor, err := builders.NewActor(*echoURL, mainActorFn, mainActorState{})
	if err != nil {
		fmt.Println("Error creating actor:", err)
		return
	}

	// Ensure we stop the actor when done
	defer func() {
		done, _ := mainActor.Stop()
		<-done
	}()

	// Read user input
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Send messages for echo with child preparation, send 'exit' to quit.")

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" {
			break
		}

		// Create and deliver the initial message
		initialMsg := echoMessage{
			sender:  *echoURL,
			mexType: messageTypeRequest,
			content: input,
		}
		mainActor.Deliver(initialMsg)
	}
}
