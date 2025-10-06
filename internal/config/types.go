package config

// ProjectConfig represents the complete project configuration
type ProjectConfig struct {
	Version        int              `yaml:"version"`
	Project        ProjectInfo      `yaml:"project"`
	Git            GitConfig        `yaml:"git"`
	Browser        BrowserConfig    `yaml:"browser"`
	Ticket         *TicketConfig    `yaml:"ticket,omitempty"`
	Templates      *Templates       `yaml:"templates,omitempty"`
	BranchPatterns *BranchPatterns  `yaml:"branch_patterns,omitempty"`
}

// ProjectInfo contains basic project information
type ProjectInfo struct {
	Name  string   `yaml:"name"`
	Paths []string `yaml:"paths"`
}

// GitConfig contains git-related configuration
type GitConfig struct {
	Provider   string          `yaml:"provider"`
	Remote     string          `yaml:"remote"`
	BaseBranch string          `yaml:"base_branch"`
	GitHub     *GitHubConfig   `yaml:"github,omitempty"`
	GitLab     *GitLabConfig   `yaml:"gitlab,omitempty"`
	Bitbucket  *BitbucketConfig `yaml:"bitbucket,omitempty"`
}

// GitHubConfig contains GitHub-specific settings
type GitHubConfig struct {
	Owner    string `yaml:"owner"`
	Repo     string `yaml:"repo"`
	TokenEnv string `yaml:"token_env"`
}

// GitLabConfig contains GitLab-specific settings
type GitLabConfig struct {
	ProjectID int    `yaml:"project_id"`
	TokenEnv  string `yaml:"token_env"`
}

// BitbucketConfig contains Bitbucket-specific settings
type BitbucketConfig struct {
	Workspace string `yaml:"workspace"`
	RepoSlug  string `yaml:"repo_slug"`
	TokenEnv  string `yaml:"token_env"`
}

// BrowserConfig contains browser-related settings
type BrowserConfig struct {
	Type    string `yaml:"type"`
	Profile string `yaml:"profile,omitempty"`
}

// TicketConfig contains ticket system configuration
type TicketConfig struct {
	System  string       `yaml:"system"`
	BaseURL string       `yaml:"base_url"`
	Jira    *JiraConfig  `yaml:"jira,omitempty"`
}

// JiraConfig contains Jira-specific settings
type JiraConfig struct {
	BoardID  string `yaml:"board_id"`
	TokenEnv string `yaml:"token_env,omitempty"`
}

// Templates contains PR/MR template strings
type Templates struct {
	PRTitle string `yaml:"pr_title"`
	PRBody  string `yaml:"pr_body"`
}

// BranchPatterns contains regex patterns for branch parsing
type BranchPatterns struct {
	TicketID string `yaml:"ticket_id"`
}

// GlobalConfig represents the global configuration
type GlobalConfig struct {
	Version  int      `yaml:"version"`
	Defaults Defaults `yaml:"defaults,omitempty"`
}

// Defaults contains default settings
type Defaults struct {
	Browser *BrowserConfig `yaml:"browser,omitempty"`
	Git     *GitDefaults   `yaml:"git,omitempty"`
}

// GitDefaults contains default git settings
type GitDefaults struct {
	Remote     string `yaml:"remote,omitempty"`
	BaseBranch string `yaml:"base_branch,omitempty"`
}
