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

// =============================================================================
// CONTABO CLUSTER CONDITIONS
// =============================================================================

// ContaboCluster condition types.
const (
	// ReadyCondition indicates the cluster infrastructure is ready.
	ReadyCondition = clusterv1.ReadyCondition

	// ControlPlaneEndpointReadyCondition indicates the control plane endpoint is ready.
	ControlPlaneEndpointReadyCondition = "ControlPlaneEndpointReady"

	// ClusterPrivateNetworkReadyCondition indicates the cluster private networks are ready.
	ClusterPrivateNetworkReadyCondition = "ClusterPrivateNetworkReady"

	// ClusterSecretsReadyCondition indicates the cluster secrets are ready.
	ClusterSecretsReadyCondition = "ClusterSecretsReady"
)

// ContaboCluster condition reasons.
const (
	// CreatingReason indicates that the cluster infrastructure is being created.
	CreatingReason = "Creating"

	// AvailableReason indicates that the cluster infrastructure is ready and available.
	AvailableReason = clusterv1.AvailableReason
)

// Control plane endpoint condition reasons.
const (
	// ControlPlaneEndpointCreatingReason indicates the control plane endpoint is being created.
	ControlPlaneEndpointCreatingReason = "ControlPlaneEndpointCreating"

	// ControlPlaneEndpointReadyReason indicates the control plane endpoint is ready.
	ControlPlaneEndpointReadyReason = "ControlPlaneEndpointReady"

	// ControlPlaneEndpointFailedReason indicates the control plane endpoint failed.
	ControlPlaneEndpointFailedReason = "ControlPlaneEndpointFailed"

	// WaitingForControlPlaneEndpointReason indicates waiting for the control plane endpoint to be set.
	WaitingForControlPlaneEndpointReason = "WaitingForControlPlaneEndpoint"
)

// Cluster private network condition reasons.
const (
	// ClusterPrivateNetworkCreatingReason indicates cluster private networks are being created.
	ClusterPrivateNetworkCreatingReason = "ClusterPrivateNetworkCreating"

	// ClusterPrivateNetworkReadyReason indicates cluster private networks are ready.
	ClusterPrivateNetworkReadyReason = "ClusterPrivateNetworkReady"

	// ClusterPrivateNetworkFailedReason indicates cluster private networks failed.
	ClusterPrivateNetworkFailedReason = "ClusterPrivateNetworkFailed"

	// ClusterPrivateNetworkDeletingReason indicates cluster private networks are being deleted.
	ClusterPrivateNetworkDeletingReason = "ClusterPrivateNetworkDeleting"

	// ClusterPrivateNetworkSkippedReason indicates cluster private network configuration was skipped.
	ClusterPrivateNetworkSkippedReason = "ClusterPrivateNetworkSkipped"
)

// Cluster secrets condition reasons.
const (
	// ClusterSecretsCreatingReason indicates cluster secrets are being created.
	ClusterSecretsCreatingReason = "ClusterSecretsCreating"

	// ClusterSecretsReadyReason indicates cluster secrets are ready.
	ClusterSecretsReadyReason = "ClusterSecretsReady"

	// ClusterSecretsFailedReason indicates cluster secrets failed.
	ClusterSecretsFailedReason = "ClusterSecretsFailed"

	// ClusterSecretsDeletingReason indicates cluster secrets are being deleted.
	ClusterSecretsDeletingReason = "ClusterSecretsDeleting"

	// ClusterSecretsSkippedReason indicates cluster secrets configuration was skipped.
	ClusterSecretsSkippedReason = "ClusterSecretsSkipped"
)

// =============================================================================
// CONTABO MACHINE CONDITIONS
// =============================================================================

// ContaboMachine condition types.
const (
	// MachineReadyCondition indicates the machine infrastructure is ready.
	MachineReadyCondition = clusterv1.MachineReadyCondition

	// InstanceReadyCondition indicates the Contabo instance is ready.
	InstanceReadyCondition = "InstanceReady"

	// MachinePrivateNetworksReadyCondition indicates the machine private networks are ready.
	MachinePrivateNetworksReadyCondition = "MachinePrivateNetworksReady"

	// MachineSecretsReadyCondition indicates the machine secrets are ready.
	MachineSecretsReadyCondition = "MachineSecretsReady"

	// MachineSSHKeysReadyCondition indicates the machine SSH keys are configured.
	MachineSSHKeysReadyCondition = "MachineSSHKeysReady"

	// ClusterInfrastructureReadyCondition indicates the cluster infrastructure dependency is ready.
	ClusterInfrastructureReadyCondition = "ClusterInfrastructureReady"
)

