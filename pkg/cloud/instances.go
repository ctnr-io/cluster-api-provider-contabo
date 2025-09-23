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
	"strconv"
	"strings"
	"time"
)

// DisplayName state management constants
const (
	// CAPCPrefix is the prefix for all CAPC-managed instances
	CAPCPrefix = "capc-"
	// AvailableState indicates an instance is available for use
	AvailableState = CAPCPrefix + "available"
	// InUseState indicates an instance is in use but not bound to a specific cluster
	InUseState = CAPCPrefix + "in-use"
	// ClusterPrefix is used to bind an instance to a specific cluster
	ClusterPrefix = CAPCPrefix + "cluster-"
	// MaxDisplayNameLength is the maximum length for displayName (based on OpenAPI spec)
	MaxDisplayNameLength = 255
)

// InstanceState represents the possible states of a CAPC-managed instance
type InstanceState int

const (
	StateUnknown InstanceState = iota
	StateAvailable
	StateInUse
	StateClusterBound
)

// GetInstanceState extracts the state from an instance's displayName
func GetInstanceState(displayName string) InstanceState {
	if displayName == AvailableState {
		return StateAvailable
	}
	if displayName == InUseState {
		return StateInUse
	}
	if strings.HasPrefix(displayName, ClusterPrefix) {
		return StateClusterBound
	}
	return StateUnknown
}

// GetClusterName extracts the cluster name from a cluster-bound instance's displayName
func GetClusterName(displayName string) string {
	if strings.HasPrefix(displayName, ClusterPrefix) {
		return strings.TrimPrefix(displayName, ClusterPrefix)
	}
	return ""
}

// CreateClusterDisplayName creates a displayName for a cluster-bound instance
func CreateClusterDisplayName(clusterName string) string {
	displayName := ClusterPrefix + clusterName
	if len(displayName) > MaxDisplayNameLength {
		// Truncate cluster name to fit within the limit
		maxClusterNameLength := MaxDisplayNameLength - len(ClusterPrefix)
		displayName = ClusterPrefix + clusterName[:maxClusterNameLength]
	}
	return displayName
}

// IsManagedByCAPC checks if an instance is managed by CAPC based on its displayName
func IsManagedByCAPC(displayName string) bool {
	return strings.HasPrefix(displayName, CAPCPrefix)
}

// Instance represents a Contabo VPS instance
type Instance struct {
	TenantID    string            `json:"tenantId"`
	CustomerID  string            `json:"customerId"`
	InstanceID  int64             `json:"instanceId"`
	Name        string            `json:"name"`
	DisplayName string            `json:"displayName"`
	Status      string            `json:"status"`
	Region      string            `json:"region"`
	RegionName  string            `json:"regionName"`
	DataCenter  string            `json:"dataCenter"`
	ProductID   string            `json:"productId"`
	ImageID     string            `json:"imageId"`
	IPConfig    InstanceIPConfig  `json:"ipConfig"`
	MacAddress  string            `json:"macAddress"`
	RAMMb       float64           `json:"ramMb"`
	CPUCores    int               `json:"cpuCores"`
	OSType      string            `json:"osType"`
	DiskMb      float64           `json:"diskMb"`
	SSHKeys     []int64           `json:"sshKeys"`
	CreatedDate time.Time         `json:"createdDate"`
	CancelDate  *time.Time        `json:"cancelDate,omitempty"`
	Image       InstanceImage     `json:"image"`
	Product     InstanceProduct   `json:"product"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// InstanceIPConfig represents the IP configuration of an instance
type InstanceIPConfig struct {
	V4 InstanceIPv4Config `json:"v4"`
	V6 InstanceIPv6Config `json:"v6"`
}

// InstanceIPv4Config represents IPv4 configuration
type InstanceIPv4Config struct {
	IP      string `json:"ip"`
	Gateway string `json:"gateway"`
	Netmask string `json:"netmask"`
}

// InstanceIPv6Config represents IPv6 configuration
type InstanceIPv6Config struct {
	IP      string `json:"ip"`
	Gateway string `json:"gateway"`
	Netmask string `json:"netmask"`
}

// InstanceImage represents the OS image of an instance
type InstanceImage struct {
	ImageID     string `json:"imageId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	OSType      string `json:"osType"`
}

