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
	"crypto/sha256"
	"errors"
	"fmt"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
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

	infrastructurev1beta2 "github.com/ctnr-io/cluster-api-provider-contabo/api/v1beta2"
	contaboclient "github.com/ctnr-io/cluster-api-provider-contabo/pkg/contabo/v1.0.0/client"
	"github.com/ctnr-io/cluster-api-provider-contabo/pkg/contabo/v1.0.0/models"

	corev1 "k8s.io/api/core/v1"
)

// ContaboMachineReconciler reconciles a ContaboMachine object
type ContaboMachineReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	Recorder      record.EventRecorder
	ContaboClient *contaboclient.ClientWithResponses
	patchHelper   *patch.Helper
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
	r.patchHelper, err = patch.NewHelper(contaboMachine, r.Client)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Always attempt to Patch the ContaboMachine object and status after each reconciliation
	defer func() {
		if err := r.patchHelper.Patch(ctx, contaboMachine); err != nil {
			log.Error(err, "failed to patch ContaboMachine")
		}
	}()

	// Handle deleted machines
	if !contaboMachine.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, contaboMachine, contaboCluster)
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

	// Add finalizer
	controllerutil.AddFinalizer(contaboMachine, infrastructurev1beta2.MachineFinalizer)

	// Check if cluster infrastructure is ready
	if !contaboCluster.Status.Ready {
		log.Info("Waiting for cluster infrastructure to be ready")
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:   infrastructurev1beta2.ClusterInfrastructureReadyCondition,
			Status: metav1.ConditionFalse,
			Reason: infrastructurev1beta2.WaitingForClusterInfrastructureReason,
		})
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	// Mark cluster infrastructure as ready
	meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
		Type:   infrastructurev1beta2.ClusterInfrastructureReadyCondition,
		Status: metav1.ConditionTrue,
		Reason: infrastructurev1beta2.ClusterInfrastructureReadyReason,
	})

	// Get bootstrap data secret
	if machine.Spec.Bootstrap.DataSecretName == nil {
		log.Info("Bootstrap data secret is not available yet")
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:   infrastructurev1beta2.MachineReadyCondition,
			Status: metav1.ConditionFalse,
			Reason: infrastructurev1beta2.WaitingForBootstrapDataReason,
		})
		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	}

	// Get bootstrap data
	bootstrapDataSecret := &corev1.Secret{}
	bootstrapDataSecretName := client.ObjectKey{
		Namespace: contaboMachine.Namespace,
		Name:      *machine.Spec.Bootstrap.DataSecretName,
	}
	if err := r.Get(ctx, bootstrapDataSecretName, bootstrapDataSecret); err != nil {
		log.Info("Bootstrap data secret is not available yet")
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:   infrastructurev1beta2.MachineReadyCondition,
			Status: metav1.ConditionFalse,
			Reason: infrastructurev1beta2.WaitingForBootstrapDataReason,
		})
		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	}

	if _, ok := bootstrapDataSecret.Data["value"]; !ok || len(bootstrapDataSecret.Data["value"]) == 0 {
		log.Info("Bootstrap data secret is missing 'value' key or is empty")
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:   infrastructurev1beta2.MachineReadyCondition,
			Status: metav1.ConditionFalse,
			Reason: infrastructurev1beta2.WaitingForBootstrapDataReason,
		})
		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	}

	// Check if instance was already reconciled
	if contaboMachine.Status.Instance == nil {
		var instance *infrastructurev1beta2.InstanceResponse

		// Hash name of the machine to create unique instance names
		hash := sha256.New()
		if _, err := fmt.Fprintf(hash, "%s-%s", contaboMachine.Namespace, contaboMachine.Name); err != nil {
			log.Error(err, "Failed to write to hash")
		}
		machineHash := fmt.Sprintf("%x", hash.Sum(nil))[:6]
		displayName := fmt.Sprintf("capc %s %s %s %s", contaboCluster.Spec.PrivateNetwork.Region, contaboCluster.Status.ClusterUUID, machineHash, contaboMachine.Name)[:64]
		sshKeys := []int64{contaboCluster.Status.SshKey.SecretId}
		userData := string(bootstrapDataSecret.Data["value"])
		region := (models.CreateInstanceRequestRegion)(contaboCluster.Spec.PrivateNetwork.Region)
		imageId := DefaultUbuntuImageID

		// Check if instance with the same name already exists in Contabo API
		instancesListResp, _ := r.ContaboClient.RetrieveInstancesListWithResponse(ctx, &models.RetrieveInstancesListParams{
			DisplayName: &displayName,
		})
		if instancesListResp != nil && len(instancesListResp.JSON200.Data) > 0 {
			instance = convertListInstanceResponseData(&instancesListResp.JSON200.Data[0])
			log.Info("Found existing instance in Contabo API",
				"instanceID", instance.InstanceId,
				"instanceName", displayName)
		} else {
			// Check for available reusable instances
			displayNameEmpty := ""
			page := int64(0)
			size := int64(100)
			for {
				resp, err := r.ContaboClient.RetrieveInstancesListWithResponse(ctx, &models.RetrieveInstancesListParams{
					Page:        &page,
					Size:        &size,
					DisplayName: &displayNameEmpty,
					ProductIds:  &contaboMachine.Spec.Instance.ProductID,
					Region:      &contaboCluster.Spec.PrivateNetwork.Region,
				})
				if resp != nil && resp.JSON200.Data != nil && len(resp.JSON200.Data) > 0 {
					instance = convertListInstanceResponseData(&instancesListResp.JSON200.Data[0])
					log.Info("Found reusable instance in Contabo API",
						"instanceID", instance.InstanceId,
						"instanceName", instance.Name)
					break
				}
				if err != nil || resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
					break
				}
				page += 1
			}
		}

		// Create instance if no reusable instance was found
		if instance == nil {
			log.Info("No reusable instance found in Contabo API, will create a new one",
				"instanceName", displayName,
				"productID", contaboMachine.Spec.Instance.ProductID,
				"region", contaboCluster.Spec.PrivateNetwork.Region)
			instanceCreateResp, err := r.ContaboClient.CreateInstanceWithResponse(ctx, &models.CreateInstanceParams{}, models.CreateInstanceRequest{
				DisplayName: &displayName,
				ProductId:   &contaboMachine.Spec.Instance.ProductID,
				ImageId:     &imageId,
				Region:      &region,
				SshKeys:     &sshKeys,
			})
			if err != nil || instanceCreateResp.StatusCode() < 200 || instanceCreateResp.StatusCode() >= 300 {
				message := fmt.Sprintf("Failed to create instance: %v", err)
				meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
					Type:    infrastructurev1beta2.MachineReadyCondition,
					Status:  metav1.ConditionFalse,
					Reason:  infrastructurev1beta2.InstanceCreatingReason,
					Message: message,
				})
				log.Error(err, message)
				r.Recorder.Event(contaboMachine, corev1.EventTypeWarning, infrastructurev1beta2.InstanceCreatingReason, message)
				if patchErr := r.patchHelper.Patch(ctx, contaboCluster); patchErr != nil {
					log.Error(patchErr, "Failed to patch cluster status")
				}
				return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
			}

			instanceId := instanceCreateResp.JSON201.Data[0].InstanceId
			log.Info("Created new instance in Contabo API",
				"instanceID", instanceId,
				"instanceName", displayName)

			// Wait for instance to be fully created
			for {
				time.Sleep(10 * time.Second)
				instanceGetResp, err := r.ContaboClient.RetrieveInstanceWithResponse(ctx, instanceId, &models.RetrieveInstanceParams{})
				if err != nil || instanceGetResp.StatusCode() < 200 || instanceGetResp.StatusCode() >= 300 {
					message := fmt.Sprintf("Failed to retrieve instance after creation: %v", err)
					meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
						Type:    infrastructurev1beta2.MachineReadyCondition,
						Status:  metav1.ConditionFalse,
						Reason:  infrastructurev1beta2.InstanceCreatingReason,
						Message: message,
					})
					log.Error(err, message)
					r.Recorder.Event(contaboMachine, corev1.EventTypeWarning, infrastructurev1beta2.InstanceCreatingReason, message)
					if patchErr := r.patchHelper.Patch(ctx, contaboCluster); patchErr != nil {
						log.Error(patchErr, "Failed to patch cluster status")
					}
					return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
				}
				if instanceGetResp.JSON200.Data[0].Status == models.InstanceStatusRunning {
					instance = convertInstanceResponseData(&instanceGetResp.JSON200.Data[0])
					break
				}
				log.Info("Waiting for instance to be fully created",
					"instanceID", instanceId,
					"currentStatus", instance.Status)
			}
		}

		// Update the instance display name to mark it as used
		if instance.DisplayName != displayName {
			_, err := r.ContaboClient.PatchInstanceWithResponse(ctx, instance.InstanceId, nil, models.PatchInstanceRequest{
				DisplayName: &displayName,
			})
			if err != nil {
				message := fmt.Sprintf("Failed to update instance display name: %v", err)
				meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
					Type:    infrastructurev1beta2.MachineReadyCondition,
					Status:  metav1.ConditionFalse,
					Reason:  infrastructurev1beta2.InstanceCreatingReason,
					Message: message,
				})
				log.Error(err, message)
				r.Recorder.Event(contaboMachine, corev1.EventTypeWarning, infrastructurev1beta2.InstanceFailedReason, message)
				if patchErr := r.patchHelper.Patch(ctx, contaboCluster); patchErr != nil {
					log.Error(patchErr, "Failed to patch cluster status")
				}
				return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
			}
		}

		// If error message is set, instance is not usable, update display name to "capc <Region> error <ClusterUUID>" to avoid reuse and alert user
		if instance.ErrorMessage != nil && *instance.ErrorMessage != "" {
			message := fmt.Sprintf("Instance %d has error message: %s, retrying...", instance.InstanceId, *instance.ErrorMessage)
			err := errors.New(message)
			meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
				Type:    infrastructurev1beta2.MachineReadyCondition,
				Status:  metav1.ConditionFalse,
				Reason:  infrastructurev1beta2.InstanceFailedReason,
				Message: message,
			})
			log.Error(err, "")
			r.Recorder.Event(contaboMachine, corev1.EventTypeWarning, infrastructurev1beta2.InstanceFailedReason, message)
			if patchErr := r.patchHelper.Patch(ctx, contaboCluster); patchErr != nil {
				log.Error(patchErr, "Failed to patch cluster status")
			}
			displayName := fmt.Sprintf("capc %s error %s", contaboCluster.Spec.PrivateNetwork.Region, contaboCluster.Status.ClusterUUID)
			_, err = r.ContaboClient.PatchInstanceWithResponse(ctx, instance.InstanceId, nil, models.PatchInstanceRequest{
				DisplayName: &displayName,
			})
			return ctrl.Result{}, err
		}

		// Check status of the instance, should not be error if this is the case, we update the resource status and requeue
		switch instance.Status {
		case infrastructurev1beta2.InstanceStatusError:
		case infrastructurev1beta2.InstanceStatusUnknown:
		case infrastructurev1beta2.InstanceStatusManualProvisioning:
		case infrastructurev1beta2.InstanceStatusOther:
		case infrastructurev1beta2.InstanceStatusProductNotAvailable:
		case infrastructurev1beta2.InstanceStatusVerificationRequired:
			errorMessage := ""
			if instance.ErrorMessage != nil {
				errorMessage = *instance.ErrorMessage
			}
			message := fmt.Sprintf("Instance %d is in %s state: %s", instance.InstanceId, instance.Status, errorMessage)
			err := errors.New(message)
			meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
				Type:    infrastructurev1beta2.MachineReadyCondition,
				Status:  metav1.ConditionFalse,
				Reason:  infrastructurev1beta2.InstanceFailedReason,
				Message: message,
			})
			log.Error(err, "")
			r.Recorder.Event(contaboMachine, corev1.EventTypeWarning, infrastructurev1beta2.InstanceFailedReason, message)
			if patchErr := r.patchHelper.Patch(ctx, contaboCluster); patchErr != nil {
				log.Error(patchErr, "Failed to patch cluster status")
			}
			return ctrl.Result{}, err
		case infrastructurev1beta2.InstanceStatusPendingPayment:
		case infrastructurev1beta2.InstanceStatusProvisioning:
		case infrastructurev1beta2.InstanceStatusRescue:
		case infrastructurev1beta2.InstanceStatusResetPassword:
		case infrastructurev1beta2.InstanceStatusInstalling:
		case infrastructurev1beta2.InstanceStatusUninstalled:
			message := fmt.Sprintf("Instance %d is still in %s state", instance.InstanceId, instance.Status)
			meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
				Type:    infrastructurev1beta2.MachineReadyCondition,
				Status:  metav1.ConditionFalse,
				Reason:  infrastructurev1beta2.InstanceCreatingReason,
				Message: message,
			})
			log.Info(message)
			if patchErr := r.patchHelper.Patch(ctx, contaboCluster); patchErr != nil {
				log.Error(patchErr, "Failed to patch cluster status")
			}
			return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
		case infrastructurev1beta2.InstanceStatusStopped:
			// Start the instance if it is stopped
			message := fmt.Sprintf("Instance %d is stopped, starting it...", instance.InstanceId)
			meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
				Type:    infrastructurev1beta2.MachineReadyCondition,
				Status:  metav1.ConditionFalse,
				Reason:  infrastructurev1beta2.InstanceFailedReason,
				Message: message,
			})
			log.Info(message)
			if patchErr := r.patchHelper.Patch(ctx, contaboCluster); patchErr != nil {
				log.Error(patchErr, "Failed to patch cluster status")
			}
			return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
		case infrastructurev1beta2.InstanceStatusRunning:
		}

		// If there is no instance there, should fail the reconciliation
		if instance == nil {
			message := "instance should not be nil at this point"
			err := errors.New(message)
			meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
				Type:    infrastructurev1beta2.MachineReadyCondition,
				Status:  metav1.ConditionFalse,
				Reason:  infrastructurev1beta2.InstanceFailedReason,
				Message: message,
			})
			log.Error(err, message)
			r.Recorder.Event(contaboMachine, corev1.EventTypeWarning, infrastructurev1beta2.InstanceFailedReason, message)
			if patchErr := r.patchHelper.Patch(ctx, contaboCluster); patchErr != nil {
				log.Error(patchErr, "Failed to patch cluster status")
			}
			return ctrl.Result{RequeueAfter: 30 * time.Second}, err
		}

		// Assign instance to private network if not already assigned
		var privateNetwork *models.PrivateNetworkResponse
		{
			// Retrieve private network details
			privateNetworkGetResp, err := r.ContaboClient.RetrievePrivateNetworkWithResponse(ctx, contaboCluster.Status.PrivateNetwork.PrivateNetworkId, &models.RetrievePrivateNetworkParams{})
			if err != nil || privateNetworkGetResp.StatusCode() < 200 || privateNetworkGetResp.StatusCode() >= 300 {
				message := fmt.Sprintf("Failed to retrieve private network details: %v", err)
				meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
					Type:    infrastructurev1beta2.MachineReadyCondition,
					Status:  metav1.ConditionFalse,
					Reason:  infrastructurev1beta2.InstanceFailedReason,
					Message: message,
				})
				log.Error(err, message)
				r.Recorder.Event(contaboMachine, corev1.EventTypeWarning, infrastructurev1beta2.InstanceFailedReason, message)
				if patchErr := r.patchHelper.Patch(ctx, contaboCluster); patchErr != nil {
					log.Error(patchErr, "Failed to patch cluster status")
				}
				return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
			}
			privateNetwork = &privateNetworkGetResp.JSON200.Data[0]

			// Check if instance is already part of the private network
			assignedToPrivateNetwork := false
			for _, pnInstance := range privateNetwork.Instances {
				if pnInstance.InstanceId == instance.InstanceId {
					assignedToPrivateNetwork = true
					break
				}
			}

			// Assign instance to private network if not already assigned
			if !assignedToPrivateNetwork {
				log.Info("Assigning instance to private network",
					"instanceID", instance.InstanceId,
					"privateNetworkID", privateNetwork.PrivateNetworkId)
				_, err := r.ContaboClient.AssignInstancePrivateNetwork(ctx, privateNetwork.PrivateNetworkId, instance.InstanceId, nil)
				if err != nil {
					message := fmt.Sprintf("Failed to assign instance to private network: %v", err)
					meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
						Type:    infrastructurev1beta2.MachineReadyCondition,
						Status:  metav1.ConditionFalse,
						Reason:  infrastructurev1beta2.InstanceFailedReason,
						Message: message,
					})
					log.Error(err, message)
					r.Recorder.Event(contaboMachine, corev1.EventTypeWarning, infrastructurev1beta2.InstanceFailedReason, message)
					if patchErr := r.patchHelper.Patch(ctx, contaboCluster); patchErr != nil {
						log.Error(patchErr, "Failed to patch cluster status")
					}
					return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
				}

				log.Info("Rebooting instance to apply private network changes",
					"instanceID", instance.InstanceId)
				_, err = r.ContaboClient.Restart(ctx, instance.InstanceId, nil)
				if err != nil {
					message := fmt.Sprintf("Failed to reboot instance: %v", err)
					meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
						Type:    infrastructurev1beta2.MachineReadyCondition,
						Status:  metav1.ConditionFalse,
						Reason:  infrastructurev1beta2.InstanceFailedReason,
						Message: message,
					})
					log.Error(err, message)
					r.Recorder.Event(contaboMachine, corev1.EventTypeWarning, infrastructurev1beta2.InstanceFailedReason, message)
					if patchErr := r.patchHelper.Patch(ctx, contaboCluster); patchErr != nil {
						log.Error(patchErr, "Failed to patch cluster status")
					}
					return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
				}

				// Wait for instance to be fully restarted
				for {
					time.Sleep(10 * time.Second)
					instanceGetResp, err := r.ContaboClient.RetrieveInstanceWithResponse(ctx, instance.InstanceId, &models.RetrieveInstanceParams{})
					if err != nil || instanceGetResp.StatusCode() < 200 || instanceGetResp.StatusCode() >= 300 {
						message := fmt.Sprintf("Failed to retrieve instance after reboot: %v", err)
						meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
							Type:    infrastructurev1beta2.MachineReadyCondition,
							Status:  metav1.ConditionFalse,
							Reason:  infrastructurev1beta2.InstanceFailedReason,
							Message: message,
						})
						log.Error(err, message)
						r.Recorder.Event(contaboMachine, corev1.EventTypeWarning, infrastructurev1beta2.InstanceFailedReason, message)
						if patchErr := r.patchHelper.Patch(ctx, contaboCluster); patchErr != nil {
							log.Error(patchErr, "Failed to patch cluster status")
						}
						return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
					}
					if instanceGetResp.JSON200.Data[0].Status == models.InstanceStatusRunning {
						instance = convertInstanceResponseData(&instanceGetResp.JSON200.Data[0])
						break
					}
					log.Info("Waiting for instance to be fully restarted",
						"instanceID", instance.InstanceId,
						"currentStatus", instance.Status)
				}
			}
		}

		// Update ContaboMachine addresses based on private network assignment
		{
			addresses := []clusterv1.MachineAddress{}

			// Look for internal ip v4 in private network
			for _, pnInstance := range privateNetwork.Instances {
				if pnInstance.InstanceId == instance.InstanceId {
					addresses = append(addresses, clusterv1.MachineAddress{
						Type:    clusterv1.MachineInternalIP,
						Address: pnInstance.IpConfig.V4.Ip,
					})
					break
				}
			}
			// Look for external ip v4 in instance
			addresses = append(addresses, clusterv1.MachineAddress{
				Type:    clusterv1.MachineExternalIP,
				Address: instance.IpConfig.V4.Ip,
			})
			// Look for external ip v6 in instance
			addresses = append(addresses, clusterv1.MachineAddress{
				Type:    clusterv1.MachineExternalIP,
				Address: instance.IpConfig.V6.Ip,
			})

			// Add hostname entry
			addresses = append(addresses, clusterv1.MachineAddress{
				Type:    clusterv1.MachineHostName,
				Address: instance.Name,
			})

			// Update status addresses
			contaboMachine.Status.Addresses = addresses

			// Patch the ContaboMachine to save the updated addresses
			if patchErr := r.patchHelper.Patch(ctx, contaboCluster); patchErr != nil {
				log.Error(patchErr, "Failed to patch cluster status")
			}
		}

		// Reinstall instance with cloud-init bootstrap data
		resp, err := r.ContaboClient.ReinstallInstanceWithResponse(ctx, instance.InstanceId, &models.ReinstallInstanceParams{}, models.ReinstallInstanceRequest{
			SshKeys:      &sshKeys,
			DefaultUser:  (*models.ReinstallInstanceRequestDefaultUser)(instance.DefaultUser),
			ImageId:      contaboMachine.Status.Instance.ImageId,
			RootPassword: nil,
			UserData:     &userData,
		})
		if err != nil || resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
			message := fmt.Sprintf("Failed to reinstall instance: %v", err)
			meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
				Type:    infrastructurev1beta2.MachineReadyCondition,
				Status:  metav1.ConditionFalse,
				Reason:  infrastructurev1beta2.InstanceReinstallingReason,
				Message: message,
			})
			log.Error(err, message)
			r.Recorder.Event(contaboMachine, corev1.EventTypeWarning, infrastructurev1beta2.InstanceReinstallingReason, message)
			if patchErr := r.patchHelper.Patch(ctx, contaboCluster); patchErr != nil {
				log.Error(patchErr, "Failed to patch cluster status")
			}
			return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
		}

		// Wait for instance to be running after reinstall
		for {
			time.Sleep(10 * time.Second)
			instanceGetResp, err := r.ContaboClient.RetrieveInstanceWithResponse(ctx, instance.InstanceId, &models.RetrieveInstanceParams{})
			if err != nil || instanceGetResp.StatusCode() < 200 || instanceGetResp.StatusCode() >= 300 {
				message := fmt.Sprintf("Failed to retrieve instance after reinstall: %v", err)
				meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
					Type:    infrastructurev1beta2.MachineReadyCondition,
					Status:  metav1.ConditionFalse,
					Reason:  infrastructurev1beta2.InstanceReinstallingReason,
					Message: message,
				})
				log.Error(err, message)
				r.Recorder.Event(contaboMachine, corev1.EventTypeWarning, infrastructurev1beta2.InstanceReinstallingReason, message)
				if patchErr := r.patchHelper.Patch(ctx, contaboCluster); patchErr != nil {
					log.Error(patchErr, "Failed to patch cluster status")
				}
				return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
			}
			if instanceGetResp.JSON200.Data[0].Status == models.InstanceStatusRunning {
				instance = convertInstanceResponseData(&instanceGetResp.JSON200.Data[0])
				break
			}
			log.Info("Waiting for instance to be running after reinstall",
				"instanceID", instance.InstanceId,
				"currentStatus", instance.Status)
		}

		// Remove undesirable routes (e.g. old private network) from the ip table routes of the instance
		{
			// Connect via ssh and remove all ips that are not part of the private network or the main public ip
			if err := r.runMachineInstanceSshCommand(
				ctx,
				contaboMachine,
				contaboCluster,
				instance,
				"ip route | grep -v default | grep -v "+privateNetwork.Cidr+" | awk '{print $1}' | xargs -r -n1 sudo ip route del",
			); err != nil {
				log.Error(err, "Failed to clean up instance routes")
			}
		}

		// Wait for cloud-init to finish
		{
			log.Info("Waiting for cloud-init to finish on instance",
				"instanceID", instance.InstanceId,
				"instanceIP", instance.IpConfig.V4.Ip)

			if err := r.runMachineInstanceSshCommand(
				ctx,
				contaboMachine,
				contaboCluster,
				instance,
				"cloud-init status --wait",
			); err != nil {
				message := fmt.Sprintf("Failed to wait for cloud-init on instance: %v", err)
				meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
					Type:    infrastructurev1beta2.MachineReadyCondition,
					Status:  metav1.ConditionFalse,
					Reason:  infrastructurev1beta2.MachineSSHKeysFailedReason,
					Message: message,
				})
				log.Error(err, message)
				r.Recorder.Event(contaboMachine, corev1.EventTypeWarning, infrastructurev1beta2.MachineSshKeyFailedReason, message)
				if patchErr := r.patchHelper.Patch(ctx, contaboCluster); patchErr != nil {
					log.Error(patchErr, "Failed to patch cluster status")
				}
				return ctrl.Result{RequeueAfter: 30 * time.Second}, err
			}
			log.Info("cloud-init has finished on instance",
				"instanceID", instance.InstanceId)
		}

		// Update ContaboMachine status with instance details
		contaboMachine.Status.Instance = instance
		contaboMachine.Status.Ready = true
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:   infrastructurev1beta2.MachineReadyCondition,
			Status: metav1.ConditionTrue,
			Reason: infrastructurev1beta2.InstanceReadyReason,
		})
		r.Recorder.Event(contaboMachine, corev1.EventTypeNormal, infrastructurev1beta2.InstanceReadyReason, "Instance is ready")

		log.Info("Instance is ready",
			"instanceID", instance.InstanceId,
			"instanceIP", instance.IpConfig.V4.Ip)
	}
	// else {
	// TODO: Handle updates to existing instances (e.g., changing instance type)
	// }

	return ctrl.Result{}, nil
}

