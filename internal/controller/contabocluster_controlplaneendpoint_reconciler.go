package controller

import (
	"context"
	"fmt"
	"net"
	"time"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	infrastructurev1beta2 "github.com/ctnr-io/cluster-api-provider-contabo/api/v1beta2"

	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/utils/ptr"
)

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
			Annotations: map[string]string{
				clusterv1.ClusterNameAnnotation: contaboCluster.Name,
			},
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
	endpointSliceV4Name := fmt.Sprintf("%s-ipv4-apiserver", contaboCluster.Name)
	endpointSliceV6Name := fmt.Sprintf("%s-ipv6-apiserver", contaboCluster.Name)

	// Collect all control plane instance IPs
	var endpointsV4 []discoveryv1.Endpoint
	var endpointsV6 []discoveryv1.Endpoint
	for _, machine := range controlPlaneContaboMachineList.Items {
		if machine.Status.Instance != nil && machine.Status.Instance.IpConfig.V4.Ip != "" {
			endpointsV4 = append(endpointsV4, discoveryv1.Endpoint{
				Addresses: []string{machine.Status.Instance.IpConfig.V4.Ip},
				Conditions: discoveryv1.EndpointConditions{
					Ready: ptr.To(true),
				},
			})
		}
		if machine.Status.Instance != nil && machine.Status.Instance.IpConfig.V6.Ip != "" {
			// Convert from 2001:0db8:85a3:0000:0000:0000:0000:0001 to 2001:db8:85a3::1
			canonicalIpV6 := net.ParseIP(machine.Status.Instance.IpConfig.V6.Ip).String()
			endpointsV6 = append(endpointsV6, discoveryv1.Endpoint{
				Addresses: []string{canonicalIpV6},
				Conditions: discoveryv1.EndpointConditions{
					Ready: ptr.To(true),
				},
			})
		}
	}

	endpointPort := contaboCluster.Spec.ControlPlaneEndpoint.Port

	endpointSliceV4 := &discoveryv1.EndpointSlice{
		ObjectMeta: metav1.ObjectMeta{
			Name:      endpointSliceV4Name,
			Namespace: contaboCluster.Namespace,
			Annotations: map[string]string{
				clusterv1.ClusterNameAnnotation: contaboCluster.Name,
			},
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
		Endpoints:   endpointsV4,
		Ports: []discoveryv1.EndpointPort{
			{
				Name:     ptr.To("https"),
				Port:     ptr.To(endpointPort),
				Protocol: ptr.To(corev1.ProtocolTCP),
			},
		},
	}

	endpointSliceV6 := &discoveryv1.EndpointSlice{
		ObjectMeta: metav1.ObjectMeta{
			Name:      endpointSliceV6Name,
			Namespace: contaboCluster.Namespace,
			Annotations: map[string]string{
				clusterv1.ClusterNameAnnotation: contaboCluster.Name,
			},
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
		AddressType: discoveryv1.AddressTypeIPv6,
		Endpoints:   endpointsV6,
		Ports: []discoveryv1.EndpointPort{
			{
				Name:     ptr.To("https"),
				Port:     ptr.To(endpointPort),
				Protocol: ptr.To(corev1.ProtocolTCP),
			},
		},
	}

	// Try to get existing endpoint slice v4
	existingEndpointSliceV4 := &discoveryv1.EndpointSlice{}
	err := r.Get(ctx, client.ObjectKey{
		Name:      endpointSliceV4Name,
		Namespace: contaboCluster.Namespace,
	}, existingEndpointSliceV4)

	if err != nil {
		if apierrors.IsNotFound(err) {
			// Create the endpoint slice
			log.Info("Creating control plane endpoint slice", "endpointSliceV4Name", endpointSliceV4Name, "endpointCount", len(endpointsV4))
			if err := r.Create(ctx, endpointSliceV4); err != nil {
				return fmt.Errorf("failed to create control plane endpoint slice: %w", err)
			}
			log.Info("Created control plane endpoint slice", "endpointSliceV4Name", endpointSliceV4Name)
		} else {
			return fmt.Errorf("failed to get control plane endpoint slice: %w", err)
		}
	} else {
		// Update the endpoint slice if needed
		existingEndpointSliceV4.Endpoints = endpointSliceV4.Endpoints
		existingEndpointSliceV4.Ports = endpointSliceV4.Ports
		log.Info("Updating control plane endpoint slice", "endpointSliceV4Name", endpointSliceV4Name, "endpointCount", len(endpointsV4))
		if err := r.Update(ctx, existingEndpointSliceV4); err != nil {
			return fmt.Errorf("failed to update control plane endpoint slice: %w", err)
		}
	}

	// Try to get existing endpoint slice v6
	existingEndpointSliceV6 := &discoveryv1.EndpointSlice{}
	err = r.Get(ctx, client.ObjectKey{
		Name:      endpointSliceV6Name,
		Namespace: contaboCluster.Namespace,
	}, existingEndpointSliceV6)

	if err != nil {
		if apierrors.IsNotFound(err) {
			// Create the endpoint slice
			log.Info("Creating control plane endpoint slice", "endpointSliceV6Name", endpointSliceV6Name, "endpointCount", len(endpointsV6))
			if err := r.Create(ctx, endpointSliceV6); err != nil {
				return fmt.Errorf("failed to create control plane endpoint slice: %w", err)
			}
			log.Info("Created control plane endpoint slice", "endpointSliceV6Name", endpointSliceV6Name)
		} else {
			return fmt.Errorf("failed to get control plane endpoint slice: %w", err)
		}
	} else {
		// Update the endpoint slice if needed
		existingEndpointSliceV6.Endpoints = endpointSliceV6.Endpoints
		existingEndpointSliceV6.Ports = endpointSliceV6.Ports
		log.Info("Updating control plane endpoint slice", "endpointSliceV6Name", endpointSliceV6Name, "endpointCount", len(endpointsV6))
		if err := r.Update(ctx, existingEndpointSliceV6); err != nil {
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
