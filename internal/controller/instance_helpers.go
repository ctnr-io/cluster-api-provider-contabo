package controller

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ctnr-io/cluster-api-provider-contabo/pkg/contabo/v1.0.0/models"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	infrastructurev1beta2 "github.com/ctnr-io/cluster-api-provider-contabo/api/v1beta2"

	_ "embed"
)

// getExistingInstance attempts to find an existing instance either from status or by display name
func (r *ContaboMachineReconciler) getExistingInstance(
	ctx context.Context,
	contaboMachine *infrastructurev1beta2.ContaboMachine,
	contaboCluster *infrastructurev1beta2.ContaboCluster,
) (*infrastructurev1beta2.ContaboInstanceStatus, error) {
	displayName := FormatDisplayName(contaboMachine, contaboCluster)

	// Optimize by checking status first, but check display name
	if contaboMachine.Status.Instance != nil {
		// Get the latest status from the instance
		instanceResp, err := r.ContaboClient.RetrieveInstanceWithResponse(ctx, contaboMachine.Status.Instance.InstanceId, nil)
		if err != nil || instanceResp.StatusCode() < 200 || instanceResp.StatusCode() >= 300 {
			return nil, fmt.Errorf("failed to find instance %d from Contabo API", contaboMachine.Status.Instance.InstanceId)
		}
		instance := convertInstanceResponseData(&instanceResp.JSON200.Data[0])
		if instance != nil && instance.DisplayName == displayName {
			return instance, nil
		}
	}

	// Try to find instance by display name
	instanceListResp, err := r.ContaboClient.RetrieveInstancesListWithResponse(ctx, &models.RetrieveInstancesListParams{
		DisplayName: &displayName,
	})
	if err == nil && len(instanceListResp.JSON200.Data) > 0 {
		return convertListInstanceResponseData(&instanceListResp.JSON200.Data[0]), nil
	}

	return nil, nil
}

// getUsedInstanceNames returns a map of instance names that are currently in use by ContaboMachines
func (r *ContaboMachineReconciler) getUsedInstanceNames(ctx context.Context) (map[string]*infrastructurev1beta2.ContaboMachine, error) {
	log := logf.FromContext(ctx)
	usedNames := make(map[string]*infrastructurev1beta2.ContaboMachine)

	// List all ContaboMachines across all namespaces
	var contaboMachineList infrastructurev1beta2.ContaboMachineList
	if err := r.List(ctx, &contaboMachineList); err != nil {
		log.Error(err, "Failed to list ContaboMachines")
		return nil, fmt.Errorf("failed to list ContaboMachines: %w", err)
	}

	// Build a map of instance names that are in use
	for _, machine := range contaboMachineList.Items {
		if machine.Spec.Instance.Name != nil && *machine.Spec.Instance.Name != "" {
			usedNames[*machine.Spec.Instance.Name] = &machine
		}
	}

	log.V(1).Info("Found used instance names", "count", len(usedNames), "names", usedNames)
	return usedNames, nil
}

