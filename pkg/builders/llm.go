package builders

import (
	"net/url"

	g "github.com/morphy76/lang-actor/internal/graph"
	"github.com/morphy76/lang-actor/pkg/graph"
	"github.com/morphy76/lang-actor/pkg/graph/ollama"
)

// NewOllamaNode creates a new instance of the Ollama node.
//
// Parameters:
//   - forGraph (graph.Graph): The graph to which the Ollama node belongs.
//   - url (*url.URL): The URL of the Ollama API.
//   - kind (ollama.KindWithModel): behavior and model to use
//   - options (...ollama.NodeOption): additional options, see github.com/morphy76/lang-actor/pkg/graph/ollama package
//
// Returns:
//   - (graph.Node): The created Ollama node.
//   - (error): An error if the node creation fails.
func NewOllamaNode(
	forGraph graph.Graph,
	url *url.URL,
	kind ollama.KindWithModel,
	options ...ollama.NodeOption,
) (graph.Node, error) {
	return g.NewOllamaNode(forGraph, url, kind, options...)
}
