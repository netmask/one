<div align="center">

# 🚀 One CLI

**The unified command-line tool for freelancers**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-macOS%20%7C%20Linux%20%7C%20Windows-lightgrey?style=for-the-badge)]()

*Work across multiple projects with a single, beautiful interface for PRs, tickets, and more.*

[Installation](#-installation) • [Quick Start](#-quick-start) • [Features](#-features) • [Documentation](#-documentation)

</div>

---

## 💡 What is One CLI?

One CLI is designed for **freelancers and consultants** who juggle multiple clients and projects. Instead of remembering different workflows, credentials, and commands for each project, One provides a **unified interface** that adapts to your current project automatically.

### The Problem

```bash
# Project A (Client 1) - GitHub + Jira
cd ~/client1/project-a
git checkout -b PROJ-123-feature
# ... make changes ...
git push origin PROJ-123-feature
# Open browser, create PR manually, copy ticket link...

# Project B (Client 2) - GitLab + Linear
cd ~/client2/project-b  
git checkout -b LIN-456-feature
# ... different workflow ...
# Different credentials, different browser profile...
```

### The Solution

```bash
# Any project, anywhere
cd ~/any/project
one start PROJ-123    # Automatically: checkout, pull, create branch
# ... make changes ...
one pr                # Automatically: push, create PR, open in browser
```

One CLI **automatically detects** your project, uses the right credentials, opens the correct browser profile, and creates PRs in the right format.

---

## ✨ Features

<table>
<tr>
<td width="50%">

### 🎨 **Beautiful TUI**
- Interactive forms powered by [Huh](https://github.com/charmbracelet/huh)
- Real-time progress with [Bubble Tea](https://github.com/charmbracelet/bubbletea)  
- Markdown rendering via [Glamour](https://github.com/charmbracelet/glamour)

</td>
<td width="50%">

### 🔒 **Secure by Default**
- Credentials in OS keyring (never plaintext)
- OAuth device flow for GitHub
- Per-project credential isolation
- No manual token management

</td>
</tr>
<tr>
<td>

### 🌐 **Multi-Provider**
- **Git**: GitHub, GitLab, Bitbucket
- **Tickets**: Jira, Linear, GitHub Issues
- **Browsers**: Chrome, Firefox, Safari
  - Shows actual profile names & emails
  - No more "Profile 1" confusion!

</td>
<td>

### ⚡ **Smart & Fast**
- Auto-detects Git provider, owner, repo
- Auto-detects browser profiles with emails
- Hooks for linting, testing, automation
- < 100ms startup time
- Single binary, no dependencies
- Cross-platform (macOS/Linux/Windows)

</td>
</tr>
</table>

---

## 🎬 Demo

### Initialize a Project

```bash
$ one init
✓ Detected Git remote: github (acme-corp/main-app)

┌─ Project Name ────────────────────────────┐
│ Acme Corp                                 │
└───────────────────────────────────────────┘

┌─ Git Provider ────────────────────────────┐
│ > GitHub          (Detected: github)      │
│   GitLab                                  │
│   Bitbucket                               │
└───────────────────────────────────────────┘

✓ Configuration saved successfully!

🔐 Starting GitHub authentication...
Please visit: https://github.com/login/device
And enter code: ABCD-1234

Waiting for authorization......

✓ Successfully authenticated with GitHub!
```

### Start Working

```bash
$ one start PROJ-1234
Starting new task...

  Project: Acme Corp

  ✓ Checked out main
  ✓ Pulled from origin/main
  ✓ Fetching ticket info...
  ✓ Created branch PROJ-1234-add-user-authentication

✓ Ready to work on PROJ-1234-add-user-authentication!

  When done, run: one pr
```

### Create a PR

```bash
$ one pr
Creating pull request...

  Project: Acme Corp
  Branch: PROJ-1234-add-authentication
  Ticket ID: PROJ-1234

  ✓ Pushed to origin
  ✓ PR created: https://github.com/acme/app/pull/123

⚡ Running after_pr hooks...

  [1/1] Notify team
  ✓ Success (0.3s)

Opening in browser (Work Profile - john@acme.com)...

Done! 🚀
```

**Everything automated:** Git detection, OAuth, hooks, browser profiles!

---

## 📦 Installation

### Using Go

```bash
go install github.com/yourusername/one@latest
```

### From Source

```bash
git clone https://github.com/yourusername/one.git
cd one
make build
sudo make install
```

### Homebrew (Coming Soon)

```bash
brew install one-cli
```

---

## 🚀 Quick Start

### 1. Navigate to Your Project

```bash
cd /path/to/your/project
```

### 2. Initialize Configuration

```bash
one init
```

This interactive wizard will:
- ✅ **Auto-detect** your Git remote (GitHub/GitLab/Bitbucket)
- ✅ **Pre-fill** owner and repository from remote URL
- ✅ **Detect** your default branch (main/master/develop)
- ✅ **Authenticate** with GitHub via OAuth (optional)
- ✅ **Configure** browser profiles for work/personal separation
- ✅ **Set up** ticket system integration (Jira/Linear)

### 3. Start Working

```bash
# Start a new task
one start PROJ-1234

# Make your changes
git add .
git commit -m "Add feature"

# Create a PR (push + create + open in browser)
one pr
```

That's it! One CLI handles the rest. 🎉

---

## 🎯 Key Commands

| Command | Description |
|---------|-------------|
| `one init` | Interactive setup with Git auto-detection |
| `one start <ticket-id>` | Start working on a task (checkout, pull, create branch) |
| `one pr` | Create and open a pull request |
| `one ticket <ticket-id>` | Open ticket in browser |
| `one config list` | List all configured projects |
| `one config show` | Show current project config |
| `one profiles` | List browser profiles with emails |
| `one help` | Beautiful formatted help |
| `one docs` | View documentation |

---

## 🎨 Screenshots

### Interactive Setup
![Init Command](https://via.placeholder.com/800x400/282a36/bd93f9?text=Interactive+Setup+with+Huh+Forms)

### Real-time Progress
![Start Command](https://via.placeholder.com/800x400/282a36/50fa7b?text=Real-time+Progress+with+Bubble+Tea)

### Beautiful Documentation
![Help Command](https://via.placeholder.com/800x400/282a36/f8f8f2?text=Markdown+Rendering+with+Glamour)

---

## 📚 Documentation

- **[Quick Start Guide](QUICKSTART.md)** - Get started in 5 minutes
- **[Features](FEATURES.md)** - Complete feature list with examples
- **[Implementation](IMPLEMENTATION.md)** - Technical deep dive
- **[Specification](SPECIFICATION.md)** - Complete technical specification

---

## 🔧 Configuration

One CLI uses YAML configuration files stored in `~/.config/one/projects/`.

### Example Configuration

```yaml
version: 1

project:
  name: "Acme Corp"
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
  profile: "Work Profile"  # Keep work and personal separate!

ticket:
  system: jira
  base_url: https://acme.atlassian.net
  
  jira:
    board_id: ACME

templates:
  pr_title: "[{ticket_id}] {branch_name}"
  pr_body: |
    ## Changes
    - 
    
    ## Ticket
    {ticket_url}

branch_patterns:
  ticket_id: "^([A-Z]+-\\d+)"

# Hooks: Run commands before/after PR creation
hooks:
  before_pr:
    - name: "Lint code"
      command: "bundle exec rubocop"
      fail_on_error: true  # Stop if lint fails
    
    - name: "Run tests"
      command: "bundle exec rspec"
      fail_on_error: true  # Stop if tests fail
  
  after_pr:
    - name: "Notify team"
      command: 'curl -X POST $SLACK_WEBHOOK -d "{\"text\":\"PR created!\"}"'
      fail_on_error: false  # Continue even if notification fails
```

See [examples/](examples/) for more configurations and [HOOKS.md](HOOKS.md) for complete hooks documentation.

---

## 🤝 Why One CLI?

### For Freelancers & Consultants

- **Multiple Clients**: Different projects, different providers, one workflow
- **Context Switching**: Automatic project detection, no mental overhead
- **Professional**: Browser profiles keep work separate from personal
- **Secure**: OS keyring integration, never store tokens in plaintext

### For Teams

- **Onboarding**: New team members set up in minutes with `one init`
- **Consistency**: Everyone uses the same workflow across projects
- **Standards**: Enforce PR templates and branch naming conventions
- **Flexibility**: Each project can have its own configuration

---

## 🛠️ Technology Stack

Built with the [Charm](https://charm.sh/) ecosystem:

- **[Huh](https://github.com/charmbracelet/huh)** - Beautiful interactive forms
- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)** - Terminal UI framework
- **[Glamour](https://github.com/charmbracelet/glamour)** - Markdown rendering
- **[Lipgloss](https://github.com/charmbracelet/lipgloss)** - Terminal styling

Plus:
- **[Cobra](https://github.com/spf13/cobra)** - CLI framework
- **[go-git](https://github.com/go-git/go-git)** - Pure Go git implementation
- **[go-keyring](https://github.com/zalando/go-keyring)** - Cross-platform keyring

---

## 🌟 Features Deep Dive

### 🎯 Auto-Detection

One CLI is smart about your environment:

- **Git Remote**: Detects GitHub/GitLab/Bitbucket from `git remote`
- **Owner & Repo**: Parses SSH and HTTPS remote URLs
- **Default Branch**: Detects main/master/develop
- **Project Path**: Uses current directory automatically

### 🔐 GitHub OAuth Device Flow

Secure authentication without copying tokens:

```bash
$ one init
# ... configuration ...

┌─ Authenticate with GitHub now? ───────────┐
│ > Yes, authenticate                       │
│   Skip for now                            │
└───────────────────────────────────────────┘

🔐 Starting GitHub authentication...
Please visit: https://github.com/login/device
And enter code: ABCD-1234

Waiting for authorization......
✓ Successfully authenticated with GitHub!
```

Your token is securely stored in:
- **macOS**: Keychain
- **Linux**: Secret Service (gnome-keyring/kwallet)
- **Windows**: Credential Manager

### 🌐 Browser Profiles

Keep work and personal browsing separate:

```yaml
browser:
  type: chrome
  profile: "Work Profile"
```

One CLI opens PRs in your work browser profile with your work Google account logged in. No more confusion!

### 📝 Template System

Customize PR titles and descriptions per project:

```yaml
templates:
  pr_title: "[{ticket_id}] {branch_name}"
  pr_body: |
    ## 📝 Summary
    
    
    ## 🔗 Related
    - Ticket: {ticket_url}
    - Branch: {branch_name}
    - Author: {author}
    
    ## ✅ Checklist
    - [ ] Tests added
    - [ ] Docs updated
```

Available variables:
- `{ticket_id}` - Extracted from branch name
- `{branch_name}` - Current branch
- `{ticket_url}` - Generated ticket URL
- `{author}` - Git user name
- `{email}` - Git user email
- `{date}` - ISO 8601 date

---

## 🎓 Use Cases

### Freelance Developer with Multiple Clients

```bash
# Client A - GitHub + Jira
cd ~/clients/acme
one start ACME-123  # Uses Acme config automatically

# Client B - GitLab + Linear
cd ~/clients/beta
one start BET-456   # Uses Beta config automatically

# Personal Project - GitHub
cd ~/personal/side-project
one start FEAT-789  # Uses personal config automatically
```

### Consultant Switching Contexts

```bash
# Morning: Enterprise client (strict PR template)
cd ~/work/enterprise
one start ENT-1001
one pr  # Opens in work Chrome profile

# Afternoon: Startup client (simple workflow)
cd ~/work/startup
one start DEV-42
one pr  # Opens in different Chrome profile

# Evening: Open source (personal)
cd ~/oss/project
one start ISS-99
one pr  # Opens in personal browser
```

---

## 📊 Comparison

| Feature | One CLI | Manual Workflow | Other Tools |
|---------|---------|-----------------|-------------|
| Multi-project support | ✅ Automatic | ❌ Manual | ⚠️ Limited |
| Auto-detection | ✅ Git remote | ❌ None | ❌ None |
| GitHub OAuth | ✅ Device flow | ❌ Manual tokens | ❌ Manual |
| Browser profiles | ✅ Yes | ❌ Manual | ❌ No |
| Ticket integration | ✅ Multiple systems | ❌ Manual | ⚠️ Single |
| Cross-platform | ✅ Full | ✅ Yes | ⚠️ Limited |
| Secure storage | ✅ OS keyring | ❌ Plaintext | ⚠️ Varies |
| Setup time | ✅ < 2 minutes | ❌ Per project | ⚠️ Complex |

---

## 🤔 FAQ

<details>
<summary><strong>Do I need to configure each project manually?</strong></summary>

No! Run `one init` in your project directory and it will:
- Auto-detect your Git provider from the remote URL
- Pre-fill owner and repository information
- Detect your default branch
- Optionally authenticate via OAuth

You just confirm the detected values.
</details>

<details>
<summary><strong>How does One CLI know which project I'm working on?</strong></summary>

One CLI checks your current working directory against configured project paths. It uses prefix matching, so any subdirectory of a configured project will work.

```bash
cd ~/projects/acme/src/api  # Matches ~/projects/acme
one pr                      # Uses Acme config automatically
```
</details>

<details>
<summary><strong>Is it safe to store credentials?</strong></summary>

Yes! One CLI uses your operating system's native secure storage:
- **macOS**: Keychain Services
- **Linux**: Secret Service (gnome-keyring/kwallet)
- **Windows**: Credential Manager

Credentials are never stored in plaintext files. You can also use environment variables as a fallback.
</details>

<details>
<summary><strong>Can I use it with self-hosted Git servers?</strong></summary>

Yes! One CLI supports:
- GitHub Enterprise
- Self-hosted GitLab instances
- Bitbucket Server

Just specify the base URL in your configuration.
</details>

<details>
<summary><strong>Does it work with my ticket system?</strong></summary>

One CLI supports:
- **Jira Cloud** (with API integration)
- **Linear** (URL generation)
- **GitHub Issues** (URL generation)

More integrations coming soon! Or use custom ticket URLs.
</details>

---

## 🛣️ Roadmap

- [ ] **OAuth flows** for GitLab and Bitbucket
- [ ] **Browser profile auto-detection**
- [ ] **Shell completions** (bash/zsh/fish)
- [ ] **Draft PR** support
- [ ] **Linear API** integration (not just URLs)
- [ ] **Azure DevOps** support
- [ ] **PR review** commands
- [ ] **Multi-repo** operations
- [ ] **GitHub App** instead of OAuth app
- [ ] **Team features** (shared configs)

---

## 🤝 Contributing

Contributions are welcome! Whether it's:
- 🐛 Bug reports
- 💡 Feature requests
- 📝 Documentation improvements
- 🔧 Code contributions

Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

---

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

---

## 💖 Acknowledgments

Built with the amazing [Charm](https://charm.sh/) ecosystem:
- [Huh](https://github.com/charmbracelet/huh) - Interactive forms
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework  
- [Glamour](https://github.com/charmbracelet/glamour) - Markdown rendering
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Styling

---

<div align="center">

**Made with ❤️ for freelancers and consultants**

[⭐ Star on GitHub](https://github.com/yourusername/one) • [🐛 Report Bug](https://github.com/yourusername/one/issues) • [💡 Request Feature](https://github.com/yourusername/one/issues)

</div>
