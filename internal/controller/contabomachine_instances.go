/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	infrastructurev1beta2 "github.com/ctnr-io/cluster-api-provider-contabo/api/v1beta2"
	"github.com/ctnr-io/cluster-api-provider-contabo/pkg/contabo/v1.0.0/models"
	corev1 "k8s.io/api/core/v1"
)

// reconcileInstance reconciles the Contabo instance
func (r *ContaboMachineReconciler) reconcileInstance(ctx context.Context, machine *clusterv1.Machine, contaboMachine *infrastructurev1beta2.ContaboMachine, contaboCluster *infrastructurev1beta2.ContaboCluster) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// Check if bootstrap data is available
	if machine.Spec.Bootstrap.DataSecretName == nil {
		log.Info("Waiting for bootstrap data")
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:   infrastructurev1beta2.InstanceReadyCondition,
			Status: metav1.ConditionFalse,
			Reason: infrastructurev1beta2.InstanceWaitingForBootstrapDataReason,
		})
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	// Check if cluster infrastructure is ready
	if !meta.IsStatusConditionTrue(contaboCluster.Status.Conditions, infrastructurev1beta2.ReadyCondition) {
		log.Info("Waiting for cluster infrastructure to be ready")
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:   infrastructurev1beta2.InstanceReadyCondition,
			Status: metav1.ConditionFalse,
			Reason: infrastructurev1beta2.InstanceWaitingForClusterInfrastructureReason,
		})
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	// Check if machine secrets are ready
	if !meta.IsStatusConditionTrue(contaboMachine.Status.Conditions, infrastructurev1beta2.MachineSecretsReadyCondition) {
		log.Info("Waiting for machine secrets to be ready")
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:   infrastructurev1beta2.InstanceReadyCondition,
			Status: metav1.ConditionFalse,
			Reason: infrastructurev1beta2.InstanceWaitingForMachineSecretsReason,
		})
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	// Check if private networks are ready
	if !meta.IsStatusConditionTrue(contaboMachine.Status.Conditions, infrastructurev1beta2.MachinePrivateNetworksReadyCondition) {
		log.Info("Waiting for private networks to be ready")
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:   infrastructurev1beta2.InstanceReadyCondition,
			Status: metav1.ConditionFalse,
			Reason: infrastructurev1beta2.InstanceWaitingForMachineSecretsReason,
		})
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	// Get the current instance condition
	instanceCondition := meta.FindStatusCondition(contaboMachine.Status.Conditions, infrastructurev1beta2.InstanceReadyCondition)

	// State machine for instance lifecycle
	switch {
	case instanceCondition == nil:
		// Initial state: Ensure instance (reuse or create)
		return r.ensureInstance(ctx, machine, contaboMachine, contaboCluster)

	case instanceCondition.Reason == infrastructurev1beta2.InstanceCreatingReason:
		// Check if instance creation is complete
		return r.checkInstanceCreation(ctx, contaboMachine)

	case instanceCondition.Reason == infrastructurev1beta2.InstanceProvisioningReason:
		// Check if instance is provisioned and ready
		return r.checkInstanceProvisioning(ctx, contaboMachine)

	case instanceCondition.Reason == infrastructurev1beta2.InstanceInstallingReason:
		// Check if instance installation/configuration is complete
		return r.checkInstanceInstallation(ctx, contaboMachine)

	case instanceCondition.Reason == infrastructurev1beta2.InstanceReadyReason:
		// Instance is ready, nothing to do
		return ctrl.Result{}, nil

	case instanceCondition.Reason == infrastructurev1beta2.InstanceFailedReason:
		// Instance failed, try to recover or ensure new instance
		return r.handleInstanceFailure(ctx, machine, contaboMachine, contaboCluster)

	default:
		// Unknown state, restart the process
		log.Info("Unknown instance state, restarting ensure process")
		return r.ensureInstance(ctx, machine, contaboMachine, contaboCluster)
	}
}

// ensureInstance ensures an instance is available (reuse existing or create new)
func (r *ContaboMachineReconciler) ensureInstance(ctx context.Context, machine *clusterv1.Machine, contaboMachine *infrastructurev1beta2.ContaboMachine, contaboCluster *infrastructurev1beta2.ContaboCluster) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	log.Info("Ensuring Contabo instance (reuse or create)")
	meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
		Type:   infrastructurev1beta2.InstanceReadyCondition,
		Status: metav1.ConditionFalse,
		Reason: infrastructurev1beta2.InstanceCreatingReason,
	})

	// First try to find an available instance
	availableInstance, err := r.findAvailableInstance(ctx, contaboMachine, contaboCluster)
	if err != nil {
		log.Error(err, "Failed to search for available instances")
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceFailedReason,
			Message: fmt.Sprintf("Failed to search for available instances: %s", err.Error()),
		})
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	if availableInstance != nil {
		// Reuse existing available instance
		log.Info("Reusing available instance", "instanceID", availableInstance.InstanceId)
		return r.reuseInstance(ctx, machine, contaboMachine, contaboCluster, *availableInstance)
	}

	// No available instance found, create new one
	log.Info("No available instances found, creating new instance")

	// Get bootstrap data for cloud-init
	bootstrapData, err := r.getBootstrapData(ctx, machine)
	if err != nil {
		log.Error(err, "Failed to get bootstrap data")
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceFailedReason,
			Message: fmt.Sprintf("Failed to get bootstrap data: %s", err.Error()),
		})
		return ctrl.Result{}, err
	}

	// Build create instance request
	createRequest := r.buildCreateInstanceRequest(contaboMachine, contaboCluster, bootstrapData)

	// Call Contabo API to create instance
	params := &models.CreateInstanceParams{}
	resp, err := r.ContaboClient.CreateInstance(ctx, params, createRequest)
	if err != nil {
		log.Error(err, "Failed to create instance")
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceFailedReason,
			Message: fmt.Sprintf("Failed to create instance: %s", err.Error()),
		})
		return ctrl.Result{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 && resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		err := fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
		log.Error(err, "Instance creation failed")
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceFailedReason,
			Message: fmt.Sprintf("Instance creation failed: %s", err.Error()),
		})
		return ctrl.Result{}, err
	}

	// Parse the response
	var createResponse models.CreateInstanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&createResponse); err != nil {
		log.Error(err, "Failed to decode instance creation response")
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceFailedReason,
			Message: fmt.Sprintf("Failed to decode response: %s", err.Error()),
		})
		return ctrl.Result{}, err
	}

	if len(createResponse.Data) == 0 {
		err := fmt.Errorf("no instance data in creation response")
		log.Error(err, "Instance creation failed")
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceFailedReason,
			Message: err.Error(),
		})
		return ctrl.Result{}, err
	}

	// Extract instance details and update machine status
	createdInstance := createResponse.Data[0]
	providerID := BuildProviderID(createdInstance.InstanceId)
	contaboMachine.Spec.ProviderID = providerID

	log.Info("Instance created successfully",
		"instanceId", createdInstance.InstanceId,
		"providerID", providerID)

	// Move to provisioning state
	meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
		Type:   infrastructurev1beta2.InstanceReadyCondition,
		Status: metav1.ConditionFalse,
		Reason: infrastructurev1beta2.InstanceProvisioningReason,
	})

	return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
}

