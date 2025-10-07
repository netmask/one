package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"one/internal/api"
	"one/internal/auth"
	"one/internal/browser"
	"one/internal/config"
	"one/internal/git"
	"one/internal/hooks"
	"one/internal/template"
)

var prCmd = &cobra.Command{
	Use:   "pr",
	Short: "Create and open a pull request",
	Long:  `Creates a pull request on the configured Git provider and opens it in the browser.`,
	RunE:  runPR,
}

func init() {
	rootCmd.AddCommand(prCmd)
	prCmd.Flags().StringP("title", "t", "", "Custom PR title")
	prCmd.Flags().StringP("description", "d", "", "Custom PR description")
	prCmd.Flags().Bool("no-browser", false, "Skip opening browser")
}

type prModel struct {
	cfg         *config.ProjectConfig
	repo        *git.Repository
	customTitle string
	customDesc  string
	noBrowser   bool
	status      string
	err         error
	done        bool
	prURL       string
	branchName  string
	ticketID    string
}

type prCreatedMsg struct {
	url string
	err error
}

func runPR(cmd *cobra.Command, args []string) error {
	customTitle, _ := cmd.Flags().GetString("title")
	customDesc, _ := cmd.Flags().GetString("description")
	noBrowser, _ := cmd.Flags().GetBool("no-browser")

	// Load config
	cfg, err := config.LoadProjectConfig()
	if err != nil {
		return err
	}

	// Open repository
	repo, err := git.OpenRepository()
	if err != nil {
		return err
	}

	// Get current branch for hooks
	branch, err := repo.CurrentBranch()
	if err != nil {
		return err
	}

	// Run before_pr hooks (if any)
	if cfg.Hooks != nil && len(cfg.Hooks.BeforePR) > 0 {
		if err := hooks.ExecuteHooks(cfg.Hooks.BeforePR, "before_pr"); err != nil {
			return err
		}
	}

	m := &prModel{
		cfg:         cfg,
		repo:        repo,
		customTitle: customTitle,
		customDesc:  customDesc,
		noBrowser:   noBrowser,
		status:      "Creating pull request...",
		branchName:  branch,
	}

	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return err
	}

	fm := finalModel.(*prModel)
	if fm.err != nil {
		if fm.err.Error() == "HOOKS_FIRST" {
			// This was just to get branch info, now actually create PR
			return runPRActual(cfg, repo, branch, customTitle, customDesc, noBrowser)
		}
		return fm.err
	}

	// Run after_pr hooks (if any)
	if cfg.Hooks != nil && len(cfg.Hooks.AfterPR) > 0 {
		if err := hooks.ExecuteHooks(cfg.Hooks.AfterPR, "after_pr"); err != nil {
			// Don't fail the whole command if after hooks fail
			fmt.Printf("Warning: after_pr hooks failed: %v\n", err)
		}
	}

	return nil
}

