package cmd

import (
	"fmt"
	"os"

	"github.com/makesikann/cmt/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration settings",
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Update a specific setting (e.g., api-key, language, model)",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Println("Error: Could not load config:", err)
			os.Exit(1)
		}

		key := args[0]
		val := args[1]

		switch key {
		case "provider":
			validProviders := []string{"gemini", "openai", "anthropic", "ollama", "groq"}
			isValid := false
			for _, p := range validProviders {
				if val == p {
					isValid = true
					break
				}
			}
			if !isValid {
				fmt.Printf("Error: invalid provider '%s'. Valid options: %v\n", val, validProviders)
				os.Exit(1)
			}
			cfg.Provider = val
		case "api-key", "api_key":
			cfg.ApiKey = val
		case "openai-api-key", "openai_api_key":
			cfg.OpenAIApiKey = val
		case "anthropic-api-key", "anthropic_api_key":
			cfg.AnthropicApiKey = val
		case "groq-api-key", "groq_api_key":
			cfg.GroqApiKey = val
		case "ollama-endpoint", "ollama_endpoint":
			cfg.OllamaEndpoint = val
		case "language":
			cfg.Language = val
		case "model":
			cfg.Model = val
		case "style":
			if val != "short" && val != "long" {
				fmt.Println("Error: style must be 'short' or 'long'")
				os.Exit(1)
			}
			cfg.Style = val
		default:
			fmt.Printf("Unknown setting key: %s\n", key)
			os.Exit(1)
		}

		if err := config.SaveConfig(cfg); err != nil {
			fmt.Println("Error: Could not save setting:", err)
			os.Exit(1)
		}
		fmt.Printf("Setting %s updated to '%s'.\n", key, val)
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Println("Error: Could not read config:", err)
			os.Exit(1)
		}

		fmt.Println("--- cmt Configuration ---")
		fmt.Printf("Current Provider: %s\n", cfg.Provider)
		fmt.Printf("Language        : %s\n", cfg.Language)
		fmt.Printf("Model           : %s\n", cfg.Model)
		fmt.Printf("Style           : %s\n", cfg.Style)
		fmt.Printf("Max Diff Lines  : %d\n", cfg.MaxDiffLines)
		fmt.Printf("Auto-Confirm    : %t\n", cfg.AutoConfirm)
		fmt.Println("\n--- API Keys & Endpoints ---")
		fmt.Printf("Gemini Key      : %s\n", maskKey(cfg.ApiKey))
		fmt.Printf("OpenAI Key      : %s\n", maskKey(cfg.OpenAIApiKey))
		fmt.Printf("Anthropic Key   : %s\n", maskKey(cfg.AnthropicApiKey))
		fmt.Printf("Groq Key        : %s\n", maskKey(cfg.GroqApiKey))
		fmt.Printf("Ollama Endpoint : %s\n", cfg.OllamaEndpoint)
	},
}

func maskKey(key string) string {
	if len(key) <= 8 {
		return "******"
	}
	return key[:4] + "..." + key[len(key)-4:]
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configShowCmd)
}
