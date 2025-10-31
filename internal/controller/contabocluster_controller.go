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
	"github.com/google/uuid"

	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"golang.org/x/crypto/ssh"

	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/utils/ptr"
)

// ContaboClusterReconciler reconciles a ContaboCluster object
type ContaboClusterReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	Recorder      record.EventRecorder
	ContaboClient *contaboclient.ClientWithResponses
	patchHelper   *patch.Helper
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=contaboclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=contaboclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=contaboclusters/finalizers,verbs=update
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters;clusters/status,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=discovery.k8s.io,resources=endpointslices,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *ContaboClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	log.Info("Reconciling ContaboCluster", "namespace", req.Namespace, "name", req.Name)

	// Fetch the ContaboCluster instance
	contaboCluster := &infrastructurev1beta2.ContaboCluster{}
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

	// Initialize the patch helper
	r.patchHelper, err = patch.NewHelper(contaboCluster, r.Client)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Always attempt to Patch the ContaboCluster object and status after each reconciliation
	defer func() {
		if patchErr := r.patchHelper.Patch(ctx, contaboCluster); patchErr != nil {
			log.Error(err, "failed to patch ContaboCluster")
		}
		// Wait to be sure patch is applied
		time.Sleep(1 * time.Second)
	}()

	// Handle deleted clusters
	if !contaboCluster.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, contaboCluster)
	}

	// Handle non-deleted clusters
	return r.reconcileApply(ctx, contaboCluster)
}

func (r *ContaboClusterReconciler) reconcileApply(ctx context.Context, contaboCluster *infrastructurev1beta2.ContaboCluster) (ctrl.Result, error) {
	// Initialize basic cluster setup
	controllerutil.AddFinalizer(contaboCluster, infrastructurev1beta2.ClusterFinalizer)

	// Ensure cluster has a unique UUID for global identification
	r.ensureClusterUUID(ctx, contaboCluster)

	// Check if private network was created
	if err := r.reconcilePrivateNetwork(ctx, contaboCluster); err != nil {
		return ctrl.Result{}, err
	}

	// Check if SSH key was created
	if err := r.reconcileSSHKey(ctx, contaboCluster); err != nil {
		return ctrl.Result{}, err
	}

	// Mark cluster infrastructure as ready after private network and SSH keys are created
	// This allows KubeadmControlPlane to proceed with creating control plane machines
	r.markClusterReadyAndProvisioned(ctx, contaboCluster)

	// Create and/or use control plane endpoint proxy service and endpoint slices to be able to connect w/ kubeconfig
	if result, err := r.reconcileControlPlaneEndpoint(ctx, contaboCluster); err != nil || result.RequeueAfter != 0 {
		return result, err
	}

	return ctrl.Result{}, nil
}

// markClusterReady sets the cluster infrastructure as ready after private network and SSH keys are created
// This signals to KubeadmControlPlane that it can start creating control plane machines
func (r *ContaboClusterReconciler) markClusterReadyAndProvisioned(ctx context.Context, contaboCluster *infrastructurev1beta2.ContaboCluster) {
	log := logf.FromContext(ctx)

	// Only mark as ready if not already ready and private network + SSH key are available
	if !contaboCluster.Status.Ready && contaboCluster.Status.PrivateNetwork != nil && contaboCluster.Status.SshKey != nil {
		log.Info("Private network and SSH key ready, marking ContaboCluster infrastructure as ready")

		meta.SetStatusCondition(&contaboCluster.Status.Conditions, metav1.Condition{
			Type:    infrastructurev1beta2.ClusterReadyCondition,
			Status:  metav1.ConditionTrue,
			Reason:  infrastructurev1beta2.ClusterAvailableReason,
			Message: "ContaboCluster infrastructure is ready",
		})

		contaboCluster.Status.Ready = true

		// This is what the Cluster controller checks via the v1beta2 contract
		// which KubeadmControlPlane waits for before creating machines
		if contaboCluster.Status.Initialization == nil {
			contaboCluster.Status.Initialization = &infrastructurev1beta2.ContaboClusterInitializationStatus{}
		}
		contaboCluster.Status.Initialization.Provisioned = true
	}
}

