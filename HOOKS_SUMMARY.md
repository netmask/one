# Hooks Feature Summary

## What Are Hooks?

Hooks allow you to run arbitrary shell commands before and after creating pull requests. This enables powerful automation and validation workflows.

## Why Hooks?

### The Problem

**Before Hooks:**
```bash
# Manual workflow - easy to forget steps!
$ bundle exec rubocop  # Did you remember to lint?
$ bundle exec rspec    # Did you run tests?
$ one pr              # Create PR
$ curl -X POST ...    # Manually notify team
```

### The Solution

**With Hooks:**
```yaml
hooks:
  before_pr:
    - name: "Lint"
      command: "bundle exec rubocop"
      fail_on_error: true
    
    - name: "Test"
      command: "bundle exec rspec"
      fail_on_error: true
  
  after_pr:
    - name: "Notify"
      command: 'curl -X POST $SLACK_WEBHOOK ...'
      fail_on_error: false
```

```bash
# Just one command - hooks run automatically!
$ one pr
```

## Key Features

### ✅ Validation (before_pr)

**Stop PR creation if checks fail:**
- Linting fails → No PR
- Tests fail → No PR
- Build fails → No PR

**Example:**
```yaml
before_pr:
  - name: "Run tests"
    command: "npm test"
    fail_on_error: true  # Stop if tests fail
```

### 🤖 Automation (after_pr)

**Auto-run after successful PR creation:**
- Send Slack notifications
- Deploy preview environments
- Update Jira tickets
- Log to analytics

**Example:**
```yaml
after_pr:
  - name: "Deploy preview"
    command: "./deploy-preview.sh"
    fail_on_error: false  # Just warn if fails
```

### 📊 Real-time Feedback

See exactly what's happening:

```
⚡ Running before_pr hooks...

  [1/3] Lint with RuboCop
        Check code style and best practices
        $ bundle exec rubocop

  ✓ Success (took 1.2s)

  [2/3] Run tests
        Run the full test suite
        $ bundle exec rspec

  ✓ Success (took 15.3s)

  [3/3] Security audit
        $ npm audit --audit-level=high

  ✓ Success (took 0.8s)

✓ All before_pr hooks completed

Creating pull request...
```

### 🎯 Smart Error Handling

**Fail fast:**
```yaml
- name: "Critical check"
  command: "npm test"
  fail_on_error: true  # Stop everything if this fails
```

**Continue on error:**
```yaml
- name: "Optional check"
  command: "npm run coverage"
  fail_on_error: false  # Just warn, don't stop
```

## Common Use Cases

### 1. Ruby on Rails - Quality Gate

```yaml
hooks:
  before_pr:
    - name: "RuboCop"
      command: "bundle exec rubocop"
      fail_on_error: true
    
    - name: "RSpec"
      command: "bundle exec rspec"
      fail_on_error: true
    
    - name: "Brakeman security scan"
      command: "bundle exec brakeman"
      fail_on_error: true
```

**Result**: No PR unless code is clean, tested, and secure.

### 2. JavaScript - Pre-commit Validation

```yaml
hooks:
  before_pr:
    - name: "ESLint"
      command: "npm run lint"
      fail_on_error: true
    
    - name: "TypeScript"
      command: "npm run typecheck"
      fail_on_error: true
    
    - name: "Jest"
      command: "npm test"
      fail_on_error: true
    
    - name: "Build"
      command: "npm run build"
      fail_on_error: true
```

**Result**: PR only created if code is valid, typed, tested, and builds.

### 3. Python - Complete Pipeline

```yaml
hooks:
  before_pr:
    - name: "Black formatting"
      command: "black --check ."
      fail_on_error: true
    
    - name: "Flake8 linting"
      command: "flake8 ."
      fail_on_error: true
    
    - name: "MyPy type checking"
      command: "mypy ."
      fail_on_error: true
    
    - name: "Pytest with coverage"
      command: "pytest --cov --cov-report=term-missing"
      fail_on_error: true
  
  after_pr:
    - name: "Generate docs"
      command: "make docs"
      fail_on_error: false
    
    - name: "Deploy staging"
      command: "./scripts/deploy-staging.sh"
      fail_on_error: false
```

**Result**: Validated code + automated deployment.

### 4. Go - Fast Checks

```yaml
hooks:
  before_pr:
    - name: "gofmt"
      command: "test -z $(gofmt -l .)"
      fail_on_error: true
    
    - name: "golangci-lint"
      command: "golangci-lint run"
      fail_on_error: true
    
    - name: "go test"
      command: "go test -v -race -cover ./..."
      fail_on_error: true
```

