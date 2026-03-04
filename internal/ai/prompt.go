package ai

import "fmt"

func BuildPrompt(diff, logs, language string) string {
	return fmt.Sprintf(`You are an expert software developer, and your task is to generate a git commit message.
Please strictly follow these rules:
1. Use the Conventional Commits format (feat:, fix:, chore:, refactor:, docs:, etc.).
2. The first line (subject) must not exceed 72 characters.
3. Leave a blank line after the subject and add a detailed description (body) in bullet points if necessary.
4. Output ONLY the commit message. Do NOT include any introductory sentences, feedback, or markdown code blocks (***).
5. Set the language of the commit message to '%s'.

Recent logs of the repository for reference:
%s

Staged changes to analyze:
%s`, language, logs, diff)
}
