package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
)

var helpCmd = &cobra.Command{
	Use:   "help",
	Short: "Show detailed help with examples",
	Long:  `Display beautifully formatted help documentation with examples.`,
	RunE:  runHelp,
}

func init() {
	rootCmd.AddCommand(helpCmd)
}

func runHelp(cmd *cobra.Command, args []string) error {
	helpContent := getHelpContent()

	// Render with Glamour
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(100),
	)
	if err != nil {
		// Fallback to plain text if rendering fails
		fmt.Println(helpContent)
		return nil
	}

	rendered, err := r.Render(helpContent)
	if err != nil {
		fmt.Println(helpContent)
		return nil
	}

	fmt.Print(rendered)
	return nil
}

func getHelpContent() string {
	return strings.TrimSpace(`
# One CLI - Complete Guide

**Version**: 0.2.0  
A unified command-line tool for freelancers working across multiple projects.

---

## 🚀 Quick Start

### 1. Initialize Your First Project

` + "```bash" + `
cd /path/to/your/project
one init
` + "```" + `

This will guide you through an interactive setup process.

### 2. Start Working on a Task

` + "```bash" + `
one start PROJ-1234
` + "```" + `

This automatically:
- ✓ Checks out your base branch
- ✓ Pulls latest changes
- ✓ Fetches ticket title from Jira/Linear
- ✓ Creates a new branch with sanitized name
- ✓ Checks out the new branch

### 3. Create a Pull Request

` + "```bash" + `
one pr
` + "```" + `

This automatically:
- ✓ Pushes your branch
- ✓ Extracts ticket ID from branch name
- ✓ Creates PR with templated title/description
- ✓ Opens PR in your browser (with profile!)

---

## 📋 Available Commands

### **one init** [-n NAME]
Initialize a new project configuration interactively.

**Example:**
` + "```bash" + `
one init
one init --name "My Project"
` + "```" + `

---

### **one start** <TICKET-ID> [-d DESCRIPTION]
Start working on a new task.

**Examples:**
` + "```bash" + `
one start PROJ-1234
one start PROJ-1234 --description "Add user authentication"
` + "```" + `

---

### **one pr** [-t TITLE] [-d DESCRIPTION] [--no-browser]
Create and open a pull request.

**Examples:**
` + "```bash" + `
one pr
one pr --title "feat: Add OAuth support"
one pr --no-browser
` + "```" + `

---

### **one ticket** <TICKET-ID>
Open a ticket in your browser.

**Example:**
` + "```bash" + `
one ticket PROJ-1234
` + "```" + `

---

### **one config list**
List all configured projects.

### **one config show**
Show current project configuration.

---

## 🎨 Template Variables

Use these in your PR templates:

| Variable | Example |
|----------|---------|
| ` + "`{ticket_id}`" + ` | PROJ-1234 |
| ` + "`{branch_name}`" + ` | proj-1234-add-feature |
| ` + "`{ticket_url}`" + ` | https://jira.../PROJ-1234 |
| ` + "`{author}`" + ` | John Doe |
| ` + "`{email}`" + ` | john@example.com |
| ` + "`{date}`" + ` | 2025-10-06 |
| ` + "`{base_branch}`" + ` | main |

---

## 📁 Configuration Example

` + "```yaml" + `
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

templates:
  pr_title: "[{ticket_id}] {branch_name}"
  pr_body: |
    ## Changes
    - 
    
    ## Ticket
    {ticket_url}

branch_patterns:
  ticket_id: "^([A-Z]+-\\d+)"
` + "```" + `

---

## 🔒 Authentication

Credentials are stored securely in your system keyring:

- **macOS**: Keychain
- **Linux**: Secret Service (gnome-keyring/kwallet)
- **Windows**: Credential Manager

**Environment Variable Fallback:**
` + "```bash" + `
export GITHUB_TOKEN="ghp_xxxxxxxxxxxx"
export GITLAB_TOKEN="glpat-xxxxxxxxxxxx"
export JIRA_TOKEN="email@example.com:api_token"
` + "```" + `

---

## 🌐 Supported Integrations

### Git Providers
- ✓ GitHub
- ✓ GitLab
- ✓ Bitbucket

### Ticket Systems
- ✓ Jira (with API integration)
- ✓ Linear (URL generation)
- ✓ GitHub Issues (URL generation)

### Browsers
- ✓ Chrome (with profile support)
- ✓ Firefox (with profile support)
- ✓ Safari (macOS only)

---

## 💡 Tips & Tricks

### Multiple Projects
You can have different configurations for different projects:

` + "```bash" + `
~/.config/one/projects/
├── acme-corp.yml
├── beta-startup.yml
└── personal.yml
` + "```" + `

One CLI automatically detects which project you're in based on the current directory.

### Browser Profiles
Keep work and personal browsing separate by using different browser profiles:

` + "```yaml" + `
browser:
  type: chrome
  profile: "Work Profile"
` + "```" + `

### Custom Templates
Customize PR templates per project:

` + "```yaml" + `
templates:
  pr_title: "[{ticket_id}] {branch_name}"
  pr_body: |
    ## 📝 Summary
    
    
    ## 🔗 Related
    - Ticket: {ticket_url}
    - Branch: {branch_name}
    
    ## ✅ Checklist
    - [ ] Tests added
    - [ ] Docs updated
` + "```" + `

---

## 📚 More Information

- See ` + "`examples/`" + ` directory for configuration examples
- See ` + "`SPECIFICATION.md`" + ` for complete technical details
- Run ` + "`one [command] --help`" + ` for command-specific help

---

**Built with ❤️ using Bubble Tea & Glow**
`)
}
