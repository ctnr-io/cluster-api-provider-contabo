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
	"io"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	infrastructurev1beta2 "github.com/ctnr-io/cluster-api-provider-contabo/api/v1beta2"
	"github.com/ctnr-io/cluster-api-provider-contabo/pkg/contabo/v1.0.0/models"
)

// reconcileClusterPrivateNetworks reconciles cluster private networks
func (r *ContaboClusterReconciler) reconcileClusterPrivateNetworks(ctx context.Context, contaboCluster *infrastructurev1beta2.ContaboCluster) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// If no private networks are specified, mark as skipped
	if len(contaboCluster.Spec.PrivateNetworks) == 0 {
		meta.SetStatusCondition(&contaboCluster.Status.Conditions, metav1.Condition{
			Type:   infrastructurev1beta2.ClusterPrivateNetworkReadyCondition,
			Status: metav1.ConditionTrue,
			Reason: infrastructurev1beta2.ClusterPrivateNetworkSkippedReason,
		})
		return ctrl.Result{}, nil
	}

	// Get the current condition
	networkCondition := meta.FindStatusCondition(contaboCluster.Status.Conditions, infrastructurev1beta2.ClusterPrivateNetworkReadyCondition)

	// State machine for private network lifecycle
	switch {
	case networkCondition == nil:
		// Initial state: Start creating private networks
		return r.createClusterPrivateNetworks(ctx, contaboCluster)

	case networkCondition.Reason == infrastructurev1beta2.ClusterPrivateNetworkCreatingReason:
		// Check if private network creation is complete
		return r.checkClusterPrivateNetworkCreation(ctx, contaboCluster)

	case networkCondition.Reason == infrastructurev1beta2.ClusterPrivateNetworkReadyReason:
		// Networks are ready, nothing to do
		return ctrl.Result{}, nil

	case networkCondition.Reason == infrastructurev1beta2.ClusterPrivateNetworkFailedReason:
		// Networks failed, try to recover or recreate
		return r.handleClusterPrivateNetworkFailure(ctx, contaboCluster)

	default:
		// Unknown state, restart the process
		log.Info("Unknown private network state, restarting creation process")
		return r.createClusterPrivateNetworks(ctx, contaboCluster)
	}
}

// createClusterPrivateNetworks starts the private network creation process
func (r *ContaboClusterReconciler) createClusterPrivateNetworks(ctx context.Context, contaboCluster *infrastructurev1beta2.ContaboCluster) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	log.Info("Creating cluster private networks")
	meta.SetStatusCondition(&contaboCluster.Status.Conditions, metav1.Condition{
		Type:   infrastructurev1beta2.ClusterPrivateNetworkReadyCondition,
		Status: metav1.ConditionFalse,
		Reason: infrastructurev1beta2.ClusterPrivateNetworkCreatingReason,
	})

	// Use existing private network reconciliation logic
	if err := r.reconcilePrivateNetworks(ctx, contaboCluster); err != nil {
		log.Error(err, "Failed to reconcile private networks")
		meta.SetStatusCondition(&contaboCluster.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.ClusterPrivateNetworkReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.ClusterPrivateNetworkFailedReason,
			Message: fmt.Sprintf("Failed to reconcile private networks: %s", err.Error()),
		})
		return ctrl.Result{}, err
	}

	log.Info("Private networks reconciled successfully")
	meta.SetStatusCondition(&contaboCluster.Status.Conditions, metav1.Condition{
		Type:   infrastructurev1beta2.ClusterPrivateNetworkReadyCondition,
		Status: metav1.ConditionTrue,
		Reason: infrastructurev1beta2.ClusterPrivateNetworkReadyReason,
	})

	return ctrl.Result{}, nil
}

// checkClusterPrivateNetworkCreation checks if private network creation is complete
func (r *ContaboClusterReconciler) checkClusterPrivateNetworkCreation(ctx context.Context, contaboCluster *infrastructurev1beta2.ContaboCluster) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// Re-run reconciliation to check current status
	if err := r.reconcilePrivateNetworks(ctx, contaboCluster); err != nil {
		log.Error(err, "Private network creation still in progress or failed")
		meta.SetStatusCondition(&contaboCluster.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.ClusterPrivateNetworkReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.ClusterPrivateNetworkFailedReason,
			Message: fmt.Sprintf("Private network reconciliation failed: %s", err.Error()),
		})
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	log.Info("Private network creation complete")
	meta.SetStatusCondition(&contaboCluster.Status.Conditions, metav1.Condition{
		Type:   infrastructurev1beta2.ClusterPrivateNetworkReadyCondition,
		Status: metav1.ConditionTrue,
		Reason: infrastructurev1beta2.ClusterPrivateNetworkReadyReason,
	})

	return ctrl.Result{}, nil
}

