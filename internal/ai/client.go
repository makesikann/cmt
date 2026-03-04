package ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type Client struct {
	genaiClient *genai.Client
	model       string
	language    string
}

func NewClient(apiKey string, model string, language string) (*Client, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("could not create gemini client: %v", err)
	}

	return &Client{
		genaiClient: client,
		model:       model,
		language:    language,
	}, nil
}

func (c *Client) GenerateCommitMessage(diff string, logs string) (string, error) {
	ctx := context.Background()
	prompt := BuildPrompt(diff, logs, c.language)

	model := c.genaiClient.GenerativeModel(c.model)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("could not generate commit message: %v", err)
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("could not retrieve valid content from API")
	}

	var msgBuilder strings.Builder
	for _, part := range resp.Candidates[0].Content.Parts {
		if text, ok := part.(genai.Text); ok {
			msgBuilder.WriteString(string(text))
		}
	}

	msg := strings.TrimSpace(msgBuilder.String())
	if msg == "" {
		return "", fmt.Errorf("no usable commit message returned")
	}

	// Remove markdown codeblocks if AI wraps it
	msg = strings.TrimPrefix(msg, "```")
	msg = strings.TrimPrefix(msg, "git")
	msg = strings.TrimPrefix(msg, "\n")
	msg = strings.TrimSuffix(msg, "```")
	msg = strings.TrimSpace(msg)

	return msg, nil
}
