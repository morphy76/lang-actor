// > ollama serve
// time=2025-08-04T15:57:04.244+02:00 level=INFO source=routes.go:1291 msg="Listening on 127.0.0.1:11434 (version 0.10.1)"
// > ollama list
// NAME                  ID              SIZE      MODIFIED
// Almawave/Velvet:2B    720611f74c11    4.5 GB    5 months ago

package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"

	"github.com/google/uuid"
	"github.com/morphy76/lang-actor/pkg/builders"

	g "github.com/morphy76/lang-actor/pkg/graph"
)

// graphState holds the streaming text state
type graphState struct {
	content string
}

// MergeChange implements the graph.State interface
func (s *graphState) MergeChange(purpose any, value any) error {
	if v, ok := value.(string); ok {
		s.content = v
	}
	return nil
}

// graphConfig holds configuration for text streaming
type graphConfig struct {
}

func startStateMonitor(graph g.Graph) {
	go func() {
		for state := range graph.StateChangedCh() {
			if graphState, ok := state.(*graphState); ok {
				fmt.Print(graphState.content)
			}
		}
	}()
}

func main() {
	graph, err := builders.NewGraph(&graphState{content: ""}, &graphConfig{})
	if err != nil {
		fmt.Printf("Error creating graph: %v\n", err)
		return
	}

	rootNode, err := builders.NewRootNode(graph)
	if err != nil {
		fmt.Printf("Error creating root node: %v\n", err)
		return
	}

	debugNode, err := builders.NewDebugNode(graph)
	if err != nil {
		fmt.Printf("Error creating child node: %v\n", err)
		return
	}

	endNode, err := builders.NewEndNode(graph)
	if err != nil {
		fmt.Printf("Error creating end node: %v\n", err)
		return
	}

	ollamaURL, err := url.Parse("http://localhost:11434")
	if err != nil {
		fmt.Printf("❌ Error parsing Ollama API URL: %v\n", err)
		return
	}

	ollamaNode, err := builders.NewOllamaNode(graph, ollamaURL)
	if err != nil {
		fmt.Printf("❌ Error creating Ollama node: %v\n", err)
		return
	}

	err = rootNode.OneWayRoute("leavingStart", ollamaNode)
	if err != nil {
		fmt.Printf("Error creating route from root to child: %v\n", err)
		return
	}
	err = ollamaNode.OneWayRoute("leavingOllama", debugNode)
	if err != nil {
		fmt.Printf("Error creating route from root to child: %v\n", err)
		return
	}
	err = debugNode.OneWayRoute("leavingDebug", endNode)
	if err != nil {
		fmt.Printf("Error creating route from child to end: %v\n", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Press a enter to start")

	startStateMonitor(graph)

	reader.ReadString('\n')

	err = rootNode.Accept(uuid.NewString())
	if err != nil {
		fmt.Printf("Error accepting message: %v\n", err)
		return
	}

	fmt.Println("End of the graph")
}