// findReusableInstance looks for available instances that can be reused
func (r *ContaboMachineReconciler) findReusableInstance(
	ctx context.Context,
	contaboMachine *infrastructurev1beta2.ContaboMachine,
	contaboCluster *infrastructurev1beta2.ContaboCluster,
) (*infrastructurev1beta2.ContaboInstanceStatus, error) {
	log := logf.FromContext(ctx)
	displayNameEmpty := ""
	page := int64(1)
	size := int64(100)

	// Get map of instance names already in use by other ContaboMachines
	usedInstanceNames, err := r.getUsedInstanceNames(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get used instance names: %w", err)
	}

	for {
		resp, err := r.ContaboClient.RetrieveInstancesListWithResponse(ctx, &models.RetrieveInstancesListParams{
			Page:        &page,
			Size:        &size,
			DisplayName: &displayNameEmpty,
			ProductIds:  contaboMachine.Spec.Instance.ProductId,
			Region:      &contaboCluster.Spec.PrivateNetwork.Region,
			Name:        contaboMachine.Spec.Instance.Name,
		})
		if err != nil {
			body := []byte{}
			if resp != nil && resp.Body != nil {
				body = resp.Body
			}
			log.Info("Failed to find instance from Contabo API when looking for reusable instances",
				"error", err, "body", string(body))
			return nil, fmt.Errorf("failed to list instances: %w", err)
		}

		if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
			if resp.StatusCode() == 429 {
				time.Sleep(60 * time.Second)
				continue // Retry after rate limit cooldown
			}
			if resp.StatusCode() == 404 {
				return nil, nil // No instances found
			}
			return nil, fmt.Errorf("failed to list instances, status code: %d", resp.StatusCode())
		}

		if resp.JSON200 != nil && resp.JSON200.Data != nil {
			if len(resp.JSON200.Data) == 0 {
				return nil, nil
			}

			// Find instance with empty display name that is not already in use
			for i := range resp.JSON200.Data {
				instance := &resp.JSON200.Data[i]

				// Check if instance has empty display name and is not cancelled
				if instance.DisplayName != displayNameEmpty || instance.CancelDate != nil {
					continue
				}

				// Check that no other ContaboMachine is using this instance name
				if usedInstanceNames[instance.Name] != nil && usedInstanceNames[instance.Name].UID != contaboMachine.UID {
					log.V(1).Info("Skipping instance with name already in use by another ContaboMachine",
						"instanceID", instance.InstanceId,
						"instanceName", instance.Name)
					continue
				}

				convertedInstance := convertListInstanceResponseData(instance)

				// Reset the instance by removing any private network assignments
				if err := r.resetInstance(ctx, contaboMachine, convertedInstance, nil); err != nil {
					log.Error(err, "Failed to reset instance",
						"instanceID", convertedInstance.InstanceId)
					continue
				}

				return convertedInstance, nil
			}
		}

		page++
		time.Sleep(5 * time.Second) // Wait to prevent rate limiting
	}
}

// createNewInstance creates a new instance based on the provisioning type
func (r *ContaboMachineReconciler) createNewInstance(
	ctx context.Context,
	contaboMachine *infrastructurev1beta2.ContaboMachine,
	contaboCluster *infrastructurev1beta2.ContaboCluster,
) (*infrastructurev1beta2.ContaboInstanceStatus, error) {
	log := logf.FromContext(ctx)

	// Check provisioning type
	switch {
	case contaboMachine.Spec.Instance.ProvisioningType == nil:
		return nil, fmt.Errorf("no provisioning type specified")
	case *contaboMachine.Spec.Instance.ProvisioningType == infrastructurev1beta2.ContaboInstanceProvisioningTypeReuseOnly:
		log.Info("No reusable instance found in Contabo API, user must intervene")
		return nil, fmt.Errorf("no reusable instance found for reuse-only provisioning type")
	case *contaboMachine.Spec.Instance.ProvisioningType == infrastructurev1beta2.ContaboInstanceProvisioningTypeReuseOrCreate:
		log.Info("No reusable instance found in Contabo API, will create a new one",
			"productID", contaboMachine.Spec.Instance.ProductId,
			"region", contaboCluster.Spec.PrivateNetwork.Region)

		if contaboMachine.Spec.Instance.Name != nil && *contaboMachine.Spec.Instance.Name != "" {
			msg := fmt.Sprintf("instance name must not be specified to create a new instance: %s", *contaboMachine.Spec.Instance.Name)
			if err := r.resetInstance(ctx, contaboMachine, nil, ptr.To(msg)); err != nil {
				log.Error(err, "Failed to reset instance creation attempt with invalid name",
					"instanceName", *contaboMachine.Spec.Instance.Name)
			}
			return nil, errors.New(msg)
		}

		sshKeys := []int64{contaboCluster.Status.SshKey.SecretId}
		imageId := DefaultUbuntuImageID
		region := *ConvertRegionToCreateInstanceRegion(contaboCluster.Spec.PrivateNetwork.Region)

		instanceCreateResp, err := r.ContaboClient.CreateInstanceWithResponse(ctx, &models.CreateInstanceParams{}, models.CreateInstanceRequest{
			ProductId: contaboMachine.Spec.Instance.ProductId,
			Period:    1,
			ImageId:   &imageId,
			Region:    &region,
			SshKeys:   &sshKeys,
			AddOns: &models.CreateInstanceAddons{
				PrivateNetworking: ptr.To(map[string]interface{}{}),
			},
			DisplayName: ptr.To(FormatDisplayName(contaboMachine, contaboCluster)),
			DefaultUser: ptr.To(models.CreateInstanceRequestDefaultUserAdmin),
		})
		if err != nil || instanceCreateResp.StatusCode() < 200 || instanceCreateResp.StatusCode() >= 300 {
			log.Error(err, "Failed to create instance in Contabo API",
				"statusCode", instanceCreateResp.StatusCode(),
				"body", string(instanceCreateResp.Body))
			return nil, fmt.Errorf("failed to create instance: %w", err)
		}

		instanceId := instanceCreateResp.JSON201.Data[0].InstanceId
		log.Info("Created new instance in Contabo API",
			"instanceID", instanceId)

		retrieveInstanceResponse, err := r.ContaboClient.RetrieveInstanceWithResponse(ctx, instanceId, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve newly created instance: %w", err)
		}

		instance := convertInstanceResponseData(&retrieveInstanceResponse.JSON200.Data[0])

		return instance, nil
	default:
		return nil, fmt.Errorf("unknown Instance provisioningType: %v", contaboMachine.Spec.Instance.ProvisioningType)
	}
}

