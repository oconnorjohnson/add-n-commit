package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config holds the application configuration
type Config struct {
	OpenAIKey        string `json:"openai_key"`
	Model            string `json:"model"`
	DefaultMode      string `json:"default_mode"`      // "all", "by-file", "interactive"
	AutoStageAll     bool   `json:"auto_stage_all"`    // Whether to auto-stage all files
	Temperature      float32 `json:"temperature"`
	SystemPromptAll  string `json:"system_prompt_all"`
	SystemPromptFile string `json:"system_prompt_file"`
}

// Default returns the default configuration
func Default() *Config {
	return &Config{
		Model:            "o4-mini",
		DefaultMode:      "interactive",
		AutoStageAll:     false,
		Temperature:      1.0,
		SystemPromptAll:  "You are a helpful AI that writes clear and concise Git commit messages based on diffs.",
		SystemPromptFile: "You are a helpful AI that writes concise Git commit messages per file.",
	}
}

// Load loads configuration from file
func Load() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(homeDir, ".config", "anc", "config.json")
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		// Try environment variable for API key
		cfg := Default()
		if apiKey := os.Getenv("OPENAI_API_KEY"); apiKey != "" {
			cfg.OpenAIKey = apiKey
		}
		return cfg, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Override with environment variable if set
	if apiKey := os.Getenv("OPENAI_API_KEY"); apiKey != "" {
		cfg.OpenAIKey = apiKey
	}

	return &cfg, nil
}

// Save saves the configuration to file
func (c *Config) Save() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(homeDir, ".config", "anc")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	configPath := filepath.Join(configDir, "config.json")
	
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0600)
} 