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

		ctx := context.Background()
		req := &api.ChatRequest{
			Model:    "Almawave/Velvet:2B",
			Messages: messages,
		}

		respFunc := func(resp api.ChatResponse) error {
			self.State().GraphState().MergeChange(nil, resp.Message.Content)
			return nil
		}

		err = ollamaClient.Chat(ctx, req, respFunc)
		if err != nil {
			log.Fatal(err)
		}

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