// updateInstanceState handles display name updates and private networking setup
func (r *ContaboMachineReconciler) updateInstanceState(
	ctx context.Context,
	contaboMachine *infrastructurev1beta2.ContaboMachine,
	contaboCluster *infrastructurev1beta2.ContaboCluster,
	instance *infrastructurev1beta2.ContaboInstanceStatus,
) error {
	log := logf.FromContext(ctx)
	displayName := FormatDisplayName(contaboMachine, contaboCluster)

	// Update display name if needed
	if instance.DisplayName != displayName {
		log.Info("Updating instance display name to mark it as used",
			"instanceID", instance.InstanceId,
			"oldDisplayName", instance.DisplayName,
			"newDisplayName", displayName)
		_, err := r.ContaboClient.PatchInstanceWithResponse(ctx, instance.InstanceId, nil, models.PatchInstanceRequest{
			DisplayName: &displayName,
		})
		if err != nil {
			return fmt.Errorf("failed to update instance display name: %w", err)
		}
	}

	// Add private networking if not already added
	privateNetworkFound := false
	for _, addons := range instance.AddOns {
		if addons.Id == 1477 { // Private Networking addon ID
			privateNetworkFound = true
			break
		}
	}

	if !privateNetworkFound {
		log.Info("Adding private networking to instance",
			"instanceID", instance.InstanceId)
		_, err := r.ContaboClient.UpgradeInstance(ctx, instance.InstanceId, nil, models.UpgradeInstanceJSONRequestBody{
			PrivateNetworking: ptr.To(map[string]interface{}{}),
		})
		if err != nil {
			return fmt.Errorf("failed to add private networking to instance: %w", err)
		}
	}

	return nil
}

