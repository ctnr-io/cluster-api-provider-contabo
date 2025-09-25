/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/log"
)

var authLog = log.Log.WithName("contabo-auth")

// OAuth2TokenResponse represents the response from Contabo's OAuth2 token endpoint
type OAuth2TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// TokenManager manages OAuth2 token lifecycle with automatic refresh
type TokenManager struct {
	mu           sync.RWMutex
	clientID     string
	clientSecret string
	apiUser      string
	apiPassword  string
	accessToken  string
	expiresAt    time.Time
	tokenURL     string
}

// NewTokenManager creates a new token manager for Contabo OAuth2 authentication
func NewTokenManager(clientID, clientSecret, apiUser, apiPassword string) *TokenManager {
	return &TokenManager{
		clientID:     clientID,
		clientSecret: clientSecret,
		apiUser:      apiUser,
		apiPassword:  apiPassword,
		tokenURL:     "https://auth.contabo.com/auth/realms/contabo/protocol/openid-connect/token",
	}
}

// GetToken returns a valid access token, refreshing if necessary
func (tm *TokenManager) GetToken() (string, error) {
	tm.mu.RLock()
	// Check if token is still valid
	if tm.accessToken != "" && time.Now().Before(tm.expiresAt) {
		token := tm.accessToken
		tm.mu.RUnlock()
		return token, nil
	}
	tm.mu.RUnlock()

	// Need to refresh token
	return tm.refreshToken()
}

// refreshToken obtains a new access token
func (tm *TokenManager) refreshToken() (string, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Double-check in case another goroutine already refreshed
	if tm.accessToken != "" && time.Now().Before(tm.expiresAt) {
		return tm.accessToken, nil
	}

	authLog.Info("Refreshing OAuth2 access token")

	data := url.Values{}
	data.Set("client_id", tm.clientID)
	data.Set("client_secret", tm.clientSecret)
	data.Set("username", tm.apiUser)
	data.Set("password", tm.apiPassword)
	data.Set("grant_type", "password")

	// Create HTTP request with proper headers
	req, err := http.NewRequest("POST", tm.tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create OAuth2 request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to request OAuth2 token: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			authLog.Error(closeErr, "failed to close response body")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		// Read all the response body
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("OAuth2 token request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp OAuth2TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode OAuth2 token response: %w", err)
	}

	// Update token and expiration time
	tm.accessToken = tokenResp.AccessToken
	tm.expiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	authLog.Info("Successfully refreshed OAuth2 access token", "expiresAt", tm.expiresAt)

	return tm.accessToken, nil
}

// IsTokenValid returns true if the current token is valid (with 5 minute buffer)
func (tm *TokenManager) IsTokenValid() bool {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.accessToken != "" && time.Now().Before(tm.expiresAt)
}

// GetExpirationTime returns when the current token expires
func (tm *TokenManager) GetExpirationTime() time.Time {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.expiresAt
}