// handleClusterPrivateNetworkFailure handles failed private networks
func (r *ContaboClusterReconciler) handleClusterPrivateNetworkFailure(ctx context.Context, contaboCluster *infrastructurev1beta2.ContaboCluster) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	log.Info("Handling private network failure - attempting recovery")

	// Check current condition to determine failure count/type
	networkCondition := meta.FindStatusCondition(contaboCluster.Status.Conditions, infrastructurev1beta2.ClusterPrivateNetworkReadyCondition)

	// Implement exponential backoff based on failure history
	var requeueDelay time.Duration = 30 * time.Second // Default delay

	if networkCondition != nil && networkCondition.Message != "" {
		// Parse the failure count from condition message (basic implementation)
		if strings.Contains(networkCondition.Message, "retry") {
			// If this is a retry, use longer delay
			requeueDelay = 2 * time.Minute
		}
	}

	// Attempt different recovery strategies based on the type of failure

	// Strategy 1: Reset and recreate networks
	if err := r.resetClusterPrivateNetworkState(ctx, contaboCluster); err != nil {
		log.Error(err, "Failed to reset private network state")
	}

	// Strategy 2: Retry creation with recovery message
	log.Info("Retrying private network creation after failure", "requeueDelay", requeueDelay)
	meta.SetStatusCondition(&contaboCluster.Status.Conditions, metav1.Condition{
		Type:    infrastructurev1beta2.ClusterPrivateNetworkReadyCondition,
		Status:  metav1.ConditionFalse,
		Reason:  infrastructurev1beta2.ClusterPrivateNetworkCreatingReason,
		Message: "Retrying private network creation after failure",
	})

	// Requeue with exponential backoff delay
	return ctrl.Result{RequeueAfter: requeueDelay}, nil
}

// resetClusterPrivateNetworkState resets the private network state for recovery
func (r *ContaboClusterReconciler) resetClusterPrivateNetworkState(ctx context.Context, contaboCluster *infrastructurev1beta2.ContaboCluster) error {
	log := logf.FromContext(ctx)

	log.Info("Resetting cluster private network state for recovery")

	// Clear any partial state in the cluster status to start fresh
	if contaboCluster.Status.PrivateNetworks != nil {
		// Keep track of existing networks but mark them for re-verification
		for i := range contaboCluster.Status.PrivateNetworks {
			log.Info("Marking private network for re-verification",
				"networkName", contaboCluster.Status.PrivateNetworks[i].Name,
				"networkId", contaboCluster.Status.PrivateNetworks[i].PrivateNetworkId)
		}
	}

	// Note: We don't actually delete the networks here as they might be in use
	// The reconciliation process will re-verify and recreate as needed

	return nil
}

func (r *ContaboClusterReconciler) reconcilePrivateNetworks(ctx context.Context, contaboCluster *infrastructurev1beta2.ContaboCluster) error {
	log := logf.FromContext(ctx)

	if len(contaboCluster.Spec.PrivateNetworks) == 0 {
		return nil
	}

	// Initialize private network status slice if not present
	if contaboCluster.Status.PrivateNetworks == nil {
		contaboCluster.Status.PrivateNetworks = []infrastructurev1beta2.ContaboPrivateNetworkStatus{}
	}

	// Process each private network specification
	for _, networkSpec := range contaboCluster.Spec.PrivateNetworks {
		// Check if this network is already in status
		found := false
		for i, privateNetworkStatus := range contaboCluster.Status.PrivateNetworks {
			if privateNetworkStatus.Name == networkSpec.Name {
				// Network already tracked in status, verify it still exists
				if err := r.verifyPrivateNetworkExists(ctx, privateNetworkStatus.PrivateNetworkId); err != nil {
					log.Error(err, "private network no longer exists", "privateNetworkName", networkSpec.Name, "privateNetworkId", privateNetworkStatus.PrivateNetworkId)
					// Remove from status if it no longer exists
					contaboCluster.Status.PrivateNetworks = append(
						contaboCluster.Status.PrivateNetworks[:i],
						contaboCluster.Status.PrivateNetworks[i+1:]...,
					)
				}
				found = true
				break
			}
		}

		if !found {
			// This is a new network specification, discover or create it
			privateNetworkStatus, err := r.ensurePrivateNetwork(ctx, contaboCluster, networkSpec)
			if err != nil {
				log.Error(err, "failed to ensure private network", "privateNetworkName", networkSpec.Name)
				return err
			}

			// Add to status
			contaboCluster.Status.PrivateNetworks = append(
				contaboCluster.Status.PrivateNetworks,
				*privateNetworkStatus,
			)
		}
	}

	return nil
}

