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
	"fmt"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	infrastructurev1beta2 "github.com/ctnr-io/cluster-api-provider-contabo/api/v1beta2"
)

// reconcileControlPlaneEndpoint reconciles the control plane endpoint
func (r *ContaboClusterReconciler) reconcileControlPlaneEndpoint(ctx context.Context, contaboCluster *infrastructurev1beta2.ContaboCluster) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// Get the current condition
	endpointCondition := meta.FindStatusCondition(contaboCluster.Status.Conditions, infrastructurev1beta2.ControlPlaneEndpointReadyCondition)

	// State machine for control plane endpoint lifecycle
	switch {
	case endpointCondition == nil:
		// Initial state: Check if endpoint is already set
		return r.checkControlPlaneEndpoint(ctx, contaboCluster)

	case endpointCondition.Reason == infrastructurev1beta2.WaitingForControlPlaneEndpointReason:
		// Still waiting for endpoint to be set by control plane machine
		return r.checkControlPlaneEndpoint(ctx, contaboCluster)

	case endpointCondition.Reason == infrastructurev1beta2.ControlPlaneEndpointReadyReason:
		// Endpoint is ready, nothing to do
		return ctrl.Result{}, nil

	case endpointCondition.Reason == infrastructurev1beta2.ControlPlaneEndpointFailedReason:
		// Endpoint failed, try to recover
		return r.handleControlPlaneEndpointFailure(ctx, contaboCluster)

	default:
		// Unknown state, restart the process
		log.Info("Unknown control plane endpoint state, restarting check process")
		return r.checkControlPlaneEndpoint(ctx, contaboCluster)
	}
}

// checkControlPlaneEndpoint checks if control plane endpoint is set
func (r *ContaboClusterReconciler) checkControlPlaneEndpoint(ctx context.Context, contaboCluster *infrastructurev1beta2.ContaboCluster) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	if contaboCluster.Spec.ControlPlaneEndpoint.Host == "" {
		log.Info("Control plane endpoint not set, waiting for control plane machine")
		meta.SetStatusCondition(&contaboCluster.Status.Conditions, metav1.Condition{
			Type:   infrastructurev1beta2.ControlPlaneEndpointReadyCondition,
			Status: metav1.ConditionFalse,
			Reason: infrastructurev1beta2.WaitingForControlPlaneEndpointReason,
		})
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	log.Info("Control plane endpoint is set", "host", contaboCluster.Spec.ControlPlaneEndpoint.Host)
	meta.SetStatusCondition(&contaboCluster.Status.Conditions, metav1.Condition{
		Type:   infrastructurev1beta2.ControlPlaneEndpointReadyCondition,
		Status: metav1.ConditionTrue,
		Reason: infrastructurev1beta2.ControlPlaneEndpointReadyReason,
	})

	return ctrl.Result{}, nil
}

// handleControlPlaneEndpointFailure handles failed control plane endpoint
func (r *ContaboClusterReconciler) handleControlPlaneEndpointFailure(ctx context.Context, contaboCluster *infrastructurev1beta2.ContaboCluster) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	log.Info("Handling control plane endpoint failure - attempting recovery")

	// Check current condition to determine failure type
	endpointCondition := meta.FindStatusCondition(contaboCluster.Status.Conditions, infrastructurev1beta2.ControlPlaneEndpointReadyCondition)

	// Implement recovery strategies

	// Strategy 1: Check if endpoint was cleared and needs to be reset
	if contaboCluster.Spec.ControlPlaneEndpoint.Host == "" {
		log.Info("Control plane endpoint is empty - waiting for control plane machine to set it")
		meta.SetStatusCondition(&contaboCluster.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.ControlPlaneEndpointReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.WaitingForControlPlaneEndpointReason,
			Message: "Control plane endpoint cleared, waiting for control plane machine",
		})
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	// Strategy 2: Validate the current endpoint
	if err := r.validateControlPlaneEndpoint(ctx, contaboCluster); err != nil {
		log.Error(err, "Control plane endpoint validation failed")

		// Determine requeue delay based on failure history
		var requeueDelay time.Duration = 1 * time.Minute // Default delay for endpoint issues

		if endpointCondition != nil && strings.Contains(endpointCondition.Message, "retry") {
			requeueDelay = 3 * time.Minute // Longer delay for repeated failures
		}

		meta.SetStatusCondition(&contaboCluster.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.ControlPlaneEndpointReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.ControlPlaneEndpointFailedReason,
			Message: fmt.Sprintf("Endpoint validation failed, retry scheduled: %s", err.Error()),
		})

		return ctrl.Result{RequeueAfter: requeueDelay}, nil
	}

	// Strategy 3: Reset the endpoint state and recheck
	log.Info("Resetting control plane endpoint state for recovery")
	return r.checkControlPlaneEndpoint(ctx, contaboCluster)
}

// validateControlPlaneEndpoint validates the control plane endpoint
func (r *ContaboClusterReconciler) validateControlPlaneEndpoint(ctx context.Context, contaboCluster *infrastructurev1beta2.ContaboCluster) error {
	log := logf.FromContext(ctx)

	endpoint := contaboCluster.Spec.ControlPlaneEndpoint

	// Basic validation checks
	if endpoint.Host == "" {
		return fmt.Errorf("control plane endpoint host is empty")
	}

	if endpoint.Port == 0 {
		return fmt.Errorf("control plane endpoint port is not set")
	}

	// Validate host format (basic check)
	if !strings.Contains(endpoint.Host, ".") && !strings.Contains(endpoint.Host, ":") {
		// Not an IP or FQDN, might be invalid
		log.V(1).Info("Control plane endpoint host format may be invalid", "host", endpoint.Host)
	}

	// Validate port range
	if endpoint.Port < 1 || endpoint.Port > 65535 {
		return fmt.Errorf("control plane endpoint port %d is out of valid range", endpoint.Port)
	}

	log.V(1).Info("Control plane endpoint validation passed",
		"host", endpoint.Host,
		"port", endpoint.Port)

	return nil
}