// findAvailableInstance searches for available instances that can be reused
func (r *ContaboMachineReconciler) findAvailableInstance(ctx context.Context, contaboMachine *infrastructurev1beta2.ContaboMachine, contaboCluster *infrastructurev1beta2.ContaboCluster) (*models.ListInstancesResponseData, error) {
	log := logf.FromContext(ctx)

	// Generate request ID
	requestID := GenerateRequestID()

	// List all instances to find available ones
	params := &models.RetrieveInstancesListParams{
		XRequestId: requestID,
	}
	resp, err := r.ContaboClient.RetrieveInstancesList(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve instance list: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected response status: %d", resp.StatusCode)
	}

	var instanceList models.ListInstancesResponse
	if err := json.NewDecoder(resp.Body).Decode(&instanceList); err != nil {
		return nil, fmt.Errorf("failed to decode instance list: %w", err)
	}

	// Look for instances in available state with matching specs
	for _, instance := range instanceList.Data {
		// Check if instance is in available or reusable state
		if instance.DisplayName != "" {
			state := GetInstanceState(instance.DisplayName)
			instanceStatus := string(instance.Status)

			// Consider instance available for reuse if display name indicates available state
			isReusable := state == StateAvailable

			if isReusable {
				// Check if instance matches required specs
				if r.instanceMatchesSpecs(instance, contaboMachine, contaboCluster) {
					log.Info("Found matching reusable instance", "instanceID", instance.InstanceId, "displayName", instance.DisplayName, "status", instanceStatus, "state", state)
					return &instance, nil
				}
			}
		}
	}

	log.Info("No available instances found matching specifications")
	return nil, nil
}

// instanceMatchesSpecs checks if an instance matches the required specifications
func (r *ContaboMachineReconciler) instanceMatchesSpecs(instance models.ListInstancesResponseData, contaboMachine *infrastructurev1beta2.ContaboMachine, contaboCluster *infrastructurev1beta2.ContaboCluster) bool {
	// Check product ID
	if contaboMachine.Spec.Instance.ProductId != nil && instance.ProductId != *contaboMachine.Spec.Instance.ProductId {
		return false
	}

	// Check region
	expectedRegion := ConvertRegionToCreateInstanceRegion(contaboCluster.Spec.Region)
	if expectedRegion != nil && instance.Region != string(*expectedRegion) {
		return false
	}

	// Instance matches requirements
	return true
}

// reuseInstance configures an existing available instance for the machine
func (r *ContaboMachineReconciler) reuseInstance(ctx context.Context, machine *clusterv1.Machine, contaboMachine *infrastructurev1beta2.ContaboMachine, contaboCluster *infrastructurev1beta2.ContaboCluster, instance models.ListInstancesResponseData) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// Update display name to show it's being provisioned for this cluster
	clusterUUID := GetClusterUUID(contaboCluster)
	clusterID := BuildShortClusterID(clusterUUID)
	newDisplayName := BuildInstanceDisplayNameWithState(instance.InstanceId, StateProvisioning, clusterID)

	// Update instance display name via API
	requestID := GenerateRequestID()
	updateParams := &models.PatchInstanceParams{
		XRequestId: requestID,
	}
	updateRequest := models.PatchInstanceRequest{
		DisplayName: &newDisplayName,
	}

	resp, err := r.ContaboClient.PatchInstance(ctx, instance.InstanceId, updateParams, updateRequest)
	if err != nil {
		log.Error(err, "Failed to update instance display name", "instanceID", instance.InstanceId)
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}
	resp.Body.Close()

	// Set provider ID
	providerID := BuildProviderID(instance.InstanceId)
	contaboMachine.Spec.ProviderID = providerID

	log.Info("Successfully reused instance",
		"instanceId", instance.InstanceId,
		"providerID", providerID,
		"newDisplayName", newDisplayName)

	// Configure the instance with bootstrap data and move to provisioning
	return r.configureReusedInstance(ctx, machine, contaboMachine, contaboCluster, instance)
}

// configureReusedInstance configures a reused instance with bootstrap data
func (r *ContaboMachineReconciler) configureReusedInstance(ctx context.Context, machine *clusterv1.Machine, contaboMachine *infrastructurev1beta2.ContaboMachine, contaboCluster *infrastructurev1beta2.ContaboCluster, instance models.ListInstancesResponseData) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// Get bootstrap data
	bootstrapData, err := r.getBootstrapData(ctx, machine)
	if err != nil {
		log.Error(err, "Failed to get bootstrap data for reused instance")
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceFailedReason,
			Message: fmt.Sprintf("Failed to get bootstrap data: %s", err.Error()),
		})
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	// Apply bootstrap data to reused instance by reinstalling with new user data
	// This ensures the bootstrap script is applied to configure the node properly
	err = r.applyBootstrapDataToInstance(ctx, &instance, contaboMachine, contaboCluster, bootstrapData)
	if err != nil {
		log.Error(err, "Failed to apply bootstrap data to reused instance", "instanceID", instance.InstanceId)
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceFailedReason,
			Message: fmt.Sprintf("Failed to apply bootstrap data: %s", err.Error()),
		})
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	log.Info("Configuring reused instance", "instanceID", instance.InstanceId)

	// Move to provisioning state
	meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
		Type:   infrastructurev1beta2.InstanceReadyCondition,
		Status: metav1.ConditionFalse,
		Reason: infrastructurev1beta2.InstanceProvisioningReason,
	})

	return ctrl.Result{RequeueAfter: 15 * time.Second}, nil
}

