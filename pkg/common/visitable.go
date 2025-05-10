package common

type Visitable interface {
	// Visit visits the node and applies the given function.
	//
	// Parameters:
	//   - fn (VisitFn): The function to apply to the node.
	Visit(fn VisitFn)
}

// VisitFn is a function type that takes a Visitable as an argument.
type VisitFn func(visitable Visitable)
