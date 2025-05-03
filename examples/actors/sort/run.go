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

type sorterState struct {
	sortedInts []int
}

var staticinputMessageAssertion framework.Message = (*inputMessage)(nil)

type inputMessageType int

const (
	inputMessageTypeSort inputMessageType = iota
	inputMessageTypeMerge
	inputMessageTypeEnd
	inputMessageTypeLeave
)

type inputMessage struct {
	ints    []int
	mexType inputMessageType
}

func (m inputMessage) Sender() url.URL {
	return url.URL{}
}

func (m inputMessage) Mutation() bool {
	return m.mexType == inputMessageTypeMerge
}

func sortFn(leaveCh chan bool) framework.ProcessingFn[sorterState] {
	return func(msg framework.Message, me framework.Actor[sorterState]) (sorterState, error) {
		useMsg := msg.(inputMessage)
		switch useMsg.mexType {
		case inputMessageTypeSort:
			return sortFnSort(leaveCh, useMsg, me)
		case inputMessageTypeMerge:
			return sortFnMerge(useMsg, me)
		case inputMessageTypeEnd:
			fmt.Println("Sort completed:", useMsg.ints)
			me.Deliver(inputMessage{mexType: inputMessageTypeLeave})
			return sorterState{}, nil
		case inputMessageTypeLeave:
			fmt.Println("Exiting...")
			leaveCh <- true
			return sorterState{}, nil
		default:
			return sorterState{}, fmt.Errorf("unknown message type")
		}
	}
}

func sortFnMerge(
	msg inputMessage,
	me framework.Actor[sorterState],
) (sorterState, error) {
	rv := sorterState{}
	if me.State().sortedInts == nil {
		rv.sortedInts = make([]int, len(msg.ints))
		copy(rv.sortedInts, msg.ints)
	} else {
		rv.sortedInts = make([]int, 0, len(me.State().sortedInts)+len(msg.ints))
		i, j := 0, 0

		for i < len(me.State().sortedInts) && j < len(msg.ints) {
			if me.State().sortedInts[i] <= msg.ints[j] {
				rv.sortedInts = append(rv.sortedInts, me.State().sortedInts[i])
				i++
			} else {
				rv.sortedInts = append(rv.sortedInts, msg.ints[j])
				j++
			}
		}

		// Append remaining elements from either slice
		rv.sortedInts = append(rv.sortedInts, me.State().sortedInts[i:]...)
		rv.sortedInts = append(rv.sortedInts, msg.ints[j:]...)

		parent, found := me.GetParent()
		if found {
			me.Send(inputMessage{ints: rv.sortedInts, mexType: inputMessageTypeMerge}, parent)
		} else {
			me.Deliver(inputMessage{ints: rv.sortedInts, mexType: inputMessageTypeEnd})
		}
	}

	return rv, nil
}

func sortFnSort(
	leaveCh chan bool,
	msg inputMessage,
	me framework.Actor[sorterState],
) (sorterState, error) {
	if len(msg.ints) < 2 {
		parent, found := me.GetParent()
		if found {
			err := me.Send(inputMessage{ints: msg.ints, mexType: inputMessageTypeMerge}, parent)
			if err != nil {
				fmt.Println("Error sending merge upstream", err)
			}
		} else {
			err := me.Deliver(inputMessage{ints: msg.ints, mexType: inputMessageTypeEnd})
			if err != nil {
				fmt.Println("Error sending end to me", err)
			}
		}
		return sorterState{}, nil
	} else if len(msg.ints) == 2 {
		useInts := msg.ints
		if useInts[0] > useInts[1] {
			useInts[0], useInts[1] = useInts[1], useInts[0]
		}
		parent, found := me.GetParent()
		if found {
			err := me.Send(inputMessage{ints: msg.ints, mexType: inputMessageTypeMerge}, parent)
			if err != nil {
				fmt.Println("Error sending merge upstream", err)
			}
		} else {
			err := me.Deliver(inputMessage{ints: msg.ints, mexType: inputMessageTypeEnd})
			if err != nil {
				fmt.Println("Error sending end to me", err)
			}
		}
		return sorterState{}, nil
	}

	mid := len(msg.ints) / 2
	fmt.Println("Splitting", msg.ints, "at", mid)

	childState := sorterState{}

	leftActor, err := builders.SpawnChild(me, sortFn(leaveCh), childState)
	if err != nil {
		fmt.Println("Error creating left actor:", err)
		me.Deliver(inputMessage{mexType: inputMessageTypeLeave})
	}

	err = leftActor.Deliver(inputMessage{ints: msg.ints[:mid], mexType: inputMessageTypeSort})
	if err != nil {
		fmt.Println("Error sending sort left downstream", err)
		me.Deliver(inputMessage{mexType: inputMessageTypeLeave})
	}
	rightActor, err := builders.SpawnChild(me, sortFn(leaveCh), childState)
	if err != nil {
		fmt.Println("Error creating left actor:", err)
		me.Deliver(inputMessage{mexType: inputMessageTypeLeave})
	}

	err = rightActor.Deliver(inputMessage{ints: msg.ints[mid:], mexType: inputMessageTypeSort})
	if err != nil {
		fmt.Println("Error sending sort right downstream", err)
		me.Deliver(inputMessage{mexType: inputMessageTypeLeave})
	}

	return sorterState{}, nil
}

func main() {

	leaveCh := make(chan bool)

	sorterURL, _ := url.Parse("actor://sorter")
	sorterActor, err := builders.NewActor(*sorterURL, sortFn(leaveCh), sorterState{})
	if err != nil {
		fmt.Println("Error creating actor:", err)
		return
	}

	defer func() {
		done, _ := sorterActor.Stop()
		<-done
	}()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Type a list of numbers and press enter to sort them.")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	parts := strings.Split(input, " ")
	partsAsIntegers := convertElementsToIntegers(parts)
	fmt.Println("partsAsIntegers", partsAsIntegers)

	sorterActor.Deliver(inputMessage{
		ints:    partsAsIntegers,
		mexType: inputMessageTypeSort,
	})

	<-leaveCh
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
