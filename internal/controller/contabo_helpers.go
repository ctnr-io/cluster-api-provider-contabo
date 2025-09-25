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

	infrastructurev1beta1 "github.com/ctnr-io/cluster-api-provider-contabo/api/v1beta1"
	"github.com/ctnr-io/cluster-api-provider-contabo/pkg/contabo/models"
)

const (
	// ProviderIDPrefix is the prefix for Contabo provider IDs
	ProviderIDPrefix = "contabo://"

	// StateAvailable indicates an instance is available for cluster assignment
	StateAvailable = "available"

	// StateClusterBound indicates an instance is bound to a cluster
	StateClusterBound = "cluster-bound"

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
// Display names are formatted as: "<original-name>-<state>-<cluster-name>"
func GetInstanceState(displayName string) string {
	parts := strings.Split(displayName, "-")
	if len(parts) >= 2 {
		// Look for known state values
		for i, part := range parts {
			if part == StateAvailable || part == StateClusterBound {
				return part
			}
			// Also check for "state" pattern
			if i > 0 && (part == "available" || part == "bound") {
				if part == "available" {
					return StateAvailable
				}
				return StateClusterBound
			}
		}
	}
	return StateAvailable // default to available if no state found
}

// MapInstanceStatusToMachineState maps Contabo instance status to ContaboMachineInstanceState
func MapInstanceStatusToMachineState(status models.InstanceStatus) infrastructurev1beta1.ContaboMachineInstanceState {
	switch status {
	case models.InstanceStatusRunning:
		return infrastructurev1beta1.ContaboMachineInstanceStateRunning
	case models.InstanceStatusStopped:
		return infrastructurev1beta1.ContaboMachineInstanceStateStopped
	case models.InstanceStatusProvisioning, models.InstanceStatusInstalling:
		return infrastructurev1beta1.ContaboMachineInstanceStatePending
	case models.InstanceStatusError, models.InstanceStatusUnknown:
		return infrastructurev1beta1.ContaboMachineInstanceStateUnknown
	default:
		return infrastructurev1beta1.ContaboMachineInstanceStateUnknown
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

// BuildInstanceDisplayName constructs a display name for an instance based on state and cluster
func BuildInstanceDisplayName(baseName, state, clusterName string) string {
	if state == StateClusterBound && clusterName != "" {
		return fmt.Sprintf("%s-%s-%s", baseName, state, clusterName)
	}
	return fmt.Sprintf("%s-%s", baseName, state)
}