func runPRActual(cfg *config.ProjectConfig, repo *git.Repository, branch, customTitle, customDesc string, noBrowser bool) error {
	// Check if clean
	clean, err := repo.IsClean()
	if err != nil {
		return fmt.Errorf("failed to check git status: %w", err)
	}
	if !clean {
		return fmt.Errorf("working directory is not clean")
	}

	// Extract ticket ID
	var ticketID string
	if cfg.BranchPatterns != nil && cfg.BranchPatterns.TicketID != "" {
		id, err := git.ParseTicketID(branch, cfg.BranchPatterns.TicketID)
		if err == nil {
			ticketID = id
		}
	}

	fmt.Println()
	fmt.Println("Creating pull request...")
	fmt.Println()
	fmt.Printf("  Project: %s\n", cfg.Project.Name)
	fmt.Printf("  Branch: %s\n", branch)
	if ticketID != "" {
		fmt.Printf("  Ticket ID: %s\n", ticketID)
	}
	fmt.Println()

	// Push to remote
	fmt.Println("Pushing to remote...")
	if err := repo.Push(cfg.Git.Remote, branch); err != nil {
		return fmt.Errorf("failed to push: %w", err)
	}
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render("  ✓ Pushed to " + cfg.Git.Remote))
	fmt.Println()

	// Build template context
	ctx := template.Context{
		"ticket_id":   ticketID,
		"branch_name": branch,
		"base_branch": cfg.Git.BaseBranch,
		"date":        template.GetCurrentDate(),
	}

	if cfg.Ticket != nil && ticketID != "" {
		ctx["ticket_url"] = template.BuildTicketURL(cfg.Ticket.System, cfg.Ticket.BaseURL, ticketID)
	}

	// Generate title and body
	title := customTitle
	if title == "" && cfg.Templates != nil {
		title = template.Render(cfg.Templates.PRTitle, ctx)
	}
	if title == "" {
		title = branch
	}

	body := customDesc
	if body == "" && cfg.Templates != nil {
		body = template.Render(cfg.Templates.PRBody, ctx)
	}

	// Get token
	token, err := getTokenForPR(cfg)
	if err != nil {
		return err
	}

	// Create PR based on provider
	fmt.Println("Creating PR...")
	var prURL string
	switch cfg.Git.Provider {
	case "github":
		if cfg.Git.GitHub == nil {
			return fmt.Errorf("GitHub configuration missing")
		}
		client := api.NewGitHubClient(token)
		prURL, err = client.CreatePullRequest(
			cfg.Git.GitHub.Owner,
			cfg.Git.GitHub.Repo,
			title,
			body,
			branch,
			cfg.Git.BaseBranch,
		)
	case "gitlab":
		if cfg.Git.GitLab == nil {
			return fmt.Errorf("GitLab configuration missing")
		}
		client := api.NewGitLabClient(token)
		prURL, err = client.CreateMergeRequest(
			cfg.Git.GitLab.ProjectID,
			title,
			body,
			branch,
			cfg.Git.BaseBranch,
		)
	default:
		return fmt.Errorf("unsupported provider: %s", cfg.Git.Provider)
	}

	if err != nil {
		return err
	}

	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render("  ✓ PR created: " + prURL))
	fmt.Println()

	// Open in browser
	if !noBrowser {
		fmt.Println("Opening in browser...")
		if err := browser.OpenURL(cfg.Browser.Type, cfg.Browser.Profile, prURL); err != nil {
			// Non-fatal
		}
		fmt.Println()
	}

	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true).Render("Done! 🚀"))
	fmt.Println()

	// Run after_pr hooks (if any)
	if cfg.Hooks != nil && len(cfg.Hooks.AfterPR) > 0 {
		if err := hooks.ExecuteHooks(cfg.Hooks.AfterPR, "after_pr"); err != nil {
			// Don't fail the whole command if after hooks fail
			fmt.Printf("Warning: after_pr hooks failed: %v\n", err)
		}
	}

	return nil
}

func getTokenForPR(cfg *config.ProjectConfig) (string, error) {
	// Try keyring first
	token, err := auth.GetToken(cfg.Git.Provider, cfg.Project.Name)
	if err == nil {
		return token.AccessToken, nil
	}

	// Try environment variable
	var envVar string
	switch cfg.Git.Provider {
	case "github":
		if cfg.Git.GitHub != nil {
			envVar = cfg.Git.GitHub.TokenEnv
		}
	case "gitlab":
		if cfg.Git.GitLab != nil {
			envVar = cfg.Git.GitLab.TokenEnv
		}
	case "bitbucket":
		if cfg.Git.Bitbucket != nil {
			envVar = cfg.Git.Bitbucket.TokenEnv
		}
	}

	if envVar != "" {
		if token := os.Getenv(envVar); token != "" {
			return token, nil
		}
	}

	return "", fmt.Errorf("not authenticated with %s", cfg.Git.Provider)
}

func (m *prModel) Init() tea.Cmd {
	return m.createPR()
}

func (m *prModel) createPR() tea.Cmd {
	return func() tea.Msg {
		// Get current branch
		branch, err := m.repo.CurrentBranch()
		if err != nil {
			return prCreatedMsg{err: err}
		}
		m.branchName = branch

		// Check if clean
		clean, err := m.repo.IsClean()
		if err != nil {
			return prCreatedMsg{err: err}
		}
		if !clean {
			return prCreatedMsg{err: fmt.Errorf("working directory is not clean")}
		}

		// Extract ticket ID
		if m.cfg.BranchPatterns != nil && m.cfg.BranchPatterns.TicketID != "" {
			ticketID, err := git.ParseTicketID(branch, m.cfg.BranchPatterns.TicketID)
			if err == nil {
				m.ticketID = ticketID
			}
		}

		// Run before_pr hooks (outside of Bubble Tea to show output)
		// We need to exit the TUI temporarily
		return prCreatedMsg{err: fmt.Errorf("HOOKS_FIRST")}
	}
}