// buildCreateInstanceRequest builds the request for creating a new instance
func (r *ContaboMachineReconciler) buildCreateInstanceRequest(contaboMachine *infrastructurev1beta2.ContaboMachine, contaboCluster *infrastructurev1beta2.ContaboCluster, bootstrapData string) models.CreateInstanceRequest {
	// Build instance display name with cluster context
	clusterUUID := GetClusterUUID(contaboCluster)
	clusterID := BuildShortClusterID(clusterUUID)
	displayName := BuildInstanceDisplayNameWithState(0, StateProvisioning, clusterID) // Will be updated with actual ID

	createRequest := models.CreateInstanceRequest{
		DisplayName:   &displayName,
		ProductId:     contaboMachine.Spec.Instance.ProductId,
		ImageId:       contaboMachine.Spec.Instance.ImageId,
		Region:        ConvertRegionToCreateInstanceRegion(contaboCluster.Spec.Region),
		UserData:      &bootstrapData,
		ApplicationId: contaboMachine.Spec.Instance.ApplicationId,
		DefaultUser:   nil,
		Period:        contaboMachine.Spec.Instance.Period,
		RootPassword:  contaboMachine.Spec.Instance.RootPassword,
		SshKeys:       contaboMachine.Spec.Instance.SshKeys,
		License:       nil,
		AddOns:        nil,
	}

	// Set default user
	if contaboMachine.Spec.Instance.DefaultUser != nil {
		defaultUser := models.CreateInstanceRequestDefaultUser(*contaboMachine.Spec.Instance.DefaultUser)
		createRequest.DefaultUser = &defaultUser
	}

	// Set license if specified
	if contaboMachine.Spec.Instance.License != nil {
		license := models.CreateInstanceRequestLicense(*contaboMachine.Spec.Instance.License)
		createRequest.License = &license
	}

	// Set addons if specified
	if contaboMachine.Spec.Instance.AddOns != nil {
		createRequest.AddOns = &models.CreateInstanceAddons{}
		if contaboMachine.Spec.Instance.AddOns.AdditionalIps != nil {
			value := map[string]interface{}{}
			createRequest.AddOns.AdditionalIps = &value
		}
		if contaboMachine.Spec.Instance.AddOns.AddonsIds != nil {
			var addonsIds []models.AddOnRequest
			for _, value := range *contaboMachine.Spec.Instance.AddOns.AddonsIds {
				addonsIds = append(addonsIds, models.AddOnRequest{
					Id:       value.Id,
					Quantity: value.Quantity,
				})
			}
			createRequest.AddOns.AddonsIds = &addonsIds
		}
		if contaboMachine.Spec.Instance.AddOns.Backup != nil {
			value := map[string]interface{}{}
			createRequest.AddOns.Backup = &value
		}
		if contaboMachine.Spec.Instance.AddOns.CustomImage != nil {
			value := map[string]interface{}{}
			createRequest.AddOns.CustomImage = &value
		}
		if contaboMachine.Spec.Instance.AddOns.ExtraStorage != nil {
			createRequest.AddOns.ExtraStorage = &models.ExtraStorageRequest{}
			if contaboMachine.Spec.Instance.AddOns.ExtraStorage.Nvme != nil {
				createRequest.AddOns.ExtraStorage.Nvme = contaboMachine.Spec.Instance.AddOns.ExtraStorage.Nvme
			}
			if contaboMachine.Spec.Instance.AddOns.ExtraStorage.Ssd != nil {
				createRequest.AddOns.ExtraStorage.Ssd = contaboMachine.Spec.Instance.AddOns.ExtraStorage.Ssd
			}
		}
		if contaboMachine.Spec.Instance.AddOns.PrivateNetworking != nil {
			value := map[string]interface{}{}
			createRequest.AddOns.PrivateNetworking = &value
		}
	}

	// Combine Secrets from cluster and machine specs
	var allSecrets []infrastructurev1beta2.ContaboSecretStatus

	// Add cluster-level secrets
	if len(contaboCluster.Status.Secrets) > 0 {
		allSecrets = append(allSecrets, contaboCluster.Status.Secrets...)
	}

	// Add machine-specific secrets
	if len(contaboMachine.Status.Secrets) > 0 {
		allSecrets = append(allSecrets, contaboMachine.Status.Secrets...)
	}

	// Retrieve and set password secret if specified
	for _, secret := range allSecrets {
		if secret.Type == infrastructurev1beta2.SecretResponseTypePassword {
			rootPassword := int64(secret.SecretId)
			createRequest.RootPassword = &rootPassword
			break
		}
	}

	// Retrieve and set SSH key secrets if specified
	var sshKeys []int64
	for _, secret := range allSecrets {
		if secret.Type == infrastructurev1beta2.SecretResponseTypeSsh {
			sshKeys = append(sshKeys, int64(secret.SecretId))
		}
	}
	// Remove duplicates
	sshKeys = removeDuplicateInt64s(sshKeys)
	if len(sshKeys) > 0 {
		createRequest.SshKeys = &sshKeys
	}

	// Enable private networking addon if specified
	if len(contaboCluster.Spec.PrivateNetworks)+len(contaboMachine.Spec.PrivateNetworks) > 0 {
		// Initialize AddOns if not set
		if createRequest.AddOns == nil {
			createRequest.AddOns = &models.CreateInstanceAddons{}
		}
		// Enable private networking (empty object as per API docs)
		createRequest.AddOns.PrivateNetworking = &map[string]interface{}{}
	} else {
		// Disable private networking if no networks are specified
		createRequest.AddOns.PrivateNetworking = nil
	}

	return createRequest
}

