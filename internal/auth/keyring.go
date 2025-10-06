package auth

import (
	"encoding/json"
	"fmt"

	"github.com/zalando/go-keyring"
)

const serviceName = "one-cli"

// Token represents a stored authentication token
type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type"`
	ExpiresAt    *int64 `json:"expires_at,omitempty"`
}

// StoreToken stores a token in the system keyring
func StoreToken(provider, projectName string, token *Token) error {
	account := fmt.Sprintf("%s:%s", provider, projectName)

	data, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	err = keyring.Set(serviceName, account, string(data))
	if err != nil {
		return fmt.Errorf("failed to store token: %w", err)
	}

	return nil
}

// GetToken retrieves a token from the system keyring
func GetToken(provider, projectName string) (*Token, error) {
	account := fmt.Sprintf("%s:%s", provider, projectName)

	data, err := keyring.Get(serviceName, account)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	var token Token
	if err := json.Unmarshal([]byte(data), &token); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token: %w", err)
	}

	return &token, nil
}

// DeleteToken removes a token from the system keyring
func DeleteToken(provider, projectName string) error {
	account := fmt.Sprintf("%s:%s", provider, projectName)

	err := keyring.Delete(serviceName, account)
	if err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}

	return nil
}

// HasToken checks if a token exists in the system keyring
func HasToken(provider, projectName string) bool {
	_, err := GetToken(provider, projectName)
	return err == nil
}