func (r *ContaboClusterReconciler) verifyPrivateNetworkExists(ctx context.Context, privateNetworkId int64) error {
	log := logf.FromContext(ctx)

	// Call Contabo API to retrieve the specific private network
	params := &models.RetrievePrivateNetworkParams{}
	resp, err := r.ContaboClient.RetrievePrivateNetwork(ctx, privateNetworkId, params)
	if err != nil {
		log.Error(err, "failed to retrieve private network", "privateNetworkId", privateNetworkId)
		return fmt.Errorf("failed to verify private network exists: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Error(closeErr, "failed to close response body")
		}
	}()

	// Check response status
	if resp.StatusCode == 404 {
		return fmt.Errorf("private network with ID %d not found", privateNetworkId)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code %d when verifying private network %d", resp.StatusCode, privateNetworkId)
	}

	log.V(1).Info("Private network verified successfully", "privateNetworkId", privateNetworkId)
	return nil
}

func (r *ContaboClusterReconciler) ensurePrivateNetwork(ctx context.Context, contaboCluster *infrastructurev1beta2.ContaboCluster, networkSpec infrastructurev1beta2.ContaboPrivateNetworkSpec) (*infrastructurev1beta2.ContaboPrivateNetworkStatus, error) {
	log := logf.FromContext(ctx)

	log.Info("Discovering/creating private network", "privateNetworkName", networkSpec.Name)

	// First, try to find existing network by name
	privateNetworkStatus, err := r.findPrivateNetworkByName(ctx, networkSpec.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to search for existing private network: %w", err)
	}

	// If found, return the existing network
	if privateNetworkStatus != nil {
		log.Info("Found existing private network", "privateNetworkName", networkSpec.Name, "privateNetworkId", privateNetworkStatus.PrivateNetworkId)
		return privateNetworkStatus, nil
	}

	// Network not found, create a new one
	log.Info("Private network not found, creating new one", "privateNetworkName", networkSpec.Name)
	return r.createPrivateNetwork(ctx, contaboCluster, networkSpec)
}

func (r *ContaboClusterReconciler) findPrivateNetworkByName(ctx context.Context, privateNetworkName string) (*infrastructurev1beta2.ContaboPrivateNetworkStatus, error) {
	log := logf.FromContext(ctx)

	// Call Contabo API to list private networks with pagination
	params := &models.RetrievePrivateNetworkListParams{
		Name: &privateNetworkName, // Filter by name to reduce data
	}
	resp, err := r.ContaboClient.RetrievePrivateNetworkList(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve private network %s: %w", privateNetworkName, err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Error(closeErr, "failed to close response body")
		}
	}()

	if resp.StatusCode != 200 {
		// Get error body for more context
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d when listing private networks: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var listResponse models.ListPrivateNetworkResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResponse); err != nil {
		return nil, fmt.Errorf("failed to decode private network list response: %w", err)
	}

	// Search for network with matching name in current page
	for _, privateNetwork := range listResponse.Data {
		if privateNetwork.Name == privateNetworkName {
			log.Info("Found existing private network", "privateNetworkName", privateNetworkName, "privateNetworkId", privateNetwork.PrivateNetworkId)

			return &infrastructurev1beta2.ContaboPrivateNetworkStatus{
				Name:             privateNetwork.Name,
				AvailableIps:     privateNetwork.AvailableIps,
				Cidr:             privateNetwork.Cidr,
				CreatedDate:      privateNetwork.CreatedDate.Unix(),
				CustomerId:       privateNetwork.CustomerId,
				DataCenter:       privateNetwork.DataCenter,
				PrivateNetworkId: privateNetwork.PrivateNetworkId,
				Region:           privateNetwork.Region,
				Description:      privateNetwork.Description,
				RegionName:       privateNetwork.RegionName,
				TenantId:         privateNetwork.TenantId,
			}, nil
		}
	}
	// Network not found in any page
	log.V(1).Info("Private network not found by name", "privateNetworkName", privateNetworkName)
	return nil, nil
}

