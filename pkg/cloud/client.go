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

package cloud

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	// DefaultBaseURL is the default Contabo API base URL
	DefaultBaseURL = "https://api.contabo.com/v1"
	// DefaultTimeout is the default HTTP client timeout
	DefaultTimeout = 30 * time.Second
)

// Client represents a Contabo API client
type Client struct {
	baseURL    string
	httpClient *http.Client
	apiToken   string
}

// ClientOption is a function type for configuring the Client
type ClientOption func(*Client)

// WithBaseURL sets the base URL for the client
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithHTTPClient sets the HTTP client
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// NewClient creates a new Contabo API client
func NewClient(apiToken string, opts ...ClientOption) *Client {
	client := &Client{
		baseURL:  DefaultBaseURL,
		apiToken: apiToken,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

// APIError represents an error response from the Contabo API
type APIError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	ErrorCode  string `json:"error"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("contabo api error: status=%d, message=%s, error=%s", e.StatusCode, e.Message, e.ErrorCode)
}

// doRequest performs an HTTP request to the Contabo API
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	reqURL, err := url.JoinPath(c.baseURL, path)
	if err != nil {
		return nil, fmt.Errorf("failed to build request URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "cluster-api-provider-contabo")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	if resp.StatusCode >= 400 {
		defer func() {
			if closeErr := resp.Body.Close(); closeErr != nil {
				// Log the close error but don't override the main error
				fmt.Printf("Warning: failed to close response body: %v\n", closeErr)
			}
		}()
		var apiErr APIError
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
			// If we can't decode the error, create a generic one
			apiErr = APIError{
				StatusCode: resp.StatusCode,
				Message:    resp.Status,
				ErrorCode:  "unknown error",
			}
		}
		return nil, &apiErr
	}

	return resp, nil
}

// parseResponse parses the response body into the provided interface
func parseResponse(resp *http.Response, v interface{}) error {
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			// Log the close error but don't override the main error
			fmt.Printf("Warning: failed to close response body: %v\n", closeErr)
		}
	}()
	return json.NewDecoder(resp.Body).Decode(v)
}