// ensureClusterUUID ensures the cluster has a unique UUID and returns it
func (r *ContaboClusterReconciler) ensureClusterUUID(ctx context.Context, contaboCluster *infrastructurev1beta2.ContaboCluster) string {
	log := logf.FromContext(ctx)

	if contaboCluster.Status.ClusterUUID == "" {
		clusterUUID := uuid.New().String()
		contaboCluster.Status.ClusterUUID = clusterUUID

		// Set initial creating state
		meta.SetStatusCondition(&contaboCluster.Status.Conditions, metav1.Condition{
			Type:   infrastructurev1beta2.ClusterReadyCondition,
			Status: metav1.ConditionFalse,
			Reason: infrastructurev1beta2.ClusterCreatingReason,
		})
		log.Info("Assigned new cluster UUID", "clusterUUID", clusterUUID)
		return clusterUUID
	}

	clusterUUID := contaboCluster.Status.ClusterUUID

	// Set initial updating state
	if !meta.IsStatusConditionTrue(contaboCluster.Status.Conditions, infrastructurev1beta2.ClusterReadyCondition) {
		meta.SetStatusCondition(&contaboCluster.Status.Conditions, metav1.Condition{
			Type:   infrastructurev1beta2.ClusterReadyCondition,
			Status: metav1.ConditionFalse,
			Reason: infrastructurev1beta2.ClusterUpdatingReason,
		})
		log.Info("Cluster already has a UUID, continuing reconciliation", "clusterUUID", clusterUUID)
	}

	return clusterUUID
}

// reconcilePrivateNetwork ensures the private network exists and is configured
func (r *ContaboClusterReconciler) reconcilePrivateNetwork(ctx context.Context, contaboCluster *infrastructurev1beta2.ContaboCluster) error {
	log := logf.FromContext(ctx)

	log.Info("Reconciling private network for ContaboCluster", "cluster", contaboCluster.Name)

	// Check if private network was created
	if contaboCluster.Status.PrivateNetwork == nil {
		var privateNetwork *models.PrivateNetworkResponse
		privateNetworkName := FormatPrivateNetworkName(contaboCluster)

		// Check if private network with the same name already exists in Contabo API
		resp, _ := r.ContaboClient.RetrievePrivateNetworkListWithResponse(ctx, &models.RetrievePrivateNetworkListParams{
			Name: &privateNetworkName,
		})

		if resp != nil && resp.JSON200 != nil && resp.JSON200.Data != nil && len(resp.JSON200.Data) > 0 {
			privateNetwork = (*models.PrivateNetworkResponse)(&resp.JSON200.Data[0])
		} else {
			log.Info("Private network not found in Contabo API, creating new one", "privateNetworkName", privateNetworkName)

			// Create private network if not found
			description := "Private network created by Cluster API Provider Contabo"
			privateNetworkCreateResp, err := r.ContaboClient.CreatePrivateNetworkWithResponse(ctx, nil, models.CreatePrivateNetworkJSONRequestBody{
				Name:        privateNetworkName,
				Description: &description,
				Region:      &contaboCluster.Spec.PrivateNetwork.Region,
			})
			if err != nil || privateNetworkCreateResp.StatusCode() < 200 || privateNetworkCreateResp.StatusCode() >= 300 {
				return r.handleError(
					ctx,
					contaboCluster,
					err,
					infrastructurev1beta2.ClusterPrivateNetworkReadyCondition,
					infrastructurev1beta2.ClusterPrivateNetworkFailedReason,
					"Failed to create private network",
				)
			}

			privateNetwork = &privateNetworkCreateResp.JSON201.Data[0]
		}
		// Update status with private network info
		contaboCluster.Status.PrivateNetwork = &infrastructurev1beta2.ContaboPrivateNetworkStatus{
			Name:             privateNetwork.Name,
			PrivateNetworkId: privateNetwork.PrivateNetworkId,
			Region:           privateNetwork.Region,
			AvailableIps:     privateNetwork.AvailableIps,
			Cidr:             privateNetwork.Cidr,
			CreatedDate:      privateNetwork.CreatedDate.UTC().Unix(),
			Instances:        []infrastructurev1beta2.Instances{},
			CustomerId:       privateNetwork.CustomerId,
			TenantId:         privateNetwork.TenantId,
			Description:      privateNetwork.Description,
			DataCenter:       privateNetwork.DataCenter,
			RegionName:       privateNetwork.RegionName,
		}
		meta.SetStatusCondition(&contaboCluster.Status.Conditions, metav1.Condition{
			Type:   infrastructurev1beta2.ClusterPrivateNetworkReadyCondition,
			Status: metav1.ConditionTrue,
			Reason: infrastructurev1beta2.ClusterAvailableReason,
		})
		log.Info("Using private network", "privateNetworkName", privateNetwork.Name, "privateNetworkId", privateNetwork.PrivateNetworkId)
	}
	// else {
	// TODO: Check if private network configuration matches spec and update if necessary
	// }

	return nil
}

