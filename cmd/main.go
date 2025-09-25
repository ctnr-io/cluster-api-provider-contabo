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

package main

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/google/uuid"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/metrics/filters"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	infrastructurev1beta1 "github.com/ctnr-io/cluster-api-provider-contabo/api/v1beta1"
	"github.com/ctnr-io/cluster-api-provider-contabo/internal/controller"
	"github.com/ctnr-io/cluster-api-provider-contabo/pkg/contabo/auth"
	contaboclient "github.com/ctnr-io/cluster-api-provider-contabo/pkg/contabo/client"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(clusterv1.AddToScheme(scheme))

	utilruntime.Must(infrastructurev1beta1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

// nolint:gocyclo
func main() {
	var metricsAddr string
	var metricsCertPath, metricsCertName, metricsCertKey string
	var webhookCertPath, webhookCertName, webhookCertKey string
	var enableLeaderElection bool
	var probeAddr string
	var secureMetrics bool
	var enableHTTP2 bool
	var contaboClientID string
	var contaboClientSecret string
	var contaboAPIUser string
	var contaboAPIPassword string
	var leaderElectionID string
	var tlsOpts []func(*tls.Config)
	flag.StringVar(&metricsAddr, "metrics-bind-address", "0", "The address the metrics endpoint binds to. "+
		"Use :8443 for HTTPS or :8080 for HTTP, or leave as 0 to disable the metrics service.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.BoolVar(&secureMetrics, "metrics-secure", true,
		"If set, the metrics endpoint is served securely via HTTPS. Use --metrics-secure=false to use HTTP instead.")
	flag.StringVar(&webhookCertPath, "webhook-cert-path", "", "The directory that contains the webhook certificate.")
	flag.StringVar(&webhookCertName, "webhook-cert-name", "tls.crt", "The name of the webhook certificate file.")
	flag.StringVar(&webhookCertKey, "webhook-cert-key", "tls.key", "The name of the webhook key file.")
	flag.StringVar(&metricsCertPath, "metrics-cert-path", "",
		"The directory that contains the metrics server certificate.")
	flag.StringVar(&metricsCertName, "metrics-cert-name", "tls.crt", "The name of the metrics server certificate file.")
	flag.StringVar(&metricsCertKey, "metrics-cert-key", "tls.key", "The name of the metrics server key file.")
	flag.BoolVar(&enableHTTP2, "enable-http2", false,
		"If set, HTTP/2 will be enabled for the metrics and webhook servers")
	flag.StringVar(&contaboClientID, "contabo-client-id", "",
		"The Contabo OAuth2 client ID. Can also be set via CONTABO_CLIENT_ID environment variable.")
	flag.StringVar(&contaboClientSecret, "contabo-client-secret", "",
		"The Contabo OAuth2 client secret. Can also be set via CONTABO_CLIENT_SECRET environment variable.")
	flag.StringVar(&contaboAPIUser, "contabo-api-user", "",
		"The Contabo API username. Can also be set via CONTABO_API_USER environment variable.")
	flag.StringVar(&contaboAPIPassword, "contabo-api-password", "",
		"The Contabo API password. Can also be set via CONTABO_API_PASSWORD environment variable.")
	flag.StringVar(&leaderElectionID, "leader-election-id", "",
		"Leader election ID. If not specified, a dynamic ID will be generated based on namespace and controller name.")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	// Get Contabo OAuth2 credentials from environment if not provided via flags
	if contaboClientID == "" {
		contaboClientID = os.Getenv("CONTABO_CLIENT_ID")
	}
	if contaboClientSecret == "" {
		contaboClientSecret = os.Getenv("CONTABO_CLIENT_SECRET")
	}
	if contaboAPIUser == "" {
		contaboAPIUser = os.Getenv("CONTABO_API_USER")
	}
	if contaboAPIPassword == "" {
		contaboAPIPassword = os.Getenv("CONTABO_API_PASSWORD")
	}

	// Validate OAuth2 credentials
	if contaboClientID == "" || contaboClientSecret == "" || contaboAPIUser == "" || contaboAPIPassword == "" {
		setupLog.Error(fmt.Errorf("contabo OAuth2 credentials are required"),
			"set via flags or environment variables: CONTABO_CLIENT_ID, CONTABO_CLIENT_SECRET, CONTABO_API_USER, CONTABO_API_PASSWORD")
		os.Exit(1)
	}

	// Create OAuth2 token manager for automatic token refresh
	tokenManager := auth.NewTokenManager(contaboClientID, contaboClientSecret, contaboAPIUser, contaboAPIPassword)

	// Test initial token acquisition
	_, err := tokenManager.GetToken()
	if err != nil {
		setupLog.Error(err, "failed to get initial OAuth2 access token")
		os.Exit(1)
	}

	// Initialize Contabo OpenAPI client with token manager
	contaboClient, err := contaboclient.NewClient(
		"https://api.contabo.com",
		contaboclient.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			token, err := tokenManager.GetToken()
			if err != nil {
				return fmt.Errorf("failed to get access token: %w", err)
			}
			req.Header.Set("Authorization", "Bearer "+token)
			return nil
		}),
	)
	if err != nil {
		setupLog.Error(err, "unable to create Contabo API client")
		os.Exit(1)
	}

	// if the enable-http2 flag is false (the default), http/2 should be disabled
	// due to its vulnerabilities. More specifically, disabling http/2 will
	// prevent from being vulnerable to the HTTP/2 Stream Cancellation and
	// Rapid Reset CVEs. For more information see:
	// - https://github.com/advisories/GHSA-qppj-fm5r-hxr3
	// - https://github.com/advisories/GHSA-4374-p667-p6c8
	disableHTTP2 := func(c *tls.Config) {
		setupLog.Info("disabling http/2")
		c.NextProtos = []string{"http/1.1"}
	}

	if !enableHTTP2 {
		tlsOpts = append(tlsOpts, disableHTTP2)
	}

	// Initial webhook TLS options
	webhookTLSOpts := tlsOpts
	webhookServerOptions := webhook.Options{
		TLSOpts: webhookTLSOpts,
	}

	if len(webhookCertPath) > 0 {
		setupLog.Info("Initializing webhook certificate watcher using provided certificates",
			"webhook-cert-path", webhookCertPath, "webhook-cert-name", webhookCertName, "webhook-cert-key", webhookCertKey)

		webhookServerOptions.CertDir = webhookCertPath
		webhookServerOptions.CertName = webhookCertName
		webhookServerOptions.KeyName = webhookCertKey
	}

	webhookServer := webhook.NewServer(webhookServerOptions)

	// Metrics endpoint is enabled in 'config/default/kustomization.yaml'. The Metrics options configure the server.
	// More info:
	// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/metrics/server
	// - https://book.kubebuilder.io/reference/metrics.html
	metricsServerOptions := metricsserver.Options{
		BindAddress:   metricsAddr,
		SecureServing: secureMetrics,
		TLSOpts:       tlsOpts,
	}

	if secureMetrics {
		// FilterProvider is used to protect the metrics endpoint with authn/authz.
		// These configurations ensure that only authorized users and service accounts
		// can access the metrics endpoint. The RBAC are configured in 'config/rbac/kustomization.yaml'. More info:
		// https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/metrics/filters#WithAuthenticationAndAuthorization
		metricsServerOptions.FilterProvider = filters.WithAuthenticationAndAuthorization
	}

	// If the certificate is not specified, controller-runtime will automatically
	// generate self-signed certificates for the metrics server. While convenient for development and testing,
	// this setup is not recommended for production.
	//
	// TODO(user): If you enable certManager, uncomment the following lines:
	// - [METRICS-WITH-CERTS] at config/default/kustomization.yaml to generate and use certificates
	// managed by cert-manager for the metrics server.
	// - [PROMETHEUS-WITH-CERTS] at config/prometheus/kustomization.yaml for TLS certification.
	if len(metricsCertPath) > 0 {
		setupLog.Info("Initializing metrics certificate watcher using provided certificates",
			"metrics-cert-path", metricsCertPath, "metrics-cert-name", metricsCertName, "metrics-cert-key", metricsCertKey)

		metricsServerOptions.CertDir = metricsCertPath
		metricsServerOptions.CertName = metricsCertName
		metricsServerOptions.KeyName = metricsCertKey
	}

	// Generate dynamic leader election ID if not provided
	leaderElectionNamespace := getLeaderElectionNamespace()
	finalLeaderElectionID := generateLeaderElectionID(leaderElectionID, leaderElectionNamespace)
	setupLog.Info("Using leader election ID", "leaderElectionID", finalLeaderElectionID)

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                  scheme,
		Metrics:                 metricsServerOptions,
		WebhookServer:           webhookServer,
		HealthProbeBindAddress:  probeAddr,
		LeaderElection:          enableLeaderElection,
		LeaderElectionID:        finalLeaderElectionID,
		LeaderElectionNamespace: leaderElectionNamespace,
		// LeaderElectionReleaseOnCancel defines if the leader should step down voluntarily
		// when the Manager ends. This requires the binary to immediately end when the
		// Manager is stopped, otherwise, this setting is unsafe. Setting this significantly
		// speeds up voluntary leader transitions as the new leader don't have to wait
		// LeaseDuration time first.
		//
		// In the default scaffold provided, the program ends immediately after
		// the manager stops, so would be fine to enable this option. However,
		// if you are doing or is intended to do any operation such as perform cleanups
		// after the manager stops then its usage might be unsafe.
		// LeaderElectionReleaseOnCancel: true,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err := (&controller.ContaboClusterReconciler{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		Recorder:      mgr.GetEventRecorderFor("contabocluster-controller"),
		ContaboClient: contaboClient,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ContaboCluster")
		os.Exit(1)
	}
	if err := (&controller.ContaboMachineReconciler{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		Recorder:      mgr.GetEventRecorderFor("contabomachine-controller"),
		ContaboClient: contaboClient,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ContaboMachine")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

// generateLeaderElectionID creates a unique leader election ID using UUID hash
func generateLeaderElectionID(customID string, namespace string) string {
	// If a custom ID is provided, use it
	if customID != "" {
		return customID
	}

	// Generate a UUID for uniqueness
	id := uuid.New()

	// Create hash of UUID to get shorter identifier (8 chars)
	hash := sha256.Sum256([]byte(id.String()))
	shortID := fmt.Sprintf("%x", hash)[:8]

	// Create leader election ID with optional namespace and short hash
	if namespace != "" {
		return fmt.Sprintf("contabo-%s-%s.cluster.x-k8s.io", namespace, shortID)
	}
	return fmt.Sprintf("contabo-%s.cluster.x-k8s.io", shortID)
}

// getLeaderElectionNamespace determines the namespace for leader election
func getLeaderElectionNamespace() string {
	// Try to get namespace from environment (set by Kubernetes deployment)
	if namespace := os.Getenv("CONTROLLER_NAMESPACE"); namespace != "" {
		return namespace
	}

	// Fallback: try to read from service account namespace
	if data, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace"); err == nil {
		return string(data)
	}

	// Return empty string to use default namespace behavior
	return ""
}
