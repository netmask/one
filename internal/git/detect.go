package git

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
)

// RemoteInfo contains parsed information about a Git remote
type RemoteInfo struct {
	Provider  string // github, gitlab, bitbucket
	Owner     string // organization or username
	Repo      string // repository name
	URL       string // full remote URL
	ProjectID int    // for GitLab (0 if not applicable)
}

// DetectRemote attempts to detect and parse the Git remote information
func DetectRemote(remoteName string) (*RemoteInfo, error) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return nil, fmt.Errorf("not a git repository: %w", err)
	}

	remote, err := repo.Remote(remoteName)
	if err != nil {
		return nil, fmt.Errorf("remote '%s' not found: %w", remoteName, err)
	}

	if len(remote.Config().URLs) == 0 {
		return nil, fmt.Errorf("no URLs configured for remote '%s'", remoteName)
	}

	remoteURL := remote.Config().URLs[0]
	return ParseRemoteURL(remoteURL)
}

// ParseRemoteURL parses a Git remote URL and extracts provider information
func ParseRemoteURL(remoteURL string) (*RemoteInfo, error) {
	info := &RemoteInfo{
		URL: remoteURL,
	}

	// Handle SSH URLs (git@github.com:owner/repo.git)
	sshPattern := regexp.MustCompile(`^(?:ssh://)?git@([^:]+):(.+?)(?:\.git)?$`)
	if matches := sshPattern.FindStringSubmatch(remoteURL); matches != nil {
		host := matches[1]
		path := matches[2]

		info.Provider = detectProvider(host)
		parts := strings.Split(path, "/")
		if len(parts) >= 2 {
			info.Owner = parts[0]
			info.Repo = parts[1]
		}
		return info, nil
	}

	// Handle HTTPS URLs
	parsedURL, err := url.Parse(remoteURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse remote URL: %w", err)
	}

	info.Provider = detectProvider(parsedURL.Host)

	// Extract owner and repo from path
	path := strings.TrimPrefix(parsedURL.Path, "/")
	path = strings.TrimSuffix(path, ".git")
	parts := strings.Split(path, "/")

	if len(parts) >= 2 {
		info.Owner = parts[0]
		info.Repo = parts[1]
	}

	return info, nil
}

// detectProvider determines the Git provider from the hostname
func detectProvider(host string) string {
	host = strings.ToLower(host)

	if strings.Contains(host, "github") {
		return "github"
	}
	if strings.Contains(host, "gitlab") {
		return "gitlab"
	}
	if strings.Contains(host, "bitbucket") {
		return "bitbucket"
	}

	return "unknown"
}

// GetDefaultBranch attempts to detect the default branch
func GetDefaultBranch() (string, error) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return "main", nil // fallback
	}

	// Try to get HEAD
	head, err := repo.Head()
	if err != nil {
		return "main", nil // fallback
	}

	branchName := head.Name().Short()

	// Common default branches
	if branchName == "main" || branchName == "master" || branchName == "develop" {
		return branchName, nil
	}

	// Default to main
	return "main", nil
}
