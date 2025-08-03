// With Full Vibes (Github Copilot using Claude 3.7 Sonnet)
// prompt:
// create a new example replicating the advanced example but using distinct fork and join nodes instead of the forkjoinnode, call it forkgraph and rename the advanced example in forknode

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
func (s *graphState) AppendGraphState(purpose any, value any) error {
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
	taskFn := func(msg f.Message, self f.Actor[g.NodeState]) (g.NodeState, error) {
		fmt.Println("Counter node processing message")

		// Get graph state and config from the node state
		graphState, ok := self.State().GraphState().(*graphState)
		if !ok {
			fmt.Println("Error: Could not cast to graphState")
			self.State().Outcome() <- "error"
			return self.State(), nil
		}

		config, ok := self.State().GraphConfig().(graphConfig)
		if !ok {
			fmt.Println("Error: Could not cast to graphConfig")
			self.State().Outcome() <- "error"
			return self.State(), nil
		}

		// Increment counter and check if we should continue or exit cycle
		graphState.AppendGraphState("count", nil)
		fmt.Printf("Counter iteration: %d/%d\n", graphState.Counter, config.MaxIterations)

		if graphState.Counter < config.MaxIterations {
			// Continue cycling
			self.State().Outcome() <- "iterate"
		} else {
			// Exit the cycle
			fmt.Println("Maximum iterations reached, proceeding to next stage")
			self.State().Outcome() <- "complete"
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

// Creates a processing function for branch nodes after fork
func createProcessingFn(id string) f.ProcessingFn[g.NodeState] {
	return func(msg f.Message, self f.Actor[g.NodeState]) (g.NodeState, error) {
		fmt.Printf("Process '%s' executing\n", id)

		// Simulate work with different durations based on ID length
		processingTime := time.Duration(len(id)*100) * time.Millisecond
		time.Sleep(processingTime)

		result := fmt.Sprintf("completed in %v", processingTime)
		fmt.Printf("Process '%s' %s\n", id, result)

		// Update state with result
		graphState := self.State().GraphState().(*graphState)
		graphState.AppendGraphState(id, result)

		// Signal completion to parent
		self.State().Outcome() <- "done"

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

// createProcessBranchNodes creates all the process branch nodes
func createProcessBranchNodes(graph g.Graph, processNames []string) ([]g.Node, error) {
	branches := make([]g.Node, 0, len(processNames))

	for _, name := range processNames {
		address, err := url.Parse("graph://nodes/process/" + name + "/" + uuid.NewString())
		if err != nil {
			return nil, fmt.Errorf("error creating address for %s: %v", name, err)
		}

		branchNode, err := builders.NewCustomNode(
			graph,
			address,
			createProcessingFn(name),
			false,
		)
		if err != nil {
			return nil, fmt.Errorf("error creating branch node for %s: %v", name, err)
		}

		branches = append(branches, branchNode)
	}

	return branches, nil
}

// connectForkToBranches connects the fork node to all branch nodes
func connectForkToBranches(forkNode g.Node, branches []g.Node) error {
	for i, branch := range branches {
		routeName := fmt.Sprintf("branch%d", i+1)
		err := forkNode.OneWayRoute(routeName, branch)
		if err != nil {
			return fmt.Errorf("error connecting fork to branch %d: %v", i+1, err)
		}
	}
	return nil
}

// connectBranchesToJoin connects all branch nodes to the join node
func connectBranchesToJoin(branches []g.Node, joinNode g.Node) error {
	for i, branch := range branches {
		err := branch.OneWayRoute("done", joinNode)
		if err != nil {
			return fmt.Errorf("error connecting branch %d to join: %v", i+1, err)
		}
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

// setupGraph creates the graph structure and connects the nodes
func setupGraph(processNames []string) (g.Graph, g.Node, error) {
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
		return nil, nil, fmt.Errorf("error creating graph: %v", err)
	}

	// Create nodes
	rootNode, counterNode, forkNode, joinNode, endNode, err := createAllNodes(graph)
	if err != nil {
		return nil, nil, err
	}

	// Connect all nodes together
	err = connectAllNodes(rootNode, counterNode, forkNode, joinNode, endNode, processNames)
	if err != nil {
		return nil, nil, err
	}

	return graph, rootNode, nil
}

// createAllNodes creates all nodes in the graph
func createAllNodes(graph g.Graph) (g.Node, g.Node, g.Node, g.Node, g.Node, error) {
	rootNode, err := builders.NewRootNode(graph)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("error creating root node: %v", err)
	}

	counterNode, err := NewCounterNode(graph)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("error creating counter node: %v", err)
	}

	// Create fork node
	forkNode, err := builders.NewForkNode(graph)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("error creating fork node: %v", err)
	}

	// Create join node that corresponds to the fork
	joinNode, err := builders.NewJoinNode(graph, forkNode)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("error creating join node: %v", err)
	}

	endNode, err := builders.NewEndNode(graph)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("error creating end node: %v", err)
	}

	return rootNode, counterNode, forkNode, joinNode, endNode, nil
}

// connectAllNodes connects the nodes with appropriate routes
func connectAllNodes(rootNode, counterNode, forkNode, joinNode, endNode g.Node, processNames []string) error {
	// Create branch nodes for each process
	graph := rootNode.GetResolver().(g.Graph)
	branches, err := createProcessBranchNodes(graph, processNames)
	if err != nil {
		return err
	}

	// Root to counter
	err = rootNode.OneWayRoute("start", counterNode)
	if err != nil {
		return fmt.Errorf("error creating route from root to counter: %v", err)
	}

	// Cyclic route for counter
	err = counterNode.OneWayRoute("iterate", counterNode)
	if err != nil {
		return fmt.Errorf("error creating cyclic route for counter: %v", err)
	}

	// Counter to fork node when cycle completes
	err = counterNode.OneWayRoute("complete", forkNode)
	if err != nil {
		return fmt.Errorf("error creating route from counter to fork: %v", err)
	}

	// Error route to end
	err = counterNode.OneWayRoute("error", endNode)
	if err != nil {
		return fmt.Errorf("error creating error route: %v", err)
	}

	// Connect fork to all branches
	err = connectForkToBranches(forkNode, branches)
	if err != nil {
		return err
	}

	// Connect all branches to join
	err = connectBranchesToJoin(branches, joinNode)
	if err != nil {
		return err
	}

	// Join to end
	err = joinNode.OneWayRoute("rejoining", endNode)
	if err != nil {
		return fmt.Errorf("error creating route from join to end: %v", err)
	}

	return nil
}

func main() {
	fmt.Println("Starting Fork Graph Example with Separate Fork and Join Nodes")

	// Read process names from user input
	processNames, err := readProcessNamesFromInput()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Using process names: %v\n", processNames)

	// Setup the graph and get the entry point
	graph, rootNode, err := setupGraph(processNames)
	if err != nil {
		fmt.Printf("Graph setup error: %v\n", err)
		return
	}

	// Wait for user input to start
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\nPress Enter to start the graph execution")
	reader.ReadString('\n')

	// Start the graph execution
	msg := uuid.NewString()
	fmt.Printf("Starting graph with message: %s\n\n", msg)

	stateChangesCh := graph.StateChangedCh()
	go func() {
		for state := range stateChangesCh {
			fmt.Printf("State changed: %+v\n", state)
		}
	}()
	err = rootNode.Accept(msg)
	if err != nil {
		fmt.Printf("Error accepting message: %v\n", err)
		return
	}

	// Print results
	printResults(graph)
	fmt.Println("\nEnd of graph execution")
}
