package graph

import (
	"net/url"
	"sync"

	"github.com/morphy76/lang-actor/internal/routing"

	c "github.com/morphy76/lang-actor/pkg/common"
	g "github.com/morphy76/lang-actor/pkg/graph"
)

var staticGraphConfigAssetion g.GraphConfiguration = (*graphConfig)(nil)
var staticGraphStateAssetion g.GraphState = (*graphState)(nil)

type graphConfig struct {
	cfg map[string]any
}

// Keys returns the keys of the graph configuration.
func (g *graphConfig) Keys() []string {
	keys := make([]string, 0, len(g.cfg))
	for k := range g.cfg {
		keys = append(keys, k)
	}
	return keys
}

// Value retrieves the value associated with the given key in the graph configuration.
func (g *graphConfig) Value(key string) (any, bool) {
	if value, exists := g.cfg[key]; exists {
		return value, true
	}
	return nil, false
}

type graphState struct {
	lock sync.Mutex

	state map[string]any
}

// Set stores a value in the graph state associated with the given key.
func (g *graphState) Set(key string, value any) {
	g.lock.Lock()
	defer g.lock.Unlock()

	if g.state == nil {
		g.state = make(map[string]any)
	}
	g.state[key] = value
}

// Value retrieves the value associated with the given key in the graph state.
func (g *graphState) Value(key string) (any, bool) {
	g.lock.Lock()
	defer g.lock.Unlock()

	if value, exists := g.state[key]; exists {
		return value, true
	}
	return nil, false
}

// Keys retrieves all keys present in the graph state.
func (g *graphState) Keys() []string {
	g.lock.Lock()
	defer g.lock.Unlock()

	keys := make([]string, 0, len(g.state))
	for k := range g.state {
		keys = append(keys, k)
	}
	return keys
}

// NewGraph creates a new instance of the actor graph.
func NewGraph(
	graphName string,
	initialState map[string]any,
	configs map[string]any,
) (g.Graph, error) {

	graphURL, err := url.Parse("graph://" + graphName)
	if err != nil {
		return nil, err
	}

	useCfg := &graphConfig{
		cfg: make(map[string]any, len(configs)),
	}
	for k, v := range configs {
		useCfg.cfg[k] = v
	}

	graph := &graph{
		resolvables: make(map[url.URL]*c.Addressable),
		graphURL:    *graphURL,
		config:      useCfg,
		status: &graphState{
			lock: sync.Mutex{},

			state: make(map[string]any, len(initialState)),
		},
		addressBook: routing.NewAddressBook(),
	}

	for k, v := range initialState {
		graph.status.Set(k, v)
	}

	return graph, nil
}
