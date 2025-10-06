# Contributing to One CLI

First off, thank you for considering contributing to One CLI! 🎉

## How Can I Contribute?

### 🐛 Reporting Bugs

Before creating bug reports, please check existing issues. When creating a bug report, include:

- **Clear title and description**
- **Steps to reproduce**
- **Expected vs actual behavior**
- **Environment details** (OS, Go version, One CLI version)
- **Configuration** (sanitized, without tokens)
- **Error messages** (full output if possible)

### 💡 Suggesting Features

Feature suggestions are welcome! Please:

- **Check existing issues** first
- **Describe the feature** clearly
- **Explain the use case**
- **Consider alternatives** you've thought about
- **Be open to discussion**

### 🔧 Pull Requests

1. **Fork the repo** and create your branch from `main`
2. **Make your changes**
3. **Test thoroughly** on your platform
4. **Update documentation** if needed
5. **Write clear commit messages**
6. **Submit the PR** with a clear description

## Development Setup

### Prerequisites

- Go 1.21 or higher
- Git
- Make (optional, but recommended)

### Getting Started

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/one.git
cd one

# Build
make build

# Or manually
go build -o one

# Test
./one --version
```

### Project Structure

```
one/
├── main.go                # Entry point
├── cmd/                   # Command implementations
│   ├── init.go           # Interactive setup (Huh)
│   ├── start.go          # Start task (Huh + Bubble Tea)
│   ├── pr.go             # Create PR (Bubble Tea)
│   └── ...
├── internal/             # Internal packages
│   ├── config/          # Configuration system
│   ├── git/             # Git operations + detection
│   ├── auth/            # Authentication + OAuth
│   ├── browser/         # Browser launching
│   ├── api/             # API clients
│   └── template/        # Template rendering
└── examples/            # Configuration examples
```

### Code Style

- **Go fmt** - Run `go fmt ./...` before committing
- **Clear names** - Use descriptive variable and function names
- **Comments** - Document exported functions and complex logic
- **Errors** - Return errors, don't panic (except in main)
- **Testing** - Add tests for new features (when test suite exists)

### Commit Messages

Use clear, descriptive commit messages:

```
Good:
  Add GitHub OAuth device flow authentication
  Fix branch name sanitization for special characters
  Update README with auto-detection examples

Bad:
  Fix bug
  Update code
  Changes
```

## Areas for Contribution

### 🚀 High Priority

- **OAuth flows** for GitLab and Bitbucket
- **Unit tests** - Test coverage currently at 0%
- **Shell completions** - Bash, Zsh, Fish
- **Browser profile detection** - Auto-discover profiles
- **Error handling** - Improve error messages

### 🎯 Medium Priority

- **Linear API integration** - Full API, not just URL generation
- **Azure DevOps support** - New Git provider
- **Draft PR support** - Create draft PRs
- **PR review commands** - `one review` to review PRs
- **Stash support** - Implement stash functionality in `one start`

### 💎 Nice to Have

- **Multi-repo support** - Operations across multiple repos
- **Time tracking** - Integration with time tracking tools
- **Slack/Discord notifications** - Notify on PR creation
- **GitHub App** - Use GitHub App instead of OAuth
- **Team features** - Shared configurations

### 📚 Documentation

- **Tutorial videos** - Screen recordings of workflows
- **Blog posts** - Use case examples
- **Translations** - Translate README to other languages
- **API documentation** - Document internal packages

## Testing

Currently, One CLI has no automated tests. This is a great area for contribution!

### What to Test

1. **Unit Tests**
   - Configuration parsing
   - Git remote detection
   - Branch name sanitization
   - Template rendering
   - Path matching

2. **Integration Tests**
   - Git operations
   - API clients (mocked)
   - Keyring operations
   - Browser launching

3. **E2E Tests**
   - Full workflows
   - Cross-platform compatibility

### Running Tests (when available)

```bash
make test
# or
go test ./...
```

## Building for Multiple Platforms

```bash
# Build for all platforms
make build-all

# This creates:
# - one-darwin-amd64 (macOS Intel)
# - one-darwin-arm64 (macOS Apple Silicon)
# - one-linux-amd64 (Linux)
# - one-linux-arm64 (Linux ARM)
# - one-windows-amd64.exe (Windows)
```

## Adding a New Command

1. **Create the command file** in `cmd/`
2. **Use Huh for input** if interactive
3. **Use Bubble Tea for progress** if long-running
4. **Register in** `cmd/root.go`
5. **Update documentation**

Example:

```go
// cmd/mycommand.go
package cmd

import (
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var myCmd = &cobra.Command{
	Use:   "my-command",
	Short: "Description",
	RunE:  runMyCommand,
}

func init() {
	rootCmd.AddCommand(myCmd)
}

func runMyCommand(cmd *cobra.Command, args []string) error {
	// Use Huh for forms
	var input string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter something").
				Value(&input),
		),
	)
	
	if err := form.Run(); err != nil {
		return err
	}
	
	// Implementation...
	return nil
}
```

## Adding a New Git Provider

1. **Add provider detection** in `internal/git/detect.go`
2. **Add config types** in `internal/config/types.go`
3. **Add API client** in `internal/api/provider.go`
4. **Update init command** in `cmd/init.go`
5. **Update PR command** in `cmd/pr.go`
6. **Add example config** in `examples/`

## Adding a New Ticket System

1. **Add to config types** in `internal/config/types.go`
2. **Add URL generation** in `internal/template/render.go`
3. **Add API client** (if fetching titles) in `internal/api/`
4. **Update init command** in `cmd/init.go`
5. **Update start command** in `cmd/start.go`

## Code of Conduct

### Our Pledge

We pledge to make participation in our project a harassment-free experience for everyone.

### Our Standards

**Positive behavior:**
- Using welcoming and inclusive language
- Being respectful of differing viewpoints
- Gracefully accepting constructive criticism
- Focusing on what is best for the community

**Unacceptable behavior:**
- Trolling, insulting/derogatory comments
- Public or private harassment
- Publishing others' private information
- Other conduct which could reasonably be considered inappropriate

### Enforcement

Instances of abusive, harassing, or otherwise unacceptable behavior may be reported to the project maintainers. All complaints will be reviewed and investigated.

## Questions?

- **GitHub Issues** - For bugs and features
- **GitHub Discussions** - For questions and ideas
- **Email** - For private concerns

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

**Thank you for contributing to One CLI!** 🚀
