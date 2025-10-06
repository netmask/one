package git

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
)

// Repository wraps go-git repository operations
type Repository struct {
	repo *git.Repository
}

// OpenRepository discovers and opens the git repository
func OpenRepository() (*Repository, error) {
	repo, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return nil, fmt.Errorf("not a git repository: %w", err)
	}

	return &Repository{repo: repo}, nil
}

// CurrentBranch returns the current branch name
func (r *Repository) CurrentBranch() (string, error) {
	ref, err := r.repo.Head()
	if err != nil {
		return "", fmt.Errorf("failed to get HEAD: %w", err)
	}

	if !ref.Name().IsBranch() {
		return "", fmt.Errorf("HEAD is detached")
	}

	return ref.Name().Short(), nil
}

// IsClean checks if the working directory is clean
func (r *Repository) IsClean() (bool, error) {
	worktree, err := r.repo.Worktree()
	if err != nil {
		return false, fmt.Errorf("failed to get worktree: %w", err)
	}

	status, err := worktree.Status()
	if err != nil {
		return false, fmt.Errorf("failed to get status: %w", err)
	}

	return status.IsClean(), nil
}

// CheckoutBranch checks out an existing branch
func (r *Repository) CheckoutBranch(branchName string) error {
	worktree, err := r.repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branchName),
	})
	if err != nil {
		return fmt.Errorf("failed to checkout branch: %w", err)
	}

	return nil
}

// CreateBranch creates a new branch from the current HEAD
func (r *Repository) CreateBranch(branchName string) error {
	head, err := r.repo.Head()
	if err != nil {
		return fmt.Errorf("failed to get HEAD: %w", err)
	}

	refName := plumbing.NewBranchReferenceName(branchName)
	ref := plumbing.NewHashReference(refName, head.Hash())

	err = r.repo.Storer.SetReference(ref)
	if err != nil {
		return fmt.Errorf("failed to create branch: %w", err)
	}

	return nil
}

// Pull performs a fast-forward pull from the remote
func (r *Repository) Pull(remoteName, branchName string) error {
	worktree, err := r.repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	err = worktree.Pull(&git.PullOptions{
		RemoteName:    remoteName,
		ReferenceName: plumbing.NewBranchReferenceName(branchName),
		SingleBranch:  true,
		Force:         false,
	})

	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("failed to pull: %w", err)
	}

	return nil
}

// Push pushes the current branch to the remote
func (r *Repository) Push(remoteName, branchName string) error {
	err := r.repo.Push(&git.PushOptions{
		RemoteName: remoteName,
		RefSpecs: []config.RefSpec{
			config.RefSpec(fmt.Sprintf("refs/heads/%s:refs/heads/%s", branchName, branchName)),
		},
	})

	if err != nil {
		return fmt.Errorf("failed to push: %w", err)
	}

	return nil
}

// ParseTicketID extracts the ticket ID from a branch name using a regex pattern
func ParseTicketID(branchName, pattern string) (string, error) {
	if pattern == "" {
		return "", fmt.Errorf("no pattern provided")
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", fmt.Errorf("invalid regex pattern: %w", err)
	}

	matches := re.FindStringSubmatch(branchName)
	if len(matches) < 2 {
		return "", fmt.Errorf("no ticket ID found in branch name")
	}

	return matches[1], nil
}

// SanitizeBranchName cleans a string to be safe for use as a branch name
func SanitizeBranchName(text string) string {
	// Convert to lowercase
	text = strings.ToLower(text)

	// Replace whitespace with hyphens
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, "-")

	// Remove special characters (keep alphanumeric, hyphens, underscores)
	text = regexp.MustCompile(`[^a-z0-9-_]`).ReplaceAllString(text, "")

	// Remove consecutive hyphens
	text = regexp.MustCompile(`-+`).ReplaceAllString(text, "-")

	// Trim hyphens from ends
	text = strings.Trim(text, "-")

	// Truncate if too long
	maxLength := 50
	if len(text) > maxLength {
		text = text[:maxLength]
		text = strings.TrimRight(text, "-")
	}

	return text
}
