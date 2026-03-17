package ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type OpenAIProvider struct {
	client   *openai.Client
	model    string
	language string
	style    string
}

func NewOpenAIProvider(apiKey string, model string, language string, style string, baseURL string) *OpenAIProvider {
	config := openai.DefaultConfig(apiKey)
	if baseURL != "" {
		config.BaseURL = baseURL
	}
	client := openai.NewClientWithConfig(config)

	return &OpenAIProvider{
		client:   client,
		model:    model,
		language: language,
		style:    style,
	}
}

func (p *OpenAIProvider) GenerateCommitMessage(diff string, logs string) (string, error) {
	prompt := BuildPrompt(diff, logs, p.language, p.style)

	resp, err := p.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: p.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return "", fmt.Errorf("openai error: %v", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from openai")
	}

	msg := strings.TrimSpace(resp.Choices[0].Message.Content)

	// Clean markdown if present
	msg = strings.TrimPrefix(msg, "```")
	msg = strings.TrimPrefix(msg, "git")
	msg = strings.TrimPrefix(msg, "\n")
	msg = strings.TrimSuffix(msg, "```")
	msg = strings.TrimSpace(msg)

	return msg, nil
}
