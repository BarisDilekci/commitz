# 🚀 Commitz

> Smart, interactive Git commit message generator that follows [Conventional Commits](https://www.conventionalcommits.org/) standards.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## ✨ Features

- 🎯 **Interactive Mode** - Beautiful CLI prompts to guide you through creating perfect commits
- 🤖 **Smart Suggestions** - Automatically analyzes your changes and suggests meaningful commit messages
- 🎨 **Emoji Support** - Add expressive emojis to your commits (optional)
- 📦 **Scope Detection** - Automatically extracts scope from branch names and project structure
- ✅ **Validation** - Ensures your commit messages follow best practices
- 🔍 **Dry Run** - Preview commits before creating them
- ⚡ **Fast & Lightweight** - Written in Go, no heavy dependencies

## 📦 Installation

### Using Go Install

```bash
go install github.com/barisdilekci/commitz@latest
```

### From Source

```bash
git clone https://github.com/barisdilekci/commitz.git
cd commitz
go build
go install
```

### Manual Download

Download the latest binary from [releases](https://github.com/barisdilekci/commitz/releases) and add it to your PATH.

## 🎬 Quick Start

```bash
# Stage your changes
git add .

# Run commitz in interactive mode
commitz -i -e

# Follow the prompts and you're done! 🎉
```

## 📖 Usage

### Interactive Mode (Recommended)

```bash
commitz -i
```

This will guide you through:
1. **Selecting commit type** (feat, fix, docs, etc.)
2. **Choosing scope** (from branch or project structure)
3. **Writing summary** (with smart suggestions)
4. **Adding description** (optional)
5. **Confirming and committing**

### Quick Mode

```bash
# Auto-detect everything
commitz

# Specify type
commitz -t feat -e

# Specify type and scope
commitz -t fix -s auth -e

# Dry run (preview only)
commitz -i -e -d
```

## 🎨 Commit Types

| Type | Emoji | Description |
|------|-------|-------------|
| `feat` | ✨ | A new feature |
| `fix` | 🐛 | A bug fix |
| `docs` | 📝 | Documentation changes |
| `style` | 💄 | Code style changes (formatting, etc.) |
| `refactor` | ♻️ | Code refactoring |
| `perf` | ⚡ | Performance improvements |
| `test` | ✅ | Adding or updating tests |
| `build` | 🔨 | Build system or dependency changes |
| `ci` | 👷 | CI/CD configuration changes |
| `chore` | 🧹 | Other changes (maintenance, etc.) |

## 🎯 Examples

### Interactive Mode
```bash
$ commitz -i -e

? Select commit type:
  ▸ ✨ feat - A new feature
    🐛 fix - A bug fix
    📝 docs - Documentation only changes
    ...

? Select scope:
  main (from branch)
  ▸ cmd
    pkg
    Skip (no scope)

? Commit summary (suggestion: add interactive mode): add user authentication

Suggested commit message:
  ✨ feat(auth): add user authentication

? Add detailed description? (y/N): y

Enter description:
- Implement JWT-based authentication
- Add login and registration endpoints
- Include password hashing with bcrypt

? Proceed with commit (y/N): y
✓ Commit successful! 🎉
```

### Quick Mode
```bash
$ commitz -t feat -s api -e

Suggested commit message:
  ✨ feat(api): add new endpoints

Proceed with commit? [Y/n]: y
✓ Commit successful! 🎉
```

## 🔧 Command-line Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--interactive` | `-i` | Enable interactive mode with prompts |
| `--type` | `-t` | Specify commit type (feat, fix, docs, etc.) |
| `--scope` | `-s` | Specify commit scope |
| `--emoji` | `-e` | Add emoji to commit message |
| `--dry-run` | `-d` | Preview commit without creating it |
| `--help` | `-h` | Show help message |

## 🎓 How It Works

### Smart Suggestions

Commitz analyzes your git diff to generate intelligent commit message suggestions:

- **File analysis**: Examines modified files and their paths
- **Content analysis**: Looks for keywords in added/modified code
- **Pattern recognition**: Identifies common patterns (tests, docs, fixes)
- **Context awareness**: Uses branch names and project structure

### Scope Detection

Automatically detects scope from:
1. **Branch names**: `feature/auth` → scope: `feature`
2. **Project structure**: Scans for common directories (cmd, pkg, api, etc.)
3. **Manual input**: You can always specify your own scope

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes using commitz! (`commitz -i -e`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 Conventional Commits

This tool follows the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

Example:
```
feat(auth): add JWT authentication

- Implement token generation and validation
- Add middleware for protected routes
- Include refresh token mechanism

Closes #123
```

## 🐛 Troubleshooting

### "No staged changes found"
Make sure you've staged your changes with `git add` before running commitz.

```bash
git add .
commitz -i
```

### "Not a git repository"
Run commitz from within a git repository.

```bash
cd your-project
commitz -i
```

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 💖 Acknowledgments

- [Conventional Commits](https://www.conventionalcommits.org/) for the specification
- [Cobra](https://github.com/spf13/cobra) for the CLI framework
- [promptui](https://github.com/manifoldco/promptui) for interactive prompts
- [color](https://github.com/fatih/color) for colored terminal output

## 🌟 Star History

If you find this project useful, please consider giving it a star! ⭐

---

<div align="center">
Made with ❤️ by <a href="https://github.com/barisdilekci">Barış Dilekci</a>
</div>