// reconcileSSHKey ensures the SSH key exists and is configured
func (r *ContaboClusterReconciler) reconcileSSHKey(ctx context.Context, contaboCluster *infrastructurev1beta2.ContaboCluster) error {
	log := logf.FromContext(ctx)

	log.Info("Reconciling SSH key for ContaboCluster", "cluster", contaboCluster.Name)

	// Check if SSH key was created
	if contaboCluster.Status.SshKey == nil {
		var sshKey *models.SecretResponse
		sshKeyName := FormatSshKeyName(contaboCluster)
		sshKeySecretName := FormatSshKeySecretName(contaboCluster)

		// Check if SSH key with the same name already exists in Contabo API
		resp, _ := r.ContaboClient.RetrieveSecretListWithResponse(ctx, &models.RetrieveSecretListParams{
			Name: &sshKeyName,
		})

		if resp != nil && resp.JSON200 != nil && resp.JSON200.Data != nil && len(resp.JSON200.Data) > 0 {
			sshKey = &resp.JSON200.Data[0]
		} else {
			log.Info("SSH key not found in Contabo API, creating new one", "sshKeyName", sshKeyName)

			// Generate an ssh key pair
			privateKey, publicKey, err := generateSSHKeyPair()
			if err != nil {
				return r.handleError(
					ctx,
					contaboCluster,
					err,
					infrastructurev1beta2.ClusterSshKeyReadyCondition,
					infrastructurev1beta2.ClusterSshKeyFailedReason,
					"Failed to generate SSH key pair",
				)
			}

			// Delete secret if already exists
			err = r.Get(ctx, client.ObjectKey{
				Name:      sshKeySecretName,
				Namespace: contaboCluster.Namespace,
			}, &corev1.Secret{})
			if err == nil {
				// Secret already exists, delete it first
				err = r.Delete(ctx, &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      sshKeySecretName,
						Namespace: contaboCluster.Namespace,
					},
				})
				if err != nil {
					return r.handleError(
						ctx,
						contaboCluster,
						err,
						infrastructurev1beta2.ClusterSshKeyReadyCondition,
						infrastructurev1beta2.ClusterSshKeyFailedReason,
						"Failed to delete existing SSH key secret",
					)
				}
			}

			// Create new secret
			secretData := map[string][]byte{
				"id_rsa":     []byte(privateKey),
				"id_rsa.pub": []byte(publicKey),
			}

			// Verify data integrity before storing in secret
			log.Info("Storing SSH keys in Kubernetes secret",
				"secretName", sshKeySecretName)

			if err := r.Create(ctx, &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      sshKeySecretName,
					Namespace: contaboCluster.Namespace,
				},
				Data: secretData,
			}); err != nil {
				return r.handleError(
					ctx,
					contaboCluster,
					err,
					infrastructurev1beta2.ClusterSshKeyReadyCondition,
					infrastructurev1beta2.ClusterSshKeyFailedReason,
					"Failed to create SSH key secret",
				)
			}

			// Create SSH key if not found
			// Trim the public key for Contabo API (remove trailing newline)
			trimmedPublicKey := strings.TrimSpace(publicKey)

			// Log public key formatting for Contabo API
			log.Info("Submitting SSH public key to Contabo API",
				"sshKeyName", sshKeyName)

			sshKeyCreateResp, err := r.ContaboClient.CreateSecretWithResponse(ctx, &models.CreateSecretParams{}, models.CreateSecretRequest{
				Name:  sshKeyName,
				Value: trimmedPublicKey,
				Type:  "ssh",
			})
			if err != nil || sshKeyCreateResp.StatusCode() < 200 || sshKeyCreateResp.StatusCode() >= 300 {
				return r.handleError(
					ctx,
					contaboCluster,
					err,
					infrastructurev1beta2.ClusterSshKeyReadyCondition,
					infrastructurev1beta2.ClusterSshKeyFailedReason,
					fmt.Sprintf("Failed to submit SSH public key to Contabo API: %s", sshKeyName),
				)
			}
			sshKey = &sshKeyCreateResp.JSON201.Data[0]
		}

		// Update status with SSH key info
		contaboCluster.Status.SshKey = &infrastructurev1beta2.ContaboSshKeyStatus{
			Name:     sshKey.Name,
			SecretId: int64(sshKey.SecretId),
			Value:    sshKey.Value,
		}
		meta.SetStatusCondition(&contaboCluster.Status.Conditions, metav1.Condition{
			Type:   infrastructurev1beta2.ClusterSshKeyReadyCondition,
			Status: metav1.ConditionTrue,
			Reason: infrastructurev1beta2.ClusterAvailableReason,
		})
		log.Info("Using SSH key", "sshKeyName", sshKey.Name, "sshKeyId", sshKey.SecretId)
	}
	// else {
	// TODO: Check if SSH key configuration matches spec and update if necessary
	// }

	return nil
}

