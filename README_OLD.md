# One CLI

> A unified command-line tool for freelancers working across multiple projects

One CLI provides a single interface for creating pull requests, managing authentication, opening tickets, and working with multiple browser profiles across different projects.

## Features

- 🚀 **Zero manual configuration** - Everything via interactive prompts
- 🔒 **Secure by default** - Uses OS keyring for credential storage (no plaintext!)
- 🌍 **Cross-platform** - Works on macOS, Linux, and Windows
- ⚡ **Fast** - Built with Go for speed
- 🎨 **Beautiful TUI** - Powered by Huh forms, Bubble Tea & Glamour rendering
- 🔧 **Multi-provider support** - GitHub, GitLab, Bitbucket
- 🎫 **Ticket integration** - Jira, Linear, GitHub Issues

## Installation

### From Source

```bash
git clone https://github.com/yourusername/one.git
cd one
go build -o one
sudo mv one /usr/local/bin/
```

### Using Go Install

```bash
go install github.com/yourusername/one@latest
```

## Quick Start

### 1. Initialize a Project

```bash
cd /path/to/your/project
one init
```

This will guide you through an interactive setup to configure:
- Project name and paths
- Git provider (GitHub/GitLab/Bitbucket)
- Browser and profile
- Ticket system (optional)
- PR templates (optional)

### 2. Start Working on a Task

```bash
one start PROJ-1234
```

This will:
- Checkout your base branch (e.g., `main`)
- Pull the latest changes
- Fetch the ticket title from your ticket system
- Create a new branch (e.g., `PROJ-1234-add-user-authentication`)
- Checkout the new branch

### 3. Create a Pull Request

```bash
one pr
```

This will:
- Push your branch to the remote
- Extract the ticket ID from your branch name
- Generate PR title and description from templates
- Create the PR via the provider API
- Open the PR in your configured browser profile

### 4. Open a Ticket

```bash
one ticket PROJ-1234
```

Opens the ticket in your configured browser.

## Configuration

Configurations are stored in `~/.config/one/`:

```
~/.config/one/
├── config.yml          # Optional global defaults
└── projects/
    ├── project-a.yml   # Project-specific configs
    └── project-b.yml
```

### Example Project Configuration

```yaml
version: 1

project:
  name: "Acme Corp Project"
  paths:
    - "/Users/john/Projects/acme-app"

git:
  provider: github
  remote: origin
  base_branch: main
  
  github:
    owner: acme-corp
    repo: main-app
    token_env: GITHUB_TOKEN

browser:
  type: chrome
  profile: "Work Profile"

ticket:
  system: jira
  base_url: https://acme.atlassian.net
  
  jira:
    board_id: ACME
    token_env: JIRA_TOKEN

templates:
  pr_title: "[{ticket_id}] {branch_name}"
  pr_body: |
    ## Changes
    - 
    
    ## Related Ticket
    {ticket_url}

branch_patterns:
  ticket_id: "^([A-Z]+-\\d+)"
```

## Commands

### `one help`
Display beautifully formatted help with examples and syntax highlighting

### `one docs`
View documentation with beautiful rendering

**Options:**
- `-s, --spec` - Show full technical specification
- `-e, --examples` - Show configuration examples

### `one init`
Initialize a new project configuration

**Options:**
- `-n, --name <NAME>` - Project name (skips prompt)

### `one start <TICKET-ID>`
Start working on a new task

**Options:**
- `-d, --description <TEXT>` - Custom branch description

### `one pr`
Create and open a pull request

**Options:**
- `-t, --title <TEXT>` - Custom PR title
- `-d, --description <TEXT>` - Custom PR description
- `--no-browser` - Skip opening browser

### `one ticket <TICKET-ID>`
Open a ticket in the browser

### `one config list`
List all configured projects (with beautiful markdown rendering)

### `one config show`
Show current project configuration (with syntax-highlighted YAML)

## Template Variables

Available variables for PR title and body templates:

- `{ticket_id}` - Extracted ticket ID (e.g., PROJ-1234)
- `{branch_name}` - Current branch name
- `{ticket_url}` - Full URL to ticket
- `{author}` - Git user name
- `{email}` - Git user email
- `{date}` - Current date (ISO 8601)
- `{base_branch}` - Target branch

## Authentication

One CLI stores credentials securely in your system's keyring:

- **macOS**: Keychain
- **Linux**: Secret Service (gnome-keyring, kwallet)
- **Windows**: Credential Manager

You can also use environment variables as a fallback:

```bash
export GITHUB_TOKEN="ghp_xxxxxxxxxxxx"
export JIRA_TOKEN="email@example.com:api_token"
```

## Browser Profiles

One CLI supports multiple browser profiles for keeping work and personal browsing separate:

**Supported browsers:**
- Chrome - Full profile support
- Firefox - Full profile support
- Safari - macOS only, no profile support

## Architecture

The codebase is organized into clean, focused modules:

```
one/
├── main.go                 # Entry point
├── cmd/                    # Command implementations
│   ├── root.go            # CLI setup
│   ├── init.go            # Interactive project setup
│   ├── start.go           # Start task command
│   ├── pr.go              # Create PR command
│   ├── ticket.go          # Open ticket command
│   └── config.go          # Config management
├── internal/
│   ├── config/            # Configuration system
│   ├── git/               # Git operations
│   ├── auth/              # Authentication & keyring
│   ├── browser/           # Browser launching
│   ├── api/               # API clients (GitHub, GitLab, Jira)
│   └── template/          # Template rendering
```

## Requirements

- **Go**: 1.21 or later
- **Git**: 2.0 or later
- **OS**: macOS 10.15+, Linux (Ubuntu 20.04+, Fedora 35+), Windows 10+

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details

## Credits

Built with:
- [Huh](https://github.com/charmbracelet/huh) - Interactive forms and prompts
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
- [Glamour](https://github.com/charmbracelet/glamour) - Markdown rendering
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [go-git](https://github.com/go-git/go-git) - Git operations
- [go-keyring](https://github.com/zalando/go-keyring) - Secure credential storage

## See Also

- [SPECIFICATION.md](SPECIFICATION.md) - Complete technical specification
