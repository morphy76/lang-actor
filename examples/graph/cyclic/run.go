package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/morphy76/lang-actor/examples/graph/cyclic/nodes"
	"github.com/morphy76/lang-actor/pkg/builders"
)

func main() {

	rootNode, err := builders.NewRootNode()
	if err != nil {
		fmt.Printf("Error creating root node: %v\n", err)
		return
	}

	debugNode, err := builders.NewDebugNode()
	if err != nil {
		fmt.Printf("Error creating debug node: %v\n", err)
		return
	}

	counterNode, err := nodes.NewCounterNode()
	if err != nil {
		fmt.Printf("Error creating counter node: %v\n", err)
		return
	}

	endNode, endCh, err := builders.NewEndNode()
	if err != nil {
		fmt.Printf("Error creating end node: %v\n", err)
		return
	}

	err = rootNode.OneWayRoute("leavingStart", counterNode)
	if err != nil {
		fmt.Printf("Error creating route from root to child: %v\n", err)
		return
	}
	err = counterNode.OneWayRoute("iterate", counterNode)
	if err != nil {
		fmt.Printf("Error creating route from child to itself: %v\n", err)
		return
	}
	err = counterNode.OneWayRoute("leavingCounter", debugNode)
	if err != nil {
		fmt.Printf("Error creating route from child to debug node: %v\n", err)
		return
	}
	err = debugNode.OneWayRoute("leavingDebug", endNode)
	if err != nil {
		fmt.Printf("Error creating route from child to end: %v\n", err)
		return
	}

	config := make(map[string]any)

	graphStatus := nodes.GraphStatus{
		Counter: 0,
	}

	graph, err := builders.NewGraph(rootNode, graphStatus, config)
	if err != nil {
		fmt.Printf("Error creating graph: %v\n", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Press a enter to start")

	reader.ReadString('\n')

	err = graph.Accept(uuid.NewString())
	if err != nil {
		fmt.Printf("Error accepting message: %v\n", err)
		return
	}

	<-endCh
	fmt.Println("End of the graph having counter value:", graphStatus.Counter)
}
