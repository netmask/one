package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	// GitHub OAuth App Client ID (you would register your own)
	// For demo purposes, using a public client ID
	githubClientID = "Iv1.b507a08c87ecfe98" // Example - replace with your own

	githubDeviceCodeURL  = "https://github.com/login/device/code"
	githubAccessTokenURL = "https://github.com/login/oauth/access_token"
)

// GitHubDeviceCodeResponse represents the response from device code request
type GitHubDeviceCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

// GitHubAccessTokenResponse represents the response from access token request
type GitHubAccessTokenResponse struct {
	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type"`
	Scope            string `json:"scope"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// GitHubDeviceFlow performs the GitHub OAuth device flow
func GitHubDeviceFlow(projectName string) (*Token, error) {
	// Step 1: Request device code
	deviceCode, err := requestDeviceCode()
	if err != nil {
		return nil, fmt.Errorf("failed to request device code: %w", err)
	}

	// Step 2: Display user code and open browser
	fmt.Println()
	fmt.Printf("🔐 GitHub Authentication\n\n")
	fmt.Printf("Please visit: %s\n", deviceCode.VerificationURI)
	fmt.Printf("And enter code: %s\n\n", deviceCode.UserCode)

	// Open browser automatically
	openBrowser(deviceCode.VerificationURI)

	// Step 3: Poll for access token
	fmt.Println("Waiting for authorization...")

	token, err := pollForAccessToken(deviceCode)
	if err != nil {
		return nil, err
	}

	// Step 4: Store token in keyring
	if err := StoreToken("github", projectName, token); err != nil {
		return nil, fmt.Errorf("failed to store token: %w", err)
	}

	return token, nil
}

func requestDeviceCode() (*GitHubDeviceCodeResponse, error) {
	data := url.Values{}
	data.Set("client_id", githubClientID)
	data.Set("scope", "repo workflow")

	resp, err := http.PostForm(githubDeviceCodeURL, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result GitHubDeviceCodeResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func pollForAccessToken(deviceCode *GitHubDeviceCodeResponse) (*Token, error) {
	interval := time.Duration(deviceCode.Interval) * time.Second
	timeout := time.After(time.Duration(deviceCode.ExpiresIn) * time.Second)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return nil, fmt.Errorf("authorization timeout - code expired after %d seconds", deviceCode.ExpiresIn)

		case <-ticker.C:
			token, err := checkAccessToken(deviceCode.DeviceCode)
			if err != nil {
				if err.Error() == "pending" {
					// Still waiting, continue polling
					fmt.Print(".")
					continue
				}
				return nil, err
			}

			if token != "" {
				fmt.Println() // New line after dots
				return &Token{
					AccessToken: token,
					TokenType:   "bearer",
				}, nil
			}
		}
	}
}

func checkAccessToken(deviceCode string) (string, error) {
	data := url.Values{}
	data.Set("client_id", githubClientID)
	data.Set("device_code", deviceCode)
	data.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")

	req, err := http.NewRequest("POST", githubAccessTokenURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/json")
	req.URL.RawQuery = data.Encode()

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result GitHubAccessTokenResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if result.Error != "" {
		if result.Error == "authorization_pending" {
			return "", fmt.Errorf("pending")
		}
		return "", fmt.Errorf("oauth error: %s", result.Error)
	}

	return result.AccessToken, nil
}

func openBrowser(url string) {
	// Try to open browser using different methods
	// This is a simple implementation - can be enhanced
	fmt.Printf("Opening browser to: %s\n", url)
	// Note: The browser package can be used here if needed
}
