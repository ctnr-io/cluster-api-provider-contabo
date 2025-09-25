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

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/conditions"
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

	infrastructurev1beta1 "github.com/ctnr-io/cluster-api-provider-contabo/api/v1beta1"
	contaboclient "github.com/ctnr-io/cluster-api-provider-contabo/pkg/contabo/client"
	"github.com/ctnr-io/cluster-api-provider-contabo/pkg/contabo/models"
)

// Helper functions for provider ID handling and state management

// parseProviderID extracts the instance ID from a provider ID string
func parseProviderID(providerID string) (int64, error) {
	// Provider ID format: "contabo://instanceId"
	const prefix = "contabo://"
	if !strings.HasPrefix(providerID, prefix) {
		return 0, fmt.Errorf("invalid provider ID format: %s", providerID)
	}

	instanceIDStr := strings.TrimPrefix(providerID, prefix)
	instanceID, err := strconv.ParseInt(instanceIDStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid instance ID in provider ID: %w", err)
	}

	return instanceID, nil
}

// buildProviderID creates a provider ID from an instance ID
func buildProviderID(instanceID int64) string {
	return fmt.Sprintf("contabo://%d", instanceID)
}

// mapInstanceStatus maps OpenAPI instance status to our CRD status
func mapInstanceStatus(status models.InstanceStatus) infrastructurev1beta1.ContaboMachineInstanceState {
	switch status {
	case models.InstanceStatusRunning:
		return infrastructurev1beta1.ContaboMachineInstanceStateRunning
	case models.InstanceStatusStopped:
		return infrastructurev1beta1.ContaboMachineInstanceStateStopped
	case models.InstanceStatusInstalling, models.InstanceStatusProvisioning, models.InstanceStatusManualProvisioning:
		return infrastructurev1beta1.ContaboMachineInstanceStatePending
	case models.InstanceStatusError:
		return infrastructurev1beta1.ContaboMachineInstanceStateTerminated
	default:
		return infrastructurev1beta1.ContaboMachineInstanceStateUnknown
	}
}

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
	contaboMachine := &infrastructurev1beta1.ContaboMachine{}
	if err := r.Get(ctx, req.NamespacedName, contaboMachine); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Fetch the Machine
	machine, err := util.GetOwnerMachine(ctx, r.Client, contaboMachine.ObjectMeta)
	if err != nil {
		return ctrl.Result{}, err
	}
	if machine == nil {
		log.Info("Machine Controller has not yet set OwnerRef")
		return ctrl.Result{}, nil
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

	contaboCluster := &infrastructurev1beta1.ContaboCluster{}
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

//nolint:unparam // reconcileNormal may return different ctrl.Result values in future implementations
func (r *ContaboMachineReconciler) reconcileNormal(ctx context.Context, machine *clusterv1.Machine, contaboMachine *infrastructurev1beta1.ContaboMachine, contaboCluster *infrastructurev1beta1.ContaboCluster) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// If the ContaboMachine doesn't have our finalizer, add it.
	controllerutil.AddFinalizer(contaboMachine, infrastructurev1beta1.MachineFinalizer)

	// Check if the infrastructure is ready, otherwise return and wait for the cluster object to be updated
	if !contaboCluster.Status.Ready {
		log.Info("ContaboCluster is not ready yet")
		conditions.MarkFalse(contaboMachine, infrastructurev1beta1.InstanceReadyCondition, infrastructurev1beta1.WaitingForClusterInfrastructureReason, clusterv1.ConditionSeverityInfo, "")
		return ctrl.Result{}, nil
	}

	// Make sure bootstrap data is available and populated
	if machine.Spec.Bootstrap.DataSecretName == nil {
		log.Info("Bootstrap data secret reference is not yet available")
		conditions.MarkFalse(contaboMachine, infrastructurev1beta1.InstanceReadyCondition, infrastructurev1beta1.WaitingForBootstrapDataReason, clusterv1.ConditionSeverityInfo, "")
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
		conditions.MarkFalse(contaboMachine, infrastructurev1beta1.InstanceReadyCondition, infrastructurev1beta1.InstanceProvisioningFailedReason, clusterv1.ConditionSeverityError, "Failed to reconcile instance: %s", err.Error())
		return ctrl.Result{}, err
	}

	// Apply cluster membership via displayName
	clusterName := machine.Spec.ClusterName
	if err := r.ensureInstanceClusterBinding(ctx, instance, clusterName); err != nil {
		log.Error(err, "failed to bind instance to cluster")
		// Don't fail the whole operation for displayName issues, just log the error
	}

	// Update the machine status
	r.updateMachineStatus(contaboMachine, instance)

	// Set the machine as ready
	contaboMachine.Status.Ready = true
	conditions.MarkTrue(contaboMachine, infrastructurev1beta1.InstanceReadyCondition)

	return ctrl.Result{}, nil
}

