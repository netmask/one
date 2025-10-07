# Hooks Documentation

Hooks allow you to run arbitrary commands before and after creating pull requests. This is perfect for running linters, tests, formatters, or any custom automation.

## Configuration

Add a `hooks` section to your project configuration:

```yaml
hooks:
  before_pr:
    - name: "Hook name"
      description: "Optional description"
      command: "command to run"
      fail_on_error: true  # Stop PR creation if this fails
  
  after_pr:
    - name: "Another hook"
      command: "another command"
      fail_on_error: false  # Continue even if this fails
```

## Hook Types

### `before_pr` Hooks

Run **before** the PR is created. Perfect for validation:
- ✅ Linting (RuboCop, ESLint, Pylint)
- ✅ Tests (RSpec, Jest, PyTest)
- ✅ Type checking (TypeScript, MyPy)
- ✅ Code formatting (Prettier, Black, gofmt)
- ✅ Security scanning
- ✅ Build validation

**If a `before_pr` hook fails with `fail_on_error: true`, the PR will NOT be created.**

### `after_pr` Hooks

Run **after** the PR is successfully created. Perfect for automation:
- ✅ Send notifications (Slack, Discord, Email)
- ✅ Update project management (Jira, Linear)
- ✅ Deploy preview environments
- ✅ Update documentation
- ✅ Create QA tickets
- ✅ Log to analytics

**`after_pr` hooks don't stop the PR creation, they just warn if they fail.**

## Hook Properties

| Property | Required | Description |
|----------|----------|-------------|
| `name` | ✅ Yes | Human-readable name for the hook |
| `command` | ✅ Yes | Shell command to execute |
| `description` | No | Explanation of what the hook does |
| `fail_on_error` | ✅ Yes | If true, stop execution on failure |

## Examples

### Ruby on Rails Project

```yaml
hooks:
  before_pr:
    - name: "Lint with RuboCop"
      description: "Check code style and best practices"
      command: "bundle exec rubocop"
      fail_on_error: true
    
    - name: "Run tests"
      description: "Run the full test suite"
      command: "bundle exec rspec"
      fail_on_error: true
    
    - name: "Check migrations"
      description: "Verify database migrations are clean"
      command: "bundle exec rails db:migrate:status"
      fail_on_error: false
  
  after_pr:
    - name: "Notify team"
      description: "Send Slack notification"
      command: 'curl -X POST $SLACK_WEBHOOK -d "{\"text\":\"PR created by $USER\"}"'
      fail_on_error: false
```

### JavaScript/TypeScript Project

```yaml
hooks:
  before_pr:
    - name: "Lint code"
      description: "Run ESLint"
      command: "npm run lint"
      fail_on_error: true
    
    - name: "Type check"
      description: "Run TypeScript compiler"
      command: "npm run typecheck"
      fail_on_error: true
    
    - name: "Run tests"
      description: "Run Jest test suite"
      command: "npm test -- --coverage"
      fail_on_error: true
    
    - name: "Build"
      description: "Verify production build works"
      command: "npm run build"
      fail_on_error: true
  
  after_pr:
    - name: "Deploy preview"
      description: "Deploy to Vercel preview"
      command: "./scripts/deploy-preview.sh"
      fail_on_error: false
```

### Python Project

```yaml
hooks:
  before_pr:
    - name: "Format check"
      description: "Check code formatting with Black"
      command: "black --check ."
      fail_on_error: true
    
    - name: "Lint"
      description: "Run Flake8"
      command: "flake8 ."
      fail_on_error: true
    
    - name: "Type check"
      description: "Run MyPy"
      command: "mypy ."
      fail_on_error: true
    
    - name: "Tests"
      description: "Run pytest"
      command: "pytest --cov"
      fail_on_error: true
  
  after_pr:
    - name: "Update docs"
      description: "Regenerate API documentation"
      command: "make docs"
      fail_on_error: false
```

### Go Project

```yaml
hooks:
  before_pr:
    - name: "Format"
      description: "Check code formatting"
      command: "gofmt -l . | grep -q . && exit 1 || exit 0"
      fail_on_error: true
    
    - name: "Lint"
      description: "Run golangci-lint"
      command: "golangci-lint run"
      fail_on_error: true
    
    - name: "Tests"
      description: "Run tests with coverage"
      command: "go test -v -cover ./..."
      fail_on_error: true
    
    - name: "Build"
      description: "Verify build"
      command: "go build -o /tmp/app"
      fail_on_error: true
```

## Advanced Examples

### Conditional Hooks

Run different commands based on environment variables:

```yaml
hooks:
  before_pr:
    - name: "Environment-aware tests"
      command: "if [ \"$CI\" = \"true\" ]; then npm run test:ci; else npm test; fi"
      fail_on_error: true
```

### Multi-step Hooks

Combine multiple commands:

```yaml
hooks:
  before_pr:
    - name: "Setup and test"
      description: "Install deps and run tests"
      command: "npm ci && npm test"
      fail_on_error: true
```

### With External Scripts

Call custom scripts:

```yaml
hooks:
  before_pr:
    - name: "Pre-PR checks"
      description: "Run all quality checks"
      command: "./scripts/pre-pr-checks.sh"
      fail_on_error: true
  
  after_pr:
    - name: "Post-PR automation"
      command: "./scripts/post-pr-automation.sh \"$PR_URL\""
      fail_on_error: false
```

### Slack Notifications