// getBootstrapData retrieves the bootstrap data for cloud-init
func (r *ContaboMachineReconciler) getBootstrapData(ctx context.Context, machine *clusterv1.Machine) (string, error) {
	if machine.Spec.Bootstrap.DataSecretName == nil {
		return "", fmt.Errorf("bootstrap data secret name is nil")
	}

	secret := &corev1.Secret{}
	key := client.ObjectKey{Namespace: machine.Namespace, Name: *machine.Spec.Bootstrap.DataSecretName}
	if err := r.Get(ctx, key, secret); err != nil {
		return "", fmt.Errorf("failed to retrieve bootstrap data secret: %w", err)
	}

	bootstrapData, exists := secret.Data["value"]
	if !exists {
		return "", fmt.Errorf("bootstrap data not found in secret")
	}

	return string(bootstrapData), nil
}

// checkInstanceCreation checks if instance creation is complete
func (r *ContaboMachineReconciler) checkInstanceCreation(ctx context.Context, contaboMachine *infrastructurev1beta2.ContaboMachine) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	if contaboMachine.Spec.ProviderID == "" {
		log.Info("No provider ID set yet, requeuing")
		return ctrl.Result{RequeueAfter: 5 * time.Second}, nil
	}

	instanceID, err := ParseProviderID(contaboMachine.Spec.ProviderID)
	if err != nil {
		log.Error(err, "Failed to parse provider ID")
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceFailedReason,
			Message: fmt.Sprintf("Invalid provider ID: %v", err),
		})
		return ctrl.Result{}, err
	}

	// Get instance details from Contabo API
	resp, err := r.ContaboClient.RetrieveInstance(ctx, instanceID, &models.RetrieveInstanceParams{})
	if err != nil {
		log.Error(err, "Failed to retrieve instance from Contabo API", "instanceID", instanceID)

		// Check if this is a 404 (instance not found yet)
		if resp != nil && resp.StatusCode == 404 {
			log.Info("Instance not found yet, still creating", "instanceID", instanceID)
			return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
		}

		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceFailedReason,
			Message: fmt.Sprintf("Failed to retrieve instance: %v", err),
		})
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Error(nil, "Unexpected response from Contabo API", "instanceID", instanceID, "statusCode", resp.StatusCode)
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	// Parse the response
	var instanceResp models.ListInstancesResponse
	if err := json.NewDecoder(resp.Body).Decode(&instanceResp); err != nil {
		log.Error(err, "Failed to parse instance response", "instanceID", instanceID)
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	if len(instanceResp.Data) == 0 {
		log.Error(nil, "No instance data returned", "instanceID", instanceID)
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	instance := instanceResp.Data[0]
	log.Info("Instance found", "instanceID", instanceID, "status", instance.Status)

	// Check instance status during creation phase
	switch string(instance.Status) {
	case "creating":
		log.Info("Instance still creating", "instanceID", instanceID)
		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil

	case "provisioning":
		log.Info("Instance creation complete, moving to provisioning", "instanceID", instanceID)
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:   infrastructurev1beta2.InstanceReadyCondition,
			Status: metav1.ConditionFalse,
			Reason: infrastructurev1beta2.InstanceProvisioningReason,
		})
		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil

	case "manual_provisioning":
		log.Info("Instance requires manual provisioning", "instanceID", instanceID)
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceFailedReason,
			Message: "Instance requires manual provisioning - contact Contabo support",
		})
		return ctrl.Result{RequeueAfter: 120 * time.Second}, nil

	case "installing":
		log.Info("Instance provisioning complete, moving to installing", "instanceID", instanceID)
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:   infrastructurev1beta2.InstanceReadyCondition,
			Status: metav1.ConditionFalse,
			Reason: infrastructurev1beta2.InstanceInstallingReason,
		})
		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil

	case "running":
		log.Info("Instance is running, checking installation status", "instanceID", instanceID)
		return r.checkInstanceInstallation(ctx, contaboMachine)

	case "stopped":
		log.Info("Instance is stopped during creation, waiting for it to start", "instanceID", instanceID)
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceProvisioningReason,
			Message: "Instance stopped during creation, waiting for start",
		})
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil

	case "error", "failed":
		log.Error(nil, "Instance creation failed", "instanceID", instanceID, "status", instance.Status)
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceFailedReason,
			Message: fmt.Sprintf("Instance creation failed with status: %s", string(instance.Status)),
		})
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil

	case "uninstalled":
		log.Info("Instance is uninstalled, needs reinstallation", "instanceID", instanceID)
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceProvisioningReason,
			Message: "Instance is uninstalled, waiting for reinstallation",
		})
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil

	case "product_not_available":
		log.Error(nil, "Product not available for instance", "instanceID", instanceID)
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceFailedReason,
			Message: "Product not available - may need different configuration",
		})
		return ctrl.Result{RequeueAfter: 300 * time.Second}, nil

	case "verification_required":
		log.Info("Instance requires verification", "instanceID", instanceID)
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceFailedReason,
			Message: "Instance requires account verification - contact Contabo support",
		})
		return ctrl.Result{RequeueAfter: 300 * time.Second}, nil

	case "pending_payment":
		log.Error(nil, "Instance pending payment", "instanceID", instanceID)
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceFailedReason,
			Message: "Instance pending payment - check Contabo account billing",
		})
		return ctrl.Result{RequeueAfter: 300 * time.Second}, nil

	case "rescue":
		log.Info("Instance is in rescue mode", "instanceID", instanceID)
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceProvisioningReason,
			Message: "Instance in rescue mode, waiting for normal boot",
		})
		return ctrl.Result{RequeueAfter: 60 * time.Second}, nil

	case "reset_password":
		log.Info("Instance password is being reset", "instanceID", instanceID)
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceProvisioningReason,
			Message: "Instance password being reset, waiting for completion",
		})
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil

	case "unknown":
		log.Info("Instance in unknown state, waiting", "instanceID", instanceID)
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceProvisioningReason,
			Message: "Instance in unknown state, waiting for stabilization",
		})
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil

	case "other":
		log.Info("Instance in other state, waiting", "instanceID", instanceID)
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceProvisioningReason,
			Message: "Instance in other state, waiting for known state",
		})
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil

	default:
		log.Info("Instance in unexpected state", "instanceID", instanceID, "status", instance.Status)
		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	}
}

