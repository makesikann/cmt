package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type AnthropicProvider struct {
	apiKey   string
	model    string
	language string
	style    string
}

type anthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	Messages  []anthropicMessage `json:"messages"`
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

func NewAnthropicProvider(apiKey string, model string, language string, style string) *AnthropicProvider {
	return &AnthropicProvider{
		apiKey:   apiKey,
		model:    model,
		language: language,
		style:    style,
	}
}

func (p *AnthropicProvider) GenerateCommitMessage(diff string, logs string) (string, error) {
	prompt := BuildPrompt(diff, logs, p.language, p.style)

	reqBody := anthropicRequest{
		Model:     p.model,
		MaxTokens: 1024,
		Messages: []anthropicMessage{
			{Role: "user", Content: prompt},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", p.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		var errResp anthropicResponse
		json.Unmarshal(body, &errResp)
		if errResp.Error.Message != "" {
			return "", fmt.Errorf("anthropic api error: %s", errResp.Error.Message)
		}
		return "", fmt.Errorf("anthropic api error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var anthropicResp anthropicResponse
	if err := json.Unmarshal(body, &anthropicResp); err != nil {
		return "", err
	}

	if len(anthropicResp.Content) == 0 {
		return "", fmt.Errorf("no response content from anthropic")
	}

	msg := strings.TrimSpace(anthropicResp.Content[0].Text)

	// Clean markdown if present
	msg = strings.TrimPrefix(msg, "```")
	msg = strings.TrimPrefix(msg, "git")
	msg = strings.TrimPrefix(msg, "\n")
	msg = strings.TrimSuffix(msg, "```")
	msg = strings.TrimSpace(msg)

	return msg, nil
}
