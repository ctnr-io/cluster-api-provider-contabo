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
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
)

const (
	// ClusterFinalizer allows the controller to clean up resources associated with ContaboCluster before
	// removing it from the apiserver.
	ClusterFinalizer = "contabocluster.infrastructure.cluster.x-k8s.io"

	// MachineFinalizer allows the controller to clean up resources associated with ContaboMachine before
	// removing it from the apiserver.
	MachineFinalizer = "contabomachine.infrastructure.cluster.x-k8s.io"
)

// Condition types and reasons for ContaboCluster.
const (
	// ReadyCondition indicates the cluster infrastructure is ready.
	ReadyCondition = clusterv1.ReadyCondition

	// NetworkInfrastructureReadyCondition indicates the network infrastructure is ready.
	NetworkInfrastructureReadyCondition = "NetworkInfrastructureReady"

	// CreatingReason indicates that the cluster infrastructure is being created.
	CreatingReason = "Creating"

	// NetworkInfrastructureFailedReason indicates that the network infrastructure failed to be created.
	NetworkInfrastructureFailedReason = "NetworkInfrastructureFailed"

	// NetworkInfrastructureCreatingReason indicates that the network infrastructure is being created.
	NetworkInfrastructureCreatingReason = "NetworkInfrastructureCreating"

	// NetworkInfrastructureReadyReason indicates that the network infrastructure is ready.
	NetworkInfrastructureReadyReason = "NetworkInfrastructureReady"

	// NetworkInfrastructureSkippedReason indicates that network infrastructure configuration was skipped.
	NetworkInfrastructureSkippedReason = "NetworkInfrastructureSkipped"

	// WaitingForControlPlaneEndpointReason indicates that the cluster is waiting for the control plane endpoint to be set.
	WaitingForControlPlaneEndpointReason = "WaitingForControlPlaneEndpoint"
)

// Condition types and reasons for ContaboMachine.
const (
	// MachineReadyCondition indicates the machine infrastructure is ready.
	MachineReadyCondition = clusterv1.MachineReadyCondition

	// InstanceReadyCondition indicates the instance is ready.
	InstanceReadyCondition = "InstanceReady"

	// InstanceProvisioningReason indicates that the instance is being provisioned.
	InstanceProvisioningReason = "InstanceProvisioning"

	// InstanceProvisioningFailedReason indicates that the instance failed to be provisioned.
	InstanceProvisioningFailedReason = "InstanceProvisioningFailed"

	// InstanceNotFoundReason indicates that the instance was not found.
	InstanceNotFoundReason = "InstanceNotFound"

	// InstanceDeletingReason indicates that the instance is being deleted.
	InstanceDeletingReason = "InstanceDeleting"

	// InstanceDeletionFailedReason indicates that the instance failed to be deleted.
	InstanceDeletionFailedReason = "InstanceDeletionFailed"

	// WaitingForClusterInfrastructureReason indicates that the machine is waiting for cluster infrastructure to be ready.
	WaitingForClusterInfrastructureReason = "WaitingForClusterInfrastructure"

	// WaitingForBootstrapDataReason indicates that the machine is waiting for bootstrap data to be ready.
	WaitingForBootstrapDataReason = "WaitingForBootstrapData"
)
