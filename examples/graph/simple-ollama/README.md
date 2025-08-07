# Simple Ollama Example

This example demonstrates how to create a conversation graph using the lang-actor framework with Ollama for AI text generation. The graph creates a philosophical dialogue between an AI "teacher" that generates questions and an AI "student" that responds to them.

## Prerequisites

Before running this example, you need to have Ollama running with a compatible model:

```bash
# Start Ollama server
ollama serve

# Check available models
ollama list

# If needed, pull the required model (or use your preferred model)
ollama pull Almawave/Velvet:2B
```

## How It Works

The example creates a graph with the following flow:

1. **Root Node** → **Question Generator** → **Answer Generator** → **Debug Node** → **End Node**

### Graph Components

- **Question Generator**: An Ollama node configured as a "high school teacher of philosophy" that generates random philosophical questions
- **Answer Generator**: An Ollama node configured as a "high school student of philosophy" that responds to the generated questions
- **Debug Node**: Displays the final state for debugging purposes
- **State Monitor**: A goroutine that prints the streaming text in real-time

### Graph State

The `graphState` struct maintains:

- `question`: Current question being generated
- `answer`: Current answer being generated  
- `TotalQuestion`: Complete generated question
- `TotalAnswer`: Complete generated answer

## Template System for UserUtterance

The example showcases the template system used in `WithUserUtterance()`. This system allows dynamic content injection from the graph state into Ollama prompts.

### Template Syntax

Templates use the `{{.attributeName}}` syntax:

```go
ollama.WithUserUtterance("{{.question}}")
```

### Template Processing

1. **Template Detection**: The system checks if the UserUtterance string starts with `{{.` and ends with `}}`
2. **Attribute Extraction**: It extracts the attribute name between the delimiters
3. **State Lookup**: It calls `ReadAttribute(attributeName)` on the current graph state
4. **Content Replacement**: The template is replaced with the actual value from the state

### Example in Code

```go
// This template:
ollama.WithUserUtterance("{{.question}}")

// Gets transformed to the actual question content at runtime:
// "What is the meaning of life in modern society?"
```

### Supported Templates

Any attribute that can be read from your graph state via the `ReadAttribute()` method can be used in templates:

- `{{.question}}` - Current question content
- `{{.answer}}` - Current answer content  
- `{{.totalQuestion}}` - Complete question text
- `{{.totalAnswer}}` - Complete answer text

### Static vs Dynamic Content

You can also use static content without templates:

```go
// Static content - used as-is
ollama.WithUserUtterance("Please respond to this question")

// Dynamic content - replaced with state value
ollama.WithUserUtterance("{{.question}}")
```

## Running the Example

1. Ensure Ollama is running with the `Almawave/Velvet:2B` model (or modify the code to use your preferred model)
2. Run the example:

   ```bash
   go run run.go
   ```

3. Press Enter when prompted to start the conversation
4. Watch as the AI generates a philosophical question and then responds to it

## Configuration Options

The Ollama nodes support various configuration options:

- `WithModel()`: Specify the model to use
- `WithSystem()`: Set the system prompt/persona
- `WithUserUtterance()`: Set user input (supports templates)
- `WithPrompt()`: Set the generation prompt (for generate mode)
- `WithStream()`: Enable streaming response

## Output

The example streams the text in real-time, showing both the question generation and the answer generation as they happen. The final state is displayed by the debug node.

Example output:

```text
Press enter to start
What does it mean to live authentically in today's digital age?
Living authentically in today's digital age means staying true to your core values while navigating the complex landscape of social media and technology...
```

This example demonstrates how to build interactive AI conversations using the lang-actor graph framework with Ollama integration. The specific personas ("high school teacher" and "student of philosophy") and prompts used are for demonstration purposes and can be customized to match your specific use case requirements while leveraging the framework's template system and state management capabilities.

## Disclaimer

This document has been generated with Full Vibes (Github Copilot using Claude 4 Sonnet).

Prompt:

```text
add a readme in the simple ollama example to explain the example itself and in particular the template system used for UserUtternace where {{. }} define boundaries
```
