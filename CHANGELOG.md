# Changelog

All notable changes to One CLI will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.0] - 2025-10-06

### Added
- 🎨 **Huh Forms Integration** - Beautiful interactive forms for all input
- 🤖 **Auto-Detection** - Automatically detect Git remote, owner, repo, and branch
  - Detects GitHub, GitLab, Bitbucket from remote URL
  - Parses SSH and HTTPS remote URLs
  - Pre-fills configuration during `one init`
- 🔐 **GitHub OAuth Device Flow** - Secure authentication without manual tokens
  - Interactive OAuth flow during `one init`
  - Stores credentials in OS keyring
  - No need to copy/paste tokens manually
- 💅 **Glamour Rendering** - Beautiful markdown throughout
  - `one help` - Formatted help with examples
  - `one docs` - Documentation viewer
  - `one config list/show` - Pretty config display
- ⏱️ **Real-time Progress** - Bubble Tea progress displays
  - `one start` - Live branch creation feedback
  - `one pr` - PR creation progress
- ✅ **Form Validation** - Built-in validation for all inputs
- 🎯 **Smart Defaults** - Sensible defaults for everything

### Commands
- `one init` - Interactive project configuration with auto-detection
- `one start <ticket-id>` - Start working on a task
- `one pr` - Create and open a pull request  
- `one ticket <ticket-id>` - Open ticket in browser
- `one config list` - List all projects (with Glamour)
- `one config show` - Show current config (with Glamour)
- `one help` - Beautiful formatted help
- `one docs` - Documentation viewer
  - `--spec` flag to view full specification
  - `--examples` flag to view configuration examples

### Features
- 🌐 Multi-provider support (GitHub, GitLab, Bitbucket)
- 🎫 Multi-ticket-system support (Jira, Linear, GitHub Issues)
- 🖥️ Browser profile support (Chrome, Firefox, Safari)
- 📝 Template system for PR customization
- 🔒 Secure credential storage in OS keyring
- 🚀 Cross-platform (macOS, Linux, Windows)
- ⚡ Fast startup (< 100ms)
- 📦 Single binary deployment

### Technical
- Built with Go 1.21+
- Uses Charm ecosystem:
  - Huh v0.7.0 for forms
  - Bubble Tea v1.3.10 for TUI
  - Glamour v0.10.0 for markdown
  - Lipgloss v1.1.0 for styling
- Uses go-git v5.16.3 for Git operations
- Uses go-keyring v0.2.6 for secure storage
- Clean, modular architecture (~2,300 lines)

### Documentation
- README.md - Main documentation
- QUICKSTART.md - 5-minute getting started
- FEATURES.md - Complete feature showcase
- IMPLEMENTATION.md - Technical deep dive
- PROJECT_STATUS.md - Current status
- SPECIFICATION.md - Complete spec
- Examples in `examples/` directory

### Infrastructure
- Makefile for build automation
- Multi-platform build support
- Example configurations included
- MIT License

## [0.1.0] - Initial Concept

### Planned
- Basic CLI structure
- Configuration system
- Git operations
- API integrations

---

## Future Releases

### [0.3.0] - Planned
- OAuth flows for GitLab and Bitbucket
- Browser profile auto-discovery
- Shell completions (bash/zsh/fish)
- Draft PR support
- Unit tests

### [0.4.0] - Planned
- Linear API integration (not just URLs)
- Azure DevOps support
- PR review commands
- Multi-repo operations

### [1.0.0] - Planned
- Stable API
- Comprehensive test coverage
- CI/CD pipeline
- Homebrew formula
- Docker image
- Team features
