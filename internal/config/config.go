package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

type Config struct {
	ApiKey       string `toml:"api_key"`
	Language     string `toml:"language"`
	Model        string `toml:"model"`
	MaxDiffLines int    `toml:"max_diff_lines"`
	AutoConfirm  bool   `toml:"auto_confirm"`
	Style        string `toml:"style"` // "short" or "long"
}

func getConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "cmt"), nil
}

func getConfigPath() (string, error) {
	dir, err := getConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.toml"), nil
}

func LoadConfig() (*Config, error) {
	path, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	var cfg Config
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// File does not exist, ask for API key
		cfg = promptForKeyAndCreateDefault()
		if err := SaveConfig(&cfg); err != nil {
			return nil, err
		}
		return &cfg, nil
	}

	// File exists, decode it
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, err
	}

	// Check if API key is empty
	if cfg.ApiKey == "" {
		cfg.ApiKey = promptForAPIKey()
		if err := SaveConfig(&cfg); err != nil {
			return nil, err
		}
	}

	return &cfg, nil
}

func SaveConfig(cfg *Config) error {
	dir, err := getConfigDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	path, err := getConfigPath()
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	return encoder.Encode(cfg)
}

func promptForKeyAndCreateDefault() Config {
	apiKey := promptForAPIKey()
	return Config{
		ApiKey:       apiKey,
		Language:     "en",
		Model:        "gemini-2.5-flash",
		MaxDiffLines: 500,
		AutoConfirm:  false,
		Style:        "long",
	}
}

func promptForAPIKey() string {
	fmt.Println("Gemini API key not found.")
	fmt.Print("Enter your key: ")
	reader := bufio.NewReader(os.Stdin)
	key, _ := reader.ReadString('\n')
	key = strings.TrimSpace(key)
	fmt.Println("Saved to ~/.config/cmt/config.toml.")
	return key
}
