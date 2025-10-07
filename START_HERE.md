# 👋 Start Here - New to One CLI?

## 🚀 Quick Links

**Just want to use it?** → [README.md](README.md)  
**Want to get started in 5 minutes?** → [QUICKSTART.md](QUICKSTART.md)  
**Want to see all features?** → [FEATURES.md](FEATURES.md)  
**Want to see hooks in action?** → [HOOKS.md](HOOKS.md)

---

## 📖 Documentation Navigation

### 👤 For Users

**Start Here:**
1. **[README.md](README.md)** ⭐ - Main documentation
   - What is One CLI?
   - Installation
   - Quick start guide
   - Feature overview

2. **[QUICKSTART.md](QUICKSTART.md)** - 5-minute tutorial
   - Step-by-step setup
   - First commands
   - Real examples

3. **[FEATURES.md](FEATURES.md)** - Complete feature showcase
   - Auto-detection
   - Browser profiles
   - Hooks system
   - OAuth authentication

**Specific Features:**
- **[HOOKS.md](HOOKS.md)** - Run commands before/after PR
- **[examples/](examples/)** - Configuration examples

### 👨‍💻 For Developers

**Technical Details:**
1. **[SPECIFICATION.md](SPECIFICATION.md)** - Original specification
2. **[IMPLEMENTATION.md](IMPLEMENTATION.md)** - How it's built
3. **[CONTRIBUTING.md](CONTRIBUTING.md)** - How to contribute

**Project Status:**
- **[PROJECT_STATUS.md](PROJECT_STATUS.md)** - What's complete
- **[CHANGELOG.md](CHANGELOG.md)** - Version history
- **[RELEASE_NOTES.md](RELEASE_NOTES.md)** - v0.2.0 details

### 🗺️ For Navigation

- **[INDEX.md](INDEX.md)** - Complete file index
- **[OVERVIEW.md](OVERVIEW.md)** - High-level overview
- **[FINAL_SUMMARY.md](FINAL_SUMMARY.md)** - Implementation summary
- **[SUCCESS.md](SUCCESS.md)** - Success metrics

---

## ⚡ 30-Second Overview

**One CLI** is a unified command-line tool for freelancers working across multiple projects.

**Core workflow:**
```bash
cd /your/project
one init     # Setup in < 2 minutes (auto-detects everything!)
one start TICKET-123
# ... make changes ...
one pr       # Create PR (runs hooks automatically!)
```

**Key features:**
- 🤖 Auto-detects Git provider, owner, repo
- 👤 Smart browser profiles (shows emails)
- 🪝 Hooks for linting, testing, automation
- 🔐 GitHub OAuth (no manual tokens)
- 🎨 Beautiful TUI (Huh + Bubble Tea + Glamour)

---

## 🎯 What Problem Does It Solve?

### The Problem

As a freelancer working on multiple client projects:
- Different Git providers (GitHub, GitLab, Bitbucket)
- Different ticket systems (Jira, Linear)
- Different workflows
- Different credentials
- Different browser profiles

**Result:** Mental overhead, context switching, mistakes.

### The Solution

**One workflow for everything:**
- `one start` - Always works the same way
- `one pr` - Always works the same way
- Automatic project detection
- Automatic credential selection
- Automatic browser profile

**Result:** Zero mental overhead. Just work.

---

## 🎨 What Does It Look Like?

### Setup (Auto-Detected)

```
$ one init

✓ Detected Git remote: github (acme-corp/app)

┌─ Project Name ──────────────┐
│ Acme Corp                   │
└─────────────────────────────┘

┌─ Select Chrome Profile ─────┐
│ > Work (john@company.com)   │
│   Personal (john@gmail.com) │
└─────────────────────────────┘

✓ Configuration saved!
```

### Working (With Hooks)

```
$ one pr

⚡ Running before_pr hooks...

  [1/2] Lint code
  ✓ Success (1.2s)

  [2/2] Run tests
  ✓ Success (15.3s)

✓ All hooks completed

  ✓ PR created!
  
Done! 🚀
```

---

## 🛠️ Installation

### Option 1: Build from Source (Recommended)

```bash
git clone https://github.com/yourusername/one.git
cd one
make build
sudo make install
```

### Option 2: Using Go

```bash
go install github.com/yourusername/one@latest
```

### Option 3: Manual

```bash
go build -o one
sudo mv one /usr/local/bin/
```

---

## 🎓 Learning Path

### Level 1: Basic User (15 minutes)

1. Read [README.md](README.md) overview
2. Follow [QUICKSTART.md](QUICKSTART.md)
3. Run `one init` in your project
4. Try `one start` and `one pr`

### Level 2: Power User (30 minutes)

1. Read [FEATURES.md](FEATURES.md)
2. Set up [hooks](HOOKS.md) for your workflow
3. Configure templates
4. Set up multiple projects

### Level 3: Advanced (1 hour)

1. Read [IMPLEMENTATION.md](IMPLEMENTATION.md)
2. Understand the architecture
3. Contribute via [CONTRIBUTING.md](CONTRIBUTING.md)
4. Add new providers/features

---

## 🎯 Choose Your Path

### I want to...

**...use it right now**
→ [QUICKSTART.md](QUICKSTART.md)

**...see what it can do**
→ [FEATURES.md](FEATURES.md)

**...understand hooks**
→ [HOOKS.md](HOOKS.md)

**...see examples**
→ [examples/](examples/)

**...understand how it works**
→ [IMPLEMENTATION.md](IMPLEMENTATION.md)

**...contribute**
→ [CONTRIBUTING.md](CONTRIBUTING.md)

**...see all files**
→ [INDEX.md](INDEX.md)

---

## ❓ FAQ

**Q: Do I need to configure everything manually?**  
A: No! `one init` auto-detects your Git remote and pre-fills everything.

**Q: How does it know which project I'm in?**  
A: It matches your current directory to configured project paths.

**Q: Is it secure?**  
A: Yes! Uses OS keyring (Keychain/Secret Service/Credential Manager).

**Q: Does it work with my setup?**  
A: Probably! Supports GitHub/GitLab/Bitbucket and macOS/Linux/Windows.

**Q: Can I use it without hooks?**  
A: Yes! Hooks are completely optional.

**Q: What about browser profiles?**  
A: It auto-detects them and shows emails so you know which is which.

---

## 🎉 Ready?

```bash
cd /your/project
one init
```

That's it! Welcome to One CLI. 🚀

---

## 📞 Need Help?

- **Documentation** - You're reading it!
- **Examples** - See [examples/](examples/)
- **Issues** - File on GitHub
- **Questions** - Check the FAQ above

---

**Made with ❤️ for freelancers**

[⭐ Star on GitHub](https://github.com/yourusername/one) • [🐛 Report Bug](https://github.com/yourusername/one/issues) • [💡 Request Feature](https://github.com/yourusername/one/issues)
