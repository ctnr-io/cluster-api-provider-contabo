/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS 	// Wait for the instance to be ready
	instance, err := r.waitForInstanceReady(ctx, stubInstance.InstanceId, "provisioning")
	if err != nil {
		return nil, fmt.Errorf("failed to wait for instance ready: %w", err)
	}

	return instance, nilIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/cluster-api/util/predicates"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	infrastructurev1beta2 "github.com/ctnr-io/cluster-api-provider-contabo/api/v1beta2"
	contaboclient "github.com/ctnr-io/cluster-api-provider-contabo/pkg/contabo/client"
	"github.com/ctnr-io/cluster-api-provider-contabo/pkg/contabo/models"
)

// ContaboMachineReconciler reconciles a ContaboMachine object
type ContaboMachineReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	Recorder      record.EventRecorder
	ContaboClient *contaboclient.Client
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=contabomachines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=contabomachines/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=contabomachines/finalizers,verbs=update
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machines;machines/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters;clusters/status,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *ContaboMachineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// Fetch the ContaboMachine instance
	contaboMachine := &infrastructurev1beta2.ContaboMachine{}
	if err := r.Get(ctx, req.NamespacedName, contaboMachine); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Fetch the Machine
	machine, err := util.GetOwnerMachine(ctx, r.Client, contaboMachine.ObjectMeta)
	if err != nil {
		return ctrl.Result{}, err
	}
	if machine == nil {
		log.Info("Machine Controller has not yet set OwnerRef, requeuing",
			"contaboMachine", contaboMachine.Name,
			"namespace", contaboMachine.Namespace,
			"ownerReferences", contaboMachine.OwnerReferences)

		// Requeue after 10 seconds to allow Machine controller to set OwnerRef
		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	}

	log = log.WithValues("machine", machine.Name)

	// Fetch the Cluster
	cluster, err := util.GetClusterFromMetadata(ctx, r.Client, machine.ObjectMeta)
	if err != nil {
		log.Info("Machine is missing cluster label or cluster does not exist")
		return ctrl.Result{}, nil
	}

	if annotations.IsPaused(cluster, contaboMachine) {
		log.Info("ContaboMachine or linked Cluster is marked as paused. Won't reconcile")
		return ctrl.Result{}, nil
	}

	log = log.WithValues("cluster", cluster.Name)

	contaboCluster := &infrastructurev1beta2.ContaboCluster{}
	contaboClusterName := client.ObjectKey{
		Namespace: contaboMachine.Namespace,
		Name:      cluster.Spec.InfrastructureRef.Name,
	}
	if err := r.Get(ctx, contaboClusterName, contaboCluster); err != nil {
		log.Info("ContaboCluster is not available yet")
		return ctrl.Result{}, nil
	}

	// Initialize the patch helper
	patchHelper, err := patch.NewHelper(contaboMachine, r.Client)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Always attempt to Patch the ContaboMachine object and status after each reconciliation
	defer func() {
		if err := patchHelper.Patch(ctx, contaboMachine); err != nil {
			log.Error(err, "failed to patch ContaboMachine")
		}
	}()

	// Handle deleted machines
	if !contaboMachine.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, machine, contaboMachine, contaboCluster)
	}

	// Handle non-deleted machines
	return r.reconcileNormal(ctx, machine, contaboMachine, contaboCluster)
}

