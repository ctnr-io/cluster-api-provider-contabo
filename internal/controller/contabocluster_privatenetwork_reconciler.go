package controller

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	infrastructurev1beta2 "github.com/ctnr-io/cluster-api-provider-contabo/api/v1beta2"
	"github.com/ctnr-io/cluster-api-provider-contabo/pkg/contabo/v1.0.0/models"
)

// reconcilePrivateNetwork ensures the private network exists and is configured
func (r *ContaboClusterReconciler) reconcilePrivateNetwork(ctx context.Context, contaboCluster *infrastructurev1beta2.ContaboCluster) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	log.Info("Reconciling private network for ContaboCluster", "cluster", contaboCluster.Name)

	// Check if private network was created
	privateNetworkName := FormatPrivateNetworkName(contaboCluster)
	if contaboCluster.Spec.PrivateNetwork.Name != "" {
		privateNetworkName = contaboCluster.Spec.PrivateNetwork.Name
	}

	// Check if private network with the same name already exists in Contabo API
	resp, err := r.ContaboClient.RetrievePrivateNetworkListWithResponse(ctx, &models.RetrievePrivateNetworkListParams{
		Name: &privateNetworkName,
	})
	if err != nil || resp.JSON200 == nil || resp.JSON200.Data == nil || len(resp.JSON200.Data) == 0 {
		log.Info("Private network not found in Contabo API, creating new one", "privateNetworkName", privateNetworkName)

		// Create private network if not found
		description := "Private network created by Cluster API Provider Contabo"
		privateNetworkCreateResp, err := r.ContaboClient.CreatePrivateNetworkWithResponse(ctx, nil, models.CreatePrivateNetworkJSONRequestBody{
			Name:        privateNetworkName,
			Description: &description,
			Region:      &contaboCluster.Spec.PrivateNetwork.Region,
		})
		if err != nil || privateNetworkCreateResp.StatusCode() < 200 || privateNetworkCreateResp.StatusCode() >= 300 {
			return ctrl.Result{}, r.handleError(
				ctx,
				contaboCluster,
				err,
				infrastructurev1beta2.ClusterPrivateNetworkReadyCondition,
				infrastructurev1beta2.ClusterPrivateNetworkFailedReason,
				"Failed to create private network",
			)
		}

		// Requeue to retry after private network creation
		return ctrl.Result{RequeueAfter: 5 * time.Second}, nil
	}

	log.Info("Found existing private network in Contabo API", "privateNetworkName", privateNetworkName)

	privateNetwork := &resp.JSON200.Data[0]

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

	log.Info("Private network reconciled", "privateNetworkName", privateNetwork.Name, "privateNetworkId", privateNetwork.PrivateNetworkId)

	return ctrl.Result{}, nil
}