//nolint:unparam // reconcileDelete may return different ctrl.Result values in future implementations
func (r *ContaboMachineReconciler) reconcileDelete(ctx context.Context, machine *clusterv1.Machine, contaboMachine *infrastructurev1beta1.ContaboMachine, contaboCluster *infrastructurev1beta1.ContaboCluster) (ctrl.Result, error) {
	_ = contaboCluster // may be used in future for cluster-specific cleanup logic
	log := logf.FromContext(ctx)

	log.Info("Reconciling ContaboMachine delete - setting instance available for reuse")

	// Set instance back to available state for reuse
	if contaboMachine.Spec.ProviderID != nil {
		instanceID, err := parseProviderID(*contaboMachine.Spec.ProviderID)
		if err == nil {
			if err := r.removeClusterBinding(ctx, instanceID); err != nil {
				log.Error(err, "failed to set instance state to available")
				// Don't fail the whole operation for displayName issues
			}
		}
	}

	// Note: We don't actually delete/cancel the instance, just mark it available
	// The instance remains available for reuse with the "capc-available" displayName
	log.Info("Instance state updated to available, instance ready for reuse")

	// Remove our finalizer from the list and update it.
	controllerutil.RemoveFinalizer(contaboMachine, infrastructurev1beta1.MachineFinalizer)

	return ctrl.Result{}, nil
}