func (r *ContaboMachineReconciler) reconcileNormal(ctx context.Context, machine *clusterv1.Machine, contaboMachine *infrastructurev1beta2.ContaboMachine, contaboCluster *infrastructurev1beta2.ContaboCluster) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// Automatically add cluster label to ContaboMachine for proper mapping
	if contaboMachine.Labels == nil {
		contaboMachine.Labels = make(map[string]string)
	}
	if machine.Spec.ClusterName != "" {
		contaboMachine.Labels[clusterv1.ClusterNameLabel] = machine.Spec.ClusterName
		log.V(4).Info("Added cluster label to ContaboMachine",
			"clusterName", machine.Spec.ClusterName,
			"label", clusterv1.ClusterNameLabel)
	}

	// If the ContaboMachine doesn't have our finalizer, add it.
	controllerutil.AddFinalizer(contaboMachine, infrastructurev1beta2.MachineFinalizer)

	// Always ensure status is initialized
	if contaboMachine.Status.Conditions == nil {
		contaboMachine.Status.Conditions = []metav1.Condition{}
	}

	// Check if the infrastructure is ready, otherwise return and wait for the cluster object to be updated
	if !contaboCluster.Status.Ready {
		log.Info("ContaboCluster is not ready yet")
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:   infrastructurev1beta2.InstanceReadyCondition,
			Status: metav1.ConditionFalse,
			Reason: infrastructurev1beta2.WaitingForClusterInfrastructureReason,
		})
		return ctrl.Result{}, nil
	}

	// Make sure bootstrap data is available and populated
	if machine.Spec.Bootstrap.DataSecretName == nil {
		log.Info("Bootstrap data secret reference is not yet available")
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:   infrastructurev1beta2.InstanceReadyCondition,
			Status: metav1.ConditionFalse,
			Reason: infrastructurev1beta2.WaitingForBootstrapDataReason,
		})
		return ctrl.Result{}, nil
	}

	// Get bootstrap data from the secret
	userData, err := r.getBootstrapData(ctx, machine)
	if err != nil {
		log.Error(err, "failed to get bootstrap data")
		return ctrl.Result{}, err
	}

	// Create or update the instance
	instance, err := r.reconcileInstance(ctx, machine, contaboMachine, contaboCluster, userData)
	if err != nil {
		log.Error(err, "failed to reconcile instance")

		// If we have an instance ID, try to reset its state to error or available
		if contaboMachine.Spec.ProviderID != nil {
			if instanceID, parseErr := ParseProviderID(*contaboMachine.Spec.ProviderID); parseErr == nil {
				// Mark instance as error if it was a serious failure (installation/configuration)
				if strings.Contains(err.Error(), "install") || strings.Contains(err.Error(), "configure") || strings.Contains(err.Error(), "timeout") {
					log.Info("Marking instance as error due to installation/configuration failure", "instanceId", instanceID)
					if updateErr := r.updateInstanceDisplayName(ctx, instanceID, StateError, ""); updateErr != nil {
						log.Error(updateErr, "failed to mark instance as error")
					}
				} else {
					// For other failures, mark as available for retry
					log.Info("Marking instance as available for retry", "instanceId", instanceID)
					if updateErr := r.updateInstanceDisplayName(ctx, instanceID, StateAvailable, ""); updateErr != nil {
						log.Error(updateErr, "failed to mark instance as available")
					}
				}
			}
		}

		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceProvisioningFailedReason,
			Message: fmt.Sprintf("Failed to reconcile instance: %s", err.Error()),
		})
		return ctrl.Result{}, err
	}

	// Apply cluster membership via displayName using cluster UUID
	clusterUUID := GetClusterUUID(contaboCluster)
	clusterID := BuildShortClusterID(clusterUUID)
	if err := r.ensureInstanceClusterBinding(ctx, instance, clusterID); err != nil {
		log.Error(err, "failed to bind instance to cluster")
		// Don't fail the whole operation for displayName issues, just log the error
	}

	// Update the machine status
	r.updateMachineStatus(contaboMachine, instance)

	// Set the machine as ready
	contaboMachine.Status.Ready = true
	meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
		Type:   infrastructurev1beta2.InstanceReadyCondition,
		Status: metav1.ConditionTrue,
		Reason: infrastructurev1beta2.InstanceAvailableReason,
	})

	return ctrl.Result{}, nil
}

