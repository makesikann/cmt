package git

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// CheckRepo verifies whether the current directory is a git repository.
func CheckRepo() error {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	if err := cmd.Run(); err != nil {
		return errors.New("this directory is not a git repository")
	}
	return nil
}

// GetRecentLogs retrieves the last 5 commit logs.
func GetRecentLogs() (string, error) {
	cmd := exec.Command("git", "log", "--oneline", "-5")
	var out bytes.Buffer
	cmd.Stdout = &out
	// It's acceptable if there are no commits yet.
	if err := cmd.Run(); err != nil {
		return "", nil
	}
	return strings.TrimSpace(out.String()), nil
}

// GetStagedDiff retrieves the staged differences, intelligently truncating them if they are too large.
func GetStagedDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--staged")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("could not retrieve staged changes: %v", err)
	}

	diffRaw := out.String()
	if diffRaw == "" {
		return "", errors.New("no staged changes. please use 'git add'")
	}

	return truncateDiff(diffRaw)
}

func truncateDiff(diffRaw string) (string, error) {
	lines := strings.Split(diffRaw, "\n")
	totalLines := len(lines)

	if totalLines <= 500 {
		return diffRaw, nil
	}

	if totalLines > 1000 {
		return getNumStat()
	}

	// Between 500 and 1000: max 100 lines per file
	return truncatePerFile(lines), nil
}

func getNumStat() (string, error) {
	cmd := exec.Command("git", "diff", "--staged", "--stat")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return fmt.Sprintf("Large diff detected (1000+ lines). Showing file statistics only:\n\n%s", out.String()), nil
}

func truncatePerFile(lines []string) string {
	var result []string
	var currentFileLines []string
	var lineCount int

	for _, line := range lines {
		if strings.HasPrefix(line, "diff --git") {
			// Process previous file
			if len(currentFileLines) > 0 {
				result = append(result, currentFileLines...)
				if lineCount > 100 {
					result = append(result, fmt.Sprintf("... (%d satır kırpıldı) ...", lineCount-100))
				}
			}
			// Reset for new file
			currentFileLines = []string{line}
			lineCount = 0
		} else {
			if lineCount < 100 {
				currentFileLines = append(currentFileLines, line)
			}
			lineCount++
		}
	}
	// Append the last file
	if len(currentFileLines) > 0 {
		result = append(result, currentFileLines...)
		if lineCount > 100 {
			result = append(result, fmt.Sprintf("... (%d satır kırpıldı) ...", lineCount-100))
		}
	}

	return strings.Join(result, "\n")
}
