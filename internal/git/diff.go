package git

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// GitService defines the interface for git operations
type GitService interface {
	GetDiff(ctx context.Context) (string, error)
	StageAll(ctx context.Context) error
	StageFiles(ctx context.Context, files []string) error
}

// gitServiceImpl implements GitService
type gitServiceImpl struct {
	excludeFiles []string
}

// NewGitService creates a new GitService
func NewGitService(excludeFiles ...string) GitService {
	return &gitServiceImpl{
		excludeFiles: append(defaultExcludeFiles, excludeFiles...),
	}
}

// defaultExcludeFiles defines common files and folders to ignore in git diffs
var defaultExcludeFiles = []string{
	"package-lock.json", "pnpm-lock.yaml", "yarn.lock", "*.lock",
	"*.min.js", "*.bundle.js",
	"node_modules/*", "dist/*", "build/*",
	"*.log", "*.bak", "*.swp",
	// "*.png", "*.jpg", "*.jpeg", "*.gif", "*.svg", "*.ico",
	// "*.woff", "*.woff2", "*.ttf", "*.eot",
	// "*.pdf", "*.zip", "*.gz",
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
	var builder strings.Builder
	lines := strings.Split(diff, "\n")
	inHunk := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		switch {
		case strings.HasPrefix(trimmed, "index "),
			strings.HasPrefix(trimmed, "--- "),
			strings.HasPrefix(trimmed, "+++ "):
			continue
		case strings.HasPrefix(trimmed, "@@"):
			inHunk = true
			parts := strings.SplitN(trimmed, "@@", 3)
			if len(parts) >= 3 {
				builder.WriteString("@@" + strings.TrimSpace(parts[2]) + "\n")
			}
		case strings.HasPrefix(trimmed, " ") && inHunk:
			continue
		default:
			builder.WriteString(trimmed + "\n")
		}
	}

	return strings.TrimSpace(builder.String())
}

// GetDiffStaged retrieves the optimized git diff for staged changes with exclusions
func GetDiffStaged(ctx context.Context, extraExcludeFiles []string) (string, error) {
	rootPath, err := getGitRoot()
	if err != nil {
		return "", err
	}

	// Construct git diff command
	args := []string{
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
	args = append(args, getExcludeFileArgs(extraExcludeFiles)...)

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

	cmd := exec.CommandContext(ctx, "git", append([]string{"add"}, files...)...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git add failed: %s", string(output))
	}

	return nil
}
