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

	config := make(map[string]any)
	config["testCfg1"] = uuid.NewString()
	config["testCfg2"] = uuid.NewString()
	config["testCfg3"] = uuid.NewString()
	config["testCfg4"] = time.Now()
	whateverURL, _ := url.Parse("https://example.com:8080/ctx?id=1234")
	config["testCfg5"] = whateverURL

	initialState := make(map[string]any)
	initialState["testState1"] = uuid.NewString()
	initialState["testState2"] = uuid.NewString()
	initialState["testState3"] = uuid.NewString()
	initialState["testState4"] = time.Now()
	whateverURL, _ = url.Parse("https://example.com:8080/ctx?id=1234")
	initialState["testState5"] = whateverURL

	graph, err := builders.NewGraph(config)
	if err != nil {
		fmt.Printf("Error creating graph: %v\n", err)
		return
	}

	rootNode, err := builders.NewRootNode(graph)
	if err != nil {
		fmt.Printf("Error creating root node: %v\n", err)
		return
	}

	childNode, err := builders.NewDebugNode(graph)
	if err != nil {
		fmt.Printf("Error creating child node: %v\n", err)
		return
	}

	endNode, err := builders.NewEndNode(graph)
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

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Press a enter to start")

	reader.ReadString('\n')

	err = rootNode.Accept(uuid.NewString())
	if err != nil {
		fmt.Printf("Error accepting message: %v\n", err)
		return
	}

	fmt.Println("End of the graph")
}