func (r *ContaboMachineReconciler) reconcileDelete(ctx context.Context, contaboMachine *infrastructurev1beta2.ContaboMachine, contaboCluster *infrastructurev1beta2.ContaboCluster) (ctrl.Result, error) {
	_ = contaboCluster // may be used in future for cluster-specific cleanup logic
	log := logf.FromContext(ctx)

	log.Info("Reconciling ContaboMachine delete - setting instance available for reuse")

	// Update machine condition
	meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
		Type:   infrastructurev1beta2.MachineReadyCondition,
		Status: metav1.ConditionFalse,
		Reason: infrastructurev1beta2.InstanceDeletingReason,
	})
	if patchErr := r.patchHelper.Patch(ctx, contaboCluster); patchErr != nil {
		log.Error(patchErr, "Failed to patch cluster status")
	}
	log.Info("Machine marked for deletion, proceeding with instance cleanup",
		"instanceID", contaboMachine.Status.Instance.InstanceId)

	// First, stop the instance
	_, err := r.ContaboClient.Stop(ctx, contaboMachine.Status.Instance.InstanceId, nil)
	if err != nil {
		message := fmt.Sprintf("Failed to stop instance: %v", err)
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.MachineReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceFailedReason,
			Message: message,
		})
		log.Error(err, message)
		r.Recorder.Event(contaboMachine, corev1.EventTypeWarning, infrastructurev1beta2.InstanceFailedReason, message)
		if patchErr := r.patchHelper.Patch(ctx, contaboCluster); patchErr != nil {
			log.Error(patchErr, "Failed to patch cluster status")
		}
		return ctrl.Result{RequeueAfter: 30 * time.Second}, err
	}

	// Unassign instance from private network
	privateNetworkGetResp, err := r.ContaboClient.RetrievePrivateNetworkWithResponse(ctx, contaboCluster.Status.PrivateNetwork.PrivateNetworkId, &models.RetrievePrivateNetworkParams{})
	if err != nil || privateNetworkGetResp.StatusCode() < 200 || privateNetworkGetResp.StatusCode() >= 300 {
		message := fmt.Sprintf("Failed to retrieve private network details: %v", err)
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.MachineReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceFailedReason,
			Message: message,
		})
		log.Error(err, message)
		r.Recorder.Event(contaboMachine, corev1.EventTypeWarning, infrastructurev1beta2.InstanceFailedReason, message)
		if patchErr := r.patchHelper.Patch(ctx, contaboCluster); patchErr != nil {
			log.Error(patchErr, "Failed to patch cluster status")
		}
		return ctrl.Result{RequeueAfter: 30 * time.Second}, err
	}

	// Update instance display name to mark it as reusable
	displayName := ""
	_, err = r.ContaboClient.PatchInstanceWithResponse(ctx, contaboMachine.Status.Instance.InstanceId, nil, models.PatchInstanceRequest{
		DisplayName: &displayName,
	})
	if err != nil {
		message := fmt.Sprintf("Failed to update instance display name: %v", err)
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.MachineReadyCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrastructurev1beta2.InstanceFailedReason,
			Message: message,
		})
		log.Error(err, message)
		r.Recorder.Event(contaboMachine, corev1.EventTypeWarning, infrastructurev1beta2.InstanceFailedReason, message)
		if patchErr := r.patchHelper.Patch(ctx, contaboCluster); patchErr != nil {
			log.Error(patchErr, "Failed to patch cluster status")
		}
		return ctrl.Result{RequeueAfter: 30 * time.Second}, err
	}

	// Remove finalizer
	controllerutil.RemoveFinalizer(contaboMachine, infrastructurev1beta2.MachineFinalizer)
	log.Info("Removed finalizer from ContaboMachine")

	return ctrl.Result{}, nil
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
		Complete(r)
}

