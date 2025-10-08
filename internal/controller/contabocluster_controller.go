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

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *ContaboClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

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
		return ctrl.Result{}, r.reconcileDelete(ctx, contaboCluster)
	}

	// Handle non-deleted clusters
	return r.reconcileApply(ctx, contaboCluster)
}

func (r *ContaboClusterReconciler) reconcileApply(ctx context.Context, contaboCluster *infrastructurev1beta2.ContaboCluster) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	log.Info("Starting ContaboCluster reconciliation", "cluster", contaboCluster.Name)

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

	// Check if control plane endpoint is set
	if result := r.reconcileControlPlaneEndpoint(ctx, contaboCluster); result.RequeueAfter != 0 {
		return result, nil
	}

	// Mark cluster as ready if all components are ready
	r.markClusterReady(ctx, contaboCluster)

	return ctrl.Result{}, nil
}

// markClusterReady sets the ClusterReadyCondition to true if all components are ready
func (r *ContaboClusterReconciler) markClusterReady(ctx context.Context, contaboCluster *infrastructurev1beta2.ContaboCluster) {
	log := logf.FromContext(ctx)

	log.Info("All components are ready, marking ContaboCluster as ready")

	meta.SetStatusCondition(&contaboCluster.Status.Conditions, metav1.Condition{
		Type:   infrastructurev1beta2.ClusterReadyCondition,
		Status: metav1.ConditionTrue,
		Reason: infrastructurev1beta2.ClusterAvailableReason,
	})

	contaboCluster.Status.Ready = true

	contaboCluster.Status.Initialization = &infrastructurev1beta2.ContaboClusterInitializationStatus{
		Provisioned: true,
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
func (r *ContaboClusterReconciler) reconcileControlPlaneEndpoint(ctx context.Context, contaboCluster *infrastructurev1beta2.ContaboCluster) ctrl.Result {
	log := logf.FromContext(ctx)

	// TODO: Handle case where control plane endpoint is set but not reachable

	// Check if control plane endpoint is set
	if contaboCluster.Spec.ControlPlaneEndpoint != nil {
		return ctrl.Result{}
	}

	// Retrieve control plane endpoint from the first ContaboMachine in the cluster
	machineList := &infrastructurev1beta2.ContaboMachineList{}
	if err := r.List(ctx, machineList, client.InNamespace(contaboCluster.Namespace), client.MatchingLabels{
		clusterv1.ClusterNameLabel:         contaboCluster.Name,
		clusterv1.MachineControlPlaneLabel: "true",
	}); err != nil || len(machineList.Items) == 0 {
		log.Info("No control plane machines found yet, requeuing")
		meta.SetStatusCondition(&contaboCluster.Status.Conditions, metav1.Condition{
			Type:   clusterv1.ClusterControlPlaneAvailableCondition,
			Status: metav1.ConditionFalse,
			Reason: clusterv1.WaitingForControlPlaneInitializedReason,
		})
		return ctrl.Result{RequeueAfter: 10 * time.Second}
	}

	firstControlPlaneMachine := machineList.Items[0]
	if firstControlPlaneMachine.Status.Instance == nil || firstControlPlaneMachine.Status.Instance.IpConfig.V4.Ip == "" {
		log.Info("Control plane machine instance or IP not yet available, requeuing", "machine", firstControlPlaneMachine.Name)
		return ctrl.Result{RequeueAfter: 10 * time.Second}
	}

	// Set control plane endpoint
	contaboCluster.Spec.ControlPlaneEndpoint = &clusterv1.APIEndpoint{
		Host: firstControlPlaneMachine.Status.Instance.IpConfig.V4.Ip,
		Port: 6443, // Default Kubernetes API server port
	}

	meta.SetStatusCondition(&contaboCluster.Status.Conditions, metav1.Condition{
		Type:   clusterv1.ClusterControlPlaneAvailableCondition,
		Status: metav1.ConditionTrue,
		Reason: clusterv1.ClusterControlPlaneMachinesReadyReason,
	})
	log.Info("Using control plane endpoint", "controlPlaneEndpoint", contaboCluster.Spec.ControlPlaneEndpoint)

	return ctrl.Result{}
}

// reconcileDelete may return different ctrl.Result values in future implementations
func (r *ContaboClusterReconciler) reconcileDelete(ctx context.Context, contaboCluster *infrastructurev1beta2.ContaboCluster) error {
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
						return r.handleError(
							ctx,
							contaboCluster,
							err,
							infrastructurev1beta2.ClusterPrivateNetworkReadyCondition,
							infrastructurev1beta2.ClusterPrivateNetworkFailedReason,
							"Failed to unassign instances from private network",
						)
					}
					log.Info("Unassigned instance from private network", "instanceID", instance.InstanceId, "privateNetworkId", privateNetwork.PrivateNetworkId)
				}
			}

			// Delete private network
			if _, err := r.ContaboClient.DeletePrivateNetwork(ctx, privateNetwork.PrivateNetworkId, nil); err != nil {
				return r.handleError(
					ctx,
					contaboCluster,
					err,
					infrastructurev1beta2.ClusterPrivateNetworkReadyCondition,
					infrastructurev1beta2.ClusterPrivateNetworkFailedReason,
					"Failed to delete private network",
				)
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
				return r.handleError(
					ctx,
					contaboCluster,
					err,
					infrastructurev1beta2.ClusterSshKeyReadyCondition,
					infrastructurev1beta2.ClusterSshKeyFailedReason,
					"Failed to delete SSH key",
				)
			}
		}

		// Update status to remove SSH key
		sshKeyName := contaboCluster.Status.SshKey.Name
		contaboCluster.Status.SshKey = nil
		// This fail with conflict
		// meta.RemoveStatusCondition(&contaboCluster.Status.Conditions, infrastructurev1beta2.ClusterSshKeyReadyCondition)
		log.Info("Deleted SSH key", "name", sshKeyName)
	}

	// Remove our finalizer from the list and update it.
	controllerutil.RemoveFinalizer(contaboCluster, infrastructurev1beta2.ClusterFinalizer)

	return nil
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
