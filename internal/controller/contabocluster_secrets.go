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

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	infrastructurev1beta2 "github.com/ctnr-io/cluster-api-provider-contabo/api/v1beta2"
	"github.com/ctnr-io/cluster-api-provider-contabo/pkg/contabo/v1.0.0/models"
)

// reconcileClusterSecrets reconciles cluster secrets
func (r *ContaboClusterReconciler) reconcileClusterSecrets(ctx context.Context, contaboCluster *infrastructurev1beta2.ContaboCluster) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// If no secrets are specified, mark as skipped
	if len(contaboCluster.Spec.Secrets) == 0 {
		log.Info("No cluster secrets specified, skipping")
		meta.SetStatusCondition(&contaboCluster.Status.Conditions, metav1.Condition{
			Type:   infrastructurev1beta2.ClusterSecretsReadyCondition,
			Status: metav1.ConditionTrue,
			Reason: infrastructurev1beta2.ClusterSecretsSkippedReason,
		})
		return ctrl.Result{}, nil
	}

	// Reconcile each secret specified in the cluster spec
	log.Info("Reconciling cluster secrets", "secretCount", len(contaboCluster.Spec.Secrets))

	// Initialize cluster secrets status if needed
	if contaboCluster.Status.Secrets == nil {
		contaboCluster.Status.Secrets = []infrastructurev1beta2.ContaboSecretStatus{}
	}

	// Process each secret specification
	for _, secretSpec := range contaboCluster.Spec.Secrets {
		secretStatus, err := r.ensureClusterSecret(ctx, contaboCluster, secretSpec)
		if err != nil {
			log.Error(err, "Failed to ensure cluster secret", "secretName", secretSpec.Name)
			meta.SetStatusCondition(&contaboCluster.Status.Conditions, metav1.Condition{
				Type:    infrastructurev1beta2.ClusterSecretsReadyCondition,
				Status:  metav1.ConditionFalse,
				Reason:  infrastructurev1beta2.ClusterSecretsFailedReason,
				Message: fmt.Sprintf("Failed to ensure secret %s: %s", secretSpec.Name, err.Error()),
			})
			return ctrl.Result{}, err
		}

		// Update status with the secret info
		r.updateClusterSecretStatus(contaboCluster, *secretStatus)
	}

	log.Info("All cluster secrets reconciled successfully")
	meta.SetStatusCondition(&contaboCluster.Status.Conditions, metav1.Condition{
		Type:   infrastructurev1beta2.ClusterSecretsReadyCondition,
		Status: metav1.ConditionTrue,
		Reason: infrastructurev1beta2.ClusterSecretsReadyReason,
	})

	return ctrl.Result{}, nil
}

// ensureClusterSecret ensures a secret exists in Contabo, creating it if necessary
func (r *ContaboClusterReconciler) ensureClusterSecret(ctx context.Context, contaboCluster *infrastructurev1beta2.ContaboCluster, secretSpec infrastructurev1beta2.ContaboSecretSpec) (*infrastructurev1beta2.ContaboSecretStatus, error) {
	log := logf.FromContext(ctx)

	log.Info("Ensuring cluster secret", "secretName", secretSpec.Name, "secretType", secretSpec.Type)

	// First, try to find existing secret by name
	secretStatus, err := r.findSecretByName(ctx, secretSpec.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to search for existing secret: %w", err)
	}

	// If found, return the existing secret
	if secretStatus != nil {
		log.Info("Found existing secret", "secretName", secretSpec.Name, "secretId", secretStatus.SecretId)
		return secretStatus, nil
	}

	// Secret not found, create a new one
	log.Info("Secret not found, creating new one", "secretName", secretSpec.Name)
	return r.createSecret(ctx, contaboCluster, secretSpec)
}

