package utils

import (
	"fmt"
	"strings"
)

func GetPromptForSingleCommit(diff, commitType, customMessageConvention, language, scope string) string {
	language = strings.ToLower(strings.TrimSpace(language))
	if language == "" {
		language = "en"
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
		language,
		diff,
		getTypeInstruction(commitType, scope),
		getConventionInstruction(customMessageConvention))
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