// Instance condition reasons.
const (
	// InstanceWaitingForClusterInfrastructureReason indicates waiting for cluster infrastructure to be ready.
	InstanceWaitingForClusterInfrastructureReason = "InstanceWaitingForClusterInfrastructure"

	// InstanceWaitingForBootstrapDataReason indicates waiting for bootstrap data to be ready.
	InstanceWaitingForBootstrapDataReason = "InstanceWaitingForBootstrapData"

	// InstanceWaitingForMachineSecretsReason indicates waiting for machine secrets to be ready.
	InstanceWaitingForMachineSecretsReason = "InstanceWaitingForMachineSecrets"

	// InstanceWaitingForPrivateNetworksReason indicates waiting for machine private networks to be ready.
	InstanceWaitingForPrivateNetworksReason = "InstanceWaitingForPrivateNetworks"

	// InstanceCreatingReason indicates the instance is being created.
	InstanceCreatingReason = "InstanceCreating"

	// InstanceProvisioningReason indicates the instance is being provisioned.
	InstanceProvisioningReason = "InstanceProvisioning"

	// InstanceInstallingReason indicates the instance is being installed/configured.
	InstanceInstallingReason = "InstanceInstalling"

	// InstanceReadyReason indicates the instance is ready.
	InstanceReadyReason = "InstanceReady"

	// InstanceFailedReason indicates the instance failed to be created or provisioned.
	InstanceFailedReason = "InstanceFailed"

	// InstanceDeletingReason indicates the instance is being deleted.
	InstanceDeletingReason = "InstanceDeleting"

	// InstanceNotFoundReason indicates the instance was not found.
	InstanceNotFoundReason = "InstanceNotFound"

	// InstanceReinstallingReason indicates the instance is being reinstalled.
	InstanceReinstallingReason = "InstanceReinstalling"
)

// Machine private network condition reasons.
const (
	// MachinePrivateNetworkCreatingReason indicates machine private networks are being created.
	MachinePrivateNetworkCreatingReason = "MachinePrivateNetworkCreating"

	// MachinePrivateNetworkAttachingReason indicates machine private networks are being attached.
	MachinePrivateNetworkAttachingReason = "MachinePrivateNetworkAttaching"

	// MachinePrivateNetworkReadyReason indicates machine private networks are ready and attached.
	MachinePrivateNetworkReadyReason = "MachinePrivateNetworkReady"

	// MachinePrivateNetworkFailedReason indicates machine private networks failed.
	MachinePrivateNetworkFailedReason = "MachinePrivateNetworkFailed"

	// MachinePrivateNetworkDetachingReason indicates machine private networks are being detached.
	MachinePrivateNetworkDetachingReason = "MachinePrivateNetworkDetaching"

	// MachinePrivateNetworkSkippedReason indicates machine private network configuration was skipped.
	MachinePrivateNetworkSkippedReason = "MachinePrivateNetworkSkipped"
)

// Machine secrets condition reasons.
const (
	// MachineSecretsCreatingReason indicates machine secrets are being created.
	MachineSecretsCreatingReason = "MachineSecretsCreating"

	// MachineSecretsReadyReason indicates machine secrets are ready.
	MachineSecretsReadyReason = "MachineSecretsReady"

	// MachineSecretsFailedReason indicates machine secrets failed.
	MachineSecretsFailedReason = "MachineSecretsFailed"

	// MachineSecretsDeletingReason indicates machine secrets are being deleted.
	MachineSecretsDeletingReason = "MachineSecretsDeleting"

	// MachineSecretsSkippedReason indicates machine secrets configuration was skipped.
	MachineSecretsSkippedReason = "MachineSecretsSkipped"
)

// Machine SSH keys condition reasons.
const (
	// MachineSSHKeysConfiguringReason indicates machine SSH keys are being configured.
	MachineSSHKeysConfiguringReason = "MachineSSHKeysConfiguring"

	// MachineSSHKeysReadyReason indicates machine SSH keys are configured and ready.
	MachineSSHKeysReadyReason = "MachineSSHKeysReady"

	// MachineSSHKeysFailedReason indicates machine SSH keys configuration failed.
	MachineSSHKeysFailedReason = "MachineSSHKeysFailed"

	// MachineSSHKeysSkippedReason indicates machine SSH keys configuration was skipped.
	MachineSSHKeysSkippedReason = "MachineSSHKeysSkipped"

	// MachineSSHKeysUpdatingReason indicates machine SSH keys are being updated.
	MachineSSHKeysUpdatingReason = "MachineSSHKeysUpdating"
)

// Cluster infrastructure dependency condition reasons.
const (
	// WaitingForClusterInfrastructureReason indicates waiting for cluster infrastructure to be ready.
	WaitingForClusterInfrastructureReason = "WaitingForClusterInfrastructure"

	// WaitingForBootstrapDataReason indicates waiting for bootstrap data to be ready.
	WaitingForBootstrapDataReason = "WaitingForBootstrapData"

	// ClusterInfrastructureReadyReason indicates the cluster infrastructure is ready.
	ClusterInfrastructureReadyReason = "ClusterInfrastructureReady"

	// ClusterInfrastructureFailedReason indicates the cluster infrastructure failed.
	ClusterInfrastructureFailedReason = "ClusterInfrastructureFailed"
)