// reconcileControlPlaneEndpoint ensures the control plane endpoint is set
func (r *ContaboClusterReconciler) reconcileControlPlaneEndpoint(ctx context.Context, contaboCluster *infrastructurev1beta2.ContaboCluster) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	log.Info("Reconciling control plane endpoint for ContaboCluster", "cluster", contaboCluster.Name)

	meta.SetStatusCondition(&contaboCluster.Status.Conditions, metav1.Condition{
		Type:   clusterv1.ClusterControlPlaneAvailableCondition,
		Status: metav1.ConditionTrue,
		Reason: clusterv1.ClusterControlPlaneMachinesReadyReason,
	})

	// Retrieve control plane endpoint from the first ContaboMachine in the cluster
	controlPlaneMachines := &infrastructurev1beta2.ContaboMachineList{}
	if err := r.List(ctx, controlPlaneMachines, client.InNamespace(contaboCluster.Namespace), client.MatchingLabels{clusterv1.ClusterNameLabel: contaboCluster.Name}, client.HasLabels{clusterv1.MachineControlPlaneLabel}); err != nil || len(controlPlaneMachines.Items) == 0 {
		log.Info("No control plane machines found yet, requeuing")
		meta.SetStatusCondition(&contaboCluster.Status.Conditions, metav1.Condition{
			Type:   clusterv1.ClusterControlPlaneAvailableCondition,
			Status: metav1.ConditionFalse,
			Reason: clusterv1.WaitingForControlPlaneInitializedReason,
		})
		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	}

	// reconcile controlplane endpoint service and endpoint slices for the control plane endpoint
	if err := r.reconcileControlPlaneService(ctx, contaboCluster); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.reconcileControlPlaneEndpointSlices(ctx, contaboCluster, controlPlaneMachines); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *ContaboClusterReconciler) reconcileControlPlaneService(ctx context.Context, contaboCluster *infrastructurev1beta2.ContaboCluster) error {
	log := logf.FromContext(ctx)

	// Service name must match the control plane endpoint host for DNS resolution
	serviceName := fmt.Sprintf("%s-apiserver", contaboCluster.Name)
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: contaboCluster.Namespace,
			Labels: map[string]string{
				clusterv1.ClusterNameLabel: contaboCluster.Name,
				"component":                "apiserver",
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: contaboCluster.APIVersion,
					Kind:       contaboCluster.Kind,
					Name:       contaboCluster.Name,
					UID:        contaboCluster.UID,
					Controller: ptr.To(true),
				},
			},
		},
		Spec: corev1.ServiceSpec{
			Type:      corev1.ServiceTypeClusterIP,
			ClusterIP: "None", // Headless service
			Ports: []corev1.ServicePort{
				{
					Name:     "https",
					Port:     contaboCluster.Spec.ControlPlaneEndpoint.Port,
					Protocol: corev1.ProtocolTCP,
				},
			},
		},
	}

	// Try to get existing service
	existingService := &corev1.Service{}
	err := r.Get(ctx, client.ObjectKey{
		Name:      serviceName,
		Namespace: contaboCluster.Namespace,
	}, existingService)

	if err != nil {
		if apierrors.IsNotFound(err) {
			// Create the service
			log.Info("Creating control plane endpoint service", "serviceName", serviceName)
			if err := r.Create(ctx, service); err != nil {
				return fmt.Errorf("failed to create control plane endpoint service: %w", err)
			}
			log.Info("Created control plane endpoint service", "serviceName", serviceName)
		} else {
			return fmt.Errorf("failed to get control plane endpoint service: %w", err)
		}
	} else {
		// Update the service if needed
		existingService.Spec.Ports = service.Spec.Ports
		log.Info("Updating control plane endpoint service", "serviceName", serviceName)
		if err := r.Update(ctx, existingService); err != nil {
			return fmt.Errorf("failed to update control plane endpoint service: %w", err)
		}
	}

	return nil
}

