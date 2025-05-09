package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/morphy76/lang-actor/pkg/builders"
)

func main() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Press a enter to start")

	reader.ReadString('\n')

	rootNode, err := builders.NewRootNode()
	if err != nil {
		fmt.Printf("Error creating root node: %v\n", err)
		return
	}

	childNode, err := builders.NewDebugNode()
	if err != nil {
		fmt.Printf("Error creating child node: %v\n", err)
		return
	}

	endNode, endCh, err := builders.NewEndNode()
	if err != nil {
		fmt.Printf("Error creating end node: %v\n", err)
		return
	}

	err = rootNode.OneWayRoute("leavingStart", childNode)
	if err != nil {
		fmt.Printf("Error creating route from root to child: %v\n", err)
		return
	}
	err = childNode.OneWayRoute("leavingDebug", endNode)
	if err != nil {
		fmt.Printf("Error creating route from child to end: %v\n", err)
		return
	}

	config := make(map[string]any)
	config["test"] = uuid.NewString()
	config["test2"] = uuid.NewString()
	config["test3"] = uuid.NewString()

	graph, err := builders.NewGraph(rootNode, config)
	if err != nil {
		fmt.Printf("Error creating graph: %v\n", err)
		return
	}

	graph.Accept(uuid.NewString())

	<-endCh
	fmt.Println("End of the graph")
}
