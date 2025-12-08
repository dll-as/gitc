package ai

import (
	"context"
	"fmt"
	"time"
)

// AIProvider defines the interface for AI providers
type AIProvider interface {
	GenerateCommitMessage(ctx context.Context, diff string, opts MessageOptions) (string, error)
}

// MessageOptions contains only fields needed to generate a commit message
type MessageOptions struct {
	Model            string
	Language         string
	CommitType       string
	Scope            string
	CustomConvention string
	MaxLength        int
	Temperature      float64
	MaxRedirects     int
}

// Config holds the full application configuration
type Config struct {
	Provider   string
	APIKey     string
	URL        string
	Timeout    time.Duration
	Proxy      string
	UseGitmoji bool

	Message MessageOptions
}

// validateConfig performs basic validation of the AI configuration.
// Returns an error if required fields are missing or invalid.
func (c *Config) Validate() error {
	if c.Provider == "" {
		return fmt.Errorf("provider is required")
	}
	if c.APIKey == "" {
		return fmt.Errorf("API key is required")
	}
	if c.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}
	if c.Message.MaxLength <= 0 {
		return fmt.Errorf("max length must be positive")
	}
	if c.Message.Temperature < 0 || c.Message.Temperature > 2 {
		return fmt.Errorf("temperature must be between 0 and 2")
	}
	return nil
}
