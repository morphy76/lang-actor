# Design principles of the actor model

## Graph

- [] Create a graph model for node routes so that actors can define processes.

### LangGraph-inspired Actor Graph Model

#### Core Components

1. **ActorNode**: Actors serving as graph nodes with specific responsibilities
   - Processing actors: Transform input and produce output
   - Decision actors: Determine the next actor in the flow
   - Aggregator actors: Combine results from multiple actors

2. **Edge Types**:
   - Static edges: Fixed connections between actors
   - Conditional edges: Paths determined by message content
   - Dynamic edges: Created/destroyed during runtime

3. **Graph Structure**:
   - DAG (Directed Acyclic Graph): For sequential workflows
   - Cyclic Graph: For iterative processes with retry logic

#### Flow Control Patterns

1. **Sequential Flow**: Messages flow linearly from one actor to the next
2. **Conditional Branching**: Based on actor decisions
3. **Parallel Processing**: Multiple actors process simultaneously
4. **Join Patterns**: Synchronizing results from multiple actor paths
5. **Iterative Processing**: Creating cycles with termination conditions

#### Implementation Strategy

1. **Graph Definition**:
   ```go
   type ActorGraph struct {
       Nodes map[url.URL]ActorRef
       Edges map[url.URL][]Edge
   }

   type Edge struct {
       From      url.URL
       To        url.URL
       Condition EdgeCondition // Optional function to determine if edge should be followed
   }

   type EdgeCondition func(Message) bool
   ```

2. **Message Routing**:
   - Use a Router actor that understands the graph structure
   - Each message includes its current position in the graph
   - Router determines next destination based on edge conditions

3. **State Management**:
   - Actor states can be persisted between steps
   - Graph execution state tracks overall progress
   - Support for resuming execution from any point

4. **Execution Models**:
   - Reactive: Pure message-driven execution
   - Orchestrated: Central coordinator manages flow

5. **Composition Patterns**:
   - Subgraphs: Encapsulate complex logic as reusable components
   - Graph templates: Define reusable workflow patterns
