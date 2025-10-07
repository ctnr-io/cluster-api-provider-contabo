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

package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"dario.cat/mergo"
	infrastructurev1beta2 "github.com/ctnr-io/cluster-api-provider-contabo/api/v1beta2"
	"github.com/ctnr-io/cluster-api-provider-contabo/pkg/contabo/v1.0.0/models"
	"github.com/google/uuid"
	"go.yaml.in/yaml/v2"
)

const (
	// ProviderIDPrefix is the prefix for Contabo provider IDs
	ProviderIDPrefix = "contabo://"

	// StateAvailable indicates an instance is available for cluster assignment
	StateAvailable = "capc-available"

	// StateProvisioning indicates an instance is being provisioned for a cluster
	StateProvisioning = "capc-provisioning"

	// StateClusterBound indicates an instance is successfully bound to a cluster
	StateClusterBound = "capc-cluster-bound"

	// StateError indicates an instance failed provisioning and has an error
	StateError = "capc-error"

	// ClusterUUIDLabel is the label key used to store the unique cluster UUID
	ClusterUUIDLabel = "cluster.x-k8s.io/capc-uuid"

	// DefaultUbuntuImageID is the standardized Ubuntu image used for all cluster nodes
	// Using a fixed image ensures consistency, security, and compatibility across the cluster
	// This should be Ubuntu 24.04 LTS - you may need to adjust this ID based on available images in Contabo
	DefaultUbuntuImageID = "d64d5c6c-9dda-4e38-8174-0ee282474d8a"
)

// BuildProviderID constructs a provider ID from an instance ID
func BuildProviderID(instanceName string) string {
	return fmt.Sprintf("%s%s", ProviderIDPrefix, instanceName)
}

// ParseProviderID extracts the instance ID from a provider ID
func ParseProviderID(providerID string) (int64, error) {
	if !strings.HasPrefix(providerID, ProviderIDPrefix) {
		return 0, fmt.Errorf("invalid provider ID format: %s", providerID)
	}

	instanceIDStr := strings.TrimPrefix(providerID, ProviderIDPrefix)
	instanceID, err := strconv.ParseInt(instanceIDStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse instance ID from provider ID %s: %w", providerID, err)
	}

	return instanceID, nil
}

// ConvertRegionToCreateInstanceRegion converts a string region to the OpenAPI enum type

// ConvertRegionToCreateInstanceRegion converts a string region to the OpenAPI enum type
func ConvertRegionToCreateInstanceRegion(region string) *models.CreateInstanceRequestRegion {
	switch strings.ToUpper(region) {
	case "EU":
		r := models.EU
		return &r
	case "US-EAST":
		r := models.USEast
		return &r
	case "US-WEST":
		r := models.USWest
		return &r
	case "US-CENTRAL":
		r := models.USCentral
		return &r
	case "AUS":
		r := models.AUS
		return &r
	case "SIN":
		r := models.SIN
		return &r
	case "JPN":
		r := models.JPN
		return &r
	case "UK":
		r := models.UK
		return &r
	case "IND":
		r := models.IND
		return &r
	default:
		// Default to EU if unknown region
		r := models.EU
		return &r
	}
}

// GetClusterNameFromDisplayName extracts cluster name from display name
// Returns empty string if not in provisioning or bound state
// EnsureClusterUUID ensures the cluster has a unique UUID label
// If the cluster doesn't have a UUID label, generates a new UUID v4 and returns it
// If it already has one, returns the existing UUID
func EnsureClusterUUID(cluster *infrastructurev1beta2.ContaboCluster) string {
	if cluster.Labels == nil {
		cluster.Labels = make(map[string]string)
	}

	if existingUUID, exists := cluster.Labels[ClusterUUIDLabel]; exists && existingUUID != "" {
		return existingUUID
	}

	// Generate new UUID v4
	newUUID := uuid.New().String()
	cluster.Labels[ClusterUUIDLabel] = newUUID
	return newUUID
}

// GetClusterUUID retrieves the cluster UUID from labels
// Returns empty string if no UUID is set
func GetClusterUUID(cluster *infrastructurev1beta2.ContaboCluster) string {
	if cluster.Labels == nil {
		return ""
	}
	return cluster.Labels[ClusterUUIDLabel]
}

// BuildShortClusterID generates a short identifier from cluster UUID
func BuildShortClusterID(clusterUUID string) string {
	if len(clusterUUID) >= 4 {
		return clusterUUID[:4]
	}
	return clusterUUID
}

// GenerateRequestID generates a UUID v4 for Contabo API request tracking
func GenerateRequestID() string {
	return uuid.New().String()
}

// decodeHTTPResponse decodes the body of an HTTP response into a target struct
func DecodeHTTPResponse[T any](resp *http.Response, err error) (*T, error) {
	var result *T
	if err != nil {
		return result, fmt.Errorf("HTTP request error: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return result, fmt.Errorf("HTTP request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return result, fmt.Errorf("failed to decode HTTP response: %w: %s", err, string(bodyBytes))
	}

	err = resp.Body.Close()
	if err != nil {
		return result, fmt.Errorf("failed to close response body: %w", err)
	}

	return result, nil
}

func mergeCloudConfig(file1, file2 []byte) ([]byte, error) {
	// Unmarshal both files into maps
	var config1, config2 map[string]interface{}
	if err := yaml.Unmarshal(file1, &config1); err != nil {
		return nil, fmt.Errorf("failed to unmarshal file1: %v", err)
	}
	if err := yaml.Unmarshal(file2, &config2); err != nil {
		return nil, fmt.Errorf("failed to unmarshal file2: %v", err)
	}

	// Merge config2 into config1 (config1 takes precedence)
	if err := mergo.Merge(&config1, config2, mergo.WithAppendSlice); err != nil {
		return nil, fmt.Errorf("failed to merge configs: %v", err)
	}

	// Marshal the merged config back to YAML
	return yaml.Marshal(config1)
}
