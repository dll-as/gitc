package cmd

import (
	"fmt"

	"github.com/dll-as/gitc/internal/git"
	"github.com/dll-as/gitc/pkg/config"
	"github.com/urfave/cli/v2"
)

// Version defines the current version of the gitc tool.
const Version = "0.3.0"

var appInstance *App

// Commands defines the CLI application configuration.
var Commands = &cli.App{
	Name:    "gitc",
	Usage:   "Generate AI-powered commit messages",
	Version: Version,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "all",
			Aliases: []string{"a"},
			Usage:   "Stage all changes before generating commit message (equivalent to 'git add .')",
			EnvVars: []string{"GITC_STAGE_ALL"},
		},
		&cli.StringFlag{
			Name:  "provider",
			Value: "openai",
			Usage: "AI provider to use (openai, anthropic)",
		},
		&cli.StringFlag{
			Name:  "model",
			Value: "gpt-4o-mini",
			Usage: "Specify the OpenAI model",
		},
		&cli.StringFlag{
			Name:  "lang",
			Value: "en",
			Usage: "Set commit message language (en, fa, ru, etc.)",
		},
		&cli.IntFlag{
			Name:  "timeout",
			Value: 10,
			Usage: "Set request timeout in seconds",
		},
		&cli.IntFlag{
			Name:  "maxLength",
			Value: 200,
			Usage: "Set maximum output length of AI response",
		},
		&cli.StringFlag{
			Name:    "api-key",
			Aliases: []string{"k"},
			Usage:   "API key for the AI provider",
			EnvVars: []string{"AI_API_KEY"},
		},
		&cli.StringFlag{
			Name:    "proxy",
			Aliases: []string{"p"},
			Usage:   "Proxy URL for API requests (e.g., http://proxy.example.com:8080)",
			EnvVars: []string{"GITC_PROXY"},
		},
		&cli.StringFlag{
			Name:    "commit-type",
			Aliases: []string{"t"},
			Usage:   "Commit type for Conventional Commits (e.g., feat, fix, docs)",
			EnvVars: []string{"GITC_COMMIT_TYPE"},
		},
		&cli.StringFlag{
			Name:    "custom-convention",
			Aliases: []string{"C"},
			Usage:   "Custom commit message convention in JSON format (e.g., '{\"prefix\": \"JIRA-123\"}')",
			EnvVars: []string{"GITC_CUSTOM_CONVENTION"},
		},
		&cli.StringFlag{
			Name:    "scope",
			Aliases: []string{"s"},
			Usage:   "Add scope to commit type (e.g., feat(auth): ...)",
		},
		&cli.BoolFlag{
			Name:    "emoji",
			Aliases: []string{"g"},
			Usage:   "Add Gitmoji to the commit message based on commit type",
			EnvVars: []string{"GITC_GITMOJI"},
		},
		&cli.BoolFlag{
			Name:  "no-emoji",
			Usage: "Disable Gitmoji in the commit message (overrides --emoji)",
		},
		&cli.IntFlag{
			Name:    "max-redirects",
			Aliases: []string{"r"},
			Value:   5,
			Usage:   "Maximum number of HTTP redirects to follow",
			EnvVars: []string{"GITC_MAX_REDIRECTS"},
		},
		&cli.Float64Flag{
			Name:  "temperature",
			Value: 0.7,
			Usage: "Control creativity (0.0 = deterministic, 1.0 = very creative)",
		},
		&cli.BoolFlag{
			Name:    "dry-run",
			Aliases: []string{"d"},
			Usage:   "Show the processed diff and prompt without sending to AI",
			EnvVars: []string{"GITC_DRY_RUN"},
		},
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Path to config file",
			EnvVars: []string{"GITC_CONFIG_PATH"},
		},
	},
	Before: func(c *cli.Context) error {
		// Set config path if provided via flag or environment variable
		if configPath := c.String("config"); configPath != "" {
			config.SetConfigPath(configPath)
		}

		// Load config
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Initialize dependencies
		gitService := git.NewGitService()
		appInstance = NewApp(gitService, cfg)
		return nil
	},
	Action: func(c *cli.Context) error {
		return appInstance.CommitAction(c)
	},
	Commands: []*cli.Command{
		{
			Name:    "config",
			Aliases: []string{"cfg"},
			Usage:   "Configure AI provider settings",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "provider",
					Aliases: []string{"ai"},
					Usage:   "AI provider to use (openai, anthropic)",
				},
				&cli.StringFlag{
					Name:  "model",
					Usage: "Specify the OpenAI model",
				},
				&cli.StringFlag{
					Name:  "lang",
					Usage: "Set commit message language (en, fa, ru, etc.)",
				},
				&cli.StringFlag{
					Name:    "proxy",
					Aliases: []string{"p"},
					Usage:   "Proxy URL for API requests (e.g., http://proxy.example.com:8080)",
				},
				&cli.IntFlag{
					Name:  "timeout",
					Usage: "Set request timeout in seconds",
				},
				&cli.IntFlag{
					Name:  "maxLength",
					Usage: "Set maximum output length of AI response",
				},
				&cli.IntFlag{
					Name:    "max-redirects",
					Aliases: []string{"r"},
					Usage:   "Set maximum number of HTTP redirects",
				},
				&cli.StringFlag{
					Name:    "api-key",
					Aliases: []string{"k"},
					Usage:   "API key for the AI provider",
				},
				&cli.StringFlag{
					Name:    "commit-type",
					Aliases: []string{"t"},
					Usage:   "Commit type for Conventional Commits (e.g., feat, fix, docs)",
				},
				&cli.StringFlag{
					Name:    "custom-convention",
					Aliases: []string{"C"},
					Usage:   "Custom commit message convention in JSON format (e.g., '{\"prefix\": \"JIRA-123\"}')",
				},
				&cli.BoolFlag{
					Name:    "emoji",
					Aliases: []string{"g"},
					Usage:   "Add Gitmoji to the commit message based on commit type",
				},
				&cli.BoolFlag{
					Name:  "no-emoji",
					Usage: "Disable Gitmoji in the commit message",
				},
				&cli.StringFlag{
					Name:    "config",
					Aliases: []string{"c"},
					Usage:   "Path to config file",
					EnvVars: []string{"GITC_CONFIG_PATH"},
				},
			},
			Action: func(c *cli.Context) error {
				// Set config path if provided
				if configPath := c.String("config"); configPath != "" {
					config.SetConfigPath(configPath)
				}

				// Load config
				cfg, err := config.Load()
				if err != nil {
					return fmt.Errorf("failed to load config: %w", err)
				}

				// Initialize dependencies
				gitService := git.NewGitService()
				app := NewApp(gitService, cfg)
				return app.ConfigAction(c)
			},
		}, {
			Name:  "reset-config",
			Usage: "Reset gitc configuration to default values",
			Action: func(c *cli.Context) error {
				if err := config.Reset(); err != nil {
					return fmt.Errorf("failed to reset config: %w", err)
				}

				fmt.Println("✅ Configuration reset to defaults.")
				return nil
			},
		},
	},
}
