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
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/record"
	"k8s.io/utils/ptr"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/cluster-api/util/predicates"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	infrastructurev1beta2 "github.com/ctnr-io/cluster-api-provider-contabo/api/v1beta2"
	contaboclient "github.com/ctnr-io/cluster-api-provider-contabo/pkg/contabo/v1.0.0/client"
	"github.com/ctnr-io/cluster-api-provider-contabo/pkg/contabo/v1.0.0/models"

	corev1 "k8s.io/api/core/v1"

	_ "embed"
)

//go:embed templates/controlplane.cloud-config.yaml
var controlplaneCloudConfig string

//go:embed templates/worker.cloud-config.yaml
var workerCloudConfig string

// ContaboMachineReconciler reconciles a ContaboMachine object
type ContaboMachineReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	Recorder      record.EventRecorder
	ContaboClient *contaboclient.ClientWithResponses
	// instanceReuseMutex protects against concurrent instance reuse
	instanceReuseMutex sync.Mutex
	// indexAssignmentMutex protects against concurrent index assignment
	indexAssignmentMutex sync.Mutex
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=contabomachines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=contabomachines/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=contabomachines/finalizers,verbs=update
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machines;machines/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters;clusters/status,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=create;update;delete;get;list;watch

// SetupWithManager sets up the controller with the Manager.
func (r *ContaboMachineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrastructurev1beta2.ContaboMachine{}).
		WithOptions(controller.Options{MaxConcurrentReconciles: 10}).
		WithEventFilter(predicates.ResourceNotPaused(mgr.GetScheme(), ctrl.LoggerFrom(context.TODO()))).
		// Uncomment to reconcile based on Machine, currently this is not what we need
		// Watches(
		// 	&clusterv1.Machine{},
		// 	handler.EnqueueRequestsFromMapFunc(util.MachineToInfrastructureMapFunc(infrastructurev1beta2.GroupVersion.WithKind("ContaboMachine"))),
		// 	builder.WithPredicates(predicates.ResourceNotPaused(mgr.GetScheme(), ctrl.LoggerFrom(context.TODO()))),
		// ).
		Complete(r)
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *ContaboMachineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	log.Info("Reconciling ContaboMachine", "namespace", req.Namespace, "name", req.Name)

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
		return ctrl.Result{RequeueAfter: 15 * time.Second}, nil
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
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
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

	// Always apply a patch at the end of reconciliation to avoid losing status updates
	// Also prevent concurrency errors
	defer func() {
		// Retrieve latest ContaboMachine for patching
		latest := &infrastructurev1beta2.ContaboMachine{}
		if err := r.Get(ctx, req.NamespacedName, latest); err != nil {
			log.Error(err, "failed to get ContaboMachine for patching")
			return
		}

		// Create patch helper from latest
		patchHelper, err := patch.NewHelper(latest, r.Client)
		if err != nil {
			log.Error(err, "failed to create patch helper for ContaboMachine", "machine", latest.Name, "cluster", cluster.Name)
			return
		}

		// Update spec, status, labels
		latest.Spec = contaboMachine.Spec
		latest.Status = contaboMachine.Status
		for key, value := range contaboMachine.Labels {
			latest.Labels[key] = value
		}
		if controllerutil.ContainsFinalizer(contaboMachine, infrastructurev1beta2.MachineFinalizer) {
			controllerutil.AddFinalizer(latest, infrastructurev1beta2.MachineFinalizer)
		} else {
			controllerutil.RemoveFinalizer(latest, infrastructurev1beta2.MachineFinalizer)
		}
		if patchErr := patchHelper.Patch(ctx, latest); patchErr != nil {
			log.Error(patchErr, "failed to patch ContaboMachine", "machine", latest.Name, "cluster", cluster.Name)
			return
		}
	}()

	// Handle deleted machines
	if !contaboMachine.DeletionTimestamp.IsZero() {
		r.reconcileDelete(ctx, contaboMachine, contaboCluster)
		return ctrl.Result{}, nil
	}

	// Handle non-deleted machines
	result, err := r.reconcileNormal(ctx, machine, contaboMachine, contaboCluster)
	if err != nil {
		log.Error(err, "Reconciliation failed")
		return result, err
	}

	return result, err
}

// Reconcile Normal will set required spec & status for CAPI and CABPK to work properly, the other statuses is handled by other methods
func (r *ContaboMachineReconciler) reconcileNormal(ctx context.Context, machine *clusterv1.Machine, contaboMachine *infrastructurev1beta2.ContaboMachine, contaboCluster *infrastructurev1beta2.ContaboCluster) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// If contabomachine have failure reason/message, do not proceed
	if contaboMachine.Status.FailureReason != nil || contaboMachine.Status.FailureMessage != nil {
		log.Info("ContaboMachine is in failure state, skipping reconciliation",
			"failureReason", contaboMachine.Status.FailureReason,
			"failureMessage", contaboMachine.Status.FailureMessage,
		)
		return ctrl.Result{}, nil
	}

	// Setup the resource
	if result := r.setupContaboMachine(ctx, machine, contaboMachine, contaboCluster); result.RequeueAfter > 0 {
		return result, nil
	}

	contaboMachine.Status.Initialization = &infrastructurev1beta2.ContaboMachineInitializationStatus{
		Provisioned:  false,
		ErrorMessage: nil,
	}
	contaboMachine.Spec.ProviderID = nil
	contaboMachine.Status.Ready = false
	contaboMachine.Status.Available = false
	meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
		Type:   infrastructurev1beta2.InstanceReadyCondition,
		Status: metav1.ConditionFalse,
		Reason: "",
	})

	// Provision the instance for CAPI
	result, err := r.provisionInstance(ctx, contaboMachine, contaboCluster)
	if err != nil {
		contaboMachine.Status.Initialization = &infrastructurev1beta2.ContaboMachineInitializationStatus{
			Provisioned:  false,
			ErrorMessage: ptr.To(err.Error()),
		}
		return result, err
	}
	if result.RequeueAfter > 0 {
		contaboMachine.Status.Initialization = &infrastructurev1beta2.ContaboMachineInitializationStatus{
			Provisioned:  false,
			ErrorMessage: nil,
		}
		return result, nil
	}
	contaboMachine.Status.Initialization = &infrastructurev1beta2.ContaboMachineInitializationStatus{
		Provisioned: true,
	}

	// Set provider ID for CAPI
	contaboMachine.Spec.ProviderID = ptr.To(BuildProviderID(contaboMachine.Status.Instance.Name))

	// Update machine address for CAPI machine
	if len(contaboMachine.Status.Addresses) == 0 {
		if err := r.reconcileContaboMachineAddresses(ctx, contaboMachine, contaboCluster); err != nil {
			return ctrl.Result{}, err
		}
	}

	// Set ContaboMachine ready for bootstrap data to be available by CABPK
	contaboMachine.Status.Ready = true

	// Bootstrap the instance if not already done
	// Validate instance status
	result, err = r.bootstrapInstance(ctx, machine, contaboMachine, contaboCluster)
	if err != nil {
		return result, err
	}
	if result.RequeueAfter > 0 {
		return result, nil
	}

	contaboMachine.Status.Available = true

	// Update ContaboMachine status with instance details
	meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
		Type:   infrastructurev1beta2.InstanceReadyCondition,
		Status: metav1.ConditionTrue,
		Reason: infrastructurev1beta2.InstanceReadyReason,
	})
	log.Info("Instance is ready",
		"instanceID", contaboMachine.Status.Instance.InstanceId,
		"instanceIP", contaboMachine.Status.Instance.IpConfig.V4.Ip)

	// ContaboMachine is ready
	return ctrl.Result{}, nil
}

