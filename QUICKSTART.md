# One CLI - Quick Start Guide

Get started with One CLI in 5 minutes!

## Installation

### Build from source

```bash
git clone https://github.com/yourusername/one.git
cd one
make build
sudo make install
```

Or manually:

```bash
go build -o one
sudo mv one /usr/local/bin/
```

## First Time Setup

### 1. Navigate to Your Project

```bash
cd /path/to/your/project
```

### 2. Initialize Configuration

```bash
one init
```

You'll be guided through:
- Project name
- Git provider (GitHub/GitLab/Bitbucket)
- Browser preferences
- Ticket system (optional)

The configuration is saved to `~/.config/one/projects/your-project.yml`

### 3. Set Up Authentication (Optional)

If using keyring authentication:

```bash
one login github
```

Or use environment variables:

```bash
export GITHUB_TOKEN="ghp_xxxxxxxxxxxx"
export JIRA_TOKEN="email@example.com:api_token"
```

## Daily Workflow

### Start a New Task

```bash
one start PROJ-1234
```

This will:
1. Checkout your base branch (main/develop)
2. Pull latest changes
3. Fetch the ticket title from Jira/Linear
4. Create a branch: `PROJ-1234-add-user-authentication`
5. Checkout the new branch

### Work on Your Changes

```bash
# Make your changes
git add .
git commit -m "Add user authentication"
```

### Create a Pull Request

```bash
one pr
```

This will:
1. Push your branch to the remote
2. Extract ticket ID from branch name
3. Generate PR title/body from templates
4. Create the PR via API
5. Open it in your browser (with your work profile!)

### Open a Ticket

```bash
one ticket PROJ-1234
```

Opens the ticket in your configured browser.

## Configuration Examples

### GitHub + Jira

```yaml
version: 1
project:
  name: "Acme Corp"
  paths:
    - "/Users/john/Projects/acme"

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
    
    ## Ticket
    {ticket_url}

branch_patterns:
  ticket_id: "^([A-Z]+-\\d+)"
```

### GitLab + Linear (Minimal)

```yaml
version: 1
project:
  name: "Beta Startup"
  paths:
    - "/Users/jane/beta"

git:
  provider: gitlab
  remote: origin
  base_branch: develop
  gitlab:
    project_id: 42
    token_env: GITLAB_TOKEN

browser:
  type: firefox
  profile: "default"

ticket:
  system: linear
  base_url: https://linear.app/beta
```

## Tips & Tricks

### Multiple Projects

One CLI automatically detects which project you're in:

```bash
# Working on project A
cd ~/Projects/acme
one start ACME-123  # Uses acme config

# Working on project B
cd ~/Projects/beta
one start BET-456   # Uses beta config
```

### Browser Profiles

Keep work and personal separate:

```yaml
browser:
  type: chrome
  profile: "Work Profile"  # Use your work Google account
```

### Custom Templates

Customize per project:

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
    - [ ] Self-reviewed
```

### View All Configs

```bash
one config list     # See all projects
one config show     # Show current project
```

### Get Help

```bash
one help            # Beautiful formatted help
one docs            # Full documentation
one docs --spec     # Technical specification
one docs --examples # Configuration examples
```

## Troubleshooting

### "No project configuration found"

Make sure you're in a directory that matches a configured project path.

```bash
one config list  # See configured projects
one init         # Create new config
```

### "Not authenticated"

Set up authentication:

```bash
# Via keyring (recommended)
one login github

# Or via environment variable
export GITHUB_TOKEN="ghp_xxxxxxxxxxxx"
```

### "Working directory not clean"

Commit or stash your changes:

```bash
git add .
git commit -m "WIP"
# or
git stash
```

## Next Steps

- Read the [full specification](SPECIFICATION.md) for all features
- Check [examples/](examples/) for more configurations
- Star the repo if you find it useful! ⭐

---

**Built with ❤️ using Bubble Tea, Glamour, and Go**
