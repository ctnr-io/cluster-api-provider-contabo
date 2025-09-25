//go:build e2e
// +build e2e

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

package e2e

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ctnr-io/cluster-api-provider-contabo/test/utils"
)

// namespace where the project is deployed in
const namespace = "cluster-api-provider-contabo-system"

// serviceAccountName created for the project
const serviceAccountName = "cluster-api-provider-contabo-controller-manager"

// metricsServiceName is the name of the metrics service of the project
const metricsServiceName = "cluster-api-provider-contabo-controller-manager-metrics-service"

// metricsRoleBindingName is the name of the RBAC that will be created to allow get the metrics data
const metricsRoleBindingName = "cluster-api-provider-contabo-metrics-binding"

var _ = Describe("Manager", Ordered, func() {
	var controllerPodName string

	// Before running the tests, set up the environment by creating the namespace,
	// enforce the restricted security policy to the namespace, installing CRDs,
	// and deploying the controller.
	BeforeAll(func() {
		By("creating manager namespace")
		cmd := exec.Command("kubectl", "create", "ns", namespace)
		_, err := utils.Run(cmd)
		Expect(err).NotTo(HaveOccurred(), "Failed to create namespace")

		By("labeling the namespace to enforce the restricted security policy")
		cmd = exec.Command("kubectl", "label", "--overwrite", "ns", namespace,
			"pod-security.kubernetes.io/enforce=restricted")
		_, err = utils.Run(cmd)
		Expect(err).NotTo(HaveOccurred(), "Failed to label namespace with restricted policy")

		By("installing CRDs")
		cmd = exec.Command("make", "install")
		_, err = utils.Run(cmd)
		Expect(err).NotTo(HaveOccurred(), "Failed to install CRDs")

		By("deploying the controller-manager")
		cmd = exec.Command("make", "deploy", fmt.Sprintf("IMG=%s", projectImage))
		_, err = utils.Run(cmd)
		Expect(err).NotTo(HaveOccurred(), "Failed to deploy the controller-manager")

		By("skipping CAPI core installation - testing provider controller in isolation")
	})

	// After all tests have been executed, clean up by undeploying the controller, uninstalling CRDs,
	// and deleting the namespace.
	AfterAll(func() {
		By("cleaning up the curl pod for metrics")
		cmd := exec.Command("kubectl", "delete", "pod", "curl-metrics", "-n", namespace)
		_, _ = utils.Run(cmd)

		By("cleaning up the metrics ClusterRoleBinding")
		cmd = exec.Command("kubectl", "delete", "clusterrolebinding", metricsRoleBindingName, "--ignore-not-found=true")
		_, _ = utils.Run(cmd)

		By("cleaning up CAPI core components")
		cmd = exec.Command("clusterctl", "delete", "--all", "--include-crd", "--include-namespace")
		_, _ = utils.Run(cmd)

		By("undeploying the controller-manager")
		cmd = exec.Command("make", "undeploy")
		_, _ = utils.Run(cmd)

		By("uninstalling CRDs")
		cmd = exec.Command("make", "uninstall")
		_, _ = utils.Run(cmd)

		By("removing manager namespace")
		cmd = exec.Command("kubectl", "delete", "ns", namespace)
		_, _ = utils.Run(cmd)
	})

	// After each test, check for failures and collect logs, events,
	// and pod descriptions for debugging.
	AfterEach(func() {
		specReport := CurrentSpecReport()
		if specReport.Failed() {
			By("Fetching controller manager pod logs")
			cmd := exec.Command("kubectl", "logs", controllerPodName, "-n", namespace)
			controllerLogs, err := utils.Run(cmd)
			if err == nil {
				_, _ = fmt.Fprintf(GinkgoWriter, "Controller logs:\n %s", controllerLogs)
			} else {
				_, _ = fmt.Fprintf(GinkgoWriter, "Failed to get Controller logs: %s", err)
			}

			By("Fetching Kubernetes events")
			cmd = exec.Command("kubectl", "get", "events", "-n", namespace, "--sort-by=.lastTimestamp")
			eventsOutput, err := utils.Run(cmd)
			if err == nil {
				_, _ = fmt.Fprintf(GinkgoWriter, "Kubernetes events:\n%s", eventsOutput)
			} else {
				_, _ = fmt.Fprintf(GinkgoWriter, "Failed to get Kubernetes events: %s", err)
			}

			By("Fetching curl-metrics logs")
			cmd = exec.Command("kubectl", "logs", "curl-metrics", "-n", namespace)
			metricsOutput, err := utils.Run(cmd)
			if err == nil {
				_, _ = fmt.Fprintf(GinkgoWriter, "Metrics logs:\n %s", metricsOutput)
			} else {
				_, _ = fmt.Fprintf(GinkgoWriter, "Failed to get curl-metrics logs: %s", err)
			}

			By("Fetching controller manager pod description")
			cmd = exec.Command("kubectl", "describe", "pod", controllerPodName, "-n", namespace)
			podDescription, err := utils.Run(cmd)
			if err == nil {
				fmt.Println("Pod description:\n", podDescription)
			} else {
				fmt.Println("Failed to describe controller pod")
			}
		}
	})

	SetDefaultEventuallyTimeout(2 * time.Minute)
	SetDefaultEventuallyPollingInterval(time.Second)

	Context("Manager", func() {
		It("should run successfully", func() {
			By("validating that the controller-manager pod is running as expected")
			verifyControllerUp := func(g Gomega) {
				// Get the name of the controller-manager pod
				cmd := exec.Command("kubectl", "get",
					"pods", "-l", "control-plane=controller-manager",
					"-o", "go-template={{ range .items }}"+
						"{{ if not .metadata.deletionTimestamp }}"+
						"{{ .metadata.name }}"+
						"{{ \"\\n\" }}{{ end }}{{ end }}",
					"-n", namespace,
				)

				podOutput, err := utils.Run(cmd)
				g.Expect(err).NotTo(HaveOccurred(), "Failed to retrieve controller-manager pod information")
				podNames := utils.GetNonEmptyLines(podOutput)
				g.Expect(podNames).To(HaveLen(1), "expected 1 controller pod running")
				controllerPodName = podNames[0]
				g.Expect(controllerPodName).To(ContainSubstring("controller-manager"))

				// Validate the pod's status
				cmd = exec.Command("kubectl", "get",
					"pods", controllerPodName, "-o", "jsonpath={.status.phase}",
					"-n", namespace,
				)
				output, err := utils.Run(cmd)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(output).To(Equal("Running"), "Incorrect controller-manager pod status")
			}
			Eventually(verifyControllerUp).Should(Succeed())
		})

		It("should ensure the metrics endpoint is serving metrics", func() {
			By("creating a ClusterRoleBinding for the service account to allow access to metrics")
			cmd := exec.Command("kubectl", "create", "clusterrolebinding", metricsRoleBindingName,
				"--clusterrole=cluster-api-provider-contabo-metrics-reader",
				fmt.Sprintf("--serviceaccount=%s:%s", namespace, serviceAccountName),
			)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred(), "Failed to create ClusterRoleBinding")

			By("validating that the metrics service is available")
			cmd = exec.Command("kubectl", "get", "service", metricsServiceName, "-n", namespace)
			_, err = utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred(), "Metrics service should exist")

			By("getting the service account token")
			token, err := serviceAccountToken()
			Expect(err).NotTo(HaveOccurred())
			Expect(token).NotTo(BeEmpty())

			By("waiting for the metrics endpoint to be ready")
			verifyMetricsEndpointReady := func(g Gomega) {
				cmd := exec.Command("kubectl", "get", "endpoints", metricsServiceName, "-n", namespace)
				output, err := utils.Run(cmd)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(output).To(ContainSubstring("8443"), "Metrics endpoint is not ready")
			}
			Eventually(verifyMetricsEndpointReady).Should(Succeed())

			By("verifying that the controller manager is serving the metrics server")
			verifyMetricsServerStarted := func(g Gomega) {
				cmd := exec.Command("kubectl", "logs", controllerPodName, "-n", namespace)
				output, err := utils.Run(cmd)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(output).To(ContainSubstring("controller-runtime.metrics\tServing metrics server"),
					"Metrics server not yet started")
			}
			Eventually(verifyMetricsServerStarted).Should(Succeed())

			By("creating the curl-metrics pod to access the metrics endpoint")
			cmd = exec.Command("kubectl", "run", "curl-metrics", "--restart=Never",
				"--namespace", namespace,
				"--image=curlimages/curl:latest",
				"--overrides",
				fmt.Sprintf(`{
					"spec": {
						"containers": [{
							"name": "curl",
							"image": "curlimages/curl:latest",
							"command": ["/bin/sh", "-c"],
							"args": ["curl -v -k -H 'Authorization: Bearer %s' https://%s.%s.svc.cluster.local:8443/metrics"],
							"securityContext": {
								"readOnlyRootFilesystem": true,
								"allowPrivilegeEscalation": false,
								"capabilities": {
									"drop": ["ALL"]
								},
								"runAsNonRoot": true,
								"runAsUser": 1000,
								"seccompProfile": {
									"type": "RuntimeDefault"
								}
							}
						}],
						"serviceAccountName": "%s"
					}
				}`, token, metricsServiceName, namespace, serviceAccountName))
			_, err = utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred(), "Failed to create curl-metrics pod")

			By("waiting for the curl-metrics pod to complete.")
			verifyCurlUp := func(g Gomega) {
				cmd := exec.Command("kubectl", "get", "pods", "curl-metrics",
					"-o", "jsonpath={.status.phase}",
					"-n", namespace)
				output, err := utils.Run(cmd)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(output).To(Equal("Succeeded"), "curl pod in wrong status")
			}
			Eventually(verifyCurlUp, 5*time.Minute).Should(Succeed())

			By("getting the metrics by checking curl-metrics logs")
			metricsOutput := getMetricsOutput()
			Expect(metricsOutput).To(ContainSubstring(
				"controller_runtime_reconcile_total",
			))
		})

		// +kubebuilder:scaffold:e2e-webhooks-checks
	})

	Describe("Cluster Lifecycle", Ordered, func() {
		var testNamespace string
		var manifestFile string

		BeforeEach(func() {
			testNamespace = "contabo-cluster-e2e-test"
		})

		AfterEach(func() {
			By("cleaning up cluster resources gracefully")

			if manifestFile != "" {
				// First try to delete via manifest file
				cmd := exec.Command("kubectl", "delete", "-f", manifestFile, "--ignore-not-found=true", "--timeout=60s")
				_, _ = utils.Run(cmd)
				_ = os.Remove(manifestFile)
			}

			// Give some time for finalizers to be processed
			time.Sleep(5 * time.Second)

			// If resources are still stuck, try direct deletion
			cmd := exec.Command("kubectl", "delete", "cluster", "test-cluster-e2e", "-n", testNamespace, "--ignore-not-found=true", "--timeout=30s")
			_, _ = utils.Run(cmd)

			cmd = exec.Command("kubectl", "delete", "contabocluster", "test-cluster-e2e", "-n", testNamespace, "--ignore-not-found=true", "--timeout=30s")
			_, _ = utils.Run(cmd)

			// Force remove finalizers if still stuck
			cmd = exec.Command("kubectl", "patch", "cluster", "test-cluster-e2e", "-n", testNamespace, "--type=merge", "-p", `{"metadata":{"finalizers":null}}`, "--ignore-not-found=true")
			_, _ = utils.Run(cmd)

			cmd = exec.Command("kubectl", "patch", "contabocluster", "test-cluster-e2e", "-n", testNamespace, "--type=merge", "-p", `{"metadata":{"finalizers":null}}`, "--ignore-not-found=true")
			_, _ = utils.Run(cmd)

			// Clean up test namespace
			cmd = exec.Command("kubectl", "delete", "namespace", testNamespace, "--ignore-not-found=true", "--timeout=120s")
			_, _ = utils.Run(cmd)
		})

		It("should successfully reconcile ContaboCluster lifecycle", func() {
			By("creating a test namespace for cluster resources")
			cmd := exec.Command("kubectl", "create", "namespace", testNamespace)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred(), "Failed to create test namespace")

			By("applying a complete ContaboCluster with private network configuration")
			clusterManifest := fmt.Sprintf(`
apiVersion: cluster.x-k8s.io/v1beta1
kind: Cluster
metadata:
  name: test-cluster-e2e
  namespace: %s
spec:
  clusterNetwork:
    pods:
      cidrBlocks:
      - 192.168.0.0/16
  infrastructureRef:
    apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
    kind: ContaboCluster
    name: test-cluster-e2e
    namespace: %s
---
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: ContaboCluster
metadata:
  name: test-cluster-e2e
  namespace: %s
spec:
  region: "EU"
  controlPlaneEndpoint:
    host: "10.0.0.100"
    port: 6443
  network:
    privateNetworks:
    - name: "test-e2e-private-network"
`, testNamespace, testNamespace, testNamespace)

			manifestFile = "/tmp/contabo-cluster-e2e-test.yaml"
			err = os.WriteFile(manifestFile, []byte(clusterManifest), 0644)
			Expect(err).NotTo(HaveOccurred(), "Failed to write cluster manifest")

			cmd = exec.Command("kubectl", "apply", "-f", manifestFile)
			_, err = utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred(), "Failed to apply cluster manifest")

			By("debugging CAPI components status")
			cmd = exec.Command("kubectl", "get", "pods", "--all-namespaces", "-l", "cluster-api")
			capiPods, err := utils.Run(cmd)
			if err == nil {
				fmt.Printf("CAPI pods found:\n%s\n", capiPods)
			} else {
				fmt.Printf("No CAPI pods found or error: %s\n", err)
			}

			By("checking for CAPI CRDs")
			cmd = exec.Command("kubectl", "get", "crd", "-o", "name")
			crds, err := utils.Run(cmd)
			if err == nil {
				if len(crds) > 0 {
					fmt.Printf("Found CRDs, checking for CAPI ones...\n")
					cmd = exec.Command("kubectl", "get", "crd", "-o", "name")
					output, _ := utils.Run(cmd)
					if output != "" {
						fmt.Printf("Sample CRDs: %s\n", output[:min(500, len(output))])
					}
				}
			}

			By("verifying the Cluster resource is created")
			Eventually(func(g Gomega) {
				cmd := exec.Command("kubectl", "get", "cluster", "test-cluster-e2e", "-n", testNamespace, "-o", "jsonpath={.metadata.name}")
				output, err := utils.Run(cmd)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(output).To(Equal("test-cluster-e2e"), "Cluster resource should be created")
			}, 30*time.Second, 5*time.Second).Should(Succeed())

			By("verifying the ContaboCluster resource is created")
			Eventually(func(g Gomega) {
				cmd := exec.Command("kubectl", "get", "contabocluster", "test-cluster-e2e", "-n", testNamespace, "-o", "jsonpath={.metadata.name}")
				output, err := utils.Run(cmd)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(output).To(Equal("test-cluster-e2e"), "ContaboCluster resource should be created")
			}, 30*time.Second, 5*time.Second).Should(Succeed())

			By("waiting for ContaboCluster controller to process the resource")
			Eventually(func(g Gomega) {
				cmd := exec.Command("kubectl", "get", "contabocluster", "test-cluster-e2e",
					"-n", testNamespace, "-o", "jsonpath={.status}")
				output, err := utils.Run(cmd)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(output).ToNot(BeEmpty(), "ContaboCluster status should be populated by controller")
			}, 2*time.Minute, 10*time.Second).Should(Succeed())

			By("verifying private network configuration is processed")
			Eventually(func(g Gomega) {
				cmd := exec.Command("kubectl", "get", "contabocluster", "test-cluster-e2e",
					"-n", testNamespace, "-o", "jsonpath={.spec.network.privateNetworks[0].name}")
				output, err := utils.Run(cmd)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(output).To(Equal("test-e2e-private-network"), "Private network specification should be correct")
			}, 30*time.Second, 5*time.Second).Should(Succeed())

			By("checking controller logs for cluster reconciliation")
			Eventually(func(g Gomega) {
				cmd := exec.Command("kubectl", "logs", controllerPodName, "-n", namespace, "--tail=100")
				logs, err := utils.Run(cmd)
				g.Expect(err).NotTo(HaveOccurred())

				// Look for cluster reconciliation activity
				g.Expect(logs).To(Or(
					ContainSubstring("Reconciling ContaboCluster"),
					ContainSubstring("reconciling contabocluster"),
					ContainSubstring("test-cluster-e2e"),
				), "Controller should show cluster reconciliation activity")
			}, 2*time.Minute, 10*time.Second).Should(Succeed())

			By("verifying cluster finalizer is added")
			Eventually(func(g Gomega) {
				cmd := exec.Command("kubectl", "get", "contabocluster", "test-cluster-e2e",
					"-n", testNamespace, "-o", "jsonpath={.metadata.finalizers}")
				output, err := utils.Run(cmd)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(output).To(ContainSubstring("contabocluster.infrastructure.cluster.x-k8s.io"),
					"ContaboCluster should have finalizer added")
			}, 1*time.Minute, 5*time.Second).Should(Succeed())

			By("checking network reconciliation logs")
			Eventually(func(g Gomega) {
				cmd := exec.Command("kubectl", "logs", controllerPodName, "-n", namespace, "--tail=200")
				logs, err := utils.Run(cmd)
				g.Expect(err).NotTo(HaveOccurred())

				// Look for network reconciliation
				g.Expect(logs).To(Or(
					ContainSubstring("reconciling private networks"),
					ContainSubstring("Network infrastructure reconciled"),
					ContainSubstring("failed to reconcile network"),
					ContainSubstring("Discovering/creating private network"),
				), "Controller should attempt network reconciliation")
			}, 2*time.Minute, 10*time.Second).Should(Succeed())

			By("verifying conditions are set on ContaboCluster")
			Eventually(func(g Gomega) {
				cmd := exec.Command("kubectl", "get", "contabocluster", "test-cluster-e2e",
					"-n", testNamespace, "-o", "jsonpath={.status.conditions}")
				output, err := utils.Run(cmd)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(output).ToNot(BeEmpty(), "ContaboCluster should have status conditions")
			}, 1*time.Minute, 10*time.Second).Should(Succeed())

			By("testing cluster deletion flow")
			cmd = exec.Command("kubectl", "delete", "cluster", "test-cluster-e2e", "-n", testNamespace)
			_, err = utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred(), "Failed to delete cluster")

			By("verifying ContaboCluster enters deletion flow")
			Eventually(func(g Gomega) {
				cmd := exec.Command("kubectl", "get", "contabocluster", "test-cluster-e2e",
					"-n", testNamespace, "-o", "jsonpath={.metadata.deletionTimestamp}")
				output, err := utils.Run(cmd)
				if err != nil {
					// Resource may already be deleted
					return
				}
				g.Expect(output).ToNot(BeEmpty(), "ContaboCluster should have deletion timestamp set")
			}, 1*time.Minute, 5*time.Second).Should(Succeed())

			By("checking deletion reconciliation logs")
			Eventually(func(g Gomega) {
				cmd := exec.Command("kubectl", "logs", controllerPodName, "-n", namespace, "--tail=100")
				logs, err := utils.Run(cmd)
				g.Expect(err).NotTo(HaveOccurred())

				// Look for deletion activity
				g.Expect(logs).To(Or(
					ContainSubstring("Reconciling ContaboCluster delete"),
					ContainSubstring("reconcileDelete"),
					ContainSubstring("deleteNetwork"),
				), "Controller should process cluster deletion")
			}, 2*time.Minute, 10*time.Second).Should(Succeed())

			By("verifying resources are eventually cleaned up")
			Eventually(func(g Gomega) {
				cmd := exec.Command("kubectl", "get", "contabocluster", "test-cluster-e2e", "-n", testNamespace)
				_, err := utils.Run(cmd)
				g.Expect(err).To(HaveOccurred(), "ContaboCluster should be deleted")
			}, 2*time.Minute, 10*time.Second).Should(Succeed())

			By("verifying controller metrics include ContaboCluster reconciliation")
			metricsOutput := getMetricsOutput()
			Expect(metricsOutput).To(ContainSubstring("contabocluster"),
				"Metrics should include ContaboCluster controller activity")
		})
	})
})