// setupContaboMachine handles initial machine setup and validation
func (r *ContaboMachineReconciler) setupContaboMachine(ctx context.Context, machine *clusterv1.Machine, contaboMachine *infrastructurev1beta2.ContaboMachine, contaboCluster *infrastructurev1beta2.ContaboCluster) ctrl.Result {
	log := logf.FromContext(ctx)

	if contaboMachine.Status.Instance == nil {
		// Add finalizer (currently there to prevent update it when patching but we need to refactor)
		controllerutil.AddFinalizer(contaboMachine, infrastructurev1beta2.MachineFinalizer)
	}

	// Automatically add cluster label to ContaboMachine for proper mapping
	if contaboMachine.Labels == nil {
		contaboMachine.Labels = make(map[string]string)
	}
	// Set machine labels to contaboMachine
	for key, value := range machine.Labels {
		contaboMachine.Labels[key] = value
	}
	// Set cluster name label to contaboMachine
	if machine.Spec.ClusterName != "" {
		contaboMachine.Labels[clusterv1.ClusterNameLabel] = machine.Spec.ClusterName
		log.V(4).Info("Added cluster label to ContaboMachine",
			"clusterName", machine.Spec.ClusterName,
			"label", clusterv1.ClusterNameLabel)
	}

	// Assign a unique index to this machine within the cluster
	if err := r.assignMachineIndex(ctx, contaboMachine, machine.Spec.ClusterName); err != nil {
		log.Error(err, "Failed to assign machine index")
		return ctrl.Result{RequeueAfter: 15 * time.Second}
	}

	// Check if cluster private network is ready
	if contaboCluster.Status.PrivateNetwork == nil || meta.IsStatusConditionFalse(contaboCluster.Status.Conditions, infrastructurev1beta2.ClusterPrivateNetworkReadyCondition) {
		log.Info("Waiting for cluster private network to be ready")
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:   infrastructurev1beta2.ClusterPrivateNetworkReadyCondition,
			Status: metav1.ConditionFalse,
			Reason: infrastructurev1beta2.InstanceWaitingForPrivateNetworksReason,
		})
		return ctrl.Result{RequeueAfter: 30 * time.Second}
	}
	meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
		Type:   infrastructurev1beta2.ClusterPrivateNetworkReadyCondition,
		Status: metav1.ConditionTrue,
		Reason: infrastructurev1beta2.ClusterPrivateNetworkReadyCondition,
	})

	// Check if cluster ssh keys is ready
	if contaboCluster.Status.SshKey == nil || meta.IsStatusConditionFalse(contaboCluster.Status.Conditions, infrastructurev1beta2.ClusterSshKeyReadyCondition) {
		log.Info("Waiting for cluster ssh key to be ready")
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:   infrastructurev1beta2.ClusterSshKeyReadyCondition,
			Status: metav1.ConditionFalse,
			Reason: infrastructurev1beta2.InstanceWaitingForSshKeyReason,
		})
		return ctrl.Result{RequeueAfter: 30 * time.Second}
	}
	meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
		Type:   infrastructurev1beta2.ClusterSshKeyReadyCondition,
		Status: metav1.ConditionTrue,
		Reason: infrastructurev1beta2.ClusterSshKeyReadyReason,
	})

	return ctrl.Result{}
}