func (r *ContaboMachineReconciler) reconcileDelete(ctx context.Context, machine *clusterv1.Machine, contaboMachine *infrastructurev1beta2.ContaboMachine, contaboCluster *infrastructurev1beta2.ContaboCluster) (ctrl.Result, error) {
	_ = contaboCluster // may be used in future for cluster-specific cleanup logic
	log := logf.FromContext(ctx)

	log.Info("Reconciling ContaboMachine delete - setting instance available for reuse")

	// Set instance back to available state for reuse
	if contaboMachine.Spec.ProviderID != nil {
		instanceID, err := ParseProviderID(*contaboMachine.Spec.ProviderID)
		if err == nil {
			if err := r.removeClusterBinding(ctx, instanceID); err != nil {
				log.Error(err, "failed to set instance state to available")
				// Don't fail the whole operation for displayName issues
			}
		}
	}

	// Note: We don't actually delete/cancel the instance, just mark it available
	// The instance remains available for reuse with the "<id>-avl" displayName format
	log.Info("Instance state updated to available, instance ready for reuse")

	// Remove our finalizer from the list and update it.
	controllerutil.RemoveFinalizer(contaboMachine, infrastructurev1beta2.MachineFinalizer)

	return ctrl.Result{}, nil
}

func (r *ContaboMachineReconciler) reconcileInstance(ctx context.Context, machine *clusterv1.Machine, contaboMachine *infrastructurev1beta2.ContaboMachine, contaboCluster *infrastructurev1beta2.ContaboCluster, userData string) (*models.InstanceResponse, error) {
	_ = machine        // may be used in future for machine-specific instance configuration
	_ = contaboCluster // may be used in future for cluster-specific instance configuration
	log := logf.FromContext(ctx)

	// If we already have a provider ID, fetch the existing instance
	if contaboMachine.Spec.ProviderID != nil {
		instanceID, err := ParseProviderID(*contaboMachine.Spec.ProviderID)
		if err != nil {
			return nil, fmt.Errorf("failed to parse provider ID: %w", err)
		}

		resp, err := r.ContaboClient.RetrieveInstance(ctx, instanceID, &models.RetrieveInstanceParams{})
		if err != nil {
			return nil, fmt.Errorf("failed to call retrieve instance API: %w", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			log.Info("Instance not found, will find or create a new one")
			contaboMachine.Spec.ProviderID = nil
		} else if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to get existing instance, status: %d", resp.StatusCode)
		} else {
			// Parse the response to get the instance data
			var instanceResp models.FindInstanceResponse
			if err := json.NewDecoder(resp.Body).Decode(&instanceResp); err != nil {
				return nil, fmt.Errorf("failed to decode instance response: %w", err)
			}
			if len(instanceResp.Data) == 0 {
				return nil, fmt.Errorf("no instance data in response")
			}
			log.Info("Found existing instance", "instanceId", instanceResp.Data[0].InstanceId)
			return &instanceResp.Data[0], nil
		}
	}

	// Try to find an available instance for reuse first
	availableInstance, err := r.findAvailableInstance(ctx, contaboMachine)
	if err != nil {
		log.Error(err, "failed to search for available instances")
		// Continue to create new instance
	} else if availableInstance != nil {
		log.Info("Found available instance for reuse", "instanceId", availableInstance.InstanceId)

		// Mark instance as provisioning
		clusterUUID := GetClusterUUID(contaboCluster)
		clusterID := BuildShortClusterID(clusterUUID)
		if err := r.updateInstanceDisplayName(ctx, availableInstance.InstanceId, StateProvisioning, clusterID); err != nil {
			log.Error(err, "failed to mark instance as provisioning, continuing anyway")
		}

		// Reinstall the instance with the correct image and user data
		if err := r.reinstallInstance(ctx, availableInstance.InstanceId, contaboMachine, userData); err != nil {
			log.Error(err, "failed to reinstall available instance, will create new one")
			// Continue to create new instance
		} else {
			// Set the provider ID to the reused instance
			providerID := BuildProviderID(availableInstance.InstanceId)
			contaboMachine.Spec.ProviderID = &providerID

			// Wait for reinstallation and cloud-init to complete
			instance, err := r.waitForInstanceReady(ctx, availableInstance.InstanceId, "reinstalling and configuring")
			if err != nil {
				return nil, fmt.Errorf("failed to wait for reinstalled instance: %w", err)
			}

			log.Info("Successfully reused and reinstalled instance", "instanceId", availableInstance.InstanceId)
			return instance, nil
		}
	}

	// Create a new instance if no available instance found or reuse failed
	log.Info("Creating new instance")
	meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
		Type:   infrastructurev1beta2.InstanceReadyCondition,
		Status: metav1.ConditionFalse,
		Reason: infrastructurev1beta2.InstanceProvisioningReason,
	})

	// Use standardized Ubuntu image for all cluster machines for consistency and security
	defaultImage := DefaultUbuntuImageID

	createReq := &models.CreateInstanceRequest{
		ImageId:     &defaultImage, // Always use Ubuntu 22.04 LTS for cluster nodes
		ProductId:   &contaboMachine.Spec.InstanceType,
		Region:      ConvertRegionToCreateInstanceRegion(contaboMachine.Spec.Region),
		DisplayName: &contaboMachine.Name, // Will be updated after creation with proper format
		UserData:    &userData,
		Period:      1, // Default period
	}

	// Add SSH keys if specified
	if len(contaboMachine.Spec.SSHKeys) > 0 {
		// Convert SSH key names to IDs (this would need to be implemented based on your SSH key management)
		sshKeys := r.convertSSHKeyNamesToIDs(contaboMachine.Spec.SSHKeys)
		createReq.SshKeys = &sshKeys
	}

	// Create instance using OpenAPI client
	createResp, err := r.ContaboClient.CreateInstance(ctx, &models.CreateInstanceParams{}, *createReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call create instance API: %w", err)
	}
	defer func() {
		if closeErr := createResp.Body.Close(); closeErr != nil {
			log.Error(closeErr, "failed to close create instance response body")
		}
	}()

	if createResp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create instance, status: %d", createResp.StatusCode)
	}

	// Parse the response to get the instance data
	var instanceCreateResp models.CreateInstanceResponse
	if err := json.NewDecoder(createResp.Body).Decode(&instanceCreateResp); err != nil {
		return nil, fmt.Errorf("failed to decode create instance response: %w", err)
	}

	if len(instanceCreateResp.Data) == 0 {
		return nil, fmt.Errorf("no instance data in create response")
	}

	createdInstance := &instanceCreateResp.Data[0]

	// Set the provider ID
	providerID := BuildProviderID(createdInstance.InstanceId)
	contaboMachine.Spec.ProviderID = &providerID

	log.Info("Created instance", "instanceId", createdInstance.InstanceId)

	// Mark new instance as provisioning
	clusterUUID := GetClusterUUID(contaboCluster)
	clusterID := BuildShortClusterID(clusterUUID)
	if err := r.updateInstanceDisplayName(ctx, createdInstance.InstanceId, StateProvisioning, clusterID); err != nil {
		log.Error(err, "failed to mark new instance as provisioning, continuing anyway")
	}

	// Wait for the new instance to be ready
	instance, err := r.waitForInstanceReady(ctx, createdInstance.InstanceId, "provisioning")
	if err != nil {
		return nil, fmt.Errorf("failed to wait for new instance: %w", err)
	}

	return instance, nil
}

