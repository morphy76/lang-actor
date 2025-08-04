// With Full Vibes (Github Copilot using Claude 3.7 Sonnet)
// prompt:
// create a graph example with the following rules:
// - does not use internal packages
// - use cyclic capabilities
// - use forkjoin
// - replicate the client pattern of the existing simple example
// - read processNames from the user input as a comma separated list of strings

package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/morphy76/lang-actor/pkg/builders"
	f "github.com/morphy76/lang-actor/pkg/framework"
	g "github.com/morphy76/lang-actor/pkg/graph"
)

// Define our graph state with counter for cycles and results map
type graphState struct {
	Counter   int
	Results   map[string]string
	StartTime time.Time
}

// Implement the required graph.State interface method
func (s *graphState) MergeChange(purpose any, value any) error {
	if purpose == "count" {
		s.Counter++
	} else if value != nil {
		key, ok := purpose.(string)
		if ok {
			strValue, ok := value.(string)
			if ok {
				if s.Results == nil {
					s.Results = make(map[string]string)
				}
				s.Results[key] = strValue
			}
		}
	}
	return nil
}

// Define our graph configuration with iteration count and process names
type graphConfig struct {
	MaxIterations int
	ProcessNames  []string
}

// NewCounterNode demonstrates cyclic capabilities
func NewCounterNode(forGraph g.Graph) (g.Node, error) {
	address, err := url.Parse("graph://nodes/counter/" + uuid.NewString())
	if err != nil {
		return nil, err
	}

	// Define the processing function
	taskFn := func(msg f.Message, self f.Actor[g.NodeRef]) (g.NodeRef, error) {
		fmt.Println("Counter node processing message")

		// Get graph state and config from the node state
		graphState, ok := self.State().GraphState().(*graphState)
		if !ok {
			fmt.Println("Error: Could not cast to graphState")
			self.State().ProceedOntoRoute() <- "error"
			return self.State(), nil
		}

		config, ok := self.State().GraphConfig().(graphConfig)
		if !ok {
			fmt.Println("Error: Could not cast to graphConfig")
			self.State().ProceedOntoRoute() <- "error"
			return self.State(), nil
		}

		// Increment counter and check if we should continue or exit cycle
		graphState.MergeChange("count", nil)
		fmt.Printf("Counter iteration: %d/%d\n", graphState.Counter, config.MaxIterations)

		if graphState.Counter < config.MaxIterations {
			// Continue cycling
			self.State().ProceedOntoRoute() <- "iterate"
		} else {
			// Exit the cycle
			fmt.Println("Maximum iterations reached, proceeding to next stage")
			self.State().ProceedOntoRoute() <- "complete"
		}

		return self.State(), nil
	}

	// Create custom node
	return builders.NewCustomNode(
		forGraph,
		address,
		taskFn,
		false,
	)
}

// Creates a processing function for fork-join child nodes
func createProcessingFn(id string) f.ProcessingFn[g.NodeRef] {
	return func(msg f.Message, self f.Actor[g.NodeRef]) (g.NodeRef, error) {
		fmt.Printf("Process '%s' executing\n", id)

		// Simulate work with different durations based on ID length
		processingTime := time.Duration(len(id)*100) * time.Millisecond
		time.Sleep(processingTime)

		result := fmt.Sprintf("completed in %v", processingTime)
		fmt.Printf("Process '%s' %s\n", id, result)

		// Update state with result
		graphState := self.State().GraphState().(*graphState)
		graphState.MergeChange(id, result)

		// Signal completion to parent
		self.State().ProceedOntoRoute() <- fmt.Sprintf("%s-done", id)

		return self.State(), nil
	}
}

// readProcessNamesFromInput prompts for and parses comma-separated process names
func readProcessNamesFromInput() ([]string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter process names as a comma-separated list (e.g. alpha,beta,gamma,delta): ")
	input, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	// Clean input and split by comma
	input = input[:len(input)-1] // Remove newline character
	if len(input) == 0 {
		fmt.Println("Using default process names")
		input = "alpha,beta,gamma,delta"
	}

	// Split into process names
	processNames := []string{}
	for _, name := range strings.Split(input, ",") {
		trimmedName := strings.TrimSpace(name)
		if trimmedName != "" {
			processNames = append(processNames, trimmedName)
		}
	}

	if len(processNames) == 0 {
		return nil, fmt.Errorf("no valid process names provided")
	}

	return processNames, nil
}