// checkInstanceProvisioning checks if instance provisioning is complete
func (r *ContaboMachineReconciler) checkInstanceProvisioning(ctx context.Context, contaboMachine *infrastructurev1beta2.ContaboMachine) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	if contaboMachine.Spec.ProviderID == "" {
		log.Info("No provider ID set, cannot check provisioning status")
		return ctrl.Result{RequeueAfter: 5 * time.Second}, nil
	}

	instanceID, err := ParseProviderID(contaboMachine.Spec.ProviderID)
	if err != nil {
		log.Error(err, "Failed to parse provider ID")
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceFailedReason,
			Message: fmt.Sprintf("Invalid provider ID: %v", err),
		})
		return ctrl.Result{}, err
	}

	// Get instance details from Contabo API
	resp, err := r.ContaboClient.RetrieveInstance(ctx, instanceID, &models.RetrieveInstanceParams{})
	if err != nil {
		log.Error(err, "Failed to retrieve instance from Contabo API", "instanceID", instanceID)
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceFailedReason,
			Message: fmt.Sprintf("Failed to retrieve instance: %v", err),
		})
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Error(nil, "Unexpected response from Contabo API", "instanceID", instanceID, "statusCode", resp.StatusCode)
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	// Parse the response
	var instanceResp models.ListInstancesResponse
	if err := json.NewDecoder(resp.Body).Decode(&instanceResp); err != nil {
		log.Error(err, "Failed to parse instance response", "instanceID", instanceID)
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	if len(instanceResp.Data) == 0 {
		log.Error(nil, "No instance data returned", "instanceID", instanceID)
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	instance := instanceResp.Data[0]
	log.Info("Instance provisioning check", "instanceID", instanceID, "status", instance.Status)

	// Check instance status during provisioning phase
	switch string(instance.Status) {
	case "provisioning":
		log.Info("Instance still provisioning", "instanceID", instanceID)
		return ctrl.Result{RequeueAfter: 15 * time.Second}, nil

	case "manual_provisioning":
		log.Info("Instance requires manual provisioning during provisioning phase", "instanceID", instanceID)
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceFailedReason,
			Message: "Instance requires manual provisioning - contact Contabo support",
		})
		return ctrl.Result{RequeueAfter: 120 * time.Second}, nil

	case "installing":
		log.Info("Instance provisioning complete, moving to installing", "instanceID", instanceID)

		// Assign private networks if specified
		if len(contaboMachine.Spec.PrivateNetworks) > 0 {
			// Set condition to indicate private networks are being attached
			meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
				Type:   infrastructurev1beta2.MachinePrivateNetworksReadyCondition,
				Status: metav1.ConditionFalse,
				Reason: infrastructurev1beta2.MachinePrivateNetworkAttachingReason,
			})

			result, err := r.assignInstancePrivateNetworks(ctx, contaboMachine)
			if err != nil {
				log.Error(err, "Failed to assign private networks", "instanceID", instanceID)
				meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
					Type:    infrastructurev1beta2.MachinePrivateNetworksReadyCondition,
					Status:  metav1.ConditionFalse,
					Reason:  infrastructurev1beta2.MachinePrivateNetworkFailedReason,
					Message: fmt.Sprintf("Failed to assign private networks: %v", err),
				})
				meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
					Type:    infrastructurev1beta2.InstanceReadyCondition,
					Status:  metav1.ConditionFalse,
					Reason:  infrastructurev1beta2.InstanceFailedReason,
					Message: fmt.Sprintf("Failed to assign private networks: %v", err),
				})
				return result, err
			}

			// Mark private networks as ready after successful assignment
			meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
				Type:   infrastructurev1beta2.MachinePrivateNetworksReadyCondition,
				Status: metav1.ConditionTrue,
				Reason: infrastructurev1beta2.MachinePrivateNetworkReadyReason,
			})
		}

		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:   infrastructurev1beta2.InstanceReadyCondition,
			Status: metav1.ConditionFalse,
			Reason: infrastructurev1beta2.InstanceInstallingReason,
		})
		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil

	case "running":
		log.Info("Instance is running, checking installation status", "instanceID", instanceID)
		return r.checkInstanceInstallation(ctx, contaboMachine)

	case "stopped":
		log.Info("Instance stopped during provisioning", "instanceID", instanceID)
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceProvisioningReason,
			Message: "Instance stopped during provisioning, waiting for restart",
		})
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil

	case "error", "failed":
		log.Error(nil, "Instance provisioning failed", "instanceID", instanceID, "status", instance.Status)
		errorMessage := "unknown error"
		if instance.ErrorMessage != nil {
			errorMessage = *instance.ErrorMessage
		}
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceFailedReason,
			Message: fmt.Sprintf("Instance provisioning failed with status: %s, message: %s", string(instance.Status), errorMessage),
		})
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil

	case "uninstalled":
		log.Info("Instance uninstalled during provisioning, waiting for reinstall", "instanceID", instanceID)
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceProvisioningReason,
			Message: "Instance uninstalled during provisioning, waiting for reinstall",
		})
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil

	case "product_not_available":
		log.Error(nil, "Product not available during provisioning", "instanceID", instanceID)
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceFailedReason,
			Message: "Product not available during provisioning",
		})
		return ctrl.Result{RequeueAfter: 300 * time.Second}, nil

	case "verification_required":
		log.Info("Verification required during provisioning", "instanceID", instanceID)
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceFailedReason,
			Message: "Account verification required during provisioning",
		})
		return ctrl.Result{RequeueAfter: 300 * time.Second}, nil

	case "pending_payment":
		log.Error(nil, "Payment pending during provisioning", "instanceID", instanceID)
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceFailedReason,
			Message: "Payment pending during provisioning",
		})
		return ctrl.Result{RequeueAfter: 300 * time.Second}, nil

	case "rescue":
		log.Info("Instance in rescue mode during provisioning", "instanceID", instanceID)
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceProvisioningReason,
			Message: "Instance in rescue mode, waiting for normal boot",
		})
		return ctrl.Result{RequeueAfter: 60 * time.Second}, nil

	case "reset_password":
		log.Info("Instance password reset during provisioning", "instanceID", instanceID)
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceProvisioningReason,
			Message: "Password reset in progress during provisioning",
		})
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil

	case "unknown":
		log.Info("Instance in unknown state during provisioning", "instanceID", instanceID)
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceProvisioningReason,
			Message: "Instance in unknown state during provisioning",
		})
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil

	case "other":
		log.Info("Instance in other state during provisioning", "instanceID", instanceID)
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceProvisioningReason,
			Message: "Instance in other state during provisioning",
		})
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil

	case "creating":
		log.Info("Instance moved back to creating state, checking creation status", "instanceID", instanceID)
		return r.checkInstanceCreation(ctx, contaboMachine)

	default:
		log.Info("Instance in unexpected state during provisioning", "instanceID", instanceID, "status", instance.Status)
		return ctrl.Result{RequeueAfter: 15 * time.Second}, nil
	}
}

