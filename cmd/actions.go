package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/dll-as/gitc/internal/ai"
	"github.com/dll-as/gitc/internal/ai/generic"
	"github.com/dll-as/gitc/internal/git"
	"github.com/dll-as/gitc/pkg/config"
	"github.com/dll-as/gitc/pkg/utils"
	"github.com/urfave/cli/v2"
)

// App encapsulates the core application logic and dependencies for gitc.
// It provides methods for AI configuration, commit message generation, and Git operations.
type App struct {
	gitService git.GitService
	config     config.Config
}

// NewApp creates a new App instance
func NewApp(gitService git.GitService, cfg *config.Config) *App {
	return &App{
		gitService: gitService,
		config:     *cfg,
	}
}

// ConfigureAI builds and validates the AI configuration from CLI context.
// It merges CLI flags with default values and performs validation.
func (a *App) ConfigureAI(c *cli.Context) (*ai.Config, error) {
	cfg := &ai.Config{
		Provider:   c.String("provider"),
		APIKey:     c.String("api-key"),
		Timeout:    time.Duration(c.Int("timeout")) * time.Second,
		Proxy:      c.String("proxy"),
		UseGitmoji: !c.Bool("no-emoji") && c.Bool("emoji"),
		URL:        c.String("url"),

		Message: ai.MessageOptions{
			Model:            c.String("model"),
			Language:         c.String("lang"),
			CommitType:       c.String("commit-type"),
			Scope:            c.String("scope"),
			CustomConvention: c.String("custom-convention"),
			MaxLength:        c.Int("max-length"),
			Temperature:      c.Float64("temperature"),
			MaxRedirects:     c.Int("max-redirects"),
		},
	}

	// Apply default values for unset fields
	a.applyConfigDefaults(cfg)

	// Validate the configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid AI configuration: %w", err)
	}

	return cfg, nil
}

