# Cluster API Provider Contabo - Instance Tagging

## Overview

Since Contabo instances are billed monthly or annually regardless of usage, this provider includes a tagging system to track cluster membership and instance availability. This helps with cost management and resource allocation.

## How It Works

### Cluster Membership Tags

When an instance is added to a cluster:
- A cluster-specific tag is created and applied (e.g., `cluster-api-cluster-my-cluster`)
- Any existing "available" tag is removed
- The instance is marked as actively used in a cluster

When an instance is removed from a cluster:
- The cluster-specific tag is removed
- An "available" tag is applied to mark the instance as ready for reuse
- The instance can be identified for potential cluster assignment

### Tag Naming Convention

- **Cluster Tag**: `{clusterTag}-{clusterName}` (default: `cluster-api-cluster-{clusterName}`)
- **Available Tag**: `{availableTag}` (default: `cluster-api-available`)

## Configuration

### Cluster-Level Configuration

Configure tagging at the cluster level in your `ContaboCluster` resource:

```yaml
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: ContaboCluster
metadata:
  name: my-cluster
spec:
  region: "EU"
  clusterTagging:
    enabled: true                    # Enable/disable tagging (default: true)
    clusterTag: "cluster-api-cluster" # Base name for cluster tags
    availableTag: "cluster-api-available" # Tag for available instances
    tagColor: "#1E90FF"              # Color for tags in Contabo UI
  # ... other cluster config
```

### Machine-Level Configuration

Override cluster settings at the machine level in your `ContaboMachine` resource:

```yaml
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: ContaboMachine
metadata:
  name: my-machine
spec:
  instanceType: "M"
  image: "ubuntu-22.04"
  region: "EU"
  clusterTagging:
    enabled: true
    clusterTag: "custom-cluster-tag"
    availableTag: "custom-available-tag"
    tagColor: "#32CD32"
  # ... other machine config
```

### Configuration Precedence

1. **Machine-level** `clusterTagging` (highest priority)
2. **Cluster-level** `clusterTagging`
3. **Default values** (if neither is specified)

### Default Values

```yaml
enabled: true
clusterTag: "cluster-api-cluster"
availableTag: "cluster-api-available"
tagColor: "#1E90FF"
```

## Use Cases

### Cost Management

1. **Track Active Usage**: Quickly identify which instances are actively used in clusters
2. **Find Available Instances**: Locate instances tagged as "available" for new cluster assignments
3. **Billing Allocation**: Associate running costs with specific clusters

### Resource Management

1. **Cluster Migration**: Move instances between clusters by updating tags
2. **Capacity Planning**: See cluster resource distribution across instances
3. **Cleanup Operations**: Identify instances that can be safely terminated

### Operations

1. **Monitoring**: Track cluster membership changes over time
2. **Automation**: Build scripts that respect cluster membership tags
3. **Compliance**: Ensure proper resource tagging for organizational policies

## Best Practices

### Tag Naming

- Use consistent naming conventions across your organization
- Include environment indicators in tag names if needed
- Keep tag names descriptive but concise

### Instance Lifecycle

1. **Pre-provisioned Instances**: Tag new instances as "available" when first created
2. **Cluster Assignment**: Let the provider automatically manage cluster tags
3. **Decommissioning**: Remove cluster tags before terminating instances

### Monitoring

Monitor tag changes to track:
- Instance assignment patterns
- Cluster growth and shrinkage
- Resource utilization trends

## Troubleshooting

### Tagging Disabled

If tagging isn't working, check:
1. `clusterTagging.enabled` is set to `true` (default)
2. Contabo API token has tag management permissions
3. Check controller logs for tagging errors

### Missing Tags

Common causes:
- API errors during tag creation (check logs)
- Timing issues during instance provisioning
- Manual tag removal in Contabo UI

### Tag Cleanup

If instances have incorrect tags:
1. Use Contabo UI or API to manually correct tags
2. Restart the affected machine controllers
3. Check cluster and machine configurations

## Examples

### Basic Cluster with Tagging

```yaml
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: ContaboCluster
metadata:
  name: production-cluster
spec:
  region: "EU"
  clusterTagging:
    enabled: true
    clusterTag: "prod-cluster"
    availableTag: "prod-available"
    tagColor: "#FF4500"
---
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: ContaboMachine
metadata:
  name: production-worker
spec:
  instanceType: "L"
  image: "ubuntu-22.04"
  region: "EU"
  # Inherits cluster tagging configuration
```

### Environment-Specific Tagging

```yaml
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: ContaboCluster
metadata:
  name: dev-cluster
spec:
  region: "EU"
  clusterTagging:
    enabled: true
    clusterTag: "dev-cluster"
    availableTag: "dev-available"
    tagColor: "#32CD32"
```

### Disabled Tagging

```yaml
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: ContaboCluster
metadata:
  name: test-cluster
spec:
  region: "EU"
  clusterTagging:
    enabled: false  # No automatic tagging
```