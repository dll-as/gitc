package git

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

type Service interface {
	GetDiff(ctx context.Context) (string, error)
	StageAll(ctx context.Context) error
	StageFiles(ctx context.Context, files []string) error
}

// gitServiceImpl implements GitService
type gitServiceImpl struct {
	excludeFiles []string
}

// defaultExcludeFiles defines common files and folders to ignore in git diffs.
var defaultExcludeFiles = []string{
	"package-lock.json", "pnpm-lock.yaml", "yarn.lock", "*.lock",
	"*.min.js", "*.bundle.js",
	"node_modules/*", "dist/*", "build/*",
	"*.log", "*.bak", "*.swp",
}

// baseArgs is the constant portion of the git diff command, allocated once.
var baseArgs = []string{
	"diff",
	"--staged",
	"--diff-algorithm=minimal",
	"--unified=3",
	"--ignore-all-space",
	"--ignore-blank-lines",
	"--no-color",
	"--no-ext-diff",
	"--no-renames",
	"--ignore-submodules",
}

// NewGitService creates a new GitService
func NewService(excludeFiles []string) Service {
	return &gitServiceImpl{
		excludeFiles: append(defaultExcludeFiles, excludeFiles...),
	}
}

// GetDiff retrieves the git diff for staged changes
func (g *gitServiceImpl) GetDiff(ctx context.Context) (string, error) {
	return GetDiffStaged(ctx, g.excludeFiles)
}

// getGitRoot retrieves the root directory of the git repository
func getGitRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("not a git repository: %w", err)
	}

	return strings.TrimSpace(out.String()), nil
}

// getExcludeFileArgs converts exclude paths into git diff exclude args
func getExcludeFileArgs(excludeFiles []string) []string {
	args := make([]string, len(excludeFiles))
	for i, f := range excludeFiles {
		args[i] = fmt.Sprintf(":(exclude)%s", f)
	}
	return args
}

// processDiff applies cleanup to reduce unnecessary lines
func processDiff(diff string) string {
	// Pre-allocate roughly half the input size to avoid repeated grows.
	b := strings.Builder{}
	b.Grow(len(diff) / 2)

	inHunk := false
	for len(diff) > 0 {
		// Advance line by line without allocating a []string.
		var line string
		if i := strings.IndexByte(diff, '\n'); i >= 0 {
			line = diff[:i]
			diff = diff[i+1:]
		} else {
			line = diff
			diff = ""
		}

		if line == "" {
			continue
		}

		switch {
		case strings.HasPrefix(line, "index "),
			strings.HasPrefix(line, "--- "),
			strings.HasPrefix(line, "+++ "):
			// Metadata lines — not useful for the AI.
			continue

		case strings.HasPrefix(line, "@@"):
			inHunk = true
			// Keep only the trailing function/symbol context if present.
			// Format: "@@ -a,b +c,d @@ optional context"
			parts := strings.SplitN(line, "@@", 3)
			if len(parts) >= 3 && strings.TrimSpace(parts[2]) != "" {
				b.WriteString("@@ ")
				b.WriteString(strings.TrimSpace(parts[2]))
				b.WriteByte('\n')
			}
			// If there's no function context, skip the @@ line entirely.

		case inHunk && len(line) > 0 && line[0] == ' ':
			// Unchanged context line — skip to reduce token count.
			continue

		default:
			b.WriteString(line)
			b.WriteByte('\n')
		}
	}

	return strings.TrimSpace(b.String())
}

// GetDiffStaged retrieves the optimized git diff for staged changes with exclusions
func GetDiffStaged(ctx context.Context, extraExcludeFiles []string) (string, error) {
	rootPath, err := getGitRoot()
	if err != nil {
		return "", err
	}

	// Construct git diff command
	args := append(baseArgs, getExcludeFileArgs(extraExcludeFiles)...)

	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = rootPath
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("git diff timed out: %w", ctx.Err())
		}
		return "", fmt.Errorf("failed to get staged diff: %w", err)
	}

	rawDiff := strings.TrimSpace(out.String())
	if rawDiff == "" {
		return "", errors.New("no staged changes found")
	}

	// Process diff to remove unnecessary lines
	optimizedDiff := processDiff(rawDiff)
	if optimizedDiff == "" {
		return "", errors.New("no meaningful staged changes after processing")
	}

	return optimizedDiff, nil
}

// StageAll stages all changes in the working directory (equivalent to 'git add .').
func (s *gitServiceImpl) StageAll(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "git", "add", ".")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git add . failed: %s", string(output))
	}

	return nil
}

// StageFiles stages specific files (equivalent to 'git add file1 file2 ...')
func (s *gitServiceImpl) StageFiles(ctx context.Context, files []string) error {
	if len(files) == 0 {
		return nil
	}

	args := make([]string, 0, 2+len(files))
	args = append(args, "add", "--")
	args = append(args, files...)
	cmd := exec.CommandContext(ctx, "git", args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git add failed: %s", strings.TrimSpace(string(out)))
	}

	return nil
}
