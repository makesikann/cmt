package commit

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	message string
	choice  string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "i", "I":
			m.choice = "cancel"
			return m, tea.Quit
		case "c", "C", "enter":
			m.choice = "commit"
			return m, tea.Quit
		case "e", "E":
			m.choice = "edit"
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	s := "┌─────────────────────────────────────────────────┐\n"

	lines := strings.Split(m.message, "\n")
	for _, l := range lines {
		// Just a simple UI box
		if len(l) > 47 {
			s += fmt.Sprintf("│ %-47s │\n", l[:44]+"...")
		} else {
			s += fmt.Sprintf("│ %-47s │\n", l)
		}
	}
	s += "└─────────────────────────────────────────────────┘\n\n"
	s += "[C]ommit  [E]dit  [I]gnore  > "
	return s
}

// ConfirmMessage interactively shows the AI created commit message
// and lets the user choose between saving it, editing it, or canceling.
func ConfirmMessage(message string) (string, error) {
	m := model{message: message}
	p := tea.NewProgram(m)

	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	fm := finalModel.(model)

	switch fm.choice {
	case "commit":
		return PerformCommit(fm.message)
	case "edit":
		return editMessage(fm.message)
	case "cancel":
		fmt.Println("Cancelled.")
		os.Exit(0)
	default:
		fmt.Println("Cancelled.")
		os.Exit(0)
	}
	return "", nil
}

func PerformCommit(message string) (string, error) {
	cmd := exec.Command("git", "commit", "-m", message)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git commit error: %v\nOutput: %s", err, string(out))
	}
	fmt.Println("Git commit successful.")
	return string(out), nil
}

func editMessage(message string) (string, error) {
	tmpFile, err := os.CreateTemp("", "cmt_EDITMSG_*.txt")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(message)
	if err != nil {
		return "", err
	}
	tmpFile.Close()

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "notepad" // Fallback fallback for Windows Default
	}

	cmd := exec.Command(editor, tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("editör çalıştırılamadı: %v", err)
	}

	bytes, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return "", err
	}

	newMessage := strings.TrimSpace(string(bytes))
	if newMessage == "" {
		fmt.Println("Message left empty. Cancelled.")
		os.Exit(0)
	}

	return PerformCommit(newMessage)
}
