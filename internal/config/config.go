package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/urfave/cli/v2"
)

type Backend string

const (
	BackendOpenAI    Backend = "openai-compatible"
	BackendOllama    Backend = "ollama"
	BackendAnthropic Backend = "anthropic"
)

type Config struct {
	Backend BackendConfig `json:"backend"`
	Prompt  PromptConfig  `json:"prompt"`
	Git     GitConfig     `json:"git"`
}

type BackendConfig struct {
	Backend Backend `json:"backend"`
	BaseURL string  `json:"base_url"`
	APIKey  string  `json:"api_key,omitempty"`
	Model   string  `json:"model"`

	Timeout    time.Duration `json:"timeout"`
	MaxRetries int           `json:"max_retries"`
	Proxy      string        `json:"proxy"`

	// Custom provider options
	ExtraHeaders map[string]string `json:"extra_headers,omitempty"`
	ExtraBody    map[string]any    `json:"extra_body,omitempty"`
	ExtraQuery   map[string]string `json:"extra_query,omitempty"`
}

type PromptConfig struct {
	Language    string  `json:"language"`
	Convention  string  `json:"convention"`
	CommitType  string  `json:"commit_type,omitempty"`
	Scope       string  `json:"scope,omitempty"`
	UseGitmoji  bool    `json:"use_gitmoji"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
	TopP        float64 `json:"top_p"`
	Debug       bool    `json:"debug"`
}

type GitConfig struct {
	MaxDiffSize int `json:"max_diff_size"`

	ExcludePatterns []string `json:"exclude_patterns,omitempty"`
}

// DefaultConfig returns a default config with fallback values
func DefaultConfig() *Config {
	return &Config{
		Backend: BackendConfig{
			Backend: BackendOpenAI,
			BaseURL: "https://api.openai.com/v1",
			APIKey:  "",
			Model:   "gpt-4o-mini",

			Timeout:    30 * time.Second,
			MaxRetries: 3,
			Proxy:      "",

			ExtraHeaders: nil,
			ExtraBody:    nil,
			ExtraQuery:   nil,
		},
		Prompt: PromptConfig{
			Language:    "en",
			Convention:  "conventional",
			CommitType:  "",
			Scope:       "",
			UseGitmoji:  false,
			MaxTokens:   1024,
			Temperature: 0.3,
			TopP:        1.0,
			Debug:       false,
		},
		Git: GitConfig{
			ExcludePatterns: []string{
				"package-lock.json", "pnpm-lock.yaml", "yarn.lock", "*.lock",
				"*.min.js", "*.bundle.js",
				"node_modules/*", "dist/*", "build/*",
				"*.log", "*.bak", "*.swp",
			},
			MaxDiffSize: 100000, // 100KB
		},
	}
}

// Load loads configuration from file and environment
func Load() (*Config, error) {
	cfg := DefaultConfig()

	// Load from config file
	path, err := getConfigPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get config path: %w", err)
	}

	data, err := os.ReadFile(path)
	switch {
	case errors.Is(err, os.ErrNotExist):
		if err = Save(cfg); err != nil {
			return nil, err
		}

		return cfg, nil
	case err != nil:
		return nil, err
	}

	if err = json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config file: %w", err)
	}

	// if err = cfg.Validate(); err != nil {
	// 	return nil, fmt.Errorf("validate config: %w", err)
	// }

	return cfg, nil
}

func (cfg *Config) ApplyCLI(c *cli.Context) {
	if c.IsSet("backend") {
		cfg.Backend.Backend = Backend(c.String("backend"))
	}

	if c.IsSet("base-url") {
		cfg.Backend.BaseURL = c.String("base-url")
	}

	if c.IsSet("api-key") {
		cfg.Backend.APIKey = c.String("api-key")
	}

	if c.IsSet("model") {
		cfg.Backend.Model = c.String("model")
	}

	if c.IsSet("timeout") {
		cfg.Backend.Timeout = time.Duration(c.Int("timeout")) * time.Second
	}

	if c.IsSet("proxy") {
		cfg.Backend.Proxy = c.String("proxy")
	}

	if c.IsSet("lang") {
		cfg.Prompt.Language = c.String("lang")
	}

	if c.IsSet("temperature") {
		cfg.Prompt.Temperature = c.Float64("temperature")
	}

	if c.IsSet("max-tokens") {
		cfg.Prompt.MaxTokens = c.Int("max-tokens")
	}

	if c.IsSet("scope") {
		cfg.Prompt.Scope = c.String("scope")
	}

	if c.IsSet("commit-type") {
		cfg.Prompt.CommitType = c.String("commit-type")
	}

	if c.IsSet("emoji") {
		cfg.Prompt.UseGitmoji = c.Bool("emoji")
	}

	if c.IsSet("no-emoji") {
		cfg.Prompt.UseGitmoji = !c.Bool("no-emoji")
	}

	if c.IsSet("debug") {
		cfg.Prompt.Debug = c.Bool("debug")
	}
}

// getConfigPath returns the path to the config file
// configPath points to ~/.gitc/config.json by default (hidden config file)
func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".gitc", "config.json"), nil
}

// Save saves configuration to file
func Save(cfg *Config) error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	if err = os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func Reset() error {
	return Save(DefaultConfig())
}

func (c *Config) Validate() error {
	if c.Backend.Backend == "" {
		return errors.New("backend is required")
	}

	switch c.Backend.Backend {
	case BackendOpenAI, BackendOllama:
	default:
		return fmt.Errorf("unsupported backend: %s", c.Backend.Backend)
	}

	if c.Backend.BaseURL == "" {
		return errors.New("base_url is required")
	}

	if c.Backend.Model == "" {
		return errors.New("model is required")
	}

	if c.Backend.Backend != "ollama" &&
		c.Backend.APIKey == "" {
		return errors.New("api_key is required")
	}

	if c.Backend.Timeout <= 0 {
		c.Backend.Timeout = 30 * time.Second
	}

	if c.Backend.MaxRetries < 0 {
		c.Backend.MaxRetries = 3
	}

	if c.Prompt.MaxTokens <= 0 {
		c.Prompt.MaxTokens = 1024
	}

	if c.Prompt.Temperature < 0 ||
		c.Prompt.Temperature > 2 {
		return errors.New("temperature must be between 0 and 2")
	}

	if c.Prompt.TopP <= 0 ||
		c.Prompt.TopP > 1 {
		return errors.New("top_p must be between 0 and 1")
	}

	if c.Git.MaxDiffSize <= 0 {
		c.Git.MaxDiffSize = 100000
	}

	return nil
}
