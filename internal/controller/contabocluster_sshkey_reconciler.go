package controller

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	infrastructurev1beta2 "github.com/ctnr-io/cluster-api-provider-contabo/api/v1beta2"
	"github.com/ctnr-io/cluster-api-provider-contabo/pkg/contabo/v1.0.0/models"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"
)

// reconcileSSHKey ensures the SSH key exists and is configured
func (r *ContaboClusterReconciler) reconcileSSHKey(ctx context.Context, contaboCluster *infrastructurev1beta2.ContaboCluster) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	log.Info("Reconciling SSH key for ContaboCluster", "cluster", contaboCluster.Name)

	if result, err := r.reconcileKubernetesSSHKeySecret(ctx, contaboCluster); err != nil || result.RequeueAfter != 0 {
		return result, err
	}

	if result, err := r.reconcileContaboSSHKeySecret(ctx, contaboCluster); err != nil || result.RequeueAfter != 0 {
		return result, err
	}

	meta.SetStatusCondition(&contaboCluster.Status.Conditions, metav1.Condition{
		Type:   infrastructurev1beta2.ClusterSshKeyReadyCondition,
		Status: metav1.ConditionTrue,
		Reason: infrastructurev1beta2.ClusterAvailableReason,
	})

	return ctrl.Result{}, nil
}

// reconcileKubernetesSSHKeySecret ensures the SSH key secret exists in Kubernetes
func (r *ContaboClusterReconciler) reconcileKubernetesSSHKeySecret(ctx context.Context, contaboCluster *infrastructurev1beta2.ContaboCluster) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	sshKeyKubernetesName := FormatSshKeyKubernetesName(contaboCluster)

	sshKeySecret := &corev1.Secret{}

	// Check if secret already exists in Kubernetes
	err := r.Get(ctx, client.ObjectKey{
		Name:      sshKeyKubernetesName,
		Namespace: contaboCluster.Namespace,
	}, sshKeySecret)
	if err != nil {
		log.Info("SSH key secret not found in Kubernetes, creating new one", "secretName", sshKeyKubernetesName)

		// Generate an ssh key pair
		privateKey, publicKey, err := generateSSHKeyPair()
		if err != nil {
			return ctrl.Result{}, r.handleError(
				ctx,
				contaboCluster,
				err,
				infrastructurev1beta2.ClusterSshKeyReadyCondition,
				infrastructurev1beta2.ClusterSshKeyFailedReason,
				"Failed to generate SSH key pair",
			)
		}

		// Create new secret
		secretData := map[string][]byte{
			"id_rsa":     []byte(privateKey),
			"id_rsa.pub": []byte(publicKey),
		}
		err = r.Create(ctx, &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      sshKeyKubernetesName,
				Namespace: contaboCluster.Namespace,
				Annotations: map[string]string{
					clusterv1.ClusterNameAnnotation: contaboCluster.Name,
				},
				Labels: map[string]string{
					clusterv1.ClusterNameLabel: contaboCluster.Name,
					"component":                "ssh-key",
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
			Type: clusterv1.ClusterSecretType,
			Data: secretData,
		})

		if err != nil {
			return ctrl.Result{}, r.handleError(
				ctx,
				contaboCluster,
				err,
				infrastructurev1beta2.ClusterSshKeyReadyCondition,
				infrastructurev1beta2.ClusterSshKeyFailedReason,
				"Failed to create SSH key secret",
			)
		}

		// Requeue to allow time for the secret to be created
		return ctrl.Result{RequeueAfter: 15 * time.Second}, nil
	}

	return ctrl.Result{}, nil
}