// findAvailableInstance searches for an available instance that can be reused
func (r *ContaboMachineReconciler) findAvailableInstance(ctx context.Context, contaboMachine *infrastructurev1beta2.ContaboMachine) (*models.InstanceResponse, error) {
	log := logf.FromContext(ctx)

	log.Info("Searching for available instances in region", "region", contaboMachine.Spec.Region, "productType", contaboMachine.Spec.InstanceType)

	// Iterate through all pages to find available instances
	page := int64(0)
	pageSize := int64(100) // Use a reasonable page size

	for {
		// Retrieve list of instances using OpenAPI client with pagination
		resp, err := r.ContaboClient.RetrieveInstancesList(ctx, &models.RetrieveInstancesListParams{
			Page: &page,
			Size: &pageSize,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve instances list (page %d): %w", page, err)
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to retrieve instances list (page %d), status: %d", page, resp.StatusCode)
		}

		// Parse the response
		var instancesResp models.ListInstancesResponse
		if err := json.NewDecoder(resp.Body).Decode(&instancesResp); err != nil {
			return nil, fmt.Errorf("failed to decode instances response (page %d): %w", page, err)
		}

		log.Info("Processing instances page", "page", page, "instanceCount", len(instancesResp.Data))

		// Filter instances to find available ones that match our requirements
		for _, instance := range instancesResp.Data {
			// Check if instance is in the correct region
			if !r.instanceMatchesRegion(&instance, contaboMachine.Spec.Region) {
				continue
			}

			// Check if instance has the correct product type (instance type)
			if !r.instanceMatchesProductType(&instance, contaboMachine.Spec.InstanceType) {
				continue
			}

			// Check if instance is available (based on display name)
			if r.isInstanceAvailable(&instance) {
				log.Info("Found available instance for reuse",
					"instanceId", instance.InstanceId,
					"displayName", instance.DisplayName,
					"region", instance.Region,
					"productId", instance.ProductId,
					"page", page)

				// Convert ListInstancesResponseData to InstanceResponse for return
				instanceResponse := r.convertToInstanceResponse(&instance)
				return instanceResponse, nil
			}
		}

		// Check if we've reached the last page
		totalPages := int64(instancesResp.UnderscorePagination.TotalPages)
		if page >= totalPages || len(instancesResp.Data) == 0 {
			log.Info("Reached last page, no more instances to check", "currentPage", page, "totalPages", totalPages)
			break
		}

		// Move to next page
		page++
		log.Info("Moving to next page", "nextPage", page, "totalPages", totalPages)
	}

	log.Info("No available instances found matching criteria after checking all pages",
		"region", contaboMachine.Spec.Region,
		"instanceType", contaboMachine.Spec.InstanceType,
		"totalPagesChecked", page)
	return nil, nil
}

// reinstallInstance reinstalls an instance with a new image and user data
func (r *ContaboMachineReconciler) reinstallInstance(ctx context.Context, instanceID int64, contaboMachine *infrastructurev1beta2.ContaboMachine, userData string) error {
	log := logf.FromContext(ctx)

	log.Info("Reinstalling instance with Ubuntu image", "instanceId", instanceID, "imageId", DefaultUbuntuImageID)

	// Prepare reinstall request
	reinstallReq := &models.ReinstallInstanceRequest{
		ImageId:  DefaultUbuntuImageID,
		UserData: &userData,
	}

	// Add SSH keys if specified
	if len(contaboMachine.Spec.SSHKeys) > 0 {
		sshKeys := r.convertSSHKeyNamesToIDs(contaboMachine.Spec.SSHKeys)
		reinstallReq.SshKeys = &sshKeys
	}

	// Call reinstall API
	reinstallResp, err := r.ContaboClient.ReinstallInstance(ctx, instanceID, &models.ReinstallInstanceParams{}, *reinstallReq)
	if err != nil {
		return fmt.Errorf("failed to call reinstall instance API: %w", err)
	}
	defer func() {
		if closeErr := reinstallResp.Body.Close(); closeErr != nil {
			log.Error(closeErr, "failed to close reinstall response body")
		}
	}()

	if reinstallResp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to reinstall instance, status: %d", reinstallResp.StatusCode)
	}

	log.Info("Successfully initiated instance reinstall", "instanceId", instanceID)
	return nil
}

