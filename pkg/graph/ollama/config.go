package ollama

import (
	"fmt"
)

// TODO
type Kind string

const (
	Chat     Kind = "chat"
	Generate Kind = "generate"
)

type KindWithModel struct {
	Kind
	Model string
}

type NodeOption struct {
	Stream        bool
	System        string
	UserUtterance string
	Prompt        string
}

func ChatWithModel(model string) KindWithModel {
	return KindWithModel{
		Kind:  Chat,
		Model: model,
	}
}

func GenerateWithModel(model string) KindWithModel {
	return KindWithModel{
		Kind:  Generate,
		Model: model,
	}
}

func WithPrompt(prompt string) NodeOption {
	return NodeOption{
		Prompt: prompt,
	}
}

func WithUserUtterance(userUtterance string) NodeOption {
	return NodeOption{
		UserUtterance: userUtterance,
	}
}

func WithStream() NodeOption {
	return NodeOption{
		Stream: true,
	}
}

func WithSystem(system string) NodeOption {
	return NodeOption{
		System: system,
	}
}

func MergeNodeOptions(options ...NodeOption) (NodeOption, error) {
	var result NodeOption
	streamSet := false
	systemSet := false
	userUtteranceSet := false
	promptSet := false

	for _, option := range options {
		if option.Stream && streamSet {
			return NodeOption{}, fmt.Errorf("conflicting Stream options")
		}
		if option.Stream {
			result.Stream = option.Stream
			streamSet = true
		}

		if option.System != "" && systemSet {
			return NodeOption{}, fmt.Errorf("conflicting System options")
		}
		if option.System != "" {
			result.System = option.System
			systemSet = true
		}

		if option.UserUtterance != "" && userUtteranceSet {
			return NodeOption{}, fmt.Errorf("conflicting UserUtterance options")
		}
		if option.UserUtterance != "" {
			result.UserUtterance = option.UserUtterance
			userUtteranceSet = true
		}

		if option.Prompt != "" && promptSet {
			return NodeOption{}, fmt.Errorf("conflicting Prompt options")
		}
		if option.Prompt != "" {
			result.Prompt = option.Prompt
			promptSet = true
		}
	}

	return result, nil
}
