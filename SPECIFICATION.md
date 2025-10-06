# One CLI - Complete Technical Specification

**Version**: 0.2.0  
**Last Updated**: 2025-10-02  
**Purpose**: Language-agnostic specification for implementing One CLI in any programming language

---

## Table of Contents

1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Configuration System](#configuration-system)
4. [Commands Specification](#commands-specification)
5. [Authentication System](#authentication-system)
6. [Git Integration](#git-integration)
7. [Browser Integration](#browser-integration)
8. [API Integrations](#api-integrations)
9. [Template System](#template-system)
10. [Error Handling](#error-handling)
11. [Data Structures](#data-structures)
12. [Algorithms](#algorithms)
13. [Security Requirements](#security-requirements)
14. [Platform Requirements](#platform-requirements)
15. [Testing Requirements](#testing-requirements)

---

## 1. Overview

### 1.1 Purpose

One CLI is a unified command-line tool for freelancers working across multiple projects. It provides a single interface for:
- Creating pull requests across different Git providers
- Managing authentication securely
- Opening tickets in different tracking systems
- Working with multiple browser profiles
- Managing per-project configurations

### 1.2 Design Goals

- **Zero manual configuration**: Everything via interactive prompts
- **Secure by default**: No plaintext credentials storage
- **Cross-platform**: Works on macOS, Linux, and Windows
- **Fast**: < 100ms startup time
- **Extensible**: Easy to add new providers/browsers/ticket systems

### 1.3 Core Concepts

- **Project**: A collection of paths with shared configuration
- **Provider**: Git hosting service (GitHub, GitLab, Bitbucket)
- **Ticket System**: Issue tracking (Jira, Linear, GitHub Issues)
- **Browser Profile**: Isolated browser session with separate cookies/cache
- **Keyring**: OS-provided secure credential storage

---

## 2. Architecture

### 2.1 High-Level Architecture

```
┌─────────────────────────────────────────────────────────┐
│                     CLI Interface                        │
│                    (Command Router)                      │
└─────────────────┬───────────────────────────────────────┘
                  │
    ┌─────────────┴─────────────┐
    │                           │
┌───▼────────┐          ┌──────▼──────┐
│  Commands  │          │   Config    │
│  Module    │◄─────────┤   Loader    │
└───┬────────┘          └──────┬──────┘
    │                          │
    │    ┌──────────────┬──────┴────────┬─────────────┐
    │    │              │               │             │
┌───▼────▼───┐  ┌──────▼──────┐  ┌────▼─────┐  ┌───▼─────┐
│    Git     │  │   Browser   │  │   Auth   │  │   API   │
│  Operations│  │  Launcher   │  │ Manager  │  │ Clients │
└────────────┘  └─────────────┘  └──────────┘  └─────────┘
```

### 2.2 Module Responsibilities

#### 2.2.1 CLI Interface
- Parse command-line arguments
- Route to appropriate command handlers
- Display help and version information
- Handle global flags

#### 2.2.2 Commands Module
- Implement each command (start, pr, init, etc.)
- Coordinate between other modules
- Handle user interaction and prompts
- Format output for display

#### 2.2.3 Config Loader
- Discover configuration files
- Parse YAML configurations
- Match current directory to project
- Merge global and project configs

#### 2.2.4 Git Operations
- Repository discovery
- Branch management
- Push/pull operations
- Stash management
- Branch name parsing

#### 2.2.5 Browser Launcher
- Detect installed browsers
- Find browser profiles
- Launch URLs with specific profiles

#### 2.2.6 Auth Manager
- OAuth flows (device flow for GitHub)
- Keyring storage operations
- Token retrieval and validation
- Per-project credential isolation

#### 2.2.7 API Clients
- GitHub API (REST v3)
- GitLab API (REST v4)
- Bitbucket API (REST 2.0)
- Jira API (REST v2)
- HTTP request handling with retries

### 2.3 Data Flow

#### 2.3.1 Command Execution Flow

```
User Input → CLI Parser → Command Handler → Business Logic
     ↓
Config Discovery → Load Project Config → Load Credentials
     ↓
Execute Operation → API Calls / Git Ops → Update State
     ↓
Format Output → Display to User
```

#### 2.3.2 Configuration Discovery Flow

```
Current Directory → Check ~/.config/one/projects/*.yml
     ↓
For Each Project Config:
    Check if current_dir matches any path in config.project.paths
     ↓
If Match Found → Load Config
If No Match → Return Error
```

---

## 3. Configuration System

### 3.1 File Locations

#### 3.1.1 Global Configuration
- **Path**: `~/.config/one/config.yml`
- **Purpose**: Default settings for all projects
- **Optional**: Yes

#### 3.1.2 Project Configurations
- **Path**: `~/.config/one/projects/*.yml`
- **Purpose**: Per-project settings
- **Required**: Yes (at least one)

### 3.2 Global Configuration Schema

```yaml
version: integer (required, currently 1)

defaults: object (optional)
  browser: object (optional)
    type: string (chrome|firefox|safari)
  
  git: object (optional)
    remote: string (default: "origin")
    base_branch: string (default: "main")
```

**Example:**
```yaml
version: 1
defaults:
  browser:
    type: chrome
  git:
    remote: origin
    base_branch: main
```

### 3.3 Project Configuration Schema

```yaml
version: integer (required, currently 1)

project: object (required)
  name: string (required)
  paths: array<string> (required, at least one)

git: object (required)
  provider: string (required, github|gitlab|bitbucket)
  remote: string (required)
  base_branch: string (required)
  
  github: object (optional, required if provider=github)
    owner: string (required)
    repo: string (required)
    token_env: string (required, environment variable name)
  
  gitlab: object (optional, required if provider=gitlab)
    project_id: integer (required)
    token_env: string (required)
  
  bitbucket: object (optional, required if provider=bitbucket)
    workspace: string (required)
    repo_slug: string (required)
    token_env: string (required)

browser: object (required)
  type: string (required, chrome|firefox|safari)
  profile: string (optional)

ticket: object (optional)
  system: string (required, jira|linear|github)
  base_url: string (required, URL without trailing slash)
  
  jira: object (optional, required if system=jira)
    board_id: string (required)
    token_env: string (optional)

templates: object (optional)
  pr_title: string (required)
  pr_body: string (required, multiline)

branch_patterns: object (optional)
  ticket_id: string (required, regex with one capture group)
```

**Example:**
```yaml
version: 1
project:
  name: "Acme Corp Project"
  paths:
    - "/Users/john/Projects/acme-app"
    - "/Users/john/Projects/acme-frontend"

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

### 3.4 Configuration Loading Algorithm

```
function load_config():
    current_dir = get_current_directory()
    projects_dir = get_projects_directory()
    
    if not projects_dir.exists():
        return error("No projects configured")
    
    for config_file in projects_dir.list_files("*.yml", "*.yaml"):
        project_config = parse_yaml(config_file)
        
        for path in project_config.project.paths:
            if current_dir.starts_with(path):
                return project_config
    
    return error("No project found for current directory")
```

### 3.5 Path Matching Rules

- **Exact match**: `/Users/john/project` matches `/Users/john/project`
- **Prefix match**: `/Users/john/project` matches `/Users/john/project/src/main`
- **Case-sensitive**: Paths are case-sensitive on Unix, case-insensitive on Windows
- **Normalized**: Paths should be normalized (no `..` or `.`)

---

## 4. Commands Specification

### 4.1 Command: `one start <TICKET-ID>`

**Purpose**: Start working on a new task

**Arguments**:
- `ticket_id` (required): Ticket identifier (e.g., PROJ-1234)

**Options**:
- `-d, --description <TEXT>`: Custom branch description (overrides ticket title)

**Algorithm**:
```
1. Load project configuration
2. Open git repository
3. Check if working directory is clean:
   - If not clean:
     - Prompt user: "Stash changes? [y/N]"
     - If yes: stash with message "one start: auto-stash"
     - If no: exit with error
4. Checkout base branch (from config.git.base_branch)
5. Pull latest changes:
   - Fetch from remote
   - Fast-forward merge
   - If merge conflicts: exit with error
6. Fetch ticket information:
   - If config.ticket exists:
     - Make API call to ticket system
     - Extract ticket title
   - If custom description provided: use that
   - Otherwise: prompt user for description
7. Generate branch name:
   - Format: "{ticket_id}-{sanitized_title}"
   - Sanitize: lowercase, replace spaces with hyphens, remove special chars
   - Truncate if > 50 chars
8. Create new branch
9. Checkout new branch
10. Display success message
```

**Success Output**:
```
Starting new task...

  Project: Acme Project

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

**Error Cases**:
- No git repository found
- No project configuration found
- Working directory not clean and user declines stash
- Remote branch doesn't exist
- Ticket system API error (warn but continue)
- Branch already exists

### 4.2 Command: `one pr`

**Purpose**: Create and open a pull request

**Options**:
- `-t, --title <TEXT>`: Custom PR title (overrides template)
- `-d, --description <TEXT>`: Custom PR description (overrides template)
- `--no-browser`: Skip opening browser

**Algorithm**:
```
1. Load project configuration
2. Open git repository
3. Get current branch name
4. Extract ticket ID from branch name:
   - Use config.branch_patterns.ticket_id regex
   - Extract first capture group
5. Check if working directory is clean:
   - If not: exit with error
6. Push current branch to remote:
   - Push with tracking (-u flag)
7. Generate PR title:
   - If custom title provided: use that
   - Otherwise: render template with variables
8. Generate PR body:
   - If custom description provided: use that
   - Otherwise: render template with variables
9. Create PR via API:
   - Call appropriate provider API
   - Handle API errors
10. Get PR URL from response
11. If not --no-browser:
    - Open URL in configured browser with profile
12. Display success message with PR URL
```

**Template Variables**:
- `{ticket_id}`: Extracted ticket ID (e.g., PROJ-1234)
- `{branch_name}`: Current branch name
- `{ticket_url}`: Full URL to ticket
- `{author}`: Git user name
- `{date}`: Current date (ISO 8601)
- `{base_branch}`: Target branch

**Success Output**:
```
Creating pull request...

  Project: Acme Project
  Branch: PROJ-1234-add-authentication
  Ticket ID: PROJ-1234

Pushing to remote...
  ✓ Pushed to origin

PR Details:
  Title: [PROJ-1234] add-authentication
  Body:
    ## Changes
    - Added OAuth flow
    
    ## Related Ticket
    https://acme.atlassian.net/browse/PROJ-1234

Creating PR...
  ✓ PR created: https://github.com/acme/app/pull/123

Opening in browser...

Done! 🚀
```

**Error Cases**:
- Not in a git repository
- Working directory not clean
- No remote configured
- Branch not tracking remote
- API authentication failed
- API rate limit exceeded
- No project configuration found

### 4.3 Command: `one init`

**Purpose**: Initialize a new project configuration

**Options**:
- `-n, --name <NAME>`: Project name (skips prompt)

**Algorithm**:
```
1. Prompt for project name (if not provided)
2. Prompt for project path (default: current directory)
3. Prompt for git provider (github/gitlab/bitbucket)
4. Prompt for git remote name (default: origin)
5. Prompt for base branch (default: main)
6. Provider-specific prompts:
   - GitHub: owner, repo, token_env
   - GitLab: project_id, token_env
   - Bitbucket: workspace, repo_slug, token_env
7. Prompt for browser type (default: chrome)
8. If chrome/firefox:
   - Detect available profiles
   - Display list with numbers
   - Prompt for selection
9. Prompt for ticket ID regex pattern (default: ^([A-Z]+-\d+))
10. Generate configuration object
11. Serialize to YAML
12. Write to ~/.config/one/projects/{sanitized-name}.yml
13. Prompt for OAuth login:
    - Check if already logged in
    - If not: "Login to {provider} now? [Y/n]"
    - If yes: execute login flow
14. Display success message with next steps
```

**Interactive Example**:
```
Initialize New Project Configuration

Project name: Acme App
Project path [/Users/john/Projects/acme]: ⏎
Git provider (github/gitlab/bitbucket) [github]: ⏎
Git remote name [origin]: ⏎
Base branch [main]: ⏎
GitHub owner: acme-corp
GitHub repo: main-app
Token environment variable [GITHUB_TOKEN]: ⏎

Browser Configuration
Browser (chrome/firefox/safari) [chrome]: ⏎

Available profiles:
  1. Default
  2. Work Profile
  3. Personal

Select profile number (or press Enter to skip): 2

Ticket ID regex pattern [^([A-Z]+-\d+)]: ⏎

✓ Configuration saved to: ~/.config/one/projects/acme-app.yml

Authentication Setup

Would you like to set up authentication now?
This will securely store your credentials in the system keyring.

Login to github now? [Y/n]: y

GitHub Authentication

Please visit: https://github.com/login/device
Enter code: ABCD-1234

Waiting for authorization...

✓ Successfully authenticated!

Setup Complete! 🎉

Next steps:
  1. one start TICKET-123 to start a new task
  2. one pr to create a pull request
  3. one ticket TICKET-123 to open a ticket
```

### 4.4 Command: `one ticket <TICKET-ID>`

**Purpose**: Open a ticket in the browser

**Arguments**:
- `ticket_id` (required): Ticket identifier

**Algorithm**:
```
1. Load project configuration
2. Check if config.ticket exists:
   - If not: exit with error
3. Generate ticket URL based on system:
   - jira: {base_url}/browse/{ticket_id}
   - linear: {base_url}/issue/{ticket_id}
   - github: {base_url}/issues/{ticket_id}
4. Open URL in configured browser with profile
5. Display success message
```

### 4.5 Command: `one config list`

**Purpose**: List all configured projects

**Algorithm**:
```
1. Get projects directory (~/.config/one/projects/)
2. List all .yml and .yaml files
3. For each file:
   - Parse YAML
   - Extract project.name and git.provider
   - Display with bullet point
4. If no projects found: display help message
```

**Output**:
```
Configured Projects:

  ● Acme App
    Provider: github
    Paths:
      - /Users/john/Projects/acme

  ● Beta Inc
    Provider: gitlab
    Paths:
      - /Users/john/Projects/beta

  ● Personal Project
    Provider: github
    Paths:
      - /Users/john/Projects/personal
```

### 4.6 Command: `one config show`

**Purpose**: Show current project configuration

**Algorithm**:
```
1. Discover project config for current directory
2. If found:
   - Display formatted configuration
3. If not found:
   - Display error message
   - Suggest running "one init"
```

### 4.7 Command: `one config validate`

**Purpose**: Validate all configuration files

**Algorithm**:
```
1. Get projects directory
2. For each config file:
   - Try to parse YAML
   - Validate against schema
   - Check required fields
   - Display ✓ or ✗ with filename
3. Display summary: "Valid: X, Invalid: Y"
```

### 4.8 Command: `one login [SERVICE]`

**Purpose**: Authenticate with services

**Arguments**:
- `service` (optional): Service name (github/gitlab/jira/bitbucket)

**Algorithm (Interactive)**:
```
1. If service not specified:
   - Display menu:
     1. GitHub
     2. GitLab
     3. Jira
     4. Bitbucket
   - Prompt for selection
2. Get current project name (if in project directory)
3. Execute service-specific login flow
4. Store credentials in keyring with key: "{service}:{project_name}"
5. Display success message
```

**GitHub OAuth Flow** (Device Flow):
```
1. POST to https://github.com/login/device/code
   - Body: client_id={CLIENT_ID}
   - Response: device_code, user_code, verification_uri, interval
2. Display: "Visit: {verification_uri}, Enter: {user_code}"
3. Open browser to verification_uri automatically
4. Poll https://github.com/login/oauth/access_token every {interval} seconds:
   - Body: client_id, device_code, grant_type
   - Until: access_token received or error
5. Store token in keyring
```

**GitLab/Jira/Bitbucket Flow** (Manual Token):
```
1. Display instructions for creating token
2. Prompt for token/credentials
3. Validate by making test API call
4. Store in keyring
```

### 4.9 Command: `one logout [SERVICE]`

**Purpose**: Remove stored credentials

**Arguments**:
- `service` (optional): Service to logout from

**Algorithm**:
```
1. Get current project name
2. If service specified:
   - Delete token for "{service}:{project_name}"
3. If service not specified:
   - Delete all tokens for current project
   - Delete all default tokens
4. Display success message
```

### 4.10 Command: `one status`

**Purpose**: Show authentication status

**Algorithm**:
```
1. Get current project name (if in project)
2. For each service (github, gitlab, jira, bitbucket):
   - Check if token exists in keyring
   - Display "✓ Logged in" or "○ Not logged in"
3. Display help message if not logged in
```

### 4.11 Command: `one profiles`

**Purpose**: List available browser profiles

**Algorithm**:
```
1. Detect Chrome profiles:
   - Find Chrome user data directory
   - List directories: Default, Profile 1, Profile 2, etc.
   - Read profile names from Preferences JSON
2. Detect Firefox profiles:
   - Find Firefox directory
   - Parse profiles.ini
   - Extract profile names
3. Display grouped by browser
4. Display usage example
```

**Output**:
```
Available Browser Profiles

  ● Chrome
    - Default
    - Work Profile
    - Personal

  ● Firefox
    - default-release
    - dev-edition

Use profile names in your project configuration:
  browser:
    type: chrome
    profile: "Work Profile"
```

---

## 5. Authentication System

### 5.1 Credential Storage

#### 5.1.1 Keyring Integration

**Platforms**:
- **macOS**: Keychain (`security` command or Keychain API)
- **Linux**: Secret Service (D-Bus, libsecret)
- **Windows**: Credential Manager (CredWrite/CredRead API)

**Storage Format**:
```
Service: "one-cli"
Account: "{provider}:{project_name}"
Password: JSON string with token data
```

**Example**:
```
Service: one-cli
Account: github:acme-app
Password: {
  "access_token": "ghp_xxxxxxxxxxxx",
  "token_type": "bearer",
  "expires_at": null,
  "refresh_token": null
}
```

#### 5.1.2 Token Structure

```json
{
  "access_token": "string (required)",
  "refresh_token": "string (optional)",
  "token_type": "string (required, usually 'bearer')",
  "expires_at": "integer (optional, Unix timestamp)"
}
```

#### 5.1.3 Keyring Operations

**Store Token**:
```
function store_token(service, account, token_json):
    entry = keyring.create_entry(service, account)
    entry.set_password(token_json)
```

**Retrieve Token**:
```
function get_token(service, account):
    entry = keyring.get_entry(service, account)
    return entry.get_password()
```

**Delete Token**:
```
function delete_token(service, account):
    entry = keyring.get_entry(service, account)
    entry.delete_password()
```

**Check Token Exists**:
```
function has_token(service, account):
    try:
        get_token(service, account)
        return true
    catch error:
        return false
```

### 5.2 OAuth Flows

#### 5.2.1 GitHub Device Flow

**Specifications**:
- OAuth 2.0 Device Authorization Grant (RFC 8628)
- Client ID: Public (no client secret required)
- Scopes: `repo`, `workflow`

**Algorithm**:
```
1. Request Device Code:
   POST https://github.com/login/device/code
   Headers:
     Accept: application/json
   Body:
     client_id: {CLIENT_ID}
   
   Response:
     device_code: string
     user_code: string (e.g., "ABCD-1234")
     verification_uri: string (e.g., "https://github.com/login/device")
     expires_in: integer (seconds)
     interval: integer (polling interval in seconds)

2. Display to User:
   - Show verification_uri
   - Show user_code
   - Open verification_uri in browser automatically

3. Poll for Token:
   Loop every {interval} seconds:
     POST https://github.com/login/oauth/access_token
     Headers:
       Accept: application/json
     Body:
       client_id: {CLIENT_ID}
       device_code: {device_code}
       grant_type: urn:ietf:params:oauth:grant-type:device_code
     
     Response (pending):
       error: "authorization_pending"
     
     Response (slow down):
       error: "slow_down"
       → Increase interval by 5 seconds
     
     Response (success):
       access_token: string
       token_type: string
       scope: string
     
     Response (error):
       error: "expired_token" | "access_denied" | "unsupported_grant_type"
       → Exit with error

4. Store Token:
   - Create token object
   - Store in keyring
```

**Timeout**: Exit if no response after `expires_in` seconds

**Error Handling**:
- `authorization_pending`: Continue polling
- `slow_down`: Increase interval
- `expired_token`: Ask user to try again
- `access_denied`: User rejected authorization
- Network errors: Retry up to 3 times

#### 5.2.2 Token Usage Priority

When making API calls, check for tokens in this order:

```
1. Keyring token for current project: "{provider}:{project_name}"
2. Keyring token for default: "{provider}:default"
3. Environment variable: {config.git.{provider}.token_env}
4. Exit with error: "Not authenticated"
```

### 5.3 Token Validation

Before using a token, validate it:

```
function validate_token(provider, token):
    if token.expires_at exists:
        if current_time >= token.expires_at:
            return false
    
    # Test API call
    response = api_call(provider, "/user", token)
    return response.status == 200
```

---

## 6. Git Integration

### 6.1 Git Operations

#### 6.1.1 Repository Discovery

```
function discover_repository():
    current_dir = get_current_directory()
    
    while current_dir != root:
        git_dir = current_dir / ".git"
        if git_dir.exists():
            return open_repository(git_dir)
        current_dir = current_dir.parent()
    
    return error("Not a git repository")
```

#### 6.1.2 Get Current Branch

```
function current_branch(repo):
    head_ref = repo.head()
    
    if head_ref.is_branch():
        return head_ref.shorthand()
    else:
        return error("HEAD is detached")
```

#### 6.1.3 Check if Clean

```
function is_clean(repo):
    statuses = repo.statuses()
    return statuses.is_empty()
```

#### 6.1.4 Stash Changes

```
function stash_save(repo, message):
    signature = repo.signature()  # From git config user.name/email
    repo.stash_save(signature, message, STASH_DEFAULT)
```

#### 6.1.5 Checkout Branch

```
function checkout_branch(repo, branch_name):
    # Get branch reference
    ref = repo.revparse("refs/heads/" + branch_name)
    
    # Checkout
    repo.checkout_tree(ref)
    repo.set_head("refs/heads/" + branch_name)
```

#### 6.1.6 Create Branch

```
function create_branch(repo, branch_name):
    # Get current commit
    head = repo.head()
    commit = head.peel_to_commit()
    
    # Create branch
    repo.branch(branch_name, commit, force=false)
```

#### 6.1.7 Pull (Fast-Forward)

```
function pull(repo, remote_name, branch_name):
    # Fetch
    remote = repo.find_remote(remote_name)
    remote.fetch([branch_name])
    
    # Get fetch head
    fetch_head = repo.find_reference("FETCH_HEAD")
    fetch_commit = repo.reference_to_annotated_commit(fetch_head)
    
    # Analyze
    analysis = repo.merge_analysis([fetch_commit])
    
    if analysis.is_up_to_date():
        return  # Already up to date
    
    if analysis.is_fast_forward():
        # Fast-forward merge
        ref = repo.find_reference("refs/heads/" + branch_name)
        ref.set_target(fetch_commit.id())
        repo.set_head("refs/heads/" + branch_name)
        repo.checkout_head(force=true)
    else:
        return error("Cannot fast-forward, manual merge required")
```

#### 6.1.8 Push

```
function push(repo, remote_name, branch_name):
    remote = repo.find_remote(remote_name)
    
    refspec = "refs/heads/" + branch_name + ":refs/heads/" + branch_name
    
    remote.push([refspec], options={
        callbacks: {
            credentials: get_credentials_callback()
        }
    })
```

### 6.2 Branch Name Parsing

```
function parse_ticket_id(branch_name, pattern):
    regex = compile(pattern)
    match = regex.match(branch_name)
    
    if match and match.groups.length > 0:
        return match.groups[0]
    
    return null
```

**Example**:
- Branch: `PROJ-1234-add-feature`
- Pattern: `^([A-Z]+-\d+)`
- Result: `PROJ-1234`

### 6.3 Branch Name Sanitization

```
function sanitize_for_branch(text):
    # Convert to lowercase
    text = text.to_lowercase()
    
    # Replace spaces with hyphens
    text = text.replace(" ", "-")
    
    # Remove special characters (keep alphanumeric and hyphens)
    text = regex_replace(text, "[^a-z0-9-]", "")
    
    # Remove consecutive hyphens
    text = regex_replace(text, "-+", "-")
    
    # Remove leading/trailing hyphens
    text = text.trim("-")
    
    # Truncate if too long
    if text.length > 50:
        text = text[0:50]
    
    return text
```

---

## 7. Browser Integration

### 7.1 Browser Profile Detection

#### 7.1.1 Chrome Profile Discovery

**Locations**:
- **macOS**: `~/Library/Application Support/Google/Chrome/`
- **Linux**: `~/.config/google-chrome/`
- **Windows**: `%LOCALAPPDATA%\Google\Chrome\User Data\`

**Algorithm**:
```
function list_chrome_profiles():
    chrome_dir = get_chrome_user_data_dir()
    profiles = []
    
    # Default profile
    default_path = chrome_dir / "Default"
    if default_path.exists():
        name = read_profile_name(default_path / "Preferences")
        profiles.append({
            name: name or "Default",
            path: default_path
        })
    
    # Numbered profiles
    for dir in chrome_dir.list_directories():
        if dir.name.starts_with("Profile "):
            prefs_path = dir / "Preferences"
            name = read_profile_name(prefs_path)
            profiles.append({
                name: name or dir.name,
                path: dir
            })
    
    return profiles
```

**Reading Profile Name**:
```
function read_profile_name(preferences_file):
    if not preferences_file.exists():
        return null
    
    json = parse_json(read_file(preferences_file))
    return json.profile.name or null
```

#### 7.1.2 Firefox Profile Discovery

**Locations**:
- **macOS**: `~/Library/Application Support/Firefox/`
- **Linux**: `~/.mozilla/firefox/`
- **Windows**: `%APPDATA%\Mozilla\Firefox\`

**Algorithm**:
```
function list_firefox_profiles():
    firefox_dir = get_firefox_profiles_dir()
    profiles_ini = firefox_dir / "profiles.ini"
    
    if not profiles_ini.exists():
        return []
    
    profiles = []
    current_name = null
    
    for line in read_lines(profiles_ini):
        line = line.trim()
        
        if line.starts_with("Name="):
            current_name = line[5:]
        
        if line.starts_with("Path="):
            path = line[5:]
            if current_name:
                profiles.append({
                    name: current_name,
                    path: firefox_dir / path
                })
                current_name = null
    
    return profiles
```

### 7.2 Browser Launching

#### 7.2.1 Chrome

**Command Format**:
- **macOS**: `open -a "Google Chrome" --args --profile-directory="{profile}" "{url}"`
- **Linux**: `google-chrome --profile-directory="{profile}" "{url}"`
- **Windows**: `chrome.exe --profile-directory="{profile}" "{url}"`

**Profile Directory Names**:
- Default profile: `Default`
- Other profiles: `Profile 1`, `Profile 2`, etc. (internal directory name, not display name)

**Note**: Must use internal directory name from `list_chrome_profiles()`, not display name!

#### 7.2.2 Firefox

**Command Format**:
- **macOS**: `open -a "Firefox" --args -P "{profile}" "{url}"`
- **Linux**: `firefox -P "{profile}" "{url}"`
- **Windows**: `firefox.exe -P "{profile}" "{url}"`

**Profile Name**: Use display name from `profiles.ini`

#### 7.2.3 Safari

**Command Format**:
- **macOS**: `open -a "Safari" "{url}"`

**Note**: Safari doesn't support profile switching via command line

### 7.3 URL Opening Algorithm

```
function open_url(browser_config, url):
    browser_type = browser_config.type
    profile = browser_config.profile
    
    if browser_type == "chrome":
        command = build_chrome_command(url, profile)
    elif browser_type == "firefox":
        command = build_firefox_command(url, profile)
    elif browser_type == "safari":
        command = build_safari_command(url)
    else:
        return error("Unsupported browser: " + browser_type)
    
    execute_command(command, background=true)
```

---

## 8. API Integrations

### 8.1 GitHub API

#### 8.1.1 Base URL
- `https://api.github.com`

#### 8.1.2 Authentication
```
Headers:
  Authorization: Bearer {access_token}
  Accept: application/vnd.github+json
  User-Agent: one-cli
```

#### 8.1.3 Create Pull Request

**Endpoint**: `POST /repos/{owner}/{repo}/pulls`

**Request Body**:
```json
{
  "title": "string (required)",
  "body": "string (optional)",
  "head": "string (required, branch name)",
  "base": "string (required, target branch)",
  "draft": false
}
```

**Response** (201 Created):
```json
{
  "html_url": "https://github.com/owner/repo/pull/123",
  "number": 123,
  "state": "open",
  ...
}
```

**Error Responses**:
- 401: Authentication failed
- 403: Rate limit exceeded
- 422: Validation failed (e.g., PR already exists)

#### 8.1.4 Get Issue

**Endpoint**: `GET /repos/{owner}/{repo}/issues/{issue_number}`

**Response** (200 OK):
```json
{
  "title": "string",
  "body": "string",
  "state": "open|closed",
  ...
}
```

### 8.2 GitLab API

#### 8.2.1 Base URL
- **GitLab.com**: `https://gitlab.com/api/v4`
- **Self-hosted**: `https://{instance}/api/v4`

#### 8.2.2 Authentication
```
Headers:
  PRIVATE-TOKEN: {access_token}
```

#### 8.2.3 Create Merge Request

**Endpoint**: `POST /projects/{project_id}/merge_requests`

**Request Body**:
```json
{
  "source_branch": "string (required)",
  "target_branch": "string (required)",
  "title": "string (required)",
  "description": "string (optional)"
}
```

**Response** (201 Created):
```json
{
  "web_url": "https://gitlab.com/owner/repo/-/merge_requests/123",
  "iid": 123,
  "state": "opened",
  ...
}
```

### 8.3 Bitbucket API

#### 8.3.1 Base URL
- `https://api.bitbucket.org/2.0`

#### 8.3.2 Authentication
```
Headers:
  Authorization: Bearer {access_token}
```

#### 8.3.3 Create Pull Request

**Endpoint**: `POST /repositories/{workspace}/{repo_slug}/pullrequests`

**Request Body**:
```json
{
  "title": "string (required)",
  "description": "string (optional)",
  "source": {
    "branch": {
      "name": "string (required)"
    }
  },
  "destination": {
    "branch": {
      "name": "string (required)"
    }
  }
}
```

**Response** (201 Created):
```json
{
  "links": {
    "html": {
      "href": "https://bitbucket.org/workspace/repo/pull-requests/123"
    }
  },
  "id": 123,
  "state": "OPEN",
  ...
}
```

### 8.4 Jira API

#### 8.4.1 Base URL
- `https://{instance}.atlassian.net/rest/api/2`

#### 8.4.2 Authentication
```
Headers:
  Authorization: Basic {base64(email:token)}
```

#### 8.4.3 Get Issue

**Endpoint**: `GET /issue/{issue_key}`

**Response** (200 OK):
```json
{
  "key": "PROJ-1234",
  "fields": {
    "summary": "Issue title",
    "description": "Issue description",
    "status": {
      "name": "In Progress"
    },
    ...
  }
}
```

### 8.5 HTTP Client Requirements

#### 8.5.1 Timeouts
- Connection timeout: 10 seconds
- Read timeout: 30 seconds

#### 8.5.2 Retries
- Retry on network errors: 3 attempts
- Exponential backoff: 1s, 2s, 4s
- Don't retry on 4xx errors (except 429)

#### 8.5.3 Rate Limiting
- Handle 429 responses
- Respect `Retry-After` header
- Implement exponential backoff

#### 8.5.4 User Agent
- Format: `one-cli/{version}`
- Example: `one-cli/0.2.0`

---

## 9. Template System

### 9.1 Template Variables

#### 9.1.1 Available Variables

| Variable | Type | Description | Example |
|----------|------|-------------|---------|
| `{ticket_id}` | string | Extracted ticket ID | `PROJ-1234` |
| `{branch_name}` | string | Current branch name | `proj-1234-add-feature` |
| `{ticket_url}` | string | Full URL to ticket | `https://jira.../PROJ-1234` |
| `{author}` | string | Git user name | `John Doe` |
| `{email}` | string | Git user email | `john@example.com` |
| `{date}` | string | Current date (ISO 8601) | `2025-10-02` |
| `{base_branch}` | string | Target branch | `main` |

#### 9.1.2 Variable Resolution

```
function resolve_variables(template, context):
    result = template
    
    for key, value in context:
        placeholder = "{" + key + "}"
        result = result.replace(placeholder, value)
    
    return result
```

#### 9.1.3 Context Building

```
function build_template_context(config, repo, ticket_id):
    context = {}
    
    # Basic info
    context["ticket_id"] = ticket_id or ""
    context["branch_name"] = repo.current_branch()
    context["base_branch"] = config.git.base_branch
    
    # Git info
    signature = repo.signature()
    context["author"] = signature.name
    context["email"] = signature.email
    
    # Date
    context["date"] = current_date_iso8601()
    
    # Ticket URL
    if config.ticket and ticket_id:
        context["ticket_url"] = build_ticket_url(config.ticket, ticket_id)
    else:
        context["ticket_url"] = ""
    
    return context
```

### 9.2 Ticket URL Generation

```
function build_ticket_url(ticket_config, ticket_id):
    base_url = ticket_config.base_url
    system = ticket_config.system
    
    if system == "jira":
        return base_url + "/browse/" + ticket_id
    elif system == "linear":
        return base_url + "/issue/" + ticket_id
    elif system == "github":
        return base_url + "/issues/" + ticket_id
    else:
        return ""
```

### 9.3 Default Templates

If `config.templates` not specified:

**PR Title**:
```
{branch_name}
```

**PR Body**:
```
(empty)
```

---

## 10. Error Handling

### 10.1 Error Types

#### 10.1.1 Configuration Errors
- `ConfigNotFound`: No project config for current directory
- `ConfigInvalid`: YAML parse error or invalid schema
- `ConfigMissing`: Required field missing

#### 10.1.2 Git Errors
- `NotGitRepo`: Not in a git repository
- `WorkingDirectoryNotClean`: Uncommitted changes
- `BranchNotFound`: Branch doesn't exist
- `PushFailed`: Failed to push to remote
- `MergeConflict`: Cannot fast-forward merge

#### 10.1.3 Authentication Errors
- `NotAuthenticated`: No credentials found
- `AuthenticationFailed`: Invalid credentials
- `TokenExpired`: Token no longer valid

#### 10.1.4 API Errors
- `APIError`: Generic API error
- `RateLimitExceeded`: Too many requests
- `ResourceNotFound`: PR/issue not found
- `ValidationError`: Invalid request data

#### 10.1.5 System Errors
- `KeyringError`: Cannot access system keyring
- `BrowserNotFound`: Configured browser not installed
- `NetworkError`: Connection failed

### 10.2 Error Messages

#### 10.2.1 Message Format

```
Error: {brief_description}

{detailed_explanation}

Suggestion: {how_to_fix}
```

**Example**:
```
Error: No project configuration found

Could not find a project configuration matching the current directory:
  /Users/john/Projects/unknown

Suggestion: Run 'one init' to create a project configuration.
```

#### 10.2.2 Error Message Guidelines

- **Be specific**: Include actual values (paths, branch names, etc.)
- **Be actionable**: Suggest how to fix the problem
- **Be friendly**: Use conversational tone
- **Include context**: Show what was being attempted

### 10.3 Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Configuration error |
| 3 | Git error |
| 4 | Authentication error |
| 5 | API error |
| 6 | User cancelled |

---

## 11. Data Structures

### 11.1 Core Types

#### 11.1.1 ProjectConfig
```
struct ProjectConfig:
    version: integer
    project: ProjectInfo
    git: GitConfig
    browser: BrowserConfig
    ticket: TicketConfig? (optional)
    templates: Templates? (optional)
    branch_patterns: BranchPatterns? (optional)
```

#### 11.1.2 ProjectInfo
```
struct ProjectInfo:
    name: string
    paths: array<string>
```

#### 11.1.3 GitConfig
```
struct GitConfig:
    provider: string (github|gitlab|bitbucket)
    remote: string
    base_branch: string
    github: GitHubConfig? (optional)
    gitlab: GitLabConfig? (optional)
    bitbucket: BitbucketConfig? (optional)
```

#### 11.1.4 BrowserConfig
```
struct BrowserConfig:
    type: string (chrome|firefox|safari)
    profile: string? (optional)
```

#### 11.1.5 TicketConfig
```
struct TicketConfig:
    system: string (jira|linear|github)
    base_url: string
    jira: JiraConfig? (optional)
```

#### 11.1.6 Templates
```
struct Templates:
    pr_title: string
    pr_body: string
```

#### 11.1.7 BranchPatterns
```
struct BranchPatterns:
    ticket_id: string (regex)
```

#### 11.1.8 StoredToken
```
struct StoredToken:
    access_token: string
    refresh_token: string? (optional)
    token_type: string
    expires_at: integer? (optional, Unix timestamp)
```

#### 11.1.9 BrowserProfile
```
struct BrowserProfile:
    name: string
    path: string (file system path)
```

---

## 12. Algorithms

### 12.1 Path Matching

**Input**: 
- `current_dir`: Current working directory
- `project_paths`: Array of paths from config

**Output**: Boolean (match found)

**Algorithm**:
```
function matches_path(current_dir, project_paths):
    # Normalize current directory
    current_normalized = normalize_path(current_dir)
    
    for path in project_paths:
        # Normalize config path
        path_normalized = normalize_path(path)
        
        # Check if current_dir starts with path
        if current_normalized.starts_with(path_normalized):
            return true
    
    return false

function normalize_path(path):
    # Resolve symlinks
    path = resolve_symlinks(path)
    
    # Remove trailing slashes
    path = path.trim_end("/")
    
    # On Windows: normalize separators
    if windows:
        path = path.replace("\\", "/")
    
    return path
```

### 12.2 Config Discovery

**Input**: Current working directory

**Output**: ProjectConfig or error

**Algorithm**:
```
function discover_config(current_dir):
    projects_dir = get_projects_directory()
    
    # List all config files
    config_files = []
    for file in projects_dir.list_files():
        if file.extension in ["yml", "yaml"]:
            config_files.append(file)
    
    # Try each config
    for config_file in config_files:
        try:
            config = parse_yaml(config_file)
            
            if matches_path(current_dir, config.project.paths):
                return config
        except ParseError:
            # Skip invalid configs
            continue
    
    # No match found
    return error("No project configuration found")
```

### 12.3 Ticket URL Generation

**Input**:
- `ticket_config`: TicketConfig
- `ticket_id`: string

**Output**: URL string

**Algorithm**:
```
function generate_ticket_url(ticket_config, ticket_id):
    base = ticket_config.base_url.trim_end("/")
    
    if ticket_config.system == "jira":
        return base + "/browse/" + ticket_id
    elif ticket_config.system == "linear":
        return base + "/issue/" + ticket_id
    elif ticket_config.system == "github":
        return base + "/issues/" + ticket_id
    else:
        return ""
```

### 12.4 Branch Name Sanitization

**Input**: Raw text (e.g., "Add User Authentication & Login")

**Output**: Branch-safe string (e.g., "add-user-authentication-login")

**Algorithm**:
```
function sanitize_branch_name(text):
    # Step 1: Convert to lowercase
    text = text.to_lowercase()
    
    # Step 2: Replace whitespace with hyphens
    text = regex_replace(text, "\\s+", "-")
    
    # Step 3: Remove special characters (keep alphanumeric, hyphens, underscores)
    text = regex_replace(text, "[^a-z0-9-_]", "")
    
    # Step 4: Remove consecutive hyphens
    text = regex_replace(text, "-+", "-")
    
    # Step 5: Trim hyphens from ends
    text = text.trim("-")
    
    # Step 6: Truncate if too long
    max_length = 50
    if text.length > max_length:
        text = text.substring(0, max_length)
        text = text.trim("-")  # Remove trailing hyphen if truncated
    
    return text
```

---

## 13. Security Requirements

### 13.1 Credential Storage

**Requirements**:
- MUST use OS-provided secure storage (keyring)
- MUST NOT store credentials in plaintext files
- MUST NOT store credentials in config files
- MUST encrypt credentials at rest (OS handles this)
- MUST isolate credentials per project

### 13.2 API Token Handling

**Requirements**:
- MUST use HTTPS for all API calls
- MUST send tokens in Authorization header (not URL)
- MUST NOT log tokens
- MUST NOT display tokens in UI
- MUST validate SSL certificates

### 13.3 OAuth Security

**Requirements**:
- MUST use OAuth 2.0 Device Flow for CLI apps
- MUST use state parameter to prevent CSRF
- MUST NOT include client secret (use public client)
- MUST validate redirect URIs
- MUST timeout authorization requests

### 13.4 File Permissions

**Requirements**:
- Config files: `0600` (rw-------)
- Config directory: `0700` (rwx------)
- MUST NOT include sensitive data in config files

### 13.5 Input Validation

**Requirements**:
- MUST validate all user input
- MUST sanitize file paths
- MUST validate URLs before opening
- MUST sanitize branch names
- MUST validate regex patterns

---

## 14. Platform Requirements

### 14.1 Supported Platforms

#### 14.1.1 Operating Systems
- **macOS**: 10.15+ (Catalina and later)
- **Linux**: Ubuntu 20.04+, Fedora 35+, Debian 11+
- **Windows**: Windows 10+, Windows Server 2019+

#### 14.1.2 CPU Architectures
- **x86_64** (AMD64): Primary support
- **ARM64** (Apple Silicon, ARM servers): Primary support
- **x86** (32-bit): Not supported

### 14.2 Dependencies

#### 14.2.1 Required System Dependencies
- **Git**: 2.0+ (for git operations)
- **OpenSSL**: 1.1.1+ or 3.0+ (for HTTPS)

#### 14.2.2 Optional System Dependencies
- **Chrome**: Latest stable (for Chrome support)
- **Firefox**: Latest stable (for Firefox support)
- **Safari**: Bundled with macOS (for Safari support)

#### 14.2.3 Platform-Specific Dependencies

**macOS**:
- Security framework (for Keychain access)
- CoreFoundation (bundled with OS)

**Linux**:
- libsecret (for Secret Service)
- D-Bus (for Secret Service communication)
- gnome-keyring or kwallet (keyring backend)

**Windows**:
- Windows Credential Manager (built-in)
- Wincred API (built-in)

### 14.3 File System Paths

#### 14.3.1 Configuration Directory

**Location**:
- **macOS/Linux**: `$XDG_CONFIG_HOME/one/` or `~/.config/one/`
- **Windows**: `%APPDATA%\one\`

**Structure**:
```
~/.config/one/
├── config.yml          (optional global config)
└── projects/
    ├── project-a.yml
    ├── project-b.yml
    └── project-c.yml
```

#### 14.3.2 Browser Profiles

**Chrome**:
- **macOS**: `~/Library/Application Support/Google/Chrome/`
- **Linux**: `~/.config/google-chrome/`
- **Windows**: `%LOCALAPPDATA%\Google\Chrome\User Data\`

**Firefox**:
- **macOS**: `~/Library/Application Support/Firefox/`
- **Linux**: `~/.mozilla/firefox/`
- **Windows**: `%APPDATA%\Mozilla\Firefox\`

#### 14.3.3 Keyring Storage

**macOS**:
- **Location**: `~/Library/Keychains/login.keychain-db`
- **Access**: Via Security framework

**Linux**:
- **Location**: Varies by backend
- **Access**: Via Secret Service D-Bus API
- **Backends**: gnome-keyring, kwallet, pass

**Windows**:
- **Location**: Windows Credential Manager
- **Access**: Via Wincred API

### 14.4 Environment Variables

#### 14.4.1 Configuration
- `XDG_CONFIG_HOME`: Override config directory (Linux/macOS)
- `HOME`: User home directory

#### 14.4.2 Authentication (Fallback)
- `GITHUB_TOKEN`: GitHub personal access token
- `GITLAB_TOKEN`: GitLab personal access token
- `JIRA_TOKEN`: Jira API token (email:token format)
- `BITBUCKET_TOKEN`: Bitbucket app password

#### 14.4.3 Git
- `GIT_AUTHOR_NAME`: Git author name
- `GIT_AUTHOR_EMAIL`: Git author email
- `GIT_SSH_COMMAND`: Custom SSH command

---

## 15. Testing Requirements

### 15.1 Unit Tests

#### 15.1.1 Coverage Requirements
- **Minimum coverage**: 80%
- **Critical paths**: 100% (auth, git operations, API calls)

#### 15.1.2 Test Categories

**Configuration Parsing**:
- Valid YAML parsing
- Invalid YAML handling
- Missing required fields
- Schema validation

**Branch Name Parsing**:
- Standard patterns (PROJ-1234)
- Edge cases (no match, multiple matches)
- Invalid regex patterns

**Template Rendering**:
- Variable substitution
- Missing variables
- Special characters
- Multiline templates

**Path Matching**:
- Exact matches
- Prefix matches
- No matches
- Symlink resolution

**Branch Name Sanitization**:
- Special characters removal
- Whitespace handling
- Truncation
- Empty strings

### 15.2 Integration Tests

#### 15.2.1 Git Operations
- Repository discovery
- Branch creation
- Branch checkout
- Push/pull operations
- Stash operations

#### 15.2.2 API Integrations
- GitHub PR creation (mocked)
- GitLab MR creation (mocked)
- Bitbucket PR creation (mocked)
- Jira issue fetching (mocked)

#### 15.2.3 Keyring Operations
- Store token
- Retrieve token
- Delete token
- Handle missing tokens

### 15.3 End-to-End Tests

#### 15.3.1 Complete Workflows

**Workflow 1: New Task**:
```
1. Initialize project (one init)
2. Start task (one start PROJ-1234)
3. Verify branch created
4. Verify ticket fetched
```

**Workflow 2: Create PR**:
```
1. Create branch with ticket ID
2. Make commits
3. Create PR (one pr)
4. Verify PR created
5. Verify browser opened
```

**Workflow 3: Authentication**:
```
1. Login (one login github)
2. Verify token stored
3. Check status (one status)
4. Logout (one logout github)
5. Verify token removed
```

### 15.4 Manual Testing Checklist

#### 15.4.1 Cross-Platform
- [ ] Works on macOS
- [ ] Works on Linux (Ubuntu, Fedora)
- [ ] Works on Windows

#### 15.4.2 Git Providers
- [ ] GitHub PR creation works
- [ ] GitLab MR creation works
- [ ] Bitbucket PR creation works

#### 15.4.3 Browsers
- [ ] Chrome with default profile
- [ ] Chrome with custom profile
- [ ] Firefox with profile
- [ ] Safari (macOS only)

#### 15.4.4 Ticket Systems
- [ ] Jira ticket fetching
- [ ] GitHub issue fetching
- [ ] Linear URL generation

#### 15.4.5 Authentication
- [ ] GitHub OAuth device flow
- [ ] GitLab token authentication
- [ ] Jira token authentication
- [ ] Bitbucket app password

#### 15.4.6 Error Handling
- [ ] Invalid config file
- [ ] No git repository
- [ ] Working directory not clean
- [ ] Network errors
- [ ] API rate limits
- [ ] Invalid credentials

---

## 16. Implementation Guidelines

### 16.1 Language-Specific Recommendations

#### 16.1.1 Rust (Current Implementation)
- Use `clap` for CLI parsing
- Use `tokio` for async operations
- Use `git2` for git operations
- Use `keyring` for credential storage
- Use `serde` + `serde_yaml` for config

#### 16.1.2 Go
- Use `cobra` for CLI
- Use `go-git` for git operations
- Use `keyring` library
- Use `gopkg.in/yaml.v3` for YAML
- Use `go-github` for GitHub API

#### 16.1.3 Python
- Use `click` or `typer` for CLI
- Use `GitPython` for git operations
- Use `keyring` library
- Use `PyYAML` for config
- Use `requests` for HTTP

#### 16.1.4 Node.js/TypeScript
- Use `commander` for CLI
- Use `simple-git` for git operations
- Use `keytar` for keyring
- Use `js-yaml` for config
- Use `axios` for HTTP

#### 16.1.5 Java/Kotlin
- Use `picocli` for CLI
- Use `JGit` for git operations
- Use platform-specific keyring APIs
- Use `SnakeYAML` for config
- Use `OkHttp` for HTTP

### 16.2 Performance Targets

- **Startup time**: < 100ms
- **Config discovery**: < 10ms
- **Git operations**: < 500ms (local)
- **API calls**: < 2s (network dependent)
- **Memory usage**: < 50MB

### 16.3 Binary Size Targets

- **Optimal**: < 5MB
- **Acceptable**: < 10MB
- **Maximum**: < 20MB

### 16.4 Code Organization

```
src/
├── main                    # Entry point
├── cli/                    # CLI parsing
├── commands/               # Command implementations
│   ├── start
│   ├── pr
│   ├── init
│   ├── ticket
│   ├── config
│   └── login
├── config/                 # Configuration handling
├── auth/                   # Authentication
│   ├── oauth
│   └── keyring
├── git/                    # Git operations
├── browser/                # Browser integration
│   └── profiles
├── api/                    # API clients
│   ├── github
│   ├── gitlab
│   ├── bitbucket
│   └── jira
├── template/               # Template rendering
└── error/                  # Error types
```

---

## 17. Future Enhancements

### 17.1 Planned Features

#### 17.1.1 High Priority
- Draft PR support
- PR review commands
- Shell completions (bash/zsh/fish)
- Auto-update mechanism

#### 17.1.2 Medium Priority
- Linear API integration (fetch ticket titles)
- Azure DevOps support
- Multi-repo PR creation
- Custom PR templates per project

#### 17.1.3 Low Priority
- TUI (Terminal UI) mode
- Commit message templates
- Git hooks installer
- Slack/Discord notifications
- Time tracking integration

### 17.2 Extension Points

#### 17.2.1 Plugin System
Consider adding plugin support:
- Custom git providers
- Custom ticket systems
- Custom browsers
- Custom commands

#### 17.2.2 Configuration Hooks
Allow scripts to run at certain points:
- Before PR creation
- After PR creation
- On branch creation
- On ticket open

---

## 18. Glossary

| Term | Definition |
|------|------------|
| **Project** | A collection of paths with shared configuration |
| **Provider** | Git hosting service (GitHub, GitLab, Bitbucket) |
| **Ticket System** | Issue tracking system (Jira, Linear, GitHub Issues) |
| **Browser Profile** | Isolated browser session with separate cookies/cache |
| **Keyring** | OS-provided secure credential storage |
| **OAuth** | Industry-standard protocol for authorization |
| **Device Flow** | OAuth flow designed for CLI apps without browser redirect |
| **Template** | Text with placeholders for variable substitution |
| **Sanitize** | Clean text to be safe for specific use (e.g., branch names) |
| **Fast-Forward** | Git merge that only moves branch pointer forward |
| **Stash** | Temporarily save uncommitted changes |
| **Base Branch** | Target branch for PRs (usually `main` or `master`) |

---

## 19. References

### 19.1 Standards & RFCs
- [RFC 8628 - OAuth 2.0 Device Authorization Grant](https://datatracker.ietf.org/doc/html/rfc8628)
- [RFC 6749 - OAuth 2.0 Authorization Framework](https://datatracker.ietf.org/doc/html/rfc6749)
- [RFC 7519 - JSON Web Token (JWT)](https://datatracker.ietf.org/doc/html/rfc7519)

### 19.2 API Documentation
- [GitHub REST API](https://docs.github.com/en/rest)
- [GitLab REST API](https://docs.gitlab.com/ee/api/)
- [Bitbucket REST API](https://developer.atlassian.com/cloud/bitbucket/rest/)
- [Jira REST API](https://developer.atlassian.com/cloud/jira/platform/rest/v2/)

### 19.3 Platform Documentation
- [Secret Service API (Linux)](https://specifications.freedesktop.org/secret-service/)
- [macOS Keychain Services](https://developer.apple.com/documentation/security/keychain_services)
- [Windows Credential Manager](https://docs.microsoft.com/en-us/windows/win32/api/wincred/)

### 19.4 Git Documentation
- [Git Documentation](https://git-scm.com/doc)
- [libgit2 Documentation](https://libgit2.org/)

---

## Appendix A: Example Configurations

### A.1 GitHub + Jira

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
  profile: "Work"

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
    
    ## Jira Ticket
    {ticket_url}

branch_patterns:
  ticket_id: "^(ACME-\\d+)"
```

### A.2 GitLab + Linear

```yaml
version: 1
project:
  name: "Beta Startup"
  paths:
    - "/Users/jane/beta-app"

git:
  provider: gitlab
  remote: origin
  base_branch: develop
  
  gitlab:
    project_id: 42
    token_env: GITLAB_TOKEN

browser:
  type: firefox
  profile: "default-release"

ticket:
  system: linear
  base_url: https://linear.app/beta

templates:
  pr_title: "{ticket_id} - {branch_name}"
  pr_body: "Closes {ticket_url}"

branch_patterns:
  ticket_id: "^(BET-\\d+)"
```

### A.3 Bitbucket + Minimal Config

```yaml
version: 1
project:
  name: "Personal Project"
  paths:
    - "/Users/alex/personal"

git:
  provider: bitbucket
  remote: origin
  base_branch: main
  
  bitbucket:
    workspace: alex-workspace
    repo_slug: personal-repo
    token_env: BITBUCKET_TOKEN

browser:
  type: safari

# No ticket system
# No templates (use defaults)
# No branch patterns (no ticket extraction)
```

---

## Document History

| Version | Date | Changes |
|---------|------|---------|
| 1.0 | 2025-10-02 | Initial specification |

---

**End of Specification**

This specification is designed to be implementation-agnostic and provides all necessary details to implement One CLI in any programming language while maintaining compatibility and feature parity.
