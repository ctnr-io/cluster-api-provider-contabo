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
	"github.com/google/uuid"
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
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
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
		return r.reconcileDelete(ctx, contaboCluster), nil
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
	if result, err := r.reconcilePrivateNetwork(ctx, contaboCluster); err != nil || result.RequeueAfter != 0 {
		return ctrl.Result{}, err
	}

	// Check if SSH key was created
	if result, err := r.reconcileSSHKey(ctx, contaboCluster); err != nil || result.RequeueAfter != 0 {
		return result, err
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

	if contaboCluster.Spec.ClusterUUID == "" {
		clusterUUID := uuid.New().String()
		contaboCluster.Spec.ClusterUUID = clusterUUID

		// Set initial creating state
		meta.SetStatusCondition(&contaboCluster.Status.Conditions, metav1.Condition{
			Type:   infrastructurev1beta2.ClusterReadyCondition,
			Status: metav1.ConditionFalse,
			Reason: infrastructurev1beta2.ClusterCreatingReason,
		})
		log.Info("Assigned new cluster UUID", "clusterUUID", clusterUUID)
		return clusterUUID
	}

	clusterUUID := contaboCluster.Spec.ClusterUUID

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

// reconcileDelete may return different ctrl.Result values in future implementations
func (r *ContaboClusterReconciler) reconcileDelete(ctx context.Context, contaboCluster *infrastructurev1beta2.ContaboCluster) ctrl.Result {
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
		sshKeyContaboName := contaboCluster.Status.SshKey.Name
		contaboCluster.Status.SshKey = nil
		// This fail with conflict
		// meta.RemoveStatusCondition(&contaboCluster.Status.Conditions, infrastructurev1beta2.ClusterSshKeyReadyCondition)
		log.Info("Deleted SSH key", "name", sshKeyContaboName)
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
		return ctrl.Result{RequeueAfter: 10 * time.Second}
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

	return ctrl.Result{}
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
