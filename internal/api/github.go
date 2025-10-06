package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const githubBaseURL = "https://api.github.com"

// GitHubClient handles GitHub API operations
type GitHubClient struct {
	token      string
	httpClient *http.Client
}

// NewGitHubClient creates a new GitHub API client
func NewGitHubClient(token string) *GitHubClient {
	return &GitHubClient{
		token: token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CreatePullRequest creates a new pull request
func (c *GitHubClient) CreatePullRequest(owner, repo, title, body, head, base string) (string, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/pulls", githubBaseURL, owner, repo)

	reqBody := map[string]interface{}{
		"title": title,
		"body":  body,
		"head":  head,
		"base":  base,
		"draft": false,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "one-cli/0.2.0")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("GitHub API error: status %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	htmlURL, ok := result["html_url"].(string)
	if !ok {
		return "", fmt.Errorf("unexpected response format")
	}

	return htmlURL, nil
}

// GetIssue fetches issue information
func (c *GitHubClient) GetIssue(owner, repo, issueNumber string) (string, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/issues/%s", githubBaseURL, owner, repo, issueNumber)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "one-cli/0.2.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API error: status %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	title, ok := result["title"].(string)
	if !ok {
		return "", fmt.Errorf("unexpected response format")
	}

	return title, nil
}