func (r *ContaboClusterReconciler) reconcileControlPlaneEndpointSlices(ctx context.Context, contaboCluster *infrastructurev1beta2.ContaboCluster, controlPlaneContaboMachineList *infrastructurev1beta2.ContaboMachineList) error {
	log := logf.FromContext(ctx)

	// Service name must match the control plane endpoint host for DNS resolution
	serviceName := fmt.Sprintf("%s-apiserver", contaboCluster.Name)
	endpointSliceName := fmt.Sprintf("%s-apiserver", contaboCluster.Name)

	// Collect all control plane instance IPs
	var endpoints []discoveryv1.Endpoint
	for _, machine := range controlPlaneContaboMachineList.Items {
		if machine.Status.Instance != nil && machine.Status.Instance.IpConfig.V4.Ip != "" {
			endpoints = append(endpoints, discoveryv1.Endpoint{
				Addresses: []string{machine.Status.Instance.IpConfig.V4.Ip},
				Conditions: discoveryv1.EndpointConditions{
					Ready: ptr.To(true),
				},
			})
		}
	}

	endpointPort := contaboCluster.Spec.ControlPlaneEndpoint.Port

	endpointSlice := &discoveryv1.EndpointSlice{
		ObjectMeta: metav1.ObjectMeta{
			Name:      endpointSliceName,
			Namespace: contaboCluster.Namespace,
			Labels: map[string]string{
				clusterv1.ClusterNameLabel:   contaboCluster.Name,
				"component":                  "apiserver",
				discoveryv1.LabelServiceName: serviceName,
				discoveryv1.LabelManagedBy:   "cluster-api-provider-contabo",
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: contaboCluster.APIVersion,
					Kind:       contaboCluster.Kind,
					Name:       contaboCluster.Name,
					UID:        contaboCluster.UID,
					Controller: ptr.To(true),
				},
			},
		},
		AddressType: discoveryv1.AddressTypeIPv4,
		Endpoints:   endpoints,
		Ports: []discoveryv1.EndpointPort{
			{
				Name:     ptr.To("https"),
				Port:     ptr.To(endpointPort),
				Protocol: ptr.To(corev1.ProtocolTCP),
			},
		},
	}

	// Try to get existing endpoint slice
	existingEndpointSlice := &discoveryv1.EndpointSlice{}
	err := r.Get(ctx, client.ObjectKey{
		Name:      endpointSliceName,
		Namespace: contaboCluster.Namespace,
	}, existingEndpointSlice)

	if err != nil {
		if apierrors.IsNotFound(err) {
			// Create the endpoint slice
			log.Info("Creating control plane endpoint slice", "endpointSliceName", endpointSliceName, "endpointCount", len(endpoints))
			if err := r.Create(ctx, endpointSlice); err != nil {
				return fmt.Errorf("failed to create control plane endpoint slice: %w", err)
			}
			log.Info("Created control plane endpoint slice", "endpointSliceName", endpointSliceName)
		} else {
			return fmt.Errorf("failed to get control plane endpoint slice: %w", err)
		}
	} else {
		// Update the endpoint slice if needed
		existingEndpointSlice.Endpoints = endpointSlice.Endpoints
		existingEndpointSlice.Ports = endpointSlice.Ports
		log.Info("Updating control plane endpoint slice", "endpointSliceName", endpointSliceName, "endpointCount", len(endpoints))
		if err := r.Update(ctx, existingEndpointSlice); err != nil {
			return fmt.Errorf("failed to update control plane endpoint slice: %w", err)
		}
	}

	return nil
}

