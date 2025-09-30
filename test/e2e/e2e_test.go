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
	"strings"
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
	})

	// After all tests have been executed, clean up by undeploying the controller, uninstalling CRDs,
	// and deleting the namespace.
	AfterAll(func() {
		// By("cleaning up the curl pod for metrics")
		// cmd := exec.Command("kubectl", "delete", "pod", "curl-metrics", "-n", namespace)
		// _, _ = utils.Run(cmd)

		// By("cleaning up the metrics ClusterRoleBinding")
		// cmd = exec.Command("kubectl", "delete", "clusterrolebinding", metricsRoleBindingName, "--ignore-not-found=true")
		// _, _ = utils.Run(cmd)

		// // By("cleaning up CAPI core components")
		// // clusterctlPath := os.Getenv("CLUSTERCTL")
		// // if clusterctlPath == "" {
		// //   clusterctlPath = "clusterctl" // fallback to system clusterctl
		// // }
		// // cmd = exec.Command(clusterctlPath, "delete", "--all", "--include-crd", "--include-namespace")
		// // _, _ = utils.Run(cmd)

		// By("undeploying the controller-manager")
		// cmd = exec.Command("make", "undeploy")
		// _, _ = utils.Run(cmd)

		// By("uninstalling CRDs")
		// cmd = exec.Command("make", "uninstall")
		// _, _ = utils.Run(cmd)

		// By("removing manager namespace")
		// cmd = exec.Command("kubectl", "delete", "ns", namespace)
		// _, _ = utils.Run(cmd)
	})

	// After each test, check for failures and collect logs, events,
	// and pod descriptions for debugging.
	AfterEach(func() {
		specReport := CurrentSpecReport()
		if specReport.Failed() {
			// var cmd *exec.Cmd
			// By("Fetching controller manager pod logs")
			// cmd = exec.Command("kubectl", "logs", controllerPodName, "-n", namespace)
			// controllerLogs, err := utils.Run(cmd)
			// if err == nil {
			// 	_, _ = fmt.Fprintf(GinkgoWriter, "Controller logs:\n %s", controllerLogs)
			// } else {
			// 	_, _ = fmt.Fprintf(GinkgoWriter, "Failed to get Controller logs: %s", err)
			// }

			// By("Fetching Kubernetes events")
			// cmd = exec.Command("kubectl", "get", "events", "-n", namespace, "--sort-by=.lastTimestamp")
			// eventsOutput, err := utils.Run(cmd)
			// if err == nil {
			// 	_, _ = fmt.Fprintf(GinkgoWriter, "Kubernetes events:\n%s", eventsOutput)
			// } else {
			// 	_, _ = fmt.Fprintf(GinkgoWriter, "Failed to get Kubernetes events: %s", err)
			// }

			// By("Fetching curl-metrics logs")
			// cmd = exec.Command("kubectl", "logs", "curl-metrics", "-n", namespace)
			// metricsOutput, err := utils.Run(cmd)
			// if err == nil {
			// 	_, _ = fmt.Fprintf(GinkgoWriter, "Metrics logs:\n %s", metricsOutput)
			// } else {
			// 	_, _ = fmt.Fprintf(GinkgoWriter, "Failed to get curl-metrics logs: %s", err)
			// }

			// By("Fetching controller manager pod description")
			// cmd = exec.Command("kubectl", "describe", "pod", controllerPodName, "-n", namespace)
			// podDescription, err := utils.Run(cmd)
			// if err == nil {
			// 	fmt.Println("Pod description:\n", podDescription)
			// } else {
			// 	fmt.Println("Failed to describe controller pod")
			// }
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

	Describe("Cluster Lifecycle", func() {
		const testNamespace = "contabo-e2e-test"

		// Helper function to run kubectl commands
		kubectl := func(args ...string) (string, error) {
			return utils.Run(exec.Command("kubectl", args...))
		}

		// Helper to apply manifests directly via kubectl
		applyManifest := func(manifest string) {
			manifestFile := "/tmp/contabo-e2e-test.yaml"
			Expect(os.WriteFile(manifestFile, []byte(manifest), 0644)).To(Succeed())
			defer os.Remove(manifestFile)

			_, err := kubectl("apply", "-f", manifestFile)
			Expect(err).NotTo(HaveOccurred())
		}

		// Helper to wait for resource existence
		waitForResource := func(resource, name string, timeout time.Duration) {
			Eventually(func() error {
				_, err := kubectl("get", resource, name, "-n", testNamespace)
				return err
			}, timeout, 2*time.Second).Should(Succeed())
		}

		// Helper to wait for resource to be ready with proper conditions
		waitForResourceReady := func(resource, name, readyConditionPath string, timeout time.Duration) {
			Eventually(func() string {
				output, _ := kubectl("get", resource, name, "-n", testNamespace, "-o", "jsonpath="+readyConditionPath)
				return strings.TrimSpace(output)
			}, timeout, 5*time.Second).Should(Equal("True"))
		}

		// Helper to check conditions are properly set
		checkConditions := func(resource, name string) {
			Eventually(func() string {
				output, _ := kubectl("get", resource, name, "-n", testNamespace, "-o", "jsonpath={.status.conditions}")
				return output
			}, 30*time.Second, 2*time.Second).ShouldNot(BeEmpty())
		}

		BeforeEach(func() {
			By("cleaning up test resources")
			kubectl("delete", "namespace", testNamespace, "--ignore-not-found=true", "--timeout=60s")
			By("creating test namespace")
			kubectl("create", "namespace", testNamespace)
		})

		It("creates ContaboCluster and control plane with V76 product", func() {
			clusterManifest, err := os.ReadFile("test/fixtures/cluster.yaml")
			if err != nil {
				panic(err)
			}
			
			By("applying cluster and control plane manifests")
			applyManifest(string(clusterManifest))

			By("verifying cluster resources are created")
			waitForResource("cluster", "test-cluster", 30*time.Second)
			waitForResource("contabocluster", "test-cluster", 30*time.Second)
			waitForResource("machine", "test-control-plane", 30*time.Second)
			waitForResource("contabomachine", "test-control-plane", 30*time.Second)

			By("verifying V76 product type is configured for control plane machine")
			Eventually(func() string {
				output, _ := kubectl("get", "contabomachine", "test-control-plane", "-n", testNamespace, "-o", "jsonpath={.spec.instance.productId}")
				return output
			}, 30*time.Second, 2*time.Second).Should(Equal("V76"))

			By("verifying ContaboCluster conditions are properly set")
			checkConditions("contabocluster", "test-cluster")

			By("waiting for ContaboCluster to be ready")
			waitForResourceReady("contabocluster", "test-cluster", "{.status.conditions[?(@.type=='Ready')].status}", 3*time.Minute)

			By("verifying ContaboCluster ready status is true")
			Eventually(func() string {
				output, _ := kubectl("get", "contabocluster", "test-cluster", "-n", testNamespace, "-o", "jsonpath={.status.ready}")
				return strings.TrimSpace(output)
			}, 30*time.Second, 5*time.Second).Should(Equal("true"))

			By("verifying ContaboMachine conditions are properly set")
			checkConditions("contabomachine", "test-control-plane")

			By("waiting for ContaboMachine to be ready")
			waitForResourceReady("contabomachine", "test-control-plane", "{.status.conditions[?(@.type=='Ready')].status}", 3*time.Minute)

			By("verifying network infrastructure is ready")
			Eventually(func() string {
				output, _ := kubectl("get", "contabocluster", "test-cluster", "-n", testNamespace, "-o", "jsonpath={.status.conditions[?(@.type=='NetworkInfrastructureReady')].status}")
				return strings.TrimSpace(output)
			}, 2*time.Minute, 10*time.Second).Should(Equal("True"))

			By("checking controller logs show successful reconciliation")
			Eventually(func() string {
				logs, _ := kubectl("logs", controllerPodName, "-n", namespace, "--tail=50")
				return logs
			}, 1*time.Minute, 5*time.Second).Should(And(
				ContainSubstring("test-cluster"),
				ContainSubstring("reconciled successfully"),
			))
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