func (r *ContaboMachineReconciler) reconcileInstance(ctx context.Context, machine *clusterv1.Machine, contaboMachine *infrastructurev1beta1.ContaboMachine, contaboCluster *infrastructurev1beta1.ContaboCluster, userData string) (*models.InstanceResponse, error) {
	_ = machine        // may be used in future for machine-specific instance configuration
	_ = contaboCluster // may be used in future for cluster-specific instance configuration
	log := logf.FromContext(ctx)

	// If we already have a provider ID, fetch the existing instance
	if contaboMachine.Spec.ProviderID != nil {
		instanceID, err := parseProviderID(*contaboMachine.Spec.ProviderID)
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
	conditions.MarkFalse(contaboMachine, infrastructurev1beta1.InstanceReadyCondition, infrastructurev1beta1.InstanceProvisioningReason, clusterv1.ConditionSeverityInfo, "")

	// TODO: Create the proper OpenAPI request
	// Use standardized Ubuntu image for all cluster machines for consistency and security
	defaultImage := DefaultUbuntuImageID
	createReq := &models.CreateInstanceRequest{
		ImageId:     &defaultImage, // Always use Ubuntu 22.04 LTS for cluster nodes
		ProductId:   &contaboMachine.Spec.InstanceType,
		Region:      ConvertRegionToCreateInstanceRegion(contaboMachine.Spec.Region),
		DisplayName: &contaboMachine.Name,
		UserData:    &userData,
		Period:      1, // Default period
	}

	// Add SSH keys if specified
	if len(contaboMachine.Spec.SSHKeys) > 0 {
		// Convert SSH key names to IDs (this would need to be implemented based on your SSH key management)
		sshKeys := r.convertSSHKeyNamesToIDs(contaboMachine.Spec.SSHKeys)
		createReq.SshKeys = &sshKeys
	}

	// TODO: Implement using OpenAPI client CreateInstance
	// createResp, err := r.ContaboClient.CreateInstance(ctx, &models.CreateInstanceParams{}, models.CreateInstanceJSONRequestBody(*createReq))
	log.Info("Create instance - not yet implemented with OpenAPI client")

	// For now, return a stub instance to get things compiling
	stubInstance := &models.InstanceResponse{
		InstanceId: 12345, // stub value
		Status:     models.InstanceStatusRunning,
		IpConfig: models.IpConfig{
			V4: models.IpV4{
				Ip: "192.168.1.100", // stub value
			},
		},
	}

	// Set the provider ID
	providerID := BuildProviderID(stubInstance.InstanceId)
	contaboMachine.Spec.ProviderID = &providerID

	log.Info("Created instance (stub)", "instanceId", stubInstance.InstanceId)

	// Wait for the new instance to be ready
	instance, err := r.waitForInstanceReady(ctx, stubInstance.InstanceId, "provisioning")
	if err != nil {
		return nil, fmt.Errorf("failed to wait for new instance: %w", err)
	}

	return instance, nil
}

// findAvailableInstance searches for an available instance that can be reused
func (r *ContaboMachineReconciler) findAvailableInstance(ctx context.Context, contaboMachine *infrastructurev1beta1.ContaboMachine) (*models.InstanceResponse, error) {
	log := logf.FromContext(ctx)

	log.Info("Searching for available instances in region", "region", contaboMachine.Spec.Region, "productType", contaboMachine.Spec.InstanceType)

	// TODO: Implement using OpenAPI client RetrieveInstancesList
	// This would involve calling the API to list instances and filtering for available ones
	log.Info("Find available instance - not yet implemented with OpenAPI client")

	// For now, return nil to indicate no available instance found
	return nil, nil
}

// reinstallInstance reinstalls an instance with a new image and user data
func (r *ContaboMachineReconciler) reinstallInstance(ctx context.Context, instanceID int64, contaboMachine *infrastructurev1beta1.ContaboMachine, userData string) error {
	log := logf.FromContext(ctx)

	log.Info("Reinstalling instance with Ubuntu image", "instanceId", instanceID, "imageId", DefaultUbuntuImageID)

	// TODO: Implement using OpenAPI client ReinstallInstance
	// This would involve calling the ReinstallInstance API endpoint
	log.Info("Reinstall instance - not yet implemented with OpenAPI client")

	return nil
}

// waitForInstanceReady waits for an instance to be in running state and cloud-init to complete
func (r *ContaboMachineReconciler) waitForInstanceReady(ctx context.Context, instanceID int64, operation string) (*models.InstanceResponse, error) {
	log := logf.FromContext(ctx)

	log.Info("Waiting for instance to be ready", "instanceId", instanceID, "operation", operation)

	// TODO: Implement using OpenAPI client RetrieveInstance
	// This would involve polling the RetrieveInstance API until status is ready
	log.Info("Wait for instance ready - not yet implemented with OpenAPI client")

	// For now, return a stub ready instance
	return &models.InstanceResponse{
		InstanceId: instanceID,
		Status:     models.InstanceStatusRunning,
		IpConfig: models.IpConfig{
			V4: models.IpV4{
				Ip: "192.168.1.100", // stub value
			},
		},
	}, nil
}

// ensureInstanceClusterBinding binds an instance to a cluster via displayName
func (r *ContaboMachineReconciler) ensureInstanceClusterBinding(ctx context.Context, instance *models.InstanceResponse, clusterName string) error {
	log := logf.FromContext(ctx)

	// Set the instance displayName to bind it to the cluster
	log.Info("Binding instance to cluster", "instanceId", instance.InstanceId, "clusterName", clusterName)

	// TODO: Implement using OpenAPI client - update instance displayName to indicate cluster binding
	// This would involve calling UpdateInstance API to change the displayName to include cluster info
	log.Info("Instance cluster binding - not yet implemented with OpenAPI client")
	return nil
}

// removeClusterBinding removes cluster binding and sets instance to available state for reuse
// The instance itself is preserved and made available for reuse rather than being deleted/cancelled
func (r *ContaboMachineReconciler) removeClusterBinding(ctx context.Context, instanceID int64) error {
	log := logf.FromContext(ctx)

	log.Info("Setting instance state to available for reuse", "instanceId", instanceID)

	// TODO: Implement using OpenAPI client - update instance displayName to mark as available
	// This would involve calling UpdateInstance API to change the displayName to indicate availability
	log.Info("Remove cluster binding - not yet implemented with OpenAPI client")
	return nil
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

func (r *ContaboMachineReconciler) updateMachineStatus(contaboMachine *infrastructurev1beta1.ContaboMachine, instance *models.InstanceResponse) {
	// Update instance state
	state := infrastructurev1beta1.ContaboMachineInstanceState(instance.Status)
	contaboMachine.Status.InstanceState = &state

	// Update network status
	if contaboMachine.Status.Network == nil {
		contaboMachine.Status.Network = &infrastructurev1beta1.ContaboMachineNetworkStatus{}
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
		For(&infrastructurev1beta1.ContaboMachine{}).
		WithOptions(controller.Options{MaxConcurrentReconciles: 10}).
		WithEventFilter(predicates.ResourceNotPaused(ctrl.LoggerFrom(context.TODO()))).
		Watches(
			&clusterv1.Machine{},
			handler.EnqueueRequestsFromMapFunc(util.MachineToInfrastructureMapFunc(infrastructurev1beta1.GroupVersion.WithKind("ContaboMachine"))),
			builder.WithPredicates(predicates.ResourceNotPaused(ctrl.LoggerFrom(context.TODO()))),
		).
		Watches(
			&infrastructurev1beta1.ContaboCluster{},
			handler.EnqueueRequestsFromMapFunc(r.ContaboClusterToContaboMachines),
		).
		Named("contabomachine").
		Complete(r)
}

// ContaboClusterToContaboMachines is a handler.ToRequestsFunc to be used to enqueue
// requests for reconciliation of ContaboMachines.
func (r *ContaboMachineReconciler) ContaboClusterToContaboMachines(ctx context.Context, o client.Object) []reconcile.Request {
	result := []reconcile.Request{}

	cluster, ok := o.(*infrastructurev1beta1.ContaboCluster)
	if !ok {
		return result
	}

	log := ctrl.LoggerFrom(ctx).WithValues("objectMapper", "contaboClusterToContaboMachine", "namespace", cluster.Namespace, "contaboCluster", cluster.Name)

	// Don't handle deleted ContaboClusters
	if !cluster.DeletionTimestamp.IsZero() {
		return result
	}

	clusterName, ok := cluster.Labels[clusterv1.ClusterNameLabel]
	if !ok {
		log.Info("ContaboCluster does not have cluster label")
		return result
	}

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
