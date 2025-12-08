package utils

import (
	"fmt"
	"strings"

	"github.com/dll-as/gitc/internal/ai"
)

// commitType, customMessageConvention, language, scope string
func GetPromptForSingleCommit(diff string, opts ai.MessageOptions) string {
	opts.Language = strings.ToLower(strings.TrimSpace(opts.Language))
	if opts.Language == "" {
		opts.Language = "en"
	}

	return fmt.Sprintf(`Write a concise Git commit message in %s based on this diff:

	%s

	Format:
	Line 1: <type>: <summary> (≤50 chars)
	Line 2: (blank)
	Line 3+: (optional) details (≤100 chars per line)

	Rules:
	- Use imperative mood (e.g. Add, Fix, Refactor)
	- Be clear and specific
	- %s
	- %s
	- No emoji, quotes, Markdown, or explanations

	Examples:
	feat: add JWT middleware

	Add access token check to protected routes.

	fix: prevent crash on nil DB config

	Add nil check before DB usage.`,
		opts.Language,
		diff,
		getTypeInstruction(opts.CommitType, opts.Scope),
		getConventionInstruction(opts.CustomConvention))
}

func getTypeInstruction(commitType, scope string) string {
	switch {
	case commitType != "" && scope != "":
		return fmt.Sprintf("Use exactly: %s(%s): <summary>", commitType, scope)

	case commitType != "":
		return fmt.Sprintf("Use exactly: %s: <summary>", commitType)

	case scope != "":
		return fmt.Sprintf("Choose type but MUST use scope: (%s)", scope)

	default:
		return "Use Conventional Commits format"
	}
}

func getConventionInstruction(convention string) string {
	if convention != "" {
		return fmt.Sprintf("Follow custom convention: %s", convention)
	}
	return "Follow Conventional Commits"
}
