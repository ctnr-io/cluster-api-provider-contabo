# cluster-api-provider-contabo

A Kubernetes Cluster API provider for Contabo cloud infrastructure. This provider enables you to create and manage Kubernetes clusters on Contabo's VPS infrastructure using the Cluster API framework.

## Description

This Cluster API provider allows you to provision and manage Kubernetes clusters on Contabo's cloud infrastructure. It integrates with the Cluster API ecosystem to provide a consistent way to deploy workload clusters across different cloud providers.

The provider includes:
- **ContaboCluster**: Manages cluster-wide infrastructure (networking, load balancers)
- **ContaboMachine**: Manages individual virtual machines (VPS instances)
- **ContaboMachineTemplate**: Template for creating machines with consistent configuration

## Getting Started

### Prerequisites
- go version v1.24.0+
- docker version 17.03+
- kubectl version v1.11.3+
- Access to a Kubernetes v1.11.3+ cluster
- Contabo API token (obtainable from the Contabo customer portal)

### Quick Start

1. **Install clusterctl**
   ```sh
   curl -L https://github.com/kubernetes-sigs/cluster-api/releases/latest/download/clusterctl-linux-amd64 -o clusterctl
   chmod +x clusterctl
   sudo mv clusterctl /usr/local/bin/clusterctl
   ```

2. **Set up environment variables**
   ```sh
   export CONTABO_API_TOKEN="your-contabo-api-token"
   export CLUSTER_NAME="my-contabo-cluster"
   export KUBERNETES_VERSION="v1.28.0"
   export CONTROL_PLANE_MACHINE_COUNT=1
   export WORKER_MACHINE_COUNT=2
   ```

3. **Initialize the management cluster**
   ```sh
   clusterctl init --infrastructure contabo
   ```

4. **Generate cluster configuration**
   ```sh
   clusterctl generate cluster $CLUSTER_NAME \
     --infrastructure contabo \
     --kubernetes-version $KUBERNETES_VERSION \
     --control-plane-machine-count $CONTROL_PLANE_MACHINE_COUNT \
     --worker-machine-count $WORKER_MACHINE_COUNT \
     > $CLUSTER_NAME.yaml
   ```

5. **Create the cluster**
   ```sh
   kubectl apply -f $CLUSTER_NAME.yaml
   ```

6. **Get the kubeconfig for the new cluster**
   ```sh
   clusterctl get kubeconfig $CLUSTER_NAME > $CLUSTER_NAME.kubeconfig
   ```

### Development Setup

### To Deploy on the cluster

**Set up your Contabo API token:**
```sh
export CONTABO_API_TOKEN="your-contabo-api-token"
```

**Build and push your image to the location specified by `IMG`:**

```sh
make docker-build docker-push IMG=<some-registry>/cluster-api-provider-contabo:tag
```

**NOTE:** This image ought to be published in the personal registry you specified.
And it is required to have access to pull the image from the working environment.
Make sure you have the proper permission to the registry if the above commands don’t work.

**Install the CRDs into the cluster:**

```sh
make install
```

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=<some-registry>/cluster-api-provider-contabo:tag
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin
privileges or be logged in as admin.

**Create instances of your solution**
You can apply the samples (examples) from the config/sample:

```sh
kubectl apply -k config/samples/
```

**Important:** Before applying the samples, make sure to:
1. Set the `CONTABO_API_TOKEN` environment variable in the manager deployment
2. Update the sample configurations with your actual Contabo region and SSH keys
3. Ensure you have valid Contabo instance types and image IDs

>**NOTE**: Ensure that the samples have valid values for your Contabo account.

### To Uninstall
**Delete the instances (CRs) from the cluster:**

```sh
kubectl delete -k config/samples/
```

**Delete the APIs(CRDs) from the cluster:**