// InstanceProduct represents the product/plan of an instance
type InstanceProduct struct {
	ProductID   string `json:"productId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CPUCores    int    `json:"cpuCores"`
	RAM         int    `json:"ram"`
	Disk        int    `json:"disk"`
}

// CreateInstanceRequest represents a request to create a new instance
type CreateInstanceRequest struct {
	ImageID     string            `json:"imageId"`
	ProductID   string            `json:"productId"`
	Region      string            `json:"region"`
	DisplayName string            `json:"displayName"`
	UserData    string            `json:"userData,omitempty"`
	SSHKeys     []int64           `json:"sshKeys,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// CreateInstanceResponse represents the response from creating an instance
type CreateInstanceResponse struct {
	InstanceID  int64  `json:"instanceId"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Status      string `json:"status"`
}

// PatchInstanceRequest represents a request to update an instance
type PatchInstanceRequest struct {
	DisplayName string `json:"displayName"`
}

// ListInstancesResponse represents the response from listing instances
type ListInstancesResponse struct {
	Data []Instance `json:"data"`
	Meta struct {
		Pagination struct {
			CurrentPage int `json:"currentPage"`
			From        int `json:"from"`
			LastPage    int `json:"lastPage"`
			PerPage     int `json:"perPage"`
			To          int `json:"to"`
			Total       int `json:"total"`
		} `json:"pagination"`
	} `json:"_links"`
	Pagination struct {
		Size          int `json:"size"`
		TotalElements int `json:"totalElements"`
		TotalPages    int `json:"totalPages"`
		Page          int `json:"page"`
	} `json:"_pagination"`
}

// ListInstancesOptions contains options for listing instances based on Contabo OpenAPI spec
type ListInstancesOptions struct {
	// Pagination
	Page    *int     `json:"page,omitempty"`    // Number of page to be fetched
	Size    *int     `json:"size,omitempty"`    // Number of elements per page
	OrderBy []string `json:"orderBy,omitempty"` // Specify fields and ordering (field:ASC|DESC)

	// Filters
	Name         string `json:"name,omitempty"`         // The name of the instance
	DisplayName  string `json:"displayName,omitempty"`  // The display name of the instance
	DataCenter   string `json:"dataCenter,omitempty"`   // The data center of the instance
	Region       string `json:"region,omitempty"`       // The Region of the instance
	InstanceIds  string `json:"instanceIds,omitempty"`  // Comma separated instances identifiers
	Status       string `json:"status,omitempty"`       // The status of the instance
	ProductIds   string `json:"productIds,omitempty"`   // Identifiers of the instance products
	AddOnIds     string `json:"addOnIds,omitempty"`     // Identifiers of Addons the instances have
	ProductTypes string `json:"productTypes,omitempty"` // Comma separated instance's category
	IpConfig     *bool  `json:"ipConfig,omitempty"`     // Filter instances that have an ip config
	Search       string `json:"search,omitempty"`       // Full text search by name, displayName, ipAddress
}

// InstancesService handles instance-related API operations
type InstancesService struct {
	client *Client
}

// NewInstancesService creates a new instances service
func NewInstancesService(client *Client) *InstancesService {
	return &InstancesService{client: client}
}

// List retrieves instances based on the provided options
// Uses the exact query parameters supported by the Contabo OpenAPI spec
func (s *InstancesService) List(ctx context.Context, opts *ListInstancesOptions) ([]Instance, error) {
	path := "/v1/compute/instances"

	// Build query parameters based on OpenAPI spec
	queryParams := make([]string, 0)

	if opts != nil {
		// Pagination parameters
		if opts.Page != nil {
			queryParams = append(queryParams, fmt.Sprintf("page=%d", *opts.Page))
		}
		if opts.Size != nil {
			queryParams = append(queryParams, fmt.Sprintf("size=%d", *opts.Size))
		}

		// OrderBy parameter (array of strings)
		for _, order := range opts.OrderBy {
			queryParams = append(queryParams, fmt.Sprintf("orderBy=%s", order))
		}

		// Filter parameters
		if opts.Name != "" {
			queryParams = append(queryParams, fmt.Sprintf("name=%s", opts.Name))
		}
		if opts.DisplayName != "" {
			queryParams = append(queryParams, fmt.Sprintf("displayName=%s", opts.DisplayName))
		}
		if opts.DataCenter != "" {
			queryParams = append(queryParams, fmt.Sprintf("dataCenter=%s", opts.DataCenter))
		}
		if opts.Region != "" {
			queryParams = append(queryParams, fmt.Sprintf("region=%s", opts.Region))
		}
		if opts.InstanceIds != "" {
			queryParams = append(queryParams, fmt.Sprintf("instanceIds=%s", opts.InstanceIds))
		}
		if opts.Status != "" {
			queryParams = append(queryParams, fmt.Sprintf("status=%s", opts.Status))
		}
		if opts.ProductIds != "" {
			queryParams = append(queryParams, fmt.Sprintf("productIds=%s", opts.ProductIds))
		}
		if opts.AddOnIds != "" {
			queryParams = append(queryParams, fmt.Sprintf("addOnIds=%s", opts.AddOnIds))
		}
		if opts.ProductTypes != "" {
			queryParams = append(queryParams, fmt.Sprintf("productTypes=%s", opts.ProductTypes))
		}
		if opts.IpConfig != nil {
			queryParams = append(queryParams, fmt.Sprintf("ipConfig=%t", *opts.IpConfig))
		}
		if opts.Search != "" {
			queryParams = append(queryParams, fmt.Sprintf("search=%s", opts.Search))
		}
	}

	// Add query parameters to path
	if len(queryParams) > 0 {
		path += "?" + strings.Join(queryParams, "&")
	}

	resp, err := s.client.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list instances: %w", err)
	}

	var listResp ListInstancesResponse
	if err := parseResponse(resp, &listResp); err != nil {
		return nil, fmt.Errorf("failed to parse list instances response: %w", err)
	}

	return listResp.Data, nil
}

// ListAll retrieves all instances by automatically paginating through all pages
// This is a convenience method for when you need all instances without manual pagination
func (s *InstancesService) ListAll(ctx context.Context, filters *ListInstancesOptions) ([]Instance, error) {
	var allInstances []Instance
	page := 1
	size := 100      // Use larger page size for efficiency
	maxPages := 1000 // Safety limit to prevent infinite loops

	for {
		// Create options for this page
		opts := &ListInstancesOptions{
			Page: &page,
			Size: &size,
		}

		// Copy filters if provided
		if filters != nil {
			opts.Name = filters.Name
			opts.DisplayName = filters.DisplayName
			opts.DataCenter = filters.DataCenter
			opts.Region = filters.Region
			opts.InstanceIds = filters.InstanceIds
			opts.Status = filters.Status
			opts.ProductIds = filters.ProductIds
			opts.AddOnIds = filters.AddOnIds
			opts.ProductTypes = filters.ProductTypes
			opts.IpConfig = filters.IpConfig
			opts.Search = filters.Search
			opts.OrderBy = filters.OrderBy
		}

		// Get this page
		instances, err := s.List(ctx, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to list instances page %d: %w", page, err)
		}

		allInstances = append(allInstances, instances...)

		// Check if we've reached the last page (no more data)
		if len(instances) == 0 || len(instances) < size {
			break
		}

		// Safety check to prevent infinite loops
		if page >= maxPages {
			return nil, fmt.Errorf("reached maximum page limit (%d) while fetching instances", maxPages)
		}

		page++
	}

	return allInstances, nil
}

// Get retrieves a specific instance by ID
func (s *InstancesService) Get(ctx context.Context, instanceID int64) (*Instance, error) {
	path := fmt.Sprintf("/v1/compute/instances/%d", instanceID)
	resp, err := s.client.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get instance: %w", err)
	}

	var instance Instance
	if err := parseResponse(resp, &instance); err != nil {
		return nil, fmt.Errorf("failed to parse get instance response: %w", err)
	}

	return &instance, nil
}

// Create creates a new instance
func (s *InstancesService) Create(ctx context.Context, req *CreateInstanceRequest) (*CreateInstanceResponse, error) {
	resp, err := s.client.doRequest(ctx, http.MethodPost, "/v1/compute/instances", req)
	if err != nil {
		return nil, fmt.Errorf("failed to create instance: %w", err)
	}

	var createResp CreateInstanceResponse
	if err := parseResponse(resp, &createResp); err != nil {
		return nil, fmt.Errorf("failed to parse create instance response: %w", err)
	}

	return &createResp, nil
}

// CancelInstanceRequest represents a request to cancel an instance
type CancelInstanceRequest struct {
	// Add any required fields for cancel request per OpenAPI spec
}

// ReinstallInstanceRequest represents a request to reinstall an instance with a new image
type ReinstallInstanceRequest struct {
	ImageID  string `json:"imageId"`
	UserData string `json:"userData,omitempty"`
}

// Cancel cancels an instance by ID (according to Contabo API terminology)
func (s *InstancesService) Cancel(ctx context.Context, instanceID int64) error {
	path := fmt.Sprintf("/v1/compute/instances/%d/cancel", instanceID)
	req := &CancelInstanceRequest{}
	_, err := s.client.doRequest(ctx, http.MethodPost, path, req)
	if err != nil {
		return fmt.Errorf("failed to cancel instance: %w", err)
	}

	return nil
}

// Reinstall reinstalls an instance with a new image and user data
func (s *InstancesService) Reinstall(ctx context.Context, instanceID int64, req *ReinstallInstanceRequest) error {
	path := fmt.Sprintf("/v1/compute/instances/%d/reinstall", instanceID)
	_, err := s.client.doRequest(ctx, http.MethodPost, path, req)
	if err != nil {
		return fmt.Errorf("failed to reinstall instance: %w", err)
	}

	return nil
}

// Delete is an alias for Cancel to maintain compatibility
func (s *InstancesService) Delete(ctx context.Context, instanceID int64) error {
	return s.Cancel(ctx, instanceID)
}

// Start starts an instance
func (s *InstancesService) Start(ctx context.Context, instanceID int64) error {
	path := fmt.Sprintf("/v1/compute/instances/%d/actions/start", instanceID)
	_, err := s.client.doRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return fmt.Errorf("failed to start instance: %w", err)
	}

	return nil
}

// Stop stops an instance
func (s *InstancesService) Stop(ctx context.Context, instanceID int64) error {
	path := fmt.Sprintf("/v1/compute/instances/%d/actions/stop", instanceID)
	_, err := s.client.doRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return fmt.Errorf("failed to stop instance: %w", err)
	}

	return nil
}

// Restart restarts an instance
func (s *InstancesService) Restart(ctx context.Context, instanceID int64) error {
	path := fmt.Sprintf("/v1/compute/instances/%d/actions/restart", instanceID)
	_, err := s.client.doRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return fmt.Errorf("failed to restart instance: %w", err)
	}

	return nil
}

// ParseProviderID extracts the instance ID from a provider ID string
func ParseProviderID(providerID string) (int64, error) {
	if providerID == "" {
		return 0, fmt.Errorf("provider ID is empty")
	}

	// Provider ID format: contabo://instanceId
	const prefix = "contabo://"
	if len(providerID) <= len(prefix) {
		return 0, fmt.Errorf("invalid provider ID format: %s", providerID)
	}

	instanceIDStr := providerID[len(prefix):]
	instanceID, err := strconv.ParseInt(instanceIDStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse instance ID from provider ID %s: %w", providerID, err)
	}

	return instanceID, nil
}

// UpdateDisplayName updates the displayName of an instance using the PATCH endpoint
func (s *InstancesService) UpdateDisplayName(ctx context.Context, instanceID int64, displayName string) error {
	path := fmt.Sprintf("/v1/compute/instances/%d", instanceID)

	req := &PatchInstanceRequest{
		DisplayName: displayName,
	}

	_, err := s.client.doRequest(ctx, http.MethodPatch, path, req)
	if err != nil {
		return fmt.Errorf("failed to update displayName for instance %d: %w", instanceID, err)
	}

	return nil
}

// SetInstanceState updates an instance's displayName to reflect its state
func (s *InstancesService) SetInstanceState(
	ctx context.Context,
	instanceID int64,
	state InstanceState,
	clusterName string,
) error {
	var displayName string

	switch state {
	case StateAvailable:
		displayName = AvailableState
	case StateInUse:
		displayName = InUseState
	case StateClusterBound:
		if clusterName == "" {
			return fmt.Errorf("cluster name is required for StateClusterBound")
		}
		displayName = CreateClusterDisplayName(clusterName)
	default:
		return fmt.Errorf("invalid state: %v", state)
	}

	return s.UpdateDisplayName(ctx, instanceID, displayName)
}

// BuildProviderID builds a provider ID string from an instance ID
func BuildProviderID(instanceID int64) string {
	return fmt.Sprintf("contabo://%d", instanceID)
}