// Convert OAPI Instance models to CAPC Instance models
func convertListInstanceResponseData(instanceList *models.ListInstancesResponseData) *infrastructurev1beta2.ContaboInstanceStatus {
	if instanceList == nil {
		return nil
	}
	addons := make([]infrastructurev1beta2.AddOnResponse, len(instanceList.AddOns))
	for i, addon := range instanceList.AddOns {
		addons[i] = infrastructurev1beta2.AddOnResponse{
			Id:       addon.Id,
			Quantity: addon.Quantity,
		}
	}

	addionalIps := make([]infrastructurev1beta2.AdditionalIp, len(instanceList.AdditionalIps))
	for i, ip := range instanceList.AdditionalIps {
		addionalIps[i] = infrastructurev1beta2.AdditionalIp{
			V4: infrastructurev1beta2.IpV4{
				Gateway:     ip.V4.Gateway,
				Ip:          ip.V4.Ip,
				NetmaskCidr: ip.V4.NetmaskCidr,
			},
		}
	}

	var defaultUser infrastructurev1beta2.InstanceResponseDefaultUser
	if instanceList.DefaultUser != nil {
		defaultUser = infrastructurev1beta2.InstanceResponseDefaultUser(string(*instanceList.DefaultUser))
	}

	ipConfig := infrastructurev1beta2.IpConfig{
		V4: infrastructurev1beta2.IpV4{
			Gateway:     instanceList.IpConfig.V4.Gateway,
			Ip:          instanceList.IpConfig.V4.Ip,
			NetmaskCidr: instanceList.IpConfig.V4.NetmaskCidr,
		},
		V6: infrastructurev1beta2.IpV6{
			Gateway:     instanceList.IpConfig.V6.Gateway,
			Ip:          instanceList.IpConfig.V6.Ip,
			NetmaskCidr: instanceList.IpConfig.V6.NetmaskCidr,
		},
	}

	productType := infrastructurev1beta2.InstanceResponseProductType(string(instanceList.ProductType))
	tenantId := infrastructurev1beta2.InstanceResponseTenantId(string(instanceList.TenantId))
	status := infrastructurev1beta2.InstanceStatus(string(instanceList.Status))

	return &infrastructurev1beta2.ContaboInstanceStatus{
		AddOns:        addons,
		AdditionalIps: addionalIps,
		CancelDate:    instanceList.CancelDate.Format(time.RFC3339),
		CpuCores:      instanceList.CpuCores,
		CreatedDate:   instanceList.CreatedDate.Unix(),
		CustomerId:    instanceList.CustomerId,
		DataCenter:    instanceList.DataCenter,
		DefaultUser:   &defaultUser,
		DiskMb:        instanceList.DiskMb,
		DisplayName:   instanceList.DisplayName,
		ErrorMessage:  instanceList.ErrorMessage,
		ImageId:       instanceList.ImageId,
		InstanceId:    instanceList.InstanceId,
		IpConfig:      ipConfig,
		MacAddress:    instanceList.MacAddress,
		Name:          instanceList.Name,
		OsType:        instanceList.OsType,
		ProductId:     instanceList.ProductId,
		ProductName:   instanceList.ProductName,
		ProductType:   productType,
		RamMb:         instanceList.RamMb,
		Region:        instanceList.Region,
		RegionName:    instanceList.RegionName,
		SshKeys:       instanceList.SshKeys,
		Status:        status,
		TenantId:      tenantId,
		VHostId:       instanceList.VHostId,
		VHostName:     instanceList.VHostName,
		VHostNumber:   instanceList.VHostNumber,
	}
}

