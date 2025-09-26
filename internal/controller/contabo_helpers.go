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
	"fmt"
	"strconv"
	"strings"

	infrastructurev1beta2 "github.com/ctnr-io/cluster-api-provider-contabo/api/v1beta2"
	"github.com/ctnr-io/cluster-api-provider-contabo/pkg/contabo/v1.0.0/models"
	"github.com/google/uuid"
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

	// Display name state strings (used in parsing) - kept short to save space
	stateAvailable    = "avl"
	stateProvisioning = "prv"
	stateBound        = "bnd"
	stateError        = "err"

	// DefaultUbuntuImageID is the standardized Ubuntu image used for all cluster nodes
	// Using a fixed image ensures consistency, security, and compatibility across the cluster
	// This should be Ubuntu 24.04 LTS - you may need to adjust this ID based on available images in Contabo
	DefaultUbuntuImageID = "d64d5c6c-9dda-4e38-8174-0ee282474d8a"
)

// BuildProviderID constructs a provider ID from an instance ID
func BuildProviderID(instanceID int64) string {
	return fmt.Sprintf("%s%d", ProviderIDPrefix, instanceID)
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

// GetInstanceState extracts the state from an instance display name
// Format: "<instance-id>-<state>-<cluster-id>"
// Returns: "capc-available", "capc-provisioning", "capc-cluster-bound", or "capc-error"
func GetInstanceState(displayName string) string {
	parts := strings.Split(displayName, "-")
	if len(parts) >= 2 {
		// Try parts[1] first (format: instanceID-state-cluster)
		switch parts[1] {
		case stateAvailable:
			return StateAvailable
		case stateProvisioning:
			return StateProvisioning
		case stateBound:
			return StateClusterBound
		case stateError:
			return StateError
		}
	}
	return StateAvailable // default
}

// MapInstanceStatusToMachineState maps Contabo instance status to ContaboMachineInstanceState
func MapInstanceStatusToMachineState(status models.InstanceStatus) infrastructurev1beta2.ContaboMachineInstanceState {
	switch status {
	case models.InstanceStatusRunning:
		return infrastructurev1beta2.ContaboMachineInstanceStateRunning
	case models.InstanceStatusStopped:
		return infrastructurev1beta2.ContaboMachineInstanceStateStopped
	case models.InstanceStatusProvisioning, models.InstanceStatusInstalling:
		return infrastructurev1beta2.ContaboMachineInstanceStatePending
	case models.InstanceStatusError, models.InstanceStatusUnknown:
		return infrastructurev1beta2.ContaboMachineInstanceStateUnknown
	default:
		return infrastructurev1beta2.ContaboMachineInstanceStateUnknown
	}
}

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

// BuildInstanceDisplayName creates a descriptive display name for a Contabo instance
// Format: "<cluster-name>-<short-cluster-id>-<machine-name>" (shortened for space)
func BuildInstanceDisplayName(cluster *infrastructurev1beta2.ContaboCluster, machineName string) string {
	clusterUUID := GetClusterUUID(cluster)
	shortID := BuildShortClusterID(clusterUUID)
	return fmt.Sprintf("%s-%s-%s", cluster.Name, shortID, machineName)
}

// BuildInstanceDisplayNameWithState creates a display name including instance state
// Format: "<instance-id>-<state>-<cluster-id>" (shortened for space)
func BuildInstanceDisplayNameWithState(instanceID int64, state, clusterID string) string {
	stateShort := mapStateToShort(state)
	if stateShort != "" && clusterID != "" {
		return fmt.Sprintf("%d-%s-%s", instanceID, stateShort, clusterID)
	} else if stateShort != "" {
		return fmt.Sprintf("%d-%s", instanceID, stateShort)
	}
	return fmt.Sprintf("%d", instanceID)
}

// mapStateToShort converts full state names to short versions for display names
func mapStateToShort(state string) string {
	switch state {
	case StateAvailable:
		return stateAvailable
	case StateProvisioning:
		return stateProvisioning
	case StateClusterBound:
		return stateBound
	case StateError:
		return stateError
	default:
		return stateAvailable
	}
}

// GetClusterNameFromDisplayName extracts cluster name from display name
// Returns empty string if not in provisioning or bound state
// GetClusterIDFromDisplayName extracts cluster ID from display name
// Returns empty string if not in provisioning or bound state
// Expects format: "<instance-id>-<state>-<cluster-id>"
func GetClusterIDFromDisplayName(displayName string) string {
	parts := strings.Split(displayName, "-")
	if len(parts) >= 3 && (parts[1] == stateProvisioning || parts[1] == stateBound) {
		return strings.Join(parts[2:], "-")
	}
	return ""
}

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
