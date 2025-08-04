# lang-actor

Lang-Actor is a ligh5. **Type-Safe Message Processing**:
   - Generic typed actors and message handlers
   - State always updated through message processingight, Go-based implementation of the Actor Model and Graph Model, providing robust foundations for building concurrent and flow-based applications with clearly defined boundaries and communication patterns.

## Actor Model Framework

### Description

The Actor Model in Lang-Actor follows these core principles:

- Each actor has a unique address (URI);
- Actors communicate exclusively through asynchronous message passing;
- Actors maintain private state that is modified by processing messages;
- Actors can create child actors, forming hierarchical supervision trees;
- Each actor processes messages sequentially from its mailbox.

### Main Capabilities

1. **Hierarchical Actor System**:
   - Actors can create and manage child actors
   - Parent-child relationships for structured supervision

2. **Flexible Message Routing**:
   - Unique URI-based addressing scheme
   - Support for local ("actor://") communication with potential for extending to other protocols

3. **Configurable Mailboxes**:
   - Multiple backpressure policies:
     - Block: Wait when mailbox is full
     - Fail: Immediately fail when mailbox is full
     - Unbounded: No capacity limit
     - DropNewest: Reject new messages when full
     - DropOldest: Discard oldest messages to make room

4. **Type-Safe Message Processing**:
   - Generic typed actors and message handlers
   - State always updated through message processing

5. **Lifecycle Management**:
   - Actors can be started, stopped, and monitored
   - Graceful shutdown with message draining

### Simple Usage Example

Here's a minimal example of how to create and use an actor:

```go
package main

import (
 "fmt"
 "net/url"
 "time"

 "github.com/morphy76/lang-actor/pkg/builders"
 "github.com/morphy76/lang-actor/pkg/framework"
)

// Define actor state
type counterState struct {
 count int
}

// Define a message type
type incrementMessage struct {
 senderURL url.URL
 amount    int
}

func (m incrementMessage) Sender() url.URL {
 return m.senderURL
}

// Define message processing function
var counterFn framework.ProcessingFn[counterState] = func(
 msg framework.Message,
 self framework.Actor[counterState],
) (counterState, error) {
 if incMsg, ok := msg.Payload().(incrementMessage); ok {
  newCount := self.State().count + incMsg.amount
  fmt.Printf("Counter incremented to: %d\n", newCount)
  return counterState{count: newCount}, nil
 }
 return self.State(), nil
}

func main() {
 // Create actor
 actorURL, _ := url.Parse("actor://counter")
 counterActor, err := builders.NewActor(*actorURL, counterFn, counterState{count: 0})
 if err != nil {
  fmt.Println("Error creating actor:", err)
  return
 }

 // Ensure actor stops when program ends
 defer func() {
  done, _ := counterActor.Stop()
  <-done
 }()

 // Send messages to the actor
 for i := 1; i <= 5; i++ {
  msg := incrementMessage{
   senderURL: *actorURL,
   amount:    i,
  }
  counterActor.Deliver(msg, nil)
  time.Sleep(500 * time.Millisecond)
 }
}
```

This simple example creates a counter actor that processes increment messages and keeps track of a running total.

### More Complex Examples

For more advanced usage patterns, refer to the examples in the repository:

1. **Counter** (`examples/actors/counter/run.go`): The simple counter example shown above.

2. **Echo Actor** (`examples/actors/echo/run.go`): Demonstrates basic message handling by echoing received messages.

3. **Ping-Pong** (`examples/actors/pingpong/run.go`): Shows communication between multiple actors.

4. **Echo With Child** (`examples/actors/echowithchild/run.go`): Demonstrates parent-child actor relationships.

5. **Calculator** (`examples/actors/calculator/run.go`): Implements a simple calculator using actors.

6. **Self-Ping-Pong** (`examples/actors/selfpingpong/run.go`): Shows how actors can send messages to themselves.

7. **Sort** (`examples/actors/sort/run.go`): Demonstrates more complex state management and processing.

These examples demonstrate various aspects of the framework including actor creation, message passing, state management, and actor hierarchies.

