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
	"fmt"
	"os"
	"os/exec"
	"sync"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ctnr-io/cluster-api-provider-contabo/test/utils"
)

var (
	// Optional Environment Variables:
	// - CERT_MANAGER_INSTALL_SKIP=true: Skips CertManager installation during test setup.
	// These variables are useful if CertManager is already installed, avoiding
	// re-installation and conflicts.
	skipCertManagerInstall = os.Getenv("CERT_MANAGER_INSTALL_SKIP") == "true"
	// isCertManagerAlreadyInstalled will be set true when CertManager CRDs be found on the cluster
	isCertManagerAlreadyInstalled = false

	// projectImage is the name of the image which will be build and loaded
	// with the code source changes to be tested.
	projectImage = "cluster-api-provider-contabo:e2e-test"
)

// TestE2E runs the end-to-end (e2e) test suite for the project. These tests execute in an isolated,
// temporary environment to validate project changes with the purpose of being used in CI jobs.
// The default setup requires Kind, builds/loads the Manager Docker image locally, and installs
// CertManager.
func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	_, _ = fmt.Fprintf(GinkgoWriter, "Starting cluster-api-provider-contabo integration test suite\n")
	RunSpecs(t, "e2e suite")
}

func ParallelRun(cmds []*exec.Cmd) {
	var wg sync.WaitGroup
	for _, cmd := range cmds {
		wg.Add(1)
		go func(c *exec.Cmd) {
			defer wg.Done()
			_, err := utils.Run(c)
			if err != nil {
				_, _ = fmt.Fprintf(GinkgoWriter, "Error running command %v: %v\n", c, err)
			}
		}(cmd)
	}
	wg.Wait()
}

var _ = BeforeSuite(func() {
	var err error
	// Clean up any remaining Cluster API resources before uninstalling operators
	_, _ = fmt.Fprintf(GinkgoWriter, "Cleaning up remaining Cluster API resources...\n")

	// Delete all cluster resources across all namespaces
	ParallelRun([]*exec.Cmd{
		exec.Command("kubectl", "delete", "contaboclusters", "--all", "-n",  "contabo-e2e-test", "--ignore-not-found=true"),
		exec.Command("kubectl", "delete", "contabomachines", "--all", "-n",  "contabo-e2e-test", "--ignore-not-found=true"),
		exec.Command("kubectl", "delete", "machines", "--all", "-n",  "contabo-e2e-test", "--ignore-not-found=true"),
		exec.Command("kubectl", "delete", "machinesets", "--all", "-n",  "contabo-e2e-test", "--ignore-not-found=true"),
		exec.Command("kubectl", "delete", "clusters", "--all", "-n",  "contabo-e2e-test", "--ignore-not-found=true"),
	})
	
	var cmd *exec.Cmd

	// Remove clusterctl artifacts if any
	cmd = exec.Command("clusterctl", "delete", "--all", "--include-crd", "--include-namespace", "--ignore-not-found=true")
	_, _ = utils.Run(cmd)

	By("cleaning up the curl pod for metrics")
	cmd = exec.Command("kubectl", "delete", "pod", "curl-metrics", "-n", namespace)
	_, _ = utils.Run(cmd)

	By("cleaning up the metrics ClusterRoleBinding")
	cmd = exec.Command("kubectl", "delete", "clusterrolebinding", metricsRoleBindingName, "--ignore-not-found=true")
	_, _ = utils.Run(cmd)

	// By("cleaning up CAPI core components")
	// clusterctlPath := os.Getenv("CLUSTERCTL")
	// if clusterctlPath == "" {
	//   clusterctlPath = "clusterctl" // fallback to system clusterctl
	// }
	// cmd = exec.Command(clusterctlPath, "delete", "--all", "--include-crd", "--include-namespace")
	// _, _ = utils.Run(cmd)

	By("undeploying the controller-manager")
	cmd = exec.Command("make", "undeploy")
	_, _ = utils.Run(cmd)

	By("uninstalling CRDs")
	cmd = exec.Command("make", "uninstall")
	_, _ = utils.Run(cmd)

	By("removing manager namespace")
	cmd = exec.Command("kubectl", "delete", "ns", namespace)
	_, _ = utils.Run(cmd)

	// // Teardown CertManager after the suite if not skipped and if it was not already installed
	// if !skipCertManagerInstall && !isCertManagerAlreadyInstalled {
	// 	_, _ = fmt.Fprintf(GinkgoWriter, "Uninstalling CertManager...\n")
	// 	utils.UninstallCertManager()
	// }

	By("building the manager(Operator) image and loading into Kind")
	cmd = exec.Command("make", "docker-build-kind", fmt.Sprintf("IMG=%s", projectImage))
	_, err = utils.Run(cmd)
	ExpectWithOffset(1, err).NotTo(HaveOccurred(), "Failed to build and load the manager(Operator) image")

	// The tests-e2e are intended to run on a temporary cluster that is created and destroyed for testing.
	// To prevent errors when tests run in environments with CertManager already installed,
	// we check for its presence before execution.
	if !skipCertManagerInstall {
		By("checking if cert manager is installed already")
		isCertManagerAlreadyInstalled = utils.IsCertManagerCRDsInstalled()
		if !isCertManagerAlreadyInstalled {
			_, _ = fmt.Fprintf(GinkgoWriter, "Installing CertManager...\n")
			Expect(utils.InstallCertManager()).To(Succeed(), "Failed to install CertManager")
		} else {
			_, _ = fmt.Fprintf(GinkgoWriter, "WARNING: CertManager is already installed. Skipping installation...\n")
		}
	}

	By("installing Cluster API core components")
	clusterctlPath := os.Getenv("CLUSTERCTL")
	if clusterctlPath == "" {
		clusterctlPath = "clusterctl" // fallback to system clusterctl
	}

	// // Clean up any existing CAPI installations first
	// cmd = exec.Command(clusterctlPath, "delete", "--all", "--include-crd", "--include-namespace")
	// _, _ = utils.Run(cmd) // Ignore errors if nothing exists

	// Initialize with latest version that supports v1beta2
	cmd = exec.Command(clusterctlPath, "init", "--core", "cluster-api", "--bootstrap", "kubeadm", "--addon", "helm")
	_, err = utils.Run(cmd)
	Expect(err).NotTo(HaveOccurred(), "Failed to install CAPI core components")
})

var _ = AfterSuite(func() {
	// Delete all cluster resources across all namespaces
	// ParallelRun([]*exec.Cmd{
	// 	exec.Command("kubectl", "delete", "contaboclusters", "--all", "-n",  "contabo-e2e-test", "--ignore-not-found=true", "--timeout=30s"),
	// 	exec.Command("kubectl", "delete", "contabomachines", "--all", "-n",  "contabo-e2e-test", "--ignore-not-found=true", "--timeout=30s"),
	// 	exec.Command("kubectl", "delete", "machines", "--all", "-n",  "contabo-e2e-test", "--ignore-not-found=true", "--timeout=30s"),
	// 	exec.Command("kubectl", "delete", "machinesets", "--all", "-n",  "contabo-e2e-test", "--ignore-not-found=true", "--timeout=30s"),
	// 	exec.Command("kubectl", "delete", "clusters", "--all", "-n",  "contabo-e2e-test", "--ignore-not-found=true", "--timeout=30s"),
	// })
	
	// Wait for reconciliation
	// time.Sleep(30 * time.Second)
})