func (r *ContaboClusterReconciler) deleteControlPlaneService(ctx context.Context, contaboCluster *infrastructurev1beta2.ContaboCluster) error {
	log := logf.FromContext(ctx)

	// Service name matches the control plane endpoint host
	serviceName := fmt.Sprintf("%s-apiserver", contaboCluster.Name)

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: contaboCluster.Namespace,
		},
	}

	err := r.Delete(ctx, service)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("Control plane endpoint service already deleted", "serviceName", serviceName)
			return nil
		}
		return fmt.Errorf("failed to delete control plane endpoint service: %w", err)
	}

	log.Info("Deleted control plane endpoint service", "serviceName", serviceName)
	return nil
}

func (r *ContaboClusterReconciler) deleteControlPlaneEndpointSlices(ctx context.Context, contaboCluster *infrastructurev1beta2.ContaboCluster) error {
	log := logf.FromContext(ctx)

	endpointSliceName := fmt.Sprintf("%s-apiserver", contaboCluster.Name)

	endpointSlice := &discoveryv1.EndpointSlice{
		ObjectMeta: metav1.ObjectMeta{
			Name:      endpointSliceName,
			Namespace: contaboCluster.Namespace,
		},
	}

	err := r.Delete(ctx, endpointSlice)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("Control plane endpoint slice already deleted", "endpointSliceName", endpointSliceName)
			return nil
		}
		return fmt.Errorf("failed to delete control plane endpoint slice: %w", err)
	}

	log.Info("Deleted control plane endpoint slice", "endpointSliceName", endpointSliceName)
	return nil
}