// generateCommitMessage creates a commit message using AI based on the provided git diff.
// It handles AI provider initialization, timeout management, and Gitmoji formatting.
func (a *App) generateCommitMessage(ctx context.Context, diff string, cfg *ai.Config) (string, error) {
	provider, err := a.initAIProvider(cfg)
	if err != nil {
		return "", fmt.Errorf("failed to initialize AI provider: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()

	msg, err := provider.GenerateCommitMessage(ctx, diff, cfg.Message)
	if err != nil {
		return "", fmt.Errorf("failed to generate commit message: %w", err)
	}

	// Apply Gitmoji if enabled
	if cfg.UseGitmoji {
		msg = utils.AddGitmojiToCommitMessage(msg)
	}

	return msg, nil
}

// formatGitCommand formats the git commit command for display based on message content.
// Handles both single-line and multi-line commit messages.
func formatGitCommand(msg string) string {
	lines := strings.Split(msg, "\n")
	nonEmptyLines := make([]string, 0, len(lines))

	// Filter out empty lines
	for _, line := range lines {
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			nonEmptyLines = append(nonEmptyLines, trimmed)
		}
	}

	if len(nonEmptyLines) == 0 {
		return "git commit -m \"\""
	}

	if len(nonEmptyLines) == 1 {
		return fmt.Sprintf("git commit -m \"%s\"", nonEmptyLines[0])
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("git commit -m \"%s\"", nonEmptyLines[0]))

	for _, line := range nonEmptyLines[1:] {
		builder.WriteString(fmt.Sprintf(" \\\n    -m \"%s\"", line))
	}

	return builder.String()
}

// CommitAction handles the generation of commit messages
func (a *App) CommitAction(c *cli.Context) error {
	// Stage all changes if --all (-a) flag is set
	if c.Bool("all") {
		if err := a.gitService.StageAll(c.Context); err != nil {
			return fmt.Errorf("❌ failed to stage changes: %v", err)
		}

		fmt.Println("✅ All changes staged successfully")
	}

	// Fetch git diff for staged changes
	diff, err := a.gitService.GetDiff(c.Context)
	if err != nil {
		return fmt.Errorf("❌ failed to get git diff: %v", err)
	} else if diff == "" {
		return fmt.Errorf("❌ nothing staged for commit")
	}

	// Configure AI settings
	cfg, err := a.ConfigureAI(c)
	if err != nil {
		return fmt.Errorf("❌ failed to build AI config: %w", err)
	}

	// preview prompt + config, no API call
	if c.Bool("dry-run") {
		utils.PrintDryRun(diff, cfg)
		return nil
	}

	// Generate commit message
	msg, err := a.generateCommitMessage(c.Context, diff, cfg)
	if err != nil {
		return fmt.Errorf("❌ failed to generate commit message: %w", err)
	}

	// Display the generated command
	fmt.Println("✅ Commit message generated. You can now run:")
	fmt.Printf("   %s\n", formatGitCommand(msg))

	return nil
}

// ConfigAction handles updating and saving application configuration.
func (a *App) ConfigAction(c *cli.Context) error {
	newCfg := a.config

	a.updateConfigFromFlags(&newCfg, c)

	if err := newCfg.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	if err := config.Save(&newCfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	a.config = newCfg

	fmt.Println("Configuration updated successfully")
	return nil
}

// applyConfigDefaults sets sensible default values for unset AI configuration fields.
func (a *App) applyConfigDefaults(cfg *ai.Config) {
	if cfg.Provider == "" {
		cfg.Provider = a.config.Provider
	}
	if cfg.Message.Model == "" {
		switch cfg.Provider {
		case "openai":
			cfg.Message.Model = "gpt-4o-mini"
		case "grok":
			cfg.Message.Model = "grok-3"
		case "deepseek":
			cfg.Message.Model = "deepseek-rag"
		default:
			cfg.Message.Model = a.config.Model
		}
	}
	if cfg.APIKey == "" {
		cfg.APIKey = a.config.APIKey
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = time.Duration(a.config.Timeout) * time.Second
	}
	if cfg.Message.MaxLength == 0 {
		cfg.Message.MaxLength = a.config.MaxLength
	}
	if cfg.Message.Language == "" {
		cfg.Message.Language = a.config.Language
	}
	if cfg.Message.MaxRedirects == 0 {
		cfg.Message.MaxRedirects = a.config.MaxRedirects
	}
	if cfg.Message.Temperature == 0 {
		cfg.Message.Temperature = a.config.Temperature
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
			cfg.URL = a.config.URL
		}
	}
}

// initAIProvider initializes the appropriate AI provider based on configuration.
func (a *App) initAIProvider(cfg *ai.Config) (ai.AIProvider, error) {
	return generic.NewGenericProvider(cfg.APIKey, cfg.Proxy, cfg.URL, cfg.Provider)
}

// updateConfigFromFlags updates the configuration with values from CLI flags.
// Only updates fields that are explicitly set in the context.
func (a *App) updateConfigFromFlags(cfg *config.Config, c *cli.Context) {
	if provider := c.String("provider"); provider != "" {
		cfg.Provider = provider
	}
	if model := c.String("model"); model != "" {
		cfg.Model = model
	}
	if apiKey := c.String("api-key"); apiKey != "" {
		cfg.APIKey = apiKey
	}
	if lang := c.String("lang"); lang != "" {
		cfg.Language = lang
	}
	if timeout := c.Int("timeout"); timeout != 0 {
		cfg.Timeout = timeout
	}
	if maxLength := c.Int("maxLength"); maxLength != 0 {
		cfg.MaxLength = maxLength
	}
	if proxy := c.String("proxy"); proxy != "" {
		cfg.Proxy = proxy
	}
	if commitType := c.String("commit-type"); commitType != "" {
		cfg.CommitType = commitType
	}
	if customConvention := c.String("custom-convention"); customConvention != "" {
		cfg.CustomConvention = customConvention
	}
	if c.IsSet("no-emoji") {
		cfg.UseGitmoji = !c.Bool("no-emoji")
	} else if c.IsSet("emoji") {
		cfg.UseGitmoji = c.Bool("emoji")
	}
	if maxRedirects := c.Int("max-redirects"); maxRedirects != 0 {
		cfg.MaxRedirects = maxRedirects
	}
	if url := c.String("url"); url != "" {
		cfg.URL = url
	}
}