// waitForInstanceReady waits for an instance to be in running state and cloud-init to complete
func (r *ContaboMachineReconciler) waitForInstanceReady(ctx context.Context, instanceID int64, operation string) (*models.InstanceResponse, error) {
	log := logf.FromContext(ctx)

	log.Info("Waiting for instance to be ready", "instanceId", instanceID, "operation", operation)

	// Poll instance status until ready
	maxAttempts := 60 // 10 minutes with 10-second intervals
	for attempt := 0; attempt < maxAttempts; attempt++ {
		resp, err := r.ContaboClient.RetrieveInstance(ctx, instanceID, &models.RetrieveInstanceParams{})
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve instance status: %w", err)
		}
		defer func() {
			if closeErr := resp.Body.Close(); closeErr != nil {
				log.Error(closeErr, "failed to close retrieve instance response body")
			}
		}()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to get instance status, status: %d", resp.StatusCode)
		}

		// Parse the response
		var instanceResp models.FindInstanceResponse
		if err := json.NewDecoder(resp.Body).Decode(&instanceResp); err != nil {
			return nil, fmt.Errorf("failed to decode instance response: %w", err)
		}

		if len(instanceResp.Data) == 0 {
			return nil, fmt.Errorf("no instance data in response")
		}

		instance := &instanceResp.Data[0]

		// Check if instance is running
		if instance.Status == models.InstanceStatusRunning {
			log.Info("Instance is ready", "instanceId", instanceID, "operation", operation)
			return instance, nil
		}

		// Log current status
		log.Info("Instance not ready yet", "instanceId", instanceID, "status", instance.Status, "attempt", attempt+1)

		// Wait before next attempt
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context cancelled while waiting for instance: %w", ctx.Err())
		case <-time.After(10 * time.Second):
			// Continue polling
		}
	}

	return nil, fmt.Errorf("timeout waiting for instance to be ready after %d attempts", maxAttempts)
}

