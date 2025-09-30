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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ContaboMachineSpec defines the desired state of ContaboMachine
type ContaboMachineSpec struct {
	// ProviderID is the unique identifier as specified by the cloud provider.
	// +optional
	ProviderID string `json:"providerID"`

	// Instance is the type of instance to create.
	Instance ContaboInstanceSpec `json:"instance"`
}

// ContaboMachineStatus defines the observed state of ContaboMachine.
type ContaboMachineStatus struct {
	// Ready is true when the provider resource is ready.
	// +optional
	Ready bool `json:"ready"`

	// Instance is the current state of the Contabo instance.
	Instance *ContaboInstanceStatus `json:"instance,omitempty"`

	// Addresses contains the Contabo instance associated addresses.
	Addresses []clusterv1.MachineAddress `json:"addresses,omitempty"`

	// Conditions defines current service state of the ContaboMachine.
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// Initialization
	Initialization *ContaboMachineInitializationStatus `json:"initialization,omitempty"`

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

type ContaboMachineInitializationStatus struct {
	// Provisioned indicates if the initialization is complete
	Provisioned bool `json:"provisioned"`

	// ErrorMessage provides details in case of initialization failure
	ErrorMessage *string `json:"errorMessage,omitempty"`
}

// ContaboInstanceSpec defines the desired state of a Contabo instance
type ContaboInstanceSpec struct {
	// ProductID is the Contabo product ID (instance type)
	// +kubebuilder:validation:Required
	ProductId string `json:"productId"`
}

// ContaboInstanceStatus defines the observed state of a Contabo instance
type ContaboInstanceStatus = InstanceResponse

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".metadata.labels.cluster\\.x-k8s\\.io/cluster-name",description="Cluster to which this ContaboMachine belongs"
// +kubebuilder:printcolumn:name="InstanceType",type="string",JSONPath=".spec.instanceType",description="Contabo instance type"
// +kubebuilder:printcolumn:name="State",type="string",JSONPath=".status.instanceState",description="Contabo instance state"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="Machine ready status"
// +kubebuilder:printcolumn:name="ProviderID",type="string",JSONPath=".spec.providerID",description="Contabo instance ID"
// +kubebuilder:printcolumn:name="Machine",type="string",JSONPath=".metadata.ownerReferences[?(@.kind==\"Machine\")].name",description="Machine object which owns this ContaboMachine"
// +kubebuilder:resource:path=contabomachines,scope=Namespaced,categories=cluster-api
// +kubebuilder:storageversion

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
	conditions := make(clusterv1.Conditions, len(m.Status.Conditions))
	for i, condition := range m.Status.Conditions {
		conditions[i] = clusterv1.Condition{
			Type:               clusterv1.ConditionType(condition.Type),
			Status:             corev1.ConditionStatus(condition.Status),
			LastTransitionTime: condition.LastTransitionTime,
			Reason:             condition.Reason,
			Message:            condition.Message,
		}
	}
	return conditions
}

// SetConditions sets the conditions of the ContaboMachine.
func (m *ContaboMachine) SetConditions(conditions clusterv1.Conditions) {
	m.Status.Conditions = make([]metav1.Condition, len(conditions))
	for i, condition := range conditions {
		m.Status.Conditions[i] = metav1.Condition{
			Type:               string(condition.Type),
			Status:             metav1.ConditionStatus(condition.Status),
			LastTransitionTime: condition.LastTransitionTime,
			Reason:             condition.Reason,
			Message:            condition.Message,
		}
	}
}
