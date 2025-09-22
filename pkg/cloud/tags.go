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
)

// Tag represents a Contabo tag according to OpenAPI spec
type Tag struct {
	TenantID    string `json:"tenantId"`
	CustomerID  string `json:"customerId"`
	TagID       int64  `json:"tagId"`
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description"`
}

// CreateTagRequest represents a request to create a new tag
type CreateTagRequest struct {
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description,omitempty"`
}

// CreateTagResponseData represents the data part of a create tag response
type CreateTagResponseData struct {
	TenantID   string `json:"tenantId"`
	CustomerID string `json:"customerId"`
	TagID      int64  `json:"tagId"`
}

// CreateTagResponse represents the response from creating a tag
type CreateTagResponse struct {
	Data  []CreateTagResponseData `json:"data"`
	Links struct {
		Self string `json:"self"`
	} `json:"_links"`
}

// ListTagsResponse represents the response from listing tags
type ListTagsResponse struct {
	Data  []Tag `json:"data"`
	Links struct {
		Self     string `json:"self,omitempty"`
		First    string `json:"first,omitempty"`
		Previous string `json:"previous,omitempty"`
		Next     string `json:"next,omitempty"`
		Last     string `json:"last,omitempty"`
	} `json:"_links"`
}

// TagsService handles tag-related API operations
type TagsService struct {
	client *Client
}

// NewTagsService creates a new tags service
func NewTagsService(client *Client) *TagsService {
	return &TagsService{client: client}
}

// List retrieves a list of tags
func (s *TagsService) List(ctx context.Context) ([]Tag, error) {
	resp, err := s.client.doRequest(ctx, http.MethodGet, "/v1/tags", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}

	var listResp ListTagsResponse
	if err := parseResponse(resp, &listResp); err != nil {
		return nil, fmt.Errorf("failed to parse list tags response: %w", err)
	}

	return listResp.Data, nil
}

// Get retrieves a specific tag by ID
func (s *TagsService) Get(ctx context.Context, tagID int64) (*Tag, error) {
	path := fmt.Sprintf("/v1/tags/%d", tagID)
	resp, err := s.client.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get tag: %w", err)
	}

	var findResp struct {
		Data []Tag `json:"data"`
	}
	if err := parseResponse(resp, &findResp); err != nil {
		return nil, fmt.Errorf("failed to parse get tag response: %w", err)
	}

	if len(findResp.Data) == 0 {
		return nil, fmt.Errorf("tag with ID %d not found", tagID)
	}

	return &findResp.Data[0], nil
}

// Create creates a new tag
func (s *TagsService) Create(ctx context.Context, req *CreateTagRequest) (*Tag, error) {
	resp, err := s.client.doRequest(ctx, http.MethodPost, "/v1/tags", req)
	if err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}

	var createResp CreateTagResponse
	if err := parseResponse(resp, &createResp); err != nil {
		return nil, fmt.Errorf("failed to parse create tag response: %w", err)
	}

	if len(createResp.Data) == 0 {
		return nil, fmt.Errorf("no tag data returned from create request")
	}

	// Get the created tag to return full details
	tagID := createResp.Data[0].TagID
	return s.Get(ctx, tagID)
}

// Delete deletes a tag by ID
func (s *TagsService) Delete(ctx context.Context, tagID int64) error {
	path := fmt.Sprintf("/v1/tags/%d", tagID)
	_, err := s.client.doRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}

	return nil
}

// FindByName finds a tag by name
func (s *TagsService) FindByName(ctx context.Context, name string) (*Tag, error) {
	tags, err := s.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}

	for _, tag := range tags {
		if tag.Name == name {
			return &tag, nil
		}
	}

	return nil, fmt.Errorf("tag with name '%s' not found", name)
}

// EnsureTag ensures a tag exists, creating it if necessary
func (s *TagsService) EnsureTag(ctx context.Context, name, color string) (*Tag, error) {
	// Try to find existing tag first
	tag, err := s.FindByName(ctx, name)
	if err == nil {
		return tag, nil
	}

	// Tag doesn't exist, create it
	createReq := &CreateTagRequest{
		Name:  name,
		Color: color,
	}

	return s.Create(ctx, createReq)
}