// getBootstrapData retrieves and validates bootstrap data
func (r *ContaboMachineReconciler) getBootstrapData(ctx context.Context, machine *clusterv1.Machine, contaboMachine *infrastructurev1beta2.ContaboMachine, contaboCluster *infrastructurev1beta2.ContaboCluster) (string, ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// Get bootstrap data secret
	if machine.Spec.Bootstrap.DataSecretName == nil {
		log.Info("Bootstrap data secret is not available yet")
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:   infrastructurev1beta2.BootstrapDataAvailableCondition,
			Status: metav1.ConditionFalse,
			Reason: infrastructurev1beta2.WaitingForBootstrapDataReason,
		})
		return "", ctrl.Result{RequeueAfter: 15 * time.Second}, nil
	}

	// Get bootstrap data
	bootstrapDataSecret := &corev1.Secret{}
	bootstrapDataSecretName := client.ObjectKey{
		Namespace: contaboMachine.Namespace,
		Name:      *machine.Spec.Bootstrap.DataSecretName,
	}
	if err := r.Get(ctx, bootstrapDataSecretName, bootstrapDataSecret); err != nil {
		return "", ctrl.Result{}, r.handleError(
			ctx,
			contaboMachine,
			err,
			infrastructurev1beta2.WaitingForBootstrapDataReason,
			"Failed to get bootstrap data secret",
		)
	}

	if _, ok := bootstrapDataSecret.Data["value"]; !ok || len(bootstrapDataSecret.Data["value"]) == 0 {
		log.Info("Bootstrap data secret is missing 'value' key or is empty")
		meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
			Type:   infrastructurev1beta2.InstanceReadyCondition,
			Status: metav1.ConditionFalse,
			Reason: infrastructurev1beta2.WaitingForBootstrapDataReason,
		})
		return "", ctrl.Result{RequeueAfter: 15 * time.Second}, nil
	}

	cloudConfig := workerCloudConfig
	if contaboMachine.Labels[clusterv1.MachineControlPlaneLabel] != "false" {
		cloudConfig = controlplaneCloudConfig
	}

	// Merge cloud-config with bootstrap data
	mergedConfig, err := mergeCloudConfig([]byte(cloudConfig), bootstrapDataSecret.Data["value"])
	if err != nil {
		return "", ctrl.Result{}, r.handleError(
			ctx,
			contaboMachine,
			err,
			infrastructurev1beta2.BootstrapDataMergeFailedReason,
			"Failed to merge cloud-config with bootstrap data",
		)
	}

	// Replace template variables in cloud-config
	mergedConfigStr := string(mergedConfig)
	kubadmVersion := strings.Join(strings.Split(machine.Spec.Version, ".")[:2], ".")
	internalIpV4CIDR := contaboCluster.Status.PrivateNetwork.Cidr
	internalIpV4 := ""
	for _, address := range contaboMachine.Status.Addresses {
		if address.Type == clusterv1.MachineInternalIP {
			internalIpV4 = address.Address
			break
		}
	}
	if internalIpV4 == "" {
		return "", ctrl.Result{}, fmt.Errorf("failed to find internal ipv4 address for instance %d", contaboMachine.Status.Instance.InstanceId)
	}

	// Replace all variables
	mergedConfigStr = strings.ReplaceAll(mergedConfigStr, "${KUBEADM_VERSION}", kubadmVersion)
	mergedConfigStr = strings.ReplaceAll(mergedConfigStr, "${INTERNAL_IPV4_CIDR}", internalIpV4CIDR)
	mergedConfigStr = strings.ReplaceAll(mergedConfigStr, "${INTERNAL_IPV4}", net.ParseIP(internalIpV4).String())
	mergedConfigStr = strings.ReplaceAll(mergedConfigStr, "${EXTERNAL_IPV4}", net.ParseIP(contaboMachine.Status.Instance.IpConfig.V4.Ip).String())
	mergedConfigStr = strings.ReplaceAll(mergedConfigStr, "${EXTERNAL_IPV6}", net.ParseIP(contaboMachine.Status.Instance.IpConfig.V6.Ip).String())
	mergedConfigStr = strings.ReplaceAll(mergedConfigStr, "${PROVIDER_ID}", *contaboMachine.Spec.ProviderID)
	mergedConfigStr = strings.ReplaceAll(mergedConfigStr, "${CLUSTER_UUID}", contaboCluster.Spec.ClusterUUID)

	meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
		Type:   infrastructurev1beta2.BootstrapDataAvailableCondition,
		Status: metav1.ConditionTrue,
		Reason: infrastructurev1beta2.BootstrapDataAvailableReason,
	})
	log.Info("Successfully retrieved and merged bootstrap data")

	return mergedConfigStr, ctrl.Result{}, nil
}

// reconcileInstance handles creation and setup of a new instance
func (r *ContaboMachineReconciler) provisionInstance(ctx context.Context, contaboMachine *infrastructurev1beta2.ContaboMachine, contaboCluster *infrastructurev1beta2.ContaboCluster) (ctrl.Result, error) {
	// Check if SSH key is available before accessing it
	if contaboCluster.Status.SshKey == nil {
		return ctrl.Result{RequeueAfter: 30 * time.Second}, fmt.Errorf("SSH key not yet available, waiting")
	}

	// Find or create instance
	result, err := r.findOrCreateInstance(ctx, contaboMachine, contaboCluster)
	if err != nil || result.RequeueAfter > 0 {
		return result, err
	}

	// Validate instance status
	if result, err := r.validateInstanceStatus(ctx, contaboMachine); err != nil || result.RequeueAfter > 0 {
		return result, err
	}

	// Handle private network assignment
	if result, err := r.reconcilePrivateNetworkAssignment(ctx, contaboMachine, contaboCluster); err != nil || result.RequeueAfter > 0 {
		return result, err
	}

	// Validate instance status again after private network assignment
	if result, err := r.validateInstanceStatus(ctx, contaboMachine); err != nil || result.RequeueAfter > 0 {
		return result, err
	}

	return ctrl.Result{}, nil
}

// findOrCreateInstance finds an existing or creates a new instance
func (r *ContaboMachineReconciler) findOrCreateInstance(ctx context.Context, contaboMachine *infrastructurev1beta2.ContaboMachine, contaboCluster *infrastructurev1beta2.ContaboCluster) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// Use mutex to prevent concurrent instance reuse
	r.instanceReuseMutex.Lock()
	defer r.instanceReuseMutex.Unlock()

	// Try to find existing instance
	instance, err := r.getExistingInstance(ctx, contaboMachine, contaboCluster)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to find existing instance: %w", err)
	}
	if instance != nil {
		contaboMachine.Status.Instance = instance
		log.Info("Found existing instance", "instanceID", instance.InstanceId)
		return ctrl.Result{}, nil
	}

	// Look for reusable instance if none found
	if contaboMachine.Status.Instance == nil {
		instance, err = r.findReusableInstance(ctx, contaboMachine, contaboCluster)
		if err != nil {
			return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
		}
		if instance != nil {
			contaboMachine.Status.Instance = instance
			log.Info("Found reusable instance", "instanceID", instance.InstanceId)
		}
	}

	// Create new instance if none found and provisioning type allows
	if contaboMachine.Status.Instance == nil {
		instance, err = r.createNewInstance(ctx, contaboMachine, contaboCluster)
		if err != nil {
			log.Error(err, "Failed to create new instance")
			// Set Failure condition instead
			contaboMachine.Status.FailureReason = ptr.To(string(infrastructurev1beta2.InstanceCreatingReason))
			contaboMachine.Status.FailureMessage = ptr.To("Failed to create new instance: " + err.Error())
			// Return error to prevent calling validateInstanceStatus with nil instance
			return ctrl.Result{}, fmt.Errorf("failed to create new instance: %w", err)
		}
		if instance != nil {
			contaboMachine.Status.Instance = instance
			log.Info("Created new instance", "instanceID", instance.InstanceId)
		}
	}

	// Update instance state (display name and networking)
	if err := r.updateInstanceState(ctx, contaboMachine, contaboCluster, contaboMachine.Status.Instance); err != nil {
		return ctrl.Result{RequeueAfter: 30 * time.Second}, r.handleError(
			ctx,
			contaboMachine,
			err,
			infrastructurev1beta2.InstanceCreatingReason,
			"Failed to update instance state",
		)
	}

	// Force requeue to retrieve instance via display name
	return ctrl.Result{RequeueAfter: 15 * time.Second}, nil
}

