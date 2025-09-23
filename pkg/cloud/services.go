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
	"context"
	"fmt"
	"net/http"
	"time"
)

// Image represents a Contabo OS image
type Image struct {
	ImageID        string  `json:"imageId"`
	TenantID       string  `json:"tenantId"`
	CustomerID     string  `json:"customerId"`
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	URL            string  `json:"url"`
	SizeMb         float64 `json:"sizeMb"`
	UploadedSizeMb float64 `json:"uploadedSizeMb"`
	OSType         string  `json:"osType"`
	Version        string  `json:"version"`
	Format         string  `json:"format"`
	Status         string  `json:"status"`
	ErrorMessage   string  `json:"errorMessage"`
	StandardImage  bool    `json:"standardImage"`
}

// ListImagesResponse represents the response from listing images
type ListImagesResponse struct {
	Data  []Image `json:"data"`
	Links struct {
		Self     string `json:"self,omitempty"`
		First    string `json:"first,omitempty"`
		Previous string `json:"previous,omitempty"`
		Next     string `json:"next,omitempty"`
		Last     string `json:"last,omitempty"`
	} `json:"_links"`
}

// SSHKey represents a Contabo secret/SSH key
type SSHKey struct {
	TenantID   string    `json:"tenantId"`
	CustomerID string    `json:"customerId"`
	SecretID   int64     `json:"secretId"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	Value      string    `json:"value"`
	CreatedAt  time.Time `json:"createdAt"`
}

// ListSSHKeysResponse represents the response from listing SSH keys/secrets
type ListSSHKeysResponse struct {
	Data  []SSHKey `json:"data"`
	Links struct {
		Self     string `json:"self,omitempty"`
		First    string `json:"first,omitempty"`
		Previous string `json:"previous,omitempty"`
		Next     string `json:"next,omitempty"`
		Last     string `json:"last,omitempty"`
	} `json:"_links"`
}

// ImagesService handles image-related API operations
type ImagesService struct {
	client *Client
}

// NewImagesService creates a new images service
func NewImagesService(client *Client) *ImagesService {
	return &ImagesService{client: client}
}

// List retrieves a list of available images
func (s *ImagesService) List(ctx context.Context) ([]Image, error) {
	resp, err := s.client.doRequest(ctx, http.MethodGet, "/v1/compute/images", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %w", err)
	}

	var listResp ListImagesResponse
	if err := parseResponse(resp, &listResp); err != nil {
		return nil, fmt.Errorf("failed to parse list images response: %w", err)
	}

	return listResp.Data, nil
}

// SSHKeysService handles SSH key-related API operations via secrets endpoint
type SSHKeysService struct {
	client *Client
}

// NewSSHKeysService creates a new SSH keys service
func NewSSHKeysService(client *Client) *SSHKeysService {
	return &SSHKeysService{client: client}
}

// List retrieves a list of SSH keys (using secrets endpoint)
func (s *SSHKeysService) List(ctx context.Context) ([]SSHKey, error) {
	resp, err := s.client.doRequest(ctx, http.MethodGet, "/v1/secrets", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list SSH keys: %w", err)
	}

	var listResp ListSSHKeysResponse
	if err := parseResponse(resp, &listResp); err != nil {
		return nil, fmt.Errorf("failed to parse list SSH keys response: %w", err)
	}

	// Filter only SSH key type secrets
	var sshKeys []SSHKey
	for _, secret := range listResp.Data {
		if secret.Type == "ssh" {
			sshKeys = append(sshKeys, secret)
		}
	}

	return sshKeys, nil
}

// ContaboService is the main service that provides access to all Contabo API endpoints
type ContaboService struct {
	client    *Client
	Instances *InstancesService
	Images    *ImagesService
	SSHKeys   *SSHKeysService
	Tags      *TagsService
}

// NewContaboService creates a new Contabo service with all sub-services
func NewContaboService(apiToken string, opts ...ClientOption) *ContaboService {
	client := NewClient(apiToken, opts...)

	return &ContaboService{
		client:    client,
		Instances: NewInstancesService(client),
		Images:    NewImagesService(client),
		SSHKeys:   NewSSHKeysService(client),
		Tags:      NewTagsService(client),
	}
}