// reconcileDelete may return different ctrl.Result values in future implementations
func (r *ContaboClusterReconciler) reconcileDelete(ctx context.Context, contaboCluster *infrastructurev1beta2.ContaboCluster) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	log.Info("Reconciling ContaboCluster delete")

	// Update cluster condition
	meta.SetStatusCondition(&contaboCluster.Status.Conditions, metav1.Condition{
		Type:   infrastructurev1beta2.ClusterReadyCondition,
		Status: metav1.ConditionFalse,
		Reason: infrastructurev1beta2.ClusterDeletingReason,
	})
	log.Info("Cluster marked for deletion, proceeding with resource cleanup")

	// Delete network infrastructure
	if contaboCluster.Status.PrivateNetwork != nil {
		meta.SetStatusCondition(&contaboCluster.Status.Conditions, metav1.Condition{
			Type:   infrastructurev1beta2.ClusterPrivateNetworkReadyCondition,
			Status: metav1.ConditionFalse,
			Reason: infrastructurev1beta2.ClusterPrivateNetworkDeletingReason,
		})
		log.Info("Deleting private network", "privateNetworkId", contaboCluster.Status.PrivateNetwork.PrivateNetworkId)
		// Check if private network exists in Contabo API
		resp, err := r.ContaboClient.RetrievePrivateNetworkWithResponse(ctx, contaboCluster.Status.PrivateNetwork.PrivateNetworkId, nil)
		if err != nil || resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
			// If the private network is not found, we can assume it has already been deleted
			log.Info("Private network not found in Contabo API, assuming already deleted", "privateNetworkId", contaboCluster.Status.PrivateNetwork.PrivateNetworkId)
		} else {
			privateNetwork := resp.JSON200.Data[0]

			// Unassign all instances from the private network
			if len(privateNetwork.Instances) > 0 {
				for _, instance := range privateNetwork.Instances {
					if _, err := r.ContaboClient.UnassignInstancePrivateNetwork(ctx, contaboCluster.Status.PrivateNetwork.PrivateNetworkId, instance.InstanceId, nil); err != nil {
						log.Error(err, "Failed to unassign instance from private network, continuing with deletion", "instanceID", instance.InstanceId, "privateNetworkId", privateNetwork.PrivateNetworkId)
					}
					log.Info("Unassigned instance from private network", "instanceID", instance.InstanceId, "privateNetworkId", privateNetwork.PrivateNetworkId)
					// Restart instance to apply network changes
					if _, err := r.ContaboClient.Restart(ctx, instance.InstanceId, nil); err != nil {
						log.Error(err, "Failed to restart instance after unassigning from private network, continuing with deletion", "instanceID", instance.InstanceId, "privateNetworkId", privateNetwork.PrivateNetworkId, "error")
					}
					log.Info("Restarted instance after unassigning from private network", "instanceID", instance.InstanceId, "privateNetworkId", privateNetwork.PrivateNetworkId)
				}
			}

			// Delete private network
			if _, err := r.ContaboClient.DeletePrivateNetwork(ctx, privateNetwork.PrivateNetworkId, nil); err != nil {
				log.Error(err, "Failed to delete private network, requeuing", "privateNetworkId", privateNetwork.PrivateNetworkId)
			}

			// Update status to remove private network
			contaboCluster.Status.PrivateNetwork = nil
			// This fail with conflict
			// meta.RemoveStatusCondition(&contaboCluster.Status.Conditions, infrastructurev1beta2.ClusterPrivateNetworkReadyCondition)
			log.Info("Deleted private network", "privateNetworkId", privateNetwork.PrivateNetworkId)
		}
	}

	// Delete Ssh Key
	if contaboCluster.Status.SshKey != nil {
		meta.SetStatusCondition(&contaboCluster.Status.Conditions, metav1.Condition{
			Type:   infrastructurev1beta2.ClusterSshKeyReadyCondition,
			Status: metav1.ConditionFalse,
			Reason: infrastructurev1beta2.ClusterSshKeyDeletingReason,
		})
		log.Info("Deleting SSH key", "sshKeyID", contaboCluster.Status.SshKey.SecretId)

		// Check if SSH key exists in Contabo API
		resp, err := r.ContaboClient.RetrieveSecretWithResponse(ctx, contaboCluster.Status.SshKey.SecretId, nil)
		if err != nil || resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
			// If the SSH key is not found, we can assume it has already been deleted
			log.Info("SSH key not found in Contabo API, assuming already deleted", "sshKeyID", contaboCluster.Status.SshKey.SecretId)
		} else {
			// Delete SSH key
			if _, err := r.ContaboClient.DeleteSecret(ctx, contaboCluster.Status.SshKey.SecretId, nil); err != nil {
				log.Error(err, "Failed to delete SSH key, requeuing", "sshKeyID", contaboCluster.Status.SshKey.SecretId)
			}
		}

		// Update status to remove SSH key
		sshKeyName := contaboCluster.Status.SshKey.Name
		contaboCluster.Status.SshKey = nil
		// This fail with conflict
		// meta.RemoveStatusCondition(&contaboCluster.Status.Conditions, infrastructurev1beta2.ClusterSshKeyReadyCondition)
		log.Info("Deleted SSH key", "name", sshKeyName)
	}

	// Remove our finalizer from the list and update it if there is no more contabomachines
	// 1. Get all contabomachines
	contaboMachineList := &infrastructurev1beta2.ContaboMachineList{}
	if err := r.List(ctx, contaboMachineList, client.InNamespace(contaboCluster.Namespace), client.MatchingLabels{
		clusterv1.ClusterNameLabel: contaboCluster.Name,
	}); err != nil {
		log.Error(err, "Failed to list ContaboMachines, continuing with deletion")
	}

	// 2. If there are still contabomachines, requeue the deletion
	if len(contaboMachineList.Items) > 0 {
		err := fmt.Errorf("there are still %d ContaboMachines in the cluster", len(contaboMachineList.Items))
		log.Error(err, "There are still ContaboMachines in the cluster, requeuing deletion", "contaboMachines", len(contaboMachineList.Items))
		return ctrl.Result{RequeueAfter: 10 * time.Second}, err
	}

	// Rmeove controlplane service and endpointslices
	if err := r.deleteControlPlaneService(ctx, contaboCluster); err != nil {
		log.Error(err, "Failed to delete control plane endpoint service, continuing with deletion")
	}

	if err := r.deleteControlPlaneEndpointSlices(ctx, contaboCluster); err != nil {
		log.Error(err, "Failed to delete control plane endpoint slices, continuing with deletion")
	}

	// 3. If there are no more contabomachines, remove the finalizer
	log.Info("No more ContaboMachines in the cluster, removing finalizer")
	controllerutil.RemoveFinalizer(contaboCluster, infrastructurev1beta2.ClusterFinalizer)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ContaboClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrastructurev1beta2.ContaboCluster{}).
		WithOptions(controller.Options{MaxConcurrentReconciles: 10}).
		WithEventFilter(predicates.ResourceNotPaused(mgr.GetScheme(), ctrl.LoggerFrom(context.TODO()))).
		Watches(
			&clusterv1.Cluster{},
			handler.EnqueueRequestsFromMapFunc(util.ClusterToInfrastructureMapFunc(context.TODO(), infrastructurev1beta2.GroupVersion.WithKind("ContaboCluster"), mgr.GetClient(), &infrastructurev1beta2.ContaboCluster{})),
			builder.WithPredicates(predicates.ClusterUnpaused(mgr.GetScheme(), ctrl.LoggerFrom(context.TODO()))),
		).
		Named("contabocluster").
		Complete(r)
}