// createGraphNodes creates all nodes in the graph
func createGraphNodes(graph g.Graph, processNames []string) (g.Node, g.Node, g.Node, g.Node, error) {
	rootNode, err := builders.NewRootNode(graph)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("error creating root node: %v", err)
	}

	counterNode, err := NewCounterNode(graph)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("error creating counter node: %v", err)
	}

	// Create processing functions for the fork-join node
	var processingFuncs []f.ProcessingFn[g.NodeRef]
	for _, name := range processNames {
		processingFuncs = append(processingFuncs, createProcessingFn(name))
	}

	forkJoinNode, err := builders.NewForkJoingNode(
		graph,
		false,
		processingFuncs...,
	)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("error creating fork-join node: %v", err)
	}

	endNode, err := builders.NewEndNode(graph)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("error creating end node: %v", err)
	}

	return rootNode, counterNode, forkJoinNode, endNode, nil
}

// connectNodes sets up the routes between nodes
func connectNodes(rootNode, counterNode, forkJoinNode, endNode g.Node) error {
	// Root to counter
	err := rootNode.OneWayRoute("start", counterNode)
	if err != nil {
		return fmt.Errorf("error creating route from root to counter: %v", err)
	}

	// Cyclic route for counter
	err = counterNode.OneWayRoute("iterate", counterNode)
	if err != nil {
		return fmt.Errorf("error creating cyclic route for counter: %v", err)
	}

	// Counter to fork-join when cycle completes
	err = counterNode.OneWayRoute("complete", forkJoinNode)
	if err != nil {
		return fmt.Errorf("error creating route from counter to fork-join: %v", err)
	}

	// Error route to end
	err = counterNode.OneWayRoute("error", endNode)
	if err != nil {
		return fmt.Errorf("error creating error route: %v", err)
	}

	// Fork-join to end
	err = forkJoinNode.OneWayRoute("", endNode)
	if err != nil {
		return fmt.Errorf("error creating route from fork-join to end: %v", err)
	}

	return nil
}

// printResults displays the final state and results
func printResults(graph g.Graph) {
	finalState, ok := graph.State().(*graphState)
	if ok {
		fmt.Println("\nGraph execution completed")
		fmt.Printf("Iterations: %d\n", finalState.Counter)
		fmt.Printf("Total execution time: %v\n", time.Since(finalState.StartTime))

		fmt.Println("\nProcess results:")
		for k, v := range finalState.Results {
			fmt.Printf("  %s: %v\n", k, v)
		}
	}
}

func main() {
	fmt.Println("Starting Advanced Graph Example with Cyclic and Fork-Join Capabilities")

	// Read process names from user input
	processNames, err := readProcessNamesFromInput()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Using process names: %v\n", processNames)

	// Set up configuration and initial state
	config := graphConfig{
		MaxIterations: 3,
		ProcessNames:  processNames,
	}

	initialState := &graphState{
		Counter:   0,
		Results:   make(map[string]string),
		StartTime: time.Now(),
	}

	// Create graph
	graph, err := builders.NewGraph(initialState, config)
	if err != nil {
		fmt.Printf("Error creating graph: %v\n", err)
		return
	}

	// Create all nodes
	rootNode, counterNode, forkJoinNode, endNode, err := createGraphNodes(graph, processNames)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	// Connect the nodes
	err = connectNodes(rootNode, counterNode, forkJoinNode, endNode)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	// Wait for user input to start
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\nPress Enter to start the graph execution")
	reader.ReadString('\n')

	// Start the graph execution
	msg := uuid.NewString()
	fmt.Printf("Starting graph with message: %s\n\n", msg)

	err = rootNode.Accept(msg)
	if err != nil {
		fmt.Printf("Error accepting message: %v\n", err)
		return
	}

	// Print results
	printResults(graph)
	fmt.Println("\nEnd of graph execution")
}
