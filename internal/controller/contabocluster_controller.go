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

package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

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

	infrastructurev1beta1 "github.com/ctnr-io/cluster-api-provider-contabo/api/v1beta1"
	contaboclient "github.com/ctnr-io/cluster-api-provider-contabo/pkg/contabo/client"
	"github.com/ctnr-io/cluster-api-provider-contabo/pkg/contabo/models"
)

// ContaboClusterReconciler reconciles a ContaboCluster object
type ContaboClusterReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	Recorder      record.EventRecorder
	ContaboClient *contaboclient.Client
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=contaboclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=contaboclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=contaboclusters/finalizers,verbs=update
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters;clusters/status,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *ContaboClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// Fetch the ContaboCluster instance
	contaboCluster := &infrastructurev1beta1.ContaboCluster{}
	if err := r.Get(ctx, req.NamespacedName, contaboCluster); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Fetch the Cluster
	cluster, err := util.GetOwnerCluster(ctx, r.Client, contaboCluster.ObjectMeta)
	if err != nil {
		return ctrl.Result{}, err
	}
	if cluster == nil {
		log.Info("Cluster Controller has not yet set OwnerRef, requeuing",
			"contaboCluster", contaboCluster.Name,
			"namespace", contaboCluster.Namespace,
			"ownerReferences", contaboCluster.OwnerReferences)

		// Requeue after 10 seconds to allow Cluster controller to set OwnerRef
		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	}

	if annotations.IsPaused(cluster, contaboCluster) {
		log.Info("ContaboCluster or linked Cluster is marked as paused. Won't reconcile")
		return ctrl.Result{}, nil
	}

	log = log.WithValues("cluster", cluster.Name)

	// Initialize the patch helper
	patchHelper, err := patch.NewHelper(contaboCluster, r.Client)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Always attempt to Patch the ContaboCluster object and status after each reconciliation
	defer func() {
		if err := patchHelper.Patch(ctx, contaboCluster); err != nil {
			log.Error(err, "failed to patch ContaboCluster")
		}
	}()

	// Handle deleted clusters
	if !contaboCluster.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, cluster, contaboCluster)
	}

	// Handle non-deleted clusters
	return r.reconcileNormal(ctx, contaboCluster)
}

func (r *ContaboClusterReconciler) reconcileNormal(ctx context.Context, contaboCluster *infrastructurev1beta1.ContaboCluster) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// If the ContaboCluster doesn't have our finalizer, add it.
	controllerutil.AddFinalizer(contaboCluster, infrastructurev1beta1.ClusterFinalizer)

	// Always ensure status is initialized
	if contaboCluster.Status.Conditions == nil {
		contaboCluster.Status.Conditions = []clusterv1.Condition{}
	}

	// Ensure cluster has a unique UUID for global identification
	clusterUUID := EnsureClusterUUID(contaboCluster)
	log.Info("Cluster UUID ensured", "uuid", clusterUUID)

	// Set the cluster in a progressing state
	conditions.MarkFalse(contaboCluster, infrastructurev1beta1.ReadyCondition, infrastructurev1beta1.CreatingReason, clusterv1.ConditionSeverityInfo, "")

	log.Info("Starting ContaboCluster reconciliation",
		"cluster", contaboCluster.Name,
		"namespace", contaboCluster.Namespace,
		"region", contaboCluster.Spec.Region)

	// Reconcile network infrastructure
	if err := r.reconcileNetwork(ctx, contaboCluster); err != nil {
		log.Error(err, "failed to reconcile network")
		conditions.MarkFalse(contaboCluster, infrastructurev1beta1.ReadyCondition, infrastructurev1beta1.NetworkInfrastructureFailedReason, clusterv1.ConditionSeverityError, "Failed to reconcile network: %s", err.Error())
		return ctrl.Result{}, err
	}

	// Set up the control plane endpoint if not already set
	if contaboCluster.Spec.ControlPlaneEndpoint.Host == "" {
		log.Info("Control plane endpoint not set, will be set by the first control plane machine")
		// The control plane endpoint will be set by the first control plane machine
		// For now, we'll mark the cluster as not ready and return
		conditions.MarkFalse(contaboCluster, infrastructurev1beta1.ReadyCondition, infrastructurev1beta1.WaitingForControlPlaneEndpointReason, clusterv1.ConditionSeverityInfo, "")
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	// Mark the cluster as ready
	contaboCluster.Status.Ready = true
	conditions.MarkTrue(contaboCluster, infrastructurev1beta1.ReadyCondition)

	return ctrl.Result{}, nil
}

