# One CLI - Features Showcase

## 🎨 Interactive Forms with Huh

One CLI uses [Huh](https://github.com/charmbracelet/huh) for all interactive forms, providing a beautiful, accessible user experience.

### Project Initialization (`one init`)

The init command features a multi-step form wizard:

**Step 1: Basic Information**
- Project name (with validation)
- Project path (defaults to current directory)
- Git provider selection (GitHub/GitLab/Bitbucket)
- Remote name and base branch

**Step 2: Provider Configuration**
- GitHub: Organization/username and repository
- GitLab: Numeric project ID
- Bitbucket: Workspace and repository slug

**Step 3: Browser & Tickets**
- Browser selection (Chrome/Firefox/Safari)
- Optional browser profile
- Optional ticket system integration (Jira/Linear/GitHub Issues)

**Features:**
- ✅ Real-time validation
- ✅ Keyboard navigation
- ✅ Clear descriptions and placeholders
- ✅ Conditional forms (only show relevant fields)
- ✅ Beautiful, consistent UI

### Task Starting (`one start`)

Interactive prompts for edge cases:

**Uncommitted Changes Dialog:**
```
┌─ Working directory has uncommitted changes ─┐
│                                              │
│ Would you like to stash them?               │
│                                              │
│  ● Yes, stash changes                        │
│    No, cancel                                │
└──────────────────────────────────────────────┘
```

**Branch Description Prompt:**
When no ticket system is configured:
```
┌─ Branch Description ─────────────────────────┐
│                                              │
│ Describe what you'll be working on          │
│                                              │
│ ▸ add user authentication                   │
└──────────────────────────────────────────────┘
```

## 🎭 Progress Display with Bubble Tea

### Real-time Feedback

Commands show live progress for long-running operations:

**Starting a Task:**
```
Starting new task...

  Project: Acme Corp

Checking out base branch...
  Branch: main
  ✓ Checked out main

Pulling latest changes...
  ✓ Pulled from origin/main

Fetching ticket info...
  Ticket: Add user authentication

Creating new branch...
  Branch: PROJ-1234-add-user-authentication

✓ Ready to work on PROJ-1234-add-user-authentication!

  When done, run one pr to create a PR
```

**Creating a PR:**
```
Creating pull request...

  Project: Acme Corp
  Branch: PROJ-1234-add-authentication
  Ticket ID: PROJ-1234

  ✓ Pushed to origin
  ✓ PR created: https://github.com/acme/app/pull/123

Opening in browser...

Done! 🚀
```

## 📚 Beautiful Documentation with Glamour

### Rendered Markdown

All documentation commands use [Glamour](https://github.com/charmbracelet/glamour) for beautiful terminal rendering:

**`one help`** - Complete guide with:
- Syntax-highlighted code blocks
- Formatted tables
- Styled headings and lists
- Properly wrapped text

**`one docs --spec`** - Full specification with:
- Table of contents
- Code examples
- Diagrams (ASCII art)
- Formatted tables

**`one docs --examples`** - Configuration examples with:
- YAML syntax highlighting
- Descriptive headings
- Code block formatting

**`one config show`** - Current config with:
- YAML syntax highlighting
- Project name heading
- Formatted structure

**`one config list`** - All projects with:
- Formatted project cards
- Bullet lists
- Inline code formatting

### Auto-styled Output

Glamour automatically adapts to your terminal theme:
- Dark mode → optimized colors for dark backgrounds
- Light mode → optimized colors for light backgrounds
- Word wrapping for readability

## 🔐 Security Features

### Keyring Integration

**OS-Native Storage:**
- macOS: Keychain Services
- Linux: Secret Service (gnome-keyring/kwallet)
- Windows: Credential Manager

**Per-Project Isolation:**
- Each project can have its own tokens
- Format: `{provider}:{project_name}`
- Example: `github:acme-app`, `jira:acme-app`

**Environment Variable Fallback:**
```bash
export GITHUB_TOKEN="ghp_xxxxxxxxxxxx"
export GITLAB_TOKEN="glpat-xxxxxxxxxxxx"
export JIRA_TOKEN="email@example.com:api_token"
```

## 🌐 Multi-Provider Support

### Git Providers

**GitHub:**
- REST API v3
- Pull request creation
- Issue fetching
- Token authentication

**GitLab:**
- REST API v4
- Merge request creation
- Project-based authentication
- Token authentication

**Bitbucket:**
- REST API 2.0
- Pull request creation
- Workspace-based organization
- Token authentication

### Ticket Systems

**Jira:**
- Full API integration
- Fetch ticket titles
- Generate ticket URLs
- Basic auth with email + token

**Linear:**
- URL generation
- Ticket link creation
- Branch pattern matching

**GitHub Issues:**
- URL generation
- Issue link creation
- Integrated with GitHub provider

## 🎯 Template System

### Variable Substitution

Available variables for PR templates:

| Variable | Description | Example |
|----------|-------------|---------|
| `{ticket_id}` | Extracted from branch | `PROJ-1234` |
| `{branch_name}` | Current branch | `proj-1234-add-feature` |
| `{ticket_url}` | Full ticket URL | `https://jira.../PROJ-1234` |
| `{author}` | Git user name | `John Doe` |
| `{email}` | Git user email | `john@example.com` |
| `{date}` | ISO 8601 date | `2025-10-06` |
| `{base_branch}` | Target branch | `main` |

### Example Templates

**Simple:**
```yaml
templates:
  pr_title: "[{ticket_id}] {branch_name}"
  pr_body: "Closes {ticket_url}"
```

**Detailed:**
```yaml
templates:
  pr_title: "[{ticket_id}] {branch_name}"
  pr_body: |
    ## 📝 Summary
    
    
    ## 🔗 Related
    - Ticket: {ticket_url}
    - Branch: {branch_name}
    - Author: {author}
    - Date: {date}
    
    ## ✅ Checklist
    - [ ] Tests added/updated
    - [ ] Documentation updated
    - [ ] Code reviewed
    - [ ] Ready for deployment
```

## 🖥️ Browser Integration

### Multi-Browser Support

**Chrome:**
- Full profile support
- Command-line profile switching
- Works with Chrome/Chromium/Brave

**Firefox:**
- Full profile support
- Profile manager integration
- Works with Firefox/Firefox Developer Edition

**Safari:**
- Basic support (macOS only)
- No profile switching (Safari limitation)

### Profile Configuration

```yaml
browser:
  type: chrome
  profile: "Work Profile"  # Uses your work Google account
```

**Benefits:**
- Keep work and personal separate
- Different cookies/history per profile
- Different extensions per profile
- Multiple logged-in accounts

## 🔄 Git Workflow

### Branch Management

**Automatic Branch Creation:**
1. Checks working directory is clean
2. Checks out base branch
3. Pulls latest changes
4. Creates sanitized branch name
5. Checks out new branch

**Branch Naming:**
- Format: `{ticket-id}-{sanitized-title}`
- Lowercase
- Spaces → hyphens
- Special characters removed
- Max 50 characters

**Example:**
- Input: `PROJ-1234: Add User Authentication & Login`
- Output: `proj-1234-add-user-authentication-login`

### Push & PR Creation

**Automatic Push:**
- Checks working directory is clean
- Pushes current branch
- Sets up tracking

**PR Creation:**
- Extracts ticket ID from branch
- Renders title/body templates
- Creates PR via API
- Opens in configured browser

## 🎨 Styling & Themes

### Consistent Visual Language

**Success:** Green (✓)
**Error:** Red (✗)
**Info:** Blue (ℹ)
**Warning:** Yellow (⚠)

**Progress Indicators:**
- Spinners for long operations
- Checkmarks for completed steps
- Clear status messages

**Typography:**
- Bold for titles
- Italic for descriptions
- Code blocks for commands
- Proper spacing and alignment

## 📊 Configuration Management

### YAML-Based Configuration

**Validation:**
- Required field checking
- Type validation
- Format validation

**Discovery:**
- Automatic path matching
- Symlink resolution
- Prefix matching for subdirectories

**Multiple Projects:**
```
~/.config/one/projects/
├── acme-corp.yml
├── beta-startup.yml
└── personal.yml
```

**Automatic Detection:**
One CLI automatically knows which project you're in based on your current directory.

## 🚀 Performance

**Benchmarks:**
- Startup time: < 100ms
- Config discovery: < 10ms
- Git operations: < 500ms
- API calls: < 2s (network dependent)
- Memory usage: < 50MB

**Optimizations:**
- Single binary (no runtime dependencies)
- Compiled Go code (fast)
- Efficient YAML parsing
- Minimal allocations

## 🌍 Cross-Platform

**Tested On:**
- ✅ macOS (Intel & Apple Silicon)
- ✅ Linux (Ubuntu, Fedora, Arch)
- ✅ Windows 10/11

**Platform-Specific:**
- Browser launching (different commands per OS)
- Keyring integration (OS-native)
- Path handling (case sensitivity)
- File permissions

## 🪝 Hooks System

Run arbitrary commands before/after creating PRs.

### Before PR Hooks

Perfect for validation:
- **Linting**: RuboCop, ESLint, Pylint
- **Tests**: RSpec, Jest, PyTest
- **Type checking**: TypeScript, MyPy
- **Formatting**: Prettier, Black, gofmt
- **Security**: npm audit, bundler-audit

**If a before_pr hook fails with `fail_on_error: true`, the PR won't be created!**

### After PR Hooks

Perfect for automation:
- **Notifications**: Slack, Discord, email
- **Project management**: Update Jira, Linear
- **Deploy**: Preview environments
- **Documentation**: Auto-generate docs
- **Analytics**: Log to tracking systems

### Example

```yaml
hooks:
  before_pr:
    - name: "Lint code"
      command: "bundle exec rubocop"
      fail_on_error: true  # Stop if fails
    
    - name: "Run tests"
      command: "bundle exec rspec"
      fail_on_error: true
  
  after_pr:
    - name: "Notify team"
      command: 'curl -X POST $SLACK_WEBHOOK -d "{\"text\":\"PR created!\"}"'
      fail_on_error: false  # Just warn if fails
```

### Real-time Output

```
⚡ Running before_pr hooks...

  [1/2] Lint code
        $ bundle exec rubocop

  ✓ Success (took 1.2s)

  [2/2] Run tests
        $ bundle exec rspec

  ✓ Success (took 15.3s)

✓ All before_pr hooks completed
```

See [HOOKS.md](HOOKS.md) for complete documentation.

---

## 📦 Installation

**Single Binary:**
```bash
make build
sudo make install
```

**No Dependencies:**
- No Python/Node/Ruby required
- No system packages required
- Works out of the box

**Multi-Platform Builds:**
```bash
make build-all
```

Produces binaries for:
- darwin/amd64, darwin/arm64
- linux/amd64, linux/arm64
- windows/amd64

---

**Built with ❤️ using Huh, Bubble Tea, Glamour, and Go**
