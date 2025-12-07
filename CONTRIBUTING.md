# 🤝 Contributing to `gitc`

Thank you for your interest in contributing to **`gitc`**, an AI-powered CLI tool for generating meaningful commit messages from your Git diffs. Your contributions help make `gitc` better for everyone!

We welcome all types of contributions, including bug fixes, feature additions, documentation improvements, and tests. This guide outlines how you can get involved.

## 📌 Table of Contents
- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Contributions](#making-contributions)
- [Commit Message Guidelines](#commit-message-guidelines)
- [Pull Request Process](#pull-request-process)
- [Feature Suggestions & Bug Reports](#feature-suggestions--bug-reports)
- [Community](#community)

## 📜 Code of Conduct
We are committed to fostering a welcoming and inclusive community. Please adhere to our [Contributor Covenant Code of Conduct](https://www.contributor-covenant.org/version/2/1/code_of_conduct/). Treat everyone with respect and kindness.

## ⚙️ Getting Started
To start contributing, follow these steps:

1. **Fork the repository** on GitHub.
2. **Clone your fork** to your local machine:
   ```bash
   git clone https://github.com/<your-username>/gitc.git
   cd gitc
   ```
3. **Create a new branch** for your changes:
   ```bash
   git checkout -b feature/<your-feature-name>
   ```

## 🛠 Development Setup
To set up the development environment, ensure you have the following prerequisites:
- **Go** ≥ 1.20
- **Git**
- (Optional) An **OpenAI API Key** for testing AI-powered features

### Steps:
1. **Install dependencies**:
   ```bash
   go mod tidy
   ```
2. **Build the project**:
   ```bash
   go build -o gitc ./cmd/gitc
   ```
3. **Run the tool**:
   ```bash
   ./gitc --help
   ```
4. **Run tests**:
   ```bash
   go test ./...
   ```

## ✍️ Making Contributions
We value contributions in many forms, including:
- 🐛 **Fixing bugs**
- 📄 **Improving documentation**
- 🚀 **Adding new features**
- ✅ **Writing or improving tests**
- 💡 **Suggesting ideas** via discussions

If you're unsure where to start, check out our [GitHub Discussions](https://github.com/dll-as/gitc/discussions) or [Issues](https://github.com/dll-as/gitc/issues) for inspiration.

## 🧾 Commit Message Guidelines
We use [Conventional Commits](https://www.conventionalcommits.org) to maintain a clear and semantic commit history. This ensures our changelog is easy to read and understand.

### Examples:
```bash
feat(config): add support for Gemini AI provider
fix(cli): resolve crash when config file is missing
docs(readme): clarify installation steps
```

> **Pro Tip**: You can use `gitc` itself to generate commit messages! For example:
> ```bash
> gitc -a --commit-type feat
> ```

### Commit Types:
- `feat`: New features or enhancements
- `fix`: Bug fixes
- `docs`: Documentation updates
- `test`: Adding or improving tests
- `refactor`: Code refactoring without changing functionality
- `chore`: Miscellaneous tasks (e.g., updating dependencies)

## 🚀 Pull Request Process
To submit your changes:
1. Ensure your pull request (PR) targets the **`main`** branch.
2. Verify that all tests pass:
   ```bash
   go test ./...
   ```
3. Write a clear and descriptive PR title and description.
4. Reference related issues (e.g., `Closes #42` or `Fixes #42`).
5. Add relevant labels (e.g., `bug`, `feature`, `documentation`).
6. Submit your PR and respond to any feedback during the code review process.

## 💡 Feature Suggestions & Bug Reports
- **Report a bug**: Open an issue on our [Issues page](https://github.com/dll-as/gitc/issues).
- **Suggest a feature**: Start a discussion on our [Discussions page](https://github.com/dll-as/gitc/discussions).
- **New to contributing?** Look for issues labeled [`good first issue`](https://github.com/dll-as/gitc/labels/good%20first%20issue) to get started.

## 🫶 Community
Join our community to connect with other contributors:
- Participate in [GitHub Discussions](https://github.com/dll-as/gitc/discussions) to share ideas or ask questions.
- Stay updated by starring the repository and following along!

## 🙌 Thank You!
Your contributions—whether code, documentation, or ideas—make `gitc` better for everyone. We're excited to have you as part of our community! 🚀