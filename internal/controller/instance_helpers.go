package controller

import (
	"context"
	"fmt"
	"time"

	infrastructurev1beta2 "github.com/ctnr-io/cluster-api-provider-contabo/api/v1beta2"
	"github.com/ctnr-io/cluster-api-provider-contabo/pkg/contabo/v1.0.0/models"
	"k8s.io/utils/ptr"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// getExistingInstance attempts to find an existing instance either from status or by display name
func (r *ContaboMachineReconciler) getExistingInstance(
	ctx context.Context,
	contaboMachine *infrastructurev1beta2.ContaboMachine,
	contaboCluster *infrastructurev1beta2.ContaboCluster,
) (*infrastructurev1beta2.ContaboInstanceStatus, error) {
	if contaboMachine.Status.Instance != nil {
		// Get the latest status from the instance
		instanceResp, err := r.ContaboClient.RetrieveInstanceWithResponse(ctx, contaboMachine.Status.Instance.InstanceId, nil)
		if err != nil || instanceResp.StatusCode() < 200 || instanceResp.StatusCode() >= 300 {
			return nil, fmt.Errorf("failed to find instance %d from Contabo API", contaboMachine.Status.Instance.InstanceId)
		}
		return convertInstanceResponseData(&instanceResp.JSON200.Data[0]), nil
	}

	// Try to find instance by display name
	displayName := FormatDisplayName(contaboMachine, contaboCluster)
	instanceListResp, err := r.ContaboClient.RetrieveInstancesListWithResponse(ctx, &models.RetrieveInstancesListParams{
		DisplayName: &displayName,
	})
	if err == nil && len(instanceListResp.JSON200.Data) > 0 {
		return convertListInstanceResponseData(&instanceListResp.JSON200.Data[0]), nil
	}

	return nil, nil
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

	for {
		resp, err := r.ContaboClient.RetrieveInstancesListWithResponse(ctx, &models.RetrieveInstancesListParams{
			Page:        &page,
			Size:        &size,
			DisplayName: &displayNameEmpty,
			ProductIds:  &contaboMachine.Spec.Instance.ProductId,
			Region:      &contaboCluster.Spec.PrivateNetwork.Region,
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

			// Find instance with empty display name
			for i := range resp.JSON200.Data {
				if resp.JSON200.Data[i].DisplayName == displayNameEmpty && resp.JSON200.Data[i].CancelDate == nil {
					instance := convertListInstanceResponseData(&resp.JSON200.Data[i])

					// Reset the instance by removing any private network assignments
					if err := r.resetInstance(ctx, instance); err != nil {
						log.Error(err, "Failed to reset instance",
							"instanceID", instance.InstanceId)
						continue
					}
					return instance, nil
				}
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
) error {
	log := logf.FromContext(ctx)

	// Check provisioning type
	switch {
	case contaboMachine.Spec.Instance.ProvisioningType == nil:
		return nil // No provisioning type specified
	case *contaboMachine.Spec.Instance.ProvisioningType == infrastructurev1beta2.ContaboInstanceProvisioningTypeReuseOnly:
		log.Info("No reusable instance found in Contabo API, user must intervene")
		return nil
	case *contaboMachine.Spec.Instance.ProvisioningType == infrastructurev1beta2.ContaboInstanceProvisioningTypeReuseOrCreate:
		log.Info("No reusable instance found in Contabo API, will create a new one",
			"productID", contaboMachine.Spec.Instance.ProductId,
			"region", contaboCluster.Spec.PrivateNetwork.Region)

		sshKeys := []int64{contaboCluster.Status.SshKey.SecretId}
		imageId := DefaultUbuntuImageID
		region := *ConvertRegionToCreateInstanceRegion(contaboCluster.Spec.PrivateNetwork.Region)

		instanceCreateResp, err := r.ContaboClient.CreateInstanceWithResponse(ctx, &models.CreateInstanceParams{}, models.CreateInstanceRequest{
			ProductId: &contaboMachine.Spec.Instance.ProductId,
			Period:    1,
			ImageId:   &imageId,
			Region:    &region,
			SshKeys:   &sshKeys,
			AddOns: &models.CreateInstanceAddons{
				PrivateNetworking: ptr.To(map[string]interface{}{}),
			},
		})
		if err != nil || instanceCreateResp.StatusCode() < 200 || instanceCreateResp.StatusCode() >= 300 {
			log.Error(err, "Failed to create instance in Contabo API",
				"statusCode", instanceCreateResp.StatusCode(),
				"body", string(instanceCreateResp.Body))
			return fmt.Errorf("failed to create instance: %w", err)
		}

		instanceId := instanceCreateResp.JSON201.Data[0].InstanceId
		log.Info("Created new instance in Contabo API",
			"instanceID", instanceId)

		return nil
	default:
		return fmt.Errorf("unknown Instance provisioningType: %v", contaboMachine.Spec.Instance.ProvisioningType)
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
