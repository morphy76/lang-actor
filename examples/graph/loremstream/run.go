// With Full Vibes (Github Copilot using Claude 4 Sonnet)
// prompt:
// add a graph example that, using a custom node between a graph start and a graph end has the following behavior:
//
// - the custom node generate a lorem ipsjm text
// - while generating it updates the graph state
// - the graph state is streamed to stdio output
//
// the final feeling when running the graph is that the end user read the lorem ipsum string streamed to console while the cutsom node generate parts of it
//
// sort of token streaming

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

// Lorem ipsum text to stream
const loremIpsum = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium, totam rem aperiam, eaque ipsa quae ab illo inventore veritatis et quasi architecto beatae vitae dicta sunt explicabo. Nemo enim ipsam voluptatem quia voluptas sit aspernatur aut odit aut fugit, sed quia consequuntur magni dolores eos qui ratione voluptatem sequi nesciunt."

// graphState holds the streaming text state
type graphState struct {
	StreamedText []string
	TotalText    string
	IsComplete   bool
}

// MergeChange implements the graph.State interface
func (s *graphState) MergeChange(purpose any, value any) error {
	switch purpose.(string) {
	case "chunk":
		if chunk, ok := value.(string); ok {
			s.StreamedText = append(s.StreamedText, chunk)
		}
	case "complete":
		s.IsComplete = true
		s.TotalText = strings.Join(s.StreamedText, "")
		fmt.Printf("\n\n📜 Total Lorem Ipsum Text:\n%s\n", s.TotalText)
	case "start":
		fmt.Print("\n🔄 Starting lorem ipsum generation...\n\n📝 ")
	}
	return nil
}

func (s *graphState) ReadAttribute(name string) any {
	switch name {
	case "StreamedText":
		return s.StreamedText
	case "TotalText":
		return s.TotalText
	case "IsComplete":
		return s.IsComplete
	default:
		return nil

	}
}

// Unwrap implements the graph.State interface
func (s *graphState) Unwrap() g.State {
	return s
}

// GraphConfig holds configuration for text streaming
type GraphConfig struct {
	ChunkSize   int
	StreamDelay time.Duration
	TokenMode   bool
}

// NewLoremGeneratorNode creates a custom node that generates lorem ipsum text with streaming
func NewLoremGeneratorNode(forGraph g.Graph) (g.Node, error) {
	address, err := url.Parse("graph://nodes/lorem-generator/" + uuid.NewString())
	if err != nil {
		return nil, err
	}

	// Define the processing function for streaming lorem ipsum
	taskFn := func(msg f.Message, self f.Actor[g.NodeRef]) (g.NodeRef, error) {
		// Get config
		config, ok := self.State().GraphConfig().(GraphConfig)
		if !ok {
			fmt.Printf("❌ Error: Could not cast to GraphConfig, got type: %T\n", self.State().GraphConfig())
			self.State().ProceedOntoRoute() <- "error"
			return self.State(), nil
		}

		// Signal start of streaming
		self.State().GraphState().MergeChange("start", nil)

		// Split text into chunks based on configuration
		text := loremIpsum
		var chunks []string

		if config.TokenMode {
			// Token-based streaming (word by word)
			words := strings.Fields(text)
			for _, word := range words {
				chunks = append(chunks, word+" ")
			}
		} else {
			// Character-based streaming
			for i := 0; i < len(text); i += config.ChunkSize {
				end := i + config.ChunkSize
				if end > len(text) {
					end = len(text)
				}
				chunks = append(chunks, text[i:end])
			}
		}

		// Stream each chunk with delay
		for i, chunk := range chunks {
			self.State().GraphState().MergeChange("chunk", chunk)

			// Add delay between chunks to simulate streaming
			if i < len(chunks)-1 { // Don't delay after the last chunk
				time.Sleep(config.StreamDelay)
			}
		}

		// Mark as complete
		self.State().GraphState().MergeChange("complete", true)

		fmt.Printf("\n\n✅ Lorem ipsum generation completed!\n")
		fmt.Printf("📊 Total chunks streamed: %d\n", len(chunks))
		if !config.TokenMode {
			fmt.Printf("⚡ Characters per chunk: %d\n", config.ChunkSize)
		}
		fmt.Printf("⏱️  Stream delay: %v\n", config.StreamDelay)

		// Proceed to next node
		self.State().ProceedOntoRoute() <- "generated"

		return self.State(), nil
	}

	// Create custom node
	return builders.NewCustomNode(
		forGraph,
		address,
		taskFn,
	)
}

// startStateMonitor starts a goroutine to monitor state changes
func startStateMonitor(graph g.Graph) {
	go func() {
		for state := range graph.StateChangedCh() {
			if graphState, ok := state.(*graphState); ok {
				if len(graphState.StreamedText) > 0 {
					fmt.Print(graphState.StreamedText[len(graphState.StreamedText)-1])
				}
			}
		}
	}()
}

func main() {
	fmt.Println("🎭 Lorem Ipsum Streaming Graph Example")
	fmt.Println("=====================================")

	// Get user preferences
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Choose streaming mode:\n1. Character-based streaming (default)\n2. Token-based streaming (word by word)\nEnter choice (1 or 2): ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	var config GraphConfig
	if choice == "2" {
		config = GraphConfig{
			TokenMode:   true,
			StreamDelay: 120 * time.Millisecond, // Delay between words
		}
		fmt.Println("📝 Token-based streaming selected (word by word)")
	} else {
		config = GraphConfig{
			ChunkSize:   2,                     // Characters per chunk
			StreamDelay: 60 * time.Millisecond, // Delay between chunks
			TokenMode:   false,
		}
		fmt.Println("📝 Character-based streaming selected")
	}

	// Create initial state
	initialState := &graphState{
		StreamedText: []string{},
		TotalText:    "",
		IsComplete:   false,
	}

	// Create graph
	graph, err := builders.NewGraph(initialState, config)
	if err != nil {
		fmt.Printf("❌ Error creating graph: %v\n", err)
		return
	}

	// Start state monitoring
	startStateMonitor(graph)

	// Create nodes
	rootNode, err := builders.NewRootNode(graph)
	if err != nil {
		fmt.Printf("❌ Error creating root node: %v\n", err)
		return
	}

	loremNode, err := NewLoremGeneratorNode(graph)
	if err != nil {
		fmt.Printf("❌ Error creating lorem generator node: %v\n", err)
		return
	}

	endNode, err := builders.NewEndNode(graph)
	if err != nil {
		fmt.Printf("❌ Error creating end node: %v\n", err)
		return
	}

	// Connect nodes with routes
	err = rootNode.OneWayRoute("start", loremNode)
	if err != nil {
		fmt.Printf("❌ Error creating route from root to lorem generator: %v\n", err)
		return
	}

	err = loremNode.OneWayRoute("generated", endNode)
	if err != nil {
		fmt.Printf("❌ Error creating route from lorem generator to end: %v\n", err)
		return
	}

	// Wait for user to start the process
	fmt.Print("\nPress Enter to start streaming lorem ipsum text...")
	reader.ReadString('\n')

	// Start the graph execution
	err = rootNode.Accept(uuid.NewString())
	if err != nil {
		fmt.Printf("❌ Error starting graph execution: %v\n", err)
		return
	}

	// Wait a bit for processing to complete
	time.Sleep(2 * time.Second)

	fmt.Println("\n🎉 Graph execution completed!")
	fmt.Println("=====================================")
}
