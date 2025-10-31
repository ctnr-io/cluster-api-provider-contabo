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

		// Wait a bit in case user want to CTRL-C to inspect the cluster
		time.Sleep(10 * time.Second)

		// Delete all cluster resources
		ParallelRun([]*exec.Cmd{
			exec.Command("kubectl", "delete", "clusters", "--all", "-n", "contabo-e2e-test", "--ignore-not-found=true"),
			exec.Command("kubectl", "delete", "kubeadmcontrolplanes", "--all", "-n", "contabo-e2e-test", "--ignore-not-found=true"),
			exec.Command("kubectl", "delete", "machinedeployments", "--all", "-n", "contabo-e2e-test", "--ignore-not-found=true"),
			exec.Command("kubectl", "delete", "contaboclusters", "--all", "-n", "contabo-e2e-test", "--ignore-not-found=true"),
			exec.Command("kubectl", "delete", "contabomachines", "--all", "-n", "contabo-e2e-test", "--ignore-not-found=true"),
			exec.Command("kubectl", "delete", "contabomachinetemplates", "--all", "-n", "contabo-e2e-test", "--ignore-not-found=true"),
			exec.Command("kubectl", "delete", "machines", "--all", "-n", "contabo-e2e-test", "--ignore-not-found=true"),
			exec.Command("kubectl", "delete", "machinesets", "--all", "-n", "contabo-e2e-test", "--ignore-not-found=true"),
			exec.Command("kubectl", "delete", "kubeadmconfigs", "--all", "-n", "contabo-e2e-test", "--ignore-not-found=true"),
			exec.Command("kubectl", "delete", "kubeadmconfigtemplates", "--all", "-n", "contabo-e2e-test", "--ignore-not-found=true"),
			exec.Command("kubectl", "delete", "services", "--all", "-n", "contabo-e2e-test", "--ignore-not-found=true"),
			exec.Command("kubectl", "delete", "endpointslices", "--all", "-n", "contabo-e2e-test", "--ignore-not-found=true"),
		})

		// Wait for reconciliation
		time.Sleep(30 * time.Second)
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

	SetDefaultEventuallyTimeout(15 * time.Minute)
	SetDefaultEventuallyPollingInterval(5 * time.Second)

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
		cntb := func(args ...string) (string, error) {
			return utils.Run(exec.Command("cntb", args...))
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
			}, 120*time.Second, 10*time.Second).ShouldNot(BeEmpty())
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
			Expect(err).NotTo(HaveOccurred())

			// Get cluster UUID from ContaboCluster status
			By("getting cluster UUID from ContaboCluster")
			var clusterUUID string
			Eventually(func() string {
				output, _ := kubectl("get", "contabocluster", "test-cluster", "-n", testNamespace, "-o", "jsonpath={.status.clusterUUID}")
				clusterUUID = strings.TrimSpace(output)
				return clusterUUID
			}, 2*time.Minute, 5*time.Second).ShouldNot(BeEmpty())

			By("verifying cluster resources are created")
			waitForResource("cluster", "test-cluster", 30*time.Second)
			waitForResource("contabocluster", "test-cluster", 30*time.Second)
			waitForResource("kubeadmcontrolplane", "test-control-plane", 30*time.Second)
			waitForResource("machinedeployment", "test-worker", 30*time.Second)

			By("waiting for KubeadmControlPlane to create control plane machines")
			var controlPlaneMachineName string
			Eventually(func() string {
				output, _ := kubectl("get", "machines", "-n", testNamespace, "-l", "cluster.x-k8s.io/control-plane", "-o", "jsonpath={.items[0].metadata.name}")
				controlPlaneMachineName = strings.TrimSpace(output)
				return controlPlaneMachineName
			}, 2*time.Minute, 5*time.Second).ShouldNot(BeEmpty())

			By("waiting for MachineDeployment to create worker machines")
			var workerMachineName string
			Eventually(func() string {
				output, _ := kubectl("get", "machines", "-n", testNamespace, "-l", "cluster.x-k8s.io/deployment-name=test-worker", "-o", "jsonpath={.items[0].metadata.name}")
				workerMachineName = strings.TrimSpace(output)
				return workerMachineName
			}, 2*time.Minute, 5*time.Second).ShouldNot(BeEmpty())

			By("getting ContaboMachine names created by KCP and MachineDeployment")
			var controlPlaneContaboMachineName string
			Eventually(func() string {
				output, _ := kubectl("get", "contabomachines", "-n", testNamespace, "-l", "cluster.x-k8s.io/control-plane", "-o", "jsonpath={.items[0].metadata.name}")
				controlPlaneContaboMachineName = strings.TrimSpace(output)
				return controlPlaneContaboMachineName
			}, 2*time.Minute, 5*time.Second).ShouldNot(BeEmpty())

			var workerContaboMachineName string
			Eventually(func() string {
				output, _ := kubectl("get", "contabomachines", "-n", testNamespace, "-l", "cluster.x-k8s.io/deployment-name=test-worker", "-o", "jsonpath={.items[0].metadata.name}")
				workerContaboMachineName = strings.TrimSpace(output)
				return workerContaboMachineName
			}, 2*time.Minute, 5*time.Second).ShouldNot(BeEmpty())

			By("verifying V76 product type is configured for control plane machine")
			Eventually(func() string {
				output, _ := kubectl("get", "contabomachine", controlPlaneContaboMachineName, "-n", testNamespace, "-o", "jsonpath={.spec.instance.productId}")
				return output
			}, 2*time.Minute, 5*time.Second).Should(Equal("V76"))

			By("verifying ContaboCluster conditions are properly set")
			checkConditions("contabocluster", "test-cluster")

			By("waiting for ContaboCluster to be ready")
			waitForResourceReady("contabocluster", "test-cluster", "{.status.conditions[?(@.type=='Ready')].status}", 3*time.Minute)

			By("verifying ContaboCluster ready status is true")
			Eventually(func() string {
				output, _ := kubectl("get", "contabocluster", "test-cluster", "-n", testNamespace, "-o", "jsonpath={.status.ready}")
				return strings.TrimSpace(output)
			}, 2*time.Minute, 5*time.Second).Should(Equal("true"))

			By("verifying ContaboMachine conditions are properly set")
			checkConditions("contabomachine", controlPlaneContaboMachineName)
			checkConditions("contabomachine", workerContaboMachineName)

			By("waiting for ContaboMachine Initialization.Provisioned to be true (control plane)")
			Eventually(func() string {
				output, _ := kubectl("get", "contabomachine", controlPlaneContaboMachineName, "-n", testNamespace, "-o", "jsonpath={.status.initialization.provisioned}")
				return strings.ToLower(strings.TrimSpace(output))
			}, 5*time.Minute, 5*time.Second).Should(Equal("true"))

			By("waiting for ContaboMachine Initialization.Provisioned to be true (worker)")
			Eventually(func() string {
				output, _ := kubectl("get", "contabomachine", workerContaboMachineName, "-n", testNamespace, "-o", "jsonpath={.status.initialization.provisioned}")
				return strings.ToLower(strings.TrimSpace(output))
			}, 5*time.Minute, 5*time.Second).Should(Equal("true"))

			By("waiting for ContaboMachine Addresses to be set (control plane)")
			Eventually(func() string {
				output, _ := kubectl("get", "contabomachine", controlPlaneContaboMachineName, "-n", testNamespace, "-o", "jsonpath={.status.addresses}")
				return output
			}, 5*time.Minute, 5*time.Second).ShouldNot(BeEmpty())

			By("waiting for ContaboMachine Addresses to be set (worker)")
			Eventually(func() string {
				output, _ := kubectl("get", "contabomachine", workerContaboMachineName, "-n", testNamespace, "-o", "jsonpath={.status.addresses}")
				return output
			}, 5*time.Minute, 5*time.Second).ShouldNot(BeEmpty())

			By("waiting for ContaboMachine Ready to be true (control plane)")
			Eventually(func() string {
				output, _ := kubectl("get", "contabomachine", controlPlaneContaboMachineName, "-n", testNamespace, "-o", "jsonpath={.status.ready}")
				return strings.ToLower(strings.TrimSpace(output))
			}, 10*time.Minute, 5*time.Second).Should(Equal("true"))

			By("waiting for ContaboMachine Ready to be true (worker)")
			Eventually(func() string {
				output, _ := kubectl("get", "contabomachine", workerContaboMachineName, "-n", testNamespace, "-o", "jsonpath={.status.ready}")
				return strings.ToLower(strings.TrimSpace(output))
			}, 10*time.Minute, 5*time.Second).Should(Equal("true"))

			By("verifying control plane endpoint is set")
			var controlPlaneEndpointHost string
			Eventually(func() string {
				output, _ := kubectl("get", "contabocluster", "test-cluster", "-n", testNamespace, "-o", "jsonpath={.spec.controlPlaneEndpoint.host}")
				controlPlaneEndpointHost = strings.TrimSpace(output)
				return controlPlaneEndpointHost
			}, 2*time.Minute, 5*time.Second).ShouldNot(BeEmpty())

			By("verifying control plane endpoint Service is created")
			serviceName := fmt.Sprintf("%s-apiserver", "test-cluster")
			Eventually(func() error {
				_, err := kubectl("get", "service", serviceName, "-n", testNamespace)
				return err
			}, 2*time.Minute, 5*time.Second).Should(Succeed())

			By("verifying control plane endpoint Service is headless")
			Eventually(func() string {
				output, _ := kubectl("get", "service", serviceName, "-n", testNamespace, "-o", "jsonpath={.spec.clusterIP}")
				return strings.TrimSpace(output)
			}, 2*time.Minute, 5*time.Second).Should(Equal("None"))

			By("verifying control plane endpoint Service has correct port")
			Eventually(func() string {
				output, _ := kubectl("get", "service", serviceName, "-n", testNamespace, "-o", "jsonpath={.spec.ports[0].port}")
				return strings.TrimSpace(output)
			}, 2*time.Minute, 5*time.Second).Should(Equal("6443"))

			By("verifying control plane endpoint EndpointSlice is created")
			endpointSliceName := fmt.Sprintf("%s-apiserver", "test-cluster")
			Eventually(func() error {
				_, err := kubectl("get", "endpointslice", endpointSliceName, "-n", testNamespace)
				return err
			}, 2*time.Minute, 5*time.Second).Should(Succeed())


			By("verifying EndpointSlice has correct labels")
			Eventually(func() string {
				output, _ := kubectl("get", "endpointslice", endpointSliceName, "-n", testNamespace, "-o", "jsonpath={.metadata.labels.kubernetes\\.io/service-name}")
				return strings.TrimSpace(output)
			}, 2*time.Minute, 5*time.Second).Should(Equal(serviceName))

			By("waiting for workload cluster kubeconfig secret to be available")
			var kubeconfigB64 string
			Eventually(func() string {
				output, _ := kubectl("get", "secret", "test-cluster-kubeconfig", "-n", testNamespace, "-o", "jsonpath={.data.value}")
				kubeconfigB64 = output
				return output
			}, 10*time.Minute, 5*time.Second).ShouldNot(ContainSubstring("not found"))
			By("fetching workload cluster kubeconfig and connecting to cluster")

			// Write kubeconfig to file
			kubeconfigPath := "/tmp/test-cluster.kubeconfig"
			os.WriteFile(kubeconfigPath+".b64", []byte(kubeconfigB64), 0600)
			decoded, err := utils.Base64Decode(kubeconfigB64)
			Expect(err).NotTo(HaveOccurred())
			err = os.WriteFile(kubeconfigPath, []byte(decoded), 0600)
			Expect(err).NotTo(HaveOccurred())

			// Wait for control plane node to exists
			By("waiting for control plane node to exist in workload cluster")
			Eventually(func() string {
				cmd := exec.Command("kubectl", "--kubeconfig", kubeconfigPath, "get", "nodes", "-o", "jsonpath={.items[0].metadata.labels}")
				output, err := utils.Run(cmd)
				if err != nil {
					return ""
				}
				return output
			}).Should(ContainSubstring("node-role.kubernetes.io/control-plane"))

			// Install Cilium CNI
			By("installing Cilium CNI in workload cluster")
			cmd := exec.Command("./bin/cilium", "install", "--version", "1.18.2", "--kubeconfig", kubeconfigPath)
			_, err = utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())

			// Wait for all nodes to be Ready
			By("waiting for all nodes to be Ready in workload cluster")
			Eventually(func() string {
				cmd := exec.Command("kubectl", "--kubeconfig", kubeconfigPath, "get", "nodes", "-o", "jsonpath={.items[*].status.conditions[?(@.type=='Ready')].status}")
				output, err := utils.Run(cmd)
				if err != nil {
					return ""
				}
				return output
			}).Should(ContainSubstring("True True"))

			// Get the node names
			By("getting the node names in workload cluster")
			cmd = exec.Command("kubectl", "--kubeconfig", kubeconfigPath, "get", "nodes", "-o", "jsonpath={.items[*].metadata.name}")
			output, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred())
			nodes := strings.Fields(output)
			Expect(nodes).To(HaveLen(2))
			workerNodeName := nodes[0]
			controlPlaneNodeName := nodes[1]

			// Verify external IPs are set on all nodes
			By("verifying external IPs are set on all nodes")
			for _, nodeName := range nodes {
				Eventually(func() string {
					cmd := exec.Command("kubectl", "--kubeconfig", kubeconfigPath, "get", "node", nodeName, "-o", "jsonpath={.status.addresses[?(@.type==\"ExternalIP\")].address}")
					output, err := utils.Run(cmd)
					if err != nil {
						return ""
					}
					return strings.TrimSpace(output)
				}, 2*time.Minute, 5*time.Second).ShouldNot(BeEmpty(), fmt.Sprintf("Node %s should have an ExternalIP set", nodeName))
			}

			// Test deletion
			By("deleting the test cluster and machines")
			_, err = kubectl("delete", "machinedeployment", "test-worker", "-n", testNamespace, "--wait=false")
			Expect(err).NotTo(HaveOccurred())
			_, err = kubectl("delete", "kubeadmcontrolplane", "test-control-plane", "-n", testNamespace, "--wait=false")
			Expect(err).NotTo(HaveOccurred())
			_, err = kubectl("delete", "cluster", "test-cluster", "-n", testNamespace, "--wait=false")
			Expect(err).NotTo(HaveOccurred())

			By("waiting for ContaboMachines to be deleted")
			Eventually(func() error {
				_, err := kubectl("get", "contabomachine", workerContaboMachineName, "-n", testNamespace)
				return err
			}, 2*time.Minute, 5*time.Second).ShouldNot(Succeed())
			Eventually(func() error {
				_, err := kubectl("get", "contabomachine", controlPlaneContaboMachineName, "-n", testNamespace)
				return err
			}, 2*time.Minute, 5*time.Second).ShouldNot(Succeed())

			By("verifying control plane endpoint Service is deleted")
			Eventually(func() error {
				_, err := kubectl("get", "service", controlPlaneEndpointHost, "-n", testNamespace)
				return err
			}, 2*time.Minute, 5*time.Second).ShouldNot(Succeed())

			By("verifying control plane endpoint EndpointSlice is deleted")
			Eventually(func() error {
				_, err := kubectl("get", "endpointslice", endpointSliceName, "-n", testNamespace)
				return err
			}, 2*time.Minute, 5*time.Second).ShouldNot(Succeed())

			By("checking instance state with cntb CLI")
			cmd = exec.Command("cntb", "get", "instances", "--output", "json")
			output, err = utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred(), "cntb CLI failed: %s", output)

			type instance struct {
				ID          int64  `json:"id"`
				Name        string `json:"name"`
				DisplayName string `json:"displayName"`
			}
			var instances []instance
			err = json.Unmarshal([]byte(output), &instances)
			Expect(err).NotTo(HaveOccurred(), "Failed to parse cntb output: %s", output)

			// Both test-control-plane and test-worker should be reset
			for _, inst := range instances {
				if inst.Name == controlPlaneNodeName || inst.Name == workerNodeName {
					Expect(inst.DisplayName).To(Equal(""), "Instance %s should have display name empty", inst.Name)
				}
			}

			// Check that private network is deleted
			By("checking private network with cntb CLI")
			output, err = cntb("get", "privateNetworks", "--output", "json")
			Expect(err).NotTo(HaveOccurred(), "cntb CLI failed: %s", output)

			type privateNetwork struct {
				ID   int64  `json:"id"`
				Name string `json:"name"`
			}
			var privateNetworks []privateNetwork
			err = json.Unmarshal([]byte(output), &privateNetworks)
			Expect(err).NotTo(HaveOccurred(), "Failed to parse cntb output: %s", output)

			// Private network name should be [capc] <clusterUUID>
			expectedNetworkName := fmt.Sprintf("[capc] %s", clusterUUID)
			for _, network := range privateNetworks {
				Expect(network.Name).NotTo(Equal(expectedNetworkName),
					"Private network %s still exists after cluster deletion", expectedNetworkName)
			}

			// Check that SSH key is deleted
			By("checking SSH key with cntb CLI")
			output, err = cntb("get", "secrets", "--output", "json")
			Expect(err).NotTo(HaveOccurred(), "cntb CLI failed: %s", output)

			type sshKey struct {
				ID   int64  `json:"id"`
				Name string `json:"name"`
			}
			var sshKeys []sshKey
			err = json.Unmarshal([]byte(output), &sshKeys)
			Expect(err).NotTo(HaveOccurred(), "Failed to parse cntb output: %s", output)

			// SSH key name should be [capc] <clusterUUID>
			expectedKeyName := expectedNetworkName // same naming pattern as network
			for _, key := range sshKeys {
				Expect(key.Name).NotTo(Equal(expectedKeyName),
					"SSH key %s still exists after cluster deletion", expectedKeyName)
			}
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