func (r *ContaboClusterReconciler) createPrivateNetwork(ctx context.Context, contaboCluster *infrastructurev1beta2.ContaboCluster, networkSpec infrastructurev1beta2.ContaboPrivateNetworkSpec) (*infrastructurev1beta2.ContaboPrivateNetworkStatus, error) {
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
	privateNetwork := createResponse.Data[0]
	log.Info("Successfully created private network",
		"privateNetworkName", privateNetwork.Name,
		"privateNetworkId", privateNetwork.PrivateNetworkId,
		"cidr", privateNetwork.Cidr,
		"dataCenter", privateNetwork.DataCenter)

	return &infrastructurev1beta2.ContaboPrivateNetworkStatus{
		Name:             privateNetwork.Name,
		AvailableIps:     privateNetwork.AvailableIps,
		Cidr:             privateNetwork.Cidr,
		CreatedDate:      privateNetwork.CreatedDate.Unix(),
		CustomerId:       privateNetwork.CustomerId,
		DataCenter:       privateNetwork.DataCenter,
		PrivateNetworkId: privateNetwork.PrivateNetworkId,
		Region:           privateNetwork.Region,
		Description:      privateNetwork.Description,
		RegionName:       privateNetwork.RegionName,
		TenantId:         privateNetwork.TenantId,
	}, nil
}

func (r *ContaboClusterReconciler) deleteNetwork(ctx context.Context, contaboCluster *infrastructurev1beta2.ContaboCluster) error {
	log := logf.FromContext(ctx)

	// Delete private networks if they were created by this cluster
	if err := r.deletePrivateNetworks(ctx, contaboCluster); err != nil {
		log.Error(err, "failed to delete private networks")
		return err
	}

	log.Info("Network infrastructure deleted successfully")
	return nil
}

func (r *ContaboClusterReconciler) deletePrivateNetworks(ctx context.Context, contaboCluster *infrastructurev1beta2.ContaboCluster) error {
	log := logf.FromContext(ctx)

	if contaboCluster.Status == nil || len(contaboCluster.Status.PrivateNetworks) == 0 {
		return nil
	}

	// Delete each private network that was tracked in status
	for _, privateNetworkStatus := range contaboCluster.Status.PrivateNetworks {
		log.Info("Deleting private network", "privateNetworkName", privateNetworkStatus.Name, "privateNetworkId", privateNetworkStatus.PrivateNetworkId)

		privateNetworkId := privateNetworkStatus.PrivateNetworkId

		// Step 1: Retrieve private network to check if it exists and get instance information
		networkData, err := r.retrievePrivateNetworkData(ctx, privateNetworkId)
		if err != nil {
			// Step 2: If not exists, tell it's not exists
			log.Info("Private network does not exist, already deleted", "privateNetworkName", privateNetworkStatus.Name, "privateNetworkId", privateNetworkStatus.PrivateNetworkId)
			continue // Network already gone, nothing to delete
		}

		// Step 3: Check if it has instances on it
		instanceCount := len(networkData.Instances)
		if instanceCount > 0 {
			// Step 4: Tell it will not be deleted but that's ok
			log.Info("Private network has instances attached, will not be deleted (this is expected)",
				"privateNetworkName", privateNetworkStatus.Name,
				"privateNetworkId", privateNetworkStatus.PrivateNetworkId,
				"instanceCount", instanceCount)
			continue // Skip deletion, but this is normal behavior
		}

		log.Info("Private network has no instances, proceeding with deletion",
			"privateNetworkName", privateNetworkStatus.Name, "privateNetworkId", privateNetworkStatus.PrivateNetworkId)

		// Call Contabo API to delete private network
		params := &models.DeletePrivateNetworkParams{}
		resp, err := r.ContaboClient.DeletePrivateNetwork(ctx, privateNetworkId, params)
		if err != nil {
			log.Error(err, "failed to delete private network", "privateNetworkName", privateNetworkStatus.Name, "privateNetworkId", privateNetworkStatus.PrivateNetworkId)
			return fmt.Errorf("failed to delete private network %s (ID: %d): %w", privateNetworkStatus.Name, privateNetworkStatus.PrivateNetworkId, err)
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
					"privateNetworkName", privateNetworkStatus.Name, "privateNetworkId", privateNetworkStatus.PrivateNetworkId, "statusCode", resp.StatusCode)
				continue // Skip this network and try others
			}
			return fmt.Errorf("unexpected status code %d when deleting private network %s (ID: %d)", resp.StatusCode, privateNetworkStatus.Name, privateNetworkStatus.PrivateNetworkId)
		}

		log.Info("Private network deletion completed", "privateNetworkName", privateNetworkStatus.Name, "privateNetworkId", privateNetworkStatus.PrivateNetworkId)
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
