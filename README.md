# lang-actor

Lang-Actor is a lightweight, Go-based implementation of the Actor Model, a computational model in which "actors" serve as the universal primitives of concurrent computation. This framework provides a robust foundation for building concurrent, message-driven applications with clearly defined boundaries and communication patterns.

## Description

The Actor Model in Lang-Actor follows these core principles:

- Each actor has a unique address (URI);
- Actors communicate exclusively through asynchronous message passing;
- Actors maintain private state that can only be modified by processing messages;
- Actors can create child actors, forming hierarchical supervision trees;
- Each actor processes messages sequentially from its mailbox.

## Main Capabilities

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
   - State mutation controlled through message processing

5. **Lifecycle Management**:
   - Actors can be started, stopped, and monitored
   - Graceful shutdown with message draining

## Simple Usage Example

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

func (m incrementMessage) Mutation() bool {
 return true
}

// Define message processing function
var counterFn framework.ProcessingFn[counterState] = func(
 msg framework.Message,
 self framework.Actor[counterState],
) (counterState, error) {
 if incMsg, ok := msg.(incrementMessage); ok {
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
  counterActor.Deliver(msg)
  time.Sleep(500 * time.Millisecond)
 }
}
```

This simple example creates a counter actor that processes increment messages and keeps track of a running total.

## More Complex Examples

For more advanced usage patterns, refer to the examples in the repository:

1. **Echo Actor** (`examples/actors/echo/run.go`): Demonstrates basic message handling by echoing received messages.

2. **Ping-Pong** (`examples/actors/pingpong/run.go`): Shows communication between multiple actors.

3. **Echo With Child** (`examples/actors/echowithchild/run.go`): Demonstrates parent-child actor relationships.

4. **Calculator** (`examples/actors/calculator/run.go`): Implements a simple calculator using actors.

5. **Self-Ping-Pong** (`examples/actors/selfpingpong/run.go`): Shows how actors can send messages to themselves.

6. **Sort** (`examples/actors/sort/run.go`): Demonstrates more complex state management and processing.

These examples demonstrate various aspects of the framework including actor creation, message passing, state management, and actor hierarchies.

## Design Principles

The framework follows these key design principles:

- Type safety through Go's generics
- Clear separation of concerns between actors
- Message-driven communication
- Hierarchical organization of actors
- Flexible backpressure policies for mailboxes
- Graceful handling of actor lifecycle

## Disclaimer

This document has been generated with Full Vibes (Github Copilot using Claude 3.7 Sonnet).

Prompt:

```text
Document the actor model framework with:
- description
- main capabilities
- simple usage example
- reference to the examples for more complex cases
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please read the [CONTRIBUTING.md](CONTRIBUTING.md) file for details on how to contribute to this project.