// reconcileContaboSSHKeySecret ensures the SSH key exists in Contabo API
func (r *ContaboClusterReconciler) reconcileContaboSSHKeySecret(ctx context.Context, contaboCluster *infrastructurev1beta2.ContaboCluster) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	var sshKey *models.SecretResponse
	sshKeyContaboName := FormatSshKeyContaboName(contaboCluster)
	sshKeyKubernetesName := FormatSshKeyKubernetesName(contaboCluster)

	sshKeySecret := &corev1.Secret{}

	// Retrieve the SSH key secret from Kubernetes
	err := r.Get(ctx, client.ObjectKey{
		Name:      sshKeyKubernetesName,
		Namespace: contaboCluster.Namespace,
	}, sshKeySecret)
	if err != nil {
		return ctrl.Result{}, r.handleError(
			ctx,
			contaboCluster,
			err,
			infrastructurev1beta2.ClusterSshKeyReadyCondition,
			infrastructurev1beta2.ClusterSshKeyFailedReason,
			"Failed to retrieve SSH key secret from Kubernetes",
		)
	}

	log.Info("Found existing SSH key secret in Kubernetes", "secretName", sshKeyKubernetesName)

	publicKey := string(sshKeySecret.Data["id_rsa.pub"])

	// Check if SSH key status is already set and valid
	if contaboCluster.Status.SshKey != nil {
		// Verify the SSH key still exists in Contabo API
		sshKeyRetrieveResp, err := r.ContaboClient.RetrieveSecretWithResponse(ctx, contaboCluster.Status.SshKey.SecretId, &models.RetrieveSecretParams{})
		if err == nil && sshKeyRetrieveResp.StatusCode() >= 200 && sshKeyRetrieveResp.StatusCode() < 300 {
			// SSH key exists and status is correct, nothing to do
			log.V(1).Info("SSH key already configured correctly in status", 
				"sshKeyID", contaboCluster.Status.SshKey.SecretId,
				"sshKeyName", contaboCluster.Status.SshKey.Name)
			return ctrl.Result{}, nil
		}
		// SSH key in status doesn't exist anymore, clear it and continue to create/find new one
		log.Info("SSH key in status no longer exists in Contabo API, will find or create new one",
			"oldSSHKeyID", contaboCluster.Status.SshKey.SecretId)
		contaboCluster.Status.SshKey = nil
	}

	// Check if SSH key with the same name already exists in Contabo API
	resp, err := r.ContaboClient.RetrieveSecretListWithResponse(ctx, &models.RetrieveSecretListParams{
		Name: &sshKeyContaboName,
	})

	if err != nil || resp.JSON200 == nil || resp.JSON200.Data == nil || len(resp.JSON200.Data) == 0 {
		log.Info("SSH key not found in Contabo API, creating new one", "sshKeyContaboName", sshKeyContaboName)

		// Create SSH key if not found
		// Trim the public key for Contabo API (remove trailing newline)
		trimmedPublicKey := strings.TrimSpace(publicKey)

		sshKeyCreateResp, err := r.ContaboClient.CreateSecretWithResponse(ctx, &models.CreateSecretParams{}, models.CreateSecretRequest{
			Name:  sshKeyContaboName,
			Value: trimmedPublicKey,
			Type:  "ssh",
		})
		if err != nil || sshKeyCreateResp.StatusCode() < 200 || sshKeyCreateResp.StatusCode() >= 300 {
			return ctrl.Result{}, r.handleError(
				ctx,
				contaboCluster,
				err,
				infrastructurev1beta2.ClusterSshKeyReadyCondition,
				infrastructurev1beta2.ClusterSshKeyFailedReason,
				fmt.Sprintf("Failed to submit SSH public key to Contabo API: %s", sshKeyContaboName),
			)
		}

		// Retrieve the created SSH key for further processing
		sshKeyRetrieveResp, err := r.ContaboClient.RetrieveSecretWithResponse(ctx, int64(sshKeyCreateResp.JSON201.Data[0].SecretId), &models.RetrieveSecretParams{})
		if err != nil || sshKeyRetrieveResp.StatusCode() < 200 || sshKeyRetrieveResp.StatusCode() >= 300 {
			return ctrl.Result{}, r.handleError(
				ctx,
				contaboCluster,
				err,
				infrastructurev1beta2.ClusterSshKeyReadyCondition,
				infrastructurev1beta2.ClusterSshKeyFailedReason,
				"Failed to retrieve created SSH key from Contabo API",
			)
		}
		
		sshKey = &sshKeyRetrieveResp.JSON200.Data[0]
		log.Info("Created new SSH key in Contabo API", "sshKeyID", sshKey.SecretId, "sshKeyName", sshKey.Name)
	} else {
		log.Info("SSH key already exists in Contabo API", "sshKeyContaboName", sshKeyContaboName, "sshKeyID", resp.JSON200.Data[0].SecretId)
		sshKey = &resp.JSON200.Data[0]
	}

	// Update status with SSH key info (only if not already set or different)
	if contaboCluster.Status.SshKey == nil || contaboCluster.Status.SshKey.SecretId != int64(sshKey.SecretId) {
		contaboCluster.Status.SshKey = &infrastructurev1beta2.ContaboSshKeyStatus{
			Name:     sshKey.Name,
			SecretId: int64(sshKey.SecretId),
			Value:    sshKey.Value,
		}
		log.Info("Updated SSH key status", "sshKeyName", sshKey.Name, "sshKeyID", sshKey.SecretId)
	}

	log.Info("SSH key reconciled successfully", "sshKeyContaboName", sshKey.Name, "sshKeyId", sshKey.SecretId)

	return ctrl.Result{}, nil
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
