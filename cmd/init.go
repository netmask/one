package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"one/internal/auth"
	browserpkg "one/internal/browser"
	"one/internal/config"
	"one/internal/git"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new project configuration",
	Long:  `Interactive setup to create a new project configuration.`,
	RunE:  runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringP("name", "n", "", "Project name (skips prompt)")
}

func runInit(cmd *cobra.Command, args []string) error {
	projectName, _ := cmd.Flags().GetString("name")

	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Try to detect Git remote information
	var detectedRemote *git.RemoteInfo
	var detectedBaseBranch string

	remoteInfo, err := git.DetectRemote("origin")
	if err == nil {
		detectedRemote = remoteInfo
		fmt.Printf("✓ Detected Git remote: %s (%s/%s)\n\n", remoteInfo.Provider, remoteInfo.Owner, remoteInfo.Repo)
	}

	baseBranch, err := git.GetDefaultBranch()
	if err == nil {
		detectedBaseBranch = baseBranch
	} else {
		detectedBaseBranch = "main"
	}

	// Form values
	var (
		name         string = projectName
		path         string = currentDir
		provider     string
		remote       string = "origin"
		baseBranchV  string = detectedBaseBranch
		owner        string
		repo         string
		projectID    string
		workspace    string
		repoSlug     string
		tokenEnv     string
		browser      string
		profile      string
		hasTicket    bool
		ticketSystem string
		ticketURL    string
		boardID      string
		pattern      string = "^([A-Z]+-\\d+)"
		authNow      bool
	)

	// Pre-fill with detected values
	if detectedRemote != nil {
		provider = detectedRemote.Provider
		owner = detectedRemote.Owner
		repo = detectedRemote.Repo
		workspace = detectedRemote.Owner
		repoSlug = detectedRemote.Repo
	}

	// Main form
	var mainFormGroups []*huh.Group

	// Project info group
	mainFormGroups = append(mainFormGroups, huh.NewGroup(
		huh.NewInput().
			Title("Project Name").
			Value(&name).
			Validate(func(s string) error {
				if s == "" {
					return fmt.Errorf("project name is required")
				}
				return nil
			}),

		huh.NewInput().
			Title("Project Path").
			Description("Path to your project directory").
			Value(&path).
			Validate(func(s string) error {
				if s == "" {
					return fmt.Errorf("path is required")
				}
				return nil
			}),
	))

	// Git configuration group
	providerSelect := huh.NewSelect[string]().
		Title("Git Provider").
		Options(
			huh.NewOption("GitHub", "github"),
			huh.NewOption("GitLab", "gitlab"),
			huh.NewOption("Bitbucket", "bitbucket"),
		).
		Value(&provider)

	// If we detected a provider, show it in the description
	if detectedRemote != nil {
		providerSelect.Description(fmt.Sprintf("Detected: %s", detectedRemote.Provider))
	}

	mainFormGroups = append(mainFormGroups, huh.NewGroup(
		providerSelect,

		huh.NewInput().
			Title("Remote Name").
			Value(&remote).
			Placeholder("origin"),

		huh.NewInput().
			Title("Base Branch").
			Description(fmt.Sprintf("Detected: %s", detectedBaseBranch)).
			Value(&baseBranchV).
			Placeholder(detectedBaseBranch),
	))

	form := huh.NewForm(mainFormGroups...)

	if err := form.Run(); err != nil {
		return err
	}

	// Provider-specific form
	var providerForm *huh.Form

	switch provider {
	case "github":
		tokenEnv = "GITHUB_TOKEN"
		providerForm = huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("GitHub Owner").
					Description(getDetectedDesc(detectedRemote, "owner")).
					Value(&owner).
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("owner is required")
						}
						return nil
					}),

				huh.NewInput().
					Title("Repository Name").
					Description(getDetectedDesc(detectedRemote, "repo")).
					Value(&repo).
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("repository is required")
						}
						return nil
					}),

				huh.NewInput().
					Title("Token Environment Variable").
					Value(&tokenEnv).
					Placeholder("GITHUB_TOKEN"),
			),
		)

	case "gitlab":
		tokenEnv = "GITLAB_TOKEN"
		providerForm = huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("GitLab Project ID").
					Description("Numeric project ID from GitLab").
					Value(&projectID).
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("project ID is required")
						}
						if _, err := strconv.Atoi(s); err != nil {
							return fmt.Errorf("project ID must be a number")
						}
						return nil
					}),

				huh.NewInput().
					Title("Token Environment Variable").
					Value(&tokenEnv).
					Placeholder("GITLAB_TOKEN"),
			),
		)

	case "bitbucket":
		tokenEnv = "BITBUCKET_TOKEN"
		providerForm = huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Bitbucket Workspace").
					Value(&workspace).
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("workspace is required")
						}
						return nil
					}),

				huh.NewInput().
					Title("Repository Slug").
					Value(&repoSlug).
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("repository slug is required")
						}
						return nil
					}),

				huh.NewInput().
					Title("Token Environment Variable").
					Value(&tokenEnv).
					Placeholder("BITBUCKET_TOKEN"),
			),
		)
	}

	if providerForm != nil {
		if err := providerForm.Run(); err != nil {
			return err
		}
	}

	// Browser and ticket configuration
	
	// Detect available browser profiles
	var chromeProfiles []browserpkg.Profile
	var firefoxProfiles []browserpkg.Profile
	
	if cp, err := browserpkg.ListChromeProfiles(); err == nil {
		chromeProfiles = cp
	}
	if fp, err := browserpkg.ListFirefoxProfiles(); err == nil {
		firefoxProfiles = fp
	}

	configForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Browser").
				Options(
					huh.NewOption("Chrome", "chrome"),
					huh.NewOption("Firefox", "firefox"),
					huh.NewOption("Safari", "safari"),
				).
				Value(&browser),
		),
	)

	if err := configForm.Run(); err != nil {
		return err
	}

	// Show profile selection based on browser choice
	if browser == "chrome" && len(chromeProfiles) > 0 {
		var selectedProfile string
		var profileOptions []huh.Option[string]
		
		// Add "Skip (use default)" option
		profileOptions = append(profileOptions, huh.NewOption("Skip (use default profile)", ""))
		
		// Add each Chrome profile
		for _, p := range chromeProfiles {
			label := browserpkg.FormatProfile(p)
			profileOptions = append(profileOptions, huh.NewOption(label, p.Directory))
		}
		
		profileForm := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Select Chrome Profile").
					Description("Choose which Google account to use").
					Options(profileOptions...).
					Value(&selectedProfile),
			),
		)
		
		if err := profileForm.Run(); err != nil {
			return err
		}
		
		profile = selectedProfile
	} else if browser == "firefox" && len(firefoxProfiles) > 0 {
		var selectedProfile string
		var profileOptions []huh.Option[string]
		
		// Add "Skip (use default)" option
		profileOptions = append(profileOptions, huh.NewOption("Skip (use default profile)", ""))
		
		// Add each Firefox profile
		for _, p := range firefoxProfiles {
			profileOptions = append(profileOptions, huh.NewOption(p.Name, p.Directory))
		}
		
		profileForm := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Select Firefox Profile").
					Options(profileOptions...).
					Value(&selectedProfile),
			),
		)
		
		if err := profileForm.Run(); err != nil {
			return err
		}
		
		profile = selectedProfile
	}

	// Continue with ticket configuration
	ticketForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Configure Ticket System?").
				Description("Set up Jira, Linear, or GitHub Issues integration").
				Value(&hasTicket),
		),

		huh.NewGroup(
			huh.NewConfirm().
				Title("Configure Ticket System?").
				Description("Set up Jira, Linear, or GitHub Issues integration").
				Value(&hasTicket),
		),

		huh.NewGroup(
			huh.NewConfirm().
				Title("Authenticate with GitHub now?").
				Description("Use OAuth device flow to securely authenticate").
				Value(&authNow).
				Affirmative("Yes, authenticate").
				Negative("Skip for now"),
		).WithHideFunc(func() bool {
			return provider != "github"
		}),
	)

	if err := ticketForm.Run(); err != nil {
		return err
	}

	// Ticket system configuration (if requested)
	if hasTicket {
		ticketForm := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Ticket System").
					Options(
						huh.NewOption("Jira", "jira"),
						huh.NewOption("Linear", "linear"),
						huh.NewOption("GitHub Issues", "github"),
					).
					Value(&ticketSystem),

				huh.NewInput().
					Title("Base URL").
					Description("e.g., https://company.atlassian.net").
					Value(&ticketURL).
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("base URL is required")
						}
						return nil
					}),
			),
		)

		if err := ticketForm.Run(); err != nil {
			return err
		}

		// Jira-specific configuration
		if ticketSystem == "jira" {
			jiraForm := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("Jira Board ID").
						Description("e.g., PROJ").
						Value(&boardID),

					huh.NewInput().
						Title("Branch Pattern").
						Description("Regex to extract ticket ID from branch name").
						Value(&pattern).
						Placeholder("^([A-Z]+-\\d+)"),
				),
			)

			if err := jiraForm.Run(); err != nil {
				return err
			}
		}
	}

	// Build configuration
	cfg := &config.ProjectConfig{
		Version: 1,
		Project: config.ProjectInfo{
			Name:  name,
			Paths: []string{path},
		},
		Git: config.GitConfig{
			Provider:   provider,
			Remote:     remote,
			BaseBranch: baseBranchV,
		},
		Browser: config.BrowserConfig{
			Type:    browser,
			Profile: profile,
		},
	}

	// Add provider-specific config
	switch provider {
	case "github":
		cfg.Git.GitHub = &config.GitHubConfig{
			Owner:    owner,
			Repo:     repo,
			TokenEnv: tokenEnv,
		}
	case "gitlab":
		pid, _ := strconv.Atoi(projectID)
		cfg.Git.GitLab = &config.GitLabConfig{
			ProjectID: pid,
			TokenEnv:  tokenEnv,
		}
	case "bitbucket":
		cfg.Git.Bitbucket = &config.BitbucketConfig{
			Workspace: workspace,
			RepoSlug:  repoSlug,
			TokenEnv:  tokenEnv,
		}
	}

	// Add ticket configuration if requested
	if hasTicket {
		cfg.Ticket = &config.TicketConfig{
			System:  ticketSystem,
			BaseURL: ticketURL,
		}

		if ticketSystem == "jira" {
			cfg.Ticket.Jira = &config.JiraConfig{
				BoardID: boardID,
			}
		}

		if pattern != "" {
			cfg.BranchPatterns = &config.BranchPatterns{
				TicketID: pattern,
			}
		}
	}

	// Save configuration
	filename := git.SanitizeBranchName(name)
	if err := config.SaveProjectConfig(cfg, filename); err != nil {
		return err
	}

	configDir, _ := config.GetConfigDir()
	savedPath := fmt.Sprintf("%s/projects/%s.yml", configDir, filename)

	// Success message
	successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	fmt.Println()
	fmt.Println(successStyle.Render("✓ Configuration saved successfully!"))
	fmt.Println()
	fmt.Printf("  Location: %s\n", savedPath)
	fmt.Println()

	// Authenticate with GitHub if requested
	if authNow && provider == "github" {
		fmt.Println("🔐 Starting GitHub authentication...")
		fmt.Println()

		_, err := auth.GitHubDeviceFlow(name)
		if err != nil {
			fmt.Printf("⚠️  Authentication failed: %v\n", err)
			fmt.Printf("You can authenticate later by setting the %s environment variable.\n\n", tokenEnv)
		} else {
			fmt.Println()
			fmt.Println(successStyle.Render("✓ Successfully authenticated with GitHub!"))
			fmt.Println()
		}
	} else if !authNow && provider == "github" {
		fmt.Println("Authentication skipped.")
		fmt.Println("To authenticate later, set the environment variable:")
		fmt.Printf("  export %s=\"your-token-here\"\n", tokenEnv)
		fmt.Println()
	}

	fmt.Println("Next steps:")
	fmt.Println("  1. Start working on a task:")
	fmt.Println("     one start TICKET-123")
	fmt.Println()
	fmt.Println("  2. Create a pull request:")
	fmt.Println("     one pr")
	fmt.Println()

	return nil
}

// getDetectedDesc returns a description showing detected values
func getDetectedDesc(remote *git.RemoteInfo, field string) string {
	if remote == nil {
		return ""
	}

	switch field {
	case "owner":
		if remote.Owner != "" {
			return fmt.Sprintf("Detected: %s", remote.Owner)
		}
	case "repo":
		if remote.Repo != "" {
			return fmt.Sprintf("Detected: %s", remote.Repo)
		}
	}

	return ""
}