// findSecretByName searches for a secret by name using the Contabo API
func (r *ContaboClusterReconciler) findSecretByName(ctx context.Context, secretName string) (*infrastructurev1beta2.ContaboSecretStatus, error) {
	log := logf.FromContext(ctx)

	// Call Contabo API to list secrets with name filter
	params := &models.RetrieveSecretListParams{
		Name: &secretName, // Filter by name to reduce data
	}
	resp, err := r.ContaboClient.RetrieveSecretList(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve secret list: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Error(closeErr, "failed to close response body")
		}
	}()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d when listing secrets: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var listResponse models.ListSecretResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResponse); err != nil {
		return nil, fmt.Errorf("failed to decode secret list response: %w", err)
	}

	// Search for secret with matching name
	for _, secret := range listResponse.Data {
		if secret.Name == secretName {
			log.Info("Found existing secret", "secretName", secretName, "secretId", secret.SecretId)

			return &infrastructurev1beta2.ContaboSecretStatus{
				Name:       secret.Name,
				SecretId:   secret.SecretId,
				Type:       infrastructurev1beta2.SecretResponseType(secret.Type),
				CreatedAt:  secret.CreatedAt.UTC().Unix(),
				UpdatedAt:  secret.UpdatedAt.UTC().Unix(),
				CustomerId: secret.CustomerId,
				TenantId:   secret.TenantId,
				Value:      secret.Value,
			}, nil
		}
	}

	// Secret not found
	log.V(1).Info("Secret not found by name", "secretName", secretName)
	return nil, nil
}

// createSecret creates a new secret via the Contabo API
func (r *ContaboClusterReconciler) createSecret(ctx context.Context, contaboCluster *infrastructurev1beta2.ContaboCluster, secretSpec infrastructurev1beta2.ContaboSecretSpec) (*infrastructurev1beta2.ContaboSecretStatus, error) {
	log := logf.FromContext(ctx)

	// Convert API types to Contabo client types
	createRequest := models.CreateSecretRequest{
		Name:  secretSpec.Name,
		Value: secretSpec.Value,
		Type:  models.CreateSecretRequestType(secretSpec.Type),
	}

	// Call Contabo API to create secret
	params := &models.CreateSecretParams{}
	resp, err := r.ContaboClient.CreateSecret(ctx, params, createRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to create secret: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Error(closeErr, "failed to close response body")
		}
	}()

	if resp.StatusCode != 201 && resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d when creating secret: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var createResponse models.CreateSecretResponse
	if err := json.NewDecoder(resp.Body).Decode(&createResponse); err != nil {
		return nil, fmt.Errorf("failed to decode create secret response: %w", err)
	}

	if len(createResponse.Data) == 0 {
		return nil, fmt.Errorf("no data returned from create secret API")
	}

	// Extract the created secret details
	secret := createResponse.Data[0]
	log.Info("Successfully created secret",
		"secretName", secret.Name,
		"secretId", secret.SecretId,
		"secretType", secret.Type)

	return &infrastructurev1beta2.ContaboSecretStatus{
		Name:       secret.Name,
		SecretId:   secret.SecretId,
		Type:       infrastructurev1beta2.SecretResponseType(secret.Type),
		CreatedAt:  secret.CreatedAt.UTC().Unix(),
		UpdatedAt:  secret.UpdatedAt.UTC().Unix(),
		CustomerId: secret.CustomerId,
		TenantId:   secret.TenantId,
		Value:      secret.Value,
	}, nil
}

// updateClusterSecretStatus updates the cluster's secret status
func (r *ContaboClusterReconciler) updateClusterSecretStatus(contaboCluster *infrastructurev1beta2.ContaboCluster, secretStatus infrastructurev1beta2.ContaboSecretStatus) {
	// Check if this secret is already in status
	found := false
	for i, existingStatus := range contaboCluster.Status.Secrets {
		if existingStatus.Name == secretStatus.Name {
			// Update existing status
			contaboCluster.Status.Secrets[i] = secretStatus
			found = true
			break
		}
	}

	if !found {
		// Add new secret status
		contaboCluster.Status.Secrets = append(contaboCluster.Status.Secrets, secretStatus)
	}
}
