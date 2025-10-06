package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// JiraClient handles Jira API operations
type JiraClient struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// NewJiraClient creates a new Jira API client
func NewJiraClient(baseURL, token string) *JiraClient {
	return &JiraClient{
		baseURL: baseURL,
		token:   token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetIssue fetches issue information
func (c *JiraClient) GetIssue(issueKey string) (string, error) {
	url := fmt.Sprintf("%s/rest/api/2/issue/%s", c.baseURL, issueKey)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Basic "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Jira API error: status %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	fields, ok := result["fields"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("unexpected response format")
	}

	summary, ok := fields["summary"].(string)
	if !ok {
		return "", fmt.Errorf("unexpected response format")
	}

	return summary, nil
}