// ensureInstanceClusterBinding binds an instance to a cluster via displayName
func (r *ContaboMachineReconciler) ensureInstanceClusterBinding(ctx context.Context, instance *models.InstanceResponse, clusterID string) error {
	log := logf.FromContext(ctx)

	// Set the instance displayName to bind it to the cluster
	log.Info("Binding instance to cluster", "instanceId", instance.InstanceId, "clusterID", clusterID)

	// Update instance display name to indicate successful cluster binding
	return r.updateInstanceDisplayName(ctx, instance.InstanceId, StateClusterBound, clusterID)
}

// removeClusterBinding removes cluster binding and sets instance to available state for reuse
// The instance itself is preserved and made available for reuse rather than being deleted/cancelled
func (r *ContaboMachineReconciler) removeClusterBinding(ctx context.Context, instanceID int64) error {
	log := logf.FromContext(ctx)

	log.Info("Setting instance state to available for reuse", "instanceId", instanceID)

	// Update instance display name to mark as available for reuse
	return r.updateInstanceDisplayName(ctx, instanceID, StateAvailable, "")
}

func (r *ContaboMachineReconciler) getBootstrapData(ctx context.Context, machine *clusterv1.Machine) (string, error) {
	if machine.Spec.Bootstrap.DataSecretName == nil {
		return "", fmt.Errorf("bootstrap data secret name is not set")
	}

	secret := &corev1.Secret{}
	key := client.ObjectKey{Namespace: machine.Namespace, Name: *machine.Spec.Bootstrap.DataSecretName}
	if err := r.Get(ctx, key, secret); err != nil {
		return "", fmt.Errorf("failed to get bootstrap data secret: %w", err)
	}

	userData, exists := secret.Data["value"]
	if !exists {
		return "", fmt.Errorf("bootstrap data secret does not contain 'value' key")
	}

	return base64.StdEncoding.EncodeToString(userData), nil
}

func (r *ContaboMachineReconciler) updateMachineStatus(contaboMachine *infrastructurev1beta2.ContaboMachine, instance *models.InstanceResponse) {
	// Update instance state
	state := infrastructurev1beta2.ContaboMachineInstanceState(instance.Status)
	contaboMachine.Status.InstanceState = &state

	// Update network status
	if contaboMachine.Status.Network == nil {
		contaboMachine.Status.Network = &infrastructurev1beta2.ContaboMachineNetworkStatus{}
	}
	if instance.IpConfig.V4.Ip != "" {
		contaboMachine.Status.Network.PrivateIP = &instance.IpConfig.V4.Ip
	}

	// Update addresses
	contaboMachine.Status.Addresses = []clusterv1.MachineAddress{}
	if instance.IpConfig.V4.Ip != "" {
		contaboMachine.Status.Addresses = append(contaboMachine.Status.Addresses, clusterv1.MachineAddress{
			Type:    clusterv1.MachineInternalIP,
			Address: instance.IpConfig.V4.Ip,
		})
		// For now, assume the same IP is used for external access
		contaboMachine.Status.Addresses = append(contaboMachine.Status.Addresses, clusterv1.MachineAddress{
			Type:    clusterv1.MachineExternalIP,
			Address: instance.IpConfig.V4.Ip,
		})
	}
}

