# Changelog

## [0.4.0] - 2025-12-08
### Added
- **`gitc file1 file2 ...`** — Auto-stage files and generate commit message instantly
  No more manual `git add` needed! Just run `gitc bot.py src/main.go` and get AI-powered commit.
- `--dry-run` (`-d`) flag — Preview exact prompt + config without any API call
  Perfect for debugging, testing prompts, or avoiding accidental costs.
- `--scope` (`-s`) flag — Full Conventional Commits scope support
  e.g. `gitc --scope auth login.go` → `feat(auth): add JWT login endpoint`
- `--temperature` flag — Control AI creativity (now default `0.7`, recommended `0.0` for teams)

### Changed
- Major refactor: `App.config` is now a **value type** (immutable)
  Prevents accidental mutation — config is copied on update, safe and predictable.
- Simplified and strengthened configuration validation
- Improved prompt engineering for better scope handling and consistency

### Fixed
- Fixed CLI flag name from `--maxLength` → `--max-length` (consistency with other flags)
- Fixed config mutation bugs in `ConfigAction`

---

## [0.3.0] - 2025-06-25
### Added
- Pretty-printed `git commit` command output with `-m` flags and line continuation (`\`).
- Automatic detection of single-line vs multi-line commit messages with clean display formatting.

### Changed
- Refactored commit message prompt for clarity, brevity, and better LLM compatibility.
- Improved formatting rules for both single-line and multi-line commit messages.
- Enforced strict summary/body structure with consistent guidelines.

---

## [0.2.0] - 2025-05-15
### Added
- Experimental support for Grok (xAI) and DeepSeek AI providers.
- New `--url` flag for custom API endpoints.
- Interactive mode for commit message preview and editing.

### Changed
- Updated README with improved provider status and clarity.
- Revised config structure to remove `open_ai` field.

### Fixed
- API key persistence issues in configuration.
- Improved validation for configuration settings.
- Fixed missing space between version number and summary in commit messages.
- Eliminated redundant spacing and formatting edge cases in output.

---

## [0.1.1] - 2025-04-01
- Initial release with OpenAI support.
