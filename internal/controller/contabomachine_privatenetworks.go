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
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	infrastructurev1beta2 "github.com/ctnr-io/cluster-api-provider-contabo/api/v1beta2"
	"github.com/ctnr-io/cluster-api-provider-contabo/pkg/contabo/v1.0.0/models"
)

// reconcileMachinePrivateNetworks reconciles machine private networks
func (r *ContaboMachineReconciler) reconcileMachinePrivateNetworks(ctx context.Context, contaboMachine *infrastructurev1beta2.ContaboMachine, contaboCluster *infrastructurev1beta2.ContaboCluster) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

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

	// Initialize machine private networks status
	contaboMachine.Status.PrivateNetworks = []infrastructurev1beta2.ContaboPrivateNetworkStatus{}

	// If no private networks are specified, mark as skipped
	if len(contaboMachine.Spec.PrivateNetworks) == 0 {
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:   infrastructurev1beta2.MachinePrivateNetworksReadyCondition,
			Status: metav1.ConditionTrue,
			Reason: infrastructurev1beta2.MachinePrivateNetworkSkippedReason,
		})
		return ctrl.Result{}, nil
	}

	// Reconcile each private network specified in the machine spec
	log.Info("Reconciling machine private networks", "networkCount", len(contaboMachine.Spec.PrivateNetworks))

	// Initialize machine private networks status if needed
	if contaboMachine.Status.PrivateNetworks == nil {
		contaboMachine.Status.PrivateNetworks = []infrastructurev1beta2.ContaboPrivateNetworkStatus{}
	}

	// Process each private network specification
	for _, networkSpec := range contaboMachine.Spec.PrivateNetworks {
		privateNetwork, err := r.ensureMachinePrivateNetwork(ctx, contaboMachine, contaboCluster, networkSpec)
		if err != nil {
			return ctrl.Result{}, err
		}
		r.updateMachinePrivateNetworkStatus(contaboMachine, *privateNetwork)
	}

	log.Info("All machine private networks reconciled successfully")
	meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
		Type:   infrastructurev1beta2.MachinePrivateNetworksReadyCondition,
		Status: metav1.ConditionTrue,
		Reason: infrastructurev1beta2.MachinePrivateNetworkReadyReason,
	})

	return ctrl.Result{}, nil
}

// updateMachinePrivateNetworkStatus updates the machine's private network status
func (r *ContaboMachineReconciler) updateMachinePrivateNetworkStatus(contaboMachine *infrastructurev1beta2.ContaboMachine, networkStatus infrastructurev1beta2.ContaboPrivateNetworkStatus) {
	// Check if this network is already in status
	found := false
	for i, existingNetwork := range contaboMachine.Status.PrivateNetworks {
		if existingNetwork.Name == networkStatus.Name {
			contaboMachine.Status.PrivateNetworks[i] = networkStatus
			found = true
			break
		}
	}

	if !found {
		contaboMachine.Status.PrivateNetworks = append(contaboMachine.Status.PrivateNetworks, networkStatus)
	}
}

// ensureMachinePrivateNetwork ensures a private network exists in Contabo, creating it if necessary
func (r *ContaboMachineReconciler) ensureMachinePrivateNetwork(ctx context.Context, contaboMachine *infrastructurev1beta2.ContaboMachine, contaboCluster *infrastructurev1beta2.ContaboCluster, networkSpec infrastructurev1beta2.ContaboPrivateNetworkSpec) (*infrastructurev1beta2.ContaboPrivateNetworkStatus, error) {
	log := logf.FromContext(ctx)

	log.Info("Ensuring machine private network", "networkName", networkSpec.Name)

	existingNetwork := r.findClusterPrivateNetwork(contaboCluster, networkSpec.Name)
	if existingNetwork != nil {
		log.Info("Found existing private network in cluster status", "networkName", networkSpec.Name)
		return existingNetwork, nil
	}

	// If not found in cluster status, try to find it via Contabo API
	networkStatus, err := r.findPrivateNetworkByName(ctx, networkSpec.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to search for existing private network: %w", err)
	}

	// If found via API, return it
	if networkStatus != nil {
		log.Info("Found existing private network via Contabo API", "networkName", networkSpec.Name, "networkId", networkStatus.PrivateNetworkId)
		return networkStatus, nil
	}

	// Private network not found, create a new one
	log.Info("Private network not found, creating new one", "networkName", networkSpec.Name)
	return r.createMachinePrivateNetwork(ctx, contaboMachine, networkSpec)

}

