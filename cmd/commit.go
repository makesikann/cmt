package cmd

import (
	"fmt"
	"os"

	"github.com/makesikann/cmt/internal/ai"
	"github.com/makesikann/cmt/internal/commit"
	"github.com/makesikann/cmt/internal/config"
	"github.com/makesikann/cmt/internal/git"
	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Generate a commit message with AI",
	Run: func(cmd *cobra.Command, args []string) {
		isShort, _ := cmd.Flags().GetBool("short")
		isLong, _ := cmd.Flags().GetBool("long")

		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Println("Error: Could not load config:", err)
			os.Exit(1)
		}

		if err := git.CheckRepo(); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		fmt.Println("Analyzing changes...")
		diff, err := git.GetStagedDiff()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		logs, err := git.GetRecentLogs()
		if err != nil {
			// Non-fatal error
		}

		style := cfg.Style
		if isShort {
			style = "short"
		} else if isLong {
			style = "long"
		}

		var aiProvider ai.Provider
		// err is already declared above

		switch cfg.Provider {
		case "gemini":
			aiProvider, err = ai.NewGeminiProvider(cfg.ApiKey, cfg.Model, cfg.Language, style)
		default:
			// Default to gemini for backward compatibility if provider is not set
			aiProvider, err = ai.NewGeminiProvider(cfg.ApiKey, cfg.Model, cfg.Language, style)
		}

		if err != nil {
			fmt.Println("Error: Could not create AI provider:", err)
			os.Exit(1)
		}

		fmt.Printf("Generating commit message with %s...\n", cfg.Provider)
		msg, err := aiProvider.GenerateCommitMessage(diff, logs)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		if cfg.AutoConfirm {
			// Auto confirm
			_, err = commit.PerformCommit(msg)
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
			return
		}

		_, err = commit.ConfirmMessage(msg)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
}

func init() {
	commitCmd.Flags().BoolP("short", "s", false, "Generate a single-line commit message")
	commitCmd.Flags().BoolP("long", "l", false, "Generate a detailed commit message (default)")
	rootCmd.AddCommand(commitCmd)
}
