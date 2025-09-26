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
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ContaboMachineTemplateSpec defines the desired state of ContaboMachineTemplate
type ContaboMachineTemplateSpec struct {
	Template ContaboMachineTemplateResource `json:"template"`
}

// ContaboMachineTemplateResource describes the data needed to create a ContaboMachine from a template
type ContaboMachineTemplateResource struct {
	Spec ContaboMachineSpec `json:"spec"`
}

// ContaboMachineTemplateStatus defines the observed state of ContaboMachineTemplate.
// NOTE: ContaboMachineTemplate is a template resource and does not have a status.
type ContaboMachineTemplateStatus struct {
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=contabomachinetemplates,scope=Namespaced,categories=cluster-api
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Time duration since creation of ContaboMachineTemplate"

// ContaboMachineTemplate is the Schema for the contabomachinetemplates API
type ContaboMachineTemplate struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec defines the desired state of ContaboMachineTemplate
	// +required
	Spec ContaboMachineTemplateSpec `json:"spec"`
}

// +kubebuilder:object:root=true

// ContaboMachineTemplateList contains a list of ContaboMachineTemplate
type ContaboMachineTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ContaboMachineTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ContaboMachineTemplate{}, &ContaboMachineTemplateList{})
}