## Graph Model Framework

### Description

The Graph Model in Lang-Actor provides a powerful foundation for building flow-based computational graphs. It represents computation as a directed graph where:

- Nodes represent processing units with specific responsibilities
- Edges define the flow of messages between nodes
- Each node can process messages and route them to connected nodes
- The graph maintains shared state accessible to all nodes

This model is particularly useful for workflows, data processing pipelines, and complex business processes with branching logic.

### Main Capabilities

1. **Node Types and Hierarchies**:
   - Root nodes as entry points for external messages
   - Debug nodes for monitoring and logging
   - End nodes as terminal points for flows
   - Custom nodes for specialized processing
   - Fork and Join nodes for parallel processing

2. **Flow Control Patterns**:
   - Sequential message passing
   - Conditional routing based on message content
   - Cyclic flows for iterative processing
   - Fork-join patterns for parallel execution

3. **State Management**:
   - Shared graph-wide state
   - Node-specific attributes
   - Type-safe state updates through message processing

4. **Flexibility and Extensibility**:
   - URI-based addressing scheme similar to the actor model
   - Custom node types through composition
   - Support for transient and persistent nodes

5. **Execution Control**:
   - Deterministic message flow
   - Error handling and recovery
   - Graceful termination

### Simple Usage Example

Here's a minimal example of how to create and use a graph:

```go
package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/morphy76/lang-actor/pkg/builders"
)

// Define graph state
type graphState struct {
	stateAsMap map[string]any
}

func (s graphState) AppendGraphState(purpose any, value any) error {
	return nil
}

func main() {
	// Create graph configuration and initial state
	config := make(map[string]any)
	config["startTime"] = time.Now()
	
	initialStateMap := make(map[string]any)
	initialStateMap["processedCount"] = 0
	
	initialState := graphState{
		stateAsMap: initialStateMap,
	}

	// Create graph
	graph, err := builders.NewGraph(initialState, config)
	if err != nil {
		fmt.Printf("Error creating graph: %v\n", err)
		return
	}

	// Create nodes
	rootNode, err := builders.NewRootNode(graph)
	if err != nil {
		fmt.Printf("Error creating root node: %v\n", err)
		return
	}

	processingNode, err := builders.NewDebugNode(graph)
	if err != nil {
		fmt.Printf("Error creating processing node: %v\n", err)
		return
	}

	endNode, err := builders.NewEndNode(graph)
	if err != nil {
		fmt.Printf("Error creating end node: %v\n", err)
		return
	}

	// Connect nodes with routes
	err = rootNode.OneWayRoute("start", processingNode)
	if err != nil {
		fmt.Printf("Error creating route: %v\n", err)
		return
	}
	
	err = processingNode.OneWayRoute("complete", endNode)
	if err != nil {
		fmt.Printf("Error creating route: %v\n", err)
		return
	}

	// Run the graph with a message
	fmt.Println("Starting graph execution...")
	err = rootNode.Accept(uuid.NewString())
	if err != nil {
		fmt.Printf("Error accepting message: %v\n", err)
		return
	}

	fmt.Println("Graph execution complete")
}
```

This simple example creates a graph with three nodes (root, processing, and end) and connects them in sequence.

### More Complex Examples

For more advanced graph usage patterns, refer to the examples in the repository:

1. **Simple Graph** (`examples/graph/simple/run.go`): Demonstrates basic graph creation with sequential message flow through multiple nodes.

2. **Fork Node** (`examples/graph/forknode/run.go`): Shows how to use a fork-join node for parallel processing of tasks with a combined join operation to synchronize results.

3. **Fork Graph** (`examples/graph/forkgraph/run.go`): Implements parallel processing using separate fork and join nodes, providing more control over the parallel execution pattern.

4. **Advanced Graph** (`examples/graph/advanced/run.go`): Demonstrates complex graph structures with conditional routing, cyclic flows, and rich state management.

These examples showcase the flexibility of the graph model for creating sophisticated computational flows with branching, parallel execution, and state management.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please read the [CONTRIBUTING.md](CONTRIBUTING.md) file for details on how to contribute to this project.
