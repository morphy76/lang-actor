package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/morphy76/lang-actor/pkg/builders"
)

func main() {

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
	config["test4"] = time.Now()
	whateverURL, _ := url.Parse("https://example.com:8080/ctx?id=1234")
	config["test5"] = whateverURL

	graph, err := builders.NewGraph(rootNode, "", config)
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
	fmt.Println("End of the graph")
}
