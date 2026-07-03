package core

import (
	"context"
	"fmt"
	"strings"

	"github.com/dll-as/gitc/internal/ai/backend"
	"github.com/dll-as/gitc/internal/config"
)

type CommitGenerator struct {
	provider backend.Provider
	config   *config.Config
}

func NewCommitGenerator(provider backend.Provider, cfg *config.Config) *CommitGenerator {
	return &CommitGenerator{provider: provider, config: cfg}
}

func (g *CommitGenerator) Generate(ctx context.Context, diff string) (string, error) {
	diff = strings.TrimSpace(diff)
	if diff == "" {
		return "", fmt.Errorf("git diff is empty")
	}

	prompt := g.buildPrompt(diff)

	resp, err := g.provider.Generate(ctx, &backend.GenerateRequest{
		Prompt:      prompt,
		Model:       g.config.Backend.Model,
		MaxTokens:   g.config.Prompt.MaxTokens,
		Temperature: g.config.Prompt.Temperature,
		TopP:        g.config.Prompt.TopP,
	})
	if err != nil {
		return "", err
	}

	msg := strings.TrimSpace(resp.Message)

	// Apply Gitmoji if enabled
	if g.config.Prompt.UseGitmoji {
		msg = AddGitmojiToCommitMessage(msg)
	}

	return msg, nil
}
