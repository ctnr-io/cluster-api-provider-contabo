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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ContaboClusterSpec defines the desired state of ContaboCluster
type ContaboClusterSpec struct {
	// Region is the Contabo region where the cluster will be deployed
	// +kubebuilder:validation:Required
	Region string `json:"region"`

	// ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.
	// +optional
	ControlPlaneEndpoint clusterv1.APIEndpoint `json:"controlPlaneEndpoint,omitempty"`

	// Network configuration for the cluster
	// +optional
	Network *ContaboNetworkSpec `json:"network,omitempty"`

	// ClusterTagging configures cluster-wide tagging for instance management
	// +optional
	ClusterTagging *ContaboClusterTaggingSpec `json:"clusterTagging,omitempty"`
}

// ContaboNetworkSpec defines network configuration for Contabo resources
type ContaboNetworkSpec struct {
	// VPC ID for the cluster resources
	// +optional
	VPCID *string `json:"vpcId,omitempty"`

	// Subnet configuration
	// +optional
	Subnets []ContaboSubnetSpec `json:"subnets,omitempty"`
}

// ContaboSubnetSpec defines a subnet configuration
type ContaboSubnetSpec struct {
	// Name of the subnet
	Name string `json:"name"`

	// CIDR block for the subnet
	CIDR string `json:"cidr"`

	// Availability zone for the subnet
	// +optional
	AvailabilityZone *string `json:"availabilityZone,omitempty"`
}

// ContaboClusterTaggingSpec defines cluster membership tagging configuration
type ContaboClusterTaggingSpec struct {
	// Enabled specifies whether cluster membership tagging is enabled
	// +kubebuilder:default=true
	// +optional
	Enabled *bool `json:"enabled,omitempty"`

	// ClusterTag is the tag name to use for cluster membership
	// +kubebuilder:default="cluster-api-cluster"
	// +optional
	ClusterTag string `json:"clusterTag,omitempty"`

	// AvailableTag is the tag name to use when marking instances as available
	// +kubebuilder:default="cluster-api-available"
	// +optional
	AvailableTag string `json:"availableTag,omitempty"`

	// TagColor is the color to use for cluster membership tags
	// +kubebuilder:default="#1E90FF"
	// +optional
	TagColor string `json:"tagColor,omitempty"`
}

// ContaboClusterStatus defines the observed state of ContaboCluster.
type ContaboClusterStatus struct {
	// Ready denotes that the cluster (infrastructure) is ready.
	// +optional
	Ready bool `json:"ready"`

	// Network describes the cluster network configuration
	// +optional
	Network *ContaboNetworkStatus `json:"network,omitempty"`

	// Conditions defines current service state of the ContaboCluster.
	// +optional
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`

	// FailureDomains is a list of failure domains that machines can be placed in.
	// +optional
	FailureDomains clusterv1.FailureDomains `json:"failureDomains,omitempty"`
}

// ContaboNetworkStatus defines the network status
type ContaboNetworkStatus struct {
	// VPCID is the ID of the VPC
	// +optional
	VPCID *string `json:"vpcId,omitempty"`

	// Subnets is a list of subnets
	// +optional
	Subnets []ContaboSubnetStatus `json:"subnets,omitempty"`
}

// ContaboSubnetStatus defines the status of a subnet
type ContaboSubnetStatus struct {
	// Name of the subnet
	Name string `json:"name"`

	// ID of the subnet
	ID string `json:"id"`

	// CIDR block of the subnet
	CIDR string `json:"cidr"`

	// AvailabilityZone of the subnet
	// +optional
	AvailabilityZone *string `json:"availabilityZone,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".metadata.labels.cluster\\.x-k8s\\.io/cluster-name",description="Cluster to which this ContaboCluster belongs"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="Cluster infrastructure is ready"
// +kubebuilder:printcolumn:name="VPC",type="string",JSONPath=".status.network.vpcId",description="VPC ID"
// +kubebuilder:printcolumn:name="Endpoint",type="string",JSONPath=".spec.controlPlaneEndpoint.host",description="API Endpoint",priority=1

// ContaboCluster is the Schema for the contaboclusters API
type ContaboCluster struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of ContaboCluster
	// +required
	Spec ContaboClusterSpec `json:"spec"`

	// status defines the observed state of ContaboCluster
	// +optional
	Status ContaboClusterStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// ContaboClusterList contains a list of ContaboCluster
type ContaboClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ContaboCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ContaboCluster{}, &ContaboClusterList{})
}

// GetConditions returns the conditions of the ContaboCluster.
func (c *ContaboCluster) GetConditions() clusterv1.Conditions {
	return c.Status.Conditions
}

// SetConditions sets the conditions of the ContaboCluster.
func (c *ContaboCluster) SetConditions(conditions clusterv1.Conditions) {
	c.Status.Conditions = conditions
}
