package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const gitlabBaseURL = "https://gitlab.com/api/v4"

// GitLabClient handles GitLab API operations
type GitLabClient struct {
	token      string
	httpClient *http.Client
}

// NewGitLabClient creates a new GitLab API client
func NewGitLabClient(token string) *GitLabClient {
	return &GitLabClient{
		token: token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CreateMergeRequest creates a new merge request
func (c *GitLabClient) CreateMergeRequest(projectID int, title, description, sourceBranch, targetBranch string) (string, error) {
	url := fmt.Sprintf("%s/projects/%d/merge_requests", gitlabBaseURL, projectID)

	reqBody := map[string]interface{}{
		"source_branch": sourceBranch,
		"target_branch": targetBranch,
		"title":         title,
		"description":   description,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("PRIVATE-TOKEN", c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("GitLab API error: status %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	webURL, ok := result["web_url"].(string)
	if !ok {
		return "", fmt.Errorf("unexpected response format")
	}

	return webURL, nil
}
