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

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/cluster-api/util/predicates"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	infrastructurev1beta1 "github.com/ctnr-io/cluster-api-provider-contabo/api/v1beta1"
	"github.com/ctnr-io/cluster-api-provider-contabo/pkg/cloud"
)

// ContaboClusterReconciler reconciles a ContaboCluster object
type ContaboClusterReconciler struct {
	client.Client
	Scheme         *runtime.Scheme
	Recorder       record.EventRecorder
	ContaboService *cloud.ContaboService
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
	contaboCluster := &infrastructurev1beta1.ContaboCluster{}
	if err := r.Get(ctx, req.NamespacedName, contaboCluster); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Fetch the Cluster
	cluster, err := util.GetOwnerCluster(ctx, r.Client, contaboCluster.ObjectMeta)
	if err != nil {
		return ctrl.Result{}, err
	}
	if cluster == nil {
		log.Info("Cluster Controller has not yet set OwnerRef")
		return ctrl.Result{}, nil
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
	return r.reconcileNormal(ctx, cluster, contaboCluster)
}

func (r *ContaboClusterReconciler) reconcileNormal(ctx context.Context, cluster *clusterv1.Cluster, contaboCluster *infrastructurev1beta1.ContaboCluster) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// If the ContaboCluster doesn't have our finalizer, add it.
	controllerutil.AddFinalizer(contaboCluster, infrastructurev1beta1.ClusterFinalizer)

	// Set the cluster in a progressing state
	conditions.MarkFalse(contaboCluster, infrastructurev1beta1.ReadyCondition, infrastructurev1beta1.CreatingReason, clusterv1.ConditionSeverityInfo, "")

	// Reconcile network infrastructure
	if err := r.reconcileNetwork(ctx, cluster, contaboCluster); err != nil {
		log.Error(err, "failed to reconcile network")
		conditions.MarkFalse(contaboCluster, infrastructurev1beta1.ReadyCondition, infrastructurev1beta1.NetworkInfrastructureFailedReason, clusterv1.ConditionSeverityError, "Failed to reconcile network: %s", err.Error())
		return ctrl.Result{}, err
	}

	// Set up the control plane endpoint if not already set
	if contaboCluster.Spec.ControlPlaneEndpoint.Host == "" {
		log.Info("Control plane endpoint not set, will be set by the first control plane machine")
		// The control plane endpoint will be set by the first control plane machine
		// For now, we'll mark the cluster as not ready and return
		conditions.MarkFalse(contaboCluster, infrastructurev1beta1.ReadyCondition, infrastructurev1beta1.WaitingForControlPlaneEndpointReason, clusterv1.ConditionSeverityInfo, "")
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	// Mark the cluster as ready
	contaboCluster.Status.Ready = true
	conditions.MarkTrue(contaboCluster, infrastructurev1beta1.ReadyCondition)

	return ctrl.Result{}, nil
}

//nolint:unparam // reconcileDelete may return different ctrl.Result values in future implementations
func (r *ContaboClusterReconciler) reconcileDelete(ctx context.Context, cluster *clusterv1.Cluster, contaboCluster *infrastructurev1beta1.ContaboCluster) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	log.Info("Reconciling ContaboCluster delete")

	// Delete network infrastructure
	if err := r.deleteNetwork(ctx, cluster, contaboCluster); err != nil {
		log.Error(err, "failed to delete network")
		return ctrl.Result{}, err
	}

	// Remove our finalizer from the list and update it.
	controllerutil.RemoveFinalizer(contaboCluster, infrastructurev1beta1.ClusterFinalizer)

	return ctrl.Result{}, nil
}

//nolint:unparam // reconcileNetwork may return different errors in future implementations
func (r *ContaboClusterReconciler) reconcileNetwork(ctx context.Context, cluster *clusterv1.Cluster, contaboCluster *infrastructurev1beta1.ContaboCluster) error {
	_ = cluster // may be used in future for cluster-specific network configuration
	log := logf.FromContext(ctx)

	// For now, we'll just initialize the network status
	// In a real implementation, you would create VPCs, subnets, etc.
	if contaboCluster.Status.Network == nil {
		contaboCluster.Status.Network = &infrastructurev1beta1.ContaboNetworkStatus{}
	}

	log.Info("Network infrastructure reconciled successfully")
	return nil
}

func (r *ContaboClusterReconciler) deleteNetwork(ctx context.Context, cluster *clusterv1.Cluster, contaboCluster *infrastructurev1beta1.ContaboCluster) error {
	log := logf.FromContext(ctx)

	// For now, we'll just log the deletion
	// In a real implementation, you would delete VPCs, subnets, etc.
	log.Info("Network infrastructure deleted successfully")
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ContaboClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrastructurev1beta1.ContaboCluster{}).
		WithOptions(controller.Options{MaxConcurrentReconciles: 10}).
		WithEventFilter(predicates.ResourceNotPaused(ctrl.LoggerFrom(context.TODO()))).
		Watches(
			&clusterv1.Cluster{},
			handler.EnqueueRequestsFromMapFunc(util.ClusterToInfrastructureMapFunc(context.TODO(), infrastructurev1beta1.GroupVersion.WithKind("ContaboCluster"), mgr.GetClient(), &infrastructurev1beta1.ContaboCluster{})),
			builder.WithPredicates(predicates.ClusterUnpaused(ctrl.LoggerFrom(context.TODO()))),
		).
		Named("contabocluster").
		Complete(r)
}
