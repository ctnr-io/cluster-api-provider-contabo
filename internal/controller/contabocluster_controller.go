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
)

// ContaboClusterReconciler reconciles a ContaboCluster object
type ContaboClusterReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	Recorder      record.EventRecorder
	ContaboClient *contaboclient.Client
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

	log = log.WithValues("cluster", cluster.Name)

	// Initialize the patch helper
	patchHelper, err := patch.NewHelper(contaboCluster, r.Client)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Always attempt to Patch the ContaboCluster object and status after each reconciliation
	defer func() {
		if err := patchHelper.Patch(ctx, contaboCluster); err != nil {
			log.Error(err, "failed to patch ContaboCluster")
		}
	}()

	// Handle deleted clusters
	if !contaboCluster.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, cluster, contaboCluster)
	}

	// Handle non-deleted clusters
	return r.reconcileNormal(ctx, contaboCluster)
}

func (r *ContaboClusterReconciler) reconcileNormal(ctx context.Context, contaboCluster *infrastructurev1beta2.ContaboCluster) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	log.Info("Starting ContaboCluster reconciliation state machine", "cluster", contaboCluster.Name)

	// Initialize basic cluster setup
	if err := r.initializeCluster(ctx, contaboCluster); err != nil {
		return ctrl.Result{}, err
	}

	// State Machine: Process each component in dependency order
	// 1. Reconcile secrets (needed for other components)
	if result, err := r.reconcileClusterSecrets(ctx, contaboCluster); err != nil || result.Requeue || result.RequeueAfter > 0 {
		return result, err
	}

	// 2. Reconcile private networks (needed before control plane setup)
	if result, err := r.reconcileClusterPrivateNetworks(ctx, contaboCluster); err != nil || result.Requeue || result.RequeueAfter > 0 {
		return result, err
	}

	// 3. Reconcile control plane endpoint (depends on networks)
	if result, err := r.reconcileControlPlaneEndpoint(ctx, contaboCluster); err != nil || result.Requeue || result.RequeueAfter > 0 {
		return result, err
	}

	// Final state: Set cluster as ready if all components are ready
	return r.reconcileClusterReady(ctx, contaboCluster)
}

// initializeCluster sets up basic cluster metadata and finalizers
func (r *ContaboClusterReconciler) initializeCluster(ctx context.Context, contaboCluster *infrastructurev1beta2.ContaboCluster) error {
	log := logf.FromContext(ctx)

	// Add finalizer
	controllerutil.AddFinalizer(contaboCluster, infrastructurev1beta2.ClusterFinalizer)

	// Initialize status conditions if needed
	if contaboCluster.Status.Conditions == nil {
		contaboCluster.Status.Conditions = []metav1.Condition{}
	}

	// Ensure cluster has a unique UUID for global identification
	clusterUUID := EnsureClusterUUID(contaboCluster)
	log.Info("Cluster UUID ensured", "uuid", clusterUUID)

	// Set initial creating state
	meta.SetStatusCondition(&contaboCluster.Status.Conditions, metav1.Condition{
		Type:   infrastructurev1beta2.ReadyCondition,
		Status: metav1.ConditionFalse,
		Reason: infrastructurev1beta2.CreatingReason,
	})

	return nil
}

// reconcileClusterReady sets the overall cluster readiness based on all component states
func (r *ContaboClusterReconciler) reconcileClusterReady(ctx context.Context, contaboCluster *infrastructurev1beta2.ContaboCluster) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// Check if all required conditions are ready
	secretsReady := meta.IsStatusConditionTrue(contaboCluster.Status.Conditions, infrastructurev1beta2.ClusterSecretsReadyCondition)
	networksReady := meta.IsStatusConditionTrue(contaboCluster.Status.Conditions, infrastructurev1beta2.ClusterPrivateNetworkReadyCondition)
	endpointReady := meta.IsStatusConditionTrue(contaboCluster.Status.Conditions, infrastructurev1beta2.ControlPlaneEndpointReadyCondition)

	if secretsReady && networksReady && endpointReady {
		// All components are ready
		contaboCluster.Status.Ready = true
		meta.SetStatusCondition(&contaboCluster.Status.Conditions, metav1.Condition{
			Type:   infrastructurev1beta2.ReadyCondition,
			Status: metav1.ConditionTrue,
			Reason: infrastructurev1beta2.AvailableReason,
		})
		log.Info("ContaboCluster is ready")
	} else {
		// Still waiting for some components
		contaboCluster.Status.Ready = false
		meta.SetStatusCondition(&contaboCluster.Status.Conditions, metav1.Condition{
			Type:   infrastructurev1beta2.ReadyCondition,
			Status: metav1.ConditionFalse,
			Reason: "ComponentsNotReady",
		})
		log.Info("ContaboCluster not ready yet",
			"secretsReady", secretsReady,
			"networksReady", networksReady,
			"endpointReady", endpointReady)
	}

	return ctrl.Result{}, nil
}

//nolint:unparam // reconcileDelete may return different ctrl.Result values in future implementations
func (r *ContaboClusterReconciler) reconcileDelete(ctx context.Context, cluster *clusterv1.Cluster, contaboCluster *infrastructurev1beta2.ContaboCluster) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	log.Info("Reconciling ContaboCluster delete")

	// Delete network infrastructure
	if err := r.deleteNetwork(ctx, contaboCluster); err != nil {
		log.Error(err, "failed to delete network")
		return ctrl.Result{}, err
	}

	// Remove our finalizer from the list and update it.
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
