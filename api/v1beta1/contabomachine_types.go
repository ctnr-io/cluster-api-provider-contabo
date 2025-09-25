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

// ContaboMachineSpec defines the desired state of ContaboMachine
type ContaboMachineSpec struct {
	// ProviderID is the unique identifier as specified by the cloud provider.
	// +optional
	ProviderID *string `json:"providerID,omitempty"`

	// InstanceType is the Contabo instance type for the machine.
	// +kubebuilder:validation:Required
	InstanceType string `json:"instanceType"`

	// Region is the Contabo region where the machine will be created.
	// +kubebuilder:validation:Required
	Region string `json:"region"`

	// Note: Image is standardized to Ubuntu 22.04 LTS and not configurable per machine

	// SSHKeys is a list of SSH key names to be added to the machine.
	// +optional
	SSHKeys []string `json:"sshKeys,omitempty"`

	// UserData is the cloud-init user data to use for the machine.
	// +optional
	UserData *string `json:"userData,omitempty"`

	// AdditionalMetadata is additional metadata to be added to the machine.
	// +optional
	AdditionalMetadata map[string]string `json:"additionalMetadata,omitempty"`

	// AdditionalTags is additional tags to be added to the machine.
	// +optional
	AdditionalTags map[string]string `json:"additionalTags,omitempty"`

	// Network configuration for the machine
	// +optional
	Network *ContaboMachineNetworkSpec `json:"network,omitempty"`
}

// ContaboMachineNetworkSpec defines network configuration for a machine
type ContaboMachineNetworkSpec struct {
	// SubnetName is the name of the subnet to place the machine in
	// +optional
	SubnetName *string `json:"subnetName,omitempty"`

	// PrivateIP is the private IP address to assign to the machine
	// +optional
	PrivateIP *string `json:"privateIP,omitempty"`

	// PublicIP specifies whether to assign a public IP to the machine
	// +optional
	PublicIP *bool `json:"publicIP,omitempty"`
}

// ContaboMachineStatus defines the observed state of ContaboMachine.
type ContaboMachineStatus struct {
	// Ready is true when the provider resource is ready.
	// +optional
	Ready bool `json:"ready"`

	// InstanceState is the state of the Contabo instance.
	// +optional
	InstanceState *ContaboMachineInstanceState `json:"instanceState,omitempty"`

	// Network describes the machine network configuration
	// +optional
	Network *ContaboMachineNetworkStatus `json:"network,omitempty"`

	// Addresses contains the Contabo instance associated addresses.
	Addresses []clusterv1.MachineAddress `json:"addresses,omitempty"`

	// Conditions defines current service state of the ContaboMachine.
	// +optional
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`

	// FailureReason will be set in the event that there is a terminal problem
	// reconciling the Machine and will contain a succinct value suitable
	// for machine interpretation.
	//
	// This field should not be set for transitive errors that a controller
	// faces that are expected to be fixed automatically over
	// time (like service outages), but instead indicate that something is
	// fundamentally wrong with the Machine's spec or the configuration of
	// the controller, and that manual intervention is required. Examples
	// of terminal errors would be invalid combinations of settings in the
	// spec, values that are unsupported by the controller, or the
	// responsible controller itself being critically misconfigured.
	//
	// Any transient errors that occur during the reconciliation of Machines
	// can be added as events to the Machine object and/or logged in the
	// controller's output.
	// +optional
	FailureReason *string `json:"failureReason,omitempty"`

	// FailureMessage will be set in the event that there is a terminal problem
	// reconciling the Machine and will contain a more verbose string suitable
	// for logging and human consumption.
	//
	// This field should not be set for transitive errors that a controller
	// faces that are expected to be fixed automatically over
	// time (like service outages), but instead indicate that something is
	// fundamentally wrong with the Machine's spec or the configuration of
	// the controller, and that manual intervention is required. Examples
	// of terminal errors would be invalid combinations of settings in the
	// spec, values that are unsupported by the controller, or the
	// responsible controller itself being critically misconfigured.
	//
	// Any transient errors that occur during the reconciliation of Machines
	// can be added as events to the Machine object and/or logged in the
	// controller's output.
	// +optional
	FailureMessage *string `json:"failureMessage,omitempty"`
}

// ContaboMachineInstanceState describes the instance state
type ContaboMachineInstanceState string

const (
	// ContaboMachineInstanceStateRunning represents a running instance
	ContaboMachineInstanceStateRunning ContaboMachineInstanceState = "running"
	// ContaboMachineInstanceStatePending represents a pending instance
	ContaboMachineInstanceStatePending ContaboMachineInstanceState = "pending"
	// ContaboMachineInstanceStateStopped represents a stopped instance
	ContaboMachineInstanceStateStopped ContaboMachineInstanceState = "stopped"
	// ContaboMachineInstanceStateTerminated represents a terminated instance
	ContaboMachineInstanceStateTerminated ContaboMachineInstanceState = "terminated"
	// ContaboMachineInstanceStateUnknown represents an unknown instance state
	ContaboMachineInstanceStateUnknown ContaboMachineInstanceState = "unknown"
)

// ContaboMachineNetworkStatus defines the network status
type ContaboMachineNetworkStatus struct {
	// PrivateIP is the private IP address of the machine
	// +optional
	PrivateIP *string `json:"privateIP,omitempty"`

	// PublicIP is the public IP address of the machine
	// +optional
	PublicIP *string `json:"publicIP,omitempty"`

	// SubnetID is the ID of the subnet the machine is in
	// +optional
	SubnetID *string `json:"subnetID,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".metadata.labels.cluster\\.x-k8s\\.io/cluster-name",description="Cluster to which this ContaboMachine belongs"
// +kubebuilder:printcolumn:name="InstanceType",type="string",JSONPath=".spec.instanceType",description="Contabo instance type"
// +kubebuilder:printcolumn:name="State",type="string",JSONPath=".status.instanceState",description="Contabo instance state"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="Machine ready status"
// +kubebuilder:printcolumn:name="ProviderID",type="string",JSONPath=".spec.providerID",description="Contabo instance ID"
// +kubebuilder:printcolumn:name="Machine",type="string",JSONPath=".metadata.ownerReferences[?(@.kind==\"Machine\")].name",description="Machine object which owns this ContaboMachine"

// ContaboMachine is the Schema for the contabomachines API
type ContaboMachine struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of ContaboMachine
	// +required
	Spec ContaboMachineSpec `json:"spec"`

	// status defines the observed state of ContaboMachine
	// +optional
	Status ContaboMachineStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// ContaboMachineList contains a list of ContaboMachine
type ContaboMachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ContaboMachine `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ContaboMachine{}, &ContaboMachineList{})
}

// GetConditions returns the conditions of the ContaboMachine.
func (m *ContaboMachine) GetConditions() clusterv1.Conditions {
	return m.Status.Conditions
}

// SetConditions sets the conditions of the ContaboMachine.
func (m *ContaboMachine) SetConditions(conditions clusterv1.Conditions) {
	m.Status.Conditions = conditions
}