// reconcilePrivateNetworkAssignment handles private network assignment for the instance
func (r *ContaboMachineReconciler) reconcilePrivateNetworkAssignment(ctx context.Context, contaboMachine *infrastructurev1beta2.ContaboMachine, contaboCluster *infrastructurev1beta2.ContaboCluster) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// Assign instance to private network if not already assigned
	var privateNetwork *models.PrivateNetworkResponse

	// Retrieve private network details
	privateNetworkGetResp, err := r.ContaboClient.RetrievePrivateNetworkWithResponse(ctx, contaboCluster.Status.PrivateNetwork.PrivateNetworkId, &models.RetrievePrivateNetworkParams{})
	if err != nil || privateNetworkGetResp.StatusCode() < 200 || privateNetworkGetResp.StatusCode() >= 300 {
		return ctrl.Result{RequeueAfter: 30 * time.Second}, r.handleError(
			ctx,
			contaboMachine,
			err,
			infrastructurev1beta2.InstanceFailedReason,
			fmt.Sprintf("Failed to retrieve private network details for ID %d", contaboCluster.Status.PrivateNetwork.PrivateNetworkId),
		)
	}
	privateNetwork = &privateNetworkGetResp.JSON200.Data[0]

	// Check if instance is already part of the private network
	assignedToPrivateNetwork := false
	for _, pnInstance := range privateNetwork.Instances {
		if pnInstance.InstanceId == contaboMachine.Status.Instance.InstanceId {
			assignedToPrivateNetwork = true
			break
		}
	}

	// Assign instance to private network if not already assigned
	if !assignedToPrivateNetwork {
		log.Info("Assigning instance to private network",
			"instanceID", contaboMachine.Status.Instance.InstanceId,
			"privateNetworkID", privateNetwork.PrivateNetworkId)
		_, err := r.ContaboClient.AssignInstancePrivateNetwork(ctx, privateNetwork.PrivateNetworkId, contaboMachine.Status.Instance.InstanceId, nil)
		if err != nil {
			return ctrl.Result{RequeueAfter: 30 * time.Second}, r.handleError(
				ctx,
				contaboMachine,
				err,
				infrastructurev1beta2.InstanceFailedReason,
				"Failed to assign instance to private network",
			)
		}

		log.Info("Rebooting instance to apply private network changes",
			"instanceID", contaboMachine.Status.Instance.InstanceId)
		_, err = r.ContaboClient.Restart(ctx, contaboMachine.Status.Instance.InstanceId, nil)
		if err != nil {
			return ctrl.Result{RequeueAfter: 30 * time.Second}, r.handleError(
				ctx,
				contaboMachine,
				err,
				infrastructurev1beta2.InstanceFailedReason,
				"Failed to reboot instance after private network assignment",
			)
		}

		// Requeue to wait for the instance to be fully restarted
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	return ctrl.Result{}, nil
}

// reconcileContaboMachineAddresses updates the ContaboMachine addresses based on private network assignment
func (r *ContaboMachineReconciler) reconcileContaboMachineAddresses(ctx context.Context, contaboMachine *infrastructurev1beta2.ContaboMachine, contaboCluster *infrastructurev1beta2.ContaboCluster) error {
	log := logf.FromContext(ctx)

	log.Info("Updating contabo machine addresses")

	// Get the private network
	privateNetworkGetResp, err := r.ContaboClient.RetrievePrivateNetworkWithResponse(ctx, contaboCluster.Status.PrivateNetwork.PrivateNetworkId, &models.RetrievePrivateNetworkParams{})
	if err != nil {
		return err
	}
	privateNetwork := &privateNetworkGetResp.JSON200.Data[0]

	addresses := []clusterv1.MachineAddress{}

	// Look for internal ip v4 in private network
	for _, pnInstance := range privateNetwork.Instances {
		if pnInstance.InstanceId == contaboMachine.Status.Instance.InstanceId {
			addresses = append(addresses, clusterv1.MachineAddress{
				Type:    clusterv1.MachineInternalIP,
				Address: pnInstance.PrivateIpConfig.V4[0].Ip,
			})
			break
		}
	}
	// Look for external ip v6 in instance
	addresses = append(addresses, clusterv1.MachineAddress{
		Type:    clusterv1.MachineExternalIP,
		Address: contaboMachine.Status.Instance.IpConfig.V6.Ip,
	})
	// Look for external ip v4 in instance
	addresses = append(addresses, clusterv1.MachineAddress{
		Type:    clusterv1.MachineExternalIP,
		Address: contaboMachine.Status.Instance.IpConfig.V4.Ip,
	})

	// Add hostname entry
	addresses = append(addresses, clusterv1.MachineAddress{
		Type:    clusterv1.MachineHostName,
		Address: contaboMachine.Status.Instance.Name,
	})

	// Update status addresses
	contaboMachine.Status.Addresses = addresses

	return nil
}