```sh
make uninstall
```

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```

## Configuration

### API Types

#### ContaboCluster
Manages cluster-wide infrastructure including networking and control plane endpoint.

**Key fields:**
- `region`: Contabo region (e.g., "EU", "US-central", "US-east", "US-west", "SIN")
- `controlPlaneEndpoint`: Kubernetes API server endpoint
- `network`: Network configuration including subnets

#### ContaboMachine
Manages individual VPS instances.

**Key fields:**
- `instanceType`: Contabo VPS plan (e.g., "S", "M", "L", "XL", "XXL")
- `image`: OS image ID (e.g., "ubuntu-22.04", "ubuntu-20.04", "centos-8")
- `region`: Contabo region where the instance will be created
- `sshKeys`: List of SSH key names/IDs
- `userData`: Cloud-init user data (base64 encoded)

#### ContaboMachineTemplate
Template for creating machines with consistent configuration.

### Environment Variables

- `CONTABO_API_TOKEN`: Your Contabo API token (required)

### Available Regions
- `EU`: European data centers
- `US-central`: US Central data centers
- `US-east`: US East data centers
- `US-west`: US West data centers
- `SIN`: Singapore data centers

### Available Instance Types
- `S`: Small VPS
- `M`: Medium VPS  
- `L`: Large VPS
- `XL`: Extra Large VPS
- `XXL`: Extra Extra Large VPS

See the Contabo documentation for current specifications and pricing.

## Project Distribution

### Creating a Release

Following the options to release and provide this solution to the users.

### By providing a bundle with all YAML files

1. Build the installer for the image built and published in the registry:

```sh
make build-installer IMG=<some-registry>/cluster-api-provider-contabo:tag
```

**NOTE:** The makefile target mentioned above generates an 'install.yaml'
file in the dist directory. This file contains all the resources built
with Kustomize, which are necessary to install this project without its
dependencies.

2. Using the installer

Users can just run 'kubectl apply -f <URL for YAML BUNDLE>' to install
the project, i.e.:

```sh
kubectl apply -f https://raw.githubusercontent.com/<org>/cluster-api-provider-contabo/<tag or branch>/dist/install.yaml
```

### By providing a Helm Chart

1. Build the chart using the optional helm plugin

```sh
kubebuilder edit --plugins=helm/v1-alpha
```

2. See that a chart was generated under 'dist/chart', and users
can obtain this solution from there.

**NOTE:** If you change the project, you need to update the Helm Chart
using the same command above to sync the latest changes. Furthermore,
if you create webhooks, you need to use the above command with
the '--force' flag and manually ensure that any custom configuration
previously added to 'dist/chart/values.yaml' or 'dist/chart/manager/manager.yaml'
is manually re-applied afterwards.

## Contributing

We welcome contributions to the Cluster API Contabo Provider! Here's how you can contribute:

### Development Workflow

1. **Fork the repository** and clone your fork
2. **Create a feature branch** from main
3. **Make your changes** and add tests
4. **Run tests** and ensure they pass:
   ```sh
   make test
   ```
5. **Build and test locally**:
   ```sh
   make build
   make docker-build
   ```
6. **Submit a pull request** with a clear description of your changes

### Code Guidelines

- Follow Go best practices and conventions
- Add unit tests for new functionality
- Update documentation for any API changes
- Ensure all CI checks pass

### Testing

Run the full test suite:
```sh
make test
```

Run end-to-end tests (requires a management cluster):
```sh
make test-e2e
```

## Troubleshooting

### Common Issues

**Authentication errors:**
- Verify your `CONTABO_API_TOKEN` is correct and has sufficient permissions
- Check that the token is properly set in the environment

**Instance creation failures:**
- Verify the instance type is available in your selected region
- Check that the image ID is valid and available
- Ensure your Contabo account has sufficient quota

**Network connectivity issues:**
- Verify security groups and firewall rules allow necessary traffic
- Check that the control plane endpoint is accessible

For more help, please open an issue in the GitHub repository.

**NOTE:** Run `make help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

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