// serviceAccountToken returns a token for the specified service account in the given namespace.
// It uses the Kubernetes TokenRequest API to generate a token by directly sending a request
// and parsing the resulting token from the API response.
func serviceAccountToken() (string, error) {
	const tokenRequestRawString = `{
		"apiVersion": "authentication.k8s.io/v1",
		"kind": "TokenRequest"
	}`

	// Temporary file to store the token request
	secretName := fmt.Sprintf("%s-token-request", serviceAccountName)
	tokenRequestFile := filepath.Join("/tmp", secretName)
	err := os.WriteFile(tokenRequestFile, []byte(tokenRequestRawString), os.FileMode(0o644))
	if err != nil {
		return "", err
	}

	var out string
	verifyTokenCreation := func(g Gomega) {
		// Execute kubectl command to create the token
		cmd := exec.Command("kubectl", "create", "--raw", fmt.Sprintf(
			"/api/v1/namespaces/%s/serviceaccounts/%s/token",
			namespace,
			serviceAccountName,
		), "-f", tokenRequestFile)

		output, err := cmd.CombinedOutput()
		g.Expect(err).NotTo(HaveOccurred())

		// Parse the JSON output to extract the token
		var token tokenRequest
		err = json.Unmarshal(output, &token)
		g.Expect(err).NotTo(HaveOccurred())

		out = token.Status.Token
	}
	Eventually(verifyTokenCreation).Should(Succeed())

	return out, err
}

// getMetricsOutput retrieves and returns the logs from the curl pod used to access the metrics endpoint.
func getMetricsOutput() string {
	By("getting the curl-metrics logs")
	cmd := exec.Command("kubectl", "logs", "curl-metrics", "-n", namespace)
	metricsOutput, err := utils.Run(cmd)
	Expect(err).NotTo(HaveOccurred(), "Failed to retrieve logs from curl pod")
	Expect(metricsOutput).To(ContainSubstring("< HTTP/1.1 200 OK"))
	return metricsOutput
}

// tokenRequest is a simplified representation of the Kubernetes TokenRequest API response,
// containing only the token field that we need to extract.
type tokenRequest struct {
	Status struct {
		Token string `json:"token"`
	} `json:"status"`
}

// getProjectImage returns the project image name for the e2e tests
func getProjectImage() string {
	projectImage := os.Getenv("IMG")
	if projectImage == "" {
		projectImage = "cluster-api-provider-contabo:latest"
	}
	return projectImage
}