//nolint:unparam // reconcileDelete may return different ctrl.Result values in future implementations
func (r *ContaboClusterReconciler) reconcileDelete(ctx context.Context, cluster *clusterv1.Cluster, contaboCluster *infrastructurev1beta1.ContaboCluster) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	log.Info("Reconciling ContaboCluster delete")

	// Delete network infrastructure
	if err := r.deleteNetwork(ctx, contaboCluster); err != nil {
		log.Error(err, "failed to delete network")
		return ctrl.Result{}, err
	}

	// Remove our finalizer from the list and update it.
	controllerutil.RemoveFinalizer(contaboCluster, infrastructurev1beta1.ClusterFinalizer)

	return ctrl.Result{}, nil
}

func (r *ContaboClusterReconciler) reconcileNetwork(ctx context.Context, contaboCluster *infrastructurev1beta1.ContaboCluster) error {
	log := logf.FromContext(ctx)

	// Initialize network status if not present
	if contaboCluster.Status.Network == nil {
		contaboCluster.Status.Network = &infrastructurev1beta1.ContaboNetworkStatus{}
	}

	// If no network spec is provided, skip network reconciliation
	if contaboCluster.Spec.Network == nil {
		log.Info("No network specification provided, skipping network reconciliation")
		return nil
	}

	// Reconcile private networks
	if err := r.reconcilePrivateNetworks(ctx, contaboCluster); err != nil {
		log.Error(err, "failed to reconcile private networks")
		return err
	}

	log.Info("Network infrastructure reconciled successfully")
	return nil
}

func (r *ContaboClusterReconciler) reconcilePrivateNetworks(ctx context.Context, contaboCluster *infrastructurev1beta1.ContaboCluster) error {
	log := logf.FromContext(ctx)

	if contaboCluster.Spec.Network == nil || len(contaboCluster.Spec.Network.PrivateNetworks) == 0 {
		return nil
	}

	// Initialize private network status slice if not present
	if contaboCluster.Status.Network.PrivateNetworks == nil {
		contaboCluster.Status.Network.PrivateNetworks = []infrastructurev1beta1.ContaboPrivateNetworkStatus{}
	}

	// Process each private network specification
	for _, networkSpec := range contaboCluster.Spec.Network.PrivateNetworks {
		// Check if this network is already in status
		found := false
		for i, networkStatus := range contaboCluster.Status.Network.PrivateNetworks {
			if networkStatus.Name == networkSpec.Name {
				// Network already tracked in status, verify it still exists
				if err := r.verifyPrivateNetworkExists(ctx, networkStatus.ID); err != nil {
					log.Error(err, "private network no longer exists", "networkName", networkSpec.Name, "networkID", networkStatus.ID)
					// Remove from status if it no longer exists
					contaboCluster.Status.Network.PrivateNetworks = append(
						contaboCluster.Status.Network.PrivateNetworks[:i],
						contaboCluster.Status.Network.PrivateNetworks[i+1:]...,
					)
				}
				found = true
				break
			}
		}

		if !found {
			// This is a new network specification, discover or create it
			networkStatus, err := r.ensurePrivateNetwork(ctx, contaboCluster, networkSpec)
			if err != nil {
				log.Error(err, "failed to ensure private network", "networkName", networkSpec.Name)
				return err
			}

			// Add to status
			contaboCluster.Status.Network.PrivateNetworks = append(
				contaboCluster.Status.Network.PrivateNetworks,
				*networkStatus,
			)
		}
	}

	return nil
}

