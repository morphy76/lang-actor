package graph

import (
	"net/url"
	"sync"

	f "github.com/morphy76/lang-actor/pkg/framework"
)

func newNode[T any](task f.Actor[T], address url.URL) *node {
	return &node{
		lock:    &sync.Mutex{},
		edges:   make(map[string]edge, 0),
		actor:   task,
		address: address,
	}
}
