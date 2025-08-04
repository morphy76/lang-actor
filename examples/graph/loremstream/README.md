# Lorem Ipsum Streaming Graph Example

This example demonstrates a graph-based token streaming implementation using the lang-actor framework. It showcases how to create a custom node that generates lorem ipsum text while streaming it to stdout in real-time and updating the graph state.

## Features

- **Real-time Streaming**: Text is streamed to console as it's generated
- **Two Streaming Modes**:
  - Character-based streaming (2 characters per chunk)
  - Token-based streaming (word by word)
- **Graph State Updates**: Each chunk updates the graph state in real-time
- **Visual Feedback**: Progress indicators and completion statistics

## How It Works

1. **Graph Setup**: Creates a graph with start → custom node → end
2. **Custom Node**: The `LoremGeneratorNode` implements the streaming logic
3. **State Management**: Uses `GraphState` to track streaming progress
4. **Real-time Updates**: Each text chunk is immediately:
   - Printed to stdout for user visibility
   - Added to the graph state via `MergeChange`
5. **Streaming Effect**: Configurable delays between chunks create the streaming illusion

## Architecture

```text
Root Node → Lorem Generator Node → End Node
     ↓              ↓                 ↓
   Start        Generate &         Complete
              Stream Text
```

## Running the Example

```bash
cd examples/graph/loremstream
go run run.go
```

Choose between:

1. **Character-based streaming**: Fast, 2-character chunks with 60ms delay
2. **Token-based streaming**: Word-by-word with 120ms delay

## Key Implementation Details

- **Custom Node**: Uses `builders.NewCustomNode()` with a custom processing function
- **State Interface**: Implements `graph.State` with `MergeChange()` method
- **Configuration**: Uses typed configuration struct for streaming parameters
- **Non-blocking**: Streaming happens within the node's processing function
- **Graph State**: Updates propagate through the framework's state management

This example demonstrates how to build interactive, real-time applications using the lang-actor graph framework while maintaining clean separation of concerns and leveraging the framework's state management capabilities.
