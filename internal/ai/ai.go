package ai

import (
	"context"
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
