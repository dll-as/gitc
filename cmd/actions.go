package cmd

import (
	"fmt"
	"strings"

	"github.com/dll-as/gitc/internal/ai"
	"github.com/dll-as/gitc/internal/ai/backend"
	"github.com/dll-as/gitc/internal/config"
	"github.com/dll-as/gitc/internal/core"
	"github.com/dll-as/gitc/internal/git"
	"github.com/urfave/cli/v2"
)

type App struct {
	git     git.Service
	config  *config.Config
	service *core.CommitGenerator
	AI      backend.Provider
}

func NewApp(cfg *config.Config) (*App, error) {
	gitSvc := git.NewService(cfg.Git.ExcludePatterns)

	provider, err := ai.New(cfg)
	if err != nil {
		return nil, err
	}

	gen := core.NewCommitGenerator(provider, cfg)

	return &App{
		config:  cfg,
		git:     gitSvc,
		AI:      provider,
		service: gen,
	}, nil
}

// CommitAction handles the generation of commit messages
func (a *App) CommitAction(c *cli.Context) error {
	if c.NArg() > 0 {
		files := c.Args().Slice()
		if err := a.git.StageFiles(c.Context, files); err != nil {
			return fmt.Errorf("failed to stage files: %w", err)
		}

		fmt.Printf("Staged %d file(s): %s\n", len(files), strings.Join(files, ", "))
	} else if c.Bool("all") { // Stage all changes if --all (-a) flag is set
		if err := a.git.StageAll(c.Context); err != nil {
			return fmt.Errorf("❌ failed to stage changes: %v", err)
		}

		fmt.Println("✅ All changes staged successfully")
	}

	// Fetch git diff for staged changes
	diff, err := a.git.GetDiff(c.Context)
	if err != nil {
		return fmt.Errorf("❌ failed to get git diff: %v", err)
	} else if diff == "" {
		return fmt.Errorf("❌ nothing staged for commit")
	}

	if c.Bool("dry-run") {
		// utils.PrintDryRun(diff, cfg)
		return nil
	}

	// Generate commit message
	msg, err := a.service.Generate(c.Context, diff)
	if err != nil {
		return fmt.Errorf("❌ failed to generate commit message: %w", err)
	}

	// Display the generated command
	fmt.Println("✅ Commit message generated. You can now run:")
	fmt.Printf("   %s\n", a.service.FormatGitCommand(msg))

	return nil
}

func (a *App) ConfigAction(c *cli.Context) error {
	// Merge CLI flags into current config
	a.config.ApplyCLI(c)

	// Validate final config
	if err := a.config.Validate(); err != nil {
		return fmt.Errorf("validate config: %w", err)
	}

	// Save config
	if err := config.Save(a.config); err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	fmt.Println("✓ Configuration saved successfully")
	return nil
}

func (a *App) ResetAction(c *cli.Context) error {

	return nil
}
