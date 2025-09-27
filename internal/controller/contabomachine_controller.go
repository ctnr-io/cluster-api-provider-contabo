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
	"github.com/ctnr-io/cluster-api-provider-contabo/pkg/contabo/v1.0.0/models"
)

// ContaboMachineReconciler reconciles a ContaboMachine object
type ContaboMachineReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	Recorder      record.EventRecorder
	ContaboClient *contaboclient.Client
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
	patchHelper, err := patch.NewHelper(contaboMachine, r.Client)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Always attempt to Patch the ContaboMachine object and status after each reconciliation
	defer func() {
		if err := patchHelper.Patch(ctx, contaboMachine); err != nil {
			log.Error(err, "failed to patch ContaboMachine")
		}
	}()

	// Handle deleted machines
	if !contaboMachine.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, machine, contaboMachine, contaboCluster)
	}

	// Handle non-deleted machines
	return r.reconcileNormal(ctx, machine, contaboMachine, contaboCluster)
}

func (r *ContaboMachineReconciler) reconcileNormal(ctx context.Context, machine *clusterv1.Machine, contaboMachine *infrastructurev1beta2.ContaboMachine, contaboCluster *infrastructurev1beta2.ContaboCluster) (ctrl.Result, error) {
	// Initialize basic machine setup
	if err := r.initializeMachine(ctx, machine, contaboMachine); err != nil {
		return ctrl.Result{}, err
	}

	// State Machine: Process each component in dependency order
	// 1. Check cluster infrastructure readiness
	if result, err := r.reconcileClusterInfrastructure(ctx, contaboCluster, contaboMachine); err != nil || result.Requeue || result.RequeueAfter > 0 {
		return result, err
	}

	// 2. Reconcile secrets (needed for instance creation)
	if result, err := r.reconcileMachineSecrets(ctx, contaboMachine, contaboCluster); err != nil || result.Requeue || result.RequeueAfter > 0 {
		return result, err
	}

	// 3. Reconcile private networks (needed before attaching to instance)
	if result, err := r.reconcileMachinePrivateNetworks(ctx, contaboMachine, contaboCluster); err != nil || result.Requeue || result.RequeueAfter > 0 {
		return result, err
	}

	// 6. Reconcile instance (depends on secrets, networks, and SSH keys)
	if result, err := r.reconcileInstance(ctx, machine, contaboMachine, contaboCluster); err != nil || result.Requeue || result.RequeueAfter > 0 {
		return result, err
	}

	// Final state: Set machine as ready if all components are ready
	return r.reconcileMachineReady(ctx, contaboMachine)
}

// initializeMachine sets up basic machine metadata and finalizers
func (r *ContaboMachineReconciler) initializeMachine(ctx context.Context, machine *clusterv1.Machine, contaboMachine *infrastructurev1beta2.ContaboMachine) error {
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

	// Initialize status conditions if needed
	if contaboMachine.Status.Conditions == nil {
		contaboMachine.Status.Conditions = []metav1.Condition{}
	}

	return nil
}

// reconcileClusterInfrastructure checks cluster infrastructure dependencies
func (r *ContaboMachineReconciler) reconcileClusterInfrastructure(ctx context.Context, contaboCluster *infrastructurev1beta2.ContaboCluster, contaboMachine *infrastructurev1beta2.ContaboMachine) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

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

	return ctrl.Result{}, nil
}

// reconcileMachineReady sets the overall machine readiness based on all component states
func (r *ContaboMachineReconciler) reconcileMachineReady(ctx context.Context, contaboMachine *infrastructurev1beta2.ContaboMachine) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// Check if all required conditions are ready
	clusterReady := meta.IsStatusConditionTrue(contaboMachine.Status.Conditions, infrastructurev1beta2.ClusterInfrastructureReadyCondition)
	secretsReady := meta.IsStatusConditionTrue(contaboMachine.Status.Conditions, infrastructurev1beta2.MachineSecretsReadyCondition)
	networksReady := meta.IsStatusConditionTrue(contaboMachine.Status.Conditions, infrastructurev1beta2.MachinePrivateNetworksReadyCondition)
	instanceReady := meta.IsStatusConditionTrue(contaboMachine.Status.Conditions, infrastructurev1beta2.InstanceReadyCondition)

	if clusterReady && secretsReady && networksReady && instanceReady {
		// All components are ready
		contaboMachine.Status.Ready = true
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:   infrastructurev1beta2.MachineReadyCondition,
			Status: metav1.ConditionTrue,
			Reason: infrastructurev1beta2.AvailableReason,
		})
		log.Info("ContaboMachine is ready")
	} else {
		// Still waiting for some components
		contaboMachine.Status.Ready = false
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:   infrastructurev1beta2.MachineReadyCondition,
			Status: metav1.ConditionFalse,
			Reason: infrastructurev1beta2.CreatingReason,
		})
		log.Info("ContaboMachine not ready yet",
			"clusterReady", clusterReady,
			"secretsReady", secretsReady,
			"networksReady", networksReady,
			"instanceReady", instanceReady)
	}

	return ctrl.Result{}, nil
}

func (r *ContaboMachineReconciler) reconcileDelete(ctx context.Context, machine *clusterv1.Machine, contaboMachine *infrastructurev1beta2.ContaboMachine, contaboCluster *infrastructurev1beta2.ContaboCluster) (ctrl.Result, error) {
	_ = contaboCluster // may be used in future for cluster-specific cleanup logic
	log := logf.FromContext(ctx)

	log.Info("Reconciling ContaboMachine delete - setting instance available for reuse")

	// Set instance back to available state for reuse
	instanceID, err := ParseProviderID(contaboMachine.Spec.ProviderID)
	if err == nil {
		if err := r.removeClusterBinding(ctx, instanceID); err != nil {
			log.Error(err, "failed to set instance state to available")
			// Don't fail the whole operation for displayName issues
		}
	}

	// Note: We don't actually delete/cancel the instance, just mark it available
	// The instance remains available for reuse with the "<id>-avl" displayName format
	log.Info("Instance state updated to available, instance ready for reuse")

	// Remove our finalizer from the list and update it.
	controllerutil.RemoveFinalizer(contaboMachine, infrastructurev1beta2.MachineFinalizer)

	return ctrl.Result{}, nil
}

// removeClusterBinding removes cluster binding and sets instance to available state for reuse
// The instance itself is preserved and made available for reuse rather than being deleted/cancelled
func (r *ContaboMachineReconciler) removeClusterBinding(ctx context.Context, instanceID int64) error {
	log := logf.FromContext(ctx)

	log.Info("Setting instance state to available for reuse", "instanceId", instanceID)

	// Create new display name with available state and no cluster ID
	newDisplayName := BuildInstanceDisplayNameWithState(instanceID, StateAvailable, "")

	// Update instance display name via Contabo API
	updateReq := models.PatchInstanceRequest{
		DisplayName: &newDisplayName,
	}

	resp, err := r.ContaboClient.PatchInstance(ctx, instanceID, &models.PatchInstanceParams{}, updateReq)
	if err != nil {
		return fmt.Errorf("failed to call patch instance API: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Error(closeErr, "failed to close patch instance response body")
		}
	}()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to update instance display name, status: %d", resp.StatusCode)
	}

	log.Info("Successfully updated instance display name", "instanceId", instanceID, "newDisplayName", newDisplayName)
	return nil
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
