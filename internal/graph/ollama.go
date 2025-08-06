package graph

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/ollama/ollama/api"
	ollamaAPI "github.com/ollama/ollama/api"

	f "github.com/morphy76/lang-actor/pkg/framework"
	g "github.com/morphy76/lang-actor/pkg/graph"
)

type ollamaNode struct {
	node
}

// NewOllamaNode creates a new instance of the Ollama node.
func NewOllamaNode(
	forGraph g.Graph,
	url *url.URL,
) (g.Node, error) {
	baseName := fmt.Sprintf("graph://nodes/llm/ollama/%s_%s", url.Host, url.Port())
	address, err := url.Parse(baseName)
	if err != nil {
		return nil, err
	}

	// 03c7fec5967149cca9a85c6baa41787c.2PT6vp8oZ5GRBSUzysjqGV3l

	ollamaClient := ollamaAPI.NewClient(url, http.DefaultClient)
	taskFn := func(msg f.Message, self f.Actor[g.NodeRef]) (g.NodeRef, error) {

		messages := []ollamaAPI.Message{
			{
				Role:    "system",
				Content: "Provide very brief, concise responses",
			},
			{
				Role:    "user",
				Content: "Name some unusual animals",
			},
			{
				Role:    "assistant",
				Content: "Monotreme, platypus, echidna",
			},
			{
				Role:    "user",
				Content: "which of these is the most dangerous?",
			},
		}

		ctx, cancel := context.WithCancel(context.Background())
		stream := new(bool)
		*stream = false
		req := &api.ChatRequest{
			Model:    "Almawave/Velvet:2B",
			Messages: messages,
			Stream:   stream,
		}

		respFunc := func(resp api.ChatResponse) error {
			self.State().GraphState().MergeChange(nil, resp.Message.Content)
			if resp.Done {
				cancel()
			}
			return nil
		}

		err = ollamaClient.Chat(ctx, req, respFunc)
		if err != nil {
			log.Fatalf("❌ Error calling Ollama API: %v\n", err)
		}

		<-ctx.Done()
		self.State().ProceedOntoRoute() <- g.WhateverOutcome

		return self.State(), nil
	}

	baseNode, err := newNode(forGraph, *address, taskFn)
	if err != nil {
		return nil, err
	}

	return &ollamaNode{
		node: *baseNode,
	}, nil
}