// Convert OAPI Instance models to CAPC Instance models
func convertInstanceResponseData(instanceList *models.InstanceResponse) *infrastructurev1beta2.ContaboInstanceStatus {
	if instanceList == nil {
		return nil
	}
	addons := make([]infrastructurev1beta2.AddOnResponse, len(instanceList.AddOns))
	for i, addon := range instanceList.AddOns {
		addons[i] = infrastructurev1beta2.AddOnResponse{
			Id:       addon.Id,
			Quantity: addon.Quantity,
		}
	}

	addionalIps := make([]infrastructurev1beta2.AdditionalIp, len(instanceList.AdditionalIps))
	for i, ip := range instanceList.AdditionalIps {
		addionalIps[i] = infrastructurev1beta2.AdditionalIp{
			V4: infrastructurev1beta2.IpV4{
				Gateway:     ip.V4.Gateway,
				Ip:          ip.V4.Ip,
				NetmaskCidr: ip.V4.NetmaskCidr,
			},
		}
	}

	var defaultUser infrastructurev1beta2.InstanceResponseDefaultUser
	if instanceList.DefaultUser != nil {
		defaultUser = infrastructurev1beta2.InstanceResponseDefaultUser(string(*instanceList.DefaultUser))
	}

	ipConfig := infrastructurev1beta2.IpConfig{
		V4: infrastructurev1beta2.IpV4{
			Gateway:     instanceList.IpConfig.V4.Gateway,
			Ip:          instanceList.IpConfig.V4.Ip,
			NetmaskCidr: instanceList.IpConfig.V4.NetmaskCidr,
		},
		V6: infrastructurev1beta2.IpV6{
			Gateway:     instanceList.IpConfig.V6.Gateway,
			Ip:          instanceList.IpConfig.V6.Ip,
			NetmaskCidr: instanceList.IpConfig.V6.NetmaskCidr,
		},
	}

	productType := infrastructurev1beta2.InstanceResponseProductType(string(instanceList.ProductType))
	tenantId := infrastructurev1beta2.InstanceResponseTenantId(string(instanceList.TenantId))
	status := infrastructurev1beta2.InstanceStatus(string(instanceList.Status))

	return &infrastructurev1beta2.ContaboInstanceStatus{
		AddOns:        addons,
		AdditionalIps: addionalIps,
		CancelDate:    instanceList.CancelDate.Format(time.RFC3339),
		CpuCores:      instanceList.CpuCores,
		CreatedDate:   instanceList.CreatedDate.Unix(),
		CustomerId:    instanceList.CustomerId,
		DataCenter:    instanceList.DataCenter,
		DefaultUser:   &defaultUser,
		DiskMb:        instanceList.DiskMb,
		DisplayName:   instanceList.DisplayName,
		ErrorMessage:  instanceList.ErrorMessage,
		ImageId:       instanceList.ImageId,
		InstanceId:    instanceList.InstanceId,
		IpConfig:      ipConfig,
		MacAddress:    instanceList.MacAddress,
		Name:          instanceList.Name,
		OsType:        instanceList.OsType,
		ProductId:     instanceList.ProductId,
		ProductName:   instanceList.ProductName,
		ProductType:   productType,
		RamMb:         instanceList.RamMb,
		Region:        instanceList.Region,
		RegionName:    instanceList.RegionName,
		SshKeys:       instanceList.SshKeys,
		Status:        status,
		TenantId:      tenantId,
		VHostId:       instanceList.VHostId,
		VHostName:     instanceList.VHostName,
		VHostNumber:   instanceList.VHostNumber,
	}
}