// findClusterPrivateNetwork finds a private network in the cluster status by name
func (r *ContaboMachineReconciler) findClusterPrivateNetwork(contaboCluster *infrastructurev1beta2.ContaboCluster, networkName string) *infrastructurev1beta2.ContaboPrivateNetworkStatus {
	if contaboCluster.Status.PrivateNetworks == nil {
		return nil
	}

	for _, network := range contaboCluster.Status.PrivateNetworks {
		if network.Name == networkName {
			return &network
		}
	}

	return nil
}

// findPrivateNetworkByName searches for a private network by name using the Contabo API
func (r *ContaboMachineReconciler) findPrivateNetworkByName(ctx context.Context, networkName string) (*infrastructurev1beta2.ContaboPrivateNetworkStatus, error) {
	log := logf.FromContext(ctx)

	// Call Contabo API to list private networks with name filter
	params := &models.RetrievePrivateNetworkListParams{
		Name: &networkName, // Filter by name to reduce data
	}
	resp, err := r.ContaboClient.RetrievePrivateNetworkList(ctx, params)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d when listing secrets: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var listResponse models.ListPrivateNetworkResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResponse); err != nil {
		return nil, fmt.Errorf("failed to decode secret list response: %w", err)
	}

	// Search for the network with the exact name
	for _, network := range listResponse.Data {
		if network.Name == networkName {
			return &infrastructurev1beta2.ContaboPrivateNetworkStatus{
				Name:             network.Name,
				Description:      network.Description,
				Region:           network.Region,
				AvailableIps:     network.AvailableIps,
				Cidr:             network.Cidr,
				CreatedDate:      network.CreatedDate.Unix(),
				CustomerId:       network.CustomerId,
				DataCenter:       network.DataCenter,
				PrivateNetworkId: network.PrivateNetworkId,
				RegionName:       network.RegionName,
				TenantId:         network.TenantId,
				Instances:        nil, // Instances are not populated in this context
			}, nil
		}
	}

	log.Info("No existing private network found with name", "networkName", networkName)
	return nil, nil
}

// createMachinePrivateNetwork creates a new private network in Contabo
func (r *ContaboMachineReconciler) createMachinePrivateNetwork(ctx context.Context, contaboMachine *infrastructurev1beta2.ContaboMachine, networkSpec infrastructurev1beta2.ContaboPrivateNetworkSpec) (*infrastructurev1beta2.ContaboPrivateNetworkStatus, error) {
	log := logf.FromContext(ctx)

	log.Info("Creating new private network via Contabo API", "networkName", networkSpec.Name)

	// Prepare the request payload
	createRequest := models.CreatePrivateNetworkRequest{
		Name:        networkSpec.Name,
		Description: networkSpec.Description,
		Region:      networkSpec.Region,
	}

	// Call Contabo API to create the private network
	resp, err := r.ContaboClient.CreatePrivateNetwork(ctx, &models.CreatePrivateNetworkParams{}, createRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to create private network: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Error(closeErr, "failed to close response body")
		}
	}()

	if resp.StatusCode != 201 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d when creating private network: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var networkResponse models.PrivateNetworkResponse
	if err := json.NewDecoder(resp.Body).Decode(&networkResponse); err != nil {
		return nil, fmt.Errorf("failed to decode private network creation response: %w", err)
	}

	log.Info("Successfully created private network", "networkName", networkResponse.Name, "networkId", networkResponse.PrivateNetworkId)

	return &infrastructurev1beta2.ContaboPrivateNetworkStatus{
		Name:             networkResponse.Name,
		Description:      networkResponse.Description,
		Region:           networkResponse.Region,
		AvailableIps:     networkResponse.AvailableIps,
		Cidr:             networkResponse.Cidr,
		CreatedDate:      networkResponse.CreatedDate.Unix(),
		CustomerId:       networkResponse.CustomerId,
		DataCenter:       networkResponse.DataCenter,
		PrivateNetworkId: networkResponse.PrivateNetworkId,
		RegionName:       networkResponse.RegionName,
		TenantId:         networkResponse.TenantId,
		Instances:        nil, // Instances are not populated in this context
	}, nil
}
