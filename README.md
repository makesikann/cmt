# cmt - AI-Powered Git Commit Message Generator

cmt is a Go-based CLI tool that uses Google Gemini AI to accelerate the git commit process for developers. It automatically generates messages following the Conventional Commits standards by analyzing staged changes (git diff --staged) and git log history, allowing the user to review and edit them interactively.

![cmt Demo](assets/demo.gif)


## Features

- **Smart Generation with Gemini**: Generates professional commit messages using the latest `gemini-2.5-flash` model.
- **Adjustable Styles**: Support for `short` (subject only) or `long` (subject + detailed body) commit messages.
- **Interactive Confirmation**: A stylish in-terminal preview and confirmation dialog built with Bubble Tea. Options: Commit, Edit, Cancel.
- **Secure Configuration**: Your Gemini API Key and preferences are stored securely in `~/.config/cmt/config.toml`.
- **Smart Diff Trimming**: Handles large file changes efficiently by using segmented analysis for diffs.
- **Premium Console UI**: Dynamically wrapped text boxes for better readability. No emojis for a clean, professional look.

## Installation

### For Windows

Install quickly using PowerShell:

```powershell
irm https://raw.githubusercontent.com/makesikann/cmt/main/install.ps1 | iex
```

### macOS/Linux

Using Homebrew Tap:

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

- **Short Style**: `cmt commit --short` or `cmt commit -s`
- **Long Style**: `cmt commit --long` or `cmt commit -l`

The generated message will be displayed in a dynamic preview box. You can proceed with **[C]ommit**, **[E]dit** the message, or **[I]gnore/Cancel**.

### Managing Configuration

Use the `cmt config` command to manage application settings:

```bash
# List settings
cmt config show

# Update a specific setting
cmt config set api-key YOUR-NEW-KEY
cmt config set style short    # or long
cmt config set language en    # e.g., turkish, english
cmt config set model gemini-2.5-flash
```

## License

This project is licensed under the [MIT License](LICENSE).

## Roadmap
- `--auto` flag for non-interactive commits.
- Git hook installation support (`cmt hook install`).
- Multi-SDK support (OpenAI, Anthropic, etc.).