// checkInstanceInstallation checks if instance installation is complete
func (r *ContaboMachineReconciler) checkInstanceInstallation(ctx context.Context, contaboMachine *infrastructurev1beta2.ContaboMachine) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	if contaboMachine.Spec.ProviderID == "" {
		log.Info("No provider ID set, cannot check installation status")
		return ctrl.Result{RequeueAfter: 5 * time.Second}, nil
	}

	instanceID, err := ParseProviderID(contaboMachine.Spec.ProviderID)
	if err != nil {
		log.Error(err, "Failed to parse provider ID")
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceFailedReason,
			Message: fmt.Sprintf("Invalid provider ID: %v", err),
		})
		return ctrl.Result{}, err
	}

	// Get instance details from Contabo API
	resp, err := r.ContaboClient.RetrieveInstance(ctx, instanceID, &models.RetrieveInstanceParams{})
	if err != nil {
		log.Error(err, "Failed to retrieve instance from Contabo API", "instanceID", instanceID)
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceFailedReason,
			Message: fmt.Sprintf("Failed to retrieve instance: %v", err),
		})
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Error(nil, "Unexpected response from Contabo API", "instanceID", instanceID, "statusCode", resp.StatusCode)
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	// Parse the response
	var instanceResp models.ListInstancesResponse
	if err := json.NewDecoder(resp.Body).Decode(&instanceResp); err != nil {
		log.Error(err, "Failed to parse instance response", "instanceID", instanceID)
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	if len(instanceResp.Data) == 0 {
		log.Error(nil, "No instance data returned", "instanceID", instanceID)
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	instance := instanceResp.Data[0]
	log.Info("Instance installation check", "instanceID", instanceID, "status", instance.Status)

	// Check instance status during installation phase
	switch string(instance.Status) {
	case "installing":
		log.Info("Instance still installing", "instanceID", instanceID)
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil

	case "running":
		// TODO: Add SSH connectivity check here
		// TODO: Check that cloud-init has completed successfully
		// Instance is running, verify cloud-init has completed
		log.Info("Instance is running, checking if ready for use", "instanceID", instanceID)

		// Check if this is a new transition to running state
		readyCondition := meta.FindStatusCondition(contaboMachine.Status.Conditions, infrastructurev1beta2.InstanceReadyCondition)
		if readyCondition == nil || readyCondition.Reason != infrastructurev1beta2.InstanceReadyReason {
			// First time we see it as running, give it some time to complete cloud-init
			log.Info("Instance just started running, waiting for cloud-init to complete", "instanceID", instanceID)
			meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
				Type:    infrastructurev1beta2.InstanceReadyCondition,
				Status:  metav1.ConditionFalse,
				Reason:  infrastructurev1beta2.InstanceInstallingReason,
				Message: "Waiting for cloud-init to complete",
			})
			return ctrl.Result{RequeueAfter: 60 * time.Second}, nil
		}

		// If we've been waiting and it's still running, assume it's ready
		log.Info("Instance installation complete, instance ready", "instanceID", instanceID)
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:   infrastructurev1beta2.InstanceReadyCondition,
			Status: metav1.ConditionTrue,
			Reason: infrastructurev1beta2.InstanceReadyReason,
		})
		return ctrl.Result{}, nil

	case "stopped":
		log.Info("Instance stopped during installation, waiting for restart", "instanceID", instanceID)
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceInstallingReason,
			Message: "Instance stopped during installation, waiting for restart",
		})
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil

	case "error", "failed":
		log.Error(nil, "Instance installation failed", "instanceID", instanceID, "status", instance.Status)
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceFailedReason,
			Message: fmt.Sprintf("Instance installation failed with status: %s", string(instance.Status)),
		})
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil

	case "uninstalled":
		log.Info("Instance uninstalled during installation, waiting for reinstall", "instanceID", instanceID)
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceInstallingReason,
			Message: "Instance uninstalled during installation, waiting for reinstall",
		})
		return ctrl.Result{RequeueAfter: 60 * time.Second}, nil

	case "product_not_available":
		log.Error(nil, "Product not available during installation", "instanceID", instanceID)
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceFailedReason,
			Message: "Product not available during installation",
		})
		return ctrl.Result{RequeueAfter: 300 * time.Second}, nil

	case "verification_required":
		log.Info("Verification required during installation", "instanceID", instanceID)
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceFailedReason,
			Message: "Account verification required during installation",
		})
		return ctrl.Result{RequeueAfter: 300 * time.Second}, nil

	case "pending_payment":
		log.Error(nil, "Payment pending during installation", "instanceID", instanceID)
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceFailedReason,
			Message: "Payment pending during installation",
		})
		return ctrl.Result{RequeueAfter: 300 * time.Second}, nil

	case "rescue":
		log.Info("Instance in rescue mode during installation", "instanceID", instanceID)
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceInstallingReason,
			Message: "Instance in rescue mode, waiting for normal boot",
		})
		return ctrl.Result{RequeueAfter: 60 * time.Second}, nil

	case "reset_password":
		log.Info("Instance password reset during installation", "instanceID", instanceID)
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceInstallingReason,
			Message: "Password reset in progress during installation",
		})
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil

	case "unknown":
		log.Info("Instance in unknown state during installation", "instanceID", instanceID)
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceInstallingReason,
			Message: "Instance in unknown state during installation",
		})
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil

	case "other":
		log.Info("Instance in other state during installation", "instanceID", instanceID)
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceInstallingReason,
			Message: "Instance in other state during installation",
		})
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil

	case "manual_provisioning":
		log.Info("Manual provisioning required during installation", "instanceID", instanceID)
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceFailedReason,
			Message: "Manual provisioning required during installation",
		})
		return ctrl.Result{RequeueAfter: 120 * time.Second}, nil

	case "provisioning":
		log.Info("Instance moved back to provisioning state, checking provisioning status", "instanceID", instanceID)
		return r.checkInstanceProvisioning(ctx, contaboMachine)

	case "creating":
		log.Info("Instance moved back to creating state, checking creation status", "instanceID", instanceID)
		return r.checkInstanceCreation(ctx, contaboMachine)

	default:
		log.Info("Instance in unexpected state during installation", "instanceID", instanceID, "status", instance.Status)
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}
}