func (r *ContaboClusterReconciler) verifyPrivateNetworkExists(ctx context.Context, networkID string) error {
	log := logf.FromContext(ctx)

	// Convert networkID string to int64 as required by the API
	privateNetworkId, err := strconv.ParseInt(networkID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid network ID format: %w", err)
	}

	// Call Contabo API to retrieve the specific private network
	params := &models.RetrievePrivateNetworkParams{}
	resp, err := r.ContaboClient.RetrievePrivateNetwork(ctx, privateNetworkId, params)
	if err != nil {
		log.Error(err, "failed to retrieve private network", "networkID", networkID)
		return fmt.Errorf("failed to verify private network exists: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Error(closeErr, "failed to close response body")
		}
	}()

	// Check response status
	if resp.StatusCode == 404 {
		return fmt.Errorf("private network with ID %s not found", networkID)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code %d when verifying private network %s", resp.StatusCode, networkID)
	}

	log.V(1).Info("Private network verified successfully", "networkID", networkID)
	return nil
}

func (r *ContaboClusterReconciler) ensurePrivateNetwork(ctx context.Context, contaboCluster *infrastructurev1beta1.ContaboCluster, networkSpec infrastructurev1beta1.ContaboPrivateNetworkSpec) (*infrastructurev1beta1.ContaboPrivateNetworkStatus, error) {
	log := logf.FromContext(ctx)

	log.Info("Discovering/creating private network", "networkName", networkSpec.Name)

	// First, try to find existing network by name
	networkStatus, err := r.findPrivateNetworkByName(ctx, networkSpec.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to search for existing private network: %w", err)
	}

	// If found, return the existing network
	if networkStatus != nil {
		log.Info("Found existing private network", "networkName", networkSpec.Name, "networkID", networkStatus.ID)
		return networkStatus, nil
	}

	// Network not found, create a new one
	log.Info("Private network not found, creating new one", "networkName", networkSpec.Name)
	return r.createPrivateNetwork(ctx, contaboCluster, networkSpec)
}

func (r *ContaboClusterReconciler) findPrivateNetworkByName(ctx context.Context, networkName string) (*infrastructurev1beta1.ContaboPrivateNetworkStatus, error) {
	log := logf.FromContext(ctx)

	// Iterate through all pages to find the network
	page := int64(1)
	size := int64(100) // Maximum page size

	for {
		// Call Contabo API to list private networks with pagination
		params := &models.RetrievePrivateNetworkListParams{
			Page: &page,
			Size: &size,
		}
		resp, err := r.ContaboClient.RetrievePrivateNetworkList(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve private network list (page %d): %w", page, err)
		}
		defer func() {
			if closeErr := resp.Body.Close(); closeErr != nil {
				log.Error(closeErr, "failed to close response body")
			}
		}()

		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("unexpected status code %d when listing private networks (page %d)", resp.StatusCode, page)
		}

		// Parse the response
		var listResponse models.ListPrivateNetworkResponse
		if err := json.NewDecoder(resp.Body).Decode(&listResponse); err != nil {
			return nil, fmt.Errorf("failed to decode private network list response (page %d): %w", page, err)
		}

		// Search for network with matching name in current page
		for _, network := range listResponse.Data {
			if network.Name == networkName {
				log.Info("Found existing private network", "networkName", networkName, "networkID", network.PrivateNetworkId)

				return &infrastructurev1beta1.ContaboPrivateNetworkStatus{
					Name:       network.Name,
					ID:         strconv.FormatInt(network.PrivateNetworkId, 10),
					CIDR:       network.Cidr,
					DataCenter: network.DataCenter,
					Region:     network.Region,
				}, nil
			}
		}

		// Check if there are more pages
		if len(listResponse.Data) < int(size) {
			// This was the last page
			break
		}

		// Move to next page
		page++
	}

	// Network not found in any page
	log.V(1).Info("Private network not found by name", "networkName", networkName)
	return nil, nil
}

