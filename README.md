# ✨ gitc - AI-Powered Git Commit Messages

[![Go Reference](https://pkg.go.dev/badge/github.com/dll-as/gitc)](https://pkg.go.dev/github.com/dll-as/gitc)
[![Go Version](https://img.shields.io/github/go-mod/go-version/dll-as/gitc?logo=go)](go.mod)
[![Sourcegraph](https://sourcegraph.com/github.com/dll-as/gitc/-/badge.svg)](https://sourcegraph.com/github.com/dll-as/gitc?badge)
[![Discussions](https://img.shields.io/github/discussions/dll-as/gitc?color=58a6ff&label=Discussions&logo=github)](https://github.com/dll-as/gitc/discussions)
[![Downloads](https://img.shields.io/github/downloads/dll-as/gitc/total?color=blue)](https://github.com/dll-as/gitc/releases)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/dll-as/gitc)

<div align="center">
  <a href="#-installation">Installation</a> •
  <a href="#-features">Features</a> •
  <a href="#-configuration">Configuration</a> •
  <a href="#-basic-usage">Usage</a> •
  <a href="#-full-options">Full Options</a> •
  <a href="#-ai-providers">AI Providers</a>
</div>
<br>

> `gitc` is a fast, lightweight CLI tool that uses AI to generate clear, consistent, and standards-compliant commit messages — directly from your Git diffs. With built-in support for [Conventional Commits](https://www.conventionalcommits.org), [Gitmoji](https://gitmoji.dev), and fully customizable rules, `gitc` helps you and your team write better commits, faster

# 🚀 Features
`gitc` is a lightweight CLI tool that leverages AI to craft clear, standards-compliant Git commit messages from your diffs. Supporting [Conventional Commits](https://www.conventionalcommits.org), [Gitmoji](https://gitmoji.dev), and custom rules, it saves time and boosts commit quality for you and your team.

- 🧠 **AI-Powered Commits**
  - Generates context-aware commit messages using OpenAI, Grok (xAI), DeepSeek, or Ollama (local).
  - Supports multiple languages (e.g., English, Persian, Russian) for global teams.
  - Extensible for future AI providers.

- 🏠 **Local AI Support**
  - Run Ollama locally for private, offline commit generation - no API key required!
  - Connect to any Ollama instance via custom URLs.

- 📝 **Standards & Customization**
  - Follows [Conventional Commits](https://www.conventionalcommits.org) (`feat`, `fix`, `docs`, etc.) for semantic versioning.
  - Adds [Gitmoji](https://gitmoji.dev) emojis for visual flair (e.g., ✨, 🚑).
  - Customizable prefixes (e.g., JIRA IDs) via JSON.

- 🔧 **Git Integration**
  - Processes staged Git diffs, ignoring irrelevant files (`node_modules/*`, `*.lock`).
  - Configurable file exclusions for focused commits.

- ⚙️ **Flexible Configuration**
  - Supports CLI flags, environment variables, and `~/.gitc/config.json`.
  - Includes proxy support, adjustable timeouts, and redirect limits.
  - No API key required for local Ollama provider.

- ⚡️ **Performance & Reliability**
  - Fast JSON parsing with [sonic](https://github.com/bytedance/sonic) and HTTP requests with [fasthttp](https://github.com/valyala/fasthttp).
  - Robust error handling for reliable operation.
  - Unified interface supporting both cloud and local AI providers.

- 🧪 Debug & Dry-Run
  - Preview prompts and configs without API calls — perfect for tuning without burning tokens.

## 📦 Installation
### Prerequisites:
  - Go: Version **1.18** or higher (required for building from source).
  - Git: Required for retrieving staged changes.
  - API Key: Required for cloud AI providers (OpenAI, Grok, DeepSeek).
  - Ollama: Optional for local AI (install from [ollama.ai](https://ollama.ai))

#### Quick Install:
  ```bash
  go install github.com/dll-as/gitc@latest
  ```

### Manual Install
  1. Download binary from [releases](https://github.com/dll-as/gitc/releases)
  2. `chmod +x gitc`
  3. Move to `/usr/local/bin`

### Verify Installation
  After installation, verify the tool is installed correctly and check its version:
  ```bash
  gitc --version
  ```

# 💻 Basic Usage
```bash
# 1. Stage your changes
git add . # or gitc -a

# 2. Generate perfect commit message
gitc

# Stage specific files and generate
gitc bot.py
gitc src/utils.go main.go

# Use local Ollama for private commits
gitc --provider ollama

# Pro Tip: Add emojis and specify language
gitc --emoji --lang fa

# Custom commit type
gitc --commit-type fix

# Debug mode: See what prompt would be sent without API cost
gitc --dry-run
```

## Environment Variables
```bash
export AI_API_KEY="sk-your-key-here"  # For cloud providers
export GITC_LANGUAGE="fa"
export GITC_MODEL="gpt-4"
export GITC_PROVIDER="ollama"  # Use local Ollama by default
```

# ⚙️ Configuration
Config File (`~/.gitc/config.json`) :
```json
{
  "provider": "openai",
  "max_length": 200,
  "temperature": 0.7,
  "proxy": "",
  "language": "en",
  "timeout": 10,
  "commit_type": "",
  "custom_convention": "",
  "use_gitmoji": false,
  "max_redirects": 5
}
```

### Update Configuration
```bash
gitc config --provider ollama --model llama3.2
gitc config --api-key "sk-your-key-here" --model "gpt-4o-mini" --lang en
```


# 📚 Full Options
The following CLI flags are available for the `ai-commit` command and its `config` subcommand. All flags can also be set via environment variables or the `~/.gitc/config.json` file.

| Flag | Alias | Description | Default | Environment Variable | Example |
|------|-------|-------------|---------|----------------------|---------|
| `--all` | `-a` | Stage all changes before generating commit message (equivalent to `git add .`) | `false` | `GITC_STAGE_ALL` | `-all` or `-a`
| `--provider` | - | AI provider to use (e.g., `openai`, `grok`, `deepseek`, `ollama`) | `openai` | `AI_PROVIDER` | `--provider ollama` |
| `--url` | `-u` | Custom API URL for the AI provider | Provider-specific | `GITC_API_URL` | `--url http://localhost:11434/api/chat`
| `--model` | - | AI model for commit message generation | Provider-specific | `GITC_MODEL` | `--model llama3.2` |
| `--lang` | - | Language for commit messages (e.g., `en`, `fa`, `ru`) | `en` | `GITC_LANGUAGE` | `--lang fa` |
| `--timeout` | - | Request timeout in seconds | `10` | - | `--timeout 15` |
| `--max-length` | - | Maximum length of the commit message | `200` | - | `--max-length 150` |
| `--temperature` | - | Control AI creativity (0.0 = fully deterministic, 1.0 = very creative) | `0.7` | - | `--temperature 0.8`
| `--api-key` | `-k` | API key for the AI provider (not required for Ollama) | - | `AI_API_KEY` | `--api-key sk-xxx` |
| `--proxy` | `-p` | Proxy URL for API requests | - | `GITC_PROXY` | `--proxy http://proxy.example.com:8080` |
| `--commit-type` | `-t` | Commit type for Conventional Commits (e.g., `feat`, `fix`) | - | `GITC_COMMIT_TYPE` | `--commit-type feat` |
| `--scope`      | `-s` | Add scope to the commit type (e.g. `auth`, `ui`, `db`) — works with or without `--commit-type` | - | - | `--scope auth` or `-s ui` |
| `--custom-convention` | `-C` | Custom commit message convention (JSON format) | - | `GITC_CUSTOM_CONVENTION` | `--custom-convention '{"prefix": "JIRA-123"}'` |
| `--emoji` | `-g` | Add Gitmoji to the commit message | `false` | `GITC_GITMOJI` | `--emoji` |
| `--no-emoji` | - | Disables Gitmoji in commit messages (overrides `--emoji` and config file) | `false` | - | `--no-emoji`
| `--dry-run` | -d | Preview the exact prompt and config sent to AI without making an API request (great for debugging prompts and avoiding costs) | `false` | `GITC_DRY_RUN` | --dry-run or -d
| `--max-redirects` | `-r` | Maximum number of HTTP redirects | `5` | `GITC_MAX_REDIRECTS` | `--max-redirects 10` |
| `--config` | `-c` | Path to the configuration file | `~/.gitc/config.json` | `GITC_CONFIG_PATH` | `--config ./my-config.json` |

> [!NOTE]
> - **Ollama Note**: No API key required when using `--provider ollama`
> - **Default URLs**:
>   - OpenAI: `https://api.openai.com/v1/chat/completions`
>   - Grok: `https://api.x.ai/v1/chat/completions`
>   - DeepSeek: `https://api.deepseek.com/v1/chat/completions`
>   - Ollama: `http://localhost:11434/api/chat`
> - **Default Models**:
>   - OpenAI: `gpt-4o-mini`
>   - Ollama: `llama3.2`
> - Flags for the `config` subcommand are similar but exclude defaults, as they override the config file.
> - **Flags** > **Environment Variables** > **Config File** — This is the order of precedence when multiple settings are provided.
> - The `--custom-convention` flag expects a JSON string with a `prefix` field (e.g., `{"prefix": "JIRA-123"}`).
> - The `--version` flag displays the current tool version (e.g., `0.5.0`) and can be used to verify installation.
> - The `--all` flag (alias `-a`) stages all changes in the working directory before generating the commit message, streamlining the workflow. For example, `gitc -a --emoji` stages all changes and generates a commit message with Gitmoji.
> - Environment variables take precedence over config file settings but are overridden by CLI flags.
> - You can reset all configuration values to their defaults by using `gitc reset-config`.

## 🤖 AI Providers
`gitc` supports both cloud-based and local AI providers, giving you flexibility and privacy options.

| Provider | Supported Models | Required Configuration | API Key | Status |
| --- | --- | --- | --- | --- |
| **OpenAI** | `gpt-4o`, `gpt-4o-mini`, `gpt-3.5-turbo` | `model`, `url` (optional) | ✅ Required | ✅ Stable |
| **Grok (xAI)** | `grok-3`, `grok-2` | `model`, `url` | ✅ Required | 🧪 Experimental |
| **DeepSeek** | `deepseek-rag`, `deepseek-coder` | `model`, `url` | ✅ Required | 🧪 Experimental |
| **Ollama (Local)** | Any Ollama model (`llama3.2`, `codellama`, `mistral`, etc.) | `model`, `url` (optional) | ❌ Not Required | ✅ Stable |
> **Local AI Benefits**: Ollama provides complete privacy, works offline, and has no API costs. Perfect for sensitive projects or when internet access is limited.

## 🤝 Contributing

We welcome contributions! Please check out the [contributing guide](CONTRIBUTING.md) before making a PR.

## ⭐️ Star History
[![Star History Chart](https://api.star-history.com/svg?repos=dll-as/gitc&type=Date)](https://www.star-history.com/#dll-as/gitc&Date)