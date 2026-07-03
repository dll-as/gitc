package core

import (
	"fmt"
	"strings"
)

func (g *CommitGenerator) buildPrompt(diff string) string {
	var b strings.Builder

	b.WriteString("Write a Git commit message for this diff.\n\n")
	b.WriteString(fmt.Sprintf("Language: %s\n", g.config.Prompt.Language))

	if g.config.Prompt.CommitType != "" {
		b.WriteString(fmt.Sprintf("Type: %s\n", g.config.Prompt.CommitType))
	}

	if g.config.Prompt.Scope != "" {
		b.WriteString(fmt.Sprintf("Scope: %s\n", g.config.Prompt.Scope))
	}

	if g.config.Prompt.UseGitmoji {
		b.WriteString("Prefix the subject with a Gitmoji.\n")
	}

	b.WriteString(`
Format:
	Line 1: <type>: <summary> (≤50 chars)
	Line 2: (blank)
	Line 3+: (optional) details (≤100 chars per line)

Rules: raw text only, no markdown, no explanation.
Examples:
	feat: add JWT middleware\nAdd access token check to protected routes.
	fix: prevent crash on nil DB config\nAdd nil check before DB usage.
`)

	b.WriteString("\nDiff:\n")
	b.WriteString(diff)

	return b.String()
}

// FormatGitCommand converts a commit message into a git commit command.
func (g *CommitGenerator) FormatGitCommand(msg string) string {
	msg = strings.TrimSpace(msg)
	if msg == "" {
		return `git commit -m ""`
	}

	msg = strings.ReplaceAll(msg, "\\n", "\n")
	lines := strings.Split(msg, "\n")

	// Drop trailing blank lines
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}

	if len(lines) == 0 {
		return `git commit -m ""`
	}

	if len(lines) == 1 {
		return fmt.Sprintf("git commit -m %q", lines[0])
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("git commit -m %q", lines[0]))
	for _, line := range lines[1:] {
		b.WriteString(fmt.Sprintf(" \\\n    -m %q", line))
	}
	return b.String()
}
