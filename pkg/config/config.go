package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bytedance/sonic"
)

// Config holds the main configuration structure for the gitc CLI tool
type Config struct {
	Provider         string  `json:"provider"`
	APIKey           string  `json:"api_key"`
	Model            string  `json:"model"`
	URL              string  `json:"url"`
	MaxLength        int     `json:"max_length"`
	Temperature      float64 `json:"temperature"`
	Proxy            string  `json:"proxy"`
	Language         string  `json:"language"`
	Timeout          int     `json:"timeout"`
	CommitType       string  `json:"commit_type"`
	CustomConvention string  `json:"custom_convention"`
	UseGitmoji       bool    `json:"use_gitmoji"`
	MaxRedirects     int     `json:"max_redirects"`
}

// DefaultConfig returns a default config with fallback values
func DefaultConfig() *Config {
	return &Config{
		Provider:         "openai",
		APIKey:           os.Getenv("AI_API_KEY"),
		Model:            "gpt-4o-mini",
		URL:              "https://api.openai.com/v1/chat/completions",
		MaxLength:        250,
		Temperature:      0.7,
		Proxy:            "",
		Language:         "en",
		Timeout:          10,
		CommitType:       "",
		CustomConvention: "",
		UseGitmoji:       false,
		MaxRedirects:     5,
	}
}

// configPath points to ~/.gitc/config.json by default (hidden config file)
var configPath = filepath.Join(userHomeDir(), ".gitc", "config.json")

// SetConfigPath updates the configuration file path
func SetConfigPath(path string) {
	configPath = path
}

// Load loads the configuration from file or creates a default one if it doesn't exist
func Load() (*Config, error) {
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve config path: %w", err)
	}

	data, err := os.ReadFile(absPath)
	if os.IsNotExist(err) {
		cfg := DefaultConfig()
		if err := Save(cfg); err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}
		return cfg, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := sonic.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Apply defaults for unset fields
	defaults := DefaultConfig()
	if cfg.Provider == "" {
		cfg.Provider = defaults.Provider
	}
	if cfg.APIKey == "" {
		cfg.APIKey = defaults.APIKey
	}
	if cfg.Model == "" {
		switch cfg.Provider {
		case "openai":
			cfg.Model = "gpt-4o-mini"
		case "grok":
			cfg.Model = "grok-3"
		case "deepseek":
			cfg.Model = "deepseek-rag"
		default:
			cfg.Model = defaults.Model
		}
	}
	if cfg.URL == "" {
		switch cfg.Provider {
		case "openai":
			cfg.URL = "https://api.openai.com/v1/chat/completions"
		case "grok":
			cfg.URL = "https://api.x.ai/v1/chat/completions"
		case "deepseek":
			cfg.URL = "https://api.deepseek.com/v1/chat/completions"
		default:
			cfg.URL = defaults.URL
		}
	}
	if cfg.MaxLength == 0 {
		cfg.MaxLength = defaults.MaxLength
	}
	if cfg.Language == "" {
		cfg.Language = defaults.Language
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = defaults.Timeout
	}
	if cfg.MaxRedirects == 0 {
		cfg.MaxRedirects = defaults.MaxRedirects
	}

	return &cfg, nil
}

// Save saves the configuration to file
func Save(cfg *Config) error {
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return fmt.Errorf("failed to resolve config path: %w", err)
	}

	dir := filepath.Dir(absPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := sonic.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	} else if err := os.WriteFile(absPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Reset overwrites the config file with default values
func Reset() error {
	return Save(DefaultConfig())
}

// userHomeDir gets the current user's home directory
func userHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic("cannot determine user home directory: " + err.Error())
	}
	return home
}