// bootstrapInstance reinstalls instance with cloud-init bootstrap data
func (r *ContaboMachineReconciler) bootstrapInstance(ctx context.Context, machine *clusterv1.Machine, contaboMachine *infrastructurev1beta2.ContaboMachine, contaboCluster *infrastructurev1beta2.ContaboCluster) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// Validate instance status before bootstrapping
	result, err := r.validateInstanceStatus(ctx, contaboMachine)
	if err != nil || result.RequeueAfter > 0 {
		return result, err
	}

	// Get and validate bootstrap data
	bootstrapData, result, err := r.getBootstrapData(ctx, machine, contaboMachine, contaboCluster)
	if err != nil || result.RequeueAfter > 0 {
		return result, err
	}

	// TODO: upgrade instance if product id is not the same

	// Reinstall instance with cloud-init bootstrap data if not already done
	meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
		Type:   infrastructurev1beta2.InstanceBootstrapCondition,
		Status: metav1.ConditionFalse,
		Reason: infrastructurev1beta2.InstanceReinstallingReason,
	})
	// Retrieve SSH key IDs
	sshKeys := []int64{contaboCluster.Status.SshKey.SecretId}

	// Required for clusterctl move to not reinstall instances again
	log.Info("Check if instance is already reinstalled with bootstrap userdata by checking its cluster uuid", "instanceID", contaboMachine.Status.Instance.InstanceId)
	output, err := r.runMachineInstanceSshCommand(
		ctx,
		contaboMachine,
		contaboCluster,
		"sudo cat /etc/cluster-uuid",
	)
	if err != nil {
		log.Info("SSH command failed, will retrying...", "error", err.Error())
		return ctrl.Result{RequeueAfter: 15 * time.Second}, nil
	} else if output != nil && strings.TrimSpace(*output) == contaboCluster.Spec.ClusterUUID {
		log.Info("Instance already has the correct clusterUUID, skipping reinstall",
			"instanceID", contaboMachine.Status.Instance.InstanceId)
	} else {
		// Reinstall instance with cloud-init bootstrap data
		log.Info("Reinstalling instance with SSH keys and bootstrap data",
			"instanceId", contaboMachine.Status.Instance.InstanceId,
			"sshKeyIds", sshKeys,
			"defaultUser", contaboMachine.Status.Instance.DefaultUser,
			"imageId", contaboMachine.Status.Instance.ImageId)

		resp, err := r.ContaboClient.ReinstallInstanceWithResponse(ctx, contaboMachine.Status.Instance.InstanceId, &models.ReinstallInstanceParams{}, models.ReinstallInstanceRequest{
			SshKeys:      &sshKeys,
			DefaultUser:  ptr.To(models.ReinstallInstanceRequestDefaultUserAdmin),
			ImageId:      contaboMachine.Status.Instance.ImageId,
			RootPassword: nil,
			UserData:     &bootstrapData,
		})
		if err != nil || resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
			meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
				Type:   infrastructurev1beta2.InstanceBootstrapCondition,
				Status: metav1.ConditionFalse,
				Reason: infrastructurev1beta2.InstanceReinstallingFailedReason,
			})
			return ctrl.Result{RequeueAfter: 30 * time.Second}, r.handleError(
				ctx,
				contaboMachine,
				err,
				infrastructurev1beta2.InstanceReinstallingFailedReason,
				fmt.Sprintf("Failed to reinstall instance, statusCode: %d", resp.StatusCode()),
			)
		}
		log.Info("Reinstall instance request sent successfully",
			"instanceId", contaboMachine.Status.Instance.InstanceId)

		// Refresh instance status after reinstall
		instanceResp, err := r.ContaboClient.RetrieveInstanceWithResponse(ctx, contaboMachine.Status.Instance.InstanceId, nil)
		if err == nil && instanceResp.StatusCode() >= 200 && instanceResp.StatusCode() < 300 {
			contaboMachine.Status.Instance = convertInstanceResponseData(&instanceResp.JSON200.Data[0])
		}

		// Validate instance status after reinstall
		result, err = r.validateInstanceStatus(ctx, contaboMachine)
		if err != nil || result.RequeueAfter > 0 {
			return result, err
		}
	}

	// Wait for cloud-init to finish
	meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
		Type:   infrastructurev1beta2.InstanceBootstrapCondition,
		Status: metav1.ConditionFalse,
		Reason: infrastructurev1beta2.InstanceWaitingForCloudInitReason,
	})
	// Wait for cloud-init to finish
	log.Info("Waiting for cloud-init to finish on instance",
		"instanceID", contaboMachine.Status.Instance.InstanceId,
		"instanceIP", contaboMachine.Status.Instance.IpConfig.V4.Ip)

	cloudInitStatus, err := r.runMachineInstanceSshCommand(
		ctx,
		contaboMachine,
		contaboCluster,
		"cloud-init status",
	)
	if cloudInitStatus == nil && err != nil {
		log.Info("SSH command failed, will retry", "error", err.Error(), "requeueAfter", "10s")
		return ctrl.Result{RequeueAfter: 15 * time.Second}, nil
	}
	if cloudInitStatus == nil {
		log.Info("SSH command returned empty output, will retry", "requeueAfter", "10s")
		return ctrl.Result{RequeueAfter: 15 * time.Second}, nil
	}
	// Check if cloud-init finished successfully
	if strings.Contains(*cloudInitStatus, "status: error") {
		message := "cloud-init failed, check the /var/log/cloud-init.log and /var/log/cloud-init-output.log files on the instance for more details"
		// Try to find error from cloud-init logs
		cloudInitLogsError, _ := r.runMachineInstanceSshCommand(
			ctx,
			contaboMachine,
			contaboCluster,
			"sudo cat /var/log/cloud-init-output.log | grep '[CAPC] Error'",
		)
		if cloudInitLogsError != nil && *cloudInitLogsError != "" {
			message = "cloud-init error: " + strings.Split(*cloudInitLogsError, "[CAPC] Error")[1]
		}
		err := errors.New(message)
		log.Error(err, message)
		if err := r.resetInstance(ctx, contaboMachine, contaboMachine.Status.Instance, &message); err != nil {
			log.Error(err, "Failed to reset instance after cloud-init failure",
				"instanceID", contaboMachine.Status.Instance.InstanceId)
		}
		return ctrl.Result{RequeueAfter: 15 * time.Second}, err
	}

	if strings.Contains(*cloudInitStatus, "status: running") {
		log.Info("cloud-init is still running, will retry", "output", *cloudInitStatus, "requeueAfter", "20s")
		return ctrl.Result{RequeueAfter: 20 * time.Second}, nil
	}

	log.Info("cloud-init has finished on instance",
		"instanceID", contaboMachine.Status.Instance.InstanceId)

	if result, err := r.initializeNode(ctx, machine, contaboMachine, contaboCluster); err != nil || result.RequeueAfter > 0 {
		return result, err
	}

	meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
		Type:   infrastructurev1beta2.InstanceBootstrapCondition,
		Status: metav1.ConditionTrue,
		Reason: infrastructurev1beta2.InstanceBootstrapedReason,
	})

	return ctrl.Result{}, nil
}

