package cmd

import (
	"fmt"

	"github.com/dll-as/gitc/internal/config"
	"github.com/urfave/cli/v2"
)

// Version defines the current version of the gitc tool.
const Version = "1.0.0"

var app *App

// Commands defines the CLI application configuration.
var Commands = &cli.App{
	Name:    "gitc",
	Usage:   "Generate AI-powered commit messages",
	Version: Version,
	Flags: []cli.Flag{
		// Git
		&cli.BoolFlag{
			Name:    "all",
			Aliases: []string{"a"},
			Usage:   "Stage all changes before generating commit message",
			EnvVars: []string{"GITC_STAGE_ALL"},
		},

		// Backend
		&cli.StringFlag{
			Name:    "backend",
			Usage:   "AI backend (openai-compatible, ollama, anthropic)",
			EnvVars: []string{"GITC_BACKEND"},
		},
		&cli.StringFlag{
			Name:    "base-url",
			Usage:   "Backend API URL",
			EnvVars: []string{"GITC_BASE_URL"},
		},
		&cli.StringFlag{
			Name:    "api-key",
			Aliases: []string{"k"},
			Usage:   "API key",
			EnvVars: []string{"AI_API_KEY"},
		},

		&cli.StringFlag{
			Name:    "model",
			Usage:   "Model name",
			EnvVars: []string{"GITC_MODEL"},
		},

		&cli.IntFlag{
			Name:    "timeout",
			Usage:   "HTTP timeout (seconds)",
			EnvVars: []string{"GITC_TIMEOUT"},
		},

		&cli.IntFlag{
			Name:    "max-retries",
			Usage:   "Maximum retry attempts",
			EnvVars: []string{"GITC_MAX_RETRIES"},
		},

		&cli.StringFlag{
			Name:    "proxy",
			Usage:   "HTTP proxy",
			EnvVars: []string{"GITC_PROXY"},
		},

		// Prompt
		&cli.StringFlag{
			Name:    "lang",
			Usage:   "Commit language",
			EnvVars: []string{"GITC_LANGUAGE"},
		},

		&cli.StringFlag{
			Name:    "convention",
			Usage:   "Commit convention",
			EnvVars: []string{"GITC_CONVENTION"},
		},

		&cli.StringFlag{
			Name:    "commit-type",
			Aliases: []string{"t"},
			Usage:   "Commit type",
			EnvVars: []string{"GITC_COMMIT_TYPE"},
		},

		&cli.StringFlag{
			Name:    "scope",
			Aliases: []string{"s"},
			Usage:   "Commit scope",
			EnvVars: []string{"GITC_SCOPE"},
		},

		&cli.BoolFlag{
			Name:    "emoji",
			Aliases: []string{"g"},
			Usage:   "Enable Gitmoji",
			EnvVars: []string{"GITC_GITMOJI"},
		},

		&cli.BoolFlag{
			Name:  "no-emoji",
			Usage: "Disable Gitmoji",
		},

		&cli.IntFlag{
			Name:    "max-tokens",
			Usage:   "Maximum output tokens",
			EnvVars: []string{"GITC_MAX_TOKENS"},
		},

		&cli.Float64Flag{
			Name:    "temperature",
			Usage:   "Sampling temperature",
			EnvVars: []string{"GITC_TEMPERATURE"},
		},

		&cli.Float64Flag{
			Name:    "top-p",
			Usage:   "Top-p sampling",
			EnvVars: []string{"GITC_TOP_P"},
		},

		// Git
		&cli.IntFlag{
			Name:    "max-diff-size",
			Usage:   "Maximum git diff size",
			EnvVars: []string{"GITC_MAX_DIFF_SIZE"},
		},

		// Misc
		// &cli.BoolFlag{
		// 	Name:    "dry-run",
		// 	Aliases: []string{"d"},
		// 	Usage:   "Print prompt without calling AI",
		// 	EnvVars: []string{"GITC_DRY_RUN"},
		// },
		&cli.BoolFlag{
			Name:    "debug",
			Aliases: []string{"d"},
			Usage:   "Enable debug mode (show raw API responses)",
			EnvVars: []string{"GITC_DEBUG"},
		},

		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Config file path",
			EnvVars: []string{"GITC_CONFIG_PATH"},
		},
	},
	Before: func(c *cli.Context) error {
		if c.Args().First() == "config" || c.Args().First() == "reset-config" {
			return nil
		}

		cfg, err := config.Load()
		if err != nil {
			return err
		}

		cfg.ApplyCLI(c)

		if err := cfg.Validate(); err != nil {
			return err
		}

		app, err = NewApp(cfg)
		if err != nil {
			return err
		}

		return nil
	},
	Action: func(c *cli.Context) error {
		return app.CommitAction(c)
	},
	Commands: []*cli.Command{
		{
			Name:    "config",
			Aliases: []string{"cfg"},
			Usage:   "Configure AI provider settings",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "backend",
					Usage:   "AI backend (openai-compatible, ollama, anthropic)",
					EnvVars: []string{"GITC_BACKEND"},
				},
				&cli.StringFlag{
					Name:    "base-url",
					Usage:   "Backend API URL",
					EnvVars: []string{"GITC_BASE_URL"},
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
					Name:  "max-tokens",
					Usage: "Maximum output tokens",
				},
				&cli.IntFlag{
					Name:  "max-length",
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
				&cli.BoolFlag{
					Name:  "debug",
					Usage: "Enable debug mode",
				},
				&cli.StringFlag{
					Name:    "config",
					Aliases: []string{"c"},
					Usage:   "Path to config file",
					EnvVars: []string{"GITC_CONFIG_PATH"},
				},
			},
			Action: func(c *cli.Context) error {
				cfg, err := config.Load()
				if err != nil {
					return err
				}

				app, err = NewApp(cfg)
				if err != nil {
					return err
				}

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