func (r *ContaboClusterReconciler) createPrivateNetwork(ctx context.Context, contaboCluster *infrastructurev1beta1.ContaboCluster, networkSpec infrastructurev1beta1.ContaboPrivateNetworkSpec) (*infrastructurev1beta1.ContaboPrivateNetworkStatus, error) {
	log := logf.FromContext(ctx)

	// Prepare create request
	createRequest := models.CreatePrivateNetworkRequest{
		Name: networkSpec.Name,
	}

	// Set optional description if provided in spec (for future extension)
	description := fmt.Sprintf("Private network created by cluster-api-provider-contabo for network: %s", networkSpec.Name)
	createRequest.Description = &description

	// Use the cluster's region
	createRequest.Region = &contaboCluster.Spec.Region

	// Call Contabo API to create private network
	params := &models.CreatePrivateNetworkParams{}
	resp, err := r.ContaboClient.CreatePrivateNetwork(ctx, params, createRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to create private network: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Error(closeErr, "failed to close response body")
		}
	}()

	if resp.StatusCode != 201 && resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code %d when creating private network", resp.StatusCode)
	}

	// Parse the response
	var createResponse models.CreatePrivateNetworkResponse
	if err := json.NewDecoder(resp.Body).Decode(&createResponse); err != nil {
		return nil, fmt.Errorf("failed to decode create private network response: %w", err)
	}

	if len(createResponse.Data) == 0 {
		return nil, fmt.Errorf("no data returned from create private network API")
	}

	// Extract the created network details
	createdNetwork := createResponse.Data[0]
	log.Info("Successfully created private network",
		"networkName", createdNetwork.Name,
		"networkID", createdNetwork.PrivateNetworkId,
		"cidr", createdNetwork.Cidr,
		"dataCenter", createdNetwork.DataCenter)

	return &infrastructurev1beta1.ContaboPrivateNetworkStatus{
		Name:       createdNetwork.Name,
		ID:         strconv.FormatInt(createdNetwork.PrivateNetworkId, 10),
		CIDR:       createdNetwork.Cidr,
		DataCenter: createdNetwork.DataCenter,
		Region:     createdNetwork.Region,
	}, nil
}

func (r *ContaboClusterReconciler) deleteNetwork(ctx context.Context, contaboCluster *infrastructurev1beta1.ContaboCluster) error {
	log := logf.FromContext(ctx)

	// Delete private networks if they were created by this cluster
	if err := r.deletePrivateNetworks(ctx, contaboCluster); err != nil {
		log.Error(err, "failed to delete private networks")
		return err
	}

	log.Info("Network infrastructure deleted successfully")
	return nil
}