// handleInstanceFailure handles failed instances
func (r *ContaboMachineReconciler) handleInstanceFailure(ctx context.Context, machine *clusterv1.Machine, contaboMachine *infrastructurev1beta2.ContaboMachine, contaboCluster *infrastructurev1beta2.ContaboCluster) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	log.Info("Handling instance failure - attempting recovery")

	// Check if we have a provider ID (existing instance)
	if contaboMachine.Spec.ProviderID != "" {
		instanceID, err := ParseProviderID(contaboMachine.Spec.ProviderID)
		if err != nil {
			log.Error(err, "Failed to parse provider ID for failed instance")
		} else {
			log.Info("Attempting to clean up failed instance", "instanceID", instanceID)

			// Try to get instance details to check if it still exists
			resp, err := r.ContaboClient.RetrieveInstance(ctx, instanceID, &models.RetrieveInstanceParams{})
			if err == nil && resp.StatusCode == 200 {
				defer resp.Body.Close()

				// Parse response to check instance state
				var instanceResp models.ListInstancesResponse
				if err := json.NewDecoder(resp.Body).Decode(&instanceResp); err == nil && len(instanceResp.Data) > 0 {
					instance := instanceResp.Data[0]
					log.Info("Failed instance still exists", "instanceID", instanceID, "status", instance.Status)

					// If instance is in a terminal failure state, attempt reinstallation
					if string(instance.Status) == "error" || string(instance.Status) == "failed" {
						log.Info("Attempting to reinstall failed instance", "instanceID", instanceID)

						// Try to reinstall the instance (this would require implementing reinstall logic)
						// For now, we'll mark it for recreation by clearing the provider ID
						contaboMachine.Spec.ProviderID = ""

						meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
							Type:    infrastructurev1beta2.InstanceReadyCondition,
							Status:  metav1.ConditionFalse,
							Reason:  infrastructurev1beta2.InstanceCreatingReason,
							Message: "Recreating instance after failure",
						})

						// Retry creation after clearing provider ID
						return ctrl.Result{RequeueAfter: 60 * time.Second}, nil
					}
				}
			} else {
				// Instance no longer exists or can't be reached, clear provider ID
				log.Info("Failed instance no longer exists, will create new one", "instanceID", instanceID)
				contaboMachine.Spec.ProviderID = ""
			}
		}
	}

	// Check retry count to prevent infinite loops
	failureCondition := meta.FindStatusCondition(contaboMachine.Status.Conditions, infrastructurev1beta2.InstanceReadyCondition)
	if failureCondition != nil && failureCondition.Reason == infrastructurev1beta2.InstanceFailedReason {
		// Check how long we've been in failure state
		timeSinceFailure := time.Since(failureCondition.LastTransitionTime.Time)

		// If we've been failing for more than 30 minutes, increase retry interval
		if timeSinceFailure > 30*time.Minute {
			log.Info("Instance has been failing for extended period, using longer retry interval", "timeSinceFailure", timeSinceFailure)

			meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
				Type:    infrastructurev1beta2.InstanceReadyCondition,
				Status:  metav1.ConditionFalse,
				Reason:  infrastructurev1beta2.InstanceCreatingReason,
				Message: "Retrying instance creation after extended failure",
			})

			return r.ensureInstance(ctx, machine, contaboMachine, contaboCluster)
		}
	}

	// Set appropriate condition for retry
	meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
		Type:    infrastructurev1beta2.InstanceReadyCondition,
		Status:  metav1.ConditionFalse,
		Reason:  infrastructurev1beta2.InstanceCreatingReason,
		Message: "Retrying instance creation after failure",
	})

	// Retry ensuring instance after a delay
	log.Info("Retrying instance ensure after failure")
	return r.ensureInstance(ctx, machine, contaboMachine, contaboCluster)
}

