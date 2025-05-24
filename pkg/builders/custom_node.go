package builders

import (
	"net/url"

	"github.com/google/uuid"

	f "github.com/morphy76/lang-actor/internal/framework"
	g "github.com/morphy76/lang-actor/internal/graph"
	"github.com/morphy76/lang-actor/pkg/framework"
	"github.com/morphy76/lang-actor/pkg/graph"
)

// NewCustomNode creates a new instance of a custom node.
//
// Type parameters:
//   - T: The type of the node state.
//
// Parameters:
//   - address (*url.URL): The URL address of the node.
//   - taskFn (framework.ProcessingFn[T]): The processing function for the node.
//   - nodeState (T): The initial state of the node.
//   - transient (bool): Whether the node is transient or not.
//
// Returns:
//   - (graph.Node): The created custom node.
//   - (error): An error if the node creation fails.
func NewCustomNode[T any](
	address *url.URL,
	taskFn framework.ProcessingFn[T],
	nodeState T,
	transient bool,
) (graph.Node, error) {

	actorAddress, err := url.Parse("actor://" + address.Host + address.Path + "/" + uuid.NewString())
	if err != nil {
		return nil, err
	}

	customTask, err := f.NewActor(
		*actorAddress,
		taskFn,
		nodeState,
		transient,
	)
	if err != nil {
		return nil, err
	}

	baseNode := g.NewNodeWithActor(*address, customTask)

	return baseNode, nil
}