// generateSSHKeyPair generates a new RSA SSH key pair and returns the public and private keys as strings
func generateSSHKeyPair() (string, string, error) {
	// generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}

	// write private key as PEM
	var privKeyBuf strings.Builder

	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	if err := pem.Encode(&privKeyBuf, privateKeyPEM); err != nil {
		return "", "", err
	}

	// generate and write public key
	pub, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", err
	}

	var pubKeyBuf strings.Builder
	pubKeyBuf.Write(ssh.MarshalAuthorizedKey(pub))

	return privKeyBuf.String(), pubKeyBuf.String(), nil
}

// handleError centralizes error handling with status condition, logging, event recording, and patching
func (r *ContaboClusterReconciler) handleError(ctx context.Context, contaboCluster *infrastructurev1beta2.ContaboCluster, err error, conditionType string, reason string, message string) error {
	log := logf.FromContext(ctx)

	// Set status condition
	meta.SetStatusCondition(&contaboCluster.Status.Conditions, metav1.Condition{
		Type:    conditionType,
		Status:  metav1.ConditionFalse,
		Reason:  reason,
		Message: message,
	})

	// Log the error
	log.Error(err, message)

	// Record event
	// r.Recorder.Event(contaboCluster, corev1.EventTypeWarning, reason, message)

	return fmt.Errorf("%s: %w", message, err)
}