// applyBootstrapDataToInstance applies bootstrap data to a reused instance via reinstall
func (r *ContaboMachineReconciler) applyBootstrapDataToInstance(ctx context.Context, instance *models.ListInstancesResponseData, contaboMachine *infrastructurev1beta2.ContaboMachine, contaboCluster *infrastructurev1beta2.ContaboCluster, bootstrapData string) error {
	log := logf.FromContext(ctx)

	// Build create request to extract SSH keys and other settings
	createRequest := r.buildCreateInstanceRequest(contaboMachine, contaboCluster, bootstrapData)

	reinstallRequest := models.ReinstallInstanceRequest{
		ApplicationId: createRequest.ApplicationId,
		RootPassword:  createRequest.RootPassword,
		SshKeys:       createRequest.SshKeys,
		UserData:      createRequest.UserData,
		DefaultUser:   nil,
		ImageId:       contaboMachine.Status.Instance.ImageId,
	}

	// Set default User
	if createRequest.DefaultUser != nil {
		defaultUser := models.ReinstallInstanceRequestDefaultUser(string(*createRequest.DefaultUser))
		reinstallRequest.DefaultUser = &defaultUser
	}

	// Set image ID
	if createRequest.ImageId != nil {
		reinstallRequest.ImageId = *createRequest.ImageId
	}

	// Handle Backup addon if specified
	if createRequest.AddOns.Backup != nil {
		// Add addons for backup
		resp, err := r.ContaboClient.UpgradeInstance(ctx, instance.InstanceId, &models.UpgradeInstanceParams{}, models.UpgradeInstanceRequest{
			Backup: &models.Backup{},
		})
		if err != nil {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to upgrade instance %d: %w\n%s", instance.InstanceId, err, body)
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("upgrade instance returned non-success status: %d\n%s", resp.StatusCode, body)
		}
		log.Info("Backup addon applied to instance", "instanceID", instance.InstanceId)
	}

	// FIXME: don't know how to set custom image addon
	// // Handle Custom Image addon if specified
	// if (createRequest.AddOns.CustomImage != nil) {
	// 	// Add addons for custom image
	// 	resp, err := r.ContaboClient.UpgradeInstance(ctx, instance.InstanceId, &models.UpgradeInstanceParams{}, models.UpgradeInstanceRequest{
	// 		CustomImage: &map[string]interface{}{},
	// 	})
	// 	if err != nil {
	// 		body, _ := io.ReadAll(resp.Body)
	// 		return fmt.Errorf("failed to upgrade instance %d: %w\n%s", instance.InstanceId, err, body)
	// 	}
	// 	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
	// 		body, _ := io.ReadAll(resp.Body)
	// 		return fmt.Errorf("upgrade instance returned non-success status: %d\n%s", resp.StatusCode, body)
	// 	}
	// 	log.Info("Custom Image addon applied to instance", "instanceID", instance.InstanceId)
	// }

	// FIXME: don't know how to set additional storage addon
	// Handle Extra Storage addon
	// if createRequest.AddOns.ExtraStorage != nil {
	// 	// Add addons for extra storage
	// 	resp, err := r.ContaboClient.UpgradeInstance(ctx, instance.InstanceId, &models.UpgradeInstanceParams{}, models.UpgradeInstanceRequest{
	// 		ExtraStorage: &models.ExtraStorage{},
	// 	})
	// 	if err != nil {
	// 		body, _ := io.ReadAll(resp.Body)
	// 		return fmt.Errorf("failed to upgrade instance %d: %w\n%s", instance.InstanceId, err, body)
	// 	}
	// 	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
	// 		body, _ := io.ReadAll(resp.Body)
	// 		return fmt.Errorf("upgrade instance returned non-success status: %d\n%s", resp.StatusCode, body)
	// 	}
	// 	log.Info("Extra storage addon applied to instance", "instanceID", instance.InstanceId)
	// }

	// Handle private networking addon if specified
	if createRequest.AddOns.PrivateNetworking != nil {
		// Add addons for private networking
		resp, err := r.ContaboClient.UpgradeInstance(ctx, instance.InstanceId, &models.UpgradeInstanceParams{}, models.UpgradeInstanceRequest{
			PrivateNetworking: &models.PrivateNetworkingUpgradeRequest{},
		})
		if err != nil {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to upgrade instance %d: %w\n%s", instance.InstanceId, err, body)
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("upgrade instance returned non-success status: %d\n%s", resp.StatusCode, body)
		}
		log.Info("Private networking addon applied to instance", "instanceID", instance.InstanceId)
	}

	log.Info("Reinstalling instance with bootstrap data", reinstallRequest)

	// Perform reinstall
	params := &models.ReinstallInstanceParams{}
	resp, err := r.ContaboClient.ReinstallInstance(ctx, instance.InstanceId, params, reinstallRequest)
	if err != nil {
		return fmt.Errorf("failed to reinstall instance %d: %w", instance.InstanceId, err)
	}

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("reinstall instance returned non-success status: %d", resp.StatusCode)
	}

	log.Info("Instance reinstall initiated successfully", "instanceID", instance.InstanceId)
	return nil
}


// assignPrivateNetworks assigns private networks to the instance via Contabo API
func (r *ContaboMachineReconciler) assignInstancePrivateNetworks(ctx context.Context, contaboMachine *infrastructurev1beta2.ContaboMachine) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// Get instance ID from provider ID
	instanceID, err := ParseProviderID(contaboMachine.Spec.ProviderID)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to parse provider ID: %w", err)
	}

	log.Info("Assigning private networks to instance", "instanceID", instanceID, "networkCount", len(contaboMachine.Spec.PrivateNetworks))

	// Assign each private network to the instance
	for _, networkSpec := range contaboMachine.Spec.PrivateNetworks {
		// Find network status to get the network ID
		var networkID int64
		found := false
		for _, networkStatus := range contaboMachine.Status.PrivateNetworks {
			if networkStatus.Name == networkSpec.Name {
				networkID = networkStatus.PrivateNetworkId
				found = true
				break
			}
		}

		if !found {
			return ctrl.Result{}, fmt.Errorf("private network %s not found in machine status", networkSpec.Name)
		}

		log.Info("Assigning private network to instance", "instanceID", instanceID, "networkID", networkID, "networkName", networkSpec.Name)

		// Prepare request parameters
		requestID := GenerateRequestID()
		params := &models.AssignInstancePrivateNetworkParams{
			XRequestId: requestID,
		}

		// Make API call to assign private network
		resp, err := r.ContaboClient.AssignInstancePrivateNetwork(ctx, networkID, instanceID, params)
		if err != nil {
			log.Error(err, "Failed to assign private network", "instanceID", instanceID, "networkID", networkID, "networkName", networkSpec.Name)
			return ctrl.Result{RequeueAfter: 30 * time.Second}, fmt.Errorf("failed to assign private network %s: %w", networkSpec.Name, err)
		}

		resp.Body.Close()

		if resp.StatusCode != 201 && resp.StatusCode != 200 {
			log.Error(nil, "Unexpected response when assigning private network", "instanceID", instanceID, "networkID", networkID, "statusCode", resp.StatusCode)
			return ctrl.Result{RequeueAfter: 30 * time.Second}, fmt.Errorf("failed to assign private network %s, status code: %d", networkSpec.Name, resp.StatusCode)
		}

		log.Info("Successfully assigned private network", "instanceID", instanceID, "networkID", networkID, "networkName", networkSpec.Name)
	}

	return ctrl.Result{}, nil
}