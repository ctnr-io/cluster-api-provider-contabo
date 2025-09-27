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

package v1beta2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
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

	// PrivateNetworks specifies the private network configuration for the cluster.
	// +optional
	PrivateNetworks []ContaboPrivateNetworkSpec `json:"privateNetworks,omitempty"`

	// Secrets references the secret to create to configure the cluster.
	// +optional
	Secrets []ContaboSecretSpec `json:"secrets,omitempty"`
}

// ContaboClusterStatus defines the observed state of ContaboCluster.
type ContaboClusterStatus struct {
	// Ready denotes that the cluster (infrastructure) is ready.
	// +optional
	Ready bool `json:"ready"`

	// PrivateNetworkds contains the discovered information about private networks
	// +optional
	PrivateNetworks []ContaboPrivateNetworkStatus `json:"privateNetworks,omitempty"`

	// Secrets contains the references to secrets used by the machine.
	// +optional
	Secrets []ContaboSecretStatus `json:"secrets,omitempty"`

	// Conditions defines current service state of the ContaboCluster.
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// FailureDomains is a list of failure domains that machines can be placed in.
	// +optional
	FailureDomains []clusterv1.FailureDomain `json:"failureDomains,omitempty"`
}

// ContaboPrivateNetworkSpec defines the desired state of a Contabo private network
type ContaboPrivateNetworkSpec = CreatePrivateNetworkRequest

// ContaboPrivateNetworkStatus defines the observed state of a Contabo private network
type ContaboPrivateNetworkStatus = PrivateNetworkResponse

// ContaboSecretSpec defines the desired state of a Contabo secret
type ContaboSecretSpec = CreateSecretRequest

// ContaboSecretStatus defines the observed state of a Contabo secret
type ContaboSecretStatus = SecretResponse

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".metadata.labels.cluster\\.x-k8s\\.io/cluster-name",description="Cluster to which this ContaboCluster belongs"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="Cluster infrastructure is ready"
// +kubebuilder:printcolumn:name="Private Networks",type="string",JSONPath=".status.privateNetworks[].name",description="Private Networks"
// +kubebuilder:printcolumn:name="Endpoint",type="string",JSONPath=".spec.controlPlaneEndpoint.host",description="API Endpoint",priority=1
// +kubebuilder:resource:path=contaboclusters,scope=Namespaced,categories=cluster-api
// +kubebuilder:storageversion

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
	Status *ContaboClusterStatus `json:"status,omitempty,omitzero"`
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
func (c *ContaboCluster) GetConditions() []metav1.Condition {
	return c.Status.Conditions
}

// SetConditions sets the conditions of the ContaboCluster.
func (c *ContaboCluster) SetConditions(conditions []metav1.Condition) {
	c.Status.Conditions = conditions
}
