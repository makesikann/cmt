# cmt - AI-Powered Git Commit Message Generator

cmt is a Go-based CLI tool that uses Google Gemini AI to accelerate the git commit process for developers. It automatically generates messages following the Conventional Commits standards by analyzing staged changes (git diff --staged) and git log history, allowing the user to review and edit them interactively.

## Features

- Smart Generation with Gemini: Generates the most suitable commit message draft by analyzing your git diff and commit history.
- Interactive Confirmation: A stylish in-terminal preview and confirmation dialog built with Bubble Tea. Options: Commit, Edit, Cancel.
- Secure Configuration: Your Gemini API Key, language preferences, and other settings are stored securely in ~/.config/cmt/config.toml. Interactive setup starts on first run.
- Smart Diff Trimming: Handles large file changes efficiently by using segmented analysis for diffs over 500 or 1000 lines.

## Installation

### For Windows

Install quickly using PowerShell:

```powershell
irm https://raw.githubusercontent.com/makesikann/cmt/main/install.ps1 | iex
```

### macOS/Linux (Homebrew Tap Example)
```bash
brew install makesikann/tap/cmt
```

*Alternatively (Build with Go):*
```bash
go install github.com/makesikann/cmt@latest
```

## Usage

Once you have staged changes in your repository, simply run:

```bash
cmt commit
```

- If you haven't defined an API Key yet, cmt will prompt you for a Gemini API Key.
- The generated message will be displayed. You can commit it to your repo using [C]ommit, [E]dit, or [I]gnore/Cancel.

### Managing Configuration

Use the cmt config command to manage application settings:

```bash
# List settings
cmt config show

# Update a specific setting
cmt config set api-key YOUR-NEW-KEY
cmt config set language en
cmt config set model gemini-2.0-flash
```

## Roadmap
- --auto flag for non-interactive commits.
- Git hook installation support (cmt hook install).
- Multi-SDK support (OpenAI, Anthropic, etc.).
