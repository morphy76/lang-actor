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
	ints []int
}

var staticChatMessageAssertion framework.Message = (*inputMessage)(nil)

type inputMessage struct {
}

func (m inputMessage) Sender() url.URL {
	return url.URL{}
}

func (m inputMessage) Mutation() bool {
	return true
}

var sortFn framework.ProcessingFn[actorState] = func(
	msg framework.Message,
	me framework.ActorView[actorState],
) (actorState, error) {
	if len(me.State().ints) == 1 {
		return actorState{me.State().ints}, nil
	} else if len(me.State().ints) == 2 {
		useInts := me.State().ints
		if useInts[0] > useInts[1] {
			useInts[0], useInts[1] = useInts[1], useInts[0]
		}
		return actorState{ints: useInts}, nil
	} else {
		// split the array in half
		mid := len(me.State().ints) / 2
		left := me.State().ints[:mid]
		right := me.State().ints[mid:]
		leftActor, _ := builders.SpawnChild(me.(framework.ActorRef), sortFn, actorState{ints: left})
		rightActor, _ := builders.SpawnChild(me.(framework.ActorRef), sortFn, actorState{ints: right})
		leftActor.Deliver(inputMessage{})
		rightActor.Deliver(inputMessage{})

	}

	return actorState{}, nil
}

func main() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Type a list of numbers and press enter to sort them.")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	parts := strings.Split(input, " ")
	partsAsIntegers := convertElementsToIntegers(parts)

	sorterURL, _ := url.Parse("actor://echo")
	sorterActor, err := builders.NewActor(*sorterURL, sortFn, actorState{ints: partsAsIntegers})
	if err != nil {
		fmt.Println("Error creating actor:", err)
		return
	}

	sorterActor.Start()
	defer func() {
		done, _ := sorterActor.Stop()
		<-done
		fmt.Println("Sorted integers:", sorterActor.State().ints)
	}()

	sorterActor.Deliver(inputMessage{})
}

func convertElementsToIntegers(parts []string) []int {
	var integers []int
	for _, part := range parts {
		var integer int
		fmt.Sscanf(part, "%d", &integer)
		integers = append(integers, integer)
	}
	return integers
}
