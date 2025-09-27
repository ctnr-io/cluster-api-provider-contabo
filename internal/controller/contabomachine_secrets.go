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

// reconcileMachineSecrets reconciles secrets for the machine
func (r *ContaboMachineReconciler) reconcileMachineSecrets(ctx context.Context, contaboMachine *infrastructurev1beta2.ContaboMachine, contaboCluster *infrastructurev1beta2.ContaboCluster) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// Initialize machine secrets status
	contaboMachine.Status.Secrets = []infrastructurev1beta2.ContaboSecretStatus{}

	// Merge cluster and machine secrets
	allSecrets := []infrastructurev1beta2.ContaboSecretSpec{}
	
	// Add cluster-wide secrets first
	if len(contaboCluster.Spec.Secrets) > 0 {
		allSecrets = append(allSecrets, contaboCluster.Spec.Secrets...)
	}

	// Add machine-specific secrets
	if len(contaboMachine.Spec.Secrets) > 0 {
		allSecrets = append(allSecrets, contaboMachine.Spec.Secrets...)
	}

	// If no secrets are specified, mark as skipped
	if len(allSecrets) == 0 {

		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:   infrastructurev1beta2.MachineSecretsReadyCondition,
			Status: metav1.ConditionTrue,
			Reason: infrastructurev1beta2.MachineSecretsSkippedReason,
		})
		return ctrl.Result{}, nil
	}

	// Reconcile each secret specified in the machine spec
	log.Info("Reconciling machine secrets", "secretCount", len(allSecrets))

	// Process each secret specification
	for _, secretSpec := range allSecrets {
		secretStatus, err := r.ensureMachineSecret(ctx, contaboMachine, secretSpec)
		if err != nil {
			log.Error(err, "Failed to ensure machine secret", "secretName", secretSpec.Name)
			meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
				Type:    infrastructurev1beta2.MachineSecretsReadyCondition,
				Status:  metav1.ConditionFalse,
				Reason:  infrastructurev1beta2.MachineSecretsFailedReason,
				Message: fmt.Sprintf("Failed to ensure secret %s: %s", secretSpec.Name, err.Error()),
			})
			return ctrl.Result{}, err
		}

		// Update status with the secret info
		contaboMachine.Status.Secrets = append(contaboMachine.Status.Secrets, *secretStatus)
	}

	log.Info("All machine secrets reconciled successfully")
	meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
		Type:   infrastructurev1beta2.MachineSecretsReadyCondition,
		Status: metav1.ConditionTrue,
		Reason: infrastructurev1beta2.MachineSecretsReadyReason,
	})

	return ctrl.Result{}, nil
}

// ensureMachineSecret ensures a secret exists in Contabo for the machine
func (r *ContaboMachineReconciler) ensureMachineSecret(ctx context.Context, contaboMachine *infrastructurev1beta2.ContaboMachine, secretSpec infrastructurev1beta2.ContaboSecretSpec) (*infrastructurev1beta2.ContaboSecretStatus, error) {
	log := logf.FromContext(ctx)

	log.Info("Ensuring machine secret", "secretName", secretSpec.Name, "secretType", secretSpec.Type)

	// First, try to find existing secret by name
	secretStatus, err := r.findMachineSecretByName(ctx, secretSpec.Name)
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
	return r.createMachineSecret(ctx, contaboMachine, secretSpec)
}

// findMachineSecretByName searches for a secret by name using the Contabo API
func (r *ContaboMachineReconciler) findMachineSecretByName(ctx context.Context, secretName string) (*infrastructurev1beta2.ContaboSecretStatus, error) {
	// Call Contabo API to list secrets with name filter
	params := &models.RetrieveSecretListParams{
		Name: &secretName,
	}
	resp, err := r.ContaboClient.RetrieveSecretList(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve secret list: %w", err)
	}
	defer resp.Body.Close()

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

	return nil, nil
}

// createMachineSecret creates a new secret via the Contabo API for the machine
func (r *ContaboMachineReconciler) createMachineSecret(ctx context.Context, contaboMachine *infrastructurev1beta2.ContaboMachine, secretSpec infrastructurev1beta2.ContaboSecretSpec) (*infrastructurev1beta2.ContaboSecretStatus, error) {
	createRequest := models.CreateSecretRequest{
		Name:  secretSpec.Name,
		Value: secretSpec.Value,
		Type:  models.CreateSecretRequestType(secretSpec.Type),
	}

	params := &models.CreateSecretParams{}
	resp, err := r.ContaboClient.CreateSecret(ctx, params, createRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to create secret: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 && resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d when creating secret: %s", resp.StatusCode, string(body))
	}

	var createResponse models.CreateSecretResponse
	if err := json.NewDecoder(resp.Body).Decode(&createResponse); err != nil {
		return nil, fmt.Errorf("failed to decode create secret response: %w", err)
	}

	if len(createResponse.Data) == 0 {
		return nil, fmt.Errorf("no data returned from create secret API")
	}

	secret := createResponse.Data[0]
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