func (r *ContaboMachineReconciler) runMachineInstanceSshCommand(ctx context.Context, contaboMachine *infrastructurev1beta2.ContaboMachine, contaboCluster *infrastructurev1beta2.ContaboCluster, instance *infrastructurev1beta2.InstanceResponse, command string) error {
	log := logf.FromContext(ctx)

	// Get ssh-key from secret
	sshPrivateKeySecret := &corev1.Secret{}
	sshPrivateKeySecretName := client.ObjectKey{
		Namespace: contaboMachine.Namespace,
		Name:      contaboCluster.Status.SshKey.Name,
	}
	if err := r.Get(ctx, sshPrivateKeySecretName, sshPrivateKeySecret); err != nil {
		return fmt.Errorf("failed to get SSH private key secret %s/%s: %v", sshPrivateKeySecretName.Namespace, sshPrivateKeySecretName.Name, err)
	}
	// Connect to the instance via SSH and wait for cloud-init to finish
	// We try to connect every 10 seconds for up to 15 minutes
	sshPrivateKey, ok := sshPrivateKeySecret.Data["id_rsa"]
	if !ok || len(sshPrivateKey) == 0 {
		return fmt.Errorf("SSH private key secret is missing 'id_rsa' key or is empty")
	}
	// SSH client configuration
	host := instance.IpConfig.V4.Ip
	user := string(infrastructurev1beta2.InstanceResponseDefaultUserAdmin)
	if instance.DefaultUser != nil {
		user = string(*instance.DefaultUser)
	}
	signer, err := ssh.ParsePrivateKey(sshPrivateKey)
	if err != nil {
		return fmt.Errorf("failed to parse SSH private key: %v", err)
	}
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // For testing; use proper host key validation in production
		Timeout:         10 * time.Second,
	}
	// SSH client configuration
	// Retry logic for SSH connection
	var sshClient *ssh.Client
	maxRetries := 10
	retryDelay := 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		sshClient, err = ssh.Dial("tcp", net.JoinHostPort(host, "22"), config)
		if err == nil {
			break
		}
		log.Info("SSH connection failed, retrying in %v... (attempt %d/%d)", retryDelay, i+1, maxRetries)
		time.Sleep(retryDelay)
	}
	if err != nil {
		return fmt.Errorf("failed to establish SSH connection to %s: %v", host, err)
	}
	defer func() {
		if closeErr := sshClient.Close(); closeErr != nil {
			log.Error(closeErr, "Failed to close SSH client")
		}
	}()

	// Create a session to run the cloud-init status command
	session, err := sshClient.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session: %v", err)
	}
	defer func() {
		if closeErr := session.Close(); closeErr != nil {
			log.Error(closeErr, "Failed to close SSH session")
		}
	}()

	// Run the command and capture output
	_, err = session.CombinedOutput(command)
	if err != nil {
		return fmt.Errorf("failed to run command '%s': %v", command, err)
	}
	return nil
}