func (m *prModel) createPRAfterHooks() tea.Cmd {
	return func() tea.Msg {
		branch := m.branchName

		// Push to remote
		if err := m.repo.Push(m.cfg.Git.Remote, branch); err != nil {
			return prCreatedMsg{err: fmt.Errorf("failed to push: %w", err)}
		}

		// Build template context
		ctx := template.Context{
			"ticket_id":   m.ticketID,
			"branch_name": branch,
			"base_branch": m.cfg.Git.BaseBranch,
			"date":        template.GetCurrentDate(),
		}

		if m.cfg.Ticket != nil && m.ticketID != "" {
			ctx["ticket_url"] = template.BuildTicketURL(m.cfg.Ticket.System, m.cfg.Ticket.BaseURL, m.ticketID)
		}

		// Generate title and body
		title := m.customTitle
		if title == "" && m.cfg.Templates != nil {
			title = template.Render(m.cfg.Templates.PRTitle, ctx)
		}
		if title == "" {
			title = branch
		}

		body := m.customDesc
		if body == "" && m.cfg.Templates != nil {
			body = template.Render(m.cfg.Templates.PRBody, ctx)
		}

		// Get token
		token, err := m.getToken()
		if err != nil {
			return prCreatedMsg{err: err}
		}

		// Create PR based on provider
		var prURL string
		switch m.cfg.Git.Provider {
		case "github":
			if m.cfg.Git.GitHub == nil {
				return prCreatedMsg{err: fmt.Errorf("GitHub configuration missing")}
			}
			client := api.NewGitHubClient(token)
			prURL, err = client.CreatePullRequest(
				m.cfg.Git.GitHub.Owner,
				m.cfg.Git.GitHub.Repo,
				title,
				body,
				branch,
				m.cfg.Git.BaseBranch,
			)
		case "gitlab":
			if m.cfg.Git.GitLab == nil {
				return prCreatedMsg{err: fmt.Errorf("GitLab configuration missing")}
			}
			client := api.NewGitLabClient(token)
			prURL, err = client.CreateMergeRequest(
				m.cfg.Git.GitLab.ProjectID,
				title,
				body,
				branch,
				m.cfg.Git.BaseBranch,
			)
		default:
			return prCreatedMsg{err: fmt.Errorf("unsupported provider: %s", m.cfg.Git.Provider)}
		}

		if err != nil {
			return prCreatedMsg{err: err}
		}

		return prCreatedMsg{url: prURL}
	}
}

func (m *prModel) getToken() (string, error) {
	// Try keyring first
	token, err := auth.GetToken(m.cfg.Git.Provider, m.cfg.Project.Name)
	if err == nil {
		return token.AccessToken, nil
	}

	// Try environment variable
	var envVar string
	switch m.cfg.Git.Provider {
	case "github":
		if m.cfg.Git.GitHub != nil {
			envVar = m.cfg.Git.GitHub.TokenEnv
		}
	case "gitlab":
		if m.cfg.Git.GitLab != nil {
			envVar = m.cfg.Git.GitLab.TokenEnv
		}
	case "bitbucket":
		if m.cfg.Git.Bitbucket != nil {
			envVar = m.cfg.Git.Bitbucket.TokenEnv
		}
	}

	if envVar != "" {
		if token := os.Getenv(envVar); token != "" {
			return token, nil
		}
	}

	return "", fmt.Errorf("not authenticated with %s", m.cfg.Git.Provider)
}

func (m *prModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			m.done = true
			return m, tea.Quit
		}

	case prCreatedMsg:
		if msg.err != nil {
			m.err = msg.err
			m.done = true
			return m, tea.Quit
		}

		m.prURL = msg.url

		// Open in browser
		if !m.noBrowser {
			if err := browser.OpenURL(m.cfg.Browser.Type, m.cfg.Browser.Profile, msg.url); err != nil {
				// Non-fatal, just continue
			}
		}

		m.status = "Done!"
		m.done = true
		return m, tea.Quit
	}

	return m, nil
}

func (m *prModel) View() string {
	if m.done {
		if m.err != nil {
			errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
			return errorStyle.Render(fmt.Sprintf("Error: %v\n", m.err))
		}

		successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
		titleStyle := lipgloss.NewStyle().Bold(true)

		var s string
		s += titleStyle.Render("Creating pull request...") + "\n\n"
		s += fmt.Sprintf("  Project: %s\n", m.cfg.Project.Name)
		s += fmt.Sprintf("  Branch: %s\n", m.branchName)
		if m.ticketID != "" {
			s += fmt.Sprintf("  Ticket ID: %s\n", m.ticketID)
		}
		s += "\n"
		s += successStyle.Render("  ✓ Pushed to "+m.cfg.Git.Remote) + "\n"
		s += successStyle.Render("  ✓ PR created: "+m.prURL) + "\n\n"
		if !m.noBrowser {
			s += "Opening in browser...\n\n"
		}
		s += successStyle.Render("Done! 🚀") + "\n"

		return s
	}

	return fmt.Sprintf("%s\n\nPress Ctrl+C to cancel", m.status)
}