**Result**: Formatted, linted, tested code.

### 5. Team Automation

```yaml
hooks:
  after_pr:
    - name: "Slack notification"
      description: "Notify #engineering channel"
      command: |
        curl -X POST $SLACK_WEBHOOK \
          -H 'Content-Type: application/json' \
          -d '{"text":"🚀 New PR: '"$BRANCH_NAME"'"}'
      fail_on_error: false
    
    - name: "Jira update"
      description: "Link PR to Jira ticket"
      command: "./scripts/link-pr-to-jira.sh"
      fail_on_error: false
    
    - name: "Deploy preview"
      description: "Deploy to preview environment"
      command: "vercel deploy --prebuilt"
      fail_on_error: false
```

**Result**: Team notified, tickets updated, preview deployed - all automatic!

## Configuration

Add to your `~/.config/one/projects/project.yml`:

```yaml
hooks:
  before_pr:
    - name: "Name of check"
      description: "What it does (optional)"
      command: "shell command to run"
      fail_on_error: true  # true = stop on fail, false = warn only
  
  after_pr:
    - name: "Name of automation"
      command: "another command"
      fail_on_error: false  # Usually false for after_pr
```

## Best Practices

### 1. Order by Speed

Fast checks first, slow checks last:

```yaml
before_pr:
  - name: "Format check"  # < 1s
    command: "prettier --check ."
    fail_on_error: true
  
  - name: "Lint"  # ~5s
    command: "eslint ."
    fail_on_error: true
  
  - name: "Tests"  # ~30s
    command: "npm test"
    fail_on_error: true
```

**Why?** Fail fast on simple issues.

### 2. Clear Names

Good:
```yaml
- name: "RuboCop: Check code style"
```

Bad:
```yaml
- name: "Check"
```

### 3. Critical = fail_on_error: true

```yaml
- name: "Tests"
  command: "npm test"
  fail_on_error: true  # MUST pass
```

### 4. Nice-to-have = fail_on_error: false

```yaml
- name: "Generate docs"
  command: "make docs"
  fail_on_error: false  # Would be nice, but not critical
```

### 5. Use Scripts for Complex Logic

Bad:
```yaml
command: "cd app && npm test && cd ../api && go test && ..."
```

Good:
```yaml
command: "./scripts/run-all-tests.sh"
```

## Impact

### Time Saved

**Manual process:**
- Remember to lint: 30s
- Remember to test: 30s
- Run commands manually: 2min
- Create PR: 1min
- Remember to notify: 1min
- **Total: 5 minutes**

**With hooks:**
- Run one pr: automatic
- **Total: 0 seconds of manual work**

### Quality Improved

- ✅ Never forget to lint
- ✅ Never forget to test
- ✅ Never create broken PRs
- ✅ Team always notified
- ✅ Consistent process

### Flexibility

Different projects, different needs:

**Project A (strict):**
```yaml
before_pr:
  - Lint (fail)
  - Test (fail)
  - Security scan (fail)
  - Coverage check (fail)
```

**Project B (lenient):**
```yaml
before_pr:
  - Quick lint (fail)
  - Basic tests (fail)
```

**Project C (custom):**
```yaml
before_pr:
  - Custom validation script (fail)
after_pr:
  - Deploy to staging (warn)
  - Notify via webhook (warn)
```

## Technical Details

### Shell Execution

Hooks run in a shell:
- **macOS/Linux**: `/bin/sh -c "command"`
- **Windows**: `cmd /C "command"`

This means you can use:
- Pipes: `npm test | grep -v deprecated`
- Environment variables: `$SLACK_WEBHOOK`
- Conditionals: `if [ "$CI" = "true" ]; then ...; fi`
- Scripts: `./scripts/custom.sh`

### Working Directory

Hooks run in your current directory (the project root).

### Environment Variables

All your environment variables are available:
- `$USER`, `$HOME`, etc.
- Custom vars: `$SLACK_WEBHOOK`, `$JIRA_TOKEN`
- One CLI provides: `$BRANCH_NAME` (coming soon)

### Exit Codes

- **0** = Success
- **Non-zero** = Failure

## Examples Directory

See [examples/hooks-example.yml](examples/hooks-example.yml) for a complete working example.

## Full Documentation

See [HOOKS.md](HOOKS.md) for complete documentation including:
- All properties
- More examples (10+ languages)
- Advanced patterns
- Troubleshooting
- Best practices

---

**Hooks make One CLI not just a PR tool, but a complete quality gate and automation platform!** 🚀