func (r *ContaboMachineReconciler) convertSSHKeyNamesToIDs(keyNames []string) []int64 {
	// This is a placeholder implementation
	// In a real implementation, you would look up SSH keys by name and return their IDs
	var ids []int64
	for _, name := range keyNames {
		// Try to parse as an ID first
		if id, err := strconv.ParseInt(name, 10, 64); err == nil {
			ids = append(ids, id)
		}
	}
	return ids
}

// SetupWithManager sets up the controller with the Manager.
func (r *ContaboMachineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrastructurev1beta2.ContaboMachine{}).
		WithOptions(controller.Options{MaxConcurrentReconciles: 10}).
		WithEventFilter(predicates.ResourceNotPaused(mgr.GetScheme(), ctrl.LoggerFrom(context.TODO()))).
		Watches(
			&clusterv1.Machine{},
			handler.EnqueueRequestsFromMapFunc(util.MachineToInfrastructureMapFunc(infrastructurev1beta2.GroupVersion.WithKind("ContaboMachine"))),
			builder.WithPredicates(predicates.ResourceNotPaused(mgr.GetScheme(), ctrl.LoggerFrom(context.TODO()))),
		).
		Watches(
			&infrastructurev1beta2.ContaboCluster{},
			handler.EnqueueRequestsFromMapFunc(r.ContaboClusterToContaboMachines),
		).
		Named("contabomachine").
		Complete(r)
}

// ContaboClusterToContaboMachines is a handler.ToRequestsFunc to be used to enqueue
// requests for reconciliation of ContaboMachines.
func (r *ContaboMachineReconciler) ContaboClusterToContaboMachines(ctx context.Context, o client.Object) []reconcile.Request {
	result := []reconcile.Request{}

	cluster, ok := o.(*infrastructurev1beta2.ContaboCluster)
	if !ok {
		return result
	}

	log := ctrl.LoggerFrom(ctx).WithValues("objectMapper", "contaboClusterToContaboMachine", "namespace", cluster.Namespace, "contaboCluster", cluster.Name)

	// Don't handle deleted ContaboClusters
	if !cluster.DeletionTimestamp.IsZero() {
		return result
	}

	// clusterName, ok := cluster.Labels[clusterv1.ClusterNameLabel]
	// if !ok {
	// 	log.Info("ContaboCluster does not have cluster label")
	// 	return result
	// }
	clusterName := cluster.Name

	machineList := &clusterv1.MachineList{}
	if err := r.List(ctx, machineList, client.InNamespace(cluster.Namespace), client.MatchingLabels{clusterv1.ClusterNameLabel: clusterName}); err != nil {
		log.Error(err, "Failed to get owned Machines")
		return result
	}

	for _, machine := range machineList.Items {
		if machine.Spec.InfrastructureRef.Name == "" {
			continue
		}
		name := client.ObjectKey{Namespace: machine.Namespace, Name: machine.Spec.InfrastructureRef.Name}
		result = append(result, reconcile.Request{NamespacedName: name})
	}

	return result
}

// instanceMatchesRegion checks if the instance is in the specified region
func (r *ContaboMachineReconciler) instanceMatchesRegion(instance *models.ListInstancesResponseData, region string) bool {
	return instance.Region == region
}

// instanceMatchesProductType checks if the instance has the correct product type (instance type)
func (r *ContaboMachineReconciler) instanceMatchesProductType(instance *models.ListInstancesResponseData, instanceType string) bool {
	// Convert the instance type to match what we expect
	// The instanceType comes from the ContaboMachine spec, and ProductId contains the product identifier
	return instance.ProductId == instanceType
}