```yaml
hooks:
  after_pr:
    - name: "Notify Slack"
      description: "Send PR notification to team channel"
      command: |
        curl -X POST $SLACK_WEBHOOK \
          -H 'Content-Type: application/json' \
          -d "{
            \"text\": \"🚀 New PR created\",
            \"attachments\": [{
              \"title\": \"$BRANCH_NAME\",
              \"color\": \"good\"
            }]
          }"
      fail_on_error: false
```

### Jira Integration

```yaml
hooks:
  after_pr:
    - name: "Link PR to Jira"
      description: "Add PR link as comment to Jira ticket"
      command: "./scripts/link-pr-to-jira.sh \"$TICKET_ID\" \"$PR_URL\""
      fail_on_error: false
```

## Output Display

Hooks show their execution in real-time:

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

  [3/3] Check migrations
        $ bundle exec rails db:migrate:status

  ✓ Success (took 0.5s)

✓ All before_pr hooks completed
```

## Error Handling

### When a hook fails with `fail_on_error: true`:

```
⚡ Running before_pr hooks...

  [1/2] Lint code
        $ npm run lint

  ✗ Failed

    src/app.js:10:5 - error: Unexpected console statement

  Duration: 0.8s

Error: before_pr hook 'Lint code' failed
```

**The PR is NOT created.**

### When a hook fails with `fail_on_error: false`:

```
⚡ Running after_pr hooks...

  [1/1] Deploy preview
        $ ./deploy-preview.sh

  ✗ Failed

⚠️ Hook 'Deploy preview' failed but continuing (fail_on_error: false)

✓ All after_pr hooks completed
```

**The PR is already created, and execution continues.**

## Best Practices

### 1. Fast Feedback

Put fast checks first:
```yaml
before_pr:
  - name: "Format" # Fast (< 1s)
    command: "prettier --check ."
    fail_on_error: true
  
  - name: "Lint" # Medium (~5s)
    command: "eslint ."
    fail_on_error: true
  
  - name: "Tests" # Slow (~30s)
    command: "npm test"
    fail_on_error: true
```

### 2. Clear Names and Descriptions

Good:
```yaml
- name: "TypeScript type check"
  description: "Ensure no type errors before PR"
  command: "tsc --noEmit"
```

Bad:
```yaml
- name: "Check"
  command: "tsc --noEmit"
```

### 3. Fail Fast

Set `fail_on_error: true` for critical checks:
```yaml
- name: "Tests"
  command: "npm test"
  fail_on_error: true  # Don't create PR if tests fail
```

### 4. Be Lenient with After Hooks

`after_pr` hooks should usually be non-blocking:
```yaml
after_pr:
  - name: "Send notification"
    command: "curl ..."
    fail_on_error: false  # Don't fail if notification fails
```

### 5. Use Scripts for Complex Logic

Instead of:
```yaml
command: "cd frontend && npm test && cd ../backend && go test && ..."
```

Use:
```yaml
command: "./scripts/run-all-tests.sh"
```

## Environment Variables

Hooks have access to all environment variables. Useful ones:

- `$USER` - Current user
- `$HOME` - Home directory
- Any variables you set (e.g., `$SLACK_WEBHOOK`, `$JIRA_TOKEN`)

## Exit Codes

- **Exit code 0** = Success
- **Any other exit code** = Failure

Example:
```bash
#!/bin/bash
# scripts/check-todos.sh

if grep -r "TODO" src/; then
  echo "Found TODO comments - please address them"
  exit 1  # Fail
fi

exit 0  # Success
```

## Debugging Hooks

### 1. Test Hooks Manually

Before adding to config, test the command:
```bash
bundle exec rubocop  # Does this work?
```

### 2. Check Exit Codes

```bash
npm test
echo $?  # Should be 0 for success
```

### 3. Verbose Mode

Add verbosity to your commands:
```yaml
command: "npm test -- --verbose"
```

### 4. Echo Commands

See what's being executed:
```yaml
command: "set -x && npm test"
```

## Common Use Cases

### Prevent PRs with failing tests

```yaml
before_pr:
  - name: "Tests must pass"
    command: "npm test"
    fail_on_error: true
```

### Auto-format before PR

```yaml
before_pr:
  - name: "Format code"
    command: "prettier --write . && git add -u"
    fail_on_error: false
```

### Notify team

```yaml
after_pr:
  - name: "Slack notification"
    command: "./scripts/notify-slack.sh"
    fail_on_error: false
```

### Security scanning

```yaml
before_pr:
  - name: "Security audit"
    command: "npm audit --audit-level=high"
    fail_on_error: true
```

### Update changelog

```yaml
after_pr:
  - name: "Update CHANGELOG"
    command: "./scripts/update-changelog.sh"
    fail_on_error: false
```

## Troubleshooting

### Hook not running?

1. Check YAML syntax
2. Verify `hooks` is at root level
3. Ensure proper indentation

### Command not found?

1. Check PATH includes the tool
2. Use full path: `/usr/local/bin/rubocop`
3. Or activate environment first:
   ```yaml
   command: "source ~/.bashrc && rubocop"
   ```

### Hook always fails?

1. Test command manually
2. Check exit code: `echo $?`
3. Look at error output in hook execution

---

## See Also

- [Examples](examples/hooks-example.yml) - Complete hook examples
- [Configuration](README.md#configuration) - Project configuration guide
- [Template Variables](README.md#template-variables) - Available template variables
