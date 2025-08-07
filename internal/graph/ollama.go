package graph

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
	ollamaAPI "github.com/ollama/ollama/api"

	f "github.com/morphy76/lang-actor/pkg/framework"
	g "github.com/morphy76/lang-actor/pkg/graph"
	"github.com/morphy76/lang-actor/pkg/graph/ollama"
)

type ollamaNode struct {
	node
}

// NewOllamaNode creates a new instance of the Ollama node.
func NewOllamaNode(
	forGraph g.Graph,
	url *url.URL,
	kind ollama.KindWithModel,
	options ...ollama.NodeOption,
) (g.Node, error) {
	baseName := fmt.Sprintf("graph://nodes/llm/ollama/%s/%s_%s/%s", kind.Kind, url.Host, url.Port(), uuid.NewString())
	address, err := url.Parse(baseName)
	if err != nil {
		return nil, err
	}

	// Merge all options into a single NodeOption
	var mergedOption ollama.NodeOption
	if len(options) > 0 {
		var err error
		mergedOption, err = ollama.MergeNodeOptions(options...)
		if err != nil {
			return nil, fmt.Errorf("failed to merge node options: %w", err)
		}
	}

	ollamaClient := ollamaAPI.NewClient(url, http.DefaultClient)
	var taskFn f.ProcessingFn[g.NodeRef]

	switch kind.Kind {
	case ollama.Generate:
		taskFn = func(msg f.Message, self f.Actor[g.NodeRef]) (g.NodeRef, error) {

			ctx, cancel := context.WithCancel(context.Background())

			req := &ollamaAPI.GenerateRequest{
				Model:  kind.Model,
				Prompt: mergedOption.Prompt,
				Stream: &mergedOption.Stream,
			}

			if mergedOption.System != "" {
				req.System = mergedOption.System
			}

			respFunc := func(resp ollamaAPI.GenerateResponse) error {
				if resp.Response != "" {
					self.State().GraphState().MergeChange(ollama.Generate, resp.Response)
				}
				if resp.Done {
					cancel()
				}
				return nil
			}

			err := ollamaClient.Generate(ctx, req, respFunc)
			if err != nil {
				return self.State(), fmt.Errorf("error calling Ollama Generate API: %w", err)
			}

			<-ctx.Done()
			self.State().ProceedOntoRoute() <- g.WhateverOutcome
			return self.State(), nil
		}
	case ollama.Chat:
		taskFn = func(msg f.Message, self f.Actor[g.NodeRef]) (g.NodeRef, error) {
			messages := []ollamaAPI.Message{}

			if mergedOption.System != "" {
				messages = append(messages, ollamaAPI.Message{
					Role:    "system",
					Content: mergedOption.System,
				})
			}

			var userUtterance string
			if mergedOption.UserUtterance != "" {
				if strings.HasPrefix(mergedOption.UserUtterance, "{{.") && strings.HasSuffix(mergedOption.UserUtterance, "}}") {
					attrName := strings.TrimPrefix(mergedOption.UserUtterance, "{{.")
					attrName = strings.TrimSuffix(attrName, "}}")
					// TODO manage an array of a chat turns
					userUtterance = self.State().GraphState().ReadAttribute(attrName).(string)
				} else {
					userUtterance = mergedOption.UserUtterance
				}

				messages = append(messages, ollamaAPI.Message{
					Role:    "user",
					Content: userUtterance,
				})
			}

			ctx, cancel := context.WithCancel(context.Background())

			req := &ollamaAPI.ChatRequest{
				Model:    kind.Model,
				Messages: messages,
				Stream:   &mergedOption.Stream,
			}

			respFunc := func(resp ollamaAPI.ChatResponse) error {
				if resp.Message.Content != "" {
					self.State().GraphState().MergeChange(ollama.Chat, resp.Message.Content)
				}
				if resp.Done {
					cancel()
				}
				return nil
			}

			err := ollamaClient.Chat(ctx, req, respFunc)
			if err != nil {
				return self.State(), fmt.Errorf("error calling Ollama Chat API: %w", err)
			}

			<-ctx.Done()
			self.State().ProceedOntoRoute() <- g.WhateverOutcome
			return self.State(), nil
		}
	}

	baseNode, err := newNode(forGraph, *address, taskFn)
	if err != nil {
		return nil, err
	}

	return &ollamaNode{
		node: *baseNode,
	}, nil
}