// isInstanceAvailable checks if an instance is marked as available based on its display name
func (r *ContaboMachineReconciler) isInstanceAvailable(instance *models.ListInstancesResponseData) bool {
	return GetInstanceState(instance.DisplayName) == StateAvailable
}

// convertToInstanceResponse converts ListInstancesResponseData to InstanceResponse
func (r *ContaboMachineReconciler) convertToInstanceResponse(listData *models.ListInstancesResponseData) *models.InstanceResponse {
	// Convert the product type enum
	var productType models.InstanceResponseProductType
	switch listData.ProductType {
	case models.ListInstancesResponseDataProductTypeHdd:
		productType = models.InstanceResponseProductTypeHdd
	case models.ListInstancesResponseDataProductTypeNvme:
		productType = models.InstanceResponseProductTypeNvme
	case models.ListInstancesResponseDataProductTypeSsd:
		productType = models.InstanceResponseProductTypeSsd
	case models.ListInstancesResponseDataProductTypeVds:
		productType = models.InstanceResponseProductTypeVds
	default:
		productType = models.InstanceResponseProductTypeVds // Default fallback
	}

	// Convert default user type
	var defaultUser *models.InstanceResponseDefaultUser
	if listData.DefaultUser != nil {
		defaultUserValue := models.InstanceResponseDefaultUser(*listData.DefaultUser)
		defaultUser = &defaultUserValue
	}

	// Convert tenant ID type
	tenantId := models.InstanceResponseTenantId(listData.TenantId)

	return &models.InstanceResponse{
		AddOns:        listData.AddOns,
		AdditionalIps: listData.AdditionalIps,
		CancelDate:    listData.CancelDate,
		CpuCores:      listData.CpuCores,
		CreatedDate:   listData.CreatedDate,
		CustomerId:    listData.CustomerId,
		DataCenter:    listData.DataCenter,
		DefaultUser:   defaultUser,
		DiskMb:        listData.DiskMb,
		DisplayName:   listData.DisplayName,
		ErrorMessage:  listData.ErrorMessage,
		ImageId:       listData.ImageId,
		InstanceId:    listData.InstanceId,
		IpConfig:      listData.IpConfig,
		MacAddress:    listData.MacAddress,
		Name:          listData.Name,
		OsType:        listData.OsType,
		ProductId:     listData.ProductId,
		ProductName:   listData.ProductName,
		ProductType:   productType,
		RamMb:         listData.RamMb,
		Region:        listData.Region,
		RegionName:    listData.RegionName,
		SshKeys:       listData.SshKeys,
		Status:        listData.Status,
		TenantId:      tenantId,
		VHostId:       listData.VHostId,
		VHostName:     listData.VHostName,
	}
}

// updateInstanceDisplayName updates the display name of an instance to reflect its current state
func (r *ContaboMachineReconciler) updateInstanceDisplayName(ctx context.Context, instanceID int64, state, clusterID string) error {
	log := logf.FromContext(ctx)

	newDisplayName := BuildInstanceDisplayNameWithState(instanceID, state, clusterID)

	log.Info("Update instance display name",
		"instanceId", instanceID,
		"newDisplayName", newDisplayName,
		"state", state,
		"clusterID", clusterID)

	// Prepare update request
	updateReq := &models.PatchInstanceRequest{
		DisplayName: &newDisplayName,
	}

	// Call update API
	updateResp, err := r.ContaboClient.PatchInstance(ctx, instanceID, &models.PatchInstanceParams{}, *updateReq)
	if err != nil {
		return fmt.Errorf("failed to call update instance API: %w", err)
	}
	defer func() {
		if closeErr := updateResp.Body.Close(); closeErr != nil {
			log.Error(closeErr, "failed to close update instance response body")
		}
	}()

	if updateResp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update instance display name, status: %d", updateResp.StatusCode)
	}

	log.Info("Successfully updated instance display name", "instanceId", instanceID, "newDisplayName", newDisplayName)
	return nil
}