// initializeNode reinstalls instance with cloud-init bootstrap data
func (r *ContaboMachineReconciler) initializeNode(ctx context.Context, machine *clusterv1.Machine, contaboMachine *infrastructurev1beta2.ContaboMachine, contaboCluster *infrastructurev1beta2.ContaboCluster) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// Configure node
	// 1. Get the kubeconfig
	// 2. Connect to cluster
	// 3. Wait for node to exists and be ready
	// 4. Update node external IPs if not already set
	// 5. Set Ready condition to true
	// 6. Remove uninitialized taint
	meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
		Type:   infrastructurev1beta2.InstanceBootstrapCondition,
		Status: metav1.ConditionFalse,
		Reason: infrastructurev1beta2.InstanceConfigureNodeReason,
	})
	log.Info("Waiting for node to exists in the cluster",
		"instanceID", contaboMachine.Status.Instance.InstanceId,
		"instanceIP", contaboMachine.Status.Instance.IpConfig.V4.Ip)

	kubeconfig, err := r.getKubeconfig(ctx, contaboCluster)
	if err != nil {
		log.Info("Waiting to get kubeconfig from owner Cluster",
			"clusterName", machine.Spec.ClusterName,
			"namespace", machine.Namespace)
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeconfig))
	if err != nil {
		log.Error(err, "Failed to create REST config from kubeconfig")
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	k8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Error(err, "Failed to create Kubernetes client from REST config")
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	nodeName, err := ParseProviderID(*contaboMachine.Spec.ProviderID)
	if err != nil {
		log.Error(err, "Failed to parse node name from provider ID",
			"providerID", *contaboMachine.Spec.ProviderID)
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	log.Info("Node found in cluster, updating ips",
		"nodeName", nodeName,
		"instanceID", contaboMachine.Status.Instance.InstanceId)

	// Build a map to deduplicate addresses by (Type, Address)
	// Collect all addresses
	addresses := []corev1.NodeAddress{}
	for _, addr := range contaboMachine.Status.Addresses {
		addresses = append(addresses, corev1.NodeAddress{
			Type:    corev1.NodeAddressType(addr.Type),
			Address: addr.Address,
		})
	}

	// Patch node with merged addresses
	patchNode := &corev1.Node{
		Status: corev1.NodeStatus{
			Addresses: addresses,
		},
	}

	patchBytes, err := json.Marshal(patchNode)
	if err != nil {
		log.Error(err, "Failed to marshal node status patch")
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	log.Info("Patching node with IP addresses",
		"nodeName", nodeName,
		"ips", addresses)

	_, err = k8sClient.CoreV1().Nodes().Patch(ctx, nodeName, types.MergePatchType, patchBytes, metav1.PatchOptions{}, "status")
	if err != nil {
		log.Error(err, "Failed to patch node with IPs", "nodeName", nodeName)
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	log.Info("Successfully patched node with IPs", "nodeName", nodeName, "ips", addresses)

	// Set node providerId and Remove unitialized taint if present
	node, err := k8sClient.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("Node not found in cluster yet, will retry", "nodeName", nodeName)
			return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
		}
		log.Info("Node not found in cluster, kubelet seems not ready, will retry", "nodeName", nodeName)
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	node.Spec.ProviderID = BuildProviderID(contaboMachine.Status.Instance.Name)

	taintRemoved := false
	for i, taint := range node.Spec.Taints {
		if taint.Key == "node.cloudprovider.kubernetes.io/uninitialized" {
			node.Spec.Taints = append(node.Spec.Taints[:i], node.Spec.Taints[i+1:]...)
			taintRemoved = true
			break
		}
	}
	if taintRemoved {
		_, err = k8sClient.CoreV1().Nodes().Update(ctx, node, metav1.UpdateOptions{})
		if err != nil {
			log.Error(err, "Failed to initialize node", "nodeName", nodeName)
			return ctrl.Result{}, fmt.Errorf("failed to initialize node %s: %w", nodeName, err)
		}
	}
	log.Info("Successfully initialize node", "nodeName", nodeName)

	return ctrl.Result{}, nil
}

func (r *ContaboMachineReconciler) reconcileDelete(ctx context.Context, contaboMachine *infrastructurev1beta2.ContaboMachine, contaboCluster *infrastructurev1beta2.ContaboCluster) {
	_ = contaboCluster // may be used in future for cluster-specific cleanup logic
	log := logf.FromContext(ctx)

	log.Info("Reconciling ContaboMachine delete - setting instance available for reuse")

	// Update machine condition
	meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
		Type:   infrastructurev1beta2.InstanceReadyCondition,
		Status: metav1.ConditionFalse,
		Reason: infrastructurev1beta2.InstanceDeletingReason,
	})
	log.Info("Machine marked for deletion, proceeding with instance cleanup",
		"name", contaboMachine.Name)

	// Retrieve instance details
	if contaboMachine.Status.Instance == nil {
		log.Info("Instance is already nil, assuming it is deleted, removing finalizer",
			"name", contaboMachine.Name)
		controllerutil.RemoveFinalizer(contaboMachine, infrastructurev1beta2.MachineFinalizer)
		return
	}

	// First, stop the instance
	_, err := r.ContaboClient.Stop(ctx, contaboMachine.Status.Instance.InstanceId, nil)
	if err != nil {
		log.Error(err, "Failed to stop instance during deletion",
			"instanceID", contaboMachine.Status.Instance.InstanceId)
	}

	// Reset the instance by removing it from any private networks
	if err := r.resetInstance(ctx, contaboMachine, contaboMachine.Status.Instance, nil); err != nil {
		log.Error(err, "Failed to reset instance during deletion",
			"instanceID", contaboMachine.Status.Instance.InstanceId)
	}

	// Update instance display name to mark it as reusable
	displayName := ""
	_, err = r.ContaboClient.PatchInstanceWithResponse(ctx, contaboMachine.Status.Instance.InstanceId, nil, models.PatchInstanceRequest{
		DisplayName: &displayName,
	})
	if err != nil {
		log.Error(err, "Failed to update instance display name during deletion",
			"instanceID", contaboMachine.Status.Instance.InstanceId,
			"newDisplayName", displayName)
	}

	// Remove finalizer
	controllerutil.RemoveFinalizer(contaboMachine, infrastructurev1beta2.MachineFinalizer)
	log.Info("Removed finalizer from ContaboMachine")
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
	var cancelDate *string
	if instanceList.CancelDate != nil {
		cancelDate = ptr.To(instanceList.CancelDate.Format(time.RFC3339))
	}

	if instanceList.SshKeys == nil {
		instanceList.SshKeys = []int64{}
	}

	return &infrastructurev1beta2.ContaboInstanceStatus{
		AddOns:        addons,
		AdditionalIps: addionalIps,
		CancelDate:    cancelDate,
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
	var cancelDate *string
	if instanceList.CancelDate != nil {
		cancelDate = ptr.To(instanceList.CancelDate.Format(time.RFC3339))
	}

	if instanceList.SshKeys == nil {
		instanceList.SshKeys = []int64{}
	}

	return &infrastructurev1beta2.ContaboInstanceStatus{
		AddOns:        addons,
		AdditionalIps: addionalIps,
		CancelDate:    cancelDate,
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

func (r *ContaboMachineReconciler) runMachineInstanceSshCommand(ctx context.Context, contaboMachine *infrastructurev1beta2.ContaboMachine, contaboCluster *infrastructurev1beta2.ContaboCluster, command string) (*string, error) {
	log := logf.FromContext(ctx)

	// Get ssh-key from secret
	sshKeySecret := &corev1.Secret{}
	sshKeySecretMetadata := client.ObjectKey{
		Namespace: contaboMachine.Namespace,
		Name:      FormatSshKeyKubernetesName(contaboCluster),
	}
	if err := r.Get(ctx, sshKeySecretMetadata, sshKeySecret); err != nil {
		return nil, fmt.Errorf("failed to get SSH private key secret %s/%s: %v", sshKeySecretMetadata.Namespace, sshKeySecretMetadata.Name, err)
	}
	// Connect to the instance via SSH and wait for cloud-init to finish
	// We try to connect every 10 seconds for up to 15 minutes
	sshPrivateKey, ok := sshKeySecret.Data["id_rsa"]
	if !ok || len(sshPrivateKey) == 0 {
		return nil, fmt.Errorf("SSH private key secret is missing 'id_rsa' key or is empty")
	}

	// Also get the public key from secret for verification
	sshPublicKeyFromSecret, ok := sshKeySecret.Data["id_rsa.pub"]
	if !ok || len(sshPublicKeyFromSecret) == 0 {
		return nil, fmt.Errorf("SSH public key secret is missing 'id_rsa.pub' key or is empty")
	}

	log.Info("Retrieved SSH keys from secret",
		"privateKeyLength", len(sshPrivateKey),
		"publicKeyLength", len(sshPublicKeyFromSecret),
		"secretName", sshKeySecretMetadata.Name)
	// SSH client configuration

	host := contaboMachine.Status.Instance.IpConfig.V4.Ip
	user := string(infrastructurev1beta2.InstanceResponseDefaultUserAdmin)
	if contaboMachine.Status.Instance.DefaultUser != nil {
		user = string(*contaboMachine.Status.Instance.DefaultUser)
	}

	// Define potential usernames to try if authentication fails
	log.Info("SSH connection setup", "user", user)
	log.Info("Attempting to parse SSH private key", "keyLength", len(sshPrivateKey), "host", host, "user", user)

	signer, err := ssh.ParsePrivateKey(sshPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SSH private key (length: %d): %v", len(sshPrivateKey), err)
	}

	// SSH client configuration with user fallback
	var sshClient *ssh.Client

	log.Info("Starting SSH connection attempts", "host", host, "primaryUser", user)

	// Try primary user first, then fallback users for auth errors
	currentUser := user

	config := &ssh.ClientConfig{
		User: currentUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	sshClient, err = ssh.Dial("tcp", net.JoinHostPort(host, "22"), config)
	if err != nil {
		// Provide more specific error information for better requeue handling
		errStr := err.Error()
		isAuthError := strings.Contains(errStr, "handshake failed") || strings.Contains(errStr, "unable to authenticate") || strings.Contains(errStr, "no supported methods remain")
		isNetworkError := strings.Contains(errStr, "connection refused") || strings.Contains(errStr, "timeout")

		// Update ssh key of the contabo machine auth failed
		if isAuthError {
			sshKeys := []int64{contaboCluster.Status.SshKey.SecretId}
			log.Info("SSH connection method issue, update instance ssh key", "host", host, "user", user, "sshKeys", sshKeys)
			_, err := r.ContaboClient.ResetPasswordAction(ctx, contaboMachine.Status.Instance.InstanceId, nil, models.InstancesResetPasswordActionsRequest{
				SshKeys: &sshKeys,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to update instance SSH keys after connection method issue: %v", err)
			}
		}
		if isAuthError {
			return nil, fmt.Errorf("SSH authentication failed to %s with user %s: %v", host, currentUser, err)
		}
		if isNetworkError {
			return nil, fmt.Errorf("SSH network connection failed to %s (instance may not be ready): %v", host, err)
		}
		return nil, fmt.Errorf("SSH connection failed to %s: %v", host, err)
	}
	defer func() {
		if closeErr := sshClient.Close(); closeErr != nil {
			// log.Error(closeErr, "Failed to close SSH client")
			return
		}
	}()

	log.Info("SSH connection established successfully", "host", host, "user", user)

	// Create a session to run the cloud-init status command
	session, err := sshClient.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH session: %v", err)
	}

	// Run the command and capture output
	output, err := session.CombinedOutput(command)
	if err != nil {
		log.Info("SSH command execution failed", "command", command, "error", err.Error())
		return ptr.To(string(output)), nil
	}
	return ptr.To(string(output)), nil
}

// handleError centralizes error handling with status condition, logging, event recording, and patching
func (r *ContaboMachineReconciler) handleError(ctx context.Context, contaboMachine *infrastructurev1beta2.ContaboMachine, err error, reason string, message string) error {
	log := logf.FromContext(ctx)

	// Set status condition (always uses MachineReadyCondition)
	meta.SetStatusCondition(&contaboMachine.Status.Conditions, metav1.Condition{
		Type:    infrastructurev1beta2.InstanceReadyCondition,
		Status:  metav1.ConditionFalse,
		Reason:  reason,
		Message: message,
	})

	// Log the error
	log.Error(err, message)

	// Record event
	// r.Recorder.Event(contaboMachine, corev1.EventTypeWarning, reason, message)

	return fmt.Errorf("%s: %w", message, err)
}

// getKubeconfig retrieves the kubeconfig from the owner Cluster's secret
// resetInstance prepares an instance for reuse by removing it from any private networks
func (r *ContaboMachineReconciler) resetInstance(ctx context.Context, contaboMachine *infrastructurev1beta2.ContaboMachine, instance *infrastructurev1beta2.ContaboInstanceStatus, errorMessage *string) error {
	log := logf.FromContext(ctx)

	hasErrorMessage := errorMessage != nil || (instance != nil && instance.ErrorMessage != nil)

	// Set error on contabo machine status
	if instance == nil {
		log.Info("Instance is nil, nothing to reset")
		return nil
	}

	// Set error on contabo instance displayName
	displayName := ""
	if errorMessage != nil {
		displayName = Truncate(fmt.Sprintf("[capc] %d %s", instance.InstanceId, *errorMessage), 255) // Contabo display name max length is 255 characters
	} else if instance.ErrorMessage != nil {
		displayName = Truncate(fmt.Sprintf("[capc] %d %s", instance.InstanceId, *instance.ErrorMessage), 255) // Contabo display name max length is 255 characters
	}
	_, err := r.ContaboClient.PatchInstanceWithResponse(ctx, instance.InstanceId, nil, models.PatchInstanceRequest{
		DisplayName: &displayName,
	})
	if err != nil {
		log.Error(err, "Failed to update instance display name to avoid reuse",
			"instanceID", instance.InstanceId,
			"newDisplayName", displayName)
	}

	if hasErrorMessage {
		log.Info("Instance marked with error message, skipping private network unassignment, for further investigation",
			"instanceID", instance.InstanceId)
		return nil
	}

	// TODO: code is ugly, needs refactoring
	// Check if instance has any private networks assigned
	if len(instance.AddOns) > 0 {
		for _, addon := range instance.AddOns {
			// ID 1477 is the Private Networking addon in Contabo API
			if addon.Id == 1477 {
				// Instance has private networking, need to check which networks
				for page := int64(1); ; page++ {
					var privateNetworksResp *contaboclient.RetrievePrivateNetworkListResponse
					// Retrieve all private networks and check if instance is part of any
					privateNetworksResp, err := r.ContaboClient.RetrievePrivateNetworkListWithResponse(ctx, &models.RetrievePrivateNetworkListParams{
						Page: &page,
						Size: ptr.To(int64(100)),
					})
					// If there was an error or no more private networks, break the loop
					if err != nil || privateNetworksResp.StatusCode() < 200 || privateNetworksResp.StatusCode() >= 300 || len(privateNetworksResp.JSON200.Data) == 0 {
						break
					}
					// If we successfully retrieved the private networks, check if instance is part of any
					for _, network := range privateNetworksResp.JSON200.Data {
						for _, instanceInNetwork := range network.Instances {
							if instanceInNetwork.InstanceId == instance.InstanceId {
								log.Info("Found instance in private network, unassigning",
									"instanceID", instance.InstanceId,
									"networkID", network.PrivateNetworkId)
								unassignResp, err := r.ContaboClient.UnassignInstancePrivateNetworkWithResponse(ctx, network.PrivateNetworkId, instance.InstanceId, nil)
								if err == nil && unassignResp.StatusCode() >= 200 && unassignResp.StatusCode() < 300 {
									log.Info("Successfully unassigned private network from instance",
										"instanceID", instance.InstanceId,
										"networkID", network.PrivateNetworkId)
									time.Sleep(1 * time.Second)

								} else {
									log.Error(err, "Failed to unassign private network from instance",
										"instanceID", instance.InstanceId,
										"networkID", network.PrivateNetworkId)
								}
								break
							}
						}
					}
				}
				break // No need to check other addons once we found private networking
			}
		}
	}

	time.Sleep(1 * time.Second) // Wait a bit to ensure unassignment is processed

	// Reinstall to clear any residual configuration
	_, err = r.ContaboClient.ReinstallInstance(ctx, instance.InstanceId, &models.ReinstallInstanceParams{}, models.ReinstallInstanceRequest{
		ImageId:     DefaultUbuntuImageID,
		DefaultUser: ptr.To(models.ReinstallInstanceRequestDefaultUserAdmin),
	})
	if err != nil {
		log.Error(err, "Failed to reinstall instance after unassigning private network",
			"instanceID", instance.InstanceId)
	} else {
		log.Info("Reinstalled instance after unassigning private network",
			"instanceID", instance.InstanceId)
	}

	// Remove Instance from Status
	contaboMachine.Status = infrastructurev1beta2.ContaboMachineStatus{}

	// Remove ProviderID
	contaboMachine.Spec.ProviderID = nil

	return err
}

func (r *ContaboMachineReconciler) getKubeconfig(ctx context.Context, contaboCluster *infrastructurev1beta2.ContaboCluster) (string, error) {
	kubeconfigSecret := &corev1.Secret{}
	kubeconfigSecretName := types.NamespacedName{
		Namespace: contaboCluster.Namespace,
		Name:      contaboCluster.Name + "-kubeconfig",
	}
	if err := r.Get(ctx, kubeconfigSecretName, kubeconfigSecret); err != nil {
		return "", err
	}

	kubeconfigData, ok := kubeconfigSecret.Data["value"]
	if !ok {
		return "", errors.New("kubeconfig secret is missing 'value' key")
	}

	return string(kubeconfigData), nil
}
