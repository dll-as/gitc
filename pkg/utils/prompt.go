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

// formatGitCommand formats the git commit command for display based on message content.
// Handles both single-line and multi-line commit messages.
func FormatGitCommand(msg string) string {
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