// validateInstanceStatus validates the instance status and handles error conditions
func (r *ContaboMachineReconciler) validateInstanceStatus(ctx context.Context, contaboMachine *infrastructurev1beta2.ContaboMachine) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// If error message is set, instance is not usable, update display name to "capc <Region> error <ClusterUUID>" to avoid reuse and alert user
	if contaboMachine.Status.Instance.ErrorMessage != nil && *contaboMachine.Status.Instance.ErrorMessage != "" {
		log.Info("Instance has error message, marking as failed",
			"instanceID", contaboMachine.Status.Instance.InstanceId,
			"errorMessage", *contaboMachine.Status.Instance.ErrorMessage)

		// Reset instance
		err := r.resetInstance(ctx, contaboMachine, contaboMachine.Status.Instance, contaboMachine.Status.Instance.ErrorMessage)
		if err != nil {
			log.Error(err, "Failed to reset instance",
				"instanceID", contaboMachine.Status.Instance.InstanceId)
		}

		// Remove instance form the ContaboMachine status to avoid further processing
		instance := contaboMachine.Status.Instance
		contaboMachine.Status.Instance = nil

		return ctrl.Result{RequeueAfter: 5 * time.Second}, r.handleError(
			ctx,
			contaboMachine,
			errors.New(*instance.ErrorMessage),
			infrastructurev1beta2.InstanceFailedReason,
			fmt.Sprintf("Instance %d has error message: %s, retrying...", instance.InstanceId, *instance.ErrorMessage),
		)
	}

	// Check status of the instance, should not be error if this is the case, we update the resource status and requeue
	switch contaboMachine.Status.Instance.Status {
	case infrastructurev1beta2.InstanceStatusError:
	case infrastructurev1beta2.InstanceStatusUnknown:
	case infrastructurev1beta2.InstanceStatusManualProvisioning:
	case infrastructurev1beta2.InstanceStatusOther:
	case infrastructurev1beta2.InstanceStatusProductNotAvailable:
	case infrastructurev1beta2.InstanceStatusVerificationRequired:
		errorMessage := "Instance in error state"
		if contaboMachine.Status.Instance.ErrorMessage != nil {
			errorMessage = *contaboMachine.Status.Instance.ErrorMessage
		}
		// Update display name to avoid reuse
		log.Info("Instance is in error state, marking as failed",
			"instanceID", contaboMachine.Status.Instance.InstanceId,
			"status", contaboMachine.Status.Instance.Status,
			"errorMessage", errorMessage)

		// Reset instance
		err := r.resetInstance(ctx, contaboMachine, contaboMachine.Status.Instance, &errorMessage)
		if err != nil {
			log.Error(err, "Failed to reset instance",
				"instanceID", contaboMachine.Status.Instance.InstanceId)
		}

		// Remove instance form the ContaboMachine status to avoid further processing
		instance := contaboMachine.Status.Instance
		contaboMachine.Status.Instance = nil

		return ctrl.Result{}, r.handleError(
			ctx,
			contaboMachine,
			errors.New(errorMessage),
			infrastructurev1beta2.InstanceFailedReason,
			fmt.Sprintf("Instance %d is in %s states", instance.InstanceId, instance.Status),
		)
	case infrastructurev1beta2.InstanceStatusPendingPayment:
	case infrastructurev1beta2.InstanceStatusProvisioning:
	case infrastructurev1beta2.InstanceStatusRescue:
	case infrastructurev1beta2.InstanceStatusResetPassword:
	case infrastructurev1beta2.InstanceStatusUninstalled:
		return ctrl.Result{RequeueAfter: 10 * time.Second}, r.handleError(
			ctx,
			contaboMachine,
			errors.New("instance is not ready"),
			infrastructurev1beta2.InstanceCreatingReason,
			fmt.Sprintf("Instance %d is in %s state", contaboMachine.Status.Instance.InstanceId, contaboMachine.Status.Instance.Status),
		)
	case infrastructurev1beta2.InstanceStatusInstalling:
		message := fmt.Sprintf("Instance %d is installing, waiting for it to be running...", contaboMachine.Status.Instance.InstanceId)
		log.Info(message)
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceCreatingReason,
			Message: message,
		})
		return ctrl.Result{RequeueAfter: 20 * time.Second}, nil
	case infrastructurev1beta2.InstanceStatusStopped:
		message := fmt.Sprintf("Instance %d is stopped, starting it...", contaboMachine.Status.Instance.InstanceId)
		log.Info(message)
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceReadyReason,
			Message: message,
		})
		// Start the instance if it is stopped
		_, err := r.ContaboClient.Start(ctx, contaboMachine.Status.Instance.InstanceId, nil)
		if err != nil {
			return ctrl.Result{RequeueAfter: 10 * time.Second}, r.handleError(
				ctx,
				contaboMachine,
				err,
				infrastructurev1beta2.InstanceFailedReason,
				fmt.Sprintf("Failed to start instance %d", contaboMachine.Status.Instance.InstanceId),
			)
		}
		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	case infrastructurev1beta2.InstanceStatusRunning:
		message := fmt.Sprintf("Instance %d is running", contaboMachine.Status.Instance.InstanceId)
		log.Info(message)
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.InstanceReadyCondition,
			Status:  metav1.ConditionTrue,
			Reason:  infrastructurev1beta2.InstanceReadyReason,
			Message: message,
		})
	}

	// If there is no instance there, should fail the reconciliation
	if contaboMachine.Status.Instance == nil {
		return ctrl.Result{}, r.handleError(
			ctx,
			contaboMachine,
			errors.New("instance is nil"),
			infrastructurev1beta2.InstanceFailedReason,
			"Instance should not be nil at this point",
		)
	}

	// Instance is valid and running
	meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
		Type:   infrastructurev1beta2.InstanceReadyCondition,
		Status: metav1.ConditionTrue,
		Reason: infrastructurev1beta2.InstanceReadyReason,
	})
	meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
		Type:   clusterv1.ReadyCondition,
		Status: metav1.ConditionFalse,
		Reason: infrastructurev1beta2.InstanceReadyReason,
	})

	return ctrl.Result{}, nil
}
