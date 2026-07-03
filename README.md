# ✨ gitc — AI-Powered Git Commit Messages

[![Go Reference](https://pkg.go.dev/badge/github.com/dll-as/gitc)](https://pkg.go.dev/github.com/dll-as/gitc)
[![Go Version](https://img.shields.io/github/go-mod/go-version/dll-as/gitc?logo=go)](go.mod)
[![Sourcegraph](https://sourcegraph.com/github.com/dll-as/gitc/-/badge.svg)](https://sourcegraph.com/github.com/dll-as/gitc?badge)
[![Discussions](https://img.shields.io/github/discussions/dll-as/gitc?color=58a6ff&label=Discussions&logo=github)](https://github.com/dll-as/gitc/discussions)
[![Downloads](https://img.shields.io/github/downloads/dll-as/gitc/total?color=blue)](https://github.com/dll-as/gitc/releases)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/dll-as/gitc)
[![Zero Dependencies](https://img.shields.io/badge/dependencies-zero-brightgreen)](go.mod)

<div align="center">
  <a href="#-installation">Installation</a> •
  <a href="#-features">Features</a> •
  <a href="#-configuration">Configuration</a> •
  <a href="#-basic-usage">Usage</a> •
  <a href="#-full-options">Full Options</a> •
  <a href="#-ai-providers">AI Providers</a>
</div>
<br>

> `gitc` is a fast, lightweight CLI tool that uses AI to generate clear, consistent, and standards-compliant commit messages — directly from your Git diffs. With built-in support for [Conventional Commits](https://www.conventionalcommits.org), [Gitmoji](https://gitmoji.dev), and fully customizable rules, `gitc` helps you and your team write better commits, faster.

## 🚀 Features

- 🧠 **AI-Powered Commits** — Generates context-aware commit messages via OpenAI-compatible APIs, Ollama (local), or Anthropic.
- 🏠 **Local AI Support** — Run Ollama locally for private, offline commit generation with no API key required.
- 📝 **Standards & Customization** — Follows [Conventional Commits](https://www.conventionalcommits.org) and optionally adds [Gitmoji](https://gitmoji.dev) emojis.
- 🔧 **Smart Git Integration** — Processes only staged diffs and automatically filters out noise (`node_modules/*`, `*.lock`, etc.).
- ⚙️ **Flexible Configuration** — CLI flags, environment variables, and `~/.gitc/config.json` all supported, with a clear precedence order.
- 🧪 **Dry-Run Mode** — Preview the prompt that would be sent to AI without making any API call.


## 📦 Installation
### Prerequisites:
  - Go: Version **1.25** or higher (required for building from source).
  - Git: Required for retrieving staged changes.
  - No external runtime dependencies — a single static binary
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
  "backend": {
    "backend": "openai-compatible",
    "base_url": "https://api.openai.com/v1",
    "api_key": "",
    "model": "gpt-4o-mini",
    "timeout": 30000000000,
    "max_retries": 3,
    "proxy": ""
  },
  "prompt": {
    "language": "en",
    "convention": "conventional",
    "use_gitmoji": false,
    "max_tokens": 512,
    "temperature": 0.3,
    "top_p": 1
  },
  "git": {
    "max_diff_size": 100000,
    "exclude_patterns": [
      "package-lock.json",
      "pnpm-lock.yaml",
      "yarn.lock",
      "*.lock",
      "*.min.js",
      "*.bundle.js",
      "node_modules/*",
      "dist/*",
      "build/*",
      "*.log",
      "*.bak",
      "*.swp"
    ]
  }
}
```

### Update Configuration
```bash
gitc config --backend ollama --model llama3.2
gitc config --api-key "sk-..." --model "gpt-4o-mini" --base-url https://api.openai.com/v1
gitc config --backend openai-compatible --base-url https://openrouter.ai/api/v1 --api-key sk-or-v1-... --model openrouter/free
```

### Reset to defaults

```bash
gitc reset-config
```

## 🤖 AI Providers

| Backend value | Description | API Key | Notes |
|---|---|---|---|
| `openai-compatible` | OpenAI or any OpenAI-compatible API (OpenRouter, DeepSeek, Grok, etc.) | ✅ Required | Default backend |
| `ollama` | Local Ollama instance | ❌ Not required | Works offline |
| `anthropic` | Anthropic Claude (via OpenAI-compatible adapter) | ✅ Required | — |

### Provider examples

**OpenAI:**
```bash
gitc config --backend openai-compatible \
  --base-url https://api.openai.com/v1 \
  --model gpt-4o-mini \
  --api-key sk-...
```

**OpenRouter (free models):**
```bash
gitc config --backend openai-compatible \
  --base-url https://openrouter.ai/api/v1 \
  --model openrouter/free \
  --api-key sk-or-v1-...
```

**Ollama (local):**
```bash
gitc config --backend ollama \
  --base-url http://localhost:11434 \
  --model llama3.2
```

**Anthropic:**
```bash
gitc config --backend anthropic \
  --base-url https://api.anthropic.com/v1 \
  --model claude-3-5-haiku-20241022 \
  --api-key sk-ant-...
```

## 🤝 Contributing
Contributions are welcome! Please check out the [contributing guide](CONTRIBUTING.md) before making a PR.