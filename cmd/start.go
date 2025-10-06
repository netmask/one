package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"one/internal/api"
	"one/internal/auth"
	"one/internal/config"
	"one/internal/git"
)

var startCmd = &cobra.Command{
	Use:   "start <TICKET-ID>",
	Short: "Start working on a new task",
	Long:  `Creates a new branch and optionally fetches ticket information.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runStart,
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringP("description", "d", "", "Custom branch description")
}

type startModel struct {
	ticketID    string
	description string
	cfg         *config.ProjectConfig
	repo        *git.Repository
	status      string
	err         error
	done        bool
	branchName  string
}

type ticketFetchedMsg struct {
	title string
	err   error
}

func runStart(cmd *cobra.Command, args []string) error {
	ticketID := args[0]
	description, _ := cmd.Flags().GetString("description")

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

	// Check if clean
	clean, err := repo.IsClean()
	if err != nil {
		return fmt.Errorf("failed to check git status: %w", err)
	}

	// If not clean, ask user what to do
	if !clean {
		var shouldStash bool
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Working directory has uncommitted changes").
					Description("Would you like to stash them?").
					Affirmative("Yes, stash changes").
					Negative("No, cancel").
					Value(&shouldStash),
			),
		)

		if err := form.Run(); err != nil {
			return err
		}

		if !shouldStash {
			return fmt.Errorf("cannot start new task with uncommitted changes")
		}

		// TODO: Implement stash functionality
		fmt.Println("⚠️  Stash functionality not yet implemented. Please commit or stash manually.")
		return fmt.Errorf("working directory not clean")
	}

	// If no description and no ticket system configured, prompt for it
	if description == "" && cfg.Ticket == nil {
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Branch Description").
					Description("Describe what you'll be working on").
					Value(&description).
					Placeholder("add user authentication"),
			),
		)

		if err := form.Run(); err != nil {
			return err
		}
	}

	m := &startModel{
		ticketID:    ticketID,
		description: description,
		cfg:         cfg,
		repo:        repo,
		status:      "Initializing...",
	}

	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return err
	}

	fm := finalModel.(*startModel)
	if fm.err != nil {
		return fm.err
	}

	return nil
}

func (m *startModel) Init() tea.Cmd {
	return m.startTask()
}

func (m *startModel) startTask() tea.Cmd {
	return func() tea.Msg {
		// Checkout base branch
		if err := m.repo.CheckoutBranch(m.cfg.Git.BaseBranch); err != nil {
			return ticketFetchedMsg{err: fmt.Errorf("failed to checkout base branch: %w", err)}
		}

		// Pull latest changes
		if err := m.repo.Pull(m.cfg.Git.Remote, m.cfg.Git.BaseBranch); err != nil {
			// Ignore if already up to date
		}

		// Fetch ticket title if needed
		var title string
		if m.description != "" {
			title = m.description
		} else if m.cfg.Ticket != nil {
			fetchedTitle, err := m.fetchTicketTitle()
			if err != nil {
				// Use ticket ID as fallback
				title = m.ticketID
			} else {
				title = fetchedTitle
			}
		} else {
			title = m.ticketID
		}

		return ticketFetchedMsg{title: title}
	}
}

func (m *startModel) fetchTicketTitle() (string, error) {
	if m.cfg.Ticket == nil {
		return "", fmt.Errorf("no ticket system configured")
	}

	// Try to get token
	token, err := auth.GetToken(m.cfg.Ticket.System, m.cfg.Project.Name)
	if err != nil {
		// Try environment variable
		if m.cfg.Ticket.Jira != nil && m.cfg.Ticket.Jira.TokenEnv != "" {
			if envToken := os.Getenv(m.cfg.Ticket.Jira.TokenEnv); envToken != "" {
				token = &auth.Token{AccessToken: envToken}
			}
		}
	}

	if token == nil {
		return "", fmt.Errorf("not authenticated with %s", m.cfg.Ticket.System)
	}

	switch m.cfg.Ticket.System {
	case "jira":
		client := api.NewJiraClient(m.cfg.Ticket.BaseURL, token.AccessToken)
		return client.GetIssue(m.ticketID)
	default:
		return "", fmt.Errorf("ticket system %s not supported for fetching", m.cfg.Ticket.System)
	}
}

func (m *startModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			m.done = true
			return m, tea.Quit
		}

	case ticketFetchedMsg:
		if msg.err != nil {
			m.err = msg.err
			m.done = true
			return m, tea.Quit
		}

		// Create branch name
		sanitized := git.SanitizeBranchName(msg.title)
		m.branchName = fmt.Sprintf("%s-%s", m.ticketID, sanitized)

		// Create and checkout branch
		if err := m.repo.CreateBranch(m.branchName); err != nil {
			m.err = fmt.Errorf("failed to create branch: %w", err)
			m.done = true
			return m, tea.Quit
		}

		if err := m.repo.CheckoutBranch(m.branchName); err != nil {
			m.err = fmt.Errorf("failed to checkout branch: %w", err)
			m.done = true
			return m, tea.Quit
		}

		m.status = "Done!"
		m.done = true
		return m, tea.Quit
	}

	return m, nil
}

func (m *startModel) View() string {
	if m.done {
		if m.err != nil {
			errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
			return errorStyle.Render(fmt.Sprintf("Error: %v\n", m.err))
		}

		successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
		titleStyle := lipgloss.NewStyle().Bold(true)

		var s string
		s += titleStyle.Render("Starting new task...") + "\n\n"
		s += fmt.Sprintf("  Project: %s\n\n", m.cfg.Project.Name)
		s += fmt.Sprintf("Checking out base branch...\n")
		s += fmt.Sprintf("  Branch: %s\n", m.cfg.Git.BaseBranch)
		s += successStyle.Render("  ✓ Checked out "+m.cfg.Git.BaseBranch) + "\n\n"
		s += fmt.Sprintf("Pulling latest changes...\n")
		s += successStyle.Render("  ✓ Pulled from "+m.cfg.Git.Remote+"/"+m.cfg.Git.BaseBranch) + "\n\n"
		s += fmt.Sprintf("Creating new branch...\n")
		s += fmt.Sprintf("  Branch: %s\n\n", m.branchName)
		s += successStyle.Render(fmt.Sprintf("✓ Ready to work on %s!\n", m.branchName))
		s += "\n  When done, run 'one pr' to create a PR\n"

		return s
	}

	return fmt.Sprintf("%s\n\nPress Ctrl+C to cancel", m.status)
}
