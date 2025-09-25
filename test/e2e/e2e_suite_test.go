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
	"testing"
	"time"

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

var _ = BeforeSuite(func() {
	By("building the manager(Operator) image and loading into Kind")
	cmd := exec.Command("make", "docker-build-kind", fmt.Sprintf("IMG=%s", projectImage))
	_, err := utils.Run(cmd)
	ExpectWithOffset(1, err).NotTo(HaveOccurred(), "Failed to build and load the manager(Operator) image")

	// The tests-e2e are intended to run on a temporary cluster that is created and destroyed for testing.
	// To prevent errors when tests run in environments with CertManager already installed,
	// we check for its presence before execution.
	// Setup CertManager FIRST (CAPI operator needs it)
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

	// Install CAPI operator (after cert-manager)
	By("installing CAPI operator via Helm")
	err = utils.InstallCAPIOperator()
	ExpectWithOffset(1, err).NotTo(HaveOccurred(), "Failed to install CAPI operator")

	// Install Cluster API core components (after CAPI operator is ready)
	By("installing Cluster API core components")
	err = utils.InstallClusterAPICRDs()
	ExpectWithOffset(1, err).NotTo(HaveOccurred(), "Failed to install Cluster API core components")
})

var _ = AfterSuite(func() {
	// Clean up any remaining Cluster API resources before uninstalling operators
	_, _ = fmt.Fprintf(GinkgoWriter, "Cleaning up remaining Cluster API resources...\n")
	
	// Delete all Cluster resources across all namespaces
	cmd := exec.Command("kubectl", "delete", "clusters", "--all", "--all-namespaces", "--ignore-not-found=true", "--timeout=120s")
	_, _ = utils.Run(cmd)
	
	// Delete all ContaboCluster resources across all namespaces
	cmd = exec.Command("kubectl", "delete", "contaboclusters", "--all", "--all-namespaces", "--ignore-not-found=true", "--timeout=120s")
	_, _ = utils.Run(cmd)
	
	// Delete any Machine and MachineSet resources that might exist
	cmd = exec.Command("kubectl", "delete", "machines", "--all", "--all-namespaces", "--ignore-not-found=true", "--timeout=60s")
	_, _ = utils.Run(cmd)
	
	cmd = exec.Command("kubectl", "delete", "machinesets", "--all", "--all-namespaces", "--ignore-not-found=true", "--timeout=60s")
	_, _ = utils.Run(cmd)
	
	// Wait a bit for finalizers to be processed
	time.Sleep(10 * time.Second)
	
	// Force delete any remaining resources if they're stuck
	cmd = exec.Command("kubectl", "patch", "clusters", "--all", "--all-namespaces", "--type=merge", "-p", `{"metadata":{"finalizers":null}}`, "--ignore-not-found=true")
	_, _ = utils.Run(cmd)
	
	cmd = exec.Command("kubectl", "patch", "contaboclusters", "--all", "--all-namespaces", "--type=merge", "-p", `{"metadata":{"finalizers":null}}`, "--ignore-not-found=true")
	_, _ = utils.Run(cmd)
	
	cmd = exec.Command("kubectl", "patch", "machines", "--all", "--all-namespaces", "--type=merge", "-p", `{"metadata":{"finalizers":null}}`, "--ignore-not-found=true")
	_, _ = utils.Run(cmd)
	
	cmd = exec.Command("kubectl", "patch", "machinesets", "--all", "--all-namespaces", "--type=merge", "-p", `{"metadata":{"finalizers":null}}`, "--ignore-not-found=true")
	_, _ = utils.Run(cmd)

	// Debug: List remaining CAPI resources before deletion
	_, _ = fmt.Fprintf(GinkgoWriter, "Checking remaining CAPI resources...\n")
	cmd = exec.Command("kubectl", "get", "coreproviders", "--all-namespaces", "--ignore-not-found")
	if output, err := utils.Run(cmd); err == nil && output != "" {
		_, _ = fmt.Fprintf(GinkgoWriter, "Remaining CoreProviders:\n%s\n", output)
	}
	
	cmd = exec.Command("kubectl", "get", "validatingwebhookconfigurations", "-o", "name")
	if output, err := utils.Run(cmd); err == nil {
		_, _ = fmt.Fprintf(GinkgoWriter, "Webhook configurations:\n%s\n", output)
	}

	// Uninstall CAPI operator and components after resources are cleaned up
	_, _ = fmt.Fprintf(GinkgoWriter, "Uninstalling CAPI operator...\n")
	utils.UninstallCAPIOperator()

	// Teardown CertManager after the suite if not skipped and if it was not already installed
	if !skipCertManagerInstall && !isCertManagerAlreadyInstalled {
		_, _ = fmt.Fprintf(GinkgoWriter, "Uninstalling CertManager...\n")
		utils.UninstallCertManager()
	}
})