func (r *ContaboClusterReconciler) deletePrivateNetworks(ctx context.Context, contaboCluster *infrastructurev1beta1.ContaboCluster) error {
	log := logf.FromContext(ctx)

	if contaboCluster.Status.Network == nil || len(contaboCluster.Status.Network.PrivateNetworks) == 0 {
		return nil
	}

	// Delete each private network that was tracked in status
	for _, networkStatus := range contaboCluster.Status.Network.PrivateNetworks {
		log.Info("Deleting private network", "networkName", networkStatus.Name, "networkID", networkStatus.ID)

		// Convert networkID string to int64 as required by the API
		privateNetworkId, err := strconv.ParseInt(networkStatus.ID, 10, 64)
		if err != nil {
			log.Error(err, "invalid network ID format during deletion", "networkID", networkStatus.ID)
			continue // Skip this network and try others
		}

		// Step 1: Retrieve private network to check if it exists and get instance information
		networkData, err := r.retrievePrivateNetworkData(ctx, privateNetworkId)
		if err != nil {
			// Step 2: If not exists, tell it's not exists
			log.Info("Private network does not exist, already deleted", "networkName", networkStatus.Name, "networkID", networkStatus.ID)
			continue // Network already gone, nothing to delete
		}

		// Step 3: Check if it has instances on it
		instanceCount := len(networkData.Instances)
		if instanceCount > 0 {
			// Step 4: Tell it will not be deleted but that's ok
			log.Info("Private network has instances attached, will not be deleted (this is expected)",
				"networkName", networkStatus.Name,
				"networkID", networkStatus.ID,
				"instanceCount", instanceCount)
			continue // Skip deletion, but this is normal behavior
		}

		log.Info("Private network has no instances, proceeding with deletion",
			"networkName", networkStatus.Name, "networkID", networkStatus.ID)

		// Call Contabo API to delete private network
		params := &models.DeletePrivateNetworkParams{}
		resp, err := r.ContaboClient.DeletePrivateNetwork(ctx, privateNetworkId, params)
		if err != nil {
			log.Error(err, "failed to delete private network", "networkName", networkStatus.Name, "networkID", networkStatus.ID)
			return fmt.Errorf("failed to delete private network %s (ID: %s): %w", networkStatus.Name, networkStatus.ID, err)
		}
		defer func() {
			if closeErr := resp.Body.Close(); closeErr != nil {
				log.Error(closeErr, "failed to close response body")
			}
		}()

		if resp.StatusCode != 204 && resp.StatusCode != 200 && resp.StatusCode != 404 {
			// Check if error is due to instances still being attached
			if resp.StatusCode == 422 || resp.StatusCode == 400 {
				log.Info("Private network cannot be deleted - likely has instances still attached",
					"networkName", networkStatus.Name, "networkID", networkStatus.ID, "statusCode", resp.StatusCode)
				continue // Skip this network and try others
			}
			return fmt.Errorf("unexpected status code %d when deleting private network %s (ID: %s)", resp.StatusCode, networkStatus.Name, networkStatus.ID)
		}

		log.Info("Private network deletion completed", "networkName", networkStatus.Name, "networkID", networkStatus.ID)
	}

	return nil
}

func (r *ContaboClusterReconciler) retrievePrivateNetworkData(ctx context.Context, privateNetworkId int64) (*models.PrivateNetworkResponse, error) {
	log := logf.FromContext(ctx)

	// Step 1: Retrieve private network
	params := &models.RetrievePrivateNetworkParams{}
	resp, err := r.ContaboClient.RetrievePrivateNetwork(ctx, privateNetworkId, params)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve private network: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Error(closeErr, "failed to close response body")
		}
	}()

	// Step 2: If not exists (404), return error
	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("private network not found")
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code %d when retrieving private network", resp.StatusCode)
	}

	// Parse the response
	var networkResponse models.FindPrivateNetworkResponse
	if err := json.NewDecoder(resp.Body).Decode(&networkResponse); err != nil {
		return nil, fmt.Errorf("failed to decode private network response: %w", err)
	}

	if len(networkResponse.Data) == 0 {
		return nil, fmt.Errorf("no private network data found in response")
	}

	log.V(1).Info("Successfully retrieved private network data",
		"privateNetworkId", privateNetworkId,
		"instanceCount", len(networkResponse.Data[0].Instances))

	return &networkResponse.Data[0], nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ContaboClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrastructurev1beta1.ContaboCluster{}).
		WithOptions(controller.Options{MaxConcurrentReconciles: 10}).
		WithEventFilter(predicates.ResourceNotPaused(ctrl.LoggerFrom(context.TODO()))).
		Watches(
			&clusterv1.Cluster{},
			handler.EnqueueRequestsFromMapFunc(util.ClusterToInfrastructureMapFunc(context.TODO(), infrastructurev1beta1.GroupVersion.WithKind("ContaboCluster"), mgr.GetClient(), &infrastructurev1beta1.ContaboCluster{})),
			builder.WithPredicates(predicates.ClusterUnpaused(ctrl.LoggerFrom(context.TODO()))),
		).
		Named("contabocluster").
		Complete(r)
}
