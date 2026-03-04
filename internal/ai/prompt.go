package ai

import "fmt"

func BuildPrompt(diff, logs, language, style string) string {
	styleGuide := "Include a detailed body with bullet points if necessary."
	if style == "short" {
		styleGuide = "Output ONLY a single line (subject). Do NOT include a body or bullet points."
	}

	return fmt.Sprintf(`You are an expert software developer generating a git commit message.
Follow Conventional Commits. Output ONLY the message.

Rules:
1. First line (subject) must be under 72 chars.
2. %s
3. Language: %s

Context (Recent Logs):
%s

Changes:
%s`, styleGuide, language, logs, diff)
}
