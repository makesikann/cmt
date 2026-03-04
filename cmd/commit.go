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

		aiClient, err := ai.NewClient(cfg.ApiKey, cfg.Model, cfg.Language)
		if err != nil {
			fmt.Println("Error: Could not create AI client:", err)
			os.Exit(1)
		}

		fmt.Println("Generating commit message with Gemini...")
		msg, err := aiClient.GenerateCommitMessage(diff, logs)
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
	rootCmd.AddCommand(commitCmd)
}
