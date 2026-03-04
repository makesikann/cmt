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
		return nil, fmt.Errorf("gemini istemcisi oluşturulamadı: %v", err)
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
		return "", fmt.Errorf("commit mesajı üretilemedi: %v", err)
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("API'den geçerli içerik alınamadı")
	}

	var msgBuilder strings.Builder
	for _, part := range resp.Candidates[0].Content.Parts {
		if text, ok := part.(genai.Text); ok {
			msgBuilder.WriteString(string(text))
		}
	}

	msg := strings.TrimSpace(msgBuilder.String())
	if msg == "" {
		return "", fmt.Errorf("kullanılabilir commit mesajı döndürülmedi")
	}

	// Remove markdown codeblocks if AI wraps it
	msg = strings.TrimPrefix(msg, "```")
	msg = strings.TrimPrefix(msg, "git")
	msg = strings.TrimPrefix(msg, "\n")
	msg = strings.TrimSuffix(msg, "```")
	msg = strings.TrimSpace(msg)

	return msg, nil
}
