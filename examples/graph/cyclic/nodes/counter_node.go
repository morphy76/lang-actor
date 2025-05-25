package nodes

import (
	"net/url"

	"github.com/google/uuid"
	"github.com/morphy76/lang-actor/pkg/builders"
	f "github.com/morphy76/lang-actor/pkg/framework"
	g "github.com/morphy76/lang-actor/pkg/graph"
)

type GraphStatus struct {
	Counter int
}

func NewCounterNode(forGraph g.Graph) (g.Node, error) {

	address, err := url.Parse("graph://nodes/counter/" + uuid.NewString())
	if err != nil {
		return nil, err
	}

	taskFn := func(msg f.Message, self f.Actor[g.NodeState]) (g.NodeState, error) {

		// statusResponse, okStatus := msg.(*g.StatusMessage)

		// if okStatus {
		// 	curVal := statusResponse.Value.(*GraphStatus).Counter
		// 	if curVal < 10 {
		// 		// g.NewStatusMessageUpdate(self.Address(), &GraphStatus{
		// 		// 	Counter: curVal + 1,
		// 		// }).Deliver(statusResponse.Sender())
		// 		// route to 'iterate'

		// 		// err = useDebugNode.ProceedOnAnyRoute(self.State().originalMessage)
		// 		// if err != nil {
		// 		// 	return self.State(), err
		// 		// }
		// 	} else {
		// 		// route to 'leavingCounter'

		// 		// err = useDebugNode.ProceedOnAnyRoute(self.State().originalMessage)
		// 		// if err != nil {
		// 		// 	return self.State(), err
		// 		// }
		// 	}
		// } else {
		// 	// requestStatus := g.NewStatusMessageRequest(self.Address())
		// 	// statusNodes := useDebugNode.GetResolver().Query("graph", "nodes", "status")
		// 	// if len(statusNodes) == 0 {
		// 	// 	return self.State(), errors.Join(g.ErrorInvalidRouting, fmt.Errorf("no status node found"))
		// 	// }
		// 	// statusNodes[0].Deliver(&requestStatus)
		// 	return nodeState{
		// 		originalMessage: msg,
		// 	}, nil
		// }

		return self.State(), nil
	}

	return builders.NewCustomNode(
		forGraph,
		address,
		taskFn,
		true,
	)
}
